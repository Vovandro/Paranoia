package telemetry

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
)

type Metrics struct {
	app interfaces.IService
	cfg MetricsConfig
}

type MetricsConfig struct {
	Type string
	Name string
	Port string
}

func NewMetrics(cfg MetricsConfig) *Metrics {
	return &Metrics{cfg: cfg}
}

func (t *Metrics) Init(app interfaces.IService) error {
	t.app = app

	var exporter metric.Reader
	var err error

	switch t.cfg.Type {
	case "prometheus":
		exporter, err = prometheus.New()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsuported metrics type: %s", t.cfg.Type)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(t.cfg.Name)

	go t.serveMetrics()

	//opt := api.WithDescription("test description")
	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	counter.Add(ctx, 5)

	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		n := -10. + rand.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.Float64Histogram(
		"baz",
		api.WithDescription("a histogram with custom buckets and rename"),
		api.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	)
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 136)
	histogram.Record(ctx, 64)
	histogram.Record(ctx, 701)
	histogram.Record(ctx, 830)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()

	return nil
}

func (t *Metrics) serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/", promhttp.Handler())
	err := http.ListenAndServe(":"+t.cfg.Port, nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
