package server

type RouteFunc func(ctx *Context)

type Router struct {
	data map[string]map[string]RouteFunc
}

func NewRouter() *Router {
	return &Router{
		data: make(map[string]map[string]RouteFunc, 5),
	}
}

func (t *Router) PushRoute(method string, path string, handler RouteFunc) {
	if _, ok := t.data[method]; !ok {
		t.data[method] = make(map[string]RouteFunc, 20)
	}

	t.data[method][path] = handler
}

func (t *Router) Find(method string, path string) RouteFunc {
	if _, ok := t.data[method]; !ok {
		return nil
	}

	if handler, ok := t.data[method][path]; ok {
		return handler
	}

	return nil
}

func (t *Router) Get(path string, handler RouteFunc) {
	t.PushRoute("GET", path, handler)
}

func (t *Router) Post(path string, handler RouteFunc) {
	t.PushRoute("POST", path, handler)
}

func (t *Router) Put(path string, handler RouteFunc) {
	t.PushRoute("PUT", path, handler)
}

func (t *Router) Delete(path string, handler RouteFunc) {
	t.PushRoute("DELETE", path, handler)
}

func (t *Router) Rpc(path string, handler RouteFunc) {
	t.PushRoute("RPC", path, handler)
}
