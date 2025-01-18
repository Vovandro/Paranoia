package kafka_client

import (
	"fmt"
	"io"
)

type Response struct {
	Body       io.Reader
	Header     map[string][]string
	Err        error
	RetryCount int
	Code       int
}

func (t *Response) GetBody() ([]byte, error) {
	if t.Body == nil {
		return []byte{}, fmt.Errorf("body cannot be nil")
	}

	return io.ReadAll(t.Body)
}

func (t *Response) GetLazyBody() io.Reader {
	return t.Body
}

func (t *Response) GetHeader() map[string][]string {
	return t.Header
}

func (t *Response) Error() error {
	return t.Err
}

func (t *Response) GetRetries() int {
	return t.RetryCount
}

func (t *Response) GetCode() int {
	return t.Code
}
