package service

import "gitlab.com/devpro_studio/Paranoia/interfaces"

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
