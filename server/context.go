package server

import (
	"net"
	"net/http"
	"sync"
)

type Context struct {
	Request  Request
	Response Response
}

type Request struct {
	Body    []byte
	Headers http.Header
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
				Body:    make([]byte, 0),
				Headers: http.Header{},
			},
			Response: Response{
				Body:       make([]byte, 0),
				StatusCode: 200,
			},
		}
	},
}

func FromHttp(request *http.Request) *Context {
	ctx := ContextPool.Get().(*Context)

	_, err := request.Body.Read(ctx.Request.Body)

	if err != nil {
		return nil
	}

	ctx.Request.Headers = request.Header

	return ctx
}
