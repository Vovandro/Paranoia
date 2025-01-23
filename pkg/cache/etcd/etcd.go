package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/devpro_studio/go_utils/decode"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strconv"
	"strings"
	"time"
)

type Etcd struct {
	name   string
	config Config

	client *clientv3.Client

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type Config struct {
	Hosts     string `yaml:"hosts"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	KeyPrefix string `yaml:"key_prefix"`
}

func New(name string) *Etcd {
	return &Etcd{
		name: name,
	}
}

func (t *Etcd) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	t.client, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(t.config.Hosts, ","),
		Username:    t.config.Username,
		Password:    t.config.Password,
		DialTimeout: time.Second * 3,
	})

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("redis." + t.name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("redis." + t.name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("redis." + t.name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("redis." + t.name + ".timeWrite")

	return nil
}

func (t *Etcd) Stop() error {
	return t.client.Close()
}

func (t *Etcd) Name() string {
	return t.name
}

func (t *Etcd) Type() string {
	return "cache"
}

func (t *Etcd) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	res, err := t.client.Get(ctx, t.config.KeyPrefix+key)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return false
	}

	if len(res.Kvs) == 0 {
		return false
	}

	return true
}

func (t *Etcd) Set(ctx context.Context, key string, args string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)
	lease, _ := t.client.Grant(ctx, int64(timeout.Seconds()))

	_, err := t.client.Put(ctx, t.config.KeyPrefix+key, args, clientv3.WithLease(lease.ID))
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	return err
}

func (t *Etcd) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	s := time.Now()
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	data.(map[string]any)[key2] = args

	err = t.SetMap(ctx, t.config.KeyPrefix+key, data, timeout)
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

	err = t.Set(ctx, t.config.KeyPrefix+key, string(data), timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Etcd) Get(ctx context.Context, key string) ([]byte, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(ctx, t.config.KeyPrefix+key)

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if err != nil {
		return []byte(""), err
	}

	if len(item.Kvs) == 0 {
		return []byte(""), ErrKeyNotFound
	}

	return item.Kvs[0].Value, nil
}

func (t *Etcd) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

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

	data, err := t.Get(ctx, t.config.KeyPrefix+key)

	if err != nil {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return "", err
	}

	res := make(map[string]any)
	err = json.Unmarshal(data, &res)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return "", ErrTypeMismatch
	}

	return res, nil
}

func (t *Etcd) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()

	v, _ := t.Get(ctx, t.config.KeyPrefix+key)

	conv, _ := strconv.ParseInt(string(v), 10, 64)
	val += conv

	err := t.Set(ctx, t.config.KeyPrefix+key, fmt.Sprint(val), timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return val, err
}

func (t *Etcd) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()

	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

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
	err = t.SetMap(ctx, t.config.KeyPrefix+key, data, timeout)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return data.(map[string]any)[key2].(int64), err
}

func (t *Etcd) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return t.Increment(ctx, t.config.KeyPrefix+key, val*-1, timeout)
}

func (t *Etcd) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return t.IncrementIn(ctx, t.config.KeyPrefix+key, key2, -1*val, timeout)
}

func (t *Etcd) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	_, err := t.client.Delete(ctx, t.config.KeyPrefix+key)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return ErrKeyNotFound
	}

	return err
}

func (t *Etcd) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	resp, err := t.client.Get(ctx, t.config.KeyPrefix+key)
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
