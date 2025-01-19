package etcd

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestEtcd_Has(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
			map[string]string{"k2": "v1"},
			"k22222",
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for k, v := range tt.store {
				t.Set(context.Background(), k, v, time.Minute)
			}

			if got := t.Has(context.Background(), tt.key); got != tt.want {
				t1.Errorf("Has() = %v, want %v", got, tt.want)
			}

			for k, _ := range tt.store {
				t.Delete(context.Background(), k)
			}
		})
	}
}

func TestEtcd_Base(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
	}
	err := t.Init(nil)

	if err != nil {
		t1.Errorf("Init() = %v", err)
	}

	defer t.Stop()

	type item struct {
		key     string
		val     string
		timeout time.Duration
	}

	tests := []struct {
		name     string
		store    []item
		sleep    time.Duration
		keyCheck string
		want     string
	}{
		{
			"base test set",
			[]item{
				{
					"k3",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k3",
			"v1",
		},
		{
			"test not exists",
			[]item{
				{
					"k4",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k222222222",
			"",
		},
		{
			"base test timeout",
			[]item{
				{
					"k5",
					"v1",
					time.Second,
				},
			},
			time.Second * 4,
			"k5",
			"",
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.Set(context.Background(), v.key, v.val, v.timeout)
			}

			time.Sleep(tt.sleep)

			got, err := t.Get(context.Background(), tt.keyCheck)

			if err != nil && tt.want != "" {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if (got == nil && tt.want != "") || (got != nil && string(got) != tt.want) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_In(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
	}
	err := t.Init(nil)
	defer t.Stop()

	if err != nil {
		t1.Errorf("Init() = %v", err)
	}

	type item struct {
		key     string
		key2    string
		val     string
		timeout time.Duration
	}

	tests := []struct {
		name      string
		store     []item
		sleep     time.Duration
		keyCheck  string
		key2Check string
		want      string
	}{
		{
			"base test set",
			[]item{
				{
					"k6",
					"k2",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k6",
			"k2",
			"v1",
		},
		{
			"test not found",
			[]item{
				{
					"k7",
					"k2",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k21231231",
			"k2",
			"",
		},
		{
			"test not found key 2",
			[]item{
				{
					"k8",
					"k2",
					"v1",
					time.Minute,
				},
			},
			time.Microsecond,
			"k8",
			"k3",
			"",
		},
		{
			"base test multiple set one",
			[]item{
				{
					"k9",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k9",
					"k3",
					"v2",
					time.Minute,
				},
			},
			time.Microsecond,
			"k9",
			"k2",
			"v1",
		},
		{
			"base test multiple set two",
			[]item{
				{
					"k10",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k10",
					"k3",
					"v2",
					time.Minute,
				},
			},
			time.Microsecond,
			"k10",
			"k3",
			"v2",
		},
		{
			"test multiple int set one",
			[]item{
				{
					"k11",
					"k2",
					"v1",
					time.Minute,
				},
				{
					"k11",
					"k3",
					"3",
					time.Minute,
				},
			},
			time.Microsecond,
			"k11",
			"k2",
			"v1",
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.SetIn(context.Background(), v.key, v.key2, v.val, v.timeout)
			}

			time.Sleep(tt.sleep)

			got, err := t.GetIn(context.Background(), tt.keyCheck, tt.key2Check)

			if err != nil && tt.want != "" {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_Map(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
					"k12",
					map[string]interface{}{"k": "v"},
					time.Minute,
				},
			},
			"k12",
			map[string]interface{}{"k": "v"},
		},
		{
			"base test int",
			[]item{
				{
					"k13",
					map[string]interface{}{"k": "v", "k2": 5},
					time.Minute,
				},
			},
			"k13",
			map[string]interface{}{"k": "v", "k2": 5.0},
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			for _, v := range tt.store {
				t.SetMap(context.Background(), v.key, v.val, v.timeout)
			}

			got, err := t.GetMap(context.Background(), tt.keyCheck)

			if err != nil && tt.want != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Check = %v, want %v", got, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_GetMapInvalid(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
	}
	err := t.Init(nil)
	defer t.Stop()

	if err != nil {
		t1.Errorf("Init() = %v", err)
	}

	t1.Run("test get type mismatch", func(t1 *testing.T) {
		t.Set(context.Background(), "k14", "test", time.Minute)

		_, err := t.GetMap(context.Background(), "k14")

		if err == nil {
			t1.Errorf("Failed test type mismatch")
		}

		t.Delete(context.Background(), "k14")
	})
}

func TestEtcd_Increment(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
					"k15",
					1,
					time.Minute,
				},
			},
			"k15",
			1,
		},
		{
			"test multiple increment",
			[]item{
				{
					"k16",
					1,
					time.Minute,
				},
				{
					"k16",
					5,
					time.Minute,
				},
			},
			"k16",
			6,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				lastVal, _ = t.Increment(context.Background(), v.key, v.val, v.timeout)
			}

			got, err := t.Get(context.Background(), tt.keyCheck)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
			}

			a, err := strconv.ParseInt(string(got), 10, 64)

			if a != tt.want || lastVal != tt.want {
				t1.Errorf("Check = %v, last = %v, want %v", a, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_Decrement(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
					"k17",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k17",
					1,
					time.Minute,
				},
			},
			"k17",
			9,
			false,
		},
		{
			"test multiple decrement",
			[]item{
				{
					"k18",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k18",
					1,
					time.Minute,
				},
				{
					"k18",
					2,
					time.Minute,
				},
			},
			"k18",
			7,
			false,
		},
		{
			"unsupported decrement start",
			[]item{},
			[]item{
				{
					"k19",
					1,
					time.Minute,
				},
			},
			"k19",
			-1,
			true,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				t.Set(context.Background(), v.key, fmt.Sprint(v.val), v.timeout)
			}
			for _, v := range tt.dec {
				lastVal, _ = t.Decrement(context.Background(), v.key, v.val, v.timeout)
			}

			got, err := t.Get(context.Background(), tt.keyCheck)

			if err != nil {
				if tt.wantErr {
					return
				} else {
					t1.Errorf("Check error = %v, want %v", err, tt.want)
					return
				}
			}

			a, err := strconv.ParseInt(strings.TrimSpace(string(got)), 10, 64)

			if err != nil {
				t1.Errorf("Convert error = %v, want %v", err, tt.want)
			}

			if a != tt.want || tt.want != lastVal {
				t1.Errorf("Check = %v, last = %v, want %v", a, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_IncrementIn(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
					"k20",
					"k2",
					1,
					time.Minute,
				},
			},
			"k20",
			"k2",
			1,
		},
		{
			"test multiple increment",
			[]item{
				{
					"k21",
					"k2",
					1,
					time.Minute,
				},
				{
					"k21",
					"k2",
					5,
					time.Minute,
				},
			},
			"k21",
			"k2",
			6,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				lastVal, _ = t.IncrementIn(context.Background(), v.key, v.key2, v.val, v.timeout)
			}

			got, err := t.GetIn(context.Background(), tt.keyCheck, tt.key2Check)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
				return
			}

			if int64(got.(float64)) != tt.want || lastVal != tt.want {
				t1.Errorf("Check = %v, last = %v, want %v", got, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_DecrementIn(t1 *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t1.Skip()
		return
	}

	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")

	t := &Etcd{
		name: "test",
		config: Config{
			Hosts: host + ":2379",
		},
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
					"k22",
					"k2",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k22",
					"k2",
					1,
					time.Minute,
				},
			},
			"k22",
			"k2",
			9,
		},
		{
			"test multiple decrement",
			[]item{
				{
					"k23",
					"k2",
					10,
					time.Minute,
				},
			},
			[]item{
				{
					"k23",
					"k2",
					1,
					time.Minute,
				},
				{
					"k23",
					"k2",
					2,
					time.Minute,
				},
			},
			"k23",
			"k2",
			7,
		},
		{
			"test decrement start",
			[]item{},
			[]item{
				{
					"k24",
					"k2",
					1,
					time.Minute,
				},
			},
			"k24",
			"k2",
			-1,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var lastVal int64

			for _, v := range tt.store {
				t.SetIn(context.Background(), v.key, v.key2, v.val, v.timeout)
			}
			for _, v := range tt.dec {
				lastVal, _ = t.DecrementIn(context.Background(), v.key, v.key2, v.val, v.timeout)
			}

			got, err := t.GetIn(context.Background(), tt.keyCheck, tt.key2Check)

			if err != nil {
				t1.Errorf("Check error = %v, want %v", err, tt.want)
				return
			}

			if int64(got.(float64)) != tt.want || tt.want != lastVal {
				t1.Errorf("Check = %v, last = %v, want %v", got, lastVal, tt.want)
			}

			for _, v := range tt.store {
				t.Delete(context.Background(), v.key)
			}
		})
	}
}

func TestEtcd_String(t1 *testing.T) {
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
			t := &Etcd{
				name: tt.fields.Name,
			}
			if got := t.Name(); got != tt.want {
				t1.Errorf("name() = %v, want %v", got, tt.want)
			}
		})
	}
}
