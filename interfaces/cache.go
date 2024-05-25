package interfaces

import (
	"time"
)

type ICache interface {
	Init(app IService) error
	Stop() error
	String() string

	Has(key string) bool
	Set(key string, args any, timeout time.Duration) error
	SetIn(key string, key2 string, args any, timeout time.Duration) error
	SetMap(key string, args any, timeout time.Duration) error
	Get(key string) (any, error)
	GetIn(key string, key2 string) (any, error)
	GetMap(key string) (any, error)
	Increment(key string, val int64, timeout time.Duration) error
	IncrementIn(key string, key2 string, val int64, timeout time.Duration) error
	Decrement(key string, val int64, timeout time.Duration) error
	DecrementIn(key string, key2 string, val int64, timeout time.Duration) error
	Delete(key string) error
}
