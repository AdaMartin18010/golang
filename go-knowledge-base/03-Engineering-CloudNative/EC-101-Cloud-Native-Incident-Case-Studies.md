# EC-101: Cloud-Native Incident Case Studies

> **维度**: Engineering & Cloud Native | **级别**: S (15+ KB)
> **标签**: #kubernetes #containers #cloud-native #production-incidents #postmortem
> **权威来源**: Industry Postmortems, CNCF Case Studies, Real-world Incidents

---

## Overview

This document contains 10 detailed cloud-native production incident case studies, covering Kubernetes failures, container runtime issues, service mesh problems, and cloud infrastructure outages. Each case study includes incident description, root cause analysis, timeline, resolution steps, lessons learned, and prevention recommendations.

---

## Case Study 1: Kubernetes Control Plane Overload

### 1.1 Incident Description

**System**: Kubernetes cluster (500 nodes, 10,000 pods) running microservices
**Impact**: API server unresponsive, new pods can't be scheduled
**Duration**: 1 hour 45 minutes
**Date**: March 2024

A misconfigured CronJob with `concurrencyPolicy: Allow` created thousands of pods per minute. The Kubernetes API server was overwhelmed by LIST and WATCH requests, causing control plane failures and inability to schedule new workloads.

### 1.2 Root Cause Analysis

```
Failure Cascade:
1. CronJob configured to run every minute with 5-minute execution time
2. concurrencyPolicy: Allow (should be Forbid/Replace)
3. CronJob created new Job before previous completed
4. Jobs accumulated: 5 concurrent jobs per schedule
5. Each Job created 100 pods
6. After 1 hour: 30,000+ pods in Pending/Running state
7. Controller-manager overwhelmed creating pod records
8. API server etcd queries saturated
9. Scheduler couldn't process pending pods fast enough
10. Cluster effectively locked up
```

### 1.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 08:00:00 | CronJob deployed with wrong configuration |
| 08:01:00 | First Job created, pods starting |
| 08:06:00 | 6th Job created (overlapping) |
| 09:00:00 | 60+ concurrent Jobs, 6,000+ pods |
| 09:30:00 | API server response time > 10s |
| 09:45:00 | etcd warnings about request duration |
| 10:00:00 | Controller-manager restarts (OOM) |
| 10:15:00 | New pods stuck in Pending |
| 10:30:00 | Incident declared |
| 11:00:00 | Emergency CronJob deletion |
| 11:45:00 | Pod cleanup completed |

### 1.4 Resolution Steps

```bash
# Emergency pod deletion script
kubectl get jobs --all-namespaces -o json | \
    jq -r '.items[] | select(.spec.suspend != true) | \
    "kubectl delete job \(.metadata.name) -n \(.metadata.namespace)"' | bash

# Scale down CronJob
kubectl patch cronjob problematic-cronjob -p '{"spec":{"suspend":true}}'

# Delete orphan pods
kubectl get pods --all-namespaces --field-selector=status.phase=Pending | \
    grep "generated-from-cronjob" | \
    awk '{print "kubectl delete pod " $2 " -n " $1}' | bash
```

```go
// Resource quota enforcement
type CronJobValidator struct {
    maxConcurrentJobs int
    maxPodsPerJob     int
}

func (v *CronJobValidator) Validate(cronJob *batchv1.CronJob) error {
    // Check concurrency policy
    if cronJob.Spec.ConcurrencyPolicy == batchv1.AllowConcurrent {
        return fmt.Errorf("concurrencyPolicy: Allow is forbidden")
    }
    
    // Check job template pod count
    jobSpec := cronJob.Spec.JobTemplate.Spec.Template.Spec
    podCount := int32(0)
    
    for _, container := range jobSpec.Containers {
        if container.Resources.Requests.Cpu().MilliValue() > 1000 {
            return fmt.Errorf("container CPU request exceeds limit")
        }
    }
    
    // Check schedule frequency
    schedule := cronJob.Spec.Schedule
    if isHighFrequency(schedule) && cronJob.Spec.ConcurrencyPolicy == batchv1.AllowConcurrent {
        return fmt.Errorf("high-frequency CronJob must use Forbid or Replace")
    }
    
    return nil
}
```

### 1.5 Lessons Learned

1. **Default concurrencyPolicy** should be Forbid, not Allow
2. **Rate limiting** needed for cronjob creation
3. **Resource quotas** must include job/pod limits
4. **API server protection** requires priority levels

### 1.6 Prevention Recommendations

