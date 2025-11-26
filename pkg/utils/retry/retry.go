package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// Strategy 重试策略
type Strategy interface {
	NextDelay(attempt int) time.Duration
	ShouldRetry(attempt int, err error) bool
}

// ExponentialBackoff 指数退避策略
type ExponentialBackoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	MaxAttempts  int
}

// NextDelay 计算下次延迟
func (e *ExponentialBackoff) NextDelay(attempt int) time.Duration {
	delay := float64(e.InitialDelay) * math.Pow(e.Multiplier, float64(attempt-1))
	if delay > float64(e.MaxDelay) {
		delay = float64(e.MaxDelay)
	}
	return time.Duration(delay)
}

// ShouldRetry 判断是否应该重试
func (e *ExponentialBackoff) ShouldRetry(attempt int, err error) bool {
	return attempt < e.MaxAttempts
}

// LinearBackoff 线性退避策略
type LinearBackoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Increment    time.Duration
	MaxAttempts  int
}

// NextDelay 计算下次延迟
func (l *LinearBackoff) NextDelay(attempt int) time.Duration {
	delay := l.InitialDelay + l.Increment*time.Duration(attempt-1)
	if delay > l.MaxDelay {
		delay = l.MaxDelay
	}
	return delay
}

// ShouldRetry 判断是否应该重试
func (l *LinearBackoff) ShouldRetry(attempt int, err error) bool {
	return attempt < l.MaxAttempts
}

// FixedBackoff 固定延迟策略
type FixedBackoff struct {
	Delay       time.Duration
	MaxAttempts int
}

// NextDelay 计算下次延迟
func (f *FixedBackoff) NextDelay(attempt int) time.Duration {
	return f.Delay
}

// ShouldRetry 判断是否应该重试
func (f *FixedBackoff) ShouldRetry(attempt int, err error) bool {
	return attempt < f.MaxAttempts
}

// RetryableError 可重试错误
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// NewRetryableError 创建可重试错误
func NewRetryableError(err error) *RetryableError {
	return &RetryableError{Err: err}
}

// RetryFunc 重试函数类型
type RetryFunc func(ctx context.Context) error

// Retry 执行重试
func Retry(ctx context.Context, strategy Strategy, fn RetryFunc) error {
	attempt := 0
	var lastErr error

	for {
		attempt++
		lastErr = fn(ctx)

		if lastErr == nil {
			return nil
		}

		// 检查是否可重试
		if !strategy.ShouldRetry(attempt, lastErr) {
			return fmt.Errorf("max attempts reached: %w", lastErr)
		}

		// 检查context是否取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		// 等待延迟
		delay := strategy.NextDelay(attempt)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
}

// RetryWithCallback 带回调的重试
func RetryWithCallback(
	ctx context.Context,
	strategy Strategy,
	fn RetryFunc,
	onRetry func(attempt int, err error),
) error {
	attempt := 0
	var lastErr error

	for {
		attempt++
		lastErr = fn(ctx)

		if lastErr == nil {
			return nil
		}

		// 检查是否可重试
		if !strategy.ShouldRetry(attempt, lastErr) {
			return fmt.Errorf("max attempts reached: %w", lastErr)
		}

		// 执行回调
		if onRetry != nil {
			onRetry(attempt, lastErr)
		}

		// 检查context是否取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		// 等待延迟
		delay := strategy.NextDelay(attempt)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
}

// DefaultExponentialBackoff 默认指数退避策略
func DefaultExponentialBackoff() Strategy {
	return &ExponentialBackoff{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		MaxAttempts:  5,
	}
}

// DefaultLinearBackoff 默认线性退避策略
func DefaultLinearBackoff() Strategy {
	return &LinearBackoff{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Increment:    100 * time.Millisecond,
		MaxAttempts:  5,
	}
}

// DefaultFixedBackoff 默认固定延迟策略
func DefaultFixedBackoff() Strategy {
	return &FixedBackoff{
		Delay:       1 * time.Second,
		MaxAttempts: 3,
	}
}
