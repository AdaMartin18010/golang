# EC-076: Kubernetes 1.34 New Features - Comprehensive Guide

> **维度**: Engineering Cloud-Native
> **级别**: S (16+ KB)
> **标签**: #kubernetes #k8s-1.34 #DRA #gpu #sidecar #oci-artifacts #pod-certificates
> **权威来源**:
>
> - [Kubernetes 1.34 Release Notes](https://kubernetes.io/blog/2025/08/27/kubernetes-v1-34-release/) - Official Kubernetes Blog
> - [KEP-4381: DRA Structured Parameters](https://github.com/kubernetes/enhancements/issues/4381) - Kubernetes Enhancement Proposals
> - [KEP-4639: OCI Artifact Volumes](https://github.com/kubernetes/enhancements/issues/4639) - Kubernetes Enhancement Proposals
> - [KEP-4317: Pod Certificates](https://github.com/kubernetes/enhancements/issues/4317) - Kubernetes Enhancement Proposals
> - [CNCF Cloud Native Landscape](https://landscape.cncf.io/) - Cloud Native Computing Foundation

---

## 1. Kubernetes 1.34 Overview

Kubernetes v1.34, codenamed **"Of Wind & Will (O' WaW)"**, released on August 27, 2025, represents one of the most impactful releases since Kubernetes 1.28 introduced sidecar containers. This release delivers **58 enhancements** including **23 features graduating to Stable (GA)**, **22 features entering Beta**, and **13 new Alpha capabilities**.

### 1.1 Release Highlights

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Kubernetes 1.34 Release Statistics                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Total Enhancements:     58                                             │
│  ├─ Graduating to Stable: 23  ████████████████████████████████  40%    │
│  ├─ Entering Beta:        22  ████████████████████████████      38%    │
│  └─ New Alpha Features:   13  ████████████████                  22%    │
│                                                                         │
│  Key Focus Areas:                                                       │
│  • AI/ML Workload Support (DRA GA)                                     │
│  • GPU Resource Optimization                                           │
│  • Container Lifecycle Management                                      │
│  • Security & Authentication                                           │
│  • Storage & Volume Management                                         │
│  • Scheduling Performance                                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Feature Maturity Matrix

| Feature | Stage | Feature Gate | Default | Since |
|---------|-------|--------------|---------|-------|
| Dynamic Resource Allocation (DRA) | Stable | N/A | Enabled | 1.30 → 1.34 |
| Native Sidecar Containers | Stable | SidecarContainers | Enabled | 1.28 → 1.33 |
| OCI Artifact Volumes | Beta | ImageVolume | Enabled | 1.33 → 1.34 |
| PodCertificateRequest | Alpha | PodCertificateRequest | Disabled | 1.34 |
| Container Restart Rules | Alpha | ContainerRestartRules | Disabled | 1.34 |
| Per-Container Restart Policy | Alpha | ContainerRestartRules | Disabled | 1.34 |

---

## 2. Dynamic Resource Allocation (DRA) for GPUs - GA

Dynamic Resource Allocation (DRA) has graduated to **General Availability (GA)** in Kubernetes 1.34, providing a revolutionary way to manage specialized hardware resources like GPUs, TPUs, FPGAs, and custom accelerators.

### 2.1 DRA Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Dynamic Resource Allocation Architecture                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Kubernetes Control Plane                        │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐   │   │
│  │  │  Scheduler   │  │   API        │  │  Device Controller   │   │   │
│  │  │  (DRA-aware) │◄─┤  Server      │◄─┤  (ResourceSlice)     │   │   │
│  │  └──────┬───────┘  └──────────────┘  └──────────────────────┘   │   │
│  │         │                                                        │   │
│  │         │ Evaluates ResourceClaims with CEL expressions          │   │
│  │         ▼                                                        │   │
│  │  ┌──────────────────────────────────────────────────────────┐   │   │
│  │  │  ResourceClaim Template → ResourceClaim → DeviceClass    │   │   │
│  │  │  resource.k8s.io/v1 (stable in 1.34)                     │   │   │
│  │  └──────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│  ┌───────────────────────────┼──────────────────────────────────────┐  │
│  │                           ▼  Worker Node                           │  │
│  │  ┌──────────────────────────────────────────────────────────┐   │  │
│  │  │  Kubelet (DRA Plugin)                                    │   │  │
│  │  │  ┌──────────────────────────────────────────────────┐    │   │  │
│  │  │  │  DRA Device Driver (GPU/FPGA vendor-specific)    │    │   │  │
│  │  │  │  • ResourceSlice publishing                      │    │   │  │
│  │  │  │  • Device preparation & configuration            │    │   │  │
│  │  │  │  • CDI (Container Device Interface) injection    │    │   │  │
│  │  │  └──────────────────────────────────────────────────┘    │   │  │
│  │  └──────────────────────────────────────────────────────────┘   │  │
│  │                              │                                    │  │
│  │                              ▼                                    │  │
│  │  ┌──────────────────────────────────────────────────────────┐   │  │
│  │  │  Allocated Devices (GPUs, TPUs, FPGAs, NICs)            │   │  │
│  │  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐           │   │  │
│  │  │  │ GPU 0  │ │ GPU 1  │ │ GPU 2  │ │ GPU 3  │           │   │  │
│  │  │  │MIG 1g.5│ │MIG 2g.10│ │MIG 3g.20│ │ Full  │           │   │  │
│  │  │  └────────┘ └────────┘ └────────┘ └────────┘           │   │  │
│  │  └──────────────────────────────────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Core DRA API Types

Kubernetes 1.34 introduces the stable `resource.k8s.io/v1` API with these key resources:

```yaml
# DeviceClass - Defines device categories and their properties
apiVersion: resource.k8s.io/v1
kind: DeviceClass
metadata:
  name: nvidia-gpu-h100
spec:
  selectors:
    - cel:
        expression: "device.attributes['nvidia.com'].productName == 'NVIDIA H100'"
  # Configurable per-device parameters
  config:
    - opaque:
        driver: nvidia.com
        parameters:
          apiVersion: nvidia.com/v1
          kind: GpuConfig
          memoryClock: 2619
          computeMode: Default
---
# ResourceClaimTemplate - Template for workload resource claims
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: ml-training-gpu
  namespace: ai-workloads
spec:
  spec:
    devices:
      requests:
        - name: gpu
          deviceClassName: nvidia-gpu-h100
          allocationMode: ExactCount
          count: 2
      constraints:
        - matchAttribute: nvidia.com/gpu.memory
          minimum: 80Gi
      # CEL-based device filtering
      selectors:
        - cel:
            expression: |
              device.attributes['nvidia.com'].driverVersion >= '535.104.05'
---
# ResourceSlice - Published by device drivers to advertise devices
apiVersion: resource.k8s.io/v1
kind: ResourceSlice
metadata:
  name: node-gpu-pool-worker-1
  ownerReferences:
    - apiVersion: nvidia.com/v1
      kind: GpuDriver
      name: nvidia-driver-worker-1
spec:
  nodeName: worker-1
  driver: nvidia.com
  pool:
    name: h100-pool
    generation: 1
  devices:
    - name: gpu-0
      basic:
        attributes:
          nvidia.com/gpu.product:
            string: NVIDIA H100 80GB
          nvidia.com/gpu.memory:
            quantity: 80Gi
          nvidia.com/gpu.compute.major:
            int: 9
          nvidia.com/gpu.compute.minor:
            int: 0
        capacity:
          nvidia.com/gpu: "1"
          nvidia.com/gpu.memory: 80Gi
          nvidia.com/mig.capable: "true"
```

### 2.3 GPU Workload Example with DRA

```yaml
# Pod with DRA GPU resource claims
apiVersion: v1
kind: Pod
metadata:
  name: distributed-ml-training
  namespace: ai-workloads
spec:
  resourceClaims:
    - name: training-gpus
      resourceClaimTemplateName: ml-training-gpu
    - name: rdma-nic
      resourceClaimTemplateName: mellanox-cx7
  containers:
    - name: trainer
      image: nvcr.io/nvidia/pytorch:24.06-py3
      resources:
        requests:
          memory: 256Gi
          cpu: "64"
        limits:
          memory: 256Gi
          cpu: "64"
      env:
        - name: NVIDIA_VISIBLE_DEVICES
          value: "all"
        - name: NCCL_IB_HCA
          value: "mlx5_0,mlx5_1"
      volumeMounts:
        - name: dra-devices
          mountPath: /var/run/nvidia-devices
  volumes:
    - name: dra-devices
      cdi:
        driver: nvidia.com
        claimName: training-gpus
```

### 2.4 DRA Benefits for GPU Workloads

| Capability | Device Plugin | DRA (1.34 GA) |
|------------|---------------|---------------|
| **Partial GPU Allocation** | ❌ No | ✅ MIG, SR-IOV support |
| **Device Initialization** | ❌ No | ✅ Pre-configure FPGAs, reset GPUs |
| **Flexible Device Sharing** | ❌ Static | ✅ Consumable capacity (beta) |
| **Priority-based Allocation** | ❌ No | ✅ Prioritized device lists |
| **Admin Access Control** | ❌ No | ✅ Namespace-scoped admin access |
| **Cross-Node Resources** | ❌ Node-local | ✅ Future: fabric-attached devices |
| **CEL-based Filtering** | ❌ No | ✅ Complex selection logic |

### 2.5 GPU Cost Optimization with DRA

```
┌─────────────────────────────────────────────────────────────────────────┐
│              DRA GPU Cost Optimization Example                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Traditional GPU Allocation (Pre-DRA):                                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Job A: Needs 40GB GPU memory → Gets 1x A100 (80GB)            │   │
│  │  Job B: Needs 30GB GPU memory → Gets 1x A100 (80GB)            │   │
│  │  Job C: Needs 20GB GPU memory → Gets 1x A100 (80GB)            │   │
│  │                                                                  │   │
│  │  Total: 3x A100 (240GB allocated, 90GB utilized)               │   │
│  │  Efficiency: 37.5%                                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  With DRA MIG Support (1.34 GA):                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Single A100 80GB partitioned with MIG:                         │   │
│  │  ┌─────────────┬─────────────┬─────────────┐                    │   │
│  │  │ MIG 3g.40gb │ MIG 2g.20gb │ MIG 1g.10gb │  + 10gb spare      │   │
│  │  ├─────────────┼─────────────┼─────────────┤                    │   │
│  │  │   Job A     │   Job B     │   Job C     │                    │   │
│  │  │  (40GB)     │  (20GB)     │  (10GB)     │                    │   │
│  │  └─────────────┴─────────────┴─────────────┘                    │   │
│  │                                                                  │   │
│  │  Total: 1x A100 (70GB allocated, 70GB utilized)                 │   │
│  │  Efficiency: 100%                                               │   │
│  │  Cost Savings: ~66% reduction in GPU instances                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Real-World Impact (Source: ScaleOps Analysis):                        │
│  • 20-35% GPU cost reduction through flexible allocation               │
│  • 40-60% faster pod scheduling for AI/ML workloads                    │
│  • 85% reduction in GPU-related scheduling failures                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Native Sidecar Containers - Stable

Native sidecar containers, introduced as alpha in Kubernetes 1.28 and promoted to stable in 1.33, are fully supported in 1.34 with the `SidecarContainers` feature gate enabled by default.

### 3.1 Native Sidecar vs Traditional Sidecar

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Native Sidecar Container Lifecycle                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Traditional Sidecar (Pre-1.28):                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod Start:                                                      │   │
│  │  1. Init Containers (sequential)                                 │   │
│  │     └── init-1 → init-2 → ...                                    │   │
│  │  2. Regular Containers (parallel)                                │   │
│  │     ├── sidecar-1  ────────────────────────────────────────┐    │   │
│  │     ├── sidecar-2  ────────────────────────────────────┐   │    │   │
│  │     └── main-app     ◄── Problem: Starts before sidecars ready   │   │
│  │                                                                  │   │
│  │  Issues:                                                         │   │
│  │  • Race conditions between main app and sidecars                │   │
│  │  • No guaranteed startup order                                  │   │
│  │  • Sidecar restarts affect main app                             │   │
│  │  • Job completion waits forever (sidecars never exit)           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Native Sidecar (1.34 with SidecarContainers):                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod Start:                                                      │   │
│  │  1. Init Containers (sequential)                                 │   │
│  │     └── init-1 → init-2 → ...                                    │   │
│  │  2. Sidecar Containers (restartPolicy: Always)                   │   │
│  │     ├── envoy-sidecar   ─────────────────────────────────────┐   │   │
│  │     ├── istio-proxy     ─────────────────────────────────┐   │   │   │
│  │     │   Wait for all sidecars to be Ready...             │   │   │   │
│  │     │            ▼                                       │   │   │   │
│  │  3. Regular Containers                                     │   │   │   │
│  │     └── main-app  ◄── Guaranteed: Starts AFTER all sidecars ready  │   │
│  │                                                                  │   │
│  │  Pod Termination (Jobs):                                         │   │
│  │  1. Main container exits                                         │   │
│  │  2. Sidecars receive TERM signal → Graceful shutdown             │   │
│  │  3. Job completes successfully                                   │   │
│  │                                                                  │   │
│  │  Benefits:                                                       │   │
│  │  • Guaranteed startup ordering                                  │   │
│  │  • Independent restart policies                                 │   │
│  │  • Proper Job completion handling                               │   │
│  │  • Health checks block main container start                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Native Sidecar Pod Specification

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app-with-native-sidecars
  labels:
    app: web-service
    version: v2.0.0
spec:
  # Standard init containers (run first, must complete)
  initContainers:
    - name: db-migration
      image: myapp/db-migrator:v2.0.0
      restartPolicy: Never  # Traditional init container behavior
      env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url

    # NATIVE SIDECAR: restartPolicy: Always makes this a sidecar
    - name: istio-proxy
      image: istio/proxyv2:1.22.0
      restartPolicy: Always  # ← Key marker for native sidecar
      args:
        - proxy
        - sidecar
        - --configPath
        - /etc/istio/proxy
      ports:
        - containerPort: 15090
          name: http-envoy-prom
          protocol: TCP
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 2000m
          memory: 1Gi
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
        readOnlyRootFilesystem: true
      # Sidecar startup probes block main container start
      startupProbe:
        httpGet:
          path: /healthz/ready
          port: 15021
        initialDelaySeconds: 1
        periodSeconds: 1
        failureThreshold: 30
      livenessProbe:
        httpGet:
          path: /healthz/live
          port: 15021
        periodSeconds: 10

  # Regular containers start AFTER all native sidecars are Ready
  containers:
    - name: web-api
      image: myapp/api-server:v2.0.0
      ports:
        - containerPort: 8080
          name: http
      env:
        - name: SERVICE_MESH_ENABLED
          value: "true"
        # Application can assume istio-proxy is ready
        - name: PROXY_ADMIN_PORT
          value: "15000"
      resources:
        requests:
          cpu: 500m
          memory: 512Mi
```

### 3.3 Sidecar Container Patterns

```yaml
# Pattern 1: Service Mesh with Istio
apiVersion: v1
kind: Pod
metadata:
  name: mesh-enabled-service
spec:
  initContainers:
    - name: istio-init
      image: istio/proxyv2:1.22.0
      restartPolicy: Never  # Traditional init for iptables setup
      securityContext:
        capabilities:
          add:
            - NET_ADMIN
        privileged: true
    - name: istio-proxy
      image: istio/proxyv2:1.22.0
      restartPolicy: Always  # Native sidecar
      startupProbe:
        httpGet:
          path: /healthz/ready
          port: 15021
        failureThreshold: 30
  containers:
    - name: my-service
      image: my-service:latest

---
# Pattern 2: Logging/Monitoring Sidecar
apiVersion: v1
kind: Pod
metadata:
  name: app-with-observability
spec:
  initContainers:
    - name: log-shipper
      image: fluent/fluent-bit:3.0
      restartPolicy: Always
      volumeMounts:
        - name: app-logs
          mountPath: /var/log/app
        - name: fluent-config
          mountPath: /fluent-bit/etc/
      startupProbe:
        exec:
          command: ["/bin/sh", "-c", "test -S /fluent-bit/tmp/fluent-bit.sock"]
        failureThreshold: 10
  containers:
    - name: application
      image: myapp:latest
      volumeMounts:
        - name: app-logs
          mountPath: /var/log

---
# Pattern 3: Configuration Watcher Sidecar
apiVersion: v1
kind: Pod
metadata:
  name: dynamic-config-app
spec:
  initContainers:
    - name: config-reloader
      image: configmap-reload:v0.12.0
      restartPolicy: Always
      env:
        - name: CONFIGMAP_NAME
          value: "app-config"
        - name: WEBHOOK_URL
          value: "http://localhost:8080/-/reload"
      volumeMounts:
        - name: config-volume
          mountPath: /config
  containers:
    - name: prometheus
      image: prometheus:v2.53.0
      volumeMounts:
        - name: config-volume
          mountPath: /etc/prometheus
```

---

## 4. OCI Artifact Volumes - Beta

Kubernetes 1.34 introduces **OCI Artifact Volumes** (graduating to Beta), enabling pods to mount content directly from OCI-compliant registries as read-only volumes.

### 4.1 OCI Artifact Volume Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│              OCI Artifact Volume Architecture                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Traditional Approach (Pre-1.34):                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐     ┌──────────────────────────────────────┐  │   │
│  │  │   Main App   │◄────│  Init Container (downloader)         │  │   │
│  │  │              │     │  • Downloads ML models from S3       │  │   │
│  │  │              │     │  • Downloads configs from Git        │  │   │
│  │  │              │     │  • Downloads certs from vault        │  │   │
│  │  │              │     │  • Writes to emptyDir                │  │   │
│  │  └──────────────┘     └──────────────────────────────────────┘  │   │
│  │         │                                              │         │   │
│  │         └──────────────────────────────────────────────┘         │   │
│  │                        emptyDir volume                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│  Problems:                                                              │
│  • Init containers add startup latency                                 │
│  • Complex credential management                                       │
│  • No image layer caching benefits                                     │
│  • Bloated init container images                                       │
│                                                                         │
│  OCI Artifact Volume Approach (1.34 Beta):                              │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐     ┌──────────────────────────────────────┐  │   │
│  │  │   Main App   │◄────│  OCI Artifact Volume (kubelet)       │  │   │
│  │  │              │     │  • Mounts OCI image directly         │  │   │
│  │  │              │     │  • Uses image pull secrets           │  │   │
│  │  │              │     │  • Layer caching enabled             │  │   │
│  │  │              │     │  • Read-only mount                   │  │   │
│  │  └──────────────┘     └──────────────────────────────────────┘  │   │
│  │         │                                              │         │   │
│  │         └──────────────────────────────────────────────┘         │   │
│  │              OCI Registry (registry.hub.docker.com)               │   │
│  │              └── artifacts/my-model:v1.0 (non-runnable)           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│  Benefits:                                                              │
│  • Faster pod startup (no init container download)                     │
│  • Unified artifact distribution via registries                        │
│  • Layer deduplication and caching                                     │
│  • Standard image pull secrets                                         │
│  • Supply chain security (SLSA, Sigstore)                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 OCI Artifact Volume Examples

```yaml
# Example 1: ML Model Serving with OCI Artifact Volume
apiVersion: v1
kind: Pod
metadata:
  name: model-server
spec:
  volumes:
    - name: model-volume
      image:  # New volume source type in 1.34
        reference: myregistry.io/ml-models/llama-3-70b:v1.0.0
        pullPolicy: IfNotPresent
    - name: config-volume
      image:
        reference: myregistry.io/configs/model-server-config:v2.1
        pullPolicy: Always
  containers:
    - name: inference-server
      image: vllm/vllm-openai:v0.5.0
      volumeMounts:
        - name: model-volume
          mountPath: /models
          readOnly: true
        - name: config-volume
          mountPath: /etc/model-config
          readOnly: true
      resources:
        limits:
          nvidia.com/gpu: "1"
---
# Example 2: Security Profiles via OCI Artifacts
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
  annotations:
    seccomp.security.alpha.kubernetes.io/pod: localhost/profiles/custom.json
spec:
  volumes:
    - name: seccomp-profiles
      image:
        reference: security.io/profiles/seccomp-bundles:v2024.06
        pullPolicy: IfNotPresent
  containers:
    - name: app
      image: myapp:latest
      securityContext:
        seccompProfile:
          type: Localhost
          localhostProfile: profiles/restricted.json  # From OCI volume
      volumeMounts:
        - name: seccomp-profiles
          mountPath: /profiles
          readOnly: true
---
# Example 3: WebAssembly Modules as OCI Artifacts
apiVersion: v1
kind: Pod
metadata:
  name: wasm-runtime
spec:
  volumes:
    - name: wasm-modules
      image:
        reference: wasmhub.io/filters/auth-filter:v1.2.3
  containers:
    - name: envoy
      image: envoyproxy/envoy:v1.30
      volumeMounts:
        - name: wasm-modules
          mountPath: /etc/envoy/wasm
          readOnly: true
```

### 4.3 OCI Artifact Volume Specification

```yaml
# Full specification for OCI Artifact Volumes
apiVersion: v1
kind: Pod
metadata:
  name: comprehensive-oci-example
spec:
  imagePullSecrets:
    - name: registry-credentials
  volumes:
    - name: artifacts
      image:
        # Required: Full image reference
        reference: "registry.example.com/artifacts/my-bundle:v1.0"

        # Optional: Pull policy (IfNotPresent, Always, Never)
        pullPolicy: IfNotPresent

        # Optional: Selector for multi-arch images
        selector:
          # Match specific architecture
          architecture: amd64
          # Match specific OS
          os: linux

        # Optional: Additional image pull options
        pullSecret:
          # Use specific secret instead of pod-level secrets
          name: artifact-pull-secret
  containers:
    - name: consumer
      image: consumer:latest
      volumeMounts:
        - name: artifacts
          # Mounts as read-only by design
          mountPath: /artifacts
          readOnly: true
```

---

## 5. PodCertificateRequest - Alpha

Kubernetes 1.34 introduces **PodCertificateRequest** (alpha), enabling pods to obtain short-lived X.509 certificates directly from the Kubernetes API server without external certificate management systems.

### 5.1 PodCertificateRequest Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│              PodCertificateRequest Architecture                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  External CA Approach (Traditional):                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐                                               │   │
│  │  │   App        │◄── Must implement CSR to external CA         │   │
│  │  │              │    (Vault, cert-manager, etc.)               │   │
│  │  │              │                                               │   │
│  │  │  ┌────────┐  │    ┌──────────────┐    ┌─────────────────┐   │   │
│  │  │  │ Cert   │◄─┼────│  cert-manager│◄───│  Vault/External │   │   │
│  │  │  │ Store  │  │    │              │    │  CA             │   │   │
│  │  │  └────────┘  │    └──────────────┘    └─────────────────┘   │   │
│  │  └──────────────┘                                               │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│  Complexity:                                                            │
│  • External CA deployment required                                     │
│  • Custom sidecar or init container for CSR                            │
│  • Certificate rotation logic in app or sidecar                        │
│  • Network egress to CA required                                       │
│                                                                         │
│  PodCertificateRequest Approach (1.34 Alpha):                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                             │   │
│  │  ┌──────────────┐                                               │   │
│  │  │   App        │◄── Uses certificates from projected volume    │   │
│  │  │              │    (automatically rotated)                    │   │
│  │  │              │                                               │   │
│  │  │  ┌────────┐  │    ┌──────────────────────────────────────┐  │   │
│  │  │  │ /var/  │  │    │  Projected Volume                     │  │   │
│  │  │  │run/    │◄─┼────│  • podCertificate source (NEW)        │  │   │
│  │  │  │secrets/│  │    │  • kube-root-ca.crt                   │  │   │
│  │  │  │k8s.io/ │  │    │  • namespace downward API             │  │   │
│  │  │  │service │  │    └──────────────────────────────────────┘  │   │
│  │  │  │account│  │                                           │     │   │
│  │  │  └────────┘  │                                           │     │   │
│  │  └──────────────┘                                           │     │   │
│  │         │                                                    │     │   │
│  │         │ PodCertificateRequest CRD (created by kubelet)     │     │   │
│  │         ▼                                                    │     │   │
│  │  ┌────────────────────────────────────────────────────────┐  │     │   │
│  │  │  Kubernetes API Server                                  │  │     │   │
│  │  │  • Signs certificates with built-in CA                  │  │     │   │
│  │  │  • Supports external signers (via signerName)           │  │◄────┘   │
│  │  │  • Short-lived, auto-rotated certificates               │  │         │
│  │  │  • OIDC-compliant workload identity                     │  │         │
│  │  └────────────────────────────────────────────────────────┘  │         │
│  └─────────────────────────────────────────────────────────────────┘     │
│                                                                         │
│  Built-in Signers:                                                      │
│  • kubernetes.io/kube-apiserver-client-pod  (Pod identity)              │
│  • kubernetes.io/kubelet-serving           (Kubelet server certs)       │
│  • Custom external signers via CertificateSigningRequest API            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 PodCertificateRequest Example

```yaml
# Pod with automatic X.509 certificate provisioning
apiVersion: v1
kind: Pod
metadata:
  name: mtls-client
  namespace: secure-apps
spec:
  # Disable default service account token (use mTLS instead)
  automountServiceAccountToken: false

  volumes:
    - name: pod-certs
      projected:
        sources:
          # NEW in 1.34: Pod certificate provisioned by API server
          - podCertificate:
              signerName: "kubernetes.io/kube-apiserver-client-pod"
              credentialBundlePath: "credentialbundle.pem"
              # Certificate will be valid for:
              # - pod name
              # - pod UID
              # - service account

          # Include cluster CA for verifying server certificates
          - configMap:
              name: kube-root-ca.crt
              items:
                - key: ca.crt
                  path: kube-apiserver-root-certificate.pem

          # Include namespace for context
          - downwardAPI:
              items:
                - path: namespace
                  fieldRef:
                    fieldPath: metadata.namespace

  containers:
    - name: secure-client
      image: curlimages/curl:latest
      command: ["sleep", "infinity"]
      volumeMounts:
        - name: pod-certs
          mountPath: /var/run/secrets/kubernetes.io/serviceaccount
          readOnly: true
      env:
        # Certificate path for applications
        - name: CLIENT_CERT_PATH
          value: "/var/run/secrets/kubernetes.io/serviceaccount/credentialbundle.pem"
        - name: CA_CERT_PATH
          value: "/var/run/secrets/kubernetes.io/serviceaccount/kube-apiserver-root-certificate.pem"

---
# Direct PodCertificateRequest (advanced use cases)
apiVersion: certificates.k8s.io/v1
kind: PodCertificateRequest
metadata:
  name: my-pod-cert
  namespace: default
spec:
  # Signer selection
  signerName: "kubernetes.io/kube-apiserver-client-pod"

  # Requested DNS names (optional)
  dnsNames:
    - my-pod.default.svc.cluster.local
    - my-pod

  # Requested IP addresses (optional)
  ipAddresses:
    - 10.244.1.10

  # Certificate validity duration (signer may override)
  expirationSeconds: 3600  # 1 hour

  # Pod reference (must match requesting pod)
  podRef:
    name: my-pod
    uid: "pod-uid-here"
```

### 5.3 Use Cases for Pod Certificates

```yaml
# Use Case 1: Pod-to-Pod mTLS
apiVersion: v1
kind: Pod
metadata:
  name: service-a
spec:
  volumes:
    - name: tls-certs
      projected:
        sources:
          - podCertificate:
              signerName: "kubernetes.io/kube-apiserver-client-pod"
  containers:
    - name: app
      image: myapp:latest
      ports:
        - containerPort: 8443
      volumeMounts:
        - name: tls-certs
          mountPath: /etc/tls
      env:
        - name: TLS_CERT_FILE
          value: /etc/tls/credentialbundle.pem
        - name: TLS_CA_FILE
          value: /etc/tls/kube-apiserver-root-certificate.pem

---
# Use Case 2: SPIFFE/SPIRE Integration
apiVersion: v1
kind: Pod
metadata:
  name: spiffe-workload
spec:
  volumes:
    - name: spiffe-certs
      projected:
        sources:
          - podCertificate:
              # External signer for SPIFFE identities
              signerName: "spire.server/spiffe-ca"
              credentialBundlePath: "svid.pem"
  containers:
    - name: workload
      image: spiffe-workload:latest
      volumeMounts:
        - name: spiffe-certs
          mountPath: /spiffe-socket

---
# Use Case 3: Database Client Certificates
apiVersion: v1
kind: Pod
metadata:
  name: db-client
spec:
  volumes:
    - name: db-certs
      projected:
        sources:
          - podCertificate:
              signerName: "postgres-operator.io/client-ca"
              credentialBundlePath: "client-cert.pem"
  containers:
    - name: app
      image: app:latest
      env:
        - name: DATABASE_SSL_CERT
          value: /etc/db/certs/client-cert.pem
      volumeMounts:
        - name: db-certs
          mountPath: /etc/db/certs
```

---

## 6. Container Restart Rules - Alpha

Kubernetes 1.34 introduces fine-grained container restart policies with the `ContainerRestartRules` feature gate (alpha), allowing per-container restart policies and exit code-based restart rules.

### 6.1 Container Restart Architecture

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ml-training-worker
  labels:
    job: distributed-training
    worker: "3"
spec:
  # Pod-level policy (still required)
  restartPolicy: Never

  containers:
    - name: training-worker
      image: ml-platform.io/trainer:v1.0.0

      # NEW in 1.34: Per-container restart policy
      restartPolicy: Never  # Container won't restart by default

      # NEW in 1.34: Exit code-based restart rules
      restartPolicyRules:
        # Restart on retriable errors
        - action: Restart
          exitCodes:
            operator: In
            values:
              - 42   # CUDA out of memory - reduce batch size
              - 43   # Network timeout during gradient sync
              - 44   # Temporary storage I/O error

        # Do NOT restart on fatal errors
        - action: DoNotRestart
          exitCodes:
            operator: In
            values:
              - 1    # General error
              - 2    # Misuse of shell builtins
              - 126  # Command invoked cannot execute
              - 127  # Command not found

      resources:
        limits:
          nvidia.com/gpu: "1"
          memory: 32Gi
          cpu: "16"

      env:
        - name: BATCH_SIZE
          value: "auto"
        - name: CHECKPOINT_DIR
          value: "/checkpoints"
```

### 6.2 Container Restart Policy Examples

```yaml
# Example 1: Database Migration with One-Time Init
apiVersion: v1
kind: Pod
metadata:
  name: app-with-migration
spec:
  restartPolicy: Always  # Pod-level policy

  initContainers:
    - name: db-migrator
      image: myapp/db-migrator:v1.0.0
      restartPolicy: Never  # Run once, fail fast
      command: ["/migrate", "--target-version=v2.0"]
      env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url

  containers:
    - name: web-api
      image: myapp/api:v2.0.0
      restartPolicy: Always  # Always restart main app

---
# Example 2: Microservices with Different Restart Behaviors
apiVersion: v1
kind: Pod
metadata:
  name: payment-service-stack
spec:
  containers:
    - name: payment-api
      image: payments.io/api:v1.2.3
      restartPolicy: Always  # Critical service

    - name: metrics-exporter
      image: prom/node-exporter:v1.3.1
      restartPolicy: OnFailure  # Restart only on failure

    - name: log-forwarder
      image: fluent/fluent-bit:3.0
      restartPolicy: OnFailure
      restartPolicyRules:
        - action: Restart
          exitCodes:
            operator: In
            values: [1, 2]  # Elasticsearch connection issues
        - action: DoNotRestart
          exitCodes:
            operator: In
            values: [125]  # Configuration syntax error

---
# Example 3: CI/CD Runner with Specific Exit Code Handling
apiVersion: v1
kind: Pod
metadata:
  name: ci-runner
spec:
  restartPolicy: Never
  containers:
    - name: runner
      image: gitlab-runner:v16.0
      restartPolicy: Never
      restartPolicyRules:
        # Retry on infrastructure issues
        - action: Restart
          exitCodes:
            operator: In
            values:
              - 130  # Job cancelled (retry)
              - 137  # OOMKilled (retry with backoff)
        # Fail fast on build errors
        - action: DoNotRestart
          exitCodes:
            operator: NotIn
            values: [130, 137]
```

---

## 7. Features Graduating to Stable (23 Total)

Kubernetes 1.34 marks a maturity milestone with 23 features graduating to stable:

### 7.1 Complete List of GA Features

| # | Feature | KEP | SIG | Description |
|---|---------|-----|-----|-------------|
| 1 | **DRA: Structured Parameters** | #4381 | Node | Dynamic Resource Allocation for GPUs/Accelerators |
| 2 | **API Server Tracing** | #647 | API Machinery | OpenTelemetry tracing for API server |
| 3 | **Consistent Reads from Cache** | #2340 | API Machinery | Eliminates stale reads, 40-60% faster list ops |
| 4 | **VolumeAttributesClass** | #3751 | Storage | Dynamic volume performance modification |
| 5 | **RecoverVolumeExpansionFailure** | #1790 | Storage | Recovery from failed PVC expansion |
| 6 | **Job Pod Replacement Policy** | #3939 | Apps | Control pod replacement timing for Jobs |
| 7 | **AppArmor Support** | #2254 | Node | Native AppArmor profile management |
| 8 | **Node Memory Swap Support** | #2400 | Node | Limited swap usage for pods |
| 9 | **Kubelet OpenTelemetry Tracing** | #2832 | Node | Distributed tracing for kubelet |
| 10 | **Sleep Action for Lifecycle Hooks** | #3960 | Node | Sleep duration in preStop/postStart |
| 11 | **Relaxed DNS Search Validation** | #4427 | Network | Underscores in DNS search paths |
| 12 | **Decouple TaintManager** | #3902 | Scheduling | Separate TaintEvictionController |
| 13 | **QueueingHint in Scheduler** | #4247 | Scheduling | Per-plugin callback functions |
| 14 | **Direct Service Return (DSR)** | #5146 | Windows | DSR and overlay networking |
| 15 | **Ordered Namespace Deletion** | #5080 | API Machinery | Secure resource deletion ordering |
| 16 | **Streaming Encoding for LIST** | #5116 | API Machinery | 85% etcd load reduction |
| 17 | **Structured Authentication Config** | #3331 | Auth | Structured authn configuration |
| 18 | **Authorize with Selectors** | #4601 | Auth | Field/label selector authorization |
| 19 | **Resilient WatchCache Init** | #4568 | API Machinery | Watch cache initialization improvements |
| 20 | **Only Allow Anonymous for Configured Endpoints** | #5113 | Auth | Restricted anonymous access |
| 21 | **Relaxed Environment Variables** | #4369 | Node | Printable ASCII in env vars |
| 22 | **Discover Cgroup Driver from CRI** | #4569 | Node | Auto-detect cgroup driver |
| 23 | **Zero Sleep Action** | #5092 | Node | Allow zero duration in sleep hook |

### 7.2 Key Performance Improvements

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Kubernetes 1.34 Performance Improvements                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Snapshottable API Server Cache (Consistent Reads)                   │
│     • 40-60% faster LIST operations                                    │
│     • 85% reduction in etcd load                                       │
│     • Eliminates stale read issues                                     │
│     • Enables larger cluster scales                                    │
│                                                                         │
│  2. Streaming Encoding for LIST Responses                               │
│     • Memory usage reduction for large lists                           │
│     • Faster response serialization                                    │
│     • Reduced GC pressure on API server                                │
│                                                                         │
│  3. QueueingHint for Scheduler                                          │
│     • 30-50% reduction in unnecessary pod requeues                     │
│     • Faster scheduling for DRA workloads                              │
│     • Better handling of resource constraints                          │
│                                                                         │
│  4. DRA GPU Allocation                                                  │
│     • 20-35% GPU cost reduction through MIG/sharing                    │
│     • 40-60% faster pod scheduling for AI/ML                           │
│     • 85% reduction in GPU scheduling failures                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Upgrade Guide and Compatibility

### 8.1 Pre-Upgrade Checklist

```bash
#!/bin/bash
# Kubernetes 1.34 Pre-Upgrade Checklist

echo "=== Kubernetes 1.34 Pre-Upgrade Checklist ==="

# 1. Check for deprecated APIs
echo "Checking for deprecated APIs..."
kubectl get --raw /apis | jq -r '.groups[].versions[].groupVersion' | \
  grep -E "(flowcontrol|resource|coordination)" | sort -u

# 2. Check DRA feature usage
echo "Checking DRA ResourceClaims..."
kubectl get resourceclaims --all-namespaces 2>/dev/null || echo "No ResourceClaims found"

# 3. Verify VolumeAttributesClass migration
echo "Checking VolumeAttributesClass..."
kubectl get volumeattributesclasses --all-namespaces 2>/dev/null || echo "No VACs found"

# 4. Check AppArmor annotations (deprecated in favor of spec)
echo "Checking AppArmor annotations..."
kubectl get pods --all-namespaces -o json | \
  jq -r '.items[] | select(.metadata.annotations["container.apparmor.security.beta.kubernetes.io/*"]) | \
  "Pod: " + .metadata.name + " in " + .metadata.namespace'

# 5. Check for PSP usage (removed in 1.25+, ensure PSA migration)
echo "Checking Pod Security Standards..."
kubectl get namespaces -o json | \
  jq -r '.items[] | select(.metadata.labels["pod-security.kubernetes.io/enforce"]) | \
  "Namespace: " + .metadata.name + " - Level: " + .metadata.labels["pod-security.kubernetes.io/enforce"]'

echo "=== Pre-Upgrade Check Complete ==="
```

### 8.2 Feature Gate Activation

```yaml
# Enable alpha features (for testing only)
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
featureGates:
  ContainerRestartRules: true  # Alpha - Per-container restart policies
  PodCertificateRequest: true  # Alpha - Pod X.509 certificates

---
# API server configuration
apiVersion: apiserver.config.k8s.io/v1
kind: FeatureGate
featureGates:
  # Beta features (enabled by default, can be disabled)
  ImageVolume: true  # OCI Artifact Volumes

  # Alpha features
  ContainerRestartRules: true
  PodCertificateRequest: true
```

---

## 9. Best Practices and Recommendations

### 9.1 DRA Best Practices

```yaml
# 1. Use DeviceClasses for standardization
apiVersion: resource.k8s.io/v1
kind: DeviceClass
metadata:
  name: production-gpu
description: "Production GPU class with NVIDIA A100 or H100"
spec:
  selectors:
    - cel:
        expression: |
          device.attributes['nvidia.com'].productName.matches('NVIDIA (A100|H100)')

# 2. Implement resource quotas for DRA claims
apiVersion: v1
kind: ResourceQuota
metadata:
  name: gpu-quota
spec:
  hard:
    count/resourceclaims.resource.k8s.io: "10"
    nvidia.com/gpu: "8"

# 3. Use ResourceClaimTemplates for consistency
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: training-gpu-template
spec:
  spec:
    devices:
      requests:
        - name: gpu
          deviceClassName: production-gpu
          allocationMode: ExactCount
          count: 1
```

### 9.2 Native Sidecar Best Practices

```yaml
# 1. Always use startup probes for sidecars
initContainers:
  - name: envoy-sidecar
    image: envoyproxy/envoy:v1.30
    restartPolicy: Always
    startupProbe:
      httpGet:
        path: /ready
        port: 9901
      failureThreshold: 30  # Wait up to 30 seconds
      periodSeconds: 1

# 2. Set appropriate resource limits
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 2000m
    memory: 512Mi

# 3. Use proper security contexts
securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 65534
```

---

## 10. References and Further Reading

1. **Official Kubernetes 1.34 Release Notes**
   - <https://kubernetes.io/blog/2025/08/27/kubernetes-v1-34-release/>

2. **DRA Documentation**
   - <https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/>
   - <https://kubernetes.io/blog/2025/09/01/kubernetes-v1-34-dra-updates/>

3. **Native Sidecar Containers**
   - <https://kubernetes.io/blog/2025/08/29/kubernetes-v1-34-per-container-restart-policy/>
   - KEP-753: Sidecar Containers

4. **OCI Artifact Volumes**
   - KEP-4639: VolumeSource OCI Artifact
   - <https://github.com/kubernetes/enhancements/issues/4639>

5. **Pod Certificates**
   - KEP-4317: PodCertificateRequest
   - <https://github.com/kubernetes/enhancements/issues/4317>

6. **CNCF Resources**
   - <https://www.cncf.io/blog/2025/09/02/kubernetes-1-34-deep-dive/>
   - <https://landscape.cncf.io/>

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Kubernetes Version: 1.34.0+*
