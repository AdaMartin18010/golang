# HTTP/3 中间件系统 - 完整实现指南

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [HTTP/3 中间件系统 - 完整实现指南](#http3-中间件系统-完整实现指南)
  - [1. 概述](#1-概述)
  - [2. 中间件架构设计](#2-中间件架构设计)
  - [3. 五大核心中间件](#3-五大核心中间件)
  - [4. 中间件链管理](#4-中间件链管理)
  - [5. 性能优化](#5-性能优化)
  - [6. 生产实践](#6-生产实践)
  - [7. 测试与调试](#7-测试与调试)
  - [8. 最佳实践](#8-最佳实践)

---

## 1. 概述

### 1.1 HTTP/3中间件的价值

HTTP/3基于QUIC协议，提供了更好的性能和可靠性。中间件系统为HTTP/3服务提供：

```text
中间件价值体系:

┌─────────────────────────────────────┐
│         HTTP/3 请求                 │
└─────────────────────────────────────┘
            ↓
┌─────────────────────────────────────┐
│    中间件链 (Middleware Chain)      │
├─────────────────────────────────────┤
│  1. 日志中间件                      │
│     └─ 请求/响应日志记录            │
│                                     │
│  2. 限流中间件                      │
│     └─ 防止恶意请求                 │
│                                     │
│  3. CORS中间件                      │
│     └─ 跨域资源共享                 │
│                                     │
│  4. 压缩中间件                      │
│     └─ 减少传输数据量               │
│                                     │
│  5. 认证中间件                      │
│     └─ 身份验证和授权               │
└─────────────────────────────────────┘
            ↓
┌─────────────────────────────────────┐
│       业务处理器 (Handler)          │
└─────────────────────────────────────┘
            ↓
┌─────────────────────────────────────┐
│         HTTP/3 响应                 │
└─────────────────────────────────────┘
```

**核心优势**:

- ✅ 模块化设计，易于扩展
- ✅ 链式调用，灵活组合
- ✅ 性能开销低（<5%）
- ✅ 易于测试和维护

---

### 1.2 技术栈

| 组件 | 技术 | 版本要求 |
|------|------|---------|
| HTTP/3服务器 | quic-go | v0.40+ |
| 限流 | golang.org/x/time/rate | latest |
| 压缩 | compress/gzip | stdlib |
| 监控 | prometheus/client_golang | v1.17+ |

---

## 2. 中间件架构设计

### 2.1 核心接口定义

```go
// pkg/http3/middleware/middleware.go

package middleware

import (
    "net/http"
)

// Middleware 中间件接口
// 所有中间件必须实现此接口
type Middleware interface {
    // Handle 处理HTTP请求
    // next: 下一个处理器
    // 返回: 包装后的处理器
    Handle(next http.Handler) http.Handler
}

// HandlerFunc 函数式中间件
// 允许直接使用函数作为中间件
type HandlerFunc func(http.Handler) http.Handler

// Handle 实现Middleware接口
func (f HandlerFunc) Handle(next http.Handler) http.Handler {
    return f(next)
}
```

**设计理念**:

- 简单的接口，易于实现
- 支持函数式编程
- 符合Go标准库风格

---

### 2.2 中间件链设计

```go
// pkg/http3/middleware/chain.go

package middleware

import (
    "net/http"
)

// Chain 中间件链
// 按顺序执行多个中间件
type Chain struct {
    middlewares []Middleware
}

// NewChain 创建中间件链
func NewChain(middlewares ...Middleware) *Chain {
    return &Chain{
        middlewares: middlewares,
    }
}

// Append 追加中间件
func (c *Chain) Append(m Middleware) *Chain {
    c.middlewares = append(c.middlewares, m)
    return c
}

// Extend 扩展中间件链
func (c *Chain) Extend(chain *Chain) *Chain {
    c.middlewares = append(c.middlewares, chain.middlewares...)
    return c
}

// Then 应用中间件链到处理器
// 中间件按LIFO（后进先出）顺序执行
func (c *Chain) Then(h http.Handler) http.Handler {
    // 从后向前应用中间件
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        h = c.middlewares[i].Handle(h)
    }
    return h
}

// ThenFunc 应用到HandlerFunc
func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
    return c.Then(fn)
}
```

**执行流程**:

```text
请求流向:
Request → M1 → M2 → M3 → Handler → M3 → M2 → M1 → Response

中间件执行顺序:
1. M1.Handle() 开始
2. M2.Handle() 开始
3. M3.Handle() 开始
4. Handler 处理业务逻辑
5. M3.Handle() 完成
6. M2.Handle() 完成
7. M1.Handle() 完成
```

---

### 2.3 上下文管理

```go
// pkg/http3/middleware/Context.go

package middleware

import (
    "Context"
    "net/http"
)

type contextKey string

const (
    RequestIDKey  contextKey = "request_id"
    UserIDKey     contextKey = "user_id"
    StartTimeKey  contextKey = "start_time"
)

// GetRequestID 从上下文获取请求ID
func GetRequestID(r *http.Request) string {
    if id, ok := r.Context().Value(RequestIDKey).(string); ok {
        return id
    }
    return ""
}

// SetRequestID 设置请求ID到上下文
func SetRequestID(r *http.Request, id string) *http.Request {
    ctx := Context.WithValue(r.Context(), RequestIDKey, id)
    return r.WithContext(ctx)
}

// GetUserID 从上下文获取用户ID
func GetUserID(r *http.Request) string {
    if id, ok := r.Context().Value(UserIDKey).(string); ok {
        return id
    }
    return ""
}

// SetUserID 设置用户ID到上下文
func SetUserID(r *http.Request, id string) *http.Request {
    ctx := Context.WithValue(r.Context(), UserIDKey, id)
    return r.WithContext(ctx)
}
```

---

## 3. 五大核心中间件

### 3.1 日志中间件

#### 3.1.1 设计目标

**功能需求**:

- ✅ 记录请求方法、路径、状态码
- ✅ 记录请求处理时间
- ✅ 记录HTTP协议版本（HTTP/3）
- ✅ 可配置日志格式
- ✅ 支持结构化日志

#### 3.1.2 完整实现

```go
// pkg/http3/middleware/logging.go

package middleware

import (
    "bufio"
    "errors"
    "log"
    "net"
    "net/http"
    "time"
)

// LoggingMiddleware HTTP/3日志中间件
type LoggingMiddleware struct {
    logger *log.Logger
    config LoggingConfig
}

// LoggingConfig 日志配置
type LoggingConfig struct {
    // IncludeHeaders 是否记录请求头
    IncludeHeaders bool

    // IncludeQuery 是否记录查询参数
    IncludeQuery bool

    // IncludeBody 是否记录请求体（谨慎使用）
    IncludeBody bool

    // SlowRequestThreshold 慢请求阈值
    SlowRequestThreshold time.Duration
}

// DefaultLoggingConfig 默认日志配置
var DefaultLoggingConfig = LoggingConfig{
    IncludeHeaders:       false,
    IncludeQuery:         true,
    IncludeBody:          false,
    SlowRequestThreshold: 1 * time.Second,
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
    return NewLoggingMiddlewareWithConfig(logger, DefaultLoggingConfig)
}

// NewLoggingMiddlewareWithConfig 创建带配置的日志中间件
func NewLoggingMiddlewareWithConfig(logger *log.Logger, config LoggingConfig) *LoggingMiddleware {
    if logger == nil {
        logger = log.Default()
    }

    return &LoggingMiddleware{
        logger: logger,
        config: config,
    }
}

// Handle 实现Middleware接口
func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装ResponseWriter以捕获状态码和字节数
        wrapped := &loggingResponseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
            bytesWritten:   0,
        }

        // 调用下一个处理器
        next.ServeHTTP(wrapped, r)

        // 计算处理时间
        duration := time.Since(start)

        // 构造日志消息
        logMsg := m.buildLogMessage(r, wrapped, duration)

        // 记录日志
        if duration > m.config.SlowRequestThreshold {
            m.logger.Printf("[SLOW] %s", logMsg)
        } else {
            m.logger.Printf("[HTTP/3] %s", logMsg)
        }
    })
}

// buildLogMessage 构造日志消息
func (m *LoggingMiddleware) buildLogMessage(
    r *http.Request,
    w *loggingResponseWriter,
    duration time.Duration,
) string {
    msg := fmt.Sprintf(
        "%s %s - Status: %d - Duration: %v - Proto: %s - Size: %d bytes",
        r.Method,
        r.URL.Path,
        w.statusCode,
        duration,
        r.Proto, // "HTTP/3"
        w.bytesWritten,
    )

    if m.config.IncludeQuery && r.URL.RawQuery != "" {
        msg += fmt.Sprintf(" - Query: %s", r.URL.RawQuery)
    }

    if m.config.IncludeHeaders {
        msg += fmt.Sprintf(" - Headers: %v", r.Header)
    }

    return msg
}

// loggingResponseWriter 包装ResponseWriter以捕获状态码和字节数
type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.bytesWritten += n
    return n, err
}

// Hijack 实现http.Hijacker接口（用于WebSocket）
func (w *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    h, ok := w.ResponseWriter.(http.Hijacker)
    if !ok {
        return nil, nil, errors.New("hijack not supported")
    }
    return h.Hijack()
}

// Flush 实现http.Flusher接口
func (w *loggingResponseWriter) Flush() {
    if f, ok := w.ResponseWriter.(http.Flusher); ok {
        f.Flush()
    }
}
```

#### 3.1.3 使用示例

```go
// 基础使用
logger := log.New(os.Stdout, "[HTTP3] ", log.LstdFlags)
loggingMW := middleware.NewLoggingMiddleware(logger)

// 自定义配置
config := middleware.LoggingConfig{
    IncludeQuery:         true,
    IncludeHeaders:       true,
    SlowRequestThreshold: 500 * time.Millisecond,
}
loggingMW := middleware.NewLoggingMiddlewareWithConfig(logger, config)

// 应用到处理器
handler := loggingMW.Handle(yourHandler)
```

**输出示例**:

```text
[HTTP/3] 2025/10/24 10:30:45 GET /api/users - Status: 200 - Duration: 15ms - Proto: HTTP/3 - Size: 1024 bytes
[HTTP/3] 2025/10/24 10:30:46 POST /api/orders - Status: 201 - Duration: 45ms - Proto: HTTP/3 - Size: 512 bytes
[SLOW] 2025/10/24 10:30:47 GET /api/reports - Status: 200 - Duration: 1.2s - Proto: HTTP/3 - Size: 10240 bytes
```

---

### 3.2 限流中间件

#### 3.2.1 设计目标

**功能需求**:

- ✅ 基于IP的限流
- ✅ 令牌桶算法
- ✅ 可配置速率和突发量
- ✅ 自动清理过期限流器
- ✅ 支持自定义键提取

#### 3.2.2 完整实现

```go
// pkg/http3/middleware/ratelimit.go

package middleware

import (
    "fmt"
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

// RateLimitMiddleware HTTP/3限流中间件
type RateLimitMiddleware struct {
    mu           sync.RWMutex
    limiters     map[string]*rateLimiterEntry
    rate         rate.Limit
    burst        int
    keyExtractor KeyExtractor
    cleanupTicker *time.Ticker
    stopCleanup  Channel struct{}
}

// rateLimiterEntry 限流器条目
type rateLimiterEntry struct {
    limiter    *rate.Limiter
    lastAccess time.Time
}

// KeyExtractor 键提取函数
type KeyExtractor func(*http.Request) string

// DefaultKeyExtractor 默认键提取器（基于IP）
func DefaultKeyExtractor(r *http.Request) string {
    return getClientIP(r)
}

// NewRateLimitMiddleware 创建限流中间件
// rps: 每秒请求数
// burst: 突发请求数
func NewRateLimitMiddleware(rps int, burst int) *RateLimitMiddleware {
    return NewRateLimitMiddlewareWithExtractor(rps, burst, DefaultKeyExtractor)
}

// NewRateLimitMiddlewareWithExtractor 创建带自定义键提取器的限流中间件
func NewRateLimitMiddlewareWithExtractor(
    rps int,
    burst int,
    extractor KeyExtractor,
) *RateLimitMiddleware {
    m := &RateLimitMiddleware{
        limiters:      make(map[string]*rateLimiterEntry),
        rate:          rate.Limit(rps),
        burst:         burst,
        keyExtractor:  extractor,
        cleanupTicker: time.NewTicker(1 * time.Minute),
        stopCleanup:   make(Channel struct{}),
    }

    // 启动清理goroutine
    go m.cleanupExpiredLimiters()

    return m
}

// Handle 实现Middleware接口
func (m *RateLimitMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 提取键
        key := m.keyExtractor(r)

        // 获取或创建限流器
        limiter := m.getLimiter(key)

        // 检查是否超过限流
        if !limiter.Allow() {
            // 设置Retry-After头
            w.Header().Set("Retry-After", "1")
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        // 调用下一个处理器
        next.ServeHTTP(w, r)
    })
}

// getLimiter 获取或创建限流器
func (m *RateLimitMiddleware) getLimiter(key string) *rate.Limiter {
    m.mu.RLock()
    entry, exists := m.limiters[key]
    m.mu.RUnlock()

    if exists {
        // 更新最后访问时间
        m.mu.Lock()
        entry.lastAccess = time.Now()
        m.mu.Unlock()

        return entry.limiter
    }

    // 创建新限流器
    m.mu.Lock()
    defer m.mu.Unlock()

    // 双重检查
    if entry, exists := m.limiters[key]; exists {
        return entry.limiter
    }

    limiter := rate.NewLimiter(m.rate, m.burst)
    m.limiters[key] = &rateLimiterEntry{
        limiter:    limiter,
        lastAccess: time.Now(),
    }

    return limiter
}

// cleanupExpiredLimiters 清理过期的限流器
func (m *RateLimitMiddleware) cleanupExpiredLimiters() {
    for {
        select {
        case <-m.cleanupTicker.C:
            m.cleanup()
        case <-m.stopCleanup:
            return
        }
    }
}

// cleanup 执行清理
func (m *RateLimitMiddleware) cleanup() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    threshold := 5 * time.Minute

    for key, entry := range m.limiters {
        if now.Sub(entry.lastAccess) > threshold {
            delete(m.limiters, key)
        }
    }
}

// Close 关闭中间件
func (m *RateLimitMiddleware) Close() {
    m.cleanupTicker.Stop()
    close(m.stopCleanup)
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
    // 优先从X-Forwarded-For获取
    if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
        return strings.Split(ip, ",")[0]
    }

    // 从X-Real-IP获取
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }

    // 直接从RemoteAddr获取
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

#### 3.2.3 高级用法

```go
// 基于用户ID的限流
userKeyExtractor := func(r *http.Request) string {
    // 从认证token提取用户ID
    userID := extractUserID(r)
    return fmt.Sprintf("user:%s", userID)
}

rateLimitMW := middleware.NewRateLimitMiddlewareWithExtractor(
    100,  // 100 rps per user
    10,   // burst 10
    userKeyExtractor,
)

// 基于API路径的限流
pathKeyExtractor := func(r *http.Request) string {
    ip := getClientIP(r)
    path := r.URL.Path
    return fmt.Sprintf("%s:%s", ip, path)
}

rateLimitMW := middleware.NewRateLimitMiddlewareWithExtractor(
    50,   // 50 rps per IP per path
    5,    // burst 5
    pathKeyExtractor,
)
```

---

### 3.3 CORS中间件

#### 3.3.1 设计目标

**功能需求**:

- ✅ 配置允许的Origin
- ✅ 配置允许的HTTP方法
- ✅ 配置允许的请求头
- ✅ 支持预检请求（OPTIONS）
- ✅ 配置凭证支持
- ✅ 配置缓存时间

#### 3.3.2 完整实现

```go
// pkg/http3/middleware/cors.go

package middleware

import (
    "net/http"
    "strconv"
    "strings"
)

// CORSConfig CORS配置
type CORSConfig struct {
    // AllowOrigins 允许的Origin列表
    // 使用"*"允许所有Origin
    AllowOrigins []string

    // AllowMethods 允许的HTTP方法
    AllowMethods []string

    // AllowHeaders 允许的请求头
    AllowHeaders []string

    // ExposeHeaders 暴露的响应头
    ExposeHeaders []string

    // AllowCredentials 是否允许凭证
    AllowCredentials bool

    // MaxAge 预检请求缓存时间（秒）
    MaxAge int
}

// DefaultCORSConfig 默认CORS配置
var DefaultCORSConfig = CORSConfig{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowHeaders:     []string{"Accept", "Content-Type", "Authorization", "X-Request-ID"},
    ExposeHeaders:    []string{},
    AllowCredentials: false,
    MaxAge:           3600, // 1小时
}

// CORSMiddleware HTTP/3 CORS中间件
type CORSMiddleware struct {
    config CORSConfig
}

// NewCORSMiddleware 创建CORS中间件
func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
    // 设置默认值
    if len(config.AllowMethods) == 0 {
        config.AllowMethods = DefaultCORSConfig.AllowMethods
    }

    if len(config.AllowHeaders) == 0 {
        config.AllowHeaders = DefaultCORSConfig.AllowHeaders
    }

    if config.MaxAge == 0 {
        config.MaxAge = DefaultCORSConfig.MaxAge
    }

    return &CORSMiddleware{
        config: config,
    }
}

// NewDefaultCORSMiddleware 创建默认CORS中间件
func NewDefaultCORSMiddleware() *CORSMiddleware {
    return NewCORSMiddleware(DefaultCORSConfig)
}

// Handle 实现Middleware接口
func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        // 检查是否允许该Origin
        if m.isOriginAllowed(origin) {
            // 设置CORS响应头
            w.Header().Set("Access-Control-Allow-Origin", m.getAllowOriginHeader(origin))

            if m.config.AllowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }

            w.Header().Set("Access-Control-Allow-Methods",
                strings.Join(m.config.AllowMethods, ", "))

            w.Header().Set("Access-Control-Allow-Headers",
                strings.Join(m.config.AllowHeaders, ", "))

            if len(m.config.ExposeHeaders) > 0 {
                w.Header().Set("Access-Control-Expose-Headers",
                    strings.Join(m.config.ExposeHeaders, ", "))
            }

            w.Header().Set("Access-Control-Max-Age",
                strconv.Itoa(m.config.MaxAge))
        }

        // 处理预检请求
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        // 调用下一个处理器
        next.ServeHTTP(w, r)
    })
}

// isOriginAllowed 检查Origin是否允许
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
    if origin == "" {
        return false
    }

    if len(m.config.AllowOrigins) == 0 {
        return false
    }

    for _, allowed := range m.config.AllowOrigins {
        if allowed == "*" {
            return true
        }

        if allowed == origin {
            return true
        }

        // 支持通配符匹配（简单实现）
        if strings.HasPrefix(allowed, "*.") {
            domain := strings.TrimPrefix(allowed, "*")
            if strings.HasSuffix(origin, domain) {
                return true
            }
        }
    }

    return false
}

// getAllowOriginHeader 获取Access-Control-Allow-Origin头的值
func (m *CORSMiddleware) getAllowOriginHeader(origin string) string {
    // 如果允许凭证，不能使用通配符
    if m.config.AllowCredentials {
        return origin
    }

    // 如果配置了通配符，直接返回通配符
    for _, allowed := range m.config.AllowOrigins {
        if allowed == "*" {
            return "*"
        }
    }

    return origin
}
```

#### 3.3.3 使用示例

```go
// 允许特定Origin
config := middleware.CORSConfig{
    AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           86400, // 24小时
}
corsMW := middleware.NewCORSMiddleware(config)

// 允许所有Origin（开发环境）
devConfig := middleware.CORSConfig{
    AllowOrigins: []string{"*"},
}
corsMW := middleware.NewCORSMiddleware(devConfig)

// 通配符域名匹配
wildcardConfig := middleware.CORSConfig{
    AllowOrigins: []string{"*.example.com"},
}
corsMW := middleware.NewCORSMiddleware(wildcardConfig)
```

---

### 3.4 压缩中间件

#### 3.4.1 设计目标

**功能需求**:

- ✅ Gzip压缩支持
- ✅ 可配置压缩级别
- ✅ 自动检测客户端支持
- ✅ 排除不可压缩内容
- ✅ 性能优化

#### 3.4.2 完整实现

```go
// pkg/http3/middleware/compression.go

package middleware

import (
    "compress/gzip"
    "io"
    "net/http"
    "strings"
    "sync"
)

// CompressionMiddleware HTTP/3压缩中间件
type CompressionMiddleware struct {
    level      int
    minLength  int
    pool       *sync.Pool
    shouldSkip func(*http.Request) bool
}

// CompressionConfig 压缩配置
type CompressionConfig struct {
    // Level 压缩级别 (1-9)
    Level int

    // MinLength 最小压缩长度（字节）
    MinLength int

    // ShouldSkip 是否跳过压缩的判断函数
    ShouldSkip func(*http.Request) bool
}

// DefaultCompressionConfig 默认压缩配置
var DefaultCompressionConfig = CompressionConfig{
    Level:     gzip.DefaultCompression,
    MinLength: 1024, // 1KB
    ShouldSkip: func(r *http.Request) bool {
        // 默认不跳过
        return false
    },
}

// NewCompressionMiddleware 创建压缩中间件
func NewCompressionMiddleware(level int) *CompressionMiddleware {
    return NewCompressionMiddlewareWithConfig(CompressionConfig{
        Level:     level,
        MinLength: DefaultCompressionConfig.MinLength,
    })
}

// NewCompressionMiddlewareWithConfig 创建带配置的压缩中间件
func NewCompressionMiddlewareWithConfig(config CompressionConfig) *CompressionMiddleware {
    if config.Level < gzip.BestSpeed || config.Level > gzip.BestCompression {
        config.Level = gzip.DefaultCompression
    }

    if config.MinLength == 0 {
        config.MinLength = DefaultCompressionConfig.MinLength
    }

    if config.ShouldSkip == nil {
        config.ShouldSkip = DefaultCompressionConfig.ShouldSkip
    }

    return &CompressionMiddleware{
        level:      config.Level,
        minLength:  config.MinLength,
        shouldSkip: config.ShouldSkip,
        pool: &sync.Pool{
            New: func() interface{} {
                w, _ := gzip.NewWriterLevel(nil, config.Level)
                return w
            },
        },
    }
}

// Handle 实现Middleware接口
func (m *CompressionMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查是否应该跳过压缩
        if m.shouldSkip(r) {
            next.ServeHTTP(w, r)
            return
        }

        // 检查客户端是否支持gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        // 从池中获取gzip writer
        gz := m.pool.Get().(*gzip.Writer)
        defer m.pool.Put(gz)

        gz.Reset(w)
        defer gz.Close()

        // 包装ResponseWriter
        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Del("Content-Length") // 压缩后长度会变化

        gzw := &gzipResponseWriter{
            ResponseWriter: w,
            Writer:         gz,
        }

        // 调用下一个处理器
        next.ServeHTTP(gzw, r)
    })
}

// gzipResponseWriter 包装ResponseWriter支持压缩
type gzipResponseWriter struct {
    http.ResponseWriter
    io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

// WriteHeader 重写WriteHeader以处理压缩相关头
func (w *gzipResponseWriter) WriteHeader(statusCode int) {
    // 删除可能冲突的头
    w.Header().Del("Content-Length")
    w.ResponseWriter.WriteHeader(statusCode)
}

// Flush 实现http.Flusher接口
func (w *gzipResponseWriter) Flush() {
    if gz, ok := w.Writer.(*gzip.Writer); ok {
        gz.Flush()
    }

    if f, ok := w.ResponseWriter.(http.Flusher); ok {
        f.Flush()
    }
}
```

#### 3.4.3 高级配置

```go
// 排除特定内容类型
config := middleware.CompressionConfig{
    Level:     gzip.BestSpeed,
    MinLength: 1024,
    ShouldSkip: func(r *http.Request) bool {
        // 跳过图片、视频等已压缩的内容
        contentType := r.Header.Get("Content-Type")
        return strings.HasPrefix(contentType, "image/") ||
               strings.HasPrefix(contentType, "video/") ||
               strings.HasPrefix(contentType, "audio/")
    },
}
compressionMW := middleware.NewCompressionMiddlewareWithConfig(config)

// 排除小文件
config := middleware.CompressionConfig{
    MinLength: 2048, // 2KB
}
```

---

### 3.5 认证中间件

#### 3.5.1 设计目标

**功能需求**:

- ✅ Bearer Token认证
- ✅ Basic认证
- ✅ JWT验证
- ✅ 可自定义验证逻辑
- ✅ 白名单路径

#### 3.5.2 完整实现

```go
// pkg/http3/middleware/auth.go

package middleware

import (
    "crypto/subtle"
    "encoding/base64"
    "net/http"
    "strings"
)

// AuthMiddleware HTTP/3认证中间件
type AuthMiddleware struct {
    scheme       string
    validateFunc ValidateFunc
    skipPaths    map[string]bool
}

// ValidateFunc 验证函数
type ValidateFunc func(token string) (userID string, ok bool)

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(scheme string, validateFunc ValidateFunc) *AuthMiddleware {
    return &AuthMiddleware{
        scheme:       scheme,
        validateFunc: validateFunc,
        skipPaths:    make(map[string]bool),
    }
}

// SkipPath 添加跳过认证的路径
func (m *AuthMiddleware) SkipPath(path string) *AuthMiddleware {
    m.skipPaths[path] = true
    return m
}

// Handle 实现Middleware接口
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查是否跳过认证
        if m.skipPaths[r.URL.Path] {
            next.ServeHTTP(w, r)
            return
        }

        // 获取Authorization头
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            m.unauthorized(w, "Missing authorization header")
            return
        }

        // 解析认证方案和token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 {
            m.unauthorized(w, "Invalid authorization header format")
            return
        }

        scheme, token := parts[0], parts[1]

        // 检查认证方案
        if !strings.EqualFold(scheme, m.scheme) {
            m.unauthorized(w, fmt.Sprintf("Expected %s authentication", m.scheme))
            return
        }

        // 验证token
        userID, ok := m.validateFunc(token)
        if !ok {
            m.unauthorized(w, "Invalid token")
            return
        }

        // 将用户ID设置到上下文
        r = SetUserID(r, userID)

        // 调用下一个处理器
        next.ServeHTTP(w, r)
    })
}

