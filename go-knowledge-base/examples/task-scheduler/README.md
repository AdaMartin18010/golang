# 分布式任务调度器示例 (Distributed Task Scheduler)

> **维度**: 示例项目 (Example Project)
> **分类**: 分布式系统实现
> **难度**: 高级
> **技术栈**: Go, etcd, PostgreSQL, Redis
> **最后更新**: 2026-04-02

---

## 1. 问题陈述 (Problem Statement)

### 1.1 分布式任务调度的挑战

在分布式系统中，任务调度面临以下核心挑战：

| 挑战 | 描述 | 影响 |
|------|------|------|
| **单点故障** | 调度器宕机导致任务无法执行 | 系统不可用 |
| **任务重复** | 网络分区导致任务被多个节点执行 | 数据不一致 |
| **负载不均** | 部分节点过载，部分节点空闲 | 资源浪费 |
| **任务依赖** | 复杂任务间依赖关系管理 | 执行顺序错乱 |
| **故障恢复** | 执行中节点故障后的任务处理 | 任务丢失 |

### 1.2 设计目标

```
分布式任务调度器目标:
┌─────────────────────────────────────────────────────────┐
│  1. 高可用 (High Availability)                          │
│     → 调度器集群，单点故障自动切换                      │
├─────────────────────────────────────────────────────────┤
│  2.  Exactly-Once 语义                                  │
│     → 任务不丢失、不重复执行                            │
├─────────────────────────────────────────────────────────┤
│  3. 水平扩展 (Horizontal Scalability)                   │
│     → 调度器和执行器均可动态扩缩容                      │
├─────────────────────────────────────────────────────────┤
│  4. 灵活调度策略                                        │
│     → 支持多种调度算法和优先级                          │
├─────────────────────────────────────────────────────────┤
│  5. 可观测性 (Observability)                            │
│     → 任务状态跟踪、执行日志、监控告警                  │
└─────────────────────────────────────────────────────────┘
```

### 1.3 非功能性需求

| 需求 | 目标值 | 约束 |
|------|--------|------|
| 调度延迟 | P99 < 100ms | 从任务提交到开始执行 |
| 故障切换 | < 5s | Leader 失效到新的 Leader 选出 |
| 任务吞吐 | > 10,000 TPS | 集群总处理能力 |
| 数据持久化 | 零丢失 | 提交即持久化 |

---

## 2. 形式化方法 (Formal Approach)

### 2.1 分布式调度模型

```
系统组件形式化定义:

Scheduler Cluster S = {s₁, s₂, ..., sₙ}  // 调度器节点集合
Worker Cluster W = {w₁, w₂, ..., wₘ}     // 工作节点集合
Task Queue Q = [t₁, t₂, ..., tₖ]         // 待执行任务队列

Leader Election:
  - 任何时候最多一个 Leader: |{s ∈ S | s.role = LEADER}| ≤ 1
  - 必须有一个 Leader 才能调度: ∃s ∈ S: s.role = LEADER

Task State Machine:
  PENDING → SCHEDULED → RUNNING → COMPLETED
     ↓         ↓           ↓
   CANCELLED  FAILED    TIMEOUT

一致性保证:
  - 每个任务最多被一个 Worker 执行
  - 任务状态转换原子性
  - Leader 故障时未完成任务的重新调度
```

### 2.2 领导选举算法

**基于 etcd 的分布式领导选举**:

```
算法: etcd Leader Election

1. 每个候选者尝试创建一个带有 TTL 的租约 (Lease)
2. 将自身信息写入固定键 (如 /scheduler/leader) 并绑定租约
3. 创建成功则成为 Leader，失败则成为 Follower
4. Leader 定期续期租约，保持领导地位
5. 租约过期后，其他节点竞争成为新 Leader

形式化保证:
  - 安全性: 同一时刻最多一个 Leader (etcd CAS 保证)
  - 活性: 原 Leader 故障后，新 Leader 一定产生
  - 租约期限: TTL 期间内 Leader 有效
```

### 2.3 任务调度策略

