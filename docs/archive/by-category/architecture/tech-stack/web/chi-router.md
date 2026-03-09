# 1. 🌐 Chi Router 深度解析

> **简介**: 本文档详细阐述了 Chi Router 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26

---

## 📋 目录

- [1. 🌐 Chi Router 深度解析](#1--chi-router-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 基础路由配置](#131-基础路由配置)
    - [1.3.2 中间件详细使用](#132-中间件详细使用)
    - [1.3.3 路由参数绑定和验证](#133-路由参数绑定和验证)
    - [1.3.4 请求上下文传递](#134-请求上下文传递)
    - [1.3.5 文件上传处理](#135-文件上传处理)
    - [1.3.6 WebSocket 集成](#136-websocket-集成)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 中间件使用最佳实践](#141-中间件使用最佳实践)
    - [1.4.2 路由分组最佳实践](#142-路由分组最佳实践)
    - [1.4.3 参数验证最佳实践](#143-参数验证最佳实践)
    - [1.4.4 错误处理最佳实践](#144-错误处理最佳实践)
    - [1.4.5 性能优化最佳实践](#145-性能优化最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Chi 是什么？**

Chi 是一个轻量级、可组合的 Go HTTP 路由器，专注于提供简洁、高性能的路由功能。

**核心特性**:

- ✅ **轻量级**: 代码量小，依赖少
- ✅ **标准库兼容**: 完全基于 `net/http`，兼容所有标准库中间件
- ✅ **高性能**: 路由匹配速度快
- ✅ **中间件支持**: 丰富的中间件生态
- ✅ **路由组**: 支持路由分组和嵌套

---

## 1.2 选型论证

**为什么选择 Chi？**

**论证矩阵**:

| 评估维度 | 权重 | Chi | Gin | Echo | 说明 |
|---------|------|-----|-----|------|------|
| **标准库兼容** | 30% | 10 | 3 | 3 | Chi 完全基于标准库 |
| **学习成本** | 25% | 10 | 7 | 7 | Chi API 与标准库一致 |
| **性能** | 20% | 8 | 10 | 9 | 性能足够，不是瓶颈 |
| **功能丰富度** | 15% | 7 | 10 | 10 | 功能足够 |
| **维护成本** | 10% | 10 | 7 | 7 | 代码量小，易维护 |
| **加权总分** | - | **8.85** | 7.15 | 7.20 | Chi 得分最高 |

**核心优势**:

1. **标准库兼容性（权重 30%）**:
   - Chi 完全基于 `net/http`，可以使用所有标准库功能
   - 中间件生态丰富，兼容所有 `net/http` 中间件
   - 迁移成本极低，从标准库迁移几乎无缝

2. **学习成本低（权重 25%）**:
   - 团队成员都熟悉标准库，无需额外培训
   - API 与标准库一致，降低学习曲线
   - 文档简洁清晰，易于理解

3. **维护成本低（权重 10%）**:
   - 代码量小，易于理解和维护
   - 依赖少，减少安全风险
   - 更新频率低，稳定性好

**为什么不选择其他框架？**

1. **Gin**:
   - ✅ 性能优秀，功能丰富
   - ❌ 自定义路由，不兼容标准库
   - ❌ 学习成本高，需要学习新的 API
   - ❌ 中间件生态不如标准库丰富

2. **Echo**:
   - ✅ 功能丰富，性能优秀
   - ❌ 不兼容标准库
   - ❌ 学习成本高
   - ❌ 与 Gin 类似，无显著优势

---

## 1.3 实际应用

### 1.3.1 基础路由配置

**完整路由配置示例**:

```go
// internal/interfaces/http/chi/router.go
package chi

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // 全局中间件（按顺序执行）
    r.Use(middleware.RequestID)      // 为每个请求生成唯一 ID
    r.Use(middleware.RealIP)         // 获取真实 IP 地址
    r.Use(middleware.Logger)         // 请求日志
    r.Use(middleware.Recoverer)      // Panic 恢复
    r.Use(middleware.Compress(5))    // 响应压缩
    r.Use(middleware.Timeout(60 * time.Second)) // 请求超时

    // API 路由
    r.Route("/api/v1", func(r chi.Router) {
        r.Mount("/users", userRoutes())
        r.Mount("/workflows", workflowRoutes())
        r.Mount("/health", healthRoutes())
    })

    // 静态文件服务
    r.Mount("/static", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

    return r
}
```

### 1.3.2 中间件详细使用

**认证中间件**:

```go
// 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 从 Header 获取 Token
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // 验证 Token
        claims, err := validateJWT(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // 将用户信息添加到上下文
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        ctx = context.WithValue(ctx, "userRole", claims.Role)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 权限检查中间件
func RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := r.Context().Value("userRole").(string)

            if !hasPermission(userRole, permission) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**限流中间件**:

```go
import "golang.org/x/time/rate"

// 限流中间件
func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// 使用限流
func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // 创建限流器：每秒 100 个请求
    limiter := rate.NewLimiter(100, 100)
    r.Use(RateLimitMiddleware(limiter))

    return r
}
```

**CORS 中间件**:

```go
// CORS 中间件
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 1.3.3 路由参数绑定和验证

**URL 参数获取**:

```go
// 获取 URL 参数
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")

    // 验证 UUID 格式
    if _, err := uuid.Parse(userID); err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid user ID"))
        return
    }

    user, err := h.service.GetUser(r.Context(), userID)
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }

    Success(w, http.StatusOK, user)
}

// 获取查询参数
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    // 获取查询参数
    page := r.URL.Query().Get("page")
    pageSize := r.URL.Query().Get("page_size")

    // 解析和验证
    pageNum, _ := strconv.Atoi(page)
    if pageNum < 1 {
        pageNum = 1
    }

    size, _ := strconv.Atoi(pageSize)
    if size < 1 || size > 100 {
        size = 20
    }

    users, err := h.service.ListUsers(r.Context(), pageNum, size)
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }

    Success(w, http.StatusOK, users)
}
```

**请求体绑定**:

```go
// 请求体绑定和验证
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    // 绑定请求体
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid JSON"))
        return
    }

    // 验证请求参数
    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        Error(w, http.StatusBadRequest, errors.NewValidationError(err.Error()))
        return
    }

    user, err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }

    Success(w, http.StatusCreated, user)
}
```

### 1.3.4 请求上下文传递

**上下文传递示例**:

```go
// 在中间件中设置上下文值
func RequestContextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 获取请求 ID
        requestID := middleware.GetReqID(r.Context())

        // 创建新的上下文，添加请求信息
        ctx := r.Context()
        ctx = context.WithValue(ctx, "requestID", requestID)
        ctx = context.WithValue(ctx, "startTime", time.Now())
        ctx = context.WithValue(ctx, "clientIP", r.RemoteAddr)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 在 Handler 中使用上下文
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 从上下文获取请求 ID
    requestID := r.Context().Value("requestID").(string)

    // 在日志中使用请求 ID
    logger.Info("Creating user",
        "requestID", requestID,
        "path", r.URL.Path,
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)
    // ...
}
```

### 1.3.5 文件上传处理

**文件上传示例**:

```go
// 文件上传 Handler
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // 限制上传文件大小（10MB）
    r.ParseMultipartForm(10 << 20)

    // 获取上传的文件
    file, handler, err := r.FormFile("file")
    if err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("No file uploaded"))
        return
    }
    defer file.Close()

    // 验证文件类型
    if !isValidFileType(handler.Filename) {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid file type"))
        return
    }

    // 保存文件
    filePath := fmt.Sprintf("./uploads/%s", handler.Filename)
    dst, err := os.Create(filePath)
    if err != nil {
        Error(w, http.StatusInternalServerError, errors.NewInternalError("Failed to save file"))
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        Error(w, http.StatusInternalServerError, errors.NewInternalError("Failed to save file"))
        return
    }

    Success(w, http.StatusOK, map[string]string{
        "filename": handler.Filename,
        "size":     fmt.Sprintf("%d", handler.Size),
    })
}
```

### 1.3.6 WebSocket 集成

**WebSocket 集成示例**:

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 生产环境需要验证 Origin
    },
}

// WebSocket Handler
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // 升级到 WebSocket 连接
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error("WebSocket upgrade failed", "error", err)
        return
    }
    defer conn.Close()

    // 处理 WebSocket 消息
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            logger.Error("WebSocket read error", "error", err)
            break
        }

        // 处理消息
        response := h.processMessage(message)

        // 发送响应
        if err := conn.WriteMessage(messageType, response); err != nil {
            logger.Error("WebSocket write error", "error", err)
            break
        }
    }
}

// 路由配置
func websocketRoutes() chi.Router {
    r := chi.NewRouter()
    handler := NewWebSocketHandler()

    r.Get("/ws", handler.HandleWebSocket)

    return r
}
```

---

## 1.4 最佳实践

### 1.4.1 中间件使用最佳实践

**为什么需要中间件？**

中间件是处理横切关注点（Cross-Cutting Concerns）的最佳方式，可以统一处理日志、认证、追踪、限流等通用逻辑，避免在每个 Handler 中重复编写相同代码。

**中间件设计原则**:

1. **单一职责**: 每个中间件只负责一个功能
2. **可组合性**: 中间件可以组合使用
3. **可测试性**: 中间件可以独立测试
4. **性能考虑**: 避免在中间件中执行耗时操作

**实际应用示例**:

```go
// 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // 验证 token
        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // 将 userID 添加到上下文
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装 ResponseWriter 以捕获状态码
        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(ww, r)

        duration := time.Since(start)
        logger.Info("HTTP request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", ww.statusCode,
            "duration", duration,
        )
    })
}

