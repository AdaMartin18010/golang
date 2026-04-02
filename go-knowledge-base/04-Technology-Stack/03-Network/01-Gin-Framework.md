# Gin Web 框架

> **分类**: 开源技术堆栈

---

## 快速开始

```go
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run()  // 默认 :8080
}
```

---

## 路由

```go
// 参数
r.GET("/user/:name", func(c *gin.Context) {
    name := c.Param("name")
    c.String(200, "Hello %s", name)
})

// 查询
r.GET("/welcome", func(c *gin.Context) {
    firstname := c.DefaultQuery("firstname", "Guest")
    lastname := c.Query("lastname")
    c.String(200, "Hello %s %s", firstname, lastname)
})

// POST 表单
r.POST("/form", func(c *gin.Context) {
    message := c.PostForm("message")
    c.JSON(200, gin.H{"message": message})
})
```

---

## 中间件

```go
// 全局中间件
r.Use(gin.Logger())
r.Use(gin.Recovery())

// 自定义中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(401)
            return
        }
        c.Next()
    }
}

// 路由组
authorized := r.Group("/admin")
authorized.Use(AuthMiddleware())
{
    authorized.GET("/dashboard", dashboardHandler)
}
```

---

## 模型绑定

```go
type Login struct {
    User     string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

func login(c *gin.Context) {
    var json Login
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    if json.User != "manu" || json.Password != "123" {
        c.JSON(401, gin.H{"status": "unauthorized"})
        return
    }

    c.JSON(200, gin.H{"status": "you are logged in"})
}
```

---

## 性能

Gin 使用 httprouter，性能优异：

- 零内存分配
- 快速路由匹配
- 中间件链式调用
