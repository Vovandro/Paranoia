package kafka_client

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"gitlab.com/devpro_studio/Paranoia/client"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type KafkaClient struct {
	Name     string
	Config   KafkaClientConfig
	app      interfaces.IEngine
	producer *otelkafka.Producer

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type KafkaClientConfig struct {
	Hosts      string `yaml:"hosts"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	RetryCount int    `yaml:"retry_count"`
}

func NewKafkaClient(name string, cfg KafkaClientConfig) *KafkaClient {
	return &KafkaClient{
		Name:   name,
		Config: cfg,
	}
}

func (t *KafkaClient) Init(app interfaces.IEngine) error {
	var err error

	t.app = app
	cfg := kafka.ConfigMap{
		"bootstrap.servers": t.Config.Hosts,
	}

	if t.Config.Username != "" {
		_ = cfg.SetKey("sasl.mechanisms", "PLAIN")
		_ = cfg.SetKey("security.protocol", "SASL_SSL")
		_ = cfg.SetKey("sasl.username", t.Config.Username)
		_ = cfg.SetKey("sasl.password", t.Config.Password)
	}

	t.producer, err = otelkafka.NewProducer(&cfg)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_kafka." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_kafka." + t.Name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_kafka." + t.Name + ".retry")

	return nil
}

func (t *KafkaClient) Stop() error {
	t.producer.Close()
	return nil
}

func (t *KafkaClient) String() string {
	return t.Name
}

func (t *KafkaClient) Fetch(ctx context.Context, _ string, topic string, data []byte, headers map[string][]string) chan interfaces.IClientResponse {
	resp := make(chan interfaces.IClientResponse)

	go func(resp chan interfaces.IClientResponse, ctx context.Context, topic string, data []byte, headers map[string][]string) {
		defer func(s time.Time) {
			t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
		}(time.Now())
		t.counter.Add(context.Background(), 1)

		res := &client.Response{}
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

		otel.GetTextMapPropagator().Inject(ctx, otelkafka.NewMessageCarrier(&request))

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

			res.Code = 200

			break
		}

		t.retryCounter.Record(context.Background(), int64(res.RetryCount))

		resp <- res
	}(resp, ctx, topic, data, headers)

	return resp
}
