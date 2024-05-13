package interfaces

type IStorage interface {
	Init(app IService) error
	Stop() error
	String() string

	Has(name string) bool
	Put(name string, data []byte) error
	Read(name string) ([]byte, error)
	Delete(name string) error
	List() []string
}
