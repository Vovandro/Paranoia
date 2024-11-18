package middleware

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

func HandlerFromList(middlewares []interfaces.IMiddleware) func(routeFunc interfaces.RouteFunc) interfaces.RouteFunc {
	return func(next interfaces.RouteFunc) interfaces.RouteFunc {
		if len(middlewares) == 0 {
			return next
		}

		handler := middlewares[len(middlewares)-1].Invoke(HandlerFromList(middlewares[:len(middlewares)-1])(next))
		for i := len(middlewares) - 2; i >= 0; i-- {
			handler = middlewares[i].Invoke(handler)
		}
		return handler
	}
}

func HandlerFromStrings(app interfaces.IEngine, middlewares []string) func(routeFunc interfaces.RouteFunc) interfaces.RouteFunc {
	m := make([]interfaces.IMiddleware, 0, len(middlewares))

	for _, middleware := range middlewares {
		t := app.GetMiddleware(middleware)

		if t == nil {
			app.GetLogger().Warn(context.Background(), "middleware not found: "+middleware)
			continue
		}
		m = append(m, t)
	}

	return HandlerFromList(m)
}
