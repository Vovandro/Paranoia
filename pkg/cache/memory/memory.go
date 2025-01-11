package memory

import (
	"container/heap"
	"context"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"hash/crc32"
	"sync"
	"time"
)

type cacheMemory struct {
	data     map[string]*cacheItem
	mutex    sync.RWMutex
	timeHeap TimeHeap
}

type Memory struct {
	Name   string
	Config Config
	pool   sync.Pool
	done   chan interface{}
	data   []cacheMemory

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type Config struct {
	TimeClear  time.Duration `yaml:"time_clear"`
	ShardCount int           `yaml:"shard_count"`
}

type cacheItem struct {
	data    any
	timeout time.Time
}

func NewMemory(name string) *Memory {
	return &Memory{
		Name: name,
	}
}

func (t *Memory) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.Config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.Config.ShardCount < 1 {
		t.Config.ShardCount = 1
	}

	t.pool.New = func() any {
		return &cacheItem{}
	}

	t.data = make([]cacheMemory, t.Config.ShardCount)

	for i := 0; i < t.Config.ShardCount; i++ {
		t.data[i] = cacheMemory{
			data:     make(map[string]*cacheItem, 100),
			mutex:    sync.RWMutex{},
			timeHeap: make(TimeHeap, 0, 100),
		}

		heap.Init(&t.data[i].timeHeap)
	}

	t.done = make(chan interface{})

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
	if t.Config.TimeClear == 0 {
		<-t.done
		return
	}

	for {
		select {
		case <-t.done:
			return

		case <-time.After(t.Config.TimeClear):
			now := time.Now()
			for i := 0; i < t.Config.ShardCount; i++ {
				t.data[i].mutex.Lock()
				for {
					it := t.data[i].timeHeap.Top()
					if it != nil && it.(*TimeHeapItem).time.Before(now) {
						item := heap.Pop(&t.data[i].timeHeap).(*TimeHeapItem)
						if val, ok := t.data[i].data[item.key]; ok {
							delete(t.data[i].data, item.key)
							t.pool.Put(val)
						}
					} else {
						break
					}
				}
				t.data[i].mutex.Unlock()
			}
		}
	}
}

func (t *Memory) String() string {
	return t.Name
}

func (t *Memory) getShardNum(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % t.Config.ShardCount
}

func (t *Memory) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.RLock()
	defer t.data[shard].mutex.RUnlock()
	val, ok := t.data[shard].data[key]

	if ok && val.timeout.After(time.Now()) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return true
	}

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	return false
}

func (t *Memory) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if ok {
		val.timeout = time.Now().Add(timeout)
		val.data = args
	} else {
		val = t.pool.Get().(*cacheItem)
		val.timeout = time.Now().Add(timeout)
		val.data = args
		t.data[shard].data[key] = val
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: val.timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return nil
}

func (t *Memory) SetIn(ctx context.Context, key string, key2 string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if ok {
		if _, ok := val.data.(map[string]any); ok {
			val.data.(map[string]any)[key2] = args
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return ErrTypeMismatch
		}

		val.timeout = time.Now().Add(timeout)
	} else {
		val = t.pool.Get().(*cacheItem)
		val.timeout = time.Now().Add(timeout)
		val.data = make(map[string]any)
		val.data.(map[string]any)[key2] = args
		t.data[shard].data[key] = val
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: val.timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())

	return nil
}

func (t *Memory) SetMap(ctx context.Context, key string, args any, timeout time.Duration) error {
	return t.Set(ctx, key, args, timeout)
}

func (t *Memory) Get(ctx context.Context, key string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.RLock()
	defer t.data[shard].mutex.RUnlock()

	val, ok := t.data[shard].data[key]

	if ok && val.timeout.After(time.Now()) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return val.data, nil
	}

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	return nil, ErrKeyNotFound
}

func (t *Memory) GetIn(ctx context.Context, key string, key2 string) (any, error) {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.RLock()
	defer t.data[shard].mutex.RUnlock()

	val, ok := t.data[shard].data[key]

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
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
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	v, ok := t.data[shard].data[key]

	if ok {
		v.timeout = time.Now().Add(timeout)

		if _, ok := v.data.(int64); ok {
			v.data = v.data.(int64) + val
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.timeout = time.Now().Add(timeout)
		v.data = val
		t.data[shard].data[key] = v
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: v.timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.data.(int64), nil
}

func (t *Memory) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	v, ok := t.data[shard].data[key]

	if ok {
		v.timeout = time.Now().Add(timeout)

		if _, ok := v.data.(map[string]any); ok {
			if _, ok := v.data.(map[string]any)[key2]; ok {
				if v2, ok := v.data.(map[string]any)[key2].(int64); ok {
					v.data.(map[string]any)[key2] = v2 + val
				} else {
					t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
					return 0, ErrTypeMismatch
				}
			} else {
				v.data.(map[string]any)[key2] = val
			}
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.timeout = time.Now().Add(timeout)
		v.data = make(map[string]any)
		v.data.(map[string]any)[key2] = val
		t.data[shard].data[key] = v
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: v.timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.data.(map[string]any)[key2].(int64), nil
}

func (t *Memory) Decrement(ctx context.Context, key string, val int64, timeout time.Duration) (int64, error) {
	return t.Increment(ctx, key, val*-1, timeout)
}

func (t *Memory) DecrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return t.IncrementIn(ctx, key, key2, val*-1, timeout)
}

func (t *Memory) Delete(ctx context.Context, key string) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if ok {
		delete(t.data[shard].data, key)
		t.pool.Put(val)
	}

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return nil
}

func (t *Memory) Expire(ctx context.Context, key string, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if !ok {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return ErrKeyNotFound
	}

	val.timeout = time.Now().Add(timeout)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return nil
}
