# Secrets Management

> **分类**: 工程与云原生
> **标签**: #secrets #vault #security #encryption #rotation
> **参考**: HashiCorp Vault, AWS Secrets Manager, Kubernetes Secrets

---

## 1. Formal Definition

### 1.1 What is Secrets Management?

Secrets management is the practice of securely storing, accessing, and distributing sensitive information such as passwords, API keys, certificates, and tokens. It encompasses the entire lifecycle of secrets including creation, rotation, revocation, and auditing.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Secrets Management Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     SECRETS MANAGEMENT SYSTEM                        │   │
│   │  (Vault, AWS Secrets Manager, Azure Key Vault, etc.)                 │   │
│   │                                                                       │   │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│   │  │   Storage    │  │   Access     │  │   Lifecycle  │              │   │
│   │  │   Engine     │  │   Control    │  │   Management │              │   │
│   │  │              │  │              │  │              │              │   │
│   │  │ • Encryption │  │ • Auth       │  │ • Rotation   │              │   │
│   │  │ • HSM        │  │ • Policies   │  │ • Revocation │              │   │
│   │  │ • Backup     │  │ • Audit      │  │ • Leasing    │              │   │
│   │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│   │                                                                       │   │
│   │  Secret Types:                                                        │   │
│   │  • Static secrets (DB passwords, API keys)                            │   │
│   │  • Dynamic secrets (on-demand credentials)                            │   │
│   │  • PKI certificates                                                   │   │
│   │  • Encryption keys                                                    │   │
│   │                                                                       │   │
│   └────────────────────────┬────────────────────────────────────────────┘   │
│                            │                                                │
│         ┌──────────────────┼──────────────────┐                             │
│         │                  │                  │                             │
│         ▼                  ▼                  ▼                             │
│   ┌──────────┐       ┌──────────┐       ┌──────────┐                       │
│   │  Humans  │       │  Apps    │       │  CI/CD   │                       │
│   │  (CLI)   │       │  (SDK)   │       │  (Token) │                       │
│   └──────────┘       └──────────┘       └──────────┘                       │
│                                                                             │
│   SECRET LIFECYCLE:                                                         │
│   Create → Store → Access → Rotate → Revoke → Audit                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Secret Types and Handling

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Secret Types Classification                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  TYPE              │  EXAMPLES              │  ROTATION │  SCOPE            │
├────────────────────┼────────────────────────┼───────────┼───────────────────│
│  Credentials       │ DB passwords           │ 90 days   │ Service-specific  │
│                    │ Service accounts       │           │                   │
│  ──────────────────┼────────────────────────┼───────────┼───────────────────│
│  API Keys          │ Third-party APIs       │ On breach │ Per environment   │
│                    │ Internal service keys  │           │                   │
│  ──────────────────┼────────────────────────┼───────────┼───────────────────│
│  Certificates      │ TLS/SSL certs          │ Pre-expiry│ Domain-based      │
│                    │ Client certs           │           │                   │
│  ──────────────────┼────────────────────────┼───────────┼───────────────────│
│  Tokens            │ JWT signing keys       │ 90 days   │ Service-wide      │
│                    │ OAuth tokens           │ On expiry │                   │
│  ──────────────────┼────────────────────────┼───────────┼───────────────────│
│  Encryption Keys   │ AES keys               │ 1 year    │ Data-class based  │
│                    │ RSA key pairs          │           │                   │
│  ──────────────────┼────────────────────────┼───────────┼───────────────────│
│  Connection Strings│ DB URLs                │ 90 days   │ Per service       │
│                    │ Message queue URLs     │           │                   │
│                                                                             │
│  NEVER STORE IN CODE:                                                       │
│  ✗ Hardcoded passwords                                                      │
│  ✗ API keys in repositories                                                 │
│  ✗ Private keys in configuration files                                      │
│  ✗ Database credentials in environment variables (unencrypted)              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 HashiCorp Vault Integration

