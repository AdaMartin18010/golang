# Go Web应用开发教程

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go Web应用开发教程](#go-web应用开发教程)
  - [📋 目录](#-目录)
  - [1. 项目概述](#1-项目概述)
    - [我们要构建什么](#我们要构建什么)
    - [技术栈](#技术栈)
  - [2. 环境准备](#2-环境准备)
    - [安装依赖](#安装依赖)
  - [3. 项目初始化](#3-项目初始化)
    - [创建项目结构](#创建项目结构)
    - [安装依赖包](#安装依赖包)
    - [配置文件](#配置文件)
    - [配置加载](#配置加载)
  - [4. 数据模型](#4-数据模型)
    - [用户模型](#用户模型)
    - [任务模型](#任务模型)
    - [数据库连接](#数据库连接)
  - [5. API实现](#5-api实现)
    - [JWT中间件](#jwt中间件)
    - [用户Handler](#用户handler)
    - [Todo Handler](#todo-handler)
    - [路由设置](#路由设置)
  - [6. 测试和部署](#6-测试和部署)
    - [API测试](#api测试)
    - [Docker部署](#docker部署)
  - [🎉 完成](#-完成)
  - [🔗 下一步](#-下一步)
  - [🔗 相关资源](#-相关资源)

---

---

## 1. 项目概述

### 我们要构建什么

一个简单的**任务管理系统（Todo App）**，包括：

- ✅ 用户注册和登录
- ✅ 创建、查看、更新、删除任务
- ✅ RESTful API
- ✅ JWT认证
- ✅ PostgreSQL数据库
- ✅ Docker部署

---

### 技术栈

- **Web框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **认证**: JWT
- **配置**: Viper
- **日志**: slog
- **容器化**: Docker

---

## 2. 环境准备

### 安装依赖

```bash
# 安装Go 1.25.3
# https://go.dev/dl/

# 验证安装
go version

# 安装PostgreSQL
# https://www.postgresql.org/download/

# 或使用Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine
```

---

## 3. 项目初始化

### 创建项目结构

```bash
mkdir todo-app && cd todo-app
go mod init github.com/username/todo-app

# 创建目录结构
mkdir -p cmd/api
mkdir -p internal/{handler,service,repository,model}
mkdir -p pkg/{config,database,middleware}
mkdir -p configs
```

---

### 安装依赖包

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get github.com/spf13/viper
go get golang.org/x/crypto/bcrypt
```

---

### 配置文件

```yaml
# configs/config.yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "todo_app"
  sslmode: "disable"

jwt:
  secret: "your-secret-key"
  expiration: 24h
```

---

### 配置加载

```go
// pkg/config/config.go
package config

import (
    "github.com/spf13/viper"
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Host string
    Port int
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

---

## 4. 数据模型

### 用户模型

```go
// internal/model/user.go
package model

import (
    "time"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Username  string         `gorm:"uniqueIndex;not null" json:"username"`
    Email     string         `gorm:"uniqueIndex;not null" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    Todos     []Todo         `gorm:"foreignKey:UserID" json:"todos,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword 密码加密
func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(u.Password),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.Password),
        []byte(password),
    )
    return err == nil
}
```

---

### 任务模型

```go
// internal/model/todo.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Todo struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Title       string         `gorm:"not null" json:"title"`
    Description string         `json:"description"`
    Completed   bool           `gorm:"default:false" json:"completed"`
    UserID      uint           `gorm:"not null" json:"user_id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

---

### 数据库连接

```go
// pkg/database/database.go
package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "todo-app/internal/model"
    "todo-app/pkg/config"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // 自动迁移
    if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
        return nil, err
    }

    return db, nil
}
```

---

## 5. API实现

### JWT中间件

```go
// pkg/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "strings"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", uint(claims["user_id"].(float64)))

        c.Next()
    }
}
```

---

### 用户Handler

```go
// internal/handler/user.go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "time"
    "todo-app/internal/model"
    "todo-app/internal/repository"
)

type UserHandler struct {
    userRepo  *repository.UserRepository
    jwtSecret string
    jwtExp    time.Duration
}

func NewUserHandler(userRepo *repository.UserRepository, jwtSecret string, jwtExp time.Duration) *UserHandler {
    return &UserHandler{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
        jwtExp:    jwtExp,
    }
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    if err := user.HashPassword(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    if err := h.userRepo.Create(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userRepo.FindByEmail(req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !user.CheckPassword(req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // 生成JWT Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(h.jwtExp).Unix(),
    })

    tokenString, err := token.SignedString([]byte(h.jwtSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}
```

---

### Todo Handler

```go
// internal/handler/todo.go
package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "todo-app/internal/model"
    "todo-app/internal/repository"
)

type TodoHandler struct {
    todoRepo *repository.TodoRepository
}

func NewTodoHandler(todoRepo *repository.TodoRepository) *TodoHandler {
    return &TodoHandler{todoRepo: todoRepo}
}

// List 获取任务列表
func (h *TodoHandler) List(c *gin.Context) {
    userID := c.GetUint("user_id")

    todos, err := h.todoRepo.FindByUserID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
        return
    }

    c.JSON(http.StatusOK, todos)
}

// Create 创建任务
func (h *TodoHandler) Create(c *gin.Context) {
    userID := c.GetUint("user_id")

    var req struct {
        Title       string `json:"title" binding:"required"`
        Description string `json:"description"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo := &model.Todo{
        Title:       req.Title,
        Description: req.Description,
        UserID:      userID,
    }

    if err := h.todoRepo.Create(todo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
        return
    }

    c.JSON(http.StatusCreated, todo)
}

// Update 更新任务
func (h *TodoHandler) Update(c *gin.Context) {
    userID := c.GetUint("user_id")
    id, _ := strconv.Atoi(c.Param("id"))

    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Completed   *bool  `json:"completed"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo, err := h.todoRepo.FindByID(uint(id))
    if err != nil || todo.UserID != userID {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    if req.Title != "" {
        todo.Title = req.Title
    }
    if req.Description != "" {
        todo.Description = req.Description
    }
    if req.Completed != nil {
        todo.Completed = *req.Completed
    }

    if err := h.todoRepo.Update(todo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
        return
    }

    c.JSON(http.StatusOK, todo)
}

// Delete 删除任务
func (h *TodoHandler) Delete(c *gin.Context) {
    userID := c.GetUint("user_id")
    id, _ := strconv.Atoi(c.Param("id"))

    todo, err := h.todoRepo.FindByID(uint(id))
    if err != nil || todo.UserID != userID {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    if err := h.todoRepo.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}
```

---

### 路由设置

```go
// cmd/api/main.go
package main

import (
    "fmt"
    "log"
    "github.com/gin-gonic/gin"
    "todo-app/internal/handler"
    "todo-app/internal/repository"
    "todo-app/pkg/config"
    "todo-app/pkg/database"
    "todo-app/pkg/middleware"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // 连接数据库
    db, err := database.Connect(&cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 初始化Repository
    userRepo := repository.NewUserRepository(db)
    todoRepo := repository.NewTodoRepository(db)

    // 初始化Handler
    userHandler := handler.NewUserHandler(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
    todoHandler := handler.NewTodoHandler(todoRepo)

    // 设置Gin路由
    r := gin.Default()

    // 公开路由
    r.POST("/api/register", userHandler.Register)
    r.POST("/api/login", userHandler.Login)

    // 受保护路由
    authorized := r.Group("/api")
    authorized.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
    {
        authorized.GET("/todos", todoHandler.List)
        authorized.POST("/todos", todoHandler.Create)
        authorized.PUT("/todos/:id", todoHandler.Update)
        authorized.DELETE("/todos/:id", todoHandler.Delete)
    }

    // 启动服务器
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

---

## 6. 测试和部署

### API测试

```bash
# 注册用户
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'

# 创建任务（需要替换TOKEN）
curl -X POST http://localhost:8080/api/todos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Complete Go tutorial"}'

# 获取任务列表
curl -X GET http://localhost:8080/api/todos \
  -H "Authorization: Bearer TOKEN"
```

---

### Docker部署

```dockerfile
# Dockerfile
FROM golang:1.25.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/configs ./configs

EXPOSE 8080

CMD ["./app"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_HOST=postgres
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: todo_app
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**运行**:

```bash
docker-compose up -d
```

---

## 🎉 完成

恭喜！你已经构建了一个完整的Go Web应用，包括：

- ✅ RESTful API
- ✅ JWT认证
- ✅ 数据库操作
- ✅ Docker部署

---

## 🔗 下一步

- 添加更多功能（任务分类、优先级、截止日期）
- 实现前端界面
- 添加单元测试和集成测试
- 部署到云平台

---

## 🔗 相关资源
