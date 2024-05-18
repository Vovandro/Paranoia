package srvCtx

import (
	"net/http"
	"net/url"
	"sync"
)

type Ctx struct {
	Request  Request
	Response Response
	Values   map[string]interface{}
}

type Request struct {
	Body     []byte
	Headers  http.Header
	Ip       string
	URI      string
	Method   string
	Host     string
	PostForm url.Values
}

type Response struct {
	Body        []byte
	ContentType string
	StatusCode  int
	Headers     http.Header
}

var ContextPool = sync.Pool{
	New: func() interface{} {
		return &Ctx{
			Request: Request{
				Body:    make([]byte, 0),
				Headers: http.Header{},
			},
			Response: Response{
				Body:        make([]byte, 0),
				StatusCode:  200,
				ContentType: "application/json; charset=utf-8",
			},
			Values: make(map[string]interface{}, 10),
		}
	},
}

func FromHttp(request *http.Request) *Ctx {
	ctx := ContextPool.Get().(*Ctx)

	request.Body.Read(ctx.Request.Body)
	ctx.Request.Headers = request.Header
	ctx.Request.Ip = request.RemoteAddr
	ctx.Request.URI = request.RequestURI
	ctx.Request.Method = request.Method
	ctx.Request.Host = request.Host
	ctx.Request.PostForm = request.PostForm

	ctx.Values = map[string]interface{}{}

	ctx.Response.Body = ctx.Response.Body[:0]

	return ctx
}
