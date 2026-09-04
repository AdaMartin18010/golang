# 健康检查包

框架级别的健康检查组件，支持 Kubernetes 健康探针（liveness、readiness、startup）。

## 📋 功能特性

- ✅ **存活探针（Liveness Probe）**: 检查应用是否存活
- ✅ **就绪探针（Readiness Probe）**: 检查应用是否准备好接收流量
- ✅ **启动探针（Startup Probe）**: 检查应用是否启动完成
- ✅ **综合健康检查**: 详细的健康状态信息
- ✅ **可扩展的检查器**: 支持注册自定义健康检查
- ✅ **超时控制**: 支持检查超时
- ✅ **定期检查**: 支持定期缓存检查结果
- ✅ **聚合检查**: 支持多个检查的聚合

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "context"
    "net/http"
    "github.com/yourusername/golang/pkg/health"
)

func main() {
    // 创建健康检查器
    checker := health.NewHealthChecker()

    // 创建 HTTP 处理器
    handler := health.NewHTTPHandler(checker)

    // 注册路由
    http.HandleFunc("/health/live", handler.LivenessHandler())
    http.HandleFunc("/health/ready", handler.ReadinessHandler())
    http.HandleFunc("/health/startup", handler.StartupHandler())
    http.HandleFunc("/health", handler.HealthHandler())

    http.ListenAndServe(":8080", nil)
}
```

### 注册自定义检查

```go
// 数据库检查
dbCheck := health.NewSimpleCheck("database", func(ctx context.Context) error {
    return db.PingContext(ctx)
})
checker.Register(dbCheck)

// Redis 检查
redisCheck := health.NewSimpleCheck("redis", func(ctx context.Context) error {
    return redisClient.Ping(ctx).Err()
})
checker.Register(redisCheck)

// 带超时的检查
timeoutCheck := health.NewTimeoutCheck(
    "external-api",
    3*time.Second,
    health.NewSimpleCheck("external-api", func(ctx context.Context) error {
        // 检查外部 API
        return nil
    }),
)
checker.Register(timeoutCheck)
```

### 定期检查（缓存结果）

```go
// 创建定期检查（每 30 秒检查一次）
periodicCheck := health.NewPeriodicCheck(
    "database",
    30*time.Second,
    health.NewSimpleCheck("database", func(ctx context.Context) error {
        return db.PingContext(ctx)
    }),
)
checker.Register(periodicCheck)
```

### 聚合检查

```go
// 创建聚合检查
aggregateCheck := health.NewAggregateCheck(
    "storage",
    dbCheck,
    redisCheck,
    s3Check,
)
checker.Register(aggregateCheck)
```

## 📚 API 参考

### HealthChecker

健康检查器，管理所有健康检查。

```go
type HealthChecker struct {
    // ...
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker() *HealthChecker

// Register 注册健康检查
func (hc *HealthChecker) Register(check Check)

// Unregister 注销健康检查
func (hc *HealthChecker) Unregister(name string) error

// Check 执行所有健康检查
func (hc *HealthChecker) Check(ctx context.Context) map[string]Result

// OverallStatus 获取整体健康状态
func (hc *HealthChecker) OverallStatus(ctx context.Context) Status
```

### Check 接口

健康检查接口。

```go
type Check interface {
    Name() string
    Check(ctx context.Context) Result
}
```

### 预定义检查类型

#### SimpleCheck

简单的健康检查，通过函数执行检查。

```go
func NewSimpleCheck(name string, checkFn func(ctx context.Context) error) *SimpleCheck
```

#### TimeoutCheck

带超时的健康检查。

```go
func NewTimeoutCheck(name string, timeout time.Duration, check Check) *TimeoutCheck
```

#### PeriodicCheck

定期健康检查，缓存结果以提高性能。

```go
func NewPeriodicCheck(name string, interval time.Duration, check Check) *PeriodicCheck
```

#### AggregateCheck

聚合多个健康检查。

```go
func NewAggregateCheck(name string, checks ...Check) *AggregateCheck
```

### HTTPHandler

HTTP 处理器，提供 Kubernetes 健康探针端点。

```go
type HTTPHandler struct {
    // ...
}

// NewHTTPHandler 创建 HTTP 处理器
func NewHTTPHandler(checker *HealthChecker) *HTTPHandler

// LivenessHandler 存活探针处理器
func (h *HTTPHandler) LivenessHandler() http.HandlerFunc

// ReadinessHandler 就绪探针处理器
func (h *HTTPHandler) ReadinessHandler() http.HandlerFunc

// StartupHandler 启动探针处理器
func (h *HTTPHandler) StartupHandler() http.HandlerFunc

// HealthHandler 综合健康检查处理器
func (h *HTTPHandler) HealthHandler() http.HandlerFunc
```

## 🔧 Kubernetes 集成

### Deployment 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
      - name: app
        image: app:latest
        ports:
        - containerPort: 8080

        # 启动探针
        startupProbe:
          httpGet:
            path: /health/startup
            port: 8080
          failureThreshold: 30
          periodSeconds: 10

        # 存活探针
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3

        # 就绪探针
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          successThreshold: 1
          failureThreshold: 3
```

## 📝 状态说明

- **healthy**: 所有检查通过
- **degraded**: 部分检查失败，但应用仍可服务
- **unhealthy**: 关键检查失败，应用不可用

## 🎯 最佳实践

1. **启动探针**: 用于慢启动应用，给应用足够的启动时间
2. **存活探针**: 检查应用是否崩溃，失败时重启容器
3. **就绪探针**: 检查应用是否准备好接收流量，失败时从 Service 中移除
4. **定期检查**: 对于耗时的检查，使用 PeriodicCheck 缓存结果
5. **超时控制**: 为所有外部依赖检查设置超时
6. **聚合检查**: 将相关的检查聚合在一起，便于管理

## 🔗 相关文档

- [Kubernetes 健康探针](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- 框架基础设施说明
