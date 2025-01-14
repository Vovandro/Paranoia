package interfaces

type IMetrics interface {
	Init(cfg map[string]interface{}) error
	Start() error
	Stop() error
	Name() string
}
