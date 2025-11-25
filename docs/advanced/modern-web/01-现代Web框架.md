# 现代Web框架

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [现代Web框架](#现代web框架)
  - [📋 目录](#-目录)
  - [1. 📖 主流Go Web框架](#1--主流go-web框架)
    - [Gin框架](#gin框架)
    - [Fiber框架](#fiber框架)
    - [Echo框架](#echo框架)
  - [🎯 框架对比](#-框架对比)
  - [💡 中间件开发](#-中间件开发)
  - [📚 相关资源](#-相关资源)

---

## 1. 📖 主流Go Web框架

### Gin框架

```go
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    // 路由组
    api := r.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }

    // 中间件
    authorized := r.Group("/admin")
    authorized.Use(AuthMiddleware())
    {
        authorized.GET("/dashboard", dashboard)
    }

    r.Run(":8080")
}

// 请求绑定
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=0,lte=130"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 创建用户逻辑...

    c.JSON(201, gin.H{
        "message": "User created",
        "data":    req,
    })
}
```

---

### Fiber框架

```go
import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New(fiber.Config{
        Prefork: true, // 多进程模式
    })

    // 中间件
    app.Use(logger.New())
    app.Use(cors.New())

    // 路由
    app.Get("/api/users", getUsers)
    app.Post("/api/users", createUser)

    // 静态文件
    app.Static("/", "./public")

    app.Listen(":3000")
}

// Fiber Handler
func getUsers(c *fiber.Ctx) error {
    users := []User{...}
    return c.JSON(users)
}
```

---

### Echo框架

```go
import "github.com/labstack/echo/v4"

func main() {
    e := echo.New()

    // 中间件
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // 路由
    e.GET("/users/:id", getUser)
    e.POST("/users", createUser)

    // 启动
    e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "id": id,
    })
}
```

---

## 🎯 框架对比

| 特性 | Gin | Fiber | Echo |
|------|-----|-------|------|
| 性能 | 高 | 极高 | 高 |
| 学习曲线 | 平缓 | 平缓 | 平缓 |
| 生态 | 丰富 | 快速增长 | 成熟 |
| 适用场景 | 通用 | 高性能API | 企业应用 |

---

## 💡 中间件开发

```go
// Gin自定义中间件
func RateLimitMiddleware(limit int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(limit), limit)

    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
r.Use(RateLimitMiddleware(100))
```

---

## 📚 相关资源
