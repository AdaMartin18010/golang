# 统一错误处理框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [统一错误处理框架](#统一错误处理框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
  - [3. 错误代码](#3-错误代码)
  - [4. 错误分类](#4-错误分类)
  - [5. 使用示例](#5-使用示例)
  - [6. 最佳实践](#6-最佳实践)

---

## 1. 概述

统一错误处理框架提供了标准化的错误处理机制，包括：

- ✅ **错误代码体系**: 统一的错误代码定义
- ✅ **错误分类**: 客户端错误、服务端错误、业务错误
- ✅ **HTTP状态码映射**: 自动映射到HTTP状态码
- ✅ **详细信息支持**: 支持添加详细的错误信息
- ✅ **追踪支持**: 支持添加追踪ID
- ✅ **可重试标记**: 标记错误是否可重试

---

## 2. 核心功能

### 2.1 AppError 结构

```go
type AppError struct {
    Code       ErrorCode              // 错误代码
    Message    string                 // 错误消息
    Cause      error                  // 底层错误
    Category   ErrorCategory          // 错误分类
    HTTPStatus int                    // HTTP状态码
    Details    map[string]interface{} // 详细信息
    Timestamp  time.Time              // 时间戳
    Retryable  bool                   // 是否可重试
    TraceID    string                 // 追踪ID
}
```

### 2.2 错误创建函数

- `NewNotFoundError(resource, id)` - 资源不存在
- `NewInvalidInputError(message)` - 无效输入
- `NewValidationError(message, details)` - 验证失败
- `NewInternalError(message, cause)` - 内部错误
- `NewUnauthorizedError(message)` - 未授权
- `NewForbiddenError(message)` - 禁止访问
- `NewConflictError(message)` - 资源冲突
- `NewTimeoutError(message)` - 超时
- `NewRateLimitError(message)` - 限流
- `NewServiceUnavailableError(message)` - 服务不可用

---

## 3. 错误代码

### 3.1 客户端错误 (4xx)

- `NOT_FOUND` - 资源不存在 (404)
- `INVALID_INPUT` - 无效输入 (400)
- `VALIDATION_ERROR` - 验证失败 (400)
- `UNAUTHORIZED` - 未授权 (401)
- `FORBIDDEN` - 禁止访问 (403)
- `CONFLICT` - 资源冲突 (409)
- `RATE_LIMIT_EXCEEDED` - 限流 (429)

### 3.2 服务端错误 (5xx)

- `INTERNAL_ERROR` - 内部错误 (500)
- `TIMEOUT` - 超时 (504)
- `SERVICE_UNAVAILABLE` - 服务不可用 (503)

---

## 4. 错误分类

- `CLIENT_ERROR` - 客户端错误 (4xx)
- `SERVER_ERROR` - 服务端错误 (5xx)
- `BUSINESS_ERROR` - 业务错误 (4xx)

---

## 5. 使用示例

### 5.1 基本使用

```go
import "github.com/yourusername/golang/pkg/errors"

// 创建错误
err := errors.NewNotFoundError("user", "123")

// 添加详细信息
err.WithDetails("field", "value")

// 添加追踪ID
err.WithTraceID("trace-123")

// 检查是否可重试
if err.IsRetryable() {
    // 重试逻辑
}

// 获取HTTP状态码
statusCode := err.HTTPStatusCode()
```

### 5.2 在Handler中使用

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    user, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        appErr := errors.FromDomainError(err)
        appErr.WithTraceID(getTraceID(r))
        response.Error(w, appErr.HTTPStatusCode(), appErr)
        return
    }

    response.Success(w, http.StatusOK, user)
}
```

### 5.3 验证错误

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest,
            errors.NewInvalidInputError("invalid request body"))
        return
    }

    // 验证
    if err := validateUser(req); err != nil {
        details := map[string]interface{}{
            "field": err.Field,
            "reason": err.Reason,
        }
        response.Error(w, http.StatusBadRequest,
            errors.NewValidationError("validation failed", details))
        return
    }

    // ...
}
```

---

## 6. 最佳实践

### 6.1 DO's ✅

1. **使用标准错误代码**: 使用预定义的错误代码
2. **添加上下文信息**: 使用 `WithDetails()` 添加详细信息
3. **添加追踪ID**: 在生产环境中添加追踪ID
4. **错误转换**: 使用 `FromDomainError()` 转换领域错误
5. **错误日志**: 记录错误日志，包含完整上下文

### 6.2 DON'Ts ❌

1. **不要暴露内部错误**: 不要直接返回底层错误给客户端
2. **不要忽略错误**: 始终处理错误
3. **不要使用panic**: 使用错误返回值而不是panic
4. **不要创建重复的错误代码**: 使用现有的错误代码

---

## 7. 相关资源

- [错误处理最佳实践](../docs/practices/engineering/05-错误处理最佳实践.md)
- [框架拓展计划](../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
