package elasticsearch8

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// IElasticSearch defines the interface for Elasticsearch operations
type IElasticSearch interface {
	// Index indexes a document and returns the document ID
	Index(ctx context.Context, index string, id string, document interface{}, refresh bool) (string, error)

	// Get retrieves a single document by ID
	Get(ctx context.Context, index string, id string) (NoSQLRow, error)

	// Search performs a search query and returns multiple hits
	Search(ctx context.Context, index []string, query *types.Query, from, size int) (NoSQLRows, error)

	// Delete removes a document by ID
	Delete(ctx context.Context, index string, id string, refresh bool) error

	// Update performs a partial update on a document
	Update(ctx context.Context, index string, id string, doc interface{}, refresh bool) error

	// GetClient returns the underlying client instance
	GetClient() interface{}
}

type NoSQLRow interface {
	Scan(dest any) error
}

type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}
