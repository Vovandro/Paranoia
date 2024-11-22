package noSql

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoDB struct {
	Name string
}

func (t *MockMongoDB) Exists(ctx context.Context, collection string, query interface{}) bool {
	return false
}

func (t *MockMongoDB) Count(ctx context.Context, collection string, query interface{}, opt *options.CountOptions) int64 {
	return 0
}

func (t *MockMongoDB) FindOne(ctx context.Context, collection string, query interface{}, opt *options.FindOneOptions) (interfaces.NoSQLRow, error) {
	return &MongoRow{nil}, nil
}

func (t *MockMongoDB) FindOneAndUpdate(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) (interfaces.NoSQLRow, error) {
	return &MongoRow{nil}, nil
}

func (t *MockMongoDB) Find(ctx context.Context, collection string, query interface{}, opt *options.FindOptions) (interfaces.NoSQLRows, error) {
	return &MongoRows{nil}, nil
}

func (t *MockMongoDB) Exec(ctx context.Context, collection string, query interface{}, opt *options.AggregateOptions) (interfaces.NoSQLRows, error) {
	return &MongoRows{nil}, nil
}

func (t *MockMongoDB) Insert(ctx context.Context, collection string, query interface{}, opt *options.InsertOneOptions) (interface{}, error) {
	return 0, nil
}

func (t *MockMongoDB) Update(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.UpdateOptions) error {
	return nil

}

func (t *MockMongoDB) Delete(ctx context.Context, collection string, query interface{}, opt *options.DeleteOptions) int64 {
	return 0
}

func (t *MockMongoDB) Batch(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error) {
	return 0, nil
}

func (t *MockMongoDB) GetDb() *mongo.Database {
	return nil
}
