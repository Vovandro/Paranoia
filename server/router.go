package server

import "gitlab.com/devpro_studio/Paranoia/interfaces"

type Router struct {
	data map[string]map[string]interfaces.RouteFunc
}

func NewRouter() *Router {
	return &Router{
		data: make(map[string]map[string]interfaces.RouteFunc, 5),
	}
}

func (t *Router) PushRoute(method string, path string, handler interfaces.RouteFunc) {
	if _, ok := t.data[method]; !ok {
		t.data[method] = make(map[string]interfaces.RouteFunc, 20)
	}

	t.data[method][path] = handler
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
