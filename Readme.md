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

<details>
<summary>Config application</summary>

The configuration system is set during the framework initialization and does not change within the project.

The framework natively supports configuration systems:
- Auto configuration from a yaml file - support for user configuration and settings for all built-in framework modules.

# Auto configuration

Allows you to initialize the application depending on the environment, loading occurs from a yaml file.

Two root tags are supported: `engine` for framework configuration and `cfg` for user configuration.

Example configuration file:

```yaml
engine:
  - type: metrics
    name: exporter
    service_name: example_app
    interval: 30s
  - type: cache
    name: primary
    time_clear: 10m

cfg:
  logLevel: WARNING
  key: val
  key_map:
    key1: val1
    key2: val2
    key_slice:
        - val1
        - val2
        - val3
```

The `cfg` block is for user configuration.

The `engine` block must include the module type, its name, and other parameters.

# Getting configuration

To get the user configuration, you need to get an instance of the configuration from the framework:

```go
cfg := app.GetConfig()
```

Common methods:

- `Has(key string) bool` - Check for the presence of a configuration.
- `GetString(key string, def string) string` - Get as a string.
- `GetBool(key string, def bool) bool` - Get as a boolean value with conversion.
- `GetInt(key string, def int) int` - Get as an integer.
- `GetFloat(key string, def float32) float32` - Get as a float.

Functions for getting maps:

- `GetMapString(key string, def map[string]string) map[string]string`
- `GetMapBool(key string, def map[string]bool) map[string]bool`
- `GetMapInt(key string, def map[string]int) map[string]int`
- `GetMapFloat(key string, def map[string]float64) map[string]float64`

Functions for getting slices:

- `GetSliceString(key string, def []string) []string`
- `GetSliceBool(key string, def []bool) []bool`
- `GetSliceInt(key string, def []int) []int`
- `GetSliceFloat(key string, def []float64) []float64`

Getting package configuration data:

- `GetConfigItem(typeName string, name string) map[string]interface{}`

</details>

<details>
<summary>Logging</summary>

The framework supports logging systems:
- Std - Output to standard output
- File - Output logs to a file
- Sentry - Logging to Sentry

# Std

Output logs to the console or other standard output

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/logger/std-log
```

Configuration:

```yaml
- type: logger
  name: std
  level: INFO
  enable: true
```

```go
app.PushPkg(std_log.New("std"))
```

# File

Output logs to a file with support for automatic file recreation at the beginning of the day. The final file name is appended with the date and the log extension.

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/logger/file-log
```

Configuration:

```yaml
- type: logger
  name: file
  level: INFO
  filename: app
  enable: true
```

```go
app.PushPkg(file_log.New("file"))
```

# Sentry

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/logger/sentry-log
```

```yaml
- type: logger
  name: sentry
  level: INFO
  sentry_url: http://sentry:9000
  app_env: dev
  sample_rate: 1.0
  trace_sample_rate: 0.1
  enable: true
```

```go
app.PushPkg(sentry_log.New("sentry"))
```

Only logs of levels other than DEBUG are written to Sentry. Through the context, it is possible to pass `span` of type `*sentry.Span`, and `tags` of type `map[string]string`.

# Cascading module nesting is possible.

For example, output to both file and console simultaneously:

```go
app.PushPkg(std_log.New("std")).
    PushPkg(file_log.New("file"))
```

</details>

<details>
<summary>Metrics</summary>

Metrics are used to monitor the operation of the application and its components.

The framework includes metric counters, to access them you need to initialize metric export, supported exporters:
- Std - output to standard output
- Prometheus - get metrics in this format via http
- OTLP - send metrics in OTLP format (http or grpc)

## Std

```yaml
- type: metrics
  name: app
  service_name: example application
  interval: 60s
```

```go
app.SetMetrics(telemetry.NewMetricStd("app"))
```

Or use name in config is "std" for auto config from framework and no use SetMetrics

## Prometheus

Get metrics in prometheus format via http.

```yaml
- type: metrics
  name: app
  service_name: example application
  port: 8090
```

```go
app.SetMetrics(telemetry.NewPrometheusMetrics("app"))
```

Or use name in config is "prometheus" for auto config from framework and no use SetMetrics

In this case, the metrics will be available at http://127.0.0.1:8090

## OTLP

Available exporters HTTP and GRPC

```yaml
- type: metrics
  name: app
  service_name: example application
  interval: 60s
