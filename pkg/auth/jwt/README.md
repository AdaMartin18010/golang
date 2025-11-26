# JWT 认证框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [JWT 认证框架](#jwt-认证框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
  - [3. 使用示例](#3-使用示例)
  - [4. 最佳实践](#4-最佳实践)

---

## 1. 概述

JWT 认证框架提供了完整的 JWT Token 管理功能：

- ✅ **Token生成**: 支持 Access Token 和 Refresh Token
- ✅ **Token验证**: 完整的Token验证机制
- ✅ **Token刷新**: 支持Token刷新
- ✅ **多种签名算法**: 支持 HS256/HS384/HS512 和 RS256/RS384/RS512
- ✅ **Claims扩展**: 支持自定义Claims

---

## 2. 核心功能

### 2.1 配置

```go
type Config struct {
    SecretKey       string        // 密钥（HMAC）
    PrivateKey      *rsa.PrivateKey // 私钥（RSA）
    PublicKey       *rsa.PublicKey  // 公钥（RSA）
    SigningMethod   string        // 签名方法
    AccessTokenTTL  time.Duration // Access Token 过期时间
    RefreshTokenTTL time.Duration // Refresh Token 过期时间
    Issuer          string        // 签发者
    Audience        string        // 受众
}
```

### 2.2 Claims

```go
type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    Email    string   `json:"email,omitempty"`
    jwt.RegisteredClaims
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/auth/jwt"

// 创建JWT管理器
config := jwt.Config{
    SecretKey:      "your-secret-key",
    SigningMethod:  "HS256",
    AccessTokenTTL: 15 * time.Minute,
    RefreshTokenTTL: 7 * 24 * time.Hour,
    Issuer:         "your-app",
    Audience:       "your-audience",
}

j, err := jwt.NewJWT(config)
if err != nil {
    // 处理错误
}

// 生成Access Token
accessToken, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")

// 生成Refresh Token
refreshToken, err := j.GenerateRefreshToken("user-123")

// 验证Token
claims, err := j.ValidateToken(accessToken)
if err != nil {
    // 处理错误
}
```

### 3.2 在Handler中使用

```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // 验证用户凭据
    user, err := h.userService.Authenticate(r.Context(), email, password)
    if err != nil {
        response.Error(w, http.StatusUnauthorized, err)
        return
    }

    // 生成Token
    accessToken, err := h.jwt.GenerateAccessToken(
        user.ID,
        user.Username,
        user.Roles,
        user.Email,
    )
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err)
        return
    }

    refreshToken, err := h.jwt.GenerateRefreshToken(user.ID)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err)
        return
    }

    response.Success(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}
```

### 3.3 Token刷新

```go
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req RefreshTokenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest,
            errors.NewInvalidInputError("invalid request"))
        return
    }

    accessToken, refreshToken, err := h.jwt.RefreshToken(req.RefreshToken)
    if err != nil {
        response.Error(w, http.StatusUnauthorized, err)
        return
    }

    response.Success(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}
```

### 3.4 RSA签名

```go
import (
    "crypto/rand"
    "crypto/rsa"
)

// 生成RSA密钥对
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
    // 处理错误
}

config := jwt.Config{
    PrivateKey:     privateKey,
    PublicKey:      &privateKey.PublicKey,
    SigningMethod:  "RS256",
    AccessTokenTTL: 15 * time.Minute,
    Issuer:         "your-app",
    Audience:       "your-audience",
}

j, err := jwt.NewJWT(config)
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用HTTPS**: 在生产环境中始终使用HTTPS
2. **安全的密钥**: 使用足够长的随机密钥
3. **合理的过期时间**: Access Token 设置较短的过期时间（15分钟）
4. **Refresh Token**: 使用Refresh Token延长会话
5. **Token存储**: 在客户端安全存储Token（HttpOnly Cookie或安全存储）

### 4.2 DON'Ts ❌

1. **不要暴露密钥**: 不要在代码中硬编码密钥
2. **不要存储敏感信息**: Token中不要存储敏感信息
3. **不要使用过长的过期时间**: Access Token过期时间不要过长
4. **不要忽略Token验证**: 始终验证Token的有效性

---

## 5. 相关资源

- [认证授权中间件](../../internal/interfaces/http/middleware/auth/README.md)
- [统一错误处理框架](../../errors/README.md)
- [框架拓展计划](../../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
