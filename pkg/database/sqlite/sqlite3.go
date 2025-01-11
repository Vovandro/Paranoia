package sqlite

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type Sqlite3 struct {
	Name   string
	Config Config
	client *sql.DB

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	Database string `yaml:"database"`
}

func NewSqlite3(name string) *Sqlite3 {
	return &Sqlite3{
		Name: name,
	}
}

func (t *Sqlite3) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.Config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.Config.Database == "" {
		return errors.New("database file name is required")
	}

	t.client, err = sql.Open("sqlite3", t.Config.Database)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("sqlite." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("sqlite." + t.Name + ".time")

	return t.client.Ping()
}

func (t *Sqlite3) Stop() error {
	return t.client.Close()
}

func (t *Sqlite3) String() string {
	return t.Name
}

func (t *Sqlite3) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
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

func (t *Sqlite3) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
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

func (t *Sqlite3) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	_, err := t.client.ExecContext(ctx, query, args...)

	return err
}

func (t *Sqlite3) GetDb() interface{} {
	return t.client
}
