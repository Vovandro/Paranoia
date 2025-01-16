package http_client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

func TestHTTPClient_Fetch(t1 *testing.T) {
	type args struct {
		method  string
		host    string
		data    []byte
		headers map[string][]string
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
				"http://127.0.0.1:8008/",
				nil,
				nil,
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
				"http://127.0.0.1:8008/test",
				[]byte("{\"id\":1}"),
				map[string][]string{
					"Content-Type": {"application/json"},
				},
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
			t := NewHTTPClient("test")

			t.Init(map[string]interface{}{
				"retry_count": tt.RetryCount,
			})

			server := &http.Server{Addr: ":8008"}

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
				server.ListenAndServe()
			}()

			got := <-t.Fetch(context.Background(), tt.args.method, tt.args.host, tt.args.data, tt.args.headers)

			body, _ := got.GetBody()
			bodyWant, _ := tt.want.GetBody()

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
