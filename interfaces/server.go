package interfaces

import (
	"context"
	"google.golang.org/grpc"
)

type RouteFunc func(c context.Context, ctx ICtx)

type IServer interface {
	Init(app IEngine) error
	Start() error
	Stop() error
	String() string
}

type IServerBase interface {
	IServer
	PushRoute(method string, path string, handler RouteFunc, middlewares []string)
}

type IServerGRPC interface {
	IServer
	RegisterService(desc *grpc.ServiceDesc, impl any)
}
