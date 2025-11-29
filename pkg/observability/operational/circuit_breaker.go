package operational

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CircuitState 熔断器状态
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"   // 关闭（正常）
	CircuitStateOpen     CircuitState = "open"     // 打开（熔断）
	CircuitStateHalfOpen CircuitState = "halfopen" // 半开（尝试恢复）
)

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name          string
	maxFailures   int
	timeout       time.Duration
	resetTimeout  time.Duration
	mu            sync.RWMutex
	state         CircuitState
	failures      int
	lastFailure   time.Time
	successCount  int
	halfOpenLimit int
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	Name          string
	MaxFailures   int           // 最大失败次数
	Timeout       time.Duration // 超时时间
	ResetTimeout  time.Duration // 重置超时时间
	HalfOpenLimit int           // 半开状态下的成功次数限制
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.MaxFailures == 0 {
		cfg.MaxFailures = 5
	}
	if cfg.ResetTimeout == 0 {
		cfg.ResetTimeout = 60 * time.Second
	}
	if cfg.HalfOpenLimit == 0 {
		cfg.HalfOpenLimit = 3
	}

	return &CircuitBreaker{
		name:          cfg.Name,
		maxFailures:   cfg.MaxFailures,
		timeout:       cfg.Timeout,
		resetTimeout:  cfg.ResetTimeout,
		state:         CircuitStateClosed,
		halfOpenLimit: cfg.HalfOpenLimit,
	}
}

// Execute 执行操作（带熔断保护）
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	state := cb.state
	cb.mu.Unlock()

	// 检查状态
	switch state {
	case CircuitStateOpen:
		// 检查是否可以进入半开状态
		cb.mu.Lock()
		if time.Since(cb.lastFailure) >= cb.resetTimeout {
			cb.state = CircuitStateHalfOpen
			cb.successCount = 0
			state = CircuitStateHalfOpen
		}
		cb.mu.Unlock()

		if state == CircuitStateOpen {
			return fmt.Errorf("circuit breaker %s is open", cb.name)
		}
	case CircuitStateHalfOpen:
		// 半开状态，允许尝试
	case CircuitStateClosed:
		// 关闭状态，正常执行
	}

	// 执行操作
	var err error
	if cb.timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, cb.timeout)
		defer cancel()
		err = fn()
		if timeoutCtx.Err() != nil {
			err = fmt.Errorf("operation timeout: %w", timeoutCtx.Err())
		}
	} else {
		err = fn()
	}

	// 更新状态
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		// 操作失败
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.state == CircuitStateHalfOpen {
			// 半开状态下失败，重新打开
			cb.state = CircuitStateOpen
		} else if cb.failures >= cb.maxFailures {
			// 达到最大失败次数，打开熔断器
			cb.state = CircuitStateOpen
		}
	} else {
		// 操作成功
		cb.failures = 0

		if cb.state == CircuitStateHalfOpen {
			cb.successCount++
			if cb.successCount >= cb.halfOpenLimit {
				// 半开状态下成功次数达到限制，关闭熔断器
				cb.state = CircuitStateClosed
				cb.successCount = 0
			}
		} else if cb.state == CircuitStateOpen {
			// 不应该发生，但处理一下
			cb.state = CircuitStateClosed
		}
	}

	return err
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitStateClosed
	cb.failures = 0
	cb.successCount = 0
	cb.lastFailure = time.Time{}
}
