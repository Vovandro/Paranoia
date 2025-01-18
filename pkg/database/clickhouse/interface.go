package clickhouse

import (
	"context"
)

// IClickHouse defines the interface for ClickHouse operations
type IClickHouse interface {
	// Query executes a query and returns multiple rows
	Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)

	// QueryRow executes a query and returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)

	// Exec executes a query without returning any rows
	Exec(ctx context.Context, query string, args ...interface{}) error

	// GetDb returns the ClickHouse client instance
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
