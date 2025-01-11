package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRow struct {
	row *mongo.SingleResult
}

type MongoRows struct {
	rows *mongo.Cursor
}

func (t *MongoRow) Scan(dest any) error {
	return t.row.Decode(dest)
}

func (t *MongoRows) Next() bool {
	return t.rows.Next(context.TODO())
}

func (t *MongoRows) Scan(dest any) error {
	return t.rows.Decode(dest)
}

func (t *MongoRows) Close() error {
	return t.rows.Close(context.TODO())
}
