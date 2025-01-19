package task

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
)

type Mock struct {
	NamePkg string
}

func (t *Mock) Init(_ interfaces.IEngine) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Start() []interfaces.ITaskRunConfiguration {
	return []interfaces.ITaskRunConfiguration{}
}

func (t *Mock) Invoke(ctx context.Context, data map[string]interface{}) {}
