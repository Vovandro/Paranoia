package cache

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Mock struct {
	Name string
}

func (t *Mock) Init(_ interfaces.IEngine) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Has(ctx context.Context, key string) bool {
	return false
}

func (t *Mock) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) Get(ctx context.Context, key string) (any, error) {
	return nil, nil
}

func (t *Mock) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	return nil, nil
}

func (t *Mock) GetMap(ctx context.Context, key string) (any, error) {
	return nil, nil
}

func (t *Mock) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (t *Mock) Delete(ctx context.Context, key string) error {
	return nil
}

func (t *Mock) Expire(ctx context.Context, key string, timeout time.Duration) error {
	return nil
}
