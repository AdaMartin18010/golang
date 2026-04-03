# EC-091-Kubernetes-Scheduler-Deep-Dive

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Kubernetes 1.35 (KEP Implementation)
> **Size**: >20KB
> **Source Reference**: k8s.io/kubernetes/pkg/scheduler

---

## 1. Scheduler Architecture

### 1.1 System Components

```
┌─────────────────────────────────────────┐
│         Kubernetes Scheduler            │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────┐    ┌──────────────┐   │
│  │ Informer    │───►│  Scheduling  │   │
│  │ (Watch API) │    │    Queue     │   │
│  └─────────────┘    └──────┬───────┘   │
│                            │            │
│                     ┌──────┴──────┐    │
│                     │   Pop()     │    │
│                     └──────┬──────┘    │
│                            │            │
│  ┌─────────────────────────┼─────────┐ │
│  │      Scheduling Cycle   │         │ │
│  │  ┌──────────────────────┘         │ │
│  │  │                               │ │
│  │  ▼                               │ │
│  │  ┌─────────────┐  ┌───────────┐  │ │
│  │  │  Pre-filter │─►│   Filter  │  │ │
│  │  │  (Snapshot) │  │  (Nodes)  │  │ │
│  │  └─────────────┘  └─────┬─────┘  │ │
│  │                         │        │ │
│  │  ┌──────────────────────┘        │ │
│  │  ▼                               │ │
│  │  ┌─────────────┐  ┌───────────┐  │ │
│  │  │   Score     │─►│   Select  │  │ │
│  │  │ (Normalize) │  │  (Winner) │  │ │
│  │  └─────────────┘  └───────────┘  │ │
│  └──────────────────────────────────┘ │
│                            │            │
│                            ▼            │
│                     ┌──────────────┐    │
│                     │   Bind()     │    │
│                     │ (API Server) │    │
│                     └──────────────┘    │
│                                         │
└─────────────────────────────────────────┘
```

### 1.2 Scheduling Framework

**Two-Phase Scheduling**:

1. **Scheduling Cycle**: Selects a node for the pod
   - Serial execution
   - Determines feasibility and scoring

2. **Binding Cycle**: Applies the decision to the cluster
   - Asynchronous
   - Can fail independently

---

## 2. Scheduling Queue

### 2.1 Queue Implementation

```go
// pkg/scheduler/internal/queue/scheduling_queue.go
type PriorityQueue struct {
    // ActiveQ: Pods ready for scheduling
    activeQ *heap.Heap

    // BackoffQ: Pods that failed scheduling, waiting to retry
    podBackoffQ *heap.Heap

    // UnschedulableQ: Pods that can't be scheduled
    unschedulableQ map[string]*framework.QueuedPodInfo

    // Plugin-specific queues (e.g., for gang scheduling)
    pluginWaitingPods map[string]*waitingPods
}
```

**Priority Calculation**:

```go
func (p *PriorityQueue) less(a, b interface{}) bool {
    pInfoA := a.(*framework.QueuedPodInfo)
    pInfoB := b.(*framework.QueuedPodInfo)

    // 1. Compare priority class
    if pInfoA.Pod.Spec.Priority != nil && pInfoB.Pod.Spec.Priority != nil {
        if *pInfoA.Pod.Spec.Priority != *pInfoB.Pod.Spec.Priority {
            return *pInfoA.Pod.Spec.Priority > *pInfoB.Pod.Spec.Priority
        }
    }

    // 2. Compare creation timestamp (FIFO for same priority)
    return pInfoA.Pod.CreationTimestamp.Before(&pInfoB.Pod.CreationTimestamp)
}
```

### 2.2 PodBackoff Mechanism

**Exponential Backoff** for failed scheduling attempts:

```go
const (
    initialBackoffDuration = 1 * time.Second
    maxBackoffDuration     = 10 * time.Second
)

func (p *PriorityQueue) calculateBackoffDuration(pod *v1.Pod) time.Duration {
    attempts := p.podAttempts[pod.Name]

    // Exponential backoff: 1s, 2s, 4s, 8s, 10s (cap)
    backoff := initialBackoffDuration * (1 << attempts)
    if backoff > maxBackoffDuration {
        backoff = maxBackoffDuration
    }

    return backoff
}
```

---

## 3. Filter Phase (Predicate)

### 3.1 Filter Plugins

**Node Fit Analysis** (`pkg/scheduler/framework/plugins/noderesources/`):

```go
// Filter nodes that have sufficient resources
func (f *Fit) Filter(ctx context.Context, state *framework.CycleState,
    pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {

    // Get pod resource requirements
    podRequests := calcResourceList(pod)

    // Check if node has enough allocatable resources
    if !hasSufficientResources(nodeInfo, podRequests) {
        return framework.NewStatus(framework.Unschedulable,
            "Insufficient resources")
    }

    // Check node conditions
    if nodeInfo.Node().Spec.Unschedulable {
        return framework.NewStatus(framework.Unschedulable,
            "Node is unschedulable")
    }

    return nil
}
```

### 3.2 Node Resources Filter

