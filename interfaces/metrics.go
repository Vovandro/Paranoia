package interfaces

type IMetrics interface {
	Init(app IEngine) error
	Start() error
	Stop() error
}
