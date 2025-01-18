package task

import (
	"context"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
)

type Mock struct {
	NamePkg string
}

func (t *Mock) Init(_ interfaces2.IEngine) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Start() []interfaces2.ITaskRunConfiguration {
	return []interfaces2.ITaskRunConfiguration{}
}

func (t *Mock) Invoke(ctx context.Context, data map[string]interface{}) {}
