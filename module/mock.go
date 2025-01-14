package module

import "gitlab.com/devpro_studio/Paranoia/interfaces"

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

func (t *Mock) Type() string {
	return string(interfaces.ModuleModule)
}
