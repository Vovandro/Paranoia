# Paranoia framework - golang microservice engine

[![pipeline status](https://gitlab.com/devpro_studio/Paranoia/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master)
[![coverage report](https://gitlab.com/devpro_studio/Paranoia/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master)
[![Go Reference](https://pkg.go.dev/badge/gitlab.com/devpro_studio/Paranoia)](https://pkg.go.dev/gitlab.com/devpro_studio/Paranoia)

## [Use GO helper library](https://gitlab.com/devpro_studio/go_utils)

## [GUI Helper project tool](https://gitlab.com/devpro_studio/paranoia-gui)

<details>
<summary>Getting Started</summary>

To install in the project, use the command

```shell
go get gitlab.com/devpro_studio/Paranoia
```

A minimal application includes the initialization of the framework:

```go
s := paranoia.New("minimal paranoia app", "cfg.yaml")
```

The first parameter is the application name, the second is the configuration system, and the third is the logging system. In this example, stub objects are used as the configuration and logging systems.

Next, the framework is populated with modules that will be used in this service, for example, add an in-memory cache at the application level:

```go
s.PushPkg(memory.New("secondary"))
```

In all engine modules, the module name and its type are used, the name must be unique within the type, and it can be used later in the code to get this module. More details about available modules and possible settings are described later in the documentation.

Next, you need to initialize the framework and start it:

```go
err := s.Init()

if (err != nil) {
    panic(err)
    return
}

defer s.Stop()
```

The minimal application is ready.

</details>

## [Documentation](https://custom-site.ru/)

## Supported:

### Database

- Postgres
- SQLite
- MySQL
- Clickhouse
- MongoDB
- Aerospike
- Elasticsearch

### Cache

- Memory
- Redis
- Memcached
- etcd

### Servers

- http
- Kafka
- RabbitMQ
- GRPC

### Server middlewares

- Restore from panic
- Register timing middleware (default use)
- Timeout request
- Authorize (JWT)
- CORS

### Clients

- http
- Kafka
- RabbitMQ
- GRPC

### Storage

- File
- S3

### Other

- Initialize base engine module from yaml config file
- Regulatory task system
- Sentry log
- JWT native support (module and middleware)
- Concurrency patterns in template

Generating RSA keys for JWT:

`openssl genrsa -out private.key 2048`

`openssl rsa -in private.key -pubout -out public.key`
