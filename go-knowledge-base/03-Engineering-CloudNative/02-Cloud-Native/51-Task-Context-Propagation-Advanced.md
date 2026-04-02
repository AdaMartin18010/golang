# 任务上下文传播高级模式 (Advanced Task Context Propagation)

> **分类**: 工程与云原生
> **标签**: #context #propagation #distributed-tracing #advanced-patterns

---

## 上下文链与延续

```go
// 上下文链管理
type ContextChain struct {
    mu       sync.RWMutex
    links    []ContextLink
    carryOver map[string]CarryOverRule
}

type ContextLink struct {
    Name    string
    Context context.Context
    Cancel  context.CancelFunc
}

type CarryOverRule struct {
    Key         string
    PropagateTo []string  // 传播目标类型
    Transform   func(interface{}) interface{}
}

// 创建上下文延续
func (cc *ContextChain) Continue(ctx context.Context, linkName string) (context.Context, context.CancelFunc) {
    cc.mu.RLock()
    defer cc.mu.RUnlock()

    // 继承上游上下文的值
    newCtx := context.Background()

    for _, link := range cc.links {
        // 传播特定键
        if value := link.Context.Value(link.Name); value != nil {
            newCtx = context.WithValue(newCtx, link.Name, value)
        }
    }

    // 添加当前链节
    newCtx, cancel := context.WithCancel(newCtx)

    cc.mu.Lock()
    cc.links = append(cc.links, ContextLink{
        Name:    linkName,
        Context: newCtx,
        Cancel:  cancel,
    })
    cc.mu.Unlock()

    return newCtx, cancel
}

// 跨进程上下文序列化
func SerializeContext(ctx context.Context) (*ContextSnapshot, error) {
    snapshot := &ContextSnapshot{
        Timestamp: time.Now(),
        Values:    make(map[string]interface{}),
    }

    // 序列化可传播的值
    if traceID := ctx.Value(TraceIDKey); traceID != nil {
        snapshot.Values["trace_id"] = traceID
    }

    if spanID := ctx.Value(SpanIDKey); spanID != nil {
        snapshot.Values["span_id"] = spanID
    }

    if tenant := ctx.Value(TenantKey); tenant != nil {
        snapshot.Values["tenant_id"] = tenant
    }

    // 序列化 baggage
    if baggage, ok := BaggageFromContext(ctx); ok {
        snapshot.Baggage = baggage.ToMap()
    }

    return snapshot, nil
}

func DeserializeContext(snapshot *ContextSnapshot) context.Context {
    ctx := context.Background()

    for key, value := range snapshot.Values {
        ctx = context.WithValue(ctx, contextKey(key), value)
    }

    if len(snapshot.Baggage) > 0 {
        baggage := BaggageFromMap(snapshot.Baggage)
        ctx = ContextWithBaggage(ctx, baggage)
    }

    return ctx
}
```

---

## 上下文感知调度

```go
// 根据上下文属性进行调度决策
type ContextAwareScheduler struct {
    defaultScheduler Scheduler
    affinityRules    []AffinityRule
}

type AffinityRule struct {
    Match   func(context.Context) bool
    Select  func([]Worker, context.Context) Worker
}

func (cas *ContextAwareScheduler) Schedule(ctx context.Context, task *Task) error {
    workers := cas.getAvailableWorkers()

    // 应用亲和性规则
    for _, rule := range cas.affinityRules {
        if rule.Match(ctx) {
            selected := rule.Select(workers, ctx)
            return cas.assignToWorker(ctx, task, selected)
        }
    }

    // 默认调度
    return cas.defaultScheduler.Schedule(ctx, task)
}

// 租户亲和性
func TenantAffinityRule() AffinityRule {
    return AffinityRule{
        Match: func(ctx context.Context) bool {
            _, ok := TenantFromContext(ctx)
            return ok
        },
        Select: func(workers []Worker, ctx context.Context) Worker {
            tenant, _ := TenantFromContext(ctx)

            // 优先选择已运行该租户任务的 worker
            for _, w := range workers {
                if w.HasTenant(tenant.TenantID) {
                    return w
                }
            }

            // 选择负载最轻的 worker
            return selectLeastLoaded(workers)
        },
    }
}

// 数据局部性亲和性
func DataLocalityRule(dataIndex DataIndex) AffinityRule {
    return AffinityRule{
        Match: func(ctx context.Context) bool {
            return ctx.Value(DataLocationKey) != nil
        },
        Select: func(workers []Worker, ctx context.Context) Worker {
            location := ctx.Value(DataLocationKey).(DataLocation)

            // 找到数据所在或最近的 worker
            nearest := workers[0]
            minLatency := time.Duration(1<<63 - 1)

            for _, w := range workers {
                latency := dataIndex.GetLatency(location, w.Location)
                if latency < minLatency {
                    minLatency = latency
                    nearest = w
                }
            }

            return nearest
        },
    }
}
```

