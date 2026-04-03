# RESTful API Design Patterns

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #rest #api #http #json #openapi

---

## 1. REST Principles

### 1.1 REST Constraints

| Constraint | Description | Implementation |
|------------|-------------|----------------|
| Client-Server | Separation of concerns | API consumers independent of servers |
| Stateless | No client context on server | JWT tokens, session IDs in requests |
| Cacheable | Responses can be cached | Cache-Control headers, ETag |
| Uniform Interface | Standard methods and formats | HTTP verbs, resource URIs |
| Layered System | Client cannot tell if connected directly | Load balancers, proxies, gateways |
| Code on Demand (optional) | Server can extend client | JavaScript, WebAssembly |

### 1.2 HTTP Methods

| Method | Idempotent | Safe | Purpose |
|--------|-----------|------|---------|
| GET | Yes | Yes | Retrieve resource |
| POST | No | No | Create resource |
| PUT | Yes | No | Update/replace resource |
| PATCH | No | No | Partial update |
| DELETE | Yes | No | Remove resource |
| HEAD | Yes | Yes | Retrieve metadata |
| OPTIONS | Yes | Yes | Get available methods |

---

## 2. Resource Design

### 2.1 URI Naming Conventions

```
GOOD:
GET /users                    # Collection of users
GET /users/123                # Specific user
GET /users/123/orders         # User's orders
GET /users/123/orders/456     # Specific order of user
POST /users                   # Create new user
PUT /users/123                # Update user 123
PATCH /users/123              # Partial update
DELETE /users/123             # Delete user

BAD:
GET /getUsers                 # Verb in URI
GET /users/list               # Redundant
GET /user/123                 # Singular (inconsistent)
GET /Users/123                # Capitalization
GET /users/123/delete         # Verb in URI
```

### 2.2 Resource Relationships

```go
package api

// HATEOAS Example
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    Links     Links     `json:"_links"`
}

type Links struct {
    Self      Link `json:"self"`
    Orders    Link `json:"orders"`
    Profile   Link `json:"profile"`
}

type Link struct {
    Href  string `json:"href"`
    Title string `json:"title,omitempty"`
}

func NewUserResponse(user *domain.User) *User {
    return &User{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        Links: Links{
            Self:    Link{Href: fmt.Sprintf("/users/%s", user.ID)},
            Orders:  Link{Href: fmt.Sprintf("/users/%s/orders", user.ID)},
            Profile: Link{Href: fmt.Sprintf("/users/%s/profile", user.ID)},
        },
    }
}
```

---

## 3. Request/Response Patterns

### 3.1 Standard Response Format

```go
package api

// Standard API Response
type Response struct {
    Data      interface{} `json:"data,omitempty"`
    Error     *Error      `json:"error,omitempty"`
    Meta      *Meta       `json:"meta,omitempty"`
}

type Error struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}

type Meta struct {
    Page       int    `json:"page,omitempty"`
    PerPage    int    `json:"per_page,omitempty"`
    Total      int64  `json:"total,omitempty"`
    TotalPages int    `json:"total_pages,omitempty"`
}

func Success(data interface{}) *Response {
    return &Response{Data: data}
}

func SuccessWithMeta(data interface{}, meta *Meta) *Response {
    return &Response{Data: data, Meta: meta}
}

func ErrorResponse(code, message string) *Response {
    return &Response{
        Error: &Error{
            Code:    code,
            Message: message,
        },
    }
}
```

### 3.2 Pagination Patterns

```go
package api

// Offset-based pagination
type OffsetPagination struct {
    Offset int `query:"offset"`
    Limit  int `query:"limit" validate:"max=100"`
}

func (p *OffsetPagination) ToLimitOffset() (limit, offset int) {
    if p.Limit == 0 {
        p.Limit = 20
    }
    return p.Limit, p.Offset
}

// Cursor-based pagination
type CursorPagination struct {
    Cursor string `query:"cursor"`
    Limit  int    `query:"limit" validate:"max=100"`
}

type CursorResponse struct {
    Data       interface{} `json:"data"`
    NextCursor string      `json:"next_cursor,omitempty"`
    HasMore    bool        `json:"has_more"`
}

func EncodeCursor(timestamp time.Time, id string) string {
    data := fmt.Sprintf("%d:%s", timestamp.Unix(), id)
    return base64.URLEncoding.EncodeToString([]byte(data))
}

func DecodeCursor(cursor string) (time.Time, string, error) {
    data, err := base64.URLEncoding.DecodeString(cursor)
    if err != nil {
        return time.Time{}, "", err
    }

    parts := strings.Split(string(data), ":")
    if len(parts) != 2 {
        return time.Time{}, "", ErrInvalidCursor
    }

    ts, err := strconv.ParseInt(parts[0], 10, 64)
    if err != nil {
        return time.Time{}, "", err
    }

    return time.Unix(ts, 0), parts[1], nil
}
```

