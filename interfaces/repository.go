package interfaces

type IRepository interface {
	Init(app IService) error
	Stop() error
	String() string
}
