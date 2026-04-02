# etcd 分布式任务调度器实现 (ETCD Distributed Task Scheduler)

> **分类**: 工程与云原生
> **标签**: #etcd #distributed-systems #task-scheduler #raft
> **参考**: etcd v3.5+, Kubernetes controller patterns, Raft consensus

---

## 架构概述

基于 etcd 的分布式任务调度器利用 etcd 的强一致性、Watch 机制和 Lease 功能实现高可用任务调度。

```
┌─────────────────────────────────────────────────────────────────┐
│                        etcd Cluster                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │   Node 1     │  │   Node 2     │  │   Node 3     │ (Raft)    │
│  │  (Leader)    │  │  (Follower)  │  │  (Follower)  │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│   Worker 1    │    │   Worker 2    │    │   Worker N    │
│  (Scheduler)  │    │  (Scheduler)  │    │  (Scheduler)  │
│               │    │               │    │               │
│ • Lease Mgmt  │    │ • Lease Mgmt  │    │ • Lease Mgmt  │
│ • Watch Tasks │    │ • Watch Tasks │    │ • Watch Tasks │
│ • Acquire Lock│    │ • Acquire Lock│    │ • Acquire Lock│
└───────────────┘    └───────────────┘    └───────────────┘
```

---

## etcd 数据模型设计

```go
// Key 结构设计
const (
    // /tasks/{taskID} - 任务元数据
    KeyTaskPrefix = "/tasks/"

    // /locks/{taskID} - 分布式锁
    KeyLockPrefix = "/locks/"

    // /nodes/{nodeID} - 节点心跳
    KeyNodePrefix = "/nodes/"

    // /assignments/{taskID} - 任务分配
    KeyAssignmentPrefix = "/assignments/"

    // /history/{date}/{taskID} - 执行历史
    KeyHistoryPrefix = "/history/"
)

type Task struct {
    ID          string            `json:"id"`
    Type        string            `json:"type"`
    Payload     []byte            `json:"payload"`
    Status      TaskStatus        `json:"status"`
    Priority    int               `json:"priority"`
    ScheduledAt time.Time         `json:"scheduled_at"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    NodeID      string            `json:"node_id,omitempty"`
    Metadata    map[string]string `json:"metadata"`
}

type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusScheduled  TaskStatus = "scheduled"
    TaskStatusRunning    TaskStatus = "running"
    TaskStatusCompleted  TaskStatus = "completed"
    TaskStatusFailed     TaskStatus = "failed"
    TaskStatusCancelled  TaskStatus = "cancelled"
)

// etcd 事务操作封装
type TaskStore struct {
    client *clientv3.Client
}

// CreateTask 使用乐观锁创建任务
func (s *TaskStore) CreateTask(ctx context.Context, task *Task) error {
    key := KeyTaskPrefix + task.ID
    value, _ := json.Marshal(task)

    // 使用事务确保幂等性
    txn := s.client.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
        Then(clientv3.OpPut(key, string(value))).
        Else(clientv3.OpGet(key))

    resp, err := txn.Commit()
    if err != nil {
        return err
    }

    if !resp.Succeeded {
        return fmt.Errorf("task %s already exists", task.ID)
    }

    return nil
}

// UpdateTaskStatus 使用 Compare-And-Swap 更新状态
func (s *TaskStore) UpdateTaskStatus(ctx context.Context, taskID string,
    oldStatus, newStatus TaskStatus, nodeID string) error {

    key := KeyTaskPrefix + taskID

    // 获取当前任务
    resp, err := s.client.Get(ctx, key)
    if err != nil {
        return err
    }

    if len(resp.Kvs) == 0 {
        return fmt.Errorf("task not found: %s", taskID)
    }

    var task Task
    if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
        return err
    }

    // 验证状态转换
    if task.Status != oldStatus {
        return fmt.Errorf("status mismatch: expected %s, got %s", oldStatus, task.Status)
    }

    // 更新状态
    task.Status = newStatus
    task.NodeID = nodeID
    task.UpdatedAt = time.Now()

    value, _ := json.Marshal(task)

    // 使用事务确保原子性
    txn := s.client.Txn(ctx).
        If(clientv3.Compare(clientv3.Value(key), "=", string(resp.Kvs[0].Value))).
        Then(clientv3.OpPut(key, string(value))).
        Else(clientv3.OpGet(key))

    txnResp, err := txn.Commit()
    if err != nil {
        return err
    }

    if !txnResp.Succeeded {
        return fmt.Errorf("task was modified by another process")
    }

    return nil
}
```

---

## 分布式锁实现

```go
// DistributedLock 基于 etcd Lease 的分布式锁
type DistributedLock struct {
    client     *clientv3.Client
    leaseID    clientv3.LeaseID
    key        string
    cancelFunc context.CancelFunc
}

