package cache

import (
	"container/heap"
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"sync"
	"time"
)

type Memory struct {
	Name     string
	Config   MemoryConfig
	data     map[string]*cacheItem
	pool     sync.Pool
	mutex    sync.RWMutex
	timeHeap TimeHeap
	done     chan interface{}

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type MemoryConfig struct {
	TimeClear time.Duration `yaml:"time_clear"`
}

type cacheItem struct {
	data    any
	timeout time.Time
}

func NewMemory(name string, cfg MemoryConfig) *Memory {
	return &Memory{
		Name:   name,
		Config: cfg,
	}
}

func (t *Memory) Init(app interfaces.IEngine) error {
	t.pool.New = func() any {
		return &cacheItem{}
	}

	t.data = make(map[string]*cacheItem, 100)
	t.timeHeap = make(TimeHeap, 0, 100)
	heap.Init(&t.timeHeap)

	t.done = make(chan interface{})

	if t.Config.TimeClear <= 0 {
		t.Config.TimeClear = time.Second * 10
	}

	t.counterRead, _ = otel.Meter("").Int64Counter("cache_memory." + t.Name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("cache_memory." + t.Name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("cache_memory." + t.Name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("cache_memory." + t.Name + ".timeWrite")

	go t.run()

	return nil
}

func (t *Memory) Stop() error {
	close(t.done)
	return nil
}

func (t *Memory) run() {
	for true {
		select {
		case <-t.done:
			return

		case <-time.After(t.Config.TimeClear):
			now := time.Now()
			t.mutex.Lock()
			for true {
				i := t.timeHeap.Top()
				if i != nil && i.(*TimeHeapItem).time.Before(now) {
					item := heap.Pop(&t.timeHeap).(*TimeHeapItem)
					if val, ok := t.data[item.key]; ok {
						delete(t.data, item.key)
						t.pool.Put(val)
					}
				} else {
					break
				}
			}
			t.mutex.Unlock()
		}
	}
}

func (t *Memory) String() string {
	return t.Name
}

func (t *Memory) Has(ctx context.Context, key string) bool {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	t.mutex.RLock()
	defer t.mutex.RUnlock()
	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		return true
	}

	return false
}

func (t *Memory) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		val.timeout = time.Now().Add(timeout)
		val.data = args
	} else {
		val = t.pool.Get().(*cacheItem)
		val.timeout = time.Now().Add(timeout)
		val.data = args
		t.data[key] = val
	}

	heap.Push(&t.timeHeap, &TimeHeapItem{
		key:  key,
		time: val.timeout,
	})

	return nil
}

func (t *Memory) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		if _, ok := val.data.(map[string]any); ok {
			val.data.(map[string]any)[key2] = args
		} else {
			return ErrTypeMismatch
		}

		val.timeout = time.Now().Add(timeout)
	} else {
		val = t.pool.Get().(*cacheItem)
		val.timeout = time.Now().Add(timeout)
		val.data = make(map[string]any)
		val.data.(map[string]any)[key2] = args
		t.data[key] = val
	}

	heap.Push(&t.timeHeap, &TimeHeapItem{
		key:  key,
		time: val.timeout,
	})

	return nil
}

func (t *Memory) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	return t.Set(ctx, key, args, timeout)
}

func (t *Memory) Get(ctx context.Context, key string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		return val.data, nil
	}

	return nil, ErrKeyNotFound
}

func (t *Memory) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	defer func(s time.Time) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterRead.Add(ctx, 1)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		if val2, ok := val.data.(map[string]any); ok {
			if v, ok := val2[key2]; ok {
				return v, nil
			} else {
				return nil, ErrKeyNotFound
			}
		} else {
			return nil, ErrTypeMismatch
		}
	}

	return nil, ErrKeyNotFound
}

func (t *Memory) GetMap(ctx context.Context, key string) (any, error) {
	return t.Get(ctx, key)
}

func (t *Memory) Increment(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	v, ok := t.data[key]

	if ok {
		v.timeout = time.Now().Add(timeout)

		if _, ok := v.data.(int64); ok {
			v.data = v.data.(int64) + val
		} else {
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.timeout = time.Now().Add(timeout)
		v.data = val
		t.data[key] = v
	}

	heap.Push(&t.timeHeap, &TimeHeapItem{
		key:  key,
		time: v.timeout,
	})

	return v.data.(int64), nil
}

func (t *Memory) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	v, ok := t.data[key]

	if ok {
		v.timeout = time.Now().Add(timeout)

		if _, ok := v.data.(map[string]any); ok {
			if _, ok := v.data.(map[string]any)[key2]; ok {
				if v2, ok := v.data.(map[string]any)[key2].(int64); ok {
					v.data.(map[string]any)[key2] = v2 + val
				} else {
					return 0, ErrTypeMismatch
				}
			} else {
				v.data.(map[string]any)[key2] = val
			}
		} else {
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.timeout = time.Now().Add(timeout)
		v.data = make(map[string]any)
		v.data.(map[string]any)[key2] = val
		t.data[key] = v
	}

	heap.Push(&t.timeHeap, &TimeHeapItem{
		key:  key,
		time: v.timeout,
	})

	return v.data.(map[string]any)[key2].(int64), nil
}

func (t *Memory) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return t.Increment(ctx, key, val*-1, timeout)
}

func (t *Memory) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return t.IncrementIn(ctx, key, key2, val*-1, timeout)
}

func (t *Memory) Delete(ctx context.Context, key string) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		delete(t.data, key)
		t.pool.Put(val)
	}

	return nil
}

func (t *Memory) Expire(ctx context.Context, key string, timeout time.Duration) error {
	defer func(s time.Time) {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	}(time.Now())
	t.counterWrite.Add(ctx, 1)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if !ok {
		return ErrKeyNotFound
	}

	val.timeout = time.Now().Add(timeout)

	return nil
}
