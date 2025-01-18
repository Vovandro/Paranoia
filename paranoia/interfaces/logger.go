package interfaces

import (
	"context"
)

type ILogger interface {
	Init(map[string]interface{}) error
	Stop() error
	Name() string
	Type() string
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Message(ctx context.Context, args ...interface{})
	Error(ctx context.Context, err error)
	Fatal(ctx context.Context, err error)
	Panic(ctx context.Context, err error)
	Parent() interface{}
	SetParent(interface{})
}
