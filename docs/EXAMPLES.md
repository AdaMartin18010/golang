# 示例代码

> **版本**: v1.0.0
> **更新日期**: 2025-01-XX
> **状态**: ✅ 生产就绪

---

## 📚 目录

- [OAuth2/OIDC 示例](#oauth2oidc-示例)
- [安全功能示例](#安全功能示例)
- [可观测性示例](#可观测性示例)
- [存储示例](#存储示例)
- [完整应用示例](#完整应用示例)

---

## OAuth2/OIDC 示例

### 基本 OAuth2 服务器

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/yourusername/golang/pkg/auth/oauth2"
    "github.com/go-chi/chi/v5"
)

func main() {
    // 创建存储
    tokenStore := oauth2.NewMemoryTokenStore()
    clientStore := oauth2.NewMemoryClientStore()
    codeStore := oauth2.NewMemoryCodeStore()

    // 创建服务器
    server := oauth2.NewServer(
        oauth2.DefaultServerConfig(),
        tokenStore,
        clientStore,
        codeStore,
    )

    // 创建路由
    router := chi.NewRouter()

    // 授权端点
    router.Get("/authorize", func(w http.ResponseWriter, r *http.Request) {
        code, err := server.AuthorizationCodeFlow(r.Context(), &oauth2.AuthorizationRequest{
            ClientID:     r.URL.Query().Get("client_id"),
            RedirectURI:  r.URL.Query().Get("redirect_uri"),
            Scope:        r.URL.Query().Get("scope"),
            State:        r.URL.Query().Get("state"),
            ResponseType: r.URL.Query().Get("response_type"),
        })
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        http.Redirect(w, r, r.URL.Query().Get("redirect_uri")+"?code="+code+"&state="+r.URL.Query().Get("state"), http.StatusFound)
    })

    // 令牌端点
    router.Post("/token", func(w http.ResponseWriter, r *http.Request) {
        token, err := server.AuthorizationCodeFlow(r.Context(), &oauth2.TokenRequest{
            GrantType:    r.FormValue("grant_type"),
            Code:         r.FormValue("code"),
            RedirectURI:  r.FormValue("redirect_uri"),
            ClientID:     r.FormValue("client_id"),
            ClientSecret: r.FormValue("client_secret"),
        })
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(token)
    })

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### OIDC 提供者

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/yourusername/golang/pkg/auth/oauth2"
    "github.com/go-chi/chi/v5"
)

func main() {
    // 创建OIDC提供者
    provider := oauth2.NewOIDCProvider(oauth2.OIDCConfig{
        Issuer: "https://example.com",
    })

    router := chi.NewRouter()

    // Discovery端点
    router.Get("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
        discovery, err := provider.GetDiscoveryConfig()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(discovery)
    })

    // UserInfo端点
    router.Get("/userinfo", func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        userInfo, err := provider.GetUserInfo(r.Context(), token)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userInfo)
    })

    // JWKS端点
    router.Get("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
        jwks, err := provider.GetJWKS()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(jwks)
    })

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

---

## 安全功能示例

### 数据加密

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建加密器
    encryptor, err := security.NewAES256EncryptorFromString("your-32-byte-secret-key-here!")
    if err != nil {
        log.Fatal(err)
    }

    // 加密数据
    plaintext := "sensitive user data"
    ciphertext, err := encryptor.EncryptString(plaintext)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Encrypted: %s\n", ciphertext)

    // 解密数据
    decrypted, err := encryptor.DecryptString(ciphertext)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Decrypted: %s\n", decrypted)

    // 字段级加密
    fieldEncryptor := security.NewFieldEncryptor(encryptor)
    encryptedEmail, err := fieldEncryptor.EncryptField("user@example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Encrypted Email: %s\n", encryptedEmail)
}
```

### 密码哈希

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建密码哈希器
    hasher := security.NewPasswordHasher(security.DefaultPasswordHashConfig())

    // 哈希密码
    password := "user-password-123"
    hash, err := hasher.Hash(password)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Password Hash: %s\n", hash)

    // 验证密码
    valid, err := hasher.Verify(password, hash)
    if err != nil {
        log.Fatal(err)
    }

    if valid {
        fmt.Println("Password is valid")
    } else {
        fmt.Println("Password is invalid")
    }

    // 验证密码强度
    validator := security.NewPasswordValidator(security.DefaultPasswordValidatorConfig())
    err = validator.Validate(password)
    if err != nil {
        fmt.Printf("Password validation failed: %v\n", err)
    } else {
        fmt.Println("Password meets requirements")
    }
}
```

### 速率限制

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/yourusername/golang/pkg/security"
    "github.com/go-chi/chi/v5"
)

