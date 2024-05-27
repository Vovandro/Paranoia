package client

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type KafkaClient struct {
	Name       string
	Hosts      string
	User       string
	Password   string
	RetryCount int
	app        interfaces.IService
	producer   *kafka.Producer
}

func (t *KafkaClient) Init(app interfaces.IService) error {
	var err error

	t.app = app
	cfg := kafka.ConfigMap{
		"bootstrap.servers": t.Hosts,
	}

	if t.User != "" {
		cfg.SetKey("sasl.mechanisms", "PLAIN")
		cfg.SetKey("security.protocol", "SASL_SSL")
		cfg.SetKey("sasl.username", t.User)
		cfg.SetKey("sasl.password", t.Password)
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

		for i := 0; i <= t.RetryCount; i++ {
			err := t.producer.Produce(&request, nil)

			if err != nil {
				res.Err = err
				res.RetryCount = i + 1

				if i == t.RetryCount {
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
