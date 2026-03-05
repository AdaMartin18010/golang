# API 参考文档

> **版本**: v1.0.0
> **更新日期**: 2025-01-XX
> **状态**: ✅ 生产就绪

---

## 📚 目录

- [API 参考文档](#api-参考文档)
  - [📚 目录](#-目录)
  - [认证和授权](#认证和授权)
    - [OAuth2/OIDC](#oauth2oidc)
      - [服务器初始化](#服务器初始化)
      - [授权码流程](#授权码流程)
      - [客户端凭证流程](#客户端凭证流程)
      - [刷新令牌流程](#刷新令牌流程)
      - [OIDC 功能](#oidc-功能)
  - [安全功能](#安全功能)
    - [数据加密](#数据加密)
    - [密钥管理](#密钥管理)
    - [密码哈希](#密码哈希)
    - [速率限制](#速率限制)
    - [CSRF 防护](#csrf-防护)
    - [XSS 防护](#xss-防护)
    - [SQL 注入防护](#sql-注入防护)
    - [审计日志](#审计日志)
    - [安全中间件](#安全中间件)
  - [可观测性](#可观测性)
    - [OTLP 集成](#otlp-集成)
    - [系统监控](#系统监控)
  - [存储](#存储)
    - [PostgreSQL 存储](#postgresql-存储)
    - [Redis 存储](#redis-存储)
  - [测试框架](#测试框架)
    - [测试框架初始化](#测试框架初始化)
    - [测试数据工厂](#测试数据工厂)
    - [测试覆盖率](#测试覆盖率)
  - [📚 更多文档](#-更多文档)

---

## 认证和授权

### OAuth2/OIDC

#### 服务器初始化

```go
import "github.com/yourusername/golang/pkg/auth/oauth2"

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
```

#### 授权码流程

```go
// 生成授权码
code, err := server.AuthorizationCodeFlow(ctx, &oauth2.AuthorizationRequest{
    ClientID:     "client-id",
    RedirectURI:  "http://localhost:3000/callback",
    Scope:        "read write",
    State:        "random-state",
    ResponseType: "code",
})

// 交换令牌
token, err := server.AuthorizationCodeFlow(ctx, &oauth2.TokenRequest{
    GrantType:    "authorization_code",
    Code:         code,
    RedirectURI:  "http://localhost:3000/callback",
    ClientID:     "client-id",
    ClientSecret: "client-secret",
})
```

#### 客户端凭证流程

```go
token, err := server.ClientCredentialsFlow(ctx, &oauth2.ClientCredentialsRequest{
    ClientID:     "client-id",
    ClientSecret: "client-secret",
    Scope:        "read write",
})
```

#### 刷新令牌流程

```go
token, err := server.RefreshTokenFlow(ctx, &oauth2.RefreshTokenRequest{
    RefreshToken: "refresh-token",
    ClientID:     "client-id",
    ClientSecret: "client-secret",
    Scope:        "read write",
})
```

#### OIDC 功能

```go
import "github.com/yourusername/golang/pkg/auth/oauth2"

// 创建OIDC提供者
provider := oauth2.NewOIDCProvider(oauth2.OIDCConfig{
    Issuer: "https://example.com",
})

// 生成ID Token
idToken, err := provider.GenerateIDToken(ctx, &oauth2.IDTokenClaims{
    Subject:  "user-123",
    Issuer:   "https://example.com",
    Audience: "client-id",
    Expiry:   time.Now().Add(1 * time.Hour),
})

// 获取用户信息
userInfo, err := provider.GetUserInfo(ctx, "access-token")

// 获取Discovery配置
discovery, err := provider.GetDiscoveryConfig()

// 获取JWKS
jwks, err := provider.GetJWKS()
```

---

## 安全功能

### 数据加密

```go
import "github.com/yourusername/golang/pkg/security"

// 创建加密器
encryptor, err := security.NewAES256EncryptorFromString("your-secret-key")

// 加密数据
ciphertext, err := encryptor.EncryptString("sensitive data")

// 解密数据
plaintext, err := encryptor.DecryptString(ciphertext)

// 字段级加密
fieldEncryptor := security.NewFieldEncryptor(encryptor)
encryptedEmail, err := fieldEncryptor.EncryptField("user@example.com")
```

### 密钥管理

```go
import "github.com/yourusername/golang/pkg/security"

// 创建密钥管理器
keyManager := security.NewKeyManager(security.KeyManagerConfig{
    StoragePath: "/path/to/keys",
})

// 生成AES密钥
aesKey, err := keyManager.GenerateAESKey(256)

// 生成RSA密钥对
rsaKeyPair, err := keyManager.GenerateRSAKeyPair(2048)

// 存储密钥
err := keyManager.StoreKey("key-id", aesKey)

// 检索密钥
key, err := keyManager.RetrieveKey("key-id")
```

### 密码哈希

```go
import "github.com/yourusername/golang/pkg/security"

// 创建密码哈希器
hasher := security.NewPasswordHasher(security.DefaultPasswordHashConfig())

// 哈希密码
hash, err := hasher.Hash("user-password")

// 验证密码
valid, err := hasher.Verify("user-password", hash)

// 密码验证器
validator := security.NewPasswordValidator(security.DefaultPasswordValidatorConfig())
err := validator.Validate("user-password")
```

### 速率限制

```go
import "github.com/yourusername/golang/pkg/security"

// IP级别速率限制
ipLimiter := security.NewIPRateLimiter(security.RateLimiterConfig{
    Limit:  100,
    Window: 1 * time.Minute,
})
defer ipLimiter.Shutdown(ctx)

allowed, err := ipLimiter.AllowIP(ctx, "192.168.1.1")

// 用户级别速率限制
userLimiter := security.NewUserRateLimiter(security.RateLimiterConfig{
    Limit:  1000,
    Window: 1 * time.Hour,
})
defer userLimiter.Shutdown(ctx)

allowed, err := userLimiter.AllowUser(ctx, "user-123")

// 端点级别速率限制
endpointLimiter := security.NewEndpointRateLimiter(security.RateLimiterConfig{
    Limit:  10,
    Window: 1 * time.Minute,
})
defer endpointLimiter.Shutdown(ctx)

allowed, err := endpointLimiter.AllowEndpoint(ctx, "/api/login", "192.168.1.1")
```

### CSRF 防护

```go
import "github.com/yourusername/golang/pkg/security"

// 创建CSRF防护
csrf := security.NewCSRFProtection(security.DefaultCSRFConfig())
defer csrf.Shutdown()

// 生成令牌
token, err := csrf.GenerateToken("session-id")

// 验证令牌
err := csrf.ValidateToken("session-id", token)
```

### XSS 防护

```go
import "github.com/yourusername/golang/pkg/security"

// 创建XSS防护
xss := security.NewXSSProtection()

// 清理输入
sanitized := xss.Sanitize("<script>alert('XSS')</script>")

// 转义HTML
escaped := xss.EscapeHTML("<div>content</div>")

// 检测XSS攻击
isXSS, err := xss.DetectXSS("<script>alert('XSS')</script>")
```

### SQL 注入防护

```go
import "github.com/yourusername/golang/pkg/security"

// 创建SQL注入防护
sqlProtection := security.NewSQLInjectionProtection(true)

// 验证输入
err := sqlProtection.ValidateInput("'; DROP TABLE users; --")

// 清理输入
sanitized := sqlProtection.SanitizeInput("'; DROP TABLE users; --")
```

### 审计日志

```go
import "github.com/yourusername/golang/pkg/security"

// 创建审计日志器
logger := security.NewAuditLogger(store)

// 记录安全事件
err := logger.LogSecurity(ctx, "user-123", "failed_login", map[string]interface{}{
    "attempts": 3,
    "ip": "192.168.1.1",
})

// 记录数据访问
err := logger.LogDataAccess(ctx, "user-123", "read", "user-data", map[string]interface{}{
    "resource": "user-profile",
})

// 查询日志
filter := &security.AuditLogFilter{
    UserID:    "user-123",
    StartTime: &startTime,
    EndTime:   &endTime,
}
logs, err := logger.QueryLogs(ctx, filter)
```

### 安全中间件

```go
import (
    "github.com/yourusername/golang/pkg/security"
    "github.com/go-chi/chi/v5"
)

// 创建安全中间件配置
config := security.SecurityMiddlewareConfig{
    SecurityHeaders: &security.DefaultSecurityHeadersConfig(),
    RateLimit: &security.RateLimiterConfig{
        Limit:  100,
        Window: 1 * time.Minute,
    },
    CSRF: &security.DefaultCSRFConfig(),
    EnableXSS: true,
}

// 创建安全中间件
middleware := security.NewSecurityMiddleware(config)
defer middleware.Shutdown()

// 使用中间件
router := chi.NewRouter()
router.Use(middleware.Middleware)
```

---

## 可观测性

### OTLP 集成

```go
import "github.com/yourusername/golang/pkg/observability"

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

// 获取追踪器
tracer := obs.GetTracer("my-component")

// 创建span
ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()

// 获取指标器
meter := obs.GetMeter("my-component")

// 创建计数器
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1)
```

### 系统监控

```go
import "github.com/yourusername/golang/pkg/observability"

// 启用系统监控
cfg := observability.Config{
    EnableSystemMonitoring: true,
    SystemCollectInterval:  10 * time.Second,
}

obs, err := observability.NewObservability(cfg)

// 获取系统监控器
systemMonitor := obs.GetSystemMonitor()

// 启动监控
err = systemMonitor.Start(ctx)

// 停止监控
err = systemMonitor.Stop(ctx)
```

---

## 存储

### PostgreSQL 存储

```go
import (
    "github.com/yourusername/golang/pkg/auth/oauth2"
    "database/sql"
)

// 创建PostgreSQL存储
tokenStore := oauth2.NewPostgresTokenStore(db)
clientStore := oauth2.NewPostgresClientStore(db)

// 保存客户端
client := &oauth2.Client{
    ID:          "client-id",
    Secret:      "client-secret",
    RedirectURI: "http://localhost:3000/callback",
}
err := clientStore.Save(ctx, client)

// 保存令牌
token := &oauth2.Token{
    AccessToken: "access-token",
    ClientID:    "client-id",
    UserID:      "user-123",
}
err := tokenStore.Save(ctx, token)
```

### Redis 存储

```go
import (
    "github.com/yourusername/golang/pkg/auth/oauth2"
    "github.com/redis/go-redis/v9"
)

// 创建Redis存储
tokenStore := oauth2.NewRedisTokenStore(redisClient)
clientStore := oauth2.NewRedisClientStore(redisClient)

// 使用方式与PostgreSQL相同
```

---

## 测试框架

### 测试框架初始化

```go
import "github.com/yourusername/golang/test/framework"

// 创建测试框架
tf, err := framework.NewTestFramework(framework.DefaultTestFrameworkConfig())
if err != nil {
    t.Skipf("Skipping integration test: %v", err)
}
defer tf.Shutdown()
```

### 测试数据工厂

```go
import "github.com/yourusername/golang/test/framework"

// 创建测试数据工厂
factory := framework.NewTestDataFactory()

// 生成测试用户
user := factory.NewUser()

// 生成OAuth2客户端
client := factory.NewOAuth2Client()

// 生成OAuth2令牌
token := factory.NewOAuth2Token()
```

### 测试覆盖率

```go
import "github.com/yourusername/golang/test/framework"

// 生成覆盖率报告
reporter := framework.NewCoverageReporter()
err := reporter.GenerateReport("coverage.out", framework.CoverageReportConfig{
    Format: framework.CoverageFormatHTML,
    Output: "coverage.html",
})
```

---

## 📚 更多文档

- [安全最佳实践](../docs/security/SECURITY-BEST-PRACTICES.md)
- [安全功能快速开始](../docs/security/SECURITY-QUICK-START.md)
- [OAuth2/OIDC 文档](../pkg/auth/oauth2/README.md)

---

**最后更新**: 2025-01-XX
