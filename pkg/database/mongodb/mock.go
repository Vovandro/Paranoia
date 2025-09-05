package mongodb

import (
	"context"
	"sync"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mock implements IMongoDB for tests. It records calls and returns configured results.
type Mock struct {
	// Configurable responses
	ExistsFunc           func(ctx context.Context, collection string, query interface{}) bool
	CountFunc            func(ctx context.Context, collection string, query interface{}, opt *options.CountOptions) int64
	FindOneFunc          func(ctx context.Context, collection string, query interface{}, opt *options.FindOneOptions) (NoSQLRow, error)
	FindOneAndUpdateFunc func(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) (NoSQLRow, error)
	FindFunc             func(ctx context.Context, collection string, query interface{}, opt *options.FindOptions) (NoSQLRows, error)
	ExecFunc             func(ctx context.Context, collection string, query interface{}, opt *options.AggregateOptions) (NoSQLRows, error)
	InsertFunc           func(ctx context.Context, collection string, query interface{}, opt *options.InsertOneOptions) (interface{}, error)
	UpdateFunc           func(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.UpdateOptions) error
	DeleteFunc           func(ctx context.Context, collection string, query interface{}, opt *options.DeleteOptions) int64
	BatchFunc            func(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error)

	// Call recording
	NamePkg string

	mu    sync.Mutex
	Calls []struct {
		Method, Collection string
		Query              interface{}
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

func (m *Mock) record(method, collection string, query interface{}) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct {
		Method, Collection string
		Query              interface{}
	}{Method: method, Collection: collection, Query: query})
	m.mu.Unlock()
}

func (m *Mock) Exists(ctx context.Context, collection string, query interface{}) bool {
	m.record("Exists", collection, query)
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, collection, query)
	}
	return false
}

func (m *Mock) Count(ctx context.Context, collection string, query interface{}, opt *options.CountOptions) int64 {
	m.record("Count", collection, query)
	if m.CountFunc != nil {
		return m.CountFunc(ctx, collection, query, opt)
	}
	return 0
}

func (m *Mock) FindOne(ctx context.Context, collection string, query interface{}, opt *options.FindOneOptions) (NoSQLRow, error) {
	m.record("FindOne", collection, query)
	if m.FindOneFunc != nil {
		return m.FindOneFunc(ctx, collection, query, opt)
	}
	return &MockRow{}, nil
}

func (m *Mock) FindOneAndUpdate(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) (NoSQLRow, error) {
	m.record("FindOneAndUpdate", collection, query)
	if m.FindOneAndUpdateFunc != nil {
		return m.FindOneAndUpdateFunc(ctx, collection, query, update, opt)
	}
	return &MockRow{}, nil
}

func (m *Mock) Find(ctx context.Context, collection string, query interface{}, opt *options.FindOptions) (NoSQLRows, error) {
	m.record("Find", collection, query)
	if m.FindFunc != nil {
		return m.FindFunc(ctx, collection, query, opt)
	}
	return &MockRows{}, nil
}

func (m *Mock) Exec(ctx context.Context, collection string, query interface{}, opt *options.AggregateOptions) (NoSQLRows, error) {
	m.record("Exec", collection, query)
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, collection, query, opt)
	}
	return &MockRows{}, nil
}

func (m *Mock) Insert(ctx context.Context, collection string, query interface{}, opt *options.InsertOneOptions) (interface{}, error) {
	m.record("Insert", collection, query)
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, collection, query, opt)
	}
	return nil, nil
}

func (m *Mock) Update(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.UpdateOptions) error {
	m.record("Update", collection, query)
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, collection, query, update, opt)
	}
	return nil
}

func (m *Mock) Delete(ctx context.Context, collection string, query interface{}, opt *options.DeleteOptions) int64 {
	m.record("Delete", collection, query)
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, collection, query, opt)
	}
	return 0
}

func (m *Mock) Batch(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error) {
	m.record("Batch", collection, query)
	if m.BatchFunc != nil {
		return m.BatchFunc(ctx, collection, query, opt)
	}
	return 0, nil
}

func (m *Mock) GetDb() *mongo.Database { return nil }

// MockRow implements NoSQLRow
type MockRow struct {
	ScanFunc func(dest any) error
}

func (r *MockRow) Scan(dest any) error {
	if r.ScanFunc != nil {
		return r.ScanFunc(dest)
	}
	return nil
}

// MockRows implements NoSQLRows
type MockRows struct {
	Values   []any
	ScanFunc func(idx int, dest any) error
	idx      int
}

func (r *MockRows) Next() bool {
	if r.idx < len(r.Values) {
		r.idx++
		return true
	}
	return false
}

func (r *MockRows) Scan(dest any) error {
	if r.idx == 0 || r.idx > len(r.Values) {
		return nil
	}
	if r.ScanFunc != nil {
		return r.ScanFunc(r.idx-1, dest)
	}

	return decode.Decode(r.Values[r.idx-1], dest, "bson", decode.DecoderStrongFoundDst)
}

func (r *MockRows) Close() error { return nil }

// Interface assertions
var _ IMongoDB = (*Mock)(nil)
var _ NoSQLRow = (*MockRow)(nil)
var _ NoSQLRows = (*MockRows)(nil)
