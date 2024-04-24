package Paranoia

type IController interface {
	Init(app *Service) error
	Stop() error
	String() string
}
