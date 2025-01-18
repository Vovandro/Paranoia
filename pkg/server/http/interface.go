package http

import (
	"context"
	"io"
	"time"
)

// IHttp defines the interface for HTTP server operations
type IHttp interface {
	// PushRoute adds a new route to the HTTP server
	PushRoute(method string, path string, handler RouteFunc, middlewares []string)
}

// IHeader defines the interface for HTTP headers
type IHeader interface {
	// Add adds a value to the header key
	Add(key, value string)

	// Set sets the header key to the specified value
	Set(key string, value string)

	// Get retrieves the first value associated with the given key
	Get(key string) string

	// Values retrieves all values associated with the given key
	Values(key string) []string

	// Del deletes the values associated with the given key
	Del(key string)

	// GetAsMap returns the header as a map
	GetAsMap() map[string][]string
}

// IQuery defines the interface for HTTP query parameters
type IQuery interface {
	// Get retrieves the value associated with the given key
	Get(key string) string
}

// ICookie defines the interface for HTTP cookies
type ICookie interface {
	// Set sets a cookie with the specified key, value, path, and expiration duration
	Set(key string, value string, path string, expires time.Duration)

	// Get retrieves the value of the cookie associated with the given key
	Get(key string) string

	// GetAsMap returns the cookies as a map
	GetAsMap() map[string]string
}

// IRequest defines the interface for an HTTP request
type IRequest interface {
	// GetBody returns the body of the request
	GetBody() io.ReadCloser

	// GetBodySize returns the size of the request body
	GetBodySize() int64

	// GetCookie returns the cookies of the request
	GetCookie() ICookie

	// GetHeader returns the headers of the request
	GetHeader() IHeader

	// GetMethod returns the HTTP method of the request
	GetMethod() string

	// GetURI returns the URI of the request
	GetURI() string

	// GetQuery returns the query parameters of the request
	GetQuery() IQuery

	// GetRemoteIP returns the remote IP address of the request
	GetRemoteIP() string

	// GetRemoteHost returns the remote host of the request
	GetRemoteHost() string

	// GetUserAgent returns the user agent of the request
	GetUserAgent() string
}

// IResponse defines the interface for an HTTP response
type IResponse interface {
	// Clear clears the response
	Clear()

	// SetBody sets the body of the response
	SetBody(data []byte)

	// GetBody returns the body of the response
	GetBody() []byte

	// SetStatus sets the status code of the response
	SetStatus(status int)

	// GetStatus returns the status code of the response
	GetStatus() int

	// Header returns the headers of the response
	Header() IHeader

	// Cookie returns the cookies of the response
	Cookie() ICookie
}

// ICtx defines the interface for an HTTP context
type ICtx interface {
	// GetRequest returns the request of the context
	GetRequest() IRequest

	// GetResponse returns the response of the context
	GetResponse() IResponse

	// GetUserValue retrieves a user value associated with the given key
	GetUserValue(key string) (interface{}, error)

	// PushUserValue associates a user value with the given key
	PushUserValue(key string, val interface{})

	// GetRouterValue retrieves a router value associated with the given key
	GetRouterValue(key string) string

	// SetRouteProps sets the router properties
	SetRouteProps(values map[string]string)
}

// RouteFunc defines the function type for a route
type RouteFunc func(c context.Context, ctx ICtx)

// IMiddleware defines the interface for middleware
type IMiddleware interface {
	// Invoke invokes the middleware with the next route function
	Invoke(next RouteFunc) RouteFunc
}
