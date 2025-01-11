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
	Name   string
	Config Config
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
		Name: name,
	}
}

func (t *ClickHouse) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.Config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.Config.Hosts == "" {
		return errors.New("hosts is required")
	}

	if t.Config.Database == "" {
		return errors.New("database is required")
	}

	t.client, err = click.Open(&click.Options{
		Addr: strings.Split(t.Config.Hosts, ","),
		Auth: click.Auth{
			Database: t.Config.Database,
			Username: t.Config.Username,
			Password: t.Config.Password,
		},
	})

	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("clickhouse." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("clickhouse." + t.Name + ".time")

	return t.client.Ping(context.Background())
}

func (t *ClickHouse) Stop() error {
	return t.client.Close()
}

func (t *ClickHouse) String() string {
	return t.Name
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
