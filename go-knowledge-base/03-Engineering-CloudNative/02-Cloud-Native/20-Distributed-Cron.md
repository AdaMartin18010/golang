# 分布式 Cron (Distributed Cron)

> **分类**: 工程与云原生  
> **标签**: #distributed-cron #leader-election #scheduler

---

## 问题分析

```
单机 Cron 的问题:
1. 单点故障 - 节点宕机则任务无法执行
2. 重复执行 - 多节点部署会导致任务重复
3. 无高可用 - 无法保证任务至少执行一次
4. 难扩展 - 无法水平扩展处理大量定时任务
```

---

## Leader 选举机制

```go
type LeaderCron struct {
    nodeID      string
    store       ElectionStore
    isLeader    bool
    cron        *cron.Cron
    mu          sync.RWMutex
    
    onLeader    func()
    onFollower  func()
}

func (lc *LeaderCron) Start(ctx context.Context) {
    // 尝试成为 Leader
    go lc.electionLoop(ctx)
}

func (lc *LeaderCron) electionLoop(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if lc.isLeader {
                // 续租
                if err := lc.renewLease(ctx); err != nil {
                    lc.stepDown()
                }
            } else {
                // 尝试成为 Leader
                if err := lc.tryAcquireLeadership(ctx); err == nil {
                    lc.becomeLeader()
                }
            }
        case <-ctx.Done():
            if lc.isLeader {
                lc.releaseLeadership(ctx)
            }
            return
        }
    }
}

func (lc *LeaderCron) tryAcquireLeadership(ctx context.Context) error {
    lease := &LeaderLease{
        NodeID:    lc.nodeID,
        AcquiredAt: time.Now(),
        ExpiresAt:  time.Now().Add(10 * time.Second),
    }
    
    // CAS 操作确保只有一个节点能成为 Leader
    return lc.store.CompareAndSwap(ctx, "cron-leader", nil, lease)
}

func (lc *LeaderCron) becomeLeader() {
    lc.mu.Lock()
    lc.isLeader = true
    lc.mu.Unlock()
    
    log.Printf("Node %s became leader", lc.nodeID)
    
    // 启动 Cron
    lc.cron.Start()
    
    if lc.onLeader != nil {
        lc.onLeader()
    }
}

func (lc *LeaderCron) stepDown() {
    lc.mu.Lock()
    lc.isLeader = false
    lc.mu.Unlock()
    
    log.Printf("Node %s stepped down", lc.nodeID)
    
    // 停止 Cron
    lc.cron.Stop()
    
    if lc.onFollower != nil {
        lc.onFollower()
    }
}
```

---

## 基于 Redis 的选举

```go
type RedisElection struct {
    client *redis.Client
    nodeID string
    key    string
    ttl    time.Duration
}

func (re *RedisElection) Acquire(ctx context.Context) error {
    // SET key value NX PX milliseconds
    ok, err := re.client.SetNX(ctx, re.key, re.nodeID, re.ttl).Result()
    if err != nil {
        return err
    }
    if !ok {
        return ErrNotAcquired
    }
    return nil
}

func (re *RedisElection) Renew(ctx context.Context) error {
    // 检查是否仍是自己的锁
    val, err := re.client.Get(ctx, re.key).Result()
    if err != nil {
        return err
    }
    if val != re.nodeID {
        return ErrLostLeadership
    }
    
    // 续期
    return re.client.Expire(ctx, re.key, re.ttl).Err()
}

func (re *RedisElection) Release(ctx context.Context) error {
    // 使用 Lua 脚本确保原子性
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    return re.client.Eval(ctx, script, []string{re.key}, re.nodeID).Err()
}
```

---

## 任务分片

```go
type ShardedCron struct {
    nodeID    string
    nodes     []string  // 所有节点
    shardFunc func(string, int) int
}

func (sc *ShardedCron) ShouldExecute(taskID string) bool {
    shard := sc.shardFunc(taskID, len(sc.nodes))
    return sc.nodes[shard] == sc.nodeID
}

// 一致性哈希
func (sc *ShardedCron) AddNode(nodeID string) {
    sc.nodes = append(sc.nodes, nodeID)
    sort.Strings(sc.nodes)
}

func (sc *ShardedCron) RemoveNode(nodeID string) {
    for i, n := range sc.nodes {
        if n == nodeID {
            sc.nodes = append(sc.nodes[:i], sc.nodes[i+1:]...)
            break
        }
    }
}

// 使用
func (sc *ShardedCron) Schedule(taskID string, job func()) {
    if !sc.ShouldExecute(taskID) {
        return  // 不是本节点执行
    }
    
    // 添加到本节点 Cron
    cron.AddFunc(spec, job)
}
```

---

## 任务状态同步

```go
type DistributedJob struct {
    ID         string
    Name       string
    Spec       string
    LastRun    time.Time
    NextRun    time.Time
    LastNode   string    // 上次执行的节点
    Status     JobStatus
}

func (dc *DistributedCron) syncJobStatus(ctx context.Context) {
    // 从存储加载任务状态
    jobs, _ := dc.store.ListJobs(ctx)
    
    for _, job := range jobs {
        // 检查是否错过执行
        if time.Now().After(job.NextRun) && dc.isLeader {
            // Leader 负责补执行
            dc.executeMissedJob(ctx, job)
        }
    }
}

func (dc *DistributedCron) executeMissedJob(ctx context.Context, job *DistributedJob) {
    // 记录执行
    execution := &JobExecution{
        JobID:     job.ID,
        NodeID:    dc.nodeID,
        StartedAt: time.Now(),
        Type:      ExecutionTypeMissed,
    }
    
    dc.store.SaveExecution(ctx, execution)
    
    // 执行
    dc.executeJob(ctx, job)
}
```

---

## 多节点协调

```go
// 使用分布式锁确保任务不重复执行
type CoordinatedJob struct {
    dc *DistributedCron
}

func (cj *CoordinatedJob) Run() {
    ctx := context.Background()
    
    // 尝试获取任务锁
    lockKey := fmt.Sprintf("job-lock:%s", cj.ID)
    lock := cj.dc.store.NewLock(lockKey, 5*time.Minute)
    
    if err := lock.Acquire(ctx); err != nil {
        // 其他节点正在执行
        log.Printf("Job %s is being executed by another node", cj.ID)
        return
    }
    defer lock.Release(ctx)
    
    // 执行
    cj.execute(ctx)
}
```
