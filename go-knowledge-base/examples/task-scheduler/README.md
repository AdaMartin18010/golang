# 分布式任务调度器示例

完整的分布式任务调度系统，包含领导选举、工作池、任务分发和故障恢复。

## 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Task Scheduler                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Scheduler Cluster                              │   │
│  │  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐           │   │
│  │  │ Scheduler-1 │     │ Scheduler-2 │     │ Scheduler-N │           │   │
│  │  │  (Leader)   │     │  (Follower) │     │  (Follower) │           │   │
│  │  │             │     │             │     │             │           │   │
│  │  │ - Elect    │◄────►│ - Elect    │◄────►│ - Elect    │           │   │
│  │  │ - Dispatch │      │ - Standby  │      │ - Standby  │           │   │
│  │  │ - Monitor  │      │            │      │            │           │   │
│  │  └──────┬──────┘     └─────────────┘      └─────────────┘           │   │
│  │         │                                                           │   │
│  │         │ Distribute Tasks                                          │   │
│  │         ▼                                                           │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                      Worker Pool                             │    │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │    │   │
│  │  │  │ Worker-1│  │ Worker-2│  │ Worker-3│  │ Worker-N│        │    │   │
│  │  │  │         │  │         │  │         │  │         │        │    │   │
│  │  │  │ Execute │  │ Execute │  │ Execute │  │ Execute │        │    │   │
│  │  │  │  Tasks  │  │  Tasks  │  │  Tasks  │  │  Tasks  │        │    │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘        │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│                              │ Coordination                                  │
│                              ▼                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Coordination Layer (etcd/Redis)                 │   │
│  │  - Leader Election                                                  │   │
│  │  - Service Discovery                                                │   │
│  │  - Distributed Locks                                                │   │
│  │  - Task Queue                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│                              │ Persist                                       │
│                              ▼                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Data Layer (PostgreSQL)                         │   │
│  │  - Task Definitions                                                 │   │
│  │  - Execution History                                                │   │
│  │  - Metrics                                                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 功能特性

- **领导选举**: 基于 etcd 的分布式领导选举
- **任务调度**: 支持多种调度策略 (轮询/最少任务/优先级)
- **故障恢复**: 自动检测工作节点故障并重分配任务
- **水平扩展**: 调度器和工作节点均可水平扩展
- **优雅关闭**: 支持优雅关闭和任务状态持久化

## 快速开始

```bash
# 启动 etcd
docker run -d --name etcd \
  -p 2379:2379 \
  quay.io/coreos/etcd:v3.5.0 \
  etcd --advertise-client-urls http://0.0.0.0:2379 \
       --listen-client-urls http://0.0.0.0:2379

# 启动 PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_DB=scheduler \
  -e POSTGRES_USER=scheduler \
  -e POSTGRES_PASSWORD=scheduler \
  -p 5432:5432 \
  postgres:18

# 启动调度器
cd cmd/scheduler
go run main.go -config ../../configs/scheduler.yaml

# 启动工作节点
cd cmd/worker
go run main.go -id worker-1
```

## 配置

```yaml
# configs/scheduler.yaml
scheduler:
  strategy: "least-tasks"  # round-robin | least-tasks | priority
  batch_size: 100
  check_interval: 5s

etcd:
  endpoints:
    - "localhost:2379"

postgres:
  dsn: "postgres://scheduler:scheduler@localhost:5432/scheduler?sslmode=disable"

redis:
  addr: "localhost:6379"
```

## API

```bash
# 提交任务
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {"to": "user@example.com", "subject": "Hello"},
    "priority": 5,
    "scheduled_at": "2026-04-02T10:00:00Z"
  }'

# 查询任务状态
curl http://localhost:8080/tasks/{task_id}

# 获取调度器状态
curl http://localhost:8080/status
```

## 生产建议

1. **高可用**: 部署 3+ 调度器节点，etcd 使用 3/5 节点集群
2. **监控**: 集成 Prometheus + Grafana
3. **告警**: 任务失败率、队列堆积、节点失联
4. **备份**: 定期备份 PostgreSQL 任务数据
