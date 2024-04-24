package interfaces

type IServer interface {
	Init(app IService) error
	Stop() error
	String() string

	Handle(ctx *Ctx) error
}

type Ctx struct {
	Header map[string]string
	Body   string
}
