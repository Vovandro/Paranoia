package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type TraceZipking struct {
	name     string
	config   TraceZipkingConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
}

type TraceZipkingConfig struct {
	ServiceName string `yaml:"service_name"`
	Url         string `yaml:"url"`
}

func NewTraceZipking(name string) *TraceZipking {
	return &TraceZipking{
		name: name,
	}
}

func (t *TraceZipking) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	t.exporter, err = zipkin.New(t.config.Url)

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
		return err
	}

	return t.exporter.Shutdown(context.TODO())
}

func (t *TraceZipking) Name() string {
	return t.name
}
