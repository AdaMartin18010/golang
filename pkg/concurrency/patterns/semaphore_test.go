package patterns

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(2)

	// 获取2个资源
	sem.Acquire()
	sem.Acquire()

	// 尝试获取第3个应该阻塞
	acquired := false
	done := make(chan bool)
	go func() {
		sem.Acquire()
		acquired = true
		done <- true
	}()

	// 等待一小段时间，第3个获取应该被阻塞
	time.Sleep(50 * time.Millisecond)
	if acquired {
		t.Error("Third acquire should be blocked")
	}

	// 释放一个，第3个应该能获取
	sem.Release()
	<-done

	if !acquired {
		t.Error("Third acquire should succeed after release")
	}

	// 清理
	sem.Release()
	sem.Release()
}

func TestSemaphoreTryAcquire(t *testing.T) {
	sem := NewSemaphore(1)

	// 第一次应该成功
	if !sem.TryAcquire() {
		t.Error("First TryAcquire should succeed")
	}

	// 第二次应该失败
	if sem.TryAcquire() {
		t.Error("Second TryAcquire should fail")
	}

	// 释放后应该成功
	sem.Release()
	if !sem.TryAcquire() {
		t.Error("TryAcquire after release should succeed")
	}

	sem.Release()
}

func TestSemaphoreWithContext(t *testing.T) {
	sem := NewSemaphore(1)
	sem.Acquire() // 占用信号量

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 应该超时
	err := sem.AcquireWithContext(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}

	sem.Release()
}

func TestParallelExecuteWithLimit(t *testing.T) {
	const numTasks = 10
	const maxConcurrent = 3

	counter := 0
	maxCounter := 0
	var mu sync.Mutex

	tasks := make([]func() error, numTasks)
	for i := 0; i < numTasks; i++ {
		tasks[i] = func() error {
			mu.Lock()
			counter++
			if counter > maxCounter {
				maxCounter = counter
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			counter--
			mu.Unlock()

			return nil
		}
	}

	errors := ParallelExecuteWithLimit(maxConcurrent, tasks)

	// 检查所有任务都完成
	for i, err := range errors {
		if err != nil {
			t.Errorf("Task %d failed: %v", i, err)
		}
	}

	// 检查并发限制
	if maxCounter > maxConcurrent {
		t.Errorf("Max concurrent tasks exceeded: %d > %d", maxCounter, maxConcurrent)
	}
}

func TestWeightedSemaphore(t *testing.T) {
	ws := NewWeightedSemaphore(10)

	// 获取权重5
	err := ws.Acquire(5)
	if err != nil {
		t.Errorf("Acquire(5) failed: %v", err)
	}

	// 再获取权重3，总共8，应该成功
	err = ws.Acquire(3)
	if err != nil {
		t.Errorf("Acquire(3) failed: %v", err)
	}

	// 尝试获取权重3，总共11，应该阻塞
	acquired := false
	go func() {
		ws.Acquire(3)
		acquired = true
	}()

	time.Sleep(50 * time.Millisecond)
	if acquired {
		t.Error("Acquire(3) should be blocked when total exceeds capacity")
	}

	// 释放5，现在应该能获取
	ws.Release(5)
	time.Sleep(50 * time.Millisecond)

	if !acquired {
		t.Error("Acquire should succeed after release")
	}

	// 清理
	ws.Release(3)
	ws.Release(3)
}

func TestWeightedSemaphoreTryAcquire(t *testing.T) {
	ws := NewWeightedSemaphore(10)

	if !ws.TryAcquire(5) {
		t.Error("TryAcquire(5) should succeed")
	}

	if !ws.TryAcquire(3) {
		t.Error("TryAcquire(3) should succeed")
	}

	// 总共8，再尝试3应该失败
	if ws.TryAcquire(3) {
		t.Error("TryAcquire(3) should fail when exceeding capacity")
	}

	ws.Release(8)
}

func TestWeightedSemaphoreExceedCapacity(t *testing.T) {
	ws := NewWeightedSemaphore(10)

	err := ws.Acquire(15)
	if err == nil {
		t.Error("Acquire exceeding capacity should fail")
	}
}

func BenchmarkSemaphore(b *testing.B) {
	sem := NewSemaphore(100)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem.Acquire()
			sem.Release()
		}
	})
}

func BenchmarkWeightedSemaphore(b *testing.B) {
	ws := NewWeightedSemaphore(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ws.Acquire(5)
			ws.Release(5)
		}
	})
}
