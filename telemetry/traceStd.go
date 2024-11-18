package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type TraceStd struct {
	cfg      TraceStdConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
	app      interfaces.IEngine
}

type TraceStdConfig struct {
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
}

func NewTraceStd(cfg TraceStdConfig) *TraceStd {
	return &TraceStd{cfg: cfg}
}

func (t *TraceStd) Init(app interfaces.IEngine) error {
	t.app = app

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	var err error
	t.exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())

	if err != nil {
		return err
	}

	t.provider = trace.NewTracerProvider(
		trace.WithBatcher(t.exporter, trace.WithBatchTimeout(t.cfg.Interval)),
	)

	otel.SetTracerProvider(t.provider)

	return nil
}

func (t *TraceStd) Start() error {
	return nil
}

func (t *TraceStd) Stop() error {
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
