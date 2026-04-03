# Network Policies

> **分类**: 工程与云原生
> **标签**: #network #kubernetes #security #microsegmentation #cni
> **参考**: Kubernetes Network Policies, Cilium, Calico, Istio

---

## 1. Formal Definition

### 1.1 What are Network Policies?

Network policies are specifications that define how groups of pods are allowed to communicate with each other and with other network endpoints. They provide a way to enforce network segmentation and micro-segmentation in Kubernetes clusters, implementing the principle of least privilege for network traffic.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Network Policy Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   WITHOUT NETWORK POLICIES              WITH NETWORK POLICIES               │
│   ━━━━━━━━━━━━━━━━━━━━━━━━━━            ━━━━━━━━━━━━━━━━━━━━━               │
│                                                                             │
│   ┌─────────────────────┐               ┌─────────────────────┐             │
│   │     Namespace       │               │     Namespace       │             │
│   │                     │               │                     │             │
│   │   ┌───┐   ┌───┐     │               │   ┌───┐   ┌───┐     │             │
│   │   │App│◄─►│DB │     │               │   │App│──►│DB │     │             │
│   │   └───┘   └───┘     │               │   └───┘   └───┘     │             │
│   │     ▲         ▲     │               │     │         ▲     │             │
│   │     │         │     │               │     │         │     │             │
│   │   ┌───┐     ┌───┐   │               │   ┌───┐     ┌───┐   │             │
│   │   │Web│◄────┤Attacker│             │   │Web│     │Cache│  │             │
│   │   └───┘     └───┘   │               │   └───┘     └───┘   │             │
│   │                     │               │                     │             │
│   │  [ALL TRAFFIC       │               │  [WHITELIST-BASED   │             │
│   │   ALLOWED]          │               │   TRAFFIC CONTROL]  │             │
│   └─────────────────────┘               └─────────────────────┘             │
│                                                                             │
│   VULNERABLE:                           SECURED:                            │
│   • Lateral movement possible           • Only authorized connections       │
│   • No traffic restrictions             • Default deny posture              │
│   • No audit trail                      • Explicit allow rules              │
│                                                                             │
│   POLICY ENFORCEMENT:                                                       │
│   • CNI Plugin: Calico, Cilium, Weave, Flannel                            │
│   • Service Mesh: Istio, Linkerd                                          │
│   • Cloud Provider: AWS Security Groups, Azure NSGs                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Network Policy Types

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Network Policy Types                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  LAYER 4 (Transport)              LAYER 7 (Application)                     │
│  ━━━━━━━━━━━━━━━━━━━━             ━━━━━━━━━━━━━━━━━━━━━                     │
│                                                                             │
│  Kubernetes Network Policies      Service Mesh Policies                     │
│  • Pod-to-pod rules               • HTTP path-based rules                   │
│  • Namespace isolation            • Method-based filtering                  │
│  • IP/CIDR blocks                 • Header-based routing                    │
│  • Port/protocol restrictions     • Rate limiting                           │
│                                                                             │
│  Example:                         Example:                                  │
│  ┌─────────────────────────┐      ┌─────────────────────────┐               │
│  │  ingress:               │      │  http:                  │               │
│  │  - from:                │      │  - match:               │               │
│  │    - podSelector:       │      │    - uri:               │               │
│  │        matchLabels:     │      │      prefix: /api/v1    │               │
│  │          app: frontend  │      │  - route:               │               │
│  │    ports:               │      │    - destination:       │               │
│  │    - protocol: TCP      │      │        host: backend    │               │
│  │      port: 8080         │      │  - fault:               │               │
│  └─────────────────────────┘      │    delay: 5s            │               │
│                                   └─────────────────────────┘               │
│                                                                             │
│  CNI-Specific Extensions                                                    │
│  • Cilium: eBPF-based filtering, DNS policies, L7                         │
│  • Calico: Global policies, host endpoints, WireGuard                     │
│  • Antrea: Tiered policies, traceflow, network visibility                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 Default Deny Policy

```yaml
# default-deny-all.yaml
# Default deny all ingress and egress traffic
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
  namespace: production
spec:
  podSelector: {}  # Applies to all pods
  policyTypes:
  - Ingress
  # No ingress rules = deny all incoming traffic
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-egress
  namespace: production
spec:
  podSelector: {}  # Applies to all pods
  policyTypes:
  - Egress
  # No egress rules = deny all outgoing traffic
```

