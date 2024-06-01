package middleware

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/srvCtx"
)

type RestoreMiddleware struct {
	logger interfaces.ILogger
}

func (t *RestoreMiddleware) Init(app interfaces.IService) error {
	t.logger = app.GetLogger()
	return nil
}

func (t *RestoreMiddleware) Stop() error {
	return nil
}

func (t *RestoreMiddleware) String() string {
	return "restore"
}

func (t *RestoreMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(ctx *srvCtx.Ctx) {
		defer func() {
			if err := recover(); err != nil {
				t.logger.Error(fmt.Errorf("%v", err))
			}
		}()

		next(ctx)
	}
}
