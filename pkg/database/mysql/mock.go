package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
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

	if len(r.Values) != len(dest) {
		return errors.New("dest length does not match values length")
	}

	for i := range dest {
		if err := assignValue(dest[i], r.Values[i]); err != nil {
			return err
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

	if len(values) != len(dest) {
		return errors.New("dest length does not match values length")
	}

	for i := range dest {
		if err := assignValue(dest[i], values[i]); err != nil {
			return err
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

// assignValue performs best-effort assignment of val into dest pointer.
// It supports common conversions used in tests, including:
// - *any / *interface{} direct assignment
// - typed pointers (e.g., *uuid.UUID, *string, *int, *int64, *time.Time)
// - pointer-to-pointer targets for nullable fields (e.g., **uuid.UUID)
// - conversions: string/[]byte -> uuid.UUID, []byte -> string, numeric -> numeric
func assignValue(dest any, val any) error {
	// Fast paths for interface pointers kept for backward compatibility
	switch d := dest.(type) {
	case *interface{}:
		*d = val
		return nil
	}

	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Ptr || dv.IsNil() {
		return errors.New("dest is not a pointer")
	}

	// Helper to assign into a non-pointer reflect.Value target
	assignTo := func(target reflect.Value, v any) error {
		if !target.CanSet() {
			return errors.New("cannot set destination value")
		}

		if v == nil {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}

		// If source is a pointer, dereference once
		rv := reflect.ValueOf(v)
		if rv.IsValid() && rv.Kind() == reflect.Ptr && !rv.IsNil() {
			rv = rv.Elem()
			v = rv.Interface()
		}

		// Special-case: uuid.UUID
		if target.Type() == reflect.TypeOf(uuid.UUID{}) {
			switch s := v.(type) {
			case uuid.UUID:
				target.Set(reflect.ValueOf(s))
				return nil
			case string:
				u, err := uuid.Parse(s)
				if err != nil {
					return err
				}
				target.Set(reflect.ValueOf(u))
				return nil
			case []byte:
				u, err := uuid.Parse(string(s))
				if err != nil {
					return err
				}
				target.Set(reflect.ValueOf(u))
				return nil
			}
		}

		// Special-case: time.Time
		if target.Type() == reflect.TypeOf(time.Time{}) {
			switch s := v.(type) {
			case time.Time:
				target.Set(reflect.ValueOf(s))
				return nil
			case string:
				// Try RFC3339; if fails, set zero value
				if ts, err := time.Parse(time.RFC3339, s); err == nil {
					target.Set(reflect.ValueOf(ts))
					return nil
				}
				target.Set(reflect.Zero(target.Type()))
				return nil
			}
		}

		// Strings: accept []byte and fmt.Stringer
		if target.Kind() == reflect.String {
			switch s := v.(type) {
			case string:
				target.SetString(s)
				return nil
			case []byte:
				target.SetString(string(s))
				return nil
			default:
				target.SetString(fmt.Sprint(v))
				return nil
			}
		}

		// Numeric conversions
		if isNumericKind(target.Kind()) {
			if rv.IsValid() && isNumericKind(rv.Kind()) {
				if rv.Type().ConvertibleTo(target.Type()) {
					target.Set(rv.Convert(target.Type()))
					return nil
				}
			}
		}

		// Direct assign or convertible types
		if rv.IsValid() {
			if rv.Type().AssignableTo(target.Type()) {
				target.Set(rv)
				return nil
			}
			if rv.Type().ConvertibleTo(target.Type()) {
				target.Set(rv.Convert(target.Type()))
				return nil
			}
		}

		return errors.New("unsupported scan assignment")
	}

	// Handle pointer-to-pointer for nullable columns
	// Example: **uuid.UUID, **string, **time.Time
	elem := dv.Elem()
	if elem.Kind() == reflect.Ptr {
		if val == nil {
			// Leave as nil
			elem.Set(reflect.Zero(elem.Type()))
			return nil
		}
		// Allocate new inner value and assign
		inner := reflect.New(elem.Type().Elem())
		if err := assignTo(inner.Elem(), val); err != nil {
			return err
		}
		elem.Set(inner)
		return nil
	}

	// Simple pointer target
	return assignTo(elem, val)
}

func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isNumericValue(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	return isNumericKind(v.Kind())
}