```

```go
app.SetMetrics(telemetry.NewMetricOtlpHttp("app"))
```

Or

```go
app.SetMetrics(telemetry.NewMetricOtlpGrpc("app"))
```

Or use name in config is "oltp_grpc"\"oltp_http" for auto config from framework and no use SetMetrics

## Base metrics

The framework already includes metric counters in most of the main modules, different modules within the same package have the same semantics.

All metrics have the format **{Module Type}**.**{Module Name}**.**{Metric Name}**

## Caching systems:

- **.countRead** - constantly increasing operation counter
- **.countWrite** - constantly increasing operation counter
- **.timeRead** - request time histogram
- **.timeWrite** - request time histogram

## Databases:

- **.count** - constantly increasing operation counter
- **.time** - request time histogram

## Server area:

- **.count** - constantly increasing operation counter
- **.count_error** - constantly increasing error counter
- **.time** - request time histogram

## Clients:

- **.count** - constantly increasing operation counter
- **.time** - request time histogram
- **.retry** - retry count histogram

</details>

<details>
<summary>Trace</summary>

The framework includes request and program execution traces. To access them, you need to initialize the export. Supported exporters:
- Std - output to standard output
- Zipkin
- Sentry
- OLTP

## Std

```yaml
- type: trace
  name: app
  service_name: example app
  interval: 60s
```

```go
app.SetTrace(telemetry.NewTraceStd("app"))
```

Or use name in config is "std" for auto config from framework and no use SetTrace

## Zipkin

```yaml
- type: trace
  name: app
  service_name: example app
  url: http://localhost
```

```go
app.SetTrace(telemetry.NewTraceZipkin("app"))
```

Or use name in config is "zipkin" for auto config from framework and no use SetTrace

## Sentry

```yaml
- type: trace
  name: app
  service_name: example app
```

```go
app.SetTrace(telemetry.NewTraceSentry("app"))
```

Or use name in config is "sentry" for auto config from framework and no use SetTrace

* To use Sentry tracing, you need to use Sentry logging *


## OTLP

Available exporters HTTP and GRPC

```yaml
- type: trace
  name: app
  service_name: example app
```


```go
app.SetTrace(telemetry.NewTraceOtlpHttp("app"))
```

or

```go
app.SetTrace(telemetry.NewTraceOtlpGrpc("app"))
```

Or use name in config is "oltp_grpc"\"oltp_http" for auto config from framework and no use SetTrace

</details>

<details>
<summary>Cache</summary>

<details>
<summary>ETCD</summary>

## Usage

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/cache/etcd
```

```yaml
- type: cache
  name: primary
  hosts: "localhost:2379"
  username: 
  password:
  key_prefix:
```

```go
app.PushPkg(etcd.New("primary")
```

Next, you can get the cache in the necessary places:

```go
cache := app.GetPkg(interfaces.PkgCache, "primary").(etcd.IEtcd)
```

## Features

etcd cannot work with maps and increments, for all functions working with maps, JSON decoding and data conversion are used.
Key renewal is only possible with the old ttl.
Only the string type is possible as a value.

In the database, data is stored and returned as a byte slice, except for getting a map, where JSON decoding is used with conversion to default types for this operation.

</details>

<details>
<summary>Memcached</summary>

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/cache/memcached
```

```yaml
- type: cache
  name: primary
  hosts: "localhost:11211"
  timeout: 3s
  key_prefix:
```

```go
app.PushPkg(memcached.New("primary"))
```

Next, you can get the cache in the necessary places:

```go
cache := app.GetPkg(interfaces.PkgCache, "primary").(memcached.IMemcached)
```

## Features

Memcached does not work with maps, for all functions working with maps, JSON decoding is used.

In the database, data is stored and returned as a byte slice, except for getting a map, where JSON decoding is used with type conversion by default for this operation.

Decrement is not supported as the first operation if the key is missing, except for nested maps.

</details>

<details>
<summary>Memory</summary>

Used for fast application-level caching. Supports all basic cache operations including timeouts and map operations.

## Usage

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/cache/memory
```

```yaml
- type: cache
  name: secondary
  time_clear: 10m
  shard_count: 10
  enable_storage: true
  storage_file: cache.back
```

```go
app.PushPkg(memory.New("secondary"))
```

The clear time sets the garbage collector pass time.
The number of shards allows you to separate locks and speed up cache operation.

Next, you can get the cache in the necessary places:

```go
cache := app.GetPkg(interfaces.PkgCache, "primary").(memory.IMemory)
```

## Features

Stores and returns data in any format, the format type does not change during storage.

</details>

<details>
<summary>Redis</summary>

## Usage

```shell
go get gitlab.com/devpro_studio/Paranoia/pkg/cache/redis
```

```yaml
- type: cache
  name: primary
  hosts: "localhost:6379"
  use_cluster: false
  db_num: 1
  timeout: 3s
  username: 
  password:
  key_prefix:
```

```go
app.PushPkg(redis.New("primary"))
```

Next, you can get the cache in the necessary places:

```go
cache := app.GetPkg(interfaces.PkgCache, "primary").(redis.IRedis)
```

## Features

In the database, data is stored and returned as strings.

Decrement less than 0 is not supported.

</details>

</details>

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
- RabbitMQ
- GRPC

### Server middlewares

- Restore from panic
- Register timing middleware (default use)
- Timeout request
- Authorize (JWT)

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
