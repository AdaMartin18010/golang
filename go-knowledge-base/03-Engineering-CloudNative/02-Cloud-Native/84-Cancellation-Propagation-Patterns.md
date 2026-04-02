# 取消传播模式 (Cancellation Propagation Patterns)

> **分类**: 工程与云原生
> **标签**: #cancellation #context #propagation #graceful-shutdown
> **参考**: Go Context, Distributed Cancellation

---

## 取消传播架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Cancellation Propagation Architecture                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Root Context (User Request)                                                │
│        │                                                                     │
│        ▼ cancel()                                                            │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │              propagated via context.WithCancel()                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│   ┌────┴────┬──────────┬──────────┬──────────┐                              │
│   ▼         ▼          ▼          ▼          ▼                              │
│   HTTP    gRPC      Message    Database    External                         │
│  Handler  Service     Queue      Query       API                            │
│   │         │          │          │          │                              │
│   ▼         ▼          ▼          ▼          ▼                              │
│  Check   Check      Check      Check      Check                             │
│  ctx.Done() ctx.Done() ctx.Done() ctx.Done() ctx.Done()                     │
│   │         │          │          │          │                              │
│   ▼         ▼          ▼          ▼          ▼                              │
│  Stop    Stop       Stop       Stop       Stop                              │
│  Process Process   Process   Process   Process                              │
│                                                                              │
│   Cancellation Strategies:                                                   │
│   1. Immediate: Stop processing immediately                                  │
│   2. Graceful: Complete current item, then stop                              │
│   3. Timeout: Stop after grace period                                        │
│   4. Forceful: Kill process after timeout                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整取消传播实现

```go
package cancellation

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// CancellationStrategy 取消策略
type CancellationStrategy int

const (
    StrategyImmediate CancellationStrategy = iota
    StrategyGraceful
    StrategyTimeout
    StrategyForceful
)

// Propagator 取消传播器
type Propagator struct {
    strategy CancellationStrategy
    gracePeriod time.Duration

    // 子组件
    children map[string]*ChildComponent
    mu       sync.RWMutex
}

// ChildComponent 子组件
type ChildComponent struct {
    Name     string
    Cancel   context.CancelFunc
    Done     <-chan struct{}
    Strategy CancellationStrategy
}

// NewPropagator 创建传播器
func NewPropagator(strategy CancellationStrategy, gracePeriod time.Duration) *Propagator {
    return &Propagator{
        strategy:    strategy,
        gracePeriod: gracePeriod,
        children:    make(map[string]*ChildComponent),
    }
}

// RegisterChild 注册子组件
func (p *Propagator) RegisterChild(name string, cancel context.CancelFunc, done <-chan struct{}) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.children[name] = &ChildComponent{
        Name:     name,
        Cancel:   cancel,
        Done:     done,
        Strategy: p.strategy,
    }
}

// Propagate 传播取消信号
func (p *Propagator) Propagate(ctx context.Context, reason string) error {
    fmt.Printf("Propagating cancellation: %s\n", reason)

    p.mu.RLock()
    children := make(map[string]*ChildComponent, len(p.children))
    for k, v := range p.children {
        children[k] = v
    }
    p.mu.RUnlock()

    var wg sync.WaitGroup
    errChan := make(chan error, len(children))

    for name, child := range children {
        wg.Add(1)
        go func(n string, c *ChildComponent) {
            defer wg.Done()

            if err := p.cancelChild(ctx, n, c); err != nil {
                errChan <- fmt.Errorf("failed to cancel %s: %w", n, err)
            }
        }(name, child)
    }

    // 等待所有子组件完成
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        close(errChan)

        // 收集错误
        var errs []error
        for err := range errChan {
            errs = append(errs, err)
        }

        if len(errs) > 0 {
            return fmt.Errorf("cancellation errors: %v", errs)
        }
        return nil

    case <-ctx.Done():
        return fmt.Errorf("cancellation propagation timed out")
    }
}

func (p *Propagator) cancelChild(ctx context.Context, name string, child *ChildComponent) error {
    switch child.Strategy {
    case StrategyImmediate:
        child.Cancel()
        <-child.Done
        return nil

    case StrategyGraceful:
        child.Cancel()
        // 等待完成或超时
        select {
        case <-child.Done:
            return nil
        case <-time.After(p.gracePeriod):
            return fmt.Errorf("graceful cancellation timeout")
        }

    case StrategyTimeout:
        child.Cancel()
        // 使用传入的上下文超时
        select {
        case <-child.Done:
            return nil
        case <-ctx.Done():
            return fmt.Errorf("cancellation timeout")
        }

    case StrategyForceful:
        child.Cancel()
        select {
        case <-child.Done:
            return nil
        case <-time.After(p.gracePeriod):
            // 强制终止（在实际实现中可能需要更强的措施）
            return fmt.Errorf("forceful cancellation required")
        }
    }

    return nil
}

// RemoveChild 移除子组件
func (p *Propagator) RemoveChild(name string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    delete(p.children, name)
}

// ContextPropagator 上下文传播器
type ContextPropagator struct {
    propagator *Propagator
}

// NewContextPropagator 创建上下文传播器
func NewContextPropagator(propagator *Propagator) *ContextPropagator {
    return &ContextPropagator{propagator: propagator}
}

// WrapHandler 包装处理器以支持取消传播
func (cp *ContextPropagator) WrapHandler(name string, handler func(context.Context) error) func(context.Context) error {
    return func(ctx context.Context) error {
        // 创建子上下文
        childCtx, cancel := context.WithCancel(ctx)
        defer cancel()

        // 注册到传播器
        cp.propagator.RegisterChild(name, cancel, childCtx.Done())
        defer cp.propagator.RemoveChild(name)

        // 监听取消信号
        done := make(chan error, 1)
        go func() {
            done <- handler(childCtx)
        }()

        select {
        case err := <-done:
            return err
        case <-ctx.Done():
            return fmt.Errorf("handler cancelled: %w", ctx.Err())
        }
    }
}
```