```yaml
# Pod security policy
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted-cronjob
spec:
  runAsUser:
    rule: MustRunAsNonRoot
  seLinux:
    rule: RunAsAny
  fsGroup:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  volumes:
  - 'configMap'
  - 'emptyDir'
  - 'projected'
  - 'secret'
  ---
# Resource quota for cronjobs
apiVersion: v1
kind: ResourceQuota
metadata:
  name: cronjob-quota
spec:
  hard:
    count/jobs.batch: "10"
    count/pods: "100"
    requests.cpu: "10"
    requests.memory: 20Gi
```

---

## Case Study 2: Container Image Pull Back-Off Storm

### 2.1 Incident Description

**System**: E-commerce platform with 200 microservices on Kubernetes
**Impact**: Rolling deployment failed, services unavailable
**Duration**: 45 minutes
**Date**: February 2024

A deployment triggered simultaneous image pulls from 500 nodes. The container registry rate limit (Docker Hub: 100 pulls/6hr) was exceeded, causing ImagePullBackOff errors. The retry storm prevented any new pods from starting.

### 2.2 Root Cause Analysis

```
Failure Chain:
1. Rolling update initiated for 50 deployments
2. MaxSurge: 25% allowed 200+ new pods simultaneously
3. All pods attempted image pull simultaneously
4. Docker Hub rate limit reached (anonymous: 100/6hr)
5. ImagePullBackOff with exponential backoff
6. Backoff capped at 5 minutes
7. Continuous retry attempts kept hitting rate limit
8. Deployment stalled, old pods remained
9. Some services running mixed versions
```

### 2.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 14:00:00 | Deployment pipeline triggered |
| 14:01:00 | First pods scheduled, image pull begins |
| 14:02:00 | Rate limit errors appearing |
| 14:05:00 | 403 Forbidden from registry |
| 14:10:00 | ImagePullBackOff on 80% of new pods |
| 14:20:00 | Deployment progress stalled |
| 14:30:00 | Manual rollback initiated |
| 14:45:00 | Services stabilized on previous version |

### 2.4 Resolution Steps

```bash
# Check image pull errors
kubectl get events --field-selector=reason=Failed | grep -i "pull"

# Update image pull secret with paid Docker Hub account
kubectl create secret docker-registry regcred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username=<username> \
    --docker-password=<password> \
    --docker-email=<email>

# Patch service account
kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "regcred"}]}'

# Force pod recreation with new secret
kubectl delete pods --all --grace-period=0 --force
```

```go
// Registry mirror configuration
type RegistryConfig struct {
    Mirrors []MirrorConfig `yaml:"mirrors"`
}

type MirrorConfig struct {
    Original string   `yaml:"original"`
    Mirrors  []string `yaml:"mirrors"`
}

// containerd config.toml
func generateContainerdConfig() string {
    return `
[plugins."io.containerd.grpc.v1.cri".registry.mirrors]
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
    endpoint = [
      "https://mirror.gcr.io",
      "https://registry-1.docker.io"
    ]
`
}
```

### 2.5 Lessons Learned

1. **Registry rate limits** are real production constraints
2. **Rolling update parameters** must consider pull capacity
3. **Registry mirrors** are essential for scale
4. **Private registries** provide better control

### 2.6 Prevention Recommendations

```yaml
# Deployment with controlled rollout
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-service
spec:
  replicas: 100
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: "5%"        # Reduced from 25%
      maxUnavailable: "5%"  # Slow rollout
  template:
    spec:
      imagePullSecrets:
      - name: regcred
      containers:
      - name: api
        image: myregistry.com/api:v1.2.3
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
```

---

## Case Study 3: Service Mesh Circuit Breaker Misconfiguration

### 3.1 Incident Description

**System**: Microservices platform with Istio service mesh (500+ services)
**Impact**: Cascading failures, complete service degradation
**Duration**: 2 hours
**Date**: January 2024

An overly aggressive circuit breaker configuration (5 errors in 1 minute) caused healthy services to be isolated. The cascading effect led to a mesh-wide outage where critical paths were blocked by incorrectly opened circuits.

### 3.2 Root Cause Analysis

```
Failure Cascade:
1. Istio DestinationRule configured aggressive circuit breaker:
   - consecutive5xxErrors: 5
   - interval: 1m
   - baseEjectionTime: 30m
2. Temporary latency spike caused timeout errors
3. Circuit breaker opened for healthy service
4. Traffic diverted to remaining instances
5. Remaining instances overloaded, more errors
6. Their circuits opened too
7. All instances ejected from pool
8. Upstream services received 503 errors
9. Their circuits opened, propagating failure
10. Complete mesh outage within 10 minutes
```

### 3.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 09:00:00 | Latency spike in database layer |
| 09:02:00 | First circuit breaker opens |
| 09:05:00 | Multiple services showing 503 errors |
| 09:10:00 | Circuit breakers opening across mesh |
| 09:15:00 | Core services unavailable |
| 09:30:00 | Incident declared |
| 10:00:00 | Istio configuration rolled back |
| 11:00:00 | Services gradually recovering |

