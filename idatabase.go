package Paranoia

type IDatabase interface {
	Init(app *Service) error
	Stop() error
	String() string

	Exists(key interface{}) bool
	Exec(args ...interface{}) error
	Query(args ...interface{}) ([]interface{}, error)
	Fetch(args ...interface{}) (interface{}, error)
	Delete(args ...interface{}) error
}
