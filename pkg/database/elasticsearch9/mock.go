package elasticsearch9

import (
	"context"
	"sync"
)

// Mock implements IElasticSearch with hooks and simple in-memory docs
type Mock struct {
	IndexFunc         func(ctx context.Context, index string, id string, document interface{}, refresh bool) (string, error)
	GetFunc           func(ctx context.Context, index string, id string) (NoSQLRow, error)
	SearchFunc        func(ctx context.Context, index []string, query map[string]any, from, size int) (NoSQLRows, error)
	SearchSourceFunc  func(ctx context.Context, index []string, query map[string]any, from, size int, include, exclude []string) (NoSQLRows, error)
	DeleteFunc        func(ctx context.Context, index string, id string, refresh bool) error
	DeleteByQueryFunc func(ctx context.Context, index []string, query map[string]any, refresh bool) error
	UpdateFunc        func(ctx context.Context, index string, id string, doc interface{}, refresh bool) error
	BulkIndexFunc     func(ctx context.Context, index string, items []BulkItem, refresh bool) (BulkIndexResult, error)

	NamePkg string

	mu   sync.Mutex
	Docs map[string]map[string]any // index -> id -> doc(any)
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

func (m *Mock) Index(ctx context.Context, index string, id string, document interface{}, refresh bool) (string, error) {
	if m.IndexFunc != nil {
		return m.IndexFunc(ctx, index, id, document, refresh)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Docs == nil {
		m.Docs = make(map[string]map[string]any)
	}
	if _, ok := m.Docs[index]; !ok {
		m.Docs[index] = make(map[string]any)
	}
	m.Docs[index][id] = document
	return id, nil
}

func (m *Mock) Get(ctx context.Context, index string, id string) (NoSQLRow, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, index, id)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Docs != nil {
		if mp, ok := m.Docs[index]; ok {
			if d, ok2 := mp[id]; ok2 {
				return &MockRow{Doc: d}, nil
			}
		}
	}
	return &MockRow{Doc: nil}, nil
}

func (m *Mock) Search(ctx context.Context, index []string, query map[string]any, from, size int) (NoSQLRows, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, index, query, from, size)
	}
	return &MockRows{}, nil
}

func (m *Mock) SearchSource(ctx context.Context, index []string, query map[string]any, from, size int, include, exclude []string) (NoSQLRows, error) {
	if m.SearchSourceFunc != nil {
		return m.SearchSourceFunc(ctx, index, query, from, size, include, exclude)
	}
	return &MockRows{}, nil
}

func (m *Mock) Delete(ctx context.Context, index string, id string, refresh bool) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, index, id, refresh)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Docs != nil {
		if mp, ok := m.Docs[index]; ok {
			delete(mp, id)
		}
	}
	return nil
}

func (m *Mock) DeleteByQuery(ctx context.Context, index []string, query map[string]any, refresh bool) error {
	if m.DeleteByQueryFunc != nil {
		return m.DeleteByQueryFunc(ctx, index, query, refresh)
	}
	return nil
}

func (m *Mock) Update(ctx context.Context, index string, id string, doc interface{}, refresh bool) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, index, id, doc, refresh)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Docs == nil {
		m.Docs = make(map[string]map[string]any)
	}
	if _, ok := m.Docs[index]; !ok {
		m.Docs[index] = make(map[string]any)
	}
	m.Docs[index][id] = doc
	return nil
}

func (m *Mock) BulkIndex(ctx context.Context, index string, items []BulkItem, refresh bool) (BulkIndexResult, error) {
	if m.BulkIndexFunc != nil {
		return m.BulkIndexFunc(ctx, index, items, refresh)
	}
	res := BulkIndexResult{IDs: make([]string, 0, len(items))}
	for _, it := range items {
		if _, err := m.Index(ctx, index, it.ID, it.Document, refresh); err == nil {
			res.IDs = append(res.IDs, it.ID)
		}
	}
	return res, nil
}

func (m *Mock) GetClient() interface{} { return nil }

type MockRow struct{ Doc any }

func (r *MockRow) Scan(dest any) error {
	if mp, ok := dest.(*any); ok {
		*mp = r.Doc
		return nil
	}
	return nil
}

type MockRows struct {
	Items []any
	idx   int
}

func (r *MockRows) Next() bool {
	if r.idx < len(r.Items) {
		r.idx++
		return true
	}
	return false
}

func (r *MockRows) Scan(dest any) error {
	if r.idx == 0 || r.idx > len(r.Items) {
		return nil
	}
	if p, ok := dest.(*any); ok {
		*p = r.Items[r.idx-1]
	}
	return nil
}

func (r *MockRows) Close() error { return nil }

var _ IElasticSearch = (*Mock)(nil)
var _ NoSQLRow = (*MockRow)(nil)
var _ NoSQLRows = (*MockRows)(nil)
