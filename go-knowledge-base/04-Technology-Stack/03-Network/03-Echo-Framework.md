# TS-NET-003: Echo Web Framework

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #echo #web-framework #golang #middleware #routing
> **权威来源**:
>
> - [Echo Documentation](https://echo.labstack.com/) - Official docs
> - [Echo GitHub](https://github.com/labstack/echo) - Source code

---

## 1. Echo Architecture

### 1.1 Core Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Echo Framework Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          Echo Instance                               │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │ Router (radix tree based, like Gin)                           │  │   │
│  │  │ - Static routes                                               │  │   │
│  │  │ - Parameter routes (:id)                                      │  │   │
│  │  │ - Wildcard routes (*)                                         │  │   │
│  │  └─────────────────────────────┬─────────────────────────────────┘  │   │
│  │                                │                                    │   │
│  │  ┌─────────────────────────────┴─────────────────────────────────┐  │   │
│  │  │                    Middleware Chain                            │  │   │
│  │  │  Pre → Router → Group → Route → Handler                      │  │   │
│  │  │                                                                │  │   │
│  │  │  Built-in:                                                     │  │   │
│  │  │  - Logger, Recover, CORS, CSRF, JWT                          │  │   │
│  │  │  - Gzip, Secure, Static, BodyLimit                           │  │   │
│  │  │  - MethodOverride, HTTPSRedirect                             │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                │                                    │   │
│  │  ┌─────────────────────────────┴─────────────────────────────────┐  │   │
│  │  │                      Context (echo.Context)                    │  │   │
│  │  │  - Request/Response                                            │  │   │
│  │  │  - Path/Query/Form params                                     │  │   │
│  │  │  - JSON/XML/HTML binding                                      │  │   │
│  │  │  - Validation                                                  │  │   │
│  │  │  - Session/Flash messages                                     │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Key Features:                                                               │
│  - Optimized router (zero dynamic memory allocation)                         │
│  - Scalable middleware system                                                │
│  - Data binding and validation                                               │
│  - JWT and OAuth2 support                                                    │
│  - WebSocket support                                                         │
│  - Automatic TLS (Let's Encrypt)                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Basic Usage

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    // Create instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    e.GET("/users/:id", getUser)
    e.POST("/users", createUser)
    e.PUT("/users/:id", updateUser)
    e.DELETE("/users/:id", deleteUser)

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}

func getUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "id": id,
    })
}

func createUser(c echo.Context) error {
    type User struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required,email"`
    }

    var u User
    if err := c.Bind(&u); err != nil {
        return err
    }
    if err := c.Validate(&u); err != nil {
        return err
    }

    return c.JSON(http.StatusCreated, u)
}

func updateUser(c echo.Context) error {
    // Implementation
    return nil
}

func deleteUser(c echo.Context) error {
    // Implementation
    return nil
}
```

---

## 3. Middleware

```go
// Custom middleware
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
        return next(c)
    }
}

// Group middleware
admin := e.Group("/admin")
admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    if username == "admin" && password == "secret" {
        return true, nil
    }
    return false, nil
}))

admin.GET("/dashboard", dashboardHandler)

// Route-specific middleware
e.GET("/special", specialHandler, middleware.CSRF())
```

---

## 4. Data Binding and Validation

```go
// Custom validator
type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    if err := cv.validator.Struct(i); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    return nil
}

// Configure
 e.Validator = &CustomValidator{validator: validator.New()}

// Binding
type User struct {
    Name     string `json:"name" form:"name" query:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=0,lte=130"`
    Password string `json:"password" validate:"required,min=8"`
}

func createUser(c echo.Context) error {
    u := new(User)
    if err := c.Bind(u); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    if err := c.Validate(u); err != nil {
        return err
    }
    return c.JSON(http.StatusCreated, u)
}
```

---

## 5. Error Handling

```go
// Custom HTTP error handler
 e.HTTPErrorHandler = func(err error, c echo.Context) {
    he, ok := err.(*echo.HTTPError)
    if !ok {
        he = &echo.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: http.StatusText(http.StatusInternalServerError),
        }
    }

    // Custom error response
    if !c.Response().Committed {
        if c.Request().Method == http.MethodHead {
            err = c.NoContent(he.Code)
        } else {
            err = c.JSON(he.Code, map[string]interface{}{
                "error": he.Message,
                "code":  he.Code,
            })
        }
        if err != nil {
            c.Echo().Logger.Error(err)
        }
    }
}
```

---

## 6. Checklist

```
Echo Best Practices:
□ Use structured logging
□ Implement proper error handling
□ Validate all input
□ Use middleware for cross-cutting concerns
□ Group routes logically
□ Configure timeouts
□ Enable CORS properly
□ Implement authentication/authorization
□ Use graceful shutdown
□ Monitor performance
```
