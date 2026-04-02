# 健康检查 (Health Checks)

> **分类**: 工程与云原生  
> **标签**: #health #kubernetes #monitoring

---

## 健康检查类型

### Liveness（存活检查）

```go
// 应用是否还在运行
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
    // 简单检查：只要能响应就 alive
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("alive"))
}
```

### Readiness（就绪检查）

```go
// 应用是否准备好接收流量
type HealthChecker struct {
    checks map[string]HealthCheck
}

type HealthCheck func(ctx context.Context) error

func (h *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    for name, check := range h.checks {
        if err := check(ctx); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "status": "not ready",
                "check":  name,
                "error":  err.Error(),
            })
            return
        }
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ready"))
}

// 注册检查
func NewHealthChecker() *HealthChecker {
    hc := &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
    
    // 数据库检查
    hc.checks["database"] = func(ctx context.Context) error {
        return db.PingContext(ctx)
    }
    
    // 缓存检查
    hc.checks["cache"] = func(ctx context.Context) error {
        return redisClient.Ping(ctx).Err()
    }
    
    // 外部服务检查
    hc.checks["external-api"] = func(ctx context.Context) error {
        req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.external.com/health", nil)
        resp, err := httpClient.Do(req)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("external API unhealthy: %d", resp.StatusCode)
        }
        return nil
    }
    
    return hc
}
```

### Startup（启动检查）

```go
// 应用是否已启动完成
var startupComplete int32

func StartupHandler(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&startupComplete) == 1 {
        w.WriteHeader(http.StatusOK)
        return
    }
    w.WriteHeader(http.StatusServiceUnavailable)
}

// 启动完成后设置
func init() {
    // 执行初始化
    initialize()
    
    atomic.StoreInt32(&startupComplete, 1)
}
```

---

## Kubernetes 配置

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        
        startupProbe:
          httpGet:
            path: /health/startup
            port: 8080
          initialDelaySeconds: 1
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 30  # 最多等待 150s
```

---

## 健康检查最佳实践

### 1. 区分检查类型

```go
// Liveness: 简单快速
http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)  // 只要不 panic 就活着
})

// Readiness: 检查依赖
http.HandleFunc("/health/ready", readinessHandler)
```

### 2. 避免级联检查

```go
// ❌ 不好：检查会触发其他服务的检查
func badCheck() error {
    // 这会触发下游服务检查它们的依赖
    return http.Get("http://downstream/health/deep")
}

// ✅ 好：只检查直接依赖
func goodCheck() error {
    return http.Get("http://downstream/health/live")
}
```

### 3. 超时控制

```go
func HealthCheckWithTimeout(check HealthCheck, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    return check(ctx)
}
```

---

## 健康聚合

```go
type HealthAggregator struct {
    services map[string]*url.URL
}

func (a *HealthAggregator) Aggregate(ctx context.Context) HealthReport {
    report := HealthReport{
        Services: make(map[string]ServiceHealth),
        Overall:  "healthy",
    }
    
    var wg sync.WaitGroup
    mu := sync.Mutex{}
    
    for name, url := range a.services {
        wg.Add(1)
        go func(n string, u *url.URL) {
            defer wg.Done()
            
            health := a.checkService(ctx, u)
            
            mu.Lock()
            report.Services[n] = health
            if health.Status != "healthy" {
                report.Overall = "degraded"
            }
            mu.Unlock()
        }(name, url)
    }
    
    wg.Wait()
    return report
}
```
