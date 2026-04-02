# Cloud-Native Engineer Learning Path

> **Version**: 1.0.0
> **Last Updated**: 2026-04-02
> **Duration**: 20 weeks (full-time) / 30 weeks (part-time)
> **Prerequisites**: Backend engineering experience, Go proficiency
> **Outcome**: Expert in cloud-native architecture, Kubernetes, and distributed platform engineering

---

## 🎯 Path Overview

### Target Competencies

Upon completion, you will be able to:

- Design and operate Kubernetes-native applications
- Build and maintain platform infrastructure
- Implement GitOps workflows and CI/CD pipelines
- Design service meshes and traffic management
- Build distributed schedulers and workflow engines
- Implement comprehensive observability stacks
- Design multi-tenant, multi-region architectures

### Prerequisites Graph

```
Backend Engineer Skills
    ├── Go proficiency
    ├── REST/gRPC APIs
    ├── Database design
    └── Testing practices
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│              CLOUD-NATIVE ENGINEER LEARNING PATH                     │
│                                                                      │
│  Phase 1: Container Foundations (Weeks 1-4)                         │
│    ├── Docker → Container Patterns → Image Security                 │
│    └── Outcome: Container-native applications                       │
│                                                                      │
│  Phase 2: Kubernetes Mastery (Weeks 5-9)                            │
│    ├── Core Concepts → Operators → Scheduling → Networking          │
│    └── Outcome: K8s-native service design                           │
│                                                                      │
│  Phase 3: Platform Engineering (Weeks 10-14)                        │
│    ├── Service Mesh → GitOps → Observability → Security             │
│    └── Outcome: Production platform                                 │
│                                                                      │
│  Phase 4: Distributed Systems (Weeks 15-20)                         │
│    ├── Schedulers → Workflows → Multi-Region → Cost Opt             │
│    └── Outcome: Large-scale distributed platforms                   │
└─────────────────────────────────────────────────────────────────────┘
    ↓
Advanced Paths
    ├── Distributed Systems Engineer
    └── Platform Architect
```

---

## 📚 Phase 1: Container Foundations (Weeks 1-4)

### Week 1: Docker Deep Dive

**Goal**: Master container fundamentals and best practices

#### Day 1-2: Container Design Principles

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-003] Container Design | 5h | Container patterns |
| [AD-008] Performance Optimization | 3h | Container perf |

**Study Notes**:

- Single process per container
- Immutable infrastructure
- Layer caching optimization
- Multi-stage builds

#### Day 3-4: Dockerfile Best Practices

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/03-Docker-Lib.md | 4h | Docker library |
| 03-Engineering-CloudNative/04-Security/07-Secure-Defaults.md | 3h | Security |

**Study Notes**:

- Minimal base images (distroless, scratch)
- Non-root user execution
- Secret management (BuildKit)
- Image scanning with Trivy/Snyk

#### Day 5-7: Container Runtime

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/03-Docker-Lib.md | 4h | Containerd integration |

**Practical Exercise**:

```dockerfile
# Build a production Go container:
# - Multi-stage build
# - Non-root user
# - Minimal attack surface
# - Health check included
```

### Week 2: Container Orchestration Concepts

**Goal**: Understand orchestration fundamentals

#### Day 1-3: Orchestration Theory

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-001] Architecture Principles | 4h | Cloud-native principles |
| [EC-002] Microservices Patterns | 4h | Service design |

**Study Notes**:

- 12-Factor App methodology
- Cattle vs pets
- Self-healing systems
- Declarative vs imperative

#### Day 4-5: Container Networking

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-006] Kubernetes Networking | 4h | Container networking |
| 04-Technology-Stack/03-Network/08-Load-Balancing.md | 3h | Load balancing |

**Study Notes**:

- Container network interfaces (CNI)
- Bridge vs overlay networks
- Service discovery
- Network policies

#### Day 6-7: Container Storage

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/03-Docker-Lib.md | 3h | Storage drivers |

**Study Notes**:

