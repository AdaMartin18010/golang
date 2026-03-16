package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	strategy := &ExponentialBackoff{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		MaxAttempts:  5,
	}

	delays := []time.Duration{
		strategy.NextDelay(1),
		strategy.NextDelay(2),
		strategy.NextDelay(3),
		strategy.NextDelay(4),
	}

	// 验证延迟递增
	for i := 1; i < len(delays); i++ {
		if delays[i] <= delays[i-1] {
			t.Errorf("Expected delay to increase, got %v <= %v", delays[i], delays[i-1])
		}
	}

	// 验证不超过最大延迟
	for _, delay := range delays {
		if delay > strategy.MaxDelay {
			t.Errorf("Expected delay <= %v, got %v", strategy.MaxDelay, delay)
		}
	}
}

func TestRetry_Success(t *testing.T) {
	ctx := context.Background()
	strategy := DefaultExponentialBackoff()

	attempts := 0
	err := Retry(ctx, strategy, func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_MaxAttempts(t *testing.T) {
	ctx := context.Background()
	strategy := &FixedBackoff{
		Delay:       10 * time.Millisecond,
		MaxAttempts: 3,
	}

	attempts := 0
	err := Retry(ctx, strategy, func(ctx context.Context) error {
		attempts++
		return errors.New("permanent error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	strategy := DefaultExponentialBackoff()

	cancel()

	err := Retry(ctx, strategy, func(ctx context.Context) error {
		return errors.New("error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRetryWithCallback(t *testing.T) {
	ctx := context.Background()
	strategy := DefaultExponentialBackoff()

	attempts := 0
	callbackCalls := 0

	err := RetryWithCallback(ctx, strategy, func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	}, func(attempt int, err error) {
		callbackCalls++
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if callbackCalls != 1 {
		t.Errorf("Expected 1 callback call, got %d", callbackCalls)
	}
}
