package rabbitmq_client

import (
	"context"
	"strings"
	"sync"
)

// Mock implements IRabbitmqClient with hooks and default queued responses.
type Mock struct {
	FetchFunc func(ctx context.Context, topic string, data []byte, headers map[string][]string) chan IClientResponse

	Responses []IClientResponse

	NamePkg string

	mu    sync.Mutex
	Calls []struct{ Topic string }
}

func (m *Mock) record(topic string) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct{ Topic string }{Topic: topic})
	m.mu.Unlock()
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

func (m *Mock) Fetch(ctx context.Context, topic string, data []byte, headers map[string][]string) chan IClientResponse {
	m.record(topic)
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, topic, data, headers)
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

var _ IRabbitmqClient = (*Mock)(nil)
