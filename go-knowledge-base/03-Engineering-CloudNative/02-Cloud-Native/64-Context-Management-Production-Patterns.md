# 上下文管理生产模式 (Context Management Production Patterns)

> **分类**: 工程与云原生
> **标签**: #context #production #patterns #observability
> **参考**: Go Context 包设计, Google Context Best Practices, OpenTelemetry

---

## 上下文传播链

```go
// ContextChain 上下文传播链管理
package contextmgmt

import (
    "context"
    "time"

    "go.opentelemetry.io/otel/trace"
)

// ContextPropagator 传播器接口
type ContextPropagator interface {
    // Inject 将上下文注入载体
    Inject(ctx context.Context, carrier interface{}) error
    // Extract 从载体提取上下文
    Extract(ctx context.Context, carrier interface{}) (context.Context, error)
}

// PropagationChain 传播链
type PropagationChain struct {
    propagators []ContextPropagator
}

func (pc *PropagationChain) Inject(ctx context.Context, carrier interface{}) error {
    for _, p := range pc.propagators {
        if err := p.Inject(ctx, carrier); err != nil {
            return err
        }
    }
    return nil
}

func (pc *PropagationChain) Extract(ctx context.Context, carrier interface{}) (context.Context, error) {
    var err error
    for _, p := range pc.propagators {
        ctx, err = p.Extract(ctx, carrier)
        if err != nil {
            return ctx, err
        }
    }
    return ctx, nil
}

// MetadataCarrier 元数据载体 (用于 gRPC/HTTP)
type MetadataCarrier map[string]string

func (m MetadataCarrier) Get(key string) string {
    return m[key]
}

func (m MetadataCarrier) Set(key, value string) {
    m[key] = value
}

func (m MetadataCarrier) Keys() []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

---

## 请求上下文管理

```go
// RequestContext 请求上下文
type RequestContext struct {
    RequestID    string
    TraceID      string
    SpanID       string
    UserID       string
    TenantID     string
    StartTime    time.Time
    Deadline     time.Time
    Metadata     map[string]string
}

// ContextManager 上下文管理器
type ContextManager struct {
    timeout      time.Duration
    propagators  *PropagationChain
}

// NewIncomingContext 创建入站请求上下文
func (cm *ContextManager) NewIncomingContext(parent context.Context,
    carrier MetadataCarrier) (context.Context, context.CancelFunc) {

    // 提取传播的数据
    ctx := parent

    // 设置请求ID
    if requestID := carrier.Get("x-request-id"); requestID != "" {
        ctx = WithRequestID(ctx, requestID)
    } else {
        ctx = WithRequestID(ctx, generateRequestID())
    }

    // 设置租户信息
    if tenantID := carrier.Get("x-tenant-id"); tenantID != "" {
        ctx = WithTenant(ctx, Tenant{ID: tenantID})
    }

    // 设置用户信息
    if userID := carrier.Get("x-user-id"); userID != "" {
        ctx = WithUser(ctx, User{ID: userID})
    }

    // 添加上下文传播
    ctx, cancel := context.WithTimeout(ctx, cm.timeout)

    return ctx, cancel
}

// NewOutgoingContext 创建出站请求上下文
func (cm *ContextManager) NewOutgoingContext(ctx context.Context) MetadataCarrier {
    carrier := make(MetadataCarrier)

    // 传播请求ID
    if requestID, ok := GetRequestID(ctx); ok {
        carrier.Set("x-request-id", requestID)
    }

    // 传播租户信息
    if tenant, ok := GetTenant(ctx); ok {
        carrier.Set("x-tenant-id", tenant.ID)
    }

    // 传播用户信息
    if user, ok := GetUser(ctx); ok {
        carrier.Set("x-user-id", user.ID)
    }

    // 传播追踪信息
    if span := trace.SpanFromContext(ctx); span != nil {
        spanContext := span.SpanContext()
        if spanContext.IsValid() {
            carrier.Set("traceparent", fmt.Sprintf("00-%s-%s-%s",
                spanContext.TraceID().String(),
                spanContext.SpanID().String(),
                spanContext.TraceFlags()))
        }
    }

    return carrier
}

// BackgroundWithValues 从上下文提取值创建后台上下文
func BackgroundWithValues(ctx context.Context) context.Context {
    bg := context.Background()

    // 复制关键值
    if requestID, ok := GetRequestID(ctx); ok {
        bg = WithRequestID(bg, requestID)
    }

    if tenant, ok := GetTenant(ctx); ok {
        bg = WithTenant(bg, tenant)
    }

    return bg
}
```

---

## 上下文装饰器模式

```go
// ContextDecorator 上下文装饰器
type ContextDecorator func(context.Context) context.Context

// DecoratorChain 装饰器链
type DecoratorChain struct {
    decorators []ContextDecorator
}

func (dc *DecoratorChain) Decorate(ctx context.Context) context.Context {
    for _, d := range dc.decorators {
        ctx = d(ctx)
    }
    return ctx
}

func (dc *DecoratorChain) Add(d ContextDecorator) {
    dc.decorators = append(dc.decorators, d)
}

