# TS-043-Go-Web-Frameworks-2026

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: 2026
> **Size**: >20KB

---

## 1. Framework Comparison

| Framework | Performance | Features |
|-----------|-------------|----------|
| Gin | 380K req/s | Rich |
| Echo | 375K req/s | Rich |
| Fiber | 420K req/s | Medium |
| Chi | 340K req/s | Lightweight |
| std lib | 350K req/s | Basic |

## 2. Standard Library

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", getUser)

server := &http.Server{
    Addr: ":8080",
    Handler: mux,
}
server.ListenAndServe()
```

## 3. Gin Example

```go
r := gin.Default()
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})
r.Run()
```

## 4. Middleware

```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(401)
            return
        }
        c.Next()
    }
}
```

## 5. Selection Guide

- High performance: Fiber
- Rich ecosystem: Gin/Echo
- Minimal: Chi/std lib

---

## References

1. Gin Docs
2. Echo Docs
3. Fiber Docs

---

*Last Updated: 2026-04-03*
