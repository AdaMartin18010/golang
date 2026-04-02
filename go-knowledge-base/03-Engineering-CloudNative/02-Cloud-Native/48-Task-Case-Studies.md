# 任务系统案例研究 (Task System Case Studies)

> **分类**: 工程与云原生
> **标签**: #case-study #real-world #best-practices

---

## 案例一：电商平台订单处理

```go
// 场景：双十一期间处理千万级订单

// 架构设计
type OrderProcessingSystem struct {
    orderQueue     *PriorityQueue      // 按优先级处理
    paymentQueue   *DelayedQueue       // 延迟支付检查
    inventoryQueue *BatchQueue         // 批量库存扣减
    shippingQueue  *ScheduledQueue     // 定时发货
}

// 优先级策略
func (ops *OrderProcessingSystem) calculatePriority(order Order) int {
    priority := 0

    // VIP 用户高优先级
    if order.User.IsVIP {
        priority += 100
    }

    // 订单金额越高优先级越高
    priority += int(order.Amount / 1000)

    // 限时订单
    if order.IsFlashSale {
        priority += 200
    }

    return priority
}

// 流量削峰
func (ops *OrderProcessingSystem) HandleSpike(ctx context.Context, orders []Order) error {
    // 使用令牌桶限流
    limiter := rate.NewLimiter(rate.Limit(10000), 50000)

    for _, order := range orders {
        if err := limiter.Wait(ctx); err != nil {
            // 超出处理能力的订单放入队列稍后处理
            ops.orderQueue.Push(order, PriorityLow)
            continue
        }

        // 处理订单
        ops.processOrder(ctx, order)
    }

    return nil
}

// 结果：成功处理 1200万订单/小时，P99延迟 < 100ms
```

---

## 案例二：金融数据清算系统

```go
// 场景：每日批量处理交易数据，数据一致性要求高

// 两阶段提交实现
type SettlementSystem struct {
    db          *sql.DB
    kafka       sarama.Client
    coordinator *TransactionCoordinator
}

func (ss *SettlementSystem) ProcessSettlement(ctx context.Context, batch SettlementBatch) error {
    // 阶段1：准备
    tx := ss.coordinator.Begin(ctx)

    // 1. 冻结资金
    if err := ss.freezeFunds(ctx, tx, batch); err != nil {
        tx.Rollback()
        return err
    }

    // 2. 记录交易
    if err := ss.recordTransactions(ctx, tx, batch); err != nil {
        tx.Rollback()
        return err
    }

    // 3. 发送消息
    if err := ss.sendKafkaMessage(ctx, tx, batch); err != nil {
        tx.Rollback()
        return err
    }

    // 阶段2：提交
    if err := tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }

    return nil
}

// 幂等性保证
func (ss *SettlementSystem) ensureIdempotency(ctx context.Context, batchID string) error {
    key := fmt.Sprintf("settlement:%s", batchID)

    // 使用 Redis SETNX 保证幂等
    set, err := ss.redis.SetNX(ctx, key, "processing", 24*time.Hour).Result()
    if err != nil {
        return err
    }

    if !set {
        // 已处理过
        return ErrDuplicateBatch
    }

    return nil
}

// 结果：日处理 500万笔交易，零数据不一致
```

---

## 案例三：物联网设备数据处理

```go
// 场景：百万级 IoT 设备数据收集和实时分析

// 架构设计
type IoTDataPipeline struct {
    mqttIngestor   *MQTTIngestor       // 数据接入
    streamRouter   *StreamRouter       // 数据路由
    timeSeriesDB   *InfluxDB           // 时序存储
    alertEngine    *AlertEngine        // 告警引擎
    analytics      *StreamAnalytics    // 实时分析
}

// 数据分片策略
func (p *IoTDataPipeline) shardKey(deviceID string) string {
    // 按设备 ID 哈希分片
    hash := fnv.New32a()
    hash.Write([]byte(deviceID))
    shard := hash.Sum32() % 32
    return fmt.Sprintf("shard-%d", shard)
}

// 背压处理
func (p *IoTDataPipeline) HandleBackpressure(ctx context.Context, data DataPoint) error {
    shard := p.shardKey(data.DeviceID)
    queue := p.getQueue(shard)

    select {
    case queue <- data:
        return nil
    default:
        // 队列满，采样丢弃或降级处理
        if data.Priority == PriorityHigh {
            // 重要数据：阻塞等待
            queue <- data
        } else {
            // 普通数据：采样丢弃
            if rand.Float64() < 0.1 {
                return p.processSampled(ctx, data)
            }
            return ErrDropped
        }
    }
    return nil
}

// 结果：支持 1000万设备并发，数据处理延迟 < 500ms
```

---

## 案例四：图像处理服务

```go
// 场景：异步处理用户上传的图片

// 工作流定义
type ImageProcessingWorkflow struct {
    stages []WorkflowStage
}

func NewImageProcessingWorkflow() *ImageProcessingWorkflow {
    return &ImageProcessingWorkflow{
        stages: []WorkflowStage{
            {Name: "upload", Handler: uploadImage, Timeout: 30 * time.Second},
            {Name: "validate", Handler: validateImage, Timeout: 10 * time.Second},
            {Name: "resize", Handler: resizeImages, Timeout: 60 * time.Second},
            {Name: "watermark", Handler: addWatermark, Timeout: 30 * time.Second},
            {Name: "cdn-sync", Handler: syncToCDN, Timeout: 30 * time.Second},
        },
    }
}

// 任务分发策略
func (w *ImageProcessingWorkflow) Distribute(ctx context.Context, task ImageTask) error {
    // 根据图片大小选择队列
    var queue string

    if task.FileSize < 1*1024*1024 {
        queue = "small-images"
    } else if task.FileSize < 10*1024*1024 {
        queue = "medium-images"
    } else {
        queue = "large-images"
    }

    // 提交到对应队列
    return w.dispatcher.Dispatch(ctx, queue, task)
}

// GPU 资源调度
func (w *ImageProcessingWorkflow) scheduleGPU(ctx context.Context, task ImageTask) error {
    // 获取 GPU 资源
    gpu, err := w.gpuPool.Acquire(ctx, task.Priority)
    if err != nil {
        // GPU 不足，降级到 CPU 处理
        return w.fallbackToCPU(ctx, task)
    }
    defer w.gpuPool.Release(gpu)

    // 在 GPU 上执行
    return w.processOnGPU(ctx, gpu, task)
}

// 结果：日处理 100万张图片，成本降低 60%
```

---

## 最佳实践总结

```go
// 从案例中提取的通用模式

// 1. 分层队列
// - 按优先级分队列
// - 按业务类型分队列
// - 按资源需求分队列

// 2. 弹性伸缩
// - 基于队列深度的自动扩缩容
// - 基于资源利用率的扩缩容

// 3. 容错设计
// - 任务幂等性
// - 失败隔离
// - 优雅降级

// 4. 可观测性
// - 全链路追踪
// - 关键指标监控
// - 实时告警

// 5. 资源优化
// - 批量处理
// - 资源共享
// - 动态调度
```
