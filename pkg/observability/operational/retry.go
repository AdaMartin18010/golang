package operational

import (
	"context"
	"fmt"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       bool // 是否添加随机抖动
}

// DefaultRetryConfig 返回默认重试配置
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}
}

// Retry 执行带重试的操作
func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.InitialDelay == 0 {
		cfg.InitialDelay = 100 * time.Millisecond
	}
	if cfg.Multiplier == 0 {
		cfg.Multiplier = 2.0
	}

	var lastErr error
	delay := cfg.InitialDelay

	for i := 0; i <= cfg.MaxRetries; i++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 执行操作
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// 如果不是最后一次重试，等待后重试
		if i < cfg.MaxRetries {
			// 计算延迟时间
			actualDelay := delay
			if cfg.Jitter {
				// 添加 ±10% 的随机抖动
				actualDelay = addJitter(delay, 0.1)
			}

			// 限制最大延迟
			if actualDelay > cfg.MaxDelay {
				actualDelay = cfg.MaxDelay
			}

			// 等待
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(actualDelay):
			}

			// 指数退避
			delay = time.Duration(float64(delay) * cfg.Multiplier)
		}
	}

	return fmt.Errorf("retry failed after %d attempts: %w", cfg.MaxRetries+1, lastErr)
}

// addJitter 添加随机抖动
func addJitter(delay time.Duration, jitterPercent float64) time.Duration {
	// 简化实现，实际应该使用随机数
	// 这里返回原值，实际实现应该添加随机抖动
	return delay
}

// RetryableError 可重试错误接口
type RetryableError interface {
	error
	IsRetryable() bool
}

// ShouldRetry 判断是否应该重试
func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if retryableErr, ok := err.(RetryableError); ok {
		return retryableErr.IsRetryable()
	}

	// 默认策略：某些错误类型应该重试
	// 这里简化实现，实际应该根据错误类型判断
	return true
}
