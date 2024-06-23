package interfaces

type IMetrics interface {
	Init(app IService) error
	Start() error
	Stop() error
}
