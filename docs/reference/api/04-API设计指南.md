# API设计指南

**难度**: 中级 | **预计阅读**: 20分钟

---

## 📋 目录


- [1. RESTful API设计](#1-restful-api设计)
  - [资源命名](#资源命名)
  - [HTTP方法](#http方法)
  - [状态码](#状态码)
- [2. 错误处理](#2-错误处理)
  - [统一错误响应](#统一错误响应)
  - [实现](#实现)
- [3. 版本控制](#3-版本控制)
  - [URL版本](#url版本)
  - [Header版本](#header版本)
- [4. 安全性](#4-安全性)
  - [认证](#认证)
  - [限流](#限流)
  - [CORS](#cors)
- [5. 最佳实践](#5-最佳实践)
  - [分页](#分页)
  - [过滤和排序](#过滤和排序)
  - [请求验证](#请求验证)
- [📚 相关资源](#-相关资源)
- [🔗 导航](#-导航)

## 1. RESTful API设计

### 资源命名

```
✅ 好的命名
GET    /api/users          # 获取用户列表
GET    /api/users/:id      # 获取单个用户
POST   /api/users          # 创建用户
PUT    /api/users/:id      # 更新用户
DELETE /api/users/:id      # 删除用户

❌ 不好的命名
GET    /api/getUsers
POST   /api/createUser
GET    /api/user_list
```

### HTTP方法

| 方法 | 用途 | 幂等性 |
|------|------|--------|
| GET | 获取资源 | ✅ |
| POST | 创建资源 | ❌ |
| PUT | 完整更新 | ✅ |
| PATCH | 部分更新 | ❌ |
| DELETE | 删除资源 | ✅ |

### 状态码

```go
// 成功
200 OK           // 请求成功
201 Created      // 创建成功
204 No Content   // 删除成功

// 客户端错误
400 Bad Request       // 请求参数错误
401 Unauthorized      // 未认证
403 Forbidden         // 无权限
404 Not Found         // 资源不存在
409 Conflict          // 资源冲突

// 服务器错误
500 Internal Server Error  // 服务器内部错误
503 Service Unavailable    // 服务不可用
```

---

## 2. 错误处理

### 统一错误响应

```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

// 示例
{
    "code": "INVALID_PARAMETER",
    "message": "邮箱格式不正确",
    "details": {
        "field": "email",
        "value": "invalid-email"
    }
}
```

### 实现

```go
func HandleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *ValidationError:
        c.JSON(400, ErrorResponse{
            Code:    "VALIDATION_ERROR",
            Message: e.Error(),
            Details: e.Fields,
        })
    case *NotFoundError:
        c.JSON(404, ErrorResponse{
            Code:    "NOT_FOUND",
            Message: e.Error(),
        })
    default:
        c.JSON(500, ErrorResponse{
            Code:    "INTERNAL_ERROR",
            Message: "服务器内部错误",
        })
    }
}
```

---

## 3. 版本控制

### URL版本

```go
// 推荐
/api/v1/users
/api/v2/users

r := gin.Default()
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsersV1)
}
v2 := r.Group("/api/v2")
{
    v2.GET("/users", getUsersV2)
}
```

### Header版本

```go
// Accept: application/vnd.api+json; version=1
func VersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("API-Version")
        if version == "" {
            version = "1"
        }
        c.Set("api_version", version)
        c.Next()
    }
}
```

---

## 4. 安全性

### 认证

```go
// JWT认证
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()
            return
        }
        
        claims, err := ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "令牌无效"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### 限流

```go
import "golang.org/x/time/rate"

func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
    limiter := rate.NewLimiter(r, b)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### CORS

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

---

## 5. 最佳实践

### 分页

```go
type PaginationParams struct {
    Page     int `form:"page" binding:"min=1"`
    PageSize int `form:"page_size" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
}

func GetUsers(c *gin.Context) {
    var params PaginationParams
    params.Page = 1
    params.PageSize = 20
    c.ShouldBindQuery(&params)
    
    // 查询数据...
}
```

### 过滤和排序

```
GET /api/users?status=active&sort=-created_at&fields=id,name,email
```

```go
type QueryParams struct {
    Status string   `form:"status"`
    Sort   string   `form:"sort"`
    Fields []string `form:"fields"`
}
```

### 请求验证

```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=0,lte=150"`
    Password string `json:"password" binding:"required,min=8"`
}

func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 处理...
}
```

---

## 📚 相关资源

- [RESTful API Design Best Practices](https://restfulapi.net/)
- [HTTP Status Codes](https://httpstatuses.com/)
- [API Security Checklist](https://github.com/shieldfy/API-Security-Checklist)

---

## 🔗 导航

- **上一页**: [03-常用第三方库](./03-常用第三方库.md)
- **相关**: [README](./README.md)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3