```go
package vault

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/vault/api"
)

// Client wraps Vault API client
type Client struct {
    client *api.Client
    config *Config
}

// Config holds Vault configuration
type Config struct {
    Address       string
    AuthMethod    string // token, kubernetes, approle, aws
    Role          string
    Token         string
    TLSConfig     *api.TLSConfig
}

// NewClient creates a new Vault client
func NewClient(cfg *Config) (*Client, error) {
    config := api.DefaultConfig()
    config.Address = cfg.Address

    if cfg.TLSConfig != nil {
        config.ConfigureTLS(cfg.TLSConfig)
    }

    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create vault client: %w", err)
    }

    vaultClient := &Client{
        client: client,
        config: cfg,
    }

    // Authenticate
    if err := vaultClient.authenticate(); err != nil {
        return nil, err
    }

    return vaultClient, nil
}

// authenticate performs authentication based on configured method
func (c *Client) authenticate() error {
    switch c.config.AuthMethod {
    case "token":
        c.client.SetToken(c.config.Token)
    case "kubernetes":
        return c.authKubernetes()
    case "approle":
        return c.authAppRole()
    default:
        return fmt.Errorf("unsupported auth method: %s", c.config.AuthMethod)
    }
    return nil
}

// authKubernetes performs Kubernetes authentication
func (c *Client) authKubernetes() error {
    // Read service account token
    jwt, err := readFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
    if err != nil {
        return fmt.Errorf("failed to read service account token: %w", err)
    }

    // Login to Vault
    path := "auth/kubernetes/login"
    secret, err := c.client.Logical().Write(path, map[string]interface{}{
        "jwt":  jwt,
        "role": c.config.Role,
    })
    if err != nil {
        return fmt.Errorf("failed to login with kubernetes auth: %w", err)
    }

    c.client.SetToken(secret.Auth.ClientToken)
    return nil
}

// authAppRole performs AppRole authentication
func (c *Client) authAppRole() error {
    // Implementation for AppRole auth
    return nil
}

// GetSecret retrieves a static secret from Vault
func (c *Client) GetSecret(ctx context.Context, path string) (map[string]interface{}, error) {
    secret, err := c.client.Logical().ReadWithContext(ctx, path)
    if err != nil {
        return nil, fmt.Errorf("failed to read secret: %w", err)
    }

    if secret == nil {
        return nil, fmt.Errorf("secret not found at path: %s", path)
    }

    return secret.Data, nil
}

// GetDynamicCredentials retrieves dynamic database credentials
func (c *Client) GetDynamicCredentials(ctx context.Context, role string) (*DynamicCredentials, error) {
    path := fmt.Sprintf("database/creds/%s", role)

    secret, err := c.client.Logical().ReadWithContext(ctx, path)
    if err != nil {
        return nil, fmt.Errorf("failed to get dynamic credentials: %w", err)
    }

    if secret == nil {
        return nil, fmt.Errorf("dynamic credentials not available for role: %s", role)
    }

    creds := &DynamicCredentials{
        Username:    secret.Data["username"].(string),
        Password:    secret.Data["password"].(string),
        LeaseID:     secret.LeaseID,
        LeaseDuration: time.Duration(secret.LeaseDuration) * time.Second,
        Renewable:   secret.Renewable,
    }

    return creds, nil
}

// DynamicCredentials holds dynamic database credentials
type DynamicCredentials struct {
    Username      string
    Password      string
    LeaseID       string
    LeaseDuration time.Duration
    Renewable     bool
}

// RenewLease renews a secret lease
func (c *Client) RenewLease(ctx context.Context, leaseID string, increment int) error {
    _, err := c.client.Sys().RenewWithContext(ctx, leaseID, increment)
    return err
}

// RevokeLease revokes a secret lease
func (c *Client) RevokeLease(ctx context.Context, leaseID string) error {
    return c.client.Sys().RevokeWithContext(ctx, leaseID)
}

// TransitEncrypt encrypts data using Vault's transit engine
func (c *Client) TransitEncrypt(ctx context.Context, keyName string, plaintext []byte) (string, error) {
    path := fmt.Sprintf("transit/encrypt/%s", keyName)

    // Base64 encode plaintext
    encoded := base64.StdEncoding.EncodeToString(plaintext)

    secret, err := c.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
        "plaintext": encoded,
    })
    if err != nil {
        return "", fmt.Errorf("failed to encrypt: %w", err)
    }

    return secret.Data["ciphertext"].(string), nil
}

// TransitDecrypt decrypts data using Vault's transit engine
func (c *Client) TransitDecrypt(ctx context.Context, keyName, ciphertext string) ([]byte, error) {
    path := fmt.Sprintf("transit/decrypt/%s", keyName)

    secret, err := c.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
        "ciphertext": ciphertext,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }

    // Base64 decode
    encoded := secret.Data["plaintext"].(string)
    return base64.StdEncoding.DecodeString(encoded)
}

// CreatePolicy creates a Vault policy
func (c *Client) CreatePolicy(name, rules string) error {
    return c.client.Sys().PutPolicy(name, rules)
}

// WithRenewal wraps a function with automatic token renewal
func (c *Client) WithRenewal(ctx context.Context, fn func() error) error {
    // Start token renewal in background
    renewer, err := c.client.NewRenewer(&api.RenewerInput{
        Grace: 5 * time.Second,
    })
    if err != nil {
        return err
    }

    go renewer.Renew()
    defer renewer.Stop()

    return fn()
}

func readFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

### 2.2 Kubernetes External Secrets

```yaml
# External Secrets Operator configuration
apiVersion: external-secrets.io/v1beta1
kind: ClusterSecretStore
metadata:
  name: vault-backend
spec:
  provider:
    vault:
      server: "https://vault.example.com"
      path: "secret"
      version: "v2"
      auth:
        kubernetes:
          mountPath: "kubernetes"
          role: "external-secrets"
          serviceAccountRef:
            name: external-secrets
            namespace: external-secrets
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: database-credentials
  namespace: production
spec:
  refreshInterval: "1h"
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-backend
  target:
    name: database-credentials
    creationPolicy: Owner
    template:
      type: Opaque
      data:
        connection-string: "postgresql://{{ .username }}:{{ .password }}@db:5432/app"
  data:
    - secretKey: username
      remoteRef:
        key: secret/data/production/database
        property: username
    - secretKey: password
      remoteRef:
        key: secret/data/production/database
        property: password
