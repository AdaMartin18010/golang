# 微服务架构

> **分类**: 工程与云原生

---

## 架构原则

```
┌─────────────┐
│ API Gateway │
└──────┬──────┘
       │
  ┌────┴────┐
  ↓         ↓
┌──────┐ ┌──────┐
│User  │ │Order │
│Svc   │ │Svc   │
└───┬──┘ └───┬──┘
    └────┬───┘
         ↓
    ┌─────────┐
    │Database │
    └─────────┘
```

---

## Go 微服务框架

### Gin + gRPC

```go
// HTTP 服务
r := gin.Default()
r.GET("/users/:id", getUser)
r.Run(":8080")

// gRPC 服务
s := grpc.NewServer()
pb.RegisterUserServiceServer(s, &userServer{})
s.Serve(lis)
```

---

## 服务治理

### 熔断与限流

```go
breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "user-service",
    MaxRequests: 3,
    Timeout:     30 * time.Second,
})

result, err := breaker.Execute(func() (interface{}, error) {
    return userClient.GetUser(ctx, req)
})
```

---

## 可观测性

### 日志

```go
logger.Info("user login",
    zap.String("user_id", userID),
    zap.Duration("latency", time.Since(start)),
)
```

### 指标

```go
var requestDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "http_request_duration_seconds",
    },
    []string{"method", "endpoint"},
)
```
