package cache

import (
	"Paranoia/interfaces"
	"fmt"
	"sync"
	"time"
)

type Memory struct {
	Name  string
	data  map[string]*cacheItem
	pool  sync.Pool
	mutex sync.RWMutex
}

type cacheItem struct {
	data    any
	timeout time.Time
}

func (t *Memory) Init(app interfaces.IService) error {
	t.pool.New = func() any {
		return &cacheItem{}
	}

	t.data = make(map[string]*cacheItem, 100)

	return nil
}

func (t *Memory) Stop() error {
	return nil
}

func (t *Memory) String() string {
	return t.Name
}

func (t *Memory) Has(key string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	val, ok := t.data[key]

	if ok {
		if val.timeout.Before(time.Now()) {
			t.mutex.Lock()
			defer t.mutex.Unlock()

			delete(t.data, key)
			t.pool.Put(val)
		}
	}

	return ok
}

func (t *Memory) Set(key string, args any, timeout time.Duration) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		val.timeout = time.Now().Add(timeout)
		val.data = args
	} else {
		val = t.pool.Get().(*cacheItem)
		val.timeout = time.Now().Add(timeout)
		val.data = args
		t.data[key] = val
	}

	return nil
}

func (t *Memory) Get(key string) (any, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	val, ok := t.data[key]

	if ok {
		if val.timeout.After(time.Now()) {
			return val.data, nil
		} else {
			t.mutex.Lock()
			defer t.mutex.Unlock()

			delete(t.data, key)
			t.pool.Put(val)
		}
	}

	return nil, fmt.Errorf("key not found")
}

func (t *Memory) Delete(key string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		delete(t.data, key)
		t.pool.Put(val)
	}

	return nil
}
