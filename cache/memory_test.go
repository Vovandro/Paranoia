package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestMemory_Delete(t1 *testing.T) {
	tests := []struct {
		name    string
		fields  string
		args    string
		wantErr bool
	}{
		{
			"delete exists",
			"k1",
			"k1",
			false,
		},
		{
			"delete not exists",
			"k1",
			"k2",
			false,
		},
		{
			"delete empty",
			"k1",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Memory{
				Name: "test",
			}
			t.data.Store(tt.fields, &cacheItem{"test", time.Now().Add(time.Minute)})

			if err := t.Delete(tt.args); (err != nil) != tt.wantErr {
				t1.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Get(t1 *testing.T) {
	type fields struct {
		key string
		val any
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		want    any
		wantErr bool
	}{
		{
			name: "test exists string",
			fields: fields{
				key: "k1",
				val: "val",
			},
			args:    "k1",
			want:    "val",
			wantErr: false,
		},
		{
			name: "test not exists string",
			fields: fields{
				key: "k1",
				val: "val",
			},
			args:    "k2",
			want:    nil,
			wantErr: true,
		},
		{
			name: "test exists array string",
			fields: fields{
				key: "k1",
				val: []string{"val", "val2"},
			},
			args:    "k1",
			want:    []string{"val", "val2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Memory{}
			t.Init(nil)

			t.data.Store(tt.fields.key, &cacheItem{tt.fields.val, time.Now().Add(time.Minute)})

			got, err := t.Get(tt.args)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Has(t1 *testing.T) {
	tests := []struct {
		name   string
		fields string
		args   string
		want   bool
	}{
		{
			"test has exists key",
			"k1",
			"k1",
			true,
		},
		{
			"test has not exists key",
			"k1",
			"k2",
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Memory{}

			t.data.Store(tt.fields, &cacheItem{"test", time.Now().Add(time.Minute)})

			if got := t.Has(tt.args); got != tt.want {
				t1.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Set(t1 *testing.T) {
	type args struct {
		key     string
		args    any
		timeout time.Duration
	}
	tests := []struct {
		name       string
		args       []args
		sleep      time.Duration
		want       any
		wantErr    bool
		wantExists bool
	}{
		{
			"base test string",
			[]args{
				{
					"k1",
					"test",
					time.Minute,
				},
			},
			time.Microsecond,
			"test",
			false,
			true,
		},
		{
			"base test array string",
			[]args{
				{
					"k1",
					[]string{"test"},
					time.Minute,
				},
			},
			time.Microsecond,
			[]string{"test"},
			false,
			true,
		},
		{
			"base test timeout",
			[]args{
				{
					"k1",
					"test",
					time.Millisecond,
				},
			},
			time.Millisecond * 10,
			nil,
			false,
			false,
		},
		{
			"base test replace",
			[]args{
				{
					"k1",
					"test",
					time.Minute,
				},
				{
					"k1",
					"test2",
					time.Minute,
				},
			},
			time.Millisecond,
			"test2",
			false,
			true,
		},
		{
			"base test replace timeout",
			[]args{
				{
					"k1",
					"test",
					time.Minute,
				},
				{
					"k1",
					"test2",
					time.Millisecond,
				},
			},
			time.Millisecond * 10,
			nil,
			false,
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Memory{}
			t.Init(nil)

			for i := 0; i < len(tt.args); i++ {
				if err := t.Set(tt.args[i].key, tt.args[i].args, tt.args[i].timeout); (err != nil) != tt.wantErr {
					t1.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, ok := t.Get(tt.args[0].key)

			if (ok == nil) != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_String(t1 *testing.T) {
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
			t := &Memory{
				Name: tt.fields.Name,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
