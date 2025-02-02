package yaml

import (
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
	"gopkg.in/yaml.v3"
	"os"
)

// Data represents the structure of the YAML data.
type Data struct {
	Engine []map[string]interface{} `yaml:"engine"` // Engine configurations.
	Cfg    map[string]interface{}   `yaml:"cfg"`    // General configurations.
}

// Yaml handles the loading and parsing of YAML configuration files.
type Yaml struct {
	cfg  AutoConfig // Configuration for the YAML file.
	data Data       // Parsed data from the YAML file.
}

// New creates a new Yaml instance with the given configuration.
func New(cfg AutoConfig) *Yaml {
	return &Yaml{
		cfg: cfg,
	}
}

// AutoConfig represents the configuration for the YAML file.
type AutoConfig struct {
	FName string `yaml:"filename"` // Filename of the YAML configuration file.
}

// loadConfig reads and parses the YAML configuration file.
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

// Init initializes the YAML configuration by loading the file.
func (t *Yaml) Init(app interfaces.IEngine) error {
	return t.loadConfig()
}

// Stop stops the YAML configuration handler.
func (t *Yaml) Stop() error {
	return nil
}

// Has checks if the given key exists in the configuration.
func (t *Yaml) Has(key string) bool {
	val, ok := t.data.Cfg[key]

	if ok && val != "" {
		return true
	}

	return false
}

// GetString returns the string value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetString(key string, def string) string {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(string)
	}

	return def
}

// GetBool returns the boolean value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetBool(key string, def bool) bool {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(bool)
	}

	return def
}

// GetInt returns the integer value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetInt(key string, def int) int {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(int)
	}

	return def
}

// GetFloat returns the float64 value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetFloat(key string, def float64) float64 {
	val, ok := t.data.Cfg[key]

	if ok {
		return val.(float64)
	}

	return def
}

// GetMapString returns the map[string]string value for the given key, or the default value if the key does not exist.
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

// GetMapBool returns the map[string]bool value for the given key, or the default value if the key does not exist.
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

// GetMapInt returns the map[string]int value for the given key, or the default value if the key does not exist.
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

// GetMapFloat returns the map[string]float64 value for the given key, or the default value if the key does not exist.
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

// GetMapInterface returns the map[string]interface{} value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetMapInterface(key string, def map[string]interface{}) map[string]interface{} {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.(map[string]interface{}); ok2 {
			return val.(map[string]interface{})
		}
	}

	return def
}

// GetSliceString returns the []string value for the given key, or the default value if the key does not exist.
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

// GetSliceBool returns the []bool value for the given key, or the default value if the key does not exist.
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

// GetSliceInt returns the []int value for the given key, or the default value if the key does not exist.
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

// GetSliceFloat returns the []float64 value for the given key, or the default value if the key does not exist.
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

// GetSliceInterface returns the []interface{} value for the given key, or the default value if the key does not exist.
func (t *Yaml) GetSliceInterface(key string, def []interface{}) []interface{} {
	val, ok := t.data.Cfg[key]

	if ok {
		if _, ok2 := val.([]interface{}); ok2 {
			return val.([]interface{})
		}
	}

	return def
}

// GetConfigItem returns the configuration item for the given type and name.
func (t *Yaml) GetConfigItem(typeName string, name string) map[string]interface{} {
	for _, item := range t.data.Engine {
		if _, ok := item["name"]; !ok {
			continue
		}

		if _, ok := item["type"]; !ok {
			continue
		}

		if item["type"] == typeName {
			if name == "" {
				return item
			} else if item["name"] == name {
				res := make(map[string]interface{}, len(item))
				for k, v := range item {
					if k == "type" || k == "name" {
						continue
					}

					res[k] = v
				}

				return res
			}
		}
	}

	return map[string]interface{}{}
}
