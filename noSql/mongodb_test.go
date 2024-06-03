package noSql

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMongoDB_Count(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			if got := t.Count(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...); got != tt.want {
				t1.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Delete(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			if got := t.Delete(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...); got != tt.want {
				t1.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Exec(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interfaces.NoSQLRows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			got, err := t.Exec(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Exists(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			if got := t.Exists(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...); got != tt.want {
				t1.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Find(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interfaces.NoSQLRows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			got, err := t.Find(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_FindOne(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("find_one_test")
	defer closeMongoTest(db)
	name := "test2"

	type args struct {
		key         interface{}
		queryInsert interface{}
		queryFind   interface{}
		args        []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    itemMongo
		wantErr bool
	}{
		{
			"base test",
			args{
				"find_one_test",
				nil,
				bson.M{"balance": 0},
				nil,
			},
			itemMongo{
				Id:        primitive.ObjectID{},
				Name:      &name,
				Balance:   0,
				CreatedAt: time.Time{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.queryInsert != nil {
				db.Insert(context.Background(), tt.args.key, tt.args.queryInsert)
			}

			got, err := db.FindOne(context.Background(), tt.args.key, tt.args.queryFind, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var item itemMongo

			got.Scan(&item)

			if *item.Name != *tt.want.Name || item.Balance != tt.want.Balance {
				t1.Errorf("FindOne() got = %v, want %v", item, tt.want)
			}
		})
	}
}

func TestMongoDB_Insert(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type args struct {
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"base test",
			args{
				"insert_test",
				map[string]interface{}{
					"foo": "bar",
				},
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := db.Insert(context.TODO(), tt.args.key, tt.args.query, tt.args.args...)

			if (err != nil) != tt.wantErr {
				t1.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t1.Errorf("Insert() want objectId")
			}
		})
	}
}

func TestMongoDB_String(t1 *testing.T) {

	t1.Run("test name", func(t1 *testing.T) {
		t := &MongoDB{
			Name: "test",
		}
		if got := t.String(); got != "test" {
			t1.Errorf("String() = %v, want %v", got, "test")
		}
	})
}

func TestMongoDB_Update(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx   context.Context
		key   interface{}
		query interface{}
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			if err := t.Update(tt.args.ctx, tt.args.key, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t1.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDB_Batch(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("insert_test")
	defer closeMongoTest(db)

	type fields struct {
		Name     string
		Database string
		Options  *options.ClientOptions
		app      interfaces.IService
		client   *mongo.Client
		db       *mongo.Database
	}
	type args struct {
		ctx    context.Context
		key    interface{}
		query  interface{}
		typeOp string
		args   []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MongoDB{
				Name:     tt.fields.Name,
				Database: tt.fields.Database,
				Options:  tt.fields.Options,
				app:      tt.fields.app,
				client:   tt.fields.client,
				db:       tt.fields.db,
			}
			got, err := t.Batch(tt.args.ctx, tt.args.key, tt.args.query, tt.args.typeOp, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("Batch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type itemMongo struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      *string            `bson:"name"`
	Balance   float64            `bson:"balance"`
	CreatedAt time.Time          `bson:"created_at"`
}

func initMongoTest(name string) *MongoDB {
	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	db := &MongoDB{
		Name:     name,
		Database: "tests",
		Options: &options.ClientOptions{
			Hosts: strings.Split(host+":27017", ","),
			Auth: &options.Credential{
				Username:   "test",
				Password:   "test",
				AuthSource: "tests",
			},
		},
	}

	err := db.Init(nil)

	if err != nil {
		panic(err)
	}

	col := db.db.Collection(name)

	_, err = col.InsertMany(context.Background(), []interface{}{
		map[string]interface{}{
			"name":       "test",
			"balance":    1.0,
			"created_at": time.Now(),
		},
		map[string]interface{}{
			"name":       "test2",
			"balance":    0.0,
			"created_at": time.Now(),
		},
		map[string]interface{}{
			"balance":    50.0,
			"created_at": time.Now(),
		},
	})

	if err != nil {
		closeMongoTest(db)
		panic(err)
	}

	return db
}

func closeMongoTest(db *MongoDB) {
	db.db.Collection(db.Name).Drop(context.Background())
	db.Stop()
}
