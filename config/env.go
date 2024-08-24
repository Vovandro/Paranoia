package config

import (
	"bytes"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	Config EnvConfig
	data   map[string]string
}

type EnvConfig struct {
	FName              string `yaml:"filename"`
	ValueItemDelimiter string `yaml:"value_delimiter"`
	ItemDelimiter      string `yaml:"item_delimiter"`
}

func NewEnv(cfg EnvConfig) *Env {
	return &Env{
		Config: cfg,
	}
}

func (t *Env) Init(app interfaces.IEngine) error {
	t.data = make(map[string]string, 20)

	if t.Config.FName == "" {
		t.Config.FName = ".env"
	}

	f, err := os.ReadFile(t.Config.FName)

	if err != nil || len(f) == 0 {
		if err == nil {
			err = fmt.Errorf("file %s is empty", t.Config.FName)
		}
	} else {
		t.ParseFile(f)
	}

	return nil
}

func (t *Env) Stop() error {
	return nil
}

func (t *Env) ParseFile(data []byte) {
	key := make([]byte, 0, 20)
	val := make([]byte, 0, 20)
	var isKey bool

	rows := bytes.Split(data, []byte("\n"))

	for _, row := range rows {
		if len(row) <= 2 {
			continue
		}

		key = key[:0]
		val = val[:0]
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

				key = append(key, c)
			} else {
				val = append(val, c)
			}
		}

		if !isKey && len(key) != 0 && len(val) != 0 {
			val = bytes.Trim(val, " \t")

			if val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}

			t.data[strings.Trim(string(key), " \t")] = string(val)
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
			fmt.Println(err)
		} else {
			return b
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		b, err := strconv.ParseBool(val)

		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		} else {
			return int(i)
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		i, err := strconv.ParseInt(val, 10, 32)

		if err != nil {
			fmt.Println(err)
		} else {
			t.data[key] = val
			return int(i)
		}
	}

	return def
}

func (t *Env) GetFloat(key string, def float64) float64 {
	val, ok := t.data[key]

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println(err)
		} else {
			return i
		}
	}

	val, ok = os.LookupEnv(key)

	if ok && val != "" {
		i, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println(err)
		} else {
			t.data[key] = val
			return i
		}
	}

	return def
}

func (t *Env) GetMapString(key string, def map[string]string) map[string]string {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make(map[string]string, len(values))

	for _, value := range values {
		v := strings.Split(value, t.Config.ValueItemDelimiter)
		if len(v) != 2 {
			fmt.Println(fmt.Errorf("invalid map row %s", value))
			continue
		}

		res[v[0]] = v[1]
	}

	return res
}

func (t *Env) GetMapBool(key string, def map[string]bool) map[string]bool {

	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make(map[string]bool, len(values))

	for _, value := range values {
		v := strings.Split(value, t.Config.ValueItemDelimiter)
		if len(v) != 2 {
			fmt.Println(fmt.Errorf("invalid map row %s", value))
			continue
		}

		res[v[0]], _ = strconv.ParseBool(v[1])
	}

	return res
}

func (t *Env) GetMapInt(key string, def map[string]int) map[string]int {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make(map[string]int, len(values))

	for _, value := range values {
		v := strings.Split(value, t.Config.ValueItemDelimiter)
		if len(v) != 2 {
			fmt.Println(fmt.Errorf("invalid map row %s", value))
			continue
		}

		i, _ := strconv.ParseInt(v[1], 10, 32)
		res[v[0]] = int(i)
	}

	return res
}

func (t *Env) GetMapFloat(key string, def map[string]float64) map[string]float64 {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make(map[string]float64, len(values))

	for _, value := range values {
		v := strings.Split(value, t.Config.ValueItemDelimiter)
		if len(v) != 2 {
			fmt.Println(fmt.Errorf("invalid map row %s", value))
			continue
		}

		res[v[0]], _ = strconv.ParseFloat(v[1], 64)
	}

	return res
}

func (t *Env) GetSliceString(key string, def []string) []string {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	return strings.Split(val, t.Config.ItemDelimiter)
}

func (t *Env) GetSliceBool(key string, def []bool) []bool {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make([]bool, 0, len(values))

	for _, value := range values {
		i, _ := strconv.ParseBool(value)
		res = append(res, i)
	}

	return res
}

func (t *Env) GetSliceInt(key string, def []int) []int {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make([]int, 0, len(values))

	for _, value := range values {
		i, _ := strconv.ParseInt(value, 10, 32)
		res = append(res, int(i))
	}

	return res
}

func (t *Env) GetSliceFloat(key string, def []float64) []float64 {
	val, ok := t.data[key]

	if !ok || val == "" {
		return def
	}

	values := strings.Split(val, t.Config.ItemDelimiter)
	res := make([]float64, 0, len(values))

	for _, value := range values {
		i, _ := strconv.ParseFloat(value, 64)
		res = append(res, i)
	}

	return res
}
