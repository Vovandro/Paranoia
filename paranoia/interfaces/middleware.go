package interfaces

type IMiddleware interface {
	Init(app IEngine, cfg map[string]interface{}) error
	Stop() error
	Name() string
	Type() string
	Invoke(next RouteFunc) RouteFunc
}