```
调度策略形式化:

1. 轮询 (Round Robin)
   next_worker = (last_worker + 1) mod |W|

2. 最少任务优先 (Least Tasks)
   selected = argmin_{w ∈ W} w.task_count

3. 负载感知 (Load Aware)
   selected = argmin_{w ∈ W} (α·w.cpu + β·w.mem + γ·w.tasks)

4. 优先级队列 (Priority Queue)
   Q 按 task.priority 排序，高优先级先调度

5. 亲和性调度 (Affinity)
   selected = {w | w.tags ∩ t.tags ≠ ∅} 优先
```

---

## 3. 实现细节 (Implementation)

### 3.1 项目结构

```
task-scheduler/
├── cmd/
│   ├── scheduler/          # 调度器主程序
│   │   └── main.go
│   └── worker/             # 工作节点主程序
│       └── main.go
├── internal/
│   ├── api/                # HTTP/gRPC API
│   │   ├── handler.go
│   │   └── server.go
│   ├── scheduler/          # 调度核心
│   │   ├── scheduler.go
│   │   ├── strategy.go
│   │   └── leader.go
│   ├── worker/             # 任务执行
│   │   ├── worker.go
│   │   └── executor.go
│   ├── storage/            # 数据存储
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── etcd.go
│   ├── model/              # 数据模型
│   │   └── task.go
│   └── config/             # 配置管理
│       └── config.go
├── configs/                # 配置文件
│   └── scheduler.yaml
├── deployments/            # 部署配置
│   ├── docker-compose.yml
│   └── k8s/
└── README.md
```

### 3.2 核心调度器实现

```go
package scheduler

import (
    "context"
    "sync"
    "time"
)

// Scheduler 任务调度器
type Scheduler struct {
    id       string
    isLeader bool

    // 依赖组件
    storage  Storage
    etcd     EtcdClient
    strategy Strategy

    // 状态管理
    workers  map[string]*WorkerInfo
    tasks    map[string]*Task
    mu       sync.RWMutex

    // 控制信号
    ctx      context.Context
    cancel   context.CancelFunc
}

type WorkerInfo struct {
    ID         string
    Address    string
    Capacity   int
    TaskCount  int
    LastHeartbeat time.Time
    Healthy    bool
}

// NewScheduler 创建调度器实例
func NewScheduler(id string, storage Storage, etcd EtcdClient) *Scheduler {
    ctx, cancel := context.WithCancel(context.Background())
    return &Scheduler{
        id:       id,
        storage:  storage,
        etcd:     etcd,
        strategy: NewLeastTasksStrategy(),
        workers:  make(map[string]*WorkerInfo),
        tasks:    make(map[string]*Task),
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Start 启动调度器
func (s *Scheduler) Start() error {
    // 1. 启动领导选举
    go s.electionLoop()

    // 2. 启动调度循环 (仅 Leader 执行调度)
    go s.scheduleLoop()

    // 3. 启动 Worker 心跳监控
    go s.workerMonitorLoop()

    return nil
}

// electionLoop 领导选举循环
func (s *Scheduler) electionLoop() {
    for {
        select {
        case <-s.ctx.Done():
            return
        default:
            if !s.isLeader {
                // 尝试成为 Leader
                if err := s.campaignLeader(); err == nil {
                    s.becomeLeader()
                }
            }
            time.Sleep(5 * time.Second)
        }
    }
}

// scheduleLoop 任务调度循环 (仅 Leader)
func (s *Scheduler) scheduleLoop() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            if !s.isLeader {
                continue
            }
            s.dispatchTasks()
        }
    }
}

// dispatchTasks 分发任务
func (s *Scheduler) dispatchTasks() {
    // 1. 获取待处理任务
    pendingTasks, err := s.storage.GetPendingTasks(s.ctx, 100)
    if err != nil {
        log.Printf("Failed to get pending tasks: %v", err)
        return
    }

    for _, task := range pendingTasks {
        // 2. 选择合适的 Worker
        worker := s.strategy.SelectWorker(s.getHealthyWorkers(), task)
        if worker == nil {
            log.Printf("No available worker for task %s", task.ID)
            continue
        }

        // 3. 尝试分配任务
        if err := s.assignTask(task, worker); err != nil {
            log.Printf("Failed to assign task %s: %v", task.ID, err)
            continue
        }

        // 4. 更新任务状态
        task.Status = TaskStatusScheduled
        task.WorkerID = worker.ID
        task.ScheduledAt = time.Now()

        if err := s.storage.UpdateTask(s.ctx, task); err != nil {
            log.Printf("Failed to update task status: %v", err)
        }
    }
}

// assignTask 分配任务到 Worker
func (s *Scheduler) assignTask(task *Task, worker *WorkerInfo) error {
    // 通过 etcd 的分布式锁确保任务只被分配一次
    lockKey := fmt.Sprintf("/tasks/locks/%s", task.ID)

    // 尝试获取锁
    if err := s.etcd.AcquireLock(s.ctx, lockKey, 10*time.Second); err != nil {
        return fmt.Errorf("failed to acquire lock: %w", err)
    }
    defer s.etcd.ReleaseLock(lockKey)

    // 双重检查任务状态
    currentTask, err := s.storage.GetTask(s.ctx, task.ID)
    if err != nil {
        return err
    }
    if currentTask.Status != TaskStatusPending {
        return fmt.Errorf("task already processed: %s", currentTask.Status)
    }

    // 发送任务到 Worker
    return s.sendTaskToWorker(task, worker)
}
```

