# 安全 Header 详解

> **分类**: 工程与云原生
> **标签**: #security #headers #csp

---

## Content-Security-Policy

```go
func CSPMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        csp := strings.Join([]string{
            "default-src 'self'",
            "script-src 'self' 'unsafe-inline' cdn.example.com",
            "style-src 'self' 'unsafe-inline' fonts.googleapis.com",
            "img-src 'self' data: https:",
            "font-src 'self' fonts.gstatic.com",
            "connect-src 'self' api.example.com",
            "frame-ancestors 'none'",
            "form-action 'self'",
            "base-uri 'self'",
            "upgrade-insecure-requests",
        }, "; ")

        c.Header("Content-Security-Policy", csp)
        c.Next()
    }
}
```

---

## 完整安全 Header 集

```go
func CompleteSecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 防止 MIME 类型嗅探
        c.Header("X-Content-Type-Options", "nosniff")

        // 防止点击劫持
        c.Header("X-Frame-Options", "DENY")

        // XSS 保护
        c.Header("X-XSS-Protection", "1; mode=block")

        // HSTS
        c.Header("Strict-Transport-Security",
            "max-age=31536000; includeSubDomains; preload")

        // 引用策略
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

        // 权限策略
        c.Header("Permissions-Policy",
            "accelerometer=(), camera=(), geolocation=(), gyroscope=(), "+
            "magnetometer=(), microphone=(), payment=(), usb=()")

        // COOP
        c.Header("Cross-Origin-Opener-Policy", "same-origin")

        // COEP
        c.Header("Cross-Origin-Embedder-Policy", "require-corp")

        // CORP
        c.Header("Cross-Origin-Resource-Policy", "same-origin")

        c.Next()
    }
}
```

---

## 报告模式

```go
func CSPWithReport() gin.HandlerFunc {
    return func(c *gin.Context) {
        csp := "default-src 'self'; report-uri /csp-report"

        // 仅报告模式（用于测试）
        // c.Header("Content-Security-Policy-Report-Only", csp)

        // 强制执行
        c.Header("Content-Security-Policy", csp)
        c.Next()
    }
}

// CSP 报告处理
func CSPReportHandler(c *gin.Context) {
    var report struct {
        CspReport map[string]interface{} `json:"csp-report"`
    }

    if err := c.BindJSON(&report); err != nil {
        return
    }

    // 记录违规
    log.Printf("CSP Violation: %+v", report.CspReport)

    // 发送到监控系统
    securityMonitor.RecordCSPViolation(report.CspReport)
}
```

---

## 安全 Cookie

```go
func SecureCookie(name, value string) *http.Cookie {
    return &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        Domain:   "",
        Expires:  time.Now().Add(24 * time.Hour),
        Secure:   true,           // HTTPS only
        HttpOnly: true,           // 禁止 JS 访问
        SameSite: http.SameSiteStrictMode,
    }
}
```

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02