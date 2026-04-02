# 分布式任务调度器架构 (Distributed Task Scheduler Architecture)

> **分类**: 工程与云原生
> **标签**: #distributed-systems #task-scheduler #architecture #consensus
> **参考**: Google Borg, Kubernetes Scheduler, Apache Mesos, HashiCorp Nomad

---

## 调度器架构演进

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Scheduler Architectures                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 单体调度器 (Monolithic)                                                  │
│  ┌─────────────────────────────────────────┐                                │
│  │          Global Scheduler               │                                │
│  │  ┌─────────┬─────────┬─────────┐        │                                │
│  │  │ Queue   │ State   │ Policy  │        │                                │
│  │  │ Manager │ Store   │ Engine  │        │                                │
│  │  └────┬────┴────┬────┴────┬────┘        │                                │
│  │       │         │         │             │                                │
│  │       └─────────┴─────────┘             │                                │
│  │                 │                       │                                │
│  └─────────────────┼──────────────────────┘                                 │
│                    │                                                        │
│       ┌────────────┼────────────┐                                           │
│       ▼            ▼            ▼                                           │
│    Node 1      Node 2      Node 3                                           │
│                                                                             │
│  2. 两层调度器 (Two-Level)                                                   │
│  ┌─────────────────┐    ┌─────────────────┐                                 │
│  │  Global Master  │    │  Local Agent    │                                 │
│  │  (Resource Offer)│───▶│  (Task Accept)  │                                │
│  └─────────────────┘    └─────────────────┘                                 │
│         │                       │                                           │
│         │    ┌──────────────────┼──────────────────┐                        │
│         │    ▼                  ▼                  ▼                        │
│         └──▶ Agent 1        Agent 2           Agent 3                       │
│                                                                              │
│  3. 共享状态调度器 (Shared-State)                                             │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                                        │
│  │Sched 1  │ │Sched 2  │ │Sched N  │  (Optimistic Concurrency)              │
│  │(Memory) │ │(Memory) │ │(Memory) │                                        │
│  └────┬────┘ └────┬────┘ └────┬────┘                                        │
│       │           │           │                                             │
│       └───────────┼───────────┘                                             │
│                   │                                                         │
│                   ▼                                                         │
│         ┌─────────────────┐                                                 │
│         │  Shared DB/etcd │  (Source of Truth)                              │
│         └─────────────────┘                                                 │
│                                                                             │
│  4. 分布式调度器 (Distributed)                                               │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                                  │
│  │Sched 1  │◀──▶│Sched 2  │◀──▶│Sched N  │  (Peer-to-Peer)                │
│  │(Shard A)│    │(Shard B)│    │(Shard C)│                                  │
│  └────┬────┘    └────┬────┘    └────┬────┘                                  │
│       │              │              │                                       │
│       └──────────────┼──────────────┘                                       │
│                      │                                                      │
│         ┌────────────┼────────────┐                                         │
│         ▼            ▼            ▼                                         │
│      Node 1       Node 2       Node 3                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心调度算法

