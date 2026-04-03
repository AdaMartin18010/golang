# FT-038-Distributed-Transactions

> **Dimension**: 01-Formal-Theory
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (2PC, 3PC, Saga, TCC)
> **Size**: >20KB

---

## 1. 分布式事务问题

### 1.1 ACID vs BASE

| 特性 | ACID (传统数据库) | BASE (分布式系统) |
|------|------------------|-------------------|
| 原子性 | 完全支持 | 最终一致 |
| 一致性 | 强一致 | 软状态 |
| 隔离性 | 完整隔离 | 无隔离 |
| 持久性 | 保证 | 最终持久 |
| 可用性 | 可能阻塞 | 优先可用 |

### 1.2 CAP下的选择

```
分布式事务需要在一致性和可用性之间权衡:

CP系统 (优先一致性):
- 两阶段提交 (2PC)
- 三阶段提交 (3PC)
- 分布式数据库 (Spanner, TiDB)

AP系统 (优先可用性):
- Saga模式
- TCC模式
- 本地消息表
- 最大努力通知
```

---

## 2. 两阶段提交 (2PC)

### 2.1 协议流程

```
Coordinator                    Participants
     │                              │
     │  Phase 1: Prepare            │
     │ ─────────────────────────────>│
     │                              │
     │         Vote YES/NO          │
     │ <─────────────────────────────│
     │                              │
     │  Phase 2: Commit/Abort       │
     │ ─────────────────────────────>│
     │                              │
     │            ACK               │
     │ <─────────────────────────────│
```

### 2.2 Go实现

```go
// 两阶段提交协调器
type Coordinator struct {
    participants []Participant
    log          TransactionLog
}

type Participant interface {
    Prepare(ctx context.Context, txID string) error
    Commit(ctx context.Context, txID string) error
    Rollback(ctx context.Context, txID string) error
}

func (c *Coordinator) Execute(ctx context.Context, tx Transaction) error {
    txID := generateTxID()

    // Phase 1: Prepare
    votes := make([]bool, len(c.participants))
    for i, p := range c.participants {
        if err := p.Prepare(ctx, txID); err != nil {
            votes[i] = false
            // 记录投票
            c.log.RecordVote(txID, i, false)
        } else {
            votes[i] = true
            c.log.RecordVote(txID, i, true)
        }
    }

    // 检查所有投票
    allYes := true
    for _, vote := range votes {
        if !vote {
            allYes = false
            break
        }
    }

    // Phase 2: Commit or Abort
    if allYes {
        c.log.RecordDecision(txID, DecisionCommit)
        for _, p := range c.participants {
            if err := p.Commit(ctx, txID); err != nil {
                // 记录需要后续处理
                c.log.RecordPendingCommit(txID)
            }
        }
    } else {
        c.log.RecordDecision(txID, DecisionAbort)
        for i, p := range c.participants {
            if votes[i] {
                p.Rollback(ctx, txID)
            }
        }
    }

    return nil
}

// 参与者实现示例
type OrderService struct {
    db *sql.DB
}

func (s *OrderService) Prepare(ctx context.Context, txID string) error {
    // 预执行操作，记录undo日志
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // 保存事务状态
    _, err = tx.Exec("INSERT INTO 2pc_transactions (tx_id, status) VALUES (?, 'prepared')", txID)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func (s *OrderService) Commit(ctx context.Context, txID string) error {
    _, err := s.db.Exec("UPDATE 2pc_transactions SET status='committed' WHERE tx_id=?", txID)
    return err
}

func (s *OrderService) Rollback(ctx context.Context, txID string) error {
    _, err := s.db.Exec("UPDATE 2pc_transactions SET status='aborted' WHERE tx_id=?", txID)
    // 执行undo操作
    return err
}
```

### 2.3 2PC的问题

| 问题 | 说明 | 解决方案 |
|------|------|---------|
| 同步阻塞 | 参与者等待协调者决定 | 超时机制 |
| 单点故障 | 协调者故障导致阻塞 | 协调者HA |
| 数据不一致 | 部分提交失败 | 事务恢复器 |
| 性能问题 | 两次网络往返 | 异步优化 |

---

## 3. 三阶段提交 (3PC)

### 3.1 协议流程

```
Coordinator                    Participants
     │                              │
     │  Phase 1: CanCommit?         │
     │ ─────────────────────────────>│
     │                              │
     │         Vote YES/NO          │
     │ <─────────────────────────────│
     │                              │
     │  Phase 2: PreCommit          │
     │ ─────────────────────────────>│
     │                              │
     │            ACK               │
     │ <─────────────────────────────│
     │                              │
     │  Phase 3: DoCommit/Abort     │
     │ ─────────────────────────────>│
     │                              │
     │            ACK               │
     │ <─────────────────────────────│
```

