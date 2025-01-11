package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type Postgres struct {
	name   string
	config Config
	client *pgx.Conn

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	URI string `yaml:"uri"`
}

func NewPostgres(name string) *Postgres {
	return &Postgres{
		name: name,
	}
}

func (t *Postgres) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.URI == "" {
		return errors.New("URI is required")
	}

	t.client, err = pgx.Connect(context.TODO(), t.config.URI)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("postgres." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("postgres." + t.name + ".time")

	return t.client.Ping(context.TODO())
}

func (t *Postgres) Stop() error {
	return t.client.Close(context.TODO())
}

func (t *Postgres) Name() string {
	return t.name
}

func (t *Postgres) Type() string {
	return "database"
}

func (t *Postgres) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
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

func (t *Postgres) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
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
