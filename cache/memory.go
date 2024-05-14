package cache

import (
	"container/heap"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"sync"
	"time"
)

type Memory struct {
	Name      string
	TimeClear time.Duration
	data      map[string]*cacheItem
	pool      sync.Pool
	mutex     sync.RWMutex
	timeHeap  TimeHeap
	done      chan interface{}
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

	if t.TimeClear <= 0 {
		t.TimeClear = time.Second * 10
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

		case <-time.After(t.TimeClear):
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

func (t *Memory) Get(key string) (any, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	val, ok := t.data[key]

	if ok && val.timeout.After(time.Now()) {
		return val.data, nil
	}

	return nil, fmt.Errorf("key not found")
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