**Resource Calculation**:

```
Node Capacity:      64 CPU, 256GB RAM
- System Reserved:  -1 CPU, -4GB RAM
- Kubelet Eviction: -2 CPU, -8GB RAM
─────────────────────────────────────
Allocatable:        61 CPU, 244GB RAM

Allocated (existing pods): 40 CPU, 180GB RAM
─────────────────────────────────────
Available:          21 CPU, 64GB RAM

Pod Request:        4 CPU, 16GB RAM
Result:             FIT (21 >= 4 && 64 >= 16)
```

### 3.3 Topology Spread Constraints

**Pod Topology Spread** (`pkg/scheduler/framework/plugins/podtopologyspread/`):

```go
func (pl *PodTopologySpread) Filter(ctx context.Context,
    state *framework.CycleState, pod *v1.Pod,
    nodeInfo *framework.NodeInfo) *framework.Status {

    for _, constraint := range pod.Spec.TopologySpreadConstraints {
        // Count pods matching label selector in this topology domain
        matchCount := countMatchingPods(nodeInfo, constraint)

        // Check if adding this pod would violate maxSkew
        maxSkew := constraint.MaxSkew
        domainCount := getDomainCount(state, constraint.TopologyKey)

        // Skew = max(count in any domain) - min(count in any domain)
        // After adding: newSkew must <= maxSkew
        if wouldViolateSkew(matchCount, domainCount, maxSkew) {
            return framework.NewStatus(framework.Unschedulable,
                "Topology spread constraint violated")
        }
    }

    return nil
}
```

---

## 4. Score Phase (Priority)

### 4.1 Scoring Framework

**Normalized Scoring** (0-100 scale):

```go
// pkg/scheduler/framework/interface.go
type ScorePlugin interface {
    Name() string
    Score(ctx context.Context, state *framework.CycleState,
        pod *v1.Pod, nodeName string) (int64, *framework.Status)
}

// Score weights (configurable)
var defaultWeights = map[string]int64{
    "NodeResourcesBalancedAllocation": 1,
    "ImageLocality":                  1,
    "InterPodAffinity":               1,
    "NodeResourcesFit":               1,
    "PodTopologySpread":              2,
    "TaintToleration":                1,
}
```

### 4.2 Node Resources Balanced Allocation

**Algorithm**: Prefer nodes that would have more balanced resource usage after scheduling.

```go
func (b *BalancedAllocation) Score(ctx context.Context,
    state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

    nodeInfo, _ := b.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)

    // Calculate resource usage ratios after scheduling
    requestedCPU := pod.Spec.Containers.Resources.Requests.Cpu().MilliValue()
    requestedMem := pod.Spec.Containers.Resources.Requests.Memory().Value()

    allocatableCPU := nodeInfo.Allocatable.MilliCPU
    allocatableMem := nodeInfo.Allocatable.Memory

    // Usage ratios
    cpuRatio := float64(nodeInfo.Requested.MilliCPU + requestedCPU) / float64(allocatableCPU)
    memRatio := float64(nodeInfo.Requested.Memory + requestedMem) / float64(allocatableMem)

    // Score = 100 - variance from mean
    // Higher score = more balanced
    mean := (cpuRatio + memRatio) / 2
    variance := (math.Abs(cpuRatio-mean) + math.Abs(memRatio-mean)) / 2

    score := int64((1.0 - variance) * 100)

    return score, nil
}
```

### 4.3 Inter-Pod Affinity

```go
func (pl *InterPodAffinity) Score(ctx context.Context, state *framework.CycleState,
    pod *v1.Pod, nodeName string) (int64, *framework.Status) {

    nodeInfo, _ := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)

    score := int64(0)

    // Process pod affinity
    if pod.Spec.Affinity != nil && pod.Spec.Affinity.PodAffinity != nil {
        for _, term := range pod.Spec.Affinity.PodAffinity.Preferred {
            // Count matching pods on this node
            matchCount := countMatchingPodsOnNode(nodeInfo, term)
            score += int64(matchCount * term.Weight)
        }
    }

    // Process pod anti-affinity (negative score)
    if pod.Spec.Affinity != nil && pod.Spec.Affinity.PodAntiAffinity != nil {
        for _, term := range pod.Spec.Affinity.PodAntiAffinity.Preferred {
            matchCount := countMatchingPodsOnNode(nodeInfo, term)
            score -= int64(matchCount * term.Weight)
        }
    }

    // Normalize to 0-100
    return normalizeScore(score, maxPossibleScore), nil
}
```

---

## 5. Scheduling Algorithms

### 5.1 DefaultProvider Algorithm

