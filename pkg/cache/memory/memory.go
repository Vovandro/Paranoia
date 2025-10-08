package memory

import (
	"container/heap"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type cacheMemory struct {
	data     map[string]*cacheItem
	mutex    sync.RWMutex
	timeHeap TimeHeap
	// LRU structures per shard
	lruList  *list.List
	lruIndex map[string]*list.Element
}

type Memory struct {
	name   string
	config Config
	pool   sync.Pool
	done   chan interface{}
	data   []cacheMemory
	// global item count across shards
	itemCount int64

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
	MaxEntries    int           `yaml:"max_entries"`
	OnLimit       string        `yaml:"on_limit"` // error | ttl | lru
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
			lruList:  list.New(),
			lruIndex: make(map[string]*list.Element, 100),
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
						// init LRU
						if _, ok := t.data[shard].lruIndex[key]; !ok {
							el := t.data[shard].lruList.PushBack(key)
							t.data[shard].lruIndex[key] = el
						}
						t.data[shard].mutex.Unlock()
						heap.Push(&t.data[shard].timeHeap, &TimeHeapItem{
							key:  key,
							time: val.Timeout,
						})
						atomic.AddInt64(&t.itemCount, 1)
					}
				}
			}
		}
	}

	// defaults
	if t.config.OnLimit == "" {
		t.config.OnLimit = "ttl"
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
				t.expireShardUntil(i, now)
			}
		}
	}
}

