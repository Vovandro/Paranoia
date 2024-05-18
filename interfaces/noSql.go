package interfaces

import "context"

type INoSql interface {
	Init(app IService) error
	Stop() error
	String() string

	Count(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64
	Exists(ctx context.Context, key interface{}, query interface{}, args ...interface{}) bool
	Insert(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interface{}, error)
	FindOne(ctx context.Context, key interface{}, query interface{}, model interface{}, args ...interface{}) error
	Find(ctx context.Context, key interface{}, query interface{}, model interface{}, args ...interface{}) error
	Exec(ctx context.Context, key interface{}, query interface{}, model interface{}, args ...interface{}) error
	Update(ctx context.Context, key interface{}, query interface{}, args ...interface{}) error
	Delete(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64
	Batch(ctx context.Context, key interface{}, query interface{}, typeOp string, args ...interface{}) (int64, error)
	GetDb() interface{}
}
