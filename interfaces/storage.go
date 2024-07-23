package interfaces

import "io"

type IStorage interface {
	Init(app IService) error
	Stop() error
	String() string

	Has(name string) bool
	Put(name string, data io.Reader) error
	StoreFolder(name string) error
	Read(name string) (io.ReadCloser, error)
	Delete(name string) error
	List(path string) ([]string, error)

	IsFolder(name string) (bool, error)
	GetSize(name string) (int64, error)
	GetModified(name string) (int64, error)
}
