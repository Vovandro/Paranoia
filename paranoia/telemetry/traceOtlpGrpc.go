package telemetry

import (
	"context"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceOtlpGrpc struct {
	name     string
	config   TraceOtlpGrpcConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
}

type TraceOtlpGrpcConfig struct {
	ServiceName string `yaml:"service_name"`
}

func NewTraceOtlpGrpc(name string) *TraceOtlpGrpc {
	return &TraceOtlpGrpc{
		name: name,
	}
}

func (t *TraceOtlpGrpc) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	t.exporter, err = otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure())

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
		return err
	}

	return t.exporter.Shutdown(context.TODO())
}

func (t *TraceOtlpGrpc) Name() string {
	return t.name
}
