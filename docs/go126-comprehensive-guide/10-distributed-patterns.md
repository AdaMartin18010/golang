# 第十章：分布式系统设计模式

> Go 微服务架构中的分布式设计模式

---

## 10.1 服务治理模式

### 10.1.1 服务发现

```go
// 服务发现接口
type ServiceDiscovery interface {
    Register(service ServiceInstance) error
    Deregister(serviceID string) error
    Discover(serviceName string) ([]ServiceInstance, error)
    Watch(serviceName string) (<-chan []ServiceInstance, error)
}

type ServiceInstance struct {
    ID       string
    Name     string
    Host     string
    Port     int
    Metadata map[string]string
    Health   HealthStatus
}

// Consul 实现
func NewConsulDiscovery(addr string) (ServiceDiscovery, error) {
    client, err := api.NewClient(&api.Config{
        Address: addr,
    })
    if err != nil {
        return nil, err
    }
    return &ConsulDiscovery{client: client}, nil
}

type ConsulDiscovery struct {
    client *api.Client
}

func (c *ConsulDiscovery) Register(service ServiceInstance) error {
    reg := &api.AgentServiceRegistration{
        ID:      service.ID,
        Name:    service.Name,
        Address: service.Host,
        Port:    service.Port,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", service.Host, service.Port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }
    return c.client.Agent().ServiceRegister(reg)
}

func (c *ConsulDiscovery) Discover(serviceName string) ([]ServiceInstance, error) {
    entries, _, err := c.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }

    var instances []ServiceInstance
    for _, entry := range entries {
        instances = append(instances, ServiceInstance{
            ID:   entry.Service.ID,
            Name: entry.Service.Service,
            Host: entry.Service.Address,
            Port: entry.Service.Port,
        })
    }
    return instances, nil
}
```

### 10.1.2 负载均衡

```go
// 负载均衡策略
type LoadBalancer interface {
    Select(instances []ServiceInstance) (ServiceInstance, error)
}

// 轮询
type RoundRobin struct {
    counter uint64
}

func (r *RoundRobin) Select(instances []ServiceInstance) (ServiceInstance, error) {
    if len(instances) == 0 {
        return ServiceInstance{}, errors.New("no instances available")
    }
    idx := atomic.AddUint64(&r.counter, 1) % uint64(len(instances))
    return instances[idx], nil
}

// 随机
type Random struct {
    rnd *rand.Rand
    mu  sync.Mutex
}

func (r *Random) Select(instances []ServiceInstance) (ServiceInstance, error) {
    if len(instances) == 0 {
        return ServiceInstance{}, errors.New("no instances available")
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    return instances[r.rnd.Intn(len(instances))], nil
}

// 加权轮询
type WeightedRoundRobin struct {
    instances []WeightedInstance
    current   int
    gcdWeight int
    maxWeight int
}

type WeightedInstance struct {
    Instance ServiceInstance
    Weight   int
    Current  int
}

func (w *WeightedRoundRobin) Select() (ServiceInstance, error) {
    for {
        w.current = (w.current + 1) % len(w.instances)
        if w.current == 0 {
            w.gcdWeight--
            if w.gcdWeight == 0 {
                w.gcdWeight = w.maxWeight
            }
        }
        if w.instances[w.current].Weight >= w.gcdWeight {
            return w.instances[w.current].Instance, nil
        }
    }
}
```

---

## 10.2 弹性模式

### 10.2.1 熔断器 (Circuit Breaker)

```go
type CircuitBreaker struct {
    name           string
    maxFailures    int
    timeout        time.Duration
    halfOpenMaxCalls int

    state          State
    failures       int
    successes      int
    lastFailureTime time.Time
    mutex          sync.RWMutex

    // 指标
    metrics        *Metrics
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (cb *CircuitBreaker) Call(fn func() error) error {
    if !cb.canCall() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canCall() bool {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailureTime) > cb.timeout {
            cb.mutex.RUnlock()
            cb.mutex.Lock()
            cb.state = StateHalfOpen
            cb.failures = 0
            cb.successes = 0
            cb.mutex.Unlock()
            cb.mutex.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        return cb.successes+cb.failures < cb.halfOpenMaxCalls
    }
    return false
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err == nil {
        cb.onSuccess()
    } else {
        cb.onFailure()
    }
}

func (cb *CircuitBreaker) onSuccess() {
    switch cb.state {
    case StateClosed:
        cb.failures = 0
    case StateHalfOpen:
        cb.successes++
        if cb.successes >= cb.halfOpenMaxCalls {
            cb.state = StateClosed
            cb.failures = 0
            cb.successes = 0
        }
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failures++
    cb.lastFailureTime = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    case StateHalfOpen:
        cb.state = StateOpen
    }
}
```

### 10.2.2 舱壁隔离 (Bulkhead)

