package interfaces

import (
	"context"
	"github.com/aerospike/aerospike-client-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type INoSql interface {
	Init(app IEngine) error
	Stop() error
	String() string

	GetDb() interface{}
}

type NoSQLRow interface {
	Scan(dest any) error
}

type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}

type IAerospike interface {
	// CRUD operations
	Exists(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool
	Count(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64
	FindOne(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (NoSQLRow, error)
	Find(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (NoSQLRows, error)
	Exec(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (NoSQLRows, error)
	Insert(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error)
	Delete(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64
	DeleteMany(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64

	// Advanced operations
	Operate(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error)

	// Access underlying database client
	GetDb() *aerospike.Client
}

type IMongoDB interface {
	// CRUD operations
	Exists(ctx context.Context, collection string, query bson.D) bool
	Count(ctx context.Context, collection string, query bson.D, opt *options.CountOptions) int64
	FindOne(ctx context.Context, collection string, query bson.D, opt *options.FindOneOptions) (NoSQLRow, error)
	Find(ctx context.Context, collection string, query bson.D, opt *options.FindOptions) (NoSQLRows, error)
	Exec(ctx context.Context, collection string, query bson.D, opt *options.AggregateOptions) (NoSQLRows, error)
	Insert(ctx context.Context, collection string, query bson.D, opt *options.InsertOneOptions) (interface{}, error)
	Update(ctx context.Context, collection string, query bson.D, update bson.D, opt *options.UpdateOptions) error
	Delete(ctx context.Context, collection string, query bson.D, opt *options.DeleteOptions) int64
	Batch(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error)

	// Access underlying database
	GetDb() *mongo.Database
}
