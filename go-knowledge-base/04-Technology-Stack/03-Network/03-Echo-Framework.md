# TS-NET-003: Echo Web Framework

> **维度**: Technology Stack > Network
> **级别**: S (18+ KB)
> **标签**: #echo #web-framework #golang #middleware #routing
> **权威来源**:
>
> - [Echo Documentation](https://echo.labstack.com/) - Official docs
> - [Echo GitHub](https://github.com/labstack/echo) - Source code

---

## 1. Echo Architecture Deep Dive

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
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Request Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Echo Request Lifecycle                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   HTTP Request                                                               │
│        │                                                                     │
│        ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 1. Pre-Router Middleware                                           │   │
│   │    - HTTPSRedirect, RemoveTrailingSlash                            │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│        ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 2. Router Match                                                    │   │
│   │    - Find handler by method + path                                 │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│        ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 3. Route Group Middleware                                          │   │
│   │    - Authentication, Rate limiting                                 │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│        ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 4. Route Middleware                                                │   │
│   │    - Validation, Body dump                                         │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│        ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 5. Handler                                                         │   │
│   │    - Business logic, Response generation                           │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│        │                                                                     │
│        ▼                                                                     │
│   HTTP Response                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Basic Usage

### 2.1 Creating Server

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

    // Global middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Routes
    e.GET("/", homeHandler)
    e.GET("/users/:id", getUser)
    e.POST("/users", createUser)
    e.PUT("/users/:id", updateUser)
    e.DELETE("/users/:id", deleteUser)

    // Route groups
    api := e.Group("/api")
    api.Use(middleware.JWT([]byte("secret")))
    api.GET("/users", getUsers)

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}
```

### 2.2 Handler Functions

```go
// Path parameters
func getUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "id": id,
    })
}

// Query parameters
func searchUsers(c echo.Context) error {
    name := c.QueryParam("name")
    page := c.QueryParam("page")
    return c.JSON(http.StatusOK, map[string]string{
        "name": name,
        "page": page,
    })
}

// Form data
func createUser(c echo.Context) error {
    name := c.FormValue("name")
    email := c.FormValue("email")
    return c.JSON(http.StatusCreated, map[string]string{
        "name":  name,
        "email": email,
    })
}

// JSON binding
func createUserJSON(c echo.Context) error {
    type User struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"gte=0,lte=130"`
    }

    var u User
    if err := c.Bind(&u); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    if err := c.Validate(&u); err != nil {
        return err
    }

    return c.JSON(http.StatusCreated, u)
}
```

---

## 3. Middleware

### 3.1 Built-in Middleware

```go
func main() {
    e := echo.New()

    // Logger
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "${time_rfc3339} | ${status} | ${latency_human} | ${method} ${uri}\n",
    }))

    // Recover from panics
    e.Use(middleware.Recover())

    // CORS
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"https://example.com"},
        AllowMethods: []string{http.MethodGet, http.MethodPost},
    }))

    // Rate limiting
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

    // Body limit
    e.Use(middleware.BodyLimit("2M"))

    // Gzip compression
    e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5,
    }))

    // Security headers
    e.Use(middleware.Secure())

    // JWT authentication
    e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte("secret"),
    }))
}
```

### 3.2 Custom Middleware

```go
// Simple middleware
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set(echo.HeaderServer, "Echo/4.0")
        return next(c)
    }
}

// Middleware with config
type AuthConfig struct {
    Skipper middleware.Skipper
    Token   string
}

func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
    if config.Skipper == nil {
        config.Skipper = middleware.DefaultSkipper
    }

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if config.Skipper(c) {
                return next(c)
            }

            token := c.Request().Header.Get("X-Auth-Token")
            if token != config.Token {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }

            return next(c)
        }
    }
}

// Usage
e.Use(AuthWithConfig(AuthConfig{
    Token: "my-secret-token",
    Skipper: func(c echo.Context) bool {
        return c.Path() == "/health"
    },
}))
```

---

## 4. Configuration Best Practices

### 4.1 Production Configuration

```go
func createServer() *echo.Echo {
    e := echo.New()

    // Hide server banner
    e.HideBanner = true

    // Custom HTTP error handler
    e.HTTPErrorHandler = customHTTPErrorHandler

    // Custom validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Custom binder
    e.Binder = &CustomBinder{}

    // Set renderer
    e.Renderer = &TemplateRenderer{}

    return e
}

// Custom error handler
func customHTTPErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := http.StatusText(code)

    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        message = he.Message.(string)
    }

    if !c.Response().Committed {
        c.JSON(code, map[string]interface{}{
            "error": message,
            "code":  code,
        })
    }
}
```

---

## 5. Comparison with Alternatives

| Framework | Performance | Features | Learning Curve | Use Case |
|-----------|-------------|----------|----------------|----------|
| **Echo** | Excellent | Rich | Low | APIs, microservices |
| **Gin** | Excellent | Rich | Low | APIs, web apps |
| **Fiber** | Fastest | Growing | Low | High performance |
| **Standard** | Good | Basic | None | Simple apps |
| **Chi** | Good | Minimal | Low | Minimalist APIs |

---

## 6. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Echo Best Practices                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Setup:                                                                      │
│  □ Use structured logging                                                    │
│  □ Implement proper error handling                                           │
│  □ Configure CORS appropriately                                              │
│  □ Enable security headers                                                   │
│                                                                              │
│  Routing:                                                                    │
│  □ Group routes logically                                                    │
│  □ Use middleware at appropriate levels                                      │
│  □ Validate all input                                                        │
│  □ Document API with OpenAPI/Swagger                                         │
│                                                                              │
│  Production:                                                                 │
│  □ Implement graceful shutdown                                               │
│  □ Configure timeouts                                                        │
│  □ Monitor performance                                                       │
│  □ Set up health checks                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18+ KB, comprehensive coverage)
