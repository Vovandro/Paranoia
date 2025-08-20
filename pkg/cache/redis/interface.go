package redis

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeMismatch = errors.New("type mismatch")

// IRedis defines the interface for interacting with Redis.
type IRedis interface {
	// Init initializes the Redis instance with the given configuration.
	Init(cfg map[string]interface{}) error

	// Stop gracefully closes the Redis connection.
	Stop() error

	// Name returns the name of the Redis instance.
	Name() string

	// Type returns the type of cache.
	Type() string

	// Has checks the existence of a key in the Redis store.
	Has(ctx context.Context, key string) bool

	// Set stores a value in Redis under the specified key with a timeout.
	Set(ctx context.Context, key string, args any, timeout time.Duration) error

	// SetIn sets a value in a nested key of a map stored at the specified key with a timeout.
	SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error

	// SetMap stores a map in Redis under the specified key with a timeout.
	SetMap(ctx context.Context, key string, args any, timeout time.Duration) error

	// Get retrieves the value stored under the given key as a string.
	Get(ctx context.Context, key string) (string, error)

	// GetIn retrieves the value for a specific nested key of a map stored at the given key.
	GetIn(ctx context.Context, key string, key2 string) (string, error)

	// GetMap retrieves a map stored under the given key.
	GetMap(ctx context.Context, key string) (map[string]string, error)

	// Increment increases the numeric value stored at the given key by the specified amount with a timeout.
	Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// IncrementIn increases the numeric value stored under a nested key of a map in Redis.
	IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// Decrement decreases the numeric value stored at the given key by the specified amount with a timeout.
	Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)

	// DecrementIn decreases the numeric value stored under a nested key of a map in Redis.
	DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)

	// IncrementMany increases multiple keys by provided deltas atomically using a pipeline.
	// Returns the resulting values per key.
	IncrementMany(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error)

	// DecrementMany decreases multiple keys by provided deltas atomically using a pipeline.
	// Returns the resulting values per key.
	DecrementMany(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error)

	// Batch allows grouping multiple operations in a single pipeline round-trip.
	// The provided function is called with a Batcher to enqueue operations.
	Batch(ctx context.Context, fn func(Batcher)) error

	// Delete removes the value stored for the given key from Redis.
	Delete(ctx context.Context, key string) error

	// Expire sets or updates the expiration time for the given key.
	Expire(ctx context.Context, key string, timeout time.Duration) error
}

// Batcher provides a minimal set of batched operations supported by Batch.
type Batcher interface {
	Set(key string, value any, timeout time.Duration)
	Del(key string)
	Expire(key string, timeout time.Duration)
	HSet(key string, values map[string]any)
	HIncrBy(key, field string, incr int64)
	IncrBy(key string, incr int64)
	DecrBy(key string, decr int64)
}
