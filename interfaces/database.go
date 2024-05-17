package interfaces

import "context"

type IDatabase interface {
	Init(app IService) error
	Stop() error
	String() string

	Query(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error
	Exec(ctx context.Context, query interface{}, args ...interface{}) error
	GetDb() interface{}
}
