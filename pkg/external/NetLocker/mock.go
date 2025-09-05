package NetLocker

import (
	"context"
	"sync"
)

// Mock implements INetLocker with hooks and in-memory lock table
type Mock struct {
	LockFunc   func(ctx context.Context, key string, timeLock int64, uniqueId *string) (bool, error)
	UnlockFunc func(ctx context.Context, key string, uniqueId *string) bool

	// In-memory locks: key -> uniqueId string
	Locks   map[string]string
	NamePkg string

	mu    sync.Mutex
	Calls []struct{ Method, Key string }
}

func (m *Mock) record(method, key string) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct{ Method, Key string }{Method: method, Key: key})
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
	return "external"
}

func (m *Mock) Lock(ctx context.Context, key string, timeLock int64, uniqueId *string) (bool, error) {
	m.record("Lock", key)
	if m.LockFunc != nil {
		return m.LockFunc(ctx, key, timeLock, uniqueId)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Locks == nil {
		m.Locks = make(map[string]string)
	}
	if _, exists := m.Locks[key]; exists {
		return false, nil
	}
	if uniqueId == nil {
		return false, nil
	}
	m.Locks[key] = *uniqueId
	return true, nil
}

func (m *Mock) Unlock(ctx context.Context, key string, uniqueId *string) bool {
	m.record("Unlock", key)
	if m.UnlockFunc != nil {
		return m.UnlockFunc(ctx, key, uniqueId)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Locks == nil {
		return false
	}
	if uniqueId == nil {
		delete(m.Locks, key)
		return true
	}
	if cur, ok := m.Locks[key]; ok && cur == *uniqueId {
		delete(m.Locks, key)
		return true
	}
	return false
}

var _ INetLocker = (*Mock)(nil)
