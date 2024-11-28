# Redis

### [Оглавление](./index.md)

## Использование

Инициализация через авто конфиг:

```yaml
engine:
  cache:
    - type: redis
      name: primary
      hosts: "localhost:6379"
      use_cluster: false
      db_num: 1
      timeout: 3s
      username: 
      password:
```

Инициализация в коде

```go
app.PushCache(cache.NewRedis("primary", cache.RedisConfig{
    Hosts: "localhost:6379",
	Timeout: time.Second * 3,
    DBNum: 1,
}))
```

Далее в необходимых местах можно получить кеш:

```go
cache := app.GetCache("primary")
```

## Особенности

В БД данные хранятся и возвращаются в виде строк.

Не поддерживается декремент менее 0.

### Далее [etcd](./cache-etcd.md)