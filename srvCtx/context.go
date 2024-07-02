package srvCtx

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io"
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
	Headers     map[string]string
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
	ctx.Request.Body, _ = io.ReadAll(request.Body)
	ctx.Request.Headers = request.Header
	ctx.Request.Ip = request.RemoteAddr
	ctx.Request.URI = request.RequestURI
	ctx.Request.Method = request.Method
	ctx.Request.Host = request.Host
	ctx.Request.PostForm = request.PostForm

	ctx.Values = map[string]interface{}{}

	ctx.Response.Body = ctx.Response.Body[:0]
	ctx.Response.StatusCode = 200

	return ctx
}

func FromKafka(msg *kafka.Message) *Ctx {
	ctx := ContextPool.Get().(*Ctx)
	ctx.Request.Body = msg.Value
	ctx.Request.Headers = make(http.Header, len(msg.Headers))

	for _, v := range msg.Headers {
		if val, ok := ctx.Request.Headers[v.Key]; ok {
			ctx.Request.Headers[v.Key] = append(val, string(v.Value))
		} else {
			ctx.Request.Headers[v.Key] = []string{string(v.Value)}
		}
	}

	ctx.Request.Ip = ""
	ctx.Request.URI = *msg.TopicPartition.Topic
	ctx.Request.Method = "KAFKA"
	ctx.Request.Host = ""
	ctx.Request.PostForm = nil

	ctx.Values = map[string]interface{}{}

	ctx.Response.Body = ctx.Response.Body[:0]
	ctx.Response.StatusCode = 200

	return ctx
}
