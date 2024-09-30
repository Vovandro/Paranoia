package telemetry

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
)

type MetricOtlpGrpc struct {
	cfg      MetricOtlpGrpcConfig
	exporter metric.Exporter
	app      interfaces.IEngine
}

type MetricOtlpGrpcConfig struct {
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
}

func NewMetricOtlpGrpc(cfg MetricOtlpGrpcConfig) *MetricOtlpGrpc {
	return &MetricOtlpGrpc{cfg: cfg}
}

func (t *MetricOtlpGrpc) Init(app interfaces.IEngine) error {
	t.app = app

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(t.cfg.Name),
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
		metric.WithReader(metric.NewPeriodicReader(t.exporter, metric.WithInterval(t.cfg.Interval))),
	)

	otel.SetMeterProvider(provider)

	return nil

}

func (t *MetricOtlpGrpc) Start() error {
	return nil
}

func (t *MetricOtlpGrpc) Stop() error {
	err := t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
