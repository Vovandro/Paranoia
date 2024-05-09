package server

import (
	"github.com/valyala/fasthttp"
	"net"
	"sync"
)

type Context struct {
	Request  Request
	Response Response
}

type Request struct {
	Body    []byte
	Headers map[string]string
	Ip      net.IP
}

type Response struct {
	Body       []byte
	StatusCode int
}

var ContextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			Request: Request{
				Headers: make(map[string]string, 10),
			},
			Response: Response{
				Body:       make([]byte, 0),
				StatusCode: 200,
			},
		}
	},
}

func FromHttp(request *fasthttp.RequestCtx) *Context {
	ctx := ContextPool.Get().(*Context)

	ctx.Request = Request{
		Body: request.Request.Body(),
	}

	ctx.Request.Headers = make(map[string]string, request.Request.Header.Len())

	request.Request.Header.VisitAll(func(key, value []byte) {
		ctx.Request.Headers[string(key)] = string(value)
	})

	return ctx
}
