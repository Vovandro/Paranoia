package server

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
	"time"
)

type Kafka struct {
	Name              string
	Hosts             string
	GroupId           string
	Topics            []string
	LimitMessageCount int64

	app      interfaces.IService
	router   *Router
	consumer *kafka.Consumer
	done     chan interface{}
}

func (t *Kafka) Init(app interfaces.IService) error {
	var err error
	t.app = app
	t.done = make(chan interface{})

	t.router = NewRouter()

	t.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": t.Hosts,
		"group.id":          t.GroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return err
	}

	return t.consumer.SubscribeTopics(t.Topics, nil)
}

func (t *Kafka) Start() error {
	go func() {
		for {
			select {
			case <-t.done:
				return

			default:
				msg, err := t.consumer.ReadMessage(time.Millisecond * 100)

				if err != nil {
					t.app.GetLogger().Error(err)
					continue
				}

				t.Handle(msg)
			}
		}
	}()

	return nil
}

func (t *Kafka) Stop() error {
	close(t.done)
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
		t.app.GetLogger().Debug(fmt.Sprintf("%d - %v, %s: %s", ctx.Response.StatusCode, time.Now().Sub(tm), "kafka", *msg.TopicPartition.Topic))
	}(time.Now())

	route := t.router.Find("KAFKA", *msg.TopicPartition.Topic)

	if route == nil {
		ctx.Response.StatusCode = 404
	} else {
		route(ctx)
	}
}
