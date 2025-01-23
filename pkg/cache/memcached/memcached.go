package memcached

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type Memcached struct {
	name   string
	config Config

	client       *memcache.Client
	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type Config struct {
	Hosts     string        `yaml:"hosts"`
	Timeout   time.Duration `yaml:"timeout"`
	KeyPrefix string        `yaml:"key_prefix"`
}

func New(name string) *Memcached {
	return &Memcached{
		name: name,
	}
}

func (t *Memcached) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Timeout == 0 {
		t.config.Timeout = 5 * time.Second
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	t.client = memcache.New(strings.Split(t.config.Hosts, ",")...)
	t.client.Timeout = t.config.Timeout

	err = t.client.Ping()

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("memcached." + t.name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("memcached." + t.name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("memcached." + t.name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("memcached." + t.name + ".timeWrite")

	return nil
}

func (t *Memcached) Stop() error {
	return t.client.Close()
}

func (t *Memcached) Name() string {
	return t.name
}

func (t *Memcached) Type() string {
	return "cache"
}

func (t *Memcached) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(t.config.KeyPrefix + key)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		return false
	}

	return item != nil
}

func (t *Memcached) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	var data []byte

	if _, ok := args.([]byte); ok {
		data = args.([]byte)
	} else if _, ok := args.(string); ok {
		data = []byte(args.(string))
	} else {
		data = []byte(fmt.Sprint(args))
	}

	err := t.client.Set(&memcache.Item{
		Key:        t.config.KeyPrefix + key,
		Value:      data,
		Expiration: int32(timeout.Seconds()),
	})
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Memcached) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	data[key2] = args

	return t.SetMap(ctx, t.config.KeyPrefix+key, data, timeout)
}

func (t *Memcached) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	data, err := json.Marshal(args)

	if err != nil {
		return err
	}

	return t.Set(ctx, t.config.KeyPrefix+key, data, timeout)
}

func (t *Memcached) Get(ctx context.Context, key string) ([]byte, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	item, err := t.client.Get(t.config.KeyPrefix + key)
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	return item.Value, nil
}

func (t *Memcached) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

	if err != nil {
		return nil, err
	}

	if val, ok := data[key2]; ok {
		return val, nil
	}

	return nil, ErrKeyNotFound
}

func (t *Memcached) GetMap(ctx context.Context, key string) (map[string]any, error) {
	data, err := t.Get(ctx, t.config.KeyPrefix+key)

	if err != nil {
		return nil, err
	}

	res := make(map[string]any)
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, ErrTypeMismatch
	}

	return res, nil
}

func (t *Memcached) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v, err := t.client.Increment(t.config.KeyPrefix+key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		return val, t.Set(ctx, t.config.KeyPrefix+key, val, timeout)
	} else if err != nil {
		return 0, err
	}

	err = t.client.Touch(t.config.KeyPrefix+key, int32(timeout.Seconds()))

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return int64(v), err
}

func (t *Memcached) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, err
	}

	if v, ok := data[key2]; ok {
		data[key2] = int64(v.(float64)) + val
	} else {
		data[key2] = val
	}

	err = t.SetMap(ctx, t.config.KeyPrefix+key, data, timeout)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return data[key2].(int64), err
}

func (t *Memcached) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v, err := t.client.Decrement(t.config.KeyPrefix+key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, ErrKeyNotFound
	} else if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, err
	}
	err = t.client.Touch(t.config.KeyPrefix+key, int32(timeout.Seconds()))

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return int64(v), err
}

func (t *Memcached) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	data, err := t.GetMap(ctx, t.config.KeyPrefix+key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return 0, err
	}

	if v, ok := data[key2]; ok {
		data[key2] = int64(v.(float64)) - val
	} else {
		data[key2] = val * -1
	}

	err = t.SetMap(ctx, t.config.KeyPrefix+key, data, timeout)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return data[key2].(int64), err
}

func (t *Memcached) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Delete(t.config.KeyPrefix + key)
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	if errors.Is(err, memcache.ErrCacheMiss) {
		return ErrKeyNotFound
	}

	return err
}

func (t *Memcached) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Touch(t.config.KeyPrefix+key, int32(timeout.Seconds()))
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return ErrKeyNotFound
		}
		return err
	}

	return nil
}
