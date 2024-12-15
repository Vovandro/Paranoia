package auto

import (
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/cache/redis"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_cfgItem_Scan(t *testing.T) {
	type args struct {
		to interface{}
	}
	tests := []struct {
		name    string
		t       cfgItem
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"test decode redis",
			cfgItem{
				"hosts":       "127.0.0.1",
				"use_cluster": "true",
				"db_num":      "5",
				"timeout":     "5s",
				"username":    "user",
				"password":    "password",
			},
			args{
				redis.RedisConfig{},
			},
			redis.RedisConfig{
				Hosts:      "127.0.0.1",
				UseCluster: true,
				DBNum:      5,
				Timeout:    5 * time.Second,
				Username:   "user",
				Password:   "password",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Scan(&tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("Scan() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}

func TestConfigEngine_LoadConfig(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		_ = os.WriteFile("./test.yaml", []byte("engine:\n cache:\n  -\n   type: memory\n   name: test_cache\n   time_clear: 10s\n   shard_count: 1"), 0666)

		t := NewAuto(AutoConfig{"./test.yaml"})

		t.app = Paranoia.New("test", nil, nil)

		if err := t.loadConfig(); err != nil {
			t1.Errorf("LoadConfig() error = %v", err)
		}

		if t.app.GetCache("test_cache") == nil {
			t1.Errorf("LoadConfig() error")
		}

		_ = os.Remove("./test.yaml")
	})
}

func TestConfigEngine_LoadType(t1 *testing.T) {
	t1.Run("test base types", func(t1 *testing.T) {
		_ = os.WriteFile("./test2.yaml", []byte("engine:\ncfg:\n s: \"str\"\n b: true\n i: 234\n f: 1.12"), 0666)

		t := NewAuto(AutoConfig{"./test2.yaml"})

		t.app = Paranoia.New("test", nil, nil)

		if err := t.loadConfig(); err != nil {
			t1.Errorf("LoadConfig() error = %v", err)
		}

		if t.GetString("s", "") != "str" {
			t1.Errorf("LoadConfig() error GetString()")
		}

		if t.GetBool("b", false) != true {
			t1.Errorf("LoadConfig() error GetBool()")
		}

		if t.GetInt("i", 0) != 234 {
			t1.Errorf("LoadConfig() error GetInt()")
		}

		if t.GetFloat("f", 0) != 1.12 {
			t1.Errorf("LoadConfig() error GetFloat()")
		}

		_ = os.Remove("./test2.yaml")
	})
}

func TestConfigEngine_LoadMaps(t1 *testing.T) {
	t1.Run("test maps", func(t1 *testing.T) {
		_ = os.WriteFile("./test3.yaml", []byte(`engine:
cfg:
 mi:
  one: 1
  two: 2
 ms:
  one: one
  two: two
 mb:
  one: true
  two: false
 mf:
  one: 1.1
  two: 2.5`), 0666)

		t := NewAuto(AutoConfig{"./test3.yaml"})

		t.app = Paranoia.New("test", nil, nil)

		if err := t.loadConfig(); err != nil {
			t1.Errorf("LoadConfig() error = %v", err)
		}

		if !reflect.DeepEqual(t.GetMapString("ms", nil), map[string]string{"one": "one", "two": "two"}) {
			t1.Errorf("LoadConfig() error GetMapString()")
		}

		if !reflect.DeepEqual(t.GetMapBool("mb", nil), map[string]bool{"one": true, "two": false}) {
			t1.Errorf("LoadConfig() error GetMapBool()")
		}

		if !reflect.DeepEqual(t.GetMapInt("mi", nil), map[string]int{"one": 1, "two": 2}) {
			t1.Errorf("LoadConfig() error GetMapInt()")
		}

		if !reflect.DeepEqual(t.GetMapFloat("mf", nil), map[string]float64{"one": 1.1, "two": 2.5}) {
			t1.Errorf("LoadConfig() error GetMapFloat()")
		}

		_ = os.Remove("./test3.yaml")
	})
}

func TestConfigEngine_LoadSlice(t1 *testing.T) {
	t1.Run("test slice", func(t1 *testing.T) {
		_ = os.WriteFile("./test4.yaml", []byte(`engine:
cfg:
 i:
  - 1
  - 2
 s:
  - one
  - two
 b:
  - true
  - false
 f:
  - 1.1
  - 2.5`), 0666)

		t := NewAuto(AutoConfig{"./test4.yaml"})

		t.app = Paranoia.New("test", nil, nil)

		if err := t.loadConfig(); err != nil {
			t1.Errorf("LoadConfig() error = %v", err)
		}

		if !reflect.DeepEqual(t.GetSliceString("s", nil), []string{"one", "two"}) {
			t1.Errorf("LoadConfig() error GetSliceString()")
		}

		if !reflect.DeepEqual(t.GetSliceBool("b", nil), []bool{true, false}) {
			t1.Errorf("LoadConfig() error GetSliceBool()")
		}

		if !reflect.DeepEqual(t.GetSliceInt("i", nil), []int{1, 2}) {
			t1.Errorf("LoadConfig() error GetSliceInt()")
		}

		if !reflect.DeepEqual(t.GetSliceFloat("f", nil), []float64{1.1, 2.5}) {
			t1.Errorf("LoadConfig() error GetSliceFloat()")
		}

		_ = os.Remove("./test4.yaml")
	})
}
