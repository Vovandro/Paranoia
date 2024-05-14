package server

import (
	"Paranoia/interfaces"
)

type Mock struct {
	Name          string
	RouteRegister func(router *Router)

	app    interfaces.IService
	router *Router
}

func (t *Mock) Init(app interfaces.IService) error {
	t.app = app
	t.router = NewRouter()

	t.RouteRegister(t.router)

	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}
