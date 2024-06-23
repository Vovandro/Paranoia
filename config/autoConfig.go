package config

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/cache"
	"gitlab.com/devpro_studio/Paranoia/client"
	"gitlab.com/devpro_studio/Paranoia/database"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/noSql"
	"gitlab.com/devpro_studio/Paranoia/server"
	"gitlab.com/devpro_studio/Paranoia/telemetry"
	"gitlab.com/devpro_studio/Paranoia/utils/decoder"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type cfgItem map[string]string
type cfgModule []cfgItem

type Data struct {
	Engine map[string]cfgModule `yaml:"engine"`
	Cfg    map[string]string    `yaml:"cfg"`
}

type AutoConfig struct {
	FName string

	app  interfaces.IService
	data Data
}

func NewAutoConfig(fName string) *AutoConfig {
	return &AutoConfig{
		FName: fName,
	}
}

func (t *AutoConfig) loadConfig() error {
	yamlFile, err := os.ReadFile(t.FName)
	if err != nil {
		return err
	}

	if t.data.Engine == nil {
		t.data.Engine = make(map[string]cfgModule, 10)
	}

	if t.data.Cfg == nil {
		t.data.Cfg = make(map[string]string, 10)
	}

	err = yaml.Unmarshal(yamlFile, &t.data)

	if err != nil {
		return err
	}

	if v, ok := t.data.Cfg["logLevel"]; ok {
		t.app.GetLogger().SetLevel(interfaces.GetLogLevel(v))
	}

	for typeModule, modules := range t.data.Engine {
		for _, module := range modules {
			name, ok := module["name"]

			if !ok {
				return fmt.Errorf("not found name module")
			}

			nameModule, ok := module["type"]

			if !ok {
				return fmt.Errorf("not found type module %s", name)
			}

			delete(module, "name")
			delete(module, "type")

			switch typeModule {
			case "cache":
				switch nameModule {
				case "memcached":
					cfg := cache.MemcachedConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewMemcached(name, cfg))

				case "memory":
					cfg := cache.MemoryConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewMemory(name, cfg))

				case "redis":
					cfg := cache.RedisConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushCache(cache.NewRedis(name, cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			case "client":
				switch nameModule {
				case "http":
					cfg := client.HTTPClientConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushClient(client.NewHTTPClient(name, cfg))

				case "kafka":
					cfg := client.KafkaClientConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushClient(client.NewKafkaClient(name, cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			case "database":
				switch nameModule {
				case "clickhouse":
					cfg := database.ClickHouseConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewClickHouse(name, cfg))

				case "postgres":
					cfg := database.PostgresConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewPostgres(name, cfg))

				case "mysql":
					cfg := database.MySQLConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewMySQL(name, cfg))

				case "sqlite3":
					cfg := database.Sqlite3Config{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushDatabase(database.NewSqlite3(name, cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			case "nosql":
				switch nameModule {
				case "aerospike":
					cfg := noSql.AerospikeConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushNoSql(noSql.NewAerospike(name, cfg))

				case "mongodb":
					cfg := noSql.MongoDBConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushNoSql(noSql.NewMongoDB(name, cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			case "server":
				switch nameModule {
				case "http":
					cfg := server.HttpConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushServer(server.NewHttp(name, cfg))

				case "kafka":
					cfg := server.KafkaConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.PushServer(server.NewKafka(name, cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			case "metrics":
				switch nameModule {
				case "prometheus":
					cfg := telemetry.MetricsPrometheusConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewPrometheusMetrics(cfg))

				case "std":
					cfg := telemetry.MetricStdConfig{}
					err = module.Scan(&cfg)
					if err != nil {
						return err
					}

					t.app.SetMetrics(telemetry.NewMetricStd(cfg))

				default:
					return fmt.Errorf("unknown module %s", nameModule)
				}

			default:
				return fmt.Errorf("unknown module type %s", typeModule)
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

func (t *AutoConfig) Init(app interfaces.IService) error {
	t.app = app

	if t.data.Engine == nil {
		t.data.Engine = make(map[string]cfgModule, 10)
	}

	if t.data.Cfg == nil {
		t.data.Cfg = make(map[string]string, 10)
	}

	return t.loadConfig()
}

func (t *AutoConfig) Stop() error {
	return nil
}

func (t *AutoConfig) Has(key string) bool {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		return true
	}

	return false
}

func (t *AutoConfig) GetString(key string, def string) string {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		return val
	}

	return def
}

func (t *AutoConfig) GetBool(key string, def bool) bool {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		b, err := strconv.ParseBool(val)

		if err == nil {
			return b
		}
	}

	return def
}

func (t *AutoConfig) GetInt(key string, def int) int {

	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		i, err := strconv.ParseInt(val, 10, 32)

		if err == nil {
			return int(i)
		}
	}

	return def
}

func (t *AutoConfig) GetFloat(key string, def float32) float32 {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 32)

		if err == nil {
			return float32(i)
		}
	}

	return def
}
