# 统一响应格式框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [统一响应格式框架](#统一响应格式框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 响应结构](#2-响应结构)
    - [2.1 成功响应](#21-成功响应)
    - [2.2 错误响应](#22-错误响应)
    - [2.3 分页响应](#23-分页响应)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 带追踪ID](#32-带追踪id)
    - [3.3 分页响应](#33-分页响应)
    - [3.4 带元数据](#34-带元数据)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

统一响应格式框架提供了标准化的 API 响应格式，包括：

- ✅ **标准响应结构**: 统一的成功和错误响应格式
- ✅ **分页支持**: 内置分页响应支持
- ✅ **元数据支持**: 支持添加请求ID、版本等元数据
- ✅ **追踪支持**: 支持添加追踪ID
- ✅ **错误集成**: 与错误处理框架无缝集成

---

## 2. 响应结构

### 2.1 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "123",
    "name": "John"
  },
  "timestamp": "2025-11-11T10:00:00Z",
  "trace_id": "trace-123",
  "meta": {
    "request_id": "req-123",
    "version": "v1.0.0"
  }
}
```

### 2.2 错误响应

```json
{
  "code": 404,
  "message": "error",
  "error": {
    "code": "NOT_FOUND",
    "message": "user with id 123 not found",
    "details": {
      "resource": "user",
      "id": "123"
    }
  },
  "timestamp": "2025-11-11T10:00:00Z",
  "trace_id": "trace-123"
}
```

### 2.3 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {"id": "1", "name": "Item 1"},
    {"id": "2", "name": "Item 2"}
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3
  },
  "timestamp": "2025-11-11T10:00:00Z"
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "net/http"
    "github.com/yourusername/golang/pkg/http/response"
)

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        response.Error(w, http.StatusNotFound, err)
        return
    }

    response.Success(w, http.StatusOK, user)
}
```

### 3.2 带追踪ID

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    traceID := getTraceID(r)

    user, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        response.ErrorWithTraceID(w, http.StatusNotFound, err, traceID)
        return
    }

    response.SuccessWithTraceID(w, http.StatusOK, user, traceID)
}
```

### 3.3 分页响应

```go
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    page := getPage(r)
    pageSize := getPageSize(r)

    users, total, err := h.userService.ListUsers(r.Context(), page, pageSize)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err)
        return
    }

    response.Paginated(w, http.StatusOK, users, page, pageSize, total)
}
```

### 3.4 带元数据

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        response.Error(w, http.StatusNotFound, err)
        return
    }

    meta := response.NewMeta(getRequestID(r), "v1.0.0")
    meta.WithExtra("server", "api-01")

    response.SuccessWithMeta(w, http.StatusOK, user, meta)
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **统一使用**: 所有API响应都使用统一格式
2. **添加追踪ID**: 在生产环境中添加追踪ID
3. **使用分页**: 列表接口使用分页响应
4. **错误集成**: 使用错误处理框架的错误类型
5. **元数据**: 添加有用的元数据信息

### 4.2 DON'Ts ❌

1. **不要混用格式**: 不要使用不同的响应格式
2. **不要暴露敏感信息**: 错误响应中不要暴露敏感信息
3. **不要忽略错误**: 始终处理错误并返回错误响应
4. **不要过度使用元数据**: 只添加必要的元数据

---

## 5. 相关资源

- [统一错误处理框架](../errors/README.md)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
