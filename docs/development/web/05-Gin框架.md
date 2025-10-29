# Gin框架

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3, Gin v1.10+

---

> **难度**: ⭐⭐⭐
> **标签**: #Web #Gin #框架 #HTTP

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3, Gin v1.10+

---


---

## 📋 目录

- [1. Gin简介](#1.-gin简介)
  - [什么是Gin](#什么是gin)
  - [安装](#安装)
- [2. 快速开始](#2.-快速开始)
  - [Hello World](#hello-world)
  - [基本路由](#基本路由)
- [3. 路由](#3.-路由)
  - [路径参数](#路径参数)
  - [查询参数](#查询参数)
  - [路由分组](#路由分组)
- [4. 参数处理](#4.-参数处理)
  - [绑定JSON](#绑定json)
  - [绑定查询参数](#绑定查询参数)
  - [绑定URI](#绑定uri)
  - [文件上传](#文件上传)
- [5. 中间件](#5.-中间件)
  - [内置中间件](#内置中间件)
  - [自定义中间件](#自定义中间件)
  - [认证中间件](#认证中间件)
- [6. 响应处理](#6.-响应处理)
  - [JSON响应](#json响应)
  - [其他响应](#其他响应)
- [7. 最佳实践](#7.-最佳实践)
  - [1. 错误处理](#1.-错误处理)
  - [2. 项目结构](#2.-项目结构)
  - [3. 路由组织](#3.-路由组织)
  - [4. 优雅关闭](#4.-优雅关闭)
- [🔗 相关资源](#相关资源)

## 1. Gin简介

### 什么是Gin

**Gin** 是Go语言的高性能Web框架：
- 快速：性能接近net/http
- 中间件支持
- JSON验证
- 路由分组
- 错误管理

### 安装

```bash
go get -u github.com/gin-gonic/gin
```

---

## 2. 快速开始

### Hello World

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    r.Run(":8080")  // 监听8080端口
}
```

---

### 基本路由

```go
func main() {
    r := gin.Default()
    
    // GET
    r.GET("/users", getUsers)
    
    // POST
    r.POST("/users", createUser)
    
    // PUT
    r.PUT("/users/:id", updateUser)
    
    // DELETE
    r.DELETE("/users/:id", deleteUser)
    
    // PATCH
    r.PATCH("/users/:id", patchUser)
    
    r.Run(":8080")
}
```

---

## 3. 路由

### 路径参数

```go
// /users/123
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})

// /files/path/to/file.txt
r.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")
    c.JSON(200, gin.H{"filepath": filepath})
})
```

---

### 查询参数

```go
// /search?q=golang&page=1
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")              // 获取query参数
    page := c.DefaultQuery("page", "1") // 带默认值
    
    c.JSON(200, gin.H{
        "query": query,
        "page":  page,
    })
})
```

---

### 路由分组

```go
func main() {
    r := gin.Default()
    
    // API v1组
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", getUsersV1)
        v1.POST("/users", createUserV1)
    }
    
    // API v2组
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", getUsersV2)
        v2.POST("/users", createUserV2)
    }
    
    // 认证组
    authorized := r.Group("/admin")
    authorized.Use(AuthMiddleware())
    {
        authorized.GET("/dashboard", dashboard)
        authorized.POST("/users", adminCreateUser)
    }
    
    r.Run(":8080")
}
```

---

## 4. 参数处理

### 绑定JSON

```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

func createUser(c *gin.Context) {
    var user User
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "message": "User created",
        "user":    user,
    })
}
```

---

### 绑定查询参数

```go
type SearchQuery struct {
    Query string `form:"q" binding:"required"`
    Page  int    `form:"page" binding:"gte=1"`
    Limit int    `form:"limit" binding:"gte=1,lte=100"`
}

func search(c *gin.Context) {
    var query SearchQuery
    
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, query)
}
```

---

### 绑定URI

```go
type UserID struct {
    ID int `uri:"id" binding:"required,gt=0"`
}

func getUser(c *gin.Context) {
    var userID UserID
    
    if err := c.ShouldBindUri(&userID); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"user_id": userID.ID})
}

// 路由
r.GET("/users/:id", getUser)
```

---

### 文件上传

```go
// 单文件上传
func uploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    dst := "./uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"filename": file.Filename})
}

// 多文件上传
func uploadFiles(c *gin.Context) {
    form, _ := c.MultipartForm()
    files := form.File["files"]
    
    for _, file := range files {
        dst := "./uploads/" + file.Filename
        c.SaveUploadedFile(file, dst)
    }
    
    c.JSON(200, gin.H{"count": len(files)})
}
```

---

## 5. 中间件

### 内置中间件

```go
func main() {
    // Default包含Logger和Recovery中间件
    r := gin.Default()
    
    // 或手动添加
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    r.Run(":8080")
}
```

---

### 自定义中间件

```go
// 日志中间件
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // 处理请求
        c.Next()
        
        // 记录日志
        latency := time.Since(start)
        status := c.Writer.Status()
        
        log.Printf("%s %s %d %v", c.Request.Method, path, status, latency)
    }
}

// CORS中间件
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}

// 使用
r.Use(Logger())
r.Use(CORS())
```

---

### 认证中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        // 验证token
        userID, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 设置用户ID到Context
        c.Set("userID", userID)
        c.Next()
    }
}

// 使用
authorized := r.Group("/api")
authorized.Use(AuthMiddleware())
{
    authorized.GET("/profile", getProfile)
}

func getProfile(c *gin.Context) {
    userID := c.GetInt("userID")
    c.JSON(200, gin.H{"user_id": userID})
}
```

---

## 6. 响应处理

### JSON响应

```go
// 返回JSON
c.JSON(200, gin.H{
    "message": "Success",
    "data":    data,
})

// 返回结构体
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

c.JSON(200, Response{
    Code:    0,
    Message: "Success",
    Data:    user,
})

// 返回格式化的JSON（带缩进）
c.IndentedJSON(200, data)
```

---

### 其他响应

```go
// 字符串
c.String(200, "Hello, %s!", name)

// XML
c.XML(200, gin.H{"message": "Hello"})

// HTML
c.HTML(200, "index.html", gin.H{
    "title": "Home",
})

// 重定向
c.Redirect(302, "/new-url")

// 文件
c.File("./files/document.pdf")

// 下载文件
c.FileAttachment("./files/document.pdf", "report.pdf")

// 静态文件服务
r.Static("/assets", "./assets")
r.StaticFile("/favicon.ico", "./favicon.ico")
```

---

## 7. 最佳实践

### 1. 错误处理

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func handleError(c *gin.Context, code int, message string) {
    c.JSON(code, APIError{
        Code:    code,
        Message: message,
    })
}

func getUser(c *gin.Context) {
    id := c.Param("id")
    
    user, err := fetchUser(id)
    if err != nil {
        handleError(c, 500, "Failed to fetch user")
        return
    }
    
    if user == nil {
        handleError(c, 404, "User not found")
        return
    }
    
    c.JSON(200, user)
}
```

---

### 2. 项目结构

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   ├── user.go
│   │   └── auth.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── model/
│   │   └── user.go
│   └── service/
│       └── user.go
├── pkg/
│   └── response/
│       └── response.go
└── go.mod
```

---

### 3. 路由组织

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // 公共路由
    r.GET("/health", healthCheck)
    
    // API v1
    v1 := r.Group("/api/v1")
    {
        // 用户相关
        users := v1.Group("/users")
        {
            users.GET("", handler.ListUsers)
            users.POST("", handler.CreateUser)
            users.GET("/:id", handler.GetUser)
            users.PUT("/:id", handler.UpdateUser)
            users.DELETE("/:id", handler.DeleteUser)
        }
        
        // 认证相关
        auth := v1.Group("/auth")
        {
            auth.POST("/login", handler.Login)
            auth.POST("/register", handler.Register)
        }
    }
    
    return r
}
```

---

### 4. 优雅关闭

```go
func main() {
    r := SetupRouter()
    
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}
```

---

## 🔗 相关资源

- [HTTP协议](./01-HTTP协议.md)
- [RESTful API设计](./02-RESTful-API设计.md)
- [中间件开发](./10-中间件开发.md)

---

**最后更新**: 2025-10-29  
**Go版本**: 1.25.3  
**Gin版本**: v1.10+

