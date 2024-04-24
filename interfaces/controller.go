package interfaces

type IController interface {
	Init(app IService) error
	Stop() error
	String() string
}
