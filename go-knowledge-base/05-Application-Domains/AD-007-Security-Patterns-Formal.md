# AD-007: Security Architecture Patterns

> **Dimension**: Application Domains
> **Level**: S (16+ KB)
> **Tags**: #security #authentication #authorization #encryption #oauth2

---

## 1. Security Architecture Fundamentals

### 1.1 CIA Triad

**Confidentiality**: Prevent unauthorized access to data
**Integrity**: Ensure data is not modified by unauthorized parties
**Availability**: Systems are accessible when needed

### 1.2 Defense in Depth

Multiple layers of security controls:

- Network layer: Firewalls, VPNs
- Application layer: Authentication, authorization
- Data layer: Encryption, access controls

---

## 2. Authentication Patterns

### 2.1 JWT Authentication

```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, username string, roles []string, secret []byte) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Roles:    roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}

func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrSignatureInvalid
}
```

### 2.2 OAuth2 / OIDC Implementation

```go
package oauth

import (
    "context"
    "golang.org/x/oauth2"
)

type OAuth2Config struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    AuthURL      string
    TokenURL     string
    Scopes       []string
}

func (c *OAuth2Config) Exchange(ctx context.Context, code string) (*Token, error) {
    config := &oauth2.Config{
        ClientID:     c.ClientID,
        ClientSecret: c.ClientSecret,
        RedirectURL:  c.RedirectURL,
        Endpoint: oauth2.Endpoint{
            AuthURL:  c.AuthURL,
            TokenURL: c.TokenURL,
        },
        Scopes: c.Scopes,
    }

    token, err := config.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }

    return &Token{
        AccessToken:  token.AccessToken,
        RefreshToken: token.RefreshToken,
        Expiry:       token.Expiry,
    }, nil
}
```

### 2.3 API Key Authentication

```go
package auth

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "sync"
    "time"
)

type APIKeyManager struct {
    keys map[string]*APIKey
    mu   sync.RWMutex
}

type APIKey struct {
    ID        string
    KeyHash   string
    OwnerID   string
    Scopes    []string
    CreatedAt time.Time
    ExpiresAt time.Time
    LastUsed  time.Time
}

func (m *APIKeyManager) ValidateKey(ctx context.Context, key string) (*APIKey, error) {
    hash := hashKey(key)

    m.mu.RLock()
    defer m.mu.RUnlock()

    for _, apiKey := range m.keys {
        if apiKey.KeyHash == hash {
            if time.Now().After(apiKey.ExpiresAt) {
                return nil, ErrKeyExpired
            }
            return apiKey, nil
        }
    }

    return nil, ErrInvalidKey
}

func hashKey(key string) string {
    h := sha256.Sum256([]byte(key))
    return hex.EncodeToString(h[:])
}
```

---

## 3. Authorization Patterns

### 3.1 RBAC (Role-Based Access Control)

```go
package auth

type RBAC struct {
    roles       map[string]*Role
    permissions map[string]*Permission
    assignments map[string][]string // userID -> roleIDs
}

type Role struct {
    ID          string
    Name        string
    Permissions []string
}

type Permission struct {
    ID     string
    Resource string
    Action   string
}

func (r *RBAC) HasPermission(userID, resource, action string) bool {
    roles := r.assignments[userID]

    for _, roleID := range roles {
        role := r.roles[roleID]
        for _, permID := range role.Permissions {
            perm := r.permissions[permID]
            if perm.Resource == resource && perm.Action == action {
                return true
            }
        }
    }

    return false
}
```

### 3.2 ABAC (Attribute-Based Access Control)

```go
package auth

type ABAC struct {
    policies []Policy
}

type Policy struct {
    ID        string
    Subject   Condition
    Resource  Condition
    Action    Condition
    Environment Condition
    Effect    Effect
}

type Condition struct {
    Attributes map[string]interface{}
    Operator   string // equals, greater_than, in, etc.
}

type Effect int

const (
    Deny Effect = iota
    Allow
)

func (a *ABAC) Evaluate(subject, resource, action map[string]interface{}) Effect {
    for _, policy := range a.policies {
        if matches(policy.Subject, subject) &&
           matches(policy.Resource, resource) &&
           matches(policy.Action, action) {
            return policy.Effect
        }
    }
    return Deny
}
```

