package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type Postgres struct {
	Name   string
	Config PostgresConfig
	app    interfaces.IEngine
	client *pgx.Conn

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type PostgresConfig struct {
	URI string `yaml:"uri"`
}

func NewPostgres(name string, cfg PostgresConfig) *Postgres {
	return &Postgres{
		Name:   name,
		Config: cfg,
	}
}

func (t *Postgres) Init(app interfaces.IEngine) error {
	t.app = app
	var err error

	t.client, err = pgx.Connect(context.TODO(), t.Config.URI)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("postgres." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("postgres." + t.Name + ".time")

	return t.client.Ping(context.TODO())
}

func (t *Postgres) Stop() error {
	return t.client.Close(context.TODO())
}

func (t *Postgres) String() string {
	return t.Name
}

func (t *Postgres) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find, err := t.client.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return &PGSQLRows{find}, err
}

func (t *Postgres) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find := t.client.QueryRow(ctx, query, args...)

	if find == nil {
		return nil, fmt.Errorf(query + " not found")
	}

	return find, nil
}

func (t *Postgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	_, err := t.client.Exec(ctx, query, args...)

	return err
}

func (t *Postgres) GetDb() interface{} {
	return t.client
}
