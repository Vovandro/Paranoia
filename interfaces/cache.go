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
	Get(key string) (any, error)
	Delete(key string) error
}
