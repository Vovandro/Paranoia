package server

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Kafka struct {
	Name              string
	Hosts             string
	GroupId           string
	Topic             string
	LimitMessageCount int64

	app      interfaces.IService
	router   *Router
	consumer *kafka.Consumer
}

func (t *Kafka) Init(app interfaces.IService) error {
	var err error
	t.app = app

	t.router = NewRouter()

	t.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": t.Hosts,
		"group.id":          t.GroupId,
		"auto.offset.reset": "earliest",
	})

	return err
}

func (t *Kafka) Start() error {
	return nil
}

func (t *Kafka) Stop() error {
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
