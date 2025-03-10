package telemetry

import (
	"context"
	"time"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type MetricOtlpHttp struct {
	config   MetricOtlpHttpConfig
	exporter metric.Exporter
	name     string
}

type MetricOtlpHttpConfig struct {
	ServiceName string        `yaml:"service_name"`
	Interval    time.Duration `yaml:"interval"`
}

func NewMetricOtlpHttp(name string) *MetricOtlpHttp {
	return &MetricOtlpHttp{
		name: name,
	}
}

func (t *MetricOtlpHttp) Init(cfg map[string]interface{}) error {
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

	t.exporter, err = otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithInsecure())

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

func (t *MetricOtlpHttp) Start() error {
	return nil
}

func (t *MetricOtlpHttp) Stop() error {
	return t.exporter.Shutdown(context.TODO())
}

func (t *MetricOtlpHttp) Name() string {
	return t.name
}
