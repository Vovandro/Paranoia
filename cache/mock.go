package cache

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Mock struct {
	app  interfaces.IEngine
	Name string
}

func (t *Mock) Init(app interfaces.IEngine) error {
	t.app = app
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Has(key string) bool {
	return false
}

func (t *Mock) Set(key string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) SetIn(key string, key2 string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) SetMap(key string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) Get(key string) (any, error) {
	return nil, nil
}

func (t *Mock) GetIn(key string, key2 string) (any, error) {
	return nil, nil
}

func (t *Mock) GetMap(key string) (any, error) {
	return nil, nil
}

func (t *Mock) Increment(key string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) Delete(key string) error {
	return nil
}

func (t *Mock) Expire(key string, timeout time.Duration) error {
	return nil
}
