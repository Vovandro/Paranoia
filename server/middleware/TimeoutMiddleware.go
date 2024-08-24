package middleware

import (
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
	return func(ctx interfaces.ICtx) {
		cancel := ctx.StartTimeout(t.Config.Timeout)
		c := make(chan interface{})
		defer close(c)
		defer cancel()

		go func() {
			next(ctx)

			if _, ok := <-c; ok {
				c <- nil
			}
		}()

		select {
		case <-ctx.Done():
			time.Sleep(time.Millisecond)
			ctx.GetResponse().SetStatus(499)
			break

		case <-c:
			break
		}
	}
}
