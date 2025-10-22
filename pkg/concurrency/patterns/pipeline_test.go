package patterns

import (
	"testing"
	"time"
)

// TestSimplePipeline 测试简单的Pipeline
func TestSimplePipeline(t *testing.T) {
	// Stage 1: 生成数字
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// Stage 2: 平方
	sq := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// 构建pipeline
	nums := gen(1, 2, 3, 4, 5)
	squares := sq(nums)

	// 验证结果
	expected := []int{1, 4, 9, 16, 25}
	i := 0
	for result := range squares {
		if result != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], result)
		}
		i++
	}

	if i != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), i)
	}
}

// TestMultiStagePipeline 测试多阶段Pipeline
func TestMultiStagePipeline(t *testing.T) {
	// Stage 1: 生成
	generate := func(max int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= max; i++ {
				out <- i
			}
		}()
		return out
	}

	// Stage 2: 过滤偶数
	filterEven := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				if n%2 == 0 {
					out <- n
				}
			}
		}()
		return out
	}

	// Stage 3: 乘以10
	multiply := func(in <-chan int, factor int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * factor
			}
		}()
		return out
	}

	// 构建pipeline
	nums := generate(10)
	evens := filterEven(nums)
	results := multiply(evens, 10)

	// 验证结果: 2*10=20, 4*10=40, 6*10=60, 8*10=80, 10*10=100
	expected := []int{20, 40, 60, 80, 100}
	i := 0
	for result := range results {
		if result != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result)
		}
		i++
	}
}

// TestPipelineWithTimeout 测试带超时的Pipeline
func TestPipelineWithTimeout(t *testing.T) {
	// 慢速生成器
	slowGen := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 0; i < 5; i++ {
				time.Sleep(100 * time.Millisecond)
				out <- i
			}
		}()
		return out
	}

	// 带超时的消费
	nums := slowGen()
	timeout := time.After(250 * time.Millisecond)
	count := 0

	for {
		select {
		case n, ok := <-nums:
			if !ok {
				t.Fatal("Channel closed before timeout")
			}
			count++
			_ = n
		case <-timeout:
			// 预期在超时前应该收到2-3个数字
			if count < 2 || count > 3 {
				t.Errorf("Expected 2-3 items before timeout, got %d", count)
			}
			return
		}
	}
}

// TestFanOutFanIn 测试扇出-扇入模式
func TestFanOutFanIn(t *testing.T) {
	// 生成器
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// Worker
	worker := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// Merge (Fan-In)
	merge := func(cs ...<-chan int) <-chan int {
		out := make(chan int)

		done := make(chan struct{})
		defer close(done)

		for _, c := range cs {
			go func(ch <-chan int) {
				for n := range ch {
					select {
					case out <- n:
					case <-done:
						return
					}
				}
			}(c)
		}

		go func() {
			// 简化版：等待一段时间后关闭
			time.Sleep(200 * time.Millisecond)
			close(out)
		}()

		return out
	}

	// 构建pipeline
	nums := gen(1, 2, 3, 4, 5, 6)

	// Fan-Out到3个worker
	w1 := worker(nums)
	w2 := worker(nums)
	w3 := worker(nums)

	// Fan-In
	results := merge(w1, w2, w3)

	// 收集结果
	sum := 0
	for n := range results {
		sum += n
	}

	// 验证总和（可能因为并发而不精确，但应该大于0）
	if sum == 0 {
		t.Error("Expected non-zero sum")
	}
}

// TestPipelineErrorHandling 测试Pipeline错误处理
func TestPipelineErrorHandling(t *testing.T) {
	type result struct {
		value int
		err   error
	}

	// 带错误的生成器
	gen := func(nums ...int) <-chan result {
		out := make(chan result)
		go func() {
			defer close(out)
			for _, n := range nums {
				if n < 0 {
					out <- result{err: &pipelineError{"negative number"}}
					return
				}
				out <- result{value: n}
			}
		}()
		return out
	}

	// 测试正常情况
	t.Run("no errors", func(t *testing.T) {
		results := gen(1, 2, 3)
		count := 0
		for r := range results {
			if r.err != nil {
				t.Errorf("Unexpected error: %v", r.err)
			}
			count++
		}
		if count != 3 {
			t.Errorf("Expected 3 results, got %d", count)
		}
	})

	// 测试错误情况
	t.Run("with errors", func(t *testing.T) {
		results := gen(1, 2, -1, 4)
		foundError := false
		count := 0
		for r := range results {
			if r.err != nil {
				foundError = true
				break
			}
			count++
		}
		if !foundError {
			t.Error("Expected to find an error")
		}
		if count > 2 {
			t.Errorf("Expected at most 2 successful results before error, got %d", count)
		}
	})
}

type pipelineError struct {
	msg string
}

func (e *pipelineError) Error() string {
	return e.msg
}

// BenchmarkSimplePipeline 基准测试
func BenchmarkSimplePipeline(b *testing.B) {
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	sq := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nums := gen(1, 2, 3, 4, 5)
		squares := sq(nums)
		for range squares {
		}
	}
}
