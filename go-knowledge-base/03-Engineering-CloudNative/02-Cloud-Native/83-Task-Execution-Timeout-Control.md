# 任务执行超时控制 (Task Execution Timeout Control)

> **分类**: 工程与云原生
> **标签**: #timeout #context #deadline #cancellation
> **参考**: Go Context, Circuit Breaker, Distributed Timeout

---

## 超时控制架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task Execution Timeout Control                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Timeout Hierarchy                                 │   │
│  │                                                                      │   │
│  │   Global Timeout (Workflow) ──► 30 minutes                          │   │
│  │          │                                                          │   │
│  │          ▼                                                          │   │
│  │   Task Timeout ──► 5 minutes                                        │   │
│  │          │                                                          │   │
│  │          ▼                                                          │   │
│  │   Operation Timeout ──► 30 seconds                                  │   │
│  │          │                                                          │   │
│  │          ▼                                                          │   │
│  │   Network Timeout ──► 10 seconds                                    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Timeout Strategies                                │   │
│  │                                                                      │   │
│  │   1. Hard Timeout: Immediate cancellation on deadline               │   │
│  │   2. Soft Timeout: Graceful shutdown period after deadline          │   │
│  │   3. Incremental Timeout: Progressive escalation                    │   │
│  │   4. Adaptive Timeout: Dynamic based on historical data             │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整超时控制实现

```go
package timeout

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// TimeoutPolicy 超时策略
type TimeoutPolicy struct {
    // 基本超时
    Timeout       time.Duration

    // 优雅关闭
    GracePeriod   time.Duration // 超时后的宽限期

    // 渐进式超时
    WarningAt     time.Duration // 提前警告时间

    // 重试
    RetryAttempts int
    RetryDelay    time.Duration

    // 行为
    HardTimeout   bool // true=立即取消, false=优雅关闭
}

// DefaultTimeoutPolicy 默认超时策略
var DefaultTimeoutPolicy = TimeoutPolicy{
    Timeout:       5 * time.Minute,
    GracePeriod:   30 * time.Second,
    WarningAt:     0,
    RetryAttempts: 0,
    RetryDelay:    0,
    HardTimeout:   false,
}

// TimeoutController 超时控制器
type TimeoutController struct {
    policy   TimeoutPolicy

    // 状态
    state    int32 // 0=idle, 1=running, 2=warning, 3=timeout, 4=cancelled

    // 控制
    ctx      context.Context
    cancel   context.CancelFunc

    // 回调
    onWarning func()
    onTimeout func()
    onCancel  func()

    mu       sync.RWMutex
}

const (
    stateIdle = iota
    stateRunning
    stateWarning
    stateTimeout
    stateCancelled
)

// NewTimeoutController 创建超时控制器
func NewTimeoutController(policy TimeoutPolicy) *TimeoutController {
    return &TimeoutController{
        policy: policy,
    }
}

// Execute 执行带超时的操作
func (tc *TimeoutController) Execute(ctx context.Context, operation func(context.Context) error) error {
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(ctx, tc.policy.Timeout)
    defer cancel()

    tc.ctx = ctx
    tc.cancel = cancel
    atomic.StoreInt32(&tc.state, stateRunning)

    // 启动警告计时器（如果配置了）
    if tc.policy.WarningAt > 0 && tc.policy.WarningAt < tc.policy.Timeout {
        go tc.warningTimer(tc.policy.WarningAt)
    }

    // 执行操作
    done := make(chan error, 1)
    go func() {
        done <- operation(ctx)
    }()

    select {
    case err := <-done:
        atomic.StoreInt32(&tc.state, stateIdle)
        return err

    case <-ctx.Done():
        // 超时处理
        return tc.handleTimeout(operation)
    }
}

// handleTimeout 处理超时
func (tc *TimeoutController) handleTimeout(operation func(context.Context) error) error {
    if tc.policy.HardTimeout {
        atomic.StoreInt32(&tc.state, stateTimeout)
        if tc.onTimeout != nil {
            tc.onTimeout()
        }
        return fmt.Errorf("operation timed out after %v", tc.policy.Timeout)
    }

    // 优雅关闭
    atomic.StoreInt32(&tc.state, stateTimeout)
    if tc.onTimeout != nil {
        tc.onTimeout()
    }

    // 给予宽限期
    if tc.policy.GracePeriod > 0 {
        graceCtx, graceCancel := context.WithTimeout(context.Background(), tc.policy.GracePeriod)
        defer graceCancel()

        done := make(chan struct{})
        go func() {
            // 等待操作完成或宽限期结束
            close(done)
        }()

        select {
        case <-done:
            return nil
        case <-graceCtx.Done():
            return fmt.Errorf("operation did not complete within grace period")
        }
    }

    return fmt.Errorf("operation timed out after %v", tc.policy.Timeout)
}

// warningTimer 警告计时器
func (tc *TimeoutController) warningTimer(after time.Duration) {
    timer := time.NewTimer(after)
    defer timer.Stop()

    select {
    case <-timer.C:
        if atomic.LoadInt32(&tc.state) == stateRunning {
            atomic.StoreInt32(&tc.state, stateWarning)
            if tc.onWarning != nil {
                tc.onWarning()
            }
        }
    case <-tc.ctx.Done():
    }
}

// Cancel 手动取消
func (tc *TimeoutController) Cancel() {
    if tc.cancel != nil {
        tc.cancel()
        atomic.StoreInt32(&tc.state, stateCancelled)
        if tc.onCancel != nil {
            tc.onCancel()
        }
    }
}

// OnWarning 设置警告回调
func (tc *TimeoutController) OnWarning(fn func()) {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    tc.onWarning = fn
}

// OnTimeout 设置超时回调
func (tc *TimeoutController) OnTimeout(fn func()) {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    tc.onTimeout = fn
}

// State 获取当前状态
func (tc *TimeoutController) State() string {
    switch atomic.LoadInt32(&tc.state) {
    case stateIdle:
        return "idle"
    case stateRunning:
        return "running"
    case stateWarning:
        return "warning"
    case stateTimeout:
        return "timeout"
    case stateCancelled:
        return "cancelled"
    default:
        return "unknown"
    }
}

// TimeoutError 超时错误
type TimeoutError struct {
    Operation string
    Timeout   time.Duration
}

func (e *TimeoutError) Error() string {
    return fmt.Sprintf("operation %s timed out after %v", e.Operation, e.Timeout)
}

func (e *TimeoutError) Timeout() bool { return true }
```

