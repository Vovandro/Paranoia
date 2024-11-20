package client

import (
	"bytes"
	"context"
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/logger"
	"gitlab.com/devpro_studio/Paranoia/server"
	"io"
	"net/http"
	"testing"
)

func TestHTTPClient_Fetch(t1 *testing.T) {
	app := Paranoia.New("test", nil, &logger.Mock{})

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
		want       interfaces.IClientResponse
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
			},
		},
		{
			name:       "base test post",
			RetryCount: 5,
			args: args{
				"POST",
				"http://127.0.0.1:8008/test",
				[]byte("{\"id\":1}"),
				nil,
			},
			want: &Response{
				bytes.NewBuffer([]byte("{\"id\":1}")),
				map[string][]string{},
				nil,
				1,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &HTTPClient{
				Config: HTTPClientConfig{
					RetryCount: tt.RetryCount,
				},
				client: http.Client{},
			}
			t.Init(app)
			s := server.Http{
				Config: server.HttpConfig{
					Port: "8008",
				},
			}
			s.Init(app)

			s.PushRoute("GET", "/", func(c context.Context, ctx interfaces.ICtx) {
				ctx.GetResponse().SetBody([]byte("{}"))
			}, nil)

			s.PushRoute("POST", "/test", func(c context.Context, ctx interfaces.ICtx) {
				d, _ := io.ReadAll(ctx.GetRequest().GetBody())
				ctx.GetResponse().SetBody(d)
			}, nil)
			s.Start()

			got := <-t.Fetch(context.Background(), tt.args.method, tt.args.host, tt.args.data, tt.args.headers)

			body, _ := got.GetBody()
			bodyWant, _ := tt.want.GetBody()

			if !bytes.Equal(body, bodyWant) {
				t1.Errorf("Fetch() = %s, want %s", body, bodyWant)
			}

			s.Stop()
		})
	}
}
