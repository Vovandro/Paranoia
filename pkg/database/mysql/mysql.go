package mysql

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type MySQL struct {
	name   string
	config Config
	client *sql.DB

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	URI string `yaml:"uri"`
}

func New(name string) *MySQL {
	return &MySQL{
		name: name,
	}
}

func (t *MySQL) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.URI == "" {
		return errors.New("URI is required")
	}

	t.client, err = sql.Open("mysql", t.config.URI)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("mysql." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("mysql." + t.name + ".time")

	return t.client.Ping()
}

func (t *MySQL) Stop() error {
	return t.client.Close()
}

func (t *MySQL) Name() string {
	return t.name
}

func (t *MySQL) Type() string {
	return "database"
}

func (t *MySQL) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find, err := t.client.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *MySQL) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find := t.client.QueryRowContext(ctx, query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *MySQL) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	_, err := t.client.ExecContext(ctx, query, args...)

	return err
}

func (t *MySQL) GetDb() interface{} {
	return t.client
}
