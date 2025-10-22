package patterns

import (
	"context"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	// 成功案例
	err := WithTimeout(context.Background(), 100*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 超时案例
	err = WithTimeout(context.Background(), 10*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestWithCancel(t *testing.T) {
	completed := false
	err, cancel := WithCancel(context.Background(), func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		completed = true
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !completed {
		t.Error("Task should have completed")
	}
	cancel()
}

func TestWithValue(t *testing.T) {
	key := "testKey"
	expectedValue := "testValue"
	var actualValue interface{}

	err := WithValue(context.Background(), key, expectedValue, func(ctx context.Context) error {
		actualValue = ctx.Value(key)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if actualValue != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, actualValue)
	}
}

func TestContextAwarePipeline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		source <- i
	}
	close(source)

	out := ContextAwarePipeline(ctx, source)

	results := []int{}
	for v := range out {
		results = append(results, v)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}
	if results[0] != 2 {
		t.Errorf("Expected first result to be 2, got %d", results[0])
	}
}

func TestContextAwarePipelineCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	source := make(chan int)
	go func() {
		for i := 1; i <= 100; i++ {
			source <- i
			time.Sleep(time.Millisecond)
		}
		close(source)
	}()

	out := ContextAwarePipeline(ctx, source)

	// 读取几个结果后取消
	count := 0
	for range out {
		count++
		if count >= 3 {
			cancel()
			break
		}
	}

	// 等待goroutine清理
	time.Sleep(50 * time.Millisecond)

	if count < 3 {
		t.Errorf("Expected at least 3 results before cancel, got %d", count)
	}
}
func TestContextAwareWorkerPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)

	results := ContextAwareWorkerPool(ctx, 3, jobs)

	count := 0
	for range results {
		count++
	}

	if count != 10 {
		t.Errorf("Expected 10 results, got %d", count)
	}
}

func TestContextAwareWorkerPoolCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	jobs := make(chan int)
	go func() {
		for i := 1; i <= 100; i++ {
			select {
			case jobs <- i:
				time.Sleep(10 * time.Millisecond)
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
	}()

	results := ContextAwareWorkerPool(ctx, 3, jobs)

	count := 0
	for range results {
		count++
	}

	// 应该在超时前只处理部分任务
	if count >= 100 {
		t.Error("Worker pool should have been cancelled before processing all jobs")
	}
}

func BenchmarkWithTimeout(b *testing.B) {
	task := func(ctx context.Context) error {
		time.Sleep(time.Microsecond)
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithTimeout(context.Background(), 10*time.Millisecond, task)
	}
}

func BenchmarkContextAwarePipeline(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		source := make(chan int, 10)
		for j := 0; j < 10; j++ {
			source <- j
		}
		close(source)

		out := ContextAwarePipeline(ctx, source)
		for range out {
		}
	}
}
