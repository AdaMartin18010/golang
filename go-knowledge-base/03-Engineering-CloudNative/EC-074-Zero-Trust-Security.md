# Zero Trust Security

> **分类**: 工程与云原生
> **标签**: #zerotrust #security #identity #microsegmentation #mfa
> **参考**: NIST SP 800-207, Google BeyondCorp, Microsoft Zero Trust

---

## 1. Formal Definition

### 1.1 What is Zero Trust?

Zero Trust is a security framework requiring all users, whether in or outside the organization's network, to be authenticated, authorized, and continuously validated before being granted access to applications and data. Zero Trust assumes there is no traditional network edge; networks can be local, cloud, or hybrid.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Zero Trust Architecture                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   TRADITIONAL SECURITY (Perimeter-Based)                                    │
│   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━                                  │
│                                                                             │
│       Internet        Corporate Network        Data Center                  │
│          │                  │                      │                        │
│          │    ┌─────────────┴─────────────┐       │                        │
│          └───►│        Firewall           │◄──────┘                        │
│               │  (Trusted inside network) │                                 │
│               └─────────────┬─────────────┘                                 │
│                             │                                               │
│                          [TRUSTED]                                          │
│                                                                             │
│   Assumption: Inside network = Safe | Outside = Dangerous                   │
│   Problem: Insider threats, lateral movement, breached perimeter            │
│                                                                             │
│   ────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│   ZERO TRUST SECURITY (Identity & Context-Based)                            │
│   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━                            │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        POLICY DECISION POINT                         │   │
│   │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐          │   │
│   │  │   Identity   │    │   Device     │    │   Context    │          │   │
│   │  │   Engine     │◄──►│   Trust      │◄──►│   Engine     │          │   │
│   │  │              │    │   Engine     │    │              │          │   │
│   │  │ • Who?       │    │ • Device     │    │ • Where?     │          │   │
│   │  │ • MFA?       │    │   health     │    │ • When?      │          │   │
│   │  │ • Privileges │    │ • Compliance │    │ • Behavior?  │          │   │
│   │  └──────────────┘    └──────────────┘    └──────────────┘          │   │
│   │                              │                                      │   │
│   │                              ▼                                      │   │
│   │                    ┌──────────────────┐                             │   │
│   │                    │  Policy Engine   │                             │   │
│   │                    │  (Authorize?)    │                             │   │
│   │                    └────────┬─────────┘                             │   │
│   └─────────────────────────────┼───────────────────────────────────────┘   │
│                                 │                                           │
│                                 │ Access Decision                           │
│                                 ▼                                           │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      POLICY ENFORCEMENT POINT                        │   │
│   │                                                                       │   │
│   │  User ──► Identity Verify ──► Device Check ──► Least Privilege ──►  │   │
│   │                               Health          Access                  │   │
│   │                                                                       │   │
│   │  [CONTINUOUS VALIDATION]                                              │   │
│   │  • Session monitoring                                                   │   │
│   │  • Behavioral analytics                                                 │   │
│   │  • Dynamic risk scoring                                                 │   │
│   │  • Just-in-time access                                                  │   │
│   │                                                                       │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   Core Principles:                                                          │
│   1. Never trust, always verify                                             │
│   2. Assume breach                                                          │
│   3. Verify explicitly                                                      │
│   4. Use least privilege access                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Zero Trust Pillars

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Zero Trust Security Pillars                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌─────────────┐                                     │
│                         │   IDENTITY  │                                     │
│                         │  (Who)      │                                     │
│                         │             │                                     │
│                         │ • SSO/MFA   │                                     │
│                         │ • Adaptive  │                                     │
│                         │   auth      │                                     │
│                         │ • Risk-based│                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│     ┌─────────────┐            │            ┌─────────────┐                │
│     │   DEVICE    │◄───────────┼───────────►│ APPLICATION │                │
│     │  (What)     │            │            │  (How)      │                │
│     │             │            │            │             │                │
│     │ • Health    │            │            │ • Micro-    │                │
│     │   check     │            │            │   services  │                │
│     │ • Compliance│            │            │ • API       │                │
│     │ • EDR/MDM   │            │            │   gateway   │                │
│     └─────────────┘            │            └─────────────┘                │
│                                │                                            │
│     ┌─────────────┐            │            ┌─────────────┐                │
│     │   NETWORK   │◄───────────┼───────────►│    DATA     │                │
│     │  (Where)    │            │            │  (What)     │                │
│     │             │            │            │             │                │
│     │ • Micro-    │            │            │ • Classification            │
│     │   segmentation         │            │ • Encryption │                │
│     │ • Encrypted │            │            │ • DLP       │                │
│     │   tunnels   │            │            │ • Access    │                │
│     └─────────────┘            │            │   controls  │                │
│                                │            └─────────────┘                │
│                         ┌──────┴──────┐                                     │
│                         │  AUTOMATION │                                     │
│                         │   & ANALYTICS                                    │
│                         │             │                                     │
│                         │ • SIEM/SOAR │                                     │
│                         │ • UEBA      │                                     │
│                         │ • Threat    │                                     │
│                         │   intel     │                                     │
│                         └─────────────┘                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 Service-to-Service Authentication

