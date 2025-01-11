package clickhouse

import (
	"context"
	"errors"
	click "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type ClickHouse struct {
	name   string
	config Config
	client driver.Conn

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hosts    string `yaml:"hosts"`
}

func NewClickHouse(name string) *ClickHouse {
	return &ClickHouse{
		name: name,
	}
}

func (t *ClickHouse) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	if t.config.Database == "" {
		return errors.New("database is required")
	}

	t.client, err = click.Open(&click.Options{
		Addr: strings.Split(t.config.Hosts, ","),
		Auth: click.Auth{
			Database: t.config.Database,
			Username: t.config.Username,
			Password: t.config.Password,
		},
	})

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("clickhouse." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("clickhouse." + t.name + ".time")

	return t.client.Ping(context.Background())
}

func (t *ClickHouse) Stop() error {
	return t.client.Close()
}

func (t *ClickHouse) Name() string {
	return t.name
}

func (t *ClickHouse) Type() string {
	return "database"
}

func (t *ClickHouse) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find, err := t.client.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *ClickHouse) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find := t.client.QueryRow(ctx, query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *ClickHouse) Exec(ctx context.Context, query string, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	err := t.client.Exec(ctx, query, args...)

	return err
}

func (t *ClickHouse) GetDb() interface{} {
	return t.client
}
