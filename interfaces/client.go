package interfaces

import "context"

type IClient interface {
	Init(app IEngine) error
	Stop() error
	String() string

	Fetch(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse
}

type IClientResponse interface {
	GetBody() []byte
	GetHeader() map[string][]string
	Error() error
	GetRetries() int
}
