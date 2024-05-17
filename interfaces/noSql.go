package interfaces

import "context"

type INoSql interface {
	Init(app IService) error
	Stop() error
	String() string

	Count(ctx context.Context, query interface{}, args ...interface{}) int64
	Exists(ctx context.Context, query interface{}, args ...interface{}) bool
	FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Update(ctx context.Context, query interface{}, args ...interface{}) error
	Delete(ctx context.Context, query interface{}, args ...interface{}) int64
	Batch(ctx context.Context, typeOp string, query interface{}, args ...interface{}) (int64, error)
	GetDb() interface{}
}
