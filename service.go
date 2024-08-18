package Paranoia

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"sync"
	"time"
)

type Service struct {
	name string

	config         interfaces.IConfig
	logger         interfaces.ILogger
	metricExporter interfaces.IMetrics

	cache       map[string]interfaces.ICache
	database    map[string]interfaces.IDatabase
	noSql       map[string]interfaces.INoSql
	controllers map[string]interfaces.IController
	modules     map[string]interfaces.IModules
	repository  map[string]interfaces.IRepository
	task        map[string]interfaces.ITask
	servers     map[string]interfaces.IServer
	clients     map[string]interfaces.IClient
	storage     map[string]interfaces.IStorage
	middlewares map[string]interfaces.IMiddleware

	taskMutex sync.RWMutex
}

func New(name string, config interfaces.IConfig, logger interfaces.ILogger) *Service {
	t := &Service{}

	t.name = name
	t.config = config
	t.logger = logger

	t.cache = make(map[string]interfaces.ICache)
	t.database = make(map[string]interfaces.IDatabase)
	t.noSql = make(map[string]interfaces.INoSql)
	t.controllers = make(map[string]interfaces.IController)
	t.modules = make(map[string]interfaces.IModules)
	t.repository = make(map[string]interfaces.IRepository)
	t.task = make(map[string]interfaces.ITask)
	t.servers = make(map[string]interfaces.IServer)
	t.storage = make(map[string]interfaces.IStorage)
	t.clients = make(map[string]interfaces.IClient)
	t.middlewares = make(map[string]interfaces.IMiddleware)

	t.taskMutex = sync.RWMutex{}

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

func (t *Service) GetLogger() interfaces.ILogger {
	return t.logger
}

func (t *Service) GetConfig() interfaces.IConfig {
	return t.config
}

func (t *Service) SetMetrics(c interfaces.IMetrics) {
	if t.metricExporter != nil {
		t.metricExporter.Stop()
	}

	t.metricExporter = c

	if t.metricExporter != nil {
		err := t.metricExporter.Init(t)

		if err != nil {
			t.logger.Error(err)
		}
	}
}

func (t *Service) PushCache(c interfaces.ICache) interfaces.IService {
	t.cache[c.String()] = c

	return t
}

func (t *Service) GetCache(key string) interfaces.ICache {
	return t.cache[key]
}

func (t *Service) PushDatabase(b interfaces.IDatabase) interfaces.IService {
	t.database[b.String()] = b

	return t
}

func (t *Service) GetDatabase(key string) interfaces.IDatabase {
	return t.database[key]
}

func (t *Service) PushNoSql(b interfaces.INoSql) interfaces.IService {
	t.noSql[b.String()] = b

	return t
}

func (t *Service) GetNoSql(key string) interfaces.INoSql {
	return t.noSql[key]
}

func (t *Service) PushController(b interfaces.IController) interfaces.IService {
	t.controllers[b.String()] = b

	return t
}

func (t *Service) GetController(key string) interfaces.IController {
	return t.controllers[key]
}

func (t *Service) PushModule(b interfaces.IModules) interfaces.IService {
	t.modules[b.String()] = b

	return t
}

func (t *Service) GetModule(key string) interfaces.IModules {
	return t.modules[key]
}

func (t *Service) PushRepository(b interfaces.IRepository) interfaces.IService {
	t.repository[b.String()] = b

	return t
}

func (t *Service) GetRepository(key string) interfaces.IRepository {
	return t.repository[key]
}

func (t *Service) PushServer(b interfaces.IServer) interfaces.IService {
	t.servers[b.String()] = b

	return t
}

func (t *Service) GetTask(key string) interfaces.ITask {
	t.taskMutex.RLock()
	defer t.taskMutex.RUnlock()

	return t.task[key]
}

func (t *Service) PushTask(b interfaces.ITask) interfaces.IService {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	if task, ok := t.task[b.String()]; ok {
		_ = task.Stop()
	}

	t.task[b.String()] = b

	_ = b.Init(t)
	b.Start()

	return t
}

func (t *Service) RemoveTask(key string) {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	if task, ok := t.task[key]; ok {
		_ = task.Stop()
		delete(t.task, key)
	}
}

func (t *Service) GetServer(key string) interfaces.IServer {
	return t.servers[key]
}

func (t *Service) PushClient(b interfaces.IClient) interfaces.IService {
	t.clients[b.String()] = b

	return t
}

func (t *Service) GetClient(key string) interfaces.IClient {
	return t.clients[key]
}

func (t *Service) PushStorage(b interfaces.IStorage) interfaces.IService {
	t.storage[b.String()] = b

	return t
}

func (t *Service) GetStorage(key string) interfaces.IStorage {
	return t.storage[key]
}

func (t *Service) PushMiddleware(b interfaces.IMiddleware) interfaces.IService {
	t.middlewares[b.String()] = b

	return t
}

func (t *Service) GetMiddleware(key string) interfaces.IMiddleware {
	return t.middlewares[key]
}

func (t *Service) Init() error {
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

	for _, server := range t.servers {
		err = server.Start()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Start()

		if err != nil {
			return err
		}
	}

	return err
}

func (t *Service) Stop() error {
	var err error = nil

	for _, server := range t.servers {
		err = server.Stop()

		if err != nil {
			t.logger.Fatal(err)
			return err
		}
	}

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

	err = t.metricExporter.Stop()

	if err != nil {
		fmt.Println(err)
		return err
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
