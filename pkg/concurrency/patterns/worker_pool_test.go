package patterns

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestBasicWorkerPool 测试基本Worker Pool
func TestBasicWorkerPool(t *testing.T) {
	numWorkers := 3
	numTasks := 10

	tasks := make(chan int, numTasks)
	results := make(chan int, numTasks)

	// 创建worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				// 处理任务
				result := task * 2
				results <- result
			}
		}(i)
	}

	// 发送任务
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}
	close(tasks)

	// 等待所有worker完成
	wg.Wait()
	close(results)

	// 验证结果
	sum := 0
	expectedSum := 0
	for i := 1; i <= numTasks; i++ {
		expectedSum += i * 2
	}

	for result := range results {
		sum += result
	}

	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}
}

// TestWorkerPoolWithContext 测试带Context的Worker Pool
func TestWorkerPoolWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	numWorkers := 2
	tasks := make(chan int)
	var processedCount atomic.Int32

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					time.Sleep(20 * time.Millisecond)
					processedCount.Add(1)
					_ = task
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// 发送任务
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case tasks <- i:
			case <-ctx.Done():
				close(tasks)
				return
			}
		}
		close(tasks)
	}()

	// 等待完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Test timeout")
	}

	// 由于超时，不是所有任务都会被处理
	processed := processedCount.Load()
	if processed == 0 {
		t.Error("Expected some tasks to be processed")
	}
	if processed > 10 {
		t.Errorf("Processed more tasks than sent: %d", processed)
	}
}

// TestWorkerPoolLoadBalancing 测试负载均衡
func TestWorkerPoolLoadBalancing(t *testing.T) {
	numWorkers := 3
	numTasks := 30

	tasks := make(chan int, numTasks)
	workerCounts := make([]atomic.Int32, numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()
			for task := range tasks {
				workerCounts[workerID].Add(1)
				time.Sleep(1 * time.Millisecond)
				_ = task
			}
		}()
	}

	// 发送任务
	for i := 0; i < numTasks; i++ {
		tasks <- i
	}
	close(tasks)

	wg.Wait()

	// 验证负载分布（每个worker应该处理大致相同数量的任务）
	for i := range workerCounts {
		workerTasks := workerCounts[i].Load()
		expectedMin := int32(numTasks/numWorkers - 2)
		expectedMax := int32(numTasks/numWorkers + 2)

		if workerTasks < expectedMin || workerTasks > expectedMax {
			t.Errorf("Worker %d processed %d tasks, expected between %d and %d",
				i, workerTasks, expectedMin, expectedMax)
		}
	}
}

// TestWorkerPoolErrorHandling 测试错误处理
func TestWorkerPoolErrorHandling(t *testing.T) {
	type task struct {
		id    int
		value int
	}

	type result struct {
		id  int
		err error
	}

	numWorkers := 2
	tasks := make(chan task, 10)
	results := make(chan result, 10)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasks {
				// 负值产生错误
				if t.value < 0 {
					results <- result{id: t.id, err: &workerError{"negative value"}}
				} else {
					results <- result{id: t.id, err: nil}
				}
			}
		}()
	}

	// 发送任务（包含一些会出错的）
	testTasks := []task{
		{1, 10},
		{2, -5}, // 错误
		{3, 20},
		{4, -10}, // 错误
		{5, 30},
	}

	for _, t := range testTasks {
		tasks <- t
	}
	close(tasks)

	// 等待并关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	errorCount := 0
	successCount := 0
	for r := range results {
		if r.err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	if errorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", errorCount)
	}
	if successCount != 3 {
		t.Errorf("Expected 3 successes, got %d", successCount)
	}
}

type workerError struct {
	msg string
}

func (e *workerError) Error() string {
	return e.msg
}

// TestWorkerPoolGracefulShutdown 测试优雅关闭
func TestWorkerPoolGracefulShutdown(t *testing.T) {
	tasks := make(chan int, 100)
	var processedCount atomic.Int32
	shutdown := make(chan struct{})

	numWorkers := 3
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					processedCount.Add(1)
					time.Sleep(5 * time.Millisecond)
					_ = task
				case <-shutdown:
					// 优雅关闭：完成当前任务后退出
					return
				}
			}
		}()
	}

	// 发送一些任务
	numTasks := 20
	for i := 0; i < numTasks; i++ {
		tasks <- i
	}

	// 等待一些任务被处理
	time.Sleep(50 * time.Millisecond)

	// 触发关闭
	close(shutdown)
	close(tasks)

	// 等待所有worker退出
	wg.Wait()

	processed := processedCount.Load()
	if processed == 0 {
		t.Error("Expected some tasks to be processed")
	}

	t.Logf("Processed %d/%d tasks before shutdown", processed, numTasks)
}

// TestWorkerPoolDynamicSize 测试动态调整worker数量
func TestWorkerPoolDynamicSize(t *testing.T) {
	tasks := make(chan int, 100)
	var processedCount atomic.Int32
	var activeWorkers atomic.Int32

	var wg sync.WaitGroup
	done := make(chan struct{})

	startWorker := func() {
		wg.Add(1)
		activeWorkers.Add(1)
		go func() {
			defer wg.Done()
			defer activeWorkers.Add(-1)
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					processedCount.Add(1)
					time.Sleep(10 * time.Millisecond)
					_ = task
				case <-done:
					return
				}
			}
		}()
	}

	// 初始2个worker
	for i := 0; i < 2; i++ {
		startWorker()
	}

	// 发送任务
	go func() {
		for i := 0; i < 50; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	// 根据负载动态增加worker
	time.Sleep(30 * time.Millisecond)
	if activeWorkers.Load() == 2 {
		// 增加1个worker
		startWorker()
	}

	// 等待完成
	wg.Wait()
	close(done)

	processed := processedCount.Load()
	if processed != 50 {
		t.Errorf("Expected 50 tasks processed, got %d", processed)
	}
}

// BenchmarkWorkerPool 基准测试
func BenchmarkWorkerPool(b *testing.B) {
	numWorkers := 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tasks := make(chan int, 100)
		results := make(chan int, 100)

		var wg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range tasks {
					results <- task * 2
				}
			}()
		}

		go func() {
			for j := 0; j < 100; j++ {
				tasks <- j
			}
			close(tasks)
		}()

		go func() {
			wg.Wait()
			close(results)
		}()

		for range results {
		}
	}
}
