package rabbitmq_client

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"time"
)

type RabbitmqClient struct {
	name   string
	config Config
	conn   *amqp.Connection
	ch     *amqp.Channel

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type Config struct {
	URI        string `yaml:"uri"`
	RetryCount int    `yaml:"retry_count"`
}

func New(name string) *RabbitmqClient {
	return &RabbitmqClient{
		name: name,
	}
}

func (t *RabbitmqClient) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.URI == "" {
		return errors.New("uri is required")
	}

	t.conn, err = amqp.Dial(t.config.URI)

	if err != nil {
		return err
	}

	t.ch, err = t.conn.Channel()

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_rabbitmq." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_rabbitmq." + t.name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_rabbitmq." + t.name + ".retry")

	return nil
}

func (t *RabbitmqClient) Stop() error {
	t.ch.Close()
	t.conn.Close()
	return nil
}

func (t *RabbitmqClient) String() string {
	return t.name
}

func (t *RabbitmqClient) Fetch(ctx context.Context, topic string, data []byte, headers map[string][]string) chan IClientResponse {
	resp := make(chan IClientResponse)

	go func(resp chan IClientResponse, ctx context.Context, topic string, data []byte, headers map[string][]string) {
		defer func(s time.Time) {
			t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
		}(time.Now())
		t.counter.Add(context.Background(), 1)

		res := &Response{}
		header := make(map[string]interface{}, len(headers))

		for k, h := range headers {
			if len(h) > 0 {
				header[k] = h[0]
			}
		}

		h := make(map[string]string)
		otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(h))

		for k, v := range h {
			header[k] = v
		}

		for i := 0; i <= t.config.RetryCount; i++ {
			err := t.ch.PublishWithContext(
				ctx,
				"",
				topic,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
					Headers:     header,
				})

			if err != nil {
				res.Err = err
				res.RetryCount = i + 1

				if i == t.config.RetryCount {
					res.Err = fmt.Errorf("request rabbitmq to topic %s", topic)
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