- Volume types (bind, volume, tmpfs)
- CSI (Container Storage Interface)
- Stateful vs stateless
- Data persistence patterns

### Week 3: Kubernetes Core Concepts

**Goal**: Master Kubernetes primitives

#### Day 1-3: Kubernetes Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-005] Kubernetes Operators | 4h | Core concepts |
| [TS-006] Kubernetes Networking | 4h | Networking model |

**Study Notes**:

- Master and worker nodes
- etcd as the brain
- API server and controllers
- Scheduler and kubelet

#### Day 4-5: Core Resources

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-099] Kubernetes CronJob Deep Dive | 4h | CronJobs |
| [EC-114] K8s CronJob Controller | 4h | Controller internals |

**Study Notes**:

- Pods: the atomic unit
- Deployments and ReplicaSets
- Services and Endpoints
- ConfigMaps and Secrets

#### Day 6-7: Advanced Resources

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-005] Kubernetes Operators | 4h | Custom resources |

**Study Notes**:

- StatefulSets for stateful apps
- DaemonSets for node agents
- Jobs and CronJobs
- CRDs (Custom Resource Definitions)

### Week 4: Kubernetes Application Design

**Goal**: Design cloud-native applications for K8s

#### Day 1-3: Pod Design Patterns

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-014] Sidecar Pattern | 4h | Sidecar implementation |
| [AD-003] Microservices Decomposition | 3h | Service boundaries |

**Study Notes**:

- Sidecar pattern (logging, proxy)
- Init containers
- Ambassador pattern
- Adapter pattern

#### Day 4-5: Configuration Management

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-040] Configuration Management | 4h | Config patterns |
| 05-Application-Domains/02-Cloud-Infrastructure/04-Helm-Charts.md | 3h | Helm templating |

**Study Notes**:

- ConfigMaps for configuration
- Secrets for sensitive data
- External configuration (external-secrets)
- Helm for packaging

#### Day 6-7: Resource Management

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-015] Resource Limits | 4h | Resource management |
| [EC-110] Resource Quota Management | 3h | Quotas |

**Study Notes**:

- Requests vs limits
- Quality of Service classes
- Resource quotas
- Limit ranges

---

## 📚 Phase 2: Kubernetes Mastery (Weeks 5-9)

### Week 5: Kubernetes Operators

**Goal**: Build custom controllers and operators

#### Day 1-3: Operator Pattern

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-005] Kubernetes Operators | 5h | Operator SDK |
| 05-Application-Domains/02-Cloud-Infrastructure/01-Kubernetes-Operators.md | 4h | Building operators |

**Study Notes**:

- Control loop pattern
- Custom Resource Definitions
- Controller-runtime library
- Operator SDK and kubebuilder

#### Day 4-5: Controller Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-114] K8s CronJob Controller | 4h | Real controller |
| [EC-099] K8s CronJob Deep Dive | 3h | Deep internals |

**Study Notes**:

- Informer pattern
- Work queues
- Reconciliation loop
- Finalizers

#### Day 6-7: Advanced Operator Patterns

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-005] Kubernetes Operators | 3h | Advanced patterns |

**Study Notes**:

- Leader election
- Webhooks (validation/mutation)
- Operator lifecycle manager
- Multi-version APIs

### Week 6: Kubernetes Scheduling

**Goal**: Understand and customize scheduling

#### Day 1-3: Scheduler Internals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-031] Scheduling Strategies | 4h | Scheduling |
| [EC-085] Resource Management | 3h | Resource scheduling |

**Study Notes**:

- Scheduling framework
- Predicates and priorities
- Node affinity/anti-affinity
- Pod affinity/anti-affinity

#### Day 4-5: Custom Scheduling

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-109] Production Task Scheduler | 4h | Custom scheduler |
| [EC-042] Scheduler Architecture | 3h | Architecture |

**Study Notes**:

- Scheduler extender
- Scheduling framework plugins
- Custom schedulers
- Gang scheduling

