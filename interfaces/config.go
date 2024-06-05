package interfaces

type IConfig interface {
	Init(logger ILogger) error
	Stop() error
	Has(key string) bool
	GetString(key string, def string) string
	GetBool(key string, def bool) bool
	GetInt(key string, def int) int
	GetFloat(key string, def float32) float32
}
