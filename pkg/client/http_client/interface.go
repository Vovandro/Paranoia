package http_client

import (
	"context"
	"io"
)

type IHTTPClient interface {
	// Fetch sends an HTTP request and returns a channel with the response
	Fetch(ctx context.Context, method string, host string, data []byte, headers map[string][]string) chan IClientResponse
}

type IClientResponse interface {
	// GetBody retrieves the response body as a byte slice. Returns an error if the body cannot be read.
	GetBody() ([]byte, error)

	// GetLazyBody returns the response body as an io.Reader for lazy reading.
	GetLazyBody() io.Reader

	// GetHeader retrieves the headers of the response as a map.
	GetHeader() map[string][]string

	// Error retrieves the error encountered during the request, if any.
	Error() error

	// GetRetries retrieves the number of retry attempts made for the request.
	GetRetries() int

	// GetCode retrieves the HTTP status code of the response.
	GetCode() int
}
