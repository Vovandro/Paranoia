package kafka_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type KafkaClient struct {
	name     string
	config   Config
	producer *otelkafka.Producer

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type Config struct {
	Hosts      string `yaml:"hosts"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	SecurityProtocol string `yaml:"security_protocol"`
	SaslMechanisms   string `yaml:"sasl_mechanisms"`
	RetryCount int    `yaml:"retry_count"`
}

func New(name string) *KafkaClient {
	return &KafkaClient{
		name: name,
	}
}

func (t *KafkaClient) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	cfgKafka := kafka.ConfigMap{
		"bootstrap.servers": t.config.Hosts,
	}

	if t.config.Username != "" {
		if t.config.SecurityProtocol != "" {
			_ = cfgKafka.SetKey("security.protocol", t.config.SecurityProtocol)
		}

		if t.config.SaslMechanisms != "" {
			_ = cfgKafka.SetKey("sasl.mechanisms", t.config.SaslMechanisms)
		}

		_ = cfgKafka.SetKey("sasl.username", t.config.Username)
		_ = cfgKafka.SetKey("sasl.password", t.config.Password)
	}

	t.producer, err = otelkafka.NewProducer(&cfgKafka)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_kafka." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_kafka." + t.name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_kafka." + t.name + ".retry")

	return nil
}

func (t *KafkaClient) Stop() error {
	t.producer.Close()
	return nil
}

func (t *KafkaClient) Name() string {
	return t.name
}

func (t *KafkaClient) Type() string {
	return "client"
}

func (t *KafkaClient) Fetch(ctx context.Context, topic string, data []byte, headers map[string][]string) chan IClientResponse {
	resp := make(chan IClientResponse)

	go func(resp chan IClientResponse, ctx context.Context, topic string, data []byte, headers map[string][]string) {
		defer func(s time.Time) {
			t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
		}(time.Now())
		t.counter.Add(context.Background(), 1)

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

		otel.GetTextMapPropagator().Inject(ctx, otelkafka.NewMessageCarrier(&request))

		for i := 0; i <= t.config.RetryCount; i++ {
			err := t.producer.Produce(&request, nil)

			if err != nil {
				res.Err = err
				res.RetryCount = i + 1

				if i == t.config.RetryCount {
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
