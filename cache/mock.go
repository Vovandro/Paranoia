package cache

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Mock struct {
	Name string
}

func (t *Mock) Init(app interfaces.IService) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) String() string {
	return t.Name
}

func (t *Mock) Has(key string) bool {
	return false
}

func (t *Mock) Set(key string, args any, timeout time.Duration) error {
	return nil
}

func (t *Mock) Get(key string) (any, error) {
	return nil, fmt.Errorf("key not found")
}

func (t *Mock) Delete(key string) error {
	return nil
}
