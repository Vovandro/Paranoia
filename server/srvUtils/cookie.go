package srvUtils

import "time"

type Cookie struct {
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

func (t *Cookie) Set(key string, value string, path string, expires time.Duration) {
	t.data[key] = cookieItem{
		Name:    key,
		Value:   value,
		Path:    path,
		Expires: expires,
	}
}

func (t *Cookie) Get(key string) string {
	if v, ok := t.data[key]; ok {
		return v.Value
	}

	return ""
}

func (t *Cookie) ToHttp(domain string, sameSite string, httpOnly bool, secure bool) []string {
	res := make([]string, 0, len(t.data))

	for _, v := range t.data {
		res = append(res, v.String(domain, sameSite, httpOnly, secure))
	}

	return res
}

func (t *Cookie) GetAsMap() map[string]string {
	res := make(map[string]string, len(t.data))

	for k, v := range t.data {
		res[k] = v.Value
	}

	return res
}
