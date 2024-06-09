package client

import (
	"bytes"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"net/http"
	"strconv"
)

type HTTPClient struct {
	Name   string
	Config HTTPClientConfig
	app    interfaces.IService
	client http.Client
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

func (t *HTTPClient) Init(app interfaces.IService) error {
	t.app = app
	t.client = http.Client{}

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
				res.Header = map[string][]string{}
				size, _ := strconv.ParseInt(do.Header.Get("Content-Length"), 10, 64)
				res.Body = make([]byte, size)
				do.Body.Read(res.Body)
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

		resp <- res
	}(resp, method, host, data, headers)

	return resp
}