### 3.4 Resolution Steps

```yaml
# Emergency circuit breaker relaxation
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: api-service
spec:
  host: api-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
    outlierDetection:
      consecutive5xxErrors: 50      # Increased from 5
      interval: 5m                   # Increased from 1m
      baseEjectionTime: 1m          # Decreased from 30m
      maxEjectionPercent: 50        # Limit ejection scope
```

```go
// Circuit breaker monitoring
type CircuitBreakerMonitor struct {
    metrics *prometheus.GaugeVec
}

func (m *CircuitBreakerMonitor) RecordState(service string, state CircuitState) {
    m.metrics.WithLabelValues(service).Set(float64(state))
}

func (m *CircuitBreakerMonitor) AlertIfTooManyOpen(threshold int) {
    openCount := 0
    for service, state := range m.getAllStates() {
        if state == Open {
            openCount++
            log.Printf("Circuit open: %s", service)
        }
    }
    
    if openCount > threshold {
        alert("Too many circuits open: %d", openCount)
    }
}
```

### 3.5 Lessons Learned

1. **Aggressive circuit breakers** cause more harm than good
2. **Cascading failures** propagate quickly in service meshes
3. **Gradual ejection** is safer than immediate isolation
4. **Max ejection percentage** limits blast radius

### 3.6 Prevention Recommendations

```yaml
# Conservative circuit breaker defaults
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: default-circuit-breaker
spec:
  host: "*.svc.cluster.local"
  trafficPolicy:
    outlierDetection:
      consecutive5xxErrors: 10
      interval: 30s
      baseEjectionTime: 30s
      maxEjectionPercent: 10
    connectionPool:
      http:
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
        maxRequestsPerConnection: 100
```

---

## Case Study 4: Horizontal Pod Autoscaler Thrashing

### 4.1 Incident Description

**System**: Video streaming service with HPA-based autoscaling
**Impact**: CPU exhaustion, service instability, increased cloud costs
**Duration**: 3 hours
**Date**: December 2023

Aggressive HPA scaling (scale down delay: 60s) caused rapid scale-up/scale-down cycles. Pods were terminated before warming up, causing continuous oscillation and CPU exhaustion from constant pod creation/destruction.

### 4.2 Root Cause Analysis

```
Oscillation Pattern:
1. Traffic spike: CPU 70%, HPA scales up (+10 pods)
2. New pods take 90s to be ready (warmup time)
3. During warmup, existing pods handle all load
4. Existing pods CPU spikes to 90%
5. HPA scales up again (+10 more pods)
6. After 60s, initial pods still not ready
7. HPA sees average CPU dropped (new pods at 0%)
8. HPA scales down (removes newest pods)
9. Original pods now ready, but some terminated
10. Cycle repeats continuously
```

### 4.3 Timeline of Events

| Time (UTC) | Event | Pods |
|------------|-------|------|
| 20:00:00 | Baseline | 10 |
| 20:05:00 | Traffic spike, scale up | 20 |
| 20:06:00 | Scale up (pods warming) | 30 |
| 20:07:00 | Scale down triggered | 25 |
| 20:08:00 | Scale up again | 35 |
| 20:09:00 | Scale down | 30 |
| 20:15:00 | Continuous oscillation | 20-40 |
| 22:30:00 | Manual HPA adjustment |
| 23:00:00 | Stabilized at 25 pods |

### 4.4 Resolution Steps

```yaml
# Stabilized HPA configuration
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: video-service
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: video-service
  minReplicas: 10
  maxReplicas: 100
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # 5 minutes
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
      - type: Pods
        value: 10
        periodSeconds: 60
      selectPolicy: Max
```

```go
// Custom metrics for better scaling decisions
type ScalingMetrics struct {
    cpuUtilization    float64
    requestLatency    float64
    queueDepth        int
    podReadyPercent   float64
}

func shouldScale(metrics ScalingMetrics) (bool, string) {
    // Don't scale if pods aren't ready
    if metrics.podReadyPercent < 80 {
        return false, "pods_still_warming"
    }
    
    // Consider latency for scale decisions
    if metrics.requestLatency > 500 && metrics.cpuUtilization > 60 {
        return true, "high_latency_and_cpu"
    }
    
    return false, "stable"
}
```

### 4.5 Lessons Learned

1. **Stabilization windows** prevent thrashing
2. **Pod warmup time** must be considered
3. **Custom metrics** provide better scaling signals
4. **Scale down should be slower** than scale up

### 4.6 Prevention Recommendations

