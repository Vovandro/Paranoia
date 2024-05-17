package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

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

func (t *Postgres) Query(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Postgres) Exec(ctx context.Context, query interface{}, args ...interface{}) error {
	_, err := t.client.Exec(ctx, query.(string), args)

	return err
}

func (t *Postgres) GetDb() interface{} {
	return t.client
}
