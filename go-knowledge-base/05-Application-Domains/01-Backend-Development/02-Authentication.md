# Authentication Patterns

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #authentication #jwt #oauth #session #security

---

## 1. Authentication Fundamentals

### 1.1 Authentication vs Authorization

**Authentication**: Verifying who you are

- Passwords, biometrics, tokens
- Establishes identity

**Authorization**: Determining what you can do

- Permissions, roles, policies
- Granted after authentication

### 1.2 Authentication Factors

| Factor | Examples |
|--------|----------|
| Something you know | Password, PIN |
| Something you have | Phone, security key |
| Something you are | Fingerprint, face |
| Somewhere you are | Location, IP |

---

## 2. Session-Based Authentication

### 2.1 Session Management

```go
package auth

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "sync"
    "time"
)

type SessionManager struct {
    store  SessionStore
    expiry time.Duration
}

type Session struct {
    ID        string
    UserID    string
    Data      map[string]interface{}
    CreatedAt time.Time
    ExpiresAt time.Time
}

func (m *SessionManager) Create(ctx context.Context, userID string) (*Session, error) {
    id, err := generateSessionID()
    if err != nil {
        return nil, err
    }

    session := &Session{
        ID:        id,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(m.expiry),
    }

    if err := m.store.Save(ctx, session); err != nil {
        return nil, err
    }

    return session, nil
}

func (m *SessionManager) Get(ctx context.Context, sessionID string) (*Session, error) {
    session, err := m.store.Get(ctx, sessionID)
    if err != nil {
        return nil, err
    }

    if time.Now().After(session.ExpiresAt) {
        m.store.Delete(ctx, sessionID)
        return nil, ErrSessionExpired
    }

    return session, nil
}

func generateSessionID() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}
```

### 2.2 Redis Session Store

```go
package auth

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
    client *redis.Client
    prefix string
}

func (s *RedisSessionStore) Save(ctx context.Context, session *Session) error {
    key := s.prefix + session.ID
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }

    ttl := time.Until(session.ExpiresAt)
    return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *RedisSessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
    key := s.prefix + sessionID
    data, err := s.client.Get(ctx, key).Bytes()
    if err != nil {
        return nil, err
    }

    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }

    return &session, nil
}

func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
    return s.client.Del(ctx, s.prefix+sessionID).Err()
}
```

---

## 3. JWT Authentication

### 3.1 JWT Implementation

```go
package auth

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token expired")
)

type JWTManager struct {
    secretKey []byte
    issuer    string
    duration  time.Duration
}

type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID, username string, roles []string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Roles:    roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.duration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    m.issuer,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secretKey)
}

func (m *JWTManager) Validate(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return m.secretKey, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, ErrInvalidToken
}
```

### 3.2 Refresh Token Pattern

```go
package auth

import (
    "crypto/rand"
    "encoding/base64"
    "time"
)

type RefreshTokenManager struct {
    store  RefreshStore
    expiry time.Duration
}

type RefreshToken struct {
    Token     string
    UserID    string
    CreatedAt time.Time
    ExpiresAt time.Time
}

func (m *RefreshTokenManager) Generate(userID string) (*RefreshToken, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return nil, err
    }

    rt := &RefreshToken{
        Token:     base64.URLEncoding.EncodeToString(b),
        UserID:    userID,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(m.expiry),
    }

    if err := m.store.Save(rt); err != nil {
        return nil, err
    }

    return rt, nil
}

func (m *RefreshTokenManager) Refresh(tokenString string) (*TokenPair, error) {
    rt, err := m.store.Get(tokenString)
    if err != nil {
        return nil, err
    }

    if time.Now().After(rt.ExpiresAt) {
        m.store.Delete(tokenString)
        return nil, ErrTokenExpired
    }

    // Generate new access token
    accessToken, err := m.jwtManager.Generate(rt.UserID, "", nil)
    if err != nil {
        return nil, err
    }

    // Rotate refresh token
    newRT, err := m.Generate(rt.UserID)
    if err != nil {
        return nil, err
    }

    // Revoke old refresh token
    m.store.Delete(tokenString)

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: newRT.Token,
    }, nil
}
```

