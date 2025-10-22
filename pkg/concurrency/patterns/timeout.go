package patterns

import (
	"context"
	"fmt"
	"time"
)

// WithTimeout 在指定时间内执行任务
func WithTimeoutFunc(timeout time.Duration, task func() (interface{}, error)) (interface{}, error) {
	done := make(chan struct{})
	var result interface{}
	var err error

	go func() {
		result, err = task()
		close(done)
	}()

	select {
	case <-done:
		return result, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("operation timed out after %v", timeout)
	}
}

// TimeoutChannel 创建一个超时channel
func TimeoutChannel(timeout time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		close(ch)
	}()
	return ch
}

// WithDeadline 在截止时间前执行任务
func WithDeadlineFunc(deadline time.Time, task func() (interface{}, error)) (interface{}, error) {
	timeout := time.Until(deadline)
	if timeout <= 0 {
		return nil, fmt.Errorf("deadline already passed")
	}
	return WithTimeoutFunc(timeout, task)
}

// TimeoutRetry 带超时的重试机制
func TimeoutRetry(maxRetries int, timeout time.Duration, task func() error) error {
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- task()
		}()

		select {
		case err := <-done:
			if err == nil {
				return nil
			}
			// 如果不是最后一次重试，继续
			if i < maxRetries-1 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			return err
		case <-ctx.Done():
			if i < maxRetries-1 {
				continue
			}
			return fmt.Errorf("all retries timed out")
		}
	}
	return fmt.Errorf("max retries exceeded")
}

// BatchWithTimeout 批量执行任务，每个任务有超时控制
func BatchWithTimeout(tasks []func() error, timeout time.Duration) []error {
	results := make([]error, len(tasks))
	done := make(chan struct {
		index int
		err   error
	})

	for i, task := range tasks {
		go func(index int, t func() error) {
			_, err := WithTimeoutFunc(timeout, func() (interface{}, error) {
				return nil, t()
			})
			done <- struct {
				index int
				err   error
			}{index, err}
		}(i, task)
	}

	for range tasks {
		result := <-done
		results[result.index] = result.err
	}

	return results
}

// CircuitBreaker 断路器，防止级联失败
type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration
	failures     int
	lastFailTime time.Time
	state        string // "closed", "open", "half-open"
}

// NewCircuitBreaker 创建断路器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

// Execute 执行任务，如果断路器打开则拒绝
func (cb *CircuitBreaker) Execute(task func() error) error {
	// 检查是否需要重置
	if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
		cb.state = "half-open"
		cb.failures = 0
	}

	// 如果断路器打开，拒绝请求
	if cb.state == "open" {
		return fmt.Errorf("circuit breaker is open")
	}

	// 执行任务
	err := task()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()

		if cb.failures >= cb.maxFailures {
			cb.state = "open"
		}
		return err
	}

	// 成功后重置
	if cb.state == "half-open" {
		cb.state = "closed"
	}
	cb.failures = 0
	return nil
}

// GetState 获取断路器状态
func (cb *CircuitBreaker) GetState() string {
	if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
		return "half-open"
	}
	return cb.state
}
