package Paranoia

type IModules interface {
	Init(app *Service) error
	Stop() error
	String() string
}
