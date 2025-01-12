package server

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/server/srvUtils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"sync"
	"time"
)

type Rabbitmq struct {
	Name string

	Config RabbitmqConfig

	app    interfaces.IEngine
	router *Router
	conn   *amqp.Connection
	ch     *amqp.Channel
	done   chan interface{}
	w      sync.WaitGroup
	md     func(interfaces.RouteFunc) interfaces.RouteFunc

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
}

type RabbitmqConfig struct {
	URI               string   `yaml:"uri"`
	Queue             string   `yaml:"queue"`
	ConsumerName      string   `yaml:"consumer_name"`
	LimitMessageCount int64    `yaml:"limit_message_count"`
	BaseMiddleware    []string `yaml:"base_middleware"`
}

func NewRabbitmq(name string, cfg RabbitmqConfig) *Rabbitmq {
	return &Rabbitmq{
		Name:   name,
		Config: cfg,
	}
}

func (t *Rabbitmq) Init(app interfaces.IEngine) error {
	var err error
	t.app = app
	t.done = make(chan interface{})

	t.router = NewRouter(app)

	if t.Config.BaseMiddleware == nil {
		t.Config.BaseMiddleware = []string{"timing"}
	}

	if len(t.Config.BaseMiddleware) > 0 {
		t.md = middleware.HandlerFromStrings(app, t.Config.BaseMiddleware)
	}

	if t.md == nil {
		t.md = func(routeFunc interfaces.RouteFunc) interfaces.RouteFunc {
			return routeFunc
		}
	}

	t.conn, err = amqp.Dial(t.Config.URI)

	if err != nil {
		return err
	}

	t.ch, err = t.conn.Channel()

	if err != nil {
		return err
	}

	_, err = t.ch.QueueDeclare(
		t.Config.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("server_kafka." + t.Name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_kafka." + t.Name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_kafka." + t.Name + ".time")

	return nil
}

func (t *Rabbitmq) Start() error {
	go func() {
		limited := make(chan interface{}, t.Config.LimitMessageCount)
		defer close(limited)

		msgs, err := t.ch.Consume(
			t.Config.Queue,        // очередь
			t.Config.ConsumerName, // consumer
			false,                 // auto-ack
			false,                 // exclusive
			false,                 // no-local
			false,                 // no-wait
			nil,                   // args
		)

		if err != nil {
			t.app.GetLogger().Error(context.Background(), err)
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
		t.app.GetLogger().Error(context.Background(), err)
	}

	if err := t.conn.Close(); err != nil {
		t.app.GetLogger().Error(context.Background(), err)
	} else {
		t.app.GetLogger().Info(context.Background(), "RabbitMQ consumer gracefully stopped.")
		time.Sleep(time.Second)
	}

	return nil
}

func (t *Rabbitmq) String() string {
	return t.Name
}

func (t *Rabbitmq) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}

func (t *Rabbitmq) Handle(msg amqp.Delivery) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	ctx := srvUtils.RabbitmqCtxPool.Get().(*srvUtils.RabbitmqCtx)
	defer srvUtils.RabbitmqCtxPool.Put(ctx)
	ctx.Fill(&msg)

	h := make(map[string]string, len(msg.Headers))
	for k, v := range msg.Headers {
		h[k] = v.(string)
	}
	consumerCtx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(h))
	consumerCtx, span := otel.Tracer("").Start(consumerCtx, t.Config.Queue)
	defer span.End()

	route, _ := t.router.Find("RABBITMQ", t.Config.Queue)

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
