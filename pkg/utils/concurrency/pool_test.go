package concurrency

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := NewPool(3, 10)
	pool.Start()
	defer pool.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		j := i
		pool.Submit(func() {
			defer wg.Done()
			_ = j
		})
	}
	wg.Wait()
}

func TestPoolWithContext(t *testing.T) {
	pool := NewPool(3, 10)
	pool.Start()
	defer pool.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := pool.SubmitWithContext(ctx, func() {
		time.Sleep(100 * time.Millisecond)
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestWorkerPool(t *testing.T) {
	processor := func(job interface{}) interface{} {
		return job.(int) * 2
	}
	pool := NewWorkerPool(3, 10, processor)
	pool.Start()
	defer pool.Stop()

	for i := 0; i < 10; i++ {
		pool.Submit(i)
	}

	for i := 0; i < 10; i++ {
		result, err := pool.GetResult()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		_ = result
	}
}

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(3)
	sem.Acquire()
	sem.Acquire()
	sem.Acquire()

	if sem.TryAcquire() {
		t.Error("Expected semaphore to be full")
	}

	sem.Release()
	if !sem.TryAcquire() {
		t.Error("Expected semaphore to have space")
	}
}

func TestSemaphoreWithContext(t *testing.T) {
	sem := NewSemaphore(1)
	sem.Acquire()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := sem.AcquireWithContext(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestMutex(t *testing.T) {
	mutex := NewMutex()
	mutex.Lock()
	defer mutex.Unlock()

	if mutex.TryLock() {
		t.Error("Expected mutex to be locked")
	}
}

func TestMutexWithContext(t *testing.T) {
	mutex := NewMutex()
	mutex.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := mutex.LockWithContext(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestOnce(t *testing.T) {
	once := &Once{}
	count := 0

	fn := func() (interface{}, error) {
		count++
		return count, nil
	}

	result1, _ := once.Do(fn)
	result2, _ := once.Do(fn)

	if result1 != 1 || result2 != 1 {
		t.Error("Expected function to be called only once")
	}
	if count != 1 {
		t.Errorf("Expected count to be 1, got %d", count)
	}
}

func TestBarrier(t *testing.T) {
	barrier := NewBarrier(3)
	var wg sync.WaitGroup
	results := make([]int, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		j := i
		go func() {
			defer wg.Done()
			results[j] = j
			barrier.Wait()
		}()
	}

	wg.Wait()
	for i := 0; i < 3; i++ {
		if results[i] != i {
			t.Errorf("Expected %d, got %d", i, results[i])
		}
	}
}

func TestWaitGroup(t *testing.T) {
	wg := NewWaitGroup(time.Second)
	wg.Add(1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}()

	err := wg.Wait()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestWaitGroupTimeout(t *testing.T) {
	wg := NewWaitGroup(100 * time.Millisecond)
	wg.Add(1)

	go func() {
		time.Sleep(200 * time.Millisecond)
		wg.Done()
	}()

	err := wg.Wait()
	if err == nil {
		t.Error("Expected timeout error")
	}
}