```yaml
# HPA behavior best practices
behavior:
  scaleDown:
    stabilizationWindowSeconds: 300
    policies:
    - type: Percent
      value: 10        # Max 10% reduction per minute
      periodSeconds: 60
  scaleUp:
    stabilizationWindowSeconds: 0  # Immediate scale up
    policies:
    - type: Percent
      value: 100       # Double if needed
      periodSeconds: 15
```

---

## Case Study 5: etcd Snapshot Restoration Failure

### 5.1 Incident Description

**System**: Kubernetes cluster control plane (3-node etcd cluster)
**Impact**: Complete cluster data loss, 12-hour recovery
**Date**: November 2023

A failed etcd snapshot restoration attempt corrupted the data directory. The backup was taken during high write load and was inconsistent. The restoration on a new cluster revealed the corruption, requiring manual data reconstruction.

### 5.2 Root Cause Analysis

```
Failure Chain:
1. etcd snapshot taken with: etcdctl snapshot save
2. Snapshot taken during heavy API server load
3. Snapshot didn't capture all recent writes
4. Original cluster lost (hardware failure)
5. New cluster provisioned
6. Snapshot restore appeared successful
7. API server couldn't read some keys (checksum errors)
8. Some deployments, secrets missing
9. Partial data loss confirmed
```

### 5.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| Day -1 23:00 | Snapshot taken (automated) |
| 00:00:00 | Hardware failure on etcd nodes |
| 00:30:00 | Control plane unavailable |
| 01:00:00 | Recovery initiated |
| 02:00:00 | Snapshot restore completed |
| 02:30:00 | API server errors detected |
| 03:00:00 | Data inconsistency confirmed |
| 14:00:00 | Manual data reconstruction completed |

### 5.4 Resolution Steps

```bash
# Proper snapshot with consistency check
ETCDCTL_API=3 etcdctl snapshot save snapshot.db
ETCDCTL_API=3 etcdctl snapshot status snapshot.db -w table

# Restore with verification
ETCDCTL_API=3 etcdctl snapshot restore snapshot.db \
    --data-dir=/var/lib/etcd-restored \
    --name=etcd-1 \
    --initial-cluster="etcd-1=https://10.0.0.1:2380" \
    --initial-cluster-token=etcd-cluster-1 \
    --initial-advertise-peer-urls=https://10.0.0.1:2380

# Verify restoration
ETCDCTL_API=3 etcdctl --endpoints=https://localhost:2379 \
    --cacert=/etc/etcd/ca.crt \
    --cert=/etc/etcd/server.crt \
    --key=/etc/etcd/server.key \
    endpoint status -w table

# Check for data integrity
ETCDCTL_API=3 etcdctl get "" --prefix --keys-only | wc -l
```

```go
// Automated backup verification
type BackupVerifier struct {
    etcdClient *clientv3.Client
}

func (bv *BackupVerifier) VerifyBackup(backupPath string) error {
    // Restore to temp directory
    tempDir, err := os.MkdirTemp("", "etcd-verify-*")
    if err != nil {
        return err
    }
    defer os.RemoveAll(tempDir)
    
    // Run etcd in temporary mode
    cmd := exec.Command("etcd", 
        "--data-dir", tempDir,
        "--listen-client-urls", "http://localhost:12379",
        "--advertise-client-urls", "http://localhost:12379",
    )
    
    // Check if we can read all keys
    cli, err := clientv3.New(clientv3.Config{
        Endpoints: []string{"localhost:12379"},
    })
    if err != nil {
        return err
    }
    
    resp, err := cli.Get(context.Background(), "", clientv3.WithPrefix())
    if err != nil {
        return fmt.Errorf("backup verification failed: %w", err)
    }
    
    log.Printf("Backup verified: %d keys", len(resp.Kvs))
    return nil
}
```

### 5.5 Lessons Learned

1. **Snapshot during low traffic** or use consistent snapshot
2. **Automated backup verification** is essential
3. **Multiple backup copies** in different locations
4. **Regular restore drills** validate backup integrity

### 5.6 Prevention Recommendations

```yaml
# Velero backup configuration
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: etcd-backup
spec:
  schedule: "0 */4 * * *"  # Every 4 hours
  template:
    includedNamespaces:
    - "*"
    excludedResources:
    - events
    - pods
    storageLocation: default
    volumeSnapshotLocations:
    - aws-default
    ttl: 720h0m0s  # 30 days
```

---

## Case Study 6: DNS Resolution Storm

### 6.1 Incident Description

**System**: Microservices platform with CoreDNS (3 replicas)
**Impact**: DNS resolution failures, service-to-service calls failing
**Duration**: 35 minutes
**Date**: October 2023

A service with a connection leak created 100,000+ goroutines, each making DNS queries. CoreDNS was overwhelmed, causing cluster-wide DNS resolution failures and cascading service degradation.

### 6.2 Root Cause Analysis

