package noSql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
	"time"
)

func TestMongoDB_Exists(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("exist_test")
	defer closeMongoTest(db)

	type args struct {
		collection  string
		queryInsert bson.M
		queryFind   bson.M
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"base test",
			args{
				"exist_test",
				nil,
				bson.M{"balance": 0},
			},
			true,
			false,
		},
		{
			"test not exists",
			args{
				"exist_test",
				nil,
				bson.M{"balance": 1000},
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.queryInsert != nil {
				db.Insert(context.Background(), tt.args.collection, tt.args.queryInsert, nil)
			}

			got := db.Exists(context.Background(), tt.args.collection, tt.args.queryFind)

			if got != tt.want {
				t1.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Find(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initMongoTest("find_test")
	defer closeMongoTest(db)
	name := "test"
	name2 := "test2"

	type args struct {
		collection  string
		queryInsert interface{}
		queryFind   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []itemMongo
		wantErr bool
	}{
		{
			"base test",
			args{
				"find_test",
				nil,
				bson.M{"balance": bson.M{"$lte": 2}},
			},
			[]itemMongo{
				{
					Id:        primitive.ObjectID{},
					Name:      &name,
					Balance:   1,
					CreatedAt: time.Time{},
				},
				{
					Id:        primitive.ObjectID{},
					Name:      &name2,
					Balance:   0,
					CreatedAt: time.Time{},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.queryInsert != nil {
				db.Insert(context.Background(), tt.args.collection, tt.args.queryInsert, nil)
			}

			got, err := db.Find(context.Background(), tt.args.collection, tt.args.queryFind, nil)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var items []itemMongo

			for got.Next() {
				var item itemMongo
				got.Scan(&item)
				items = append(items, item)
			}

			if len(items) != len(tt.want) {
				t1.Errorf("Find() got = %v, want %v", items, tt.want)
			}

			for i, item := range items {
				if *item.Name != *tt.want[i].Name || item.Balance != tt.want[i].Balance {
					t1.Errorf("Find() got = %v, want %v", item, tt.want)
				}
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
		collection  string
		queryInsert interface{}
		queryFind   interface{}
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
				db.Insert(context.Background(), tt.args.collection, tt.args.queryInsert, nil)
			}

			got, err := db.FindOne(context.Background(), tt.args.collection, tt.args.queryFind, nil)
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
		collection string
		query      interface{}
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
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := db.Insert(context.TODO(), tt.args.collection, tt.args.query, nil)

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

	db := initMongoTest("update_test")
	defer closeMongoTest(db)

	t1.Run("base test", func(t1 *testing.T) {

		err := db.Update(context.Background(), "update_test", bson.M{"balance": 50}, bson.M{"$set": bson.M{"balance": 100}}, nil)

		if err != nil {
			t1.Errorf("Update() error = %v", err)
			return
		}

		if !db.Exists(context.Background(), "update_test", bson.M{"balance": 100}) {
			t1.Errorf("Update() want = %v", true)
		}
	})
}

type itemMongo struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      *string            `bson:"name"`
	Balance   float64            `bson:"balance"`
	CreatedAt time.Time          `bson:"created_at"`
}

func initMongoTest(name string) *MongoDB {
	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	db := NewMongoDB(name, MongoDBConfig{
		Database: "tests",
		User:     "test",
		Password: "test",
		Hosts:    host + ":27017",
	})

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
