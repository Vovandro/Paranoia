package cache

import (
	"fmt"
	"goServer"
	"sync"
	"time"
)

type Memory struct {
	Name string
	data sync.Map
}

func (t *Memory) Init(app *goServer.Service) error {
	return nil
}

func (t *Memory) Stop() error {
	return nil
}

func (t *Memory) String() string {
	return t.Name
}

func (t *Memory) Has(key string) bool {
	_, ok := t.data.Load(key)
	return ok
}

func (t *Memory) Set(key string, args any, timeout time.Duration) error {
	t.data.Store(key, args)

	time.AfterFunc(timeout, func() {
		_ = t.Delete(key)
	})

	return nil
}

func (t *Memory) Get(key string) (any, error) {
	val, ok := t.data.Load(key)

	if ok {
		return val, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (t *Memory) Delete(key string) error {
	t.data.Delete(key)
	return nil
}
