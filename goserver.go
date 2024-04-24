package Paranoia

import (
	"fmt"
)

type Service struct {
	name string

	config IConfig
	logger ILogger

	cache       map[string]ICache
	brokers     map[string]IBroker
	database    map[string]IDatabase
	controllers map[string]IController
	modules     map[string]IModules
	repository  map[string]IRepository
	servers     map[string]IServer
	store       map[string]IStore
}

func (t *Service) New(name string, config IConfig, logger ILogger) *Service {
	t.name = name
	t.config = config
	t.logger = logger

	t.cache = make(map[string]ICache)
	t.brokers = make(map[string]IBroker)

	return t
}

func (t *Service) GetLogger() ILogger {
	return t.logger
}

func (t *Service) PushCache(c ICache) {
	t.cache[c.String()] = c
}

func (t *Service) GetCache(key string) ICache {
	return t.cache[key]
}

func (t *Service) PushBroker(b IBroker) {
	t.brokers[b.String()] = b
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

	for _, modules := range t.modules {
		err = modules.Init(t)

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

	return err
}