---

## 4. OAuth 2.0 / OpenID Connect

### 4.1 OAuth2 Flows

| Flow | Use Case | User Present? |
|------|----------|---------------|
| Authorization Code | Web apps, mobile apps | Yes |
| Implicit (deprecated) | SPAs | Yes |
| Password | Trusted apps | Yes |
| Client Credentials | Service-to-service | No |
| Device Code | Input-constrained devices | Yes |
| PKCE | Mobile apps | Yes |

### 4.2 OAuth2 Client

```go
package auth

import (
    "context"
    "golang.org/x/oauth2"
)

type OAuth2Client struct {
    config *oauth2.Config
}

type Token struct {
    AccessToken  string
    RefreshToken string
    TokenType    string
    Expiry       time.Time
}

func (c *OAuth2Client) GetAuthURL(state string) string {
    return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (c *OAuth2Client) Exchange(ctx context.Context, code string) (*Token, error) {
    oauth2Token, err := c.config.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }

    return &Token{
        AccessToken:  oauth2Token.AccessToken,
        RefreshToken: oauth2Token.RefreshToken,
        TokenType:    oauth2Token.TokenType,
        Expiry:       oauth2Token.Expiry,
    }, nil
}

func (c *OAuth2Client) Refresh(ctx context.Context, refreshToken string) (*Token, error) {
    token := &oauth2.Token{
        RefreshToken: refreshToken,
    }

    ts := c.config.TokenSource(ctx, token)
    newToken, err := ts.Token()
    if err != nil {
        return nil, err
    }

    return &Token{
        AccessToken:  newToken.AccessToken,
        RefreshToken: newToken.RefreshToken,
        TokenType:    newToken.TokenType,
        Expiry:       newToken.Expiry,
    }, nil
}
```

### 4.3 OIDC ID Token Verification

```go
package auth

import (
    "context"
    ""github.com/coreos/go-oidc/v3/oidc"
)

type OIDCVerifier struct {
    verifier *oidc.IDTokenVerifier
}

func NewOIDCVerifier(ctx context.Context, issuer, clientID string) (*OIDCVerifier, error) {
    provider, err := oidc.NewProvider(ctx, issuer)
    if err != nil {
        return nil, err
    }

    config := &oidc.Config{
        ClientID: clientID,
    }

    return &OIDCVerifier{
        verifier: provider.Verifier(config),
    }, nil
}

func (v *OIDCVerifier) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
    return v.verifier.Verify(ctx, rawIDToken)
}
```

---

## 5. API Key Authentication

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
    store APIKeyStore
    cache map[string]*APIKey
    mu    sync.RWMutex
}

type APIKey struct {
    ID          string
    KeyHash     string
    Name        string
    OwnerID     string
    Scopes      []string
    Permissions []string
    CreatedAt   time.Time
    ExpiresAt   *time.Time
    LastUsedAt  *time.Time
    Enabled     bool
}

func (m *APIKeyManager) Validate(ctx context.Context, key string) (*APIKey, error) {
    // Check cache first
    hash := hashKey(key)

    m.mu.RLock()
    if cached := m.cache[hash]; cached != nil {
        m.mu.RUnlock()
        if !cached.Enabled {
            return nil, ErrKeyDisabled
        }
        if cached.ExpiresAt != nil && time.Now().After(*cached.ExpiresAt) {
            return nil, ErrKeyExpired
        }
        return cached, nil
    }
    m.mu.RUnlock()

    // Check store
    apiKey, err := m.store.GetByHash(ctx, hash)
    if err != nil {
        return nil, err
    }

    // Update cache
    m.mu.Lock()
    m.cache[hash] = apiKey
    m.mu.Unlock()

    // Update last used
    now := time.Now()
    apiKey.LastUsedAt = &now
    go m.store.UpdateLastUsed(ctx, apiKey.ID, now)

    return apiKey, nil
}

