# Paranoia framework - golang microservice engine

[![pipeline status](https://gitlab.com/devpro_studio/Paranoia/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![coverage report](https://gitlab.com/devpro_studio/Paranoia/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![Latest Release](https://gitlab.com/devpro_studio/Paranoia/-/badges/release.svg)](https://gitlab.com/devpro_studio/Paranoia/-/releases)
[![Go Reference](https://pkg.go.dev/badge/gitlab.com/devpro_studio/Paranoia)](https://pkg.go.dev/gitlab.com/devpro_studio/Paranoia)

## [Documentations rus](./docs/index.md)


## [GUI Helper project tool](https://gitlab.com/devpro_studio/paranoia-gui)

## Supported:

### Database

- Postgres
- SQLite
- MySQL
- Clickhouse
- MongoDB
- Aerospike

### Cache

- Memory
- Redis
- Memcached
- etcd

### Servers

- http
- Kafka
- GRPC

### Server middlewares

- Restore from panic
- Register timing middleware (default use)
- Timeout request
- Authorize (JWT)

### Clients

- http
- Kafka
- GRPC

### Storage

- File
- S3

### Other

- Initialize base engine module from yaml config file
- Regulatory task system 
- Sentry log
- JWT native support (module and middleware)

Generating RSA keys for JWT:

`openssl genrsa -out private.key 2048`

`openssl rsa -in private.key -pubout -out public.key`
