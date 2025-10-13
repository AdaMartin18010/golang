package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Job 表示一个工作任务
type Job struct {
	ID       int
	Data     string
	Duration time.Duration
}

// Result 表示任务执行结果
type Result struct {
	JobID    int
	Output   string
	Duration time.Duration
	Error    error
}

// Worker 工作协程
type Worker struct {
	ID         int
	JobChan    <-chan Job
	ResultChan chan<- Result
	QuitChan   chan bool
}

// NewWorker 创建新的工作协程
func NewWorker(id int, jobChan <-chan Job, resultChan chan<- Result) *Worker {
	return &Worker{
		ID:         id,
		JobChan:    jobChan,
		ResultChan: resultChan,
		QuitChan:   make(chan bool),
	}
}

// Start 启动工作协程
func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case job := <-w.JobChan:
				// 执行任务
				result := w.processJob(job)
				w.ResultChan <- result
			case <-w.QuitChan:
				fmt.Printf("Worker %d quitting\n", w.ID)
				return
			}
		}
	}()
}

// Stop 停止工作协程
func (w *Worker) Stop() {
	w.QuitChan <- true
}

// processJob 处理单个任务
func (w *Worker) processJob(job Job) Result {
	start := time.Now()

	// 模拟任务处理
	time.Sleep(job.Duration)

	// 模拟处理结果
	output := fmt.Sprintf("Worker %d processed job %d: %s", w.ID, job.ID, job.Data)

	return Result{
		JobID:    job.ID,
		Output:   output,
		Duration: time.Since(start),
		Error:    nil,
	}
}

// WorkerPool 工作池
type WorkerPool struct {
	Workers    []*Worker
	JobChan    chan Job
	ResultChan chan Result
	QuitChan   chan bool
	Wg         sync.WaitGroup
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(numWorkers int) *WorkerPool {
	jobChan := make(chan Job, 100)
	resultChan := make(chan Result, 100)

	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i+1, jobChan, resultChan)
	}

	return &WorkerPool{
		Workers:    workers,
		JobChan:    jobChan,
		ResultChan: resultChan,
		QuitChan:   make(chan bool),
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
	fmt.Printf("Starting worker pool with %d workers\n", len(wp.Workers))

	// 启动所有工作协程
	for _, worker := range wp.Workers {
		worker.Start(&wp.Wg)
	}

	// 启动结果处理协程
	go wp.handleResults()
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	fmt.Println("Stopping worker pool...")

	// 停止所有工作协程
	for _, worker := range wp.Workers {
		worker.Stop()
	}

	// 等待所有工作协程结束
	wp.Wg.Wait()

	// 关闭通道
	close(wp.JobChan)
	close(wp.ResultChan)
	close(wp.QuitChan)
}

// AddJob 添加任务到工作池
func (wp *WorkerPool) AddJob(job Job) {
	select {
	case wp.JobChan <- job:
		fmt.Printf("Job %d added to queue\n", job.ID)
	case <-wp.QuitChan:
		fmt.Println("Worker pool is shutting down")
	}
}

// handleResults 处理结果
func (wp *WorkerPool) handleResults() {
	for {
		select {
		case result := <-wp.ResultChan:
			if result.Error != nil {
				fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)
			} else {
				fmt.Printf("Job %d completed: %s (took %v)\n",
					result.JobID, result.Output, result.Duration)
			}
		case <-wp.QuitChan:
			return
		}
	}
}

// 示例使用
func main() {
	// 获取CPU核心数
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d workers (CPU cores: %d)\n", numWorkers, numWorkers)

	// 创建工作池
	pool := NewWorkerPool(numWorkers)

	// 启动工作池
	pool.Start()

	// 添加一些任务
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:       i,
			Data:     fmt.Sprintf("Task data %d", i),
			Duration: time.Duration(i*100) * time.Millisecond,
		}
		pool.AddJob(job)
	}

	// 等待一段时间让任务完成
	time.Sleep(5 * time.Second)

	// 停止工作池
	pool.Stop()

	fmt.Println("Worker pool example completed")
}
