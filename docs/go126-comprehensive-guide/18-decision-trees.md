# 第十八章：决策树与设计权衡

> 系统化的决策框架和设计权衡分析

---

## 18.1 技术选型决策树

### 18.1.1 数据库访问技术选择

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    数据库访问技术决策树                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  需要 ORM 功能？                                                            │
│       │                                                                     │
│       ├── Yes ──▶ 需要代码生成？                                            │
│       │              │                                                      │
│       │              ├── Yes ──▶ Ent (类型安全，代码生成)                   │
│       │              │                                                      │
│       │              └── No ──▶ 需要复杂关联？                              │
│       │                         │                                           │
│       │                         ├── Yes ──▶ GORM (功能全面)                 │
│       │                         │                                           │
│       │                         └── No ──▶ GORM / sqlc                      │
│       │                                                                     │
│       └── No ──▶ 需要类型安全查询？                                         │
│                     │                                                       │
│                     ├── Yes ──▶ sqlx (扫描便利)                             │
│                     │                                                       │
│                     └── No ──▶ database/sql (标准库)                        │
│                                  │                                          │
│                                  └── 需要 PostgreSQL 特性？                  │
│                                       │                                     │
│                                       └── Yes ──▶ pgx (性能优先)            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 18.1.2 缓存策略选择

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      缓存策略决策树                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  数据需要跨服务共享？                                                        │
│       │                                                                     │
│       ├── Yes ──▶ 数据更新频率？                                            │
│       │              │                                                      │
│       │              ├── 高频更新 ──▶ Redis (支持过期、发布订阅)            │
│       │              │                                                      │
│       │              └── 低频更新 ──▶ Redis / Memcached                     │
│       │                                                                     │
│       └── No ──▶ 单机内存足够？                                             │
│                     │                                                       │
│                     ├── Yes ──▶ 读写频率？                                  │
│                     │              │                                        │
│                     │              ├── 读多写少 ──▶ bigcache (零 GC)        │
│                     │              │                                        │
│                     │              └── 读写均衡 ──▶ ristretto (高命中率)    │
│                     │                                                       │
│                     └── No ──▶ groupcache (分布式缓存)                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 18.1.3 通信协议选择

| 场景 | 推荐协议 | 理由 |
|------|----------|------|
| 内部服务通信 | gRPC | 高性能、强类型、流支持 |
| 浏览器客户端 | HTTP/REST | 广泛支持、易于调试 |
| 实时通信 | WebSocket | 全双工、低延迟 |
| 服务器推送 | SSE | 简单、自动重连 |
| 消息队列 | 异步消息 | 解耦、削峰填谷 |

```go
// 协议选择决策代码示例

// 高吞吐内部通信 - gRPC
func (s *Server) ProcessOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
    // gRPC 实现
}

// 外部 API - REST
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    // HTTP REST 实现
}

// 实时通知 - WebSocket
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    // WebSocket 实现
}
```

---

## 18.2 架构设计权衡矩阵

### 18.2.1 CAP 定理权衡

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                       CAP 定理权衡分析                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  系统类型              一致性      可用性      分区容错      适用场景          │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  传统关系型数据库       CP          -          -         金融交易             │
│  (PostgreSQL, MySQL)                                                        │
│                                                                             │
│  NoSQL 数据库          AP          ✓          ✓         社交应用            │
│  (Cassandra, DynamoDB)                                                      │
│                                                                             │
│  分布式协调服务         CP          -          ✓         配置管理            │
│  (etcd, ZooKeeper)                                                          │
│                                                                             │
│  缓存系统              AP          ✓          ✓         会话存储            │
│  (Redis Cluster)                                                            │
│                                                                             │
│  消息队列              AP          ✓          ✓         异步处理            │
│  (Kafka, RabbitMQ)                                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 18.2.2 一致性模型选择

```go
// 强一致性
func (s *Service) StrongConsistentWrite(ctx context.Context, data Data) error {
    // 同步复制到所有节点
    return s.writeQuorum(ctx, data, s.allNodes)
}

// 最终一致性
func (s *Service) EventualConsistentWrite(ctx context.Context, data Data) error {
    // 写入主节点，异步复制
    if err := s.writePrimary(ctx, data); err != nil {
        return err
    }
    go s.asyncReplicate(data)
    return nil
}

// 读写一致性
func (s *Service) ReadYourWrites(ctx context.Context, key string) (Data, error) {
    // 确保读到自己写的数据
    return s.readFromNode(ctx, key, s.getPreferredNode())
}
```

---

