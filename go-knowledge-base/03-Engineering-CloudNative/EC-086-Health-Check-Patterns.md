# 健康检查模式 (Health Check Patterns)

> **分类**: 工程与云原生
> **标签**: #health-check #probes #kubernetes #monitoring
> **参考**: Kubernetes Liveness/Readiness Probes, Google SRE

---

## 健康检查架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Health Check Architecture                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Probe Types                                       │   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   Liveness   │  │  Readiness   │  │  Startup     │              │   │
│  │  │              │  │              │  │              │              │   │
│  │  │ "Is process  │  │ "Is ready to │  │ "Has app     │              │   │
│  │  │  alive?"     │  │  serve?"     │  │  started?"   │              │   │
│  │  │              │  │              │  │              │              │   │
│  │  │ Failure ──►  │  │ Failure ──►  │  │ Failure ──►  │              │   │
│  │  │ Restart      │  │ Remove from  │  │ Wait         │              │   │
│  │  │ container    │  │ service pool │  │ (no action)  │              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Probe Mechanisms                                  │   │
│  │                                                                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │  HTTP    │  │   TCP    │  │  Command │  │   gRPC   │            │   │
│  │  │  GET     │  │  Socket  │  │  Exec    │  │  Call    │            │   │
│  │  │  /health │  │  Connect │  │  Custom  │  │  Health  │            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整健康检查实现

```go
package health

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Status 健康状态
type Status string

const (
    StatusHealthy   Status = "healthy"
    StatusUnhealthy Status = "unhealthy"
    StatusDegraded  Status = "degraded"
    StatusUnknown   Status = "unknown"
)

// Check 健康检查接口
type Check interface {
    Name() string
    Execute(ctx context.Context) CheckResult
}

// CheckResult 检查结果
type CheckResult struct {
    Name      string                 `json:"name"`
    Status    Status                 `json:"status"`
    Message   string                 `json:"message,omitempty"`
    Duration  time.Duration          `json:"duration"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
    checks map[string]Check
    mu     sync.RWMutex

    // 配置
    timeout time.Duration
    cache   *checkCache
}

// checkCache 检查缓存
type checkCache struct {
    results   map[string]CheckResult
    timestamp time.Time
    ttl       time.Duration
    mu        sync.RWMutex
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(timeout time.Duration, cacheTTL time.Duration) *HealthChecker {
    return &HealthChecker{
        checks:  make(map[string]Check),
        timeout: timeout,
        cache: &checkCache{
            results: make(map[string]CheckResult),
            ttl:     cacheTTL,
        },
    }
}

// Register 注册检查
func (hc *HealthChecker) Register(check Check) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    hc.checks[check.Name()] = check
}

// Unregister 注销检查
func (hc *HealthChecker) Unregister(name string) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    delete(hc.checks, name)
}

// CheckAll 执行所有检查
func (hc *HealthChecker) CheckAll(ctx context.Context) HealthReport {
    hc.mu.RLock()
    checks := make(map[string]Check, len(hc.checks))
    for k, v := range hc.checks {
        checks[k] = v
    }
    hc.mu.RUnlock()

    report := HealthReport{
        Status:    StatusHealthy,
        Timestamp: time.Now(),
        Checks:    make(map[string]CheckResult),
    }

    var wg sync.WaitGroup
    resultChan := make(chan CheckResult, len(checks))

    for _, check := range checks {
        wg.Add(1)
        go func(c Check) {
            defer wg.Done()

            // 检查缓存
            if result, ok := hc.cache.Get(c.Name()); ok {
                resultChan <- result
                return
            }

            // 执行检查
            ctx, cancel := context.WithTimeout(ctx, hc.timeout)
            defer cancel()

            start := time.Now()
            result := c.Execute(ctx)
            result.Duration = time.Since(start)

            // 更新缓存
            hc.cache.Set(c.Name(), result)

            resultChan <- result
        }(check)
    }

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    for result := range resultChan {
        report.Checks[result.Name] = result

        // 更新总体状态
        if result.Status == StatusUnhealthy {
            report.Status = StatusUnhealthy
        } else if result.Status == StatusDegraded && report.Status == StatusHealthy {
            report.Status = StatusDegraded
        }
    }

    return report
}

