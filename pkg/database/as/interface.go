package as

import (
	"context"
	"github.com/aerospike/aerospike-client-go/v7"
)

// IAerospike defines the interface for Aerospike operations
type IAerospike interface {
	// Exists checks if a key exists in the Aerospike database
	Exists(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool

	// Count returns the count of a key in the Aerospike database
	Count(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64

	// FindOne retrieves a single row from the Aerospike database
	FindOne(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (NoSQLRow, error)

	// Find retrieves multiple rows from the Aerospike database
	Find(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (NoSQLRows, error)

	// Exec executes a function on the Aerospike database
	Exec(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (NoSQLRows, error)

	// Insert inserts data into the Aerospike database
	Insert(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error)

	// Delete deletes a key from the Aerospike database
	Delete(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64

	// DeleteMany deletes multiple keys from the Aerospike database
	DeleteMany(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64

	// Operate performs batch operations on the Aerospike database
	Operate(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error)

	// GetDb returns the Aerospike client instance
	GetDb() *aerospike.Client
}

type NoSQLRow interface {
	Scan(dest any) error
}

type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}
