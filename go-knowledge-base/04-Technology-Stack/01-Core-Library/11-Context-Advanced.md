# Context 高级模式

> **分类**: 开源技术堆栈  
> **标签**: #context #advanced #patterns

---

## Context 树结构

```
Background
    └── WithCancel
            ├── WithTimeout(5s)
            │       └── WithValue(requestID)
            └── WithDeadline
                    └── WithValue(userID)
```

---

## 派生策略

### 独立超时

```go
func ParentHandler(ctx context.Context) {
    // 子操作独立超时
    dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    result, err := db.Query(dbCtx, "SELECT ...")
    
    // 另一个子操作
    cacheCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()
    cached, err := cache.Get(cacheCtx, key)
}
```

### 组合取消

```go
func CombinedContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(parent)
    
    go func() {
        select {
        case <-parent.Done():
            cancel()
        case <-time.After(timeout):
            cancel()
        }
    }()
    
    return ctx, cancel
}
```

---

## Context 值的最佳实践

### 类型安全封装

```go
package contextutil

type key int

const (
    requestIDKey key = iota
    traceIDKey
    userIDKey
)

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) string {
    id, _ := ctx.Value(requestIDKey).(string)
    return id
}

func WithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userIDKey, user)
}

func User(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(userIDKey).(*User)
    return user, ok
}
```

---

## Context 装饰器

```go
func LoggingContext(ctx context.Context, logger *zap.Logger) context.Context {
    return &loggingContext{
        Context: ctx,
        logger:  logger,
    }
}

type loggingContext struct {
    context.Context
    logger *zap.Logger
}

func (c *loggingContext) Value(key interface{}) interface{} {
    c.logger.Debug("context value accessed", zap.Any("key", key))
    return c.Context.Value(key)
}
```

---

## Context 链追踪

```go
type contextInfo struct {
    context.Context
    name string
    depth int
}

func (c *contextInfo) String() string {
    var path []string
    
    curr := c.Context
    for curr != nil {
        if info, ok := curr.(*contextInfo); ok {
            path = append(path, info.name)
        }
        curr = contextParent(curr)
    }
    
    return strings.Join(path, " -> ")
}

func WithName(ctx context.Context, name string) context.Context {
    depth := 0
    if info, ok := ctx.(*contextInfo); ok {
        depth = info.depth + 1
    }
    
    return &contextInfo{
        Context: ctx,
        name:    name,
        depth:   depth,
    }
}
```

---

## 性能优化

### 避免频繁创建

```go
// ❌ 不好
for _, item := range items {
    ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
    process(ctx, item)
    cancel()
}

// ✅ 好
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()

for _, item := range items {
    select {
    case <-parentCtx.Done():
        return parentCtx.Err()
    case <-timer.C:
        return context.DeadlineExceeded
    default:
        process(parentCtx, item)
    }
}
```
