package cache

import (
	"context"
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
		return t.cluster.HSet(context.Background(), key, key2, args, timeout).Err()
	}

	return t.client.HSet(context.Background(), key, key, args, timeout).Err()
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
	if t.UseCluster {
		return t.cluster.Get(context.Background(), key).Result()
	}

	return t.client.Get(context.Background(), key).Result()
}

func (t *Redis) GetIn(key string, key2 string) (any, error) {
	if t.UseCluster {
		return t.cluster.HGet(context.Background(), key, key2).Result()
	}

	return t.client.HGet(context.Background(), key, key2).Result()
}

func (t *Redis) GetMap(key string) (any, error) {
	if t.UseCluster {
		return t.cluster.HGetAll(context.Background(), key).Result()
	}

	return t.client.HGetAll(context.Background(), key).Result()
}

func (t *Redis) Increment(key string, val int64, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.IncrBy(context.Background(), key, val).Err()
		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.IncrBy(context.Background(), key, val).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) IncrementIn(key string, key2 string, val int64, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.HIncrBy(context.Background(), key, key2, val).Err()
		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.HIncrBy(context.Background(), key, key2, val).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Decrement(key string, val int64, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.DecrBy(context.Background(), key, val).Err()
		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.DecrBy(context.Background(), key, val).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) DecrementIn(key string, key2 string, val int64, timeout time.Duration) error {
	if t.UseCluster {
		err := t.cluster.HIncrBy(context.Background(), key, key2, val*-1).Err()
		if err != nil {
			return err
		}

		return t.cluster.Expire(context.Background(), key, timeout).Err()
	}

	err := t.client.HIncrBy(context.Background(), key, key2, val*-1).Err()
	if err != nil {
		return err
	}

	return t.client.Expire(context.Background(), key, timeout).Err()
}

func (t *Redis) Delete(key string) error {
	if t.UseCluster {
		return t.cluster.Del(context.Background(), key).Err()
	}

	return t.client.Del(context.Background(), key).Err()
}
