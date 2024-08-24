package interfaces

type IMiddleware interface {
	Init(app IEngine) error
	Stop() error
	String() string
	Invoke(next RouteFunc) RouteFunc
}
