# 分布式事务 (Distributed Transactions)

> **分类**: 成熟应用领域

---

## Saga 模式

```go
type Saga struct {
    steps []Step
}

type Step struct {
    Action    func() error
    Compensate func() error
}

func (s *Saga) Execute() error {
    completed := []int{}
    
    for i, step := range s.steps {
        if err := step.Action(); err != nil {
            // 补偿已完成的步骤
            for j := len(completed) - 1; j >= 0; j-- {
                s.steps[completed[j]].Compensate()
            }
            return err
        }
        completed = append(completed, i)
    }
    
    return nil
}

// 使用示例
saga := &Saga{
    steps: []Step{
        {
            Action:     func() error { return deductBalance(userID, amount) },
            Compensate: func() error { return refundBalance(userID, amount) },
        },
        {
            Action:     func() error { return createOrder(order) },
            Compensate: func() error { return cancelOrder(order.ID) },
        },
        {
            Action:     func() error { return reserveInventory(itemID, qty) },
            Compensate: func() error { return releaseInventory(itemID, qty) },
        },
    },
}

if err := saga.Execute(); err != nil {
    log.Fatal(err)
}
```

---

## 两阶段提交 (2PC)

```go
type Coordinator struct {
    participants []Participant
}

type Participant interface {
    Prepare() error
    Commit() error
    Rollback() error
}

func (c *Coordinator) Execute() error {
    // Phase 1: Prepare
    for _, p := range c.participants {
        if err := p.Prepare(); err != nil {
            c.rollback()
            return err
        }
    }
    
    // Phase 2: Commit
    for _, p := range c.participants {
        if err := p.Commit(); err != nil {
            // 需要人工介入
            return err
        }
    }
    
    return nil
}

func (c *Coordinator) rollback() {
    for _, p := range c.participants {
        p.Rollback()
    }
}
```

---

## 基于消息的可靠投递

```go
func PlaceOrder(ctx context.Context, order Order) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. 保存订单（状态：PENDING）
    if err := saveOrderTx(tx, order); err != nil {
        return err
    }
    
    // 2. 发送消息（本地事务表）
    if err := saveOutboxTx(tx, Event{
        Type: "ORDER_CREATED",
        Data: order,
    }); err != nil {
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return err
    }
    
    // 3. 异步发送消息到消息队列
    return publishEvent(order)
}
```