// unauthorized 返回401错误
func (m *AuthMiddleware) unauthorized(w http.ResponseWriter, message string) {
    w.Header().Set("WWW-Authenticate", fmt.Sprintf("%s realm=\"Restricted\"", m.scheme))
    http.Error(w, message, http.StatusUnauthorized)
}

// NewBearerAuthMiddleware 创建Bearer认证中间件
func NewBearerAuthMiddleware(validateFunc ValidateFunc) *AuthMiddleware {
    return NewAuthMiddleware("Bearer", validateFunc)
}

// NewBasicAuthMiddleware 创建Basic认证中间件
func NewBasicAuthMiddleware(users map[string]string) *AuthMiddleware {
    validateFunc := func(token string) (string, bool) {
        // Base64解码
        decoded, err := base64.StdEncoding.DecodeString(token)
        if err != nil {
            return "", false
        }

        // 解析username:password
        parts := strings.SplitN(string(decoded), ":", 2)
        if len(parts) != 2 {
            return "", false
        }

        username, password := parts[0], parts[1]

        // 验证用户名和密码
        expectedPassword, ok := users[username]
        if !ok {
            return "", false
        }

        // 使用常量时间比较防止时序攻击
        if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
            return "", false
        }

        return username, true
    }

    return NewAuthMiddleware("Basic", validateFunc)
}

