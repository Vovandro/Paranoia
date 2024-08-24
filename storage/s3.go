package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"io"
)

type S3 struct {
	Name   string
	app    interfaces.IEngine
	config S3Config
	client *minio.Client
}

type S3Config struct {
	URL         string `yaml:"url"`
	AccessKey   string `yaml:"access_key"`
	SecretKey   string `yaml:"secret_key"`
	UseSSL      bool   `yaml:"use_ssl"`
	ForceDelete bool   `yaml:"force_delete"`
	Location    string `yaml:"location"`
	Bucket      string `yaml:"bucket"`
}

func NewS3(name string, cfg S3Config) interfaces.IStorage {
	return &S3{Name: name, config: cfg}
}

func (t *S3) Init(app interfaces.IEngine) error {
	t.app = app
	var err error

	t.client, err = minio.New(t.config.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(t.config.AccessKey, t.config.SecretKey, ""),
		Secure: t.config.UseSSL,
	})

	if err != nil {
		return err
	}

	exists, err := t.client.BucketExists(context.Background(), t.config.Bucket)

	if err != nil {
		return err
	}

	if !exists {
		err = t.client.MakeBucket(
			context.Background(),
			t.config.Bucket,
			minio.MakeBucketOptions{Region: t.config.Location},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *S3) Stop() error {
	return nil
}

func (t *S3) String() string {
	return t.Name
}

func (t *S3) Has(name string) bool {

	object, err := t.client.StatObject(context.Background(), t.config.Bucket, name, minio.StatObjectOptions{})

	if err != nil {
		return false
	}

	return !object.IsDeleteMarker
}

func (t *S3) Put(name string, data io.Reader) error {
	_, err := t.client.PutObject(
		context.Background(),
		t.config.Bucket,
		name,
		data,
		-1,
		minio.PutObjectOptions{ContentType: "application/text"},
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *S3) StoreFolder(name string) error {
	return ErrNotSupported
}

func (t *S3) Read(name string) (io.ReadCloser, error) {
	return t.client.GetObject(context.Background(), t.config.Bucket, name, minio.GetObjectOptions{})
}

func (t *S3) Delete(name string) error {
	return t.client.RemoveObject(context.Background(), t.config.Bucket, name, minio.RemoveObjectOptions{
		ForceDelete: t.config.ForceDelete,
	})
}

func (t *S3) List(path string) ([]string, error) {
	objects := t.client.ListObjects(context.Background(), t.config.Bucket, minio.ListObjectsOptions{
		Prefix: path,
	})

	res := make([]string, 0, 10)

	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}

		res = append(res, object.Key)
	}

	return res, nil
}

func (t *S3) IsFolder(name string) (bool, error) {
	return false, ErrNotSupported
}

func (t *S3) GetSize(name string) (int64, error) {
	info, err := t.client.StatObject(context.Background(), t.config.Bucket, name, minio.StatObjectOptions{})

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.Size, nil
}

func (t *S3) GetModified(name string) (int64, error) {
	info, err := t.client.StatObject(context.Background(), t.config.Bucket, name, minio.StatObjectOptions{})

	if err != nil {
		return 0, ErrFileNotFound
	}

	return info.LastModified.Unix(), nil
}
