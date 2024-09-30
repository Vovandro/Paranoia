# Экспортер трассировок приложения

### [Оглавление](./index.md)

Во фреймворк встроены трассировки запросов и выполнения программы, для получения доступа к ним необходимо инициализировать экспорт, поддерживается экспортеры:
- Std - вывод в стандартный вывод
- Zipkin

## Std

Инициализация через авто конфигурацию, указывается интервал через который будут выводиться все в вывод:

```yaml
engine:
  trace:
    - type: std
      name: app
      interval: 60s
```

Инициализация в коде

```go
app.SetTrace(telemetry.NewTraceStd(telemetry.TraceStdConfig{
	Name: "app",
	Interval: time.Second * 60,
}))
```

## Zipkin

Инициализация через авто конфигурацию, указывается url сбора zipkin:

```yaml
engine:
  trace:
    - type: zipkin
      name: app
      url: http://localhost
```

Инициализация в коде

```go
app.SetTrace(telemetry.NewTraceZipkin(telemetry.TraceZipkinConfig{
	Name: "app",
	Url: "http://localhost",
}))
```

## Sentry

Инициализация через авто конфигурацию:

```yaml
engine:
  trace:
    - type: sentry
      name: app
```

Инициализация в коде

```go
app.SetTrace(telemetry.NewTraceSentry(telemetry.TraceSentryConfig{
	Name: "app",
}))
```

* Для работы трассировки сентри необходимо использовать логирование в сентри *


## OTLP

Доступны экспортеры HTTP и GRPC

```yaml
engine:
  trace:
    - type: otlp_http
      name: app
```

```yaml
engine:
  trace:
    - type: otlp_grpc
      name: app
```


### Далее [Системы кеширования](./cache.md)