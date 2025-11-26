package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// CircuitState 熔断器状态
type CircuitState int

const (
	StateClosed CircuitState = iota // 关闭状态（正常）
	StateOpen                       // 开启状态（熔断）
	StateHalfOpen                   // 半开状态（尝试恢复）
)

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	FailureThreshold   int           // 失败阈值
	SuccessThreshold   int           // 成功阈值（半开状态下）
	Timeout            time.Duration // 熔断持续时间
	TimeoutWindow      time.Duration // 时间窗口
	OnStateChange      func(string, CircuitState) // 状态变更回调
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name              string
	config            CircuitBreakerConfig
	state             CircuitState
	failures          int
	successes         int
	lastFailureTime   time.Time
	lastStateChange   time.Time
	mu                sync.RWMutex
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, config CircuitBreakerConfig) *CircuitBreaker {
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 5
	}
	if config.SuccessThreshold == 0 {
		config.SuccessThreshold = 2
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	if config.TimeoutWindow == 0 {
		config.TimeoutWindow = 60 * time.Second
	}

	return &CircuitBreaker{
		name:            name,
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Allow 检查是否允许请求
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	// 检查是否需要重置计数器（时间窗口）
	if cb.state == StateClosed && now.Sub(cb.lastStateChange) > cb.config.TimeoutWindow {
		cb.failures = 0
		cb.lastStateChange = now
	}

	// 根据状态决定是否允许
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否应该进入半开状态
		if now.Sub(cb.lastFailureTime) >= cb.config.Timeout {
			cb.setState(StateHalfOpen)
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// OnSuccess 记录成功
func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// 重置失败计数
		cb.failures = 0
	case StateHalfOpen:
		cb.successes++
		// 如果成功次数达到阈值，关闭熔断器
		if cb.successes >= cb.config.SuccessThreshold {
			cb.setState(StateClosed)
		}
	}
}

// OnFailure 记录失败
func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		cb.failures++
		// 如果失败次数达到阈值，开启熔断器
		if cb.failures >= cb.config.FailureThreshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		// 半开状态下失败，重新开启熔断器
		cb.setState(StateOpen)
	}
}

// setState 设置状态
func (cb *CircuitBreaker) setState(newState CircuitState) {
	if cb.state != newState {
		oldState := cb.state
		cb.state = newState
		cb.lastStateChange = time.Now()
		cb.failures = 0
		cb.successes = 0

		if cb.config.OnStateChange != nil {
			cb.config.OnStateChange(cb.name, newState)
		}
	}
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// CircuitBreakerMiddleware 熔断器中间件
func CircuitBreakerMiddleware(breaker *CircuitBreaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否允许请求
			if !breaker.Allow() {
				response.Error(w, http.StatusServiceUnavailable,
					errors.NewServiceUnavailableError("circuit breaker is open"))
				return
			}

			// 创建响应包装器以捕获状态码
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// 执行下一个处理器
			next.ServeHTTP(ww, r)

			// 根据响应状态码记录成功或失败
			if ww.statusCode >= 200 && ww.statusCode < 500 {
				breaker.OnSuccess()
			} else {
				breaker.OnFailure()
			}
		})
	}
}

// responseWriter 响应包装器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CircuitBreakerManager 熔断器管理器
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager 创建熔断器管理器
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate 获取或创建熔断器
func (m *CircuitBreakerManager) GetOrCreate(name string, config CircuitBreakerConfig) *CircuitBreaker {
	m.mu.RLock()
	breaker, exists := m.breakers[name]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if breaker, exists := m.breakers[name]; exists {
		return breaker
	}

	breaker = NewCircuitBreaker(name, config)
	m.breakers[name] = breaker
	return breaker
}

// Get 获取熔断器
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	breaker, exists := m.breakers[name]
	return breaker, exists
}
