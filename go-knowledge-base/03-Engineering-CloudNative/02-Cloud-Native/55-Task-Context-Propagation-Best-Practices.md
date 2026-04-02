# 任务上下文传播最佳实践 (Task Context Propagation Best Practices)

> **分类**: 工程与云原生
> **标签**: #best-practices #context #propagation #guidelines

---

## 上下文传播黄金法则

```go
// 法则 1: 始终传播上下文
// 好
type GoodService struct{}

func (s *GoodService) Process(ctx context.Context, req *Request) (*Response, error) {
    // 传递上下文
    data, err := s.db.Query(ctx, req.Query)
    if err != nil {
        return nil, err
    }

    // 继续传递
    result, err := s.processor.Process(ctx, data)
    return result, err
}

// 坏
type BadService struct{}

func (s *BadService) Process(ctx context.Context, req *Request) (*Response, error) {
    // ❌ 丢失上下文
    data, err := s.db.Query(context.Background(), req.Query)
    // ...
}

// 法则 2: 不要存储上下文
// 好
type GoodTask struct {
    id string
}

func (t *GoodTask) Execute(ctx context.Context) error {
    // 使用传入的上下文
    return t.doWork(ctx)
}

// 坏
type BadTask struct {
    id  string
    ctx context.Context  // ❌ 存储上下文
}

func (t *BadTask) Execute() error {
    // 使用存储的上下文
    return t.doWork(t.ctx)
}

// 法则 3: 及时取消
// 好
func GoodFunction(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()  // 确保取消

    return doWork(ctx)
}

// 法则 4: 检查取消
// 好
func GoodLoop(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        if err := doIteration(ctx, i); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 传播边界处理

```go
// 异步边界
func AsyncBoundary(ctx context.Context, work func(context.Context)) {
    // 提取需要传递的值
    traceID, _ := ctx.Value(TraceIDKey).(string)
    tenantID, _ := ctx.Value(TenantIDKey).(string)

    go func() {
        // 创建新的上下文，携带必要的值
        newCtx := context.Background()
        newCtx = context.WithValue(newCtx, TraceIDKey, traceID)
        newCtx = context.WithValue(newCtx, TenantIDKey, tenantID)

        work(newCtx)
    }()
}

// 缓存边界
func CacheBoundary(ctx context.Context, key string, loader func(context.Context) (interface{}, error)) (interface{}, error) {
    // 缓存键包含上下文中的身份标识
    tenantID, _ := ctx.Value(TenantIDKey).(string)
    cacheKey := fmt.Sprintf("%s:%s", tenantID, key)

    // 使用无限制的上下文访问缓存
    if cached, ok := cache.Get(cacheKey); ok {
        return cached, nil
    }

    // 使用原始上下文加载
    value, err := loader(ctx)
    if err != nil {
        return nil, err
    }

    cache.Set(cacheKey, value)
    return value, nil
}

// 数据库边界
type DBContext struct {
    db *sql.DB
}

func (dbc *DBContext) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    // 添加上下文元数据到查询
    if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
        // 使用 RLS 或其他租户隔离机制
        ctx = context.WithValue(ctx, "rls.tenant_id", tenantID)
    }

    return dbc.db.QueryContext(ctx, query, args...)
}
```

---

## 性能优化

```go
// 1. 减少上下文查找次数
func OptimizedHandler(ctx context.Context) {
    // 一次性提取需要的值
    tenant, _ := TenantFromContext(ctx)
    user, _ := UserFromContext(ctx)
    traceID, _ := TraceIDFromContext(ctx)

    // 后续使用局部变量
    doSomethingWithTenant(tenant)
    doSomethingWithUser(user)
    logWithTraceID(traceID)
}

// 2. 延迟传播
func LazyPropagation(parent context.Context) context.Context {
    // 只在需要时创建子上下文
    return &lazyContext{
        parent:   parent,
        computed: make(map[interface{}]interface{}),
    }
}

type lazyContext struct {
    parent   context.Context
    computed map[interface{}]interface{}
    mu       sync.RWMutex
}

func (lc *lazyContext) Value(key interface{}) interface{} {
    // 先检查缓存
    lc.mu.RLock()
    if v, ok := lc.computed[key]; ok {
        lc.mu.RUnlock()
        return v
    }
    lc.mu.RUnlock()

    // 从父上下文获取
    v := lc.parent.Value(key)

    // 缓存结果
    lc.mu.Lock()
    lc.computed[key] = v
    lc.mu.Unlock()

    return v
}

// 3. 批量传播
func BatchPropagation(ctx context.Context, tasks []Task) []context.Context {
    // 提取一次公共值
    commonValues := extractCommonValues(ctx)

    contexts := make([]context.Context, len(tasks))
    for i, task := range tasks {
        // 基于公共值创建上下文
        taskCtx := context.WithValue(context.Background(), CommonValuesKey, commonValues)
        taskCtx = context.WithValue(taskCtx, TaskIDKey, task.ID)
        contexts[i] = taskCtx
    }

    return contexts
}
```

---

## 调试与监控

```go
// 上下文调试工具
type ContextDebugger struct {
    enabled bool
}

func (cd *ContextDebugger) Dump(ctx context.Context) string {
    if !cd.enabled {
        return ""
    }

    var info []string

    // 标准值
    if deadline, ok := ctx.Deadline(); ok {
        info = append(info, fmt.Sprintf("deadline=%v", deadline))
    }

    info = append(info, fmt.Sprintf("cancelled=%v", ctx.Err() != nil))

    // 自定义值
    ctx.Value(DebugKeyFunc(func(key interface{}) {
        info = append(info, fmt.Sprintf("%v=%v", key, ctx.Value(key)))
    }))

    return strings.Join(info, ", ")
}

// 上下文传播追踪
type PropagationTracer struct {
    tracer trace.Tracer
}

func (pt *PropagationTracer) TracePropagation(ctx context.Context, from, to string) context.Context {
    ctx, span := pt.tracer.Start(ctx, fmt.Sprintf("propagate:%s->%s", from, to))
    defer span.End()

    // 记录传播的上下文内容
    span.SetAttributes(
        attribute.String("propagation.from", from),
        attribute.String("propagation.to", to),
    )

    return ctx
}

// 性能监控
type ContextMetrics struct {
    valueLookupDuration prometheus.Histogram
    contextCreationDuration prometheus.Histogram
}

func (cm *ContextMetrics) RecordValueLookup(duration time.Duration) {
    cm.valueLookupDuration.Observe(duration.Seconds())
}

func (cm *ContextMetrics) RecordContextCreation(duration time.Duration) {
    cm.contextCreationDuration.Observe(duration.Seconds())
}
```
