package interfaces

type IModules interface {
	Init(app IEngine, cfg map[string]interface{}) error
	Stop() error
	Name() string
	Type() string
}
