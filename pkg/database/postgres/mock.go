package postgres

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
)

// Mock implements IPostgres for tests. It records calls and returns configured results.
type Mock struct {
	NamePkg string

	// Configurable responses
	QueryFunc    func(ctx context.Context, query string, args ...interface{}) (SQLRows, error)
	QueryRowFunc func(ctx context.Context, query string, args ...interface{}) (SQLRow, error)
	ExecFunc     func(ctx context.Context, query string, args ...interface{}) error
	BeginTxFunc  func(ctx context.Context) (SQLTx, error)

	// Call recording
	mu      sync.Mutex
	Queries []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	QueryRows []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	Execs []struct {
		Ctx   context.Context
		Query string
		Args  []interface{}
	}
	BeginTxCalls int
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
	m.QueryRows = append(m.QueryRows, struct {
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
	m.BeginTxCalls++
	m.mu.Unlock()
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	return &MockTx{}, nil
}

func (m *Mock) GetDb() *pgx.Conn {
	return nil
}

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
	// Preloaded rows to iterate through
	Values [][]any

	// If set, called by Scan to customize behavior per current index
	ScanFunc func(idx int, dest ...any) error

	idx int
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
			switch d := dest[i].(type) {
			case *any:
				*d = values[i]
			default:
				// Best-effort plain assignment when possible
				if p, ok := dest[i].(*interface{}); ok {
					*p = values[i]
				} else {
					return errors.New("dest is not a pointer")
				}
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

	mu      sync.Mutex
	closed  bool
	actions []string
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
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	t.actions = append(t.actions, "commit")
	if t.CommitFunc != nil {
		return t.CommitFunc(ctx)
	}
	return nil
}

func (t *MockTx) Rollback(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	t.actions = append(t.actions, "rollback")
	if t.RollbackFunc != nil {
		return t.RollbackFunc(ctx)
	}
	return nil
}

// Interface assertions ensure mock implements required contracts
var _ IPostgres = (*Mock)(nil)
var _ SQLRow = (*MockRow)(nil)
var _ SQLRows = (*MockRows)(nil)
var _ SQLTx = (*MockTx)(nil)
