package interfaces

import (
	"context"
)

type RouteFunc func(c context.Context, ctx ICtx)

type IServer interface {
	Init(app IEngine) error
	Start() error
	Stop() error
	String() string
}
