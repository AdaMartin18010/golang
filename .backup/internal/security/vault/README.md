# HashiCorp Vault 集成

本包提供了与 HashiCorp Vault 集成的完整实现，支持密钥管理、动态数据库凭据、加密服务和密钥轮换功能。

## 文件结构

| 文件 | 说明 |
|------|------|
| `client.go` | Vault 客户端实现，支持多种认证方式 |
| `secret.go` | 密钥管理（KV v1/v2）、动态数据库凭据 |
| `rotation.go` | 密钥轮换策略和自动轮换调度 |
| `encryption.go` | 基于 Transit 密钥引擎的加密/解密服务 |
| `vault_test.go` | 单元测试（使用 mock 实现） |

## 功能特性

- **KV 密钥引擎**：支持 KV v1 和 v2，包括版本管理、元数据管理
- **动态密钥**：支持数据库动态凭据，自动租约管理
- **加密服务**：基于 Transit 密钥引擎，支持 AES、ChaCha20、RSA、ECDSA、Ed25519
- **密钥轮换**：支持自动和手动密钥轮换，轮换历史记录
- **多种认证**：Token、Kubernetes、AppRole、TLS 证书
- **连接管理**：连接重试、超时配置、健康检查

## 快速开始

### 1. 基本配置

```yaml
# configs/config.yaml
vault:
  address: "https://vault.example.com:8200"
  auth_method: "token"
  token: "your-vault-token"
  namespace: ""  # 企业版命名空间
  max_retries: 3
  timeout: 30s
  kv_version: 2
```

### 2. 创建客户端

```go
package main

import (
    "context"
    "log"
    "github.com/yourusername/golang/internal/security/vault"
)

func main() {
    config := &vault.Config{
        Address:    "https://vault.example.com:8200",
        AuthMethod: vault.AuthMethodToken,
        Token:      "your-vault-token",
        MaxRetries: 3,
        Timeout:    30 * time.Second,
        KVVersion:  2,
    }

    client, err := vault.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 检查健康状态
    health, err := client.Health(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Vault version: %s, Initialized: %v", health.Version, health.Initialized)
}
```

### 3. 密钥管理

```go
ctx := context.Background()

// 存储密钥
err := client.KV().Put(ctx, "secret/myapp/database", map[string]interface{}{
    "host":     "localhost",
    "port":     5432,
    "username": "dbuser",
    "password": "dbpass",
})

// 读取密钥
secret, err := client.KV().Get(ctx, "secret/myapp/database")
if err != nil {
    log.Fatal(err)
}
log.Printf("Database host: %s", secret.Data["host"])

// 列出密钥
keys, err := client.KV().List(ctx, "secret/myapp")

// 删除密钥
err := client.KV().Delete(ctx, "secret/myapp/database")
```

### 4. 动态数据库凭据

```go
// 获取动态数据库凭据
creds, err := client.DB().GetCredentials(ctx, "db-readonly-role")
if err != nil {
    log.Fatal(err)
}

// 使用凭据连接数据库
dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=mydb", 
    creds.Username, creds.Password)

// 续期租约
if creds.Renewable {
    err := client.DB().RenewLease(ctx, creds.LeaseID, 3600)
}

// 撤销租约
defer client.DB().RevokeLease(ctx, creds.LeaseID)
```

### 5. 加密服务

```go
// 创建加密密钥
err := client.Encryption().CreateKey(ctx, "my-key", vault.KeyTypeAES256GCM, nil)

// 加密数据
plaintext := []byte("sensitive data")
result, err := client.Encryption().Encrypt(ctx, "my-key", plaintext, nil)
if err != nil {
    log.Fatal(err)
}
log.Printf("Ciphertext: %s", result.Ciphertext)

// 解密数据
decrypted, err := client.Encryption().Decrypt(ctx, "my-key", result.Ciphertext, nil)

// 批量加密
items := []*vault.BatchEncryptItem{
    {ID: "1", Plaintext: []byte("data1")},
    {ID: "2", Plaintext: []byte("data2")},
}
results, err := client.Encryption().EncryptBatch(ctx, "my-key", items)

// 签名数据
signResult, err := client.Encryption().Sign(ctx, "sign-key", data, &vault.SignOptions{
    HashAlgorithm: "sha2-256",
})

// 验证签名
valid, err := client.Encryption().Verify(ctx, "sign-key", data, signResult.Signature, nil)
```

