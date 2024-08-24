package interfaces

type IModules interface {
	Init(app IEngine) error
	Stop() error
	String() string
}
