# 安全默认配置 (Secure Defaults)

> **分类**: 工程与云原生  
> **标签**: #security #configuration #hardening

---

## HTTP 服务器安全配置

```go
func SecureServer() *http.Server {
    return &http.Server{
        Addr:         ":8443",
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        MaxHeaderBytes: 1 << 20,  // 1MB
        
        TLSConfig: &tls.Config{
            MinVersion:               tls.VersionTLS12,
            PreferServerCipherSuites: true,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            },
        },
    }
}
```

---

## 安全 Header

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 防止 XSS
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // CSP
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; script-src 'self'; object-src 'none'")
        
        // HSTS
        w.Header().Set("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains; preload")
        
        // 引用策略
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // 权限策略
        w.Header().Set("Permissions-Policy", 
            "geolocation=(), microphone=(), camera=()")
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 数据库安全

```go
type DBConfig struct {
    MaxOpenConns    int           `default:"25"`
    MaxIdleConns    int           `default:"5"`
    ConnMaxLifetime time.Duration `default:"5m"`
    
    // 安全选项
    SSLMode     string `default:"require"`
    SSLRootCert string `default:"/etc/ssl/certs/ca.crt"`
}

func SecureDBConnection(cfg DBConfig) (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
        cfg.Host, cfg.User, cfg.Password, cfg.DBName, 
        cfg.SSLMode, cfg.SSLRootCert,
    )
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // 设置连接限制
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    return db, nil
}
```

---

## 日志安全

```go
func SanitizeLogFields(data map[string]interface{}) map[string]interface{} {
    sensitive := []string{"password", "token", "secret", "credit_card", "ssn"}
    
    result := make(map[string]interface{})
    for k, v := range data {
        isSensitive := false
        for _, s := range sensitive {
            if strings.Contains(strings.ToLower(k), s) {
                isSensitive = true
                break
            }
        }
        
        if isSensitive {
            result[k] = "[REDACTED]"
        } else {
            result[k] = v
        }
    }
    
    return result
}
```

---

## 配置验证

```go
type SecurityValidator struct {
    rules []ValidationRule
}

func (v *SecurityValidator) Validate(cfg Config) error {
    for _, rule := range v.rules {
        if err := rule.Check(cfg); err != nil {
            return fmt.Errorf("security violation: %w", err)
        }
    }
    return nil
}

// 内置规则
var DefaultRules = []ValidationRule{
    {
        Name: "no_default_passwords",
        Check: func(cfg Config) error {
            if cfg.Password == "admin" || cfg.Password == "password" {
                return fmt.Errorf("default password detected")
            }
            return nil
        },
    },
    {
        Name: "tls_required",
        Check: func(cfg Config) error {
            if !cfg.TLS.Enabled {
                return fmt.Errorf("TLS is required")
            }
            return nil
        },
    },
    {
        Name: "strong_crypto",
        Check: func(cfg Config) error {
            if cfg.TLS.MinVersion < tls.VersionTLS12 {
                return fmt.Errorf("TLS 1.2+ required")
            }
            return nil
        },
    },
}
```
