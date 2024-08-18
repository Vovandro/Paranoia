package interfaces

type ITask interface {
	Init(app IService) error
	Stop() error
	String() string
	Start() chan interface{}
	Invoke()
}
