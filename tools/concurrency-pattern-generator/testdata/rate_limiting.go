package main

import (
	"context"
	"time"
)

// RateLimiter 简单的限流器
type RateLimiter struct {
	ticker *time.Ticker
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(rate),
	}
}

// Wait 等待令牌
func (r *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-r.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close 关闭限流器
func (r *RateLimiter) Close() {
	r.ticker.Stop()
}

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	tokens chan struct{}
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		tokens: make(chan struct{}, capacity),
	}
	
	// 填充初始令牌
	for i := 0; i < capacity; i++ {
		tb.tokens <- struct{}{}
	}
	
	// 定期补充令牌
	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case tb.tokens <- struct{}{}:
			default:
			}
		}
	}()
	
	return tb
}

// Acquire 获取令牌
func (tb *TokenBucket) Acquire(ctx context.Context) error {
	select {
	case <-tb.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
