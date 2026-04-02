# context 包详解

> **分类**: 开源技术堆栈

---

## 核心功能

- **取消信号**: 传播取消
- **超时控制**: 设置截止时间
- **值传递**: 请求元数据

---

## 创建 Context

```go
// 根 context
ctx := context.Background()

// 带取消
ctx, cancel := context.WithCancel(ctx)
defer cancel()

// 带超时
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 带截止时间
ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1*time.Hour))
```

---

## 取消传播

```go
func process(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()  // context.Canceled or context.DeadlineExceeded
    case result := <-ch:
        return handle(result)
    }
}
```

---

## 值传递

```go
// 设置值
ctx := context.WithValue(parent, "user_id", "123")

// 获取值
userID := ctx.Value("user_id").(string)

// 建议: 使用自定义类型作为 key
type key int
const userIDKey key = 0
ctx = context.WithValue(ctx, userIDKey, "123")
```

---

## 最佳实践

1. **Context 作为第一个参数**
2. **不要存储在结构体中**
3. **值只用于请求元数据**
4. **及时调用 cancel()**
