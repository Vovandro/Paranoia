package Paranoia

type IServer interface {
	Init(app *Service) error
	Stop() error
	String() string

	Handle(ctx *Ctx) error
}

type Ctx struct {
	Header map[string]string
	Body   string
}
