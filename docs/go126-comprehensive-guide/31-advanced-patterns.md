# Go高级设计模式

> 面向复杂场景的Go高级模式与架构设计

---

## 一、函数选项模式进阶

### 1.1 带验证的选项模式

```text
高级选项模式:
────────────────────────────────────────
支持验证、默认值、选项组合

代码示例:
type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    TLSConfig    *tls.Config
    Middleware   []Middleware
}

// 选项类型
type Option func(*ServerConfig) error

// 选项实现
func WithHost(host string) Option {
    return func(c *ServerConfig) error {
        if host == "" {
            return errors.New("host cannot be empty")
        }
        c.Host = host
        return nil
    }
}

func WithPort(port int) Option {
    return func(c *ServerConfig) error {
        if port <= 0 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
        c.Port = port
        return nil
    }
}

func WithTimeout(read, write time.Duration) Option {
    return func(c *ServerConfig) error {
        if read <= 0 || write <= 0 {
            return errors.New("timeout must be positive")
        }
        c.ReadTimeout = read
        c.WriteTimeout = write
        return nil
    }
}

// 带错误处理的构造
func NewServer(opts ...Option) (*Server, error) {
    cfg := &ServerConfig{
        Host:         "0.0.0.0",
        Port:         8080,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    for _, opt := range opts {
        if err := opt(cfg); err != nil {
            return nil, err
        }
    }

    return &Server{config: cfg}, nil
}

// 使用
func advancedOptionsExample() {
    srv, err := NewServer(
        WithHost("localhost"),
        WithPort(9090),
        WithTimeout(10*time.Second, 10*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }
    _ = srv
}
```

---

## 二、依赖注入模式

### 2.1 构造函数注入

```text
依赖注入实现:
────────────────────────────────────────

代码示例:
// 接口定义
type Database interface {
    Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
    Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
}

type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
}

// 服务实现
type UserService struct {
    db     Database
    cache  Cache
    logger Logger
}

// 构造函数注入
func NewUserService(db Database, cache Cache, logger Logger) *UserService {
    return &UserService{
        db:     db,
        cache:  cache,
        logger: logger,
    }
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // 先查缓存
    if cached, err := s.cache.Get(ctx, "user:"+id); err == nil {
        var user User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }

    // 查数据库
    rows, err := s.db.Query(ctx, "SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        s.logger.Error("failed to query user", err, String("id", id))
        return nil, err
    }
    defer rows.Close()

    // 处理结果...
    return &User{}, nil
}

// Wire代码生成
//go:build wireinject
// +build wireinject

func InitializeUserService() (*UserService, error) {
    wire.Build(
        NewPostgresDB,
        NewRedisCache,
        NewZapLogger,
        NewUserService,
    )
    return nil, nil
}
```

---

## 三、事件驱动模式

### 3.1 事件总线实现

```text
事件总线:
────────────────────────────────────────

代码示例:
type Event struct {
    Type      string
    Payload   interface{}
    Timestamp time.Time
    Metadata  map[string]string
}

type Handler func(ctx context.Context, event Event) error

type EventBus struct {
    handlers map[string][]Handler
    mu       sync.RWMutex
    workers  int
    queue    chan Event
}

func NewEventBus(workers int, queueSize int) *EventBus {
    return &EventBus{
        handlers: make(map[string][]Handler),
        workers:  workers,
        queue:    make(chan Event, queueSize),
    }
}

func (b *EventBus) Subscribe(eventType string, handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(ctx context.Context, event Event) error {
    select {
    case b.queue <- event:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    default:
        return errors.New("event queue full")
    }
}

func (b *EventBus) Start(ctx context.Context) {
    for i := 0; i < b.workers; i++ {
        go b.worker(ctx)
    }
}

func (b *EventBus) worker(ctx context.Context) {
    for {
        select {
        case event := <-b.queue:
            b.processEvent(ctx, event)
        case <-ctx.Done():
            return
        }
    }
}

func (b *EventBus) processEvent(ctx context.Context, event Event) {
    b.mu.RLock()
    handlers := b.handlers[event.Type]
    b.mu.RUnlock()

    var wg sync.WaitGroup
    for _, handler := range handlers {
        wg.Add(1)
        go func(h Handler) {
            defer wg.Done()
            if err := h(ctx, event); err != nil {
                log.Printf("handler error: %v", err)
            }
        }(handler)
    }
    wg.Wait()
}

// 使用
func eventBusExample() {
    bus := NewEventBus(10, 1000)
    ctx := context.Background()

    // 订阅
    bus.Subscribe("user.created", func(ctx context.Context, e Event) error {
        log.Printf("User created: %v", e.Payload)
        return nil
    })

    bus.Subscribe("user.created", func(ctx context.Context, e Event) error {
        // 发送邮件
        return nil
    })

    bus.Start(ctx)

    // 发布
    bus.Publish(ctx, Event{
        Type:      "user.created",
        Payload:   User{ID: "123", Name: "John"},
        Timestamp: time.Now(),
    })
}
```

---

## 四、CQRS模式

### 4.1 命令查询分离