---

## 4. Data Protection

### 4.1 Encryption at Rest

```go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

type Encrypter struct {
    key []byte
}

func NewEncrypter(key []byte) *Encrypter {
    return &Encrypter{key: key}
}

func (e *Encrypter) Encrypt(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *Encrypter) Decrypt(ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, ErrCipherTooShort
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
```

### 4.2 TLS Configuration

```go
package tls

import (
    "crypto/tls"
    "time"
)

func NewConfig(certFile, keyFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
        SessionTicketsDisabled:   false,
        SessionTicketKey:         [32]byte{},
        ClientSessionCache:       tls.NewLRUClientSessionCache(128),
    }, nil
}
```

---

## 5. Secure Communication

### 5.1 mTLS (Mutual TLS)

```go
package mtls

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

func NewMutualTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }

    caCert, err := ioutil.ReadFile(caFile)
    if err != nil {
        return nil, err
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert,
        MinVersion:   tls.VersionTLS12,
    }, nil
}
```

---

## 6. Security Headers

```go
package security

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

        next.ServeHTTP(w, r)
    })
}
```

---

## 7. Secrets Management

```go
package secrets

import (
    "context"
    "encoding/json"
    "os"
)

type Manager interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key, value string) error
}

// Environment-based implementation (development only)
type EnvManager struct{}

func (e *EnvManager) Get(ctx context.Context, key string) (string, error) {
    value := os.Getenv(key)
    if value == "" {
        return "", ErrSecretNotFound
    }
    return value, nil
}

// HashiCorp Vault implementation
type VaultManager struct {
    client *vault.Client
}

func (v *VaultManager) Get(ctx context.Context, key string) (string, error) {
    secret, err := v.client.KVv2("secret").Get(ctx, key)
    if err != nil {
        return "", err
    }

    value, ok := secret.Data["value"].(string)
    if !ok {
        return "", ErrInvalidSecret
    }

    return value, nil
}
```

---

## 8. Security Best Practices

### 8.1 Input Validation

```go
package validation

import (
    "regexp"
    "unicode"
)

func ValidateUsername(username string) error {
    if len(username) < 3 || len(username) > 32 {
        return ErrInvalidLength
    }

    valid := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !valid.MatchString(username) {
        return ErrInvalidCharacters
    }

    return nil
}

func ValidatePassword(password string) error {
    if len(password) < 12 {
        return ErrPasswordTooShort
    }

    var hasUpper, hasLower, hasDigit, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
        return ErrPasswordComplexity
    }

    return nil
}
```

### 8.2 Rate Limiting

```go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

type RateLimiter struct {
    requests map[string][]time.Time
    limit    int
    window   time.Duration
    mu       sync.RWMutex
}

func (r *RateLimiter) Allow(ctx context.Context, key string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-r.window)

    // Filter requests within window
    var valid []time.Time
    for _, t := range r.requests[key] {
        if t.After(windowStart) {
            valid = append(valid, t)
        }
    }

    if len(valid) >= r.limit {
        return false
    }

    valid = append(valid, now)
    r.requests[key] = valid

    return true
}
```

---

## 9. Security Checklist

- [ ] Use HTTPS everywhere
- [ ] Implement proper authentication
- [ ] Apply least privilege principle
- [ ] Encrypt sensitive data at rest
- [ ] Use prepared statements (prevent SQL injection)
- [ ] Validate all inputs
- [ ] Implement rate limiting
- [ ] Set security headers
- [ ] Use secure password storage (bcrypt/Argon2)
- [ ] Implement audit logging
- [ ] Regular dependency updates
- [ ] Security testing (SAST/DAST)

---

## References

1. OWASP Top 10
2. NIST Cybersecurity Framework
3. OAuth 2.0 and OpenID Connect
4. Go Secure Coding Practices

---

**Quality Rating**: S (16+ KB)
**Last Updated**: 2026-04-02