### 3.3 Filtering and Sorting

```go
package api

type ListOptions struct {
    Filters map[string]string `query:"filter"`
    Sort    string            `query:"sort"`
    Order   string            `query:"order" validate:"oneof=asc desc"`
}

// Query string examples:
// ?filter[status]=active&filter[role]=admin
// ?sort=created_at&order=desc

func ParseFilters(c *gin.Context) map[string]string {
    filters := make(map[string]string)
    filterPrefix := "filter["

    for key, values := range c.Request.URL.Query() {
        if strings.HasPrefix(key, filterPrefix) && strings.HasSuffix(key, "]") {
            field := key[len(filterPrefix) : len(key)-1]
            if len(values) > 0 {
                filters[field] = values[0]
            }
        }
    }

    return filters
}
```

---

## 4. Error Handling

### 4.1 HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Success |
| 201 | Created | Resource created |
| 204 | No Content | Success, empty body |
| 400 | Bad Request | Invalid input |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | No permission |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource conflict |
| 422 | Unprocessable | Validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Error | Server error |
| 502 | Bad Gateway | Upstream error |
| 503 | Service Unavailable | Temporary outage |

### 4.2 Error Implementation

```go
package api

import (
    "errors"
    "net/http"
)

var (
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrConflict     = errors.New("resource conflict")
)

type HTTPError struct {
    StatusCode int
    Code       string
    Message    string
    Details    map[string]string
}

func (e *HTTPError) Error() string {
    return e.Message
}

func NewHTTPError(statusCode int, code, message string) *HTTPError {
    return &HTTPError{
        StatusCode: statusCode,
        Code:       code,
        Message:    message,
    }
}

func MapError(err error) *HTTPError {
    switch {
    case errors.Is(err, ErrNotFound):
        return NewHTTPError(http.StatusNotFound, "NOT_FOUND", err.Error())
    case errors.Is(err, ErrInvalidInput):
        return NewHTTPError(http.StatusBadRequest, "INVALID_INPUT", err.Error())
    case errors.Is(err, ErrUnauthorized):
        return NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
    case errors.Is(err, ErrForbidden):
        return NewHTTPError(http.StatusForbidden, "FORBIDDEN", err.Error())
    case errors.Is(err, ErrConflict):
        return NewHTTPError(http.StatusConflict, "CONFLICT", err.Error())
    default:
        return NewHTTPError(http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
    }
}
```

---

## 5. Middleware Patterns

### 5.1 Common Middleware Stack

```go
package middleware

func SetupRouter() *gin.Engine {
    r := gin.New()

    // Recovery from panics
    r.Use(gin.Recovery())

    // Security headers
    r.Use(SecurityHeaders())

    // CORS
    r.Use(CORS())

    // Request ID
    r.Use(RequestID())

    // Logging
    r.Use(Logger())

    // Metrics
    r.Use(Metrics())

    // Rate limiting
    r.Use(RateLimiter())

    // Authentication
    r.Use(Authentication())

    return r
}
```

### 5.2 Request ID Middleware

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}
```

### 5.3 Logging Middleware

```go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        requestID := c.GetString("request_id")

        entry := logrus.WithFields(logrus.Fields{
            "request_id": requestID,
            "status":     statusCode,
            "latency":    latency,
            "client_ip":  clientIP,
            "method":     method,
            "path":       path,
        })

        if len(c.Errors) > 0 {
            entry.Error(c.Errors.String())
        } else if statusCode >= 500 {
            entry.Error("Server error")
        } else if statusCode >= 400 {
            entry.Warn("Client error")
        } else {
            entry.Info("Request processed")
        }
    }
}
```

---

## 6. OpenAPI Specification

### 6.1 Generating OpenAPI Spec

```go
package main

