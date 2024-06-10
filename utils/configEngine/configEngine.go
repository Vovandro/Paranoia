package configEngine

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/cache"
	"gitlab.com/devpro_studio/Paranoia/client"
	"gitlab.com/devpro_studio/Paranoia/database"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/noSql"
	"gitlab.com/devpro_studio/Paranoia/server"
	"gitlab.com/devpro_studio/Paranoia/utils/decoder"
	"gopkg.in/yaml.v3"
	"os"
)

type cfgItem map[string]string
type cfgModule []cfgItem
type cfgType map[string]cfgModule

type ConfigEngine struct {
	FName string
	data  map[string]cfgType `yaml:",inline"`
}

func NewConfigEngine(fName string) *ConfigEngine {
	return &ConfigEngine{
		FName: fName,
	}
}

func (t *ConfigEngine) LoadConfig(app interfaces.IService) error {
	yamlFile, err := os.ReadFile(t.FName)
	if err != nil {
		return err
	}

	if t.data == nil {
		t.data = make(map[string]cfgType)
	}

	err = yaml.Unmarshal(yamlFile, &t.data)

	if err != nil {
		return err
	}

	if engine, ok := t.data["engine"]; ok {
		for typeModule, modules := range engine {

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

						app.PushCache(cache.NewMemcached(name, cfg))

					case "memory":
						cfg := cache.MemoryConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushCache(cache.NewMemory(name, cfg))

					case "redis":
						cfg := cache.RedisConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushCache(cache.NewRedis(name, cfg))

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

						app.PushClient(client.NewHTTPClient(name, cfg))

					case "kafka":
						cfg := client.KafkaClientConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushClient(client.NewKafkaClient(name, cfg))

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

						app.PushDatabase(database.NewClickHouse(name, cfg))

					case "postgres":
						cfg := database.PostgresConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushDatabase(database.NewPostgres(name, cfg))

					case "mysql":
						cfg := database.MySQLConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushDatabase(database.NewMySQL(name, cfg))

					case "sqlite3":
						cfg := database.Sqlite3Config{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushDatabase(database.NewSqlite3(name, cfg))

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

						app.PushNoSql(noSql.NewAerospike(name, cfg))

					case "mongodb":
						cfg := noSql.MongoDBConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushNoSql(noSql.NewMongoDB(name, cfg))

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

						app.PushServer(server.NewHttp(name, cfg))

					case "kafka":
						cfg := server.KafkaConfig{}
						err = module.Scan(&cfg)
						if err != nil {
							return err
						}

						app.PushServer(server.NewKafka(name, cfg))

					default:
						return fmt.Errorf("unknown module %s", nameModule)
					}

				default:
					return fmt.Errorf("unknown module type %s", typeModule)
				}
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
