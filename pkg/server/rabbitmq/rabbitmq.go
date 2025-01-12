package rabbitmq

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"sync"
	"time"
)

type Rabbitmq struct {
	name string

	config Config

	router *Router
	conn   *amqp.Connection
	ch     *amqp.Channel
	done   chan interface{}
	w      sync.WaitGroup
	md     func(RouteFunc) RouteFunc

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
}

type Config struct {
	URI               string   `yaml:"uri"`
	Queue             string   `yaml:"queue"`
	ConsumerName      string   `yaml:"consumer_name"`
	LimitMessageCount int64    `yaml:"limit_message_count"`
	BaseMiddleware    []string `yaml:"base_middleware"`
}

func NewRabbitmq(name string) *Rabbitmq {
	return &Rabbitmq{
		name: name,
	}
}

func (t *Rabbitmq) Init(cfg map[string]interface{}) error {
	middlewares := make(map[string]IMiddleware)

	if m, ok := cfg["middlewares"]; ok {
		middlewares = m.(map[string]IMiddleware)
		delete(cfg, "middlewares")
	}

	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.URI == "" {
		return errors.New("uri is required")
	}

	if t.config.ConsumerName == "" {
		return errors.New("consumer_name is required")
	}

	if t.config.Queue == "" {
		return errors.New("queue is required")
	}

	t.done = make(chan interface{})

	t.router = NewRouter(middlewares)

	if t.config.BaseMiddleware == nil {
		t.config.BaseMiddleware = []string{"timing"}
	}

	if len(t.config.BaseMiddleware) > 0 {
		t.md, err = t.router.HandlerMiddleware(t.config.BaseMiddleware)
		if err != nil {
			return err
		}
	}

	if t.md == nil {
		t.md = func(routeFunc RouteFunc) RouteFunc {
			return routeFunc
		}
	}

	t.conn, err = amqp.Dial(t.config.URI)

	if err != nil {
		return err
	}

	t.ch, err = t.conn.Channel()

	if err != nil {
		return err
	}

	_, err = t.ch.QueueDeclare(
		t.config.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("server_kafka." + t.name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_kafka." + t.name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_kafka." + t.name + ".time")

	return nil
}

func (t *Rabbitmq) Start() error {
	go func() {
		limited := make(chan interface{}, t.config.LimitMessageCount)
		defer close(limited)

		msgs, err := t.ch.Consume(
			t.config.Queue,        // очередь
			t.config.ConsumerName, // consumer
			false,                 // auto-ack
			false,                 // exclusive
			false,                 // no-local
			false,                 // no-wait
			nil,                   // args
		)

		if err != nil {
			return
		}

	forLoop:
		for {
			select {
			case <-t.done:
				break forLoop

			case msg := <-msgs:
				if msg.Body == nil {
					continue
				}

				t.w.Add(1)
				limited <- nil

				go func(delivery amqp.Delivery) {
					defer t.w.Done()
					t.Handle(delivery)
					<-limited
				}(msg)
			}
		}

		t.w.Wait()
	}()

	return nil
}

func (t *Rabbitmq) Stop() error {
	close(t.done)
	t.w.Wait()

	if err := t.ch.Close(); err != nil {
		return err
	}

	err := t.conn.Close()

	time.Sleep(time.Second)

	return nil
}

func (t *Rabbitmq) Name() string {
	return t.name
}

func (t *Rabbitmq) Type() string {
	return "server"
}

func (t *Rabbitmq) PushRoute(path string, handler RouteFunc, middlewares []string) {
	t.router.PushRoute(path, handler, middlewares)
}

func (t *Rabbitmq) Handle(msg amqp.Delivery) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	ctx := RabbitmqCtxPool.Get().(*RabbitmqCtx)
	defer RabbitmqCtxPool.Put(ctx)
	ctx.Fill(&msg)

	h := make(map[string]string, len(msg.Headers))
	for k, v := range msg.Headers {
		h[k] = v.(string)
	}
	consumerCtx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(h))
	consumerCtx, span := otel.Tracer("").Start(consumerCtx, t.config.Queue)
	defer span.End()

	route, _ := t.router.Find(t.config.Queue)

	if route == nil {
		ctx.GetResponse().SetStatus(404)
	} else {
		t.md(route)(consumerCtx, ctx)
	}

	if ctx.GetResponse().GetStatus() >= 400 {
		t.counterError.Add(context.Background(), 1)
	}

	// Подтверждение сообщения
	msg.Ack(false)
}