```go
// 舱壁模式：隔离资源
type Bulkhead struct {
    name        string
    maxConcurrent int
    maxWait     time.Duration

    semaphore   chan struct{}
    metrics     *Metrics
}

func NewBulkhead(name string, maxConcurrent int, maxWait time.Duration) *Bulkhead {
    return &Bulkhead{
        name:        name,
        maxConcurrent: maxConcurrent,
        maxWait:     maxWait,
        semaphore:   make(chan struct{}, maxConcurrent),
    }
}

func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
    select {
    case b.semaphore <- struct{}{}:
        defer func() { <-b.semaphore }()
        return fn()
    case <-time.After(b.maxWait):
        return ErrBulkheadFull
    case <-ctx.Done():
        return ctx.Err()
    }
}

// 服务级别的舱壁隔离
type ServiceBulkheads struct {
    bulkheads map[string]*Bulkhead
    mu        sync.RWMutex
}

func (s *ServiceBulkheads) Execute(serviceName string, fn func() error) error {
    s.mu.RLock()
    bulkhead, ok := s.bulkheads[serviceName]
    s.mu.RUnlock()

    if !ok {
        // 使用默认配置
        bulkhead = NewBulkhead(serviceName, 100, 5*time.Second)
        s.mu.Lock()
        s.bulkheads[serviceName] = bulkhead
        s.mu.Unlock()
    }

    return bulkhead.Execute(context.Background(), fn)
}
```

### 10.2.3 重试模式

```go
type RetryConfig struct {
    MaxAttempts     int
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
    Jitter          float64
    RetryableErrors func(error) bool
}

func Retry(ctx context.Context, config RetryConfig, fn func() error) error {
    var err error
    interval := config.InitialInterval

    for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
        err = fn()
        if err == nil {
            return nil
        }

        if !config.RetryableErrors(err) {
            return err
        }

        if attempt == config.MaxAttempts {
            break
        }

        // 计算下次重试间隔
        jitter := time.Duration(rand.Float64() * config.Jitter * float64(interval))
        sleepTime := interval + jitter

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(sleepTime):
        }

        // 指数退避
        interval = time.Duration(float64(interval) * config.Multiplier)
        if interval > config.MaxInterval {
            interval = config.MaxInterval
        }
    }

    return fmt.Errorf("max retry attempts reached: %w", err)
}

// 使用
config := RetryConfig{
    MaxAttempts:     5,
    InitialInterval: 100 * time.Millisecond,
    MaxInterval:     5 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.1,
    RetryableErrors: func(err error) bool {
        return errors.Is(err, ErrTemporary) || errors.Is(err, ErrTimeout)
    },
}

err := Retry(ctx, config, func() error {
    return callExternalService()
})
```

---

## 10.3 Saga 分布式事务模式

```go
// Saga 协调器
type Saga struct {
    steps []SagaStep
    compensations []func() error
}

type SagaStep struct {
    Action       func() error
    Compensation func() error
    Name         string
}

func NewSaga() *Saga {
    return &Saga{
        compensations: make([]func() error, 0),
    }
}

func (s *Saga) AddStep(step SagaStep) {
    s.steps = append(s.steps, step)
}

func (s *Saga) Execute() error {
    for i, step := range s.steps {
        if err := step.Action(); err != nil {
            // 执行补偿
            s.compensate(i)
            return fmt.Errorf("step %s failed: %w", step.Name, err)
        }
        s.compensations = append(s.compensations, step.Compensation)
    }
    return nil
}

func (s *Saga) compensate(failedIndex int) {
    for i := len(s.compensations) - 1; i >= 0; i-- {
        if err := s.compensations[i](); err != nil {
            log.Printf("Compensation failed: %v", err)
        }
    }
}

// 使用示例：订单处理
func ProcessOrder(order Order) error {
    saga := NewSaga()

    // Step 1: 扣减库存
    saga.AddStep(SagaStep{
        Name: "Deduct Inventory",
        Action: func() error {
            return inventoryService.Deduct(order.Items)
        },
        Compensation: func() error {
            return inventoryService.Restore(order.Items)
        },
    })

    // Step 2: 扣款
    saga.AddStep(SagaStep{
        Name: "Charge Payment",
        Action: func() error {
            return paymentService.Charge(order.UserID, order.Total)
        },
        Compensation: func() error {
            return paymentService.Refund(order.UserID, order.Total)
        },
    })

    // Step 3: 创建订单
    saga.AddStep(SagaStep{
        Name: "Create Order",
        Action: func() error {
            return orderService.Create(order)
        },
        Compensation: func() error {
            return orderService.Cancel(order.ID)
        },
    })

    return saga.Execute()
}
```

---

## 10.4 事件驱动架构

### 10.4.1 事件总线

