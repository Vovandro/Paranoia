package http

import (
	"context"
	"fmt"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"time"
)

type TimingMiddleware struct {
	name   string
	logger interfaces2.ILogger
}

func NewTimingMiddleware(name string) interfaces2.IMiddleware {
	return &TimingMiddleware{
		name: name,
	}
}

func (t *TimingMiddleware) Init(app interfaces2.IEngine, _ map[string]interface{}) error {
	t.logger = app.GetLogger()
	return nil
}

func (t *TimingMiddleware) Stop() error {
	return nil
}

func (t *TimingMiddleware) Name() string {
	return t.name
}

func (t *TimingMiddleware) Type() string {
	return "middleware"
}

func (t *TimingMiddleware) Invoke(next RouteFunc) RouteFunc {
	return func(c context.Context, ctx ICtx) {
		tm := time.Now()

		next(c, ctx)

		s := time.Now().Sub(tm)
		ctx.PushUserValue("request_time", s)
		t.logger.Debug(c, fmt.Sprintf("%d - %v, %s: %s", ctx.GetResponse().GetStatus(), s, ctx.GetRequest().GetMethod(), ctx.GetRequest().GetURI()))
	}
}
