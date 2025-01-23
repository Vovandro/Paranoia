package telemetry

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"net/http"
	"time"
)

type MetricPrometheus struct {
	name     string
	config   MetricPrometheusConfig
	server   *http.Server
	exporter metric.Reader
	meter    api.Meter
}

type MetricPrometheusConfig struct {
	ServiceName string `yaml:"service_name"`
	Port        string `yaml:"port"`
}

func NewMetricPrometheus(name string) *MetricPrometheus {
	return &MetricPrometheus{
		name: name,
	}
}

func (t *MetricPrometheus) Init(cfg map[string]interface{}) error {
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

	t.server = &http.Server{
		Addr:                         ":" + t.config.Port,
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
		return err

	case <-time.After(time.Second):
		// pass
	}

	return nil
}

func (t *MetricPrometheus) Stop() error {
	err := t.server.Shutdown(context.TODO())

	if err != nil {
		return err
	}

	return t.exporter.Shutdown(context.TODO())
}

func (t *MetricPrometheus) Name() string {
	return t.name
}