// 追踪中间件
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, span := tracer.Start(r.Context(), r.URL.Path)
        defer span.End()

        span.SetAttributes(
            attribute.String("http.method", r.Method),
            attribute.String("http.path", r.URL.Path),
        )

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 使用中间件
func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // 全局中间件（按顺序执行）
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(LoggingMiddleware)
    r.Use(TracingMiddleware)
    r.Use(middleware.Recoverer)

    // 路由
    r.Route("/api/v1", func(r chi.Router) {
        // 公共路由
        r.Post("/login", loginHandler)

        // 需要认证的路由
        r.Group(func(r chi.Router) {
            r.Use(AuthMiddleware)
            r.Mount("/users", userRoutes())
        })
    })

    return r
}
```

**最佳实践要点**:

1. **中间件顺序**: 按照执行顺序排列中间件，例如 RequestID → Logging → Tracing → Auth → Handler
2. **错误处理**: 在中间件中正确处理错误，避免错误传播到 Handler
3. **上下文传递**: 使用 context 传递中间件处理的数据（如 userID、requestID）
4. **性能优化**: 避免在中间件中执行耗时操作，如数据库查询

### 1.4.2 路由分组最佳实践

**为什么需要路由分组？**

路由分组可以提高代码的可维护性和可读性，将相关的路由组织在一起，便于管理和测试。

**路由分组设计原则**:

1. **按功能分组**: 将相同功能的路由组织在一起
2. **按权限分组**: 将需要相同权限的路由组织在一起
3. **按版本分组**: 将不同版本的 API 分组管理
4. **嵌套分组**: 支持多级嵌套，提高灵活性

**实际应用示例**:

```go
// 用户路由组
func userRoutes() chi.Router {
    r := chi.NewRouter()
    handler := handlers.NewUserHandler(userService)

    // 用户列表和创建（需要认证）
    r.Group(func(r chi.Router) {
        r.Use(AuthMiddleware)
        r.Get("/", handler.ListUsers)
        r.Post("/", handler.CreateUser)
    })

    // 用户详情、更新、删除（需要认证和权限检查）
    r.Group(func(r chi.Router) {
        r.Use(AuthMiddleware)
        r.Use(RequirePermission("user:write"))
        r.Get("/{id}", handler.GetUser)
        r.Put("/{id}", handler.UpdateUser)
        r.Delete("/{id}", handler.DeleteUser)
    })

    return r
}

