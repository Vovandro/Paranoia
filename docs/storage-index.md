# Файловая область

### [Оглавление](./index.md)


## - [File](./storage-file.md)
## - [S3](./storage-s3.md)

Общие методы:

- `Has(name string) bool` - Проверка существования файла или директории
- `Put(name string, data io.Reader) error` - Создание файла и запись данных
- `StoreFolder(name string) error` - Создание директории
- `Read(name string) (io.ReadCloser, error)` - Получение файла
- `Delete(name string) error` - Удаление файла или директории
- `List(path string) ([]string, error)` - Получение списка файлов и директорий
- `IsFolder(name string) (bool, error)` - Проверка, что путь указывает на директорию
- `GetSize(name string) (int64, error)` - Получение размера файла
- `GetModified(name string) (int64, error)` - Получение времени последней модификации файла

Общие ошибки:

- `storage.ErrFileNotFound` - Файла или директории не существует
- `storage.ErrTypeMismatch` - К данному типу действие не применимо
- `storage.ErrNotSupported` - Данное действие не применимо

Пример:

```go
storage.Put("/public/test.txt", "Hello")
```

### Далее [File](./storage-file.md)