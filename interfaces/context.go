package interfaces

import (
	"context"
	"io"
	"time"
)

type IHeader interface {
	Add(key, value string)
	Set(key string, value string)
	Get(key string) string
	Values(key string) []string
	Del(key string)
	GetAsMap() map[string]string
}

type IQuery interface {
	Get(key string) (string, error)
}

type ICookie interface {
	Set(key string, value string, path string, expires time.Duration)
	Get(key string) string
	ToHttp(domain string, sameSite string, httpOnly bool, secure bool) []string
	GetAsMap() map[string]string
}

type IRequest interface {
	GetBody() io.ReadCloser
	GetBodySize() int64
	GetCookie() ICookie

	GetHeader() IHeader
	GetMethod() string
	GetURI() string
	GetQuery() IQuery

	GetRemoteIP() string
	GetRemoteHost() string
}

type IResponse interface {
	GetBody() []byte
	GetStatus() int
	GetHeader() IHeader
	GetCookie() ICookie
}

type ICtx interface {
	GetRequest() IRequest
	GetResponse() IResponse

	Done() <-chan struct{}
	StartTimeout(t time.Duration) context.CancelFunc

	GetUserValue(key string) (interface{}, error)
	PushUserValue(key string, val interface{})
}
