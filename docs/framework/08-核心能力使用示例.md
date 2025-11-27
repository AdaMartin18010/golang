# 核心能力使用示例

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本文档提供框架核心能力的完整使用示例，展示如何在实际项目中使用这些能力。

---

## 1. 数据库抽象层使用

### 基本使用

```go
import "github.com/yourusername/golang/pkg/database"

// 创建数据库连接
db, err := database.NewDatabase(database.Config{
    Driver:       database.DriverPostgreSQL,
    DSN:          "postgres://user:password@localhost/dbname?sslmode=disable",
    MaxOpenConns: 25,
    MaxIdleConns: 5,
})
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// 执行查询
rows, err := db.Query(ctx, "SELECT id, name FROM users WHERE id = $1", 1)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ID: %d, Name: %s\n", id, name)
}
```

### 使用事务

```go
// 开始事务
tx, err := db.Begin(ctx)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback()

// 在事务中执行操作
_, err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "John")
if err != nil {
    return err
}

// 提交事务
if err := tx.Commit(); err != nil {
    return err
}
```

---

## 2. 数据转换使用

### 类型转换

```go
import "github.com/yourusername/golang/pkg/converter"

conv := converter.NewConverter()

// 转换为字符串
str := conv.ToString(123)        // "123"
str = conv.ToString(true)        // "true"

// 转换为整数
num, _ := conv.ToInt("123")      // 123
num, _ = conv.ToInt(123.45)      // 123

// 转换为浮点数
f, _ := conv.ToFloat64("123.45") // 123.45

// 转换为布尔值
b, _ := conv.ToBool("true")      // true
```

### JSON 转换

```go
// 结构体转 JSON
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

user := User{ID: 1, Name: "John"}
jsonStr, _ := conv.ToJSON(user)

// JSON 转结构体
var result User
conv.FromJSON(jsonStr, &result)
```

---

## 3. 采样机制使用

### 概率采样

```go
import "github.com/yourusername/golang/pkg/sampling"

// 创建概率采样器（50% 采样率）
sampler, _ := sampling.NewProbabilisticSampler(0.5)

// 判断是否采样
if sampler.ShouldSample(ctx) {
    // 执行采样操作
    collectData()
}

// 动态调整采样率
sampler.UpdateRate(0.1) // 降低到 10%
```

### 自适应采样

```go
// 创建自适应采样器
sampler, _ := sampling.NewAdaptiveSampler(0.5, 0.1, 1.0)

// 根据系统负载调整
adaptiveSampler := sampler.(*sampling.AdaptiveSampler)
adaptiveSampler.AdjustForLoad(0.9) // 高负载，降低采样率
```

---

## 4. 追踪和定位使用

### 基本追踪

```go
import "github.com/yourusername/golang/pkg/tracing"

tracer := tracing.NewTracer("my-service")

// 开始 Span
ctx, span := tracer.StartSpan(ctx, "operation-name")
defer span.End()

// 添加属性
tracer.AddAttributes(span, map[string]interface{}{
    "user.id": 123,
    "operation.type": "create",
})

// 记录错误
if err != nil {
    tracer.RecordError(span, err)
}
```

### 错误定位

```go
// 自动记录错误的完整上下文
tracer.LocateError(ctx, err, map[string]interface{}{
    "user.id": 123,
    "request.id": "req-456",
})
// 这会自动记录：
// - 错误信息
// - 堆栈跟踪
// - 调用位置（文件、行号、函数名）
// - 自定义属性
```

---

## 5. 反射/自解释使用

### 类型检查

```go
import "github.com/yourusername/golang/pkg/reflect"

inspector := reflect.NewInspector()

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{}
metadata := inspector.InspectType(user)

fmt.Printf("Type: %s\n", metadata.Name)
fmt.Printf("Package: %s\n", metadata.Package)
fmt.Printf("Fields: %d\n", len(metadata.Fields))
```

### 函数检查

```go
func Add(a, b int) int {
    return a + b
}

metadata := inspector.InspectFunction(Add)
fmt.Printf("Function: %s\n", metadata.Name)
fmt.Printf("Inputs: %v\n", metadata.Inputs)
fmt.Printf("Outputs: %v\n", metadata.Outputs)
```

