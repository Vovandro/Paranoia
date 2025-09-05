package http_client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type HTTPClient struct {
	name   string
	config Config
	client http.Client

	counter      metric.Int64Counter
	timeCounter  metric.Int64Histogram
	retryCounter metric.Int64Histogram
}

type Config struct {
	RetryCount int `yaml:"retry_count"`
}

func New(name string) *HTTPClient {
	return &HTTPClient{
		name: name,
	}
}

func (t *HTTPClient) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	t.client = http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	t.counter, _ = otel.Meter("").Int64Counter("client_http." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("client_http." + t.name + ".time")
	t.retryCounter, _ = otel.Meter("").Int64Histogram("client_http." + t.name + ".retry")

	return nil
}

func (t *HTTPClient) Stop() error {
	t.client.CloseIdleConnections()
	return nil
}

func (t *HTTPClient) Name() string {
	return t.name
}

func (t *HTTPClient) Type() string {
	return "client"
}

func (t *HTTPClient) Fetch(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse {
	resp := make(chan IClientResponse)

	go func(resp chan IClientResponse, ctx context.Context, method string, host string, data []byte, headers map[string][]string) {
		defer func(s time.Time) {
			t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
		}(time.Now())
		t.counter.Add(context.Background(), 1)

		res := &Response{}
		request, _ := http.NewRequestWithContext(ctx, method, host, bytes.NewBuffer(data))

		for key, value := range headers {
			if len(value) > 0 {
				request.Header.Set(key, value[0])
			}
		}

		for i := 0; i <= t.config.RetryCount; i++ {
			do, err := t.client.Do(request)

			if do != nil {
				res.Code = do.StatusCode
			}

			if err != nil {
				res.Err = err
				res.RetryCount = i + 1
				break
			}

			if do.StatusCode == 200 {
				res.RetryCount = i + 1
				res.Body = do.Body
				res.Header = map[string][]string{}
				for s, strings := range do.Header {
					res.Header[s] = strings
				}
				break
			}

			if do.StatusCode > 200 && do.StatusCode < 300 {
				res.RetryCount = i + 1
				res.Body = do.Body
				res.Header = map[string][]string{}
				for s, strings := range do.Header {
					res.Header[s] = strings
				}
				break
			}

			if (do.StatusCode >= 500 && do.StatusCode < 600) || do.StatusCode == 499 {
				if i+1 == t.config.RetryCount {
					res.RetryCount = i + 1
					res.Err = fmt.Errorf("max retry count exceeded")
				}
				continue
			}

			res.RetryCount = i + 1
			res.Err = fmt.Errorf("request status code %d", do.StatusCode)
			break
		}

		t.retryCounter.Record(context.Background(), int64(res.RetryCount))

		resp <- res
	}(resp, ctx, method, host, data, headers)

	return resp
}
