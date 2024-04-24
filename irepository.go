package Paranoia

type IRepository interface {
	Init(app *Service) error
	Stop() error
	String() string
}
