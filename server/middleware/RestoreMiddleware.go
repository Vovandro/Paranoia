package middleware

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type RestoreMiddleware struct {
	Name   string
	logger interfaces.ILogger
}

func NewRestoreMiddleware(name string) interfaces.IMiddleware {
	return &RestoreMiddleware{
		Name: name,
	}
}

func (t *RestoreMiddleware) Init(app interfaces.IEngine) error {
	t.logger = app.GetLogger()
	return nil
}

func (t *RestoreMiddleware) Stop() error {
	return nil
}

func (t *RestoreMiddleware) String() string {
	return t.Name
}

func (t *RestoreMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(ctx interfaces.ICtx) {
		defer func() {
			if err := recover(); err != nil {
				t.logger.Error(fmt.Errorf("%v", err))
			}
		}()

		next(ctx)
	}
}
