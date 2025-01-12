package file

import (
	"errors"
	"gitlab.com/devpro_studio/go_utils/decode"
	"io"
	"os"
	"path"
)

type File struct {
	name string

	config Config
}

type Config struct {
	Folder string `yaml:"folder"`
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

func (t *File) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Folder == "" {
		return errors.New("folder is required")
	}

	os.MkdirAll(t.config.Folder, 0755)

	return nil
}

func (t *File) Stop() error {
	return nil
}

func (t *File) Name() string {
	return t.name
}

func (t *File) Type() string {
	return "storage"
}

func (t *File) Has(name string) bool {
	_, err := os.Stat(path.Join(t.config.Folder, name))

	if err != nil {
		return false
	}

	return true
}

func (t *File) Put(name string, data io.Reader) error {
	f, err := os.Create(path.Join(t.config.Folder, name))

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
	return os.MkdirAll(path.Join(t.config.Folder, name), 0755)
}

func (t *File) Read(name string) (io.ReadCloser, error) {
	f, err := os.Open(path.Join(t.config.Folder, name))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (t *File) Delete(name string) error {
	return os.Remove(path.Join(t.config.Folder, name))
}

func (t *File) List(folder string) ([]string, error) {
	info, err := os.Stat(path.Join(t.config.Folder, folder))

	if err != nil {
		return nil, ErrFileNotFound
	}

	if !info.IsDir() {
		return nil, ErrTypeMismatch
	}

	dir, err := os.ReadDir(path.Join(t.config.Folder, folder))

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
	info, err := os.Stat(path.Join(t.config.Folder, name))

	if err != nil {
		return false, ErrFileNotFound
	}

	return info.IsDir(), nil
}

func (t *File) GetSize(name string) (int64, error) {
	info, err := os.Stat(path.Join(t.config.Folder, name))

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.Size(), nil
}

func (t *File) GetModified(name string) (int64, error) {
	info, err := os.Stat(path.Join(t.config.Folder, name))

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.ModTime().Unix(), nil
}
