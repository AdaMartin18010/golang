# 17.2 RESTful API项目

<!-- TOC START -->
- [17.2 RESTful API项目](#172-restful-api项目)
  - [17.2.1 📚 项目概述](#1721--项目概述)
    - [17.2.1.1 项目目标](#17211-项目目标)
    - [17.2.1.2 技术栈](#17212-技术栈)
    - [17.2.1.3 架构设计](#17213-架构设计)
  - [17.2.2 💻 核心实现](#1722--核心实现)
    - [17.2.2.1 项目结构](#17221-项目结构)
    - [17.2.2.2 数据模型](#17222-数据模型)
    - [17.2.2.3 路由设计](#17223-路由设计)
    - [17.2.2.4 中间件实现](#17224-中间件实现)
  - [17.2.3 🎯 功能特性](#1723--功能特性)
    - [17.2.3.1 用户管理](#17231-用户管理)
    - [17.2.3.2 认证授权](#17232-认证授权)
    - [17.2.3.3 数据验证](#17233-数据验证)
    - [17.2.3.4 错误处理](#17234-错误处理)
  - [17.2.4 🧪 测试实现](#1724--测试实现)
    - [17.2.4.1 单元测试](#17241-单元测试)
    - [17.2.4.2 集成测试](#17242-集成测试)
    - [17.2.4.3 性能测试](#17243-性能测试)
  - [17.2.5 🚀 部署运维](#1725--部署运维)
    - [17.2.5.1 Docker部署](#17251-docker部署)
    - [17.2.5.2 监控日志](#17252-监控日志)
    - [17.2.5.3 性能优化](#17253-性能优化)
  - [17.2.6 📚 扩展阅读](#1726--扩展阅读)
<!-- TOC END -->

## 17.2.1 📚 项目概述

### 17.2.1.1 项目目标

构建一个完整的RESTful API项目，包含：

- **完整的CRUD操作**: 用户、文章、评论等实体管理
- **认证授权系统**: JWT token认证，角色权限控制
- **数据验证**: 请求参数验证，数据完整性保证
- **错误处理**: 统一的错误响应格式
- **测试覆盖**: 完整的单元测试和集成测试
- **文档生成**: 自动生成API文档

### 17.2.1.2 技术栈

| 技术类别 | 选择方案 | 说明 |
|----------|----------|------|
| Web框架 | Gin | 高性能HTTP框架 |
| 数据库 | PostgreSQL | 关系型数据库 |
| ORM | GORM | Go语言ORM库 |
| 认证 | JWT | JSON Web Token |
| 验证 | validator | 数据验证库 |
| 测试 | testify | 测试框架 |
| 文档 | Swagger | API文档生成 |

### 17.2.1.3 架构设计

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   客户端应用    │    │   API Gateway   │    │   微服务集群    │
│                 │    │                 │    │                 │
│  Web/Mobile/CLI │───▶│  路由/认证/限流  │───▶│  User/Post/...  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐              │
                       │   数据存储层    │◀─────────────┘
                       │                 │
                       │ PostgreSQL/Redis│
                       └─────────────────┘
```

## 17.2.2 💻 核心实现

### 17.2.2.1 项目结构

```text
restful-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── auth.go
│   ├── service/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── auth.go
│   ├── repository/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── database.go
│   ├── model/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── response.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── utils/
│       ├── jwt.go
│       ├── validator.go
│       └── response.go
├── pkg/
│   └── database/
│       └── postgres.go
├── api/
│   └── docs/
├── scripts/
│   ├── migrate.go
│   └── seed.go
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

### 17.2.2.2 数据模型

```go
// internal/model/user.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=20"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
    Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
    Role      string         `json:"role" gorm:"default:user" validate:"oneof=user admin"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`
    Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

type Post struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
    Content   string         `json:"content" gorm:"type:text" validate:"required,min=1"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    Status    string         `json:"status" gorm:"default:draft" validate:"oneof=draft published archived"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

type Comment struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Content   string         `json:"content" gorm:"type:text" validate:"required,min=1,max=1000"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    PostID    uint           `json:"post_id" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Post Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
}

// 请求和响应模型
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
    Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=20"`
    Email    *string `json:"email,omitempty" validate:"omitempty,email"`
    Role     *string `json:"role,omitempty" validate:"omitempty,oneof=user admin"`
    IsActive *bool   `json:"is_active,omitempty"`
}

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

### 17.2.2.3 路由设计

```go
// internal/handler/routes.go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "restful-api/internal/middleware"
)

func SetupRoutes(r *gin.Engine, userHandler *UserHandler, postHandler *PostHandler, authHandler *AuthHandler) {
    // 中间件
    r.Use(cors.Default())
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    
    // API版本分组
    v1 := r.Group("/api/v1")
    {
        // 认证路由
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
        }
        
        // 需要认证的路由
        protected := v1.Group("/")
        protected.Use(middleware.AuthRequired())
        {
            // 用户管理
            users := protected.Group("/users")
            {
                users.GET("", userHandler.GetUsers)
                users.GET("/:id", userHandler.GetUser)
                users.PUT("/:id", userHandler.UpdateUser)
                users.DELETE("/:id", userHandler.DeleteUser)
                users.GET("/profile", userHandler.GetProfile)
                users.PUT("/profile", userHandler.UpdateProfile)
            }
            
            // 文章管理
            posts := protected.Group("/posts")
            {
                posts.GET("", postHandler.GetPosts)
                posts.POST("", postHandler.CreatePost)
                posts.GET("/:id", postHandler.GetPost)
                posts.PUT("/:id", postHandler.UpdatePost)
                posts.DELETE("/:id", postHandler.DeletePost)
                posts.GET("/:id/comments", postHandler.GetPostComments)
                posts.POST("/:id/comments", postHandler.CreateComment)
            }
        }
        
        // 公开路由
        public := v1.Group("/public")
        {
            public.GET("/posts", postHandler.GetPublicPosts)
            public.GET("/posts/:id", postHandler.GetPublicPost)
        }
    }
    
    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
}
```

### 17.2.2.4 中间件实现

```go
// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "restful-api/internal/utils"
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取Authorization头
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Authorization header required",
            })
            c.Abort()
            return
        }
        
        // 检查Bearer前缀
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Invalid authorization header format",
            })
            c.Abort()
            return
        }
        
        // 验证token
        claims, err := utils.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Invalid token",
            })
            c.Abort()
            return
        }
        
        // 将用户信息添加到上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}

func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "User not authenticated",
            })
            c.Abort()
            return
        }
        
        if userRole != role {
            c.JSON(http.StatusForbidden, gin.H{
                "success": false,
                "message": "Insufficient permissions",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// internal/middleware/logger.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        logger.Info("HTTP Request",
            zap.String("method", param.Method),
            zap.String("path", param.Path),
            zap.Int("status", param.StatusCode),
            zap.Duration("latency", param.Latency),
            zap.String("client_ip", param.ClientIP),
            zap.String("user_agent", param.Request.UserAgent()),
        )
        return ""
    })
}

// internal/middleware/cors.go
package middleware

import (
    "github.com/gin-contrib/cors"
    "time"
)

func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

## 17.2.3 🎯 功能特性

### 17.2.3.1 用户管理

```go
// internal/handler/user.go
package handler

import (
    "net/http"
    "strconv"
    "restful-api/internal/model"
    "restful-api/internal/service"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    users, total, err := h.userService.GetUsers(page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.APIResponse{
            Success: false,
            Message: "Failed to get users",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "Users retrieved successfully",
        Data: gin.H{
            "users": users,
            "total": total,
            "page":  page,
            "limit": limit,
        },
    })
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid user ID",
        })
        return
    }
    
    user, err := h.userService.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, model.APIResponse{
            Success: false,
            Message: "User not found",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "User retrieved successfully",
        Data:    user,
    })
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid user ID",
        })
        return
    }
    
    var req model.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid request data",
            Error:   err.Error(),
        })
        return
    }
    
    user, err := h.userService.UpdateUser(uint(id), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.APIResponse{
            Success: false,
            Message: "Failed to update user",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "User updated successfully",
        Data:    user,
    })
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid user ID",
        })
        return
    }
    
    if err := h.userService.DeleteUser(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, model.APIResponse{
            Success: false,
            Message: "Failed to delete user",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "User deleted successfully",
    })
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, model.APIResponse{
            Success: false,
            Message: "User not authenticated",
        })
        return
    }
    
    user, err := h.userService.GetUserByID(userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, model.APIResponse{
            Success: false,
            Message: "User not found",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "Profile retrieved successfully",
        Data:    user,
    })
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, model.APIResponse{
            Success: false,
            Message: "User not authenticated",
        })
        return
    }
    
    var req model.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid request data",
            Error:   err.Error(),
        })
        return
    }
    
    user, err := h.userService.UpdateUser(userID.(uint), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.APIResponse{
            Success: false,
            Message: "Failed to update profile",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "Profile updated successfully",
        Data:    user,
    })
}
```

### 17.2.3.2 认证授权

```go
// internal/handler/auth.go
package handler

import (
    "net/http"
    "restful-api/internal/model"
    "restful-api/internal/service"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid request data",
            Error:   err.Error(),
        })
        return
    }
    
    user, err := h.authService.Register(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Failed to register user",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, model.APIResponse{
        Success: true,
        Message: "User registered successfully",
        Data:    user,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req model.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.APIResponse{
            Success: false,
            Message: "Invalid request data",
            Error:   err.Error(),
        })
        return
    }
    
    response, err := h.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, model.APIResponse{
            Success: false,
            Message: "Invalid credentials",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "Login successful",
        Data:    response,
    })
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    // 实现token刷新逻辑
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: "Token refreshed successfully",
    })
}
```

### 17.2.3.3 数据验证

```go
// internal/utils/validator.go
package utils

