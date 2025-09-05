package redis

import (
	"context"
	"strconv"
	"sync"
	"time"
)

// BatchOp describes a single batched operation captured by MockBatcher
type BatchOp struct {
	Method  string
	Key     string
	Value   any
	Timeout time.Duration
	Values  map[string]any
	Field   string
	Delta   int64
}

// Mock implements IRedis for tests with hookable behavior and call recording.
type Mock struct {
	HasFunc           func(ctx context.Context, key string) bool
	SetFunc           func(ctx context.Context, key string, args any, timeout time.Duration) error
	SetInFunc         func(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error
	SetMapFunc        func(ctx context.Context, key string, args any, timeout time.Duration) error
	GetFunc           func(ctx context.Context, key string) (string, error)
	GetInFunc         func(ctx context.Context, key string, key2 string) (string, error)
	GetMapFunc        func(ctx context.Context, key string) (map[string]string, error)
	IncrementFunc     func(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)
	IncrementInFunc   func(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)
	DecrementFunc     func(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)
	DecrementInFunc   func(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)
	IncrementManyFunc func(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error)
	DecrementManyFunc func(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error)
	BatchFunc         func(ctx context.Context, ops []BatchOp) error
	DeleteFunc        func(ctx context.Context, key string) error
	ExpireFunc        func(ctx context.Context, key string, timeout time.Duration) error

	NamePkg string

	// In-memory stores for default behavior when hooks are not set
	Data map[string]string
	Maps map[string]map[string]string
	Nums map[string]int64

	mu      sync.Mutex
	Calls   []struct{ Method, Key, Key2 string }
	Batches [][]BatchOp
}

func (m *Mock) record(method, key, key2 string) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct{ Method, Key, Key2 string }{Method: method, Key: key, Key2: key2})
	m.mu.Unlock()
}

func (t *Mock) Init(cfg map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Type() string {
	return "cache"
}

func (m *Mock) Has(ctx context.Context, key string) bool {
	m.record("Has", key, "")
	if m.HasFunc != nil {
		return m.HasFunc(ctx, key)
	}
	if m.Data == nil {
		return false
	}
	_, ok := m.Data[key]
	return ok
}

func (m *Mock) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	m.record("Set", key, "")
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, args, timeout)
	}
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	switch v := args.(type) {
	case string:
		m.Data[key] = v
	case []byte:
		m.Data[key] = string(v)
	default:
		m.Data[key] = toString(v)
	}
	return nil
}

func (m *Mock) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	m.record("SetIn", key, key2)
	if m.SetInFunc != nil {
		return m.SetInFunc(ctx, key, key2, args, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]string)
	}
	if _, ok := m.Maps[key]; !ok {
		m.Maps[key] = make(map[string]string)
	}
	switch v := args.(type) {
	case string:
		m.Maps[key][key2] = v
	case []byte:
		m.Maps[key][key2] = string(v)
	default:
		m.Maps[key][key2] = toString(v)
	}
	return nil
}

func (m *Mock) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	m.record("SetMap", key, "")
	if m.SetMapFunc != nil {
		return m.SetMapFunc(ctx, key, args, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]string)
	}
	if mp, ok := args.(map[string]string); ok {
		m.Maps[key] = mp
		return nil
	}
	if mp, ok := args.(map[string]any); ok {
		conv := make(map[string]string, len(mp))
		for k, v := range mp {
			conv[k] = toString(v)
		}
		m.Maps[key] = conv
		return nil
	}
	return nil
}

func (m *Mock) Get(ctx context.Context, key string) (string, error) {
	m.record("Get", key, "")
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	if m.Data == nil {
		return "", ErrKeyNotFound
	}
	v, ok := m.Data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return v, nil
}

func (m *Mock) GetIn(ctx context.Context, key string, key2 string) (string, error) {
	m.record("GetIn", key, key2)
	if m.GetInFunc != nil {
		return m.GetInFunc(ctx, key, key2)
	}
	if m.Maps == nil {
		return "", ErrKeyNotFound
	}
	mp, ok := m.Maps[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	v, ok := mp[key2]
	if !ok {
		return "", ErrKeyNotFound
	}
	return v, nil
}

func (m *Mock) GetMap(ctx context.Context, key string) (map[string]string, error) {
	m.record("GetMap", key, "")
	if m.GetMapFunc != nil {
		return m.GetMapFunc(ctx, key)
	}
	if m.Maps == nil {
		return nil, ErrKeyNotFound
	}
	mp, ok := m.Maps[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	out := make(map[string]string, len(mp))
	for k, v := range mp {
		out[k] = v
	}
	return out, nil
}

func (m *Mock) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	m.record("Increment", key, "")
	if m.IncrementFunc != nil {
		return m.IncrementFunc(ctx, key, val, timeout)
	}
	if m.Nums == nil {
		m.Nums = make(map[string]int64)
	}
	m.Nums[key] += val
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	m.Data[key] = strconv.FormatInt(m.Nums[key], 10)
	return m.Nums[key], nil
}

func (m *Mock) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	m.record("IncrementIn", key, key2)
	if m.IncrementInFunc != nil {
		return m.IncrementInFunc(ctx, key, key2, val, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]string)
	}
	if _, ok := m.Maps[key]; !ok {
		m.Maps[key] = make(map[string]string)
	}
	curr := int64(0)
	if v, ok := m.Maps[key][key2]; ok {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			curr = n
		}
	}
	curr += val
	m.Maps[key][key2] = strconv.FormatInt(curr, 10)
	return curr, nil
}

