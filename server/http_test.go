package server

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/client"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/logger"
	"gitlab.com/devpro_studio/Paranoia/server/middleware"
	"testing"
)

func TestHTTP_Fetch(t1 *testing.T) {
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
		path       string
		args       args
		want       interfaces.IClientResponse
	}{
		{
			name:       "base test 200",
			RetryCount: 5,
			path:       "/test",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "base test 404",
			RetryCount: 2,
			path:       "/test",
			args: args{
				"POST",
				"http://127.0.0.1:8009/test",
				[]byte("{\"id\":1}"),
				nil,
			},
			want: &client.Response{
				nil,
				map[string][]string{},
				fmt.Errorf("not found"),
				1,
			},
		},
		{
			name:       "test dynamic",
			RetryCount: 5,
			path:       "/test/{name}",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test slash",
			RetryCount: 5,
			path:       "/test/{name}/",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test slash 2",
			RetryCount: 5,
			path:       "/test/{name}",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex/",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test default",
			RetryCount: 5,
			path:       "*",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex/",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test query parameters",
			RetryCount: 5,
			path:       "/test",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test?test_one=1&test_two=2",
				nil,
				nil,
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &client.HTTPClient{
				Config: client.HTTPClientConfig{
					RetryCount: tt.RetryCount,
				},
			}
			t.Init(app)

			s := Http{
				Config: HttpConfig{
					Port: "8009",
				},
			}
			s.Init(app)

			s.PushRoute("GET", tt.path, func(c context.Context, ctx interfaces.ICtx) {
				ctx.GetResponse().SetBody([]byte("{}"))
			}, nil)

			s.Start()

			if got := <-t.Fetch(context.Background(), tt.args.method, tt.args.host, tt.args.data, tt.args.headers); !bytes.Equal(got.GetBody(), tt.want.GetBody()) {
				t1.Errorf("Fetch() = %s, want %s", got.GetBody(), tt.want.GetBody())
			}

			s.Stop()
		})
	}
}

func TestHTTP_Middleware(t1 *testing.T) {
	app := Paranoia.New("test", nil, &logger.Mock{}).
		PushMiddleware(&middleware.TimingMiddleware{}).
		PushMiddleware(&middleware.RestoreMiddleware{})

	app.Init()

	type args struct {
		method      string
		host        string
		middlewares []string
	}
	tests := []struct {
		name       string
		RetryCount int
		args       args
		want       interfaces.IClientResponse
	}{
		{
			name:       "test no middleware",
			RetryCount: 5,
			args: args{
				"GET",
				"http://127.0.0.1:8010/test",
				[]string{},
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test one middleware",
			RetryCount: 5,
			args: args{
				"GET",
				"http://127.0.0.1:8010/test",
				[]string{"timing"},
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
		{
			name:       "test two middleware",
			RetryCount: 5,
			args: args{
				"GET",
				"http://127.0.0.1:8010/test",
				[]string{"timing", "restore"},
			},
			want: &client.Response{
				[]byte("{}"),
				map[string][]string{},
				nil,
				1,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &client.HTTPClient{
				Config: client.HTTPClientConfig{
					RetryCount: tt.RetryCount,
				},
			}
			t.Init(app)

			s := Http{
				Config: HttpConfig{
					Port:           "8010",
					BaseMiddleware: tt.args.middlewares,
				},
			}
			s.Init(app)

			s.PushRoute("GET", "/test", func(c context.Context, ctx interfaces.ICtx) {
				ctx.GetResponse().SetBody([]byte("{}"))
			}, nil)

			s.Start()

			if got := <-t.Fetch(context.Background(), tt.args.method, tt.args.host, nil, nil); !bytes.Equal(got.GetBody(), tt.want.GetBody()) {
				t1.Errorf("Fetch() = %s, want %s", got.GetBody(), tt.want.GetBody())
			}

			s.Stop()
		})
	}
}