#### Day 6-7: Scheduling in Practice

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-017] Scheduled Task Framework | 3h | Task scheduling |
| [EC-020] Distributed Cron | 3h | Cron patterns |

### Week 7: Kubernetes Networking

**Goal**: Master container networking

#### Day 1-3: Network Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-006] Kubernetes Networking | 5h | CNI deep dive |
| [TS-007] Kubernetes Networking Deep Dive | 4h | Implementation |

**Study Notes**:

- CNI plugins (Calico, Cilium)
- Service types (ClusterIP, NodePort, LB)
- kube-proxy modes (iptables, IPVS, eBPF)
- CoreDNS

#### Day 4-5: Network Policies

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-035] Multi-Tenancy Isolation | 4h | Network isolation |
| [EC-093] Task Multi-Tenancy | 3h | Tenant separation |

**Study Notes**:

- NetworkPolicy resources
- Default deny policies
- Micro-segmentation
- Cilium network policies

#### Day 6-7: Ingress and Gateway

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/04-API-Gateway.md | 3h | Gateway patterns |
| [AD-006] API Gateway Design | 3h | Design principles |

**Study Notes**:

- Ingress controllers (nginx, traefik)
- Gateway API
- TLS termination
- Path and host routing

### Week 8: Kubernetes Storage

**Goal**: Manage persistent data in K8s

#### Day 1-3: Storage Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-005] Database Patterns | 4h | Stateful patterns |
| [EC-065] Transaction Isolation | 3h | Data consistency |

**Study Notes**:

- PV and PVC lifecycle
- Storage classes
- Dynamic provisioning
- CSI drivers

#### Day 4-5: Stateful Applications

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-001] PostgreSQL Internals | 3h | DB in K8s |
| [TS-002] Redis Internals | 3h | Redis in K8s |

**Study Notes**:

- StatefulSets
- Headless services
- Pod identity
- Rolling updates for stateful apps

#### Day 6-7: Backup and DR

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-105] Disaster Recovery Planning | 4h | DR strategies |
| 05-Application-Domains/03-DevOps-Tools/14-Backup-Recovery.md | 3h | Backup tools |

### Week 9: Kubernetes Security

**Goal**: Secure Kubernetes clusters and workloads

#### Day 1-3: Cluster Security

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/06-Zero-Trust.md | 4h | Zero trust |
| [AD-007] Security Patterns | 4h | Security patterns |

**Study Notes**:

- RBAC configuration
- Service accounts
- Pod security standards
- Pod security policies

#### Day 4-5: Supply Chain Security

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-045] Task Security Hardening | 4h | Hardening |
| 03-Engineering-CloudNative/04-Security/02-Vulnerability-Management.md | 3h | Vuln management |

**Study Notes**:

- Image signing (cosign)
- Admission controllers
- OPA/Gatekeeper policies
- Secret encryption

#### Day 6-7: Runtime Security

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-104] Security Hardening Checklist | 4h | Checklist |

**Study Notes**:

- Falco for runtime detection
- seccomp and AppArmor
- gVisor/Kata containers
- Security monitoring

---

## 📚 Phase 3: Platform Engineering (Weeks 10-14)

### Week 10: Service Mesh

**Goal**: Implement service-to-service communication

#### Day 1-3: Service Mesh Fundamentals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-015] Service Mesh Formal | 4h | Theory |
| [TS-015] Service Mesh Istio | 4h | Istio deep dive |
| 04-Technology-Stack/03-Network/09-Service-Mesh.md | 3h | Patterns |

**Study Notes**:

- Sidecar proxy pattern
- Data plane vs control plane
- mTLS automatically
- Traffic management

#### Day 4-5: Istio Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/06-Service-Mesh-Control.md | 4h | Istio config |

**Study Notes**:

- VirtualServices and DestinationRules
- Traffic splitting
- Circuit breaking at mesh level
- Retries and timeouts

#### Day 6-7: Advanced Mesh Features

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-056] Distributed Tracing Deep Dive | 3h | Mesh tracing |

**Study Notes**:

