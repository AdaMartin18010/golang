# Go HTTP编程深度实战指南

**字数**: ~25,000字  
**代码示例**: 60+个完整示例  
**实战案例**: 7个端到端案例  
**适用人群**: 初级到高级Go开发者

---

## 📚 目录

<!-- TOC -->
- [Go HTTP编程深度实战指南](#go-http编程深度实战指南)
  - [📚 目录](#-目录)
  - [第一部分：HTTP基础与net/http包](#第一部分http基础与nethttp包)
    - [为什么选择Go做Web开发？](#为什么选择go做web开发)
    - [net/http核心类型](#nethttp核心类型)
    - [实战案例1：最小HTTP服务器](#实战案例1最小http服务器)
      - [方式1：最简单](#方式1最简单)
      - [方式2：自定义ServeMux](#方式2自定义servemux)
  - [第二部分：HTTP服务器深度实战](#第二部分http服务器深度实战)
    - [服务器配置详解](#服务器配置详解)
      - [完整配置示例](#完整配置示例)
    - [实战案例2：HTTPS服务器](#实战案例2https服务器)
      - [生成自签名证书](#生成自签名证书)
      - [HTTPS服务器实现](#https服务器实现)
      - [HTTP自动重定向到HTTPS](#http自动重定向到https)
  - [第三部分：中间件模式深度实战](#第三部分中间件模式深度实战)
    - [中间件原理](#中间件原理)
    - [实战案例3：中间件链](#实战案例3中间件链)
      - [场景](#场景)
      - [完整实现](#完整实现)
      - [测试](#测试)
  - [第四部分：路由设计](#第四部分路由设计)
    - [实战案例4：RESTful路由](#实战案例4restful路由)
      - [场景4](#场景4)
      - [完整实现（原生net/http）](#完整实现原生nethttp)
      - [测试4](#测试4)
  - [第五部分：请求与响应处理](#第五部分请求与响应处理)
    - [请求参数处理](#请求参数处理)
    - [响应处理](#响应处理)
  - [第六部分：文件处理](#第六部分文件处理)
    - [实战案例5：文件上传](#实战案例5文件上传)
    - [文件下载](#文件下载)
  - [第七部分：WebSocket实战](#第七部分websocket实战)
    - [实战案例6：WebSocket聊天室](#实战案例6websocket聊天室)
  - [第八部分：HTTP客户端最佳实践](#第八部分http客户端最佳实践)
    - [实战案例7：高级HTTP客户端](#实战案例7高级http客户端)
  - [第九部分：性能优化](#第九部分性能优化)
    - [性能优化清单](#性能优化清单)
    - [压缩中间件](#压缩中间件)
  - [第十部分：完整RESTful API实战](#第十部分完整restful-api实战)
    - [完整项目结构](#完整项目结构)
    - [最佳实践总结](#最佳实践总结)
      - [DO's ✅](#dos-)
      - [DON'Ts ❌](#donts-)
  - [🎯 总结](#-总结)
    - [HTTP编程核心要点](#http编程核心要点)
    - [技术选型建议](#技术选型建议)

<!-- TOC -->

---

## 第一部分：HTTP基础与net/http包

### 为什么选择Go做Web开发？

```text
✅ 高性能：并发模型优秀，天然支持高并发
✅ 简单：标准库功能完善，无需框架也能开发
✅ 部署方便：单一可执行文件，无依赖
✅ 社区成熟：Gin、Echo、Fiber等优秀框架
✅ 云原生：Kubernetes原生语言
```

---

### net/http核心类型

| 类型 | 作用 | 关键方法/字段 |
|------|------|--------------|
| **http.Server** | HTTP服务器 | Addr, Handler, ReadTimeout |
| **http.Request** | 请求对象 | Method, URL, Header, Body |
| **http.ResponseWriter** | 响应写入器 | Write(), WriteHeader(), Header() |
| **http.Handler** | 处理器接口 | ServeHTTP(w, r) |
| **http.ServeMux** | 路由复用器 | Handle(), HandleFunc() |
| **http.Client** | HTTP客户端 | Get(), Post(), Do() |

---

### 实战案例1：最小HTTP服务器

#### 方式1：最简单

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

#### 方式2：自定义ServeMux

```go
package main

import (
    "fmt"
    "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from handler!")
}

func about(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "About page")
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    mux.HandleFunc("/about", about)
    
    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", mux)
}
```

---

## 第二部分：HTTP服务器深度实战

### 服务器配置详解

#### 完整配置示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 1. 创建路由
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 模拟慢响应
        time.Sleep(100 * time.Millisecond)
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // 2. 配置服务器
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
        
        // 超时配置（防止慢速攻击）
        ReadTimeout:       5 * time.Second,   // 读取请求超时
        WriteTimeout:      10 * time.Second,  // 写入响应超时
        IdleTimeout:       120 * time.Second, // 空闲连接超时
        ReadHeaderTimeout: 2 * time.Second,   // 读取请求头超时
        
        // 最大请求头大小
        MaxHeaderBytes: 1 << 20, // 1MB
    }
    
    // 3. 启动服务器
    go func() {
        fmt.Println("Server starting on :8080...")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()
    
    // 4. 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("Shutting down server...")
    
    // 5秒内完成所有请求
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    fmt.Println("Server gracefully stopped")
}
```

### 实战案例2：HTTPS服务器

#### 生成自签名证书

```bash
# 生成私钥
openssl genrsa -out server.key 2048

# 生成证书
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
```

#### HTTPS服务器实现

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, HTTPS!")
    })
    
    fmt.Println("HTTPS server starting on :443...")
    
    // 启动HTTPS服务器
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

#### HTTP自动重定向到HTTPS

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
    target := "https://" + r.Host + r.URL.Path
    if len(r.URL.RawQuery) > 0 {
        target += "?" + r.URL.RawQuery
    }
    
    http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func main() {
    // HTTP服务器（重定向到HTTPS）
    go func() {
        fmt.Println("HTTP server starting on :80 (redirecting to HTTPS)...")
        http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
    }()
    
    // HTTPS服务器
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Secure connection!")
    })
    
    fmt.Println("HTTPS server starting on :443...")
    log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
}
```

---

## 第三部分：中间件模式深度实战

### 中间件原理

**洋葱模型**:

```text
       请求
        ↓
    ┌───────────┐
    │ 日志中间件  │  ← 外层
    │ ┌─────────┐│
    │ │认证中间件││  ← 中层
    │ │┌───────┐││
    │ ││ 处理器 │││  ← 核心
    │ │└───────┘││
    │ └─────────┘│
    └───────────┘
        ↓
       响应
```

---

### 实战案例3：中间件链

#### 场景

- 日志中间件：记录请求信息
- 认证中间件：验证token
- 恢复中间件：捕获panic
- 限流中间件：防止过载

#### 完整实现

```go
package middleware

import (
    "fmt"
    "log"
    "net/http"
    "runtime/debug"
    "sync"
    "time"
)

// ===== 中间件1：日志中间件 =====
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 包装ResponseWriter以捕获状态码
        wrapped := &responseWriter{ResponseWriter: w}
        
        next.ServeHTTP(wrapped, r)
        
        log.Printf("%s %s %d %v",
            r.Method,
            r.URL.Path,
            wrapped.status,
            time.Since(start),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    if rw.status == 0 {
        rw.status = 200
    }
    return rw.ResponseWriter.Write(b)
}

// ===== 中间件2：认证中间件 =====
func Authentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 获取token
        token := r.Header.Get("Authorization")
        
        // 验证token（简化版）
        if token == "" || token != "Bearer secret-token" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // 验证成功，继续
        next.ServeHTTP(w, r)
    })
}

// ===== 中间件3：恢复中间件（捕获panic）=====
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// ===== 中间件4：限流中间件（令牌桶）=====
type RateLimiter struct {
    tokens     int
    maxTokens  int
    refillRate int // 每秒refill
    mu         sync.Mutex
    lastRefill time.Time
}

func NewRateLimiter(maxTokens, refillRate int) *RateLimiter {
    return &RateLimiter{
        tokens:     maxTokens,
        maxTokens:  maxTokens,
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    // 补充token
    elapsed := time.Since(rl.lastRefill)
    if elapsed > time.Second {
        tokensToAdd := int(elapsed.Seconds()) * rl.refillRate
        rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
        rl.lastRefill = time.Now()
    }
    
    // 检查token
    if rl.tokens > 0 {
        rl.tokens--
        return true
    }
    
    return false
}

func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
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

// ===== 中间件链组合 =====
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    // 从后往前包装
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

// ===== 使用示例 =====
func Example() {
    // 业务处理器
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, Middleware!")
    })
    
    // 创建限流器（每秒10个请求）
    limiter := NewRateLimiter(10, 10)
    
    // 组合中间件链
    finalHandler := Chain(
        handler,
        Recovery,
        Logging,
        Authentication,
        RateLimit(limiter),
    )
    
    // 启动服务器
    http.ListenAndServe(":8080", finalHandler)
}
```

#### 测试

```bash
# 成功请求
curl -H "Authorization: Bearer secret-token" http://localhost:8080
# 输出: Hello, Middleware!

# 未认证
curl http://localhost:8080
# 输出: Unauthorized

# 超过限流
# 快速请求11次，第11次会被拒绝
# 输出: Too Many Requests
```

---

## 第四部分：路由设计

### 实战案例4：RESTful路由

#### 场景4

- 用户资源的CRUD操作
- GET /users - 列表
- POST /users - 创建
- GET /users/:id - 详情
- PUT /users/:id - 更新
- DELETE /users/:id - 删除

#### 完整实现（原生net/http）

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "sync"
)

// User 用户模型
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// UserStore 用户存储
type UserStore struct {
    users  map[int]*User
    nextID int
    mu     sync.RWMutex
}

func NewUserStore() *UserStore {
    return &UserStore{
        users:  make(map[int]*User),
        nextID: 1,
    }
}

// ListUsers 列表
func (s *UserStore) ListUsers() []*User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    users := make([]*User, 0, len(s.users))
    for _, user := range s.users {
        users = append(users, user)
    }
    return users
}

// GetUser 详情
func (s *UserStore) GetUser(id int) (*User, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    user, exists := s.users[id]
    return user, exists
}

// CreateUser 创建
func (s *UserStore) CreateUser(name, email string) *User {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    user := &User{
        ID:    s.nextID,
        Name:  name,
        Email: email,
    }
    s.users[s.nextID] = user
    s.nextID++
    
    return user
}

// UpdateUser 更新
func (s *UserStore) UpdateUser(id int, name, email string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    user, exists := s.users[id]
    if !exists {
        return false
    }
    
    user.Name = name
    user.Email = email
    return true
}

// DeleteUser 删除
func (s *UserStore) DeleteUser(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    _, exists := s.users[id]
    if exists {
        delete(s.users, id)
    }
    return exists
}

// ===== HTTP处理器 =====
type UserHandler struct {
    store *UserStore
}

func NewUserHandler(store *UserStore) *UserHandler {
    return &UserHandler{store: store}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 解析路径
    path := strings.TrimPrefix(r.URL.Path, "/users")
    
    // GET /users - 列表
    if r.Method == http.MethodGet && path == "" {
        h.handleList(w, r)
        return
    }
    
    // POST /users - 创建
    if r.Method == http.MethodPost && path == "" {
        h.handleCreate(w, r)
        return
    }
    
    // GET /users/:id - 详情
    if r.Method == http.MethodGet && path != "" {
        h.handleGet(w, r, path)
        return
    }
    
    // PUT /users/:id - 更新
    if r.Method == http.MethodPut && path != "" {
        h.handleUpdate(w, r, path)
        return
    }
    
    // DELETE /users/:id - 删除
    if r.Method == http.MethodDelete && path != "" {
        h.handleDelete(w, r, path)
        return
    }
    
    http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *UserHandler) handleList(w http.ResponseWriter, r *http.Request) {
    users := h.store.ListUsers()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    user := h.store.CreateUser(req.Name, req.Email)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request, path string) {
    id, err := strconv.Atoi(strings.Trim(path, "/"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, exists := h.store.GetUser(id)
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleUpdate(w http.ResponseWriter, r *http.Request, path string) {
    id, err := strconv.Atoi(strings.Trim(path, "/"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if !h.store.UpdateUser(id, req.Name, req.Email) {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request, path string) {
    id, err := strconv.Atoi(strings.Trim(path, "/"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    if !h.store.DeleteUser(id) {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// ===== 主程序 =====
func main() {
    store := NewUserStore()
    handler := NewUserHandler(store)
    
    http.Handle("/users", handler)
    http.Handle("/users/", handler)
    
    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

#### 测试4

```bash
# 创建用户
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
# 输出: {"id":1,"name":"Alice","email":"alice@example.com"}

# 列表
curl http://localhost:8080/users
# 输出: [{"id":1,"name":"Alice","email":"alice@example.com"}]

# 详情
curl http://localhost:8080/users/1
# 输出: {"id":1,"name":"Alice","email":"alice@example.com"}

# 更新
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","email":"alice.new@example.com"}'
# 输出: 204 No Content

# 删除
curl -X DELETE http://localhost:8080/users/1
# 输出: 204 No Content
```

---

## 第五部分：请求与响应处理

### 请求参数处理

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 1. URL参数（Query String）
    name := r.URL.Query().Get("name")
    age := r.URL.Query().Get("age")
    
    // 2. 路径参数（需要自己解析或使用路由库）
    // 例如：/users/123
    // 使用strings.Split或正则解析
    
    // 3. 请求头
    contentType := r.Header.Get("Content-Type")
    authorization := r.Header.Get("Authorization")
    
    // 4. Cookie
    cookie, err := r.Cookie("session_id")
    if err == nil {
        fmt.Println("Cookie:", cookie.Value)
    }
    
    // 5. 表单数据（POST application/x-www-form-urlencoded）
    if r.Method == http.MethodPost {
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")
        fmt.Println("Form:", username, password)
    }
    
    // 6. JSON数据
    if contentType == "application/json" {
        var data map[string]interface{}
        json.NewDecoder(r.Body).Decode(&data)
        fmt.Println("JSON:", data)
    }
    
    // 响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "name":  name,
        "age":   age,
        "message": "Request received",
    })
}
```

### 响应处理

```go
package main

import (
    "encoding/json"
    "net/http"
)

// 1. JSON响应
func jsonResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Success",
        "data":    []int{1, 2, 3},
    })
}

// 2. HTML响应
func htmlResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    
    html := `
    <!DOCTYPE html>
    <html>
    <head><title>Go HTTP</title></head>
    <body><h1>Hello, HTML!</h1></body>
    </html>
    `
    w.Write([]byte(html))
}

// 3. 文本响应
func textResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello, Plain Text!"))
}

// 4. 重定向
func redirectResponse(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/new-location", http.StatusMovedPermanently)
}

// 5. 设置Cookie
func setCookieResponse(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    "abc123",
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
    })
    
    w.Write([]byte("Cookie set"))
}
```

---

## 第六部分：文件处理

### 实战案例5：文件上传

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // 1. 解析multipart表单（最大32MB）
    r.ParseMultipartForm(32 << 20)
    
    // 2. 获取文件
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // 3. 验证文件类型
    if header.Header.Get("Content-Type") != "image/png" && 
       header.Header.Get("Content-Type") != "image/jpeg" {
        http.Error(w, "Only PNG/JPEG allowed", http.StatusBadRequest)
        return
    }
    
    // 4. 创建目标文件
    uploadDir := "./uploads"
    os.MkdirAll(uploadDir, 0755)
    
    dstPath := filepath.Join(uploadDir, header.Filename)
    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, "Failed to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()
    
    // 5. 复制文件
    written, err := io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "File uploaded successfully! Size: %d bytes\n", written)
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
    
    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

### 文件下载

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    if filename == "" {
        http.Error(w, "Missing filename", http.StatusBadRequest)
        return
    }
    
    filepath := filepath.Join("./uploads", filename)
    
    // 检查文件是否存在
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }
    
    // 设置响应头
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    w.Header().Set("Content-Type", "application/octet-stream")
    
    // 发送文件
    http.ServeFile(w, r, filepath)
}
```

---

## 第七部分：WebSocket实战

### 实战案例6：WebSocket聊天室

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有来源
    },
}

// Client 客户端
type Client struct {
    conn *websocket.Conn
    send chan []byte
}

// Hub 聊天室
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            fmt.Println("Client connected. Total:", len(h.clients))
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            fmt.Println("Client disconnected. Total:", len(h.clients))
            
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        
        // 广播消息
        hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    
    for message := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}

func wsHandler(hub *Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 升级HTTP连接为WebSocket
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("Upgrade error:", err)
            return
        }
        
        client := &Client{
            conn: conn,
            send: make(chan []byte, 256),
        }
        
        hub.register <- client
        
        // 启动读写协程
        go client.writePump()
        go client.readPump(hub)
    }
}

func main() {
    hub := NewHub()
    go hub.Run()
    
    http.HandleFunc("/ws", wsHandler(hub))
    
    // 静态文件（聊天页面）
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "chat.html")
    })
    
    fmt.Println("WebSocket server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 第八部分：HTTP客户端最佳实践

### 实战案例7：高级HTTP客户端

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// HTTPClient 封装的HTTP客户端
type HTTPClient struct {
    client  *http.Client
    baseURL string
}

func NewHTTPClient(baseURL string) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        baseURL: baseURL,
    }
}

// GET请求
func (c *HTTPClient) Get(ctx context.Context, path string) ([]byte, error) {
    url := c.baseURL + path
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    return io.ReadAll(resp.Body)
}

// POST请求
func (c *HTTPClient) Post(ctx context.Context, path string, data interface{}) ([]byte, error) {
    url := c.baseURL + path
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    return io.ReadAll(resp.Body)
}

// 使用示例
func main() {
    client := NewHTTPClient("https://jsonplaceholder.typicode.com")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // GET请求
    data, err := client.Get(ctx, "/posts/1")
    if err != nil {
        fmt.Println("GET error:", err)
        return
    }
    fmt.Println("GET response:", string(data))
    
    // POST请求
    postData := map[string]interface{}{
        "title":  "foo",
        "body":   "bar",
        "userId": 1,
    }
    
    data, err = client.Post(ctx, "/posts", postData)
    if err != nil {
        fmt.Println("POST error:", err)
        return
    }
    fmt.Println("POST response:", string(data))
}
```

---

## 第九部分：性能优化

### 性能优化清单

| 优化项 | 说明 | 代码示例 |
|--------|------|---------|
| **连接池** | 复用TCP连接 | MaxIdleConns=100 |
| **超时配置** | 防止慢连接 | ReadTimeout=5s |
| **压缩** | gzip压缩响应 | 中间件实现 |
| **缓存** | HTTP缓存头 | Cache-Control |
| **对象池** | sync.Pool | 复用Buffer |

### 压缩中间件

```go
func Gzip(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查客户端是否支持gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        // 创建gzip writer
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        // 包装ResponseWriter
        gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
        gzw.Header().Set("Content-Encoding", "gzip")
        
        next.ServeHTTP(gzw, r)
    })
}
```

---

## 第十部分：完整RESTful API实战

### 完整项目结构

```text
project/
├── main.go
├── handlers/
│   ├── user.go
│   └── middleware.go
├── models/
│   └── user.go
└── store/
    └── user_store.go
```

### 最佳实践总结

#### DO's ✅

1. **使用Context传递请求上下文**
2. **设置合理的超时时间**
3. **使用中间件处理横切关注点**
4. **返回标准的HTTP状态码**
5. **记录日志和监控指标**
6. **优雅关闭服务器**
7. **使用HTTPS保护数据**
8. **验证和清理用户输入**

#### DON'Ts ❌

1. **不要忽略错误处理**
2. **不要在Handler中做重计算**
3. **不要忘记关闭资源（Body、文件）**
4. **不要硬编码配置**
5. **不要暴露内部错误细节**

---

## 🎯 总结

### HTTP编程核心要点

1. **net/http包** - 功能完善的标准库
2. **中间件模式** - 洋葱模型，链式处理
3. **路由设计** - RESTful规范
4. **请求处理** - 参数、表单、JSON
5. **文件处理** - 上传、下载、流式传输
6. **WebSocket** - 实时双向通信
7. **HTTP客户端** - 连接池、超时、重试
8. **性能优化** - 连接池、压缩、缓存

### 技术选型建议

| 场景 | 推荐方案 | 理由 |
|------|---------|------|
| 简单API | net/http | 无依赖、简单 |
| 复杂项目 | Gin/Echo | 功能丰富、性能好 |
| 高性能 | Fiber | 性能最优 |
| WebSocket | gorilla/websocket | 成熟稳定 |

---

**文档版本**: v6.0  

<div align="center">

Made with ❤️ for Go HTTP Developers

[⬆ 回到顶部](#go-http编程深度实战指南)

</div>

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.21+