// expireShardUntil removes all expired items in shard i whose heap item time matches current timeout
func (t *Memory) expireShardUntil(i int, now time.Time) {
	sh := &t.data[i]
	sh.mutex.Lock()
	defer sh.mutex.Unlock()
	for {
		top := sh.timeHeap.Top()
		if top == nil {
			return
		}
		thi := top.(*TimeHeapItem)
		if !thi.time.Before(now) && !thi.time.Equal(now) {
			return
		}
		// pop and validate against current item timeout
		popped := heap.Pop(&sh.timeHeap).(*TimeHeapItem)
		if val, ok := sh.data[popped.key]; ok {
			// delete only if item is actually expired now AND this heap entry matches current timeout
			if !val.Timeout.After(now) && val.Timeout.Equal(popped.time) {
				delete(sh.data, popped.key)
				if el, ok := sh.lruIndex[popped.key]; ok {
					sh.lruList.Remove(el)
					delete(sh.lruIndex, popped.key)
				}
				t.pool.Put(val)
				atomic.AddInt64(&t.itemCount, -1)
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

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()
	val, ok := t.data[shard].data[key]

	if ok && val.Timeout.After(time.Now()) {
		// LRU touch
		t.touchLRULocked(shard, key)
		t.timeRead.Record(ctx, time.Since(s).Milliseconds())
		return true
	}

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	return false
}

func (t *Memory) Set(ctx context.Context, key string, args any, timeout time.Duration) error {
	s := time.Now()
	t.counterWrite.Add(ctx, 1)

	// capacity enforcement outside of shard lock to avoid deadlocks
	if err := t.ensureCapacityForInsert(); err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

	shard := t.getShardNum(key)

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if ok {
		val.Timeout = time.Now().Add(timeout)
		val.Data = args
		// LRU touch existing
		t.touchLRULocked(shard, key)
	} else {
		val = t.pool.Get().(*cacheItem)
		val.Timeout = time.Now().Add(timeout)
		val.Data = args
		t.data[shard].data[key] = val
		// LRU add
		t.addLRULocked(shard, key)
		atomic.AddInt64(&t.itemCount, 1)
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

	if err := t.ensureCapacityForInsert(); err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return err
	}

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
		t.touchLRULocked(shard, key)
	} else {
		val = t.pool.Get().(*cacheItem)
		val.Timeout = time.Now().Add(timeout)
		val.Data = make(map[string]any)
		val.Data.(map[string]any)[key2] = args
		t.data[shard].data[key] = val
		t.addLRULocked(shard, key)
		atomic.AddInt64(&t.itemCount, 1)
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

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	if ok && val.Timeout.After(time.Now()) {
		t.touchLRULocked(shard, key)
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

	t.data[shard].mutex.Lock()
	defer t.data[shard].mutex.Unlock()

	val, ok := t.data[shard].data[key]

	t.timeRead.Record(ctx, time.Since(s).Milliseconds())
	if ok && val.Timeout.After(time.Now()) {
		if val2, ok := val.Data.(map[string]any); ok {
			if v, ok := val2[key2]; ok {
				t.touchLRULocked(shard, key)
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

	if err := t.ensureCapacityForInsert(); err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, err
	}

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
		t.touchLRULocked(shard, key)
	} else {
		v = t.pool.Get().(*cacheItem)
		v.Timeout = time.Now().Add(timeout)
		v.Data = val
		t.data[shard].data[key] = v
		t.addLRULocked(shard, key)
		atomic.AddInt64(&t.itemCount, 1)
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

	if err := t.ensureCapacityForInsert(); err != nil {
		t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
		return 0, err
	}

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
		t.touchLRULocked(shard, key)
	} else {
		v = t.pool.Get().(*cacheItem)
		v.Timeout = time.Now().Add(timeout)
		v.Data = make(map[string]any)
		v.Data.(map[string]any)[key2] = val
		t.data[shard].data[key] = v
		t.addLRULocked(shard, key)
		atomic.AddInt64(&t.itemCount, 1)
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
		if el, ok := t.data[shard].lruIndex[key]; ok {
			t.data[shard].lruList.Remove(el)
			delete(t.data[shard].lruIndex, key)
		}
		atomic.AddInt64(&t.itemCount, -1)
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
	t.touchLRULocked(shard, key)

	t.timeWrite.Record(ctx, time.Since(s).Milliseconds())
	return nil
}

// --- Helpers: LRU and capacity enforcement ---

func (t *Memory) touchLRULocked(shard int, key string) {
	sh := &t.data[shard]
	if el, ok := sh.lruIndex[key]; ok {
		sh.lruList.MoveToBack(el)
	} else {
		el := sh.lruList.PushBack(key)
		sh.lruIndex[key] = el
	}
}

func (t *Memory) addLRULocked(shard int, key string) {
	sh := &t.data[shard]
	if _, ok := sh.lruIndex[key]; !ok {
		el := sh.lruList.PushBack(key)
		sh.lruIndex[key] = el
	}
}

func (t *Memory) ensureCapacityForInsert() error {
	if t.config.MaxEntries <= 0 {
		return nil
	}
	for {
		curr := atomic.LoadInt64(&t.itemCount)
		if int(curr) < t.config.MaxEntries {
			return nil
		}
		switch t.config.OnLimit {
		case "error":
			return ErrCapacityExceeded
		case "ttl":
			if !t.enforceCapacityTTL() {
				return ErrCapacityExceeded
			}
		case "lru":
			if !t.enforceCapacityLRU() {
				return ErrCapacityExceeded
			}
		default:
			// fallback to ttl
			if !t.enforceCapacityTTL() {
				return ErrCapacityExceeded
			}
		}
	}
}

// enforceCapacityTTL tries to evict at least one item by earliest TTL across shards
func (t *Memory) enforceCapacityTTL() bool {
	type cand struct {
		shard int
		key   string
		when  time.Time
	}
	var best cand
	have := false
	now := time.Now()
	// First pass: expire
	for i := 0; i < t.config.ShardCount; i++ {
		t.expireShardUntil(i, now)
	}
	// Second pass: choose earliest among remaining
	for i := 0; i < t.config.ShardCount; i++ {
		sh := &t.data[i]
		sh.mutex.Lock()
		// skip stale tops
		for {
			top := sh.timeHeap.Top()
			if top == nil {
				break
			}
			thi := top.(*TimeHeapItem)
			v, ok := sh.data[thi.key]
			if !ok || !v.Timeout.Equal(thi.time) {
				heap.Pop(&sh.timeHeap)
				continue
			}
			// valid candidate
			if !have || thi.time.Before(best.when) {
				best = cand{shard: i, key: thi.key, when: thi.time}
				have = true
			}
			break
		}
		sh.mutex.Unlock()
	}
	if !have {
		return false
	}
	// Evict candidate
	sh := &t.data[best.shard]
	sh.mutex.Lock()
	defer sh.mutex.Unlock()
	if v, ok := sh.data[best.key]; ok {
		delete(sh.data, best.key)
		if el, ok := sh.lruIndex[best.key]; ok {
			sh.lruList.Remove(el)
			delete(sh.lruIndex, best.key)
		}
		atomic.AddInt64(&t.itemCount, -1)
		t.pool.Put(v)
		// Also remove heap item (top) if matching
		top := sh.timeHeap.Top()
		if top != nil && top.(*TimeHeapItem).key == best.key {
			heap.Pop(&sh.timeHeap)
		}
		return true
	}
	return false
}

// enforceCapacityLRU evicts least recently used from any non-empty shard
func (t *Memory) enforceCapacityLRU() bool {
	// Try to expire first
	now := time.Now()
	for i := 0; i < t.config.ShardCount; i++ {
		t.expireShardUntil(i, now)
	}
	// Choose a shard with non-empty LRU head
	for i := 0; i < t.config.ShardCount; i++ {
		sh := &t.data[i]
		sh.mutex.Lock()
		front := sh.lruList.Front()
		if front == nil {
			sh.mutex.Unlock()
			continue
		}
		key := front.Value.(string)
		if v, ok := sh.data[key]; ok {
			delete(sh.data, key)
			sh.lruList.Remove(front)
			delete(sh.lruIndex, key)
			atomic.AddInt64(&t.itemCount, -1)
			t.pool.Put(v)
			// heap item will be lazily cleaned when encountered
			sh.mutex.Unlock()
			return true
		}
		// If not ok, just remove from LRU and try next
		sh.lruList.Remove(front)
		delete(sh.lruIndex, key)
		sh.mutex.Unlock()
	}
	return false
}
