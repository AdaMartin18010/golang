# AD-007: 应用安全设计模式 (Application Security Patterns)

> **维度**: Application Domains
> **级别**: S (17+ KB)
> **标签**: #security #authentication #authorization #jwt #oauth
> **权威来源**: [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/), [Security Patterns](https://www.oreilly.com/library/view/security-patterns-in/9780470858844/)

---

## 安全架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Defense in Depth                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Layer 1: 网络层                                                             │
│  ├── Firewall (WAF)                                                          │
│  ├── DDoS Protection                                                         │
│  └── TLS/mTLS                                                                │
│                                                                              │
│  Layer 2: 网关层                                                             │
│  ├── Rate Limiting                                                           │
│  ├── Authentication                                                          │
│  └── Request Validation                                                      │
│                                                                              │
│  Layer 3: 应用层                                                             │
│  ├── Authorization (RBAC/ABAC)                                               │
│  ├── Input Sanitization                                                      │
│  └── Output Encoding                                                         │
│                                                                              │
│  Layer 4: 数据层                                                             │
│  ├── Encryption at Rest                                                      │
│  ├── Encryption in Transit                                                   │
│  └── Access Control                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 认证模式

### JWT (JSON Web Token)

```go
package security

import (
    "context"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT 配置
type JWTConfig struct {
    SecretKey       []byte
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
    Issuer          string
}

// Claims 自定义声明
type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

// JWTManager JWT 管理器
type JWTManager struct {
    config *JWTConfig
}

func NewJWTManager(config *JWTConfig) *JWTManager {
    return &JWTManager{config: config}
}

// GenerateTokens 生成访问令牌和刷新令牌
func (m *JWTManager) GenerateTokens(userID, username string, roles []string) (*TokenPair, error) {
    // Access Token
    accessClaims := Claims{
        UserID:   userID,
        Username: username,
        Roles:    roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.AccessTokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    m.config.Issuer,
            Subject:   userID,
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(m.config.SecretKey)
    if err != nil {
        return nil, fmt.Errorf("sign access token: %w", err)
    }

    // Refresh Token
    refreshClaims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.RefreshTokenTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Issuer:    m.config.Issuer,
        Subject:   userID,
        ID:        generateRefreshTokenID(),
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(m.config.SecretKey)
    if err != nil {
        return nil, fmt.Errorf("sign refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int64(m.config.AccessTokenTTL.Seconds()),
    }, nil
}

// ValidateToken 验证令牌
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return m.config.SecretKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("parse token: %w", err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

// TokenPair 令牌对
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    TokenType    string `json:"token_type"`
}

func generateRefreshTokenID() string {
    // 生成唯一 ID
    return uuid.New().String()
}
```

### OAuth 2.0 / OIDC

```go
package security

import (
    "context"

    "golang.org/x/oauth2"
    "github.com/coreos/go-oidc/v3/oidc"
)

// OIDCConfig OpenID Connect 配置
type OIDCConfig struct {
    IssuerURL    string
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Scopes       []string
}

// OIDCAuthenticator OIDC 认证器
type OIDCAuthenticator struct {
    provider *oidc.Provider
    verifier *oidc.IDTokenVerifier
    oauth2Config *oauth2.Config
}

func NewOIDCAuthenticator(ctx context.Context, config *OIDCConfig) (*OIDCAuthenticator, error) {
    provider, err := oidc.NewProvider(ctx, config.IssuerURL)
    if err != nil {
        return nil, err
    }

    oauth2Config := &oauth2.Config{
        ClientID:     config.ClientID,
        ClientSecret: config.ClientSecret,
        RedirectURL:  config.RedirectURL,
        Endpoint:     provider.Endpoint(),
        Scopes:       append(config.Scopes, oidc.ScopeOpenID, "profile", "email"),
    }

    verifier := provider.Verifier(&oidc.Config{
        ClientID: config.ClientID,
    })

    return &OIDCAuthenticator{
        provider:     provider,
        verifier:     verifier,
        oauth2Config: oauth2Config,
    }, nil
}

// AuthCodeURL 获取授权 URL
func (a *OIDCAuthenticator) AuthCodeURL(state string) string {
    return a.oauth2Config.AuthCodeURL(state)
}

// Exchange 交换授权码获取令牌
func (a *OIDCAuthenticator) Exchange(ctx context.Context, code string) (*OIDCToken, error) {
    oauth2Token, err := a.oauth2Config.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }

    rawIDToken, ok := oauth2Token.Extra("id_token").(string)
    if !ok {
        return nil, fmt.Errorf("no id_token in response")
    }

    idToken, err := a.verifier.Verify(ctx, rawIDToken)
    if err != nil {
        return nil, err
    }

    var claims struct {
        Email         string `json:"email"`
        EmailVerified bool   `json:"email_verified"`
        Name          string `json:"name"`
        Picture       string `json:"picture"`
    }

    if err := idToken.Claims(&claims); err != nil {
        return nil, err
    }

    return &OIDCToken{
        AccessToken:  oauth2Token.AccessToken,
        RefreshToken: oauth2Token.RefreshToken,
        IDToken:      rawIDToken,
        Expiry:       oauth2Token.Expiry,
        Claims:       claims,
    }, nil
}
```

---

## 授权模式

### RBAC (基于角色的访问控制)

```go
package security

// RBAC 权限模型
// User -> Role -> Permission -> Resource

type Permission struct {
    Resource string // e.g., "order", "user"
    Action   string // e.g., "read", "write", "delete"
}

type Role struct {
    ID          string
    Name        string
    Permissions []Permission
}

type RBACAuthorizer struct {
    roles      map[string]*Role
    userRoles  map[string][]string // userID -> roleIDs
}

func (r *RBACAuthorizer) HasPermission(userID string, resource, action string) bool {
    roleIDs, ok := r.userRoles[userID]
    if !ok {
        return false
    }

    for _, roleID := range roleIDs {
        role, ok := r.roles[roleID]
        if !ok {
            continue
        }
        for _, perm := range role.Permissions {
            if perm.Resource == resource && perm.Action == action {
                return true
            }
        }
    }
    return false
}
```

### ABAC (基于属性的访问控制)

```go
// ABAC 更细粒度的控制
// 策略: 用户.部门 == 资源.所属部门 AND 时间.小时 >= 9 AND 时间.小时 <= 18

type ABACContext struct {
    User     UserAttributes
    Resource ResourceAttributes
    Environment EnvironmentAttributes
}

type UserAttributes struct {
    ID       string
    Department string
    Role     string
    ClearanceLevel int
}

type ResourceAttributes struct {
    ID          string
    Owner       string
    Department  string
    Classification string
}

type EnvironmentAttributes struct {
    Time     time.Time
    Location string
    Device   string
}

// Policy 策略规则
type Policy struct {
    ID       string
    Effect   string // "allow" or "deny"
    Conditions []Condition
}

type Condition struct {
    Attribute string
    Operator  string // "eq", "gt", "lt", "in"
    Value     interface{}
}
```

---

## 安全最佳实践

| 领域 | 实践 |
|------|------|
| 认证 | 使用强密码策略，支持 MFA |
| 会话 | 短 TTL 访问令牌，轮换刷新令牌 |
| 传输 | 强制 TLS 1.3，HSTS 头 |
| 存储 | bcrypt/Argon2 哈希密码，加密敏感数据 |
| 日志 | 记录安全事件，脱敏敏感信息 |
| 审计 | 不可篡改的操作日志 |

---

## OWASP Top 10 (2025)

1. Broken Access Control
2. Cryptographic Failures
3. Injection
4. Insecure Design
5. Security Misconfiguration
6. Vulnerable Components
7. Auth Failures
8. Data Integrity Failures
9. Logging Failures
10. SSRF

---

## 参考文献

1. [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/)
2. [OWASP Top 10](https://owasp.org/www-project-top-ten/)
3. [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