## 18.3 性能优化决策树

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    性能优化决策树                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  性能瓶颈类型？                                                             │
│       │                                                                     │
│       ├── CPU ──▶ 热点分析                                                  │
│       │              │                                                      │
│       │              ├── 算法复杂度高 ──▶ 算法优化                          │
│       │              │         ├── 使用合适的数据结构                       │
│       │              │         └── 降低时间复杂度                           │
│       │              │                                                      │
│       │              └── 计算密集型 ──▶ 并行化                              │
│       │                        ├── Goroutine 并发                          │
│       │                        └── 多核利用                                │
│       │                                                                     │
│       ├── 内存 ──▶ 内存分析                                                 │
│       │              │                                                      │
│       │              ├── 分配频繁 ──▶ 对象池                                 │
│       │              │              └── sync.Pool                          │
│       │              │                                                      │
│       │              ├── 逃逸分析 ──▶ 栈分配优化                            │
│       │              │              └── 避免接口、闭包捕获                  │
│       │              │                                                      │
│       │              └── GC 压力大 ──▶ 减少堆分配                           │
│       │                             └── 使用值类型                         │
│       │                                                                     │
│       ├── I/O ──▶ I/O 分析                                                  │
│       │              │                                                      │
│       │              ├── 网络 I/O ──▶ 连接池                                 │
│       │              │              ├── HTTP keep-alive                    │
│       │              │              └── 数据库连接池                       │
│       │              │                                                      │
│       │              ├── 磁盘 I/O ──▶ 异步/缓冲                             │
│       │              │              └── 批量读写                           │
│       │              │                                                      │
│       │              └── 数据库 ──▶ 查询优化                                │
│       │                             ├── 索引优化                           │
│       │                             └── 缓存层                             │
│       │                                                                     │
│       └── 延迟 ──▶ 调用链分析                                               │
│                      │                                                      │
│                      ├── 串行调用 ──▞ 并行化                                  │
│                      │                                                      │
│                      ├── 同步阻塞 ──▞ 异步化                                  │
│                      │                                                      │
│                      └── 重复计算 ──▞ 缓存结果                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 18.4 错误处理策略决策

```go
// 错误处理策略矩阵

// 1. 立即返回 - 快速失败
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    if err := validate(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("create user failed: %w", err)
    }

    return user, nil
}

// 2. 累积错误 - 批量操作
func (s *Service) BatchProcess(items []Item) (Result, []error) {
    var errors []error
    var success []Item

    for _, item := range items {
        if err := s.process(item); err != nil {
            errors = append(errors, fmt.Errorf("item %s: %w", item.ID, err))
        } else {
            success = append(success, item)
        }
    }

    return Result{Success: success}, errors
}

// 3. 降级处理 - 优雅降级
func (s *Service) GetRecommendations(ctx context.Context, userID string) ([]Recommendation, error) {
    // 尝试个性化推荐
    recs, err := s.personalizedRecs(ctx, userID)
    if err != nil {
        log.Printf("Personalized recs failed: %v, falling back to popular", err)
        // 降级到热门推荐
        return s.popularRecs(ctx)
    }
    return recs, nil
}

// 4. 重试策略 - 瞬态故障
func (s *Service) CallWithRetry(ctx context.Context, fn func() error) error {
    return retry.Retry(ctx, retry.Config{
        MaxAttempts: 3,
        Delay:       time.Second,
        Retryable: func(err error) bool {
            var netErr net.Error
            return errors.As(err, &netErr) || errors.Is(err, ErrTimeout)
        },
    }, fn)
}

// 5. 断路器 - 防止级联故障
func (s *Service) CallWithCircuitBreaker(fn func() error) error {
    return s.breaker.Execute(fn)
}
```

---

## 18.5 测试策略决策树

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                       测试策略决策树                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  测试层级                                                                   │
│       │                                                                     │
│       ├── 单元测试 ──▶ 测试范围？                                           │
│       │              │                                                      │
│       │              ├── 纯函数 ──▶ 表驱动测试                              │
│       │              │                                                      │
│       │              ├── 依赖外部 ──▶ Mock/Stub                             │
│       │              │         └── mockery / testify/mock                   │
│       │              │                                                      │
│       │              └── 并发代码 ──▶ 竞态检测                              │
│       │                        └── go test -race                            │
│       │                                                                     │
│       ├── 集成测试 ──▶ 依赖类型？                                           │
│       │              │                                                      │
│       │              ├── 数据库 ──▶ testcontainers-go                       │
│       │              │                                                      │
│       │              ├── HTTP ──▶ httptest                                  │
│       │              │                                                      │
│       │              └── 消息队列 ──▶ 嵌入式服务器                            │
│       │                                                                     │
│       └── E2E 测试 ──▶ 测试工具                                             │
│                      │                                                      │
│                      ├── API ──▶ Postman / REST Assured                     │
│                      │                                                      │
│                      ├── UI ──▶ Playwright / Selenium                       │
│                      │                                                      │
│                      └── 性能 ──▶ k6 / Locust                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 18.6 部署策略选择

| 策略 | 风险 | 复杂度 | 回滚速度 | 适用场景 |
|------|------|--------|----------|----------|
| **滚动更新** | 中 | 低 | 慢 | 标准无状态服务 |
| **蓝绿部署** | 低 | 中 | 快 | 关键业务系统 |
| **金丝雀** | 低 | 高 | 快 | 大规模用户系统 |
| **A/B 测试** | 低 | 高 | 快 | 产品功能验证 |

