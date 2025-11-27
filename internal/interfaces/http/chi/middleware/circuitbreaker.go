package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// CircuitState 是熔断器的状态类型。
//
// 熔断器有三种状态：
// - StateClosed: 关闭状态（正常）
//   所有请求都允许通过
//   记录失败次数，达到阈值后切换到开启状态
// - StateOpen: 开启状态（熔断）
//   所有请求都被拒绝
//   经过超时时间后切换到半开状态
// - StateHalfOpen: 半开状态（尝试恢复）
//   允许少量请求通过以测试服务是否恢复
//   如果成功次数达到阈值，切换到关闭状态
//   如果失败，重新切换到开启状态
type CircuitState int

const (
	StateClosed CircuitState = iota // 关闭状态（正常）
	StateOpen                       // 开启状态（熔断）
	StateHalfOpen                   // 半开状态（尝试恢复）
)

// CircuitBreakerConfig 是熔断器的配置结构。
//
// 功能说明：
// - 配置熔断器的阈值和超时时间
// - 支持状态变更回调
//
// 字段说明：
// - FailureThreshold: 失败阈值（默认：5）
//   关闭状态下，失败次数达到此值时切换到开启状态
// - SuccessThreshold: 成功阈值（默认：2）
//   半开状态下，成功次数达到此值时切换到关闭状态
// - Timeout: 熔断持续时间（默认：60秒）
//   开启状态下，经过此时间后切换到半开状态
// - TimeoutWindow: 时间窗口（默认：60秒）
//   关闭状态下，在此时间窗口内统计失败次数
// - OnStateChange: 状态变更回调函数
//   当熔断器状态发生变化时调用
//
// 使用示例：
//
//	config := middleware.CircuitBreakerConfig{
//	    FailureThreshold: 5,
//	    SuccessThreshold: 2,
//	    Timeout:          60 * time.Second,
//	    OnStateChange: func(name string, state middleware.CircuitState) {
//	        log.Printf("Circuit breaker %s changed to state %d", name, state)
//	    },
//	}
type CircuitBreakerConfig struct {
	FailureThreshold   int           // 失败阈值
	SuccessThreshold   int           // 成功阈值（半开状态下）
	Timeout            time.Duration // 熔断持续时间
	TimeoutWindow      time.Duration // 时间窗口
	OnStateChange      func(string, CircuitState) // 状态变更回调
}

// CircuitBreaker 是熔断器的实现。
//
// 功能说明：
// - 实现熔断器模式，保护下游服务
// - 支持三种状态：关闭、开启、半开
// - 自动状态转换和恢复
//
// 字段说明：
// - name: 熔断器名称（用于标识和日志）
// - config: 熔断器配置
// - state: 当前状态
// - failures: 失败次数（关闭状态下）
// - successes: 成功次数（半开状态下）
// - lastFailureTime: 上次失败时间
// - lastStateChange: 上次状态变更时间
// - mu: 读写互斥锁（保证并发安全）
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

// NewCircuitBreaker 创建并初始化熔断器。
//
// 功能说明：
// - 创建新的熔断器实例
// - 设置默认配置值
// - 初始状态为关闭（StateClosed）
//
// 参数：
// - name: 熔断器名称（用于标识和日志）
// - config: 熔断器配置
//
// 返回：
// - *CircuitBreaker: 配置好的熔断器实例
//
// 默认配置：
// - FailureThreshold: 5（如果未设置）
// - SuccessThreshold: 2（如果未设置）
// - Timeout: 60秒（如果未设置）
// - TimeoutWindow: 60秒（如果未设置）
//
// 使用示例：
//
//	config := middleware.CircuitBreakerConfig{
//	    FailureThreshold: 5,
//	    Timeout:          60 * time.Second,
//	}
//	breaker := middleware.NewCircuitBreaker("user-service", config)
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

