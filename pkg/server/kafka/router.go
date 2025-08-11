package kafka

import (
	"errors"
	"strings"
)

type dynamicItem struct {
	name string
	next dynamicRouter
}

type dynamicRouter struct {
	static  map[string]dynamicRouter
	dynamic []dynamicItem
	hande   RouteFunc
}

type Router struct {
	static     map[string]RouteFunc
	dynamic    dynamicRouter
	middleware map[string]IMiddleware
}

func NewRouter(middleware map[string]IMiddleware) *Router {
	return &Router{
		static:     make(map[string]RouteFunc, 5),
		dynamic:    dynamicRouter{},
		middleware: middleware,
	}
}

func (t *Router) PushRoute(path string, handler RouteFunc, middlewares []string) error {
	var md func(RouteFunc) RouteFunc = nil
	var err error

	if middlewares != nil {
		md, err = t.HandlerMiddleware(middlewares)
		if err != nil {
			return err
		}
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
		t.static[path] = h
	} else {
		p := strings.Split(path, "/")

		if len(p) > 1 && p[1] != "" {
			router := t.dynamic
			router.Push(p[1:], h)
			t.dynamic = router
		}
	}

	return nil
}

func (t *dynamicRouter) Push(path []string, handler RouteFunc) {
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

func (t *Router) Find(path string) (RouteFunc, map[string]string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if handler, ok := t.static[path]; ok {
		return handler, nil
	}

	if len(path) > 0 {
		p := strings.Split(path, "/")

		h, m := t.dynamic.Find(p[1:])

		if h != nil {
			return h, m
		}

	}

	if handler, ok := t.static["*/"]; ok {
		return handler, nil
	}

	return nil, nil
}

func (t *dynamicRouter) Find(path []string) (RouteFunc, map[string]string) {
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

func (t *Router) HandlerMiddleware(middlewares []string) (func(routeFunc RouteFunc) RouteFunc, error) {
	m := make([]IMiddleware, 0, len(middlewares))

	for _, md := range middlewares {

		if mid, ok := t.middleware[md]; ok {
			m = append(m, mid)
		} else {
			return nil, errors.New("middleware not found: " + md)
		}
	}

	return t.HandlerFromList(m), nil
}

func (t *Router) HandlerFromList(middlewares []IMiddleware) func(routeFunc RouteFunc) RouteFunc {
	return func(next RouteFunc) RouteFunc {
		if len(middlewares) == 0 {
			return next
		}

		handler := middlewares[len(middlewares)-1].Invoke(t.HandlerFromList(middlewares[:len(middlewares)-1])(next))
		for i := len(middlewares) - 2; i >= 0; i-- {
			handler = middlewares[i].Invoke(handler)
		}
		return handler
	}
}
