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

func (t *Etcd) Has(key string) bool {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	res, err := t.client.Get(context.Background(), key)

	if err != nil {
		return false
	}

	if len(res.Kvs) == 0 {
		return false
	}

	return true
}

func (t *Etcd) Set(key string, args any, timeout time.Duration) error {
	if _, ok := args.(string); !ok {
		return ErrTypeMismatch
	}

	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)
	lease, _ := t.client.Grant(context.Background(), int64(timeout.Seconds()))

	_, err := t.client.Put(context.Background(), key, args.(string), clientv3.WithLease(lease.ID))

	return err
}

func (t *Etcd) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	data, err := t.GetMap(key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	data.(map[string]any)[key2] = args

	return t.SetMap(key, data, timeout)
}

func (t *Etcd) SetMap(key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	data, err := json.Marshal(args)

	if err != nil {
		return err
	}

	return t.Set(key, string(data), timeout)
}

func (t *Etcd) Get(key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	item, err := t.client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if len(item.Kvs) == 0 {
		return "", ErrKeyNotFound
	}

	return string(item.Kvs[0].Value), nil
}

func (t *Etcd) GetIn(key string, key2 string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	data, err := t.GetMap(key)

	if err != nil {
		return "", err
	}

	if val, ok := data.(map[string]any)[key2]; ok {
		return val, nil
	}

	return "", ErrKeyNotFound
}

func (t *Etcd) GetMap(key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	data, err := t.Get(key)

	if err != nil {
		return "", err
	}

	res := make(map[string]any)
	err = json.Unmarshal([]byte(data.(string)), &res)

	if err != nil {
		return "", ErrTypeMismatch
	}

	return res, nil
}

func (t *Etcd) Increment(key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	v, _ := t.Get(key)

	conv, _ := strconv.ParseInt(v.(string), 10, 64)
	val += conv

	return val, t.Set(key, fmt.Sprint(val), timeout)
}

func (t *Etcd) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	data, err := t.GetMap(key)

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

	return data.(map[string]any)[key2].(int64), t.SetMap(key, data, timeout)
}

func (t *Etcd) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	return t.Increment(key, val*-1, timeout)
}

func (t *Etcd) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	data, err := t.GetMap(key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return 0, err
	}

	if v, ok := data.(map[string]any)[key2]; ok {
		data.(map[string]any)[key2] = int64(v.(float64)) - val
	} else {
		data.(map[string]any)[key2] = val * -1
	}

	return data.(map[string]any)[key2].(int64), t.SetMap(key, data, timeout)
}

func (t *Etcd) Delete(key string) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	_, err := t.client.Delete(context.Background(), key)

	if err != nil {
		return ErrKeyNotFound
	}

	return err
}

func (t *Etcd) Expire(key string, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	resp, err := t.client.Get(context.Background(), key)
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return ErrKeyNotFound
	}

	_, err = t.client.KeepAliveOnce(context.Background(), clientv3.LeaseID(resp.Kvs[0].Lease))

	return err
}