### 3.3 Worker 实现

```go
package worker

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// Worker 任务执行节点
type Worker struct {
    id       string
    address  string

    scheduler SchedulerClient
    executor  Executor

    // 任务执行管理
    tasks     map[string]*RunningTask
    mu        sync.RWMutex

    // 配置
    maxConcurrent int
    ctx           context.Context
    cancel        context.CancelFunc
}

type RunningTask struct {
    Task   *model.Task
    Cancel context.CancelFunc
    Done   chan struct{}
}

// NewWorker 创建工作节点
func NewWorker(id, address string, scheduler SchedulerClient) *Worker {
    ctx, cancel := context.WithCancel(context.Background())
    return &Worker{
        id:            id,
        address:       address,
        scheduler:     scheduler,
        executor:      NewDefaultExecutor(),
        tasks:         make(map[string]*RunningTask),
        maxConcurrent: 10,
        ctx:           ctx,
        cancel:        cancel,
    }
}

// Start 启动 Worker
func (w *Worker) Start() error {
    // 1. 注册到调度器
    if err := w.register(); err != nil {
        return fmt.Errorf("failed to register: %w", err)
    }

    // 2. 启动心跳
    go w.heartbeatLoop()

    // 3. 启动任务接收循环
    go w.taskReceiveLoop()

    return nil
}

// ExecuteTask 执行任务
func (w *Worker) ExecuteTask(task *model.Task) error {
    w.mu.Lock()
    if len(w.tasks) >= w.maxConcurrent {
        w.mu.Unlock()
        return fmt.Errorf("max concurrent tasks reached")
    }

    ctx, cancel := context.WithTimeout(w.ctx, task.Timeout)
    rt := &RunningTask{
        Task:   task,
        Cancel: cancel,
        Done:   make(chan struct{}),
    }
    w.tasks[task.ID] = rt
    w.mu.Unlock()

    // 异步执行任务
    go func() {
        defer close(rt.Done)
        defer func() {
            w.mu.Lock()
            delete(w.tasks, task.ID)
            w.mu.Unlock()
            cancel()
        }()

        // 更新任务状态为运行中
        w.scheduler.UpdateTaskStatus(task.ID, model.TaskStatusRunning, "")

        // 执行实际任务
        result, err := w.executor.Execute(ctx, task)

        // 更新任务结果
        if err != nil {
            w.scheduler.UpdateTaskStatus(task.ID, model.TaskStatusFailed, err.Error())
        } else {
            w.scheduler.UpdateTaskStatus(task.ID, model.TaskStatusCompleted, result)
        }
    }()

    return nil
}

// heartbeatLoop 心跳循环
func (w *Worker) heartbeatLoop() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-w.ctx.Done():
            return
        case <-ticker.C:
            w.mu.RLock()
            taskCount := len(w.tasks)
            w.mu.RUnlock()

            status := &WorkerStatus{
                ID:        w.id,
                Address:   w.address,
                TaskCount: taskCount,
                Timestamp: time.Now(),
            }

            if err := w.scheduler.Heartbeat(status); err != nil {
                log.Printf("Heartbeat failed: %v", err)
            }
        }
    }
}
```

### 3.4 任务模型

