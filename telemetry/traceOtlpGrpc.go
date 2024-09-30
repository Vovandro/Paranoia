package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceOtlpGrpc struct {
	cfg      TraceOtlpGrpcConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
	app      interfaces.IEngine
}

type TraceOtlpGrpcConfig struct {
	Name string `yaml:"name"`
}

func NewTraceOtlpGrpc(cfg TraceOtlpGrpcConfig) *TraceOtlpGrpc {
	return &TraceOtlpGrpc{cfg: cfg}
}

func (t *TraceOtlpGrpc) Init(app interfaces.IEngine) error {
	t.app = app

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	var err error
	t.exporter, err = otlptracegrpc.New(context.Background())

	if err != nil {
		return err
	}

	t.provider = trace.NewTracerProvider(
		trace.WithBatcher(t.exporter),
	)

	otel.SetTracerProvider(t.provider)

	return nil
}

func (t *TraceOtlpGrpc) Start() error {
	return nil
}

func (t *TraceOtlpGrpc) Stop() error {
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
