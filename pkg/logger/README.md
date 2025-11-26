# 结构化日志框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [结构化日志框架](#结构化日志框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
    - [2.1 日志配置](#21-日志配置)
    - [2.2 日志级别](#22-日志级别)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 带上下文](#32-带上下文)
    - [3.3 带字段](#33-带字段)
    - [3.4 记录请求](#34-记录请求)
    - [3.5 记录错误](#35-记录错误)
    - [3.6 自定义配置](#36-自定义配置)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

结构化日志框架基于 Go 1.21+ 的 `slog` 包，提供了增强的日志记录功能：

- ✅ **结构化日志**: 基于 slog 的结构化日志记录
- ✅ **上下文支持**: 自动提取 TraceID、SpanID、UserID
- ✅ **多种输出格式**: 支持 JSON 和文本格式
- ✅ **日志级别**: 支持 Debug、Info、Warn、Error
- ✅ **请求日志**: 内置 HTTP 请求日志记录
- ✅ **错误日志**: 内置错误日志记录

---

## 2. 核心功能

### 2.1 日志配置

```go
type Config struct {
    Level      slog.Level  // 日志级别
    Output     io.Writer   // 输出目标
    AddSource  bool        // 是否添加源码位置
    JSONFormat bool        // 是否使用 JSON 格式
}
```

### 2.2 日志级别

- `slog.LevelDebug` - 调试日志
- `slog.LevelInfo` - 信息日志
- `slog.LevelWarn` - 警告日志
- `slog.LevelError` - 错误日志

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/logger"

// 创建日志记录器
log := logger.NewLogger(slog.LevelInfo)

// 记录日志
log.Info("Application started")
log.Error("Operation failed", "error", err)
```

### 3.2 带上下文

```go
ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
log.WithContext(ctx).Info("Processing request")
```

### 3.3 带字段

```go
log.WithFields("user_id", "123", "action", "login").Info("User logged in")
```

### 3.4 记录请求

```go
log.LogRequest(ctx, "GET", "/users", 200, 100*time.Millisecond)
```

### 3.5 记录错误

```go
log.LogError(ctx, err, "Failed to process request", "user_id", "123")
```

### 3.6 自定义配置

```go
config := logger.Config{
    Level:      slog.LevelDebug,
    Output:     os.Stdout,
    AddSource:  true,
    JSONFormat: true,
}
log := logger.NewLoggerWithConfig(config)
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用结构化日志**: 使用字段而不是字符串拼接
2. **添加上下文**: 使用 `WithContext()` 添加追踪信息
3. **适当的日志级别**: 根据重要性选择合适的日志级别
4. **记录关键操作**: 记录重要的业务操作和错误
5. **性能考虑**: 避免在高频路径中记录过多日志

### 4.2 DON'Ts ❌

1. **不要记录敏感信息**: 不要记录密码、令牌等敏感信息
2. **不要过度记录**: 避免记录过多不必要的日志
3. **不要使用 panic**: 使用错误日志而不是 panic
4. **不要阻塞**: 日志记录不应该阻塞主流程

---

## 5. 相关资源

- [Go slog 文档](https://pkg.go.dev/log/slog)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
