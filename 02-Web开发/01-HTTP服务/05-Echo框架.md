# Echo框架基础

## 📚 **理论分析**

### **Echo框架简介**

- Echo是Go语言高性能、极简风格的Web框架，API风格类似Express。
- 支持高效路由、中间件、分组、RESTful API、WebSocket、静态文件服务等。
- 适合开发高性能API服务、微服务和Web应用。

### **核心原理**

- 路由基于高效的树结构，支持参数、通配符、分组
- 中间件采用链式调用，支持全局/分组/路由级中间件
- Context对象贯穿请求生命周期，便于数据传递和响应

### **主要类型与接口**

- `echo.Echo`：应用实例，负责路由和中间件管理
- `echo.Context`：请求上下文，封装请求、响应、参数、状态等
- `echo.HandlerFunc`：处理函数类型

## 💻 **代码示例**

### **最小Echo应用**

```go
package main
import "github.com/labstack/echo/v4"
func main() {
    e := echo.New()
    e.GET("/ping", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"message": "pong"})
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

### **路由与参数绑定**

```go
package main
import "github.com/labstack/echo/v4"
func main() {
    e := echo.New()
    e.GET("/user/:name", func(c echo.Context) error {
        name := c.Param("name")
        return c.String(200, "Hello "+name)
    })
    e.GET("/search", func(c echo.Context) error {
        q := c.QueryParam("q")
        return c.String(200, "Query: "+q)
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

### **中间件用法**

```go
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)
func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.GET("/", func(c echo.Context) error {
        return c.String(200, "Hello with middleware")
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

### **分组与RESTful API**

```go
package main
import "github.com/labstack/echo/v4"
func main() {
    e := echo.New()
    api := e.Group("/api/v1")
    api.GET("/users", func(c echo.Context) error {
        return c.JSON(200, map[string][]string{"users": {"Alice", "Bob"}})
    })
    api.POST("/users", func(c echo.Context) error {
        return c.JSON(201, map[string]string{"status": "created"})
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

## 🧪 **测试代码**

```go
package main
import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/labstack/echo/v4"
)
func TestPingRoute(t *testing.T) {
    e := echo.New()
    e.GET("/ping", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"message": "pong"})
    })
    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    rec := httptest.NewRecorder()
    e.ServeHTTP(rec, req)
    if rec.Code != 200 || rec.Body.String() != "{\"message\":\"pong\"}\n" {
        t.Errorf("unexpected response: %s", rec.Body.String())
    }
}
```

## 🎯 **最佳实践**

- 使用`echo.New()`自动集成日志与恢复中间件
- 路由分组便于模块化管理
- 参数校验与绑定建议用`Bind`方法
- 错误处理建议统一返回JSON结构
- 生产环境关闭debug模式，合理配置日志

## 🔍 **常见问题**

- Q: Echo和Gin有何区别？
  A: Echo更注重极简和性能，Gin生态更丰富
- Q: 如何自定义中间件？
  A: 实现`echo.MiddlewareFunc`并用`Use()`注册
- Q: 如何优雅关闭Echo服务？
  A: 通过`e.Shutdown(ctx)`实现

## 📚 **扩展阅读**

- [Echo官方文档](https://echo.labstack.com/guide)
- [Echo源码分析](https://github.com/labstack/echo)
- [Go by Example: Echo](https://gobyexample.com/echo)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
