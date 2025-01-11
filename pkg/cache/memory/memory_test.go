package memory

import (
	"context"
	"fmt"
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
				Config: Config{
					TimeClear:  time.Minute,
					ShardCount: 5,
				},
			}
			t.Init(nil)
			defer t.Stop()

			shard := t.getShardNum(tt.fields)

			t.data[shard].data[tt.fields] = &cacheItem{"test", time.Now().Add(time.Minute)}

			if err := t.Delete(context.Background(), tt.args); (err != nil) != tt.wantErr {
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
			t := &Memory{
				Config: Config{
					TimeClear:  time.Minute,
					ShardCount: 5,
				},
			}
			t.Init(nil)
			defer t.Stop()

			shard := t.getShardNum(tt.fields.key)

			t.data[shard].data[tt.fields.key] = &cacheItem{tt.fields.val, time.Now().Add(time.Minute)}

			got, err := t.Get(context.Background(), tt.args)
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
			t := &Memory{
				Config: Config{
					TimeClear:  time.Minute,
					ShardCount: 5,
				},
			}
			t.Init(nil)
			defer t.Stop()

			shard := t.getShardNum(tt.fields)

			t.data[shard].data[tt.fields] = &cacheItem{"test", time.Now().Add(time.Minute)}

			if got := t.Has(context.Background(), tt.args); got != tt.want {
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
			t := &Memory{
				Config: Config{
					TimeClear:  time.Minute,
					ShardCount: 5,
				},
			}
			t.Init(nil)
			defer t.Stop()

			for i := 0; i < len(tt.args); i++ {
				if err := t.Set(context.Background(), tt.args[i].key, tt.args[i].args, tt.args[i].timeout); (err != nil) != tt.wantErr {
					t1.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, ok := t.Get(context.Background(), tt.args[0].key)

			if (ok == nil) != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_SetIn(t1 *testing.T) {
	type args struct {
		key     string
		key2    string
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
					"k2",
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
					"k2",
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
					"k2",
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
					"k2",
					"test",
					time.Minute,
				},
				{
					"k1",
					"k2",
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
					"k2",
					"test",
					time.Minute,
				},
				{
					"k1",
					"k2",
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
			defer t.Stop()

			for i := 0; i < len(tt.args); i++ {
				if err := t.SetIn(context.Background(), tt.args[i].key, tt.args[i].key2, tt.args[i].args, tt.args[i].timeout); (err != nil) != tt.wantErr {
					t1.Errorf("SetIn() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, ok := t.GetIn(context.Background(), tt.args[0].key, tt.args[0].key2)

			if (ok == nil) != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetIn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Increment(t1 *testing.T) {
	type args struct {
		key     string
		args    int64
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
			"base test one",
			[]args{
				{
					"k1",
					1,
					time.Minute,
				},
			},
			time.Microsecond,
			int64(1),
			false,
			true,
		},
		{
			"base test two",
			[]args{
				{
					"k1",
					2,
					time.Minute,
				},
			},
			time.Microsecond,
			int64(2),
			false,
			true,
		},
		{
			"base test timeout",
			[]args{
				{
					"k1",
					1,
					time.Millisecond,
				},
			},
			time.Millisecond * 10,
			nil,
			false,
			false,
		},
		{
			"base test multiple increment",
			[]args{
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
			time.Millisecond,
			int64(3),
			false,
			true,
		},
		{
			"base test replace timeout",
			[]args{
				{
					"k1",
					1,
					time.Minute,
				},
				{
					"k1",
					1,
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
			defer t.Stop()

			for i := 0; i < len(tt.args); i++ {
				if _, err := t.Increment(context.Background(), tt.args[i].key, tt.args[i].args, tt.args[i].timeout); (err != nil) != tt.wantErr {
					t1.Errorf("Increment() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, ok := t.Get(context.Background(), tt.args[0].key)

			if (ok == nil) != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_IncrementIn(t1 *testing.T) {
	type args struct {
		key     string
		key2    string
		args    int64
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
			"base test one",
			[]args{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
			},
			time.Microsecond,
			int64(1),
			false,
			true,
		},
		{
			"base test two",
			[]args{
				{
					"k1",
					"k2",
					2,
					time.Minute,
				},
			},
			time.Microsecond,
			int64(2),
			false,
			true,
		},
		{
			"base test timeout",
			[]args{
				{
					"k1",
					"k2",
					1,
					time.Millisecond,
				},
			},
			time.Millisecond * 10,
			nil,
			false,
			false,
		},
		{
			"base test multiple increment",
			[]args{
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
			time.Millisecond,
			int64(3),
			false,
			true,
		},
		{
			"base test replace timeout",
			[]args{
				{
					"k1",
					"k2",
					1,
					time.Minute,
				},
				{
					"k1",
					"k2",
					1,
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
			defer t.Stop()

			for i := 0; i < len(tt.args); i++ {
				if _, err := t.IncrementIn(context.Background(), tt.args[i].key, tt.args[i].key2, tt.args[i].args, tt.args[i].timeout); (err != nil) != tt.wantErr {
					t1.Errorf("IncrementIn() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, ok := t.GetIn(context.Background(), tt.args[0].key, tt.args[0].key2)

			if (ok == nil) != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetIn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_ClearTimeout(t1 *testing.T) {
	t1.Run("test timeout clear", func(t1 *testing.T) {
		t := &Memory{
			Config: Config{
				TimeClear:  time.Millisecond * 10,
				ShardCount: 2,
			},
		}
		t.Init(nil)
		defer t.Stop()

		err := t.Set(context.Background(), "test", "test", time.Millisecond)

		if err != nil {
			t1.Errorf("Set() error = %v", err)
		}

		time.Sleep(time.Second)

		shard := t.getShardNum("test")

		t.data[shard].mutex.RLock()
		_, ok := t.data[shard].data["test"]
		t.data[shard].mutex.RUnlock()

		if ok {
			t1.Errorf("Get() exists")
		}
	})
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

func BenchmarkStore(b *testing.B) {
	t := Memory{
		Config: Config{
			TimeClear:  time.Second,
			ShardCount: 1,
		},
	}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, k, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
	}
}

func BenchmarkStore10(b *testing.B) {
	t := Memory{
		Config: Config{
			TimeClear:  time.Second,
			ShardCount: 10,
		},
	}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, k, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
	}
}

func BenchmarkStoreMutex(b *testing.B) {
	t := Memory{}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, k, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
	}
}

func BenchmarkCheckAndStore(b *testing.B) {
	t := Memory{}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, k, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}

		if val, ok := t.Get(context.Background(), k); ok != nil || val.(string) != k {
			b.Fatalf("Unexpected error data")
		}
	}
}

func BenchmarkCheckAndStoreMutex(b *testing.B) {
	t := Memory{}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, k, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}

		if val, ok := t.Get(context.Background(), k); ok != nil || val.(string) != k {
			b.Fatalf("Unexpected error data")
		}
	}
}

func BenchmarkCheckAndStoreArray(b *testing.B) {
	t := Memory{}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, []string{k}, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}

		if val, ok := t.Get(context.Background(), k); ok != nil || !reflect.DeepEqual(val.([]string), []string{k}) {
			b.Fatalf("Unexpected error data")
		}
	}
}

func BenchmarkCheckAndStoreArrayMutex(b *testing.B) {
	t := Memory{}
	t.Init(nil)
	defer t.Stop()

	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Set(context.Background(), k, []string{k}, time.Hour)

		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}

		if val, ok := t.Get(context.Background(), k); ok != nil || !reflect.DeepEqual(val.([]string), []string{k}) {
			b.Fatalf("Unexpected error data")
		}
	}
}
