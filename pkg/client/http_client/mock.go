package http_client

import (
	"context"
	"strings"
	"sync"
)

// Mock implements IHTTPClient with hookable behavior and default queued responses.
type Mock struct {
	// Hook to override behavior
	FetchFunc func(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse

	// Seedable queue of responses used when FetchFunc is nil
	Responses []IClientResponse

	NamePkg string

	mu    sync.Mutex
	Calls []struct{ Method, Host string }
}

func (t *Mock) Init(cfg map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Type() string {
	return "client"
}

func (m *Mock) record(method, host string) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct{ Method, Host string }{Method: method, Host: host})
	m.mu.Unlock()
}

func (m *Mock) Fetch(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse {
	m.record(method, host)
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, method, host, data, headers)
	}
	ch := make(chan IClientResponse, 1)
	m.mu.Lock()
	var resp IClientResponse
	if len(m.Responses) > 0 {
		resp = m.Responses[0]
		m.Responses = m.Responses[1:]
	} else {
		resp = &Response{Code: 200, Header: map[string][]string{}, Body: strings.NewReader("")}
	}
	m.mu.Unlock()
	ch <- resp
	return ch
}

var _ IHTTPClient = (*Mock)(nil)
