# Go设计模式深度分析

> 基于形式化方法的Go惯用模式与架构设计

---

## 一、设计模式的形式化基础

### 1.1 模式的代数表示

```
模式的抽象定义:
────────────────────────────────────────
设计模式 = (Context, Problem, Solution, Consequences)

形式化:
Context:  C  (设计情境)
Problem:  P ⊆ C × Requirements  (问题空间)
Solution: S: C → Implementation  (解决方案)
Consequences: Q: S(C) → Properties (产生属性)

Go模式的特点:
├─ 隐式接口降低耦合
├─ 一等函数替代部分OO模式
├─ 组合优于继承
└─ 并发原语原生支持
```

### 1.2 模式分类学

```
Go模式分类:
────────────────────────────────────────

创建型模式:
├── 单例 (Singleton)        → sync.Once实现
├── 工厂 (Factory)          → 函数+接口组合
├── 构建者 (Builder)        → 选项函数模式
└── 对象池 (Object Pool)    → sync.Pool

结构型模式:
├── 适配器 (Adapter)        → 接口隐式实现
├── 装饰器 (Decorator)      → 函数包装
├── 外观 (Facade)           → 包级别封装
└── 组合 (Composite)        → 嵌入类型

行为型模式:
├── 策略 (Strategy)         → 函数值
├── 观察者 (Observer)       → Channel通知
├── 命令 (Command)          → 函数闭包
└── 模板方法 (Template)     → 接口组合

并发模式:
├── Worker Pool             → Goroutine池
├── Pipeline                → Channel链
├── Fan-out/Fan-in          → 并行分发聚合
└── Context传播             → 取消信号
```

---

## 二、创建型模式

### 2.1 单例模式的形式化

```
单例公理:
────────────────────────────────────────
∀t: instance(t) = instance(t₀)  (唯一性)
instance(0) = nil → instance(1) = new(Instance)  (惰性)

Go实现:
type Singleton struct{}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}

正确性证明:
1. sync.Once保证Do内代码执行且仅执行一次
2. 第一次调用创建instance
3. 后续调用返回同一instance
∴ 满足单例定义

对比传统实现:
├─ 懒汉式: 需要双重检查锁定
├─ 饿汉式: 启动时初始化
└─ Go: sync.Once封装复杂性

代码示例:
// 正确的单例实现
type Database struct {
    conn *sql.DB
}

var (
    db   *Database
    once sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        conn, err := sql.Open("postgres", "...")
        if err != nil {
            log.Fatal(err)
        }
        db = &Database{conn: conn}
    })
    return db
}

// 反例: 错误的单例实现
var instance *Singleton

func GetInstanceWrong() *Singleton {
    if instance == nil {        // 竞态条件！
        instance = &Singleton{} // 可能创建多个实例
    }
    return instance
}
```

### 2.2 函数选项模式

```
构建者模式的函数式实现:
────────────────────────────────────────

理论基础:
Option = Server → Server  (端态变换)

实现:
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        timeout: defaultTimeout,
        logger:  defaultLogger,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

代数性质:
├─ 选项应用顺序无关 (最终值取决于最后应用)
├─ 可组合: opts = append(defaultOpts, userOpts...)
└─ 可扩展: 新增选项不破坏现有代码

代码示例:
// 完整的服务器配置选项
type Server struct {
    addr     string
    timeout  time.Duration
    logger   *log.Logger
    tlsConfig *tls.Config
    maxConns int
}

type Option func(*Server)

func WithAddress(addr string) Option {
    return func(s *Server) {
        s.addr = addr
    }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func WithTLS(config *tls.Config) Option {
    return func(s *Server) {
        s.tlsConfig = config
    }
}

func WithMaxConnections(n int) Option {
    return func(s *Server) {
        s.maxConns = n
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        addr:     ":8080",
        timeout:  30 * time.Second,
        logger:   log.Default(),
        maxConns: 100,
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

// 使用示例
func optionPatternExample() {
    // 使用默认配置
    srv1 := NewServer()

    // 自定义配置
    srv2 := NewServer(
        WithAddress(":9090"),
        WithTimeout(60*time.Second),
        WithMaxConnections(200),
    )

    // 组合选项
    baseOpts := []Option{
        WithTimeout(10 * time.Second),
        WithMaxConnections(50),
    }

    srv3 := NewServer(append(baseOpts, WithAddress(":7070"))...)

    _ = srv1
    _ = srv2
    _ = srv3
}

// 反例: 传统的复杂构造函数
func NewServerBad(addr string, timeout time.Duration, logger *log.Logger,
                  tlsConfig *tls.Config, maxConns int) *Server {
    // 参数过多，调用时难以记住顺序
    return &Server{
        addr:     addr,
        timeout:  timeout,
        logger:   logger,
        tlsConfig: tlsConfig,
        maxConns: maxConns,
    }
}
```

