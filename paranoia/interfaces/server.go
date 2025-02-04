package interfaces

type IServer interface {
	Init(cfg map[string]interface{}) error
	Start() error
	Stop() error
	Name() string
	Type() string
}
