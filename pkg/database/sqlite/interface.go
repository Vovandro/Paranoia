package sqlite

import (
	"context"
	"database/sql"
)

// ISQLite defines the interface for SQLite operations
type ISQLite interface {
	// Query executes a query and returns multiple rows
	Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)

	// QueryRow executes a query and returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)

	// Exec executes a query without returning any rows
	Exec(ctx context.Context, query string, args ...interface{}) error

	// GetDb returns the SQLite client instance
	GetDb() *sql.DB
}

type SQLRow interface {
	Scan(dest ...any) error
}

type SQLRows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
