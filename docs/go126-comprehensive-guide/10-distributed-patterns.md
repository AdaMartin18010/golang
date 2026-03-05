# 分布式系统模式

> 一致性理论与分布式算法的Go实现

---

## 一、分布式系统的理论基础

### 1.1 CAP定理的形式化

```
CAP定理 (Brewer):
────────────────────────────────────────
分布式系统最多同时满足以下两项:
├─ C (Consistency): 一致性
│   所有节点在同一时间看到相同数据
│   形式化: ∀t: read(t) → write(t₀) where t₀ ≤ t
│
├─ A (Availability): 可用性
│   每个请求都能在有限时间内获得响应
│   形式化: ∀request: P(response) = 1
│
└─ P (Partition tolerance): 分区容错
    网络分区时系统仍能运行
    形式化: partition → system_continues

不可能三角:
    C ──────── A
     ╲       ╱
      ╲   ╱
        P (必须选择)

CP系统: etcd, ZooKeeper, Consul
AP系统: Cassandra, DynamoDB, CouchDB
```

### 1.2 一致性模型谱系

```
一致性强度谱系 (强 → 弱):
────────────────────────────────────────

线性一致性 (Linearizable):
├─ 所有操作表现为在全局时间点原子执行
├─ 实时顺序保持
└─ 实现: 单领导者复制、共识算法

顺序一致性 (Sequential):
├─ 操作顺序与程序顺序一致
├─ 所有进程看到相同操作顺序
└─ 不保证实时性

因果一致性 (Causal):
├─ 因果相关的操作顺序保持一致
├─ 并发操作顺序任意
└─ 实现: 向量时钟

最终一致性 (Eventual):
├─ 无新更新时，副本最终一致
├─ 允许临时不一致
└─ 实现: 异步复制

弱一致性 (Weak):
└─ 无一致性保证
```

---

## 二、共识算法

### 2.1 Raft共识算法

```
Raft核心机制:
────────────────────────────────────────

角色状态机:
Follower ──► Candidate ──► Leader
    ▲                        │
    └────────────────────────┘
(心跳超时)              (选举成功)

安全性保证:
├─ 选举安全: 一个任期内最多一个Leader
├─ 领导者追加: Leader不会删除/修改已提交日志
├─ 日志匹配: 相同索引和任期的日志内容相同
├─ 领导者完备: 已提交日志在未来的Leader中
└─ 状态机安全: 相同索引应用相同命令

Go实现要点:
type Raft struct {
    mu        sync.Mutex
    state     State
    currentTerm int
    votedFor  int
    log       []LogEntry
    commitIndex int
    lastApplied int
    nextIndex  []int
    matchIndex []int
}
```

### 2.2 拜占庭容错 (BFT)

```
BFT问题:
────────────────────────────────────────
n个节点，f个拜占庭(恶意)节点
安全条件: n ≥ 3f + 1

PBFT算法:
├─ Request: 客户端发送请求
├─ Pre-prepare: Leader分配序号
├─ Prepare: 节点验证并广播
├─ Commit: 2f+1个prepare后提交
└─ Reply: 返回结果

Go实现考虑:
├─ 加密签名验证
├─ 消息认证码
└─ 视图更换机制
```

---

## 三、分布式事务

### 3.1 两阶段提交 (2PC)

```
2PC协议:
────────────────────────────────────────

阶段1 (投票):
Coordinator → Prepare? ──► Participant
Participant: 执行本地事务，锁定资源
Participant ── Vote YES/NO ──► Coordinator

阶段2 (提交/回滚):
若全部YES:
  Coordinator → Commit ──► Participants
  Participants: 提交本地事务，释放锁
  Participants ── ACK ──► Coordinator

若有NO:
  Coordinator → Rollback ──► All
  Participants: 回滚本地事务，释放锁

问题:
├─ 协调者单点故障
├─ 同步阻塞
└─ 脑裂风险

Go实现:
type Coordinator struct {
    participants []Participant
    timeout      time.Duration
}

func (c *Coordinator) Commit(ctx context.Context, txn Transaction) error {
    // Phase 1: Prepare
    votes := make(chan Vote, len(c.participants))
    for _, p := range c.participants {
        go func(p Participant) {
            vote, err := p.Prepare(ctx, txn)
            votes <- Vote{Participant: p, Vote: vote, Err: err}
        }(p)
    }

    // 收集投票
    for i := 0; i < len(c.participants); i++ {
        select {
        case v := <-votes:
            if v.Err != nil || v.Vote == No {
                return c.rollback(ctx)
            }
        case <-ctx.Done():
            return ctx.Err()
        }
    }

    // Phase 2: Commit
    return c.commit(ctx)
}
```

### 3.2 Saga模式

```
Saga长事务:
────────────────────────────────────────
将长事务拆分为本地事务序列
每个本地事务提交后立即释放资源
通过补偿事务处理失败

两种协调方式:
├─ 编排式 (Choreography): 事件驱动，服务自主决策
└─ 编排式 (Orchestration): 中央协调器指挥

Go实现 (编排式):
type Saga struct {
    steps []Step
    compensations []Compensation
}

func (s *Saga) Execute(ctx context.Context) error {
    completed := 0
    for i, step := range s.steps {
        if err := step.Execute(ctx); err != nil {
            // 补偿已完成的步骤
            for j := completed - 1; j >= 0; j-- {
                s.compensations[j].Execute(ctx)
            }
            return err
        }
        completed++
    }
    return nil
}
```

