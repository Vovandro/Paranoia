# Экспортер метрик приложения

### [Оглавление](./index.md)

Во фреймворк встроены счетчики метрик, для получения доступа к ним необходимо инициализировать экспорт метрик, поддерживается экспортеры:
- Std - вывод в стандартный вывод
- Prometheus - получение метрик в данном формате по http

## Std

Инициализация через авто конфигурацию, указывается интервал через который будут выводиться все метрики в вывод:

```yaml
engine:
  metrics:
    - type: std
      name: app
      interval: 60s
```

Инициализация в коде

```go
app.SetMetrics(telemetry.NewStd(telemetry.MetricStdConfig{
	Name: "app",
	Interval: time.Second * 60,
}))
```

## Prometheus

Получение метрик в формате prometheus по http.

Инициализация через авто конфиг:

```yaml
engine:
  metrics:
    - type: prometheus
      name: app
      port: 8090
```

Инициализация в коде

```go
app.SetMetrics(telemetry.NewPrometheus(telemetry.MetricPrometheusConfig{
	Name: "app",
	port: "8090",
}))
```

В данном случае метрики будут доступны по адресу http://127.0.0.1:8090


### Далее [Список встроенных метрик](./metrics-list.md)