// NewJWTAuthMiddleware 创建JWT认证中间件
func NewJWTAuthMiddleware(secretKey []byte) *AuthMiddleware {
    validateFunc := func(token string) (string, bool) {
        // JWT验证逻辑
        // 这里需要使用JWT库，如github.com/golang-jwt/jwt

        // 简化示例
        claims, err := parseAndValidateJWT(token, secretKey)
        if err != nil {
            return "", false
        }

        userID, ok := claims["user_id"].(string)
        return userID, ok
    }

    return NewAuthMiddleware("Bearer", validateFunc)
}

// parseAndValidateJWT JWT解析和验证（示例）
func parseAndValidateJWT(tokenString string, secretKey []byte) (map[string]interface{}, error) {
    // 实际实现应使用JWT库
    // 这里只是占位符
    return nil, fmt.Errorf("not implemented")
}
```

#### 3.5.3 使用示例

```go
// Bearer Token认证
validateToken := func(token string) (string, bool) {
    // 从数据库或缓存验证token
    userID, err := tokenStore.Validate(token)
    return userID, err == nil
}

authMW := middleware.NewBearerAuthMiddleware(validateToken).
    SkipPath("/health").
    SkipPath("/metrics")

// Basic认证
users := map[string]string{
    "admin": "secret123",
    "user":  "password456",
}