// Allow 检查是否允许请求通过。
//
// 功能说明：
// - 根据熔断器状态决定是否允许请求
// - 自动处理状态转换
//
// 状态行为：
// - StateClosed: 允许所有请求
//   如果时间窗口已过，重置失败计数
// - StateOpen: 拒绝所有请求
//   如果超时时间已过，切换到半开状态并允许请求
// - StateHalfOpen: 允许请求（用于测试恢复）
//
// 返回：
// - bool: 如果允许请求返回 true，否则返回 false
//
// 使用场景：
// - 在调用下游服务前检查
// - 如果返回 false，直接返回错误，不调用下游服务
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	// 检查是否需要重置计数器（时间窗口）
	// 在关闭状态下，如果时间窗口已过，重置失败计数
	if cb.state == StateClosed && now.Sub(cb.lastStateChange) > cb.config.TimeoutWindow {
		cb.failures = 0
		cb.lastStateChange = now
	}

	// 根据状态决定是否允许
	switch cb.state {
	case StateClosed:
		// 关闭状态：允许所有请求
		return true
	case StateOpen:
		// 开启状态：检查是否应该进入半开状态
		if now.Sub(cb.lastFailureTime) >= cb.config.Timeout {
			// 超时时间已过，切换到半开状态并允许请求
			cb.setState(StateHalfOpen)
			return true
		}
		// 仍在熔断期间，拒绝请求
		return false
	case StateHalfOpen:
		// 半开状态：允许请求（用于测试恢复）
		return true
	default:
		return false
	}
}

// OnSuccess 记录请求成功。
//
// 功能说明：
// - 在请求成功后调用
// - 根据状态更新计数器和状态
//
// 状态行为：
// - StateClosed: 重置失败计数（表示服务正常）
// - StateHalfOpen: 增加成功计数
//   如果成功次数达到阈值，切换到关闭状态
//
// 使用场景：
// - 在请求成功后调用
// - 用于更新熔断器状态
func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// 关闭状态：重置失败计数
		// 表示服务正常，清除之前的失败记录
		cb.failures = 0
	case StateHalfOpen:
		// 半开状态：增加成功计数
		cb.successes++
		// 如果成功次数达到阈值，关闭熔断器
		// 表示服务已恢复，可以正常处理请求
		if cb.successes >= cb.config.SuccessThreshold {
			cb.setState(StateClosed)
		}
	}
}

// OnFailure 记录请求失败。
//
// 功能说明：
// - 在请求失败后调用
// - 根据状态更新计数器和状态
//
// 状态行为：
// - StateClosed: 增加失败计数
//   如果失败次数达到阈值，切换到开启状态（熔断）
// - StateHalfOpen: 切换到开启状态
//   表示服务仍未恢复，继续熔断
//
// 使用场景：
// - 在请求失败后调用
// - 用于更新熔断器状态
func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		// 关闭状态：增加失败计数
		cb.failures++
		// 如果失败次数达到阈值，开启熔断器
		// 表示服务异常，需要熔断保护
		if cb.failures >= cb.config.FailureThreshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		// 半开状态下失败，重新开启熔断器
		// 表示服务仍未恢复，继续熔断
		cb.setState(StateOpen)
	}
}

// setState 设置熔断器状态（内部方法）。
//
// 功能说明：
// - 更新熔断器状态
// - 重置计数器和时间戳
// - 调用状态变更回调
//
// 参数：
// - newState: 新状态
//
// 注意事项：
// - 这是内部方法，不应直接调用
// - 状态变更时会重置失败和成功计数
// - 如果配置了回调函数，会调用回调通知状态变更
func (cb *CircuitBreaker) setState(newState CircuitState) {
	if cb.state != newState {
		oldState := cb.state
		cb.state = newState
		cb.lastStateChange = time.Now()
		cb.failures = 0      // 重置失败计数
		cb.successes = 0     // 重置成功计数

		// 调用状态变更回调
		if cb.config.OnStateChange != nil {
			cb.config.OnStateChange(cb.name, newState)
		}
	}
}