- Multi-cluster mesh
- External service access
- Egress gateways
- Mesh expansion

### Week 11: GitOps and CI/CD

**Goal**: Implement GitOps workflows

#### Day 1-3: GitOps Principles

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/08-GitOps.md | 4h | GitOps |
| 05-Application-Domains/03-DevOps-Tools/04-CI-CD.md | 3h | CI/CD |

**Study Notes**:

- Declarative infrastructure
- Git as single source of truth
- Automated synchronization
- Drift detection

#### Day 4-5: ArgoCD

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/04-CI-CD.md | 4h | ArgoCD setup |

**Study Notes**:

- Application definitions
- Sync waves and hooks
- App of Apps pattern
- Multi-source applications

#### Day 6-7: Progressive Delivery

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/10-Feature-Flags.md | 3h | Feature flags |

**Study Notes**:

- Argo Rollouts
- Canary deployments
- Blue-green deployments
- Analysis and promotion

### Week 12: Observability Stack

**Goal**: Build comprehensive observability

#### Day 1-3: Metrics and Alerting

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-013] Prometheus | 4h | Prometheus |
| [TS-013] Prometheus Formal | 3h | Theory |
| [EC-080] Metrics Integration | 3h | Integration |

**Study Notes**:

- Prometheus architecture
- ServiceMonitor resources
- Recording and alerting rules
- Alertmanager routing

#### Day 4-5: Logging

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-074] Context-Aware Logging | 3h | Logging |
| 05-Application-Domains/03-DevOps-Tools/05-Log-Analysis.md | 3h | Log analysis |

**Study Notes**:

- Fluent Bit / Fluentd
- Loki for log aggregation
- LogQL queries
- Structured logging

#### Day 6-7: Distributed Tracing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-060] OpenTelemetry Production | 4h | OTel |
| [EC-070] W3C Trace Context | 3h | Standards |

**Study Notes**:

- OpenTelemetry Collector
- Jaeger/Grafana Tempo
- Trace sampling strategies
- Correlating traces with logs/metrics

### Week 13: Platform Security

**Goal**: Implement comprehensive platform security

#### Day 1-3: Policy as Code

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/06-Zero-Trust.md | 4h | Zero trust arch |
| [EC-045] Security Hardening | 3h | Hardening |

**Study Notes**:

- OPA (Open Policy Agent)
- Gatekeeper constraints
- Kyverno policies
- Policy testing

#### Day 4-5: Secrets Management

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/04-Secrets-Management.md | 4h | Secrets |
| [EC-040] Configuration Management | 2h | Config |

**Study Notes**:

- External Secrets Operator
- Vault integration
- Sealed Secrets
- SOPS

#### Day 6-7: Compliance

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-104] Security Hardening Checklist | 3h | Compliance |

### Week 14: Multi-Tenancy

**Goal**: Design multi-tenant platforms

#### Day 1-3: Tenancy Models

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-035] Multi-Tenancy Isolation | 4h | Isolation |
| [EC-093] Task Multi-Tenancy | 4h | Implementation |

**Study Notes**:

- Namespace per tenant
- Virtual clusters
- SaaS tenancy models
- Data isolation strategies

#### Day 4-5: Resource Management

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-110] Resource Quota Management | 4h | Quotas |
| [EC-015] Resource Limits | 3h | Limits |

**Study Notes**:

- Resource quotas per tenant
- Limit ranges
- Fair scheduling
- Noisy neighbor prevention

#### Day 6-7: Cost Management

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/11-Cost-Management.md | 4h | Cost optimization |
| 05-Application-Domains/03-DevOps-Tools/13-Cost-Optimization.md | 3h | FinOps |

---

## 📚 Phase 4: Distributed Systems (Weeks 15-20)

### Week 15: Distributed Scheduling

**Goal**: Build distributed schedulers

#### Day 1-3: Scheduler Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-062] Distributed Task Scheduler | 5h | Architecture |
| [EC-042] Scheduler Core Architecture | 4h | Core design |

