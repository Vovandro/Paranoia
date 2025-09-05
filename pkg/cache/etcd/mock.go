package etcd

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

// Mock implements IEtcd for tests with hookable behavior and call recording.
type Mock struct {
	HasFunc         func(ctx context.Context, key string) bool
	SetFunc         func(ctx context.Context, key string, args any, timeout time.Duration) error
	SetInFunc       func(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error
	SetMapFunc      func(ctx context.Context, key string, args any, timeout time.Duration) error
	GetFunc         func(ctx context.Context, key string) ([]byte, error)
	GetInFunc       func(ctx context.Context, key string, key2 string) (any, error)
	GetMapFunc      func(ctx context.Context, key string) (map[string]any, error)
	IncrementFunc   func(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)
	IncrementInFunc func(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)
	DecrementFunc   func(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error)
	DecrementInFunc func(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error)
	DeleteFunc      func(ctx context.Context, key string) error
	ExpireFunc      func(ctx context.Context, key string, timeout time.Duration) error

	// In-memory stores for default behavior when hooks are not set
	Data map[string][]byte
	Maps map[string]map[string]any
	Nums map[string]int64

	NamePkg string

	mu    sync.Mutex
	Calls []struct {
		Method string
		Key    string
		Key2   string
	}
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
		m.Data = make(map[string][]byte)
	}
	switch v := args.(type) {
	case string:
		m.Data[key] = []byte(v)
	case []byte:
		m.Data[key] = append([]byte(nil), v...)
	default:
		b, _ := json.Marshal(v)
		m.Data[key] = b
	}
	return nil
}

func (m *Mock) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	m.record("SetIn", key, key2)
	if m.SetInFunc != nil {
		return m.SetInFunc(ctx, key, key2, args, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]any)
	}
	if _, ok := m.Maps[key]; !ok {
		m.Maps[key] = make(map[string]any)
	}
	m.Maps[key][key2] = args
	b, _ := json.Marshal(m.Maps[key])
	if m.Data == nil {
		m.Data = make(map[string][]byte)
	}
	m.Data[key] = b
	return nil
}

func (m *Mock) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	m.record("SetMap", key, "")
	if m.SetMapFunc != nil {
		return m.SetMapFunc(ctx, key, args, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]any)
	}
	if mp, ok := args.(map[string]any); ok {
		m.Maps[key] = mp
	} else {
		b, _ := json.Marshal(args)
		mp := make(map[string]any)
		_ = json.Unmarshal(b, &mp)
		m.Maps[key] = mp
	}
	b, _ := json.Marshal(m.Maps[key])
	if m.Data == nil {
		m.Data = make(map[string][]byte)
	}
	m.Data[key] = b
	return nil
}

func (m *Mock) Get(ctx context.Context, key string) ([]byte, error) {
	m.record("Get", key, "")
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	if m.Data == nil {
		return nil, ErrKeyNotFound
	}
	v, ok := m.Data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return append([]byte(nil), v...), nil
}

func (m *Mock) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	m.record("GetIn", key, key2)
	if m.GetInFunc != nil {
		return m.GetInFunc(ctx, key, key2)
	}
	if m.Maps == nil {
		return nil, ErrKeyNotFound
	}
	mp, ok := m.Maps[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	val, ok := mp[key2]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return val, nil
}

func (m *Mock) GetMap(ctx context.Context, key string) (map[string]any, error) {
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
	out := make(map[string]any, len(mp))
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
		m.Data = make(map[string][]byte)
	}
	m.Data[key] = []byte(strconv.FormatInt(m.Nums[key], 10))
	return m.Nums[key], nil
}

func (m *Mock) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	m.record("IncrementIn", key, key2)
	if m.IncrementInFunc != nil {
		return m.IncrementInFunc(ctx, key, key2, val, timeout)
	}
	if m.Maps == nil {
		m.Maps = make(map[string]map[string]any)
	}
	if _, ok := m.Maps[key]; !ok {
		m.Maps[key] = make(map[string]any)
	}
	curr := int64(0)
	if v, ok := m.Maps[key][key2]; ok {
		switch t := v.(type) {
		case int64:
			curr = t
		case float64:
			curr = int64(t)
		case int:
			curr = int64(t)
		case string:
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				curr = n
			}
		}
	}
	curr += val
	m.Maps[key][key2] = curr
	b, _ := json.Marshal(m.Maps[key])
	if m.Data == nil {
		m.Data = make(map[string][]byte)
	}
	m.Data[key] = b
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

var _ IEtcd = (*Mock)(nil)
