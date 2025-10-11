package http

import (
	"context"
	"math"
	"strconv"
	"sync"
	"time"

	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
)

type RateLimitMiddleware struct {
	name   string
	config RateLimitMiddlewareConfig

	mu      sync.RWMutex
	buckets map[string]*bucket

	keyFunc func(ICtx) (string, int)

	stopCh        chan struct{}
	cleanupTicker *time.Ticker
	wg            sync.WaitGroup
}

type RateLimitMiddlewareConfig struct {
	Requests        int           `yaml:"requests"`
	Interval        time.Duration `yaml:"interval"`
	Burst           int           `yaml:"burst"`
	KeyStrategy     string        `yaml:"key_strategy"` // ip|header|global|method_path|ip_method_path
	HeaderName      string        `yaml:"header_name"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	EvictAfter      time.Duration `yaml:"evict_after"`
}

type bucket struct {
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
	lastUsed   time.Time
}

func NewRateLimitMiddleware(name string) interfaces2.IMiddleware {
	return &RateLimitMiddleware{
		name: name,
	}
}

func (t *RateLimitMiddleware) Init(_ interfaces2.IEngine, cfg map[string]interface{}) error {
	if err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst); err != nil {
		return err
	}

	if t.config.Requests <= 0 {
		t.config.Requests = 60
	}
	if t.config.Interval <= 0 {
		t.config.Interval = time.Minute
	}
	if t.config.Burst <= 0 {
		t.config.Burst = t.config.Requests
	}
	if t.config.KeyStrategy == "" {
		t.config.KeyStrategy = "ip"
	}
	if t.config.CleanupInterval <= 0 {
		t.config.CleanupInterval = time.Minute
	}
	if t.config.EvictAfter <= 0 {
		t.config.EvictAfter = 15 * time.Minute
	}

	t.mu.Lock()
	t.buckets = make(map[string]*bucket, 128)
	t.mu.Unlock()

	t.stopCh = make(chan struct{})
	t.cleanupTicker = time.NewTicker(t.config.CleanupInterval)
	t.wg.Add(1)
	go t.cleanupLoop()

	// default key builder
	if t.keyFunc == nil {
		t.keyFunc = t.buildKeyDefault
	}

	return nil
}

func (t *RateLimitMiddleware) cleanupLoop() {
	defer t.wg.Done()
	for {
		select {
		case <-t.cleanupTicker.C:
			t.cleanupOnce()
		case <-t.stopCh:
			if t.cleanupTicker != nil {
				t.cleanupTicker.Stop()
			}
			return
		}
	}
}

func (t *RateLimitMiddleware) cleanupOnce() {
	now := time.Now()

	// Copy entries under RLock
	t.mu.RLock()
	type pair struct {
		key string
		b   *bucket
	}
	items := make([]pair, 0, len(t.buckets))
	for k, v := range t.buckets {
		items = append(items, pair{key: k, b: v})
	}
	t.mu.RUnlock()

	// Decide which to evict
	toEvict := make([]string, 0)
	for i := 0; i < len(items); i++ {
		b := items[i].b
		b.mu.Lock()
		idle := now.Sub(b.lastUsed)
		b.mu.Unlock()
		if idle >= t.config.EvictAfter {
			toEvict = append(toEvict, items[i].key)
		}
	}

	if len(toEvict) == 0 {
		return
	}

	// Delete selected keys
	t.mu.Lock()
	for _, k := range toEvict {
		delete(t.buckets, k)
	}
	t.mu.Unlock()
}

func (t *RateLimitMiddleware) Stop() error {
	if t.stopCh != nil {
		close(t.stopCh)
	}
	t.wg.Wait()
	return nil
}

func (t *RateLimitMiddleware) Name() string { return t.name }
func (t *RateLimitMiddleware) Type() string { return "middleware" }

func (t *RateLimitMiddleware) Invoke(next RouteFunc) RouteFunc {
	return func(c context.Context, ctx ICtx) {
		key, burst := t.keyFunc(ctx)
		// If burst is 0, disable rate limiting for this request
		if burst == 0 {
			next(c, ctx)
			return
		}
		now := time.Now()

		// Get or create bucket
		t.mu.RLock()
		b := t.buckets[key]
		t.mu.RUnlock()
		if b == nil {
			t.mu.Lock()
			// Re-check after acquiring write lock
			if b = t.buckets[key]; b == nil {
				b = &bucket{
					tokens:     float64(burst),
					lastRefill: now,
					lastUsed:   now,
				}
				t.buckets[key] = b
			}
			t.mu.Unlock()
		}

		rate := float64(t.config.Requests) / t.config.Interval.Seconds()
		if rate <= 0 {
			// Degenerate config: block everything immediately
			hdr := ctx.GetResponse().Header()
			hdr.Set("X-RateLimit-Limit", strconv.Itoa(t.config.Requests))
			hdr.Set("X-RateLimit-Remaining", "0")
			hdr.Set("Retry-After", "1")
			ctx.GetResponse().SetStatus(429)
			return
		}

		// Refill and consume
		b.mu.Lock()
		dt := now.Sub(b.lastRefill).Seconds()
		if dt > 0 {
			b.tokens = math.Min(float64(burst), b.tokens+dt*rate)
			b.lastRefill = now
		}

		if b.tokens >= 1.0 {
			b.tokens -= 1.0
			remaining := int(math.Floor(b.tokens))
			b.lastUsed = now
			b.mu.Unlock()

			hdr := ctx.GetResponse().Header()
			hdr.Set("X-RateLimit-Limit", strconv.Itoa(t.config.Requests))
			hdr.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			next(c, ctx)
			return
		}

		// Not enough tokens
		needed := 1.0 - b.tokens
		secondsToOne := needed / rate
		retryAfter := int(math.Ceil(secondsToOne))
		resetAt := now.Add(time.Duration(retryAfter) * time.Second)
		b.mu.Unlock()

		hdr := ctx.GetResponse().Header()
		hdr.Set("X-RateLimit-Limit", strconv.Itoa(t.config.Requests))
		hdr.Set("X-RateLimit-Remaining", "0")
		hdr.Set("Retry-After", strconv.Itoa(retryAfter))
		hdr.Set("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		ctx.GetResponse().SetStatus(429)
		// Body is optional; keeping empty by default
	}
}

// SetKeyFunc allows overriding key building logic.
// Custom function returns key and bucket capacity; capacity 0 disables rate limiting for the request.
func (t *RateLimitMiddleware) SetKeyFunc(f func(ICtx) (string, int)) {
	t.keyFunc = f
}

// buildKeyDefault implements the default key strategy and returns default burst
func (t *RateLimitMiddleware) buildKeyDefault(ctx ICtx) (string, int) {
	req := ctx.GetRequest()
	switch t.config.KeyStrategy {
	case "global":
		return "global", t.config.Burst
	case "header":
		v := req.GetHeader().Get(t.config.HeaderName)
		if v == "" {
			return "header:unknown", t.config.Burst
		}
		return "header:" + v, t.config.Burst
	case "method_path":
		return req.GetMethod() + " " + req.GetURI(), t.config.Burst
	case "ip_method_path":
		return req.GetRemoteIP() + " " + req.GetMethod() + " " + req.GetURI(), t.config.Burst
	case "ip":
		fallthrough
	default:
		return req.GetRemoteIP(), t.config.Burst
	}
}
