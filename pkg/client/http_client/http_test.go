package http_client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestHTTPClient_Fetch(t1 *testing.T) {
	type args struct {
		method     string
		host       string
		data       []byte
		headers    map[string][]string
		serverPort string
	}
	tests := []struct {
		name       string
		RetryCount int
		args       args
		want       IClientResponse
	}{
		{
			name:       "base test 200",
			RetryCount: 5,
			args: args{
				"GET",
				"http://127.0.0.1:9008/",
				nil,
				nil,
				"127.0.0.1:9008",
			},
			want: &Response{
				bytes.NewBuffer([]byte("{}")),
				map[string][]string{},
				nil,
				1,
				200,
			},
		},
		{
			name:       "base test post",
			RetryCount: 5,
			args: args{
				"POST",
				"http://127.0.0.1:9009/test",
				[]byte("{\"id\":1}"),
				map[string][]string{
					"Content-Type": {"application/json"},
				},
				"127.0.0.1:9009",
			},
			want: &Response{
				bytes.NewBuffer([]byte("{\"id\":1}")),
				map[string][]string{},
				nil,
				1,
				200,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := New("test")

			t.Init(map[string]interface{}{
				"retry_count": tt.RetryCount,
			})

			server := &http.Server{Addr: tt.args.serverPort}

			server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/":
					if r.Method == "GET" {
						w.Write([]byte("{}"))
						return
					}

				case "/test":
					if r.Method == "POST" {
						data, _ := io.ReadAll(r.Body)
						w.Write(data)
						return
					}
				}

				w.WriteHeader(http.StatusNotFound)
			})

			go func() {
				err := server.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					t1.Errorf("ListenAndServe() error = %v", err)
				}
			}()

			time.Sleep(time.Second)

			got := <-t.Fetch(context.Background(), tt.args.method, tt.args.host, tt.args.data, tt.args.headers)

			body, _ := got.GetBody()
			bodyWant, _ := tt.want.GetBody()

			if !errors.Is(got.Error(), tt.want.Error()) {
				t1.Errorf("Fetch() = %s, want %s", got.Error(), tt.want.Error())
			}

			if got.GetCode() != tt.want.GetCode() {
				t1.Errorf("Fetch() = %d, want %d", got.GetCode(), tt.want.GetCode())
			}

			if !bytes.Equal(body, bodyWant) {
				t1.Errorf("Fetch() = %s, want %s", body, bodyWant)
			}

			server.Shutdown(context.Background())
		})
	}
}
