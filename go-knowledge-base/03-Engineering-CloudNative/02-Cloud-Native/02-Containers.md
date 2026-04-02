# 容器化

> **分类**: 工程与云原生

---

## Dockerfile

### 多阶段构建

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
```

---

## 安全实践

### 非 root 用户

```dockerfile
RUN adduser -D -g '' appuser
USER appuser
```

---

## 优雅关闭

```go
func main() {
    srv := &http.Server{Addr: ":8080"}

    go srv.ListenAndServe()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}
```
