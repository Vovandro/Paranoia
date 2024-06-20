package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
)

type MetricStd struct {
	cfg      MetricStdConfig
	exporter metric.Exporter
	app      interfaces.IService
}

type MetricStdConfig struct {
	Name     string
	Interval time.Duration
}

func NewMetricStd(cfg MetricStdConfig) *MetricStd {
	return &MetricStd{cfg: cfg}
}

func (t *MetricStd) Init(app interfaces.IService) error {
	t.app = app

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(t.cfg.Name),
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
		metric.WithReader(metric.NewPeriodicReader(t.exporter, metric.WithInterval(t.cfg.Interval))),
	)

	otel.SetMeterProvider(provider)

	return nil

}

func (t *MetricStd) Start() error {
	return nil
}

func (t *MetricStd) Stop() error {
	err := t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
