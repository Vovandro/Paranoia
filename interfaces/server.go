package interfaces

type RouteFunc func(ctx ICtx)

type IServer interface {
	Init(app IService) error
	Start() error
	Stop() error
	String() string
	PushRoute(method string, path string, handler RouteFunc, middlewares []string)
}
