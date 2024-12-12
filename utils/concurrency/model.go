package concurrency

import (
	"sync"
	"time"
)

// Pipeline convert channel for function
func Pipeline[IN, OUT any](fn func(IN) OUT) func(<-chan IN) <-chan OUT {
	return func(in <-chan IN) <-chan OUT {
		out := make(chan OUT)

		go func() {
			for val := range in {
				out <- fn(val)
			}
			close(out)
		}()

		return out
	}
}

func Split[IN, OUT1, OUT2 any](in <-chan IN, fn func(IN) (OUT1, OUT2)) (<-chan OUT1, <-chan OUT2) {
	out1 := make(chan OUT1)
	out2 := make(chan OUT2)

	go func() {
		for val := range in {
			v1, v2 := fn(val)
			out1 <- v1
			out2 <- v2
		}
		close(out1)
		close(out2)
	}()

	return out1, out2
}

// FanIn split channels
func FanIn[T any](streams ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	wg.Add(len(streams))
	for _, stream := range streams {
		go func(ch <-chan T) {
			for val := range ch {
				out <- val
			}
			wg.Done()
		}(stream)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanOut split data in channel round-robin
func FanOut[T any](ch chan T, size int, ordered bool) []<-chan T {
	out := make([]chan T, size)
	for i := 0; i < size; i++ {
		out[i] = make(chan T)
	}

	go func() {
		index := 0
		for val := range ch {
			if ordered {
				out[index] <- val
				index = (index + 1) % size
			} else {
				sent := false
				for !sent {
					for i := index; i < size; i++ {
						select {
						case out[i] <- val:
							sent = true
							index = i + 1
							break
						default:
							// pass
						}

						if sent {
							break
						}
					}

					if !sent {
						index = 0
						time.Sleep(time.Millisecond * 10)
					}
				}
			}
		}

		for _, outCh := range out {
			close(outCh)
		}
	}()

	result := make([]<-chan T, size)
	for i, outCh := range out {
		result[i] = outCh
	}
	return result
}

// Batch batch data of channel
func Batch[T any](ch <-chan T, size int) <-chan []T {
	out := make(chan []T)

	go func() {
		defer close(out)

		batch := make([]T, 0, size)

		for val := range ch {
			batch = append(batch, val)
			if len(batch) == size {
				out <- batch
				batch = make([]T, 0, size)
			}
		}

		if len(batch) > 0 {
			out <- batch
		}
	}()

	return out
}

// Parallel Worker pool
func Parallel[IN, OUT any](stream <-chan IN, fn func(IN) OUT, count int) <-chan OUT {
	out := make(chan OUT)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			for item := range stream {
				out <- fn(item)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
