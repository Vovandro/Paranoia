package server

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"strings"
)

type dynamicItem struct {
	name string
	next dynamicRouter
}

type dynamicRouter struct {
	static  map[string]dynamicRouter
	dynamic []dynamicItem
	hande   interfaces.RouteFunc
}

type Router struct {
	app     interfaces.IService
	static  map[string]map[string]interfaces.RouteFunc
	dynamic map[string]dynamicRouter
}

func NewRouter(app interfaces.IService) *Router {
	return &Router{
		app:     app,
		static:  make(map[string]map[string]interfaces.RouteFunc, 5),
		dynamic: make(map[string]dynamicRouter, 5),
	}
}

func (t *Router) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	var md func(interfaces.RouteFunc) interfaces.RouteFunc = nil

	if middlewares != nil {
		md = middleware.HandlerFromStrings(t.app, middlewares)
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	idx := strings.Index(path, "{")

	h := handler

	if md != nil {
		h = md(handler)
	}

	if idx == -1 {
		if _, ok := t.static[method]; !ok {
			t.static[method] = make(map[string]interfaces.RouteFunc, 20)
		}

		t.static[method][path] = h
	} else {
		p := strings.Split(path, "/")

		if len(p) > 1 && p[1] != "" {
			if _, ok := t.dynamic[method]; !ok {
				t.dynamic[method] = dynamicRouter{
					static:  make(map[string]dynamicRouter, 5),
					dynamic: make([]dynamicItem, 0, 5),
				}
			}

			router := t.dynamic[method]
			router.Push(p[1:], h)
		}
	}
}

func (t *dynamicRouter) Push(path []string, handler interfaces.RouteFunc) {
	if len(path) == 0 || path[0] == "" {
		t.hande = handler
		return
	}

	if strings.HasPrefix(path[0], "{") && strings.HasSuffix(path[0], "}") {
		r := dynamicItem{
			name: path[0][1 : len(path[0])-1],
			next: dynamicRouter{
				static:  make(map[string]dynamicRouter, 5),
				dynamic: make([]dynamicItem, 0, 5),
			},
		}

		r.next.Push(path[1:], handler)
		t.dynamic = append(t.dynamic, r)
		return
	}

	r := dynamicRouter{
		static:  make(map[string]dynamicRouter, 5),
		dynamic: make([]dynamicItem, 0, 5),
	}

	r.Push(path[1:], handler)

	t.static[path[0]] = r
}

func (t *Router) Find(method string, path string) (interfaces.RouteFunc, map[string]string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if _, ok := t.static[method]; ok {
		if handler, ok := t.static[method][path]; ok {
			return handler, nil
		}
	}

	if len(path) > 0 {
		p := strings.Split(path, "/")

		if dr, ok := t.dynamic[method]; ok {
			h, m := dr.Find(p[1:])

			if h != nil {
				return h, m
			}
		}
	}

	if _, ok := t.static[method]; ok {
		if handler, ok := t.static[method]["*/"]; ok {
			return handler, nil
		}
	}

	return nil, nil
}

func (t *dynamicRouter) Find(path []string) (interfaces.RouteFunc, map[string]string) {
	if len(path) == 0 || path[0] == "" {
		if t.hande != nil {
			return t.hande, map[string]string{}
		}

		return nil, nil
	}

	if v, ok := t.static[path[0]]; ok {
		r, m := v.Find(path[1:])

		if r != nil {
			return r, m
		}
	}

	for _, v := range t.dynamic {
		r, m := v.next.Find(path[1:])

		if r != nil {
			m[v.name] = path[0]
			return r, m
		}
	}

	return nil, nil
}
