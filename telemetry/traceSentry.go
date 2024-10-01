package telemetry

import (
	"context"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceSentry struct {
	cfg      TraceSentryConfig
	provider *trace.TracerProvider
	app      interfaces.IEngine
}

type TraceSentryConfig struct {
	Name string `yaml:"name"`
}

func NewTraceSentry(cfg TraceSentryConfig) *TraceSentry {
	return &TraceSentry{cfg: cfg}
}

func (t *TraceSentry) Init(app interfaces.IEngine) error {
	t.app = app

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
	err := t.provider.Shutdown(context.Background())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
