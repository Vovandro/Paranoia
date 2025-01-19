package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

type MetricStd struct {
	name     string
	config   MetricStdConfig
	exporter metric.Exporter
}

type MetricStdConfig struct {
	ServiceName string        `yaml:"service_name"`
	Interval    time.Duration `yaml:"interval"`
}

func NewMetricStd(name string) *MetricStd {
	return &MetricStd{
		name: name,
	}
}

func (t *MetricStd) Init(cfg map[string]interface{}) error {
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

	t.exporter, err = stdoutmetric.New()

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

func (t *MetricStd) Start() error {
	return nil
}

func (t *MetricStd) Stop() error {
	return t.exporter.Shutdown(context.TODO())
}

func (t *MetricStd) Name() string {
	return t.name
}
