package client

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	Name   string
	Config HTTPClientConfig
	app    interfaces.IEngine
	client http.Client

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type HTTPClientConfig struct {
	RetryCount int `yaml:"retry_count"`
}

func NewHTTPClient(name string, cfg HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		Name:   name,
		Config: cfg,
	}
}

func (t *HTTPClient) Init(app interfaces.IEngine) error {
	t.app = app
	t.client = http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_http." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_http." + t.Name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_http." + t.Name + ".retry")

	return nil
}

func (t *HTTPClient) Stop() error {
	t.client.CloseIdleConnections()
	return nil
}

func (t *HTTPClient) String() string {
	return t.Name
}

func (t *HTTPClient) Fetch(method string, host string, data []byte, headers map[string][]string) chan interfaces.IClientResponse {
	resp := make(chan interfaces.IClientResponse)

	go func(resp chan interfaces.IClientResponse, method string, host string, data []byte, headers map[string][]string) {
		defer func(s time.Time) {
			t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
		}(time.Now())
		t.counter.Add(context.Background(), 1)

		res := &Response{}
		request, _ := http.NewRequest(method, host, bytes.NewBuffer(data))

		for i := 0; i <= t.Config.RetryCount; i++ {
			do, err := t.client.Do(request)
			if err != nil {
				res.Err = err
				res.RetryCount = i + 1
				break
			}

			if do.StatusCode == 200 {
				res.RetryCount = i + 1
				res.Body, _ = io.ReadAll(do.Body)
				res.Header = map[string][]string{}
				for s, strings := range do.Header {
					res.Header[s] = strings
				}
				do.Body.Close()
				break
			}

			if do.StatusCode > 200 && do.StatusCode < 300 {
				res.RetryCount = i + 1
				res.Header = map[string][]string{}
				for s, strings := range do.Header {
					res.Header[s] = strings
				}
				do.Body.Close()
				break
			}

			if (do.StatusCode >= 500 && do.StatusCode < 600) || do.StatusCode == 499 {
				if i+1 == t.Config.RetryCount {
					res.RetryCount = i + 1
					res.Err = fmt.Errorf("max retry count exceeded")
				}
				do.Body.Close()
				continue
			}

			res.RetryCount = i + 1
			res.Err = fmt.Errorf("request status code %d", do.StatusCode)
			do.Body.Close()
			break
		}

		t.retryCounter.Record(context.Background(), int64(res.RetryCount))

		resp <- res
	}(resp, method, host, data, headers)

	return resp
}