```go
package zerotrust

import (
    "context"
    "crypto/tls"
    "fmt"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
)

// ServiceIdentity represents a service's identity
type ServiceIdentity struct {
    ServiceID   string
    ServiceName string
    Namespace   string
    Cluster     string
}

// ServiceTokenManager manages service tokens
type ServiceTokenManager struct {
    privateKey []byte
    publicKey  []byte
    ttl        time.Duration
}

// ServiceClaims represents JWT claims for service identity
type ServiceClaims struct {
    ServiceID   string `json:"sid"`
    ServiceName string `json:"sname"`
    Namespace   string `json:"ns"`
    Cluster     string `json:"cluster"`
    jwt.RegisteredClaims
}

// NewServiceTokenManager creates a new token manager
func NewServiceTokenManager(privateKey, publicKey []byte, ttl time.Duration) *ServiceTokenManager {
    return &ServiceTokenManager{
        privateKey: privateKey,
        publicKey:  publicKey,
        ttl:        ttl,
    }
}

// GenerateToken generates a service token
func (m *ServiceTokenManager) GenerateToken(identity ServiceIdentity) (string, error) {
    claims := ServiceClaims{
        ServiceID:   identity.ServiceID,
        ServiceName: identity.ServiceName,
        Namespace:   identity.Namespace,
        Cluster:     identity.Cluster,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "service-identity",
            Subject:   identity.ServiceID,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(m.privateKey)
}

// ValidateToken validates a service token
func (m *ServiceTokenManager) ValidateToken(tokenString string) (*ServiceIdentity, error) {
    token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return m.publicKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    if claims, ok := token.Claims.(*ServiceClaims); ok && token.Valid {
        return &ServiceIdentity{
            ServiceID:   claims.ServiceID,
            ServiceName: claims.ServiceName,
            Namespace:   claims.Namespace,
            Cluster:     claims.Cluster,
        }, nil
    }

    return nil, fmt.Errorf("invalid token")
}

// ZeroTrustInterceptor creates a gRPC interceptor for zero trust
func ZeroTrustInterceptor(tokenManager *ServiceTokenManager) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Extract token from metadata
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, fmt.Errorf("missing metadata")
        }

        tokens := md.Get("authorization")
        if len(tokens) == 0 {
            return nil, fmt.Errorf("missing authorization token")
        }

        // Validate token
        identity, err := tokenManager.ValidateToken(tokens[0])
        if err != nil {
            return nil, fmt.Errorf("invalid token: %w", err)
        }

        // Add identity to context
        ctx = context.WithValue(ctx, "service-identity", identity)

        // Call handler
        return handler(ctx, req)
    }
}

// ZeroTrustHTTPMiddleware creates HTTP middleware for zero trust
func ZeroTrustHTTPMiddleware(tokenManager *ServiceTokenManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from header
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Error(w, "Missing authorization", http.StatusUnauthorized)
                return
            }

            // Remove "Bearer " prefix if present
            if len(token) > 7 && token[:7] == "Bearer " {
                token = token[7:]
            }

            // Validate token
            identity, err := tokenManager.ValidateToken(token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            // Add identity to context
            ctx := context.WithValue(r.Context(), "service-identity", identity)

            // Continue with request
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 2.2 mTLS Configuration

```go
package zerotrust

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "os"
)

// MTLSConfig holds mTLS configuration
type MTLSConfig struct {
    CertFile       string
    KeyFile        string
    CAFile         string
    ServerName     string
    RequireClientCert bool
}

// CreateServerTLSConfig creates TLS config for server
func CreateServerTLSConfig(config *MTLSConfig) (*tls.Config, error) {
    // Load server certificate
    cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
    if err != nil {
        return nil, fmt.Errorf("failed to load server cert: %w", err)
    }

    // Load CA certificate
    caCert, err := os.ReadFile(config.CAFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read CA cert: %w", err)
    }

    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        return nil, fmt.Errorf("failed to parse CA cert")
    }

    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
    }

    if config.RequireClientCert {
        tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
    }

    return tlsConfig, nil
}

// CreateClientTLSConfig creates TLS config for client
func CreateClientTLSConfig(config *MTLSConfig) (*tls.Config, error) {
    // Load client certificate
    cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
    if err != nil {
        return nil, fmt.Errorf("failed to load client cert: %w", err)
    }

    // Load CA certificate
    caCert, err := os.ReadFile(config.CAFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read CA cert: %w", err)
    }

    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        return nil, fmt.Errorf("failed to parse CA cert")
    }

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        ServerName:   config.ServerName,
    }, nil
}
```

---

## 3. Production-Ready Configurations

### 3.1 Istio Zero Trust Configuration

```yaml
# PeerAuthentication - Require mTLS
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT
---
# AuthorizationPolicy - Deny all by default
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
  namespace: production
