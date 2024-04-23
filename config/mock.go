package config

import (
	"goServer"
	"strconv"
)

type Mock struct {
	Data map[string]string
}

func (t *Mock) Init(app *goServer.Service) error {
	if t.Data == nil {
		t.Data = make(map[string]string)
	}

	return nil
}

func (t *Mock) Has(key string) bool {
	val, ok := t.Data[key]

	if ok && val != "" {
		return true
	}

	return false
}

func (t *Mock) GetString(key string, def string) string {
	val, ok := t.Data[key]

	if ok && val != "" {
		return val
	}

	return def
}

func (t *Mock) GetBool(key string, def bool) bool {
	val, ok := t.Data[key]

	if ok && val != "" {
		b, err := strconv.ParseBool(val)

		if err == nil {
			return b
		}
	}

	return def
}

func (t *Mock) GetInt(key string, def int) int {

	val, ok := t.Data[key]

	if ok && val != "" {
		i, err := strconv.ParseInt(val, 10, 32)

		if err == nil {
			return int(i)
		}
	}

	return def
}

func (t *Mock) GetFloat(key string, def float32) float32 {
	val, ok := t.Data[key]

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 32)

		if err == nil {
			return float32(i)
		}
	}

	return def
}
