package interfaces

type IService interface {
	Init(app IEngine) error
	Stop() error
	String() string
}