### 3.2 改进点

```go
// 3PC增加了PreCommit阶段
// 优势:
// 1. 参与者在CanCommit阶段不锁定资源
// 2. 协调者故障时，参与者可以自主决定

func (c *Coordinator) Execute3PC(ctx context.Context, tx Transaction) error {
    txID := generateTxID()

    // Phase 1: CanCommit
    // 只检查是否可以执行，不锁定

    // Phase 2: PreCommit
    // 预提交，锁定资源

    // Phase 3: DoCommit
    // 真正提交

    return nil
}

// 参与者超时处理
func (p *Participant3PC) HandleTimeout(phase string) {
    switch phase {
    case "waiting_precommit":
        // 等待PreCommit超时，可以安全中止
        p.Abort()
    case "waiting_commit":
        // 等待DoCommit超时，询问其他参与者
        if p.CanCommit() {
            p.Commit()  // 可以提交
        }
    }
}
```

---

## 4. Saga模式

### 4.1 概念

Saga将长事务拆分为多个本地事务，每个本地事务有对应的补偿操作。

```
正向流程:  T1 → T2 → T3 → ... → Tn
补偿流程:  如果T3失败，执行: C2 → C1

T1: 扣减库存      C1: 恢复库存
T2: 创建订单      C2: 取消订单
T3: 扣减余额      C3: 恢复余额
```

### 4.2 Saga实现

```go
// Saga编排器
type Saga struct {
    steps []SagaStep
    status SagaStatus
}

type SagaStep struct {
    Name        string
    Action      func() error
    Compensation func() error
}

func (s *Saga) Execute() error {
    completed := []int{}

    for i, step := range s.steps {
        if err := step.Action(); err != nil {
            // 执行补偿
            for j := len(completed) - 1; j >= 0; j-- {
                if compErr := s.steps[completed[j]].Compensation(); compErr != nil {
                    // 补偿失败，需要人工介入
                    s.logCompensationFailure(completed[j], compErr)
                }
            }
            return fmt.Errorf("step %s failed: %w", step.Name, err)
        }
        completed = append(completed, i)
    }

    return nil
}

// 订单Saga示例
func NewOrderSaga(inventorySvc InventoryService, orderSvc OrderService, paymentSvc PaymentService) *Saga {
    return &Saga{
        steps: []SagaStep{
            {
                Name: "DeductInventory",
                Action: func() error {
                    return inventorySvc.Reserve(skuID, quantity)
                },
                Compensation: func() error {
                    return inventorySvc.ReleaseReservation(skuID, quantity)
                },
            },
            {
                Name: "CreateOrder",
                Action: func() error {
                    orderID, err = orderSvc.Create(orderData)
                    return err
                },
                Compensation: func() error {
                    return orderSvc.Cancel(orderID)
                },
            },
            {
                Name: "ProcessPayment",
                Action: func() error {
                    return paymentSvc.Charge(userID, amount)
                },
                Compensation: func() error {
                    return paymentSvc.Refund(userID, amount)
                },
            },
        },
    }
}
```

### 4.3 Saga编排方式

| 方式 | 优点 | 缺点 |
|------|------|------|
| 编排 (Choreography) | 松耦合 | 难以追踪 |
| 协调 (Orchestration) | 集中控制 | 单点风险 |

```go
// 编排式: 通过事件驱动
// InventoryService扣减后发出"InventoryDeducted"事件
// OrderService监听并创建订单

// 协调式: Saga协调器集中控制
type SagaOrchestrator struct {
    inventorySvc InventoryService
    orderSvc     OrderService
    paymentSvc   PaymentService
}

func (o *SagaOrchestrator) StartOrderSaga(ctx context.Context, req OrderRequest) error {
    // 1. 扣减库存
    if err := o.inventorySvc.Reserve(ctx, req.SKUID, req.Quantity); err != nil {
        return err
    }

    // 2. 创建订单
    order, err := o.orderSvc.Create(ctx, req)
    if err != nil {
        o.inventorySvc.Release(ctx, req.SKUID, req.Quantity)
        return err
    }

    // 3. 处理支付
    if err := o.paymentSvc.Charge(ctx, req.UserID, req.Amount); err != nil {
        o.orderSvc.Cancel(ctx, order.ID)
        o.inventorySvc.Release(ctx, req.SKUID, req.Quantity)
        return err
    }

    return nil
}
```

---

## 5. TCC模式

### 5.1 概念

Try-Confirm-Cancel: 资源预留模式

```
Try: 预留资源
Confirm: 确认执行业务
Cancel: 释放预留资源

示例:
Try: 冻结账户100元
Confirm: 从冻结金额中扣除
Cancel: 解冻金额
```

### 5.2 TCC实现

