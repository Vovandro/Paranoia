package interfaces

import (
	"context"
	"google.golang.org/grpc"
)

type IClient interface {
	Init(app IEngine) error
	Stop() error
	String() string
}

type IClientBase interface {
	IClient
	Fetch(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse
}

type IClientGrpc interface {
	IClient
	GetClient() *grpc.ClientConn
}

type IClientResponse interface {
	GetBody() []byte
	GetHeader() map[string][]string
	Error() error
	GetRetries() int
}