### 2.3 对象池模式

```
对象池的形式化:
────────────────────────────────────────
Pool = (Available, InUse, Factory, Reset)
性质:
├─ 复用对象减少GC压力
├─ 限制最大对象数
└─ 对象重置后可用

Go实现:
var pool = sync.Pool{
    New: func() interface{} {
        return new(Buffer)
    },
}

func getBuffer() *Buffer {
    return pool.Get().(*Buffer)
}

func putBuffer(b *Buffer) {
    b.Reset()  // 重置状态
    pool.Put(b)
}

代码示例:
// Buffer池
type Buffer struct {
    data []byte
}

func (b *Buffer) Write(p []byte) {
    b.data = append(b.data, p...)
}

func (b *Buffer) Reset() {
    b.data = b.data[:0]  // 保留容量，重置长度
}

var bufferPool = sync.Pool{
    New: func() interface{} {
        return &Buffer{data: make([]byte, 0, 1024)}
    },
}

func processWithPool(data []byte) {
    buf := bufferPool.Get().(*Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    buf.Write(data)
    // 处理...
}

// 反例: 不使用池
func processWithoutPool(data []byte) {
    buf := &Buffer{data: make([]byte, 0, len(data))}  // 每次分配
    buf.Write(data)
    // 处理... 然后GC
}
```

---

## 三、结构型模式

### 3.1 适配器模式的隐式实现

```
适配器的形式化:
────────────────────────────────────────
Adapter: Target <: Adapter(Source)

Go特性:
接口的隐式实现消除了显式适配器需要

传统方式 (Java):
class Adapter implements Target {
    private Source source;
    public void method() { source.differentMethod(); }
}

Go方式:
type Target interface { Method() }
type Source struct{}
func (s Source) Method() { /* 直接实现 */ }

// Source 自动适配 Target，无需显式适配器

何时需要显式适配器:
├─ 第三方库类型不匹配
├─ 遗留代码接口转换
└─ 功能增强(装饰)

代码示例:
// 适配器场景: 统一不同第三方库的接口

// 第三方库A
type ThirdPartyA struct{}
func (a *ThirdPartyA) Send(data []byte) error { return nil }

// 第三方库B
type ThirdPartyB struct{}
func (b *ThirdPartyB) Write(p []byte) (n int, err error) { return len(p), nil }

// 目标接口
type Writer interface {
    Write([]byte) error
}

// 适配器A
type AdapterA struct {
    a *ThirdPartyA
}

func (ad *AdapterA) Write(p []byte) error {
    return ad.a.Send(p)
}

// 适配器B
type AdapterB struct {
    b *ThirdPartyB
}

func (ad *AdapterB) Write(p []byte) error {
    _, err := ad.b.Write(p)
    return err
}

// 统一使用
func processData(w Writer, data []byte) error {
    return w.Write(data)
}
```

### 3.2 装饰器模式的函数实现

```
装饰器的组合数学:
────────────────────────────────────────
装饰器 = 高阶函数: (a → b) → (a → b)

HTTP中间件实现:
type Handler func(w http.ResponseWriter, r *http.Request)

type Middleware func(Handler) Handler

func Logger(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL)
        next(w, r)
    }
}

func Auth(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next(w, r)
    }
}

组合:
final := Logger(Auth(Timeout(Handler)))

代数性质:
├─ 结合律: (A ∘ B) ∘ C = A ∘ (B ∘ C)
├─ 单位元: Identity装饰器
└─ 顺序重要性: 装饰器应用顺序影响行为

代码示例:
// 完整HTTP中间件链
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(wrapped, r)

        log.Printf("[%s] %s %d %v",
            r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
    })
}

func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func AuthMiddleware(apiKey string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.Header.Get("X-API-Key")
            if key != apiKey {
                http.Error(w, "Unauthorized", 401)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func Chain(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

// 使用
func decoratorExample() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })

    handler := Chain(
        RecoveryMiddleware,
        LoggingMiddleware,
        AuthMiddleware("secret-key"),
    )(mux)

    http.ListenAndServe(":8080", handler)
}
```

---

## 四、行为型模式

### 4.1 策略模式的函数式实现

