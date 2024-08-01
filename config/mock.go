package config

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Mock struct {
	Data map[string]interface{}
}

func NewMock(data map[string]interface{}) *Mock {
	return &Mock{
		Data: data,
	}
}

func (t *Mock) Init(_ interfaces.IService) error {
	if t.Data == nil {
		t.Data = make(map[string]interface{})
	}

	return nil
}

func (t *Mock) Stop() error {
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

	if ok {
		return val.(string)
	}

	return def
}

func (t *Mock) GetBool(key string, def bool) bool {
	val, ok := t.Data[key]

	if ok {
		return val.(bool)
	}

	return def
}

func (t *Mock) GetInt(key string, def int) int {

	val, ok := t.Data[key]

	if ok {
		return val.(int)
	}

	return def
}

func (t *Mock) GetFloat(key string, def float64) float64 {
	val, ok := t.Data[key]

	if ok {
		return val.(float64)
	}

	return def
}

func (t *Mock) GetMapString(key string, def map[string]string) map[string]string {
	val, ok := t.Data[key]

	if ok {
		return val.(map[string]string)
	}

	return def
}

func (t *Mock) GetMapBool(key string, def map[string]bool) map[string]bool {
	val, ok := t.Data[key]

	if ok {
		return val.(map[string]bool)
	}

	return def
}

func (t *Mock) GetMapInt(key string, def map[string]int) map[string]int {
	val, ok := t.Data[key]

	if ok {
		return val.(map[string]int)
	}

	return def
}

func (t *Mock) GetMapFloat(key string, def map[string]float64) map[string]float64 {
	val, ok := t.Data[key]

	if ok {
		return val.(map[string]float64)
	}

	return def
}

func (t *Mock) GetSliceString(key string, def []string) []string {
	val, ok := t.Data[key]

	if ok {
		return val.([]string)
	}

	return def
}

func (t *Mock) GetSliceBool(key string, def []bool) []bool {
	val, ok := t.Data[key]

	if ok {
		return val.([]bool)
	}

	return def
}

func (t *Mock) GetSliceInt(key string, def []int) []int {
	val, ok := t.Data[key]

	if ok {
		return val.([]int)
	}

	return def
}

func (t *Mock) GetSliceFloat(key string, def []float64) []float64 {
	val, ok := t.Data[key]

	if ok {
		return val.([]float64)
	}

	return def
}
