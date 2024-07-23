# S3

### [Оглавление](./index.md)

## Использование

Инициализация через авто конфиг:

```yaml
engine:
  storage:
    - type: s3
      name: public
      url: s3.ru
      access_key: "key"
      secret_key: "key"
      use_ssl: false
      force_delete: true
      location: ru-RU
      bucket: mybucket
```

Инициализация в коде

```go
app.PushStorage(storage.NewS3("public", storage.S3Config{
	URL: "s3.ru",
	AccessKey: "key",
	SecretKey: "key",
	UseSSL: false,
	ForceDelete: false,
	Location: "ru-RU",
	Bucket: "mybucket",
}))
```

Далее в необходимых местах можно получить:

```go
storage := app.GetStorage("public")
```

## Особенности

Операции над директориями не поддерживаются и возвращается ошибка ErrNotSupported
