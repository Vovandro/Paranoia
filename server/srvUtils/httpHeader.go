package srvUtils

import "net/http"

type HttpHeader struct {
	http.Header
}

func (t *HttpHeader) GetAsMap() map[string][]string {
	return t.Header
}