// 工作流路由组
func workflowRoutes() chi.Router {
    r := chi.NewRouter()
    handler := handlers.NewWorkflowHandler(workflowService)

    r.Use(AuthMiddleware)
    r.Use(RequirePermission("workflow:manage"))

    r.Post("/", handler.StartWorkflow)
    r.Get("/{id}", handler.GetWorkflow)
    r.Post("/{id}/signal", handler.SignalWorkflow)
    r.Get("/{id}/query", handler.QueryWorkflow)

    return r
}

// 版本化路由
func apiRoutes() chi.Router {
    r := chi.NewRouter()

    // v1 API
    r.Route("/v1", func(r chi.Router) {
        r.Mount("/users", userRoutes())
        r.Mount("/workflows", workflowRoutes())
    })

    // v2 API（未来版本）
    r.Route("/v2", func(r chi.Router) {
        // v2 路由
    })

    return r
}
```

**最佳实践要点**:

1. **功能内聚**: 将相关功能的路由组织在一起
2. **权限控制**: 在路由组级别应用权限中间件
3. **版本管理**: 使用路由分组管理 API 版本
4. **代码复用**: 提取公共路由逻辑，避免重复

### 1.4.3 参数验证最佳实践

**为什么需要参数验证？**

参数验证是保证 API 安全性和可靠性的重要手段，可以防止无效数据进入业务逻辑层，减少错误处理成本。

**参数验证设计原则**:

1. **早期验证**: 在 Handler 层进行参数验证，避免无效数据进入业务层
2. **统一验证**: 使用统一的验证库和验证规则
3. **清晰错误**: 返回清晰的验证错误信息
4. **类型安全**: 使用类型安全的验证方式

**实际应用示例**:

```go
// 使用 validator 库进行参数验证
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid JSON"))
        return
    }

    // 参数验证
    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        var validationErrors []string
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors, getValidationErrorMessage(err))
        }
        Error(w, http.StatusBadRequest, errors.NewValidationError(validationErrors))
        return
    }

    // 调用业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }

    Success(w, http.StatusCreated, user)
}