```
Failure Cascade:
1. Service had connection leak (not closing HTTP responses)
2. Each request created new connection
3. Each connection triggered DNS lookup
4. DNS client cache TTL: 5 seconds (too short)
5. 10,000 RPS * 5 second TTL = 50,000 DNS queries/sec
6. CoreDNS (3 replicas) couldn't handle load
7. DNS queries timing out
8. Services couldn't resolve other services
9. Circuit breakers opened on DNS failures
10. Platform-wide outage
```

### 6.3 Timeline of Events

| Time (UTC) | Event | DNS QPS |
|------------|-------|---------|
| 15:00:00 | Normal operations | 1000 |
| 15:10:00 | Connection leak begins | 5000 |
| 15:15:00 | DNS queries increasing | 20000 |
| 15:20:00 | CoreDNS showing high CPU | 50000 |
| 15:25:00 | DNS timeouts reported | 60000 |
| 15:30:00 | Service degradation begins | 80000 |
| 15:35:00 | Problematic service identified | - |
| 16:00:00 | Service rolled back, DNS recovered | 1000 |

### 6.4 Resolution Steps

```yaml
# CoreDNS scaling and caching
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health {
            lameduck 5s
        }
        ready
        # Enhanced caching
        cache {
            success 9984 300  # Cache successful responses 5 min
            denial 9984 60    # Cache NXDOMAIN 1 min
            prefetch 10       # Prefetch before TTL expires
        }
        kubernetes cluster.local in-addr.arpa ip6.arpa {
            pods insecure
            fallthrough in-addr.arpa ip6.arpa
            ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf {
            max_concurrent 1000
        }
        reload
        loadbalance
    }
```

```go
// Connection pooling with DNS caching
type HTTPClient struct {
    client *http.Client
    dialer *net.Dialer
}

func NewHTTPClient() *HTTPClient {
    // Custom resolver with caching
    resolver := &CachingResolver{
        ttl:    5 * time.Minute,
        cache:  make(map[string]cacheEntry),
    }
    
    dialer := &net.Dialer{
        Resolver: &net.Resolver{
            PreferGo: true,
            Dial: resolver.Dial,
        },
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }
    
    transport := &http.Transport{
        DialContext:           dialer.DialContext,
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   10,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }
    
    return &HTTPClient{
        client: &http.Client{
            Transport: transport,
            Timeout:   30 * time.Second,
        },
    }
}
```

### 6.5 Lessons Learned

1. **DNS is critical infrastructure** at scale
2. **Connection pooling** prevents DNS storms
3. **DNS caching** must be configured properly
4. **CoreDNS needs HPA** based on query rate

### 6.6 Prevention Recommendations

```yaml
# CoreDNS HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: coredns
  namespace: kube-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: coredns
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 30
      policies:
      - type: Pods
        value: 5
        periodSeconds: 15
```

---

## Case Study 7: Persistent Volume Data Loss

### 7.1 Incident Description

**System**: Stateful application with PostgreSQL on Kubernetes
**Impact**: Database corruption, 6 hours of data loss
**Date**: September 2023

A storage node failure combined with a misconfigured StorageClass `reclaimPolicy: Delete` caused PVC deletion to wipe the underlying PV data. The backup restore revealed the backup job had been failing silently for 2 weeks.

### 7.2 Root Cause Analysis

```
Failure Chain:
1. Storage node hardware failure
2. Pod rescheduled to different node
3. PVC couldn't attach (node affinity)
4. PVC deleted to force re-creation
5. StorageClass had reclaimPolicy: Delete
6. Underlying cloud disk deleted automatically
7. New PVC created new empty disk
8. Attempted backup restore
9. Last successful backup: 2 weeks ago
10. Silent backup job failures not alerted
```

### 7.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 03:00:00 | Storage node failure |
| 03:05:00 | PostgreSQL pod CrashLoopBackOff |
| 03:30:00 | PVC deleted and recreated |
| 03:35:00 | New empty volume attached |
| 04:00:00 | Data loss confirmed |
| 04:30:00 | Backup restoration attempted |
| 05:00:00 | Last backup: 2 weeks old |
| 06:00:00 | Incident escalated |
| 12:00:00 | Data reconstruction completed |

### 7.4 Resolution Steps

```yaml
# Safe StorageClass configuration
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: standard-retain
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-standard
  replication-type: regional-dual-replica
reclaimPolicy: Retain  # NEVER Delete for production data
allowVolumeExpansion: true
mountOptions:
  - debug
volumeBindingMode: WaitForFirstConsumer
```

