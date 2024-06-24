# Быстрый старт

### [Оглавление](./index.md)

Для более простого использования фреймворка рекомендуется использовать MVC паттерны, [пример приложения](https://gitlab.com/devpro_studio/paranoia_example). 

Рекомендуется создание контроллеров для выполнения бизнес логики, модели для создания представлений, репозитории для взаимодействия с внешними системами и получения данных моделей.

Для облегчения взаимодействия во фреймворк встроены методы инициализации данных сущностей

```
s := Paranoia.New("minimal paranoia app", config.NewMock(), logger.NewMock()).
        PushController(MyController.New()).
        PushController(HelloController.New()).
        PushRepository(UserRepository.New())
```

Контроллеры должны соответствовать интерфейсу:

```
type IController interface {
	Init(app IService) error
	Stop() error
	String() string
}
```

Репозитории:

```
type IRepository interface {
	Init(app IService) error
	Stop() error
	String() string
}
```

Пользовательские модули:

```
type IModules interface {
	Init(app IService) error
	Stop() error
	String() string
}
```

Все структуры в методе `String()` должны возвращать уникальное имя.

В методе `Init(app IService) error` допускается получение других модулей фреймворка, они в этот момент уже проинициализированы. 

В методе `Stop() error` необходимо корректно завершить работу и освободить занятые ресурсы

### Далее [Конфигурация системы](./config-index.md)