```go
package model

import (
    "encoding/json"
    "time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "PENDING"
    TaskStatusScheduled  TaskStatus = "SCHEDULED"
    TaskStatusRunning    TaskStatus = "RUNNING"
    TaskStatusCompleted  TaskStatus = "COMPLETED"
    TaskStatusFailed     TaskStatus = "FAILED"
    TaskStatusTimeout    TaskStatus = "TIMEOUT"
    TaskStatusCancelled  TaskStatus = "CANCELLED"
)

// Task 任务定义
type Task struct {
    ID          string                 `json:"id" db:"id"`
    Type        string                 `json:"type" db:"type"`
    Payload     json.RawMessage        `json:"payload" db:"payload"`
    Status      TaskStatus             `json:"status" db:"status"`
    Priority    int                    `json:"priority" db:"priority"`

    // 调度相关
    WorkerID    string                 `json:"worker_id,omitempty" db:"worker_id"`
    ScheduledAt *time.Time             `json:"scheduled_at,omitempty" db:"scheduled_at"`
    StartedAt   *time.Time             `json:"started_at,omitempty" db:"started_at"`
    CompletedAt *time.Time             `json:"completed_at,omitempty" db:"completed_at"`

    // 超时配置
    Timeout     time.Duration          `json:"timeout" db:"timeout"`

    // 重试配置
    MaxRetries  int                    `json:"max_retries" db:"max_retries"`
    RetryCount  int                    `json:"retry_count" db:"retry_count"`

    // 结果
    Result      string                 `json:"result,omitempty" db:"result"`
    Error       string                 `json:"error,omitempty" db:"error"`

    // 元数据
    Tags        []string               `json:"tags,omitempty" db:"tags"`
    CreatedAt   time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// TaskResult 任务执行结果
type TaskResult struct {
    TaskID    string    `json:"task_id"`
    Status    TaskStatus `json:"status"`
    Output    string    `json:"output,omitempty"`
    Error     string    `json:"error,omitempty"`
    Duration  time.Duration `json:"duration"`
}
```

### 3.5 配置示例

```yaml
# configs/scheduler.yaml
server:
  http:
    addr: ":8080"
  grpc:
    addr: ":9090"

scheduler:
  id: "scheduler-1"
  strategy: "least-tasks"  # round-robin | least-tasks | priority
  batch_size: 100
  check_interval: 100ms

  # 领导选举配置
  election:
    prefix: "/scheduler/leader"
    ttl: 10
    retry_interval: 5s

etcd:
  endpoints:
    - "localhost:2379"
  dial_timeout: 5s

postgres:
  host: "localhost"
  port: 5432
  database: "scheduler"
  username: "scheduler"
  password: "${DB_PASSWORD}"
  pool_size: 10

redis:
  addr: "localhost:6379"
  db: 0
  pool_size: 10

worker:
  max_concurrent: 10
  heartbeat_interval: 5s

metrics:
  enabled: true
  path: "/metrics"

logging:
  level: "info"
  format: "json"
```

---

## 4. 语义分析 (Semantic Analysis)

### 4.1 一致性模型

```
任务执行一致性保证:

1. Exactly-Once 执行语义
   - 使用 etcd 分布式锁防止重复调度
   - 状态机转换原子性 (数据库事务)
   - Worker 任务认领确认机制

2. Leader 选举一致性
   - etcd Raft 协议保证
   - 租约机制防止脑裂
   - 旧 Leader 任务接管

3. 故障恢复语义
   - Worker 故障: 任务超时后重新调度
   - Leader 故障: 新 Leader 扫描未完成任务
   - 网络分区: 租约过期触发重新选举
```

### 4.2 时序约束

```
任务调度时序:

1. 任务提交时序:
   Client ──Submit──► API ──Persist──► DB ──Ack──► Client

2. 调度时序 (Leader):
   Poll DB ──Select Worker──► Acquire Lock ──Dispatch──► Worker
                    │
                    └── Update Status ──► DB

3. 执行时序:
   Worker ──Receive──► Execute ──Report Status──► Scheduler ──Persist──► DB

时序约束:
  - 任务提交到调度: < 100ms (正常负载)
  - 任务执行超时: 可配置，默认 5m
  - 心跳间隔: 5s，超时 15s 视为失联
```

---

## 5. 权衡分析 (Trade-offs)

### 5.1 调度策略对比

