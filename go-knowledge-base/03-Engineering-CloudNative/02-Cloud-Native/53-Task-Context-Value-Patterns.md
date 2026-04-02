# 任务上下文值模式 (Task Context Value Patterns)

> **分类**: 工程与云原生
> **标签**: #context #values #patterns #type-safety

---

## 类型安全的上下文值

```go
// 使用泛型实现类型安全的上下文值
package ctxval

import "context"

// Key 是强类型的上下文键
type Key[T any] struct {
    name string
}

func NewKey[T any](name string) Key[T] {
    return Key[T]{name: name}
}

func (k Key[T]) WithValue(ctx context.Context, value T) context.Context {
    return context.WithValue(ctx, k, value)
}

func (k Key[T]) Value(ctx context.Context) (T, bool) {
    var zero T
    v := ctx.Value(k)
    if v == nil {
        return zero, false
    }
    t, ok := v.(T)
    return t, ok
}

func (k Key[T]) MustValue(ctx context.Context) T {
    v, ok := k.Value(ctx)
    if !ok {
        panic("context value not found: " + k.name)
    }
    return v
}

// 使用示例
var (
    TraceIDKey = NewKey[string]("trace_id")
    TenantKey  = NewKey[Tenant]("tenant")
    UserKey    = NewKey[User]("user")
)

func Example() {
    ctx := context.Background()

    // 类型安全地设置值
    ctx = TraceIDKey.WithValue(ctx, "abc-123")
    ctx = TenantKey.WithValue(ctx, Tenant{ID: "t-1", Name: "Acme"})

    // 类型安全地获取值
    if traceID, ok := TraceIDKey.Value(ctx); ok {
        fmt.Println(traceID) // 自动推断为 string 类型
    }

    tenant := TenantKey.MustValue(ctx) // 类型安全，编译时检查
}
```

---

## 命名空间上下文值

```go
// 避免键冲突的命名空间模式
type Namespace string

type NamespacedKey struct {
    namespace Namespace
    key       string
}

func (n Namespace) Key(k string) NamespacedKey {
    return NamespacedKey{namespace: n, key: k}
}

func (nk NamespacedKey) String() string {
    return string(nk.namespace) + "/" + nk.key
}

// 预定义命名空间
const (
    NSRequest     Namespace = "request"
    NSTenant      Namespace = "tenant"
    NSAuth        Namespace = "auth"
    NSTelemetry   Namespace = "telemetry"
    NSExecution   Namespace = "execution"
)

// 使用
func SetRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, NSRequest.Key("id"), id)
}

func GetRequestID(ctx context.Context) string {
    v, _ := ctx.Value(NSRequest.Key("id")).(string)
    return v
}

func SetUserID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, NSAuth.Key("user_id"), id)
}

func GetUserID(ctx context.Context) string {
    v, _ := ctx.Value(NSAuth.Key("user_id")).(string)
    return v
}
```

---

## 上下文值访问器模式

```go
// 统一的上下文值访问接口
type ContextAccessor interface {
    Get(ctx context.Context) (interface{}, bool)
    Set(ctx context.Context, value interface{}) context.Context
    Delete(ctx context.Context) context.Context
}

// 实现示例：带默认值的访问器
type DefaultValueAccessor struct {
    key         interface{}
    defaultValue interface{}
}

func (dva *DefaultValueAccessor) Get(ctx context.Context) (interface{}, bool) {
    v := ctx.Value(dva.key)
    if v == nil {
        return dva.defaultValue, false
    }
    return v, true
}

func (dva *DefaultValueAccessor) Set(ctx context.Context, value interface{}) context.Context {
    return context.WithValue(ctx, dva.key, value)
}

// 计算值访问器
type ComputedValueAccessor struct {
    key      interface{}
    compute  func(context.Context) interface{}
}

func (cva *ComputedValueAccessor) Get(ctx context.Context) (interface{}, bool) {
    if v := ctx.Value(cva.key); v != nil {
        return v, true
    }

    computed := cva.compute(ctx)
    return computed, true
}

// 缓存访问器
type CachedValueAccessor struct {
    key      interface{}
    cache    sync.Map
    loader   func(context.Context) (interface{}, error)
}

func (cva *CachedValueAccessor) Get(ctx context.Context) (interface{}, bool) {
    // 尝试从 context 获取
    if v := ctx.Value(cva.key); v != nil {
        return v, true
    }

    // 尝试从缓存获取
    if cached, ok := cva.cache.Load(cva.key); ok {
        return cached, true
    }

    // 加载新值
    v, err := cva.loader(ctx)
    if err != nil {
        return nil, false
    }

    cva.cache.Store(cva.key, v)
    return v, true
}
```

---

## 上下文值验证

```go
// 带验证的上下文值设置
type ValidatedValue struct {
    key       interface{}
    validator func(interface{}) error
}

func (vv *ValidatedValue) Set(ctx context.Context, value interface{}) (context.Context, error) {
    if err := vv.validator(value); err != nil {
        return ctx, fmt.Errorf("validation failed: %w", err)
    }

    return context.WithValue(ctx, vv.key, value), nil
}

// 使用示例
var ValidatedTenantID = &ValidatedValue{
    key: "tenant_id",
    validator: func(v interface{}) error {
        id, ok := v.(string)
        if !ok {
            return fmt.Errorf("tenant_id must be string")
        }

        if !tenantIDRegex.MatchString(id) {
            return fmt.Errorf("invalid tenant_id format")
        }

        return nil
    },
}

func SetTenantID(ctx context.Context, tenantID string) (context.Context, error) {
    return ValidatedTenantID.Set(ctx, tenantID)
}
```
