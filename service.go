package Paranoia

import (
	"Paranoia/interfaces"
	"fmt"
)

type Service struct {
	name string

	config interfaces.IConfig
	logger interfaces.ILogger

	cache       map[string]interfaces.ICache
	brokers     map[string]interfaces.IBroker
	database    map[string]interfaces.IDatabase
	controllers map[string]interfaces.IController
	modules     map[string]interfaces.IModules
	repository  map[string]interfaces.IRepository
	servers     map[string]interfaces.IServer
	storage     map[string]interfaces.IStorage
}

func New(name string, config interfaces.IConfig, logger interfaces.ILogger) *Service {
	t := &Service{}

	t.name = name
	t.config = config
	t.logger = logger

	t.cache = make(map[string]interfaces.ICache)
	t.brokers = make(map[string]interfaces.IBroker)
	t.database = make(map[string]interfaces.IDatabase)
	t.controllers = make(map[string]interfaces.IController)
	t.modules = make(map[string]interfaces.IModules)
	t.repository = make(map[string]interfaces.IRepository)
	t.servers = make(map[string]interfaces.IServer)
	t.storage = make(map[string]interfaces.IStorage)

	return t
}

func (t *Service) GetLogger() interfaces.ILogger {
	return t.logger
}

func (t *Service) PushCache(c interfaces.ICache) interfaces.IService {
	t.cache[c.String()] = c

	return t
}

func (t *Service) GetCache(key string) interfaces.ICache {
	return t.cache[key]
}

func (t *Service) PushBroker(b interfaces.IBroker) interfaces.IService {
	t.brokers[b.String()] = b

	return t
}

func (t *Service) GetBroker(key string) interfaces.IBroker {
	return t.brokers[key]
}

func (t *Service) PushDatabase(b interfaces.IDatabase) interfaces.IService {
	t.database[b.String()] = b

	return t
}

func (t *Service) GetDatabase(key string) interfaces.IDatabase {
	return t.database[key]
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

func (t *Service) GetServer(key string) interfaces.IServer {
	return t.servers[key]
}

func (t *Service) PushStorage(b interfaces.IStorage) interfaces.IService {
	t.storage[b.String()] = b

	return t
}

func (t *Service) GetStorage(key string) interfaces.IStorage {
	return t.storage[key]
}

func (t *Service) Init() error {
	var err error = nil

	err = t.config.Init(t)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = t.logger.Init(t)

	if err != nil {
		fmt.Println(err)
		return err
	}

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

	for _, st := range t.storage {
		err = st.Init(t)

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

	for _, broker := range t.brokers {
		err = broker.Init(t)

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

	for _, broker := range t.brokers {
		err = broker.Stop()

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

	for _, st := range t.storage {
		err = st.Stop()

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
