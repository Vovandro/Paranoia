package interfaces

import "gitlab.com/devpro_studio/Paranoia/context"

type RouteFunc func(ctx *context.Context)

type IServer interface {
	Init(app IService) error
	Stop() error
	String() string
	PushRoute(method string, path string, handler RouteFunc)
}
