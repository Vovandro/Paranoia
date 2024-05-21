package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type PGSQLRows struct {
	Rows pgx.Rows
}

func (t *PGSQLRows) Next() bool {
	return t.Rows.Next()
}

func (t *PGSQLRows) Scan(dest ...any) error {
	return t.Rows.Scan(dest...)
}

func (t *PGSQLRows) Close() error {
	t.Rows.Close()
	return nil
}

type Postgres struct {
	Name     string
	Database string
	URI      string
	app      interfaces.IService
	client   *pgx.Conn
}

func (t *Postgres) Init(app interfaces.IService) error {
	t.app = app
	var err error

	t.client, err = pgx.Connect(context.TODO(), t.URI)

	if err != nil {
		return err
	}

	return nil
}

func (t *Postgres) Stop() error {
	if err := t.client.Close(context.TODO()); err != nil {
		return err
	}

	return nil
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
