# Системы кеширования

### [Оглавление](./index.md)


## - [Memory cache](./cache-memory.md)
## - [Memcached](./cache-memcached.md)
## - [Redis](./cache-redis.md)

Общие методы:

- `Has(key string) bool` - проверяет наличие ключа.
- `Set(key string, args any, timeout time.Duration) error` - Создает или изменяет значение по ключу и устанавливает время жизни.
- `SetIn(key string, key2 string, args any, timeout time.Duration) error` - Создает или изменяет карту, в карте данные записывает под вторым ключом и устанавливает время жизни всей карты.
- `SetMap(key string, args any, timeout time.Duration) error` - Создает или изменяет значение (карту) по ключу и устанавливает время жизни.
- `Get(key string) (any, error)` - Получает данные по ключу.
- `GetIn(key string, key2 string) (any, error)` - Получает данные из карты по ключу и значение из поля по второму ключу.
- `GetMap(key string) (any, error)` - Получает карту целиком.
- `Increment(key string, val int64, timeout time.Duration) (int64, error)` - Увеличивает счетчик по ключу на определенное значение и устанавливает время жизни записи, возвращает новое значение. 
- `IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error)` - Увеличивает счетчик в карте по ключу и поле по второму ключу на определенное значение и устанавливает время жизни всей карты, возвращает новое значение.
- `Decrement(key string, val int64, timeout time.Duration) (int64, error)` - Уменьшает счетчик по ключу на определенное значение и устанавливает время жизни записи, возвращает новое значение.
- `DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error)` - Уменьшает счетчик в карте по ключу и поле по второму ключу на определенное значение и устанавливает время жизни всей карты, возвращает новое значение.
- `Delete(key string) error` - Удаляет данные из кеша по ключу.

Общие ошибки:

- `cache.ErrKeyNotFound` - Ключа не существует или время его жизни кончилось 
- `cache.ErrTypeMismatch` - Ошибка формата данных, например если попытаться установить значение в карте, а под данным ключом храниться строка.

### Далее [Memory cache](./cache-memory.md)