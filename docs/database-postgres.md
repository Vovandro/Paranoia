# PostgreSQL

### [Оглавление](./index.md)

## Использование

Инициализация через авто конфиг:

```yaml
engine:
  database:
    - type: postgres
      name: primary
      uri: "postgres://test:test@127.0.0.1:5432/test"
```

Инициализация в коде

```go
app.PushDatabase(database.NewPostgres("primary", database.PostgresConfig{
    URI: "postgres://test:test@127.0.0.1:5432/test",
}))
```

Далее в необходимых местах можно получить:

```go
db := app.GetDatabase("primary")
```

## Особенности

Параметры запроса передаются нумеровано: `"SELECT * FROM users WHERE id < $1 AND id > $2"`

### Далее [SQLite3](./database-sqlite.md)