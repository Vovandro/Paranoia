package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

type MetricOtlpGrpc struct {
	config   MetricOtlpGrpcConfig
	exporter metric.Exporter
	name     string
}

type MetricOtlpGrpcConfig struct {
	ServiceName string        `yaml:"service_name"`
	Interval    time.Duration `yaml:"interval"`
}

func NewMetricOtlpGrpc(name string) *MetricOtlpGrpc {
	return &MetricOtlpGrpc{
		name: name,
	}
}

func (t *MetricOtlpGrpc) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(t.config.ServiceName),
		))

	if err != nil {
		return err
	}

	t.exporter, err = otlpmetricgrpc.New(context.Background())

	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(t.exporter, metric.WithInterval(t.config.Interval))),
	)

	otel.SetMeterProvider(provider)

	return nil

}

func (t *MetricOtlpGrpc) Start() error {
	return nil
}

func (t *MetricOtlpGrpc) Stop() error {
	return t.exporter.Shutdown(context.TODO())
}

func (t *MetricOtlpGrpc) Name() string {
	return t.name
}
