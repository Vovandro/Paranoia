package interfaces

type IDatabase interface {
	Init(app IService) error
	Stop() error
	String() string

	Exists(key interface{}) bool
	Exec(args ...interface{}) error
	Query(args ...interface{}) ([]interface{}, error)
	Fetch(args ...interface{}) (interface{}, error)
	Delete(args ...interface{}) error
}
