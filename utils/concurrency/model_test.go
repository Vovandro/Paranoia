package concurrency

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestFanIn(t *testing.T) {
	tests := []struct {
		name     string
		inputs   [][]int
		expected []int
	}{
		{"Merge two streams", [][]int{{1, 2, 3}, {4, 5}}, []int{1, 2, 3, 4, 5}},
		{"Empty streams", [][]int{{}, {}}, []int{}},
		{"Single stream", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем входные каналы
			var inputs []<-chan int
			for _, input := range tt.inputs {
				ch := make(chan int, len(input))
				for _, v := range input {
					ch <- v
				}
				close(ch)
				inputs = append(inputs, ch)
			}

			// Тестируем FanIn
			resultCh := FanIn(inputs...)
			result := make([]int, 0)
			for v := range resultCh {
				result = append(result, v)
			}

			sort.Ints(result)
			sort.Ints(tt.expected)

			// Проверяем результат
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFanOut(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{"Distribute evenly", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 3, 5}, {2, 4}}},
		{"Single channel", []int{1, 2, 3}, 1, [][]int{{1, 2, 3}}},
		{"More channels than data", []int{1, 2}, 3, [][]int{{1}, {2}, {}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			var w sync.WaitGroup
			w.Add(tt.size)
			outChannels := FanOut(inCh, tt.size, true)
			resOut := make(chan struct {
				I    int
				Data []int
			}, tt.size)
			results := make([][]int, tt.size)

			for i, ch := range outChannels {
				go func(i int) {
					batch := make([]int, 0)
					for v := range ch {
						batch = append(batch, v)
					}
					resOut <- struct {
						I    int
						Data []int
					}{I: i, Data: batch}
					w.Done()
				}(i)
			}

			w.Wait()
			close(resOut)

			for v := range resOut {
				results[v.I] = v.Data
			}

			// Проверяем результат
			if len(results) != len(tt.expected) {
				t.Fatalf("expected %d channels, got %d", len(tt.expected), len(results))
			}

			for i := range results {
				sort.Ints(results[i])
				sort.Ints(tt.expected[i])

				if !reflect.DeepEqual(results[i], tt.expected[i]) {
					t.Errorf("channel %d: expected %v, got %v", i, tt.expected[i], results[i])
				}
			}
		})
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{"Exact batches", []int{1, 2, 3, 4, 5, 6}, 2, [][]int{{1, 2}, {3, 4}, {5, 6}}},
		{"Incomplete batch", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"Single element batches", []int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			resultCh := Batch(inCh, tt.size)

			var results [][]int
			for batch := range resultCh {
				results = append(results, batch)
			}

			// Проверяем результат
			if len(results) != len(tt.expected) {
				t.Fatalf("expected %d batches, got %d", len(tt.expected), len(results))
			}

			for i := range results {
				sort.Ints(results[i])
				sort.Ints(tt.expected[i])

				if !reflect.DeepEqual(results[i], tt.expected[i]) {
					t.Errorf("batch %d: expected %v, got %v", i, tt.expected[i], results[i])
				}
			}
		})
	}
}

func TestParallel(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		count    int
		expected []int
	}{
		{"Double values", []int{1, 2, 3, 4}, func(x int) int { return x * 2 }, 3, []int{2, 4, 6, 8}},
		{"Square values", []int{1, 2, 3}, func(x int) int { return x * x }, 2, []int{1, 4, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			outCh := Parallel(inCh, tt.fn, tt.count)

			var results []int
			for v := range outCh {
				results = append(results, v)
			}

			sort.Ints(results)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, results)
			}
		})
	}
}
