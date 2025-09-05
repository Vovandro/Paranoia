package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

// Mock implements IMySQL for tests. It records calls and returns configured results.
type Mock struct {
	QueryFunc    func(ctx context.Context, query string, args ...interface{}) (SQLRows, error)
	QueryRowFunc func(ctx context.Context, query string, args ...interface{}) (SQLRow, error)
	ExecFunc     func(ctx context.Context, query string, args ...interface{}) error
	BeginTxFunc  func(ctx context.Context) (SQLTx, error)

	NamePkg string

	mu      sync.Mutex
	Queries []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	RowQs []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	Execs []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	Begins int
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
		Ctx   context.Context
		Query string
		Args  []interface{}
	}{Ctx: ctx, Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, query, args...)
	}
	return &MockRows{}, nil
}

func (m *Mock) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	m.mu.Lock()
	m.RowQs = append(m.RowQs, struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}{Ctx: ctx, Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, query, args...)
	}
	return &MockRow{}, nil
}

func (m *Mock) Exec(ctx context.Context, query string, args ...interface{}) error {
	m.mu.Lock()
	m.Execs = append(m.Execs, struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}{Ctx: ctx, Query: query, Args: append([]interface{}{}, args...)})
	m.mu.Unlock()
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, query, args...)
	}
	return nil
}

func (m *Mock) BeginTx(ctx context.Context) (SQLTx, error) {
	m.mu.Lock()
	m.Begins++
	m.mu.Unlock()
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	return &MockTx{}, nil
}

func (m *Mock) GetDb() *sql.DB { return nil }

// MockRow implements SQLRow
type MockRow struct {
	ScanFunc func(dest ...any) error
	Values   []any
}

func (r *MockRow) Scan(dest ...any) error {
	if r.ScanFunc != nil {
		return r.ScanFunc(dest...)
	}
	for i := range dest {
		if i >= len(r.Values) {
			break
		}
		val := r.Values[i]
		if p, ok := dest[i].(*interface{}); ok {
			*p = val
		} else {
			return errors.New("dest is not a pointer")
		}
	}
	return nil
}

// MockRows implements SQLRows
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
	values := r.Values[r.idx-1]
	for i := range dest {
		if i < len(values) {
			if p, ok := dest[i].(*interface{}); ok {
				*p = values[i]
			} else {
				return errors.New("dest is not a pointer")
			}
		}
	}
	return nil
}
func (r *MockRows) Close() error { return nil }

// MockTx implements SQLTx
type MockTx struct {
	QueryFunc    func(ctx context.Context, query string, args ...interface{}) (SQLRows, error)
	QueryRowFunc func(ctx context.Context, query string, args ...interface{}) (SQLRow, error)
	ExecFunc     func(ctx context.Context, query string, args ...interface{}) error
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error
}

func (t *MockTx) Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error) {
	if t.QueryFunc != nil {
		return t.QueryFunc(ctx, query, args...)
	}
	return &MockRows{}, nil
}
func (t *MockTx) QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error) {
	if t.QueryRowFunc != nil {
		return t.QueryRowFunc(ctx, query, args...)
	}
	return &MockRow{}, nil
}
func (t *MockTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	if t.ExecFunc != nil {
		return t.ExecFunc(ctx, query, args...)
	}
	return nil
}
func (t *MockTx) Commit(ctx context.Context) error {
	if t.CommitFunc != nil {
		return t.CommitFunc(ctx)
	}
	return nil
}
func (t *MockTx) Rollback(ctx context.Context) error {
	if t.RollbackFunc != nil {
		return t.RollbackFunc(ctx)
	}
	return nil
}

// Interface assertions
var _ IMySQL = (*Mock)(nil)
var _ SQLRow = (*MockRow)(nil)
var _ SQLRows = (*MockRows)(nil)
var _ SQLTx = (*MockTx)(nil)
