package storage

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"io"
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

func (t *Mock) Has(name string) bool {
	return false
}

func (t *Mock) Put(name string, data io.Reader) error {
	return nil
}

func (t *Mock) StoreFolder(name string) error {
	return nil
}

func (t *Mock) Read(name string) (io.ReadCloser, error) {
	return nil, nil
}

func (t *Mock) Delete(name string) error {
	return nil
}

func (t *Mock) List(path string) ([]string, error) {
	return nil, nil
}

func (t *Mock) IsFolder(name string) (bool, error) {
	return false, nil
}

func (t *Mock) GetSize(name string) (int64, error) {
	return 0, nil
}

func (t *Mock) GetModified(name string) (int64, error) {
	return 0, nil
}
