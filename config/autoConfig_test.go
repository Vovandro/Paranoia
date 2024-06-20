package config

import (
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/cache"
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
				cache.RedisConfig{},
			},
			cache.RedisConfig{
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
		_ = os.WriteFile("./test.yaml", []byte("engine:\n cache:\n  -\n   type: memory\n   name: test_cache\n   time_clear: 10s"), 0666)

		t := NewAutoConfig("./test.yaml")

		t.app = Paranoia.New("test", nil, nil, nil)

		if err := t.loadConfig(); err != nil {
			t1.Errorf("LoadConfig() error = %v", err)
		}

		if t.app.GetCache("test_cache") == nil {
			t1.Errorf("LoadConfig() error")
		}

		_ = os.Remove("./test.yaml")
	})
}