import (
    "github.com/go-playground/validator/v10"
    "reflect"
    "strings"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
    
    // 注册自定义验证器
    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
}

func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

func GetValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            field := e.Field()
            tag := e.Tag()
            
            switch tag {
            case "required":
                errors[field] = field + " is required"
            case "email":
                errors[field] = field + " must be a valid email"
            case "min":
                errors[field] = field + " must be at least " + e.Param() + " characters"
            case "max":
                errors[field] = field + " must be at most " + e.Param() + " characters"
            case "oneof":
                errors[field] = field + " must be one of: " + e.Param()
            default:
                errors[field] = field + " is invalid"
            }
        }
    }
    
    return errors
}
```

### 17.2.3.4 错误处理

```go
// internal/utils/response.go
package utils

import (
    "net/http"
    "restful-api/internal/model"
    "github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}, message string) {
    c.JSON(http.StatusOK, model.APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
    response := model.APIResponse{
        Success: false,
        Message: message,
    }
    
    if err != nil {
        response.Error = err.Error()
    }
    
    c.JSON(statusCode, response)
}

func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
    c.JSON(http.StatusBadRequest, model.APIResponse{
        Success: false,
        Message: "Validation failed",
        Data:    errors,
    })
}
```

## 17.2.4 🧪 测试实现

### 17.2.4.1 单元测试

```go
// internal/service/user_test.go
package service

