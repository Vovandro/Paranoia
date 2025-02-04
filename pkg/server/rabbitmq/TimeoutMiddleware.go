package rabbitmq

import (
	"context"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
	"time"
)

type TimeoutMiddleware struct {
	name   string
	config TimeoutMiddlewareConfig
}

type TimeoutMiddlewareConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

func NewTimeoutMiddleware(name string) interfaces2.IMiddleware {
	return &TimeoutMiddleware{
		name: name,
	}
}

func (t *TimeoutMiddleware) Init(app interfaces2.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Timeout == 0 {
		t.config.Timeout = time.Second
	}

	return nil
}

func (t *TimeoutMiddleware) Stop() error {
	return nil
}

func (t *TimeoutMiddleware) Name() string {
	return t.name
}

func (t *TimeoutMiddleware) Type() string {
	return "middleware"
}

func (t *TimeoutMiddleware) Invoke(next RouteFunc) RouteFunc {
	return func(c context.Context, ctx ICtx) {
		var end context.CancelFunc
		c, end = context.WithTimeout(c, t.config.Timeout)

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
