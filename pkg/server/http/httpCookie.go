package http

import (
	"net/http"
	"time"
)

type HttpCookie struct {
	data map[string]cookieItem
}

type cookieItem struct {
	Name    string
	Value   string
	Path    string
	Expires time.Duration
}

func (t cookieItem) String(domain string, sameSite string, httpOnly bool, secure bool) string {
	s := t.Name + "=" + t.Value + "; Expires=" + time.Now().Add(t.Expires).String() + "; Path=" + t.Path +
		"; Domain=" + domain + "; SameSite=" + sameSite

	if httpOnly {
		s += "; HttpOnly"
	}

	if secure {
		s += "; Secure"
	}

	return s
}

func (t *HttpCookie) Set(key string, value string, path string, expires time.Duration) {
	t.data[key] = cookieItem{
		Name:    key,
		Value:   value,
		Path:    path,
		Expires: expires,
	}
}

func (t *HttpCookie) Get(key string) string {
	if v, ok := t.data[key]; ok {
		return v.Value
	}

	return ""
}

func (t *HttpCookie) ToHttp(domain string, sameSite string, httpOnly bool, secure bool) []string {
	res := make([]string, 0, len(t.data))

	for _, v := range t.data {
		res = append(res, v.String(domain, sameSite, httpOnly, secure))
	}

	return res
}

func (t *HttpCookie) FromHttp(cookie []*http.Cookie) {
	t.data = make(map[string]cookieItem, len(cookie))
	for _, v := range cookie {
		t.data[v.Name] = cookieItem{
			Name:    v.Name,
			Value:   v.Value,
			Path:    v.Path,
			Expires: v.Expires.Sub(time.Now()),
		}
	}
}

func (t *HttpCookie) GetAsMap() map[string]string {
	res := make(map[string]string, len(t.data))

	for k, v := range t.data {
		res[k] = v.Value
	}

	return res
}
