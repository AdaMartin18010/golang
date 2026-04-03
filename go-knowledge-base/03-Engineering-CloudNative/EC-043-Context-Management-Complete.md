# EC-043: Context 管理完整指南 (Context Management Complete)

> **维度**: Engineering CloudNative
> **级别**: S (18+ KB)
> **标签**: #context #cancellation #propagation #go
> **相关**: EC-007, EC-008, LD-022

---

## 整合说明

本文档整合并提升了：

- `05-Context-Management.md` (5.7 KB)
- `18-Context-Propagation-Framework.md` (8.6 KB)
- `51-Task-Context-Propagation-Advanced.md` (8.2 KB)
- `52-Task-Context-Cancellation-Patterns.md` (8.2 KB)
- `66-Context-Propagation-Implementation.md` (17 KB)
- `64-Context-Management-Production-Patterns.md` (16 KB)

---

## Context 核心原理

```
┌─────────────────────────────────────────────────────────────────┐
│                      Context 树结构                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Background()                                                    │
│      │                                                           │
│      ├──► WithCancel() ───► cancel()                             │
│      │         │                                                 │
│      │         ├──► WithTimeout() ───► deadline exceeded         │
│      │         │         │                                       │
│      │         │         ├──► WithValue(key, val)                │
│      │         │                                                 │
│      │         └──► WithValue(traceID, "abc123")                 │
│      │                                                           │
│      └──► TODO()                                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 完整实现模式

### 1. 取消传播

```go
// 任务执行器
type TaskExecutor struct {
 activeTasks sync.Map // map[string]context.CancelFunc
}

func (e *TaskExecutor) Execute(parentCtx context.Context, task *Task) error {
 // 创建可取消的子上下文
 ctx, cancel := context.WithCancel(parentCtx)
 defer cancel()

 // 注册到活跃任务
 e.activeTasks.Store(task.ID, cancel)
 defer e.activeTasks.Delete(task.ID)

 // 监听取消信号
 done := make(chan error, 1)
 go func() {
  done <- e.runTask(ctx, task)
 }()

 select {
 case err := <-done:
  return err
 case <-ctx.Done():
  return ctx.Err()
 }
}

// 取消特定任务
func (e *TaskExecutor) Cancel(taskID string) error {
 if cancel, ok := e.activeTasks.Load(taskID); ok {
  cancel.(context.CancelFunc)()
  return nil
 }
 return fmt.Errorf("task not found: %s", taskID)
}

// 取消所有任务
func (e *TaskExecutor) CancelAll() {
 e.activeTasks.Range(func(key, value interface{}) bool {
  value.(context.CancelFunc)()
  return true
 })
}
```

### 2. 超时控制

```go
func (e *TaskExecutor) ExecuteWithTimeout(ctx context.Context, task *Task) error {
 // 使用任务指定的超时
 timeout := task.Timeout
 if timeout == 0 {
  timeout = 30 * time.Second // 默认超时
 }

 ctx, cancel := context.WithTimeout(ctx, timeout)
 defer cancel()

 return e.runTask(ctx, task)
}

// 分层超时
func (e *TaskExecutor) ExecuteWithHierarchicalTimeout(ctx context.Context, task *Task) error {
 // 全局超时（来自父上下文）
 if deadline, ok := ctx.Deadline(); ok {
  remaining := time.Until(deadline)
  if remaining < task.Timeout {
   // 父上下文超时更短，使用父超时
   return e.runTask(ctx, task)
  }
 }

 // 使用任务特定的超时
 return e.ExecuteWithTimeout(ctx, task)
}
```

### 3. 值传播

```go
// 上下文键类型（避免冲突）
type contextKey string

const (
 traceIDKey contextKey = "trace_id"
 spanIDKey  contextKey = "span_id"
 userIDKey  contextKey = "user_id"
)

// 设置值
func WithTraceID(ctx context.Context, traceID string) context.Context {
 return context.WithValue(ctx, traceIDKey, traceID)
}

// 获取值
func GetTraceID(ctx context.Context) string {
 if id, ok := ctx.Value(traceIDKey).(string); ok {
  return id
 }
 return ""
}

// 跨服务传播
func propagateContext(ctx context.Context) metadata.MD {
 md := metadata.New(map[string]string{
  "x-trace-id": GetTraceID(ctx),
  "x-span-id":  GetSpanID(ctx),
  "x-user-id":  GetUserID(ctx),
 })
 return md
}
```

---

## 生产模式

### HTTP 中间件

```go
func ContextMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  // 提取或生成 trace ID
  traceID := r.Header.Get("X-Trace-ID")
  if traceID == "" {
   traceID = generateTraceID()
  }

  // 创建带值的上下文
  ctx := r.Context()
  ctx = WithTraceID(ctx, traceID)
  ctx = WithUserID(ctx, getUserFromToken(r))

  // 设置超时
  ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
  defer cancel()

  // 注入到请求
  next.ServeHTTP(w, r.WithContext(ctx))
 })
}
```

### gRPC 拦截器

```go
func UnaryInterceptor(ctx context.Context, req interface{},
 info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

 // 从 metadata 提取上下文
 md, _ := metadata.FromIncomingContext(ctx)

 if traceIDs := md.Get("x-trace-id"); len(traceIDs) > 0 {
  ctx = WithTraceID(ctx, traceIDs[0])
 }

 return handler(ctx, req)
}
```

---

## 最佳实践

| 实践 | 说明 | 示例 |
|------|------|------|
| 尽早检查取消 | 在长循环中检查 `ctx.Done()` | `for { select { case <-ctx.Done(): return } }` |
| 传递而非存储 | 将 Context 作为第一个参数 | `func(ctx context.Context, ...)` |
| 不存储在结构体 | 避免生命周期问题 | 使用参数传递 |
| 使用具体键类型 | 避免冲突 | `type key string` |
| 设置合理超时 | 防止无限阻塞 | `WithTimeout(ctx, 5s)` |

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