package interfaces

import "gitlab.com/devpro_studio/Paranoia/srvCtx"

type RouteFunc func(ctx *srvCtx.Ctx)

type IServer interface {
	Init(app IService) error
	Start() error
	Stop() error
	String() string
	PushRoute(method string, path string, handler RouteFunc, middlewares []string)
}
