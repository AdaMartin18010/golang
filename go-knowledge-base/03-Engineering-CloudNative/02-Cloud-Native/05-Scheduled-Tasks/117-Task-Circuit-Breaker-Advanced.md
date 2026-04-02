# 熔断器高级实现 (Circuit Breaker Advanced Implementation)

> **分类**: 工程与云原生
> **标签**: #circuit-breaker #resilience #failure-handling
> **参考**: Netflix Hystrix, Google SRE Book, Microsoft Polly

---

## 熔断器状态机

```
          成功计数 > threshold
    ┌────────────────────────────┐
    │                            │
    ▼                            │
┌────────┐    失败率 > %     ┌────────┐
│ CLOSED │ ─────────────────► │  OPEN  │
│ (正常)  │                    │ (熔断) │
└────────┘                    └────────┘
    ▲                              │
    │                              │ 超时后
    │    半开状态测试成功           ▼
    └───────────────────────── ┌─────────┐
                                 │  HALF   │
                                 │  OPEN   │
                                 │ (半开)   │
                                 └─────────┘
```

---

## 完整熔断器实现

```go
package circuitbreaker

import (
 "context"
 "errors"
 "sync"
 "sync/atomic"
 "time"
)

// State 熔断器状态
type State int32

const (
 StateClosed    State = iota // 关闭（正常）
 StateOpen                   // 打开（熔断）
 StateHalfOpen               // 半开（测试）
)

func (s State) String() string {
 switch s {
 case StateClosed:
  return "closed"
 case StateOpen:
  return "open"
 case StateHalfOpen:
  return "half-open"
 default:
  return "unknown"
 }
}

// Config 熔断器配置
type Config struct {
 // 失败阈值
 MaxFailures    int           // 最大连续失败次数
 FailureRatio   float64       // 失败率阈值（0.0-1.0）
 MinCalls       int           // 计算失败率的最小调用数

 // 超时配置
 Timeout        time.Duration // 调用超时
 ResetTimeout   time.Duration // 熔断后重置时间
 HalfOpenMaxCalls int        // 半开状态最大测试调用数

 // 成功判定
 SuccessOn     []error       // 视为成功的错误
 FailOn        []error       // 视为失败的错误
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
 name   string
 config Config

 // 状态（原子操作）
 state int32

 // 统计
 stats    *Stats
 statsMu  sync.RWMutex

 // 半开状态计数
 halfOpenCalls int32

 // 最后失败时间
 lastFailureTime int64 // UnixNano

 // 回调
 onStateChange func(name string, from, to State)
}

// Stats 统计数据
type Stats struct {
 TotalCalls    int64
 SuccessCalls  int64
 FailureCalls  int64
 ConsecutiveSuccesses int64
 ConsecutiveFailures  int64
}

// New 创建熔断器
func New(name string, config Config) *CircuitBreaker {
 return &CircuitBreaker{
  name:   name,
  config: config,
  state:  int32(StateClosed),
  stats:  &Stats{},
 }
}

// Execute 执行受保护的函数
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
 state := cb.currentState()

 switch state {
 case StateOpen:
  // 检查是否可以进入半开状态
  if cb.canTransitionToHalfOpen() {
   cb.transitionTo(StateHalfOpen)
  } else {
   return errors.New("circuit breaker is open")
  }

 case StateHalfOpen:
  // 限制半开状态并发
  if atomic.AddInt32(&cb.halfOpenCalls, 1) > int32(cb.config.HalfOpenMaxCalls) {
   atomic.AddInt32(&cb.halfOpenCalls, -1)
   return errors.New("circuit breaker is half-open, too many test calls")
  }
  defer atomic.AddInt32(&cb.halfOpenCalls, -1)
 }

 // 执行函数
 err := cb.doExecute(ctx, fn)

 // 更新统计和状态
 cb.recordResult(err)

 return err
}

// doExecute 执行函数（带超时）
func (cb *CircuitBreaker) doExecute(ctx context.Context, fn func() error) error {
 if cb.config.Timeout <= 0 {
  return fn()
 }

 ctx, cancel := context.WithTimeout(ctx, cb.config.Timeout)
 defer cancel()

 done := make(chan error, 1)
 go func() {
  done <- fn()
 }()

 select {
 case err := <-done:
  return err
 case <-ctx.Done():
  return errors.New("circuit breaker timeout")
 }
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
 cb.statsMu.Lock()
 defer cb.statsMu.Unlock()

 cb.stats.TotalCalls++

 isSuccess := cb.isSuccess(err)

 if isSuccess {
  cb.stats.SuccessCalls++
  cb.stats.ConsecutiveSuccesses++
  cb.stats.ConsecutiveFailures = 0
 } else {
  cb.stats.FailureCalls++
  cb.stats.ConsecutiveFailures++
  cb.stats.ConsecutiveSuccesses = 0
  atomic.StoreInt64(&cb.lastFailureTime, time.Now().UnixNano())
 }

 // 状态转换
 cb.updateState(isSuccess)
}

// updateState 根据结果更新状态
func (cb *CircuitBreaker) updateState(isSuccess bool) {
 state := cb.currentState()

 switch state {
 case StateClosed:
  // 检查是否需要熔断
  if cb.shouldOpen() {
   cb.transitionTo(StateOpen)
  }

 case StateHalfOpen:
  if isSuccess {
   // 测试成功，关闭熔断器
   cb.transitionTo(StateClosed)
  } else {
   // 测试失败，重新熔断
   cb.transitionTo(StateOpen)
  }

 case StateOpen:
  // 不应该发生，因为Open状态会直接返回错误
 }
}

// shouldOpen 检查是否应该熔断
func (cb *CircuitBreaker) shouldOpen() bool {
 stats := cb.stats

 // 连续失败次数超过阈值
 if stats.ConsecutiveFailures >= int64(cb.config.MaxFailures) {
  return true
 }

 // 失败率超过阈值
 if stats.TotalCalls >= int64(cb.config.MinCalls) {
  failureRatio := float64(stats.FailureCalls) / float64(stats.TotalCalls)
  if failureRatio >= cb.config.FailureRatio {
   return true
  }
 }

 return false
}

// canTransitionToHalfOpen 检查是否可以进入半开状态
func (cb *CircuitBreaker) canTransitionToHalfOpen() bool {
 lastFailure := atomic.LoadInt64(&cb.lastFailureTime)
 if lastFailure == 0 {
  return false
 }

 elapsed := time.Since(time.Unix(0, lastFailure))
 return elapsed >= cb.config.ResetTimeout
}

// transitionTo 状态转换
func (cb *CircuitBreaker) transitionTo(newState State) {
 oldState := cb.currentState()
 if oldState == newState {
  return
 }

 atomic.StoreInt32(&cb.state, int32(newState))

 // 重置统计（关闭时）
 if newState == StateClosed {
  cb.resetStats()
 }

 // 回调
 if cb.onStateChange != nil {
  cb.onStateChange(cb.name, oldState, newState)
 }
}

// resetStats 重置统计
func (cb *CircuitBreaker) resetStats() {
 cb.statsMu.Lock()
 defer cb.statsMu.Unlock()

 cb.stats = &Stats{}
}

// currentState 获取当前状态
func (cb *CircuitBreaker) currentState() State {
 return State(atomic.LoadInt32(&cb.state))
}

// isSuccess 判断是否为成功
func (cb *CircuitBreaker) isSuccess(err error) bool {
 if err == nil {
  return true
 }

 // 检查是否为指定成功的错误
 for _, successErr := range cb.config.SuccessOn {
  if errors.Is(err, successErr) {
   return true
  }
 }

 // 检查是否为指定失败的错误
 for _, failErr := range cb.config.FailOn {
  if errors.Is(err, failErr) {
   return false
  }
 }

 return false
}

// GetState 获取当前状态（用于监控）
func (cb *CircuitBreaker) GetState() State {
 return cb.currentState()
}

// GetStats 获取统计（用于监控）
func (cb *CircuitBreaker) GetStats() Stats {
 cb.statsMu.RLock()
 defer cb.statsMu.RUnlock()
 return *cb.stats
}
```

---

## 自适应熔断器

```go
// AdaptiveBreaker 自适应熔断器
type AdaptiveBreaker struct {
 *CircuitBreaker

 // 自适应配置
 adaptationInterval time.Duration
 latencyStats       *LatencyStats
}

type LatencyStats struct {
 p50   float64
 p95   float64
 p99   float64
 count int64
}

// 根据延迟自适应调整
func (ab *AdaptiveBreaker) adaptThreshold() {
 stats := ab.latencyStats

 // 如果 P99 延迟过高，降低失败阈值
 if stats.p99 > 1000 { // 1s
  ab.config.MaxFailures = max(1, ab.config.MaxFailures-1)
 }

 // 如果延迟正常，提高失败阈值
 if stats.p95 < 100 { // 100ms
  ab.config.MaxFailures = min(10, ab.config.MaxFailures+1)
 }
}
```
