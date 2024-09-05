package database

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Mock struct {
	Name string
	app  interfaces.IEngine
}

func (t *Mock) Init(app interfaces.IEngine) error {
	t.app = app
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Query(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRows, error) {
	return nil, nil
}

func (t *Mock) QueryRow(ctx context.Context, query string, args ...interface{}) (interfaces.SQLRow, error) {
	return nil, nil
}

func (t *Mock) Exec(ctx context.Context, query string, args ...interface{}) error {
	return nil
}

func (t *Mock) GetDb() interface{} {
	return nil
}
