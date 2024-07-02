package srvUtils

import "gitlab.com/devpro_studio/Paranoia/interfaces"

type HttpResponse struct {
	Body       []byte
	StatusCode int
	headers    interfaces.IHeader
	cookie     interfaces.ICookie
}

func (t *HttpResponse) Clear() {
	t.Body = t.Body[:0]
	t.StatusCode = 200
	t.headers = &HttpHeader{}
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

func (t *HttpResponse) Header() interfaces.IHeader {
	return t.headers
}

func (t *HttpResponse) Cookie() interfaces.ICookie {
	return t.cookie
}
