package interfaces

type IServer interface {
	Init(app IService) error
	Stop() error
	String() string
}