// 常用装饰器

// WithTimeoutDecorator 超时装饰器
func WithTimeoutDecorator(timeout time.Duration) ContextDecorator {
    return func(ctx context.Context) context.Context {
        ctx, _ = context.WithTimeout(ctx, timeout)
        return ctx
    }
}

// WithCancelDecorator 取消装饰器
func WithCancelDecorator() ContextDecorator {
    return func(ctx context.Context) context.Context {
        ctx, _ = context.WithCancel(ctx)
        return ctx
    }
}

// WithValueDecorator 值装饰器
func WithValueDecorator(key, value interface{}) ContextDecorator {
    return func(ctx context.Context) context.Context {
        return context.WithValue(ctx, key, value)
    }
}

// WithTelemetryDecorator 遥测装饰器
func WithTelemetryDecorator(tracer trace.Tracer) ContextDecorator {
    return func(ctx context.Context) context.Context {
        ctx, span := tracer.Start(ctx, "operation")
        defer span.End()
        return ctx
    }
}

// WithLoggerDecorator 日志装饰器
func WithLoggerDecorator(logger *zap.Logger) ContextDecorator {
    return func(ctx context.Context) context.Context {
        // 添加上下文相关的日志字段
        fields := []zap.Field{}

        if requestID, ok := GetRequestID(ctx); ok {
            fields = append(fields, zap.String("request_id", requestID))
        }

        return WithLogger(ctx, logger.With(fields...))
    }
}

