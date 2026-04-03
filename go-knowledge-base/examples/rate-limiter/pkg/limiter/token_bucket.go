package limiter

import (
	"context"
	"sync"
	"time"
)

// TokenBucket implements the token bucket rate limiting algorithm
type TokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a single request should be allowed
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN checks if n requests can be allowed
func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		return true
	}

	return false
}

// Wait blocks until a token is available or context is cancelled
func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		if tb.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tb.timeUntilNextToken()):
		}
	}
}

// refill adds tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now
}

// timeUntilNextToken returns the time until next token is available
func (tb *TokenBucket) timeUntilNextToken() time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens >= 1 {
		return 0
	}

	needed := 1.0 - tb.tokens
	seconds := needed / tb.refillRate
	return time.Duration(seconds * float64(time.Second))
}

// Tokens returns the current number of available tokens
func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}

// Capacity returns the bucket capacity
func (tb *TokenBucket) Capacity() float64 {
	return tb.capacity
}

// Rate returns the refill rate
func (tb *TokenBucket) Rate() float64 {
	return tb.refillRate
}
