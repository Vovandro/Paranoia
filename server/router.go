package server

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
)

type Router struct {
	app  interfaces.IService
	data map[string]map[string]interfaces.RouteFunc
}

func NewRouter(app interfaces.IService) *Router {
	return &Router{
		app:  app,
		data: make(map[string]map[string]interfaces.RouteFunc, 5),
	}
}

func (t *Router) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	if _, ok := t.data[method]; !ok {
		t.data[method] = make(map[string]interfaces.RouteFunc, 20)
	}

	var md func(interfaces.RouteFunc) interfaces.RouteFunc = nil

	if middlewares != nil {
		md = middleware.HandlerFromStrings(t.app, middlewares)
	}

	if md != nil {
		t.data[method][path] = md(handler)
	} else {
		t.data[method][path] = handler
	}
}

func (t *Router) Find(method string, path string) interfaces.RouteFunc {
	if _, ok := t.data[method]; !ok {
		return nil
	}

	if handler, ok := t.data[method][path]; ok {
		return handler
	}

	return nil
}