func (m *Mock) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	m.record("Decrement", key, "")
	if m.DecrementFunc != nil {
		return m.DecrementFunc(ctx, key, val, timeout)
	}
	return m.Increment(ctx, key, -val, timeout)
}

func (m *Mock) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	m.record("DecrementIn", key, key2)
	if m.DecrementInFunc != nil {
		return m.DecrementInFunc(ctx, key, key2, val, timeout)
	}
	return m.IncrementIn(ctx, key, key2, -val, timeout)
}

func (m *Mock) IncrementMany(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error) {
	if m.IncrementManyFunc != nil {
		return m.IncrementManyFunc(ctx, deltas, timeout)
	}
	if len(deltas) == 0 {
		return map[string]int64{}, nil
	}
	if m.Nums == nil {
		m.Nums = make(map[string]int64)
	}
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	out := make(map[string]int64, len(deltas))
	for k, v := range deltas {
		m.Nums[k] += v
		m.Data[k] = strconv.FormatInt(m.Nums[k], 10)
		out[k] = m.Nums[k]
	}
	return out, nil
}

func (m *Mock) DecrementMany(ctx context.Context, deltas map[string]int64, timeout time.Duration) (map[string]int64, error) {
	if m.DecrementManyFunc != nil {
		return m.DecrementManyFunc(ctx, deltas, timeout)
	}
	if len(deltas) == 0 {
		return map[string]int64{}, nil
	}
	neg := make(map[string]int64, len(deltas))
	for k, v := range deltas {
		neg[k] = -v
	}
	return m.IncrementMany(ctx, neg, timeout)
}

func (m *Mock) Delete(ctx context.Context, key string) error {
	m.record("Delete", key, "")
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key)
	}
	if m.Data != nil {
		delete(m.Data, key)
	}
	if m.Maps != nil {
		delete(m.Maps, key)
	}
	if m.Nums != nil {
		delete(m.Nums, key)
	}
	return nil
}

func (m *Mock) Expire(ctx context.Context, key string, timeout time.Duration) error {
	m.record("Expire", key, "")
	if m.ExpireFunc != nil {
		return m.ExpireFunc(ctx, key, timeout)
	}
	return nil
}

// MockBatcher captures batched operations for tests
type MockBatcher struct{ ops *[]BatchOp }

func (b *MockBatcher) Set(key string, value any, timeout time.Duration) {
	*b.ops = append(*b.ops, BatchOp{Method: "Set", Key: key, Value: value, Timeout: timeout})
}

func (b *MockBatcher) Del(key string) { *b.ops = append(*b.ops, BatchOp{Method: "Del", Key: key}) }

func (b *MockBatcher) Expire(key string, timeout time.Duration) {
	*b.ops = append(*b.ops, BatchOp{Method: "Expire", Key: key, Timeout: timeout})
}

func (b *MockBatcher) HSet(key string, values map[string]any) {
	*b.ops = append(*b.ops, BatchOp{Method: "HSet", Key: key, Values: values})
}

func (b *MockBatcher) HIncrBy(key, field string, incr int64) {
	*b.ops = append(*b.ops, BatchOp{Method: "HIncrBy", Key: key, Field: field, Delta: incr})
}

func (b *MockBatcher) IncrBy(key string, incr int64) {
	*b.ops = append(*b.ops, BatchOp{Method: "IncrBy", Key: key, Delta: incr})
}

func (b *MockBatcher) DecrBy(key string, decr int64) {
	*b.ops = append(*b.ops, BatchOp{Method: "DecrBy", Key: key, Delta: decr})
}

func (m *Mock) Batch(ctx context.Context, fn func(Batcher)) error {
	ops := make([]BatchOp, 0)
	b := &MockBatcher{ops: &ops}
	fn(b)
	m.mu.Lock()
	m.Batches = append(m.Batches, ops)
	m.mu.Unlock()
	if m.BatchFunc != nil {
		return m.BatchFunc(ctx, ops)
	}
	return nil
}

var _ IRedis = (*Mock)(nil)
var _ Batcher = (*MockBatcher)(nil)

// toString converts common types to string for default value handling
func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
