package interfaces

type IStorage interface {
	Init(app IService) error
	Stop() error
	String() string

	Has(name string) bool
	Put(name string, data []byte) error
	StoreFolder(name string) error
	Read(name string) ([]byte, error)
	Delete(name string) error
	List(path string) ([]string, error)

	IsFolder(name string) (bool, error)
	GetSize(name string) (int64, error)
	GetModified(name string) (int64, error)
}
