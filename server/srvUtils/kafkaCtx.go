package srvUtils

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"sync"
	"time"
)

type KafkaCtx struct {
	request  interfaces.IRequest
	response interfaces.IResponse
	values   map[string]interface{}
	done     chan struct{}
}

var KafkaCtxPool = sync.Pool{
	New: func() interface{} {
		return &HttpCtx{
			request:  &kafkaRequest{},
			response: &HttpResponse{},
			values:   make(map[string]interface{}, 10),
		}
	},
}

func (t *KafkaCtx) Fill(msg *kafka.Message) {
	t.request.(*kafkaRequest).Fill(msg)
	t.response.Clear()
	t.values = make(map[string]interface{}, 10)
}

func (t *KafkaCtx) GetRequest() interfaces.IRequest {
	return t.request
}

func (t *KafkaCtx) GetResponse() interfaces.IResponse {
	return t.response
}

func (t *KafkaCtx) Done() <-chan struct{} {
	return t.done
}

func (t *KafkaCtx) StartTimeout(tm time.Duration) context.CancelFunc {
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

func (t *KafkaCtx) GetUserValue(key string) (interface{}, error) {
	if val, ok := t.values[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (t *KafkaCtx) PushUserValue(key string, val interface{}) {
	t.values[key] = val
}

func (t *KafkaCtx) GetRouterValue(key string) string {
	return ""
}

func (t *KafkaCtx) SetRouteProps(values map[string]string) {}
