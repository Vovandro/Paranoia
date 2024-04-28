package cache

import (
	"Paranoia/interfaces"
	"fmt"
	"sync"
	"time"
)

type Memory struct {
	Name string
	data sync.Map
	pool sync.Pool
}

type cacheItem struct {
	data    any
	timeout time.Time
}

func (t *Memory) Init(app interfaces.IService) error {
	t.pool.New = func() any {
		return &cacheItem{}
	}

	return nil
}

func (t *Memory) Stop() error {
	return nil
}

func (t *Memory) String() string {
	return t.Name
}

func (t *Memory) Has(key string) bool {
	val, ok := t.data.Load(key)

	if ok {
		if val.(*cacheItem).timeout.Before(time.Now()) {
			t.data.Delete(key)
			t.pool.Put(val)
		}
	}

	return ok
}

func (t *Memory) Set(key string, args any, timeout time.Duration) error {
	val, ok := t.data.Load(key)

	if ok {
		val.(*cacheItem).timeout = time.Now().Add(timeout)
		val.(*cacheItem).data = args
	} else {
		val = t.pool.Get().(*cacheItem)
		val.(*cacheItem).timeout = time.Now().Add(timeout)
		val.(*cacheItem).data = args
		t.data.Store(key, val)
	}

	return nil
}

func (t *Memory) Get(key string) (any, error) {
	val, ok := t.data.Load(key)

	if ok {
		if val.(*cacheItem).timeout.After(time.Now()) {
			return val.(*cacheItem).data, nil
		} else {
			t.data.Delete(key)
			t.pool.Put(val)
		}
	}

	return nil, fmt.Errorf("key not found")
}

func (t *Memory) Delete(key string) error {
	val, ok := t.data.Load(key)

	if ok {
		t.data.Delete(key)
		t.pool.Put(val)
	}

	return nil
}
