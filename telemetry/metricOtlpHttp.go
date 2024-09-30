package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
)

type MetricOtlpHttp struct {
	cfg      MetricOtlpHttpConfig
	exporter metric.Exporter
	app      interfaces.IEngine
}

type MetricOtlpHttpConfig struct {
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
}

func NewMetricOtlpHttp(cfg MetricOtlpHttpConfig) *MetricOtlpHttp {
	return &MetricOtlpHttp{cfg: cfg}
}

func (t *MetricOtlpHttp) Init(app interfaces.IEngine) error {
	t.app = app

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(t.cfg.Name),
		))

	if err != nil {
		return err
	}

	t.exporter, err = otlpmetrichttp.New(context.Background())

	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(t.exporter, metric.WithInterval(t.cfg.Interval))),
	)

	otel.SetMeterProvider(provider)

	return nil

}

func (t *MetricOtlpHttp) Start() error {
	return nil
}

func (t *MetricOtlpHttp) Stop() error {
	err := t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
