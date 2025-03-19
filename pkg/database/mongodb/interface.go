package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IMongoDB defines the interface for MongoDB operations
type IMongoDB interface {
	// Exists checks if a document exists in the MongoDB collection
	Exists(ctx context.Context, collection string, query interface{}) bool

	// Count returns the count of documents in the MongoDB collection
	Count(ctx context.Context, collection string, query interface{}, opt *options.CountOptions) int64

	// FindOne retrieves a single document from the MongoDB collection
	FindOne(ctx context.Context, collection string, query interface{}, opt *options.FindOneOptions) (NoSQLRow, error)

	// FindOneAndUpdate retrieves and updates a single document in the MongoDB collection
	FindOneAndUpdate(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) (NoSQLRow, error)

	// Find retrieves multiple documents from the MongoDB collection
	Find(ctx context.Context, collection string, query interface{}, opt *options.FindOptions) (NoSQLRows, error)

	// Exec executes an aggregation pipeline on the MongoDB collection
	Exec(ctx context.Context, collection string, query interface{}, opt *options.AggregateOptions) (NoSQLRows, error)

	// Insert inserts a document into the MongoDB collection
	Insert(ctx context.Context, collection string, query interface{}, opt *options.InsertOneOptions) (interface{}, error)

	// Update updates multiple documents in the MongoDB collection
	Update(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.UpdateOptions) error

	// Delete deletes multiple documents from the MongoDB collection
	Delete(ctx context.Context, collection string, query interface{}, opt *options.DeleteOptions) int64

	// Batch performs batch operations on the MongoDB collection
	Batch(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error)

	// GetDb returns the MongoDB client instance
	GetDb() *mongo.Database
}

type NoSQLRow interface {
	Scan(dest any) error
}

type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}
