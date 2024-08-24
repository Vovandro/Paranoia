package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type MySQL struct {
	Name   string
	Config MySQLConfig
	app    interfaces.IEngine
	client *sql.DB

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
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

func (t *MySQL) Init(app interfaces.IEngine) error {
	t.app = app
	var err error

	if t.Config.URI == "" {
		return fmt.Errorf("URI is required")
	}

	t.client, err = sql.Open("mysql", t.Config.URI)

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("mysql." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("mysql." + t.Name + ".time")

	return t.client.Ping()
}

func (t *MySQL) Stop() error {
	return t.client.Close()
}

func (t *MySQL) String() string {
	return t.Name
}

func (t *MySQL) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
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

func (t *MySQL) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
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
