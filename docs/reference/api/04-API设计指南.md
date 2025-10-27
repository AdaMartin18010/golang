# API设计指南

**难度**: 中级 | **预计阅读**: 15分钟

---

## 📋 目录

- [1. 📖 RESTful设计原则](#1--restful设计原则)
- [2. 📚 相关资源](#2--相关资源)

---

## 1. 📖 RESTful设计原则

### 资源命名
```
✅ 好的设计:
GET    /users          # 列表
GET    /users/123      # 详情
POST   /users          # 创建
PUT    /users/123      # 更新
DELETE /users/123      # 删除

❌ 避免:
GET /getUsers
POST /createUser
```

---

### HTTP状态码
```
200 OK              - 成功
201 Created         - 创建成功
204 No Content      - 删除成功
400 Bad Request     - 请求错误
401 Unauthorized    - 未认证
403 Forbidden       - 无权限
404 Not Found       - 未找到
500 Internal Error  - 服务器错误
```

---

## 🎯 API版本控制

```go
// URL路径版本
r.Group("/api/v1")
r.Group("/api/v2")

// Header版本
if r.Header.Get("API-Version") == "2" {
    // v2处理
}
```

---

## 📝 响应格式

```go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 成功响应
c.JSON(200, APIResponse{
    Code:    0,
    Message: "success",
    Data:    users,
})

// 错误响应
c.JSON(400, APIResponse{
    Code:    1001,
    Message: "Invalid input",
})
```

---

## 🔐 认证与授权

```go
// JWT认证
Authorization: Bearer <token>

// API Key
X-API-Key: <key>

// Basic Auth
Authorization: Basic <base64(username:password)>
```

---

## 📊 分页

```go
type PaginatedResponse struct {
    Data       []interface{} `json:"data"`
    Page       int           `json:"page"`
    PerPage    int           `json:"per_page"`
    Total      int64         `json:"total"`
    TotalPages int           `json:"total_pages"`
}

// 查询参数
GET /users?page=1&per_page=20
```

---

## 🔍 过滤与排序

```go
// 过滤
GET /users?status=active&role=admin

// 排序
GET /users?sort=created_at&order=desc

// 搜索
GET /users?q=john
```

---

## 📚 相关资源

- [REST API Guidelines](https://restfulapi.net/)

**下一步**: [guides/01-学习路线图](../guides/01-学习路线图.md)

---

**最后更新**: 2025-10-28

