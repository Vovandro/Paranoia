package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"sync"
)

type KafkaCtx struct {
	request  IRequest
	response IResponse
	values   map[string]interface{}
	done     chan struct{}
}

var KafkaCtxPool = sync.Pool{
	New: func() interface{} {
		return &KafkaCtx{
			request:  &kafkaRequest{},
			response: &KafkaResponse{},
			values:   make(map[string]interface{}, 10),
		}
	},
}

func (t *KafkaCtx) Fill(msg *kafka.Message) {
	t.request.(*kafkaRequest).Fill(msg)
	t.response.Clear()
	t.values = make(map[string]interface{}, 10)
}

func (t *KafkaCtx) GetRequest() IRequest {
	return t.request
}

func (t *KafkaCtx) GetResponse() IResponse {
	return t.response
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
