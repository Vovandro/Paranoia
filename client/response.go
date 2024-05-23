package client

type Response struct {
	Body       []byte
	Header     map[string][]string
	Err        error
	RetryCount int
}

func (t *Response) GetBody() []byte {
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