basicAuthMW := middleware.NewBasicAuthMiddleware(users)

// JWT认证
jwtAuthMW := middleware.NewJWTAuthMiddleware([]byte("your-secret-key"))
```

---

## 4. 中间件链管理

### 4.1 链式组合

```go
// 创建全局中间件链
globalChain := middleware.NewChain(
    middleware.NewLoggingMiddleware(logger),
    middleware.NewRateLimitMiddleware(100, 10),
)

// 创建API特定中间件链
apiChain := middleware.NewChain(
    middleware.NewCORSMiddleware(corsConfig),
    middleware.NewCompressionMiddleware(gzip.BestSpeed),
    middleware.NewBearerAuthMiddleware(validateToken),
)

// 组合中间件链
fullChain := globalChain.Extend(apiChain)

// 应用到处理器
handler := fullChain.Then(yourHandler)
```

### 4.2 条件中间件

```go
// 条件应用中间件
func ConditionalMiddleware(condition func(*http.Request) bool, m middleware.Middleware) middleware.Middleware {
    return middleware.HandlerFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if condition(r) {
                m.Handle(next).ServeHTTP(w, r)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    })
}

// 使用示例
authMW := ConditionalMiddleware(
    func(r *http.Request) bool {
        return strings.HasPrefix(r.URL.Path, "/api/")
    },
    middleware.NewBearerAuthMiddleware(validateToken),
)
```

---

## 5. 性能优化

### 5.1 性能基准

```go
// benchmarks/middleware_bench_test.go

