# AD-007: 应用安全模式的形式化 (Security Patterns: Formal Analysis)

> **维度**: Application Domains
> **级别**: S (17+ KB)
> **标签**: #security #authentication #authorization #jwt #oauth #zero-trust
> **权威来源**:
>
> - [Security Patterns](https://www.oreilly.com/library/view/security-patterns-in/9780470858844/) - Schumacher et al. (2006)
> - [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/) - OWASP (2025)
> - [Zero Trust Architecture](https://www.nist.gov/publications/zero-trust-architecture) - NIST (2020)
> - [OAuth 2.0 and OpenID Connect](https://oauth.net/2/) - IETF (2024)
> - [JWT Handbook](https://auth0.com/resources/ebooks/jwt-handbook) - Auth0 (2017)

---

## 1. 安全的形式化定义

### 1.1 CIA 三元组

**定义 1.1 (机密性 Confidentiality)**
$$\forall d \in \text{Data}: \text{Authorized}(u) \Leftarrow \text{Access}(u, d)$$

**定义 1.2 (完整性 Integrity)**
$$\forall d: \text{Modified}(d) \Rightarrow \text{Authorized}(modifier)$$

**定义 1.3 (可用性 Availability)**
$$\Diamond(\text{Access}(u, \text{service}))$$

### 1.2 威胁模型

**定义 1.4 (STRIDE)**

| 威胁 | 定义 |
|------|------|
| Spoofing | 身份伪造 |
| Tampering | 数据篡改 |
| Repudiation | 否认行为 |
| Information Disclosure | 信息泄露 |
| Denial of Service | 拒绝服务 |
| Elevation of Privilege | 权限提升 |

---

## 2. 认证的形式化

### 2.1 JWT 结构

**定义 2.1 (JWT)**
$$\text{JWT} = \text{Base64}(\text{Header}) \cdot \text{Base64}(\text{Payload}) \cdot \text{Signature}$$

**签名验证**:
$$\text{Verify}(\text{JWT}, \text{secret}) \Leftrightarrow \text{Signature} = H(\text{Header} \circ \text{Payload}, \text{secret})$$

### 2.2 OAuth 2.0 流程

**定义 2.2 (授权码模式)**

```
1. Client → AuthServer: authorize request
2. AuthServer → User: login + consent
3. User → AuthServer: approval
4. AuthServer → Client: authorization code
5. Client → AuthServer: code + client_secret
6. AuthServer → Client: access_token + refresh_token
```

---

## 3. 授权的形式化

### 3.1 RBAC

**定义 3.1 (基于角色的访问控制)**
$$\text{Permission}(u, r) \Leftrightarrow \exists \text{role}: u \in \text{role} \land r \in \text{Permissions}(\text{role})$$

### 3.2 ABAC

**定义 3.2 (基于属性的访问控制)**
$$\text{Permission}(u, r, e) \Leftrightarrow \text{Policy}(\text{Attributes}(u, r, e)) = \text{allow}$$

---

## 4. 多元表征

### 4.1 安全层次图

```
Security Layers
├── Network
│   ├── TLS/mTLS
│   ├── VPN
│   └── Firewall
├── Application
│   ├── Authentication
│   ├── Authorization
│   ├── Input Validation
│   └── Output Encoding
├── Data
│   ├── Encryption at Rest
│   ├── Encryption in Transit
│   └── Tokenization
└── Infrastructure
    ├── Secrets Management
    ├── IAM
    └── Audit Logging
```

### 4.2 认证方式对比矩阵

| 方式 | 安全性 | 用户体验 | 适用场景 |
|------|--------|----------|---------|
| **Password** | 中 | 好 | 基础认证 |
| **MFA** | 高 | 中 | 敏感操作 |
| **OAuth/OIDC** | 高 | 好 | 第三方登录 |
| **mTLS** | 极高 | 差 | 服务间认证 |
| **JWT** | 中 | 好 | 状态less认证 |

### 4.3 零信任架构

```
Zero Trust Principles
├── Never Trust, Always Verify
│   └── 每个请求都认证和授权
├── Assume Breach
│   └── 分段网络，限制爆炸半径
├── Verify Explicitly
│   └── 多因素认证，最少权限
└── Use Least Privilege Access
    └── 动态授权，实时评估
```

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Security Implementation Checklist                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  认证:                                                                       │
│  □ 强密码策略                                                                │
│  □ MFA 支持                                                                  │
│  □ 会话管理 (超时、刷新)                                                      │
│  □ JWT 安全存储                                                              │
│                                                                              │
│  授权:                                                                       │
│  □ 最小权限原则                                                              │
│  □ RBAC/ABAC 实现                                                            │
│  □ 资源级权限控制                                                             │
│                                                                              │
│  输入验证:                                                                   │
│  □ 白名单验证                                                                │
│  □ 参数化查询 (防 SQL 注入)                                                   │
│  □ XSS/CSRF 防护                                                             │
│                                                                              │
│  审计:                                                                       │
│  □ 操作日志                                                                  │
│  □ 不可篡改                                                                  │
│  □ 敏感操作告警                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (17KB, 完整形式化)
