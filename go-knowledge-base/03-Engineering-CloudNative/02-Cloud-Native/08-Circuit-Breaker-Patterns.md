# 熔断器模式详解 (Circuit Breaker Patterns)

> **分类**: 工程与云原生
> **标签**: #circuit-breaker #resilience #pattern

---

## 熔断器状态机

```
        ┌─────────────┐
   ┌─── │   CLOSED    │ ◄── 成功计数
   │    │  (正常请求)  │
   │    └──────┬──────┘
   │           │ 失败阈值
   │           ▼
   │    ┌─────────────┐
   │    │    OPEN     │
   │    │  (拒绝请求)  │
   │    └──────┬──────┘
   │           │ 超时
   │           ▼
   │    ┌─────────────┐
   └─── │  HALF-OPEN  │
        │ (测试请求)  │
        └─────────────┘
```

---

## 完整实现

```go
type State int

const (
    StateClosed State = iota    // 正常
    StateOpen                    // 熔断
    StateHalfOpen                // 半开
)

type CircuitBreaker struct {
    name          string
    state         State
    failureCount  int
    successCount  int
    lastFailureTime time.Time

    // 配置
    maxFailures    int           // 触发熔断的失败次数
    timeout        time.Duration // 熔断持续时间
    halfOpenMaxCalls int         // 半开状态最大测试请求

    mu sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if !cb.canExecute() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return true

    case StateOpen:
        if time.Since(cb.lastFailureTime) > cb.timeout {
            cb.state = StateHalfOpen
            cb.successCount = 0
            return true
        }
        return false

    case StateHalfOpen:
        return cb.successCount < cb.halfOpenMaxCalls
    }

    return false
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()

        if cb.state == StateHalfOpen || cb.failureCount >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        cb.successCount++

        if cb.state == StateHalfOpen {
            if cb.successCount >= cb.halfOpenMaxCalls {
                cb.state = StateClosed
                cb.failureCount = 0
            }
        } else {
            cb.failureCount = 0
        }
    }
}
```

---

## 高级特性

### 慢调用熔断

```go
func (cb *CircuitBreaker) CallWithTimeout(fn func() error, maxDuration time.Duration) error {
    if !cb.canExecute() {
        return ErrCircuitOpen
    }

    done := make(chan error, 1)
    start := time.Now()

    go func() {
        done <- fn()
    }()

    select {
    case err := <-done:
        duration := time.Since(start)

        // 记录慢调用
        if duration > maxDuration {
            cb.recordSlowCall()
        }

        cb.recordResult(err)
        return err

    case <-time.After(maxDuration):
        cb.recordSlowCall()
        return ErrTimeout
    }
}
```

### 自适应熔断

```go
type AdaptiveCircuitBreaker struct {
    *CircuitBreaker

    // 滑动窗口
    window []bool  // true = 成功, false = 失败
    windowSize int
    windowPos  int
}

func (acb *AdaptiveCircuitBreaker) recordResult(err error) {
    acb.mu.Lock()
    defer acb.mu.Unlock()

    // 记录到窗口
    acb.window[acb.windowPos] = (err == nil)
    acb.windowPos = (acb.windowPos + 1) % acb.windowSize

    // 计算错误率
    failureRate := acb.calculateFailureRate()

    // 自适应调整阈值
    if failureRate > 0.5 {
        acb.maxFailures = max(1, acb.maxFailures-1)  // 更敏感
    } else if failureRate < 0.1 {
        acb.maxFailures = min(10, acb.maxFailures+1)  // 更宽容
    }

    acb.recordResult(err)
}
```

---

## 与重试结合

```go
func CallWithRetryAndCircuitBreaker(
    cb *CircuitBreaker,
    fn func() error,
    maxRetries int,
    backoff time.Duration,
) error {
    return cb.Call(func() error {
        var err error

        for i := 0; i < maxRetries; i++ {
            err = fn()
            if err == nil {
                return nil
            }

            time.Sleep(backoff * time.Duration(1<<i))
        }

        return err
    })
})
```

---

## 监控

```go
type CircuitBreakerMetrics struct {
    State           prometheus.Gauge
    FailureCount    prometheus.Counter
    SuccessCount    prometheus.Counter
    RejectedCount   prometheus.Counter
}

func (cb *CircuitBreaker) Metrics() CircuitBreakerMetrics {
    return CircuitBreakerMetrics{
        State: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "circuit_breaker_state",
            Labels: map[string]string{"name": cb.name},
        }),
    }
}
```