---

## 取消监听模式

```go
package cancellation

import (
    "context"
    "sync"
    "time"
)

// DoneNotifier 完成通知器
type DoneNotifier struct {
    done chan struct{}
    once sync.Once
}

// NewDoneNotifier 创建完成通知器
func NewDoneNotifier() *DoneNotifier {
    return &DoneNotifier{
        done: make(chan struct{}),
    }
}

// Done 返回 done channel
func (dn *DoneNotifier) Done() <-chan struct{} {
    return dn.done
}

// Notify 通知完成
func (dn *DoneNotifier) Notify() {
    dn.once.Do(func() {
        close(dn.done)
    })
}

// CancellationListener 取消监听器
type CancellationListener struct {
    callbacks []func()
    mu        sync.RWMutex
}

// NewCancellationListener 创建监听器
func NewCancellationListener() *CancellationListener {
    return &CancellationListener{
        callbacks: make([]func(), 0),
    }
}

// AddCallback 添加取消回调
func (cl *CancellationListener) AddCallback(cb func()) {
    cl.mu.Lock()
    defer cl.mu.Unlock()
    cl.callbacks = append(cl.callbacks, cb)
}

// Listen 监听上下文取消
func (cl *CancellationListener) Listen(ctx context.Context) {
    go func() {
        <-ctx.Done()
        cl.triggerCallbacks()
    }()
}

func (cl *CancellationListener) triggerCallbacks() {
    cl.mu.RLock()
    callbacks := make([]func(), len(cl.callbacks))
    copy(callbacks, cl.callbacks)
    cl.mu.RUnlock()

    for _, cb := range callbacks {
        cb()
    }
}

// GracefulOperation 优雅操作
type GracefulOperation struct {
    stopChan chan struct{}
    doneChan chan struct{}
}

// NewGracefulOperation 创建优雅操作
func NewGracefulOperation() *GracefulOperation {
    return &GracefulOperation{
        stopChan: make(chan struct{}),
        doneChan: make(chan struct{}),
    }
}

// Run 运行操作
func (go *GracefulOperation) Run(ctx context.Context, operation func(stop <-chan struct{}) error) error {
    defer close(go.doneChan)

    // 监听上下文取消
    go func() {
        <-ctx.Done()
        close(go.stopChan)
    }()

    return operation(go.stopChan)
}

// Wait 等待完成
func (go *GracefulOperation) Wait(timeout time.Duration) error {
    select {
    case <-go.doneChan:
        return nil
    case <-time.After(timeout):
        return fmt.Errorf("wait timeout")
    }
}

// CascadeCancellation 级联取消
type CascadeCancellation struct {
    cancelFuncs []context.CancelFunc
    mu          sync.Mutex
}

// NewCascadeCancellation 创建级联取消
func NewCascadeCancellation() *CascadeCancellation {
    return &CascadeCancellation{
        cancelFuncs: make([]context.CancelFunc, 0),
    }
}

// Add 添加取消函数
func (cc *CascadeCancellation) Add(cancel context.CancelFunc) {
    cc.mu.Lock()
    defer cc.mu.Unlock()
    cc.cancelFuncs = append(cc.cancelFuncs, cancel)
}

// CancelAll 取消所有
func (cc *CascadeCancellation) CancelAll() {
    cc.mu.Lock()
    cancels := make([]context.CancelFunc, len(cc.cancelFuncs))
    copy(cancels, cc.cancelFuncs)
    cc.mu.Unlock()

    for _, cancel := range cancels {
        cancel()
    }
}
```