```
策略模式的形式化:
────────────────────────────────────────
策略 = 算法族，可互相替换

传统OO实现:
interface Strategy { Execute() }
class ConcreteStrategyA implements Strategy { ... }
class ConcreteStrategyB implements Strategy { ... }

Go函数式实现:
type Strategy func(input string) string

var (
    UpperCase Strategy = strings.ToUpper
    LowerCase Strategy = strings.ToLower
    Reverse   Strategy = func(s string) string {
        runes := []rune(s)
        for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
            runes[i], runes[j] = runes[j], runes[i]
        }
        return string(runes)
    }
)

优势:
├─ 无类型层次结构
├─ 闭包捕获状态
├─ 内联策略定义
└─ 零开销抽象(内联优化)

代码示例:
// 排序策略
type SortStrategy func([]int) []int

func BubbleSort(data []int) []int {
    n := len(data)
    result := make([]int, len(data))
    copy(result, data)

    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    return result
}

func QuickSort(data []int) []int {
    if len(data) <= 1 {
        return data
    }

    result := make([]int, len(data))
    copy(result, data)

    pivot := result[len(result)/2]
    var left, right, equal []int

    for _, v := range result {
        if v < pivot {
            left = append(left, v)
        } else if v > pivot {
            right = append(right, v)
        } else {
            equal = append(equal, v)
        }
    }

    left = QuickSort(left)
    right = QuickSort(right)

    return append(append(left, equal...), right...)
}

type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) Sort(data []int) []int {
    return s.strategy(data)
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

// 使用
func strategyExample() {
    data := []int{64, 34, 25, 12, 22, 11, 90}

    sorter := &Sorter{strategy: BubbleSort}
    fmt.Println("Bubble:", sorter.Sort(data))

    sorter.SetStrategy(QuickSort)
    fmt.Println("Quick:", sorter.Sort(data))

    // 内联定义策略
    sorter.SetStrategy(func(d []int) []int {
        // 插入排序
        result := make([]int, len(d))
        copy(result, d)
        for i := 1; i < len(result); i++ {
            key := result[i]
            j := i - 1
            for j >= 0 && result[j] > key {
                result[j+1] = result[j]
                j--
            }
            result[j+1] = key
        }
        return result
    })
    fmt.Println("Insertion:", sorter.Sort(data))
}
```

### 4.2 观察者模式的Channel实现

```
观察者 ≡ Publish-Subscribe:
────────────────────────────────────────

形式化:
Subject = (Notify, Subscribe, Unsubscribe)
Observer = Event → Action

Go实现:
type Event struct {
    Type string
    Data interface{}
}

type Subject struct {
    observers map[chan Event]struct{}
    mu        sync.RWMutex
}

func (s *Subject) Subscribe() chan Event {
    ch := make(chan Event, 10)
    s.mu.Lock()
    s.observers[ch] = struct{}{}
    s.mu.Unlock()
    return ch
}

func (s *Subject) Notify(e Event) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for ch := range s.observers {
        select {
        case ch <- e:
        default:
            // 缓冲满，丢弃或处理
        }
    }
}

对比传统实现:
├─ 解耦: Channel作为中介
├─ 异步: 非阻塞通知
├─ 背压: 缓冲控制
└─ 取消: 通过Channel关闭

代码示例:
// 事件总线实现
type EventBus struct {
    subscribers map[string][]chan Event
    mu          sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]chan Event),
    }
}

func (eb *EventBus) Subscribe(eventType string, bufferSize int) chan Event {
    ch := make(chan Event, bufferSize)
    eb.mu.Lock()
    eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
    eb.mu.Unlock()
    return ch
}

func (eb *EventBus) Publish(event Event) {
    eb.mu.RLock()
    subs := eb.subscribers[event.Type]
    eb.mu.RUnlock()

    for _, ch := range subs {
        select {
        case ch <- event:
        default:
            // 阻塞时丢弃或记录
        }
    }
}

func (eb *EventBus) Unsubscribe(eventType string, ch chan Event) {
    eb.mu.Lock()
    defer eb.mu.Unlock()

    subs := eb.subscribers[eventType]
    for i, sub := range subs {
        if sub == ch {
            eb.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
            close(ch)
            break
        }
    }
}

// 使用示例
type OrderCreatedEvent struct {
    OrderID string
    Amount  float64
}

func observerExample() {
    bus := NewEventBus()

    // 订阅者1: 发送邮件
    emailCh := bus.Subscribe("order.created", 10)
    go func() {
        for event := range emailCh {
            fmt.Printf("发送邮件: %v\n", event.Data)
        }
    }()

    // 订阅者2: 更新库存
    inventoryCh := bus.Subscribe("order.created", 10)
    go func() {
        for event := range inventoryCh {
            fmt.Printf("更新库存: %v\n", event.Data)
        }
    }()

    // 发布事件
    bus.Publish(Event{
        Type: "order.created",
        Data: OrderCreatedEvent{OrderID: "123", Amount: 99.99},
    })
}
```

---

## 五、架构模式

### 5.1 依赖注入的显式实现

