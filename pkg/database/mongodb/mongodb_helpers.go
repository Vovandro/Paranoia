package mongodb

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRow struct {
	row *mongo.SingleResult
}

type MongoRows struct {
	rows *mongo.Cursor
}

func (t *MongoRow) Scan(dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	// Check if dest is a pointer
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}
	return t.row.Decode(dest)
}

func (t *MongoRows) Next() bool {
	return t.rows.Next(context.TODO())
}

func (t *MongoRows) Scan(dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	// Check if dest is a pointer
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}

	return t.rows.Decode(dest)
}

func (t *MongoRows) Close() error {
	return t.rows.Close(context.TODO())
}
