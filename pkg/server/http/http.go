package http

import (
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"time"
)

type Http struct {
	name   string
	config Config

	router *Router
	server *http.Server
	md     func(RouteFunc) RouteFunc

	counter      metric.Int64Counter
	counterError metric.Int64Counter
	timeCounter  metric.Int64Histogram
}

type Config struct {
	Port string `yaml:"port"`

	CookieDomain   string `yaml:"cookie_domain"`
	CookieSameSite string `yaml:"cookie_same_site"`
	CookieHttpOnly bool   `yaml:"cookie_http_only"`
	CookieSecure   bool   `yaml:"cookie_secure"`

	BaseMiddleware []string `yaml:"base_middleware"`
}

func NewHttp(name string) *Http {
	return &Http{
		name: name,
	}
}

func (t *Http) Init(cfg map[string]interface{}) error {
	middlewares := make(map[string]IMiddleware)

	if m, ok := cfg["middlewares"]; ok {
		middlewares = m.(map[string]IMiddleware)
		delete(cfg, "middlewares")
	}

	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Port == "" {
		t.config.Port = "80"
	}

	t.router = NewRouter(middlewares)

	if t.config.BaseMiddleware == nil {
		t.config.BaseMiddleware = []string{}
	}

	if len(t.config.BaseMiddleware) > 0 {
		t.md, err = t.router.HandlerMiddleware(t.config.BaseMiddleware)
		if err != nil {
			return err
		}
	}

	if t.md == nil {
		t.md = func(routeFunc RouteFunc) RouteFunc {
			return routeFunc
		}
	}

	t.server = &http.Server{
		Addr:                         ":" + t.config.Port,
		Handler:                      otelhttp.NewHandler(t, t.name),
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  5 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  5 * time.Second,
	}

	t.counter, _ = otel.Meter("").Int64Counter("server_http." + t.name + ".count")
	t.counterError, _ = otel.Meter("").Int64Counter("server_http." + t.name + ".count_error")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("server_http." + t.name + ".time")

	return nil

}

func (t *Http) Start() error {
	listenErr := make(chan error, 1)

	go func() {
		listenErr <- t.server.ListenAndServe()
	}()

	select {
	case err := <-listenErr:
		return err

	case <-time.After(time.Second * 5):
		// pass
	}

	return nil
}

func (t *Http) Stop() error {
	err := t.server.Shutdown(context.TODO())

	time.Sleep(time.Second)

	return err
}

func (t *Http) Name() string {
	return t.name
}

func (t *Http) Type() string {
	return "server"
}

func (t *Http) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	ctx := HttpCtxPool.Get().(*HttpCtx)
	defer HttpCtxPool.Put(ctx)
	ctx.Fill(req)

	route, props := t.router.Find(req.Method, req.URL.Path)

	if route == nil {
		ctx.GetResponse().SetStatus(404)
		w.WriteHeader(404)
	} else {
		ctx.SetRouteProps(props)

		t.md(route)(req.Context(), ctx)

		header := ctx.GetResponse().Header().GetAsMap()

		if _, ok := header["Content-Type"]; !ok {
			header["Content-Type"] = []string{"application/json; charset=utf-8"}
		}

		for k, v := range ctx.GetResponse().Header().GetAsMap() {
			for _, v2 := range v {
				w.Header().Set(k, v2)
			}
		}

		cookie := ctx.GetResponse().Cookie().(*HttpCookie).ToHttp(t.config.CookieDomain, t.config.CookieSameSite, t.config.CookieHttpOnly, t.config.CookieSecure)

		for i := 0; i < len(cookie); i++ {
			w.Header().Add("Set-Cookie", cookie[i])
		}

		body := ctx.GetResponse().GetBody()
		status := ctx.GetResponse().GetStatus()
		if body != nil && len(body) > 0 {
			w.WriteHeader(status)
			w.Write(body)
		} else {
			if status == http.StatusOK {
				status = http.StatusNoContent
			}

			w.WriteHeader(status)
		}
	}

	if ctx.GetResponse().GetStatus() >= 400 {
		t.counterError.Add(context.Background(), 1)
	}
}

func (t *Http) PushRoute(method string, path string, handler RouteFunc, middlewares []string) {
	t.router.PushRoute(method, path, handler, middlewares)
}
