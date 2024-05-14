package database

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Sqlite3 struct {
	Name     string
	Database string
	app      interfaces.IService
	client   *sql.DB
}

func (t *Sqlite3) Init(app interfaces.IService) error {
	t.app = app
	var err error

	t.client, err = sql.Open("sqlite3", t.Database)

	if err != nil {
		return err
	}

	return nil
}

func (t *Sqlite3) Stop() error {
	if err := t.client.Close(); err != nil {
		return err
	}

	return nil
}

func (t *Sqlite3) String() string {
	return t.Name
}

func (t *Sqlite3) Exists(ctx context.Context, query interface{}, args ...interface{}) bool {
	find, err := t.client.Query(query.(string), args)

	if err != nil {
		return false
	}

	val, _ := find.Columns()

	return len(val) != 0
}

func (t *Sqlite3) Count(ctx context.Context, query interface{}, args ...interface{}) int64 {
	find, err := t.client.Query(query.(string), args)

	if err != nil {
		return 0
	}

	var values map[string]interface{}

	err = find.Scan(values)

	if err != nil {
		return 0
	}

	for _, v := range values {
		if v != nil {
			return v.(int64)
		}
	}

	return 0
}

func (t *Sqlite3) FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Sqlite3) Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Sqlite3) Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Sqlite3) Update(ctx context.Context, query interface{}, args ...interface{}) error {
	_, err := t.client.Exec(query.(string), args)

	return err
}

func (t *Sqlite3) Delete(ctx context.Context, query interface{}, args ...interface{}) int64 {
	tag, err := t.client.Exec(query.(string), args)

	if err != nil {
		return 0
	}

	v, _ := tag.RowsAffected()

	return v
}

func (t *Sqlite3) GetDb() interface{} {
	return t.client
}
