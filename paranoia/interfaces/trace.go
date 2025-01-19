package interfaces

type ITrace interface {
	Init(cfg map[string]interface{}) error
	Start() error
	Stop() error
	Name() string
}
