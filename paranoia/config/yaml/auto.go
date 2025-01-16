package yaml

import (
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
	"gopkg.in/yaml.v3"
	"os"
)

type Data struct {
	Engine []map[string]interface{} `yaml:"engine"`
	Cfg    map[string]interface{}   `yaml:"cfg"`
}

type Yaml struct {
	cfg  AutoConfig
	data Data
}

func NewYaml(cfg AutoConfig) *Yaml {
	return &Yaml{
		cfg: cfg,
	}
}

type AutoConfig struct {
	FName string `yaml:"filename"`
}

func (t *Yaml) loadConfig() error {
	yamlFile, err := os.ReadFile(t.cfg.FName)
	if err != nil {
		return err
	}

	if t.data.Engine == nil {
		t.data.Engine = make([]map[string]interface{}, 10)
	}

	if t.data.Cfg == nil {
		t.data.Cfg = make(map[string]interface{}, 10)
	}

	err = yaml.Unmarshal(yamlFile, &t.data)

	return err
}

func (t *Yaml) Init(app interfaces.IEngine) error {
	return t.loadConfig()
}

func (t *Yaml) Stop() error {
	return nil
}

func (t *Yaml) Has(key string) bool {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		return true
	}

	return false
}

func (t *Yaml) GetString(key string, def string) string {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(string)
	}

	return def
}

func (t *Yaml) GetBool(key string, def bool) bool {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(bool)
	}

	return def
}

func (t *Yaml) GetInt(key string, def int) int {

	val, ok := t.data.Cfg[key]

	if ok {
		return val.(int)
	}

	return def
}

func (t *Yaml) GetFloat(key string, def float64) float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(float64)
	}

	return def
}

func (t *Yaml) GetMapString(key string, def map[string]string) map[string]string {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]string, len(val.(map[string]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetMapBool(key string, def map[string]bool) map[string]bool {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]bool, len(val.(map[string]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetMapInt(key string, def map[string]int) map[string]int {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]int, len(val.(map[string]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetMapFloat(key string, def map[string]float64) map[string]float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			v := make(map[string]float64, len(val.(map[string]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetMapInterface(key string, def map[string]interface{}) map[string]interface{} {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			return val.(map[string]interface{})
		}
	}

	return def
}

func (t *Yaml) GetSliceString(key string, def []string) []string {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]string, len(val.([]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetSliceBool(key string, def []bool) []bool {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]bool, len(val.([]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetSliceInt(key string, def []int) []int {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]int, len(val.([]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetSliceFloat(key string, def []float64) []float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			v := make([]float64, len(val.([]interface{})))

			err := decode.Decode(val, &v, "", 0)
			if err != nil {
				return def
			}

			return v
		}
	}

	return def
}

func (t *Yaml) GetSliceInterface(key string, def []interface{}) []interface{} {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			return val.([]interface{})
		}
	}

	return def
}

func (t *Yaml) GetConfigItem(typeName string, name string) map[string]interface{} {
	for _, item := range t.data.Engine {
		if _, ok := item["name"]; !ok {
			continue
		}

		if _, ok := item["type"]; !ok {
			continue
		}

		if item["type"] == typeName && item["name"] == name {
			delete(item, "type")
			delete(item, "name")
			return item
		}
	}

	return map[string]interface{}{}
}
