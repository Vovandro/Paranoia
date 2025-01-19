package http

import (
	"io"
	"net/http"
)

type HttpRequest struct {
	request *http.Request
	cookies ICookie
	headers IHeader
}

func (t *HttpRequest) Fill(request *http.Request) {
	t.request = request

	if t.cookies == nil {
		t.cookies = &HttpCookie{}
	}

	t.cookies.(*HttpCookie).FromHttp(t.request.Cookies())

	if t.headers == nil {
		t.headers = &HttpHeader{}
	}

	t.headers.(*HttpHeader).Header = t.request.Header
}

func (t *HttpRequest) GetBody() io.ReadCloser {
	return t.request.Body
}

func (t *HttpRequest) GetBodySize() int64 {
	return t.request.ContentLength
}

func (t *HttpRequest) GetCookie() ICookie {
	return t.cookies
}

func (t *HttpRequest) GetHeader() IHeader {
	return t.headers
}

func (t *HttpRequest) GetMethod() string {
	return t.request.Method
}

func (t *HttpRequest) GetURI() string {
	return t.request.RequestURI
}

func (t *HttpRequest) GetQuery() IQuery {
	if t.request.Form != nil {
		return t.request.Form
	}

	return t.request.URL.Query()
}

func (t *HttpRequest) GetRemoteIP() string {
	return t.request.RemoteAddr
}

func (t *HttpRequest) GetRemoteHost() string {
	return t.request.Host
}

func (t *HttpRequest) GetUserAgent() string {
	return t.request.Header.Get("User-Agent")
}
