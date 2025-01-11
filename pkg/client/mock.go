package client

type Mock struct {
	name string
}

func (t *Mock) Init(_ map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.name
}
