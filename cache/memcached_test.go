package cache

import (
	"bytes"
	"github.com/bradfitz/gomemcache/memcache"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMemcached_Has(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

func TestMemcached_Base(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		val     any
		timeout time.Duration
	}

	tests := []struct {
		name     string
		store    []item
		sleep    time.Duration
		keyCheck string
		want     []byte
	}{
		{
			"base test set",
			[]item{
				{
					"k1",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			[]byte("v1"),
		},
		{
			"test not exists",
			[]item{
				{
					"k1",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k2",
			nil,
		},
		{
			"base test timeout",
			[]item{
				{
					"k1",
					"v1",
					time.Second,
				},
			},
			time.Second * 2,
			"k1",
			nil,
		},
		{
			"test byte",
			[]item{
				{
					"k1",
					[]byte("v1"),
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			[]byte("v1"),
		},
		{
			"test int",
			[]item{
				{
					"k1",
					2,
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			[]byte("2"),
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.Set(v.key, v.val, v.timeout)
			}

			time.Sleep(tt.sleep)

			got, err := t.Get(tt.keyCheck)

			if err != nil && tt.want != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if (got == nil && tt.want != nil) || (got != nil && !bytes.Equal(got.([]byte), tt.want)) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_In(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		key2    string
		val     any
		timeout time.Duration
	}

	tests := []struct {
		name      string
		store     []item
		sleep     time.Duration
		keyCheck  string
		key2Check string
		want      any
	}{
		{
			"base test set",
			[]item{
				{
					"k1",
					"k2",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k2",
			"v1",
		},
		{
			"base test multiple set one",
			[]item{
				{
					"k1",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k1",
					"k3",
					"v2",
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k2",
			"v1",
		},
		{
			"base test multiple set two",
			[]item{
				{
					"k1",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k1",
					"k3",
					"v2",
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k3",
			"v2",
		},
		{
			"test integer",
			[]item{
				{
					"k1",
					"k2",
					3,
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k2",
			float64(3),
		},
		{
			"test multiple int set one",
			[]item{
				{
					"k1",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k1",
					"k3",
					3,
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k2",
			"v1",
		},
		{
			"test multiple int set two",
			[]item{
				{
					"k1",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k1",
					"k3",
					3,
					time.Minute,
				},
			},
			time.Microsecond,
			"k1",
			"k3",
			float64(3),
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.SetIn(v.key, v.key2, v.val, v.timeout)
			}

			time.Sleep(tt.sleep)

			got, err := t.GetIn(tt.keyCheck, tt.key2Check)

			if err != nil && tt.want != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_Map(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		val     map[string]interface{}
		timeout time.Duration
	}

	tests := []struct {
		name     string
		store    []item
		keyCheck string
		want     map[string]interface{}
	}{
		{
			"base test set",
			[]item{
				{
					"k1",
					map[string]interface{}{"k": "v"},
					time.Minute,
				},
			},
			"k1",
			map[string]interface{}{"k": "v"},
		},
		{
			"base test int",
			[]item{
				{
					"k1",
					map[string]interface{}{"k": "v", "k2": 5},
					time.Minute,
				},
			},
			"k1",
			map[string]interface{}{"k": "v", "k2": 5.0},
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.SetMap(v.key, v.val, v.timeout)
			}

			got, err := t.GetMap(tt.keyCheck)

			if err != nil && tt.want != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_Increment(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		val     int64
		timeout time.Duration
	}

	tests := []struct {
		name     string
		store    []item
		keyCheck string
		want     int64
	}{
		{
			"base test increment",
			[]item{
				{
					"k1",
					1,
					time.Minute,
				},
			},
			"k1",
			1,
		},
		{
			"test multiple increment",
			[]item{
				{
					"k1",
					1,
					time.Minute,
				},
				{
					"k1",
					5,
					time.Minute,
				},
			},
			"k1",
			6,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				lastVal, _ = t.Increment(v.key, v.val, v.timeout)
			}

			got, err := t.Get(tt.keyCheck)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			a, err := strconv.ParseInt(string(got.([]byte)), 10, 64)

			if a != tt.want || lastVal != tt.want {
				t1.Errorf("Check = %v, last = %v, want %v", a, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_Decrement(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		val     int64
		timeout time.Duration
	}

	tests := []struct {
		name     string
		store    []item
		dec      []item
		keyCheck string
		want     int64
		wantErr  bool
	}{
		{
			"base test decrement",
			[]item{
				{
					"k1",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k1",
					1,
					time.Minute,
				},
			},
			"k1",
			9,
			false,
		},
		{
			"test multiple decrement",
			[]item{
				{
					"k1",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k1",
					1,
					time.Minute,
				},
				{
					"k1",
					2,
					time.Minute,
				},
			},
			"k1",
			7,
			false,
		},
		{
			"unsupported decrement start",
			[]item{},
			[]item{
				{
					"k1",
					1,
					time.Minute,
				},
			},
			"k1",
			0,
			true,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				t.Set(v.key, v.val, v.timeout)
			}
			for _, v := range tt.dec {
				lastVal, _ = t.Decrement(v.key, v.val, v.timeout)
			}

			got, err := t.Get(tt.keyCheck)

			if err != nil {
				if tt.wantErr {
					return
				} else {
					t1.Errorf("Check error = %v, want %v", err, tt.want)
					return
				}
			}

			a, err := strconv.ParseInt(strings.TrimSpace(string(got.([]byte))), 10, 64)

			if err != nil {
				t1.Errorf("Convert error = %v, want %v", err, tt.want)
			}

			if a != tt.want || tt.want != lastVal {
				t1.Errorf("Check = %v, last = %v, want %v", a, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_IncrementIn(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		key2    string
		val     int64
		timeout time.Duration
	}

	tests := []struct {
		name      string
		store     []item
		keyCheck  string
		key2Check string
		want      int64
	}{
		{
			"base test increment",
			[]item{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
			},
			"k1",
			"k2",
			1,
		},
		{
			"test multiple increment",
			[]item{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
				{
					"k1",
					"k2",
					5,
					time.Minute,
				},
			},
			"k1",
			"k2",
			6,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				lastVal, _ = t.IncrementIn(v.key, v.key2, v.val, v.timeout)
			}

			got, err := t.GetIn(tt.keyCheck, tt.key2Check)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if int64(got.(float64)) != tt.want || lastVal != tt.want {
				t1.Errorf("Check = %v, last = %v, want %v", got, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_DecrementIn(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
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

	type item struct {
		key     string
		key2    string
		val     int64
		timeout time.Duration
	}

	tests := []struct {
		name      string
		store     []item
		dec       []item
		keyCheck  string
		key2Check string
		want      int64
	}{
		{
			"base test decrement",
			[]item{
				{
					"k1",
					"k2",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
			},
			"k1",
			"k2",
			9,
		},
		{
			"test multiple decrement",
			[]item{
				{
					"k1",
					"k2",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
				{
					"k1",
					"k2",
					2,
					time.Minute,
				},
			},
			"k1",
			"k2",
			7,
		},
		{
			"test decrement start",
			[]item{},
			[]item{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
			},
			"k1",
			"k2",
			-1,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				t.SetIn(v.key, v.key2, v.val, v.timeout)
			}
			for _, v := range tt.dec {
				lastVal, _ = t.DecrementIn(v.key, v.key2, v.val, v.timeout)
			}

			got, err := t.GetIn(tt.keyCheck, tt.key2Check)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
				return
			}

			if int64(got.(float64)) != tt.want || tt.want != lastVal {
				t1.Errorf("Check = %v, last = %v, want %v", got, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(v.key)
			}
		})
	}
}

func TestMemcached_String(t1 *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"base test",
			fields{
				"test",
			},
			"test",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Memcached{
				Name: tt.fields.Name,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