package benchmarks

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "your-project/pkg/http3/middleware"
)

func BenchmarkLoggingMiddleware(b *testing.B) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    logger := log.New(io.Discard, "", 0)
    loggingMW := middleware.NewLoggingMiddleware(logger)
    wrappedHandler := loggingMW.Handle(handler)

    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        wrappedHandler.ServeHTTP(rec, req)
    }
}

func BenchmarkRateLimitMiddleware(b *testing.B) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    rateLimitMW := middleware.NewRateLimitMiddleware(1000, 10)
    defer rateLimitMW.Close()
    wrappedHandler := rateLimitMW.Handle(handler)

    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        wrappedHandler.ServeHTTP(rec, req)
    }
}

func BenchmarkCompressionMiddleware(b *testing.B) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(strings.Repeat("Hello, World! ", 100)))
    })

    compressionMW := middleware.NewCompressionMiddleware(gzip.BestSpeed)
    wrappedHandler := compressionMW.Handle(handler)

    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Accept-Encoding", "gzip")
    rec := httptest.NewRecorder()

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        wrappedHandler.ServeHTTP(rec, req)
    }
}
```

**性能指标（目标）**:

| 中间件 | 延迟开销 | 内存分配 | 目标 |
|--------|---------|---------|------|
| 日志 | <1ms | 1-2次 | <2% |
| 限流 | <0.5ms | 0-1次 | <1% |
| CORS | <0.1ms | 0次 | <0.5% |
| 压缩 | <5ms | 2-3次 | <3% |
| 认证 | <1ms | 1-2次 | <2% |
| **总计** | **<7.6ms** | **4-9次** | **<8.5%** |

---

### 5.2 优化技巧

#### 对象池化

```go
// 使用sync.Pool复用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        buf := bufferPool.Get().(*bytes.Buffer)
        buf.Reset()
        defer bufferPool.Put(buf)

        // 使用buf...
    })
}
```

#### 避免不必要的分配

```go
// ❌ 每次都分配新map
headers := make(map[string]string)

