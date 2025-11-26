package concurrency

import (
	"context"
	"sync"
	"time"
)

// Pool goroutine池
type Pool struct {
	workers    int
	jobQueue   chan func()
	wg         sync.WaitGroup
	once       sync.Once
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPool 创建新的goroutine池
func NewPool(workers int, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:  workers,
		jobQueue: make(chan func(), queueSize),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动goroutine池
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker 工作协程
func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case job := <-p.jobQueue:
			if job != nil {
				job()
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// Submit 提交任务
func (p *Pool) Submit(job func()) error {
	select {
	case p.jobQueue <- job:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// SubmitWithContext 使用context提交任务
func (p *Pool) SubmitWithContext(ctx context.Context, job func()) error {
	select {
	case p.jobQueue <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// Stop 停止goroutine池
func (p *Pool) Stop() {
	p.once.Do(func() {
		close(p.jobQueue)
		p.cancel()
		p.wg.Wait()
	})
}

// Wait 等待所有任务完成
func (p *Pool) Wait() {
	p.wg.Wait()
}

// WorkerPool worker池
type WorkerPool struct {
	workers    int
	jobQueue   chan interface{}
	resultQueue chan interface{}
	processor  func(interface{}) interface{}
	wg         sync.WaitGroup
	once       sync.Once
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool 创建新的worker池
func NewWorkerPool(workers int, queueSize int, processor func(interface{}) interface{}) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:     workers,
		jobQueue:    make(chan interface{}, queueSize),
		resultQueue: make(chan interface{}, queueSize),
		processor:   processor,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动worker池
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker 工作协程
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case job := <-wp.jobQueue:
			if job != nil {
				result := wp.processor(job)
				select {
				case wp.resultQueue <- result:
				case <-wp.ctx.Done():
					return
				}
			}
		case <-wp.ctx.Done():
			return
		}
	}
}

// Submit 提交任务
func (wp *WorkerPool) Submit(job interface{}) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	}
}

// GetResult 获取结果
func (wp *WorkerPool) GetResult() (interface{}, error) {
	select {
	case result := <-wp.resultQueue:
		return result, nil
	case <-wp.ctx.Done():
		return nil, wp.ctx.Err()
	}
}

// Stop 停止worker池
func (wp *WorkerPool) Stop() {
	wp.once.Do(func() {
		close(wp.jobQueue)
		wp.cancel()
		wp.wg.Wait()
		close(wp.resultQueue)
	})
}

// Wait 等待所有任务完成
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Semaphore 信号量
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore 创建新的信号量
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire 获取信号量
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release 释放信号量
func (s *Semaphore) Release() {
	<-s.ch
}

// TryAcquire 尝试获取信号量（非阻塞）
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithContext 使用context获取信号量
func (s *Semaphore) AcquireWithContext(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Size 获取当前信号量大小
func (s *Semaphore) Size() int {
	return len(s.ch)
}

// Capacity 获取信号量容量
func (s *Semaphore) Capacity() int {
	return cap(s.ch)
}

// Mutex 互斥锁（带超时）
type Mutex struct {
	ch chan struct{}
}

// NewMutex 创建新的互斥锁
func NewMutex() *Mutex {
	return &Mutex{
		ch: make(chan struct{}, 1),
	}
}

// Lock 加锁
func (m *Mutex) Lock() {
	m.ch <- struct{}{}
}

// Unlock 解锁
func (m *Mutex) Unlock() {
	<-m.ch
}

// TryLock 尝试加锁（非阻塞）
func (m *Mutex) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// LockWithContext 使用context加锁
func (m *Mutex) LockWithContext(ctx context.Context) error {
	select {
	case m.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Once 只执行一次（带错误处理）
type Once struct {
	mu    sync.Mutex
	done  bool
	value interface{}
	err   error
}

// Do 执行函数（只执行一次）
func (o *Once) Do(fn func() (interface{}, error)) (interface{}, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.done {
		return o.value, o.err
	}
	o.value, o.err = fn()
	o.done = true
	return o.value, o.err
}

// Reset 重置Once
func (o *Once) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.done = false
	o.value = nil
	o.err = nil
}

// Barrier 屏障
type Barrier struct {
	count    int
	waiters  int
	mu       sync.Mutex
	cond     *sync.Cond
	released bool
}

// NewBarrier 创建新的屏障
func NewBarrier(count int) *Barrier {
	b := &Barrier{
		count: count,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait 等待屏障
func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.released {
		return
	}
	b.waiters++
	if b.waiters >= b.count {
		b.released = true
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
}

// Reset 重置屏障
func (b *Barrier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.released = false
	b.waiters = 0
}

// WaitGroup 等待组（带超时）
type WaitGroup struct {
	wg      sync.WaitGroup
	timeout time.Duration
}

// NewWaitGroup 创建新的等待组
func NewWaitGroup(timeout time.Duration) *WaitGroup {
	return &WaitGroup{
		timeout: timeout,
	}
}

// Add 添加计数
func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
}

// Done 完成计数
func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

// Wait 等待完成
func (wg *WaitGroup) Wait() error {
	if wg.timeout <= 0 {
		wg.wg.Wait()
		return nil
	}
	done := make(chan struct{})
	go func() {
		wg.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(wg.timeout):
		return context.DeadlineExceeded
	}
}
