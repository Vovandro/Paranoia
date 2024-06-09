package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Postgres struct {
	Name   string
	Config PostgresConfig
	app    interfaces.IService
	client *pgx.Conn
}

type PostgresConfig struct {
	URI string
}

func (t *Postgres) Init(app interfaces.IService) error {
	t.app = app
	var err error

	t.client, err = pgx.Connect(context.TODO(), t.Config.URI)

	if err != nil {
		return err
	}

	return t.client.Ping(context.TODO())
}

func (t *Postgres) Stop() error {
	return t.client.Close(context.TODO())
}

func (t *Postgres) String() string {
	return t.Name
}

func (t *Postgres) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	find, err := t.client.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return &PGSQLRows{find}, err
}

func (t *Postgres) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	find := t.client.QueryRow(ctx, query, args...)

	if find == nil {
		return nil, fmt.Errorf(query + " not found")
	}

	return find, nil
}

func (t *Postgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.client.Exec(ctx, query, args...)

	return err
}

func (t *Postgres) GetDb() interface{} {
	return t.client
}
