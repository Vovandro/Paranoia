package telemetry

import (
	"context"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceSentry struct {
	name     string
	config   TraceSentryConfig
	provider *trace.TracerProvider
}

type TraceSentryConfig struct {
	ServiceName string `yaml:"service_name"`
}

func NewTraceSentry(name string) *TraceSentry {
	return &TraceSentry{
		name: name,
	}
}

func (t *TraceSentry) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	t.provider = trace.NewTracerProvider(
		trace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	otel.SetTracerProvider(t.provider)

	return nil
}

func (t *TraceSentry) Start() error {
	return nil
}

func (t *TraceSentry) Stop() error {
	return t.provider.Shutdown(context.Background())
}

func (t *TraceSentry) Name() string {
	return t.name
}
