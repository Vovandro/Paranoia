package http

import (
	"fmt"
	"net/http"
	"sync"
)

type HttpCtx struct {
	request      IRequest
	response     IResponse
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

func (t *HttpCtx) GetRequest() IRequest {
	return t.request
}

func (t *HttpCtx) GetResponse() IResponse {
	return t.response
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