**Study Notes**:

- Scheduler components
- Task queue design
- Worker pool management
- Leader election

#### Day 4-5: ETCD Integration

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-057] ETCD Distributed Scheduler | 4h | ETCD scheduler |
| [EC-071] ETCD Coordination | 4h | Coordination |
| [EC-116] ETCD Coordination Patterns | 3h | Patterns |

**Study Notes**:

- ETCD for service discovery
- Distributed locking
- Lease management
- Watch API

#### Day 6-7: Production Scheduler

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-067] Production Task Scheduler | 4h | Production |
| [EC-109] Complete Scheduler | 4h | Full system |

### Week 16: Workflow Engines

**Goal**: Implement workflow orchestration

#### Day 1-3: Temporal Workflow

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-058] Temporal Workflow Engine | 4h | Temporal |
| [EC-069] Temporal Workflow | 4h | Patterns |
| [EC-100] Temporal Workflow | 4h | Advanced |
| [EC-115] Temporal Deep Dive | 4h | Internals |

**Study Notes**:

- Workflow as code
- Activity execution
- Durable execution
- Saga pattern with Temporal

#### Day 4-5: State Machine Workflows

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-012] State Machine Workflow | 3h | State machines |
| [EC-024] Task State Machine | 3h | Implementation |
| [EC-063] State Machine Implementation | 3h | Code |
| [EC-077] State Machine Execution | 3h | Execution |

#### Day 6-7: Cron and Scheduling

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-020] Distributed Cron | 3h | Cron |
| [EC-114] K8s CronJob Controller | 3h | Controller |

### Week 17: Event-Driven Architecture

**Goal**: Build event-driven systems

#### Day 1-3: Event Sourcing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-015] Event Sourcing | 4h | Patterns |
| [EC-034] Task Event Sourcing | 3h | Implementation |
| [EC-092] Event Sourcing Persistence | 3h | Persistence |
| [EC-111] Event Sourcing Implementation | 3h | Complete |

**Study Notes**:

- Event store design
- Event versioning
- Snapshotting
- Projection rebuild

#### Day 4-5: Kafka Deep Dive

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-003] Kafka KRaft Internals | 4h | Kafka internals |
| [TS-011] Kafka Internals | 4h | Deep dive |

**Study Notes**:

- KRaft mode (no ZooKeeper)
- Partition assignment
- Consumer groups
- Exactly-once semantics

#### Day 6-7: Event Processing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-004] Event-Driven Architecture | 4h | EDA patterns |
| 05-Application-Domains/02-Cloud-Infrastructure/07-Event-Driven-Architecture.md | 3h | Implementation |

### Week 18: Multi-Region and Edge

**Goal**: Design globally distributed systems

#### Day 1-3: Multi-Region Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/10-Multi-Cluster-Management.md | 4h | Multi-cluster |
| [EC-082] Distributed Task Sharding | 3h | Sharding |

**Study Notes**:

- Active-active vs active-passive
- Data replication strategies
- Global load balancing
- Latency optimization

#### Day 4-5: Edge Computing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/09-Edge-Computing.md | 4h | Edge patterns |

**Study Notes**:

- Edge deployment patterns
- Data synchronization
- Offline-first design
- Edge ML inference

#### Day 6-7: CRDTs

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-012] CRDTs | 4h | Theory |
| [EC-113] CRDT Conflict Resolution | 3h | Implementation |

### Week 19: Chaos Engineering

**Goal**: Implement chaos engineering practices

#### Day 1-3: Chaos Principles

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/08-Chaos-Engineering.md | 4h | Chaos engineering |
| [EC-102] Performance Benchmarking | 3h | Benchmarking |

**Study Notes**:

- Chaos Monkey principles
- Failure injection
- Blast radius control
- Abort conditions

#### Day 4-5: Chaos Tools

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/08-Chaos-Engineering.md | 4h | Litmus/Gremlin |

**Study Notes**:

- LitmusChaos
- Network chaos
- Pod chaos
- Stress testing

