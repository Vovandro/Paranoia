# File

### [Оглавление](./index.md)

## Использование

Инициализация через авто конфиг:

```yaml
engine:
  storage:
    - type: file
      name: public
      folder: tmp
```

Инициализация в коде

```go
app.PushStorage(storage.NewFile("public"), storage.FileConfig{Folder: "tmp"})
```

Далее в необходимых местах можно получить:

```go
storage := app.GetStorage("public")
```

## Особенности

Папки и файлы всегда создаются с правами по умолчанию

### Далее [S3](./storage-s3.md)