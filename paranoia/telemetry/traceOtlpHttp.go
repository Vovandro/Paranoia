package telemetry

import (
	"context"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceOtlpHttp struct {
	name     string
	config   TraceOtlpHttpConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
}

type TraceOtlpHttpConfig struct {
	ServiceName string `yaml:"service_name"`
}

func NewTraceOtlpHttp(name string) *TraceOtlpHttp {
	return &TraceOtlpHttp{
		name: name,
	}
}

func (t *TraceOtlpHttp) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	t.exporter, err = otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())

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
		return err
	}

	return t.exporter.Shutdown(context.TODO())
}

func (t *TraceOtlpHttp) Name() string {
	return t.name
}
