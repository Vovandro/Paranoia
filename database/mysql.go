package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type MySQL struct {
	Name   string
	Config MySQLConfig
	app    interfaces.IService
	client *sql.DB
}

type MySQLConfig struct {
	URI string `yaml:"uri"`
}

func NewMySQL(name string, cfg MySQLConfig) *MySQL {
	return &MySQL{
		Name:   name,
		Config: cfg,
	}
}

func (t *MySQL) Init(app interfaces.IService) error {
	t.app = app
	var err error

	if t.Config.URI == "" {
		return fmt.Errorf("URI is required")
	}

	t.client, err = sql.Open("mysql", t.Config.URI)

	if err != nil {
		return err
	}

	return t.client.Ping()
}

func (t *MySQL) Stop() error {
	return t.client.Close()
}

func (t *MySQL) String() string {
	return t.Name
}

func (t *MySQL) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	find, err := t.client.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *MySQL) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	find := t.client.QueryRowContext(ctx, query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *MySQL) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.client.ExecContext(ctx, query, args...)

	return err
}

func (t *MySQL) GetDb() interface{} {
	return t.client
}
