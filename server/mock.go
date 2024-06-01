package server

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Mock struct {
	Name          string
	RouteRegister func(router *Router)

	app    interfaces.IService
	router *Router
}

func (t *Mock) Init(app interfaces.IService) error {
	t.app = app
	t.router = NewRouter(app)

	t.RouteRegister(t.router)

	return nil
}

func (t *Mock) Start() error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}
