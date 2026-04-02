# TS-CL-003: Go net/http Package Architecture

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #http #web-server #client #middleware
> **权威来源**:
>
> - [Go net/http Package](https://golang.org/pkg/net/http/) - Go standard library
> - [HTTP Server Source](https://golang.org/src/net/http/server.go) - Go source code
> - [HTTP/2 in Go](https://godoc.org/golang.org/x/net/http2) - HTTP/2 implementation

---

## 1. HTTP Server Architecture

### 1.1 Core Components

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Go HTTP Server Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     http.Server                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  Addr: ":8080"                                                 │  │   │
│  │  │  Handler: http.Handler (multiplexer)                          │  │   │
│  │  │  ReadTimeout: 5s                                               │  │   │
│  │  │  WriteTimeout: 10s                                             │  │   │
│  │  │  IdleTimeout: 120s                                             │  │   │
│  │  │  MaxHeaderBytes: 1MB                                           │  │   │
│  │  │  TLSConfig: *tls.Config                                        │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                  TCP Listener (net.Listen)                     │  │   │
│  │  └───────────────────────────┬───────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                     ┌────────┴────────┐                            │   │
│  │                     ▼                 ▼                            │   │
│  │  ┌───────────────────────┐  ┌───────────────────────┐             │   │
│  │  │   Serve(net.Conn)     │  │   Serve(net.Conn)     │             │   │
│  │  │   (Goroutine 1)       │  │   (Goroutine 2)       │             │   │
│  │  │                       │  │                       │             │   │
│  │  │  ┌─────────────────┐  │  │  ┌─────────────────┐  │             │   │
│  │  │  │  bufio.Reader   │  │  │  │  bufio.Reader   │  │             │   │
│  │  │  │  (4KB buffer)   │  │  │  │  (4KB buffer)   │  │             │   │
│  │  │  └────────┬────────┘  │  │  │  └────────┬────────┘  │             │   │
│  │  │           ▼            │  │  │           ▼            │             │   │
│  │  │  ┌─────────────────┐  │  │  │  ┌─────────────────┐  │             │   │
│  │  │  │  http.Request   │  │  │  │  │  http.Request   │  │             │   │
│  │  │  │  (parsed)       │  │  │  │  │  (parsed)       │  │             │   │
│  │  │  └────────┬────────┘  │  │  │  │  └────────┬────────┘  │             │   │
│  │  │           ▼            │  │  │  │           ▼            │             │   │
│  │  │  ┌─────────────────┐  │  │  │  │  ┌─────────────────┐  │             │   │
│  │  │  │  Handler.ServeHTTP  │  │  │  │  │  Handler.ServeHTTP  │  │             │   │
│  │  │  │  (business logic)   │  │  │  │  │  (business logic)   │  │             │   │
│  │  │  └─────────────────┘  │  │  │  │  └─────────────────┘  │             │   │
│  │  └───────────────────────┘  │  └───────────────────────┘             │   │
│  │                              │                                      │   │
│  └──────────────────────────────┼──────────────────────────────────────┘   │
│                                 │                                          │
│                    ┌────────────┴────────────┐                            │
│                    ▼                         ▼                            │
│         ┌─────────────────────┐  ┌─────────────────────┐                  │
│         │  response.Write()   │  │  response.Write()   │                  │
│         │  (bufio.Writer)     │  │  (bufio.Writer)     │                  │
│         └─────────────────────┘  └─────────────────────┘                  │
│                                                                              │
│  One goroutine per connection (until HTTP/2)                                │
│  Efficient through goroutine scheduling                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Server Lifecycle

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Create handler
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/api/", apiHandler)

    // Create server
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    // Start server in goroutine
    go func() {
        log.Println("Server starting on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
```

---

## 2. Handler Interface

### 2.1 Handler Definition

```go
// Core handler interface - everything implements this
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Handler function adapter
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// Concrete implementations
func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello, World!"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
}
```

### 2.2 ServeMux - Request Multiplexer

```go
// ServeMux matches URL paths to handlers
mux := http.NewServeMux()

// Exact match
mux.HandleFunc("/", homeHandler)

// Path prefix match
mux.HandleFunc("/api/", apiHandler)  // Matches /api/, /api/users, etc.

// Host-based routing
mux.HandleFunc("api.example.com/", apiHandler)
mux.HandleFunc("www.example.com/", webHandler)

// Static file serving
mux.Handle("/static/", http.StripPrefix("/static/",
    http.FileServer(http.Dir("./static"))))

// Custom handler
mux.Handle("/custom", customHandler())
```

### 2.3 Middleware Pattern

```go
// Middleware function type
type Middleware func(http.Handler) http.Handler

// Chain middlewares
func Chain(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status code
        wrapped := &responseRecorder{ResponseWriter: w, statusCode: 200}

        next.ServeHTTP(wrapped, r)

        log.Printf("%s %s %d %v",
            r.Method,
            r.URL.Path,
            wrapped.statusCode,
            time.Since(start),
        )
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

// Recovery middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }

        // Validate token...
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Usage
handler := http.HandlerFunc(apiHandler)
chain := Chain(RecoveryMiddleware, LoggingMiddleware, AuthMiddleware)
http.Handle("/api/", chain(handler))
```

---

## 3. HTTP Client

### 3.1 Client Configuration

```go
// Default client (connection pooling enabled)
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// Custom client with timeouts
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        // Connection pool
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        MaxConnsPerHost:     100,
        IdleConnTimeout:     90 * time.Second,

        // TLS
        TLSHandshakeTimeout: 10 * time.Second,

        // Timeouts
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,

        // Compression
        DisableCompression: false,

        // HTTP/2
        ForceAttemptHTTP2: true,
    },
}

// Make request
req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
if err != nil {
    log.Fatal(err)
}

req.Header.Set("Accept", "application/json")
req.Header.Set("Authorization", "Bearer token")

resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// Read response
body, err := io.ReadAll(resp.Body)
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(body))
```

### 3.2 Request with Context

```go
// Request with timeout context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/data", nil)
if err != nil {
    log.Fatal(err)
}

resp, err := client.Do(req)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Request timeout")
    }
    log.Fatal(err)
}
defer resp.Body.Close()
```

---

## 4. Best Practices

```
Server Best Practices:
□ Always set timeouts (Read, Write, Idle)
□ Use graceful shutdown
□ Implement middleware for cross-cutting concerns
□ Use structured logging
□ Handle panics with recovery middleware
□ Validate all input
□ Use context for request cancellation
□ Set appropriate response headers

Client Best Practices:
□ Reuse http.Client (don't create per-request)
□ Set appropriate timeouts
□ Use connection pooling
□ Handle redirects properly
□ Close response body
□ Use context for cancellation
□ Implement retry logic
```

---

## 5. Checklist

```
HTTP Production Checklist:
□ Proper timeout configuration
□ Graceful shutdown implemented
□ Request logging enabled
□ Panic recovery in place
□ Authentication/authorization
□ Input validation
□ Rate limiting
□ CORS configuration
□ TLS configured
□ Compression enabled
□ Health check endpoint
□ Metrics collection
```
