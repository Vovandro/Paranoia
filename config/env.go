package config

import (
	"bytes"
	"goServer"
	"os"
	"strconv"
)

type Env struct {
	FName string
	data  map[string]string
	app   *goServer.Service
}

func (t *Env) Init(app *goServer.Service) error {
	t.data = make(map[string]string, 20)
	t.app = app

	if t.FName == "" {
		t.FName = ".env"
	}

	f, err := os.ReadFile(t.FName)

	if err != nil || len(f) == 0 {
		t.app.GetLogger().Error(err)
	} else {
		t.ParseFile(f)
	}

	return nil
}

func (t *Env) ParseFile(data []byte) {
	key := make([]byte, 0, 20)
	val := make([]byte, 0, 20)
	var isKey bool

	rows := bytes.Split(data, []byte("\n"))

	for _, row := range rows {
		key = key[:0]
		val = key[:0]
		isKey = true

		for _, c := range row {
			if c == '#' {
				break
			}

			if isKey {
				if c == '=' {
					isKey = false
					continue
				}

				key = append(val, c)
			} else {
				val = append(val, c)
			}
		}

		if !isKey && len(key) != 0 && len(val) != 0 {
			if val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}

			t.data[string(key)] = string(val)
		}
	}
}

func (t *Env) Has(key string) bool {
	val, ok := t.data[key]

	if ok && val != "" {
		return true
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		t.data[key] = val
		return true
	}

	return false
}

func (t *Env) GetString(key string, def string) string {
	val, ok := t.data[key]

	if ok && val != "" {
		return val
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		t.data[key] = val
		return val
	}

	return def
}

func (t *Env) GetBool(key string, def bool) bool {
	val, ok := t.data[key]

	if ok && val != "" {
		b, err := strconv.ParseBool(val)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			return b
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		b, err := strconv.ParseBool(val)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			t.data[key] = val
			return b
		}
	}

	return def
}

func (t *Env) GetInt(key string, def int) int {

	val, ok := t.data[key]

	if ok && val != "" {
		i, err := strconv.ParseInt(val, 10, 32)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			return int(i)
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		i, err := strconv.ParseInt(val, 10, 32)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			t.data[key] = val
			return int(i)
		}
	}

	return def
}

func (t *Env) GetFloat(key string, def float32) float32 {
	val, ok := t.data[key]

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 32)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			return float32(i)
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 32)

		if err != nil {
			t.app.GetLogger().Error(err)
		} else {
			t.data[key] = val
			return float32(i)
		}
	}

	return def
}
