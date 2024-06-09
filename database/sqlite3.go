package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Sqlite3 struct {
	Name   string
	Config Sqlite3Config
	app    interfaces.IService
	client *sql.DB
}

type Sqlite3Config struct {
	Database string `yaml:"database"`
}

func NewSqlite3(name string, cfg Sqlite3Config) *Sqlite3 {
	return &Sqlite3{
		Name:   name,
		Config: cfg,
	}
}

func (t *Sqlite3) Init(app interfaces.IService) error {
	t.app = app
	var err error

	if t.Config.Database == "" {
		return fmt.Errorf("database file name is required")
	}

	t.client, err = sql.Open("sqlite3", t.Config.Database)

	if err != nil {
		return err
	}

	return t.client.Ping()
}

func (t *Sqlite3) Stop() error {
	return t.client.Close()
}

func (t *Sqlite3) String() string {
	return t.Name
}

func (t *Sqlite3) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	find, err := t.client.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *Sqlite3) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	find := t.client.QueryRowContext(ctx, query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *Sqlite3) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.client.ExecContext(ctx, query, args...)

	return err
}

func (t *Sqlite3) GetDb() interface{} {
	return t.client
}
