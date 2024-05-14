package server

import (
	"context"
	"fmt"
	context2 "gitlab.com/devpro_studio/Paranoia/context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"net/http"
	"time"
)

type Http struct {
	Name string
	Port string

	app    interfaces.IService
	router *Router
	server *http.Server
}

func (t *Http) Init(app interfaces.IService) error {
	t.app = app

	t.router = NewRouter()

	t.server = &http.Server{
		Addr:                         ":" + t.Port,
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
	ctx := context2.FromHttp(req)
	defer func(tm time.Time) {
		t.app.GetLogger().Debug(fmt.Sprintf("[%d] [%v] %s: %s", ctx.Response.StatusCode, time.Now().Sub(tm), req.Method, req.RequestURI))
	}(time.Now())

	route := t.router.Find(req.Method, req.RequestURI)

	if route == nil {
		ctx.Response.StatusCode = 404
		w.WriteHeader(404)
	} else {
		route(ctx)
	}
}

func (t *Http) PushRoute(method string, path string, handler interfaces.RouteFunc) {
	t.router.PushRoute(method, path, handler)
}
