package interfaces

type IController interface {
	Init(app IEngine) error
	Stop() error
	String() string
}