### 2.2 Application-Specific Policies

```yaml
# web-tier-policy.yaml
# Allow web tier to receive traffic from ingress controller
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-tier-ingress
  namespace: production
spec:
  podSelector:
    matchLabels:
      tier: web
      app: frontend
  policyTypes:
  - Ingress
  ingress:
  # Allow from ingress controller
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    - podSelector:
        matchLabels:
          app.kubernetes.io/name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  # Allow health checks from monitoring
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring
    - podSelector:
        matchLabels:
          app: prometheus
    ports:
    - protocol: TCP
      port: 8080
---
# api-tier-policy.yaml
# Allow API tier to receive traffic from web tier and access database
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-tier-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      tier: api
      app: backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow from web tier
  - from:
    - podSelector:
        matchLabels:
          tier: web
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
  # Allow from internal tools
  - from:
    - namespaceSelector:
        matchLabels:
          name: tools
    - podSelector:
        matchLabels:
          app: admin-portal
    ports:
    - protocol: TCP
      port: 8080
  egress:
  # Allow to database
  - to:
    - podSelector:
        matchLabels:
          tier: database
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  # Allow to cache
  - to:
    - podSelector:
        matchLabels:
          tier: cache
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  # Allow to message queue
  - to:
    - podSelector:
        matchLabels:
          tier: messaging
          app: kafka
    ports:
    - protocol: TCP
      port: 9092
  # Allow DNS
  - to:
    - namespaceSelector: {}
      podSelector:
        matchLabels:
          k8s-app: kube-dns
    ports:
    - protocol: UDP
      port: 53
```

### 2.3 Cilium Layer 7 Policies

```yaml
# cilium-l7-policy.yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: api-l7-policy
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: api-gateway
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: frontend
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: GET
          path: "/api/v1/users/.*"
          headers:
          - name: Authorization
            presence: true
        - method: POST
          path: "/api/v1/orders"
          headers:
          - name: Content-Type
            value: application/json
        - method: GET
          path: "/health"
  egress:
  - toEndpoints:
    - matchLabels:
        app: user-service
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
        - method: GET
          path: "/users/.*"
  - toEndpoints:
    - matchLabels:
        app: order-service
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
        - method: "*"
          path: "/orders/.*"
---
# Cilium DNS policy
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: egress-dns
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: backend
  egress:
  - toFQDNs:
    - matchName: api.stripe.com
    - matchName: api.sendgrid.com
    toPorts:
    - ports:
      - port: "443"
        protocol: TCP
  - toEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: kube-system
        k8s-app: kube-dns
    toPorts:
    - ports:
      - port: "53"
        protocol: UDP
      rules:
        dns:
        - matchPattern: "*.stripe.com"
        - matchPattern: "*.sendgrid.com"
```

### 2.4 Calico Global Policies

```yaml
# calico-global-policy.yaml
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: default-deny
spec:
  order: 1000
  selector: all()
  types:
  - Ingress
  - Egress
  # Default deny - no rules means deny all
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: allow-dns
spec:
  order: 100
  selector: all()
  types:
  - Egress
  egress:
  - action: Allow
    protocol: UDP
    destination:
      selector: k8s-app == 'kube-dns'
      ports:
      - 53
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: allow-ingress-nginx
spec:
  order: 200
  selector: app.kubernetes.io/name == 'ingress-nginx'
  types:
  - Ingress
  ingress:
  - action: Allow
    source:
      nets:
      - 0.0.0.0/0
    destination:
      ports:
      - 80
      - 443
```

---

## 3. Production-Ready Configurations

### 3.1 Defense in Depth Network Policy Set