import (
    "testing"
    "restful-api/internal/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
    args := m.Called(id)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
    args := m.Called(username)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockUserRepository) GetUsers(page, limit int) ([]*model.User, int64, error) {
    args := m.Called(page, limit)
    return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    userService := NewUserService(mockRepo)
    
    req := &model.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Role:     "user",
        IsActive: true,
    }
    
    mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
    
    user, err := userService.CreateUser(req)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser.Username, user.Username)
    assert.Equal(t, expectedUser.Email, user.Email)
    assert.Equal(t, expectedUser.Role, user.Role)
    assert.True(t, user.IsActive)
    
    mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
    mockRepo := new(MockUserRepository)
    userService := NewUserService(mockRepo)
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Role:     "user",
        IsActive: true,
    }
    
    mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)
    
    user, err := userService.GetUserByID(1)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    
    mockRepo.AssertExpectations(t)
}
```

### 17.2.4.2 集成测试

```go
// internal/handler/user_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "restful-api/internal/model"
    "restful-api/internal/service"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUsers(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    // 创建测试路由
    r := gin.New()
    
    // 模拟服务
    mockUserService := &MockUserService{}
    userHandler := NewUserHandler(mockUserService)
    
    r.GET("/users", userHandler.GetUsers)
    
    // 创建测试请求
    req, _ := http.NewRequest("GET", "/users?page=1&limit=10", nil)
    w := httptest.NewRecorder()
    
    // 执行请求
    r.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response model.APIResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success)
}

