# 请求验证框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26.2

---

## 📋 目录

- [请求验证框架](#请求验证框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 验证规则](#2-验证规则)
    - [2.1 支持的规则](#21-支持的规则)
  - [3. 使用示例](#3-使用示例)
    - [3.1 结构体验证](#31-结构体验证)
    - [3.2 字段验证](#32-字段验证)
    - [3.3 工具函数](#33-工具函数)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

请求验证框架提供了统一的参数验证机制：

- ✅ **结构体验证**: 基于结构体标签的验证
- ✅ **多种验证规则**: 支持 required、min、max、email、url 等
- ✅ **自定义验证器**: 支持自定义验证函数
- ✅ **错误格式化**: 统一的验证错误格式
- ✅ **字段验证**: 支持单个字段验证

---

## 2. 验证规则

### 2.1 支持的规则

- `required` - 必填
- `min=n` - 最小值/最小长度
- `max=n` - 最大值/最大长度
- `len=n` - 固定长度
- `email` - 邮箱格式
- `url` - URL 格式
- `uuid` - UUID 格式
- `datetime` - 日期时间格式
- `regexp=pattern` - 正则表达式
- `in=value1|value2` - 值在指定列表中

---

## 3. 使用示例

### 3.1 结构体验证

```go
type CreateUserRequest struct {
    Name     string `validate:"required,min=2,max=50"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"required,min=18,max=100"`
    Password string `validate:"required,min=8"`
    Role     string `validate:"required,in=admin|user|guest"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest,
            errors.NewInvalidInputError("invalid request body"))
        return
    }

    validator := validator.NewValidator()
    if !validator.Validate(req) {
        validationErrors := validator.ValidateStruct(req)
        details := make(map[string]interface{})
        for _, err := range validationErrors {
            details[err.Field] = err.Message
        }
        response.Error(w, http.StatusBadRequest,
            errors.NewValidationError("validation failed", details))
        return
    }

    // 处理请求...
}
```

### 3.2 字段验证

```go
validator := validator.NewValidator()
if !validator.ValidateField(email, "required,email") {
    // 处理验证错误
}
```

### 3.3 工具函数

```go
// 验证邮箱
if !validator.IsEmail(email) {
    return errors.NewValidationError("email", "invalid email format")
}

// 验证URL
if !validator.IsURL(url) {
    return errors.NewValidationError("url", "invalid URL format")
}

// 验证UUID
if !validator.IsUUID(uuid) {
    return errors.NewValidationError("uuid", "invalid UUID format")
}

// 验证必填
if !validator.IsRequired(value) {
    return errors.NewValidationError("field", "field is required")
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用结构体标签**: 在结构体定义中使用 validate 标签
2. **组合规则**: 使用逗号分隔多个验证规则
3. **错误处理**: 使用验证错误创建统一的错误响应
4. **字段验证**: 对于简单场景使用字段验证函数
5. **自定义验证**: 对于复杂验证逻辑使用自定义验证器

### 4.2 DON'Ts ❌

1. **不要忽略验证**: 始终验证用户输入
2. **不要过度验证**: 只验证必要的字段
3. **不要暴露内部错误**: 验证错误应该对用户友好
4. **不要重复验证**: 避免在多个地方验证相同的数据

---

## 5. 相关资源

- [统一错误处理框架](../errors/README.md)
- [统一响应格式框架](../http/response/README.md)
- 框架拓展计划

---

**更新日期**: 2025-11-11