### 6. 密钥轮换

```go
// 注册轮换策略
policy := &vault.RotationPolicy{
    Name:           "db-password-rotation",
    Path:           "secret/myapp/database",
    Interval:       24 * time.Hour,
    Enabled:        true,
    AutoRotate:     true,
    RetainVersions: 10,
}
err := client.Rotation().RegisterPolicy(policy)

// 启动自动轮换
err := client.Rotation().StartAutoRotation()

// 手动轮换
err := client.Rotation().Rotate(ctx, "secret/myapp/database")

// 查看轮换历史
history := client.Rotation().GetRotationHistory("secret/myapp/database")
for _, h := range history {
    log.Printf("Rotated at %v: v%d -> v%d", h.RotatedAt, h.OldVersion, h.NewVersion)
}

// 停止自动轮换
err := client.Rotation().StopAutoRotation()
```

## 认证方式

### Token 认证

```go
config := &vault.Config{
    Address:    "https://vault.example.com:8200",
    AuthMethod: vault.AuthMethodToken,
    Token:      os.Getenv("VAULT_TOKEN"),
}
```

### Kubernetes 认证

```go
config := &vault.Config{
    Address:      "https://vault.example.com:8200",
    AuthMethod:   vault.AuthMethodKubernetes,
    K8sRole:      "my-app-role",
    K8sTokenPath: "/var/run/secrets/kubernetes.io/serviceaccount/token",
}
```

### AppRole 认证

```go
config := &vault.Config{
    Address:     "https://vault.example.com:8200",
    AuthMethod:  vault.AuthMethodAppRole,
    AppRoleID:   os.Getenv("VAULT_APPROLE_ID"),
    AppSecretID: os.Getenv("VAULT_APPROLE_SECRET_ID"),
}
```

### TLS 证书认证

```go
config := &vault.Config{
    Address:            "https://vault.example.com:8200",
    AuthMethod:         vault.AuthMethodCert,
    TLSCertPath:        "/path/to/client.crt",
    TLSKeyPath:         "/path/to/client.key",
    CACertPath:         "/path/to/ca.crt",
    InsecureSkipVerify: false,
}
```

## 环境变量

```bash
# 基本配置
export VAULT_ADDR="https://vault.example.com:8200"

# Token 认证
export VAULT_TOKEN="your-token"

# Kubernetes 认证
export VAULT_K8S_ROLE="my-app-role"

# AppRole 认证
export VAULT_APPROLE_ID="role-id"
export VAULT_APPROLE_SECRET_ID="secret-id"
```

## 测试

```bash
# 运行单元测试
go test ./internal/security/vault/...

# 运行基准测试
go test -bench=. ./internal/security/vault/...

# 运行测试并生成覆盖率报告
go test -cover -coverprofile=coverage.out ./internal/security/vault/...
go tool cover -html=coverage.out -o coverage.html
```

## 注意事项

1. **安全性**：生产环境中切勿将 Vault Token 硬编码，应使用环境变量或 Kubernetes Secrets
2. **TLS**：生产环境必须启用 TLS，并正确配置 CA 证书
3. **租约管理**：动态凭据的租约需要及时续期或撤销，避免泄露
4. **错误处理**：所有操作都可能失败，需要正确处理错误
5. **上下文控制**：使用带超时的上下文控制操作，避免长时间阻塞

## 依赖

```go
require (
    github.com/hashicorp/vault/api v1.16.0
    github.com/hashicorp/vault/api/auth/approle v0.11.0
    github.com/hashicorp/vault/api/auth/kubernetes v0.11.0
)
```

## 参考文档

- [Vault Go API 文档](https://pkg.go.dev/github.com/hashicorp/vault/api)
- [Vault 官方文档](https://www.vaultproject.io/docs)
- [KV Secrets Engine](https://www.vaultproject.io/docs/secrets/kv)
- [Transit Secrets Engine](https://www.vaultproject.io/docs/secrets/transit)
- [Database Secrets Engine](https://www.vaultproject.io/docs/secrets/databases)
