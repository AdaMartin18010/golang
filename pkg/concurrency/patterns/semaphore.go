package patterns

import (
	"context"
	"fmt"
	"sync"
)

// Semaphore 信号量实现，用于限制并发访问
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore 创建一个新的信号量
// maxConcurrency 指定最大并发数
func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, maxConcurrency),
	}
}

// Acquire 获取信号量
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release 释放信号量
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire 尝试获取信号量，如果无法立即获取则返回false
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithContext 使用Context获取信号量
func (s *Semaphore) AcquireWithContext(ctx context.Context) error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithSemaphore 使用信号量执行任务
func WithSemaphore(sem *Semaphore, task func() error) error {
	sem.Acquire()
	defer sem.Release()
	return task()
}

// ParallelExecuteWithLimit 并行执行多个任务，限制并发数
func ParallelExecuteWithLimit(maxConcurrency int, tasks []func() error) []error {
	sem := NewSemaphore(maxConcurrency)
	errors := make([]error, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t func() error) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()
			errors[index] = t()
		}(i, task)
	}

	wg.Wait()
	return errors
}

// WeightedSemaphore 加权信号量，支持不同任务占用不同资源
type WeightedSemaphore struct {
	capacity int
	current  int
	mu       sync.Mutex
	cond     *sync.Cond
}

// NewWeightedSemaphore 创建加权信号量
func NewWeightedSemaphore(capacity int) *WeightedSemaphore {
	ws := &WeightedSemaphore{
		capacity: capacity,
		current:  0,
	}
	ws.cond = sync.NewCond(&ws.mu)
	return ws
}

// Acquire 获取指定权重的资源
func (ws *WeightedSemaphore) Acquire(weight int) error {
	if weight > ws.capacity {
		return fmt.Errorf("weight %d exceeds capacity %d", weight, ws.capacity)
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	for ws.current+weight > ws.capacity {
		ws.cond.Wait()
	}

	ws.current += weight
	return nil
}

// Release 释放指定权重的资源
func (ws *WeightedSemaphore) Release(weight int) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.current -= weight
	if ws.current < 0 {
		ws.current = 0
	}

	ws.cond.Broadcast()
}

// TryAcquire 尝试获取资源
func (ws *WeightedSemaphore) TryAcquire(weight int) bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.current+weight <= ws.capacity {
		ws.current += weight
		return true
	}
	return false
}