func main() {
    // 创建速率限制器
    limiter := security.NewIPRateLimiter(security.RateLimiterConfig{
        Limit:  10,
        Window: 1 * time.Minute,
    })
    defer limiter.Shutdown(context.Background())

    // 创建速率限制中间件
    rateLimitMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            allowed, err := limiter.AllowIP(r.Context(), ip)
            if err != nil || !allowed {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }

    router := chi.NewRouter()
    router.Use(rateLimitMiddleware)

    router.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Data"))
    })

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

---

## 可观测性示例

### OTLP 集成

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourusername/golang/pkg/observability"
)

func main() {
    // 创建可观测性配置
    cfg := observability.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        OTLPEndpoint:   "http://localhost:4317",
        OTLPInsecure:   true,
        SampleRate:     1.0,
    }

    // 创建可观测性实例
    obs, err := observability.NewObservability(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer obs.Shutdown(context.Background())

    // 获取追踪器
    tracer := obs.GetTracer("my-component")

    // 创建span
    ctx, span := tracer.Start(context.Background(), "operation-name")
    defer span.End()

    // 执行操作
    time.Sleep(100 * time.Millisecond)

    // 获取指标器
    meter := obs.GetMeter("my-component")

    // 创建计数器
    counter, _ := meter.Int64Counter("requests_total")
    counter.Add(ctx, 1)
}
```

---

## 存储示例

### PostgreSQL 存储

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    "github.com/yourusername/golang/pkg/auth/oauth2"
    _ "github.com/lib/pq"
)

func main() {
    // 连接PostgreSQL
    db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建存储
    tokenStore := oauth2.NewPostgresTokenStore(db)
    clientStore := oauth2.NewPostgresClientStore(db)

    ctx := context.Background()

    // 保存客户端
    client := &oauth2.Client{
        ID:          "client-id",
        Secret:      "client-secret",
        RedirectURI: "http://localhost:3000/callback",
        Scopes:      []string{"read", "write"},
    }

    err = clientStore.Save(ctx, client)
    if err != nil {
        log.Fatal(err)
    }

    // 保存令牌
    token := &oauth2.Token{
        AccessToken:  "access-token",
        TokenType:    "Bearer",
        ExpiresIn:    3600,
        RefreshToken: "refresh-token",
        Scope:        "read write",
        ClientID:     client.ID,
        UserID:       "user-123",
        CreatedAt:    time.Now(),
        ExpiresAt:    time.Now().Add(1 * time.Hour),
    }

    err = tokenStore.Save(ctx, token)
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## 完整应用示例

### 完整的 REST API 服务器

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/yourusername/golang/pkg/auth/oauth2"
    "github.com/yourusername/golang/pkg/observability"
    "github.com/yourusername/golang/pkg/security"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    // 初始化可观测性
    obs, err := observability.NewObservability(observability.Config{
        ServiceName:    "api-server",
        ServiceVersion: "1.0.0",
        OTLPEndpoint:   "http://localhost:4317",
        OTLPInsecure:   true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer obs.Shutdown(context.Background())

    // 初始化OAuth2服务器
    tokenStore := oauth2.NewMemoryTokenStore()
    clientStore := oauth2.NewMemoryClientStore()
    codeStore := oauth2.NewMemoryCodeStore()

    server := oauth2.NewServer(
        oauth2.DefaultServerConfig(),
        tokenStore,
        clientStore,
        codeStore,
    )

    // 初始化安全中间件
    securityMiddleware := security.NewSecurityMiddleware(security.SecurityMiddlewareConfig{
        SecurityHeaders: &security.DefaultSecurityHeadersConfig(),
        RateLimit: &security.RateLimiterConfig{
            Limit:  100,
            Window: 1 * time.Minute,
        },
    })
    defer securityMiddleware.Shutdown()

    // 创建路由
    router := chi.NewRouter()

    // 中间件
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
    router.Use(securityMiddleware.Middleware)

    // API路由
    router.Route("/api", func(r chi.Router) {
        r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("OK"))
        })

        r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Data"))
        })
    })

    // OAuth2路由
    router.Route("/oauth2", func(r chi.Router) {
        r.Get("/authorize", func(w http.ResponseWriter, r *http.Request) {
            // 处理授权请求
        })

        r.Post("/token", func(w http.ResponseWriter, r *http.Request) {
            // 处理令牌请求
        })
    })

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

---

## 📚 更多文档

- [API 参考](API-REFERENCE.md)
- [安全最佳实践](security/SECURITY-BEST-PRACTICES.md)
- [快速开始指南](QUICK-START.md)

---

**最后更新**: 2025-01-XX
