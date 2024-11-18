package server

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/server/srvUtils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/net/trace"
	"sync"
	"time"
)

type Kafka struct {
	Name string

	Config KafkaConfig

	app      interfaces.IEngine
	router   *Router
	consumer *otelkafka.Consumer
	done     chan interface{}
	w        sync.WaitGroup
	md       func(interfaces.RouteFunc) interfaces.RouteFunc

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
}

type KafkaConfig struct {
	Hosts             string   `yaml:"hosts"`
	GroupId           string   `yaml:"group_id"`
	User              string   `yaml:"user"`
	Password          string   `yaml:"password"`
	Topics            []string `yaml:"topics"`
	LimitMessageCount int64    `yaml:"limit_message_count"`
	BaseMiddleware    []string `yaml:"base_middleware"`
}

func NewKafka(name string, cfg KafkaConfig) *Kafka {
	return &Kafka{
		Name:   name,
		Config: cfg,
	}
}

func (t *Kafka) Init(app interfaces.IEngine) error {
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

	cfg := kafka.ConfigMap{
		"bootstrap.servers": t.Config.Hosts,
		"group.id":          t.Config.GroupId,
		"auto.offset.reset": "earliest",
	}

	if t.Config.User != "" {
		cfg.SetKey("sasl.mechanisms", "PLAIN")
		cfg.SetKey("security.protocol", "SASL_SSL")
		cfg.SetKey("sasl.username", t.Config.User)
		cfg.SetKey("sasl.password", t.Config.Password)
	}

	t.consumer, err = otelkafka.NewConsumer(&cfg)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("server_kafka." + t.Name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_kafka." + t.Name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_kafka." + t.Name + ".time")

	return t.consumer.SubscribeTopics(t.Config.Topics, nil)
}

func (t *Kafka) Start() error {
	go func() {
		limited := make(chan interface{}, t.Config.LimitMessageCount)
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
					t.app.GetLogger().Error(context.Background(), err)
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

	if err != nil {
		t.app.GetLogger().Error(context.Background(), err)
	} else {
		t.app.GetLogger().Info(context.Background(), "kafka consumer gracefully stopped.")
		time.Sleep(time.Second)
	}

	return err
}

func (t *Kafka) String() string {
	return t.Name
}

func (t *Kafka) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}

func (t *Kafka) Handle(msg *kafka.Message) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	tr := trace.New(t.Name, msg.TopicPartition.String())
	defer tr.Finish()

	ctx := srvUtils.KafkaCtxPool.Get().(*srvUtils.KafkaCtx)
	defer srvUtils.KafkaCtxPool.Put(ctx)
	ctx.Fill(msg)

	route, _ := t.router.Find("KAFKA", *msg.TopicPartition.Topic)

	if route == nil {
		ctx.GetResponse().SetStatus(404)
	} else {
		t.md(route)(context.Background(), ctx)
	}

	if ctx.GetResponse().GetStatus() >= 400 {
		t.counterError.Add(context.Background(), 1)
	}
}
