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

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02