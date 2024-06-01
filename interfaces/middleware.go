package interfaces

type IMiddleware interface {
	Init(app IService) error
	Stop() error
	String() string
	Invoke(next RouteFunc) RouteFunc
}