---

## 6. 精细控制使用

### 功能开关

```go
import "github.com/yourusername/golang/pkg/control"

controller := control.NewFeatureController()

// 注册功能
controller.Register("feature-a", "Feature A description", true, map[string]interface{}{
    "max_requests": 100,
})

// 启用/禁用功能
if controller.IsEnabled("feature-a") {
    // 执行功能
}

// 监听配置变化
controller.Watch("feature-a", func(config interface{}) {
    fmt.Printf("Config updated: %v\n", config)
})
```

### 速率控制

```go
rateController := control.NewRateController()

// 设置速率限制（每秒最多 100 次）
rateController.SetRateLimit("api-calls", 100.0, time.Second)

// 检查是否允许
if rateController.Allow("api-calls") {
    // 执行操作
}
```

### 熔断器

```go
circuitController := control.NewCircuitController()

// 注册熔断器
circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

// 记录成功/失败
circuitController.RecordSuccess("external-api")
circuitController.RecordFailure("external-api")

// 检查是否允许
if circuitController.Allow("external-api") {
    // 执行操作
}
```

---

## 7. 增强的 OTLP 集成使用

### 完整设置

```go
import (
    "github.com/yourusername/golang/pkg/observability/otlp"
    "github.com/yourusername/golang/pkg/sampling"
)

// 创建采样器
sampler, _ := sampling.NewProbabilisticSampler(0.5)

// 创建增强的 OTLP 集成
otlp, err := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:    "my-service",
    ServiceVersion: "v1.0.0",
    Endpoint:       "localhost:4317",
    Insecure:       true,
    Sampler:        sampler,
})
if err != nil {
    log.Fatal(err)
}
defer otlp.Shutdown(context.Background())

// 使用追踪器
tracer := otlp.Tracer("my-tracer")
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

// 使用指标器
meter := otlp.Meter("my-meter")
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1)
```

---

## 8. 完整示例：可观测的 HTTP 服务

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/pkg/control"
    "github.com/yourusername/golang/pkg/database"
    "github.com/yourusername/golang/pkg/observability/otlp"
    "github.com/yourusername/golang/pkg/sampling"
    "github.com/yourusername/golang/pkg/tracing"
)

func main() {
    ctx := context.Background()

    // 1. 设置可观测性
    sampler, _ := sampling.NewProbabilisticSampler(0.5)
    otlp, _ := otlp.NewEnhancedOTLP(otlp.Config{
        ServiceName: "my-service",
        Endpoint:    "localhost:4317",
        Sampler:     sampler,
    })
    defer otlp.Shutdown(ctx)

    tracer := tracing.NewTracer("my-service")

    // 2. 设置数据库
    db, _ := database.NewDatabase(database.Config{
        Driver: database.DriverPostgreSQL,
        DSN:    "postgres://...",
    })
    defer db.Close()

    // 3. 设置精细控制
    controller := control.NewFeatureController()
    controller.Register("feature-a", "Feature A", true, nil)

    rateController := control.NewRateController()
    rateController.SetRateLimit("api", 100.0, time.Second)

    // 4. 创建 HTTP 服务
    r := chi.NewRouter()

    r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
        // 开始追踪
        ctx, span := tracer.StartSpan(r.Context(), "get-users")
        defer span.End()

        // 速率控制
        if !rateController.Allow("api") {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        // 功能开关
        if !controller.IsEnabled("feature-a") {
            http.Error(w, "Feature disabled", http.StatusServiceUnavailable)
            return
        }

        // 执行操作
        rows, err := db.Query(ctx, "SELECT * FROM users")
        if err != nil {
            tracer.LocateError(ctx, err, map[string]interface{}{
                "endpoint": "/api/users",
            })
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // 处理结果...
    })

    http.ListenAndServe(":8080", r)
}
```

---

## 📚 相关文档

- [框架核心能力总结](07-框架核心能力总结.md)
- [框架最佳实践指南](06-最佳实践指南.md)
- [框架快速开始指南](05-快速开始指南.md)

---

**最后更新**: 2025-01-XX
