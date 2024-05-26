package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"strings"
	"time"
)

type Memcached struct {
	Name  string
	Hosts string

	app    interfaces.IService
	client *memcache.Client
}

func (t *Memcached) Init(app interfaces.IService) error {
	t.app = app

	t.client = memcache.New(strings.Split(t.Hosts, ",")...)

	err := t.client.Ping()

	if err != nil {
		return err
	}

	return nil
}

func (t *Memcached) Stop() error {
	return t.client.Close()
}

func (t *Memcached) String() string {
	return t.Name
}

func (t *Memcached) Has(key string) bool {
	item, err := t.client.Get(key)

	if err != nil {
		return false
	}

	return item != nil
}

func (t *Memcached) Set(key string, args any, timeout time.Duration) error {
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

func (t *Memcached) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	data, err := t.GetMap(key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	data.(map[string]any)[key2] = args

	return t.SetMap(key, data, timeout)
}

func (t *Memcached) SetMap(key string, args any, timeout time.Duration) error {
	data, err := json.Marshal(args)

	if err != nil {
		return err
	}

	return t.Set(key, data, timeout)
}

func (t *Memcached) Get(key string) (any, error) {
	item, err := t.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	return item.Value, nil
}

func (t *Memcached) GetIn(key string, key2 string) (any, error) {
	data, err := t.GetMap(key)

	if err != nil {
		return nil, err
	}

	if val, ok := data.(map[string]any)[key2]; ok {
		return val, nil
	}

	return nil, ErrKeyNotFound
}

func (t *Memcached) GetMap(key string) (any, error) {
	data, err := t.Get(key)

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

func (t *Memcached) Increment(key string, val int64, timeout time.Duration) error {
	_, err := t.client.Increment(key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		return t.Set(key, val, timeout)
	} else if err != nil {
		return err
	}

	return t.client.Touch(key, int32(timeout.Seconds()))
}

func (t *Memcached) IncrementIn(key string, key2 string, val int64, timeout time.Duration) error {
	data, err := t.GetMap(key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	if v, ok := data.(map[string]any)[key2]; ok {
		data.(map[string]any)[key2] = v.(int64) + val
	} else {
		data.(map[string]any)[key2] = val
	}

	return t.SetMap(key, data, timeout)
}

func (t *Memcached) Decrement(key string, val int64, timeout time.Duration) error {
	_, err := t.client.Decrement(key, uint64(val))

	if errors.Is(err, memcache.ErrCacheMiss) {
		return ErrKeyNotFound
	} else if err != nil {
		return err
	}

	return t.client.Touch(key, int32(timeout.Seconds()))
}

func (t *Memcached) DecrementIn(key string, key2 string, val int64, timeout time.Duration) error {
	data, err := t.GetMap(key)

	if errors.Is(err, ErrKeyNotFound) {
		data = make(map[string]any)
	} else if err != nil {
		return err
	}

	if v, ok := data.(map[string]any)[key2]; ok {
		data.(map[string]any)[key2] = v.(int64) - val
	} else {
		data.(map[string]any)[key2] = val * -1
	}

	return t.SetMap(key, data, timeout)
}

func (t *Memcached) Delete(key string) error {
	err := t.client.Delete(key)

	if errors.Is(err, memcache.ErrCacheMiss) {
		return ErrKeyNotFound
	}

	return err
}