func TestUserHandler_CreateUser(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    r := gin.New()
    
    mockUserService := &MockUserService{}
    userHandler := NewUserHandler(mockUserService)
    
    r.POST("/users", userHandler.CreateUser)
    
    // 创建测试数据
    userData := model.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonData, _ := json.Marshal(userData)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response model.APIResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success)
}

// Mock服务实现
type MockUserService struct{}

func (m *MockUserService) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
    return &model.User{
        ID:       1,
        Username: req.Username,
        Email:    req.Email,
        Role:     "user",
        IsActive: true,
    }, nil
}

func (m *MockUserService) GetUserByID(id uint) (*model.User, error) {
    return &model.User{
        ID:       id,
        Username: "testuser",
        Email:    "test@example.com",
        Role:     "user",
        IsActive: true,
    }, nil
}

func (m *MockUserService) GetUsers(page, limit int) ([]*model.User, int64, error) {
    users := []*model.User{
        {
            ID:       1,
            Username: "testuser1",
            Email:    "test1@example.com",
            Role:     "user",
            IsActive: true,
        },
        {
            ID:       2,
            Username: "testuser2",
            Email:    "test2@example.com",
            Role:     "user",
            IsActive: true,
        },
    }
    return users, 2, nil
}

func (m *MockUserService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error) {
    return &model.User{
        ID:       id,
        Username: "updateduser",
        Email:    "updated@example.com",
        Role:     "user",
        IsActive: true,
    }, nil
}

func (m *MockUserService) DeleteUser(id uint) error {
    return nil
}
```

### 17.2.4.3 性能测试

```go
// internal/handler/benchmark_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "restful-api/internal/model"
    "github.com/gin-gonic/gin"
)

func BenchmarkUserHandler_GetUsers(b *testing.B) {
    gin.SetMode(gin.TestMode)
    
    r := gin.New()
    mockUserService := &MockUserService{}
    userHandler := NewUserHandler(mockUserService)
    
    r.GET("/users", userHandler.GetUsers)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        req, _ := http.NewRequest("GET", "/users?page=1&limit=10", nil)
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)
    }
}

func BenchmarkUserHandler_CreateUser(b *testing.B) {
    gin.SetMode(gin.TestMode)
    
    r := gin.New()
    mockUserService := &MockUserService{}
    userHandler := NewUserHandler(mockUserService)
    
    r.POST("/users", userHandler.CreateUser)
    
    userData := model.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonData, _ := json.Marshal(userData)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)
    }
}
```

## 17.2.5 🚀 部署运维

### 17.2.5.1 Docker部署

```dockerfile
# docker/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

```yaml
# docker/docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=restful_api
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
      - redis
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=restful_api
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - app
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:

networks:
  app-network:
    driver: bridge
```

### 17.2.5.2 监控日志

```go
// internal/middleware/metrics.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests in seconds",
        },
        []string{"method", "endpoint"},
    )
)

func Metrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            string(rune(status)),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration)
    }
}
```

### 17.2.5.3 性能优化

```go
// internal/middleware/cache.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cache"
    "github.com/gin-contrib/cache/persistence"
)

func Cache(store persistence.CacheStore) gin.HandlerFunc {
    return cache.CachePage(store, time.Minute*5, func(c *gin.Context) {
        c.Next()
    })
}

// 数据库连接池优化
// pkg/database/postgres.go
package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}
```

## 17.2.6 📚 扩展阅读

- [Gin框架官方文档](https://gin-gonic.com/docs/)
- [GORM官方文档](https://gorm.io/docs/)
- [JWT Go实现](https://github.com/golang-jwt/jwt)
- [Go测试最佳实践](https://github.com/golang/go/wiki/TestComments)
- [RESTful API设计指南](https://restfulapi.net/)
- [Docker最佳实践](https://docs.docker.com/develop/dev-best-practices/)
