package server

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
	"sync"
	"time"
)

type Kafka struct {
	Name string

	Config KafkaConfig

	app      interfaces.IService
	router   *Router
	consumer *kafka.Consumer
	done     chan interface{}
	w        sync.WaitGroup
	md       func(interfaces.RouteFunc) interfaces.RouteFunc
}

type KafkaConfig struct {
	Hosts             string
	GroupId           string
	User              string
	Password          string
	Topics            []string
	LimitMessageCount int64
	BaseMiddleware    []string
}

func (t *Kafka) Init(app interfaces.IService) error {
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

	t.consumer, err = kafka.NewConsumer(&cfg)

	if err != nil {
		return err
	}

	return t.consumer.SubscribeTopics(t.Config.Topics, nil)
}

func (t *Kafka) Start() error {
	go func() {
		limited := make(chan interface{}, t.Config.LimitMessageCount)
		defer close(limited)

		for {
			select {
			case <-t.done:
				break

			default:
				msg, err := t.consumer.ReadMessage(time.Millisecond * 100)
				t.w.Add(1)
				limited <- nil

				if err != nil {
					t.app.GetLogger().Error(err)
					continue
				}

				go func() {
					t.Handle(msg)
					t.w.Done()
					<-limited
				}()
			}
		}
	}()

	return nil
}

func (t *Kafka) Stop() error {
	close(t.done)
	t.w.Wait()
	err := t.consumer.Close()

	if err != nil {
		t.app.GetLogger().Error(err)
	} else {
		t.app.GetLogger().Info("kafka consumer gracefully stopped.")
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
	ctx := srvCtx.FromKafka(msg)
	defer srvCtx.ContextPool.Put(ctx)

	route := t.router.Find("KAFKA", *msg.TopicPartition.Topic)

	if route == nil {
		ctx.Response.StatusCode = 404
	} else {
		t.md(route)(ctx)
	}
}
