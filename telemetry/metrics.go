package telemetry

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"net/http"
	"time"
)

type Metrics struct {
	Name     string
	app      interfaces.IService
	cfg      MetricsConfig
	server   *http.Server
	exporter metric.Reader
	meter    api.Meter
}

type MetricsConfig struct {
	Type string
	Name string
	Port string
}

func NewMetrics(name string, cfg MetricsConfig) *Metrics {
	return &Metrics{Name: name, cfg: cfg}
}

func (t *Metrics) Init(app interfaces.IService) error {
	t.app = app

	t.server = &http.Server{
		Addr:                         ":" + t.cfg.Port,
		Handler:                      promhttp.Handler(),
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  5 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  5 * time.Second,
	}

	var err error

	switch t.cfg.Type {
	case "prometheus":
		t.exporter, err = prometheus.New()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsuported metrics type: %s", t.cfg.Type)
	}

	provider := metric.NewMeterProvider(metric.WithReader(t.exporter))
	t.meter = provider.Meter(t.cfg.Name)

	return nil

}

func (t *Metrics) Start() error {
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

func (t *Metrics) Stop() error {
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

func (t *Metrics) String() string {
	return t.Name
}
