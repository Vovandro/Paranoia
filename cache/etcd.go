package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strconv"
	"strings"
	"time"
)

type Etcd struct {
	Name   string
	Config EtcdConfig

	app    interfaces.IEngine
	client *clientv3.Client

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type EtcdConfig struct {
	Hosts    string `yaml:"hosts"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func NewEtcd(name string, cfg EtcdConfig) *Etcd {
	return &Etcd{
		Name:   name,
		Config: cfg,
	}
}

func (t *Etcd) Init(app interfaces.IEngine) error {
	t.app = app
	var err error

	t.client, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(t.Config.Hosts, ","),
		Username:    t.Config.Username,
		Password:    t.Config.Password,
		DialTimeout: time.Second * 3,
	})

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("redis." + t.Name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("redis." + t.Name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("redis." + t.Name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("redis." + t.Name + ".timeWrite")

	return nil
}

func (t *Etcd) Stop() error {
	return t.client.Close()
}

func (t *Etcd) String() string {
	return t.Name
}

func (t *Etcd) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	res, err := t.client.Get(ctx, key)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return false
	}

	if len(res.Kvs) == 0 {
		return false
	}

	return true
}

func (t *Etcd) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	if _, ok := args.(string); !ok {
		return ErrTypeMismatch
	}

	s := time.Now()
	t.counterWrite.Add(ctx, 1)
	lease, _ := t.client.Grant(ctx, int64(timeout.Seconds()))

	_, err := t.client.Put(ctx, key, args.(string), clientv3.WithLease(lease.ID))
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	return err
}

func (t *Etcd) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	s := time.Now()
	data, err := t.GetMap(ctx, key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	data.(map[string]any)[key2] = args

	err = t.SetMap(ctx, key, data, timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	return err
}

func (t *Etcd) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()

	data, err := json.Marshal(args)

	if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	err = t.Set(ctx, key, string(data), timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Etcd) Get(ctx context.Context, key string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(ctx, key)

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if err != nil {
		return "", err
	}

	if len(item.Kvs) == 0 {
		return "", ErrKeyNotFound
	}

	return string(item.Kvs[0].Value), nil
}

func (t *Etcd) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	data, err := t.GetMap(ctx, key)

	if err != nil {
		return "", err
	}

	if val, ok := data.(map[string]any)[key2]; ok {
		return val, nil
	}

	return "", ErrKeyNotFound
}

func (t *Etcd) GetMap(ctx context.Context, key string) (any, error) {
	s := time.Now()

	data, err := t.Get(ctx, key)

	if err != nil {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return "", err
	}

	res := make(map[string]any)
	err = json.Unmarshal([]byte(data.(string)), &res)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return "", ErrTypeMismatch
	}

	return res, nil
}

func (t *Etcd) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()

	v, _ := t.Get(ctx, key)

	conv, _ := strconv.ParseInt(v.(string), 10, 64)
	val += conv

	err := t.Set(ctx, key, fmt.Sprint(val), timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return val, err
}

func (t *Etcd) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()

	data, err := t.GetMap(ctx, key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return 0, err
	}

	if v, ok := data.(map[string]any)[key2]; ok {
		data.(map[string]any)[key2] = int64(v.(float64)) + val
	} else {
		data.(map[string]any)[key2] = val
	}
	err = t.SetMap(ctx, key, data, timeout)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return data.(map[string]any)[key2].(int64), err
}

func (t *Etcd) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return t.Increment(ctx, key, val*-1, timeout)
}

func (t *Etcd) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return t.IncrementIn(ctx, key, key2, -1*val, timeout)
}

func (t *Etcd) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	_, err := t.client.Delete(ctx, key)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return ErrKeyNotFound
	}

	return err
}

func (t *Etcd) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	resp, err := t.client.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return ErrKeyNotFound
	}

	_, err = t.client.KeepAliveOnce(ctx, clientv3.LeaseID(resp.Kvs[0].Lease))
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	return err
}