// Acquire 获取锁 (非阻塞)
func (l *DistributedLock) Acquire(ctx context.Context, ttl int64) (bool, error) {
    // 创建 Lease
    leaseResp, err := l.client.Grant(ctx, ttl)
    if err != nil {
        return false, err
    }

    l.leaseID = leaseResp.ID

    // 使用事务原子创建锁
    txn := l.client.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(l.key), "=", 0)).
        Then(clientv3.OpPut(l.key, "", clientv3.WithLease(leaseResp.ID))).
        Else(clientv3.OpGet(l.key))

    resp, err := txn.Commit()
    if err != nil {
        return false, err
    }

    if !resp.Succeeded {
        // 锁已被占用，撤销 lease
        l.client.Revoke(ctx, leaseResp.ID)
        return false, nil
    }

    // 启动续约 goroutine
    ctx, cancel := context.WithCancel(context.Background())
    l.cancelFunc = cancel
    go l.keepAlive(ctx)

    return true, nil
}

// keepAlive 自动续约
func (l *DistributedLock) keepAlive(ctx context.Context) {
    ch, err := l.client.KeepAlive(ctx, l.leaseID)
    if err != nil {
        log.Printf("KeepAlive failed: %v", err)
        return
    }

    for {
        select {
        case <-ctx.Done():
            return
        case ka, ok := <-ch:
            if !ok {
                log.Println("KeepAlive channel closed")
                return
            }
            log.Printf("Lease renewed: %d, TTL: %d", ka.ID, ka.TTL)
        }
    }
}

// Release 释放锁
func (l *DistributedLock) Release(ctx context.Context) error {
    if l.cancelFunc != nil {
        l.cancelFunc()
    }

    // 删除锁键，触发立即释放
    _, err := l.client.Delete(ctx, l.key)
    if err != nil {
        return err
    }

    // 撤销 Lease
    _, err = l.client.Revoke(ctx, l.leaseID)
    return err
}
```

---

## 任务调度器实现

```go
// Scheduler 分布式任务调度器
type Scheduler struct {
    client   *clientv3.Client
    nodeID   string
    leaseTTL int64

    // 任务处理
    handlers map[string]TaskHandler

    // 生命周期管理
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

type TaskHandler func(ctx context.Context, task *Task) error

// Start 启动调度器
func (s *Scheduler) Start() error {
    s.ctx, s.cancel = context.WithCancel(context.Background())

    // 1. 注册节点心跳
    if err := s.registerNode(); err != nil {
        return err
    }

    // 2. 启动任务监听
    s.wg.Add(1)
    go s.watchTasks()

    // 3. 启动任务轮询（作为 Watch 的补偿）
    s.wg.Add(1)
    go s.pollTasks()

    // 4. 启动过期任务清理
    s.wg.Add(1)
    go s.cleanupStaleTasks()

    return nil
}

// watchTasks 使用 etcd Watch 监听任务变化
func (s *Scheduler) watchTasks() {
    defer s.wg.Done()

    watchChan := s.client.Watch(s.ctx, KeyTaskPrefix, clientv3.WithPrefix())

    for resp := range watchChan {
        if resp.Canceled {
            log.Println("Watch canceled, reconnecting...")
            time.Sleep(time.Second)
            watchChan = s.client.Watch(s.ctx, KeyTaskPrefix, clientv3.WithPrefix())
            continue
        }

        for _, ev := range resp.Events {
            if ev.Type == clientv3.EventTypePut {
                var task Task
                if err := json.Unmarshal(ev.Kv.Value, &task); err != nil {
                    continue
                }

                // 只处理待执行的任务
                if task.Status == TaskStatusPending {
                    go s.tryAcquireAndExecute(&task)
                }
            }
        }
    }
}

// tryAcquireAndExecute 尝试获取锁并执行任务
func (s *Scheduler) tryAcquireAndExecute(task *Task) {
    ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
    defer cancel()

    lock := &DistributedLock{
        client: s.client,
        key:    KeyLockPrefix + task.ID,
    }

    // 尝试获取锁
    acquired, err := lock.Acquire(ctx, 60)
    if err != nil || !acquired {
        return
    }
    defer lock.Release(context.Background())

    // 更新状态为运行中
    if err := s.store.UpdateTaskStatus(ctx, task.ID, TaskStatusPending,
        TaskStatusRunning, s.nodeID); err != nil {
        return
    }

    // 执行处理
    handler, ok := s.handlers[task.Type]
    if !ok {
        log.Printf("No handler for task type: %s", task.Type)
        return
    }

    // 创建带超时的上下文
    execCtx, execCancel := context.WithTimeout(s.ctx, time.Duration(task.Timeout)*time.Second)
    defer execCancel()

    // 执行并捕获结果
    err = handler(execCtx, task)

    // 更新最终状态
    if err != nil {
        s.store.UpdateTaskStatus(ctx, task.ID, TaskStatusRunning, TaskStatusFailed, s.nodeID)
        s.recordHistory(task, TaskStatusFailed, err.Error())
    } else {
        s.store.UpdateTaskStatus(ctx, task.ID, TaskStatusRunning, TaskStatusCompleted, s.nodeID)
        s.recordHistory(task, TaskStatusCompleted, "")
    }
}

// pollTasks 轮询待处理任务（补偿机制）
func (s *Scheduler) pollTasks() {
    defer s.wg.Done()

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            s.processPendingTasks()
        }
    }
}