---

## 跨服务取消传播

```go
package cancellation

import (
    "context"
    "encoding/json"
    "net/http"
)

// CancellationHeader 取消头
type CancellationHeader struct {
    TraceID   string `json:"trace_id"`
    Reason    string `json:"reason"`
    Timestamp int64  `json:"timestamp"`
}

// HTTPPropagator HTTP 取消传播器
type HTTPPropagator struct {
    headerName string
}

// NewHTTPPropagator 创建 HTTP 传播器
func NewHTTPPropagator(headerName string) *HTTPPropagator {
    if headerName == "" {
        headerName = "X-Cancellation-Context"
    }
    return &HTTPPropagator{headerName: headerName}
}

// Inject 注入取消上下文
func (hp *HTTPPropagator) Inject(ctx context.Context, header http.Header) {
    deadline, ok := ctx.Deadline()
    if !ok {
        return
    }

    cancelHeader := CancellationHeader{
        TraceID:   getTraceID(ctx),
        Reason:    "timeout",
        Timestamp: deadline.Unix(),
    }

    data, _ := json.Marshal(cancelHeader)
    header.Set(hp.headerName, string(data))
}

// Extract 提取取消上下文
func (hp *HTTPPropagator) Extract(header http.Header) (context.Context, context.CancelFunc) {
    ctx := context.Background()

    data := header.Get(hp.headerName)
    if data == "" {
        return ctx, func() {}
    }

    var cancelHeader CancellationHeader
    if err := json.Unmarshal([]byte(data), &cancelHeader); err != nil {
        return ctx, func() {}
    }

    // 根据时间戳创建上下文
    // 实际实现中可能需要更复杂的逻辑
    return context.WithCancel(ctx)
}

// MessagePropagator 消息队列取消传播器
type MessagePropagator struct {
    headerKey string
}

// NewMessagePropagator 创建消息传播器
func NewMessagePropagator(headerKey string) *MessagePropagator {
    if headerKey == "" {
        headerKey = "cancellation_context"
    }
    return &MessagePropagator{headerKey: headerKey}
}

// InjectIntoMessage 注入到消息
func (mp *MessagePropagator) InjectIntoMessage(ctx context.Context, headers map[string]string) {
    deadline, ok := ctx.Deadline()
    if !ok {
        return
    }

    cancelHeader := CancellationHeader{
        TraceID:   getTraceID(ctx),
        Reason:    "timeout",
        Timestamp: deadline.Unix(),
    }

    data, _ := json.Marshal(cancelHeader)
    headers[mp.headerKey] = string(data)
}

// ExtractFromMessage 从消息提取
func (mp *MessagePropagator) ExtractFromMessage(headers map[string]string) (context.Context, context.CancelFunc) {
    ctx := context.Background()

    data, ok := headers[mp.headerKey]
    if !ok {
        return ctx, func() {}
    }

    var cancelHeader CancellationHeader
    if err := json.Unmarshal([]byte(data), &cancelHeader); err != nil {
        return ctx, func() {}
    }

    return context.WithCancel(ctx)
}

func getTraceID(ctx context.Context) string {
    // 从上下文中获取 trace ID
    return ""
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

    "cancellation"
)

func main() {
    // 创建传播器
    propagator := cancellation.NewPropagator(
        cancellation.StrategyGraceful,
        5*time.Second,
    )

    // 创建根上下文
    ctx, rootCancel := context.WithCancel(context.Background())

    // 启动多个 goroutine
    for i := 0; i < 3; i++ {
        childCtx, childCancel := context.WithCancel(ctx)
        propagator.RegisterChild(fmt.Sprintf("worker-%d", i), childCancel, childCtx.Done())

        go worker(childCtx, i)
    }

    // 模拟运行
    time.Sleep(2 * time.Second)

    // 触发取消
    rootCancel()

    // 传播取消
    if err := propagator.Propagate(context.Background(), "manual shutdown"); err != nil {
        fmt.Printf("Propagation error: %v\n", err)
    }

    fmt.Println("All workers stopped")
}

func worker(ctx context.Context, id int) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Printf("Worker %d working...\n", id)
        case <-ctx.Done():
            fmt.Printf("Worker %d received cancellation, cleaning up...\n", id)
            time.Sleep(500 * time.Millisecond) // 模拟清理
            fmt.Printf("Worker %d stopped\n", id)
            return
        }
    }
}
```
