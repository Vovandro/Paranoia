package memory

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"hash/crc32"
	"os"
	"sync"
	"time"
)

type cacheMemory struct {
	data     map[string]*cacheItem
	mutex    sync.RWMutex
	timeHeap TimeHeap
}

type Memory struct {
	name   string
	config Config
	pool   sync.Pool
	done   chan interface{}
	data   []cacheMemory

	counterRead  metric.Int64Counter
	counterWrite metric.Int64Counter
	timeRead     metric.Int64Histogram
	timeWrite    metric.Int64Histogram
}

type Config struct {
	TimeClear     time.Duration `yaml:"time_clear"`
	ShardCount    int           `yaml:"shard_count"`
	EnableStorage bool          `yaml:"enable_storage"`
	StorageFile   string        `yaml:"storage_file"`
}

type cacheItem struct {
	Data    any       `json:"data"`
	Timeout time.Time `json:"timeout"`
}

func New(name string) *Memory {
	return &Memory{
		name: name,
	}
}

func (t *Memory) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.ShardCount < 1 {
		t.config.ShardCount = 1
	}

	t.pool.New = func() any {
		return &cacheItem{}
	}

	t.data = make([]cacheMemory, t.config.ShardCount)

	for i := 0; i < t.config.ShardCount; i++ {
		t.data[i] = cacheMemory{
			data:     make(map[string]*cacheItem, 100),
			mutex:    sync.RWMutex{},
			timeHeap: make(TimeHeap, 0, 100),
		}

		heap.Init(&t.data[i].timeHeap)
	}

	t.done = make(chan interface{})

	t.counterRead, _ = otel.Meter("").Int64Counter("cache_memory." + t.name + ".countRead")
	t.counterWrite, _ = otel.Meter("").Int64Counter("cache_memory." + t.name + ".countWrite")
	t.timeRead, _ = otel.Meter("").Int64Histogram("cache_memory." + t.name + ".timeRead")
	t.timeWrite, _ = otel.Meter("").Int64Histogram("cache_memory." + t.name + ".timeWrite")

	if t.config.EnableStorage {
		if t.config.StorageFile == "" {
			return fmt.Errorf("storage file is empty")
		}

		if _, err := os.Stat(t.config.StorageFile); !os.IsNotExist(err) {
			// Open the storage file for reading
			file, err := os.Open(t.config.StorageFile)
			if err == nil {
				data := make(map[string]cacheItem)
				// Read all content from the file
				fileContent, err := os.ReadFile(t.config.StorageFile)
				file.Close()

				if err != nil {
					os.Remove(t.config.StorageFile)
				} else {
					// Unmarshal the content into the Data map
					err = json.Unmarshal(fileContent, &data)
					if err != nil {
						return fmt.Errorf("could not unmarshal storage file content: %w", err)
					}

					for key, val := range data {
						shard := t.getShardNum(key)
						t.data[shard].mutex.Lock()
						t.data[shard].data[key] = &val
						t.data[shard].mutex.Unlock()
						heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
							key: key,
						})
					}
				}
			}
		}
	}

	go t.run()

	return nil
}

func (t *Memory) Stop() error {
	close(t.done)
	if t.config.EnableStorage {
		data := make(map[string]cacheItem)
		for i := 0; i < t.config.ShardCount; i++ {
			for key, val := range t.data[i].data {
				data[key] = *val
			}
		}
		file, err := os.Create(t.config.StorageFile)
		if err != nil {
			return err
		}

		json.NewEncoder(file).Encode(data)

		file.Close()
	}
	return nil
}

func (t *Memory) run() {
	if t.config.TimeClear == 0 {
		<-t.done
		return
	}

	for {
		select {
		case <-t.done:
			return

		case <-time.After(t.config.TimeClear):
			now := time.Now()
			for i := 0; i < t.config.ShardCount; i++ {
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

func (t *Memory) Name() string {
	return t.name
}

func (t *Memory) Type() string {
	return "cache"
}

func (t *Memory) getShardNum(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % t.config.ShardCount
}

func (t *Memory) Has(ctx context.Context, key string) bool {
	s := time.Now()
	t.counterRead.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.RLock()
	defer t.data[shard].mutex.RUnlock()
	val, ok := t.data[shard].data[key]

	if ok && val.Timeout.After(time.Now()) {
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
		val.Timeout = time.Now().Add(timeout)
		val.Data = args
	} else {
		val = t.pool.Get().(*cacheItem)
		val.Timeout = time.Now().Add(timeout)
		val.Data = args
		t.data[shard].data[key] = val
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: val.Timeout,
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
		if _, ok := val.Data.(map[string]any); ok {
			val.Data.(map[string]any)[key2] = args
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return ErrTypeMismatch
		}

		val.Timeout = time.Now().Add(timeout)
	} else {
		val = t.pool.Get().(*cacheItem)
		val.Timeout = time.Now().Add(timeout)
		val.Data = make(map[string]any)
		val.Data.(map[string]any)[key2] = args
		t.data[shard].data[key] = val
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: val.Timeout,
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

	if ok && val.Timeout.After(time.Now()) {
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return val.Data, nil
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
	if ok && val.Timeout.After(time.Now()) {
		if val2, ok := val.Data.(map[string]any); ok {
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
		v.Timeout = time.Now().Add(timeout)

		if _, ok := v.Data.(int64); ok {
			v.Data = v.Data.(int64) + val
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.Timeout = time.Now().Add(timeout)
		v.Data = val
		t.data[shard].data[key] = v
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: v.Timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Data.(int64), nil
}

func (t *Memory) IncrementIn(ctx context.Context, key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	v, ok := t.data[shard].data[key]

	if ok {
		v.Timeout = time.Now().Add(timeout)

		if _, ok := v.Data.(map[string]any); ok {
			if _, ok := v.Data.(map[string]any)[key2]; ok {
				if v2, ok := v.Data.(map[string]any)[key2].(int64); ok {
					v.Data.(map[string]any)[key2] = v2 + val
				} else {
					t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
					return 0, ErrTypeMismatch
				}
			} else {
				v.Data.(map[string]any)[key2] = val
			}
		} else {
			t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
			return 0, ErrTypeMismatch
		}
	} else {
		v = t.pool.Get().(*cacheItem)
		v.Timeout = time.Now().Add(timeout)
		v.Data = make(map[string]any)
		v.Data.(map[string]any)[key2] = val
		t.data[shard].data[key] = v
	}

	heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
		key:  key,
		time: v.Timeout,
	})

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return v.Data.(map[string]any)[key2].(int64), nil
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

	val.Timeout = time.Now().Add(timeout)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return nil
}