---

## 四、一致性协议实现

### 4.1 向量时钟

```
向量时钟原理:
────────────────────────────────────────
VC(a)[i] = 进程a观察到进程i的事件数

偏序关系:
VC₁ ≤ VC₂ ⟺ ∀i: VC₁[i] ≤ VC₂[i]
VC₁ < VC₂ ⟺ VC₁ ≤ VC₂ ∧ VC₁ ≠ VC₂
VC₁ ∥ VC₂ ⟺ ¬(VC₁ ≤ VC₂) ∧ ¬(VC₂ ≤ VC₁)

Go实现:
type VectorClock map[string]uint64

func (vc VectorClock) Increment(node string) {
    vc[node]++
}

func (vc VectorClock) Merge(other VectorClock) {
    for node, ts := range other {
        if ts > vc[node] {
            vc[node] = ts
        }
    }
}

func (vc VectorClock) HappensBefore(other VectorClock) bool {
    for node, ts := range other {
        if vc[node] > ts {
            return false
        }
    }
    return true
}
```

### 4.2 Gossip协议

```
Gossip协议 (流行病协议):
────────────────────────────────────────

传播模型:
├─ 反熵 (Anti-entropy): 全量同步，修复不一致
├─ 谣言传播 (Rumor mongering): 增量传播新信息
└─ 聚合 (Aggregation): 分布式计算统计值

Go实现要点:
├─ 成员列表维护
├─ 随机节点选择
├─ 收敛检测
└─ 消息去重

type Gossiper struct {
    members []string
    seen    map[string]time.Time
    mu      sync.RWMutex
}

func (g *Gossiper) Spread(msg Message) {
    g.mu.RLock()
    targets := selectRandom(g.members, fanout)
    g.mu.RUnlock()

    for _, target := range targets {
        go g.send(target, msg)
    }
}
```

---

## 五、容错与恢复

### 5.1 断路器模式

```
状态机模型:
────────────────────────────────────────
Closed ──(失败阈值)──► Open ──(超时)──► HalfOpen
  ▲                                           │
  └──────────(成功)───────────────────────────┘

Go实现:
type CircuitBreaker struct {
    state       State
    failures    int
    threshold   int
    resetTimeout time.Duration
    lastFailure time.Time
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.allow() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}
```

### 5.2 重试策略

```
重试算法:
────────────────────────────────────────

指数退避:
delay = min(base × 2^attempt + jitter, maxDelay)

Go实现:
func Retry(ctx context.Context, fn func() error, opts RetryOptions) error {
    delay := opts.BaseDelay
    for attempt := 0; attempt < opts.MaxAttempts; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }

        if !opts.Retryable(err) {
            return err
        }

        select {
        case <-time.After(delay):
            delay = min(delay*2, opts.MaxDelay)
            delay += time.Duration(rand.Int63n(int64(delay/4))) // jitter
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return ErrMaxRetriesExceeded
}
```

---

## 六、服务发现与负载均衡

### 6.1 服务发现模式

```
服务发现架构:
────────────────────────────────────────
┌─────────────┐         ┌─────────────┐
│   Client    │◄───────►│  Registry   │
└──────┬──────┘         └──────┬──────┘
       │                       │
       ▼                       ▼
┌─────────────┐         ┌─────────────┐
│  Service A  │◄───────►│ Health Check│
└─────────────┘         └─────────────┘

一致性模型:
├─ 强一致性: etcd, Consul (Raft)
└─ 最终一致性: Eureka, ZooKeeper

Go实现 (客户端发现):
type Resolver struct {
    registry Registry
    cache    map[string][]Endpoint
    mu       sync.RWMutex
}

func (r *Resolver) Resolve(service string) ([]Endpoint, error) {
    r.mu.RLock()
    endpoints := r.cache[service]
    r.mu.RUnlock()

    if len(endpoints) > 0 {
        return endpoints, nil
    }

    endpoints, err := r.registry.Lookup(service)
    if err != nil {
        return nil, err
    }

    r.mu.Lock()
    r.cache[service] = endpoints
    r.mu.Unlock()

    return endpoints, nil
}
```

### 6.2 负载均衡算法

```
负载均衡策略:
────────────────────────────────────────

轮询 (Round Robin):
next = (last + 1) mod n

加权轮询:
weight = current_weight + effective_weight
if weight > max_weight: weight = 0

最少连接:
selected = argmin(connections)

一致性哈希:
hash(key) → node
虚拟节点解决不均匀问题

Go实现:
type LoadBalancer struct {
    endpoints []Endpoint
    current   uint64
}

func (lb *LoadBalancer) Next() Endpoint {
    idx := atomic.AddUint64(&lb.current, 1) % uint64(len(lb.endpoints))
    return lb.endpoints[idx]
}
```

---

*本章涵盖分布式系统的核心理论与Go实现，为构建可靠分布式服务提供基础。*
