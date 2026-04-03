package limiter

import (
	"context"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 1) // 10 capacity, 1 token/sec

	// Should allow 10 requests immediately
	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// 11th request should be denied
	if tb.Allow() {
		t.Error("Expected 11th request to be denied")
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	// Should allow 1 more request
	if !tb.Allow() {
		t.Error("Expected request after refill to be allowed")
	}
}

func TestTokenBucketAllowN(t *testing.T) {
	tb := NewTokenBucket(10, 10)

	// Should allow 5 requests at once
	if !tb.AllowN(5) {
		t.Error("Expected AllowN(5) to be allowed")
	}

	// Remaining: 5
	if !tb.AllowN(5) {
		t.Error("Expected AllowN(5) to be allowed (second)")
	}

	// Remaining: 0
	if tb.AllowN(1) {
		t.Error("Expected AllowN(1) to be denied")
	}
}

func TestTokenBucketWait(t *testing.T) {
	tb := NewTokenBucket(1, 10) // 1 capacity, 10 tokens/sec

	// Take the only token
	tb.Allow()

	// Try to wait for next token
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := tb.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected Wait to succeed, got error: %v", err)
	}

	// Should have waited approximately 100ms (for 1 token at 10 tokens/sec)
	if elapsed < 80*time.Millisecond || elapsed > 150*time.Millisecond {
		t.Errorf("Expected wait time ~100ms, got %v", elapsed)
	}
}

func TestSlidingWindow(t *testing.T) {
	sw := NewSlidingWindow(10, time.Second)

	// Should allow 10 requests
	for i := 0; i < 10; i++ {
		if !sw.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// 11th request should be denied
	if sw.Allow() {
		t.Error("Expected 11th request to be denied")
	}
}

func TestSlidingWindowWeighted(t *testing.T) {
	sw := NewSlidingWindow(10, time.Second)

	// Add 5 requests
	for i := 0; i < 5; i++ {
		sw.Allow()
	}

	current, previous, weighted := sw.Stats()

	if current != 5 {
		t.Errorf("Expected current window count to be 5, got %d", current)
	}

	if previous != 0 {
		t.Errorf("Expected previous window count to be 0, got %d", previous)
	}

	if weighted < 5 || weighted > 6 {
		t.Errorf("Expected weighted count to be ~5, got %d", weighted)
	}
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := NewTokenBucket(1000000, 1000000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}

func BenchmarkSlidingWindow(b *testing.B) {
	sw := NewSlidingWindow(1000000, time.Second)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sw.Allow()
		}
	})
}

func BenchmarkTokenBucketAllowN(b *testing.B) {
	tb := NewTokenBucket(1000000, 1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.AllowN(10)
	}
}
