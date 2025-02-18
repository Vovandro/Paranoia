package kafka

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"sync"
	"time"
)

type Kafka struct {
	name string

	config Config

	router   *Router
	consumer *otelkafka.Consumer
	done     chan interface{}
	w        sync.WaitGroup
	md       func(RouteFunc) RouteFunc

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
}

type Config struct {
	Hosts             string   `yaml:"hosts"`
	GroupId           string   `yaml:"group_id"`
	User              string   `yaml:"user"`
	Password          string   `yaml:"password"`
	SecurityProtocol  string   `yaml:"security_protocol"`
	SaslMechanisms    string   `yaml:"sasl_mechanisms"`
	Topics            []string `yaml:"topics"`
	LimitMessageCount int64    `yaml:"limit_message_count"`
	BaseMiddleware    []string `yaml:"base_middleware"`
}

func New(name string) *Kafka {
	return &Kafka{
		name: name,
	}
}

func (t *Kafka) Init(cfg map[string]interface{}) error {
	middlewares := make(map[string]IMiddleware)

	if m, ok := cfg["middlewares"]; ok {
		for k, v := range m.(map[string]interface{}) {
			if md, ok := v.(IMiddleware); ok {
				middlewares[k] = md
			}
		}
		delete(cfg, "middlewares")
	}

	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	if t.config.GroupId == "" {
		return errors.New("group_id is required")
	}

	if t.config.Topics == nil || len(t.config.Topics) == 0 {
		return errors.New("topics is required")
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

	cfgKafka := kafka.ConfigMap{
		"bootstrap.servers": t.config.Hosts,
		"group.id":          t.config.GroupId,
		"auto.offset.reset": "earliest",
	}

	if t.config.User != "" {
		if t.config.SecurityProtocol != "" {
			cfgKafka.SetKey("security.protocol", t.config.SecurityProtocol)
		}

		if t.config.SaslMechanisms != "" {
			cfgKafka.SetKey("sasl.mechanisms", t.config.SaslMechanisms)
		}

		cfgKafka.SetKey("sasl.username", t.config.User)
		cfgKafka.SetKey("sasl.password", t.config.Password)
	}

	t.consumer, err = otelkafka.NewConsumer(&cfgKafka)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("server_kafka." + t.name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_kafka." + t.name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_kafka." + t.name + ".time")

	return t.consumer.SubscribeTopics(t.config.Topics, nil)
}

func (t *Kafka) Start() error {
	go func() {
		limited := make(chan interface{}, t.config.LimitMessageCount)
		defer close(limited)

	forLoop:
		for {
			select {
			case <-t.done:
				break forLoop

			default:
				t.w.Add(1)
				limited <- nil
				msg, err := t.consumer.ReadMessage(time.Second)

				if err != nil {
					<-limited
					t.w.Done()
					continue
				}

				go func() {
					t.Handle(msg)
					<-limited
					t.w.Done()
				}()
			}
		}

		t.w.Wait()
	}()

	return nil
}

func (t *Kafka) Stop() error {
	_ = t.consumer.Unsubscribe()
	close(t.done)
	t.w.Wait()
	err := t.consumer.Close()

	time.Sleep(time.Second)

	return err
}

func (t *Kafka) Name() string {
	return t.name
}

func (t *Kafka) Type() string {
	return "server"
}

func (t *Kafka) PushRoute(path string, handler RouteFunc, middlewares []string) {
	t.router.PushRoute(path, handler, middlewares)
}

func (t *Kafka) Handle(msg *kafka.Message) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	c, tr := otel.Tracer("").Start(context.Background(), msg.TopicPartition.String())
	defer tr.End()

	ctx := KafkaCtxPool.Get().(*KafkaCtx)
	defer KafkaCtxPool.Put(ctx)
	ctx.Fill(msg)

	route, _ := t.router.Find(*msg.TopicPartition.Topic)

	if route == nil {
		ctx.GetResponse().SetStatus(404)
	} else {
		t.md(route)(c, ctx)
	}

	if ctx.GetResponse().GetStatus() >= 400 {
		t.counterError.Add(context.Background(), 1)
	}
}
