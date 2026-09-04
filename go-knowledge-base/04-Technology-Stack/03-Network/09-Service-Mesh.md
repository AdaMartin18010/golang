# TS-NET-009: Service Mesh Architecture (Istio/Linkerd)

> **维度**: Technology Stack > Network
> **级别**: S (29 KB)
> **标签**: #service-mesh #istio #linkerd #microservices #sidecar
> **权威来源**:
>
> - [Istio Documentation](https://istio.io/latest/docs/) - Istio
> - [Linkerd Documentation](https://linkerd.io/2/overview/) - Linkerd
> - [Service Mesh Interface](https://smi-spec.io/) - SMI Spec

> **Go 版本**: 1.27+
---

## 1. Service Mesh Architecture

### 1.1 Core Concept

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Service Mesh Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  WITHOUT Service Mesh:                                                      │
│  ┌─────────┐         TLS         ┌─────────┐                               │
│  │ Service │◄────────────────────►│ Service │                               │
│  │    A    │    Retry Logic      │    B    │                               │
│  └─────────┘    Circuit Breaker  └─────────┘                               │
│                 Metrics/Tracing                                             │
│                 (Implemented in each service)                               │
│                                                                              │
│  WITH Service Mesh:                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Kubernetes Pod                               │   │
│  │  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐         │   │
│  │  │   Service   │◄────►│   Sidecar   │◄────►│   Service   │         │   │
│  │  │     A       │ IPC  │   Proxy     │ mTLS │     B       │         │   │
│  │  │             │      │(Envoy/Link2d)│     │             │         │   │
│  │  └─────────────┘      └──────┬──────┘      └─────────────┘         │   │
│  │                              │                                       │   │
│  │                         ┌────┴────┐                                  │   │
│  │                         │ Control │                                  │   │
│  │                         │  Plane  │                                  │   │
│  │                         └─────────┘                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Service Mesh Layer:                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Traffic Management     Security        Observability              │   │
│  │  ├── Routing            ├── mTLS        ├── Metrics               │   │
│  │  ├── Load Balancing     ├── AuthZ       ├── Distributed Tracing   │   │
│  │  ├── Canary/Blue-Green  └── Identity    └── Logging               │   │
│  │  ├── Circuit Breaking                                               │   │
│  │  └── Retries/Timeouts                                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Sidecar Pattern

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Sidecar Pattern Details                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          Kubernetes Pod                              │   │
│  │                                                                      │   │
│  │   ┌─────────────────────────┐    ┌─────────────────────────┐        │   │
│  │   │    Application          │    │      Sidecar Proxy      │        │   │
│  │   │    Container            │    │      (Envoy)            │        │   │
│  │   │                         │    │                         │        │   │
│  │   │  Port: 8080             │◄──►│  Port: 15001 (inbound)  │        │   │
│  │   │                         │    │  Port: 15006 (outbound) │        │   │
│  │   │  Code:                  │    │                         │        │   │
│  │   │  - Business logic only  │    │  Functions:             │        │   │
│  │   │  - HTTP handlers        │    │  - Traffic routing      │        │   │
│  │   │  - gRPC services        │    │  - Load balancing       │        │   │
│  │   │                         │    │  - mTLS termination     │        │   │
│  │   │  No need to know        │    │  - AuthN/AuthZ          │        │   │
│  │   │  about the mesh!        │    │  - Circuit breaking     │        │   │
│  │   │                         │    │  - Retry/Timeout        │        │   │
│  │   └─────────────────────────┘    │  - Metrics collection   │        │   │
│  │              │                   └───────────┬─────────────┘        │   │
│  │              │                               │                      │   │
│  │              └───────────────IPC (localhost)─┘                      │   │
│  │                                                                      │   │
│  │   Shared Network Namespace: iptables rules redirect traffic         │   │
│  │   - Outbound: 8080 → 15006 (sidecar outbound)                       │   │
│  │   - Inbound:  external → 15001 → 8080 (sidecar → app)              │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Traffic Interception (iptables REDIRECT):                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  PREROUTING chain:                                                  │   │
│  │    - Destination port 80/443 → redirect to 15006 (sidecar)         │   │
│  │                                                                     │   │
│  │  OUTPUT chain:                                                      │   │
│  │    - Local process outbound → redirect to 15006                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Istio Architecture

### 2.1 Components

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Istio Architecture                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Control Plane (istiod)                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   Pilot      │  │  Citadel     │  │   Galley     │              │   │
│  │  │              │  │              │  │              │              │   │
│  │  │ - Service    │  │ - Certificate│  │ - Config     │              │   │
│  │  │   discovery  │  │   management │  │   validation │              │   │
│  │  │ - Config     │  │ - Key/Cert   │  │ - Distribution│             │   │
│  │  │   distribution│  │   rotation   │  │              │              │   │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │   │
│  │         │                 │                 │                       │   │
│  │         └─────────────────┴─────────────────┘                       │   │
│  │                           │                                         │   │
│  │                    ┌──────┴──────┐                                  │   │
│  │                    │    xDS API    │ (ADS, CDS, EDS, LDS, RDS)      │   │
│  │                    │  (Envoy API)  │                                  │   │
│  │                    └──────┬──────┘                                  │   │
│  │                           │                                         │   │
│  └───────────────────────────┼─────────────────────────────────────────┘   │
│                              │                                             │
│  ┌───────────────────────────┼─────────────────────────────────────────┐   │
│  │                      Data Plane                                        │   │
│  │  ┌────────────────────────┼─────────────────────────────────────┐   │   │
│  │  │                        ▼                                     │   │   │
│  │  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │   │   │
│  │  │   │Envoy Sidecar│  │Envoy Sidecar│  │Envoy Sidecar│         │   │   │
│  │  │   │   (Pod 1)   │  │   (Pod 2)   │  │   (Pod N)   │         │   │   │
│  │  │   └─────────────┘  └─────────────┘  └─────────────┘         │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Envoy Configuration (xDS):                                                 │
│  - CDS (Cluster Discovery Service): Service clusters                       │
│  - EDS (Endpoint Discovery Service): Service endpoints                     │
│  - LDS (Listener Discovery Service): Traffic listeners                     │
│  - RDS (Route Discovery Service): Traffic routes                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Traffic Management

```yaml
# VirtualService - Traffic routing
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: user-service
spec:
  hosts:
  - user-service
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: user-service
        subset: v2
      weight: 100
  - route:
    - destination:
        host: user-service
        subset: v1
      weight: 90
    - destination:
        host: user-service
        subset: v2
      weight: 10
    retries:
      attempts: 3
      perTryTimeout: 2s
    timeout: 10s
    fault:
      delay:
        percentage:
          value: 0.1
        fixedDelay: 5s

---
# DestinationRule - Traffic policies
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: user-service
spec:
  host: user-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
    loadBalancer:
      simple: LEAST_CONN
    outlierDetection:
      consecutiveErrors: 5
      interval: 30s
      baseEjectionTime: 30s
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
    trafficPolicy:
      connectionPool:
        http:
          http2MaxRequests: 500

---
# Gateway - Ingress traffic
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: public-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: tls-cert
    hosts:
    - "*.example.com"
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*.example.com"
    tls:
      httpsRedirect: true
```

### 2.3 Security (mTLS)

```yaml
# PeerAuthentication - mTLS policy
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT  # Require mTLS for all traffic

---
# AuthorizationPolicy - Access control
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: user-service-policy
  namespace: default
spec:
  selector:
    matchLabels:
      app: user-service
  action: ALLOW
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/default/sa/api-gateway"]
    to:
    - operation:
        methods: ["GET"]
        paths: ["/api/users/*"]
  - from:
    - source:
        principals: ["cluster.local/ns/default/sa/admin-service"]
    to:
    - operation:
        methods: ["*"]
        paths: ["/api/admin/*"]

---
# RequestAuthentication - JWT validation
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: jwt-auth
  namespace: default
spec:
  selector:
    matchLabels:
      app: user-service
  jwtRules:
  - issuer: "https://auth.example.com"
    jwksUri: "https://auth.example.com/.well-known/jwks.json"
    audiences:
    - "user-service"
    forwardOriginalToken: true
```

---

## 3. Linkerd Architecture

### 3.1 Components

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Linkerd Architecture                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Control Plane                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │   │
│  │  │  Controller │  │   Proxy     │  │   Identity  │                 │   │
│  │  │  (destination)│  │  Injector  │  │             │                 │   │
│  │  │             │  │             │  │ - CA for mTLS│                │   │
│  │  │ - Service   │  │ - Sidecar   │  │ - Certificate│                │   │
│  │  │   discovery │  │   injection │  │   issuance   │                │   │
│  │  └──────┬──────┘  └─────────────┘  └─────────────┘                 │   │
│  │         │                                                          │   │
│  │  ┌──────┴──────┐  ┌─────────────┐  ┌─────────────┐                 │   │
│  │  │   Destination│  │  Tap       │  │   Web       │                 │   │
│  │  │   API       │  │  - Real-time │  │   Dashboard │                 │   │
│  │  │             │  │    traffic   │  │             │                 │   │
│  │  └─────────────┘  │    tapping   │  └─────────────┘                 │   │
│  │                   └─────────────┘                                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                             │
│  ┌───────────────────────────┼─────────────────────────────────────────┐   │
│  │                      Data Plane                                        │   │
│  │  ┌────────────────────────┼─────────────────────────────────────┐   │   │
│  │  │                        ▼                                     │   │   │
│  │  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │   │   │
│  │  │   │Linkerd2 Proxy│  │Linkerd2 Proxy│  │Linkerd2 Proxy│        │   │   │
│  │  │   │  (Rust)     │  │  (Rust)     │  │  (Rust)     │         │   │   │
│  │  │   │             │  │             │  │             │         │   │   │
│  │  │   │ - Ultra-light│  │ - Ultra-light│  │ - Ultra-light│        │   │   │
│  │  │   │ - <10MB RSS │  │ - <10MB RSS │  │ - <10MB RSS │         │   │   │
│  │  │   │ - <1ms p99  │  │ - <1ms p99  │  │ - <1ms p99  │         │   │   │
│  │  │   │   latency   │  │   latency   │  │   latency   │         │   │   │
│  │  │   └─────────────┘  └─────────────┘  └─────────────┘         │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Linkerd2 Proxy Features:                                                   │
│  - Written in Rust for safety and performance                               │
│  - Minimal resource overhead                                                 │
│  - Built-in Prometheus metrics                                              │
│  - Automatic mTLS with zero config                                          │
│  - L7 load balancing ( EWMA, Least Loaded)                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Service Mesh Patterns

### 4.1 Traffic Splitting (Canary Deployment)

```yaml
# Istio: Canary deployment
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: canary-deployment
spec:
  hosts:
  - reviews
  http:
  - match:
    - headers:
        cookie:
          regex: ^(.*;?)?(user=tester)(;.*)?$
    route:
    - destination:
        host: reviews
        subset: v2
  - route:
    - destination:
        host: reviews
        subset: v1
      weight: 95
    - destination:
        host: reviews
        subset: v2
      weight: 5
```

### 4.2 Circuit Breaking

```yaml
# Istio: Circuit breaker
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: circuit-breaker
spec:
  host: ratings
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
        consecutiveErrors: 5  # Circuit breaker threshold
        interval: 10s
        baseEjectionTime: 30s
    outlierDetection:
      consecutiveErrors: 5
      interval: 10s
      baseEjectionTime: 30s
      maxEjectionPercent: 50
```

### 4.3 Retry and Timeout

```yaml
# Istio: Retry policy
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: retry-policy
spec:
  hosts:
  - ratings
  http:
  - route:
    - destination:
        host: ratings
        subset: v1
    retries:
      attempts: 3
      perTryTimeout: 2s
      retryOn: 5xx,connect-failure,refused-stream
    timeout: 10s
```

---

## 5. Observability

### 5.1 Metrics

```yaml
# Istio: Telemetry configuration
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: mesh-default
spec:
  metrics:
  - providers:
    - name: prometheus
    overrides:
    - match:
        metric: REQUEST_COUNT
      tagOverrides:
        destination_port:
          operation: REMOVE
  accessLogging:
  - providers:
    - name: envoy
    filter:
      expression: "response.code >= 400"
```

### 5.2 Distributed Tracing

```yaml
# Istio: Tracing configuration
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
    defaultConfig:
      tracing:
        sampling: 100.0  # Percentage
        zipkin:
          address: jaeger-collector.istio-system:9411
```

---

## 6. Go Integration

### 6.1 Service without Mesh Knowledge

```go
// Application code remains unchanged
// The mesh handles cross-cutting concerns

package main

import (
    "context"
    "log"
    "net/http"
    "time"
)

type UserService struct{}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
    // Business logic only
    // No circuit breaker code
    // No retry logic
    // No metrics code

    user, err := fetchUserFromDB(r.Context(), r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    respondJSON(w, user)
}

func main() {
    svc := &UserService{}
    http.HandleFunc("/user", svc.GetUser)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 7. Checklist

```
Service Mesh Implementation Checklist:

Planning:
□ Evaluate Istio vs Linkerd vs Consul Connect
□ Plan cluster capacity (sidecar resource overhead)
□ Design mTLS rollout strategy (PERMISSIVE → STRICT)
□ Plan traffic management policies

Installation:
□ Install control plane
□ Configure automatic sidecar injection
□ Deploy gateway for ingress
□ Verify mTLS is working

Configuration:
□ Define VirtualServices for routing
□ Configure DestinationRules for policies
□ Set up AuthorizationPolicies
□ Configure rate limiting

Observability:
□ Deploy Prometheus/Grafana
□ Configure distributed tracing
□ Set up Kiali for topology visualization
□ Configure access logging

Production:
□ Test circuit breakers
□ Verify retry policies
□ Load test with sidecars
□ Monitor resource usage
```
