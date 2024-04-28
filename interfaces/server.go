package interfaces

import "github.com/valyala/fasthttp"

type IServer interface {
	Init(app IService) error
	Stop() error
	String() string

	Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler
}
