# 中间件模式

> **分类**: 成熟应用领域

---

## Gin 中间件链

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        log.Printf("%s %s %v", c.Request.Method, c.Request.URL, latency)
    }
}

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                c.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
            }
        }()
        c.Next()
    }
}

r.Use(Logger(), Recovery())
```

---

## 标准库中间件

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL, time.Since(start))
    })
}

mux := http.NewServeMux()
mux.HandleFunc("/", handler)

wrapped := loggingMiddleware(mux)
http.ListenAndServe(":8080", wrapped)
```

---

## 中间件组合

```go
type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

handler := Chain(Logger, Recovery, Auth)(mux)
```