// CheckOne 执行单个检查
func (hc *HealthChecker) CheckOne(ctx context.Context, name string) (CheckResult, error) {
    hc.mu.RLock()
    check, ok := hc.checks[name]
    hc.mu.RUnlock()

    if !ok {
        return CheckResult{}, fmt.Errorf("check %s not found", name)
    }

    ctx, cancel := context.WithTimeout(ctx, hc.timeout)
    defer cancel()

    return check.Execute(ctx), nil
}

// HTTPHandler HTTP处理器
func (hc *HealthChecker) HTTPHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 检查特定检查
        if checkName := r.URL.Query().Get("check"); checkName != "" {
            result, err := hc.CheckOne(ctx, checkName)
            if err != nil {
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }

            statusCode := http.StatusOK
            if result.Status == StatusUnhealthy {
                statusCode = http.StatusServiceUnavailable
            } else if result.Status == StatusDegraded {
                statusCode = http.StatusOK
            }

            w.WriteHeader(statusCode)
            json.NewEncoder(w).Encode(result)
            return
        }

        // 执行所有检查
        report := hc.CheckAll(ctx)

        statusCode := http.StatusOK
        if report.Status == StatusUnhealthy {
            statusCode = http.StatusServiceUnavailable
        } else if report.Status == StatusDegraded {
            statusCode = http.StatusOK
        }

        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(report)
    }
}

// HealthReport 健康报告
type HealthReport struct {
    Status    Status                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Checks    map[string]CheckResult `json:"checks"`
}

// cache 方法
func (c *checkCache) Get(name string) (CheckResult, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if time.Since(c.timestamp) > c.ttl {
        return CheckResult{}, false
    }

    result, ok := c.results[name]
    return result, ok
}

func (c *checkCache) Set(name string, result CheckResult) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.results[name] = result
    c.timestamp = time.Now()
}
```

---

## 内置检查实现

```go
package health

import (
    "context"
    "database/sql"
    "fmt"
    "net"
    "time"
)

// DatabaseCheck 数据库检查
type DatabaseCheck struct {
    name string
    db   *sql.DB
}

// NewDatabaseCheck 创建数据库检查
func NewDatabaseCheck(name string, db *sql.DB) Check {
    return &DatabaseCheck{name: name, db: db}
}

func (dc *DatabaseCheck) Name() string {
    return dc.name
}

func (dc *DatabaseCheck) Execute(ctx context.Context) CheckResult {
    start := time.Now()

    result := CheckResult{
        Name:      dc.name,
        Timestamp: start,
    }

    // 检查连接
    if err := dc.db.PingContext(ctx); err != nil {
        result.Status = StatusUnhealthy
        result.Message = fmt.Sprintf("database ping failed: %v", err)
        return result
    }

    result.Status = StatusHealthy
    result.Message = "database connection ok"
    result.Duration = time.Since(start)

    return result
}

// HTTPCheck HTTP检查
type HTTPCheck struct {
    name    string
    url     string
    timeout time.Duration
}

// NewHTTPCheck 创建HTTP检查
func NewHTTPCheck(name, url string, timeout time.Duration) Check {
    return &HTTPCheck{name: name, url: url, timeout: timeout}
}

func (hc *HTTPCheck) Name() string {
    return hc.name
}

func (hc *HTTPCheck) Execute(ctx context.Context) CheckResult {
    start := time.Now()

    result := CheckResult{
        Name:      hc.name,
        Timestamp: start,
    }

    client := &http.Client{Timeout: hc.timeout}
    resp, err := client.Get(hc.url)
    if err != nil {
        result.Status = StatusUnhealthy
        result.Message = fmt.Sprintf("http request failed: %v", err)
        return result
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        result.Status = StatusHealthy
        result.Message = fmt.Sprintf("http status: %d", resp.StatusCode)
    } else {
        result.Status = StatusUnhealthy
        result.Message = fmt.Sprintf("http error status: %d", resp.StatusCode)
    }

    result.Duration = time.Since(start)
    return result
}

// TCPCheck TCP检查
type TCPCheck struct {
    name    string
    address string
    timeout time.Duration
}

// NewTCPCheck 创建TCP检查
func NewTCPCheck(name, address string, timeout time.Duration) Check {
    return &TCPCheck{name: name, address: address, timeout: timeout}
}

func (tc *TCPCheck) Name() string {
    return tc.name
}

