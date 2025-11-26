# 1. ⚙️ Slog 日志库深度解析

> **简介**: 本文档详细阐述了 Slog 日志库的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. ⚙️ Slog 日志库深度解析](#1-️-slog-日志库深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 基础日志使用](#131-基础日志使用)
    - [1.3.2 结构化日志](#132-结构化日志)
    - [1.3.3 日志上下文](#133-日志上下文)
    - [1.3.4 自定义 Handler](#134-自定义-handler)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 日志级别选择最佳实践](#141-日志级别选择最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Slog 是什么？**

Slog 是 Go 1.21+ 引入的标准库结构化日志包。

**核心特性**:

- ✅ **结构化日志**: 支持结构化日志
- ✅ **标准库**: 标准库，无需第三方依赖
- ✅ **Handler**: 可定制的 Handler
- ✅ **性能**: 性能优化

---

## 1.2 选型论证

**为什么选择 Slog？**

**论证矩阵**:

| 评估维度 | 权重 | Slog | Logrus | Zap | Zerolog | 说明 |
|---------|------|------|--------|-----|---------|------|
| **标准库** | 40% | 10 | 0 | 0 | 0 | Slog 是标准库 |
| **结构化日志** | 25% | 10 | 9 | 10 | 10 | Slog 支持结构化 |
| **性能** | 20% | 9 | 6 | 10 | 10 | Slog 性能优秀 |
| **易用性** | 10% | 9 | 10 | 7 | 8 | Slog API 简单 |
| **生态兼容** | 5% | 8 | 9 | 8 | 7 | Slog 兼容性好 |
| **加权总分** | - | **9.50** | 5.85 | 6.90 | 6.80 | Slog 得分最高 |

**核心优势**:

1. **标准库（权重 40%）**:
   - Go 1.21+ 标准库，稳定可靠
   - 无需第三方依赖，减少依赖风险
   - 未来 Go 日志标准，长期支持

2. **结构化日志（权重 25%）**:
   - 原生支持结构化日志
   - 支持键值对和属性
   - 与 OpenTelemetry 集成良好

3. **性能（权重 20%）**:
   - 性能优秀，零分配设计
   - 支持 Handler 定制
   - 适合生产环境

**为什么不选择其他日志库？**

1. **Logrus**:
   - ✅ 功能丰富，生态成熟
   - ❌ 非标准库，需要第三方依赖
   - ❌ 性能不如 Slog
   - ❌ 维护状态不确定

2. **Zap**:
   - ✅ 性能优秀，结构化日志支持好
   - ❌ 非标准库，需要第三方依赖
   - ❌ API 较复杂，学习成本高
   - ❌ 与标准库不兼容

3. **Zerolog**:
   - ✅ 性能优秀，API 简洁
   - ❌ 非标准库，需要第三方依赖
   - ❌ 生态不如 Slog 丰富
   - ❌ 与标准库不兼容

---

## 1.3 实际应用

### 1.3.1 基础日志使用

**基础日志示例**:

```go
// internal/infrastructure/logging/logger.go
package logging

import (
    "log/slog"
    "os"
)

func InitLogger(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{
        Level: logLevel,
    }

    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler)

    slog.SetDefault(logger)
    return logger
}
```

### 1.3.2 结构化日志

**结构化日志示例**:

```go
// 结构化日志
logger.Info("User created",
    "user_id", user.ID,
    "email", user.Email,
    "name", user.Name,
)

logger.Error("Failed to create user",
    "error", err,
    "email", req.Email,
)
```

### 1.3.3 日志上下文

**日志上下文示例**:

```go
// 使用上下文
logger := slog.Default().With(
    "request_id", requestID,
    "user_id", userID,
)

logger.InfoContext(ctx, "Processing request",
    "path", r.URL.Path,
    "method", r.Method,
)
```

### 1.3.4 自定义 Handler

**自定义 Handler 示例**:

```go
// 自定义 Handler
type CustomHandler struct {
    handler slog.Handler
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
    // 添加自定义字段
    r.AddAttrs(
        slog.String("service", "golang-service"),
        slog.String("version", "1.0.0"),
    )

    return h.handler.Handle(ctx, r)
}

// 使用自定义 Handler
handler := &CustomHandler{
    handler: slog.NewJSONHandler(os.Stdout, nil),
}
logger := slog.New(handler)
```

---

## 1.4 最佳实践

### 1.4.1 日志级别选择最佳实践

**为什么需要合理选择日志级别？**

合理选择日志级别可以提高日志的可读性和可维护性。

**日志级别选择原则**:

1. **DEBUG**: 详细的调试信息
2. **INFO**: 一般信息，如请求处理
3. **WARN**: 警告信息，如配置问题
4. **ERROR**: 错误信息，如处理失败

**实际应用示例**:

```go
// 日志级别选择最佳实践
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    logger := slog.Default().With(
        "method", r.Method,
        "path", r.URL.Path,
    )

    // DEBUG: 详细的调试信息
    logger.DebugContext(r.Context(), "Creating user",
        "email", req.Email,
        "name", req.Name,
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    if err != nil {
        // ERROR: 错误信息
        logger.ErrorContext(r.Context(), "Failed to create user",
            "error", err,
            "email", req.Email,
        )
        Error(w, http.StatusInternalServerError, err)
        return
    }

    // INFO: 一般信息
    logger.InfoContext(r.Context(), "User created successfully",
        "user_id", user.ID,
    )
    Success(w, http.StatusCreated, user)
}
```

**最佳实践要点**:

1. **DEBUG**: 用于详细的调试信息，生产环境通常关闭
2. **INFO**: 用于一般信息，如请求处理、状态变更
3. **WARN**: 用于警告信息，如配置问题、性能问题
4. **ERROR**: 用于错误信息，如处理失败、异常情况

---

## 📚 扩展阅读

- [Slog 官方文档](https://pkg.go.dev/log/slog)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Slog 日志库的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
