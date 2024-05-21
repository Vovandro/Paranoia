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

func (t *Sqlite3) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	find, err := t.client.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *Sqlite3) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	find := t.client.QueryRow(query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *Sqlite3) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.client.Query(query, args...)

	return err
}

func (t *Sqlite3) GetDb() interface{} {
	return t.client
}
