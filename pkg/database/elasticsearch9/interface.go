package elasticsearch9

import (
	"context"
)

// keep surface similar; v9 may not have typedapi stable
type IElasticSearch interface {
	Index(ctx context.Context, index string, id string, document interface{}, refresh bool) (string, error)
	Get(ctx context.Context, index string, id string) (NoSQLRow, error)
	Search(ctx context.Context, index []string, query map[string]any, from, size int) (NoSQLRows, error)
	// SearchSource performs a search query with _source filtering (includes/excludes)
	SearchSource(ctx context.Context, index []string, query map[string]any, from, size int, include, exclude []string) (NoSQLRows, error)
	Delete(ctx context.Context, index string, id string, refresh bool) error
	DeleteByQuery(ctx context.Context, index []string, query map[string]any, refresh bool) error
	Update(ctx context.Context, index string, id string, doc interface{}, refresh bool) error
	GetClient() interface{}
}

type NoSQLRow interface{ Scan(dest any) error }
type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}
