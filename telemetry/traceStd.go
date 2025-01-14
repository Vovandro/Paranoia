package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type TraceStd struct {
	name     string
	config   TraceStdConfig
	exporter trace.SpanExporter
	provider *trace.TracerProvider
}

type TraceStdConfig struct {
	ServiceName string        `yaml:"service_name"`
	Interval    time.Duration `yaml:"interval"`
}

func NewTraceStd(name string) *TraceStd {
	return &TraceStd{
		name: name,
	}
}

func (t *TraceStd) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	t.exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())

	if err != nil {
		return err
	}

	t.provider = trace.NewTracerProvider(
		trace.WithBatcher(t.exporter, trace.WithBatchTimeout(t.config.Interval)),
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
		return err
	}

	return t.exporter.Shutdown(context.TODO())
}

func (t *TraceStd) Name() string {
	return t.name
}
