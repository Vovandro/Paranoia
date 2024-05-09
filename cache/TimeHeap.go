package cache

import "time"

type TimeHeapItem struct {
	key  string
	time time.Time
}

type TimeHeap []*TimeHeapItem

func (t *TimeHeap) Len() int {
	return len(*t)
}

func (t *TimeHeap) Less(i, j int) bool {
	return (*t)[i].time.After((*t)[j].time)
}

func (t *TimeHeap) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *TimeHeap) Push(x any) {
	*t = append(*t, x.(*TimeHeapItem))
}

func (t *TimeHeap) Pop() any {
	old := *t
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*t = old[0 : n-1]
	return item
}

func (t *TimeHeap) Top() any {
	if len(*t) > 0 {
		return (*t)[len(*t)-1]
	}

	return nil
}
