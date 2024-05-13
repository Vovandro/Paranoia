package database

import (
	"Paranoia/interfaces"
	"context"
	"github.com/jackc/pgx/v5"
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

func (t *Postgres) Exists(ctx context.Context, query interface{}, args ...interface{}) bool {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return false
	}

	return len(find.RawValues()) != 0
}

func (t *Postgres) Count(ctx context.Context, query interface{}, args ...interface{}) int64 {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return 0
	}

	values, err := find.Values()

	if err != nil {
		return 0
	}

	return values[0].(int64)
}

func (t *Postgres) FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Postgres) Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Postgres) Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	find, err := t.client.Query(ctx, query.(string), args)

	if err != nil {
		return err
	}

	err = find.Scan(model)

	return err
}

func (t *Postgres) Update(ctx context.Context, query interface{}, args ...interface{}) error {
	_, err := t.client.Exec(ctx, query.(string), args)

	return err
}

func (t *Postgres) Delete(ctx context.Context, query interface{}, args ...interface{}) int64 {
	tag, err := t.client.Exec(ctx, query.(string), args)

	if err != nil {
		return 0
	}

	return tag.RowsAffected()
}

func (t *Postgres) GetDb() interface{} {
	return t.client
}
