package service

import (
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
)

type Mock struct {
	NamePkg string
}

func (t *Mock) Init(_ interfaces.IEngine, _ map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Type() string {
	return interfaces.ModuleService
}
