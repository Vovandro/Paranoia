package memory

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeMismatch = errors.New("type mismatch")

// IMemory defines the interface for the Memory struct, specifying behavior for interacting with Memory.
type IMemory interface {
	// Has checks the existence of a key in the cache.
	Has(ctx context.Context, key string) bool

	// Set stores a value in the cache under the given key with the specified Timeout.
	Set(ctx context.Context, key string, args any, timeout time.Duration) error

	// SetIn sets a value within a nested key for the given key in the cache with a Timeout.
	SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error

	// SetMap stores a map in the cache after converting it to JSON with the specified Timeout.
	SetMap(ctx context.Context, key string, args any, timeout time.Duration) error

	// Get retrieves the raw value stored under the given key.
	Get(ctx context.Context, key string) (any, error)

	// GetIn retrieves a value for a specific nested key within a map stored in the cache.
	GetIn(ctx context.Context, key string, key2 string) (any, error)

	// GetMap retrieves a map stored under the given key by decoding its JSON value.
	GetMap(ctx context.Context, key string) (any, error)

	// Increment increments the numeric value stored at the given key by the specified amount and applies a Timeout.
	Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// IncrementIn increments the numeric value stored under a nested key within a map in the cache.
	IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// Decrement decrements the numeric value stored at the given key by the specified amount and applies a Timeout.
	Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// DecrementIn decrements the numeric value stored under a nested key within a map in the cache.
	DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// Delete removes the value stored under the given key from the cache.
	Delete(ctx context.Context, key string) error

	// Expire updates the expiration time for the value stored under the given key.
	Expire(ctx context.Context, key string, timeout time.Duration) error
}
