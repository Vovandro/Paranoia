package middleware

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type TimingMiddleware struct {
	name   string
	logger interfaces.ILogger
}

func NewTimingMiddleware(name string) interfaces.IMiddleware {
	return &TimingMiddleware{
		name: name,
	}
}

func (t *TimingMiddleware) Init(app interfaces.IEngine, _ map[string]interface{}) error {
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

func (t *TimingMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(c context.Context, ctx interfaces.ICtx) {
		tm := time.Now()

		next(c, ctx)

		s := time.Now().Sub(tm)
		ctx.PushUserValue("request_time", s)
		t.logger.Debug(c, fmt.Sprintf("%d - %v, %s: %s", ctx.GetResponse().GetStatus(), s, ctx.GetRequest().GetMethod(), ctx.GetRequest().GetURI()))
	}
}
