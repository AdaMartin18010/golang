# 上下文管理 (Context Management)

> **分类**: 工程与云原生
> **标签**: #context #并发 #最佳实践

---

## 上下文传播模式

### 1. 显式传播模式

```go
// ✅ 推荐: context 作为第一个参数
func ProcessOrder(ctx context.Context, orderID string) error {
    // 向下传递
    user, err := GetUser(ctx, order.UserID)
    if err != nil {
        return err
    }

    // 继续传递
    return ChargePayment(ctx, user.ID, order.Amount)
}

func GetUser(ctx context.Context, userID string) (*User, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    // ...
}
```

### 2. 请求生命周期管理

```go
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
    // 基于请求创建上下文
    ctx := r.Context()

    // 添加请求级超时
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // 添加请求ID用于追踪
    ctx = WithRequestID(ctx, generateRequestID())

    // 处理请求
    result, err := ProcessRequest(ctx, r)
    // ...
}
```

---

## 上下文值管理

### 类型安全的键

```go
// 定义私有类型避免冲突
type contextKey int

const (
    requestIDKey contextKey = iota
    userIDKey
    traceIDKey
)

// 封装存取方法
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFrom(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}
```

### 值传播追踪

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 注入多个值
        ctx = WithRequestID(ctx, r.Header.Get("X-Request-ID"))
        ctx = WithUserID(ctx, getUserID(r))
        ctx = WithTraceID(ctx, generateTraceID())

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 取消传播策略

### 级联取消

```go
func ParentOperation(ctx context.Context) error {
    // 创建子上下文
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    errChan := make(chan error, 2)

    // 启动多个子任务
    go func() {
        errChan <- ChildOperationA(ctx)
    }()

    go func() {
        errChan <- ChildOperationB(ctx)
    }()

    // 任一任务失败则取消其他
    for i := 0; i < 2; i++ {
        if err := <-errChan; err != nil {
            cancel()  // 触发取消
            return err
        }
    }

    return nil
}
```

### 超时控制

```go
func QueryWithTimeout(ctx context.Context) error {
    // 数据库查询设置5秒超时
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    rows, err := db.QueryContext(ctx, "SELECT * FROM large_table")
    if err != nil {
        // 检查是否是超时错误
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("query timeout: %w", err)
        }
        return err
    }
    defer rows.Close()

    return processRows(ctx, rows)
}
```

---

## 上下文继承模式

### 上下文组合

```go
type CompositeContext struct {
    ctx    context.Context
    values map[interface{}]interface{}
    mu     sync.RWMutex
}

func (c *CompositeContext) Deadline() (time.Time, bool) {
    return c.ctx.Deadline()
}

func (c *CompositeContext) Done() <-chan struct{} {
    return c.ctx.Done()
}

func (c *CompositeContext) Err() error {
    return c.ctx.Err()
}

func (c *CompositeContext) Value(key interface{}) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if v, ok := c.values[key]; ok {
        return v
    }
    return c.ctx.Value(key)
}
```

---

## 最佳实践

### ✅ Do's

```go
// 1. 总是传递 context，不要存储
func Good(ctx context.Context) error {
    result, err := db.QueryContext(ctx, "SELECT ...")
    // ...
}

// 2. 及时调用 cancel
func Good() {
    ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
    defer cancel()  // 确保释放资源
    // ...
}

// 3. 检查 ctx.Err()
func Good(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-ch:
        return handle(result)
    }
}
```

### ❌ Don'ts

```go
// 1. 不要存储 context 在结构体中
type BadService struct {
    ctx context.Context  // ❌ 不要这样做
}

// 2. 不要传递 nil context
func Bad() {
    result, err := db.QueryContext(nil, "SELECT ...")  // ❌
}

// 3. 不要忽略 cancel
cancel()  // ❌ 没有 defer，可能忘记调用
```

---

## 测试上下文

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    err := LongRunningOperation(ctx)
    if err != context.DeadlineExceeded {
        t.Errorf("expected timeout, got: %v", err)
    }
}

func TestWithCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        time.Sleep(50 * time.Millisecond)
        cancel()
    }()

    err := Operation(ctx)
    if err != context.Canceled {
        t.Errorf("expected canceled, got: %v", err)
    }
}
```

---

## 调试技巧

```go
// 上下文追踪
func ContextInfo(ctx context.Context) string {
    var info []string

    if deadline, ok := ctx.Deadline(); ok {
        info = append(info, fmt.Sprintf("deadline=%v", time.Until(deadline)))
    }

    if reqID, ok := RequestIDFrom(ctx); ok {
        info = append(info, fmt.Sprintf("request_id=%s", reqID))
    }

    return strings.Join(info, ", ")
}
```
