package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type RabbitmqCtx struct {
	request  IRequest
	response IResponse
	values   map[string]interface{}
	done     chan struct{}
}

var RabbitmqCtxPool = sync.Pool{
	New: func() interface{} {
		return &RabbitmqCtx{
			request:  &rabbitmqRequest{},
			response: &RabbitmqResponse{},
			values:   make(map[string]interface{}, 10),
		}
	},
}

func (t *RabbitmqCtx) Fill(msg *amqp.Delivery) {
	t.request.(*rabbitmqRequest).Fill(msg)
	t.response.Clear()
	t.values = make(map[string]interface{}, 10)
}

func (t *RabbitmqCtx) GetRequest() IRequest {
	return t.request
}

func (t *RabbitmqCtx) GetResponse() IResponse {
	return t.response
}

func (t *RabbitmqCtx) GetUserValue(key string) (interface{}, error) {
	if val, ok := t.values[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (t *RabbitmqCtx) PushUserValue(key string, val interface{}) {
	t.values[key] = val
}

func (t *RabbitmqCtx) GetRouterValue(key string) string {
	return ""
}

func (t *RabbitmqCtx) SetRouteProps(values map[string]string) {}