```

---

## 3. Production-Ready Configurations

### 3.1 Vault HA Configuration

```hcl
# vault.hcl
# High Availability Vault configuration

storage "raft" {
  path = "/vault/data"
  node_id = "NODE_ID"

  retry_leader_election = true

  autopilot {
    cleanup_dead_servers = true
    last_contact_threshold = "10s"
    max_trailing_logs = 250
    min_quorum = 3
    server_stabilization_time = "10s"
  }
}

listener "tcp" {
  address = "0.0.0.0:8200"
  cluster_address = "0.0.0.0:8201"
  tls_cert_file = "/vault/certs/tls.crt"
  tls_key_file = "/vault/certs/tls.key"
  tls_min_version = "tls12"
}

seal "awskms" {
  region = "us-east-1"
  kms_key_id = "arn:aws:kms:us-east-1:123456789:key/VAULT_KEY_ID"
}

disable_mlock = true

api_addr = "https://VAULT_NODE_ADDRESS:8200"
cluster_addr = "https://VAULT_NODE_ADDRESS:8201"

ui = true

# Audit logging
audit_file {
  path = "/vault/audit/audit.log"
}

# Telemetry
telemetry {
  prometheus_retention_time = "30s"
  disable_hostname = true
}
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   Secrets Management Security                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ENCRYPTION                                                                 │
│  ✓ Encryption at rest (AES-256-GCM)                                         │
│  ✓ Encryption in transit (TLS 1.2+)                                         │
│  ✓ HSM integration for master key                                           │
│  ✓ Auto-unseal with cloud KMS                                               │
│                                                                             │
│  ACCESS CONTROL                                                             │
│  ✓ Principle of least privilege                                             │
│  ✓ RBAC with fine-grained policies                                          │
│  ✓ MFA for administrative access                                            │
│  ✓ Short-lived tokens with automatic expiration                             │
│                                                                             │
│  AUDITING                                                                   │
│  ✓ Comprehensive audit logging                                              │
│  ✓ SIEM integration                                                         │
│  ✓ Access pattern analysis                                                  │
│  ✓ Anomaly detection                                                        │
│                                                                             │
│  ROTATION                                                                   │
│  ✓ Automated credential rotation                                            │
│  ✓ Dynamic secrets with TTL                                                 │
│  ✓ Certificate auto-renewal                                                 │
│  ✓ Break-glass procedures                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 Secrets Manager Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Secrets Manager Comparison Matrix                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Feature         │  Vault    │  AWS SM    │  Azure KV  │  GCP SM    │  K8s │
├──────────────────┼───────────┼────────────┼────────────┼────────────┼──────│
│  Dynamic Secrets │  ★★★★★    │  ★★★☆☆     │  ★★★☆☆     │  ★★★☆☆     │  ★☆☆ │
│  PKI             │  ★★★★★    │  ★★☆☆☆     │  ★★★★☆     │  ★★★☆☆     │  ★☆☆ │
│  Encryption      │  ★★★★★    │  ★★★★☆     │  ★★★★★     │  ★★★★☆     │  ★★☆ │
│  Multi-cloud     │  ★★★★★    │  ★★☆☆☆     │  ★★☆☆☆     │  ★★☆☆☆     │  ★★☆ │
│  K8s Integration │  ★★★★★    │  ★★★★☆     │  ★★★★☆     │  ★★★★☆     │  ★★★ │
│  Cost            │  $$$      │  $         │  $         │  $         │  Free│
│                                                                             │
│  Recommendation:                                                            │
│  • Multi-cloud: HashiCorp Vault                                             │
│  • AWS-only: AWS Secrets Manager + Parameter Store                          │
│  • Azure-only: Azure Key Vault                                              │
│  • GCP-only: Google Secret Manager                                          │
│  • K8s-native: External Secrets Operator + cloud provider                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Secrets Management Best Practices                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  STORAGE                                                                    │
│  ✓ Centralized secrets management                                           │
│  ✓ Encryption at rest and in transit                                        │
│  ✓ No secrets in code or repositories                                       │
│  ✓ HSM for master key protection                                            │
│                                                                             │
│  ACCESS                                                                     │
│  ✓ Dynamic secrets over static                                              │
│  ✓ Short-lived tokens                                                       │
│  ✓ Service-specific credentials                                             │
│  ✓ Just-in-time access                                                      │
│                                                                             │
│  ROTATION                                                                   │
│  ✓ Automated rotation schedules                                             │
│  ✓ Rotation on compromise                                                   │
│  ✓ Graceful rotation (dual-credential period)                               │
│  ✓ Certificate lifecycle management                                         │
│                                                                             │
│  MONITORING                                                                 │
│  ✓ Audit all secret access                                                  │
│  ✓ Alert on unusual patterns                                                │
│  ✓ Track secret usage                                                       │
│  ✓ Regular access reviews                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. HashiCorp Vault Documentation
2. AWS Secrets Manager Best Practices
3. Azure Key Vault Security Guide
4. NIST SP 800-57 - Key Management
5. OWASP Secrets Management Cheat Sheet
