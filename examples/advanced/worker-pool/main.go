// Worker Pool示例：使用WaitGroup.Go()实现高效的工作池
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Task 表示一个工作任务
type Task struct {
	ID   int
	Data string
}

// WorkerPool 工作池
type WorkerPool struct {
	maxWorkers int
	tasks      chan Task
	wg         sync.WaitGroup
	results    chan Result
}

// Result 任务结果
type Result struct {
	TaskID   int
	Output   string
	Duration time.Duration
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		tasks:      make(chan Task, maxWorkers*2),
		results:    make(chan Result, maxWorkers*2),
	}
}

// Start 启动工作池
func (p *WorkerPool) Start() {
	for i := 0; i < p.maxWorkers; i++ {
		workerID := i
		// 使用 WaitGroup.Go() 简化 goroutine 管理
		p.wg.Go(func() {
			p.worker(workerID)
		})
	}
}

// worker 工作协程
func (p *WorkerPool) worker(id int) {
	for task := range p.tasks {
		start := time.Now()

		// 模拟工作
		output := p.processTask(task)

		p.results <- Result{
			TaskID:   task.ID,
			Output:   output,
			Duration: time.Since(start),
		}

		log.Printf("[Worker %d] Completed task %d in %v",
			id, task.ID, time.Since(start))
	}
}

// processTask 处理任务
func (p *WorkerPool) processTask(task Task) string {
	// 模拟随机处理时间
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	return fmt.Sprintf("Processed: %s", task.Data)
}

// Submit 提交任务
func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}

// Shutdown 关闭工作池
func (p *WorkerPool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

// Results 获取结果channel
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

func main() {
	// 创建工作池（4个worker）
	pool := NewWorkerPool(4)
	pool.Start()

	// 提交任务
	go func() {
		for i := 1; i <= 20; i++ {
			task := Task{
				ID:   i,
				Data: fmt.Sprintf("Task-%d", i),
			}
			pool.Submit(task)
			log.Printf("Submitted task %d", i)
		}
		pool.Shutdown()
	}()

	// 收集结果
	var totalDuration time.Duration
	count := 0

	for result := range pool.Results() {
		fmt.Printf("✅ Task %d: %s (took %v)\n",
			result.TaskID, result.Output, result.Duration)
		totalDuration += result.Duration
		count++
	}

	// 统计
	avgDuration := totalDuration / time.Duration(count)
	fmt.Printf("\n📊 Statistics:\n")
	fmt.Printf("  Total tasks: %d\n", count)
	fmt.Printf("  Total time: %v\n", totalDuration)
	fmt.Printf("  Average time: %v\n", avgDuration)
	fmt.Printf("  Workers: %d\n", pool.maxWorkers)
}