```
DI的形式化:
────────────────────────────────────────
组件定义:
Component = (Dependencies, Provided)

容器:
Container = Name → Component

解析:
Resolve: Container × Name → Instance

Go实现 (构造函数注入):
type UserService struct {
    repo UserRepository
    log  *Logger
}

func NewUserService(repo UserRepository, log *Logger) *UserService {
    return &UserService{repo: repo, log: log}
}

Wire代码生成:
//go:build wireinject
func InitializeApp() (*App, error) {
    wire.Build(
        NewUserService,
        NewDatabaseRepo,
        NewLogger,
        NewApp,
    )
    return nil, nil
}

代码示例:
// 接口定义
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

type Logger interface {
    Info(msg string)
    Error(msg string, err error)
}

// 实现
type PostgresUserRepo struct {
    db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
    return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    // 实现...
    return nil, nil
}

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
    return &ConsoleLogger{}
}

func (l *ConsoleLogger) Info(msg string) {
    fmt.Println("[INFO]", msg)
}

// 服务
type UserService struct {
    repo UserRepository
    log  Logger
}

func NewUserService(repo UserRepository, log Logger) *UserService {
    return &UserService{repo: repo, log: log}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    s.log.Info("Getting user: " + id)
    return s.repo.GetByID(ctx, id)
}

// 手动DI
func diExample() {
    db, _ := sql.Open("postgres", "...")
    repo := NewPostgresUserRepo(db)
    logger := NewConsoleLogger()
    service := NewUserService(repo, logger)

    user, _ := service.GetUser(context.Background(), "123")
    fmt.Println(user)
}
```

### 5.2 六边形架构 (Ports & Adapters)

```
架构代数:
────────────────────────────────────────
Domain = Core Logic (独立于技术细节)
Ports = Domain定义的接口
Adapters = Ports的实现

Go映射:
Domain ──────────► internal/domain/
Ports ───────────► 接口定义在domain包
Adapters ───────► internal/infrastructure/

示例:
// Domain (internal/domain/order.go)
type OrderService interface {
    Create(ctx context.Context, cmd CreateOrderCmd) (*Order, error)
}

// Port定义
type PaymentGateway interface {
    Charge(ctx context.Context, amount Money) error
}

// Adapter实现 (internal/infrastructure/stripe.go)
type StripeAdapter struct { client *stripe.Client }
func (s *StripeAdapter) Charge(...) error { ... }

边界:
├─ Domain不依赖任何外部包
├─ 依赖关系指向Domain
└─ 技术细节可替换

代码示例:
// internal/domain/order.go
type Order struct {
    ID        string
    UserID    string
    Amount    Money
    Status    OrderStatus
    CreatedAt time.Time
}

type Money struct {
    Amount   decimal.Decimal
    Currency string
}

type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
    FindByUser(ctx context.Context, userID string) ([]*Order, error)
}

type PaymentGateway interface {
    Charge(ctx context.Context, orderID string, amount Money) error
    Refund(ctx context.Context, orderID string, amount Money) error
}

type OrderService struct {
    repo    OrderRepository
    payment PaymentGateway
}

func NewOrderService(repo OrderRepository, payment PaymentGateway) *OrderService {
    return &OrderService{repo: repo, payment: payment}
}

func (s *OrderService) Create(ctx context.Context, cmd CreateOrderCmd) (*Order, error) {
    order := &Order{
        ID:     uuid.New().String(),
        UserID: cmd.UserID,
        Amount: cmd.Amount,
        Status: OrderStatusPending,
    }

    if err := s.repo.Save(ctx, order); err != nil {
        return nil, err
    }

    if err := s.payment.Charge(ctx, order.ID, order.Amount); err != nil {
        order.Status = OrderStatusFailed
        s.repo.Save(ctx, order)
        return nil, err
    }

    order.Status = OrderStatusPaid
    s.repo.Save(ctx, order)
    return order, nil
}

// internal/infrastructure/postgres/order_repo.go
type PostgresOrderRepo struct {
    db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
    return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) Save(ctx context.Context, order *Order) error {
    // SQL实现
    return nil
}

// internal/infrastructure/stripe/gateway.go
type StripeGateway struct {
    client *stripe.Client
}

func NewStripeGateway(apiKey string) *StripeGateway {
    return &StripeGateway{client: stripe.NewClient(apiKey)}
}

func (g *StripeGateway) Charge(ctx context.Context, orderID string, amount Money) error {
    // Stripe API调用
    return nil
}
```

---

## 六、模式选择决策

```
模式选择决策树:
────────────────────────────────────────

需要状态封装?
├── 是 → 需要多态行为?
│       ├── 是 → 接口 + 结构体 (策略/状态)
│       └── 否 → 函数闭包
└── 否 → 纯函数实现

需要并发?
├── 是 → 共享状态?
│       ├── 是 → Mutex保护 + 封装
│       └── 否 → Channel通信
└── 否 → 顺序实现

需要扩展点?
├── 是 → 运行时扩展?
│       ├── 是 → 接口/插件系统
│       └── 否 → 编译期泛型
└── 否 → 直接实现
```

---

*本章将经典设计模式与Go特性结合，提供了形式化的模式分析、丰富的代码示例、反例对比和选型框架。*