```go
// PVC protection webhook
type PVCProtectionWebhook struct {
    validator Validator
}

func (w *PVCProtectionWebhook) ValidateDelete(ctx context.Context, pvc *v1.PersistentVolumeClaim) error {
    // Check if PVC is in use by any pod
    pods, err := w.getPodsUsingPVC(pvc)
    if err != nil {
        return err
    }
    
    if len(pods) > 0 {
        return fmt.Errorf("PVC %s/%s is still in use by pods: %v", 
            pvc.Namespace, pvc.Name, podNames(pods))
    }
    
    // Check for recent backups
    if !w.hasRecentBackup(pvc) {
        return fmt.Errorf("PVC %s/%s has no recent backup (within 24h)",
            pvc.Namespace, pvc.Name)
    }
    
    return nil
}
```

### 7.5 Lessons Learned

1. **reclaimPolicy: Delete** is dangerous for stateful apps
2. **Backup monitoring** must verify backup success
3. **PVC finalizers** can prevent accidental deletion
4. **Regional replication** provides better durability

### 7.6 Prevention Recommendations

```yaml
# Velero backup with monitoring
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: postgres-backup
spec:
  schedule: "0 */6 * * *"  # Every 6 hours
  template:
    includedNamespaces:
    - database
    labelSelector:
      matchLabels:
        app: postgres
    storageLocation: default
    ttl: 720h0m0s
    
---
# Backup verification job
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-verifier
spec:
  schedule: "0 */6 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: verifier
            image: velero/velero:latest
            command:
            - /bin/sh
            - -c
            - |
              velero backup get --output json | \
              jq -e '.items | map(select(.status.phase != "Completed")) | length == 0'
```

---

## Case Study 8: Resource Quota Exhaustion

### 8.1 Incident Description

**System**: Multi-tenant Kubernetes cluster with namespace quotas
**Impact**: Critical services unable to scale, customer impact
**Duration**: 2 hours
**Date**: August 2023

A runaway batch job in a shared namespace consumed the entire namespace's CPU quota. Critical services in the same namespace couldn't create new pods during a traffic spike, causing customer-facing outages.

### 8.2 Root Cause Analysis

```
Resource Contention:
1. Namespace quota: 100 CPU, 200Gi memory
2. Batch job created 50 pods x 2 CPU = 100 CPU
3. Critical service needed to scale from 10 to 30 pods
4. Scale-up blocked: "exceeded quota"
5. Service couldn't handle traffic load
6. Customer requests failed with 503
7. Batch job completed after 2 hours
8. Quota freed, service scaled successfully
```

### 8.3 Timeline of Events

| Time (UTC) | Event | CPU Used |
|------------|-------|----------|
| 10:00:00 | Normal operations | 60/100 |
| 10:15:00 | Batch job started | 60/100 |
| 10:20:00 | Batch pods scheduled | 160/100 (pending) |
| 10:25:00 | Critical service needs scale | 100/100 (blocked) |
| 10:30:00 | Customer impact begins | 100/100 |
| 10:45:00 | Root cause identified | 100/100 |
| 11:00:00 | Attempted quota increase | 100/100 |
| 12:00:00 | Batch job completes | 10/100 |
| 12:05:00 | Service scaled successfully | 50/100 |

### 8.4 Resolution Steps

```yaml
# ResourceQuota with scope selectors
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-quota
spec:
  hard:
    requests.cpu: "100"
    requests.memory: 200Gi
    limits.cpu: "200"
    limits.memory: 400Gi
  scopeSelector:
    matchExpressions:
    - operator: NotIn
      scopeName: PriorityClass
      values: ["system-critical", "production-high"]
---
# Critical service priority
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: system-critical
value: 1000000
globalDefault: false
description: "Critical system services"
---
# Separate quotas for batch workloads
apiVersion: v1
kind: ResourceQuota
metadata:
  name: batch-quota
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
  scopes:
  - BestEffort
  - NotTerminating
```

```go
// Quota monitoring
type QuotaMonitor struct {
    client kubernetes.Interface
}

func (qm *QuotaMonitor) AlertOnHighUsage(namespace string, threshold float64) {
    quota, err := qm.client.CoreV1().ResourceQuotas(namespace).Get(
        context.Background(), "compute-quota", metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get quota: %v", err)
        return
    }
    
    cpuUsed := quota.Status.Used[corev1.ResourceCPU]
    cpuHard := quota.Status.Hard[corev1.ResourceCPU]
    
    usage := float64(cpuUsed.MilliValue()) / float64(cpuHard.MilliValue())
    if usage > threshold {
        alert("Quota usage high in %s: %.0f%%", namespace, usage*100)
    }
}
```

### 8.5 Lessons Learned

1. **Shared quotas** create resource contention
2. **Priority classes** should bypass quota limits
3. **Separate quotas** for batch vs service workloads
4. **Quota usage alerts** needed before exhaustion

### 8.6 Prevention Recommendations

