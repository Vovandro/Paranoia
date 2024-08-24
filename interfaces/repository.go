package interfaces

type IRepository interface {
	Init(app IEngine) error
	Stop() error
	String() string
}
