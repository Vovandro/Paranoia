package interfaces

type ITrace interface {
	Init(app IEngine) error
	Start() error
	Stop() error
}