| 策略 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| 轮询 | 简单，公平 | 不考虑负载 | 同构 Worker |
| 最少任务 | 负载均衡 | 启动延迟敏感 | 异构任务 |
| 负载感知 | 资源优化 | 需要监控 | 资源受限 |
| 优先级 | 关键任务优先 | 饥饿风险 | 多级 SLA |
| 亲和性 | 缓存友好 | 复杂度 | 数据局部性 |

### 5.2 存储选择

```
PostgreSQL vs Redis 任务队列:

PostgreSQL:
  优势: 强一致性，持久化，复杂查询
  劣势: 性能较低，扩展复杂
  适用: 任务元数据，状态持久化

Redis:
  优势: 高性能，原子操作
  劣势: 内存限制，持久化开销
  适用: 任务队列，分布式锁，缓存

本方案采用混合架构:
  - PostgreSQL: 任务定义，历史记录
  - Redis: 任务队列，Worker 状态缓存
  - etcd: 领导选举，服务发现
```

---

## 6. 视觉表示 (Visual Representations)

### 6.1 系统架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │   Web Console   │  │    API Client   │  │   CLI Tool      │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
└───────────┼────────────────────┼────────────────────┼──────────────────────┘
            │                    │                    │
            └────────────────────┼────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         API Gateway (Load Balancer)                          │
│                    ┌─────────────────────────────────────┐                  │
│                    │  - Authentication                   │                  │
│                    │  - Rate Limiting                    │                  │
│                    │  - Request Routing                  │                  │
│                    └─────────────────────────────────────┘                  │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Scheduler-1    │     │  Scheduler-2    │     │  Scheduler-3    │
│  ┌───────────┐  │     │  ┌───────────┐  │     │  ┌───────────┐  │
│  │  Leader   │  │     │  │ Follower  │  │     │  │ Follower  │  │
│  │           │  │     │  │           │  │     │  │           │  │
│  │ ┌───────┐ │  │     │  │ ┌───────┐ │  │     │  │ ┌───────┐ │  │
│  │ │Queue  │ │  │     │  │ │Queue  │ │  │     │  │ │Queue  │ │  │
│  │ │Poll   │ │  │     │  │ │(sync) │ │  │     │  │ │(sync) │ │  │
│  │ └───┬───┘ │  │     │  │ └───┬───┘ │  │     │  │ └───┬───┘ │  │
│  │     │     │  │     │  │     │     │  │     │  │     │     │  │
│  │ ┌───▼───┐ │  │     │  │ ┌───▼───┐ │  │     │  │ ┌───▼───┐ │  │
│  │ │Dispatch│ │  │     │  │ │ Stand │ │  │     │  │ │ Stand │ │  │
│  │ └───────┘ │  │     │  │ └───────┘ │  │     │  │ └───────┘ │  │
│  └───────────┘  │     │  └───────────┘  │     │  └───────────┘  │
└────────┬────────┘     └─────────────────┘     └─────────────────┘
         │
         │ etcd watch / distributed lock
         │
    ┌────┴────┬────────┬────────┬────────┐
    ▼         ▼        ▼        ▼        ▼
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐
│Worker-1│ │Worker-2│ │Worker-3│ │Worker-4│ │Worker-5│
│       │ │       │ │       │ │       │ │       │
│Execute│ │Execute│ │Execute│ │Execute│ │Execute│
│ Tasks │ │ Tasks │ │ Tasks │ │ Tasks │ │ Tasks │
└───┬───┘ └───┬───┘ └───┬───┘ └───┬───┘ └───┬───┘
    │         │         │         │         │
    └─────────┴─────────┴─────────┴─────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │        Data Layer             │
         │  ┌─────────────────────────┐  │
         │  │  PostgreSQL (Primary)   │  │
         │  │  - Task Definitions     │  │
         │  │  - Execution History    │  │
         │  └─────────────────────────┘  │
         │  ┌─────────────────────────┐  │
         │  │  etcd Cluster (3/5)     │  │
         │  │  - Leader Election      │  │
         │  │  - Distributed Locks    │  │
         │  │  - Service Discovery    │  │
         │  └─────────────────────────┘  │
         │  ┌─────────────────────────┐  │
         │  │  Redis Cluster          │  │
         │  │  - Task Queue           │  │
         │  │  - Worker Cache         │  │
         │  │  - Rate Limiting        │  │
         │  └─────────────────────────┘  │
         └───────────────────────────────┘
