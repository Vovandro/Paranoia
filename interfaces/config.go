package interfaces

type IConfig interface {
	Init(app IEngine) error
	Stop() error
	Has(key string) bool
	GetString(key string, def string) string
	GetBool(key string, def bool) bool
	GetInt(key string, def int) int
	GetFloat(key string, def float64) float64

	GetMapString(key string, def map[string]string) map[string]string
	GetMapBool(key string, def map[string]bool) map[string]bool
	GetMapInt(key string, def map[string]int) map[string]int
	GetMapFloat(key string, def map[string]float64) map[string]float64

	GetSliceString(key string, def []string) []string
	GetSliceBool(key string, def []bool) []bool
	GetSliceInt(key string, def []int) []int
	GetSliceFloat(key string, def []float64) []float64
	GetConfigItem(typeName string, name string, dst interface{}) error
}

type IConfigItem interface {
	GetConfigItem(typeName string, name string, dst interface{}) error
}
