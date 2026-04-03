# EC-077: Multi-Container Patterns - Native Sidecars, Init Containers, and Complementary Patterns

> **维度**: Engineering Cloud-Native
> **级别**: S (17+ KB)
> **标签**: #kubernetes #multi-container #sidecar #init-container #ambassador #adapter #configuration
> **权威来源**:
>
> - [Kubernetes Multi-Container Patterns](https://kubernetes.io/docs/concepts/workloads/pods/#workload-resources-for-managing-pods) - Official Documentation
> - [KEP-753: Sidecar Containers](https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/753-sidecar-containers) - Kubernetes Enhancement
> - [CNCF Cloud Native Patterns](https://www.cncf.io/phippy/) - Cloud Native Patterns
> - [Google Cloud Blog: Multi-Container Patterns](https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-organizing-containers-with-pods)

---

## 1. Introduction to Multi-Container Pods

Multi-container pods are a fundamental Kubernetes pattern that enables co-located, co-scheduled containers sharing the same network namespace, storage volumes, and lifecycle context. Kubernetes 1.34 significantly enhances this paradigm with native sidecar support and improved container lifecycle management.

### 1.1 Pod Container Types Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Kubernetes Pod Container Types                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  │                                                               │   │
│  │  ├─▶ Init Containers (Sequential, must complete)                │   │
│  │  │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │   │
│  │  │   │  init-1     │───▶│  init-2     │───▶│  init-n     │      │   │
│  │  │   │             │    │             │    │             │      │   │
│  │  │   │ restartPolicy│    │ restartPolicy│    │ restartPolicy│     │   │
│  │  │   │ = Never     │    │ = Always    │    │ = Never     │      │   │
│  │  │   │ (legacy)    │    │ (native     │    │ (legacy)    │      │   │
│  │  │   │             │    │  sidecar)   │    │             │      │   │
│  │  │   └─────────────┘    └─────────────┘    └─────────────┘      │   │
│  │  │                                                               │   │
│  │  └─▶ Application Containers (Parallel, long-running)            │   │
│  │      ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │   │
│  │      │  Sidecar    │    │  Main App   │    │  Sidecar    │      │   │
│  │      │  (Native)   │◄──►│             │◄──►│  (Native)   │      │   │
│  │      │             │    │             │    │             │      │   │
│  │      │ restartPolicy│    │ restartPolicy│    │ restartPolicy│     │   │
│  │      │ = Always    │    │ = Always    │    │ = OnFailure │      │   │
│  │      │             │    │             │    │             │      │   │
│  │      └─────────────┘    └─────────────┘    └─────────────┘      │   │
│  │                                                                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Execution Order:                                                       │
│  1. All Init Containers (including native sidecars) run sequentially    │
│  2. Init Containers with restartPolicy: Always block until Ready        │
│  3. Application Containers start only after all init containers ready   │
│  4. Native sidecars terminate automatically when main containers exit   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Communication Patterns Between Containers

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Inter-Container Communication Patterns                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pattern 1: Shared Process Namespace (Sidecar Debugging)                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod (shareProcessNamespace: true)                               │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │  Main App    │◄───────►│  Debug Side  │                     │   │
│  │  │  (PID 1)     │ signals │  car (PID 2) │                     │   │
│  │  │              │         │  (kubectl    │                     │   │
│  │  │              │         │   debug)     │                     │   │
│  │  └──────────────┘         └──────────────┘                     │   │
│  │         │                                                        │   │
│  │         ▼                                                        │   │
│  │  Shared PID namespace - sidecar can debug main process          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 2: Shared Volumes (File-based Communication)                   │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐  /shared/logs  ┌──────────────┐               │   │
│  │  │  App         │◄──────────────►│  Log Shipper │               │   │
│  │  │  (writes     │  (emptyDir)    │  (reads and  │               │   │
│  │  │   logs)      │                │   forwards)  │               │   │
│  │  └──────────────┘                └──────────────┘               │   │
│  │         │                                                        │   │
│  │         ▼                                                        │   │
│  │  emptyDir volume mounted by both containers                      │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 3: Localhost Networking (Network Communication)                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod (shared network namespace)                                  │   │
│  │  ┌──────────────┐  localhost:8080 ┌──────────────┐              │   │
│  │  │  Main App    │◄───────────────►│  Envoy       │              │   │
│  │  │  (:8080)     │                 │  (:8000)     │              │   │
│  │  │              │                 │  (mTLS,      │              │   │
│  │  │              │                 │   routing)   │              │   │
│  │  └──────────────┘                 └──────┬───────┘              │   │
│  │                                          │                       │   │
│  │                              External Network                    │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 4: Unix Domain Sockets (High-performance IPC)                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐  /tmp/app.sock  ┌──────────────┐              │   │
│  │  │  gRPC        │◄───────────────►│  gRPC        │              │   │
│  │  │  Server      │  (UDS)          │  Client      │              │   │
│  │  │  (:0)        │                 │  (sidecar)   │              │   │
│  │  └──────────────┘                 └──────────────┘              │   │
│  │         │                                                        │   │
│  │         ▼                                                        │   │
│  │  30-50% lower latency than TCP localhost                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Native Sidecar Pattern (Kubernetes 1.33+ Stable)

Native sidecar containers, promoted to stable in Kubernetes 1.33 with full 1.34 support, provide lifecycle management guarantees that were previously impossible with traditional sidecar implementations.

### 2.1 Native Sidecar Lifecycle Guarantees

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Native Sidecar Lifecycle (restartPolicy: Always)            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Phase 1: Initialization                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Time ─────────────────────────────────────────────────────▶    │   │
│  │                                                                  │   │
│  │  Init Containers:                                                │   │
│  │  ┌─────┐   ┌──────────────┐   ┌──────────────┐                  │   │
│  │  │ init│──▶│   Sidecar    │──▶│   Sidecar    │                  │   │
│  │  │ -db │   │   (proxy)    │   │   (monitor)  │                  │   │
│  │  └──┬──┘   │              │   │              │                  │   │
│  │     │      │ restartPolicy│   │ restartPolicy│                  │   │
│  │     │      │ = Always     │   │ = Always     │                  │   │
│  │     │      └──────┬───────┘   └──────┬───────┘                  │   │
│  │     │             │                  │                          │   │
│  │     ▼             ▼                  ▼                          │   │
│  │  Complete    ┌────┴────┐        ┌────┴────┐                      │   │
│  │              │Startup  │        │Startup  │                      │   │
│  │              │Probe    │        │Probe    │                      │   │
│  │              │Waiting  │        │Waiting  │                      │   │
│  │              └────┬────┘        └────┬────┘                      │   │
│  │                   │                  │                          │   │
│  │                   └────────┬─────────┘                          │   │
│  │                            │                                    │   │
│  │                            ▼ (Both Ready)                        │   │
│  │  Phase 2: Application Start ────────────────────────────────▶    │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐    │   │
│  │  │  Main App Container Starts (Guaranteed: AFTER sidecars)  │    │   │
│  │  │  ┌──────────┐  ┌──────────┐  ┌──────────┐              │    │   │
│  │  │  │ Sidecar  │  │ Main App │  │ Sidecar  │              │    │   │
│  │  │  │ (proxy)  │  │          │  │(monitor) │              │    │   │
│  │  │  │ Running  │  │ Running  │  │ Running  │              │    │   │
│  │  │  └──────────┘  └──────────┘  └──────────┘              │    │   │
│  │  └─────────────────────────────────────────────────────────┘    │   │
│  │                                                                  │   │
│  │  Phase 3: Termination (Jobs Only) ─────────────────────────▶    │   │
│  │                                                                  │   │
│  │  Main App Exits ──▶ Sidecars Receive SIGTERM ──▶ Pod Complete   │   │
│  │       │                    │                                     │   │
│  │       ▼                    ▼                                     │   │
│  │  ┌──────────┐        ┌──────────┐                               │   │
│  │  │ Main App │        │ Sidecars │                               │   │
│  │  │  Exited  │   ──▶  │ Graceful │                               │   │
│  │  │ Code: 0  │        │ Shutdown │                               │   │
│  │  └──────────┘        └──────────┘                               │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Key Guarantees:                                                        │
│  ✓ Sidecars start before main application                             │
│  ✓ Main app waits for sidecar readiness probes                        │
│  ✓ Sidecars can restart independently without affecting main app      │
│  ✓ For Jobs, sidecars terminate automatically when main container exits│
│  ✓ No race conditions during startup or shutdown                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Production-Ready Native Sidecar Implementation

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: production-service-with-sidecars
  labels:
    app: payment-api
    version: v2.3.1
    sidecar.istio.io/inject: "true"
spec:
  # Enable process sharing for debugging capabilities
  shareProcessNamespace: true

  # Traditional init containers (run first, must complete)
  initContainers:
    # 1. Database migration (must complete before anything else)
    - name: db-migrator
      image: myregistry/db-migrator:v2.3.1
      restartPolicy: Never  # Traditional init behavior
      command:
        - "/migrate"
        - "--target-version=v2.3.1"
        - "--strict"
      env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
      resources:
        limits:
          cpu: "1"
          memory: 512Mi
        requests:
          cpu: 100m
          memory: 128Mi
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        runAsNonRoot: true
        runAsUser: 1000

    # 2. Istio init container for iptables setup
    - name: istio-init
      image: istio/proxyv2:1.22.0
      restartPolicy: Never  # One-time setup
      args:
        - istio-iptables
        - -p
        - "15001"
        - -z
        - "15006"
        - -u
        - "1337"
        - -m
        - REDIRECT
        - -i
        - '*'
        - -x
        - ""
        - -b
        - "8080,9090"
        - -d
        - "15090,15021,15020"
      securityContext:
        capabilities:
          add:
            - NET_ADMIN
            - NET_RAW
        privileged: true
        runAsGroup: 0
        runAsUser: 0

    # 3. NATIVE SIDECAR: Istio proxy (restartPolicy: Always)
    - name: istio-proxy
      image: istio/proxyv2:1.22.0
      restartPolicy: Always  # ← Native sidecar marker
      args:
        - proxy
        - sidecar
        - --domain
        - $(POD_NAMESPACE).svc.cluster.local
        - --proxyLogLevel=warning
        - --proxyComponentLogLevel=misc:error
        - --log_output_level=default:info
      ports:
        - containerPort: 15090
          name: http-envoy-prom
          protocol: TCP
        - containerPort: 15021
          name: health
          protocol: TCP
        - containerPort: 15020
          name: mesh-telemetry
          protocol: TCP
      env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: ISTIO_META_CLUSTER_ID
          value: "kubernetes"
      resources:
        limits:
          cpu: "2000m"
          memory: 1Gi
        requests:
          cpu: 100m
          memory: 128Mi
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
        readOnlyRootFilesystem: true
        runAsGroup: 1337
        runAsUser: 1337
      # CRITICAL: Startup probe blocks main container until ready
      startupProbe:
        httpGet:
          path: /healthz/ready
          port: 15021
        initialDelaySeconds: 1
        periodSeconds: 1
        failureThreshold: 60  # Wait up to 60 seconds
      livenessProbe:
        httpGet:
          path: /healthz/live
          port: 15021
        periodSeconds: 10
        failureThreshold: 3
      readinessProbe:
        httpGet:
          path: /healthz/ready
          port: 15021
        periodSeconds: 5
        failureThreshold: 3
      volumeMounts:
        - name: istio-envoy
          mountPath: /etc/istio/proxy
        - name: istio-certs
          mountPath: /etc/certs
          readOnly: true

    # 4. NATIVE SIDECAR: Prometheus metrics exporter
    - name: metrics-exporter
      image: prom/node-exporter:v1.8.0
      restartPolicy: Always  # ← Native sidecar marker
      args:
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        - --path.rootfs=/host/root
        - --collector.filesystem.ignored-mount-points='^/(sys|proc|dev|host|etc)($$|/)'
      ports:
        - containerPort: 9100
          name: metrics
          protocol: TCP
      resources:
        limits:
          cpu: 250m
          memory: 180Mi
        requests:
          cpu: 50m
          memory: 64Mi
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
      startupProbe:
        httpGet:
          path: /
          port: 9100
        failureThreshold: 10
        periodSeconds: 1
      volumeMounts:
        - name: proc
          mountPath: /host/proc
          readOnly: true
        - name: sys
          mountPath: /host/sys
          readOnly: true

    # 5. NATIVE SIDECAR: Log shipper (Fluent Bit)
    - name: log-shipper
      image: fluent/fluent-bit:3.0
      restartPolicy: Always  # ← Native sidecar marker
      command:
        - /fluent-bit/bin/fluent-bit
        - -c
        - /fluent-bit/etc/fluent-bit.conf
      env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch.logging.svc.cluster.local"
        - name: FLUENT_ELASTICSEARCH_PORT
          value: "9200"
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 100m
          memory: 64Mi
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
      startupProbe:
        exec:
          command:
            - /bin/sh
            - -c
            - "test -S /fluent-bit/tmp/fluent-bit.sock"
        failureThreshold: 20
        periodSeconds: 1
      volumeMounts:
        - name: app-logs
          mountPath: /var/log/app
        - name: fluent-bit-config
          mountPath: /fluent-bit/etc

  # MAIN CONTAINER: Starts AFTER all native sidecars are Ready
  containers:
    - name: payment-api
      image: myregistry/payment-api:v2.3.1
      ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        - containerPort: 9090
          name: grpc
          protocol: TCP
      env:
        - name: SERVER_PORT
          value: "8080"
        - name: GRPC_PORT
          value: "9090"
        - name: LOG_LEVEL
          value: "info"
        - name: LOG_FORMAT
          value: "json"
        - name: SERVICE_MESH_ENABLED
          value: "true"
        # Application knows sidecars are ready when it starts
        - name: ISTIO_SIDECAR_STATUS
          value: "ready"
      resources:
        limits:
          cpu: "2"
          memory: 2Gi
        requests:
          cpu: 500m
          memory: 512Mi
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
        readOnlyRootFilesystem: true
        runAsNonRoot: true
        runAsUser: 1000
      livenessProbe:
        httpGet:
          path: /healthz/live
          port: 8080
        initialDelaySeconds: 10
        periodSeconds: 10
        failureThreshold: 3
      readinessProbe:
        httpGet:
          path: /healthz/ready
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 5
        failureThreshold: 3
      volumeMounts:
        - name: app-logs
          mountPath: /var/log/app
        - name: tmp
          mountPath: /tmp

  volumes:
    - name: istio-envoy
      emptyDir: {}
    - name: istio-certs
      secret:
        secretName: istio.default
        optional: true
    - name: app-logs
      emptyDir:
        sizeLimit: 1Gi
    - name: tmp
      emptyDir:
        sizeLimit: 100Mi
    - name: proc
      hostPath:
        path: /proc
        type: Directory
    - name: sys
      hostPath:
        path: /sys
        type: Directory
    - name: fluent-bit-config
      configMap:
        name: fluent-bit-config
```

### 2.3 Sidecar Pattern Decision Matrix

| Pattern | Use Case | Implementation | Kubernetes Version |
|---------|----------|----------------|-------------------|
| **Service Mesh** | mTLS, traffic management | Istio/Linkerd as native sidecar | 1.28+ (1.33 stable) |
| **Observability** | Metrics, logging, tracing | Prometheus exporter, Fluent Bit | 1.28+ (1.33 stable) |
| **Configuration** | Dynamic config reloading | Config reloader sidecar | 1.28+ (1.33 stable) |
| **Security** | Secrets injection, TLS | Vault agent, cert-manager | 1.28+ (1.33 stable) |
| **Data Synchronization** | File sync, caching | rsync, custom sync agents | 1.28+ (1.33 stable) |
| **Legacy Adapter** | Protocol translation | REST-to-gRPC adapter | 1.28+ (1.33 stable) |

---

## 3. Init Container Patterns

Init containers perform setup tasks before the main application starts. With Kubernetes 1.34, init containers can optionally use `restartPolicy: Always` to become native sidecars.

### 3.1 Init Container Execution Order

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Init Container Execution Flow                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pod Scheduled ──▶ Pod Initialized ──▶ Containers Running ──▶ Running  │
│       │                  │                    │                    │     │
│       ▼                  ▼                    ▼                    ▼     │
│  ┌─────────┐    ┌──────────────────┐   ┌──────────────┐   ┌──────────┐ │
│  │ kubelet │    │ Init Containers  │   │   PostStart  │   │  Main    │ │
│  │ creates │───▶│  (sequential)    │──▶│    Hooks     │──▶│  Loop    │ │
│  │   Pod   │    │                  │   │              │   │          │ │
│  └─────────┘    ├──────────────────┤   └──────────────┘   └──────────┘ │
│                 │ 1. Network Setup   │                                  │
│                 │    (CNI wait)      │                                  │
│                 │    ↓ Success       │                                  │
│                 │ 2. Permissions     │                                  │
│                 │    (fsGroup setup) │                                  │
│                 │    ↓ Success       │                                  │
│                 │ 3. Data Prep       │                                  │
│                 │    (migration)     │                                  │
│                 │    ↓ Success       │                                  │
│                 │ 4. Security Setup  │                                  │
│                 │    (certificates)  │                                  │
│                 │    ↓ Success       │                                  │
│                 │ 5. Native Sidecars │                                  │
│                 │    (restartPolicy: │                                  │
│                 │     Always - wait  │                                  │
│                 │     for Ready)     │                                  │
│                 └──────────────────┘                                  │
│                                                                         │
│  Failure Handling:                                                      │
│  • Any init container failure → Pod fails                             │
│  • restartPolicy: Never (default) → No retry                          │
│  • restartPolicy: Always (1.34+) → Container restarts, blocks start   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Common Init Container Patterns

```yaml
# Pattern 1: Network Initialization
apiVersion: v1
kind: Pod
metadata:
  name: network-dependent-app
spec:
  initContainers:
    - name: wait-for-database
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          until nc -z postgres.database.svc.cluster.local 5432; do
            echo "Waiting for database..."
            sleep 2
          done
          echo "Database is ready!"

    - name: wait-for-cache
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          until nc -z redis.cache.svc.cluster.local 6379; do
            echo "Waiting for cache..."
            sleep 1
          done
          echo "Cache is ready!"

    - name: wait-for-message-queue
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          until nc -z kafka.kafka.svc.cluster.local 9092; do
            echo "Waiting for Kafka..."
            sleep 2
          done
          echo "Kafka is ready!"
  containers:
    - name: app
      image: myapp:latest

---
# Pattern 2: Data Initialization and Migration
apiVersion: v1
kind: Pod
metadata:
  name: data-migration-app
spec:
  initContainers:
    - name: volume-permissions
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          chown -R 1000:1000 /data
          chmod 755 /data
      securityContext:
        runAsUser: 0
      volumeMounts:
        - name: data-volume
          mountPath: /data

    - name: schema-migration
      image: migrate/migrate:v4.17
      command:
        - migrate
        - -path=/migrations
        - -database=postgresql://$(DB_USER):$(DB_PASS)@postgres:5432/app?sslmode=disable
        - up
      env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: DB_PASS
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
      volumeMounts:
        - name: migrations
          mountPath: /migrations

    - name: seed-data
      image: myapp/seed-data:v1.0
      command:
        - /seed
        - --if-empty
        - /seed-data
      volumeMounts:
        - name: seed-data-volume
          mountPath: /seed-data

  containers:
    - name: app
      image: myapp:latest
      volumeMounts:
        - name: data-volume
          mountPath: /data

  volumes:
    - name: data-volume
      persistentVolumeClaim:
        claimName: app-data
    - name: migrations
      configMap:
        name: db-migrations
    - name: seed-data-volume
      configMap:
        name: seed-data

---
# Pattern 3: Security and Certificate Setup
apiVersion: v1
kind: Pod
metadata:
  name: tls-enabled-app
spec:
  initContainers:
    - name: vault-agent-init
      image: hashicorp/vault:1.16
      command:
        - vault
        - agent
        - -config=/etc/vault/config.hcl
      env:
        - name: VAULT_ADDR
          value: "https://vault.vault-system.svc.cluster.local:8200"
        - name: VAULT_ROLE
          value: "app-role"
      volumeMounts:
        - name: vault-config
          mountPath: /etc/vault
        - name: vault-agent-sa
          mountPath: /var/run/secrets/kubernetes.io/serviceaccount
        - name: shared-data
          mountPath: /vault/secrets

    - name: cert-setup
      image: cfssl/cfssl:latest
      command:
        - sh
        - -c
        - |
          cfssl gencert \
            -ca=/vault/secrets/ca.crt \
            -ca-key=/vault/secrets/ca.key \
            -config=/etc/cfssl/config.json \
            -hostname="${POD_NAME}.${POD_NAMESPACE}.pod.cluster.local,${POD_IP}" \
            /etc/cfssl/csr.json | cfssljson -bare /vault/secrets/server
      env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
      volumeMounts:
        - name: shared-data
          mountPath: /vault/secrets
        - name: cfssl-config
          mountPath: /etc/cfssl

  containers:
    - name: app
      image: myapp:latest
      volumeMounts:
        - name: shared-data
          mountPath: /etc/tls
          readOnly: true
      ports:
        - containerPort: 8443
          name: https

  volumes:
    - name: shared-data
      emptyDir:
        medium: Memory
    - name: vault-config
      configMap:
        name: vault-agent-config
    - name: vault-agent-sa
      projected:
        sources:
          - serviceAccountToken:
              path: token
              expirationSeconds: 7200
              audience: vault
    - name: cfssl-config
      configMap:
        name: cfssl-config

---
# Pattern 4: Configuration Template Rendering
apiVersion: v1
kind: Pod
metadata:
  name: config-template-app
spec:
  initContainers:
    - name: config-renderer
      image: hairyhenderson/gomplate:v3.11
      command:
        - gomplate
        - -f
        - /templates/config.yaml.tmpl
        - -o
        - /config/config.yaml
      env:
        - name: DATABASE_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: database.host
        - name: DATABASE_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: database.port
        - name: REDIS_CLUSTER
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: redis.cluster
      volumeMounts:
        - name: config-templates
          mountPath: /templates
        - name: rendered-config
          mountPath: /config

  containers:
    - name: app
      image: myapp:latest
      volumeMounts:
        - name: rendered-config
          mountPath: /etc/myapp/config.yaml
          subPath: config.yaml
          readOnly: true

  volumes:
    - name: config-templates
      configMap:
        name: app-config-templates
    - name: rendered-config
      emptyDir: {}
```

### 3.3 Init Container Best Practices

```yaml
# Best Practice 1: Keep init containers lightweight
initContainers:
  - name: lightweight-check
    image: busybox:1.36  # Small image, fast pull
    command:
      - sh
      - -c
      - 'nc -z db 5432 || exit 1'
    resources:
      limits:
        cpu: 100m
        memory: 64Mi
      requests:
        cpu: 10m
        memory: 16Mi

# Best Practice 2: Set appropriate timeouts
initContainers:
  - name: dependency-check
    image: busybox:1.36
    command:
      - sh
      - -c
      - |
        for i in $(seq 1 60); do  # 60 second timeout
          nc -z db 5432 && exit 0
          sleep 1
        done
        exit 1

# Best Practice 3: Use shared volumes for data passing
initContainers:
  - name: generate-config
    image: python:3.12-alpine
    command:
      - python
      - -c
      - |
        import json, os
        config = {
          'host': os.environ['POD_IP'],
          'pod_name': os.environ['POD_NAME']
        }
        with open('/shared/config.json', 'w') as f:
          json.dump(config, f)
    env:
      - name: POD_IP
        valueFrom:
          fieldRef:
            fieldPath: status.podIP
      - name: POD_NAME
        valueFrom:
          fieldRef:
            fieldPath: metadata.name
    volumeMounts:
      - name: shared-config
        mountPath: /shared

# Best Practice 4: Handle idempotency
initContainers:
  - name: idempotent-setup
    image: postgres:16-alpine
    command:
      - sh
      - -c
      - |
        # Check if already initialized
        if psql "$DATABASE_URL" -c "SELECT 1 FROM schema_migrations" 2>/dev/null; then
          echo "Already initialized"
          exit 0
        fi
        # Run initialization
        psql "$DATABASE_URL" -f /init/00-schema.sql
        psql "$DATABASE_URL" -f /init/01-data.sql
```

---

## 4. Ambassador Pattern

The Ambassador pattern provides a proxy for remote service access, handling connection pooling, retries, circuit breaking, and protocol adaptation.

### 4.1 Ambassador vs Sidecar Comparison

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Ambassador vs Sidecar Pattern Comparison                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  SIDECAR Pattern (Traffic Management):                                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Service    │◄───────►│   Sidecar    │◄────► External      │   │
│  │  │              │         │  (Envoy)     │        Services      │   │
│  │  │              │         │              │                     │   │
│  │  │   Manages:   │         │   Manages:   │                     │   │
│  │  │   Business   │         │   • mTLS     │                     │   │
│  │  │   Logic      │         │   • Routing  │                     │   │
│  │  │              │         │   • Auth     │                     │   │
│  │  │              │         │   • Rate     │                     │   │
│  │  │              │         │     Limiting │                     │   │
│  │  └──────────────┘         └──────────────┘                     │   │
│  │                                                                 │   │
│  │  Scope: ALL traffic (inbound + outbound)                        │   │
│  │  Deployment: One per service                                    │   │
│  │  Configuration: Control plane managed                           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  AMBASSADOR Pattern (Service Proxy):                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Service    │◄───────►│  Ambassador  │◄────► PostgreSQL    │   │
│  │  │              │         │  (PgBouncer) │        Cluster      │   │
│  │  │              │         └──────────────┘                     │   │
│  │  │              │         ┌──────────────┐                     │   │
│  │  │   Thinks it  │◄───────►│  Ambassador  │◄────► Redis        │   │
│  │  │   connects   │         │  (Twemproxy) │        Cluster      │   │
│  │  │   directly   │         └──────────────┘                     │   │
│  │  │              │                                              │   │
│  │  │   localhost  │                                              │   │
│  │  │   :5432      │                                              │   │
│  │  │   :6379      │                                              │   │
│  │  └──────────────┘                                              │   │
│  │                                                                 │   │
│  │  Scope: ONE external service type per ambassador                │   │
│  │  Deployment: Can be per-pod or shared                           │   │
│  │  Configuration: Protocol-specific                               │   │
│  │  Key Feature: Connection pooling, protocol optimization         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  When to Use Which:                                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Use SIDECAR when:          │  Use AMBASSADOR when:             │   │
│  ├─────────────────────────────┼───────────────────────────────────┤   │
│  │  • Service mesh required    │  • Direct DB/cache connection     │   │
│  │  • Universal mTLS           │  • Connection pooling critical    │   │
│  │  • Traffic routing          │  • Protocol optimization          │   │
│  │  • Cross-cutting concerns   │  • Simplified client config       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Ambassador Implementation Examples

```yaml
# Ambassador 1: PostgreSQL Connection Pool (PgBouncer)
apiVersion: v1
kind: Pod
metadata:
  name: app-with-pgbouncer
spec:
  containers:
    - name: pgbouncer-ambassador
      image: pgbouncer/pgbouncer:1.22
      ports:
        - containerPort: 6432
          name: pgbouncer
      env:
        - name: DATABASES_HOST
          value: "postgres-primary.database.svc.cluster.local"
        - name: DATABASES_PORT
          value: "5432"
        - name: DATABASES_DATABASE
          value: "myapp"
        - name: POOL_MODE
          value: "transaction"
        - name: MAX_CLIENT_CONN
          value: "10000"
        - name: DEFAULT_POOL_SIZE
          value: "20"
        - name: RESERVE_POOL_SIZE
          value: "5"
        - name: RESERVE_POOL_TIMEOUT
          value: "3"
      volumeMounts:
        - name: pgbouncer-config
          mountPath: /etc/pgbouncer
      livenessProbe:
        tcpSocket:
          port: 6432
        initialDelaySeconds: 5
        periodSeconds: 10
      readinessProbe:
        tcpSocket:
          port: 6432
        initialDelaySeconds: 2
        periodSeconds: 5
      resources:
        limits:
          cpu: 500m
          memory: 128Mi
        requests:
          cpu: 100m
          memory: 64Mi

    - name: myapp
      image: myapp:latest
      env:
        # App connects to localhost ambassador, not remote DB
        - name: DATABASE_URL
          value: "postgresql://user:pass@localhost:6432/myapp?sslmode=require"
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

---
# Ambassador 2: Redis Cluster Proxy (Twemproxy)
apiVersion: v1
kind: Pod
metadata:
  name: app-with-redis-ambassador
spec:
  containers:
    - name: twemproxy-ambassador
      image: malevich/twemproxy:latest
      ports:
        - containerPort: 6379
          name: redis
      env:
        - name: REDIS_SERVERS
          value: |
            redis-0.redis-cluster.cache.svc.cluster.local:6379:0
            redis-1.redis-cluster.cache.svc.cluster.local:6379:1
            redis-2.redis-cluster.cache.svc.cluster.local:6379:2
        - name: HASH_TAG
          value: "{}"
        - name: DISTRIBUTION
          value: "ketama"
        - name: TIMEOUT
          value: "400"
        - name: TCPKEEPALIVE
          value: "true"
      livenessProbe:
        tcpSocket:
          port: 6379
        initialDelaySeconds: 5
        periodSeconds: 10
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 100m
          memory: 64Mi

    - name: myapp
      image: myapp:latest
      env:
        # App thinks it's connecting to single Redis
        - name: REDIS_URL
          value: "redis://localhost:6379"
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

---
# Ambassador 3: Kafka Protocol Proxy
apiVersion: v1
kind: Pod
metadata:
  name: app-with-kafka-ambassador
spec:
  containers:
    - name: kafka-proxy-ambassador
      image: grepplabs/kafka-proxy:latest
      command:
        - /kafka-proxy
        - server
        - --bootstrap-server-mapping
        - "kafka-0.kafka.svc.cluster.local:9093,localhost:32400"
        - --bootstrap-server-mapping
        - "kafka-1.kafka.svc.cluster.local:9093,localhost:32401"
        - --bootstrap-server-mapping
        - "kafka-2.kafka.svc.cluster.local:9093,localhost:32402"
        - --tls-enable
        - --tls-ca-chain-cert-file=/etc/kafka/certs/ca.crt
      ports:
        - containerPort: 32400
          name: kafka-0
        - containerPort: 32401
          name: kafka-1
        - containerPort: 32402
          name: kafka-2
      volumeMounts:
        - name: kafka-certs
          mountPath: /etc/kafka/certs
          readOnly: true
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 100m
          memory: 64Mi

    - name: myapp
      image: myapp:latest
      env:
        - name: KAFKA_BROKERS
          value: "localhost:32400,localhost:32401,localhost:32402"
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

  volumes:
    - name: kafka-certs
      secret:
        secretName: kafka-client-certs
```

---

## 5. Adapter Pattern

The Adapter pattern transforms the interface of a container to match what another container expects, handling protocol translation, data format conversion, and interface normalization.

### 5.1 Adapter Pattern Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Adapter Pattern Architecture                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Scenario: Legacy System Integration                                    │
│                                                                         │
│  Legacy SOAP Service ──▶ REST API Consumer                              │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │                                                                  │   │
│  │  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐   │   │
│  │  │   Adapter    │      │   Main App   │      │   External   │   │   │
│  │  │   (SOAP to   │◄────►│   (REST      │◄────►│   Client     │   │   │
│  │  │    REST)     │      │    API)      │      │              │   │   │
│  │  │              │      │              │      │              │   │   │
│  │  │  Inbound:    │      │  Exposes:    │      │  Consumes:   │   │   │
│  │  │  localhost   │      │  /api/v1/*   │      │  REST        │   │   │
│  │  │  :8081       │      │              │      │              │   │   │
│  │  │              │      │              │      │              │   │   │
│  │  │  Outbound:   │      │              │      │              │   │   │
│  │  │  SOAP calls  │      │              │      │              │   │   │
│  │  │  to legacy   │      │              │      │              │   │   │
│  │  └──────────────┘      └──────────────┘      └──────────────┘   │   │
│  │                                                                  │   │
│  │  Flow:                                                           │   │
│  │  1. External client calls Main App REST API                      │   │
│  │  2. Main App needs legacy data                                   │   │
│  │  3. Main App calls Adapter's REST interface (localhost)          │   │
│  │  4. Adapter converts REST to SOAP                                │   │
│  │  5. Adapter calls Legacy SOAP service                            │   │
│  │  6. Response converted SOAP → JSON → returned to Main App        │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Scenario: Protocol Translation (HTTP/1 to HTTP/2)                      │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐   │   │
│  │  │   Legacy     │      │   Adapter    │      │   Modern     │   │   │
│  │  │   Client     │◄────►│   (HTTP/2    │◄────►│   gRPC       │   │   │
│  │  │   (HTTP/1.1) │      │   to HTTP/1) │      │   Service    │   │   │
│  │  └──────────────┘      └──────────────┘      └──────────────┘   │   │
│  │                                                                  │   │
│  │  The adapter translates between HTTP/1.1 and HTTP/2/gRPC        │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Adapter Implementation Examples

```yaml
# Adapter 1: REST to gRPC Translation
apiVersion: v1
kind: Pod
metadata:
  name: rest-to-grpc-adapter
spec:
  containers:
    - name: grpc-adapter
      image: envoyproxy/envoy:v1.30
      ports:
        - containerPort: 8080
          name: http
      volumeMounts:
        - name: envoy-config
          mountPath: /etc/envoy
      command:
        - envoy
        - -c
        - /etc/envoy/envoy.yaml
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 100m
          memory: 64Mi

    - name: grpc-service
      image: myapp/grpc-service:latest
      ports:
        - containerPort: 50051
          name: grpc
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

  volumes:
    - name: envoy-config
      configMap:
        name: envoy-grpc-transcoder-config

---
# Adapter 2: Metrics Format Conversion (Prometheus to CloudWatch)
apiVersion: v1
kind: Pod
metadata:
  name: metrics-adapter
spec:
  containers:
    - name: cloudwatch-adapter
      image: prometheus/cloudwatch-exporter:latest
      ports:
        - containerPort: 9106
          name: http
      env:
        - name: AWS_REGION
          value: "us-east-1"
      volumeMounts:
        - name: cloudwatch-config
          mountPath: /config
      resources:
        limits:
          cpu: 500m
          memory: 512Mi
        requests:
          cpu: 100m
          memory: 128Mi

    - name: myapp
      image: myapp:latest
      ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

  volumes:
    - name: cloudwatch-config
      configMap:
        name: cloudwatch-exporter-config

---
# Adapter 3: Log Format Normalization
apiVersion: v1
kind: Pod
metadata:
  name: log-adapter
spec:
  containers:
    - name: fluentd-adapter
      image: fluent/fluentd:v1.16
      volumeMounts:
        - name: app-logs
          mountPath: /var/log/app
        - name: fluentd-config
          mountPath: /fluentd/etc
      resources:
        limits:
          cpu: 500m
          memory: 512Mi
        requests:
          cpu: 100m
          memory: 128Mi

    - name: legacy-app
      image: legacy-app:latest
      volumeMounts:
        - name: app-logs
          mountPath: /var/log
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

  volumes:
    - name: app-logs
      emptyDir:
        sizeLimit: 1Gi
    - name: fluentd-config
      configMap:
        name: fluentd-normalization-config
```

---

## 6. Configuration Helper Pattern

The Configuration Helper pattern uses init containers or sidecars to manage dynamic configuration, enabling hot reloading and external configuration sources.

### 6.1 Configuration Helper Architectures

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Configuration Helper Patterns                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pattern 1: Init Container Configuration Render                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │   │
│  │  │  Config      │───▶│  Config      │───▶│  Main App    │       │   │
│  │  │  Renderer    │    │  Volume      │    │  (read-only) │       │   │
│  │  │  (gomplate,  │    │  (emptyDir)  │    │              │       │   │
│  │  │   consul-    │    │              │    │  Watches:    │       │   │
│  │  │   template)  │    │  Templates:  │    │  Config file │       │   │
│  │  │              │    │  • app.yaml  │    │              │       │   │
│  │  │  Sources:    │    │  • certs.pem │    │  No reload   │       │   │
│  │  │  • Env vars  │    │              │    │  capability  │       │   │
│  │  │  • Secrets   │    │              │    │              │       │   │
│  │  │  • ConfigMaps│    │              │    │              │       │   │
│  │  └──────────────┘    └──────────────┘    └──────────────┘       │   │
│  │                                                                  │   │
│  │  Use Case: One-time configuration at startup                     │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 2: Sidecar Configuration Reloader (Hot Reload)                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │   │
│  │  │  Config      │───▶│  Config      │◄──▶│  Config      │       │   │
│  │  │  Reloader    │    │  Volume      │    │  Reloader    │       │   │
│  │  │  (native     │    │  (shared)    │    │  (native     │       │   │
│  │  │   sidecar)   │    │              │    │   sidecar)   │       │   │
│  │  │              │    │  Files:      │    │              │       │   │
│  │  │  Sources:    │    │  • config.yml│    │  Watches:    │       │   │
│  │  │  • etcd      │    │  • certs     │    │  Config file │       │   │
│  │  │  • Consul    │    │              │    │              │       │   │
│  │  │  • Vault     │    │              │    │  Action:     │       │   │
│  │  │  • AWS SSM   │    │              │    │  SIGUSR1 to  │       │   │
│  │  │              │    │              │    │  main app    │       │   │
│  │  └──────────────┘    └──────────────┘    └──────────────┘       │   │
│  │                                                                  │   │
│  │  Use Case: Dynamic configuration updates without restart         │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 3: External Configuration Server                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │  Config      │◄───────►│  Main App    │                     │   │
│  │  │  Server      │   gRPC   │              │                     │   │
│  │  │  (sidecar)   │          │  Fetches     │                     │   │
│  │  │              │          │  config on   │                     │   │
│  │  │  Sources:    │          │  demand      │                     │   │
│  │  │  • etcd      │          │              │                     │   │
│  │  │  • Consul    │          │  Subscribe:  │                     │   │
│  │  │  • Custom    │          │  Config      │                     │   │
│  │  │    backend   │          │  changes     │                     │   │
│  │  └──────────────┘          └──────────────┘                     │   │
│  │                                                                  │   │
│  │  Use Case: Complex configuration with dynamic updates            │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Configuration Helper Implementation Examples

```yaml
# Configuration Helper 1: ConfigMap Watcher with Hot Reload
apiVersion: v1
kind: Pod
metadata:
  name: app-with-config-reloader
spec:
  initContainers:
    # Initial config render
    - name: config-init
      image: k8s.gcr.io/git-sync:v3.6.4
      volumeMounts:
        - name: config-dir
          mountPath: /tmp/config

  containers:
    # Config reloader sidecar (native sidecar in 1.33+)
    - name: config-reloader
      image: jimmidyson/configmap-reload:v0.12.0
      args:
        - --volume-dir=/etc/config
        - --webhook-url=http://localhost:8080/-/reload
        - --webhook-method=POST
      volumeMounts:
        - name: config-volume
          mountPath: /etc/config
        - name: config-dir
          mountPath: /etc/config-source
      resources:
        limits:
          cpu: 100m
          memory: 64Mi
      restartPolicy: Always  # Native sidecar

    - name: prometheus
      image: prom/prometheus:v2.53.0
      ports:
        - containerPort: 9090
          name: http
      volumeMounts:
        - name: config-volume
          mountPath: /etc/prometheus
        - name: prometheus-data
          mountPath: /prometheus
      resources:
        limits:
          cpu: "2"
          memory: 8Gi

  volumes:
    - name: config-volume
      emptyDir: {}
    - name: config-dir
      configMap:
        name: prometheus-config
    - name: prometheus-data
      persistentVolumeClaim:
        claimName: prometheus-data

---
# Configuration Helper 2: Vault Agent for Dynamic Secrets
apiVersion: v1
kind: Pod
metadata:
  name: app-with-vault-secrets
spec:
  initContainers:
    - name: vault-agent-init
      image: hashicorp/vault:1.16
      args:
        - agent
        - -config=/etc/vault/agent-init.hcl
        - -exit-after-auth
      env:
        - name: VAULT_ADDR
          value: "https://vault.vault-system.svc.cluster.local:8200"
      volumeMounts:
        - name: vault-config
          mountPath: /etc/vault
        - name: shared-data
          mountPath: /vault/secrets
        - name: vault-sa-token
          mountPath: /var/run/secrets/kubernetes.io/serviceaccount

  containers:
    # Vault agent as native sidecar for secret renewal
    - name: vault-agent
      image: hashicorp/vault:1.16
      args:
        - agent
        - -config=/etc/vault/agent.hcl
      env:
        - name: VAULT_ADDR
          value: "https://vault.vault-system.svc.cluster.local:8200"
      volumeMounts:
        - name: vault-config
          mountPath: /etc/vault
        - name: shared-data
          mountPath: /vault/secrets
        - name: vault-sa-token
          mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 50m
          memory: 64Mi
      restartPolicy: Always

    - name: myapp
      image: myapp:latest
      volumeMounts:
        - name: shared-data
          mountPath: /etc/secrets
          readOnly: true
      env:
        - name: DATABASE_PASSWORD_FILE
          value: /etc/secrets/database-password
        - name: API_KEY_FILE
          value: /etc/secrets/api-key
      resources:
        limits:
          cpu: "1"
          memory: 512Mi

  volumes:
    - name: vault-config
      configMap:
        name: vault-agent-config
    - name: shared-data
      emptyDir:
        medium: Memory
    - name: vault-sa-token
      projected:
        sources:
          - serviceAccountToken:
              path: token
              expirationSeconds: 7200
              audience: vault

---
# Configuration Helper 3: GitOps-based Configuration
apiVersion: v1
kind: Pod
metadata:
  name: app-with-gitops-config
spec:
  initContainers:
    - name: git-sync-init
      image: k8s.gcr.io/git-sync:v4.2.0
      args:
        - --repo=https://github.com/myorg/configs.git
        - --branch=main
        - --root=/git
        - --depth=1
        - --one-time
      volumeMounts:
        - name: git-data
          mountPath: /git
        - name: git-creds
          mountPath: /etc/git-secret
          readOnly: true

  containers:
    # Git sync as native sidecar for continuous config updates
    - name: git-sync
      image: k8s.gcr.io/git-sync:v4.2.0
      args:
        - --repo=https://github.com/myorg/configs.git
        - --branch=main
        - --root=/git
        - --depth=1
        - --wait=60
        - --webhook-url=http://localhost:8080/-/reload
      volumeMounts:
        - name: git-data
          mountPath: /git
        - name: git-creds
          mountPath: /etc/git-secret
          readOnly: true
      resources:
        limits:
          cpu: 100m
          memory: 128Mi
      restartPolicy: Always

    - name: nginx
      image: nginx:alpine
      ports:
        - containerPort: 80
          name: http
      volumeMounts:
        - name: git-data
          mountPath: /etc/nginx/conf.d
          subPath: nginx/conf.d
        - name: git-data
          mountPath: /usr/share/nginx/html
          subPath: nginx/html
      resources:
        limits:
          cpu: 500m
          memory: 256Mi

  volumes:
    - name: git-data
      emptyDir: {}
    - name: git-creds
      secret:
        secretName: git-credentials
        defaultMode: 0400
```

---

## 7. Pattern Selection Guide

### 7.1 Decision Flowchart

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Multi-Container Pattern Selection Guide                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Start: What is your primary need?                                      │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Need to add cross-cutting concerns to a service?                 │  │
│  │  (logging, monitoring, mTLS, routing)                             │  │
│  └────────────────────┬─────────────────────────────────────────────┘  │
│                       │                                                │
│              ┌────────┴────────┐                                       │
│              ▼                 ▼                                       │
│            Yes                No                                       │
│              │                 │                                       │
│              ▼                 ▼                                       │
│  ┌──────────────────┐  ┌────────────────────────────────────────────┐ │
│  │ Use SIDECAR      │  │ Need to prepare environment before app?   │ │
│  │ pattern          │  │ (database setup, migrations, waiting)     │ │
│  │                  │  └────────────────────┬───────────────────────┘ │
│  │ • Istio/Linkerd  │                       │                        │
│  │ • Fluent Bit     │              ┌────────┴────────┐                │
│  │ • Prometheus     │              ▼                 ▼                │
│  │   exporter       │            Yes                No                │
│  └──────────────────┘              │                 │                │
│                                    ▼                 ▼                │
│                         ┌──────────────────┐  ┌─────────────────────┐ │
│                         │ Use INIT         │  │ Need to proxy to    │ │
│                         │ CONTAINER        │  │ external services?  │ │
│                         │ pattern          │  │ (database, cache)   │ │
│                         │                  │  └──────────┬──────────┘ │
│                         │ • Database       │             │            │
│                         │   migrations     │    ┌────────┴────────┐   │
│                         │ • Config render  │    ▼                 ▼   │
│                         │ • Cert setup     │  Yes                No   │
│                         └──────────────────┘    │                 │   │
│                                                 ▼                 ▼   │
│                                      ┌──────────────────┐  ┌────────┐ │
│                                      │ Use AMBASSADOR   │  │ Need   │ │
│                                      │ pattern          │  │ to     │ │
│                                      │                  │  │ trans- │ │
│                                      │ • PgBouncer      │  │ form   │ │
│                                      │ • Twemproxy      │  │ data   │ │
│                                      │ • Kafka proxy    │  │ or     │ │
│                                      └──────────────────┘  │ proto- │ │
│                                                            │ cols?  │ │
│                                                            └───┬────┘ │
│                                                    ┌───────────┴────┐ │
│                                                    ▼                ▼ │
│                                                  Yes               No │
│                                                   │                 │ │
│                                                   ▼                 ▼ │
│                                        ┌──────────────────┐  ┌─────┐ │
│                                        │ Use ADAPTER      │  │ Use │ │
│                                        │ pattern          │  │Simpl│ │
│                                        │                  │  │e    │ │
│                                        │ • REST/gRPC      │  │Multi│ │
│                                        │   transcoding    │  │-Cont│ │ │
│                                        │ • Protocol       │  │ainer│ │ │
│                                        │   conversion     │  └─────┘ │
│                                        └──────────────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Pattern Comparison Matrix

| Aspect | Sidecar | Init Container | Ambassador | Adapter | Config Helper |
|--------|---------|----------------|------------|---------|---------------|
| **Lifecycle** | Parallel to app | Before app starts | Parallel to app | Parallel to app | Init or parallel |
| **Primary Use** | Cross-cutting concerns | Setup/teardown | Connection proxy | Protocol translation | Config management |
| **Restart** | Independent (1.33+) | Never (or Always for native sidecar) | Independent | Independent | As configured |
| **Network** | Same namespace | Same namespace | localhost proxy | Translation layer | Same namespace |
| **Storage** | Shared volumes | Shared volumes (pass to app) | Local connection | May translate | Shared volumes |
| **Examples** | Istio, Fluent Bit | Migration, wait | PgBouncer, Twemproxy | Envoy, gRPC gateway | Vault agent, git-sync |

---

## 8. Best Practices Summary

### 8.1 General Multi-Container Best Practices

```yaml
# 1. Set resource limits for all containers
resources:
  limits:
    cpu: "500m"
    memory: "256Mi"
  requests:
    cpu: "100m"
    memory: "64Mi"

# 2. Use appropriate security contexts
securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
  seccompProfile:
    type: RuntimeDefault

# 3. Implement proper health checks
startupProbe:
  httpGet:
    path: /health/startup
    port: 8080
  failureThreshold: 30
  periodSeconds: 1

livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  periodSeconds: 5
  failureThreshold: 3

# 4. Use shared volumes appropriately
volumes:
  - name: shared-tmp
    emptyDir:
      medium: Memory  # tmpfs for sensitive data
      sizeLimit: 100Mi

  - name: shared-data
    emptyDir:
      sizeLimit: 1Gi  # Regular emptyDir for data
```

### 8.2 Native Sidecar Migration Checklist

When migrating from traditional sidecars to native sidecars (1.33+):

- [ ] Add `restartPolicy: Always` to init container definitions
- [ ] Move sidecar containers from `containers` to `initContainers`
- [ ] Add `startupProbe` to block main container until sidecar ready
- [ ] Remove manual sidecar readiness checks from main app
- [ ] Update Job resources to rely on automatic sidecar termination
- [ ] Test rolling update behavior with new pattern
- [ ] Monitor resource usage changes

---

## 9. References

1. [Kubernetes Sidecar Containers](https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/) - Official Documentation
2. [KEP-753: Sidecar Containers](https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/753-sidecar-containers) - Enhancement Proposal
3. [Multi-Container Patterns](https://kubernetes.io/docs/concepts/workloads/pods/#how-pods-manage-multiple-containers) - Kubernetes Docs
4. [CNCF Cloud Native Patterns](https://www.cncf.io/phippy/) - Illustrated Guide
5. [Google Cloud: Organizing Containers](https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-organizing-containers-with-pods) - Best Practices

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Kubernetes Version: 1.34.0+ (Native Sidecars Stable in 1.33+)*