```go
// 金丝雀部署检查
func canaryCheck(ctx context.Context, version string) bool {
    // 检查错误率
    errorRate := metrics.GetErrorRate(version)
    if errorRate > 0.01 { // 1% 阈值
        return false
    }

    // 检查延迟
    p99Latency := metrics.GetP99Latency(version)
    if p99Latency > 500*time.Millisecond {
        return false
    }

    // 检查自定义业务指标
    if !businessMetrics.Healthy(version) {
        return false
    }

    return true
}

// 自动回滚
func autoRollback(deployment Deployment) {
    if !canaryCheck(context.Background(), deployment.Version) {
        log.Printf("Canary failed for version %s, rolling back", deployment.Version)
        kubernetes.Rollback(deployment)
    }
}
```

---

## 18.7 安全设计决策

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                       安全设计决策树                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  安全层面                                                                   │
│       │                                                                     │
│       ├── 认证 ──▶ 架构类型？                                               │
│       │              │                                                      │
│       │              ├── 微服务 ──▶ JWT / OAuth2                            │
│       │              │                                                      │
│       │              ├── 内部服务 ──▶ mTLS                                  │
│       │              │                                                      │
│       │              └── Web 应用 ──▶ Session / OIDC                        │
│       │                                                                     │
│       ├── 授权 ──▶ 模型选择？                                               │
│       │              │                                                      │
│       │              ├── 简单角色 ──▶ RBAC                                  │
│       │              │                                                      │
│       │              ├── 细粒度 ──▶ ABAC / OPA                              │
│       │              │                                                      │
│       │              └── 资源级 ──▶ Casbin                                  │
│       │                                                                     │
│       ├── 传输 ──▶ 加密要求？                                               │
│       │              │                                                      │
│       │              ├── 标准 ──▶ TLS 1.3                                   │
│       │              │                                                      │
│       │              └── 高安全 ──▶ mTLS + 证书轮转                         │
│       │                                                                     │
│       └── 数据 ──▶ 敏感级别？                                               │
│                      │                                                      │
│                      ├── 一般 ──▶ 传输加密                                  │
│                      │                                                      │
│                      ├── 敏感 ──▶ 字段级加密                                │
│                      │                                                      │
│                      └── 高度敏感 ──▶ 应用级加密 + HSM                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 18.8 监控告警决策

```go
// 告警级别定义
type AlertLevel int

const (
    AlertInfo AlertLevel = iota
    AlertWarning
    AlertCritical
    AlertEmergency
)

// 告警决策函数
func evaluateAlert(metric Metric) AlertLevel {
    switch metric.Name {
    case "error_rate":
        if metric.Value > 0.5 {
            return AlertEmergency
        } else if metric.Value > 0.1 {
            return AlertCritical
        } else if metric.Value > 0.01 {
            return AlertWarning
        }

    case "latency_p99":
        if metric.Value > 5000 { // 5s
            return AlertCritical
        } else if metric.Value > 1000 { // 1s
            return AlertWarning
        }

    case "cpu_usage":
        if metric.Value > 90 {
            return AlertCritical
        } else if metric.Value > 70 {
            return AlertWarning
        }

    case "disk_usage":
        if metric.Value > 90 {
            return AlertCritical
        } else if metric.Value > 80 {
            return AlertWarning
        }
    }

    return AlertInfo
}

// 告警抑制
func shouldSuppress(alert Alert) bool {
    // 维护窗口抑制
    if isMaintenanceWindow() {
        return true
    }

    // 依赖服务故障抑制
    if alert.Dependency != "" && isDependencyDown(alert.Dependency) {
        return true
    }

    // 告警风暴抑制
    if recentAlertCount(alert.Name) > 10 {
        return true
    }

    return false
}
```

---

## 18.9 技术债务管理决策

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                     技术债务管理决策框架                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  评估维度                              决策标准                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  影响范围                              修复优先级                           │
│  ─────────                             ───────────                          │
│  核心业务流程阻塞                        P0 - 立即修复                      │
│  性能显著下降 (>50%)                     P1 - 本周修复                      │
│  开发效率降低 (>30%)                     P2 - 本月修复                      │
│  代码可维护性差                          P3 - 下季度规划                    │
│  轻微优化空间                            P4 -  backlog                       │
│                                                                             │
│  修复成本 vs 收益分析                                                       │
│  ─────────────────────                                                      │
│                                                                             │
│  高影响 + 低成本  ──▶ 立即执行                                              │
│  高影响 + 高成本  ──▶ 分阶段执行                                            │
│  低影响 + 低成本  ──▶ 空闲时执行                                            │
│  低影响 + 高成本  ──▶ 暂缓 / 接受                                           │
│                                                                             │
│  重构策略                                                                   │
│  ─────────                                                                  │
│                                                                             │
│  绞杀者模式 (Strangler Fig)                                                 │
│    ├── 逐步替换旧系统                                                       │
│    └── 新功能走新架构                                                       │
│                                                                             │
│  Branch by Abstraction                                                      │
│    ├── 抽象层隔离变化                                                       │
│    └── 逐步替换实现                                                         │
│                                                                             │
│  特性开关 (Feature Flags)                                                   │
│    ├── 控制新功能上线                                                       │
│    └── 快速回滚能力                                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*本章提供了系统化的决策框架，帮助在复杂的技术选择中做出明智的决定。*