// 路由参数验证
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")

    // 验证 UUID 格式
    if _, err := uuid.Parse(userID); err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid user ID format"))
        return
    }

    user, err := h.service.GetUser(r.Context(), userID)
    if err != nil {
        if errors.Is(err, errors.ErrNotFound) {
            Error(w, http.StatusNotFound, err)
        } else {
            Error(w, http.StatusInternalServerError, err)
        }
        return
    }

    Success(w, http.StatusOK, user)
}
```

**最佳实践要点**:

1. **使用验证库**: 使用成熟的验证库（如 validator），避免手写验证逻辑
2. **验证规则**: 在结构体标签中定义验证规则，清晰直观
3. **错误信息**: 返回清晰的验证错误信息，帮助客户端理解问题
4. **类型转换**: 在验证后进行类型转换，确保类型安全

### 1.4.4 错误处理最佳实践

**为什么需要统一错误处理？**

统一的错误处理可以提高 API 的一致性和可维护性，便于客户端处理和错误监控。

**错误处理设计原则**:

1. **统一格式**: 所有错误使用统一的响应格式
2. **错误分类**: 区分不同类型的错误（业务错误、系统错误、验证错误）
3. **错误码**: 使用错误码标识错误类型
4. **错误日志**: 记录详细的错误日志，便于排查问题

**实际应用示例**:

```go
// 统一错误响应格式
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

// 错误处理中间件
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("Panic recovered",
                    "error", err,
                    "path", r.URL.Path,
                    "method", r.Method,
                )
                Error(w, http.StatusInternalServerError, errors.NewInternalError("Internal server error"))
            }
        }()

        next.ServeHTTP(w, r)
    })
}

