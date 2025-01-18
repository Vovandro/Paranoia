package s3

import (
	"io"
)

// IS3 defines the interface for S3 storage operations
type IS3 interface {
	// Has checks if an object with the given name exists in the S3 bucket
	Has(name string) bool

	// Put stores the data from the reader into an object with the given name in the S3 bucket
	Put(name string, data io.Reader) error

	// StoreFolder creates a folder with the given name in the S3 bucket
	StoreFolder(name string) error

	// Read returns a reader for the object with the given name in the S3 bucket
	Read(name string) (io.ReadCloser, error)

	// Delete removes the object with the given name from the S3 bucket
	Delete(name string) error

	// List returns a list of objects in the specified path of the S3 bucket
	List(path string) ([]string, error)

	// IsFolder checks if the given name is a folder in the S3 bucket
	IsFolder(name string) (bool, error)

	// GetSize returns the size of the object with the given name in the S3 bucket
	GetSize(name string) (int64, error)

	// GetModified returns the last modified time of the object with the given name in the S3 bucket
	GetModified(name string) (int64, error)
}
