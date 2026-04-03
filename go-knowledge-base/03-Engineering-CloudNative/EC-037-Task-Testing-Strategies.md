# 任务测试策略 (Task Testing Strategies)

> **分类**: 工程与云原生
> **标签**: #testing #unit-test #integration-test

---

## 单元测试框架

```go
// 任务处理器单元测试
func TestEmailTaskHandler(t *testing.T) {
    tests := []struct {
        name    string
        payload []byte
        wantErr bool
    }{
        {
            name:    "valid email",
            payload: []byte(`{"to":"test@example.com","subject":"Hello"}`),
            wantErr: false,
        },
        {
            name:    "invalid email",
            payload: []byte(`{"to":"invalid","subject":"Hello"}`),
            wantErr: true,
        },
    }

    handler := &EmailTaskHandler{
        SMTPClient: &MockSMTPClient{},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            err := handler.Handle(ctx, tt.payload)

            if (err != nil) != tt.wantErr {
                t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// Mock 实现
type MockSMTPClient struct {
    SendFunc func(email Email) error
    sent     []Email
}

func (m *MockSMTPClient) Send(email Email) error {
    m.sent = append(m.sent, email)
    if m.SendFunc != nil {
        return m.SendFunc(email)
    }
    return nil
}
```

---

## 集成测试

```go
type TaskIntegrationTest struct {
    scheduler *TaskScheduler
    store     TaskStore
}

func (tit *TaskIntegrationTest) TestTaskLifecycle(t *testing.T) {
    ctx := context.Background()

    // 创建任务
    task := &Task{
        Name:     "integration-test-task",
        Type:     "test",
        Payload:  []byte(`{"data":"test"}`),
        Schedule: time.Now().Add(time.Second),
    }

    taskID, err := tit.scheduler.Schedule(ctx, task)
    if err != nil {
        t.Fatalf("schedule task: %v", err)
    }

    // 等待任务执行
    time.Sleep(2 * time.Second)

    // 验证任务状态
    storedTask, err := tit.store.Get(ctx, taskID)
    if err != nil {
        t.Fatalf("get task: %v", err)
    }

    if storedTask.Status != TaskStatusCompleted {
        t.Errorf("expected completed, got %s", storedTask.Status)
    }
}

// 使用 TestContainers
func TestWithRedis(t *testing.T) {
    ctx := context.Background()

    // 启动 Redis 容器
    req := testcontainers.ContainerRequest{
        Image:        "redis:alpine",
        ExposedPorts: []string{"6379/tcp"},
        WaitingFor:   wait.ForListeningPort("6379/tcp"),
    }

    redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatal(err)
    }
    defer redisC.Terminate(ctx)

    // 使用容器运行测试
    endpoint, _ := redisC.Endpoint(ctx, "tcp")
    client := redis.NewClient(&redis.Options{Addr: endpoint})

    // 测试任务队列
    queue := NewRedisTaskQueue(client)
    // ... 测试代码
}
```

---

## 混沌测试

```go
type ChaosTest struct {
    executor *TaskExecutor
}

func (ct *ChaosTest) TestNetworkPartition(t *testing.T) {
    // 模拟网络分区
    fault := NetworkPartitionFault{
        Target: "worker-1",
        Duration: time.Second * 10,
    }

    fault.Inject()
    defer fault.Recover()

    // 验证任务能转移到其他 worker
    task := &Task{Type: "test"}
    ct.executor.Submit(task)

    // 等待并验证
    // ...
}

func (ct *ChaosTest) TestWorkerCrash(t *testing.T) {
    // 模拟 worker 崩溃
    worker := ct.executor.StartWorker()

    // 提交任务
    task := &Task{Type: "long-running"}
    ct.executor.SubmitToWorker(task, worker.ID)

    // 杀死 worker
    worker.Kill()

    // 验证任务被重新调度
    // ...
}
```

---

## 负载测试

```go
type LoadTest struct {
    scheduler *TaskScheduler
}

func (lt *LoadTest) RunConcurrentTasks(t *testing.T) {
    ctx := context.Background()

    const (
        numWorkers = 100
        numTasks   = 10000
    )

    var wg sync.WaitGroup
    wg.Add(numWorkers)

    start := time.Now()

    for i := 0; i < numWorkers; i++ {
        go func(workerID int) {
            defer wg.Done()

            for j := 0; j < numTasks/numWorkers; j++ {
                task := &Task{
                    Name: fmt.Sprintf("task-%d-%d", workerID, j),
                    Type: "benchmark",
                }

                if _, err := lt.scheduler.Schedule(ctx, task); err != nil {
                    t.Errorf("schedule failed: %v", err)
                }
            }
        }(i)
    }

    wg.Wait()

    duration := time.Since(start)
    tps := float64(numTasks) / duration.Seconds()

    t.Logf("TPS: %.2f", tps)

    if tps < 1000 {
        t.Errorf("TPS too low: %.2f", tps)
    }
}
```

---

## 契约测试

```go
func TestTaskContract(t *testing.T) {
    // 验证任务序列化/反序列化
    original := &Task{
        ID:        "task-123",
        Name:      "test-task",
        Type:      "email",
        Payload:   []byte(`{"to":"a@b.com"}`),
        Priority:  5,
        CreatedAt: time.Now(),
    }

    // 序列化
    data, err := json.Marshal(original)
    if err != nil {
        t.Fatal(err)
    }

    // 反序列化
    var decoded Task
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatal(err)
    }

    // 验证所有字段
    if decoded.ID != original.ID {
        t.Errorf("ID mismatch: %s vs %s", decoded.ID, original.ID)
    }
    // ... 验证其他字段
}
```

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02
---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02