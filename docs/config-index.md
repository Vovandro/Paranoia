# Конфигурация системы

### [Оглавление](./index.md)

Система конфигурации устанавливается при инициализации фреймворка и не меняется в пределах проекта

Фреймворк штатно поддерживает системы конфигурирования:
- Mock конфигурация - конфигурация прописана в исходном коде приложения
- .env файлы - минимальная поддержка env структуры, без импорта других файлов. Для получения карт и слайсов используется декодирование строки по токенам указанных в конфигурации.
- Авто конфигурация из yaml файла - поддержка плоской пользовательской конфигурации и настройки всех встроенных модулей фреймворка

# Mock конфигурация

Используется в качестве заглушки в коде или тестах

```go
s := Paranoia.New("app", config.NewMock(map[string]string{
	"key": "val",
}), nil)
```

# .env файлы

Получение данных происходит из файла или переменных операционной системы. Возможно указать любое название файла

```go
s := Paranoia.New("app", config.NewEnv(
	config.EnvConfig{FName: ".env.local"},
), nil)
```

# Авто конфигурация

Позволяет инициализировать приложение в зависимости от окружения, загрузка происходит из yaml файла. 

Поддерживается 2 корневых тега: engine - для конфигурации фреймворка и cfg для пользовательской конфигурации.

```go
s := Paranoia.New("app", config.NewAuto(
config.AutoConfig{FName: "local.yaml"},
), nil)
```

Пример файла конфигурации:

```yaml
engine:
  metrics:
    - type: std
      name: example_app
      interval: 30s
  cache:
    - type: memory
      name: cache
      time_clear: 10m

cfg:
  logLevel: WARNING
  key: val
```

`logLevel` из пользовательской конфигурации автоматически применяется при загрузке конфигурации. В блоке `cfg` доступно только плоская структура - без объектов и без массивов

В блоке `engine` следующем уровнем указывается название модуля:
- metrics
- cache
- database
- nosql
- client
- server
- middleware
- storage

Следующем уровнем в качестве массива указываются модули которые будут добавлены во фреймворк. Допускаются повторения модулей, но под уникальным именем в пределах типа пакета.
В описании каждого модуля необходимо обязательно указать тип модуля, его имя и остальные настройки модуля.

# Получение конфигурации

Для получения пользовательской конфигурации необходимо получить из фреймворка инстанс конфигурации:

```go
cfg := app.GetConfig()
```

Общие методы:

- `Has(key string) bool` - Проверка наличия конфигурации.
- `GetString(key string, def string) string` - Получить как строку.
- `GetBool(key string, def bool) bool` - Получить как булево значение с конвертацией.
- `GetInt(key string, def int) int` - Получить как целое число.
- `GetFloat(key string, def float32) float32` - Получить как дробное число.

Функции получения карт:

- `GetMapString(key string, def map[string]string) map[string]string`
- `GetMapBool(key string, def map[string]bool) map[string]bool`
- `GetMapInt(key string, def map[string]int) map[string]int`
- `GetMapFloat(key string, def map[string]float64) map[string]float64`

Функции получения массивов:

- `GetSliceString(key string, def []string) []string`
- `GetSliceBool(key string, def []bool) []bool`
- `GetSliceInt(key string, def []int) []int`
- `GetSliceFloat(key string, def []float64) []float64`

### Далее [Логирование](./logger.md)