```yaml
# namespace-isolation.yaml
# Isolate namespaces from each other by default
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: namespace-isolation
  namespace: production
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Only allow traffic from same namespace
  - from:
    - podSelector: {}
  egress:
  # Only allow traffic to same namespace
  - to:
    - podSelector: {}
  # Allow DNS to kube-system
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
      podSelector:
        matchLabels:
          k8s-app: kube-dns
    ports:
    - protocol: UDP
      port: 53
---
# allow-cross-namespace.yaml
# Explicitly allow specific cross-namespace traffic
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-monitoring
  namespace: production
spec:
  podSelector:
    matchLabels:
      monitoring: enabled
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring
    - podSelector:
        matchLabels:
          app: prometheus
    ports:
    - protocol: TCP
      port: 9090
    - protocol: TCP
      port: 9091
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Network Policy Security Checklist                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DEFAULT POSTURE                                                            │
│  ✓ Default deny all ingress traffic                                         │
│  ✓ Default deny all egress traffic                                          │
│  ✓ Explicit allow rules only                                                │
│  ✓ Regular review of exceptions                                             │
│                                                                             │
│  SEGMENTATION                                                               │
│  ✓ Namespace-level isolation                                                │
│  ✓ Pod-to-pod microsegmentation                                             │
│  ✓ Environment separation (dev/staging/prod)                                │
│  ✓ Tier-based segmentation (web/app/db)                                     │
│                                                                             │
│  MONITORING                                                                 │
│  ✓ Network flow logging                                                     │
│  ✓ Policy violation alerts                                                  │
│  ✓ Traffic visualization                                                    │
│  ✓ Anomaly detection                                                        │
│                                                                             │
│  TESTING                                                                    │
│  ✓ Validate policies in staging                                             │
│  ✓ Test connectivity between services                                       │
│  ✓ Verify policy coverage (no blind spots)                                  │
│  ✓ Regular policy audits                                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 CNI Selection Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CNI Comparison Matrix                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Feature            │  Cilium  │  Calico  │  Flannel │  Weave   │  AWS VPC │
├─────────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Network Policy     │  ★★★★★   │  ★★★★★   │  ★★★☆☆   │  ★★★☆☆   │  ★★★★☆   │
│  (L3/L4)            │          │          │          │          │          │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  L7 Policies        │  ★★★★★   │  ★★★☆☆   │  ★☆☆☆☆   │  ★☆☆☆☆   │  ★☆☆☆☆   │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Encryption         │  ★★★★★   │  ★★★★★   │  ★☆☆☆☆   │  ★★★★☆   │  N/A     │
│  (WireGuard/IPSec)  │          │          │          │          │          │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Observability      │  ★★★★★   │  ★★★★☆   │  ★★☆☆☆   │  ★★★☆☆   │  ★★☆☆☆   │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Performance        │  ★★★★★   │  ★★★★☆   │  ★★★★☆   │  ★★★☆☆   │  ★★★★★   │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Cluster Mesh       │  ★★★★★   │  ★★★★★   │  ★☆☆☆☆   │  ★★☆☆☆   │  N/A     │
│  ───────────────────┼──────────┼──────────┼──────────┼──────────┼──────────│
│  Service Mesh       │  ★★★★★   │  ★★☆☆☆   │  ★☆☆☆☆   │  ★☆☆☆☆   │  ★☆☆☆☆   │
│  Integration        │          │          │          │          │          │
│                                                                             │
│  Recommendation:                                                            │
│  • Security-first: Cilium                                                   │
│  • Multi-cluster: Calico or Cilium                                          │
│  • Simple networking: Flannel                                               │
│  • AWS-specific: VPC CNI + Security Groups for Pods                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Network Policy Best Practices Summary                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DESIGN                                                                     │
│  ✓ Start with default deny policies                                         │
│  ✓ Use explicit allow rules                                                 │
│  ✓ Label pods consistently for policy selection                             │
│  ✓ Document all policy exceptions                                           │
│                                                                             │
│  IMPLEMENTATION                                                             │
│  ✓ Enable network policies in all namespaces                                │
│  ✓ Use CNI with full policy support                                         │
│  ✓ Implement defense in depth (L3/L4/L7)                                    │
│  ✓ Test policies in non-production first                                    │
│                                                                             │
│  OPERATIONS                                                                 │
│  ✓ Monitor policy violations                                                │
│  ✓ Regular policy audits                                                    │
│  ✓ Automated policy testing                                                 │
│  ✓ Visualize network traffic                                                │
│                                                                             │
│  TROUBLESHOOTING                                                            │
│  ✓ Use policy simulation tools                                              │
│  ✓ Enable flow logging                                                      │
│  ✓ Test connectivity before enforcing                                       │
│  ✓ Keep emergency bypass procedures                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Kubernetes Network Policies Documentation
2. Cilium Documentation
3. Calico Network Policy Guide
4. NIST SP 800-207 - Zero Trust Architecture
5. OWASP Kubernetes Security Guide
