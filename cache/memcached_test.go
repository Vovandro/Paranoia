package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"os"
	"testing"
)

func TestMemcached_Has(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Memcached{
		Name:  "test",
		Hosts: host + ":11211",
	}
	err := t.Init(nil)
	defer t.Stop()

	if err != nil {
		t1.Errorf("Init() = %v", err)
	}

	tests := []struct {
		name  string
		store map[string]string
		key   string
		want  bool
	}{
		{
			"test exists",
			map[string]string{"k1": "v1"},
			"k1",
			true,
		},
		{
			"test does not exist",
			map[string]string{"k1": "v1"},
			"k2",
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for k, v := range tt.store {
				t.client.Add(&memcache.Item{
					Key:   k,
					Value: []byte(v),
				})
			}

			if got := t.Has(tt.key); got != tt.want {
				t1.Errorf("Has() = %v, want %v", got, tt.want)
			}

			for k, _ := range tt.store {
				t.client.Delete(k)
			}
		})
	}
}
