# MySQL

### [Оглавление](./index.md)

## Использование

Инициализация через авто конфиг:

```yaml
engine:
  database:
    - type: mysql
      name: primary
      uri: "test:test@(127.0.0.1:3306)/test?parseTime=true"
```

Инициализация в коде

```go
app.PushDatabase(database.NewMySQL("primary", database.MySQLConfig{
    URI: "test:test@(127.0.0.1:3306)/test?parseTime=true",
}))
```

Далее в необходимых местах можно получить:

```go
cache := app.GetDatabase("primary")
```

## Особенности

Для правильного сканирования дат необходимо передать параметр `parseTime=true`

### Далее [PostgreSQL](./database-postgres.md)