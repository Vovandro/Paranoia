package clickhouse

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

type testSQLItem struct {
	Id        int64
	Name      *string
	Balance   float64
	CreatedAt time.Time
}

func TestClickhouse_Exec(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initClickhouseTest("test_exec")
	defer closeClickhouseTest(db)

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
				query: "insert into test.test_exec values (?, ?, ?, now());",
				check: "select id from test.test_exec where id = 10;",
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

			var exist *int64

			got, err := db.QueryRow(context.Background(), tt.args.check)
			if (err != nil) != tt.wantErr {
				t1.Errorf("QueryRow() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = got.Scan(&exist)
			if err != nil {
				t1.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != (exist != nil) {
				t1.Errorf("QueryRow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClickhouse_Query(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initClickhouseTest("test_rows")
	defer closeClickhouseTest(db)

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
		want    []testSQLItem
		wantErr bool
	}{
		{
			name: "base test query",
			args: args{
				ctx:   context.Background(),
				query: "select * from test.test_rows where id <= ?;",
				args: []interface{}{
					2,
				},
			},
			want: []testSQLItem{
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

			items := make([]testSQLItem, 0, 5)

			for got.Next() {
				var item testSQLItem

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

func TestClickhouse_QueryRow(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	db := initClickhouseTest("test_row")
	defer closeClickhouseTest(db)

	name1 := "test"

	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    testSQLItem
		wantErr bool
	}{
		{
			name: "base test query",
			args: args{
				ctx:   context.Background(),
				query: "select * from test.test_row where id = 1;",
				args:  []interface{}{},
			},
			want: testSQLItem{
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

			var item testSQLItem

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

func TestClickhouse_String(t1 *testing.T) {
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
			t := &ClickHouse{
				name: tt.fields.Name,
			}
			if got := t.Name(); got != tt.want {
				t1.Errorf("name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func initClickhouseTest(name string) *ClickHouse {
	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	db := New(name)

	err := db.Init(map[string]interface{}{
		"database": "test",
		"username": "test",
		"password": "test",
		"hosts":    host + ":8124",
	})

	if err != nil {
		panic(err)
	}

	err = db.client.Exec(context.Background(), `create table test.`+name+`
(
    id         Int64,
    name       Nullable(varchar(255)),
    balance    Nullable(Float64),
    created_at timestamp
)
    engine = Memory;`)

	if err != nil {
		panic(err)
	}

	err = db.client.Exec(context.Background(), `insert into test.`+name+` (id, name, balance, created_at) values 
						 (1, 'test', 1.0, now()), (2, 'test2', 0.0, now()), (3, null, 50.0, now());`)

	if err != nil {
		closeClickhouseTest(db)
		panic(err)
	}

	return db
}

func closeClickhouseTest(db *ClickHouse) {
	db.Exec(context.Background(), "drop table if exists test."+db.name+";")
	db.Stop()
}