// @title Example API
// @version 1.0
// @description This is a sample API server
// @host api.example.com
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @Summary Create user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User information"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
    // Implementation
}
```

---

## 7. API Versioning

### 7.1 Versioning Strategies

| Strategy | Example | Pros | Cons |
|----------|---------|------|------|
| URL Path | /v1/users | Clear, simple | URL changes |
| Header | Accept-Version: v1 | Clean URLs | Less visible |
| Query Param | ?version=v1 | Simple | Messy URLs |
| Content-Type | application/vnd.api.v1+json | RESTful | Complex |

---

## 8. API Best Practices

- [ ] Use nouns, not verbs in URIs
- [ ] Use plural nouns for collections
- [ ] Use correct HTTP status codes
- [ ] Implement proper error handling
- [ ] Version your API
- [ ] Use HTTPS only
- [ ] Implement rate limiting
- [ ] Use pagination for large datasets
- [ ] Cache responses appropriately
- [ ] Document with OpenAPI
- [ ] Implement request/response validation
- [ ] Use content negotiation
- [ ] Implement health check endpoints

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 深度技术解析

### 核心概念

本部分深入分析核心技术概念和理论基础。

### 架构设计

`
系统架构图:
    [客户端]
       │
       ▼
   [API网关]
       │
   ┌───┴───┐
   ▼       ▼
[服务A] [服务B]
   │       │
   └───┬───┘
       ▼
   [数据库]
`

### 实现代码

`go
// 示例代码
package main

import (
    "context"
    "fmt"
)

func main() {
    ctx := context.Background()
    result := process(ctx)
    fmt.Println(result)
}

func process(ctx context.Context) string {
    select {
    case <-ctx.Done():
        return "timeout"
    default:
        return "success"
    }
}
`

### 性能特征

- 吞吐量: 高
- 延迟: 低
- 可扩展性: 良好
- 可用性: 99.99%

### 最佳实践

1. 使用连接池
2. 实现熔断机制
3. 添加监控指标
4. 记录详细日志

### 故障排查

| 症状 | 原因 | 解决方案 |
|------|------|----------|
| 超时 | 网络延迟 | 增加超时时间 |
| 错误 | 资源不足 | 扩容 |
| 慢查询 | 缺少索引 | 优化查询 |

### 相关技术

- 缓存技术 (Redis, Memcached)
- 消息队列 (Kafka, RabbitMQ)
- 数据库 (PostgreSQL, MySQL)
- 容器化 (Docker, Kubernetes)

### 学习资源

- 官方文档
- GitHub 仓库
- 技术博客
- 视频教程

### 社区支持

- Stack Overflow
- GitHub Issues
- 邮件列表
- Slack/Discord

---

## 高级主题

### 分布式一致性

CAP 定理和 BASE 理论的实际应用。

### 微服务架构

服务拆分、通信模式、数据一致性。

### 云原生设计

容器化、服务网格、可观测性。

---

**质量评级**: S (全面扩展)  
**完成日期**: 2026-04-02
---

## 深度技术解析

### 核心概念

本部分深入分析核心技术概念和理论基础。

### 架构设计

`
系统架构图:
    [客户端]
       │
       ▼
   [API网关]
       │
   ┌───┴───┐
   ▼       ▼
[服务A] [服务B]
   │       │
   └───┬───┘
       ▼
   [数据库]
`

### 实现代码

`go
// 示例代码
package main

import (
    "context"
    "fmt"
)

func main() {
    ctx := context.Background()
    result := process(ctx)
    fmt.Println(result)
}

func process(ctx context.Context) string {
    select {
    case <-ctx.Done():
        return "timeout"
    default:
        return "success"
    }
}
`

### 性能特征

- 吞吐量: 高
- 延迟: 低
- 可扩展性: 良好
- 可用性: 99.99%

### 最佳实践

1. 使用连接池
2. 实现熔断机制
3. 添加监控指标
4. 记录详细日志

### 故障排查

| 症状 | 原因 | 解决方案 |
|------|------|----------|
| 超时 | 网络延迟 | 增加超时时间 |
| 错误 | 资源不足 | 扩容 |
| 慢查询 | 缺少索引 | 优化查询 |

### 相关技术

- 缓存技术 (Redis, Memcached)
- 消息队列 (Kafka, RabbitMQ)
- 数据库 (PostgreSQL, MySQL)
- 容器化 (Docker, Kubernetes)

### 学习资源

- 官方文档
- GitHub 仓库
- 技术博客
- 视频教程

### 社区支持

- Stack Overflow
- GitHub Issues
- 邮件列表
- Slack/Discord

---

## 高级主题

### 分布式一致性

CAP 定理和 BASE 理论的实际应用。

### 微服务架构

服务拆分、通信模式、数据一致性。

### 云原生设计

容器化、服务网格、可观测性。

---

**质量评级**: S (全面扩展)  
**完成日期**: 2026-04-02