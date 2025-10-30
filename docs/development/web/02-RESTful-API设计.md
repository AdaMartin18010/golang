# RESTful API设计

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

> **难度**: ⭐⭐⭐
> **标签**: #RESTful #API设计 #HTTP

## 📋 目录

- [1. REST原则](#1-rest原则)
  - [REST核心概念](#rest核心概念)
  - [REST约束](#rest约束)
- [2. URL设计](#2-url设计)
  - [资源命名](#资源命名)
  - [资源层级](#资源层级)
  - [查询参数](#查询参数)
  - [Go实现示例](#go实现示例)
- [3. HTTP方法](#3-http方法)
  - [标准方法](#标准方法)
  - [方法使用示例](#方法使用示例)
- [4. 状态码](#4-状态码)
  - [常用状态码](#常用状态码)
  - [使用示例](#使用示例)
- [5. 请求和响应](#5-请求和响应)
  - [请求格式](#请求格式)
  - [响应格式](#响应格式)
  - [示例响应](#示例响应)
- [6. 最佳实践](#6-最佳实践)
  - [1. 版本控制](#1-版本控制)
  - [2. 过滤和搜索](#2-过滤和搜索)
  - [3. 分页](#3-分页)
  - [4. 字段过滤](#4-字段过滤)
  - [5. HATEOAS](#5-hateoas)
  - [6. 完整示例](#6-完整示例)
- [🔗 相关资源](#相关资源)

## 1. REST原则

### REST核心概念

- **资源（Resource）**: 网络上的一个实体
- **表现层（Representation）**: 资源的表现形式（JSON/XML）
- **状态转移（State Transfer）**: 通过HTTP方法改变资源状态
- **无状态（Stateless）**: 每个请求独立

---

### REST约束

1. **客户端-服务器**: 分离关注点
2. **无状态**: 每个请求包含完整信息
3. **可缓存**: 响应可缓存
4. **统一接口**: 统一的资源操作方式
5. **分层系统**: 客户端无需知道中间层
6. **按需代码**: 可选的代码下载

---

## 2. URL设计

### 资源命名

```
✅ 推荐：使用名词复数
GET    /users          # 获取用户列表
GET    /users/123      # 获取特定用户
POST   /users          # 创建用户
PUT    /users/123      # 更新用户
DELETE /users/123      # 删除用户

❌ 不推荐：使用动词
GET    /getUsers
POST   /createUser
PUT    /updateUser
DELETE /deleteUser
```

---

### 资源层级

```
✅ 推荐：体现资源关系
GET    /users/123/posts           # 用户的文章
GET    /users/123/posts/456       # 用户的特定文章
POST   /users/123/posts           # 为用户创建文章
GET    /posts/456/comments        # 文章的评论

❌ 不推荐：过深的嵌套（>3层）
GET    /users/123/posts/456/comments/789/likes
```

---

### 查询参数

```
# 过滤
GET /users?role=admin&status=active

# 排序
GET /users?sort=created_at&order=desc

# 分页
GET /users?page=2&limit=20
GET /users?offset=20&limit=20

# 字段选择
GET /users?fields=id,name,email

# 搜索
GET /users?q=john
```

---

### Go实现示例

```go
func listUsers(w http.ResponseWriter, r *http.Request) {
    // 解析查询参数
    query := r.URL.Query()

    // 分页
    page, _ := strconv.Atoi(query.Get("page"))
    if page < 1 {
        page = 1
    }

    limit, _ := strconv.Atoi(query.Get("limit"))
    if limit < 1 || limit > 100 {
        limit = 20
    }

    // 排序
    sort := query.Get("sort")
    if sort == "" {
        sort = "created_at"
    }

    order := query.Get("order")
    if order != "asc" && order != "desc" {
        order = "desc"
    }

    // 过滤
    filters := map[string]string{
        "role":   query.Get("role"),
        "status": query.Get("status"),
    }

    // 查询数据
    users, total, err := getUserList(page, limit, sort, order, filters)
    if err != nil {
        respondError(w, 500, "Internal server error")
        return
    }

    // 响应
    respondJSON(w, 200, map[string]interface{}{
        "data": users,
        "pagination": map[string]int{
            "page":  page,
            "limit": limit,
            "total": total,
        },
    })
}
```

---

## 3. HTTP方法

### 标准方法

| 方法 | 用途 | 幂等性 | 安全性 |
|------|------|--------|--------|
| **GET** | 获取资源 | ✅ | ✅ |
| **POST** | 创建资源 | ❌ | ❌ |
| **PUT** | 完整更新 | ✅ | ❌ |
| **PATCH** | 部分更新 | ❌ | ❌ |
| **DELETE** | 删除资源 | ✅ | ❌ |
| **HEAD** | 获取元数据 | ✅ | ✅ |
| **OPTIONS** | 获取选项 | ✅ | ✅ |

---

### 方法使用示例

```go
// GET - 获取资源
func getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := findUserByID(id)
    if err != nil {
        respondError(w, 404, "User not found")
        return
    }
    respondJSON(w, 200, user)
}

// POST - 创建资源
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        respondError(w, 400, "Invalid request body")
        return
    }

    if err := validate(user); err != nil {
        respondError(w, 422, err.Error())
        return
    }

    createdUser, err := insertUser(user)
    if err != nil {
        respondError(w, 500, "Failed to create user")
        return
    }

    respondJSON(w, 201, createdUser)
}

// PUT - 完整更新
func updateUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        respondError(w, 400, "Invalid request body")
        return
    }

    user.ID = id
    updatedUser, err := replaceUser(user)
    if err != nil {
        respondError(w, 404, "User not found")
        return
    }

    respondJSON(w, 200, updatedUser)
}

// PATCH - 部分更新
func patchUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var updates map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        respondError(w, 400, "Invalid request body")
        return
    }

    updatedUser, err := partialUpdateUser(id, updates)
    if err != nil {
        respondError(w, 404, "User not found")
        return
    }

    respondJSON(w, 200, updatedUser)
}

// DELETE - 删除资源
func deleteUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := removeUser(id); err != nil {
        respondError(w, 404, "User not found")
        return
    }

    respondJSON(w, 204, nil)
}
```

---

## 4. 状态码

### 常用状态码

**2xx 成功**:
- `200 OK`: 请求成功
- `201 Created`: 资源已创建
- `202 Accepted`: 请求已接受（异步处理）
- `204 No Content`: 成功但无返回内容

**3xx 重定向**:
- `301 Moved Permanently`: 永久重定向
- `302 Found`: 临时重定向
- `304 Not Modified`: 资源未修改

**4xx 客户端错误**:
- `400 Bad Request`: 请求格式错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 无权限
- `404 Not Found`: 资源不存在
- `405 Method Not Allowed`: 方法不允许
- `409 Conflict`: 冲突（如重复创建）
- `422 Unprocessable Entity`: 验证失败
- `429 Too Many Requests`: 请求过多

**5xx 服务器错误**:
- `500 Internal Server Error`: 服务器错误
- `502 Bad Gateway`: 网关错误
- `503 Service Unavailable`: 服务不可用
- `504 Gateway Timeout`: 网关超时

---

### 使用示例

```go
func respondError(w http.ResponseWriter, code int, message string) {
    respondJSON(w, code, map[string]string{
        "error": message,
    })
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 验证失败 - 422
    if !isValid(data) {
        respondError(w, 422, "Validation failed")
        return
    }

    // 未找到 - 404
    user, err := findUser(id)
    if err == ErrNotFound {
        respondError(w, 404, "User not found")
        return
    }

    // 服务器错误 - 500
    if err != nil {
        log.Printf("Error: %v", err)
        respondError(w, 500, "Internal server error")
        return
    }

    // 成功 - 200
    respondJSON(w, 200, user)
}
```

---

## 5. 请求和响应

### 请求格式

```go
// 创建用户请求
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

// 更新用户请求
type UpdateUserRequest struct {
    Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty" binding:"omitempty,email"`
    Role  *string `json:"role,omitempty" binding:"omitempty,oneof=user admin"`
}
```

---

### 响应格式

```go
// 统一响应格式
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 列表响应
type ListResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Page  int `json:"page"`
    Limit int `json:"limit"`
    Total int `json:"total"`
}

// 错误响应
type ErrorResponse struct {
    Error   string            `json:"error"`
    Details map[string]string `json:"details,omitempty"`
}
```

---

### 示例响应

```json
// 成功响应
{
  "code": 0,
  "message": "Success",
  "data": {
    "id": 123,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-10-28T10:00:00Z"
  }
}

// 列表响应
{
  "data": [
    {"id": 1, "name": "User 1"},
    {"id": 2, "name": "User 2"}
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}

// 错误响应
{
  "error": "Validation failed",
  "details": {
    "email": "Invalid email format",
    "password": "Password too short"
  }
}
```

---

## 6. 最佳实践

### 1. 版本控制

```
# URL版本
GET /api/v1/users
GET /api/v2/users

# Header版本
GET /api/users
Accept: application/vnd.myapp.v1+json
```

**Go实现**:
```go
r := chi.NewRouter()

// v1路由
r.Route("/api/v1", func(r chi.Router) {
    r.Get("/users", v1.ListUsers)
    r.Post("/users", v1.CreateUser)
})

// v2路由
r.Route("/api/v2", func(r chi.Router) {
    r.Get("/users", v2.ListUsers)
    r.Post("/users", v2.CreateUser)
})
```

---

### 2. 过滤和搜索

```go
// 支持多种过滤方式
GET /users?status=active
GET /users?role=admin&status=active
GET /users?created_after=2025-01-01
GET /users?q=john  # 全文搜索
```

---

### 3. 分页

```go
// 页码分页
GET /users?page=2&limit=20

// 偏移分页
GET /users?offset=20&limit=20

// 游标分页（大数据集）
GET /users?cursor=eyJpZCI6MTAwfQ==&limit=20
```

---

### 4. 字段过滤

```go
// 只返回需要的字段
GET /users?fields=id,name,email

// 排除字段
GET /users?exclude=password,salt
```

---

### 5. HATEOAS

```json
{
  "id": 123,
  "name": "John Doe",
  "email": "john@example.com",
  "_links": {
    "self": {"href": "/users/123"},
    "posts": {"href": "/users/123/posts"},
    "followers": {"href": "/users/123/followers"}
  }
}
```

---

### 6. 完整示例

```go
package api

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
)

type UserAPI struct {
    service *UserService
}

func NewUserAPI(service *UserService) *UserAPI {
    return &UserAPI{service: service}
}

func (api *UserAPI) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/", api.list)
    r.Post("/", api.create)
    r.Get("/{id}", api.get)
    r.Put("/{id}", api.update)
    r.Patch("/{id}", api.patch)
    r.Delete("/{id}", api.delete)

    return r
}

func (api *UserAPI) list(w http.ResponseWriter, r *http.Request) {
    params := ParseQueryParams(r)
    users, total, err := api.service.List(params)

    if err != nil {
        RespondError(w, 500, "Failed to fetch users")
        return
    }

    RespondJSON(w, 200, map[string]interface{}{
        "data": users,
        "pagination": map[string]int{
            "page":  params.Page,
            "limit": params.Limit,
            "total": total,
        },
    })
}

func (api *UserAPI) create(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        RespondError(w, 400, "Invalid request body")
        return
    }

    if err := Validate(req); err != nil {
        RespondValidationError(w, err)
        return
    }

    user, err := api.service.Create(req)
    if err != nil {
        RespondError(w, 500, "Failed to create user")
        return
    }

    w.Header().Set("Location", fmt.Sprintf("/users/%s", user.ID))
    RespondJSON(w, 201, user)
}
```

---

## 🔗 相关资源

- [HTTP协议](./01-HTTP协议.md)
- [Gin框架](./05-Gin框架.md)
- [API文档化](./07-API文档化.md)

---

**最后更新**: 2025-10-29
**Go版本**: 1.25.3
