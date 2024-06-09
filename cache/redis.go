package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"strings"
	"time"
)

type Redis struct {
	Name   string
	Config RedisConfig

	app    interfaces.IService
	client iRedis
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
	Hosts      string
	UseCluster bool
	DBNum      int
	Timeout    time.Duration
	Username   string
	Password   string
}

func (t *Redis) Init(app interfaces.IService) error {
	t.app = app

	if t.Config.UseCluster {
		if t.Config.DBNum != 0 {
			t.app.GetLogger().Warn("Database number not available when using Redis cluster")
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

	return nil
}

func (t *Redis) Stop() error {
	return t.client.Close()
}

func (t *Redis) String() string {
	return t.Name
}

func (t *Redis) Has(key string) bool {
	return t.client.Exists(context.Background(), key).Val() != 0
}

func (t *Redis) Set(key string, args any, timeout time.Duration) error {
	return t.client.Set(context.Background(), key, args, timeout).Err()
}

func (t *Redis) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	err := t.client.HSet(context.Background(), key, map[string]interface{}{key2: args}).Err()

	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) SetMap(key string, args any, timeout time.Duration) error {
	err := t.client.HSet(context.Background(), key, args).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Get(key string) (any, error) {
	v, err := t.client.Get(context.Background(), key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetIn(key string, key2 string) (any, error) {
	v, err := t.client.HGet(context.Background(), key, key2).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetMap(key string) (any, error) {
	v, err := t.client.HGetAll(context.Background(), key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) Increment(key string, val int64, timeout time.Duration) (int64, error) {
	v := t.client.IncrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	v := t.client.HIncrBy(context.Background(), key, key2, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	v := t.client.DecrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	v := t.client.HIncrBy(context.Background(), key, key2, val*-1)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Delete(key string) error {
	return t.client.Del(context.Background(), key).Err()
}
