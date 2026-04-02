# EC-047: Process Injector Pattern (Sidecar & DaemonSet)

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #sidecar-pattern #daemonset #process-injection #service-mesh #observability #security
> **Authoritative Sources**:
>
> - [Kubernetes Patterns](https://k8s-patterns.io/) - Ibryam & Huß (2019)
> - [Site Reliability Engineering](https://sre.google/sre-book/table-of-contents/) - Google (2017)
> - [The Datacenter as a Computer](https://www.morganclaypool.com/doi/abs/10.2200/S00874ED3V01Y201809CAC046) - Barroso et al. (2018)
> - [Service Mesh Patterns](https://layer5.io/books/service-mesh-patterns) - Layer5 (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Containerized Application Topology)**
Let $\mathcal{P} = \{P_1, P_2, ..., P_n\}$ be a set of application pods where each pod $P_i$ consists of:

- Primary container $C_{primary}$: Business logic execution
- Node resources $R_{node}$: Shared kernel, networking, storage

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Container Isolation** | $\forall C_i, C_j: C_i \perp C_j$ (namespace isolation) | Cannot directly share memory or signals |
| **Single Responsibility** | $|C_{concerns}| = 1$ per container image | Cross-cutting concerns need externalization |
| **Immutable Infrastructure** | $\nexists modify(C_{image})$ at runtime | Runtime extensions must be external |
| **Resource Boundaries** | $\sum_{C \in P} r(C) \leq R_{limit}$ | Additional processes compete for resources |
| **Security Posture** | $privilege(C_{sidecar}) \neq privilege(C_{app})$ | Privilege separation required |

### 1.2 Problem Statement

**Problem 1.1 (Cross-Cutting Concern Injection)**
Given a primary application $A$ requiring auxiliary functionality set $\Phi = \{\phi_1, \phi_2, ..., \phi_k\}$ where each $\phi_i \notin domain(A)$, design an injection mechanism $I$ such that:

$$\forall \phi \in \Phi: I(A, \phi) \Rightarrow A \models \phi \land A \perp \phi$$

**Key Challenges:**

1. **Lifecycle Coupling**: Ensure $\forall \phi: lifecycle(\phi) \propto lifecycle(A)$
2. **Resource Sharing**: Enable $\exists R_{shared}: A \leftrightarrow \phi$ for necessary communication
3. **Deployment Decoupling**: Maintain $\forall \phi: deploy(\phi) \nLeftrightarrow deploy(A)$
4. **Security Isolation**: Enforce $privilege(A) \cap privilege(\phi) = \emptyset$ where required
5. **Observability Propagation**: Guarantee $\forall event_e \in A: correlated(event_e, \phi)$

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Co-location Invariant)**
$$\forall P_i: C_{primary} \in P_i \Leftrightarrow C_{sidecar} \in P_i$$

**Requirement 1.2 (Shared Namespace)**
$$\forall P_i: network(C_{primary}) = network(C_{sidecar}) \land ipc(C_{primary}) \cap ipc(C_{sidecar}) \neq \emptyset$$

**Requirement 1.3 (Independent Lifecycle)**
$$restart(C_{sidecar}) \nRightarrow restart(C_{primary}) \land upgrade(C_{sidecar}) \nRightarrow upgrade(C_{primary})$$

---

## 2. Solution Architecture

### 2.1 Formal Process Injector Definition

**Definition 2.1 (Process Injector)**
A Process Injector $PI$ is a 6-tuple $\langle M, T, L, S, R, O \rangle$:

- $M \in \{sidecar, init, daemonset, ephemeral\}$: Injection mode
- $T = \{observability, security, networking, proxy\}$: Concern type
- $L: P \times C \to \{co\_located, node\_scoped, cluster\_scoped\}$: Location function
- $S = \{shared\_pid, shared\_net, shared\_ipc, shared\_uts, shared\_mount\}$: Namespace sharing
- $R = \{cpu, memory, storage, network\}$: Resource constraints
- $O$: Observability hooks for distributed tracing

### 2.2 Injection Pattern Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROCESS INJECTION PATTERNS                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         SIDECAR PATTERN                              │   │
│  │                                                                        │   │
│  │   ┌─────────────┐        ┌─────────────┐        ┌─────────────┐      │   │
│  │   │             │◄──────►│   Shared    │◄──────►│             │      │   │
│  │   │   Primary   │  IPC   │  Network    │ Files  │   Sidecar   │      │   │
│  │   │  Container  │        │  Namespace  │        │  Container  │      │   │
│  │   │             │        │             │        │             │      │   │
│  │   │ • App Logic │        │ • localhost │        │ • Logging   │      │   │
│  │   │ • Business  │        │ • Same IP   │        │ • Metrics   │      │   │
│  │   │   Logic     │        │ • Port      │        │ • Proxy     │      │   │
│  │   │             │        │   Binding   │        │ • Security  │      │   │
│  │   └─────────────┘        └─────────────┘        └─────────────┘      │   │
│  │                                                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        INIT CONTAINER PATTERN                        │   │
│  │                                                                        │   │
│  │   ┌─────────────┐       ┌─────────────┐       ┌─────────────┐        │   │
│  │   │   Init C1   │──────►│   Init C2   │──────►│   Primary   │        │   │
│  │   │             │       │             │       │  Container  │        │   │
│  │   │ • Config    │       │ • Security  │       │             │        │   │
│  │   │   Setup     │       │   Scan      │       │ • App runs  │        │   │
│  │   │ • Data Load │       │ • Cert Gen  │       │   after init│        │   │
│  │   └─────────────┘       └─────────────┘       └─────────────┘        │   │
│  │                                                                         │   │
│  │   Execution: Sequential → Terminates before primary starts             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        DAEMONSET PATTERN                             │   │
│  │                                                                        │   │
│  │    Node 1          Node 2          Node 3         ...    Node N       │   │
│  │   ┌───────┐       ┌───────┐       ┌───────┐          ┌───────┐       │   │
│  │   │ Pod 1 │       │ Pod 2 │       │ Pod 3 │          │ Pod N │       │   │
│  │   │ Pod 2 │       │ Pod 3 │       │ Pod 1 │          │ Pod 2 │       │   │
│  │   └───────┘       └───────┘       └───────┘          └───────┘       │   │
│  │      │               │               │                  │             │   │
│  │   ┌──┴───┐        ┌──┴───┐        ┌──┴───┐           ┌──┴───┐        │   │
│  │   │Daemon│        │Daemon│        │Daemon│           │Daemon│        │   │
│  │   │Agent │        │Agent │        │Agent │           │Agent │        │   │
│  │   │      │        │      │        │      │           │      │        │   │
│  │   │• Logs│        │• Logs│        │• Logs│           │• Logs│        │   │
│  │   │•Node │        │•Node │        │•Node │           │•Node │        │   │
│  │   │ Mon  │        │ Mon  │        │ Mon  │           │ Mon  │        │   │
│  │   └──────┘        └──────┘        └──────┘           └──────┘        │   │
│  │                                                                         │   │
│  │   Execution: One per node → Runs continuously                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      EPHEMERAL CONTAINER PATTERN                     │   │
│  │                                                                        │   │
│  │   ┌─────────────┐       ┌─────────────┐       ┌─────────────┐        │   │
│  │   │   Primary   │◄──────│  Ephemeral  │       │   Primary   │        │   │
│  │   │  Container  │ Debug │  Container  │ Debug │  Container  │        │   │
│  │   │  (running)  │──────►│  (kubectl   │──────►│  (running)  │        │   │
│  │   │             │       │   debug)    │       │             │        │   │
│  │   │             │       │             │       │             │        │   │
│  │   │ • Live      │       │ • gdb       │       │ • Profile   │        │   │
│  │   │   process   │       │ • tcpdump   │       │ • Continue  │        │   │
│  │   │   attached  │       │ • shell     │       │             │        │   │
│  │   └─────────────┘       └─────────────┘       └─────────────┘        │   │
│  │                                                                         │   │
│  │   Execution: On-demand → Shares process namespace                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Visual Representations

### 3.1 Sidecar Pattern Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              KUBERNETES POD                                  │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                           Pod Sandbox (Pause)                          │  │
│  │              (PID 1, Network Namespace, IPC Namespace)                 │  │
│  └───────────────────────────────────┬───────────────────────────────────┘  │
│                                      │                                      │
│  ════════════════════════════════════╪══════════════════════════════════    │
│                                      │                                      │
│  ┌─────────────────────┐   Shared   │   ┌─────────────────────┐            │
│  │                     │◄───────────┼───►│                     │            │
│  │  PRIMARY CONTAINER  │ Network    │   │   SIDECAR CONTAINER │            │
│  │                     │ Namespace  │   │                     │            │
│  │  ┌───────────────┐  │            │   │  ┌───────────────┐  │            │
│  │  │  Application  │  │  Loopback  │   │  │  Envoy Proxy  │  │            │
│  │  │     Code      │  │  127.0.0.1 │   │  │               │  │            │
│  │  └───────┬───────┘  │            │   │  │ • mTLS         │  │            │
│  │          │          │  IPC       │   │  │ • Rate Limit   │  │            │
│  │  ┌───────▼───────┐  │  Shared    │   │  │ • Load Balance │  │            │
│  │  │   HTTP/gRPC   │  │  Memory    │   │  │ • Retry        │  │            │
│  │  │    Client     │──┼────────────┼──►│  │ • Circuit      │  │            │
│  │  │               │  │  (emptyDir)│   │  │   Breaker      │  │            │
│  │  └───────────────┘  │            │   │  └───────┬───────┘  │            │
│  │                     │            │   │          │          │            │
│  │  ┌───────────────┐  │  Files     │   │  ┌───────▼───────┐  │            │
│  │  │  File System  │◄─┼────────────┼──►│  │  File System  │  │            │
│  │  │   (/app)      │  │  Volume    │   │  │   (/proxy)    │  │            │
│  │  └───────────────┘  │            │   │  └───────────────┘  │            │
│  │                     │            │   │                     │            │
│  │  Resources:         │            │   │  Resources:         │            │
│  │  CPU: 500m          │            │   │  CPU: 100m          │            │
│  │  Memory: 512Mi      │            │   │  Memory: 128Mi      │            │
│  └─────────────────────┘            │   └─────────────────────┘            │
│                                     │                                      │
│  ═══════════════════════════════════╪══════════════════════════════════    │
│                                     │                                      │
│                                     ▼                                      │
│                         ┌─────────────────────┐                            │
│                         │   Service Network   │                            │
│                         │  (Cluster/External) │                            │
│                         └─────────────────────┘                            │
│                                                                            │
└────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 DaemonSet Node-Level Injection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           KUBERNETES CLUSTER                                 │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                         CONTROL PLANE                                 │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │ API Server  │  │   etcd      │  │ Scheduler   │  │ Controller  │  │   │
│  │  │             │  │             │  │             │  │   Manager   │  │   │
│  │  └──────┬──────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  │         │                                                             │   │
│  │         │ DaemonSet Controller watches nodes                          │   │
│  │         │                                                             │   │
│  └─────────┼─────────────────────────────────────────────────────────────┘   │
│            │                                                                 │
│            │ creates/updates                                                  │
│            ▼                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                           WORKER NODES                                │   │
│  │                                                                        │   │
│  │  ┌─────────────────────────┐    ┌─────────────────────────┐           │   │
│  │  │       NODE-1            │    │       NODE-2            │           │   │
│  │  │  ┌───────────────────┐  │    │  ┌───────────────────┐  │           │   │
│  │  │  │   kubelet         │  │    │  │   kubelet         │  │           │   │
│  │  │  └───────────────────┘  │    │  └───────────────────┘  │           │   │
│  │  │         │               │    │         │               │           │   │
│  │  │  ┌──────┴──────┐        │    │  ┌──────┴──────┐        │           │   │
│  │  │  │  Pod A      │        │    │  │  Pod C      │        │           │   │
│  │  │  │  Pod B      │        │    │  │  Pod D      │        │           │   │
│  │  │  └─────────────┘        │    │  └─────────────┘        │           │   │
│  │  │         │               │    │         │               │           │   │
│  │  │  ┌──────┴──────────┐    │    │  ┌──────┴──────────┐    │           │   │
│  │  │  │ ┌─────────────┐ │    │    │  │ ┌─────────────┐ │    │           │   │
│  │  │  │ │  DaemonSet  │ │    │    │  │ │  DaemonSet  │ │    │           │   │
│  │  │  │ │    Pod      │ │    │    │  │ │    Pod      │ │    │           │   │
│  │  │  │ │             │ │    │    │  │ │             │ │    │           │   │
│  │  │  │ │ • Log       │ │    │    │  │ │ • Log       │ │    │           │   │
│  │  │  │ │   Collector │ │    │    │  │ │   Collector │ │    │           │   │
│  │  │  │ │ • Node      │ │    │    │  │ │ • Node      │ │    │           │   │
│  │  │  │ │   Exporter  │ │    │    │  │ │   Exporter  │ │    │           │   │
│  │  │  │ │ • Security  │ │    │    │  │ │ • Security  │ │    │           │   │
│  │  │  │ │   Agent     │ │    │    │  │ │   Agent     │ │    │           │   │
│  │  │  │ │ • CNI       │ │    │    │  │ │ • CNI       │ │    │           │   │
│  │  │  │ │   Plugin    │ │    │    │  │ │   Plugin    │ │    │           │   │
│  │  │  │ └─────────────┘ │    │    │  │ └─────────────┘ │    │           │   │
│  │  │  └─────────────────┘    │    │  └─────────────────┘    │           │   │
│  │  │                         │    │                         │           │   │
│  │  │  ┌─────────────────┐    │    │  ┌─────────────────┐    │           │   │
│  │  │  │  Node Resources │    │    │  │  Node Resources │    │           │   │
│  │  │  │  • /var/log     │    │    │  │  • /var/log     │    │           │   │
│  │  │  │  • /proc        │◄───┼────┼──┼──┤  • /proc        │    │           │   │
│  │  │  │  • /sys         │    │    │  │  • /sys         │    │           │   │
│  │  │  │  • /var/lib/...│    │    │  │  • /var/lib/...│    │           │   │
│  │  │  └─────────────────┘    │    │  └─────────────────┘    │           │   │
│  │  └─────────────────────────┘    └─────────────────────────┘           │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Process Injection Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROCESS INJECTION LIFECYCLE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DEPLOYMENT PHASE                    RUNTIME PHASE                    CLEANUP│
│                                                                             │
│  ┌─────────────┐                  ┌─────────────┐                  ┌───────┐ │
│  │             │  Injection       │             │  Monitoring      │       │ │
│  │   Define    │─────────────────►│   Active    │─────────────────►│ Stop  │ │
│  │   Pattern   │                  │   Running   │                  │       │ │
│  │             │                  │             │                  │       │ │
│  └─────────────┘                  └──────┬──────┘                  └───┬───┘ │
│         │                                │                             │     │
│         │                                │ Events                      │     │
│         ▼                                ▼                             ▼     │
│  ┌─────────────┐                  ┌─────────────┐                  ┌───────┐ │
│  │  Configure  │                  │   Health    │                  │ Signal│ │
│  │  Sidecar/   │                  │   Checks    │                  │ Prop. │ │
│  │  DaemonSet  │                  │             │                  │       │ │
│  └─────────────┘                  └──────┬──────┘                  └───────┘ │
│         │                                │                                   │
│         │                                │ Failure                           │
│         ▼                                ▼                                   │
│  ┌─────────────┐                  ┌─────────────┐                            │
│  │   Apply     │                  │  Recovery/  │                            │
│  │   to Pod/   │                  │  Restart    │                            │
│  │   Node      │                  │             │                            │
│  └─────────────┘                  └─────────────┘                            │
│         │                                                                   │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  Admission  │                                                            │
│  │  Controller │                                                            │
│  │  Mutates    │                                                            │
│  │  Pod Spec   │                                                            │
│  └─────────────┘                                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Sidecar Controller

```go
package injector

import (
 "context"
 "encoding/json"
 "fmt"
 "sync"
 "time"

 corev1 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/client-go/informers"
 "k8s.io/client-go/kubernetes"
 "k8s.io/client-go/tools/cache"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/metric"
 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
)

// InjectionMode defines how the sidecar is injected
type InjectionMode string

const (
 ModeManual    InjectionMode = "manual"
 ModeAutomatic InjectionMode = "automatic"
 ModeWebhook   InjectionMode = "webhook"
)

// SidecarConfig defines sidecar configuration
type SidecarConfig struct {
 Name            string            `json:"name"`
 Image           string            `json:"image"`
 ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`
 Resources       corev1.ResourceRequirements `json:"resources,omitempty"`
 Env             []corev1.EnvVar   `json:"env,omitempty"`
 VolumeMounts    []corev1.VolumeMount `json:"volumeMounts,omitempty"`
 Ports           []corev1.ContainerPort `json:"ports,omitempty"`
 SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`

 // Injection settings
 Mode              InjectionMode `json:"mode"`
 NamespaceSelector string        `json:"namespaceSelector,omitempty"`
 PodSelector       string        `json:"podSelector,omitempty"`

 // Lifecycle hooks
 PostStart *corev1.LifecycleHandler `json:"postStart,omitempty"`
 PreStop   *corev1.LifecycleHandler `json:"preStop,omitempty"`
}

// SidecarInjector manages automatic sidecar injection
type SidecarInjector struct {
 client        kubernetes.Interface
 config        SidecarConfig
 informer      cache.SharedIndexInformer
 logger        *zap.Logger
 tracer        trace.Tracer
 meter         metric.Meter

 // Metrics
 injectionsTotal   metric.Int64Counter
 injectionsFailed  metric.Int64Counter
 injectionDuration metric.Float64Histogram
 podsWatched       metric.Int64UpDownCounter

 mu sync.RWMutex
}

// NewSidecarInjector creates a new sidecar injector
func NewSidecarInjector(
 client kubernetes.Interface,
 config SidecarConfig,
 logger *zap.Logger,
 tracer trace.Tracer,
 meter metric.Meter,
) (*SidecarInjector, error) {
 si := &SidecarInjector{
  client: client,
  config: config,
  logger: logger,
  tracer: tracer,
  meter:  meter,
 }

 // Initialize metrics
 if meter != nil {
  var err error
  si.injectionsTotal, err = meter.Int64Counter(
   "sidecar_injections_total",
   metric.WithDescription("Total number of sidecar injections"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create injections counter: %w", err)
  }

  si.injectionsFailed, err = meter.Int64Counter(
   "sidecar_injections_failed_total",
   metric.WithDescription("Total number of failed sidecar injections"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create injection failures counter: %w", err)
  }

  si.injectionDuration, err = meter.Float64Histogram(
   "sidecar_injection_duration_seconds",
   metric.WithDescription("Duration of sidecar injection in seconds"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create injection duration histogram: %w", err)
  }

  si.podsWatched, err = meter.Int64UpDownCounter(
   "sidecar_pods_watched",
   metric.WithDescription("Number of pods being watched for injection"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create pods watched counter: %w", err)
  }
 }

 return si, nil
}

// Start begins watching pods for injection
func (si *SidecarInjector) Start(ctx context.Context) error {
 factory := informers.NewSharedInformerFactory(si.client, 30*time.Second)

 podInformer := factory.Core().V1().Pods().Informer()

 podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
  AddFunc: func(obj interface{}) {
   pod := obj.(*corev1.Pod)
   if si.shouldInject(pod) {
    si.logger.Debug("New pod detected for injection",
     zap.String("pod", pod.Name),
     zap.String("namespace", pod.Namespace))
    if err := si.inject(ctx, pod); err != nil {
     si.logger.Error("Failed to inject sidecar",
      zap.String("pod", pod.Name),
      zap.Error(err))
    }
   }
  },
  UpdateFunc: func(old, new interface{}) {
   // Handle updates if needed
  },
  DeleteFunc: func(obj interface{}) {
   if si.podsWatched != nil {
    si.podsWatched.Add(ctx, -1)
   }
  },
 })

 si.informer = podInformer
 factory.Start(ctx.Done())

 si.logger.Info("Sidecar injector started",
  zap.String("mode", string(si.config.Mode)),
  zap.String("sidecar", si.config.Name))

 return nil
}

// shouldInject determines if a pod should receive the sidecar
func (si *SidecarInjector) shouldInject(pod *corev1.Pod) bool {
 // Check if sidecar already exists
 for _, container := range pod.Spec.Containers {
  if container.Name == si.config.Name {
   return false
  }
 }

 // Check annotation
 if annotation, ok := pod.Annotations["sidecar.injector/enabled"]; ok {
  return annotation == "true"
 }

 // Check label selector if configured
 if si.config.Mode == ModeAutomatic {
  return true
 }

 return false
}

// inject adds the sidecar container to the pod
func (si *SidecarInjector) inject(ctx context.Context, pod *corev1.Pod) error {
 ctx, span := si.tracer.Start(ctx, "sidecar.inject",
  trace.WithAttributes(
   attribute.String("pod.name", pod.Name),
   attribute.String("pod.namespace", pod.Namespace),
   attribute.String("sidecar.name", si.config.Name),
  ))
 defer span.End()

 start := time.Now()

 // Create the sidecar container
 sidecar := corev1.Container{
  Name:            si.config.Name,
  Image:           si.config.Image,
  ImagePullPolicy: si.config.ImagePullPolicy,
  Resources:       si.config.Resources,
  Env:             si.config.Env,
  VolumeMounts:    si.config.VolumeMounts,
  Ports:           si.config.Ports,
  SecurityContext: si.config.SecurityContext,
 }

 if si.config.PostStart != nil {
  sidecar.Lifecycle = &corev1.Lifecycle{
   PostStart: si.config.PostStart,
  }
 }
 if si.config.PreStop != nil {
  if sidecar.Lifecycle == nil {
   sidecar.Lifecycle = &corev1.Lifecycle{}
  }
  sidecar.Lifecycle.PreStop = si.config.PreStop
 }

 // Update pod with sidecar
 pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

 // Add volume if needed (e.g., shared emptyDir)
 hasSharedVolume := false
 for _, vol := range pod.Spec.Volumes {
  if vol.Name == "shared-data" {
   hasSharedVolume = true
   break
  }
 }

 if !hasSharedVolume {
  pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
   Name: "shared-data",
   VolumeSource: corev1.VolumeSource{
    EmptyDir: &corev1.EmptyDirVolumeSource{},
   },
  })
 }

 // Apply the update
 _, err := si.client.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
 if err != nil {
  span.RecordError(err)
  if si.injectionsFailed != nil {
   si.injectionsFailed.Add(ctx, 1, metric.WithAttributes(
    attribute.String("reason", "update_failed"),
   ))
  }
  return fmt.Errorf("failed to update pod: %w", err)
 }

 duration := time.Since(start).Seconds()
 if si.injectionDuration != nil {
  si.injectionDuration.Record(ctx, duration)
 }
 if si.injectionsTotal != nil {
  si.injectionsTotal.Add(ctx, 1)
 }

 si.logger.Info("Sidecar injected successfully",
  zap.String("pod", pod.Name),
  zap.String("namespace", pod.Namespace),
  zap.String("sidecar", si.config.Name),
  zap.Float64("duration_seconds", duration))

 return nil
}
```

### 4.2 DaemonSet Controller

```go
package injector

import (
 "context"
 "fmt"
 "time"

 appsv1 "k8s.io/api/apps/v1"
 corev1 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/client-go/kubernetes"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/metric"
 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
)

// DaemonSetConfig defines DaemonSet configuration for node-level injection
type DaemonSetConfig struct {
 Name      string            `json:"name"`
 Namespace string            `json:"namespace"`
 Image     string            `json:"image"`
 Command   []string          `json:"command,omitempty"`
 Args      []string          `json:"args,omitempty"`

 // Node selection
 NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
 Tolerations    []corev1.Toleration `json:"tolerations,omitempty"`
 Affinity       *corev1.Affinity `json:"affinity,omitempty"`

 // Resource management
 Resources      corev1.ResourceRequirements `json:"resources,omitempty"`

 // Security
 Privileged     bool              `json:"privileged"`
 RunAsUser      *int64            `json:"runAsUser,omitempty"`
 RunAsGroup     *int64            `json:"runAsGroup,omitempty"`
 Capabilities   []corev1.Capability `json:"capabilities,omitempty"`

 // Volumes for host access
 HostPathMounts []HostPathMount   `json:"hostPathMounts,omitempty"`

 // Health checks
 LivenessProbe  *corev1.Probe     `json:"livenessProbe,omitempty"`
 ReadinessProbe *corev1.Probe     `json:"readinessProbe,omitempty"`

 // Update strategy
 UpdateStrategy appsv1.DaemonSetUpdateStrategy `json:"updateStrategy,omitempty"`
}

// HostPathMount defines a host path volume mount
type HostPathMount struct {
 Name      string               `json:"name"`
 HostPath  string               `json:"hostPath"`
 MountPath string               `json:"mountPath"`
 ReadOnly  bool                 `json:"readOnly"`
 Type      corev1.HostPathType  `json:"type,omitempty"`
}

// DaemonSetController manages DaemonSet lifecycle
type DaemonSetController struct {
 client    kubernetes.Interface
 config    DaemonSetConfig
 logger    *zap.Logger
 tracer    trace.Tracer
 meter     metric.Meter

 // Metrics
 daemonsetsCreated   metric.Int64Counter
 nodesCovered        metric.Int64Gauge
 podsReady           metric.Int64Gauge
}

// NewDaemonSetController creates a new DaemonSet controller
func NewDaemonSetController(
 client kubernetes.Interface,
 config DaemonSetConfig,
 logger *zap.Logger,
 tracer trace.Tracer,
 meter metric.Meter,
) (*DaemonSetController, error) {
 dc := &DaemonSetController{
  client: client,
  config: config,
  logger: logger,
  tracer: tracer,
  meter:  meter,
 }

 if meter != nil {
  var err error
  dc.daemonsetsCreated, err = meter.Int64Counter(
   "daemonset_created_total",
   metric.WithDescription("Total number of DaemonSets created"),
  )
  if err != nil {
   return nil, err
  }

  dc.nodesCovered, err = meter.Int64Gauge(
   "daemonset_nodes_covered",
   metric.WithDescription("Number of nodes covered by DaemonSet"),
  )
  if err != nil {
   return nil, err
  }

  dc.podsReady, err = meter.Int64Gauge(
   "daemonset_pods_ready",
   metric.WithDescription("Number of ready DaemonSet pods"),
  )
  if err != nil {
   return nil, err
  }
 }

 return dc, nil
}

// Deploy creates or updates the DaemonSet
func (dc *DaemonSetController) Deploy(ctx context.Context) error {
 ctx, span := dc.tracer.Start(ctx, "daemonset.deploy",
  trace.WithAttributes(
   attribute.String("daemonset.name", dc.config.Name),
   attribute.String("daemonset.namespace", dc.config.Namespace),
  ))
 defer span.End()

 // Build the DaemonSet spec
 ds := dc.buildDaemonSet()

 // Check if DaemonSet exists
 existing, err := dc.client.AppsV1().DaemonSets(dc.config.Namespace).Get(ctx, dc.config.Name, metav1.GetOptions{})
 if err == nil {
  // Update existing
  ds.ResourceVersion = existing.ResourceVersion
  _, err = dc.client.AppsV1().DaemonSets(dc.config.Namespace).Update(ctx, ds, metav1.UpdateOptions{})
  if err != nil {
   span.RecordError(err)
   return fmt.Errorf("failed to update DaemonSet: %w", err)
  }
  dc.logger.Info("DaemonSet updated",
   zap.String("name", dc.config.Name),
   zap.String("namespace", dc.config.Namespace))
 } else {
  // Create new
  _, err = dc.client.AppsV1().DaemonSets(dc.config.Namespace).Create(ctx, ds, metav1.CreateOptions{})
  if err != nil {
   span.RecordError(err)
   return fmt.Errorf("failed to create DaemonSet: %w", err)
  }
  dc.logger.Info("DaemonSet created",
   zap.String("name", dc.config.Name),
   zap.String("namespace", dc.config.Namespace))

  if dc.daemonsetsCreated != nil {
   dc.daemonsetsCreated.Add(ctx, 1)
  }
 }

 return nil
}

func (dc *DaemonSetController) buildDaemonSet() *appsv1.DaemonSet {
 // Build volumes and volume mounts
 volumes := make([]corev1.Volume, 0, len(dc.config.HostPathMounts))
 volumeMounts := make([]corev1.VolumeMount, 0, len(dc.config.HostPathMounts))

 for _, mount := range dc.config.HostPathMounts {
  volName := fmt.Sprintf("host-%s", mount.Name)
  volumes = append(volumes, corev1.Volume{
   Name: volName,
   VolumeSource: corev1.VolumeSource{
    HostPath: &corev1.HostPathVolumeSource{
     Path: mount.HostPath,
     Type: &mount.Type,
    },
   },
  })
  volumeMounts = append(volumeMounts, corev1.VolumeMount{
   Name:      volName,
   MountPath: mount.MountPath,
   ReadOnly:  mount.ReadOnly,
  })
 }

 // Build security context
 securityContext := &corev1.SecurityContext{}
 if dc.config.Privileged {
  securityContext.Privileged = &dc.config.Privileged
 }
 if dc.config.RunAsUser != nil {
  securityContext.RunAsUser = dc.config.RunAsUser
 }
 if dc.config.RunAsGroup != nil {
  securityContext.RunAsGroup = dc.config.RunAsGroup
 }
 if len(dc.config.Capabilities) > 0 {
  securityContext.Capabilities = &corev1.Capabilities{
   Add: dc.config.Capabilities,
  }
 }

 return &appsv1.DaemonSet{
  ObjectMeta: metav1.ObjectMeta{
   Name:      dc.config.Name,
   Namespace: dc.config.Namespace,
   Labels: map[string]string{
    "app.kubernetes.io/name":      dc.config.Name,
    "app.kubernetes.io/managed-by": "process-injector",
   },
  },
  Spec: appsv1.DaemonSetSpec{
   Selector: &metav1.LabelSelector{
    MatchLabels: map[string]string{
     "app": dc.config.Name,
    },
   },
   UpdateStrategy: dc.config.UpdateStrategy,
   Template: corev1.PodTemplateSpec{
    ObjectMeta: metav1.ObjectMeta{
     Labels: map[string]string{
      "app": dc.config.Name,
     },
    },
    Spec: corev1.PodSpec{
     NodeSelector:   dc.config.NodeSelector,
     Tolerations:    dc.config.Tolerations,
     Affinity:       dc.config.Affinity,
     HostNetwork:    false,
     HostPID:        dc.config.Privileged, // Access host processes if privileged
     Containers: []corev1.Container{
      {
       Name:            dc.config.Name,
       Image:           dc.config.Image,
       Command:         dc.config.Command,
       Args:            dc.config.Args,
       Resources:       dc.config.Resources,
       SecurityContext: securityContext,
       VolumeMounts:    volumeMounts,
       LivenessProbe:   dc.config.LivenessProbe,
       ReadinessProbe:  dc.config.ReadinessProbe,
      },
     },
     Volumes: volumes,
    },
   },
  },
 }
}

// WatchStatus monitors DaemonSet rollout status
func (dc *DaemonSetController) WatchStatus(ctx context.Context) error {
 ticker := time.NewTicker(10 * time.Second)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return ctx.Err()
  case <-ticker.C:
   ds, err := dc.client.AppsV1().DaemonSets(dc.config.Namespace).Get(ctx, dc.config.Name, metav1.GetOptions{})
   if err != nil {
    dc.logger.Error("Failed to get DaemonSet status", zap.Error(err))
    continue
   }

   status := ds.Status

   if dc.nodesCovered != nil {
    dc.nodesCovered.Record(ctx, int64(status.NumberAvailable))
   }
   if dc.podsReady != nil {
    dc.podsReady.Record(ctx, int64(status.NumberReady))
   }

   dc.logger.Debug("DaemonSet status",
    zap.Int32("desired", status.DesiredNumberScheduled),
    zap.Int32("current", status.CurrentNumberScheduled),
    zap.Int32("ready", status.NumberReady),
    zap.Int32("available", status.NumberAvailable),
    zap.Int32("misscheduled", status.NumberMisscheduled))
  }
 }
}
```

### 4.3 Ephemeral Container Injector

```go
package injector

import (
 "context"
 "fmt"

 corev1 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/client-go/kubernetes"
 "k8s.io/utils/pointer"
)

// EphemeralContainerConfig defines ephemeral container configuration
type EphemeralContainerConfig struct {
 Name            string   `json:"name"`
 Image           string   `json:"image"`
 Command         []string `json:"command,omitempty"`
 TargetContainer string   `json:"targetContainer,omitempty"`

 // Namespace sharing
 ShareProcessNamespace bool `json:"shareProcessNamespace"`

 // Security
 Privileged bool     `json:"privileged"`
 Sysctls    []string `json:"sysctls,omitempty"`
}

// EphemeralInjector manages ephemeral container injection
type EphemeralInjector struct {
 client kubernetes.Interface
}

// NewEphemeralInjector creates a new ephemeral container injector
func NewEphemeralInjector(client kubernetes.Interface) *EphemeralInjector {
 return &EphemeralInjector{
  client: client,
 }
}

// InjectDebugContainer adds an ephemeral debug container to a running pod
func (ei *EphemeralInjector) InjectDebugContainer(
 ctx context.Context,
 namespace, podName string,
 config EphemeralContainerConfig,
) error {
 // Get the target pod
 pod, err := ei.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
 if err != nil {
  return fmt.Errorf("failed to get pod: %w", err)
 }

 // Build ephemeral container
 ec := &corev1.EphemeralContainer{
  EphemeralContainerCommon: corev1.EphemeralContainerCommon{
   Name:    config.Name,
   Image:   config.Image,
   Command: config.Command,
   SecurityContext: &corev1.SecurityContext{
    Privileged: &config.Privileged,
   },
   Stdin:     true,
   TTY:       true,
  },
  TargetContainerName: config.TargetContainer,
 }

 // Update pod with ephemeral container
 pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, *ec)

 // Apply update
 _, err = ei.client.CoreV1().Pods(namespace).UpdateEphemeralContainers(
  ctx,
  podName,
  pod,
  metav1.UpdateOptions{},
 )
 if err != nil {
  return fmt.Errorf("failed to inject ephemeral container: %w", err)
 }

 return nil
}

// InjectProfiler injects a profiling container
func (ei *EphemeralInjector) InjectProfiler(ctx context.Context, namespace, podName, targetContainer string) error {
 config := EphemeralContainerConfig{
  Name:            "profiler",
  Image:           "gcr.io/google-containers/profiler:latest",
  Command:         []string{"/bin/sh"},
  TargetContainer: targetContainer,
  Privileged:      true,
 }
 return ei.InjectDebugContainer(ctx, namespace, podName, config)
}

// InjectNetworkDebugger injects a network debugging container
func (ei *EphemeralInjector) InjectNetworkDebugger(ctx context.Context, namespace, podName string) error {
 config := EphemeralContainerConfig{
  Name:    "network-debug",
  Image:   "nicolaka/netshoot:latest",
  Command: []string{"/bin/bash"},
  Privileged: true,
 }
 return ei.InjectDebugContainer(ctx, namespace, podName, config)
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      INJECTION FAILURE SCENARIOS                             │
├───────────────────────────────┬───────────────────┬─────────────────────────┤
│         Scenario              │     Detection     │      Mitigation         │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Admission Webhook Timeout     │ Webhook latency   │ Fail open policy +      │
│                               │ > 30s             │ Alerting                │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Sidecar Image Pull Failure    │ ImagePullBackOff  │ Image mirroring +       │
│                               │ status            │ Fallback registry       │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Resource Exhaustion           │ OOMKilled/        │ Resource quotas +       │
│                               │ CPU throttling    │ Priority classes        │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ PID Namespace Leak            │ High PID usage    │ PID limits +            │
│                               │                   │ Namespace cleanup       │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Network Policy Conflict       │ Connection        │ Policy coordination +   │
│                               │ refused           │ CNI compatibility       │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ DaemonSet Rollout Stuck       │ Pods not ready    │ MaxUnavailable surge +  │
│                               │                   │ Progress deadline       │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Security Context Violation    │ Seccomp/AppArmor  │ Security profiles +     │
│                               │ denials           │ Audit logging           │
└───────────────────────────────┴───────────────────┴─────────────────────────┘
```

### 5.2 Mitigation Strategies

```go
// InjectionPolicy defines failure handling policies
type InjectionPolicy struct {
 FailurePolicy   FailurePolicy `json:"failurePolicy"`
 TimeoutSeconds  int32         `json:"timeoutSeconds"`
 ReinvocationPolicy string     `json:"reinvocationPolicy"`
}

type FailurePolicy string

const (
 FailFailurePolicy  FailurePolicy = "Fail"
 IgnoreFailurePolicy FailurePolicy = "Ignore"
)

// HealthChecker monitors injection health
type HealthChecker struct {
 client kubernetes.Interface
 logger *zap.Logger
}

// Check performs health checks on injected containers
func (hc *HealthChecker) Check(ctx context.Context, pod *corev1.Pod) error {
 // Check sidecar health
 for _, container := range pod.Spec.Containers {
  if container.Name == "sidecar-proxy" {
   // Check if container is running
   for _, status := range pod.Status.ContainerStatuses {
    if status.Name == container.Name {
     if !status.Ready {
      return fmt.Errorf("sidecar container not ready: %s", status.State.Waiting.Reason)
     }
    }
   }
  }
 }
 return nil
}
```

---

## 6. Semantic Trade-off Analysis

### 6.1 Injection Pattern Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    INJECTION PATTERN TRADE-OFFS                              │
├─────────────────────┬───────────────┬───────────────┬───────────────────────┤
│     Dimension       │    Sidecar    │   DaemonSet   │  Ephemeral Container  │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Co-location         │ Pod-level     │ Node-level    │ Process-level         │
│ Lifecycle Coupling  │ Tight (same   │ Loose (node   │ Very tight (shares    │
│                     │ pod lifecycle)│ lifecycle)    │ PID namespace)        │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Resource Overhead   │ Per-pod       │ Per-node      │ Temporary             │
│                     │ (high)        │ (amortized)   │ (minimal)             │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Security Isolation  │ Container     │ Node-level    │ Minimal (shared PID)  │
│                     │ boundary      │ privilege     │                       │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Deployment Speed    │ Pod startup   │ Node join     │ Instant               │
│                     │ latency       │ latency       │                       │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Use Cases           │ Service mesh, │ Node          │ Debugging,            │
│                     │ logging,      │ monitoring,   │ profiling,            │
│                     │ monitoring    │ security      │ troubleshooting       │
└─────────────────────┴───────────────┴───────────────┴───────────────────────┘
```

### 6.2 Security vs Functionality Trade-offs

| Aspect | High Security | High Functionality | Balance |
|--------|---------------|-------------------|---------|
| **Privileges** | Non-root, read-only FS | Privileged, host access | Capabilities-based |
| **Network** | Deny all, explicit allow | Full access | Service mesh proxy |
| **Storage** | EmptyDir only | HostPath mounts | VolumeClaim templates |
| **Observability** | Audit logs only | Full syscall tracing | eBPF-based |

---

## 7. References

1. Burns, B., Beda, J., & Hightower, K. (2019). *Kubernetes: Up and Running*. O'Reilly Media.
2. Ibryam, B., & Huß, R. (2019). *Kubernetes Patterns*. O'Reilly Media.
3. Beyer, B., Jones, C., Petoff, J., & Murphy, N. (2016). *Site Reliability Engineering*. O'Reilly Media.
4. Barroso, L. A., Clidaras, J., & Hölzle, U. (2018). *The Datacenter as a Computer*. Morgan & Claypool.
5. Kubernetes Documentation. (2024). Sidecar Containers. kubernetes.io.
6. CNCF. (2024). Service Mesh Interface Specification.