// GetState 获取熔断器的当前状态。
//
// 功能说明：
// - 返回熔断器的当前状态
// - 线程安全
//
// 返回：
// - CircuitState: 当前状态（StateClosed、StateOpen 或 StateHalfOpen）
//
// 使用场景：
// - 监控和日志记录
// - 健康检查
// - 调试和诊断
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// CircuitBreakerMiddleware 创建熔断器中间件。
//
// 功能说明：
// - 在 HTTP 请求处理中集成熔断器
// - 根据熔断器状态决定是否处理请求
// - 根据响应状态码记录成功或失败
//
// 工作流程：
// 1. 检查熔断器是否允许请求
// 2. 如果不允许，返回 503 Service Unavailable
// 3. 如果允许，创建响应包装器
// 4. 执行下一个处理器
// 5. 根据响应状态码记录成功或失败
//
// 成功/失败判断：
// - 2xx-4xx: 记录为成功（服务正常响应）
// - 5xx: 记录为失败（服务错误）
//
// 参数：
// - breaker: 熔断器实例
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	breaker := middleware.NewCircuitBreaker("user-service", config)
//	router.Use(middleware.CircuitBreakerMiddleware(breaker))
//
// 注意事项：
// - 熔断器应该针对特定的下游服务创建
// - 可以根据不同的路由使用不同的熔断器
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
			// 用于判断请求成功或失败
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// 执行下一个处理器
			next.ServeHTTP(ww, r)

			// 根据响应状态码记录成功或失败
			// 2xx-4xx: 成功（服务正常响应）
			// 5xx: 失败（服务错误）
			if ww.statusCode >= 200 && ww.statusCode < 500 {
				breaker.OnSuccess()
			} else {
				breaker.OnFailure()
			}
		})
	}
}

// responseWriter 是 HTTP 响应写入器的包装器。
//
// 功能说明：
// - 包装 http.ResponseWriter
// - 捕获响应状态码
// - 用于熔断器判断请求成功或失败
//
// 字段说明：
// - ResponseWriter: 底层的 HTTP 响应写入器
// - statusCode: 响应状态码（默认：200）
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 记录状态码并写入响应头。
//
// 功能说明：
// - 记录状态码
// - 调用底层 ResponseWriter 的 WriteHeader
//
// 参数：
// - code: HTTP 状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CircuitBreakerManager 是熔断器管理器。
//
// 功能说明：
// - 管理多个熔断器实例
// - 支持按名称获取或创建熔断器
// - 线程安全
//
// 使用场景：
// - 管理多个下游服务的熔断器
// - 按需创建熔断器
// - 统一管理熔断器生命周期
//
// 字段说明：
// - breakers: 熔断器映射表（名称 -> 熔断器实例）
// - mu: 读写互斥锁（保证并发安全）
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager 创建熔断器管理器。
//
// 功能说明：
// - 创建新的管理器实例
// - 初始化熔断器映射表
//
// 返回：
// - *CircuitBreakerManager: 配置好的管理器实例
//
// 使用示例：
//
//	manager := middleware.NewCircuitBreakerManager()
//	breaker := manager.GetOrCreate("user-service", config)
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate 获取或创建熔断器。
//
// 功能说明：
// - 如果熔断器已存在，直接返回
// - 如果不存在，创建新的熔断器并返回
// - 使用双重检查锁定模式保证线程安全
//
// 参数：
// - name: 熔断器名称
// - config: 熔断器配置（仅在创建时使用）
//
// 返回：
// - *CircuitBreaker: 熔断器实例
//
// 使用示例：
//
//	config := middleware.CircuitBreakerConfig{
//	    FailureThreshold: 5,
//	    Timeout:          60 * time.Second,
//	}
//	breaker := manager.GetOrCreate("user-service", config)
func (m *CircuitBreakerManager) GetOrCreate(name string, config CircuitBreakerConfig) *CircuitBreaker {
	// 第一次检查（读锁）
	m.mu.RLock()
	breaker, exists := m.breakers[name]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	// 获取写锁
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查（避免重复创建）
	if breaker, exists := m.breakers[name]; exists {
		return breaker
	}

	// 创建新的熔断器
	breaker = NewCircuitBreaker(name, config)
	m.breakers[name] = breaker
	return breaker
}

// Get 获取指定名称的熔断器。
//
// 功能说明：
// - 根据名称获取熔断器
// - 如果不存在，返回 nil 和 false
//
// 参数：
// - name: 熔断器名称
//
// 返回：
// - *CircuitBreaker: 熔断器实例，如果不存在则为 nil
// - bool: 如果存在返回 true，否则返回 false
//
// 使用示例：
//
//	breaker, exists := manager.Get("user-service")
//	if exists {
//	    state := breaker.GetState()
//	}
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	breaker, exists := m.breakers[name]
	return breaker, exists
}
