# Логирование

### [Оглавление](./index.md)

Система логирования задается при инициализации приложения, возможно изменение параметров во время выполнения.


Фреймворк штатно поддерживает системы логирования:
- Mock - используется в качестве заглушек, ничего не делает
- Std - Вывод в стандартный вывод
- File - Вывод логов в файл
- Sentry - Логирование в Sentry

## Mock

Используется для заглушек в тестах

```go
s := Paranoia.New("app", nil, logger.NewMock())
```

# Std

Вывод логов в консоль или другой стандартный вывод

```go
s := Paranoia.New("app", nil, logger.NewStd(logger.StdConfig{
    Level: interfaces.INFO,
}, nil))
```

# File

Вывод логов в файл с поддержкой автоматического пересоздания файла в начале суток. К имени итогового файла добавляется дата и расширение log.

```go
s := Paranoia.New("app", nil, logger.NewFile(logger.FileConfig{
    Level: interfaces.INFO,
	FName: "app",
}, nil))
```

# Sentry

```go
s := Paranoia.New("app", nil, logger.NewSentry(logger.SentryConfig{
    Level: interfaces.INFO,
	SentryURL: "key@localhost.ru",
    AppEnv: "local",
	SampleRate: 0.1,
	TraceSampleRate: 0.5,
}, nil))
```

В sentry пишутся только логи все кроме DEBUG уровня. Через контекст возможно передача `span` с типом `*sentry.Span`, и `tags` с типом `map[string]string`

# Возможно каскадное вложение модулей.

К примеру вывод в файл и консоль одновременно:

```go

s := Paranoia.New("app", nil, 
        logger.NewFile(
            logger.FileConfig{
                Level: interfaces.INFO,
                FName: "app",
            },
            logger.NewStd(logger.StdConfig{
                Level: interfaces.INFO,
            }, nil), 
        ),
	)
```

# Пользовательское логирование

Для использования необходимо получить из фреймворка инстанс логирования:

```go
cfg := app.GetLogger()
```

Общие методы:

- `SetLevel(level LogLevel)` - Изменения уровня логирования
- `Push(level LogLevel, msg string, toParent bool)` - Прямая запись сообщений
- `Debug(args ...interface{})` - Сахар для записи сообщения
- `Info(args ...interface{})` - Сахар для записи сообщения
- `Warn(args ...interface{})` - Сахар для записи сообщения
- `Message(args ...interface{})` - Сахар для записи сообщения
- `Error(err error)` - Сахар для записи сообщения
- `Fatal(err error)` - Сахар для записи сообщения
- `Panic(err error)` - Сахар для записи сообщения

## Конфигурация из файла

Доступно автоматическая замена конфигурации по используемой конфигурации:

- `LOG_LEVEL` - уровень логирования

Индивидуально для sentry:
- `SENTRY_URL` - Sentry DSN
- `APP_ENV` - Название окружения
- `SENTRY_SAMPLE_RATE` - Процент логирования
- `SENTRY_TRACE_SAMPLE_RATE` - процент трассировки

### Далее [Экспортер метрик приложения](./metrics.md)