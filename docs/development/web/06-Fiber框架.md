# Fiber框架基础

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+ / Fiber v2.50+

---

## 📋 目录

- [📚 **理论分析**](#理论分析)
  - [**Fiber框架简介**](#fiber框架简介)
  - [**核心原理**](#核心原理)
  - [**主要类型与接口**](#主要类型与接口)
- [💻 **代码示例**](#代码示例)
  - [**最小Fiber应用**](#最小fiber应用)
  - [**路由与参数绑定**](#路由与参数绑定)
  - [**中间件用法**](#中间件用法)
  - [**分组与RESTful API**](#分组与restful-api)
- [🧪 **测试代码**](#测试代码)
- [🎯 **最佳实践**](#最佳实践)
- [🔍 **常见问题**](#常见问题)
- [📚 **扩展阅读**](#扩展阅读)

## 📚 **理论分析**

### **Fiber框架简介**

- Fiber是Go语言高性能Web框架，API风格类似Node.js的Express。
- 基于`fasthttp`库，极致追求性能，适合高并发API服务和微服务。
- 支持高效路由、中间件、分组、RESTful API、WebSocket、静态文件服务等。

### **核心原理**

- 路由基于树结构，支持参数、通配符、分组
- 中间件采用链式调用，支持全局/分组/路由级中间件
- Context对象贯穿请求生命周期，便于数据传递和响应

### **主要类型与接口**

- `fiber.App`：应用实例，负责路由和中间件管理
- `fiber.Ctx`：请求上下文，封装请求、响应、参数、状态等
- `fiber.Handler`：处理函数类型

## 💻 **代码示例**

### **最小Fiber应用**

```go
package main
import "github.com/gofiber/fiber/v2"
func main() {
    app := fiber.New()
    app.Get("/ping", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "pong"})
    })
    app.Listen(":8080")
}
```

### **路由与参数绑定**

```go
package main
import "github.com/gofiber/fiber/v2"
func main() {
    app := fiber.New()
    app.Get("/user/:name", func(c *fiber.Ctx) error {
        name := c.Params("name")
        return c.SendString("Hello " + name)
    })
    app.Get("/search", func(c *fiber.Ctx) error {
        q := c.Query("q")
        return c.SendString("Query: " + q)
    })
    app.Listen(":8080")
}
```

### **中间件用法**

```go
package main
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)
func main() {
    app := fiber.New()
    app.Use(logger.New())
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello with middleware")
    })
    app.Listen(":8080")
}
```

### **分组与RESTful API**

```go
package main
import "github.com/gofiber/fiber/v2"
func main() {
    app := fiber.New()
    api := app.Group("/api/v1")
    api.Get("/users", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"users": []string{"Alice", "Bob"}})
    })
    api.Post("/users", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "created"})
    })
    app.Listen(":8080")
}
```

## 🧪 **测试代码**

```go
package main
import (
    "net/http/httptest"
    "testing"
    "github.com/gofiber/fiber/v2"
)
func TestPingRoute(t *testing.T) {
    app := fiber.New()
    app.Get("/ping", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "pong"})
    })
    req := httptest.NewRequest("GET", "/ping", nil)
    resp, _ := app.Test(req)
    if resp.StatusCode != 200 {
        t.Errorf("unexpected status: %d", resp.StatusCode)
    }
}
```

## 🎯 **最佳实践**

- 使用`fiber.New()`自动集成日志与恢复中间件
- 路由分组便于模块化管理
- 参数校验与绑定建议用`BodyParser`方法
- 错误处理建议统一返回JSON结构
- 生产环境关闭debug模式，合理配置日志

## 🔍 **常见问题**

- Q: Fiber和Gin/Echo有何区别？
  A: Fiber基于fasthttp，极致追求性能，API风格更接近Express
- Q: 如何自定义中间件？
  A: 实现`fiber.Handler`并用`Use()`注册
- Q: 如何优雅关闭Fiber服务？
  A: 通过`app.Shutdown()`实现

## 📚 **扩展阅读**

- [Fiber官方文档](https://docs.gofiber.io/)
- [Fiber源码分析](https://github.com/gofiber/fiber)
- [Go by Example: Fiber](https://gobyexample.com/fiber)

---

**文档维护者**: Go Documentation Team
**最后更新**: 2025-10-29
**文档状态**: 完成
**适用版本**: Go 1.25.3+
