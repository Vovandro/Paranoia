package Paranoia

import (
	"fmt"
	"time"

	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
)

type Engine struct {
	name string

	starting bool

	config         interfaces.IConfig
	logger         interfaces.ILogger
	metricExporter interfaces.IMetrics
	trace          interfaces.ITrace

	task task

	cache       map[string]interfaces.ICache
	database    map[string]interfaces.IDatabase
	noSql       map[string]interfaces.INoSql
	controllers map[string]interfaces.IController
	modules     map[string]interfaces.IModules
	repository  map[string]interfaces.IRepository
	service     map[string]interfaces.IService
	servers     map[string]interfaces.IServer
	clients     map[string]interfaces.IClient
	storage     map[string]interfaces.IStorage
	middlewares map[string]interfaces.IMiddleware
}

func New(name string, config interfaces.IConfig, logger interfaces.ILogger) *Engine {
	t := &Engine{}

	t.starting = false
	t.name = name
	t.config = config
	t.logger = logger

	t.cache = make(map[string]interfaces.ICache)
	t.database = make(map[string]interfaces.IDatabase)
	t.noSql = make(map[string]interfaces.INoSql)
	t.controllers = make(map[string]interfaces.IController)
	t.modules = make(map[string]interfaces.IModules)
	t.repository = make(map[string]interfaces.IRepository)
	t.service = make(map[string]interfaces.IService)
	t.servers = make(map[string]interfaces.IServer)
	t.storage = make(map[string]interfaces.IStorage)
	t.clients = make(map[string]interfaces.IClient)
	t.middlewares = make(map[string]interfaces.IMiddleware)

	t.task.Init(t)

	if t.config != nil {
		err := t.config.Init(t)

		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if t.logger != nil {
		err := t.logger.Init(t.config)

		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return t
}

func (t *Engine) GetLogger() interfaces.ILogger {
	return t.logger
}

func (t *Engine) GetConfig() interfaces.IConfig {
	return t.config
}

func (t *Engine) SetMetrics(c interfaces.IMetrics) {
	if t.metricExporter != nil {
		_ = t.metricExporter.Stop()
	}

	t.metricExporter = c

	if t.metricExporter != nil {
		err := t.metricExporter.Init(t)

		if err != nil {
			t.logger.Error(err)
		}
	}
}

func (t *Engine) SetTrace(c interfaces.ITrace) {
	if t.trace != nil {
		_ = t.trace.Stop()
	}

	t.trace = c

	if t.trace != nil {
		err := t.trace.Init(t)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (t *Engine) PushCache(c interfaces.ICache) interfaces.IEngine {
	if _, ok := t.cache[c.String()]; ok {
		t.logger.Fatal(fmt.Errorf("cache %s already exists", c.String()))
	} else {
		t.cache[c.String()] = c
	}

	return t
}

func (t *Engine) GetCache(key string) interfaces.ICache {
	return t.cache[key]
}

func (t *Engine) PushDatabase(b interfaces.IDatabase) interfaces.IEngine {
	if _, ok := t.database[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("database %s already exists", b.String()))
	} else {
		t.database[b.String()] = b
	}

	return t
}

func (t *Engine) GetDatabase(key string) interfaces.IDatabase {
	return t.database[key]
}

func (t *Engine) PushNoSql(b interfaces.INoSql) interfaces.IEngine {
	if _, ok := t.noSql[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("nosql %s already exists", b.String()))
	} else {
		t.noSql[b.String()] = b
	}

	return t
}

func (t *Engine) GetNoSql(key string) interfaces.INoSql {
	return t.noSql[key]
}

func (t *Engine) PushController(b interfaces.IController) interfaces.IEngine {
	if _, ok := t.controllers[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("controller %s already exists", b.String()))
	} else {
		t.controllers[b.String()] = b
	}

	return t
}

func (t *Engine) GetController(key string) interfaces.IController {
	return t.controllers[key]
}

func (t *Engine) PushModule(b interfaces.IModules) interfaces.IEngine {
	if _, ok := t.modules[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("module %s already exists", b.String()))
	} else {
		t.modules[b.String()] = b
	}

	return t
}

func (t *Engine) GetModule(key string) interfaces.IModules {
	return t.modules[key]
}

func (t *Engine) PushRepository(b interfaces.IRepository) interfaces.IEngine {
	if _, ok := t.repository[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("repository %s already exists", b.String()))
	} else {
		t.repository[b.String()] = b
	}

	return t
}