```go
// 调度器核心结构
type Scheduler struct {
    // 调度策略
    strategy SchedulingStrategy

    // 状态管理
    state StateManager

    // 资源跟踪
    resourceTracker ResourceTracker

    // 任务队列
    pendingTasks TaskQueue

    // 节点管理
    nodeManager NodeManager
}

// SchedulingStrategy 调度策略接口
type SchedulingStrategy interface {
    // Filter 过滤不满足条件的节点
    Filter(task *Task, nodes []*Node) []*Node

    // Score 为节点打分
    Score(task *Task, node *Node) float64

    // Select 选择最优节点
    Select(task *Task, candidates []*Node) *Node
}

// DefaultStrategy 默认调度策略
func (s *Scheduler) DefaultStrategy() SchedulingStrategy {
    return &compositeStrategy{
        filters: []FilterFunc{
            ResourceFilter,
            AffinityFilter,
            TaintTolerationFilter,
            VolumeBindingFilter,
        },
        scorers: []ScorerFunc{
            LeastAllocatedScorer,
            BalancedResourceScorer,
            ImageLocalityScorer,
            NodeAffinityScorer,
        },
    }
}

// ResourceFilter 资源过滤器
func ResourceFilter(task *Task, node *Node) bool {
    // 检查 CPU
    if node.Allocatable.CPU - node.Used.CPU < task.Resources.CPU {
        return false
    }

    // 检查内存
    if node.Allocatable.Memory - node.Used.Memory < task.Resources.Memory {
        return false
    }

    // 检查 GPU
    if task.Resources.GPU > 0 {
        if node.Allocatable.GPU - node.Used.GPU < task.Resources.GPU {
            return false
        }
    }

    return true
}

// AffinityFilter 亲和性过滤器
func AffinityFilter(task *Task, node *Node) bool {
    // 检查节点亲和性
    if task.NodeAffinity != nil {
        if !matchNodeAffinity(task.NodeAffinity, node) {
            return false
        }
    }

    // 检查 Pod 亲和性
    if task.PodAffinity != nil {
        if !matchPodAffinity(task.PodAffinity, node) {
            return false
        }
    }

    // 检查反亲和性
    if task.PodAntiAffinity != nil {
        if !matchPodAntiAffinity(task.PodAntiAffinity, node) {
            return false
        }
    }

    return true
}

// LeastAllocatedScorer 最少资源分配打分器
// 优先选择资源使用率低的节点
func LeastAllocatedScorer(task *Task, node *Node) float64 {
    cpuFraction := float64(node.Used.CPU) / float64(node.Allocatable.CPU)
    memoryFraction := float64(node.Used.Memory) / float64(node.Allocatable.Memory)

    // 平均使用率越低，分数越高
    avgFraction := (cpuFraction + memoryFraction) / 2
    return (1 - avgFraction) * 100
}

// BalancedResourceScorer 资源均衡打分器
// 优先选择 CPU 和内存使用率均衡的节点
func BalancedResourceScorer(task *Task, node *Node) float64 {
    cpuFraction := float64(node.Used.CPU) / float64(node.Allocatable.CPU)
    memoryFraction := float64(node.Used.Memory) / float64(node.Allocatable.Memory)

    // 计算差异
    diff := math.Abs(cpuFraction - memoryFraction)

    // 差异越小，分数越高
    return (1 - diff) * 100
}

// Schedule 调度主流程
func (s *Scheduler) Schedule(ctx context.Context, task *Task) (*ScheduleResult, error) {
    // 1. 获取所有可用节点
    nodes, err := s.nodeManager.GetAvailableNodes()
    if err != nil {
        return nil, err
    }

    if len(nodes) == 0 {
        return nil, ErrNoNodesAvailable
    }

    // 2. 过滤阶段 (Filter)
    feasibleNodes := s.strategy.Filter(task, nodes)
    if len(feasibleNodes) == 0 {
        return nil, ErrNoFeasibleNodes
    }

    // 3. 打分阶段 (Score)
    if len(feasibleNodes) == 1 {
        return &ScheduleResult{Node: feasibleNodes[0]}, nil
    }

    // 为每个节点打分
    scores := make(map[string]float64)
    for _, node := range feasibleNodes {
        scores[node.ID] = s.strategy.Score(task, node)
    }

    // 4. 选择最优节点
    bestNode := s.selectBestNode(feasibleNodes, scores)

    // 5. 预占用资源 (Reserve)
    if err := s.resourceTracker.Reserve(task, bestNode); err != nil {
        return nil, err
    }

    return &ScheduleResult{
        Node: bestNode,
        Score: scores[bestNode.ID],
    }, nil
}

// selectBestNode 选择最优节点
func (s *Scheduler) selectBestNode(nodes []*Node, scores map[string]float64) *Node {
    var bestNode *Node
    var bestScore float64

    for _, node := range nodes {
        if score := scores[node.ID]; score > bestScore {
            bestScore = score
            bestNode = node
        }
    }

    return bestNode
}
```

---

## 资源分配策略