// Handler 中的错误处理
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid request body"))
        return
    }

    user, err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        // 根据错误类型返回不同的状态码
        switch {
        case errors.Is(err, errors.ErrValidation):
            Error(w, http.StatusBadRequest, err)
        case errors.Is(err, errors.ErrConflict):
            Error(w, http.StatusConflict, err)
        case errors.Is(err, errors.ErrNotFound):
            Error(w, http.StatusNotFound, err)
        default:
            logger.Error("Unexpected error",
                "error", err,
                "path", r.URL.Path,
            )
            Error(w, http.StatusInternalServerError, errors.NewInternalError("Internal server error"))
        }
        return
    }

    Success(w, http.StatusCreated, user)
}
```

**最佳实践要点**:

1. **错误分类**: 区分业务错误和系统错误，返回不同的 HTTP 状态码
2. **错误码**: 使用错误码标识错误类型，便于客户端处理
3. **错误日志**: 记录详细的错误日志，包括请求信息、错误堆栈等
4. **错误恢复**: 使用 recover 捕获 panic，避免服务崩溃

### 1.4.5 性能优化最佳实践

**为什么需要性能优化？**

性能优化可以提高应用响应速度，减少服务器负载，改善用户体验。根据生产环境的实际经验，合理的性能优化可以将响应时间减少 50-80%，将吞吐量提升 2-5 倍。

**性能优化对比**:

| 优化项 | 未优化 | 优化后 | 提升比例 |
|--------|--------|--------|---------|
| **响应时间** | 100ms | 20-50ms | +50-80% |
| **吞吐量** | 1,000 QPS | 2,000-5,000 QPS | +100-400% |
| **带宽使用** | 100% | 20-40% | -60-80% |
| **CPU 使用** | 80% | 40-50% | -37.5-50% |

**性能优化策略**:

1. **连接池**: 使用 HTTP 连接池，复用连接（减少延迟 30-50%）
2. **响应压缩**: 启用响应压缩，减少传输数据量（减少带宽 60-80%）
3. **缓存**: 对静态资源和频繁访问的数据进行缓存（提升性能 5-10 倍）
4. **异步处理**: 对耗时操作使用异步处理（提升响应速度 50-80%）
5. **中间件优化**: 优化中间件执行顺序和逻辑（减少开销 20-30%）

**完整的性能优化示例**:

```go
// 生产环境级别的性能优化配置
func NewOptimizedRouter() *chi.Mux {
    r := chi.NewRouter()

    // 1. 性能优化中间件（按顺序）
    r.Use(middleware.RequestID)           // 请求 ID（开销小）
    r.Use(middleware.RealIP)              // 真实 IP（开销小）
    r.Use(middleware.Logger)              // 日志（异步记录）
    r.Use(middleware.Recoverer)          // Panic 恢复（必须）

    // 2. 响应压缩（关键优化，减少带宽 60-80%）
    r.Use(middleware.Compress(5))        // gzip 压缩级别 5（平衡压缩率和 CPU）

    // 3. 请求超时（避免长时间占用连接）
    r.Use(middleware.Timeout(60 * time.Second))

    // 4. 请求大小限制（防止内存溢出）
    r.Use(middleware.Throttle(100))      // 限制并发请求数

    // API 路由
    r.Route("/api/v1", func(r chi.Router) {
        r.Mount("/users", userRoutes())
        r.Mount("/workflows", workflowRoutes())
    })

    // 静态资源（启用缓存）
    r.Mount("/static", staticFileHandler())

    return r
}
```

**响应压缩优化**:

```go
// 响应压缩配置（减少带宽 60-80%）
func NewCompressionMiddleware() func(http.Handler) http.Handler {
    // 压缩级别：1-9，5 是性能和压缩率的平衡点
    return middleware.Compress(5,
        "text/html",
        "text/css",
        "text/plain",
        "text/javascript",
        "application/json",
        "application/javascript",
        "application/xml",
        "application/xhtml+xml",
    )
}