func (t *Engine) GetRepository(key string) interfaces.IRepository {
	return t.repository[key]
}

func (t *Engine) PushService(b interfaces.IService) interfaces.IEngine {
	if _, ok := t.service[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("service %s already exists", b.String()))
	} else {
		t.service[b.String()] = b
	}

	return t
}

func (t *Engine) GetService(key string) interfaces.IService {
	return t.service[key]
}

func (t *Engine) PushServer(b interfaces.IServer) interfaces.IEngine {
	if _, ok := t.servers[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("server %s already exists", b.String()))
	} else {
		t.servers[b.String()] = b
	}

	return t
}

func (t *Engine) GetTask(key string) interfaces.ITask {
	return t.task.GetTask(key)
}

func (t *Engine) PushTask(b interfaces.ITask) interfaces.IEngine {
	t.task.PushTask(b, t.starting)

	return t
}

func (t *Engine) RemoveTask(key string) {
	t.task.RemoveTask(key)
}

func (t *Engine) RunTask(key string, args map[string]interface{}) error {
	return t.task.RunTask(key, args)
}

func (t *Engine) GetServer(key string) interfaces.IServer {
	return t.servers[key]
}

func (t *Engine) PushClient(b interfaces.IClient) interfaces.IEngine {
	if _, ok := t.clients[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("client %s already exists", b.String()))
	} else {
		t.clients[b.String()] = b
	}

	return t
}

func (t *Engine) GetClient(key string) interfaces.IClient {
	return t.clients[key]
}

func (t *Engine) PushStorage(b interfaces.IStorage) interfaces.IEngine {
	if _, ok := t.storage[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("storage %s already exists", b.String()))
	} else {
		t.storage[b.String()] = b
	}

	return t
}

func (t *Engine) GetStorage(key string) interfaces.IStorage {
	return t.storage[key]
}

func (t *Engine) PushMiddleware(b interfaces.IMiddleware) interfaces.IEngine {
	if _, ok := t.middlewares[b.String()]; ok {
		t.logger.Fatal(fmt.Errorf("middleware %s already exists", b.String()))
	} else {
		t.middlewares[b.String()] = b
	}

	return t
}

func (t *Engine) GetMiddleware(key string) interfaces.IMiddleware {
	return t.middlewares[key]
}

func (t *Engine) Init() error {
	var err error = nil

	for _, cache := range t.cache {
		err = cache.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, db := range t.database {
		err = db.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, db := range t.noSql {
		err = db.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, st := range t.storage {
		err = st.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, client := range t.clients {
		err = client.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	if _, ok := t.middlewares["timing"]; !ok {
		t.PushMiddleware(middleware.NewTimingMiddleware("timing"))
	}

	if _, ok := t.middlewares["restore"]; !ok {
		t.PushMiddleware(middleware.NewRestoreMiddleware("restore"))
	}

	if _, ok := t.middlewares["timeout"]; !ok {
		t.PushMiddleware(middleware.NewTimeoutMiddleware("timeout", middleware.TimeoutMiddlewareConfig{Timeout: time.Second * 60}))
	}

	for _, item := range t.middlewares {
		err = item.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, server := range t.servers {
		err = server.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, repository := range t.repository {
		err = repository.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, service := range t.service {
		err = service.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, module := range t.modules {
		err = module.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, controller := range t.controllers {
		err = controller.Init(t)

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	t.task.Start()

	for _, server := range t.servers {
		err = server.Start()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	if t.trace != nil {
		err = t.trace.Start()

		if err != nil {
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Start()

		if err != nil {
			return err
		}
	}

	t.starting = true

	return err
}

func (t *Engine) Stop() error {
	var err error = nil

	t.starting = false

	for _, server := range t.servers {
		err = server.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	t.task.Stop()

	for _, item := range t.middlewares {
		err = item.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, controller := range t.controllers {
		err = controller.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, module := range t.modules {
		err = module.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, service := range t.service {
		err = service.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, repository := range t.repository {
		err = repository.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, client := range t.clients {
		err = client.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, st := range t.storage {
		err = st.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, db := range t.noSql {
		err = db.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, db := range t.database {
		err = db.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	for _, cache := range t.cache {
		err = cache.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Stop()

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if t.trace != nil {
		err = t.trace.Stop()

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = t.config.Stop()

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = t.logger.Stop()

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
