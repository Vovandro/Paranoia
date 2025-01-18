package file

import (
	"io"
)

// IFile defines the interface for file storage operations
type IFile interface {
	// Has checks if a file with the given name exists
	Has(name string) bool

	// Put stores the data from the reader into a file with the given name
	Put(name string, data io.Reader) error

	// StoreFolder creates a folder with the given name
	StoreFolder(name string) error

	// Read returns a reader for the file with the given name
	Read(name string) (io.ReadCloser, error)

	// Delete removes the file with the given name
	Delete(name string) error

	// List returns a list of files in the specified folder
	List(folder string) ([]string, error)

	// IsFolder checks if the given name is a folder
	IsFolder(name string) (bool, error)

	// GetSize returns the size of the file with the given name
	GetSize(name string) (int64, error)

	// GetModified returns the last modified time of the file with the given name
	GetModified(name string) (int64, error)
}
