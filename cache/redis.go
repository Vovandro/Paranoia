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
	client redis.UniversalClient

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
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

func (t *Redis) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	res := t.client.Exists(ctx, key).Val() != 0
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	return res
}

func (t *Redis) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Set(ctx, key, args, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.HSet(ctx, key, map[string]interface{}{key2: args}).Err()

	if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	err = t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.HSet(ctx, key, args).Err()
	if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	err = t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) Get(ctx context.Context, key string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.Get(ctx, key).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.HGet(ctx, key, key2).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetMap(ctx context.Context, key string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.HGetAll(ctx, key).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.IncrBy(ctx, key, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.HIncrBy(ctx, key, key2, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.DecrBy(ctx, key, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.HIncrBy(ctx, key, key2, val*-1)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Del(ctx, key).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Expire(ctx, key, timeout).Err()

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redis.Nil) {
		return ErrKeyNotFound
	} else if err != nil {
		return err
	}

	return nil
}
