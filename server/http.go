package server

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/server/srvUtils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
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

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
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

	t.counter, _ = otel.Meter("").Int64Counter("server_http." + t.Name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_http." + t.Name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_http." + t.Name + ".time")

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
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	ctx := srvUtils.HttpCtxPool.Get().(*srvUtils.HttpCtx)
	defer srvUtils.HttpCtxPool.Put(ctx)
	ctx.Fill(req)

	route := t.router.Find(req.Method, req.RequestURI)

	if route == nil {
		ctx.GetResponse().SetStatus(404)
		w.WriteHeader(404)
	} else {
		t.md(route)(ctx)

		header := ctx.GetResponse().Header().GetAsMap()

		if _, ok := header["Content-Type"]; !ok {
			header["Content-Type"] = []string{"application/json; charset=utf-8"}
		}

		for k, v := range ctx.GetResponse().Header().GetAsMap() {
			for _, v2 := range v {
				w.Header().Set(k, v2)
			}
		}

		cookie := ctx.GetResponse().Cookie().(*srvUtils.HttpCookie).ToHttp(t.Config.CookieDomain, t.Config.CookieSameSite, t.Config.CookieHttpOnly, t.Config.CookieSecure)

		for i := 0; i < len(cookie); i++ {
			w.Header().Add("Set-Cookie", cookie[i])
		}

		w.WriteHeader(ctx.GetResponse().GetStatus())
		w.Write(ctx.GetResponse().GetBody())
	}

	if ctx.GetResponse().GetStatus() >= 400 {
		t.counterError.Add(context.Background(), 1)
	}
}

func (t *Http) PushRoute(method string, path string, handler interfaces.RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}