// ✅ 复用预分配的map
var headerPool = sync.Pool{
    New: func() interface{} {
        return make(map[string]string, 10)
    },
}
```

---

## 6. 生产实践

### 6.1 完整示例

```go
// cmd/server/main.go

package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "github.com/quic-go/quic-go/http3"
    "your-project/pkg/http3/middleware"
)

func main() {
    // 创建日志器
    logger := log.New(os.Stdout, "[HTTP3] ", log.LstdFlags|log.Lshortfile)

    // 配置中间件
    loggingMW := middleware.NewLoggingMiddlewareWithConfig(logger, middleware.LoggingConfig{
        IncludeQuery:         true,
        SlowRequestThreshold: 500 * time.Millisecond,
    })

    rateLimitMW := middleware.NewRateLimitMiddleware(100, 20)
    defer rateLimitMW.Close()

    corsMW := middleware.NewCORSMiddleware(middleware.CORSConfig{
        AllowOrigins:     []string{"https://app.example.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowCredentials: true,
    })

    compressionMW := middleware.NewCompressionMiddleware(gzip.BestSpeed)

    authMW := middleware.NewBearerAuthMiddleware(validateToken).
        SkipPath("/health").
        SkipPath("/metrics")

    // 创建中间件链
    chain := middleware.NewChain(
        loggingMW,
        rateLimitMW,
        corsMW,
        compressionMW,
        authMW,
    )

    // 创建路由
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/api/users", usersHandler)
    mux.HandleFunc("/api/orders", ordersHandler)

    // 应用中间件链
    handler := chain.Then(mux)

    // 启动HTTP/3服务器
    server := &http3.Server{
        Addr:    ":443",
        Handler: handler,
    }

    logger.Println("HTTP/3 server starting on :443")
    if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
        logger.Fatal(err)
    }
}

