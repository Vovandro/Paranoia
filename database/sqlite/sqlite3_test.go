package sqlite

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSqlite3_Exec(t1 *testing.T) {
	db := initSQLite3Test("test_exec")
	defer closeSQLite3Test(db)

	type args struct {
		ctx   context.Context
		query string
		check string
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "base test query",
			args: args{
				ctx:   context.Background(),
				query: "insert into test values (?, ?, ?, datetime());",
				check: "select exists(select id from test where id = 10);",
				args: []interface{}{
					10,
					"test",
					1.0,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := db.Exec(tt.args.ctx, tt.args.query, tt.args.args...)

			if (err != nil) != tt.wantErr {
				t1.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var exist bool

			got, err := db.QueryRow(context.Background(), tt.args.check)

			err = got.Scan(&exist)
			if err != nil {
				t1.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != exist {
				t1.Errorf("QueryRow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite3_Query(t1 *testing.T) {
	db := initSQLite3Test("test_rows")
	defer closeSQLite3Test(db)

	name1 := "test"
	name2 := "test2"

	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []testSQLiteItem
		wantErr bool
	}{
		{
			name: "base test query",
			args: args{
				ctx:   context.Background(),
				query: "select * from test where id <= ?;",
				args: []interface{}{
					2,
				},
			},
			want: []testSQLiteItem{
				{
					1,
					&name1,
					1,
					time.Now(),
				},
				{
					2,
					&name2,
					0,
					time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			got, err := db.Query(tt.args.ctx, tt.args.query, tt.args.args...)

			if (err != nil) != tt.wantErr {
				t1.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			items := make([]testSQLiteItem, 0, 5)

			for got.Next() {
				var item testSQLiteItem

				err = got.Scan(&item.Id, &item.Name, &item.Balance, &item.CreatedAt)
				if err != nil {
					t1.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				items = append(items, item)
				tt.want[len(items)-1].CreatedAt = item.CreatedAt
			}

			if !reflect.DeepEqual(items, tt.want) {
				t1.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite3_QueryRow(t1 *testing.T) {
	db := initSQLite3Test("test_row")
	defer closeSQLite3Test(db)

	name1 := "test"

	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    testSQLiteItem
		wantErr bool
	}{
		{
			name: "base test query",
			args: args{
				ctx:   context.Background(),
				query: "select * from test where id = 1;",
				args:  []interface{}{},
			},
			want: testSQLiteItem{
				1,
				&name1,
				1,
				time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			got, err := db.QueryRow(tt.args.ctx, tt.args.query, tt.args.args...)

			if (err != nil) != tt.wantErr {
				t1.Errorf("QueryRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var item testSQLiteItem

			err = got.Scan(&item.Id, &item.Name, &item.Balance, &item.CreatedAt)
			if err != nil {
				t1.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want.CreatedAt = item.CreatedAt

			if !reflect.DeepEqual(item, tt.want) {
				t1.Errorf("QueryRow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite3_String(t1 *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"base test",
			fields{
				"test",
			},
			"test",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Sqlite3{
				Name: tt.fields.Name,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testSQLiteItem struct {
	Id        int64
	Name      *string
	Balance   float64
	CreatedAt time.Time
}

func initSQLite3Test(name string) *Sqlite3 {
	db := NewSqlite3(name, Sqlite3Config{
		Database: name + ".db",
	})

	err := db.Init(nil)

	if err != nil {
		panic(err)
	}

	_, err = db.client.Exec("create table test (id integer primary key, name varchar(255), balance float not null, created_at datetime)")

	if err != nil {
		os.Remove(db.Name + ".db")
		panic(err)
	}

	_, err = db.client.Exec(`insert into test (id, name, balance, created_at) values 
						 (1, 'test', 1.0, datetime()), (2, 'test2', 0.0, datetime()), (3, null, 50.0, datetime());`)

	if err != nil {
		os.Remove(db.Name + ".db")
		panic(err)
	}

	return db
}

func closeSQLite3Test(db *Sqlite3) {
	db.Stop()

	os.Remove(db.Name + ".db")
}