```go
// ResourceAllocation 资源分配
type ResourceAllocation struct {
    // 分配算法
    allocator ResourceAllocator

    // 资源池
    resourcePool *ResourcePool
}

type ResourceAllocator interface {
    Allocate(task *Task, nodes []*Node) (*Allocation, error)
    Release(allocation *Allocation) error
}

// DominantResourceFairness DRF 算法实现
// 确保多资源情况下的公平分配
type DominantResourceFairness struct {
    // 记录每个用户的资源使用量
    userResources map[string]*ResourceVector

    // 总资源
    totalResources *ResourceVector
}

type ResourceVector struct {
    CPU    float64
    Memory float64
    GPU    float64
}

func (drf *DominantResourceFairness) Allocate(task *Task, nodes []*Node) (*Allocation, error) {
    // 计算每个用户的 dominant share
    userDominantShares := make(map[string]float64)

    for userID, used := range drf.userResources {
        cpuShare := used.CPU / drf.totalResources.CPU
        memShare := used.Memory / drf.totalResources.Memory
        gpuShare := used.GPU / drf.totalResources.GPU

        // 取最大值作为 dominant share
        userDominantShares[userID] = max(cpuShare, memShare, gpuShare)
    }

    // 选择 dominant share 最小的用户优先分配
    // ... 排序逻辑

    // 分配资源
    for _, node := range nodes {
        if drf.canFit(task, node) {
            // 更新用户资源使用量
            drf.userResources[task.UserID].CPU += task.Resources.CPU
            drf.userResources[task.UserID].Memory += task.Resources.Memory

            return &Allocation{
                Node:      node,
                Resources: task.Resources,
            }, nil
        }
    }

    return nil, ErrInsufficientResources
}

// BinPacking 装箱算法
// 最大化资源利用率
func BinPackingAllocate(task *Task, nodes []*Node) (*Allocation, error) {
    // 按资源使用率排序
    sort.Slice(nodes, func(i, j int) bool {
        iUtilization := nodes[i].Used.CPU / nodes[i].Allocatable.CPU
        jUtilization := nodes[j].Used.CPU / nodes[j].Allocatable.CPU
        return iUtilization > jUtilization
    })

    // 优先选择已使用较多的节点
    for _, node := range nodes {
        if canFit(task, node) {
            return &Allocation{Node: node}, nil
        }
    }

    return nil, ErrNoFit
}

// Spread 分散算法
// 提高容错性
func SpreadAllocate(task *Task, nodes []*Node) (*Allocation, error) {
    // 按资源使用率排序（升序）
    sort.Slice(nodes, func(i, j int) bool {
        iUtilization := nodes[i].Used.CPU / nodes[i].Allocatable.CPU
        jUtilization := nodes[j].Used.CPU / nodes[j].Allocatable.CPU
        return iUtilization < jUtilization
    })

    // 优先选择已使用较少的节点
    for _, node := range nodes {
        if canFit(task, node) {
            return &Allocation{Node: node}, nil
        }
    }

    return nil, ErrNoFit
}
```

---

## 抢占与重调度

```go
// Preemption 抢占机制
type PreemptionPolicy struct {
    // 可抢占的任务优先级阈值
    victimPriorityThreshold int

    // 最大抢占任务数
    maxVictims int
}

func (p *PreemptionPolicy) Preempt(task *Task, node *Node) (*PreemptionResult, error) {
    // 获取节点上的低优先级任务
    victims := p.findVictims(task, node)

    if len(victims) == 0 {
        return nil, ErrNoVictims
    }

    // 验证抢占后是否能容纳新任务
    freedResources := p.calculateFreedResources(victims)
    if !p.canAccommodate(task, node, freedResources) {
        return nil, ErrInsufficientResourcesAfterPreemption
    }

    // 执行抢占
    for _, victim := range victims {
        if err := p.evictTask(victim); err != nil {
            return nil, err
        }
    }

    return &PreemptionResult{
        Victims: victims,
        Node:    node,
    }, nil
}

func (p *PreemptionPolicy) findVictims(task *Task, node *Node) []*Task {
    var victims []*Task

    for _, runningTask := range node.RunningTasks {
        // 只抢占优先级更低的任务
        if runningTask.Priority < task.Priority {
            victims = append(victims, runningTask)
        }
    }

    // 按优先级排序（优先抢占优先级最低的）
    sort.Slice(victims, func(i, j int) bool {
        return victims[i].Priority < victims[j].Priority
    })

    // 限制最大抢占数
    if len(victims) > p.maxVictims {
        victims = victims[:p.maxVictims]
    }

    return victims
}

// Rescheduler 重调度器
type Rescheduler struct {
    // 策略配置
    policy ReschedulePolicy

    // 指标收集
    metrics MetricsCollector
}

type ReschedulePolicy struct {
    // 节点利用率阈值
    lowUtilizationThreshold  float64
    highUtilizationThreshold float64

    // 检查间隔
    interval time.Duration
}

func (r *Rescheduler) Run(ctx context.Context) {
    ticker := time.NewTicker(r.policy.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.checkAndReschedule()
        }
    }
}

func (r *Rescheduler) checkAndReschedule() {
    // 1. 查找需要重调度的任务
    tasksToMove := r.findTasksToReschedule()

    // 2. 为每个任务寻找更好的节点
    for _, task := range tasksToMove {
        betterNode := r.findBetterNode(task)
        if betterNode != nil {
            // 执行迁移
            r.migrateTask(task, betterNode)
        }
    }
}

func (r *Rescheduler) findTasksToReschedule() []*Task {
    var tasks []*Task

    // 查找热点节点上的任务
    hotNodes := r.findHotNodes()
    for _, node := range hotNodes {
        // 找出可迁移的任务
        for _, task := range node.RunningTasks {
            if task.Migratable {
                tasks = append(tasks, task)
            }
        }
    }

    return tasks
}
```

