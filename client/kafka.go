package client

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type KafkaClient struct {
	Name     string
	Config   KafkaClientConfig
	app      interfaces.IService
	producer *kafka.Producer
}

type KafkaClientConfig struct {
	Hosts      string
	User       string
	Password   string
	RetryCount int
}

func (t *KafkaClient) Init(app interfaces.IService) error {
	var err error

	t.app = app
	cfg := kafka.ConfigMap{
		"bootstrap.servers": t.Config.Hosts,
	}

	if t.Config.User != "" {
		_ = cfg.SetKey("sasl.mechanisms", "PLAIN")
		_ = cfg.SetKey("security.protocol", "SASL_SSL")
		_ = cfg.SetKey("sasl.username", t.Config.User)
		_ = cfg.SetKey("sasl.password", t.Config.Password)
	}

	t.producer, err = kafka.NewProducer(&cfg)

	if err != nil {
		return err
	}

	return nil
}

func (t *KafkaClient) Stop() error {
	t.producer.Close()
	return nil
}

func (t *KafkaClient) String() string {
	return t.Name
}

func (t *KafkaClient) Fetch(_ string, topic string, data []byte, headers map[string][]string) chan interfaces.IClientResponse {
	resp := make(chan interfaces.IClientResponse)

	go func(resp chan interfaces.IClientResponse, topic string, data []byte, headers map[string][]string) {
		res := &Response{}
		request := kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value:   data,
			Headers: make([]kafka.Header, len(headers)),
		}

		for k, h := range headers {
			for i := 0; i < len(h); i++ {
				request.Headers = append(request.Headers, kafka.Header{
					Key:   k,
					Value: []byte(h[i]),
				})
			}
		}

		for i := 0; i <= t.Config.RetryCount; i++ {
			err := t.producer.Produce(&request, nil)

			if err != nil {
				res.Err = err
				res.RetryCount = i + 1

				if i == t.Config.RetryCount {
					res.Err = fmt.Errorf("request kafka to topic %s", topic)
					break
				}

				continue
			}

			break
		}

		resp <- res
	}(resp, topic, data, headers)

	return resp
}
