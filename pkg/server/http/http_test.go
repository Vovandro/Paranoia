package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

func TestHTTP_Fetch(t1 *testing.T) {
	type args struct {
		method  string
		host    string
		data    []byte
		headers map[string][]string
	}
	tests := []struct {
		name           string
		path           string
		args           args
		wantBody       []byte
		wantStatusCode int
	}{
		{
			name: "base test 200",
			path: "/test",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
		{
			name: "base test 404",
			path: "/test",
			args: args{
				"POST",
				"http://127.0.0.1:8009/test",
				[]byte("{\"id\":1}"),
				nil,
			},
			wantBody:       []byte(""),
			wantStatusCode: 404,
		},
		{
			name: "test dynamic",
			path: "/test/{name}",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
		{
			name: "test slash",
			path: "/test/{name}/",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
		{
			name: "test slash 2",
			path: "/test/{name}",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex/",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
		{
			name: "test default",
			path: "*",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test/alex/",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
		{
			name: "test query parameters",
			path: "/test",
			args: args{
				"GET",
				"http://127.0.0.1:8009/test?test_one=1&test_two=2",
				nil,
				nil,
			},
			wantBody:       []byte("{}"),
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			s := Http{
				config: Config{
					Port: "8009",
				},
			}
			s.Init(map[string]interface{}{
				"port":        "8009",
				"middlewares": map[string]IMiddleware{},
			})

			s.PushRoute("GET", tt.path, func(c context.Context, ctx ICtx) {
				ctx.GetResponse().SetBody([]byte("{}"))
			}, nil)

			s.Start()

			req, _ := http.NewRequest(tt.args.method, tt.args.host, bytes.NewReader(tt.args.data))

			if tt.args.headers != nil {
				for key, value := range tt.args.headers {
					if len(value) > 0 {
						req.Header.Set(key, value[0])
					}
				}
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if tt.wantStatusCode != resp.StatusCode {
					t1.Errorf("want status code %d, got %d", tt.wantStatusCode, resp.StatusCode)
				}
				return
			}

			body, _ := io.ReadAll(resp.Body)

			if !bytes.Equal(body, tt.wantBody) {
				t1.Errorf("Fetch() = %s, want %s", body, tt.wantBody)
			}

			s.Stop()
		})
	}
}

func TestHTTP_Dynamic_DoNotOverwrite(t *testing.T) {
	s := Http{
		config: Config{
			Port: "8011",
		},
	}

	_ = s.Init(map[string]interface{}{
		"port":        "8011",
		"middlewares": map[string]IMiddleware{},
	})

	s.PushRoute("GET", "/api/{v}/users", func(c context.Context, ctx ICtx) {
		ctx.GetResponse().SetBody([]byte("users"))
	}, nil)

	s.PushRoute("GET", "/api/{v}/orders", func(c context.Context, ctx ICtx) {
		ctx.GetResponse().SetBody([]byte("orders"))
	}, nil)

	s.PushRoute("GET", "/api/{v}/users/profile", func(c context.Context, ctx ICtx) {
		ctx.GetResponse().SetBody([]byte("profile"))
	}, nil)

	_ = s.Start()

	// 1) Check first route
	req1, _ := http.NewRequest("GET", "http://127.0.0.1:8011/api/v1/users", nil)
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	if resp1.StatusCode != 200 || string(body1) != "users" {
		t.Fatalf("want 200 users, got %d %s", resp1.StatusCode, string(body1))
	}

	// 2) Check sibling route under same prefix
	req2, _ := http.NewRequest("GET", "http://127.0.0.1:8011/api/v1/orders", nil)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	if resp2.StatusCode != 200 || string(body2) != "orders" {
		t.Fatalf("want 200 orders, got %d %s", resp2.StatusCode, string(body2))
	}

	// 3) Check deeper branch under existing dynamic node
	req3, _ := http.NewRequest("GET", "http://127.0.0.1:8011/api/v2/users/profile", nil)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	body3, _ := io.ReadAll(resp3.Body)
	if resp3.StatusCode != 200 || string(body3) != "profile" {
		t.Fatalf("want 200 profile, got %d %s", resp3.StatusCode, string(body3))
	}

	_ = s.Stop()
}
