package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRateLimitMiddleware_Basic(t *testing.T) {
	s := Http{
		config: Config{
			Port: "8012",
		},
	}

	// Configure middleware with small limits to trigger 429 quickly
	rl := NewRateLimitMiddleware("rate_limit").(*RateLimitMiddleware)
	err := rl.Init(nil, map[string]interface{}{
		"requests":         2,
		"interval":         "1s",
		"burst":            2,
		"key_strategy":     "global",
		"cleanup_interval": "5s",
		"evict_after":      "10s",
	})
	if err != nil {
		t.Fatalf("rate limit init error: %v", err)
	}

	_ = s.Init(map[string]interface{}{
		"port": "8012",
		"middlewares": map[string]IMiddleware{
			"rate_limit": rl,
		},
		"base_middleware": []string{"rate_limit"},
	})

	s.PushRoute("GET", "/rl", func(_ context.Context, ctx ICtx) {
		ctx.GetResponse().SetBody([]byte("{}"))
	}, nil)

	_ = s.Start()
	defer s.Stop()

	do := func() (int, []byte) {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8012/rl", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request error: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	// First two requests should pass
	if code, body := do(); code != 200 || !bytes.Equal(body, []byte("{}")) {
		t.Fatalf("want 200 {}, got %d %s", code, string(body))
	}

	if code, body := do(); code != 200 || !bytes.Equal(body, []byte("{}")) {
		t.Fatalf("want 200 {}, got %d %s", code, string(body))
	}

	// Third should be limited
	if code, _ := do(); code != 429 {
		t.Fatalf("want 429, got %d", code)
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	// After refill, should allow again
	if code, body := do(); code != 200 || !bytes.Equal(body, []byte("{}")) {
		t.Fatalf("want 200 {}, got %d %s", code, string(body))
	}
}
