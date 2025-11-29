package security

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  5,
		Window: 1 * time.Second,
	}

	rl := NewRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()
	key := "test-key"

	// 允许 5 次请求
	for i := 0; i < 5; i++ {
		allowed, err := rl.Allow(ctx, key)
		if err != nil {
			t.Fatalf("Failed to allow request: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 第 6 次应该被拒绝
	allowed, err := rl.Allow(ctx, key)
	if err == nil {
		t.Error("Should return error for exceeded rate limit")
	}
	if allowed {
		t.Error("Request should be denied")
	}
}

func TestRateLimiter_Check(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  3,
		Window: 1 * time.Second,
	}

	rl := NewRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()
	key := "test-key"

	// 记录 3 次请求
	for i := 0; i < 3; i++ {
		rl.Allow(ctx, key)
	}

	// 检查应该返回 false（不记录）
	allowed, err := rl.Check(ctx, key)
	if err != nil {
		t.Fatalf("Failed to check: %v", err)
	}
	if allowed {
		t.Error("Check should return false when limit exceeded")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  2,
		Window: 1 * time.Second,
	}

	rl := NewRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()
	key := "test-key"

	// 使用所有限制
	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	// 应该被拒绝
	allowed, _ := rl.Allow(ctx, key)
	if allowed {
		t.Error("Request should be denied")
	}

	// 重置
	rl.Reset(ctx, key)

	// 应该可以再次请求
	allowed, err := rl.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Failed to allow after reset: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed after reset")
	}
}

func TestRateLimiter_GetRemaining(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  10,
		Window: 1 * time.Second,
	}

	rl := NewRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()
	key := "test-key"

	// 初始剩余应该是 10
	remaining, err := rl.GetRemaining(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get remaining: %v", err)
	}
	if remaining != 10 {
		t.Errorf("Expected remaining 10, got %d", remaining)
	}

	// 使用 3 次
	rl.Allow(ctx, key)
	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	// 剩余应该是 7
	remaining, err = rl.GetRemaining(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get remaining: %v", err)
	}
	if remaining != 7 {
		t.Errorf("Expected remaining 7, got %d", remaining)
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  2,
		Window: 100 * time.Millisecond,
	}

	rl := NewRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()
	key := "test-key"

	// 使用所有限制
	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	// 应该被拒绝
	allowed, _ := rl.Allow(ctx, key)
	if allowed {
		t.Error("Request should be denied")
	}

	// 等待窗口过期
	time.Sleep(150 * time.Millisecond)

	// 应该可以再次请求
	allowed, err := rl.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Failed to allow after window expiry: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed after window expiry")
	}
}

func TestIPRateLimiter(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  5,
		Window: 1 * time.Second,
	}

	rl := NewIPRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()

	// 测试 IP 限制
	allowed, err := rl.AllowIP(ctx, "192.168.1.1")
	if err != nil {
		t.Fatalf("Failed to allow IP: %v", err)
	}
	if !allowed {
		t.Error("IP request should be allowed")
	}
}

func TestUserRateLimiter(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  10,
		Window: 1 * time.Second,
	}

	rl := NewUserRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()

	// 测试用户限制
	allowed, err := rl.AllowUser(ctx, "user-123")
	if err != nil {
		t.Fatalf("Failed to allow user: %v", err)
	}
	if !allowed {
		t.Error("User request should be allowed")
	}
}

func TestEndpointRateLimiter(t *testing.T) {
	config := RateLimiterConfig{
		Limit:  20,
		Window: 1 * time.Second,
	}

	rl := NewEndpointRateLimiter(config)
	defer rl.Shutdown(context.Background())

	ctx := context.Background()

	// 测试端点限制
	allowed, err := rl.AllowEndpoint(ctx, "/api/users")
	if err != nil {
		t.Fatalf("Failed to allow endpoint: %v", err)
	}
	if !allowed {
		t.Error("Endpoint request should be allowed")
	}
}
