package patterns

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// 创建限流器: 每秒10个令牌，容量20
	rl := NewRateLimiter(10, 20)

	// 快速消耗所有令牌
	for i := 0; i < 20; i++ {
		if !rl.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 下一个请求应该被拒绝
	if rl.Allow() {
		t.Error("Request should be rejected when bucket is empty")
	}

	// 等待令牌补充
	time.Sleep(200 * time.Millisecond) // 应该补充2个令牌

	// 现在应该有新令牌
	if !rl.Allow() {
		t.Error("Request should be allowed after refill")
	}
}

func TestRateLimiterWait(t *testing.T) {
	rl := NewRateLimiter(100, 1) // 每秒100个，容量1

	rl.Allow() // 消耗唯一的令牌

	start := time.Now()
	rl.Wait() // 应该等待直到新令牌产生
	elapsed := time.Since(start)

	// 应该等待大约10ms (1/100秒)
	if elapsed < 5*time.Millisecond {
		t.Errorf("Wait should take at least 5ms, took %v", elapsed)
	}
}

func TestRateLimiterWaitWithContext(t *testing.T) {
	rl := NewRateLimiter(10, 1)
	rl.Allow() // 消耗令牌

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := rl.WaitWithContext(ctx)
	if err == nil {
		t.Error("Expected context timeout error")
	}
}

func TestLeakyBucket(t *testing.T) {
	// 创建漏桶: 容量10，每10ms漏一个
	lb := NewLeakyBucket(10, 10*time.Millisecond)

	// 快速填满桶
	for i := 0; i < 10; i++ {
		if !lb.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 桶满了，应该拒绝
	if lb.Allow() {
		t.Error("Request should be rejected when bucket is full")
	}

	// 等待漏水
	time.Sleep(50 * time.Millisecond) // 应该漏出5个

	// 现在应该有空间
	if !lb.Allow() {
		t.Error("Request should be allowed after leak")
	}
}

func TestSlidingWindowLimiter(t *testing.T) {
	// 创建滑动窗口限流器: 1秒内最多5个请求
	sw := NewSlidingWindowLimiter(5, 1*time.Second)

	// 快速发送5个请求
	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 第6个请求应该被拒绝
	if sw.Allow() {
		t.Error("Request should be rejected when limit exceeded")
	}

	// 等待窗口滑动
	time.Sleep(1100 * time.Millisecond)

	// 现在应该允许新请求
	if !sw.Allow() {
		t.Error("Request should be allowed after window slides")
	}
}

func TestSlidingWindowLimiterReset(t *testing.T) {
	sw := NewSlidingWindowLimiter(5, 1*time.Second)

	// 发送一些请求
	for i := 0; i < 3; i++ {
		sw.Allow()
	}

	// 重置
	sw.Reset()

	// 应该能再发送5个请求
	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("Request %d should be allowed after reset", i)
		}
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	rl := NewRateLimiter(100, 100)

	const numGoroutines = 10
	const requestsPerGoroutine = 20

	results := make(chan bool, numGoroutines*requestsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				results <- rl.Allow()
			}
		}()
	}

	allowed := 0
	for i := 0; i < numGoroutines*requestsPerGoroutine; i++ {
		if <-results {
			allowed++
		}
	}

	// 应该有大约100个请求被允许（桶的容量）
	if allowed < 90 || allowed > 110 {
		t.Errorf("Expected ~100 allowed requests, got %d", allowed)
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(10000, 10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow()
	}
}

func BenchmarkLeakyBucket(b *testing.B) {
	lb := NewLeakyBucket(10000, time.Microsecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Allow()
	}
}

func BenchmarkSlidingWindowLimiter(b *testing.B) {
	sw := NewSlidingWindowLimiter(10000, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Allow()
	}
}
