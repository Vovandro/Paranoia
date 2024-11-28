package config

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/cache"
	"gitlab.com/devpro_studio/Paranoia/client"
	"gitlab.com/devpro_studio/Paranoia/database"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/module"
	"gitlab.com/devpro_studio/Paranoia/noSql"
	"gitlab.com/devpro_studio/Paranoia/server"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"gitlab.com/devpro_studio/Paranoia/storage"
	"gitlab.com/devpro_studio/Paranoia/telemetry"
	"gitlab.com/devpro_studio/Paranoia/utils/decoder"
	"gopkg.in/yaml.v3"
	"os"
)

type cfgItem map[string]interface{}
type cfgModule []cfgItem

type Data struct {
	Engine map[string]cfgModule   `yaml:"engine"`
	Cfg    map[string]interface{} `yaml:"cfg"`
}

type Auto struct {
	cfg AutoConfig

	app  interfaces.IEngine
	data Data
}

func NewAuto(cfg AutoConfig) *Auto {
	return &Auto{
		cfg: cfg,
	}
}

type AutoConfig struct {
	FName string `yaml:"filename"`
}

func (t *Auto) loadConfig() error {
	yamlFile, err := os.ReadFile(t.cfg.FName)
	if err != nil {
		return err
	}

	if t.data.Engine == nil {
		t.data.Engine = make(map[string]cfgModule, 10)
	}

	if t.data.Cfg == nil {
		t.data.Cfg = make(map[string]interface{}, 10)
	}

	err = yaml.Unmarshal(yamlFile, &t.data)

	if err != nil {
		return err
	}

	if v, ok := t.data.Cfg["logLevel"]; ok {
		t.app.GetLogger().SetLevel(interfaces.GetLogLevel(v.(string)))
	}

	for typeModule, modules := range t.data.Engine {
		for _, mod := range modules {
			name, ok := mod["name"]

			if !ok {
				return fmt.Errorf("not found name mod")
			}

			nameModule, ok := mod["type"]

			if !ok {
				return fmt.Errorf("not found type mod %s", name)
			}

			delete(mod, "type")

			if typeModule != "metrics" && typeModule != "trace" {
				delete(mod, "name")
			}

			switch typeModule {
			case "cache":
				switch nameModule {
				case "memcached":
					cfg := cache.MemcachedConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewMemcached(name.(string), cfg))

				case "memory":
					cfg := cache.MemoryConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewMemory(name.(string), cfg))

				case "redis":
					cfg := cache.RedisConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewRedis(name.(string), cfg))

				case "etcd":
					cfg := cache.EtcdConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewEtcd(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "client":
				switch nameModule {
				case "http":
					cfg := client.HTTPClientConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushClient(client.NewHTTPClient(name.(string), cfg))

				case "kafka":
					cfg := client.KafkaClientConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushClient(client.NewKafkaClient(name.(string), cfg))

				case "grpc":
					cfg := client.GrpcClientConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushClient(client.NewGrpcClient(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "database":
				switch nameModule {
				case "clickhouse":
					cfg := database.ClickHouseConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewClickHouse(name.(string), cfg))

				case "postgres":
					cfg := database.PostgresConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewPostgres(name.(string), cfg))

				case "mysql":
					cfg := database.MySQLConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewMySQL(name.(string), cfg))

				case "sqlite3":
					cfg := database.Sqlite3Config{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewSqlite3(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "nosql":
				switch nameModule {
				case "aerospike":
					cfg := noSql.AerospikeConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushNoSql(noSql.NewAerospike(name.(string), cfg))

				case "mongodb":
					cfg := noSql.MongoDBConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushNoSql(noSql.NewMongoDB(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "server":
				switch nameModule {
				case "http":
					cfg := server.HttpConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushServer(server.NewHttp(name.(string), cfg))

				case "kafka":
					cfg := server.KafkaConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushServer(server.NewKafka(name.(string), cfg))

				case "grpc":
					cfg := server.GrpcConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushServer(server.NewGrpc(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "metrics":
				switch nameModule {
				case "prometheus":
					cfg := telemetry.MetricsPrometheusConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewPrometheusMetrics(cfg))

				case "std":
					cfg := telemetry.MetricStdConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewMetricStd(cfg))

				case "otlp_http":
					cfg := telemetry.MetricOtlpHttpConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewMetricOtlpHttp(cfg))

				case "otlp_grpc":
					cfg := telemetry.MetricOtlpGrpcConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewMetricOtlpGrpc(cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "trace":
				switch nameModule {
				case "std":
					cfg := telemetry.TraceStdConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetTrace(telemetry.NewTraceStd(cfg))

				case "zipkin":
					cfg := telemetry.TraceZipkingConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetTrace(telemetry.NewTraceZipking(cfg))

				case "sentry":
					cfg := telemetry.TraceSentryConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetTrace(telemetry.NewTraceSentry(cfg))

				case "otlp_http":
					cfg := telemetry.TraceOtlpHttpConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetTrace(telemetry.NewTraceOtlpHttp(cfg))

				case "otlp_grpc":
					cfg := telemetry.TraceOtlpGrpcConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetTrace(telemetry.NewTraceOtlpGrpc(cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "middleware":
				switch nameModule {
				case "timeout":
					cfg := middleware.TimeoutMiddlewareConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushMiddleware(middleware.NewTimeoutMiddleware(name.(string), cfg))

				case "timing":
					t.app.PushMiddleware(middleware.NewTimingMiddleware(name.(string)))

				case "restore":
					t.app.PushMiddleware(middleware.NewRestoreMiddleware(name.(string)))

				case "jwt":
					cfg := middleware.JWTMiddlewareConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushMiddleware(middleware.NewJWTMiddleware(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "storage":
				switch nameModule {
				case "file":
					cfg := storage.FileConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushStorage(storage.NewFile(name.(string), cfg))

				case "s3":
					cfg := storage.S3Config{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushStorage(storage.NewS3(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			case "mod":
				switch nameModule {
				case "jwt":
					cfg := module.JWTConfig{}
					err = mod.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushModule(module.NewJWT(name.(string), cfg))

				default:
					return fmt.Errorf("unknown mod %s", nameModule)
				}

			default:
				return fmt.Errorf("unknown mod type %s", typeModule)
			}
		}
	}

	return nil
}

func (t cfgItem) Scan(to interface{}) error {
	err := decoder.Decode(t, to, "yaml", true)

	if err != nil {
		return err
	}

	return nil
}

func (t *Auto) Init(app interfaces.IEngine) error {
	t.app = app

	if t.data.Engine == nil {
		t.data.Engine = make(map[string]cfgModule, 10)
	}

	if t.data.Cfg == nil {
		t.data.Cfg = make(map[string]interface{}, 10)
	}

	return t.loadConfig()
}

func (t *Auto) Stop() error {
	return nil
}

func (t *Auto) Has(key string) bool {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		return true
	}

	return false
}

func (t *Auto) GetString(key string, def string) string {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(string)
	}

	return def
}

func (t *Auto) GetBool(key string, def bool) bool {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(bool)
	}

	return def
}

func (t *Auto) GetInt(key string, def int) int {

	val, ok := t.data.Cfg[key]

	if ok {
		return val.(int)
	}

	return def
}

func (t *Auto) GetFloat(key string, def float64) float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(float64)
	}

	return def
}

func (t *Auto) GetMapString(key string, def map[string]string) map[string]string {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]string, len(val.(map[string]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetMapBool(key string, def map[string]bool) map[string]bool {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]bool, len(val.(map[string]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetMapInt(key string, def map[string]int) map[string]int {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]int, len(val.(map[string]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetMapFloat(key string, def map[string]float64) map[string]float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]float64, len(val.(map[string]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetSliceString(key string, def []string) []string {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]string, len(val.([]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetSliceBool(key string, def []bool) []bool {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]bool, len(val.([]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetSliceInt(key string, def []int) []int {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]int, len(val.([]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Auto) GetSliceFloat(key string, def []float64) []float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]float64, len(val.([]interface{})))

			err := decoder.Decode(val, &v, "", false)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}
