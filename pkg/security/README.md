# 安全模块

**版本**: v1.0
**更新日期**: 2025-12-03
**状态**: 🔄 实现中

---

## 🎯 功能概览

本模块提供企业级安全功能：

1. **OAuth2/OIDC** - 标准认证协议
2. **RBAC** - 基于角色的访问控制
3. **ABAC** - 基于属性的访问控制（计划中）
4. **JWT** - JSON Web Token
5. **Vault** - 密钥管理（计划中）

---

## 🏗️ 模块结构

```text
pkg/security/
├── oauth2/
│   ├── provider.go      # OAuth2 提供者 ✅
│   ├── oidc.go          # OIDC 实现 ✅
│   └── README.md
├── rbac/
│   ├── rbac.go          # RBAC 核心 ✅
│   ├── middleware.go    # HTTP 中间件 ✅
│   └── README.md
├── jwt/
│   ├── jwt.go           # JWT 实现
│   └── README.md
├── vault/
│   ├── client.go        # Vault 客户端
│   └── README.md
└── README.md            # 本文档
```

---

## 🚀 快速开始

### OAuth2/OIDC 认证

```go
import (
    "github.com/yourusername/golang/pkg/security/oauth2"
)

// 创建 OIDC 提供者 (Google)
provider, err := oauth2.NewGoogleOIDCProvider(
    ctx,
    "your-client-id",
    "your-client-secret",
    "http://localhost:8080/callback",
)

// 生成授权 URL
authURL := provider.AuthorizationURL("random-state")

// 交换授权码
token, err := provider.Exchange(ctx, code)

// 验证 ID Token
claims, err := provider.VerifyIDToken(ctx, token.IDToken)
```

### RBAC 授权

```go
import (
    "github.com/yourusername/golang/pkg/security/rbac"
)

// 创建 RBAC
rbacSystem := rbac.NewRBAC()

// 初始化默认角色
rbacSystem.InitializeDefaultRoles()

// 检查权限
hasPermission, err := rbacSystem.CheckPermission(
    ctx,
    []string{"admin"},  // 用户角色
    "user",             // 资源
    "create",           // 操作
)

// 使用中间件
middleware := rbac.NewMiddleware(rbacSystem)
router.Use(middleware.RequirePermission("user", "read"))
```

---

## 📊 实现状态

| 功能 | 状态 | 优先级 | 预计完成 |
|------|------|--------|---------|
| **OAuth2** | ✅ 基础实现 | P0 | 完成 |
| **OIDC** | ✅ 基础实现 | P0 | 完成 |
| **RBAC** | ✅ 基础实现 | P0 | 完成 |
| **RBAC 中间件** | ✅ 完成 | P0 | 完成 |
| **JWT** | ⏳ 待实现 | P0 | 本周 |
| **ABAC** | ⏳ 待实现 | P1 | 下周 |
| **Vault** | ⏳ 待实现 | P1 | 下周 |
| **测试** | ⏳ 待实现 | P0 | 本周 |

---

## 🔐 安全最佳实践

### 1. OAuth2/OIDC

- ✅ 使用 PKCE (RFC 7636)
- ✅ 验证 state 参数防止 CSRF
- ✅ 使用 HTTPS（生产环境）
- ✅ 安全存储 client_secret
- ✅ 验证 ID Token 签名

### 2. RBAC

- ✅ 最小权限原则
- ✅ 角色继承支持
- ✅ 权限细粒度控制
- ✅ 线程安全实现

### 3. 令牌管理

- ⏳ 短期访问令牌（15分钟）
- ⏳ 长期刷新令牌（7天）
- ⏳ 令牌轮换
- ⏳ 令牌撤销

---

## 📚 参考资源

### 标准和规范

- [RFC 6749 - OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7636 - PKCE](https://datatracker.ietf.org/doc/html/rfc7636)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [NIST RBAC](https://csrc.nist.gov/projects/role-based-access-control)

### Go 库

- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- [github.com/coreos/go-oidc](https://github.com/coreos/go-oidc)
- [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)

---

## 🎯 下一步

### 本周任务

1. **JWT 实现**
   - 生成和验证 JWT
   - 支持 RS256/ES256
   - 刷新令牌机制

2. **测试**
   - OAuth2/OIDC 单元测试
   - RBAC 单元测试
   - 集成测试

3. **文档**
   - 使用指南
   - 最佳实践
   - 故障排查

### 下周任务

1. **ABAC 实现**
   - 策略引擎
   - 属性评估

2. **Vault 集成**
   - 密钥存储
   - 密钥轮换

---

**状态**: 🔄 快速推进中
**目标**: 安全性 6/10 → 9/10
**优先级**: P0 (最高)
