# Логирование

### [Оглавление](./index.md)

Система логирования задается при инициализации приложения, возможно изменение параметров во время выполнения.


Фреймворк штатно поддерживает системы логирования:
- Mock - используется в качестве заглушек, ничего не делает
- Std - Вывод в стандартный вывод
- File - Вывод логов в файл

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

## Возможно каскадное вложение модулей.

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


### Далее [Экспортер метрик приложения](./metrics.md)