func (s *Scheduler) processPendingTasks() {
    ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
    defer cancel()

    // 查询所有待处理任务
    resp, err := s.client.Get(ctx, KeyTaskPrefix,
        clientv3.WithPrefix(),
        clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
    if err != nil {
        log.Printf("Failed to get tasks: %v", err)
        return
    }

    for _, kv := range resp.Kvs {
        var task Task
        if err := json.Unmarshal(kv.Value, &task); err != nil {
            continue
        }

        if task.Status == TaskStatusPending {
            go s.tryAcquireAndExecute(&task)
        }
    }
}

// registerNode 注册节点并保持心跳
func (s *Scheduler) registerNode() error {
    ctx := context.Background()

    // 创建 Lease
    leaseResp, err := s.client.Grant(ctx, s.leaseTTL)
    if err != nil {
        return err
    }

    // 注册节点
    nodeKey := KeyNodePrefix + s.nodeID
    nodeInfo := map[string]interface{}{
        "id":         s.nodeID,
        "registered": time.Now(),
        "hostname":   getHostname(),
    }
    value, _ := json.Marshal(nodeInfo)

    _, err = s.client.Put(ctx, nodeKey, string(value), clientv3.WithLease(leaseResp.ID))
    if err != nil {
        return err
    }

    // 保持心跳
    go func() {
        ch, err := s.client.KeepAlive(s.ctx, leaseResp.ID)
        if err != nil {
            log.Printf("KeepAlive failed: %v", err)
            return
        }

        for range ch {
            // 心跳续约
        }
    }()

    return nil
}

// cleanupStaleTasks 清理过期任务
func (s *Scheduler) cleanupStaleTasks() {
    defer s.wg.Done()

    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            s.doCleanup()
        }
    }
}

func (s *Scheduler) doCleanup() {
    ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
    defer cancel()

    // 获取所有节点
    resp, _ := s.client.Get(ctx, KeyNodePrefix, clientv3.WithPrefix())

    activeNodes := make(map[string]bool)
    for _, kv := range resp.Kvs {
        nodeID := strings.TrimPrefix(string(kv.Key), KeyNodePrefix)
        activeNodes[nodeID] = true
    }

    // 检查运行中的任务
    resp, _ = s.client.Get(ctx, KeyTaskPrefix,
        clientv3.WithPrefix())

    for _, kv := range resp.Kvs {
        var task Task
        if err := json.Unmarshal(kv.Value, &task); err != nil {
            continue
        }

        // 如果任务标记为运行中但节点已失效，重置为待处理
        if task.Status == TaskStatusRunning && !activeNodes[task.NodeID] {
            log.Printf("Task %s node %s is inactive, resetting to pending", task.ID, task.NodeID)
            s.store.UpdateTaskStatus(ctx, task.ID, TaskStatusRunning, TaskStatusPending, "")
        }
    }
}
```

---

## Leader 选举实现

```go
// LeaderElection 基于 etcd 的 Leader 选举
type LeaderElection struct {
    client     *clientv3.Client
    election   *concurrency.Election
    session    *concurrency.Session

    nodeID     string
    isLeader   atomic.Bool
    onLeader   func()
    onFollower func()
}

// NewLeaderElection 创建选举实例
func NewLeaderElection(client *clientv3.Client, electionKey, nodeID string) (*LeaderElection, error) {
    // 创建会话
    s, err := concurrency.NewSession(client, concurrency.WithTTL(5))
    if err != nil {
        return nil, err
    }

    // 创建选举
    e := concurrency.NewElection(s, electionKey)

    return &LeaderElection{
        client:   client,
        election: e,
        session:  s,
        nodeID:   nodeID,
    }, nil
}