```go
// TCC接口
type TCCAction interface {
    Try(ctx context.Context, bizID string, params interface{}) error
    Confirm(ctx context.Context, bizID string) error
    Cancel(ctx context.Context, bizID string) error
}

// 账户服务TCC实现
type AccountTCC struct {
    db *sql.DB
}

func (a *AccountTCC) Try(ctx context.Context, bizID string, params interface{}) error {
    p := params.(PaymentParams)

    return a.db.Transaction(func(tx *gorm.DB) error {
        // 检查余额
        var account Account
        if err := tx.First(&account, "user_id = ?", p.UserID).Error; err != nil {
            return err
        }

        if account.Balance < p.Amount {
            return ErrInsufficientBalance
        }

        // 冻结金额
        account.Balance -= p.Amount
        account.Frozen += p.Amount

        // 记录TCC事务
        tccRecord := TCCRecord{
            BizID:     bizID,
            UserID:    p.UserID,
            Amount:    p.Amount,
            Status:    TCCStatusTry,
            CreatedAt: time.Now(),
        }

        return tx.Save(&account).Create(&tccRecord).Error
    })
}

func (a *AccountTCC) Confirm(ctx context.Context, bizID string) error {
    return a.db.Transaction(func(tx *gorm.DB) error {
        var record TCCRecord
        if err := tx.First(&record, "biz_id = ?", bizID).Error; err != nil {
            return err
        }

        if record.Status != TCCStatusTry {
            return ErrInvalidTCCStatus
        }

        // 从冻结中扣减
        var account Account
        tx.First(&account, "user_id = ?", record.UserID)
        account.Frozen -= record.Amount

        record.Status = TCCStatusConfirm
        record.ConfirmedAt = time.Now()

        return tx.Save(&account).Save(&record).Error
    })
}

func (a *AccountTCC) Cancel(ctx context.Context, bizID string) error {
    return a.db.Transaction(func(tx *gorm.DB) error {
        var record TCCRecord
        if err := tx.First(&record, "biz_id = ?", bizID).Error; err != nil {
            return err
        }

        if record.Status != TCCStatusTry {
            return nil  // 已经处理过了
        }

        // 解冻金额
        var account Account
        tx.First(&account, "user_id = ?", record.UserID)
        account.Frozen -= record.Amount
        account.Balance += record.Amount

        record.Status = TCCStatusCancel
        record.CancelledAt = time.Now()

        return tx.Save(&account).Save(&record).Error
    })
}
```

---

## 6. 本地消息表

### 6.1 概念

通过本地事务保证业务操作和消息发送的原子性。

```go
// 消息表
type OutboxMessage struct {
    ID          uint64
    Topic       string
    Key         string
    Payload     string
    Status      string  // pending, sent, failed
    RetryCount  int
    CreatedAt   time.Time
    ProcessedAt *time.Time
}

// 业务服务
type OrderService struct {
    db *gorm.DB
    producer sarama.AsyncProducer
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建订单
        order := Order{...}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // 2. 写入消息表 (同一事务)
        message := OutboxMessage{
            Topic:   "order_created",
            Key:     order.ID,
            Payload: mustMarshal(order),
            Status:  "pending",
        }
        if err := tx.Create(&message).Error; err != nil {
            return err
n        }

        return nil
    })
}

// 消息投递服务
func (s *OrderService) StartMessagePoller() {
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            s.processOutbox()
        }
    }()
}

func (s *OrderService) processOutbox() {
    var messages []OutboxMessage
    s.db.Where("status = ?", "pending").Limit(100).Find(&messages)

    for _, msg := range messages {
        // 发送到Kafka
        s.producer.Input() <- &sarama.ProducerMessage{
            Topic: msg.Topic,
            Key:   sarama.StringEncoder(msg.Key),
            Value: sarama.StringEncoder(msg.Payload),
        }

        // 更新状态
        now := time.Now()
        s.db.Model(&msg).Updates(map[string]interface{}{
            "status":       "sent",
            "processed_at": &now,
        })
    }
}
```

---

## 7. 方案对比

| 方案 | 一致性 | 性能 | 复杂度 | 适用场景 |
|------|--------|------|--------|---------|
| 2PC | 强一致 | 低 | 中 | 传统金融 |
| 3PC | 强一致 | 低 | 高 | 较少使用 |
| Saga | 最终一致 | 高 | 中 | 长事务业务 |
| TCC | 最终一致 | 高 | 高 | 短事务高并发 |
| 消息表 | 最终一致 | 高 | 中 | 异步场景 |

---

## 8. 参考文献

1. "Principles of Distributed Computing" - Roger Wattenhofer
2. "Saga Pattern" - Chris Richardson
3. "Distributed Systems: Principles and Paradigms"
4. DTM (Distributed Transaction Manager) Documentation
5. Seata Documentation

---

*Last Updated: 2026-04-03*