#### Day 6-7: Game Days

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-121] Google SRE | 4h | Reliability |

### Week 20: SRE and Platform Operations

**Goal**: Operate platforms at scale

#### Day 1-3: SRE Practices

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-121] Google SRE | 5h | SRE principles |
| 05-Application-Domains/03-DevOps-Tools/12-Platform-Engineering.md | 4h | Platform eng |

**Study Notes**:

- SLIs, SLOs, SLAs
- Error budgets
- Toil reduction
- Incident management

#### Day 4-5: AIOps

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/11-AIOps.md | 4h | ML for ops |

**Study Notes**:

- Anomaly detection
- Predictive scaling
- Automated remediation
- Log analysis with ML

#### Day 6-7: Platform as Product

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/12-Platform-Engineering.md | 4h | Platform team |

---

## 🎓 Capstone Project: Cloud-Native Platform

### Project: Multi-Tenant SaaS Platform

**Architecture**:

```
┌─────────────────────────────────────────────────────────────┐
│                        Ingress/Gateway                       │
│                    (Istio Gateway + Cert-Manager)           │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    Service Mesh (Istio)                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  API Gateway │  │   Tenant    │  │   Workflow Engine   │  │
│  │   Service    │  │   Service   │  │     (Temporal)      │  │
│  └──────┬───────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                 │                     │             │
│  ┌──────▼───────┐  ┌──────▼──────┐  ┌───────────▼────────┐  │
│  │  Scheduler   │  │  Event Bus  │  │  Distributed Task  │  │
│  │  Service     │  │   (Kafka)   │  │     Scheduler      │  │
│  └──────────────┘  └─────────────┘  └────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                  Observability Stack                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐   │
│  │Prometheus│  │  Loki    │  │  Jaeger  │  │  Grafana   │   │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Requirements**:

1. **Multi-Tenancy**
   - Namespace isolation per tenant
   - Resource quotas
   - Network policies
   - Data isolation

2. **Service Mesh**
   - mTLS between services
   - Traffic splitting for canary
   - Circuit breaking
   - Observability integration

3. **Workflow Engine**
   - Temporal for durable workflows
   - Saga pattern for transactions
   - Event-driven triggers

4. **Distributed Scheduler**
   - ETCD-backed coordination
   - Multi-tenant job scheduling
   - Cron support
   - Worker pool management

5. **Observability**
   - OpenTelemetry integration
   - Prometheus metrics
   - Distributed tracing
   - Structured logging

6. **GitOps**
   - ArgoCD for deployments
   - Helm charts
   - Progressive delivery
   - Drift detection

7. **Security**
   - OPA policies
   - Vault for secrets
   - Pod security standards
   - Network segmentation

---

## ✅ Progress Tracker

| Phase | Week | Topic | Complete |
|-------|------|-------|----------|
| 1 | 1 | Docker | [ ] |
| 1 | 2 | Orchestration | [ ] |
| 1 | 3 | K8s Core | [ ] |
| 1 | 4 | K8s Apps | [ ] |
| 2 | 5 | Operators | [ ] |
| 2 | 6 | Scheduling | [ ] |
| 2 | 7 | Networking | [ ] |
| 2 | 8 | Storage | [ ] |
| 2 | 9 | Security | [ ] |
| 3 | 10 | Service Mesh | [ ] |
| 3 | 11 | GitOps | [ ] |
| 3 | 12 | Observability | [ ] |
| 3 | 13 | Platform Security | [ ] |
| 3 | 14 | Multi-Tenancy | [ ] |
| 4 | 15 | Distributed Scheduling | [ ] |
| 4 | 16 | Workflows | [ ] |
| 4 | 17 | Event-Driven | [ ] |
| 4 | 18 | Multi-Region | [ ] |
| 4 | 19 | Chaos Engineering | [ ] |
| 4 | 20 | SRE | [ ] |

---

*This learning path prepares you for platform engineering and cloud-native architecture roles. The capstone project demonstrates mastery of modern distributed systems platform design.*
