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
	Name       string
	UseCluster bool
	Hosts      string

	app     interfaces.IService
	client  *redis.Client
	cluster *redis.ClusterClient
}

func (t *Redis) Init(app interfaces.IService) error {
	t.app = app

	if t.UseCluster {
		t.cluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(t.Hosts, ","),
		})

		_, err := t.cluster.Ping(context.Background()).Result()

		if err != nil {
			return err
		}
	} else {
		t.client = redis.NewClient(&redis.Options{
			Addr: t.Hosts,
		})

		_, err := t.client.Ping(context.Background()).Result()

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Redis) Stop() error {
	if t.UseCluster {
		return t.cluster.Close()
	}

	return t.client.Close()
}

func (t *Redis) String() string {
	return t.Name
}

func (t *Redis) Has(key string) bool {
	if t.UseCluster {
		return t.cluster.Exists(context.Background(), key).Val() != 0
	}

	return t.client.Exists(context.Background(), key).Val() != 0
}

func (t *Redis) Set(key string, args any, timeout time.Duration) error {
	if t.UseCluster {
		return t.cluster.Set(context.Background(), key, args, timeout).Err()
	}

	return t.client.Set(context.Background(), key, args, timeout).Err()
}

func (t *Redis) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.HSet(context.Background(), key, map[string]interface{}{key2: args}).Err()

		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.HSet(context.Background(), key, map[string]interface{}{key2: args}).Err()

	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) SetMap(key string, args any, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.HSet(context.Background(), key, args).Err()
		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.HSet(context.Background(), key, args).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Get(key string) (any, error) {
	var v string
	var err error

	if t.UseCluster {
		v, err = t.cluster.Get(context.Background(), key).Result()
	} else {
		v, err = t.client.Get(context.Background(), key).Result()
	}

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetIn(key string, key2 string) (any, error) {
	var v string
	var err error

	if t.UseCluster {
		v, err = t.cluster.HGet(context.Background(), key, key2).Result()
	} else {
		v, err = t.client.HGet(context.Background(), key, key2).Result()
	}

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) GetMap(key string) (any, error) {
	var v map[string]string
	var err error

	if t.UseCluster {
		v, err = t.cluster.HGetAll(context.Background(), key).Result()
	} else {
		v, err = t.client.HGetAll(context.Background(), key).Result()
	}

	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (t *Redis) Increment(key string, val int64, timeout time.Duration) (int64, error) {
	if t.UseCluster {
		v := t.cluster.IncrBy(context.Background(), key, val)
		if v.Err() != nil {
			return 0, v.Err()
		}

		return v.Val(), t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	v := t.client.IncrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	if t.UseCluster {
		v := t.cluster.HIncrBy(context.Background(), key, key2, val)
		if v.Err() != nil {
			return 0, v.Err()
		}

		return v.Val(), t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	v := t.client.HIncrBy(context.Background(), key, key2, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	if t.UseCluster {
		v := t.cluster.DecrBy(context.Background(), key, val)
		if v.Err() != nil {
			return 0, v.Err()
		}

		return v.Val(), t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	v := t.client.DecrBy(context.Background(), key, val)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	if t.UseCluster {
		v := t.cluster.HIncrBy(context.Background(), key, key2, val*-1)
		if v.Err() != nil {
			return 0, v.Err()
		}

		return v.Val(), t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	v := t.client.HIncrBy(context.Background(), key, key2, val*-1)
	if v.Err() != nil {
		return 0, v.Err()
	}

	return v.Val(), t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Delete(key string) error {
	if t.UseCluster {
		return t.cluster.Del(context.Background(), key).Err()
	}

	return t.client.Del(context.Background(), key).Err()
}