spec:
  {}
---
# AuthorizationPolicy - Allow specific traffic
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: payment-service
  namespace: production
spec:
  selector:
    matchLabels:
      app: payment-service
  action: ALLOW
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/production/sa/order-service"]
    to:
    - operation:
        methods: ["POST", "GET"]
        paths: ["/api/v1/payments/*"]
    when:
    - key: request.auth.claims[scope]
      values: ["payments:write"]
  - from:
    - source:
        principals: ["cluster.local/ns/production/sa/admin-service"]
    to:
    - operation:
        methods: ["*"]
        paths: ["/api/v1/admin/*"]
---
# RequestAuthentication - JWT validation
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: jwt-auth
  namespace: production
spec:
  selector:
    matchLabels:
      app: api-gateway
  jwtRules:
  - issuer: "https://auth.example.com"
    jwksUri: "https://auth.example.com/.well-known/jwks.json"
    audiences:
    - "api.example.com"
    forwardOriginalToken: true
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Zero Trust Security Checklist                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  IDENTITY                                                                   │
│  ✓ Single Sign-On (SSO) implementation                                      │
│  ✓ Multi-Factor Authentication (MFA) required                               │
│  ✓ Risk-based adaptive authentication                                       │
│  ✓ Just-in-time privileged access                                           │
│  ✓ Regular access reviews                                                   │
│                                                                             │
│  DEVICE                                                                     │
│  ✓ Device health attestation                                                │
│  ✓ Compliance checking (patch level, AV)                                    │
│  ✓ Mobile Device Management (MDM)                                           │
│  ✓ Endpoint Detection and Response (EDR)                                    │
│                                                                             │
│  NETWORK                                                                    │
│  ✓ Micro-segmentation                                                       │
│  ✓ Software-Defined Perimeter (SDP)                                         │
│  ✓ Encrypted connections (TLS 1.3)                                          │
│  ✓ No implicit trust based on network location                              │
│                                                                             │
│  APPLICATION                                                                │
│  ✓ Service-to-service authentication (mTLS)                                 │
│  ✓ API gateway with authentication                                          │
│  ✓ Application-layer authorization                                          │
│  ✓ Rate limiting and DDoS protection                                        │
│                                                                             │
│  DATA                                                                       │
│  ✓ Data classification                                                      │
│  ✓ Encryption at rest and in transit                                        │
│  ✓ Data Loss Prevention (DLP)                                               │
│  ✓ Rights management                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 Authentication Method Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Authentication Method Matrix                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Scenario                │  Method          │  Additional Controls          │
├──────────────────────────┼──────────────────┼───────────────────────────────│
│  User access             │  OIDC/SAML SSO   │  MFA required                 │
│  (web/mobile apps)       │                  │  Risk-based step-up           │
│  ────────────────────────┼──────────────────┼───────────────────────────────│
│  Service-to-service      │  mTLS + JWT      │  Short-lived tokens           │
│  (internal APIs)         │                  │  Scope-based authorization    │
│  ────────────────────────┼──────────────────┼───────────────────────────────│
│  Machine-to-machine      │  Client certs    │  Certificate rotation         │
│  (automation)            │  or SPIFFE       │  IP allowlisting              │
│  ────────────────────────┼──────────────────┼───────────────────────────────│
│  Legacy integration      │  API keys        │  IP restrictions              │
│                          │                  │  Request signing              │
│  ────────────────────────┼──────────────────┼───────────────────────────────│
│  Admin access            │  Hardware tokens │  Break-glass procedures       │
│                          │  + MFA           │  Session recording            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Zero Trust Best Practices Summary                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PRINCIPLES                                                                 │
│  ✓ Never trust, always verify                                               │
│  ✓ Assume breach has occurred                                               │
│  ✓ Verify explicitly at every access point                                  │
│  ✓ Use least privilege access                                               │
│  ✓ Design for insider threat                                                │
│                                                                             │
│  IMPLEMENTATION                                                             │
│  ✓ Identity as the primary security perimeter                               │
│  ✓ Device health verification                                               │
│  ✓ Micro-segmentation of network                                            │
│  ✓ Encrypted communications everywhere                                        │
│  ✓ Continuous monitoring and validation                                       │
│                                                                             │
│  OPERATIONS                                                                 │
│  ✓ Regular access reviews and attestation                                   │
│  ✓ Automated threat response                                                │
│  ✓ Security analytics and behavioral monitoring                             │
│  ✓ Incident response integration                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. NIST SP 800-207 - Zero Trust Architecture
2. Google BeyondCorp Framework
3. Microsoft Zero Trust Architecture Guide
4. Forrester Zero Trust eXtended Framework
5. Gartner CARTA (Continuous Adaptive Risk and Trust Assessment)
