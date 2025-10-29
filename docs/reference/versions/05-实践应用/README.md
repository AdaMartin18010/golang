# Go新特性实践应用

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.21+

---

## 📋 目录

- [📋 概述](#概述)
- [🎯 最佳实践](#最佳实践)
  - [泛型使用最佳实践](#泛型使用最佳实践)
    - [✅ 何时使用泛型](#何时使用泛型)
    - [❌ 何时避免泛型](#何时避免泛型)
  - [结构化日志最佳实践](#结构化日志最佳实践)
    - [✅ slog使用建议](#slog使用建议)
    - [🎯 日志级别策略](#日志级别策略)
  - [for循环变量最佳实践](#for循环变量最佳实践)
    - [✅ Go 1.22+ 正确做法](#go-1.22-正确做法)
    - [⚠️ 向后兼容注意](#向后兼容注意)
- [🚀 迁移指南](#迁移指南)
  - [从Go 1.20迁移到1.21+](#从go-1.20迁移到1.21)
    - [步骤1: 更新Go版本](#步骤1-更新go版本)
    - [步骤2: 利用新特性](#步骤2-利用新特性)
    - [步骤3: 迁移到slog](#步骤3-迁移到slog)
  - [从Go 1.21迁移到1.22+](#从go-1.21迁移到1.22)
    - [主要变更](#主要变更)
  - [从Go 1.22迁移到1.23+](#从go-1.22迁移到1.23)
    - [实验性特性尝试](#实验性特性尝试)
- [⚠️ 常见问题与陷阱](#常见问题与陷阱)
  - [1. 泛型性能问题](#1.-泛型性能问题)
  - [2. slog内存分配](#2.-slog内存分配)
  - [3. for循环边界情况](#3.-for循环边界情况)
- [🔗 相关资源](#相关资源)
  - [迁移工具](#迁移工具)
  - [文档](#文档)
  - [社区](#社区)

## 📋 概述

本文档总结Go新特性（1.21-1.25）在实际项目中的应用经验，包括：

- ✅ **最佳实践**: 如何正确使用新特性
- 🔄 **迁移指南**: 平滑升级到新版本
- ⚠️ **常见陷阱**: 避免踩坑
- 💡 **实战案例**: 真实项目经验

---

## 🎯 最佳实践

### 泛型使用最佳实践

#### ✅ 何时使用泛型

**适合场景**：

```go
// 1. 容器类型
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// 2. 算法函数
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// 3. 数据结构
type Pair[K, V any] struct {
    Key   K
    Value V
}

func NewPair[K, V any](k K, v V) Pair[K, V] {
    return Pair[K, V]{Key: k, Value: v}
}
```

#### ❌ 何时避免泛型

```go
// ❌ 避免：过度抽象
// 坏例子：为单一类型使用泛型
type StringProcessor[T string] struct {  // 不必要
    data T
}

// ✅ 改进：直接使用具体类型
type StringProcessor struct {
    data string
}

// ❌ 避免：复杂的约束
type ComplexConstraint[T interface {
    comparable
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
    Method1() string
    Method2() int
}] struct {
    value T
}

// ✅ 改进：拆分或使用接口
type Processor interface {
    Method1() string
    Method2() int
}
```

---

### 结构化日志最佳实践

#### ✅ slog使用建议

```go
package main

import (
    "context"
    "log/slog"
    "os"
)

// 1. 使用全局logger（简单场景）
func simpleUsage() {
    slog.Info("application started", "version", "1.0.0")
    slog.Error("connection failed", "error", "timeout")
}

// 2. 使用自定义logger（生产环境）
func productionUsage() {
    // JSON格式，适合日志收集
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,  // 添加源代码位置
    }))
    
    slog.SetDefault(logger)
    
    // 使用上下文
    ctx := context.Background()
    logger.InfoContext(ctx, "request processed",
        slog.Group("request",
            slog.String("method", "GET"),
            slog.String("path", "/api/users"),
        ),
        slog.Int("status", 200),
        slog.Duration("latency", 45*time.Millisecond),
    )
}

// 3. 带上下文的logger
func contextLogger(ctx context.Context) *slog.Logger {
    // 从context中提取trace_id等信息
    logger := slog.Default()
    if traceID, ok := ctx.Value("trace_id").(string); ok {
        logger = logger.With("trace_id", traceID)
    }
    return logger
}

// 4. 分层日志
type Service struct {
    logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
    return &Service{
        logger: logger.With("component", "service"),
    }
}

func (s *Service) ProcessRequest(ctx context.Context, userID string) error {
    s.logger.Info("processing request", "user_id", userID)
    
    // 业务逻辑...
    
    return nil
}
```

#### 🎯 日志级别策略

```go
// 开发环境：DEBUG
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// 生产环境：INFO
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))

// 日志级别使用建议
slog.Debug("detailed debugging info")      // 开发调试
slog.Info("normal operation")              // 正常信息
slog.Warn("deprecated API called")         // 警告
slog.Error("operation failed", err)        // 错误
```

---

### for循环变量最佳实践

#### ✅ Go 1.22+ 正确做法

```go
// ✅ Go 1.22+: 直接使用，无需担心
func processItems(items []string) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func() {
            defer wg.Done()
            process(item)  // 安全：每次迭代item是新变量
        }()
    }
    wg.Wait()
}

// ✅ 闭包中使用
func createHandlers(routes []string) []http.HandlerFunc {
    handlers := make([]http.HandlerFunc, len(routes))
    for i, route := range routes {
        handlers[i] = func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Route: %s", route)  // 安全
        }
    }
    return handlers
}
```

#### ⚠️ 向后兼容注意

```go
// 如果需要兼容Go 1.21及更早版本
func backwardCompatible(items []string) {
    var wg sync.WaitGroup
    for _, item := range items {
        item := item  // 显式创建副本（兼容旧版本）
        wg.Add(1)
        go func() {
            defer wg.Done()
            process(item)
        }()
    }
    wg.Wait()
}
```

---

## 🚀 迁移指南

### 从Go 1.20迁移到1.21+

#### 步骤1: 更新Go版本

```bash
# 下载Go 1.21+
# 更新go.mod
go mod edit -go=1.21
```

#### 步骤2: 利用新特性

```go
// 使用min/max
// 旧代码
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 新代码：使用内置函数
result := min(a, b)

// 使用clear
// 旧代码
for k := range m {
    delete(m, k)
}

// 新代码
clear(m)
```

#### 步骤3: 迁移到slog

```go
// 旧代码：log包
log.Printf("user %s logged in", username)

// 新代码：slog
slog.Info("user logged in", "username", username)
```

---

### 从Go 1.21迁移到1.22+

#### 主要变更

1. **for循环语义变化**

    ```go
    // 检查代码中的for循环使用
    // 重点检查：goroutine、闭包、defer中使用循环变量

    // 可能有问题的代码（Go 1.21）
    for i := 0; i < 10; i++ {
        go func() {
            fmt.Println(i)  // Go 1.21: 可能全打印10
                            // Go 1.22: 正确打印0-9
        }()
    }
    ```

2. **HTTP路由迁移**

    ```go
    // 旧代码
    mux := http.NewServeMux()
    mux.HandleFunc("/users/", handleUsers)

    // 新代码：使用方法和路径参数
    mux.HandleFunc("GET /users/{id}", handleUser)
    mux.HandleFunc("POST /users", createUser)
    ```

---

### 从Go 1.22迁移到1.23+

#### 实验性特性尝试

```bash
# 尝试迭代器（实验性）
GOEXPERIMENT=rangefunc go test

# 如果测试通过，可以在开发环境使用
export GOEXPERIMENT=rangefunc
```

---

## ⚠️ 常见问题与陷阱

### 1. 泛型性能问题

```go
// ❌ 问题：泛型可能导致代码膨胀
func GenericProcess[T any](items []T) {
    // 每个类型实例化一次
}

// ✅ 解决：合理使用泛型
// - 对于性能关键路径，考虑使用具体类型
// - 使用基准测试验证性能
```

### 2. slog内存分配

```go
// ❌ 避免：频繁的字符串拼接
slog.Info("message", "data", fmt.Sprintf("%v", largeStruct))

// ✅ 改进：使用LogValuer接口
type User struct {
    ID   int
    Name string
}

func (u User) LogValue() slog.Value {
    return slog.GroupValue(
        slog.Int("id", u.ID),
        slog.String("name", u.Name),
    )
}

slog.Info("user data", "user", user)  // 零分配
```

### 3. for循环边界情况

```go
// ⚠️ 注意：range修改slice
for i, v := range slice {
    if condition {
        slice = append(slice, newItem)  // 可能导致无限循环
    }
}

// ✅ 解决：使用传统for循环
for i := 0; i < len(slice); i++ {
    if condition {
        slice = append(slice, newItem)
        // 不会影响当前迭代
    }
}
```

---

## 🔗 相关资源

### 迁移工具

- [go fix工具](https://pkg.go.dev/cmd/fix)
- [gopls IDE支持](https://github.com/golang/tools/tree/master/gopls)

### 文档

- [Go 1.21迁移指南](https://go.dev/doc/go1.21)
- [Go 1.22迁移指南](https://go.dev/doc/go1.22)
- [Go 1.23迁移指南](https://go.dev/doc/go1.23)

### 社区

- [Go Forum](https://forum.golangbridge.org/)
- [Reddit r/golang](https://www.reddit.com/r/golang/)
- [Go中文社区](https://gocn.vip/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025-10-29  
**文档状态**: 完成  
**适用版本**: Go 1.21+
