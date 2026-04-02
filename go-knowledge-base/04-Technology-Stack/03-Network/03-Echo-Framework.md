# Echo Web 框架

> **分类**: 开源技术堆栈

---

## 快速开始

```go
import "github.com/labstack/echo/v4"

func main() {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    e.Start(":8080")
}
```

---

## 路由

```go
// 参数
e.GET("/users/:id", getUser)

// 查询参数
name := c.QueryParam("name")

// POST
payload := &struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}{}
if err := c.Bind(payload); err != nil {
    return err
}
```

---

## 中间件

```go
// 内置中间件
e.Use(middleware.Logger())
e.Use(middleware.Recover())
e.Use(middleware.CORS())

// 自定义
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("X-Custom-Header", "value")
        return next(c)
    }
}
e.Use(ServerHeader)
```

---

## 组路由

```go
g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    if username == "admin" && password == "secret" {
        return true, nil
    }
    return false, nil
}))

g.GET("/main", mainAdmin)
```

---

## 验证

```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=130"`
}

e.POST("/users", func(c echo.Context) error {
    u := new(User)
    if err := c.Bind(u); err != nil {
        return err
    }
    if err := c.Validate(u); err != nil {
        return err
    }
    return c.JSON(http.StatusCreated, u)
})
```

---

## 与 Gin 对比

| 特性 | Echo | Gin |
|------|------|-----|
| 路由 | 基于 radix tree | 基于 httprouter |
| 性能 | 快 | 更快 |
| 中间件 | 丰富 | 丰富 |
| 验证 | 内置 | 需第三方 |
| 文档 | 详细 | 详细 |