---

## 上下文超时传播

```go
// 分布式超时管理
type TimeoutPropagator struct {
    clock Clock
}

func (tp *TimeoutPropagator) PropagateTimeout(parent context.Context, childTimeout time.Duration) (context.Context, context.CancelFunc) {
    // 获取父上下文剩余时间
    deadline, hasDeadline := parent.Deadline()

    if !hasDeadline {
        // 父上下文无超时，使用子超时
        return context.WithTimeout(parent, childTimeout)
    }

    remaining := time.Until(deadline)

    if remaining <= 0 {
        // 父上下文已过期
        ctx, cancel := context.WithCancel(parent)
        cancel()
        return ctx, func() {}
    }

    // 取较小值
    effectiveTimeout := min(remaining, childTimeout)

    // 添加传播路径信息
    ctx, cancel := context.WithTimeout(parent, effectiveTimeout)
    ctx = context.WithValue(ctx, TimeoutPathKey, TimeoutPath{
        Original:    childTimeout,
        Effective:   effectiveTimeout,
        ParentRemaining: remaining,
        PropagatedAt: tp.clock.Now(),
    })

    return ctx, cancel
}

// 自适应超时
func AdaptiveTimeout(ctx context.Context, historicalData []ExecutionTime) (context.Context, context.CancelFunc) {
    // 基于历史执行时间计算合适的超时
    avg := calculateAverage(historicalData)
    p99 := calculatePercentile(historicalData, 99)
    stdDev := calculateStdDev(historicalData)

    // 使用 p99 + 2*标准差作为超时
    timeout := p99 + 2*stdDev

    // 确保至少有一定余量
    minTimeout := avg * 3
    if timeout < minTimeout {
        timeout = minTimeout
    }

    return context.WithTimeout(ctx, timeout)
}
```

---

## 上下文安全检查

```go
// 上下文安全验证
type ContextSecurityChecker struct {
    validators []ContextValidator
}

type ContextValidator interface {
    Validate(ctx context.Context) error
    Name() string
}

// 租户隔离验证
func TenantIsolationValidator() ContextValidator {
    return &tenantValidator{}
}

type tenantValidator struct{}

func (tv *tenantValidator) Validate(ctx context.Context) error {
    tenant, ok := TenantFromContext(ctx)
    if !ok {
        return fmt.Errorf("tenant not found in context")
    }

    // 验证租户 ID 格式
    if !isValidTenantID(tenant.TenantID) {
        return fmt.Errorf("invalid tenant ID format")
    }

    // 验证租户状态
    if tenant.Status != "active" {
        return fmt.Errorf("tenant is not active: %s", tenant.Status)
    }

    return nil
}

func (tv *tenantValidator) Name() string {
    return "TenantIsolation"
}

// 链路完整性验证
func TraceIntegrityValidator() ContextValidator {
    return &traceValidator{}
}

type traceValidator struct{}

func (tv *traceValidator) Validate(ctx context.Context) error {
    traceID, ok := ctx.Value(TraceIDKey).(string)
    if !ok || traceID == "" {
        return fmt.Errorf("trace ID missing from context")
    }

    // 验证 trace ID 格式 (W3C)
    if !w3cTraceIDRegex.MatchString(traceID) {
        return fmt.Errorf("invalid trace ID format")
    }

    return nil
}

// 应用验证
func (csc *ContextSecurityChecker) Check(ctx context.Context) error {
    for _, validator := range csc.validators {
        if err := validator.Validate(ctx); err != nil {
            return fmt.Errorf("%s validation failed: %w", validator.Name(), err)
        }
    }
    return nil
}
```
