package server

import (
	"Paranoia/interfaces"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Mock struct {
	Name          string
	RouteRegister func(router *router.Router)

	app    interfaces.IService
	router *router.Router
}

func (t *Mock) Init(app interfaces.IService) error {
	t.app = app
	t.router = router.New()

	t.RouteRegister(t.router)

	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return next
}

func (t *Mock) Get(body []byte, header []string) {
	r := fasthttp.RequestCtx{
		Request:  fasthttp.Request{},
		Response: fasthttp.Response{},
	}

	r.Request.SetBody(body)

	t.Handle(func(ctx *fasthttp.RequestCtx) {})(&r)
}
