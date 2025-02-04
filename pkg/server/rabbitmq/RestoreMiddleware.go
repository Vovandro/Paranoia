package rabbitmq

import (
	"context"
	"fmt"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
)

type RestoreMiddleware struct {
	name   string
	logger interfaces2.ILogger
}

func NewRestoreMiddleware(name string) interfaces2.IMiddleware {
	return &RestoreMiddleware{
		name: name,
	}
}

func (t *RestoreMiddleware) Init(app interfaces2.IEngine, _ map[string]interface{}) error {
	t.logger = app.GetLogger()
	return nil
}

func (t *RestoreMiddleware) Stop() error {
	return nil
}

func (t *RestoreMiddleware) Name() string {
	return t.name
}

func (t *RestoreMiddleware) Type() string {
	return "middleware"
}

func (t *RestoreMiddleware) Invoke(next RouteFunc) RouteFunc {
	return func(c context.Context, ctx ICtx) {
		defer func() {
			if err := recover(); err != nil {
				t.logger.Error(context.Background(), fmt.Errorf("%v", err))
			}
		}()

		next(c, ctx)
	}
}
