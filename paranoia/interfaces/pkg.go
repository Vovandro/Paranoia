package interfaces

type IPkg interface {
	Init(cfg map[string]interface{}) error
	Stop() error
	Name() string
	Type() string
}