// 自定义压缩中间件（更精细控制）
func CustomCompressionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查客户端是否支持压缩
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        // 创建 gzip writer
        gz := gzip.NewWriter(w)
        defer gz.Close()

        // 设置响应头
        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Set("Vary", "Accept-Encoding")

        // 包装 ResponseWriter
        gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}
```

**缓存优化**:

```go
// 静态资源缓存（提升性能 5-10 倍）
func staticFileHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置缓存头
        w.Header().Set("Cache-Control", "public, max-age=3600")  // 1 小时
        w.Header().Set("ETag", generateETag(r.URL.Path))

        // 检查 If-None-Match 头（ETag 验证）
        if match := r.Header.Get("If-None-Match"); match != "" {
            if match == w.Header().Get("ETag") {
                w.WriteHeader(http.StatusNotModified)
                return
            }
        }

        http.ServeFile(w, r, r.URL.Path)
    })
}

// API 响应缓存（减少数据库查询）
func CacheMiddleware(ttl time.Duration) func(http.Handler) http.Handler {
    cache := make(map[string]*cachedResponse)
    mu := sync.RWMutex{}

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 只缓存 GET 请求
            if r.Method != http.MethodGet {
                next.ServeHTTP(w, r)
                return
            }

            // 生成缓存键
            cacheKey := r.URL.String()

            // 检查缓存
            mu.RLock()
            cached, ok := cache[cacheKey]
            mu.RUnlock()

            if ok && time.Since(cached.timestamp) < ttl {
                // 返回缓存响应
                for k, v := range cached.headers {
                    w.Header().Set(k, v)
                }
                w.WriteHeader(cached.statusCode)
                w.Write(cached.body)
                return
            }

            // 缓存响应
            cw := &cacheWriter{ResponseWriter: w}
            next.ServeHTTP(cw, r)

            // 保存到缓存
            mu.Lock()
            cache[cacheKey] = &cachedResponse{
                headers:   cw.headers,
                statusCode: cw.statusCode,
                body:      cw.body,
                timestamp: time.Now(),
            }
            mu.Unlock()
        })
    }
}

type cachedResponse struct {
    headers    map[string]string
    statusCode int
    body       []byte
    timestamp  time.Time
}

type cacheWriter struct {
    http.ResponseWriter
    headers    map[string]string
    statusCode int
    body       []byte
}

func (cw *cacheWriter) WriteHeader(code int) {
    cw.statusCode = code
    cw.headers = make(map[string]string)
    for k, v := range cw.Header() {
        cw.headers[k] = v[0]
    }
    cw.ResponseWriter.WriteHeader(code)
}

func (cw *cacheWriter) Write(b []byte) (int, error) {
    cw.body = append(cw.body, b...)
    return cw.ResponseWriter.Write(b)
}
```

**连接池优化**:

```go
// HTTP 客户端连接池配置（减少连接建立开销 30-50%）
func NewHTTPClient() *http.Client {
    transport := &http.Transport{
        MaxIdleConns:        100,              // 最大空闲连接数
        MaxIdleConnsPerHost: 10,              // 每个主机的最大空闲连接数
        IdleConnTimeout:     90 * time.Second, // 空闲连接超时
        DisableCompression:   false,          // 启用压缩
        DisableKeepAlives:    false,           // 启用 Keep-Alive
    }

    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
}

// 服务器端连接复用配置
func NewServer() *http.Server {
    return &http.Server{
        Addr:         ":8080",
        Handler:      NewOptimizedRouter(),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,  // Keep-Alive 超时
        MaxHeaderBytes: 1 << 20,        // 1MB
    }
}
```

**异步处理优化**:

```go
// 异步处理（提升响应速度 50-80%）
type AsyncHandler struct {
    workerPool chan struct{}
    wg         sync.WaitGroup
}

