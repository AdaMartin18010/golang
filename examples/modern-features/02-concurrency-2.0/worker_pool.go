package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job 表示一个工作任务
type Job[T any] struct {
	ID      string
	Data    T
	Process func(T) (T, error)
}

// Result 表示任务执行结果
type Result[T any] struct {
	JobID string
	Data  T
	Error error
}

// WorkerPool 泛型工作池
type WorkerPool[T any] struct {
	workers    int
	jobQueue   chan Job[T]
	resultChan chan Result[T]
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool[T any](workers int) *WorkerPool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool[T]{
		workers:    workers,
		jobQueue:   make(chan Job[T], 100),
		resultChan: make(chan Result[T], 100),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动工作池
func (wp *WorkerPool[T]) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker 工作协程
func (wp *WorkerPool[T]) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job := <-wp.jobQueue:
			result, err := job.Process(job.Data)
			wp.resultChan <- Result[T]{
				JobID: job.ID,
				Data:  result,
				Error: err,
			}
		case <-wp.ctx.Done():
			return
		}
	}
}

// Submit 提交任务
func (wp *WorkerPool[T]) Submit(job Job[T]) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return fmt.Errorf("job queue is full")
	}
}

// GetResult 获取结果
func (wp *WorkerPool[T]) GetResult() <-chan Result[T] {
	return wp.resultChan
}

// Shutdown 关闭工作池
func (wp *WorkerPool[T]) Shutdown() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	close(wp.resultChan)
}

// 使用示例
func main() {
	// 创建字符串处理工作池
	pool := NewWorkerPool[string](3)
	pool.Start()

	// 提交任务
	jobs := []Job[string]{
		{
			ID:   "job1",
			Data: "hello",
			Process: func(s string) (string, error) {
				time.Sleep(100 * time.Millisecond)
				return fmt.Sprintf("processed: %s", s), nil
			},
		},
		{
			ID:   "job2",
			Data: "world",
			Process: func(s string) (string, error) {
				time.Sleep(150 * time.Millisecond)
				return fmt.Sprintf("processed: %s", s), nil
			},
		},
		{
			ID:   "job3",
			Data: "golang",
			Process: func(s string) (string, error) {
				time.Sleep(200 * time.Millisecond)
				return fmt.Sprintf("processed: %s", s), nil
			},
		},
	}

	// 提交所有任务
	for _, job := range jobs {
		if err := pool.Submit(job); err != nil {
			fmt.Printf("Failed to submit job %s: %v\n", job.ID, err)
		}
	}

	// 收集结果
	go func() {
		for result := range pool.GetResult() {
			if result.Error != nil {
				fmt.Printf("Job %s failed: %v\n", result.JobID, result.Error)
			} else {
				fmt.Printf("Job %s completed: %s\n", result.JobID, result.Data)
			}
		}
	}()

	// 等待所有任务完成
	time.Sleep(1 * time.Second)
	pool.Shutdown()
}
