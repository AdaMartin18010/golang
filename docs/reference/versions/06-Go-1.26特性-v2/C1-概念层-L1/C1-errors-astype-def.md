# errors.AsType (错误类型断言增强)

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **包路径**: `errors`
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 功能概述

```
errors.AsType: Go 1.26增强的错误类型断言函数

功能:
  - 安全的错误类型转换
  - 支持泛型，避免类型断言的interface{}
  - 编译期类型检查

对比 errors.As:
  As:   errors.As(err, &target)  // target必须是指针
  AsType: errors.AsType[MyError](err)  // 泛型版本
```

### 1.2 语法形式

```go
// Go 1.26新API
func AsType[T error](err error) (T, bool)

// 使用示例
if myErr, ok := errors.AsType[*MyError](err); ok {
    // myErr的类型是*MyError，无需类型断言
    fmt.Println(myErr.Code)
}
```

---

## 二、核心优势

### 2.1 类型安全

```go
// ❌ Go 1.25: 运行时类型断言
var myErr *MyError
if errors.As(err, &myErr) {
    // 可能需要nil检查
    fmt.Println(myErr.Code)
}

// ✅ Go 1.26: 编译期类型检查
if myErr, ok := errors.AsType[*MyError](err); ok {
    // myErr保证非nil且类型正确
    fmt.Println(myErr.Code)
}
```

### 2.2 简洁性

```go
// 传统方式
var connErr *net.OpError
if errors.As(err, &connErr) {
    // ...
}

// Go 1.26方式
if connErr, ok := errors.AsType[*net.OpError](err); ok {
    // ...
}
```

---

## 三、使用场景

### 3.1 错误处理链

```go
func handleError(err error) {
    // 尝试匹配特定错误类型
    if validationErr, ok := errors.AsType[*ValidationError](err); ok {
        log.Printf("Validation failed: %v", validationErr.Fields)
        return
    }

    if networkErr, ok := errors.AsType[*NetworkError](err); ok {
        log.Printf("Network error: %v", networkErr.Endpoint)
        return
    }

    // 默认处理
    log.Printf("Unknown error: %v", err)
}
```

### 3.2 与errors.Join配合

```go
err := errors.Join(err1, err2, err3)

// 检查是否包含特定类型错误
if timeoutErr, ok := errors.AsType[*TimeoutError](err); ok {
    // 处理超时
}
```

---

**概念分类**: 标准库 - 错误处理
**Go版本**: 1.26+
**包路径**: `errors`
