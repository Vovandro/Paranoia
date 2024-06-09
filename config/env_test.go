package config

import (
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/logger"
	"os"
	"testing"
)

func TestEnv_GetBool(t1 *testing.T) {
	type fields struct {
		data map[string]string
		app  interfaces.IService
	}
	type args struct {
		key string
		def bool
	}

	app := Paranoia.New("test", nil, &logger.Mock{})

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"base get bool",
			fields{
				map[string]string{"k1": "true"},
				app,
			},
			args{
				"k1",
				false,
			},
			true,
		},
		{
			"base get bool string",
			fields{
				map[string]string{"k1": "1"},
				app,
			},
			args{
				"k1",
				false,
			},
			true,
		},
		{
			"base get bool string false",
			fields{
				map[string]string{"k1": "false"},
				app,
			},
			args{
				"k1",
				true,
			},
			false,
		},
		{
			"base get bool default",
			fields{
				map[string]string{"k1": "false"},
				app,
			},
			args{
				"k2",
				true,
			},
			true,
		},
		{
			"get invalid bool",
			fields{
				map[string]string{"k1": "test"},
				app,
			},
			args{
				"k2",
				false,
			},
			false,
		},
		{
			"get invalid bool 2",
			fields{
				map[string]string{"k1": "test"},
				app,
			},
			args{
				"k2",
				true,
			},
			true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data: tt.fields.data,
			}
			if got := t.GetBool(tt.args.key, tt.args.def); got != tt.want {
				t1.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_GetFloat(t1 *testing.T) {
	type fields struct {
		data map[string]string
		app  interfaces.IService
	}
	type args struct {
		key string
		def float32
	}

	app := Paranoia.New("test", nil, &logger.Mock{})

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		{
			"base get float",
			fields{
				map[string]string{"k1": "1.0"},
				app,
			},
			args{
				"k1",
				0,
			},
			1,
		},
		{
			"base get big float",
			fields{
				map[string]string{"k1": "0.000001"},
				app,
			},
			args{
				"k1",
				0,
			},
			0.000001,
		},
		{
			"base get float from int",
			fields{
				map[string]string{"k1": "1"},
				app,
			},
			args{
				"k1",
				0,
			},
			1,
		},
		{
			"base get float default",
			fields{
				map[string]string{"k1": "1.0"},
				app,
			},
			args{
				"k2",
				2,
			},
			2,
		},
		{
			"get invalid float",
			fields{
				map[string]string{"k1": "test"},
				app,
			},
			args{
				"k2",
				2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data: tt.fields.data,
			}
			if got := t.GetFloat(tt.args.key, tt.args.def); got != tt.want {
				t1.Errorf("GetFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_GetInt(t1 *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		key string
		def int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"base get int",
			fields{
				map[string]string{"k1": "1"},
			},
			args{
				"k1",
				0,
			},
			1,
		},
		{
			"base get int from float",
			fields{
				map[string]string{"k1": "1.1"},
			},
			args{
				"k1",
				2,
			},
			2,
		},
		{
			"base get int default",
			fields{
				map[string]string{"k1": "1"},
			},
			args{
				"k2",
				2,
			},
			2,
		},
		{
			"get invalid int",
			fields{
				map[string]string{"k1": "test"},
			},
			args{
				"k2",
				2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data:   tt.fields.data,
				logger: &logger.Mock{},
			}

			if got := t.GetInt(tt.args.key, tt.args.def); got != tt.want {
				t1.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_GetString(t1 *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		key string
		def string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"base get string",
			fields{
				map[string]string{"k1": "1r"},
			},
			args{
				"k1",
				"",
			},
			"1r",
		},
		{
			"base get string default",
			fields{
				map[string]string{"k1": "1"},
			},
			args{
				"k2",
				"2",
			},
			"2",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data:   tt.fields.data,
				logger: &logger.Mock{},
			}
			if got := t.GetString(tt.args.key, tt.args.def); got != tt.want {
				t1.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Has(t1 *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"base has exists",
			fields{
				map[string]string{"k1": "test"},
			},
			args{
				"k1",
			},
			true,
		},
		{
			"base has not exists",
			fields{
				map[string]string{"k1": "test"},
			},
			args{
				"k2",
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data: tt.fields.data,
			}
			if got := t.Has(tt.args.key); got != tt.want {
				t1.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Init(t1 *testing.T) {
	type fields struct {
		FName string
	}
	type file struct {
		FName    string
		FileData string
	}
	type want struct {
		key   string
		value string
	}

	tests := []struct {
		name    string
		fields  fields
		file    file
		wantErr bool
		want    want
	}{
		{
			"base test",
			fields{
				"./.env.test",
			},
			file{
				"./.env.test",
				"ENV=test\nAPP=app\n",
			},
			false,
			want{
				"ENV",
				"test",
			},
		},
		{
			"base test last",
			fields{
				"./.env.test",
			},
			file{
				"./.env.test",
				"ENV=test\nAPP=app\n",
			},
			false,
			want{
				"APP",
				"app",
			},
		},
		{
			"base test with comment",
			fields{
				"./.env.test",
			},
			file{
				"./.env.test",
				"ENV=test\n\tAPP=app #test comment\n",
			},
			false,
			want{
				"APP",
				"app",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewEnv(EnvConfig{
				FName: tt.fields.FName,
			})

			_ = os.WriteFile(tt.file.FName, []byte(tt.file.FileData), 0666)

			if err := t.Init(&logger.Mock{}); (err != nil) != tt.wantErr {
				t1.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got, ok := t.data[tt.want.key]; !ok || got != tt.want.value {
				t1.Errorf("check %v, want %v", got, tt.want)
			}

			_ = os.Remove(tt.file.FName)
		})
	}
}

func TestEnv_ParseFile(t1 *testing.T) {
	type args struct {
		data []byte
	}
	type want struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"base test",
			args{
				[]byte("ENV=test"),
			},
			want{
				"ENV",
				"test",
			},
		},
		{
			"test comment",
			args{
				[]byte("ENV=test #comment"),
			},
			want{
				"ENV",
				"test",
			},
		},
		{
			"test quote",
			args{
				[]byte("ENV=\"test\""),
			},
			want{
				"ENV",
				"test",
			},
		},
		{
			"test quote and comment",
			args{
				[]byte("ENV=\"test\" #comment"),
			},
			want{
				"ENV",
				"test",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Env{
				data: map[string]string{},
			}

			t.ParseFile(tt.args.data)

			if got, ok := t.data[tt.want.key]; !ok || got != tt.want.value {
				t1.Errorf("check %v, want %v", got, tt.want)
			}
		})
	}
}