```yaml
# pkg/scheduler/apis/config/scheme/defaults.go
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
  - schedulerName: default-scheduler
    plugins:
      queueSort:
        enabled:
          - name: PrioritySort
      preFilter:
        enabled:
          - name: NodeResourcesFit
          - name: NodePorts
          - name: PodTopologySpread
          - name: InterPodAffinity
          - name: VolumeBinding
      filter:
        enabled:
          - name: NodeUnschedulable
          - name: NodeName
          - name: TaintToleration
          - name: NodeAffinity
          - name: NodePorts
          - name: NodeResourcesFit
          - name: VolumeRestrictions
          - name: EBSLimits
          - name: GCEPDLimits
          - name: NodeVolumeLimits
          - name: AzureDiskLimits
          - name: VolumeBinding
          - name: VolumeZone
          - name: PodTopologySpread
          - name: InterPodAffinity
      preScore:
        enabled:
          - name: InterPodAffinity
          - name: PodTopologySpread
          - name: TaintToleration
      score:
        enabled:
          - name: NodeResourcesBalancedAllocation
            weight: 1
          - name: ImageLocality
            weight: 1
          - name: InterPodAffinity
            weight: 1
          - name: NodeResourcesFit
            weight: 1
          - name: PodTopologySpread
            weight: 2
          - name: TaintToleration
            weight: 1
```

### 5.2 Extender Support

For custom scheduling logic:

```go
// HTTP-based extender
type HTTPExtender struct {
    client      *http.Client
    filterURL   string
    prioritizeURL string
    bindURL     string
    weight      int64
}

func (h *HTTPExtender) Filter(pod *v1.Pod, nodes []*v1.Node) ([]*v1.Node, error) {
    req := &extender.FilterArgs{Pod: pod, Nodes: nodes}

    resp, err := h.post(h.filterURL, req)
    if err != nil {
        return nil, err
    }

    return resp.Nodes, nil
}
```

---

## 6. Performance Optimization

### 6.1 Snapshotting

**SnapshotSharedLister** creates a consistent view of cluster state:

```go
type Snapshot struct {
    nodeInfoList []*framework.NodeInfo
    nodeInfoMap  map[string]*framework.NodeInfo
    havePodsWithAffinityNodeInfoList []*framework.NodeInfo
    generation   int64
}

func (s *Snapshot) Refresh() {
    // Atomic snapshot of all nodes
    // Pre-computed for scheduling cycle
}
```

### 6.2 Parallel Scoring

```go
func (g *genericScheduler) prioritizeNodes(ctx context.Context,
    state *framework.CycleState, pod *v1.Pod, nodes []*framework.NodeInfo)
    (framework.NodeScoreList, error) {

    scores := make(framework.NodeScoreList, len(nodes))

    // Parallel scoring
    workqueue.ParallelizeUntil(ctx, 16, len(nodes), func(index int) {
        nodeName := nodes[index].Node().Name

        totalScore := int64(0)
        for _, plugin := range g.scorePlugins {
            score, _ := plugin.Score(ctx, state, pod, nodeName)
            totalScore += score * plugin.Weight()
        }

        scores[index] = framework.NodeScore{
            Name:  nodeName,
            Score: totalScore,
        }
    })

    return scores, nil
}
```

---

## 7. Gang Scheduling (KEP-4671, Alpha in 1.35)

### 7.1 PodGroup API

```yaml
apiVersion: scheduling.x-k8s.io/v1alpha1
kind: PodGroup
metadata:
  name: ml-training-job
spec:
  scheduleTimeoutSeconds: 300
  minMember: 4
  minResources:
    cpu: "32"
    memory: "128Gi"
    nvidia.com/gpu: "4"
```

### 7.2 Gang Scheduling Plugin

```go
// pkg/scheduler/framework/plugins/gang/gang.go
type GangScheduling struct {
    handle framework.Handle
    pgLister schedulinglisters.PodGroupLister
}

func (g *GangScheduling) PreFilter(ctx context.Context, state *framework.CycleState,
    pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {

    pgName := pod.Annotations[schedulingv1alpha1.PodGroupAnnotationKey]
    if pgName == "" {
        return nil, nil  // Not part of a pod group
    }

    pg, err := g.pgLister.PodGroups(pod.Namespace).Get(pgName)
    if err != nil {
        return nil, framework.NewStatus(framework.Error, err.Error())
    }

    // Check if all members are ready
    ready, err := g.checkAllMembersReady(pg)
    if err != nil {
        return nil, framework.NewStatus(framework.Unschedulable,
            "Waiting for all pod group members")
    }

    if !ready {
        // Permit phase: hold this pod
        g.handle.GetWaitingPod(pod.UID).Allow()
    }

    return nil, nil
}
```

---

## 8. References

1. **Kubernetes Scheduler Documentation**
   - <https://kubernetes.io/docs/concepts/scheduling-eviction/>

2. **Scheduler Framework KEP**
   - <https://github.com/kubernetes/enhancements/tree/master/keps/sig-scheduling/624-scheduling-framework>

3. **Pod Topology Spread KEP**
   - <https://github.com/kubernetes/enhancements/tree/master/keps/sig-scheduling/895-pod-topology-spread>

4. **Gang Scheduling KEP-4671**
   - <https://github.com/kubernetes/enhancements/issues/4671>

5. **Source Code**
   - <https://github.com/kubernetes/kubernetes/tree/master/pkg/scheduler>

---

*Last Updated: 2026-04-03*
*Source Version: Kubernetes 1.35*
