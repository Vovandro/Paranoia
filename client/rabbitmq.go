package client

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"time"
)

type RabbitmqClient struct {
	Name   string
	Config RabbitmqClientConfig
	app    interfaces.IEngine
	conn   *amqp.Connection
	ch     *amqp.Channel

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type RabbitmqClientConfig struct {
	URI        string `yaml:"uri"`
	RetryCount int    `yaml:"retry_count"`
}

func NewRabbitmqClient(name string, cfg RabbitmqClientConfig) *RabbitmqClient {
	return &RabbitmqClient{
		Name:   name,
		Config: cfg,
	}
}

func (t *RabbitmqClient) Init(app interfaces.IEngine) error {
	var err error

	t.app = app

	t.conn, err = amqp.Dial(t.Config.URI)

	if err != nil {
		return err
	}

	t.ch, err = t.conn.Channel()

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_rabbitmq." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_rabbitmq." + t.Name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_rabbitmq." + t.Name + ".retry")

	return nil
}

func (t *RabbitmqClient) Stop() error {
	t.ch.Close()
	t.conn.Close()
	return nil
}

func (t *RabbitmqClient) String() string {
	return t.Name
}

func (t *RabbitmqClient) Fetch(ctx context.Context, _ string, topic string, data []byte, headers map[string][]string) chan interfaces.IClientResponse {
	resp := make(chan interfaces.IClientResponse)

	go func(resp chan interfaces.IClientResponse, ctx context.Context, topic string, data []byte, headers map[string][]string) {
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

		for i := 0; i <= t.Config.RetryCount; i++ {
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

				if i == t.Config.RetryCount {
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
