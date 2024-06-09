package cache

import (
	"container/heap"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
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
}

type MemoryConfig struct {
	TimeClear time.Duration
}

type cacheItem struct {
	data    any
	timeout time.Time
}

func (t *Memory) Init(app interfaces.IService) error {
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

func (t *Memory) Has(key string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		return true
	}

	return false
}

func (t *Memory) Set(key string, args any, timeout time.Duration) error {
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

func (t *Memory) SetIn(key string, key2 string, args any, timeout time.Duration) error {
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

func (t *Memory) SetMap(key string, args any, timeout time.Duration) error {
	return t.Set(key, args, timeout)
}

func (t *Memory) Get(key string) (any, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		return val.data, nil
	}

	return nil, ErrKeyNotFound
}

func (t *Memory) GetIn(key string, key2 string) (any, error) {
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

func (t *Memory) GetMap(key string) (any, error) {
	return t.Get(key)
}

func (t *Memory) Increment(key string, val int64, timeout time.Duration) (int64, error) {
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

func (t *Memory) IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
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

func (t *Memory) Decrement(key string, val int64, timeout time.Duration) (int64, error) {
	return t.Increment(key, val*-1, timeout)
}

func (t *Memory) DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error) {
	return t.IncrementIn(key, key2, val*-1, timeout)
}

func (t *Memory) Delete(key string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	val, ok := t.data[key]

	if ok {
		delete(t.data, key)
		t.pool.Put(val)
	}

	return nil
}
