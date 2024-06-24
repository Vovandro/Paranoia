package telemetry

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net/http"
	"time"
)

type MetricPrometheus struct {
	app      interfaces.IService
	cfg      MetricsPrometheusConfig
	server   *http.Server
	exporter metric.Reader
	meter    api.Meter
}

type MetricsPrometheusConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}

func NewPrometheusMetrics(cfg MetricsPrometheusConfig) *MetricPrometheus {
	return &MetricPrometheus{cfg: cfg}
}

func (t *MetricPrometheus) Init(app interfaces.IService) error {
	t.app = app

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(t.cfg.Name),
		))

	if err != nil {
		return err
	}

	t.server = &http.Server{
		Addr:                         ":" + t.cfg.Port,
		Handler:                      promhttp.Handler(),
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  5 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  5 * time.Second,
	}

	t.exporter, err = prometheus.New()
	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(t.exporter))
	otel.SetMeterProvider(provider)

	return nil

}

func (t *MetricPrometheus) Start() error {
	listenErr := make(chan error, 1)

	go func() {
		listenErr <- t.server.ListenAndServe()
	}()

	select {
	case err := <-listenErr:
		t.app.GetLogger().Error(err)
		return err

	case <-time.After(time.Second):
		// pass
	}

	return nil
}

func (t *MetricPrometheus) Stop() error {
	err := t.server.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
		return err
	}

	err = t.exporter.Shutdown(context.TODO())

	if err != nil {
		t.app.GetLogger().Error(err)
	}

	return err
}
