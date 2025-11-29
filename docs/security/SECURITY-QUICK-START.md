# 安全功能快速开始

> **版本**: v1.0.0
> **状态**: ✅ 生产就绪

---

## 🚀 快速开始

### 1. 基本安全配置

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建安全配置
    config := security.DefaultSecurityConfig()
    
    // 验证配置
    if err := config.Validate(); err != nil {
        log.Fatal(err)
    }
    
    // 创建安全配置管理器
    manager, _ := security.NewSecurityConfigManager(config)
    
    // 使用配置...
}
```

### 2. OAuth2/OIDC 服务器

```go
package main

import (
    "github.com/yourusername/golang/pkg/auth/oauth2"
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
    
    // 使用服务器...
}
```

### 3. 数据加密

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建加密器
    encryptor, _ := security.NewAES256EncryptorFromString("your-secret-key")
    
    // 加密数据
    ciphertext, _ := encryptor.EncryptString("sensitive data")
    
    // 解密数据
    plaintext, _ := encryptor.DecryptString(ciphertext)
}
```

### 4. 密码哈希

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建密码哈希器
    hasher := security.NewPasswordHasher(security.DefaultPasswordHashConfig())
    
    // 哈希密码
    hash, _ := hasher.Hash("user-password")
    
    // 验证密码
    valid, _ := hasher.Verify("user-password", hash)
}
```

### 5. 速率限制

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建速率限制器
    limiter := security.NewIPRateLimiter(security.RateLimiterConfig{
        Limit:  100,
        Window: 1 * time.Minute,
    })
    defer limiter.Shutdown(context.Background())
    
    // 检查是否允许
    allowed, _ := limiter.AllowIP(ctx, "192.168.1.1")
}
```

### 6. CSRF 防护

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
)

func main() {
    // 创建 CSRF 防护
    csrf := security.NewCSRFProtection(security.DefaultCSRFConfig())
    defer csrf.Shutdown()
    
    // 生成令牌
    token, _ := csrf.GenerateToken(sessionID)
    
    // 验证令牌
    err := csrf.ValidateToken(sessionID, token)
}
```

### 7. 安全中间件

```go
package main

import (
    "github.com/yourusername/golang/pkg/security"
    "github.com/go-chi/chi/v5"
)

func main() {
    // 创建安全中间件
    config := security.SecurityMiddlewareConfig{
        SecurityHeaders: &security.DefaultSecurityHeadersConfig(),
        RateLimit: &security.RateLimiterConfig{
            Limit:  100,
            Window: 1 * time.Minute,
        },
    }
    
    middleware := security.NewSecurityMiddleware(config)
    defer middleware.Shutdown()
    
    // 使用中间件
    router := chi.NewRouter()
    router.Use(middleware.Middleware)
}
```

---

## 📚 更多文档

- [安全最佳实践](SECURITY-BEST-PRACTICES.md)
- [安全功能 API 文档](../../pkg/security/README.md)
- [OAuth2/OIDC 文档](../../pkg/auth/oauth2/README.md)

---

**最后更新**: 2025-01-XX

