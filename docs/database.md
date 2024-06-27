# SQL базы данных

### [Оглавление](./index.md)


## - [MySQL](./database-mysql.md)
## - [PostgreSQL](./database-postgres.md)
## - [SQLite3](./database-sqlite.md)
## - [Clickhouse](./database-clickhouse.md)

Общие методы:

- `Query(ctx context.Context, query string, args ...interface{}) (SQLRows, error)` - выполнение запроса и получение множественных строк
- `QueryRow(ctx context.Context, query string, args ...interface{}) (SQLRow, error)` - Выполнение запроса и получение 1 строки результата
- `Exec(ctx context.Context, query string, args ...interface{}) error` - Выполнение запроса без получения результатов
- `GetDb() interface{}` - Получение прямого инстанса базы данных, **использовать с осторожностью**

`SQLRows` - имеет методы:

- `Next() bool` - для цикла получения следующего элемента 
- `Scan(dest ...any) error` - Сканирование данных в переменные
- `Close() error` - Обязательно закрыть запрос по завершении обработки

`SQLRow` - Имеет аналогичный метод `Scan`, сканирование происходит в базовые структуры по 1 полю:

Пример:

```go
rows, err := db.Query(context.Background(), "SELECT id, name FROM users WHERE id < ?", 5)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	var item Model
	rows.Scan(&item.Id, &item.Name)
	items = append(items, item)
}
```

### Далее [MySQL](./database-mysql.md)