// Middleware 风格的上下文装饰 (用于 HTTP/gRPC)
func ContextMiddleware(decorators ...ContextDecorator) func(http.Handler) http.Handler {
    chain := &DecoratorChain{decorators: decorators}

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := chain.Decorate(r.Context())
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 上下文池化

```go
// ContextPool 上下文对象池
type ContextPool struct {
    pool sync.Pool
}

func NewContextPool() *ContextPool {
    return &ContextPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &pooledContext{
                    values: make(map[interface{}]interface{}),
                }
            },
        },
    }
}

func (cp *ContextPool) Get(parent context.Context) *pooledContext {
    pc := cp.pool.Get().(*pooledContext)
    pc.parent = parent
    pc.done = make(chan struct{})
    return pc
}

func (cp *ContextPool) Put(pc *pooledContext) {
    // 清理
    pc.parent = nil
    for k := range pc.values {
        delete(pc.values, k)
    }
    cp.pool.Put(pc)
}

type pooledContext struct {
    parent context.Context
    values map[interface{}]interface{}
    done   chan struct{}
    mu     sync.RWMutex
}

func (pc *pooledContext) Deadline() (time.Time, bool) {
    return pc.parent.Deadline()
}

func (pc *pooledContext) Done() <-chan struct{} {
    return pc.done
}

func (pc *pooledContext) Err() error {
    select {
    case <-pc.done:
        return context.Canceled
    default:
        return nil
    }
}

func (pc *pooledContext) Value(key interface{}) interface{} {
    pc.mu.RLock()
    defer pc.mu.RUnlock()

    if v, ok := pc.values[key]; ok {
        return v
    }
    return pc.parent.Value(key)
}

func (pc *pooledContext) SetValue(key, value interface{}) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    pc.values[key] = value
}

func (pc *pooledContext) Cancel() {
    close(pc.done)
}
```

---

## 上下文传播边界

```go
// Boundary 边界处理

// AsyncBoundary 异步边界处理
func AsyncBoundary(ctx context.Context, fn func(context.Context)) {
    // 提取需要传播的值
    values := extractContextValues(ctx)

    go func() {
        // 创建新的上下文
        newCtx := context.Background()

        // 传播值
        for _, v := range values {
            newCtx = v.Inject(newCtx)
        }

        fn(newCtx)
    }()
}

// CacheBoundary 缓存边界处理
func CacheBoundary(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
    // 提取上下文中的缓存键部分
    tenantID, _ := GetTenantID(ctx)
    cacheKey := fmt.Sprintf("%s:%s", tenantID, key)

    // 使用无上下文限制的上下文访问缓存
    return cache.Get(cacheKey)
}

// DatabaseBoundary 数据库边界处理
func DatabaseBoundary(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    // 提取 RLS (Row Level Security) 信息
    tenantID, hasTenant := GetTenantID(ctx)

    if hasTenant {
        // 添加 RLS 谓词
        query = addRLSPredicate(query, tenantID)
    }

    // 使用原始上下文执行
    return db.QueryContext(ctx, query, args...)
}

// ExternalServiceBoundary 外部服务边界处理
func ExternalServiceBoundary(ctx context.Context, timeout time.Duration,
    fn func(context.Context) error) error {

    // 为外部调用创建独立的超时上下文
    // 但不传播取消信号
    deadline, hasDeadline := ctx.Deadline()
    if hasDeadline {
        remaining := time.Until(deadline)
        if remaining < timeout {
            timeout = remaining
        }
    }

    newCtx, cancel := context.WithTimeout(WithoutCancel(ctx), timeout)
    defer cancel()

    return fn(newCtx)
}

// WithoutCancel 创建不继承取消信号的上下文
func WithoutCancel(ctx context.Context) context.Context {
    return &withoutCancelCtx{ctx}
}

type withoutCancelCtx struct {
    context.Context
}

func (c *withoutCancelCtx) Deadline() (time.Time, bool) {
    return time.Time{}, false
}

func (c *withoutCancelCtx) Done() <-chan struct{} {
    return nil
}

func (c *withoutCancelCtx) Err() error {
    return nil
}
```

---

## 上下文监控与诊断

```go
// ContextMonitor 上下文监控
type ContextMonitor struct {
    activeContexts prometheus.Gauge
    contextDuration prometheus.Histogram
    cancelledContexts prometheus.Counter
    timeoutContexts prometheus.Counter
}

func NewContextMonitor() *ContextMonitor {
    return &ContextMonitor{
        activeContexts: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "context_active_total",
            Help: "Number of active contexts",
        }),
        contextDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "context_duration_seconds",
            Help:    "Context lifetime duration",
            Buckets: prometheus.DefBuckets,
        }),
        cancelledContexts: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "context_cancelled_total",
            Help: "Total number of cancelled contexts",
        }),
        timeoutContexts: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "context_timeout_total",
            Help: "Total number of timeout contexts",
        }),
    }
}

// MonitoredContext 带监控的上下文
func (cm *ContextMonitor) MonitoredContext(parent context.Context) (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(parent)

    cm.activeContexts.Inc()
    start := time.Now()

    // 包装取消函数
    monitoredCancel := func() {
        cancel()
        cm.activeContexts.Dec()
        cm.contextDuration.Observe(time.Since(start).Seconds())
        cm.cancelledContexts.Inc()
    }

    // 监听取消
    go func() {
        <-ctx.Done()
        if ctx.Err() == context.DeadlineExceeded {
            cm.timeoutContexts.Inc()
        }
        cm.activeContexts.Dec()
        cm.contextDuration.Observe(time.Since(start).Seconds())
    }()

    return ctx, monitoredCancel
}

// ContextDebugger 上下文调试器
type ContextDebugger struct {
    enabled bool
}

func (cd *ContextDebugger) DumpContext(ctx context.Context) map[string]interface{} {
    if !cd.enabled {
        return nil
    }

    info := map[string]interface{}{
        "deadline": nil,
        "cancelled": false,
        "values":    make(map[string]interface{}),
    }

    if deadline, ok := ctx.Deadline(); ok {
        info["deadline"] = deadline.Format(time.RFC3339)
        info["remaining"] = time.Until(deadline).String()
    }

    select {
    case <-ctx.Done():
        info["cancelled"] = true
        info["error"] = ctx.Err().Error()
    default:
    }

    // 提取所有值
    values := info["values"].(map[string]interface{})
    cd.extractValues(ctx, values)

    return info
}

func (cd *ContextDebugger) extractValues(ctx context.Context, values map[string]interface{}) {
    // 遍历所有值
    // 注意：这使用了反射，仅在调试时使用
    // ...
}

// ContextChainDebugger 上下文链调试器
func (cd *ContextDebugger) DumpContextChain(ctx context.Context) []map[string]interface{} {
    var chain []map[string]interface{}

    for ctx != nil {
        info := cd.DumpContext(ctx)
        chain = append(chain, info)

        // 获取父上下文
        // 注意：这需要反射，仅用于调试
        ctx = reflectParent(ctx)
    }

    return chain
}
```

---

## 上下文测试模式

```go
// TestContext 测试上下文
type TestContext struct {
    context.Context
    mu     sync.RWMutex
    values map[interface{}]interface{}
}

func NewTestContext() *TestContext {
    return &TestContext{
        Context: context.Background(),
        values:  make(map[interface{}]interface{}),
    }
}

func (tc *TestContext) Value(key interface{}) interface{} {
    tc.mu.RLock()
    defer tc.mu.RUnlock()

    if v, ok := tc.values[key]; ok {
        return v
    }
    return tc.Context.Value(key)
}

func (tc *TestContext) SetValue(key, value interface{}) {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    tc.values[key] = value
}

// MockContext 模拟上下文
func MockContext() context.Context {
    ctx := context.Background()
    ctx = WithRequestID(ctx, "test-request-id")
    ctx = WithTenant(ctx, Tenant{ID: "test-tenant"})
    ctx = WithUser(ctx, User{ID: "test-user"})
    return ctx
}

// ContextMatcher 上下文匹配器 (用于 mock)
type ContextMatcher struct {
    expectedValues map[interface{}]interface{}
}

func (m *ContextMatcher) Matches(x interface{}) bool {
    ctx, ok := x.(context.Context)
    if !ok {
        return false
    }

    for key, expectedValue := range m.expectedValues {
        if ctx.Value(key) != expectedValue {
            return false
        }
    }

    return true
}

func ContextWithValue(key, value interface{}) *ContextMatcher {
    return &ContextMatcher{
        expectedValues: map[interface{}]interface{}{
            key: value,
        },
    }
}
```
