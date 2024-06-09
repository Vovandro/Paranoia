# Paranoia framework - golang microservice engine

[![pipeline status](https://gitlab.com/devpro_studio/Paranoia/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![coverage report](https://gitlab.com/devpro_studio/Paranoia/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![Latest Release](https://gitlab.com/devpro_studio/Paranoia/-/badges/release.svg)](https://gitlab.com/devpro_studio/Paranoia/-/releases)


## Simple start:
Import to project `go get gitlab.com/devpro_studio/Paranoia`

add to main.go

```
	s := Paranoia.
		New("base paranoia app", &config.Env{}, &logger.File{&logger.Std{}}).
		PushCache(&cache.Memory{Name: "cache"}).
		PushRepository(&myRepository{Name: "repository"}).
		PushController(&myController{Name: "controller"})
	
	err := s.Init()

	if err != nil {
		panic(err)
		return
	}
	
	defer s.Stop()
```

## Supported:

### Database

- Postgres
- SQLite
- Clickhouse
- MongoDB
- Aerospike

### Cache

- Memory
- Redis
- Memcached

### Servers

- http
- Kafka

### Server middlewares

- Restore from panic
- Register timing middleware (default use)

### Clients

- http
- Kafka

### Other

- Initialize base engine module from yaml config file