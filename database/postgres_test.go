package database

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPostgres_Exec(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initPostgresTest("test_exec")
	defer closePostgresTest(db)

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
				query: "insert into test_exec values ($1, $2, $3, now());",
				check: "select exists(select id from test_exec where id = 10);",
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

func TestPostgres_Query(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initPostgresTest("test_rows")
	defer closePostgresTest(db)

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
				query: "select * from test_rows where id <= $1;",
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

func TestPostgres_QueryRow(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initPostgresTest("test_row")
	defer closePostgresTest(db)

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
				query: "select * from test_row where id = 1;",
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

func TestPostgres_String(t1 *testing.T) {
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
			t := &Postgres{
				Name: tt.fields.Name,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func initPostgresTest(name string) *Postgres {
	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	db := &Postgres{
		Name: name,
		URI:  "postgres://test:test@" + host + ":5432/test",
	}

	err := db.Init(nil)

	if err != nil {
		panic(err)
	}

	_, err = db.client.Exec(context.Background(), "create table if not exists "+name+" (id integer primary key, name varchar(255), balance float not null, created_at timestamp)")

	if err != nil {
		closePostgresTest(db)
		panic(err)
	}

	_, err = db.client.Exec(context.Background(), `insert into `+name+` (id, name, balance, created_at) values 
						 (1, 'test', 1.0, now()), (2, 'test2', 0.0, now()), (3, null, 50.0, now());`)

	if err != nil {
		closePostgresTest(db)
		panic(err)
	}

	return db
}

func closePostgresTest(db *Postgres) {
	db.Exec(context.Background(), "drop table if exists "+db.Name+";")
	db.Stop()
}
