package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Postgres struct {
	name   string
	config Config
	pool   *pgxpool.Pool

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	URI string `yaml:"uri"`
}

func New(name string) *Postgres {
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

	poolCfg, err := pgxpool.ParseConfig(t.config.URI)
	if err != nil {
		return err
	}

	t.pool, err = pgxpool.NewWithConfig(context.TODO(), poolCfg)
	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("postgres." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("postgres." + t.name + ".time")

	return t.pool.Ping(context.TODO())
}

func (t *Postgres) Stop() error {
	t.pool.Close()
	return nil
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

	find, err := t.pool.Query(ctx, query, args...)

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

	find := t.pool.QueryRow(ctx, query, args...)

	if find == nil {
		return nil, errors.New(query + " not found")
	}

	return find, nil
}

func (t *Postgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	_, err := t.pool.Exec(ctx, query, args...)

	return err
}

func (t *Postgres) GetDb() *pgxpool.Pool {
	return t.pool
}

// BeginTx starts a transaction with metrics tracking
func (t *Postgres) BeginTx(ctx context.Context) (SQLTx, error) {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PGSQLTx{tx: tx, counter: t.counter, timeCounter: t.timeCounter}, nil
}

type PGSQLTx struct {
	tx          pgx.Tx
	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

func (p *PGSQLTx) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	rows, err := p.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PGSQLRows{Rows: rows}, nil
}

func (p *PGSQLTx) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	row := p.tx.QueryRow(ctx, query, args...)
	if row == nil {
		return nil, errors.New(query + " not found")
	}
	return row, nil
}

func (p *PGSQLTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)

	_, err := p.tx.Exec(ctx, query, args...)
	return err
}

func (p *PGSQLTx) Commit(ctx context.Context) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)
	return p.tx.Commit(ctx)
}

func (p *PGSQLTx) Rollback(ctx context.Context) error {
	defer func(s time.Time) {
		p.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	p.counter.Add(context.Background(), 1)
	return p.tx.Rollback(ctx)
}
