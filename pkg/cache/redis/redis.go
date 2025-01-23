package redis

import (
	"context"
	"errors"
	redisExt "github.com/redis/go-redis/v9"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type Redis struct {
	name   string
	config Config

	client redisExt.UniversalClient

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type Config struct {
	Hosts      string        `yaml:"hosts"`
	UseCluster bool          `yaml:"use_cluster"`
	DBNum      int           `yaml:"db_num"`
	Timeout    time.Duration `yaml:"timeout"`
	Username   string        `yaml:"username"`
	Password   string        `yaml:"password"`
	KeyPrefix  string        `yaml:"key_prefix"`
}

func New(name string) *Redis {
	return &Redis{
		name: name,
	}
}

func (t *Redis) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	if t.config.UseCluster {
		if t.config.DBNum != 0 {
			return errors.New("database number not available when using Redis cluster")
		}

		t.client = redisExt.NewClusterClient(&redisExt.ClusterOptions{
			Addrs:    strings.Split(t.config.Hosts, ","),
			Username: t.config.Username,
			Password: t.config.Password,
		})
	} else {
		t.client = redisExt.NewClient(&redisExt.Options{
			Addr:     t.config.Hosts,
			DB:       t.config.DBNum,
			Username: t.config.Username,
			Password: t.config.Password,
		})
	}

	_, err = t.client.Ping(context.Background()).Result()

	if err != nil {
		return err
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("redis." + t.name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("redis." + t.name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("redis." + t.name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("redis." + t.name + ".timeWrite")

	return nil
}

func (t *Redis) Stop() error {
	return t.client.Close()
}

func (t *Redis) Name() string {
	return t.name
}

func (t *Redis) Type() string {
	return "cache"
}

func (t *Redis) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	res := t.client.Exists(ctx, t.config.KeyPrefix+key).Val() != 0
	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	return res
}

func (t *Redis) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Set(ctx, t.config.KeyPrefix+key, args, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.HSet(ctx, t.config.KeyPrefix+key, map[string]interface{}{key2: args}).Err()

	if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	err = t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.HSet(ctx, t.config.KeyPrefix+key, args).Err()
	if err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	err = t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) Get(ctx context.Context, key string) (string, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.Get(ctx, t.config.KeyPrefix+key).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redisExt.Nil) {
		return "", ErrKeyNotFound
	} else if err != nil {
		return "", err
	}

	return v, nil
}

func (t *Redis) GetIn(ctx context.Context, key string, key2 string) (string, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.HGet(ctx, t.config.KeyPrefix+key, key2).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redisExt.Nil) {
		return "", ErrKeyNotFound
	} else if err != nil {
		return "", err
	}

	return v, nil
}

func (t *Redis) GetMap(ctx context.Context, key string) (map[string]string, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	v, err := t.client.HGetAll(ctx, t.config.KeyPrefix+key).Result()

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redisExt.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.IncrBy(ctx, t.config.KeyPrefix+key, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.HIncrBy(ctx, t.config.KeyPrefix+key, key2, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.DecrBy(ctx, t.config.KeyPrefix+key, val)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	v := t.client.HIncrBy(ctx, t.config.KeyPrefix+key, key2, val*-1)
	if v.Err() != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, v.Err()
	}

	err := t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Val(), err
}

func (t *Redis) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Del(ctx, t.config.KeyPrefix+key).Err()
	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return err
}

func (t *Redis) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	err := t.client.Expire(ctx, t.config.KeyPrefix+key, timeout).Err()

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	if errors.Is(err, redisExt.Nil) {
		return ErrKeyNotFound
	} else if err != nil {
		return err
	}

	return nil
}
