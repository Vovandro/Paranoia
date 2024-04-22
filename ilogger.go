package goServer

type ILogger interface {
	Init(app *Service) error
	Stop() error
	SetLevel(level string)
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(err error)
	Fatal(err error)
	Panic(err error)
}
