package rabbitmq_client

import (
	"context"
	"io"
)

type IRabbitmqClient interface {
	// Fetch sends a message to the specified RabbitMQ topic and returns a channel with the response
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