func (tc *TCPCheck) Execute(ctx context.Context) CheckResult {
    start := time.Now()

    result := CheckResult{
        Name:      tc.name,
        Timestamp: start,
    }

    conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
    if err != nil {
        result.Status = StatusUnhealthy
        result.Message = fmt.Sprintf("tcp connection failed: %v", err)
        return result
    }
    defer conn.Close()

    result.Status = StatusHealthy
    result.Message = "tcp connection ok"
    result.Duration = time.Since(start)

    return result
}

// CompositeCheck 组合检查
type CompositeCheck struct {
    name   string
    checks []Check
}

// NewCompositeCheck 创建组合检查
func NewCompositeCheck(name string, checks ...Check) Check {
    return &CompositeCheck{name: name, checks: checks}
}

func (cc *CompositeCheck) Name() string {
    return cc.name
}

func (cc *CompositeCheck) Execute(ctx context.Context) CheckResult {
    start := time.Now()

    result := CheckResult{
        Name:      cc.name,
        Timestamp: start,
        Metadata:  make(map[string]interface{}),
    }

    allHealthy := true
    anyUnhealthy := false

    for _, check := range cc.checks {
        checkResult := check.Execute(ctx)
        result.Metadata[check.Name()] = checkResult

        if checkResult.Status != StatusHealthy {
            allHealthy = false
        }
        if checkResult.Status == StatusUnhealthy {
            anyUnhealthy = true
        }
    }

    if anyUnhealthy {
        result.Status = StatusUnhealthy
        result.Message = "one or more checks failed"
    } else if !allHealthy {
        result.Status = StatusDegraded
        result.Message = "some checks degraded"
    } else {
        result.Status = StatusHealthy
        result.Message = "all checks passed"
    }

    result.Duration = time.Since(start)
    return result
}

// ThresholdCheck 阈值检查
type ThresholdCheck struct {
    name      string
    getter    func() float64
    threshold float64
    operator  string // >, <, >=, <=, ==
}

// NewThresholdCheck 创建阈值检查
func NewThresholdCheck(name string, getter func() float64, threshold float64, operator string) Check {
    return &ThresholdCheck{
        name:      name,
        getter:    getter,
        threshold: threshold,
        operator:  operator,
    }
}

func (tc *ThresholdCheck) Name() string {
    return tc.name
}

func (tc *ThresholdCheck) Execute(ctx context.Context) CheckResult {
    start := time.Now()

    result := CheckResult{
        Name:      tc.name,
        Timestamp: start,
    }

    value := tc.getter()
    passed := false

    switch tc.operator {
    case ">":
        passed = value > tc.threshold
    case "<":
        passed = value < tc.threshold
    case ">=":
        passed = value >= tc.threshold
    case "<=":
        passed = value <= tc.threshold
    case "==":
        passed = value == tc.threshold
    }

    if passed {
        result.Status = StatusHealthy
        result.Message = fmt.Sprintf("value %f %s %f", value, tc.operator, tc.threshold)
    } else {
        result.Status = StatusUnhealthy
        result.Message = fmt.Sprintf("value %f does not satisfy %s %f", value, tc.operator, tc.threshold)
    }

    result.Duration = time.Since(start)
    return result
}
```

---

## 使用示例

```go
package main

import (
    "database/sql"
    "net/http"
    "time"

    "health"
)

func main() {
    // 创建健康检查器
    checker := health.NewHealthChecker(5*time.Second, 10*time.Second)

    // 数据库连接
    db, _ := sql.Open("postgres", "...")

    // 注册检查
    checker.Register(health.NewDatabaseCheck("database", db))
    checker.Register(health.NewHTTPCheck("external-api", "https://api.example.com/health", 5*time.Second))
    checker.Register(health.NewTCPCheck("redis", "localhost:6379", 2*time.Second))

    // 组合检查
    checker.Register(health.NewCompositeCheck("dependencies",
        health.NewDatabaseCheck("db", db),
        health.NewTCPCheck("cache", "localhost:6379", 2*time.Second),
    ))

    // HTTP 端点
    http.HandleFunc("/health", checker.HTTPHandler())
    http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        // 只检查就绪相关的
        result, _ := checker.CheckOne(r.Context(), "dependencies")
        if result.Status == health.StatusHealthy {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        json.NewEncoder(w).Encode(result)
    })

    http.ListenAndServe(":8080", nil)
}
```