---

## 一致性保证

```go
// StateManager 状态管理
type StateManager struct {
    // 使用 etcd 存储集群状态
    etcd *clientv3.Client

    // 本地缓存
    cache Cache

    // 版本控制
    resourceVersion int64
}

// OptimisticConcurrency 乐观并发控制
func (sm *StateManager) UpdateTask(task *Task) error {
    key := "/tasks/" + task.ID

    // 获取当前版本
    resp, err := sm.etcd.Get(context.Background(), key)
    if err != nil {
        return err
    }

    if len(resp.Kvs) == 0 {
        return ErrTaskNotFound
    }

    currentVersion := resp.Kvs[0].Version

    // 准备更新
    value, _ := json.Marshal(task)

    // 使用事务确保乐观并发
    txn := sm.etcd.Txn(context.Background()).
        If(clientv3.Compare(clientv3.Version(key), "=", currentVersion)).
        Then(clientv3.OpPut(key, string(value))).
        Else(clientv3.OpGet(key))

    txnResp, err := txn.Commit()
    if err != nil {
        return err
    }

    if !txnResp.Succeeded {
        return ErrConcurrentModification
    }

    return nil
}

// Watch 监听状态变化
func (sm *StateManager) Watch(ctx context.Context) {
    watchChan := sm.etcd.Watch(ctx, "/tasks/", clientv3.WithPrefix())

    for resp := range watchChan {
        for _, ev := range resp.Events {
            switch ev.Type {
            case clientv3.EventTypePut:
                sm.handleTaskUpdate(ev.Kv)
            case clientv3.EventTypeDelete:
                sm.handleTaskDelete(ev.Kv)
            }
        }
    }
}

// DistributedLock 分布式锁
type DistributedLock struct {
    etcd  *clientv3.Client
    lease clientv3.LeaseID
}

func (dl *DistributedLock) Acquire(ctx context.Context, key string, ttl int64) error {
    // 创建租约
    resp, err := dl.etcd.Grant(ctx, ttl)
    if err != nil {
        return err
    }

    dl.lease = resp.ID

    // 尝试获取锁
    txn := dl.etcd.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
        Then(clientv3.OpPut(key, "", clientv3.WithLease(dl.lease))).
        Else(clientv3.OpGet(key))

    txnResp, err := txn.Commit()
    if err != nil {
        return err
    }

    if !txnResp.Succeeded {
        return ErrLockAlreadyHeld
    }

    // 启动续约
    go dl.keepAlive(ctx)

    return nil
}

func (dl *DistributedLock) keepAlive(ctx context.Context) {
    ch, err := dl.etcd.KeepAlive(ctx, dl.lease)
    if err != nil {
        return
    }

    for range ch {
        // 续约成功
    }
}
```

---

## 高可用设计

```go
// LeaderElection 主备选举
type LeaderElection struct {
    etcd     *clientv3.Client
    election *concurrency.Election
    session  *concurrency.Session

    isLeader atomic.Bool
    callbacks LeaderCallbacks
}

type LeaderCallbacks struct {
    OnStartedLeading func()
    OnStoppedLeading func()
}

func (le *LeaderElection) Run(ctx context.Context) {
    for {
        // 竞选 Leader
        err := le.election.Campaign(ctx, le.identity)
        if err != nil {
            time.Sleep(time.Second)
            continue
        }

        // 成为 Leader
        le.isLeader.Store(true)

        if le.callbacks.OnStartedLeading != nil {
            go le.callbacks.OnStartedLeading()
        }

        // 监听 Leader 状态
        le.observe(ctx)

        // 失去 Leadership
        le.isLeader.Store(false)

        if le.callbacks.OnStoppedLeading != nil {
            le.callbacks.OnStoppedLeading()
        }
    }
}

// SchedulerHA 高可用调度器
type SchedulerHA struct {
    // Leader 选举
    leaderElection *LeaderElection

    // 主调度器
    primary *Scheduler

    // 备用调度器 (只读)
    standby *Scheduler
}

func (ha *SchedulerHA) Schedule(ctx context.Context, task *Task) (*ScheduleResult, error) {
    // 只有 Leader 才能执行调度
    if !ha.leaderElection.IsLeader() {
        return nil, ErrNotLeader
    }

    return ha.primary.Schedule(ctx, task)
}

// 故障转移
func (ha *SchedulerHA) failover() {
    // 1. 恢复状态
    ha.primary.state.Restore()

    // 2. 重新调度挂起的任务
    pendingTasks := ha.primary.pendingTasks.GetAll()
    for _, task := range pendingTasks {
        ha.primary.Schedule(context.Background(), task)
    }
}
```
