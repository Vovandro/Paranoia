package kafka_client

import (
	"context"
	"io"
)

type IKafkaClient interface {
	// Fetch sends a Kafka request and returns a channel with the response
	Fetch(ctx context.Context, topic string, data []byte, headers map[string][]string) chan IClientResponse
}

type IClientResponse interface {
	GetBody() ([]byte, error)
	GetLazyBody() io.Reader
	GetHeader() map[string][]string
	Error() error
	GetRetries() int
	GetCode() int
}
