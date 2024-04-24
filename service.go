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
	store       map[string]interfaces.IStore
}

func New(name string, config interfaces.IConfig, logger interfaces.ILogger) *Service {
	t := &Service{}

	t.name = name
	t.config = config
	t.logger = logger

	t.cache = make(map[string]interfaces.ICache)
	t.brokers = make(map[string]interfaces.IBroker)

	return t
}

func (t *Service) GetLogger() interfaces.ILogger {
	return t.logger
}

func (t *Service) PushCache(c interfaces.ICache) {
	t.cache[c.String()] = c
}

func (t *Service) GetCache(key string) interfaces.ICache {
	return t.cache[key]
}

func (t *Service) PushBroker(b interfaces.IBroker) {
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
