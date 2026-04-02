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