```go
type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

type Event struct {
    Type      string
    Payload   interface{}
    Timestamp time.Time
    Metadata  map[string]string
}

type EventHandler func(ctx context.Context, event Event) error

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(ctx context.Context, event Event) error {
    b.mu.RLock()
    handlers := b.handlers[event.Type]
    b.mu.RUnlock()

    var wg sync.WaitGroup
    errChan := make(chan error, len(handlers))

    for _, handler := range handlers {
        wg.Add(1)
        go func(h EventHandler) {
            defer wg.Done()
            if err := h(ctx, event); err != nil {
                errChan <- err
            }
        }(handler)
    }

    wg.Wait()
    close(errChan)

    var errs []error
    for err := range errChan {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return fmt.Errorf("handler errors: %v", errs)
    }
    return nil
}

// 集成消息队列
func (b *EventBus) PublishToMQ(event Event) error {
    // 序列化事件
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    // 发布到 Kafka/NATS
    return kafkaClient.Publish(event.Type, data)
}
```

### 10.4.2 CQRS 模式

```go
// 命令端
type OrderCommandHandler struct {
    eventStore EventStore
    bus        EventBus
}

func (h *OrderCommandHandler) CreateOrder(cmd CreateOrderCommand) error {
    order := NewOrder(cmd.UserID, cmd.Items)

    // 存储事件
    if err := h.eventStore.Save(order.Events()); err != nil {
        return err
    }

    // 发布事件
    for _, event := range order.Events() {
        if err := h.bus.Publish(context.Background(), event); err != nil {
            return err
        }
    }

    return nil
}

// 查询端
type OrderQueryHandler struct {
    readModel OrderReadModel
}

func (h *OrderQueryHandler) GetOrder(orderID string) (OrderDTO, error) {
    return h.readModel.FindByID(orderID)
}

func (h *OrderQueryHandler) ListOrders(userID string) ([]OrderDTO, error) {
    return h.readModel.FindByUser(userID)
}

// 事件处理器更新读模型
type OrderProjector struct {
    readModel OrderReadModel
}

func (p *OrderProjector) HandleOrderCreated(ctx context.Context, event Event) error {
    created := event.Payload.(OrderCreatedEvent)

    dto := OrderDTO{
        ID:        created.OrderID,
        UserID:    created.UserID,
        Status:    "pending",
        Total:     created.Total,
        CreatedAt: event.Timestamp,
    }

    return p.readModel.Save(dto)
}
```

---

## 10.5 分布式追踪

```go
// OpenTelemetry 集成

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("my-service")

func ProcessOrder(ctx context.Context, order Order) error {
    // 创建 span
    ctx, span := tracer.Start(ctx, "ProcessOrder",
        trace.WithAttributes(
            attribute.String("order.id", order.ID),
            attribute.String("user.id", order.UserID),
        ),
    )
    defer span.End()

    // 子操作
    if err := validateOrder(ctx, order); err != nil {
        span.RecordError(err)
        return err
    }

    if err := saveOrder(ctx, order); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "save failed")
        return err
    }

    return nil
}

func validateOrder(ctx context.Context, order Order) error {
    ctx, span := tracer.Start(ctx, "validateOrder")
    defer span.End()

    // 验证逻辑
    if len(order.Items) == 0 {
        span.SetAttributes(attribute.Bool("order.valid", false))
        return ErrEmptyOrder
    }

    span.SetAttributes(attribute.Bool("order.valid", true))
    return nil
}

// HTTP 中间件
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 从请求头提取 trace context
        ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

        ctx, span := tracer.Start(ctx, r.URL.Path,
            trace.WithAttributes(
                attribute.String("http.method", r.Method),
                attribute.String("http.url", r.URL.String()),
            ),
        )
        defer span.End()

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 10.6 分布式锁

```go
// Redis 分布式锁
type DistributedLock struct {
    client *redis.Client
    key    string
    value  string
    ttl    time.Duration
}

func (l *DistributedLock) Lock(ctx context.Context) error {
    l.value = uuid.New().String()

    for {
        ok, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
        if err != nil {
            return err
        }
        if ok {
            // 获取锁成功，启动续期
            go l.keepAlive()
            return nil
        }

        // 等待重试
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
        }
    }
}

func (l *DistributedLock) Unlock(ctx context.Context) error {
    // 使用 Lua 脚本确保原子性
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `

    result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
    if err != nil {
        return err
    }

    if result.(int64) == 0 {
        return ErrLockNotHeld
    }
    return nil
}

func (l *DistributedLock) keepAlive() {
    ticker := time.NewTicker(l.ttl / 3)
    defer ticker.Stop()

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        l.client.Expire(ctx, l.key, l.ttl)
        cancel()
    }
}
```

---

*本章涵盖了构建可靠分布式系统的核心模式，包括弹性、事务、事件驱动和可观测性等方面。*