```yaml
# LimitRange for default resource allocation
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
spec:
  limits:
  - default:
      cpu: "500m"
      memory: "512Mi"
    defaultRequest:
      cpu: "100m"
      memory: "128Mi"
    type: Container
  - max:
      cpu: "2"
      memory: "4Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
    type: Container
```

---

## Case Study 9: Admission Webhook Timeout

### 9.1 Incident Description

**System**: Kubernetes cluster with multiple admission webhooks
**Impact**: All pod creation blocked, deployment pipeline broken
**Duration**: 1 hour 30 minutes
**Date**: July 2023

A validating admission webhook for pod security policies experienced high latency due to database connection pool exhaustion. The 30-second timeout caused all pod creation requests to fail, breaking the entire deployment pipeline.

### 9.2 Root Cause Analysis

```
Failure Cascade:
1. Security webhook connects to policy database
2. Database connection pool limited to 10 connections
3. Deployment pipeline created 50+ pods simultaneously
4. Webhook made 50 concurrent database queries
5. Connection pool exhausted, queries queued
6. Webhook requests timed out after 30s
7. API server rejected all pod creation
8. Deployments failed with "webhook timeout"
9. Rolling updates couldn't create new pods
```

### 9.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 16:00:00 | Deployment pipeline triggered |
| 16:01:00 | Pod creation requests begin |
| 16:01:30 | Webhook timeouts reported |
| 16:02:00 | All pod creation failing |
| 16:05:00 | Pipeline failures escalating |
| 16:15:00 | Root cause identified (webhook DB) |
| 16:30:00 | Webhook bypassed (emergency) |
| 16:45:00 | Database pool increased |
| 17:30:00 | Webhook restored with fixes |

### 9.4 Resolution Steps

```yaml
# Webhook with timeout and failure policy
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pod-security-webhook
webhooks:
- name: pod-security.webhook.io
  clientConfig:
    service:
      name: security-webhook
      namespace: security
      path: "/validate"
    caBundle: ${CA_BUNDLE}
  rules:
  - operations: ["CREATE", "UPDATE"]
    apiGroups: [""]
    apiVersions: ["v1"]
    resources: ["pods"]
  failurePolicy: Ignore  # Fail open on webhook failure
  timeoutSeconds: 5      # Reduced from 30
  namespaceSelector:
    matchLabels:
      security-validation: enabled
```

```go
// Webhook with caching
type SecurityWebhook struct {
    cache      *ristretto.Cache
    db         *sql.DB
    timeout    time.Duration
}

func (w *SecurityWebhook) Validate(ctx context.Context, pod *v1.Pod) error {
    cacheKey := fmt.Sprintf("policy:%s:%s", pod.Namespace, pod.Spec.ServiceAccountName)
    
    // Check cache first
    if _, found := w.cache.Get(cacheKey); found {
        return nil
    }
    
    // Query with timeout
    ctx, cancel := context.WithTimeout(ctx, w.timeout)
    defer cancel()
    
    policy, err := w.getPolicy(ctx, pod.Namespace)
    if err != nil {
        // Log but don't fail - fail open
        log.Printf("Policy check failed, allowing: %v", err)
        return nil
    }
    
    // Cache result
    w.cache.Set(cacheKey, policy, 5*time.Minute)
    
    return w.validatePod(pod, policy)
}
```

### 9.5 Lessons Learned

1. **Webhooks must fail open** for availability
2. **Short timeouts** prevent cascading failures
3. **Caching** reduces external dependencies
4. **Namespace selectors** limit webhook scope

### 9.6 Prevention Recommendations

```yaml
# Webhook monitoring
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: webhook-alerts
spec:
  groups:
  - name: webhook
    rules:
    - alert: WebhookHighLatency
      expr: |
        histogram_quantile(0.99, 
          sum(rate(apiserver_admission_webhook_admission_duration_seconds_bucket[5m])) by (name, le)
        ) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Webhook {{ $labels.name }} has high latency"
        
    - alert: WebhookErrors
      expr: |
        sum(rate(apiserver_admission_webhook_rejection_count[5m])) by (name) > 0.1
      for: 2m
      labels:
        severity: critical
```

---

## Case Study 10: Container Runtime Socket Exhaustion

### 10.1 Incident Description

**System**: Kubernetes nodes with containerd runtime
**Impact**: Pods stuck in ContainerCreating, node NotReady
**Duration**: 2 hours
**Date**: June 2023

A memory leak in containerd caused the shim process count to grow unbounded. Each shim held open a socket to the containerd daemon. When the socket limit was reached, no new containers could be created, and nodes became NotReady.

### 10.2 Root Cause Analysis