func validateToken(token string) (string, bool) {
    // 实际的token验证逻辑
    return "user123", token != ""
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"user_id": "%s", "message": "Hello!"}`, userID)
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
    // 处理订单逻辑
}
```

---

### 6.2 监控指标

```go
// pkg/http3/middleware/metrics.go

package middleware

import (
    "net/http"
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsMiddleware 监控指标中间件
type MetricsMiddleware struct {
    requests  *prometheus.CounterVec
    duration  *prometheus.HistogramVec
    inFlight  prometheus.Gauge
}

// NewMetricsMiddleware 创建监控中间件
func NewMetricsMiddleware(namespace string) *MetricsMiddleware {
    return &MetricsMiddleware{
        requests: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:      "http_requests_total",
                Help:      "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),

        duration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Name:      "http_request_duration_seconds",
                Help:      "HTTP request duration in seconds",
                Buckets:   prometheus.DefBuckets,
            },
            []string{"method", "path"},
        ),

        inFlight: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "http_requests_in_flight",
                Help:      "Current number of HTTP requests being served",
            },
        ),
    }
}

// Handle 实现Middleware接口
func (m *MetricsMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        m.inFlight.Inc()
        defer m.inFlight.Dec()

        wrapped := &metricsResponseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }

        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()

        m.requests.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(wrapped.statusCode),
        ).Inc()

        m.duration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    })
}

type metricsResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
```

