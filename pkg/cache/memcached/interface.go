package memcached

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeMismatch = errors.New("type mismatch")

// IMemcached defines the interface for the Memcached struct, specifying behavior for interacting with Memcached.
type IMemcached interface {
	// Has checks the existence of a key in the Memcached store.
	Has(ctx context.Context, key string) bool

	// Set stores a value with the specified timeout.
	Set(ctx context.Context, key string, args any, timeout time.Duration) error

	// SetIn sets a nested value under a specific key in the Memcached store.
	SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error

	// SetMap converts a map to JSON format and stores it under a key.
	SetMap(ctx context.Context, key string, args any, timeout time.Duration) error

	// Get retrieves a value for the specified key.
	Get(ctx context.Context, key string) ([]byte, error)

	// GetIn retrieves a nested value from a stored map by key.
	GetIn(ctx context.Context, key string, key2 string) (any, error)

	// GetMap retrieves a JSON-stored map and unmarshals it.
	GetMap(ctx context.Context, key string) (map[string]any, error)

	// Increment increments a numeric value at a given key by a specified amount.
	Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// IncrementIn increments a nested numeric value in a stored map at a given key.
	IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// Decrement decrements a numeric value at a given key by a specified amount.
	Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// DecrementIn decrements a nested numeric value in a stored map at a given key.
	DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// Delete removes a value from the store by key.
	Delete(ctx context.Context, key string) error

	// Expire sets or updates the expiration timeout for a value stored by key.
	Expire(ctx context.Context, key string, timeout time.Duration) error
}
