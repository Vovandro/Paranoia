package server

import (
	"Paranoia/interfaces"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"time"
)

type Http struct {
	Name          string
	Port          string
	RouteRegister func(router *router.Router)

	app    interfaces.IService
	router *router.Router
	server *fasthttp.Server
}

func (t *Http) Init(app interfaces.IService) error {
	t.app = app

	t.router = router.New()

	t.router.RedirectTrailingSlash = true
	t.router.RedirectFixedPath = true

	t.server = &fasthttp.Server{
		Handler:      t.Handle(t.router.Handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	t.RouteRegister(t.router)

	listenErr := make(chan error, 1)

	go func() {
		listenErr <- t.server.ListenAndServe(t.Port)
	}()

	select {
	case err := <-listenErr:
		t.app.GetLogger().Error(err)
		return err

	case <-time.After(time.Second):
		// pass
	}

	return nil

}

func (t *Http) Stop() error {
	t.server.DisableKeepalive = true

	err := t.server.Shutdown()

	if err != nil {
		t.app.GetLogger().Error(err)
	} else {
		t.app.GetLogger().Info("http server gracefully stopped.")
		time.Sleep(time.Second)
	}

	return err
}

func (t *Http) String() string {
	return t.Name
}

func (t *Http) Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.CompressHandler(
		func(ctx *fasthttp.RequestCtx) {
			defer func(tm time.Time) {
				t.app.GetLogger().Debug(fmt.Sprintf("[%d] [%v] %s: %s", ctx.Response.StatusCode(), time.Now().Sub(tm), ctx.Method(), ctx.RequestURI()))
			}(time.Now())

			if ctx.IsOptions() {
				t.router.GlobalOPTIONS(ctx)
				return
			}

			next(ctx)
		})
}
