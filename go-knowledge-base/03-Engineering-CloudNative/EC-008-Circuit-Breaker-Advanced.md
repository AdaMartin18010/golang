# EC-008: 熔断器高级实现 (Circuit Breaker Advanced)

> **维度**: Engineering CloudNative
> **级别**: S (15+ KB)
> **标签**: #circuit-breaker #resilience #failure-handling #adaptive
> **相关**: EC-007, EC-042, FT-015

---

## 整合说明

本文档合并了：

- `08-Circuit-Breaker-Patterns.md` (5.1 KB) - 基础模式
- `117-Task-Circuit-Breaker-Advanced.md` (8.3 KB) - 高级实现

---

## 状态机

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

## 完整实现

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
 MaxFailures      int           // 最大连续失败次数
 FailureRatio     float64       // 失败率阈值（0.0-1.0）
 MinCalls         int           // 计算失败率的最小调用数
 Timeout          time.Duration // 调用超时
 ResetTimeout     time.Duration // 熔断后重置时间
 HalfOpenMaxCalls int           // 半开状态最大测试调用数
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
 name   string
 config Config

 state           int32
 stats           *Stats
 statsMu         sync.RWMutex
 halfOpenCalls   int32
 lastFailureTime int64
 onStateChange   func(name string, from, to State)
}

// Stats 统计数据
type Stats struct {
 TotalCalls           int64
 SuccessCalls         int64
 FailureCalls         int64
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
  if cb.canTransitionToHalfOpen() {
   cb.transitionTo(StateHalfOpen)
  } else {
   return errors.New("circuit breaker is open")
  }

 case StateHalfOpen:
  if atomic.AddInt32(&cb.halfOpenCalls, 1) > int32(cb.config.HalfOpenMaxCalls) {
   atomic.AddInt32(&cb.halfOpenCalls, -1)
   return errors.New("circuit breaker is half-open, too many test calls")
  }
  defer atomic.AddInt32(&cb.halfOpenCalls, -1)
 }

 err := cb.doExecute(ctx, fn)
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

 isSuccess := err == nil

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

 cb.updateState(isSuccess)
}

// updateState 根据结果更新状态
func (cb *CircuitBreaker) updateState(isSuccess bool) {
 state := cb.currentState()

 switch state {
 case StateClosed:
  if cb.shouldOpen() {
   cb.transitionTo(StateOpen)
  }

 case StateHalfOpen:
  if isSuccess {
   cb.transitionTo(StateClosed)
  } else {
   cb.transitionTo(StateOpen)
  }
 }
}

// shouldOpen 检查是否应该熔断
func (cb *CircuitBreaker) shouldOpen() bool {
 stats := cb.stats

 if stats.ConsecutiveFailures >= int64(cb.config.MaxFailures) {
  return true
 }

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

 if newState == StateClosed {
  cb.resetStats()
 }

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
 if stats.p99 > 1000 {
  ab.config.MaxFailures = max(1, ab.config.MaxFailures-1)
 }

 // 如果延迟正常，提高失败阈值
 if stats.p95 < 100 {
  ab.config.MaxFailures = min(10, ab.config.MaxFailures+1)
 }
}
```

---

## 使用示例

```go
// 创建熔断器
cb := circuitbreaker.New("api-client", circuitbreaker.Config{
 MaxFailures:      5,
 FailureRatio:     0.5,
 MinCalls:         10,
 Timeout:          5 * time.Second,
 ResetTimeout:     30 * time.Second,
 HalfOpenMaxCalls: 3,
})

// 使用
err := cb.Execute(ctx, func() error {
 return callExternalAPI()
})
```
