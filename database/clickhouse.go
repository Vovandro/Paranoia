package database

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"strings"
)

type ClickHouse struct {
	Name   string
	Config ClickHouseConfig
	app    interfaces.IService
	client driver.Conn
}

type ClickHouseConfig struct {
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hosts    string `yaml:"hosts"`
}

func NewClickHouse(name string, cfg ClickHouseConfig) *ClickHouse {
	return &ClickHouse{
		Name:   name,
		Config: cfg,
	}
}

func (t *ClickHouse) Init(app interfaces.IService) error {
	t.app = app
	var err error

	t.client, err = clickhouse.Open(&clickhouse.Options{
		Addr: strings.Split(t.Config.Hosts, ","),
		Auth: clickhouse.Auth{
			Database: t.Config.Database,
			Username: t.Config.Username,
			Password: t.Config.Password,
		},
	})

	if err != nil {
		return err
	}

	return t.client.Ping(context.Background())
}

func (t *ClickHouse) Stop() error {
	return t.client.Close()
}

func (t *ClickHouse) String() string {
	return t.Name
}

func (t *ClickHouse) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	find, err := t.client.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return find, nil
}

func (t *ClickHouse) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	find := t.client.QueryRow(ctx, query, args...)

	if find.Err() != nil {
		return nil, find.Err()
	}

	return find, nil
}

func (t *ClickHouse) Exec(ctx context.Context, query string, args ...interface{}) error {
	err := t.client.Exec(ctx, query, args...)

	return err
}

func (t *ClickHouse) GetDb() interface{} {
	return t.client
}
