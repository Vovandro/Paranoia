package goServer

type IModules interface {
	Init(app *Service) error
	Stop() error
	String() string
}
