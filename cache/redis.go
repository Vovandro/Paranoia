package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type Redis struct {
	Name   string
	Config RedisConfig

	app    interfaces.IEngine
	client iRedis

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type iRedis interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd
	DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisConfig struct {
	Hosts      string        `yaml:"hosts"`
	UseCluster bool          `yaml:"use_cluster,omitempty"`
	DBNum      int           `yaml:"db_num,omitempty"`
	Timeout    time.Duration `yaml:"timeout"`
	Username   string        `yaml:"username,omitempty"`
	Password   string        `yaml:"password,omitempty"`
}

func NewRedis(name string, cfg RedisConfig) *Redis {
	return &Redis{
		Name:   name,
		Config: cfg,
	}
}

func (t *Redis) Init(app interfaces.IEngine) error {
	t.app = app

	if t.Config.UseCluster {
		if t.Config.DBNum != 0 {
			t.app.GetLogger().Warn(context.Background(), "Database number not available when using Redis cluster")
		}

		t.client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(t.Config.Hosts, ","),
			Username: t.Config.Username,
			Password: t.Config.Password,
		})
	} else {
		t.client = redis.NewClient(&redis.Options{
			Addr:     t.Config.Hosts,
			DB:       t.Config.DBNum,
			Username: t.Config.Username,
			Password: t.Config.Password,
		})
	}

	_, err := t.client.Ping(context.Background()).Result()

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("redis." + t.Name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("redis." + t.Name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("redis." + t.Name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("redis." + t.Name + ".timeWrite")

	return nil
}

func (t *Redis) Stop() error {
	return t.client.Close()
}

func (t *Redis) String() string {
	return t.Name
}

func (t *Redis) Has(key string) bool {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	return t.client.Exists(context.Background(), key).Val() != 0
}

func (t *Redis) Set(key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	return t.client.Set(context.Background(), key, args, timeout).Err()
}

func (t *Redis) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	err := t.client.HSet(context.Background(), key, map[string]interface{}{key2: args}).Err()

	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) SetMap(key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	err := t.client.HSet(context.Background(), key, args).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Get(key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	v, err := t.client.Get(context.Background(), key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetIn(key string, key2 string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	v, err := t.client.HGet(context.Background(), key, key2).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetMap(key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(context.Background(), 1)

	v, err := t.client.HGetAll(context.Background(), key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) Increment(key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	v := t.client.IncrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	v := t.client.HIncrBy(context.Background(), key, key2, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	v := t.client.DecrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	v := t.client.HIncrBy(context.Background(), key, key2, val*-1)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Delete(key string) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	return t.client.Del(context.Background(), key).Err()
}

func (t *Redis) Expire(key string, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(context.Background(), 1)

	err := t.client.Expire(context.Background(), key, timeout).Err()

	if errors.Is(err, redis.Nil) {
		return ErrKeyNotFound
	} else if err != nil {
		return err
	}

	return nil
}
