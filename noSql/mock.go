package noSql

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Mock struct {
	Name string
}

func (t *Mock) Init(_ interfaces.IEngine) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Exists(ctx context.Context, key interface{}, query interface{}, args ...interface{}) bool {
	return false
}

func (t *Mock) Count(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	return 0
}

func (t *Mock) FindOne(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRow, error) {
	return nil, nil
}

func (t *Mock) Find(ctx context.Context, _ interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	return nil, nil
}

func (t *Mock) Exec(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	return nil, nil
}

func (t *Mock) Insert(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (t *Mock) Update(ctx context.Context, key interface{}, query interface{}, args ...interface{}) error {
	return nil

}

func (t *Mock) Delete(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	return 0
}

func (t *Mock) Batch(ctx context.Context, key interface{}, query interface{}, typeOp string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (t *Mock) GetDb() interface{} {
	return nil
}
