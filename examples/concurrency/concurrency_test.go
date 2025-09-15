package concurrency

import (
	"runtime"
	"sync"
	"testing"
)

// 简单并发安全计数器
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc()       { c.mu.Lock(); c.n++; c.mu.Unlock() }
func (c *Counter) Value() int { c.mu.Lock(); defer c.mu.Unlock(); return c.n }

func TestCounterRaceFree(t *testing.T) {
	t.Parallel()
	var c Counter
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0) * 2
	loops := 1000
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < loops; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
	if got, want := c.Value(), workers*loops; got != want {
		t.Fatalf("unexpected value: got=%d want=%d", got, want)
	}
}

func BenchmarkCounter(b *testing.B) {
	var c Counter
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}
