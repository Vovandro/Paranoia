package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type Memcached struct {
	Name   string
	Config MemcachedConfig

	app          interfaces.IEngine
	client       *memcache.Client
	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type MemcachedConfig struct {
	Hosts   string        `yaml:"hosts"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewMemcached(name string, cfg MemcachedConfig) *Memcached {
	return &Memcached{
		Name:   name,
		Config: cfg,
	}
}

func (t *Memcached) Init(app interfaces.IEngine) error {
	t.app = app

	if t.Config.Timeout == 0 {
		t.Config.Timeout = 5 * time.Second
	}

	if t.Config.Hosts == "" {
		t.Config.Hosts = "localhost:11211"
	}

	t.client = memcache.New(strings.Split(t.Config.Hosts, ",")...)
	t.client.Timeout = t.Config.Timeout

	err := t.client.Ping()

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("memcached." + t.Name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("memcached." + t.Name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("memcached." + t.Name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("memcached." + t.Name + ".timeWrite")

	return nil
}

func (t *Memcached) Stop() error {
	return t.client.Close()
}

func (t *Memcached) String() string {
	return t.Name
}

func (t *Memcached) Has(ctx context.Context, key string) bool {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(key)

	if err != nil {
		return false
	}

	return item != nil
}

func (t *Memcached) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	var data []byte

	if _, ok := args.([]byte); ok {
		data = args.([]byte)
	} else if _, ok := args.(string); ok {
		data = []byte(args.(string))
	} else {
		data = []byte(fmt.Sprint(args))
	}

	return t.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(timeout.Seconds()),
	})
}

func (t *Memcached) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	data, err := t.GetMap(ctx, key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	data.(map[string]any)[key2] = args

	return t.SetMap(ctx, key, data, timeout)
}

func (t *Memcached) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	data, err := json.Marshal(args)

	if err != nil {
		return err
	}

	return t.Set(ctx, key, data, timeout)
}

func (t *Memcached) Get(ctx context.Context, key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	return item.Value, nil
}

func (t *Memcached) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	data, err := t.GetMap(ctx, key)

	if err != nil {
		return nil, err
	}

	if val, ok := data.(map[string]any)[key2]; ok {
		return val, nil
	}

	return nil, ErrKeyNotFound
}

func (t *Memcached) GetMap(ctx context.Context, key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	data, err := t.Get(ctx, key)

	if err != nil {
		return nil, err
	}

	res := make(map[string]any)
	err = json.Unmarshal(data.([]byte), &res)

	if err != nil {
		return nil, ErrTypeMismatch
	}

	return res, nil
}

func (t *Memcached) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	v, err := t.client.Increment(key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		return val, t.Set(ctx, key, val, timeout)
	} else if err != nil {
		return 0, err
	}

	return int64(v), t.client.Touch(key, int32(timeout.Seconds()))
}

func (t *Memcached) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

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

	return data.(map[string]any)[key2].(int64), t.SetMap(ctx, key, data, timeout)
}

func (t *Memcached) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	v, err := t.client.Decrement(key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		return 0, ErrKeyNotFound
	} else if err != nil {
		return 0, err
	}

	return int64(v), t.client.Touch(key, int32(timeout.Seconds()))
}

func (t *Memcached) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	data, err := t.GetMap(ctx, key)

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

	return data.(map[string]any)[key2].(int64), t.SetMap(ctx, key, data, timeout)
}

func (t *Memcached) Delete(ctx context.Context, key string) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	err := t.client.Delete(key)

	if errors.Is(err, memcache.ErrCacheMiss) {
		return ErrKeyNotFound
	}

	return err
}

func (t *Memcached) Expire(ctx context.Context, key string, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	err := t.client.Touch(key, int32(timeout.Seconds()))

	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return ErrKeyNotFound
		}
		return err
	}

	return nil
}
