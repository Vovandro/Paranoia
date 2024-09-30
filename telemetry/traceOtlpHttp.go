package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceOtlpHttp struct {
	cfg      TraceOtlpHttpConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
	app      interfaces.IEngine
}

type TraceOtlpHttpConfig struct {
	Name string `yaml:"name"`
}

func NewTraceOtlpHttp(cfg TraceOtlpHttpConfig) *TraceOtlpHttp {
	return &TraceOtlpHttp{cfg: cfg}
}

func (t *TraceOtlpHttp) Init(app interfaces.IEngine) error {
	t.app = app

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	var err error
	t.exporter, err = otlptracehttp.New(context.Background())

	if err != nil {
		return err
	}

	t.provider = trace.NewTracerProvider(
		trace.WithBatcher(t.exporter),
	)

	otel.SetTracerProvider(t.provider)

	return nil
}

func (t *TraceOtlpHttp) Start() error {
	return nil
}

func (t *TraceOtlpHttp) Stop() error {
	err := t.provider.Shutdown(context.Background())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	err = t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
