package server

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
	"sync"
	"time"
)

type Kafka struct {
	Name              string
	Hosts             string
	GroupId           string
	User              string
	Password          string
	Topics            []string
	LimitMessageCount int64

	app      interfaces.IService
	router   *Router
	consumer *kafka.Consumer
	done     chan interface{}
	w        sync.WaitGroup
}

func (t *Kafka) Init(app interfaces.IService) error {
	var err error
	t.app = app
	t.done = make(chan interface{})

	t.router = NewRouter()

	cfg := kafka.ConfigMap{
		"bootstrap.servers": t.Hosts,
		"group.id":          t.GroupId,
		"auto.offset.reset": "earliest",
	}

	if t.User != "" {
		cfg.SetKey("sasl.mechanisms", "PLAIN")
		cfg.SetKey("security.protocol", "SASL_SSL")
		cfg.SetKey("sasl.username", t.User)
		cfg.SetKey("sasl.password", t.Password)
	}

	t.consumer, err = kafka.NewConsumer(&cfg)

	if err != nil {
		return err
	}

	return t.consumer.SubscribeTopics(t.Topics, nil)
}

func (t *Kafka) Start() error {
	go func() {
		limited := make(chan interface{}, t.LimitMessageCount)
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

func (t *Kafka) PushRoute(method string, path string, handler interfaces.RouteFunc) {
	t.router.PushRoute(method, path, handler)
}

func (t *Kafka) Handle(msg *kafka.Message) {
	ctx := srvCtx.FromKafka(msg)
	defer srvCtx.ContextPool.Put(ctx)

	defer func(tm time.Time) {
		t.app.GetLogger().Debug(fmt.Sprintf("%d - %v, %s: %s", ctx.Response.StatusCode, time.Now().Sub(tm), "KAFKA", *msg.TopicPartition.Topic))
	}(time.Now())

	route := t.router.Find("KAFKA", *msg.TopicPartition.Topic)

	if route == nil {
		ctx.Response.StatusCode = 404
	} else {
		route(ctx)
	}
}
