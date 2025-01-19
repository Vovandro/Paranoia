package file

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestFile_Has(t1 *testing.T) {
	type fields struct {
		fName string
		data  []byte
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "exists",
			fields: fields{
				fName: "test.txt",
				data:  []byte("hello world"),
			},
			args: args{
				name: "test.txt",
			},
			want: true,
		},
		{
			name: "not exists",
			fields: fields{
				fName: "test.txt",
				data:  []byte("hello world"),
			},
			args: args{
				name: "test_2.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &File{}

			err := os.WriteFile(tt.fields.fName, tt.fields.data, 0600)

			if err != nil {
				t1.Errorf("os.WriteFile() error = %v", err)
				return
			}

			if got := t.Has(tt.args.name); got != tt.want {
				t1.Errorf("Has() = %v, want %v", got, tt.want)
			}

			_ = os.Remove(tt.fields.fName)
		})
	}
}

func TestFile_Put(t1 *testing.T) {
	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"base test",
			args{
				name: "test.txt",
				data: []byte("hello world"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &File{}

			if err := t.Put(tt.args.name, bytes.NewReader(tt.args.data)); (err != nil) != tt.wantErr {
				t1.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, err := t.Read(tt.args.name)

			if err != nil {
				t1.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}

			defer file.Close()

			d, err := io.ReadAll(file)

			if !bytes.Equal(d, tt.args.data) {
				t1.Errorf("Put() = %v, want %v", d, tt.args.data)
			}

			if _, err = t.GetSize(tt.args.name); err != nil {
				t1.Errorf("GetSize() error = %v", err)
			}

			if _, err = t.GetModified(tt.args.name); err != nil {
				t1.Errorf("GetModified() error = %v", err)
			}

			_ = t.Delete(tt.args.name)
		})
	}
}

func TestFile_StoreFolder(t1 *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		has     bool
	}{
		{
			"base test",
			args{
				"./test",
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &File{}

			if err := t.StoreFolder(tt.args.name); (err != nil) != tt.wantErr {
				t1.Errorf("StoreFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if t.Has(tt.args.name) != tt.has {
				t1.Errorf("Has() = %v, want %v", t.Has(tt.args.name), tt.has)
			}

			if ok, _ := t.IsFolder(tt.args.name); !ok {
				t1.Error("IsFolder() error")
			}

			if _, err := t.List(tt.args.name); err != nil {
				t1.Errorf("List() error = %v", err)
			}

			_ = t.Delete(tt.args.name)
		})
	}
}
