package kafka

import (
	"context"
	"io"
	"time"
)

// IKafka defines the interface for Kafka server operations
type IKafka interface {
	// PushRoute adds a new route to the Kafka server
	PushRoute(path string, handler RouteFunc, middlewares []string)
}

type IHeader interface {
	Add(key, value string)
	Set(key string, value string)
	Get(key string) string
	Values(key string) []string
	Del(key string)
	GetAsMap() map[string][]string
}

type IQuery interface {
	Get(key string) string
}

type ICookie interface {
	Set(key string, value string, path string, expires time.Duration)
	Get(key string) string
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
	GetUserAgent() string
}

type IResponse interface {
	Clear()
	SetBody(data []byte)
	GetBody() []byte
	SetStatus(status int)
	GetStatus() int

	Header() IHeader
	Cookie() ICookie
}

type ICtx interface {
	GetRequest() IRequest
	GetResponse() IResponse

	GetUserValue(key string) (interface{}, error)
	PushUserValue(key string, val interface{})

	GetRouterValue(key string) string
	SetRouteProps(values map[string]string)
}

type RouteFunc func(c context.Context, ctx ICtx)

type IMiddleware interface {
	Invoke(next RouteFunc) RouteFunc
}
