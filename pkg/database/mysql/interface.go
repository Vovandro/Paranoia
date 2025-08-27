package mysql

import (
	"context"
	"database/sql"
)

// IMySQL defines the interface for MySQL operations
type IMySQL interface {
	// Query executes a query and returns multiple rows
	Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)

	// QueryRow executes a query and returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)

	// Exec executes a query without returning any rows
	Exec(ctx context.Context, query string, args ...interface{}) error

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (SQLTx, error)

	// GetDb returns the MySQL client instance
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

// SQLTx represents a SQL transaction
type SQLTx interface {
	// Query executes a query and returns multiple rows within the transaction
	Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)

	// QueryRow executes a query and returns a single row within the transaction
	QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)

	// Exec executes a query without returning any rows within the transaction
	Exec(ctx context.Context, query string, args ...interface{}) error

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error
}
