package server

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
	"net/http"
	"time"
)

type Http struct {
	Name   string
	Config HttpConfig

	app    interfaces.IService
	router *Router
	server *http.Server
	md     func(interfaces.RouteFunc) interfaces.RouteFunc
}

type HttpConfig struct {
	Port string `yaml:"port"`

	CookieDomain   string `yaml:"cookie_domain"`
	CookieSameSite string `yaml:"cookie_same_site"`
	CookieHttpOnly bool   `yaml:"cookie_http_only"`
	CookieSecure   bool   `yaml:"cookie_secure"`

	BaseMiddleware []string `yaml:"base_middleware"`
}

func NewHttp(name string, cfg HttpConfig) *Http {
	return &Http{
		Name:   name,
		Config: cfg,
	}
}

func (t *Http) Init(app interfaces.IService) error {
	t.app = app

	t.router = NewRouter(app)

	if t.Config.BaseMiddleware == nil {
		t.Config.BaseMiddleware = []string{"timing"}
	}

	if len(t.Config.BaseMiddleware) > 0 {
		t.md = middleware.HandlerFromStrings(app, t.Config.BaseMiddleware)
	}

	if t.md == nil {
		t.md = func(routeFunc interfaces.RouteFunc) interfaces.RouteFunc {
			return routeFunc
		}
	}

	t.server = &http.Server{
		Addr:                         ":" + t.Config.Port,
		Handler:                      t,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  5 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  5 * time.Second,
	}

	return nil

}

func (t *Http) Start() error {
	listenErr := make(chan error, 1)

	go func() {
		listenErr <- t.server.ListenAndServe()
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
	err := t.server.Shutdown(context.TODO())

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

func (t *Http) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := srvCtx.FromHttp(req)
	defer srvCtx.ContextPool.Put(ctx)

	route := t.router.Find(req.Method, req.RequestURI)

	if route == nil {
		ctx.Response.StatusCode = 404
		w.WriteHeader(404)
	} else {
		t.md(route)(ctx)

		w.Header().Add("Content-Type", ctx.Response.ContentType)

		for k, v := range ctx.Response.Headers {
			w.Header().Set(k, v)
		}

		for i := 0; i < len(ctx.Response.Cookie); i++ {
			w.Header().Add("Set-Cookie", ctx.Response.Cookie[i].String(t.Config.CookieDomain, t.Config.CookieSameSite, t.Config.CookieHttpOnly, t.Config.CookieSecure))
		}

		w.WriteHeader(ctx.Response.StatusCode)
		w.Write(ctx.Response.Body)
	}
}

func (t *Http) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}
