package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
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

// BeginTx starts a transaction with metrics tracking
func (t *MySQL) BeginTx(ctx context.Context) (SQLTx, error) {
	tx, err := t.client.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &MySQLTx{tx: tx, counter: t.counter, timeCounter: t.timeCounter}, nil
}

type MySQLTx struct {
	tx          *sql.Tx
	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

func (p *MySQLTx) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	rows, err := p.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *MySQLTx) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	row := p.tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}

func (p *MySQLTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	_, err := p.tx.ExecContext(ctx, query, args...)
	return err
}

func (p *MySQLTx) Commit(ctx context.Context) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)
	return p.tx.Commit()
}

func (p *MySQLTx) Rollback(ctx context.Context) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)
	return p.tx.Rollback()
}
