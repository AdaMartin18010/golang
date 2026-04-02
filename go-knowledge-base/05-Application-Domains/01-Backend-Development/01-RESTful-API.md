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
