package middleware

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type TimeoutMiddleware struct {
	Name   string
	Config TimeoutMiddlewareConfig
}

type TimeoutMiddlewareConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

func NewTimeoutMiddleware(name string, cfg TimeoutMiddlewareConfig) interfaces.IMiddleware {
	return &TimeoutMiddleware{
		Name:   name,
		Config: cfg,
	}
}

func (t *TimeoutMiddleware) Init(app interfaces.IEngine) error {
	return nil
}

func (t *TimeoutMiddleware) Stop() error {
	return nil
}

func (t *TimeoutMiddleware) String() string {
	return t.Name
}

func (t *TimeoutMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(c context.Context, ctx interfaces.ICtx) {
		var end context.CancelFunc
		c, end = context.WithTimeout(c, t.Config.Timeout)

		done := make(chan interface{})
		defer close(done)
		defer end()

		go func() {
			next(c, ctx)

			if _, ok := <-c.Done(); ok {
				done <- nil
			}
		}()

		select {
		case <-c.Done():
			time.Sleep(time.Millisecond)
			ctx.GetResponse().SetStatus(499)
			break

		case <-done:
			break
		}
	}
}
