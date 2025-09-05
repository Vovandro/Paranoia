package as

import (
	"context"
	"sync"

	"github.com/aerospike/aerospike-client-go/v7"
)

// Mock implements IAerospike with hooks and simple in-memory behavior
type Mock struct {
	ExistsFunc     func(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool
	CountFunc      func(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64
	FindOneFunc    func(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (NoSQLRow, error)
	FindFunc       func(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (NoSQLRows, error)
	ExecFunc       func(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (NoSQLRows, error)
	InsertFunc     func(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error)
	DeleteFunc     func(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64
	DeleteManyFunc func(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64
	OperateFunc    func(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error)

	NamePkg string

	mu    sync.Mutex
	Store map[string]map[string]any // key.String() -> bins
}

func (m *Mock) nsk(key *aerospike.Key) string { return key.String() }

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

func (m *Mock) Exists(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, key, policy)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Store == nil {
		return false
	}
	_, ok := m.Store[m.nsk(key)]
	return ok
}

func (m *Mock) Count(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64 {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, key, policy)
	}
	if m.Exists(ctx, key, policy) {
		return 1
	}
	return 0
}

func (m *Mock) FindOne(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (NoSQLRow, error) {
	if m.FindOneFunc != nil {
		return m.FindOneFunc(ctx, key, policy, bins)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Store == nil {
		return &ASRow{row: &aerospike.Record{Bins: map[string]interface{}{}}}, nil
	}
	rec, ok := m.Store[m.nsk(key)]
	if !ok {
		return &ASRow{row: &aerospike.Record{Bins: map[string]interface{}{}}}, nil
	}
	return &ASRow{row: &aerospike.Record{Bins: rec}}, nil
}

func (m *Mock) Find(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (NoSQLRows, error) {
	if m.FindFunc != nil {
		return m.FindFunc(ctx, query, policy)
	}
	return &ASRows{}, nil
}

func (m *Mock) Exec(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (NoSQLRows, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, key, policy, packageName, functionName)
	}
	return nil, nil
}

func (m *Mock) Insert(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error) {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, key, query, policy)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Store == nil {
		m.Store = make(map[string]map[string]any)
	}
	bins := make(map[string]any)
	if bm, ok := query.(*aerospike.BinMap); ok {
		for k, v := range *bm {
			bins[k] = v
		}
	}
	if b, ok := query.(*aerospike.Bin); ok {
		bins[b.Name] = b.Value
	}
	if bs, ok := query.([]*aerospike.Bin); ok {
		for _, b := range bs {
			bins[b.Name] = b.Value
		}
	}
	m.Store[m.nsk(key)] = bins
	return query, nil
}

func (m *Mock) Delete(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64 {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key, policy)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Store == nil {
		return 0
	}
	if _, ok := m.Store[m.nsk(key)]; ok {
		delete(m.Store, m.nsk(key))
		return 1
	}
	return 0
}

func (m *Mock) DeleteMany(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64 {
	if m.DeleteManyFunc != nil {
		return m.DeleteManyFunc(ctx, keys, policy, policyDelete)
	}
	var n int64
	for _, k := range keys {
		n += m.Delete(ctx, k, nil)
	}
	return n
}

func (m *Mock) Operate(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error) {
	if m.OperateFunc != nil {
		return m.OperateFunc(ctx, query)
	}
	return int64(len(query)), nil
}

func (m *Mock) GetDb() *aerospike.Client { return nil }

var _ IAerospike = (*Mock)(nil)
