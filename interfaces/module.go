package interfaces

type IModules interface {
	Init(app IService) error
	Stop() error
	String() string
}