```text
CQRS实现:
────────────────────────────────────────

代码示例:
// 命令
 type CreateOrderCommand struct {
    UserID    string
    ProductID string
    Quantity  int
}

type CommandHandler interface {
    Handle(ctx context.Context, cmd interface{}) error
}

// 查询
type OrderQuery struct {
    ID     string
    UserID string
    Status string
}

type QueryHandler interface {
    Handle(ctx context.Context, query interface{}) (interface{}, error)
}

// 命令处理器
type CreateOrderHandler struct {
    eventStore EventStore
    bus        EventBus
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd interface{}) error {
    createCmd := cmd.(CreateOrderCommand)

    // 创建聚合
    order, err := NewOrder(createCmd.UserID, createCmd.ProductID, createCmd.Quantity)
    if err != nil {
        return err
    }

    // 保存事件
    events := order.GetUncommittedEvents()
    if err := h.eventStore.Save(ctx, order.ID, events); err != nil {
        return err
    }

    // 发布事件
    for _, event := range events {
        if err := h.bus.Publish(ctx, event); err != nil {
            return err
        }
    }

    return nil
}

// 查询处理器 (读模型)
type OrderQueryHandler struct {
    readDB *sql.DB
}

func (h *OrderQueryHandler) Handle(ctx context.Context, query interface{}) (interface{}, error) {
    q := query.(OrderQuery)

    var order OrderView
    err := h.readDB.QueryRowContext(ctx,
        "SELECT id, user_id, product_id, quantity, status FROM order_views WHERE id = ?",
        q.ID,
    ).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Status)

    if err != nil {
        return nil, err
    }

    return order, nil
}

// 事件投影 (同步读模型)
type OrderProjector struct {
    readDB *sql.DB
}

func (p *OrderProjector) HandleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    _, err := p.readDB.ExecContext(ctx,
        "INSERT INTO order_views (id, user_id, product_id, quantity, status) VALUES (?, ?, ?, ?, ?)",
        event.OrderID, event.UserID, event.ProductID, event.Quantity, "created",
    )
    return err
}
```

---

## 五、Saga模式

### 5.1 分布式事务

```text
Saga模式实现:
────────────────────────────────────────

代码示例:
type Saga struct {
    steps []SagaStep
}

type SagaStep struct {
    Action    func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (s *Saga) AddStep(action, compensate func(ctx context.Context) error) {
    s.steps = append(s.steps, SagaStep{
        Action:     action,
        Compensate: compensate,
    })
}

func (s *Saga) Execute(ctx context.Context) error {
    completed := make([]int, 0, len(s.steps))

    for i, step := range s.steps {
        if err := step.Action(ctx); err != nil {
            // 补偿已完成的步骤
            for j := len(completed) - 1; j >= 0; j-- {
                if compErr := s.steps[completed[j]].Compensate(ctx); compErr != nil {
                    log.Printf("compensation error: %v", compErr)
                }
            }
            return err
        }
        completed = append(completed, i)
    }

    return nil
}

// 使用
func sagaExample() {
    saga := &Saga{}

    // 步骤1: 扣减库存
    saga.AddStep(
        func(ctx context.Context) error {
            return inventoryService.Reserve(ctx, "product-1", 10)
        },
        func(ctx context.Context) error {
            return inventoryService.Release(ctx, "product-1", 10)
        },
    )

    // 步骤2: 创建订单
    saga.AddStep(
        func(ctx context.Context) error {
            return orderService.Create(ctx, "order-1", "product-1", 10)
        },
        func(ctx context.Context) error {
            return orderService.Cancel(ctx, "order-1")
        },
    )

    // 步骤3: 扣款
    saga.AddStep(
        func(ctx context.Context) error {
            return paymentService.Charge(ctx, "user-1", 100)
        },
        func(ctx context.Context) error {
            return paymentService.Refund(ctx, "user-1", 100)
        },
    )

    if err := saga.Execute(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

---

## 六、多租户架构

### 6.1 租户隔离

```text
多租户实现:
────────────────────────────────────────

代码示例:
type TenantID string

func TenantFromContext(ctx context.Context) (TenantID, error) {
    tenant, ok := ctx.Value(tenantKey{}).(TenantID)
    if !ok {
        return "", errors.New("tenant not found in context")
    }
    return tenant, nil
}

// 中间件
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := r.Header.Get("X-Tenant-ID")
        if tenantID == "" {
            http.Error(w, "tenant required", 400)
            return
        }

        ctx := context.WithValue(r.Context(), tenantKey{}, TenantID(tenantID))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 租户感知的数据库
type TenantDB struct {
    db *sql.DB
}

func (tdb *TenantDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    tenant, err := TenantFromContext(ctx)
    if err != nil {
        return nil, err
    }

    // 添加租户过滤
    query = fmt.Sprintf("SELECT * FROM (%s) WHERE tenant_id = ?", query)
    args = append(args, tenant)

    return tdb.db.QueryContext(ctx, query, args...)
}

// 行级安全 (PostgreSQL)
func setupRLS(db *sql.DB) {
    db.Exec(`
        CREATE POLICY tenant_isolation ON orders
        USING (tenant_id = current_setting('app.current_tenant')::UUID);
    `)
}
```

---

*本章提供了Go高级设计模式，面向复杂业务场景的架构设计。*