// Start 开始竞选
func (le *LeaderElection) Start(ctx context.Context) error {
    // 尝试成为 Leader
    if err := le.election.Campaign(ctx, le.nodeID); err != nil {
        return err
    }

    le.isLeader.Store(true)

    if le.onLeader != nil {
        go le.onLeader()
    }

    // 监控 Leader 变化
    go le.observe(ctx)

    return nil
}

// observe 监控 Leader 状态
func (le *LeaderElection) observe(ctx context.Context) {
    ch := le.election.Observe(ctx)

    for resp := range ch {
        if len(resp.Kvs) > 0 {
            leaderID := string(resp.Kvs[0].Value)
            isLeader := leaderID == le.nodeID

            wasLeader := le.isLeader.Swap(isLeader)

            if !wasLeader && isLeader {
                // 成为 Leader
                log.Println("Became leader")
                if le.onLeader != nil {
                    go le.onLeader()
                }
            } else if wasLeader && !isLeader {
                // 失去 Leader
                log.Println("Lost leadership")
                if le.onFollower != nil {
                    go le.onFollower()
                }
            }
        }
    }
}

// Resign 主动放弃 Leader
func (le *LeaderElection) Resign(ctx context.Context) error {
    return le.election.Resign(ctx)
}

// IsLeader 检查是否是 Leader
func (le *LeaderElection) IsLeader() bool {
    return le.isLeader.Load()
}
```

---

## 性能优化

```go
// BatchOperations 批量操作优化
type BatchOperations struct {
    client *clientv3.Client
}

// BatchCreateTasks 批量创建任务
func (b *BatchOperations) BatchCreateTasks(ctx context.Context, tasks []*Task) error {
    const batchSize = 100

    for i := 0; i < len(tasks); i += batchSize {
        end := i + batchSize
        if end > len(tasks) {
            end = len(tasks)
        }

        batch := tasks[i:end]

        // 构建事务操作
        ops := make([]clientv3.Op, 0, len(batch))
        cmps := make([]clientv3.Cmp, 0, len(batch))

        for _, task := range batch {
            key := KeyTaskPrefix + task.ID
            value, _ := json.Marshal(task)

            cmps = append(cmps, clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
            ops = append(ops, clientv3.OpPut(key, string(value)))
        }

        // 执行批量事务
        txn := b.client.Txn(ctx).If(cmps...).Then(ops...)
        if _, err := txn.Commit(); err != nil {
            return err
        }
    }

    return nil
}

// CompactHistory 压缩历史数据
func (b *BatchOperations) CompactHistory(ctx context.Context, before time.Time) error {
    // 获取历史键
    key := KeyHistoryPrefix + before.Format("2006-01-02")

    resp, err := b.client.Get(ctx, key, clientv3.WithPrefix())
    if err != nil {
        return err
    }

    // 批量删除
    for _, kv := range resp.Kvs {
        _, err := b.client.Delete(ctx, string(kv.Key))
        if err != nil {
            log.Printf("Failed to delete %s: %v", kv.Key, err)
        }
    }

    // 执行压缩
    _, err = b.client.Compact(ctx, resp.Header.Revision)
    return err
}
```

---

## 监控与可观测性

```go
// SchedulerMetrics 调度器指标
type SchedulerMetrics struct {
    taskSubmitted    prometheus.Counter
    taskCompleted    prometheus.Counter
    taskFailed       prometheus.Counter
    taskDuration     prometheus.Histogram
    leaderStatus     prometheus.Gauge
    nodeCount        prometheus.Gauge
    etcdOperationDur prometheus.Histogram
}

func NewSchedulerMetrics() *SchedulerMetrics {
    return &SchedulerMetrics{
        taskSubmitted: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "scheduler_tasks_submitted_total",
            Help: "Total number of tasks submitted",
        }),
        taskDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "scheduler_task_duration_seconds",
            Help:    "Task execution duration",
            Buckets: prometheus.DefBuckets,
        }),
        leaderStatus: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "scheduler_leader",
            Help: "Whether this instance is the leader (1) or not (0)",
        }),
    }
}

// instrumentedClient 包装 etcd 客户端以收集指标
func instrumentedClient(client *clientv3.Client, metrics *SchedulerMetrics) *clientv3.Client {
    // 添加拦截器收集操作指标
    // 实际实现需要自定义 etcd 客户端配置
    return client
}
```
