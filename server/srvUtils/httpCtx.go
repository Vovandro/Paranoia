package srvUtils

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"net/http"
	"sync"
	"time"
)

type HttpCtx struct {
	request      interfaces.IRequest
	response     interfaces.IResponse
	values       map[string]interface{}
	done         chan struct{}
	routerValues map[string]string
}

var HttpCtxPool = sync.Pool{
	New: func() interface{} {
		return &HttpCtx{
			request:      &HttpRequest{},
			response:     &HttpResponse{},
			values:       make(map[string]interface{}, 10),
			routerValues: nil,
		}
	},
}

func (t *HttpCtx) Fill(request *http.Request) {
	t.request.(*HttpRequest).Fill(request)
	t.response.Clear()
	t.values = make(map[string]interface{}, 10)
	t.routerValues = nil
}

func (t *HttpCtx) GetRequest() interfaces.IRequest {
	return t.request
}

func (t *HttpCtx) GetResponse() interfaces.IResponse {
	return t.response
}

func (t *HttpCtx) Done() <-chan struct{} {
	return t.done
}

func (t *HttpCtx) StartTimeout(tm time.Duration) context.CancelFunc {
	t.done = make(chan struct{})
	go func() {
		select {
		case <-t.Done():
			return

		case <-time.After(tm):
			close(t.done)
			t.done = nil
			return
		}
	}()

	return func() {
		if t.done != nil {
			close(t.done)
		}
	}
}

func (t *HttpCtx) GetUserValue(key string) (interface{}, error) {
	if val, ok := t.values[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (t *HttpCtx) PushUserValue(key string, val interface{}) {
	t.values[key] = val
}

func (t *HttpCtx) GetRouterValue(key string) string {
	if t.routerValues == nil {
		return ""
	}

	v, _ := t.routerValues[key]

	return v
}

func (t *HttpCtx) SetRouteProps(values map[string]string) {
	t.routerValues = values
}
