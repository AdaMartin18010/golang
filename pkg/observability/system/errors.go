package system

import (
	"fmt"
	"time"
)

// MonitorError 监控错误
type MonitorError struct {
	Component string
	Operation string
	Message   string
	Retryable bool
	Timestamp time.Time
	Err       error
}

func (e *MonitorError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s.%s: %s: %v", e.Component, e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("%s.%s: %s", e.Component, e.Operation, e.Message)
}

func (e *MonitorError) Unwrap() error {
	return e.Err
}

// IsRetryable 检查错误是否可重试
func (e *MonitorError) IsRetryable() bool {
	return e.Retryable
}

// NewMonitorError 创建监控错误
func NewMonitorError(component, operation, message string, retryable bool, err error) *MonitorError {
	return &MonitorError{
		Component: component,
		Operation: operation,
		Message:   message,
		Retryable: retryable,
		Timestamp: time.Now(),
		Err:       err,
	}
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig 返回默认重试配置
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Retry 执行带重试的操作
func Retry(operation func() error, cfg RetryConfig) error {
	var lastErr error
	delay := cfg.InitialDelay

	for i := 0; i < cfg.MaxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// 检查是否可重试
		if monitorErr, ok := err.(*MonitorError); ok && !monitorErr.IsRetryable() {
			return err
		}

		// 如果不是最后一次重试，等待后重试
		if i < cfg.MaxRetries-1 {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", cfg.MaxRetries, lastErr)
}
