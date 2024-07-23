package storage

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"io"
	"os"
)

type File struct {
	Name string
	app  interfaces.IService
}

func NewFile(name string) interfaces.IStorage {
	return &File{Name: name}
}

func (t *File) Init(app interfaces.IService) error {
	t.app = app
	return nil
}

func (t *File) Stop() error {
	return nil
}

func (t *File) String() string {
	return t.Name
}

func (t *File) Has(name string) bool {
	_, err := os.Stat(name)

	if err != nil {
		return false
	}

	return true
}

func (t *File) Put(name string, data io.Reader) error {
	f, err := os.Create(name)

	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (t *File) StoreFolder(name string) error {
	return os.MkdirAll(name, 0700)
}

func (t *File) Read(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (t *File) Delete(name string) error {
	return os.Remove(name)
}

func (t *File) List(path string) ([]string, error) {
	info, err := os.Stat(path)

	if err != nil {
		return nil, ErrFileNotFound
	}

	if !info.IsDir() {
		return nil, ErrTypeMismatch
	}

	dir, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	res := make([]string, len(dir))

	for i, d := range dir {
		res[i] = d.Name()
	}

	return res, nil
}

func (t *File) IsFolder(name string) (bool, error) {
	info, err := os.Stat(name)

	if err != nil {
		return false, ErrFileNotFound
	}

	return info.IsDir(), nil
}

func (t *File) GetSize(name string) (int64, error) {
	info, err := os.Stat(name)

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.Size(), nil
}

func (t *File) GetModified(name string) (int64, error) {
	info, err := os.Stat(name)

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.ModTime().Unix(), nil
}
