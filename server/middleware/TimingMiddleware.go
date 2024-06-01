package middleware

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
	"time"
)

type TimingMiddleware struct {
	logger interfaces.ILogger
}

func (t *TimingMiddleware) Init(app interfaces.IService) error {
	t.logger = app.GetLogger()
	return nil
}

func (t *TimingMiddleware) Stop() error {
	return nil
}

func (t *TimingMiddleware) String() string {
	return "timing"
}

func (t *TimingMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(ctx *srvCtx.Ctx) {
		tm := time.Now()

		next(ctx)

		s := time.Now().Sub(tm)
		ctx.Values["request_time"] = s
		t.logger.Debug(fmt.Sprintf("%d - %v, %s: %s", ctx.Response.StatusCode, s, ctx.Request.Method, ctx.Request.URI))
	}
}
