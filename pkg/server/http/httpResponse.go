package http

type HttpResponse struct {
	Body       []byte
	StatusCode int
	headers    IHeader
	cookie     ICookie
}

func (t *HttpResponse) Clear() {
	t.Body = t.Body[:0]
	t.StatusCode = 200
	t.headers = &HttpHeader{
		make(map[string][]string, 10),
	}
	t.cookie = &HttpCookie{}
}

func (t *HttpResponse) SetBody(data []byte) {
	t.Body = data
}

func (t *HttpResponse) GetBody() []byte {
	return t.Body
}

func (t *HttpResponse) SetStatus(status int) {
	t.StatusCode = status
}

func (t *HttpResponse) GetStatus() int {
	return t.StatusCode
}

func (t *HttpResponse) Header() IHeader {
	return t.headers
}

func (t *HttpResponse) Cookie() ICookie {
	return t.cookie
}