---

## 7. 测试与调试

### 7.1 单元测试

```go
// pkg/http3/middleware/logging_test.go

package middleware_test

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "your-project/pkg/http3/middleware"
)

func TestLoggingMiddleware(t *testing.T) {
    // 准备
    var buf bytes.Buffer
    logger := log.New(&buf, "", 0)

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    loggingMW := middleware.NewLoggingMiddleware(logger)
    wrappedHandler := loggingMW.Handle(handler)

    req := httptest.NewRequest("GET", "/test", nil)
    rec := httptest.NewRecorder()

    // 执行
    wrappedHandler.ServeHTTP(rec, req)

    // 验证
    if rec.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rec.Code)
    }

    logOutput := buf.String()
    if !strings.Contains(logOutput, "GET /test") {
        t.Errorf("Log should contain request method and path")
    }

    if !strings.Contains(logOutput, "Status: 200") {
        t.Errorf("Log should contain status code")
    }
}
```

---

## 8. 最佳实践

### 8.1 中间件顺序

```text
推荐的中间件顺序（从外到内）:

1. 日志中间件 - 最外层，记录所有请求
2. 指标中间件 - 收集性能数据
3. 恢复中间件 - 捕获panic
4. 限流中间件 - 防止滥用
5. CORS中间件 - 跨域支持
6. 压缩中间件 - 减少传输
7. 认证中间件 - 验证身份
8. 授权中间件 - 检查权限
9. 业务处理器 - 实际业务逻辑
```

### 8.2 性能建议