func hashKey(key string) string {
    h := sha256.Sum256([]byte(key))
    return hex.EncodeToString(h[:])
}
```

---

## 6. Multi-Factor Authentication

### 6.1 TOTP Implementation

```go
package auth

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base32"
    "encoding/binary"
    "fmt"
    "time"
)

type TOTP struct {
    secret []byte
    digits int
    period int
}

func NewTOTP(secret string) *TOTP {
    key, _ := base32.StdEncoding.DecodeString(secret)
    return &TOTP{
        secret: key,
        digits: 6,
        period: 30,
    }
}

func (t *TOTP) Generate() string {
    counter := uint64(time.Now().Unix() / int64(t.period))
    return t.generateForCounter(counter)
}

func (t *TOTP) Verify(code string) bool {
    counter := uint64(time.Now().Unix() / int64(t.period))

    // Allow 1 period before and after
    for i := -1; i <= 1; i++ {
        if t.generateForCounter(counter+uint64(i)) == code {
            return true
        }
    }
    return false
}

func (t *TOTP) generateForCounter(counter uint64) string {
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, counter)

    mac := hmac.New(sha1.New, t.secret)
    mac.Write(buf)
    hash := mac.Sum(nil)

    offset := hash[len(hash)-1] & 0x0F
    code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

    return fmt.Sprintf(fmt.Sprintf("%%0%dd", t.digits), code%uint32(pow10(t.digits)))
}
```

---

## 7. Password Security

### 7.1 Password Hashing

```go
package auth

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/argon2"
)

var (
    ErrPasswordMismatch = errors.New("password mismatch")
    ErrPasswordTooWeak  = errors.New("password too weak")
)

// Bcrypt hashing
type BcryptHasher struct {
    cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
    return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    return string(bytes), err
}

func (h *BcryptHasher) Verify(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Argon2id hashing (recommended)
type Argon2Hasher struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
    saltLen uint32
}

func NewArgon2Hasher() *Argon2Hasher {
    return &Argon2Hasher{
        time:    3,
        memory:  64 * 1024,
        threads: 4,
        keyLen:  32,
        saltLen: 16,
    }
}

func (h *Argon2Hasher) Hash(password string, salt []byte) []byte {
    return argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)
}
```

### 7.2 Password Validation

```go
package auth

import (
    "unicode"
    "golang.org/x/crypto/ssh/terminal"
)

type PasswordPolicy struct {
    MinLength      int
    MaxLength      int
    RequireUpper   bool
    RequireLower   bool
    RequireDigit   bool
    RequireSpecial bool
}

func (p *PasswordPolicy) Validate(password string) error {
    if len(password) < p.MinLength {
        return ErrPasswordTooShort
    }
    if p.MaxLength > 0 && len(password) > p.MaxLength {
        return ErrPasswordTooLong
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

    if p.RequireUpper && !hasUpper {
        return ErrPasswordNoUpper
    }
    if p.RequireLower && !hasLower {
        return ErrPasswordNoLower
    }
    if p.RequireDigit && !hasDigit {
        return ErrPasswordNoDigit
    }
    if p.RequireSpecial && !hasSpecial {
        return ErrPasswordNoSpecial
    }

    return nil
}
```

---

## 8. Authentication Middleware

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

func Authentication(authService AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractBearerToken(c)
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "UNAUTHORIZED",
                Message: "Missing authentication token",
            })
            return
        }

        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "UNAUTHORIZED",
                Message: "Invalid or expired token",
            })
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("claims", claims)
        c.Next()
    }
}

func extractBearerToken(c *gin.Context) string {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        return ""
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
        return ""
    }

    return parts[1]
}
```

---

## 9. Best Practices

- [ ] Use HTTPS for all authentication
- [ ] Hash passwords with Argon2id or bcrypt
- [ ] Implement rate limiting on auth endpoints
- [ ] Use short-lived access tokens
- [ ] Implement refresh token rotation
- [ ] Support MFA for sensitive operations
- [ ] Log authentication events
- [ ] Implement account lockout
- [ ] Use secure session management
- [ ] Validate all tokens server-side

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02