```
Resource Leak:
1. Containerd shim v1 used for pods
2. Memory leak in shim process (50MB/day)
3. Each shim holds socket to containerd.sock
4. Node limit: 65536 open files per process
5. After 3 weeks: 1000+ shims running
6. Socket limit approached
7. New pods couldn't create containers
8. Kubelet marked node NotReady
9. Pods evicted, rescheduled to other nodes
10. Cascade to other nodes
```

### 10.3 Timeline of Events

| Time (UTC) | Event | Shim Count |
|------------|-------|------------|
| Week -3 | Node provisioned | 50 |
| Week -2 | Shims accumulating | 400 |
| Week -1 | Memory pressure rising | 800 |
| 08:00:00 | First pod creation failures | 1000 |
| 08:30:00 | Node marked NotReady | 1000+ |
| 09:00:00 | Pods evicted | - |
| 09:30:00 | Other nodes showing same pattern | - |
| 10:00:00 | Rolling restart initiated | - |
| 11:00:00 | Containerd upgraded to shim v2 | 50 |

### 10.4 Resolution Steps

```bash
# Check shim count
ps aux | grep containerd-shim | wc -l

# Monitor containerd sockets
ss -x | grep containerd | wc -l

# Restart containerd (drain node first)
kubectl drain <node> --ignore-daemonsets
systemctl restart containerd
kubectl uncordon <node>

# Upgrade to shim v2 (memory-efficient)
# /etc/containerd/config.toml
[plugins.linux]
  shim = "containerd-shim-runc-v2"
```

```go
// Node monitoring for shim count
type ContainerdMonitor struct {
    client metrics.Client
}

func (m *ContainerdMonitor) checkShimCount(node string) error {
    // Query node-exporter for process count
    query := `count(containerd_shim_processes{node="%s"})`
    
    value, err := m.client.Query(fmt.Sprintf(query, node))
    if err != nil {
        return err
    }
    
    shimCount := int(value)
    if shimCount > 500 {
        alert("High shim count on %s: %d", node, shimCount)
    }
    
    if shimCount > 800 {
        // Trigger preventive restart
        return m.drainAndRestart(node)
    }
    
    return nil
}
```

### 10.5 Lessons Learned

1. **Container runtime metrics** need monitoring
2. **Shim v2** is more memory-efficient than v1
3. **Process limits** should trigger alerts
4. **Regular node rotation** prevents accumulation

### 10.6 Prevention Recommendations

```yaml
# Node problem detector configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-problem-detector-config
data:
  shim-abort.config: |
    {
      "plugin": "custom",
      "pluginConfig": {
        "invoke_interval": "5m",
        "timeout": "1m",
        "max_output_length": "80",
        "concurrency": 1,
        "enable_message_change_based_condition_update": false
      },
      "source": "shim-monitor",
      "metricsReporting": true,
      "conditions": [
        {
          "type": "ShimCountHigh",
          "reason": "ShimCountNormal",
          "message": "shim count is within acceptable limits"
        }
      ],
      "rules": [
        {
          "type": "permanent",
          "condition": "ShimCountHigh",
          "reason": "ShimCountHigh",
          "path": "/home/kubernetes/bin/shim-check.sh"
        }
      ]
    }
```

---

## Summary and Best Practices

### Common Failure Patterns

| Pattern | Frequency | Impact | Detectability |
|---------|-----------|--------|---------------|
| Control Plane Overload | High | Critical | Medium |
| Registry Rate Limiting | Medium | High | High |
| Circuit Breaker Cascade | Medium | Critical | Medium |
| HPA Thrashing | Medium | Medium | High |
| etcd Backup Issues | Low | Critical | Low |
| DNS Storm | Medium | High | Medium |
| PV Data Loss | Low | Critical | Low |
| Quota Exhaustion | Medium | High | Medium |
| Webhook Timeout | Medium | Critical | Medium |
| Runtime Socket Leak | Low | High | Low |

### Prevention Checklist

- [ ] Configure conservative CronJob concurrency policies
- [ ] Use private registries with proper rate limits
- [ ] Implement gradual circuit breaker ejection
- [ ] Add HPA stabilization windows
- [ ] Verify backup integrity regularly
- [ ] Configure CoreDNS caching and HPA
- [ ] Use Retain reclaim policy for critical data
- [ ] Separate quotas for different workload types
- [ ] Set webhooks to fail open with short timeouts
- [ ] Monitor container runtime resource usage

### References

1. "Kubernetes Failure Stories" - Kubernetes Community
2. "Site Reliability Engineering" - Google
3. "Cloud Native Infrastructure" - Garrison et al.
4. "The Children's Illustrated Guide to Kubernetes" - CNCF
5. Kubernetes Official Documentation - Production Best Practices

---

*Document Size: 15+ KB | Level: S | Last Updated: 2026-04-03*
