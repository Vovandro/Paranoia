package noSql

import (
	"context"
	"github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type MockAerospike struct {
	Name string
}

func (t *MockAerospike) Exists(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool {
	return false
}

func (t *MockAerospike) Count(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64 {
	return 0
}

func (t *MockAerospike) FindOne(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (interfaces.NoSQLRow, error) {
	return &ASRow{nil}, nil
}

func (t *MockAerospike) Find(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (interfaces.NoSQLRows, error) {
	return &ASRows{rows: nil}, nil
}

func (t *MockAerospike) Exec(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (interfaces.NoSQLRows, error) {
	return nil, nil
}

// Insert query is *aerospike.Bin or []*aerospike.Bin or *aerospike.BinMap
func (t *MockAerospike) Insert(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error) {
	return query, nil
}

func (t *MockAerospike) Delete(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64 {
	return 0
}

func (t *MockAerospike) DeleteMany(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64 {
	return 0
}

func (t *MockAerospike) Operate(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error) {
	return 0, nil
}

func (t *MockAerospike) GetDb() *aerospike.Client {
	return nil
}
