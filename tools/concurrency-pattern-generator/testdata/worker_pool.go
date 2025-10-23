package main

// Pattern: Worker Pool
// CSP Model: Pool = worker₁ || worker₂ || ... || worker8
//
// Safety Properties:
//   - Deadlock-free: ✓ (All workers can terminate when done channel closed)
//   - Race-free: ✓ (Channel synchronization guarantees happens-before)
//
// Theory: 文档02 第3.2节, 文档16 第1.1节
//
// Happens-Before Relations:
//   1. job sent → job received by worker
//   2. job processed → result sent
//   3. done channel closed → all workers exit
//   4. all workers exit → results channel closed
//
// Formal Proof:
//   ∀ job ∈ Jobs: 
//     sent(job) →ʰᵇ received(job) →ʰᵇ processed(job) →ʰᵇ result_sent(job)

import (
	"context"
	"sync"
)

// Job 表示待处理的任务
type Job struct {
	ID   int
	Data interface{}
}

// Result 表示处理结果
type Result struct {
	JobID int
	Data  interface{}
	Error error
}

// WorkerPool 创建一个工作池来并发处理任务
//
// 参数:
//   - ctx: 上下文，用于取消和超时控制
//   - numWorkers: worker数量
//   - jobs: 任务输入channel
//
// 返回:
//   - <-chan Result: 结果输出channel
//
// 使用示例:
//   ctx := context.Background()
//   jobs := make(chan Job, 100)
//   results := WorkerPool(ctx, 8, jobs)
//
//   // 发送任务
//   for i := 0; i < 100; i++ {
//       jobs <- Job{ID: i, Data: i}
//   }
//   close(jobs)
//
//   // 接收结果
//   for result := range results {
//       fmt.Printf("Result: %+v\n", result)
//   }
func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup

	// 启动worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			// Worker循环
			for {
				select {
				case <-ctx.Done():
					// Context取消，退出
					return
					
				case job, ok := <-jobs:
					if !ok {
						// Jobs channel关闭，退出
						return
					}
					
					// 处理任务
					result := processJob(job)
					
					// 发送结果
					select {
					case results <- result:
						// 成功发送
					case <-ctx.Done():
						// Context取消，退出
						return
					}
				}
			}
		}(i)
	}

	// 等待所有workers完成后关闭results
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// processJob 处理单个任务（用户需要实现此函数）
func processJob(job Job) Result {
	// TODO: 用户实现具体的任务处理逻辑
	return Result{
		JobID: job.ID,
		Data:  job.Data,
		Error: nil,
	}
}

// CSP形式化定义：
//
// Process Definition:
//   Worker(i) = jobs?job → process(job) → results!result → Worker(i)
//             □ done → SKIP
//             □ ctx.Done → SKIP
//
//   Pool = Worker(1) ||| Worker(2) ||| ... ||| Worker(n)
//
// Trace Semantics:
//   traces(Pool) = { t | t is interleaving of traces(Worker(i)) }
//
// Deadlock Freedom:
//   Pool is deadlock-free because:
//     1. Each worker can independently terminate
//     2. No circular dependencies
//     3. Channels are properly closed
//
// Safety:
//   ∀ t ∈ traces(Pool): 
//     safe(t) ⟺ ¬(∃ job: received(job) ∧ ¬sent(job))
