# 零信任安全架构概念体系

## 目录

- [零信任安全架构概念体系](#零信任安全架构概念体系)
  - [目录](#目录)
  - [一、核心概念本体论](#一核心概念本体论)
  - [二、零信任公理定理](#二零信任公理定理)
    - [公理 1: 零信任基本公理](#公理-1-零信任基本公理)
    - [公理 2: 最小权限公理](#公理-2-最小权限公理)
    - [定理 1: 身份验证强度定理](#定理-1-身份验证强度定理)
    - [定理 2: 动态授权安全性定理](#定理-2-动态授权安全性定理)
    - [定理 3: 密钥轮换安全定理](#定理-3-密钥轮换安全定理)
  - [三、OAuth 2.0 + OIDC 流程详解](#三oauth-20--oidc-流程详解)
    - [JWT 结构详解](#jwt-结构详解)
  - [四、RBAC vs ABAC 决策矩阵](#四rbac-vs-abac-决策矩阵)
  - [五、Vault 架构详解](#五vault-架构详解)
  - [六、RBAC 实现示例](#六rbac-实现示例)
  - [七、ABAC 实现示例](#七abac-实现示例)
  - [八、Vault 集成示例](#八vault-集成示例)
  - [九、安全架构决策树](#九安全架构决策树)
  - [十、常见反模式](#十常见反模式)
    - [反模式 1: JWT 误用](#反模式-1-jwt-误用)
    - [反模式 2: 硬编码密钥](#反模式-2-硬编码密钥)
    - [反模式 3: 权限检查缺失](#反模式-3-权限检查缺失)

## 一、核心概念本体论

```text
Zero Trust Architecture (ZTA)
├── 核心原则
│   ├── 永不信任，始终验证 (Never Trust, Always Verify)
│   ├── 假设已攻破 (Assume Breach)
│   └── 最小权限 (Least Privilege)
│
├── 认证 (Authentication)
│   ├── 因素类型
│   │   ├── 所知 (Something you know): 密码、PIN
│   │   ├── 所有 (Something you have): 手机、Token、证书
│   │   └── 所是 (Something you are): 指纹、人脸、虹膜
│   │
│   ├── 协议
│   │   ├── OAuth 2.0: 授权框架
│   │   ├── OpenID Connect (OIDC): 身份层
│   │   ├── SAML: 企业SSO
│   │   └── mTLS: 双向TLS
│   │
│   └── JWT (JSON Web Token)
│       ├── Header: 算法信息
│       ├── Payload: 声明 (Claims)
│       └── Signature: 签名验证
│
├── 授权 (Authorization)
│   ├── RBAC (Role-Based Access Control)
│   │   ├── 角色 (Role): 权限集合
│   │   ├── 权限 (Permission): 操作许可
│   │   └── 继承: 角色层次结构
│   │
│   ├── ABAC (Attribute-Based Access Control)
│   │   ├── 主体属性 (Subject): 用户角色、部门
│   │   ├── 资源属性 (Resource): 数据分类、所有者
│   │   ├── 操作属性 (Action): 读、写、删
│   │   └── 环境属性 (Environment): 时间、地点、设备
│   │
│   └── Policy Decision/Enforcement Point
│       ├── PDP (Policy Decision Point): 决策引擎
│       └── PEP (Policy Enforcement Point): 执行点
│
├── 密钥管理
│   ├── Vault (HashiCorp Vault)
│   │   ├── 动态密钥: 数据库凭据、云凭证
│   │   ├── 静态密钥: API Key、证书
│   │   ├── 加密服务: Transit引擎
│   │   └── PKI: 证书管理
│   │
│   └── 密钥轮换
│       ├── 自动轮换: 定时触发
│       ├── 版本管理: 多版本共存
│       └── 审计追踪: 完整日志
│
└── 安全监控
    ├── SIEM (Security Information and Event Management)
    ├── SOAR (Security Orchestration, Automation and Response)
    └── Threat Intelligence
```

## 二、零信任公理定理

### 公理 1: 零信任基本公理

```text
定义: 网络位置不等于信任
数学表达: ∀resource, ∀network, Trust(resource) ≠ Location(network)
推论: 内网用户和外网用户需要相同的验证强度
```

### 公理 2: 最小权限公理

```text
定义: 主体只能访问完成任务所需的最小资源集合
数学表达: ∀subject, Access(subject) = {r | Need(subject, r) ∧ TimeValid(t)}
约束: 权限随时间和上下文动态调整
```

### 定理 1: 身份验证强度定理

```text
条件: 系统实施多因素认证 (MFA)
证明:
  1. 单因素被攻破的概率: P(breach_single) = p
  2. 双因素同时被攻破的概率: P(breach_mfa) = p₁ × p₂
  3. 假设 p = 0.1, p₁ = p₂ = 0.1
  4. P(breach_mfa) = 0.01 << 0.1
结论: Security(MFA) >> Security(Single Factor)
```

### 定理 2: 动态授权安全性定理

```text
条件: 使用 ABAC 动态授权 vs 静态 RBAC
证明:
  1. RBAC: Access = f(Role), 静态，粗粒度
  2. ABAC: Access = f(S,R,A,E), 动态，细粒度
  3. 攻击面: AttackSurface(ABAC) ⊂ AttackSurface(RBAC)
结论: Security(ABAC) > Security(RBAC) for complex scenarios
```

### 定理 3: 密钥轮换安全定理

```text
条件: 定期轮换密钥，旧密钥自动失效
证明:
  1. 密钥暴露风险窗口: Window = T_compromise + T_detection + T_response
  2. 轮换周期: T_rotation << Window
  3. 如果密钥在轮换周期内暴露，下次轮换后失效
结论: 密钥轮换将长期风险转化为短期风险
```

## 三、OAuth 2.0 + OIDC 流程详解

```text
┌─────────────────────────────────────────────────────────────────────┐
│                  OAuth 2.0 Authorization Code Flow                  │
│                      + PKCE (Proof Key for Code Exchange)           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌──────────┐                                    ┌──────────┐     │
│   │  Client  │                                    │   Auth   │     │
│   │ (App)    │                                    │  Server  │     │
│   └────┬─────┘                                    └────┬─────┘     │
│        │                                               │           │
│        │ 1. Authorization Request (+ code_challenge)   │           │
│        │ ─────────────────────────────────────────────►│           │
│        │                                               │           │
│        │ 2. User Authentication & Consent              │           │
│        │ ◄─────────────────────────────────────────────│           │
│        │                                               │           │
│        │ 3. Authorization Code                         │           │
│        │ ◄─────────────────────────────────────────────│           │
│        │                                               │           │
│        │ 4. Token Request (+ code_verifier)            │           │
│        │ ─────────────────────────────────────────────►│           │
│        │                                               │           │
│        │ 5. Access Token + ID Token (OIDC)             │           │
│        │ ◄─────────────────────────────────────────────│           │
│        │                                               │           │
│   ┌────┴─────┐                                    ┌────┴─────┐     │
│   │ Resource │◄──────── 6. API Call + Token ──────►│   User   │     │
│   │  Server  │                                    │  Info    │     │
│   └──────────┘                                    └──────────┘     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### JWT 结构详解

```text
┌─────────────────────────────────────────────────────────────────────┐
│                         JWT Token 结构                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   eyJhbGciOiJSUzI1NiIs...  // Header (Base64Url)                   │
│   .                                                                   │
│   eyJzdWIiOiIxMjM0NTY3O...  // Payload (Base64Url)                  │
│   .                                                                   │
│   SflKxwRJSMeKKF2QT4fwpM...  // Signature                           │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Header:                                                           │
│   {                                                                 │
│     "alg": "RS256",      // 签名算法                               │
│     "typ": "JWT",        // Token类型                              │
│     "kid": "2024-01"     // 密钥ID                                 │
│   }                                                                 │
│                                                                     │
│   Payload (Claims):                                                 │
│   {                                                                 │
│     "sub": "user123",    // 主题 (用户ID)                          │
│     "iss": "auth.server",// 签发者                                │
│     "aud": "api.server", // 接收者                                │
│     "exp": 1704067200,   // 过期时间                               │
│     "iat": 1703980800,   // 签发时间                               │
│     "scope": "read write",// 权限范围                              │
│     "roles": ["user", "admin"]                                     │
│   }                                                                 │
│                                                                     │
│   Signature:                                                        │
│   RSASHA256(                                                        │
│     base64Url(header) + "." + base64Url(payload),                   │
│     private_key                                                     │
│   )                                                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 四、RBAC vs ABAC 决策矩阵

```text
选择访问控制模型
│
├─ 系统复杂度?
│   ├─ 简单 (少量角色，固定权限) ─────────────► RBAC
│   │   ├─ 优点: 简单易懂，性能好
│   │   └─ 场景: 内部管理系统
│   │
│   └─ 复杂 (动态权限，多维度控制) ────────────► ABAC
│       ├─ 优点: 细粒度，灵活
│       └─ 场景: 多云环境，数据分级
│
├─ 性能要求?
│   ├─ 极高性能 (>10k QPS) ───────────────────► RBAC + 缓存
│   └─ 可接受稍低性能 ────────────────────────► ABAC + OPA
│
├─ 合规要求?
│   ├─ 严格审计 (金融/医疗) ──────────────────► ABAC + 策略即代码
│   └─ 一般合规 ──────────────────────────────► RBAC
│
└─ 团队能力?
    ├─ 有安全专家团队 ────────────────────────► ABAC + 自定义策略引擎
    └─ 一般开发团队 ──────────────────────────► RBAC
```

## 五、Vault 架构详解

```text
┌─────────────────────────────────────────────────────────────────────┐
│                      HashiCorp Vault 架构                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │                        Vault Server                          │  │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │  │
│   │  │   HTTP API  │  │   Auth      │  │     Secret Engines  │  │  │
│   │  │   (REST)    │──│   Methods   │──│  ┌───────────────┐  │  │  │
│   │  └─────────────┘  │             │  │  │ KV v1/v2      │  │  │  │
│   │                   │ • Token     │  │  │ Database      │  │  │  │
│   │  ┌─────────────┐  │ • Kubernetes│  │  │ Transit       │  │  │  │
│   │  │   Storage   │  │ • AppRole   │  │  │ PKI           │  │  │  │
│   │  │  (Encrypted)│  │ • OIDC      │  │  │ AWS/Azure/GCP │  │  │  │
│   │  │             │  │ • LDAP      │  │  │ SSH           │  │  │  │
│   │  │ • Consul    │  │ • TLS Cert  │  │  │ ...           │  │  │  │
│   │  │ • Integrated│  │             │  │  └───────────────┘  │  │  │
│   │  │ • Raft      │  └─────────────┘  └─────────────────────┘  │  │
│   │  └─────────────┘                                            │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                                                                     │
│   Authentication Flow:                                              │
│   ┌──────────┐    1. Login    ┌──────────┐    2. Issue Token      │
│   │  Client  │ ──────────────►│   Vault  │ ──────────────────┐    │
│   └──────────┘                └──────────┘                   │    │
│        ▲                                                        │    │
│        │                    3. Access Secret                    │    │
│        └────────────────────────────────────────────────────────┘    │
│                              (Token + Path)                         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 六、RBAC 实现示例

```go
// RBAC 核心概念定义

// Permission 表示一个权限
type Permission struct {
    Resource string // 资源: user, order, product
    Action   string // 操作: create, read, update, delete
}

// Role 表示一个角色
type Role struct {
    ID          string
    Name        string
    Permissions []Permission
    ParentID    string // 支持角色继承
}

// RBAC 引擎
type RBAC struct {
    roles       map[string]*Role
    userRoles   map[string][]string // userID -> roleIDs
    permissions map[string]map[string]bool // 缓存: roleID -> permissionKey
}

// CheckPermission 检查用户是否有权限
func (r *RBAC) CheckPermission(userID string, resource, action string) bool {
    roleIDs := r.userRoles[userID]
    requiredPerm := resource + ":" + action

    for _, roleID := range roleIDs {
        if r.hasPermission(roleID, requiredPerm) {
            return true
        }
    }
    return false
}

// hasPermission 递归检查角色权限（支持继承）
func (r *RBAC) hasPermission(roleID, permission string) bool {
    role := r.roles[roleID]
    if role == nil {
        return false
    }

    // 检查直接权限
    if r.permissions[roleID][permission] {
        return true
    }

    // 递归检查父角色
    if role.ParentID != "" {
        return r.hasPermission(role.ParentID, permission)
    }

    return false
}

// 使用示例
func main() {
    rbac := NewRBAC()

    // 定义角色
    adminRole := &Role{
        ID:   "role:admin",
        Name: "Admin",
        Permissions: []Permission{
            {Resource: "*", Action: "*"}, // 通配符权限
        },
    }

    userRole := &Role{
        ID:       "role:user",
        Name:     "User",
        ParentID: "", // 无继承
        Permissions: []Permission{
            {Resource: "order", Action: "create"},
            {Resource: "order", Action: "read"},
        },
    }

    rbac.AddRole(adminRole)
    rbac.AddRole(userRole)
    rbac.AssignRole("user:123", "role:user")

    // 检查权限
    if rbac.CheckPermission("user:123", "order", "create") {
        fmt.Println("允许创建订单")
    }
}
```

## 七、ABAC 实现示例

```go
// ABAC 核心概念定义

// Subject 主体属性
type Subject struct {
    ID       string
    Role     string
    Department string
    ClearanceLevel int
}

// Resource 资源属性
type Resource struct {
    ID         string
    Type       string
    Owner      string
    Classification string // public, internal, confidential, secret
}

// Action 操作属性
type Action struct {
    Type string // read, write, delete, execute
}

// Environment 环境属性
type Environment struct {
    Time        time.Time
    Location    string
    DeviceTrust int
    ConnectionType string // internal, vpn, external
}

// Policy 策略
type Policy struct {
    ID          string
    Name        string
    Description string
    Effect      string // allow, deny
    Rules       []Rule
    Priority    int
}

// Rule 规则
type Rule struct {
    SubjectMatcher    func(Subject) bool
    ResourceMatcher   func(Resource) bool
    ActionMatcher     func(Action) bool
    EnvironmentMatcher func(Environment) bool
}

// ABAC 引擎
type ABAC struct {
    policies []Policy
}

// Evaluate 评估访问请求
func (a *ABAC) Evaluate(subject Subject, resource Resource, action Action, env Environment) (bool, string) {
    // 按优先级排序
    sortedPolicies := a.sortByPriority()

    for _, policy := range sortedPolicies {
        if a.matchPolicy(policy, subject, resource, action, env) {
            if policy.Effect == "deny" {
                return false, fmt.Sprintf("Policy %s denied", policy.ID)
            }
            return true, fmt.Sprintf("Policy %s allowed", policy.ID)
        }
    }

    // 默认拒绝
    return false, "No matching policy found"
}

// 策略示例
func createConfidentialDataPolicy() Policy {
    return Policy{
        ID:     "policy:confidential-data",
        Name:   "Confidential Data Access",
        Effect: "allow",
        Rules: []Rule{
            {
                // 主体：必须有机密级别 >= 3 或者是资源所有者
                SubjectMatcher: func(s Subject) bool {
                    return s.ClearanceLevel >= 3
                },
                // 资源：机密级别
                ResourceMatcher: func(r Resource) bool {
                    return r.Classification == "confidential"
                },
                // 操作：读取或写入
                ActionMatcher: func(a Action) bool {
                    return a.Type == "read" || a.Type == "write"
                },
                // 环境：工作时间，内部网络
                EnvironmentMatcher: func(e Environment) bool {
                    hour := e.Time.Hour()
                    isWorkHours := hour >= 9 && hour <= 18
                    isInternal := e.ConnectionType == "internal" || e.ConnectionType == "vpn"
                    return isWorkHours && isInternal
                },
            },
        },
        Priority: 100,
    }
}

// 使用示例
func main() {
    abac := NewABAC()
    abac.AddPolicy(createConfidentialDataPolicy())

    subject := Subject{
        ID:             "user:123",
        Role:           "analyst",
        Department:     "finance",
        ClearanceLevel: 4,
    }

    resource := Resource{
        ID:             "doc:456",
        Type:           "document",
        Owner:          "user:789",
        Classification: "confidential",
    }

    action := Action{Type: "read"}

    env := Environment{
        Time:           time.Now(),
        Location:       "office",
        DeviceTrust:    5,
        ConnectionType: "internal",
    }

    allowed, reason := abac.Evaluate(subject, resource, action, env)
    fmt.Printf("Access %v: %s\n", allowed, reason)
}
```

## 八、Vault 集成示例

```go
package vault

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/vault/api"
)

// VaultClient Vault客户端封装
type VaultClient struct {
    client *api.Client
    config *Config
}

// Config Vault配置
type Config struct {
    Address    string
    AuthMethod string // token, kubernetes, approle
    Token      string
    RoleID     string // for AppRole
    SecretID   string // for AppRole
}

// NewClient 创建Vault客户端
func NewClient(cfg *Config) (*VaultClient, error) {
    apiConfig := api.DefaultConfig()
    apiConfig.Address = cfg.Address

    client, err := api.NewClient(apiConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create vault client: %w", err)
    }

    vc := &VaultClient{client: client, config: cfg}

    // 认证
    if err := vc.authenticate(); err != nil {
        return nil, fmt.Errorf("failed to authenticate: %w", err)
    }

    return vc, nil
}

// authenticate 执行认证
func (v *VaultClient) authenticate() error {
    switch v.config.AuthMethod {
    case "token":
        v.client.SetToken(v.config.Token)
        return nil

    case "approle":
        return v.authAppRole()

    case "kubernetes":
        return v.authKubernetes()

    default:
        return fmt.Errorf("unsupported auth method: %s", v.config.AuthMethod)
    }
}

// authAppRole AppRole认证
func (v *VaultClient) authAppRole() error {
    data := map[string]interface{}{
        "role_id":   v.config.RoleID,
        "secret_id": v.config.SecretID,
    }

    secret, err := v.client.Logical().Write("auth/approle/login", data)
    if err != nil {
        return err
    }

    v.client.SetToken(secret.Auth.ClientToken)
    return nil
}

// GetSecret 从KV v2获取密钥
func (v *VaultClient) GetSecret(ctx context.Context, path string) (map[string]interface{}, error) {
    secret, err := v.client.KVv2("secret").Get(ctx, path)
    if err != nil {
        return nil, fmt.Errorf("failed to get secret: %w", err)
    }

    return secret.Data, nil
}

// GetDynamicDBCredentials 获取动态数据库凭据
func (v *VaultClient) GetDynamicDBCredentials(ctx context.Context, role string) (*DBCredentials, error) {
    secret, err := v.client.Logical().ReadWithContext(ctx, fmt.Sprintf("database/creds/%s", role))
    if err != nil {
        return nil, fmt.Errorf("failed to get db credentials: %w", err)
    }

    return &DBCredentials{
        Username:      secret.Data["username"].(string),
        Password:      secret.Data["password"].(string),
        LeaseID:       secret.LeaseID,
        LeaseDuration: time.Duration(secret.LeaseDuration) * time.Second,
    }, nil
}

// Encrypt 使用Transit引擎加密
func (v *VaultClient) Encrypt(ctx context.Context, keyName string, plaintext []byte) (string, error) {
    data := map[string]interface{}{
        "plaintext": base64.StdEncoding.EncodeToString(plaintext),
    }

    secret, err := v.client.Logical().WriteWithContext(ctx,
        fmt.Sprintf("transit/encrypt/%s", keyName), data)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt: %w", err)
    }

    return secret.Data["ciphertext"].(string), nil
}

// DBCredentials 数据库凭据
type DBCredentials struct {
    Username      string
    Password      string
    LeaseID       string
    LeaseDuration time.Duration
}

// RenewLease 续期租约
func (v *VaultClient) RenewLease(ctx context.Context, leaseID string, increment int) error {
    _, err := v.client.Sys().RenewWithContext(ctx, leaseID, increment)
    return err
}
```

## 九、安全架构决策树

```text
设计零信任安全架构
│
├─ 认证方式?
│   ├─ 传统Web应用 ───────────────► Session + Cookie
│   ├─ 移动应用 ──────────────────► OAuth 2.0 + PKCE
│   ├─ 服务间通信 ────────────────► mTLS + Service Account
│   └─ 混合场景 ──────────────────► OAuth 2.0 + JWT
│
├─ 授权粒度?
│   ├─ 粗粒度 (角色) ─────────────► RBAC
│   ├─ 细粒度 (属性) ─────────────► ABAC
│   └─ 混合 ─────────────────────► RBAC + ABAC
│
├─ 密钥管理?
│   ├─ 云环境 ────────────────────► 云厂商KMS + Vault
│   ├─ 混合云 ────────────────────► Vault集群
│   └─ 简单场景 ──────────────────► 配置文件 + 环境变量
│
├─ 会话管理?
│   ├─ 无状态 ────────────────────► JWT (短有效期)
│   ├─ 有状态 ────────────────────► Redis Session
│   └─ 混合 ─────────────────────► JWT + Refresh Token
│
└─ 审计要求?
    ├─ 严格审计 ──────────────────► 结构化日志 + SIEM
    └─ 一般审计 ──────────────────► 应用日志
```

## 十、常见反模式

### 反模式 1: JWT 误用

```go
// ❌ 错误：在 JWT 中存储敏感信息
claims := jwt.MapClaims{
    "user_id": userID,
    "password": password,      // 错误！JWT 可解码
    "credit_card": cardNumber, // 错误！
    "exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 错误！有效期太长
}

// ✅ 正确：JWT 只存身份标识，短期有效
claims := jwt.MapClaims{
    "sub": userID,
    "roles": []string{"user"},
    "iat": time.Now().Unix(),
    "exp": time.Now().Add(time.Hour * 2).Unix(), // 2小时有效期
}
```

### 反模式 2: 硬编码密钥

```go
// ❌ 错误：硬编码密钥
const APIKey = "sk-1234567890abcdef"  // 安全风险！

// ✅ 正确：从 Vault 获取
apiKey, err := vaultClient.GetSecret(ctx, "api-keys/payment-gateway")
```

### 反模式 3: 权限检查缺失

```go
// ❌ 错误：没有权限检查
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    db.Exec("DELETE FROM users WHERE id = ?", userID)  // 危险！
}

// ✅ 正确：强制权限检查
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)
    targetID := r.URL.Query().Get("id")

    // 检查权限
    if !rbac.CheckPermission(userID, "user", "delete") {
        http.Error(w, "Forbidden", 403)
        return
    }

    // 额外检查：只能删除自己或管理员
    if userID != targetID && !rbac.HasRole(userID, "admin") {
        http.Error(w, "Forbidden", 403)
        return
    }

    db.Exec("DELETE FROM users WHERE id = ?", targetID)
}
```

---

**参考来源**:

- NIST SP 800-207: Zero Trust Architecture
- OAuth 2.0 RFC 6749
- OpenID Connect Core 1.0
- Vault Documentation - HashiCorp