---

## 分层超时控制

```go
package timeout

import (
    "context"
    "fmt"
    "time"
)

// HierarchicalTimeout 分层超时控制
type HierarchicalTimeout struct {
    // 全局超时
    GlobalTimeout time.Duration

    // 各层超时
    Layers []LayerTimeout
}

// LayerTimeout 层超时配置
type LayerTimeout struct {
    Name    string
    Timeout time.Duration
    Percent float64 // 占总超时的百分比
}

// HierarchicalController 分层控制器
type HierarchicalController struct {
    globalTimeout time.Duration
    layers        []layerControl
}

type layerControl struct {
    name    string
    timeout time.Duration
    used    time.Duration
}

// NewHierarchicalController 创建分层控制器
func NewHierarchicalController(globalTimeout time.Duration, layers []LayerTimeout) *HierarchicalController {
    hc := &HierarchicalController{
        globalTimeout: globalTimeout,
    }

    for _, layer := range layers {
        var timeout time.Duration
        if layer.Timeout > 0 {
            timeout = layer.Timeout
        } else if layer.Percent > 0 {
            timeout = time.Duration(float64(globalTimeout) * layer.Percent / 100)
        }

        hc.layers = append(hc.layers, layerControl{
            name:    layer.Name,
            timeout: timeout,
        })
    }

    return hc
}

// Execute 分层执行
func (hc *HierarchicalController) Execute(ctx context.Context, operations map[string]func(context.Context) error) error {
    // 创建全局超时上下文
    ctx, cancel := context.WithTimeout(ctx, hc.globalTimeout)
    defer cancel()

    // 按顺序执行各层
    for i, layer := range hc.layers {
        op, ok := operations[layer.name]
        if !ok {
            continue
        }

        // 创建层上下文
        layerCtx, layerCancel := context.WithTimeout(ctx, layer.timeout)

        start := time.Now()
        err := op(layerCtx)
        layer.used = time.Since(start)

        layerCancel()

        if err != nil {
            return fmt.Errorf("layer %s failed: %w", layer.name, err)
        }

        // 更新剩余层的时间
        hc.adjustTimeouts(i)
    }

    return nil
}

// adjustTimeouts 调整超时
func (hc *HierarchicalController) adjustTimeouts(completedIndex int) {
    // 根据已用时间调整后续层的超时
}

// GetUsage 获取各层使用情况
func (hc *HierarchicalController) GetUsage() map[string]time.Duration {
    usage := make(map[string]time.Duration)
    for _, layer := range hc.layers {
        usage[layer.name] = layer.used
    }
    return usage
}

// AdaptiveTimeout 自适应超时
type AdaptiveTimeout struct {
    baseTimeout time.Duration

    // 历史数据
    history     []time.Duration
    maxHistory  int

    // 调整因子
    factor      float64
}

// NewAdaptiveTimeout 创建自适应超时
func NewAdaptiveTimeout(baseTimeout time.Duration, maxHistory int) *AdaptiveTimeout {
    return &AdaptiveTimeout{
        baseTimeout: baseTimeout,
        maxHistory:  maxHistory,
        factor:      1.0,
    }
}

// GetTimeout 获取当前超时
func (at *AdaptiveTimeout) GetTimeout() time.Duration {
    if len(at.history) == 0 {
        return at.baseTimeout
    }

    // 计算 P95
    p95 := at.calculateP95()

    // 自适应调整
    timeout := time.Duration(float64(p95) * 1.5) // 1.5x P95
    if timeout < at.baseTimeout {
        timeout = at.baseTimeout
    }

    return timeout
}

// RecordExecution 记录执行时间
func (at *AdaptiveTimeout) RecordExecution(duration time.Duration) {
    at.history = append(at.history, duration)
    if len(at.history) > at.maxHistory {
        at.history = at.history[1:]
    }
}

func (at *AdaptiveTimeout) calculateP95() time.Duration {
    if len(at.history) == 0 {
        return at.baseTimeout
    }

    // 排序并取 P95
    sorted := make([]time.Duration, len(at.history))
    copy(sorted, at.history)

    // 简化为取最大值
    max := sorted[0]
    for _, d := range sorted {
        if d > max {
            max = d
        }
    }

    return max
}
```

