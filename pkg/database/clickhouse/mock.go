package clickhouse

import (
	"context"
	"sync"
)

// Mock implements IClickHouse for tests with hooks and simple in-memory rows
type Mock struct {
	QueryFunc    func(ctx context.Context, query string, args ...interface{}) (SQLRows, error)
	QueryRowFunc func(ctx context.Context, query string, args ...interface{}) (SQLRow, error)
	ExecFunc     func(ctx context.Context, query string, args ...interface{}) error

	NamePkg string

	mu      sync.Mutex
	Queries []struct {
		Query string
		Args  []interface{}
	}
	RowQs []struct {
		Query string
		Args  []interface{}
	}
	Execs []struct {
		Query string
		Args  []interface{}
	}
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
	return "database"
}

func (m *Mock) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	m.mu.Lock()
	m.Queries = append(m.Queries, struct {
		Query string
		Args  []interface{}
	}{Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, query, args...)
	}
	return &MockRows{}, nil
}
func (m *Mock) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	m.mu.Lock()
	m.RowQs = append(m.RowQs, struct {
		Query string
		Args  []interface{}
	}{Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, query, args...)
	}
	return &MockRow{}, nil
}
func (m *Mock) Exec(ctx context.Context, query string, args ...interface{}) error {
	m.mu.Lock()
	m.Execs = append(m.Execs, struct {
		Query string
		Args  []interface{}
	}{Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, query, args...)
	}
	return nil
}
func (m *Mock) GetDb() interface{} { return nil }

type MockRow struct {
	ScanFunc func(dest ...any) error
	Values   []any
}

func (r *MockRow) Scan(dest ...any) error {
	if r.ScanFunc != nil {
		return r.ScanFunc(dest...)
	}
	for i := range dest {
		if i < len(r.Values) {
			if p, ok := dest[i].(*interface{}); ok {
				*p = r.Values[i]
			}
		}
	}
	return nil
}

type MockRows struct {
	Values   [][]any
	ScanFunc func(idx int, dest ...any) error
	idx      int
}

func (r *MockRows) Next() bool {
	if r.idx < len(r.Values) {
		r.idx++
		return true
	}
	return false
}
func (r *MockRows) Scan(dest ...any) error {
	if r.idx == 0 || r.idx > len(r.Values) {
		return nil
	}
	if r.ScanFunc != nil {
		return r.ScanFunc(r.idx-1, dest...)
	}
	vals := r.Values[r.idx-1]
	for i := range dest {
		if i < len(vals) {
			if p, ok := dest[i].(*interface{}); ok {
				*p = vals[i]
			}
		}
	}
	return nil
}
func (r *MockRows) Close() error { return nil }

var _ IClickHouse = (*Mock)(nil)
var _ SQLRow = (*MockRow)(nil)
var _ SQLRows = (*MockRows)(nil)
