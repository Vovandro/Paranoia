package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type TraceZipking struct {
	cfg      TraceZipkingConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
	app      interfaces.IEngine
}

type TraceZipkingConfig struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func NewTraceZipking(cfg TraceZipkingConfig) *TraceZipking {
	return &TraceZipking{cfg: cfg}
}

func (t *TraceZipking) Init(app interfaces.IEngine) error {
	t.app = app

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	var err error
	t.exporter, err = zipkin.New(t.cfg.Url)

	if err != nil {
		return err
	}

	t.provider = trace.NewTracerProvider(
		trace.WithBatcher(t.exporter, trace.WithBatchTimeout(time.Second)),
	)

	otel.SetTracerProvider(t.provider)

	return nil
}

func (t *TraceZipking) Start() error {
	return nil
}

func (t *TraceZipking) Stop() error {
	err := t.provider.Shutdown(context.Background())

	if err != nil {
		t.app.GetLogger().Error(context.Background(), err)
	}

	err = t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(context.Background(), err)
	}

	return err
}