func NewAsyncHandler(workers int) *AsyncHandler {
    return &AsyncHandler{
        workerPool: make(chan struct{}, workers),
    }
}

func (h *AsyncHandler) HandleAsync(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 快速返回
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"status":"processing"}`))

        // 异步处理
        h.workerPool <- struct{}{}  // 获取 worker
        h.wg.Add(1)

        go func() {
            defer func() {
                <-h.workerPool  // 释放 worker
                h.wg.Done()
            }()

            // 执行耗时操作
            h.processAsync(r.Context(), r)
        }()
    })
}

func (h *AsyncHandler) processAsync(ctx context.Context, r *http.Request) {
    // 执行耗时操作（如发送邮件、生成报告等）
    // ...
}

// 使用示例
asyncHandler := NewAsyncHandler(10)  // 10 个并发 worker
r.Post("/users", asyncHandler.HandleAsync(userHandler.CreateUser))
```

**中间件性能优化**:

```go
// 中间件性能优化（减少开销 20-30%）
// 1. 避免在中间件中执行耗时操作
func OptimizedLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 异步记录日志，不阻塞请求
        go func() {
            duration := time.Since(start)
            logger.Info("HTTP request",
                "method", r.Method,
                "path", r.URL.Path,
                "duration", duration,
            )
        }()

        next.ServeHTTP(w, r)
    })
}

// 2. 使用对象池减少内存分配
var responseWriterPool = sync.Pool{
    New: func() interface{} {
        return &responseWriter{}
    },
}

func PooledLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 从对象池获取
        ww := responseWriterPool.Get().(*responseWriter)
        ww.ResponseWriter = w
        ww.statusCode = http.StatusOK

        defer func() {
            // 归还到对象池
            responseWriterPool.Put(ww)
        }()

        next.ServeHTTP(ww, r)
    })
}

// 3. 条件执行中间件（避免不必要的开销）
func ConditionalMiddleware(condition func(*http.Request) bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if condition(r) {
                // 执行中间件逻辑
                // ...
            }
            next.ServeHTTP(w, r)
        })
    }
}

// 使用示例：只对 API 路由执行认证
r.Use(ConditionalMiddleware(func(r *http.Request) bool {
    return strings.HasPrefix(r.URL.Path, "/api/")
})(AuthMiddleware))
```

**性能监控**:

```go
// 性能监控中间件
func PerformanceMonitoringMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(ww, r)

        duration := time.Since(start)

        // 记录性能指标
        metrics.RecordHTTPRequest(
            r.Method,
            r.URL.Path,
            ww.statusCode,
            duration,
        )

        // 慢请求告警
        if duration > 1*time.Second {
            logger.Warn("Slow request detected",
                "method", r.Method,
                "path", r.URL.Path,
                "duration", duration,
            )
        }
    })
}
```

**性能优化最佳实践要点**:

1. **连接复用**:
   - 使用 HTTP 连接池，减少连接建立开销（减少延迟 30-50%）
   - 配置合理的连接池大小
   - 启用 Keep-Alive

2. **响应压缩**:
   - 启用 gzip 压缩，减少传输数据量（减少带宽 60-80%）
   - 选择合适的压缩级别（推荐 5）
   - 只压缩文本类型

3. **缓存策略**:
   - 对静态资源启用缓存（提升性能 5-10 倍）
   - 对 API 响应使用 ETag 验证
   - 合理设置缓存过期时间

4. **异步处理**:
   - 对耗时操作使用异步处理（提升响应速度 50-80%）
   - 使用 worker pool 控制并发
   - 快速返回，后台处理

5. **中间件优化**:
   - 优化中间件执行顺序
   - 避免在中间件中执行耗时操作
   - 使用对象池减少内存分配
   - 条件执行中间件

6. **性能监控**:
   - 监控请求响应时间
   - 识别慢请求
   - 设置性能告警阈值

---

## 📚 扩展阅读

- [Chi Router 官方文档](https://github.com/go-chi/chi)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Chi Router 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