```

### 6.2 任务状态机

```
                         ┌─────────────┐
           ┌────────────►│   PENDING   │◄────────────┐
           │             │  (等待调度)  │             │
           │             └──────┬──────┘             │
           │                    │ Submit             │
           │                    │                    │
           │                    ▼                    │
      Retry│             ┌─────────────┐             │
      (if  │   ┌─────────┤  SCHEDULED  ├─────────┐   │
      retry)│   │         │  (已分配)    │         │   │
           │   │         └──────┬──────┘         │   │
           │   │                │                │   │
           ▼   │                ▼                │   │
    ┌──────────┴┐         ┌─────────────┐        │   │
    │  FAILED   │         │   RUNNING   ├────────┘   │
    │  (失败)    │◄────────┤  (执行中)    │  Timeout   │
    └─────┬─────┘  Error  └──────┬──────┘            │
          │                      │                   │
          │                      │ Success           │
          │                      ▼                   │
          │               ┌─────────────┐            │
          └──────────────►│  COMPLETED  ├────────────┘
                          │  (已完成)    │   Archive
                          └─────────────┘

          ┌─────────────┐
          │  CANCELLED  │◄── Cancel Request
          │  (已取消)    │
          └─────────────┘
```

---

## 7. 部署与运维

### 7.1 Docker Compose 部署

```yaml
# deployments/docker-compose.yml
version: '3.8'

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    command:
      - etcd
      - --advertise-client-urls=http://0.0.0.0:2379
      - --listen-client-urls=http://0.0.0.0:2379
    ports:
      - "2379:2379"

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: scheduler
      POSTGRES_USER: scheduler
      POSTGRES_PASSWORD: scheduler
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  scheduler-1:
    build: ../
    command: ["/app/scheduler", "-config", "/etc/scheduler/config.yaml"]
    environment:
      - SCHEDULER_ID=scheduler-1
      - DB_PASSWORD=scheduler
    volumes:
      - ../configs/scheduler.yaml:/etc/scheduler/config.yaml
    ports:
      - "8080:8080"
      - "9090:9090"
    depends_on:
      - etcd
      - postgres
      - redis

  scheduler-2:
    build: ../
    command: ["/app/scheduler", "-config", "/etc/scheduler/config.yaml"]
    environment:
      - SCHEDULER_ID=scheduler-2
      - DB_PASSWORD=scheduler
    volumes:
      - ../configs/scheduler.yaml:/etc/scheduler/config.yaml
    depends_on:
      - etcd
      - postgres
      - redis

  worker-1:
    build: ../
    command: ["/app/worker", "-id", "worker-1", "-scheduler", "scheduler-1:9090"]
    depends_on:
      - scheduler-1

  worker-2:
    build: ../
    command: ["/app/worker", "-id", "worker-2", "-scheduler", "scheduler-1:9090"]
    depends_on:
      - scheduler-1

volumes:
  postgres_data:
```

### 7.2 监控指标

```go
// Prometheus 指标定义
var (
    tasksScheduled = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "scheduler_tasks_scheduled_total",
            Help: "Total tasks scheduled",
        },
        []string{"type"},
    )

    taskDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "scheduler_task_duration_seconds",
            Help: "Task execution duration",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
        },
        []string{"type", "status"},
    )

    workerTasks = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "scheduler_worker_tasks",
            Help: "Current tasks per worker",
        },
        []string{"worker_id"},
    )

    leaderStatus = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "scheduler_leader",
            Help: "1 if this instance is leader",
        },
    )
)
```

---

## 8. 相关资源

### 8.1 内部文档

- [EC-017-Scheduled-Task-Framework.md](../../../03-Engineering-CloudNative/EC-017-Scheduled-Task-Framework.md)
- [EC-019-Task-Execution-Engine.md](../../../03-Engineering-CloudNative/EC-019-Task-Execution-Engine.md)

### 8.2 外部参考

- [Kubernetes CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
- [Temporal Workflow](https://temporal.io/)
- [AWS Step Functions](https://aws.amazon.com/step-functions/)

---

*S-Level Quality Document | Generated: 2026-04-02*
