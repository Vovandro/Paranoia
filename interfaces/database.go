package interfaces

import "context"

type IDatabase interface {
	Init(app IService) error
	Stop() error
	String() string

	Exists(ctx context.Context, query interface{}, args ...interface{}) bool
	Count(ctx context.Context, query interface{}, args ...interface{}) int64
	FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Update(ctx context.Context, query interface{}, args ...interface{}) error
	Delete(ctx context.Context, query interface{}, args ...interface{}) int64
	GetDb() interface{}
}