---

## 分布式超时传播

```go
package timeout

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

// TimeoutHeader 超时头
type TimeoutHeader struct {
    Deadline  time.Time `json:"deadline"`
    Remaining int64     `json:"remaining_ms"` // 剩余毫秒
    TraceID   string    `json:"trace_id"`
}

// PropagateTimeout 传播超时
func PropagateTimeout(ctx context.Context, header http.Header) {
    deadline, ok := ctx.Deadline()
    if !ok {
        return
    }

    remaining := time.Until(deadline)
    if remaining <= 0 {
        return
    }

    timeoutHeader := TimeoutHeader{
        Deadline:  deadline,
        Remaining: remaining.Milliseconds(),
    }

    data, _ := json.Marshal(timeoutHeader)
    header.Set("X-Timeout-Context", string(data))
}

// ExtractTimeout 提取超时
func ExtractTimeout(header http.Header) (context.Context, context.CancelFunc) {
    ctx := context.Background()

    data := header.Get("X-Timeout-Context")
    if data == "" {
        return ctx, func() {}
    }

    var timeoutHeader TimeoutHeader
    if err := json.Unmarshal([]byte(data), &timeoutHeader); err != nil {
        return ctx, func() {}
    }

    remaining := time.Until(timeoutHeader.Deadline)
    if remaining <= 0 {
        return ctx, func() {}
    }

    return context.WithTimeout(ctx, remaining)
}

// TimeoutMiddleware HTTP 超时中间件
func TimeoutMiddleware(defaultTimeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 从请求中提取超时
            ctx, cancel := ExtractTimeout(r.Header)
            if ctx == r.Context() {
                // 没有传播的超时，使用默认
                ctx, cancel = context.WithTimeout(r.Context(), defaultTimeout)
            }
            defer cancel()

            // 传播超时到下游
            PropagateTimeout(ctx, w.Header())

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "timeout"
)

func main() {
    // 基本超时控制
    policy := timeout.TimeoutPolicy{
        Timeout:       5 * time.Second,
        GracePeriod:   1 * time.Second,
        WarningAt:     3 * time.Second,
        HardTimeout:   false,
    }

    controller := timeout.NewTimeoutController(policy)

    controller.OnWarning(func() {
        fmt.Println("Warning: approaching timeout")
    })

    controller.OnTimeout(func() {
        fmt.Println("Timeout occurred")
    })

    err := controller.Execute(context.Background(), func(ctx context.Context) error {
        // 模拟长时间操作
        select {
        case <-time.After(10 * time.Second):
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    })

    if err != nil {
        fmt.Printf("Execution failed: %v\n", err)
    }

    // 分层超时
    hc := timeout.NewHierarchicalController(30*time.Second, []timeout.LayerTimeout{
        {Name: "validation", Timeout: 5 * time.Second},
        {Name: "processing", Timeout: 20 * time.Second},
        {Name: "storage", Timeout: 5 * time.Second},
    })

    operations := map[string]func(context.Context) error{
        "validation": func(ctx context.Context) error {
            time.Sleep(1 * time.Second)
            return nil
        },
        "processing": func(ctx context.Context) error {
            time.Sleep(2 * time.Second)
            return nil
        },
        "storage": func(ctx context.Context) error {
            time.Sleep(500 * time.Millisecond)
            return nil
        },
    }

    if err := hc.Execute(context.Background(), operations); err != nil {
        panic(err)
    }

    usage := hc.GetUsage()
    for layer, used := range usage {
        fmt.Printf("Layer %s used: %v\n", layer, used)
    }
}
```
