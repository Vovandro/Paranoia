package interfaces

import (
	"context"
)

type IDatabase interface {
	Init(app IEngine) error
	Stop() error
	String() string

	Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
	GetDb() interface{}
}

type SQLRow interface {
	Scan(dest ...any) error
}

type SQLRows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
