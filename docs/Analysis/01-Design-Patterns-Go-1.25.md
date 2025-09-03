# 1 1 1 1 1 1 1 Go 1.25 设计模式综合分析

<!-- TOC START -->
- [1 1 1 1 1 1 1 Go 1.25 设计模式综合分析](#1-1-1-1-1-1-1-go-125-设计模式综合分析)
  - [1.1 目录](#目录)
  - [1.2 创建型模式](#创建型模式)
    - [1.2.1 单例模式 (Singleton)](#单例模式-singleton)
    - [1.2.2 工厂模式 (Factory)](#工厂模式-factory)
    - [1.2.3 建造者模式 (Builder)](#建造者模式-builder)
  - [1.3 结构型模式](#结构型模式)
    - [1.3.1 适配器模式 (Adapter)](#适配器模式-adapter)
    - [1.3.2 装饰器模式 (Decorator)](#装饰器模式-decorator)
    - [1.3.3 代理模式 (Proxy)](#代理模式-proxy)
  - [1.4 行为型模式](#行为型模式)
    - [1.4.1 观察者模式 (Observer)](#观察者模式-observer)
    - [1.4.2 策略模式 (Strategy)](#策略模式-strategy)
    - [1.4.3 命令模式 (Command)](#命令模式-command)
  - [1.5 并发型模式](#并发型模式)
    - [1.5.1 工作池模式 (Worker Pool)](#工作池模式-worker-pool)
    - [1.5.2 发布订阅模式 (Pub/Sub)](#发布订阅模式-pubsub)
    - [1.5.3 管道模式 (Pipeline)](#管道模式-pipeline)
  - [1.6 云原生模式](#云原生模式)
    - [1.6.1 健康检查模式 (Health Check)](#健康检查模式-health-check)
    - [1.6.2 配置管理模式 (Configuration Management)](#配置管理模式-configuration-management)
  - [1.7 性能优化模式](#性能优化模式)
    - [1.7.1 对象池模式 (Object Pool)](#对象池模式-object-pool)
    - [1.7.2 缓存模式 (Cache)](#缓存模式-cache)
  - [1.8 总结](#总结)
<!-- TOC END -->














## 1.1 目录

1. [创建型模式](#创建型模式)
2. [结构型模式](#结构型模式)
3. [行为型模式](#行为型模式)
4. [并发型模式](#并发型模式)
5. [云原生模式](#云原生模式)
6. [性能优化模式](#性能优化模式)

## 1.2 创建型模式

### 1.2.1 单例模式 (Singleton)

```go
// 线程安全的单例模式
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "initialized",
        }
    })
    return instance
}

// 泛型单例
type GenericSingleton[T any] struct {
    data T
}

var (
    genericInstances = make(map[string]interface{})
    genericOnce      sync.RWMutex
)

func GetGenericInstance[T any](key string, factory func() T) T {
    genericOnce.RLock()
    if instance, exists := genericInstances[key]; exists {
        genericOnce.RUnlock()
        return instance.(T)
    }
    genericOnce.RUnlock()
    
    genericOnce.Lock()
    defer genericOnce.Unlock()
    
    // 双重检查
    if instance, exists := genericInstances[key]; exists {
        return instance.(T)
    }
    
    newInstance := factory()
    genericInstances[key] = newInstance
    return newInstance
}
```

### 1.2.2 工厂模式 (Factory)

```go
// 抽象工厂模式
type Product interface {
    Use() string
}

type ProductFactory interface {
    CreateProduct() Product
}

type ConcreteProductA struct{}
type ConcreteProductB struct{}

func (p *ConcreteProductA) Use() string { return "Product A" }
func (p *ConcreteProductB) Use() string { return "Product B" }

type ConcreteFactoryA struct{}
type ConcreteFactoryB struct{}

func (f *ConcreteFactoryA) CreateProduct() Product { return &ConcreteProductA{} }
func (f *ConcreteFactoryB) CreateProduct() Product { return &ConcreteProductB{} }

// 泛型工厂
type GenericFactory[T Product] interface {
    Create() T
}

type GenericConcreteFactory[T Product] struct {
    creator func() T
}

func (f *GenericConcreteFactory[T]) Create() T {
    return f.creator()
}

func NewGenericFactory[T Product](creator func() T) GenericFactory[T] {
    return &GenericConcreteFactory[T]{creator: creator}
}
```

### 1.2.3 建造者模式 (Builder)

```go
// 建造者模式
type Computer struct {
    CPU       string
    Memory    string
    Storage   string
    Graphics  string
}

type ComputerBuilder struct {
    computer *Computer
}

func NewComputerBuilder() *ComputerBuilder {
    return &ComputerBuilder{
        computer: &Computer{},
    }
}

func (b *ComputerBuilder) SetCPU(cpu string) *ComputerBuilder {
    b.computer.CPU = cpu
    return b
}

func (b *ComputerBuilder) SetMemory(memory string) *ComputerBuilder {
    b.computer.Memory = memory
    return b
}

func (b *ComputerBuilder) SetStorage(storage string) *ComputerBuilder {
    b.computer.Storage = storage
    return b
}

func (b *ComputerBuilder) SetGraphics(graphics string) *ComputerBuilder {
    b.computer.Graphics = graphics
    return b
}

func (b *ComputerBuilder) Build() *Computer {
    return b.computer
}
```

## 1.3 结构型模式

### 1.3.1 适配器模式 (Adapter)

```go
// 适配器模式
type LegacySystem interface {
    OldMethod() string
}

type NewSystem interface {
    NewMethod() string
}

type LegacyImplementation struct{}

func (l *LegacyImplementation) OldMethod() string {
    return "legacy result"
}

type Adapter struct {
    legacy LegacySystem
}

func NewAdapter(legacy LegacySystem) *Adapter {
    return &Adapter{legacy: legacy}
}

func (a *Adapter) NewMethod() string {
    // 适配旧接口到新接口
    oldResult := a.legacy.OldMethod()
    return "adapted: " + oldResult
}
```

### 1.3.2 装饰器模式 (Decorator)

```go
// 装饰器模式
type Component interface {
    Operation() string
}

type ConcreteComponent struct{}

func (c *ConcreteComponent) Operation() string {
    return "basic operation"
}

type Decorator struct {
    component Component
}

func NewDecorator(component Component) *Decorator {
    return &Decorator{component: component}
}

func (d *Decorator) Operation() string {
    return "decorated: " + d.component.Operation()
}

// 多层装饰
type LoggingDecorator struct {
    Decorator
}

func (l *LoggingDecorator) Operation() string {
    result := l.component.Operation()
    fmt.Printf("Logging: %s\n", result)
    return result
}
```

### 1.3.3 代理模式 (Proxy)

```go
// 代理模式
type Subject interface {
    Request() string
}

type RealSubject struct{}

func (r *RealSubject) Request() string {
    return "real subject response"
}

type Proxy struct {
    realSubject *RealSubject
}

func NewProxy() *Proxy {
    return &Proxy{}
}

func (p *Proxy) Request() string {
    if p.realSubject == nil {
        p.realSubject = &RealSubject{}
    }
    
    // 前置处理
    fmt.Println("Proxy: before request")
    
    result := p.realSubject.Request()
    
    // 后置处理
    fmt.Println("Proxy: after request")
    
    return result
}
```

## 1.4 行为型模式

### 1.4.1 观察者模式 (Observer)

```go
// 观察者模式
type Observer interface {
    Update(data interface{})
}

type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
}

type ConcreteSubject struct {
    observers []Observer
    data      interface{}
    mu        sync.RWMutex
}

func (s *ConcreteSubject) Attach(observer Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers = append(s.observers, observer)
}

func (s *ConcreteSubject) Detach(observer Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, obs := range s.observers {
        if obs == observer {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            break
        }
    }
}

func (s *ConcreteSubject) Notify() {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    for _, observer := range s.observers {
        observer.Update(s.data)
    }
}

func (s *ConcreteSubject) SetData(data interface{}) {
    s.mu.Lock()
    s.data = data
    s.mu.Unlock()
    s.Notify()
}
```

### 1.4.2 策略模式 (Strategy)

```go
// 策略模式
type Strategy interface {
    Execute(data interface{}) interface{}
}

type Context struct {
    strategy Strategy
}

func NewContext(strategy Strategy) *Context {
    return &Context{strategy: strategy}
}

func (c *Context) SetStrategy(strategy Strategy) {
    c.strategy = strategy
}

func (c *Context) ExecuteStrategy(data interface{}) interface{} {
    return c.strategy.Execute(data)
}

// 具体策略
type BubbleSortStrategy struct{}

func (b *BubbleSortStrategy) Execute(data interface{}) interface{} {
    // 实现冒泡排序
    return "bubble sorted"
}

type QuickSortStrategy struct{}

func (q *QuickSortStrategy) Execute(data interface{}) interface{} {
    // 实现快速排序
    return "quick sorted"
}
```

### 1.4.3 命令模式 (Command)

```go
// 命令模式
type Command interface {
    Execute()
    Undo()
}

type Receiver struct{}

func (r *Receiver) Action() {
    fmt.Println("Receiver: performing action")
}

type ConcreteCommand struct {
    receiver *Receiver
}

func NewConcreteCommand(receiver *Receiver) *ConcreteCommand {
    return &ConcreteCommand{receiver: receiver}
}

func (c *ConcreteCommand) Execute() {
    c.receiver.Action()
}

func (c *ConcreteCommand) Undo() {
    fmt.Println("Command: undoing action")
}

type Invoker struct {
    commands []Command
}

func (i *Invoker) AddCommand(command Command) {
    i.commands = append(i.commands, command)
}

func (i *Invoker) ExecuteCommands() {
    for _, command := range i.commands {
        command.Execute()
    }
}
```

## 1.5 并发型模式

### 1.5.1 工作池模式 (Worker Pool)

```go
// 工作池模式
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultChan chan Result
    wg         sync.WaitGroup
}

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID  int
    Data   interface{}
    Error  error
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, workers*2),
        resultChan: make(chan Result, workers*2),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    
    for job := range wp.jobQueue {
        result := wp.processJob(job)
        wp.resultChan <- result
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // 处理具体任务
    return Result{
        JobID: job.ID,
        Data:  fmt.Sprintf("processed: %v", job.Data),
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Close() {
    close(wp.jobQueue)
    wp.wg.Wait()
    close(wp.resultChan)
}
```

### 1.5.2 发布订阅模式 (Pub/Sub)

```go
// 发布订阅模式
type Event struct {
    Topic   string
    Data    interface{}
    Time    time.Time
}

type Subscriber interface {
    OnEvent(event Event)
}

type EventBus struct {
    subscribers map[string][]Subscriber
    mu          sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]Subscriber),
    }
}

func (eb *EventBus) Subscribe(topic string, subscriber Subscriber) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.subscribers[topic] = append(eb.subscribers[topic], subscriber)
}

func (eb *EventBus) Publish(event Event) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()
    
    if subscribers, exists := eb.subscribers[event.Topic]; exists {
        for _, subscriber := range subscribers {
            go subscriber.OnEvent(event)
        }
    }
}

// 具体订阅者
type LoggingSubscriber struct{}

func (l *LoggingSubscriber) OnEvent(event Event) {
    fmt.Printf("Logging: Topic=%s, Data=%v, Time=%v\n",
        event.Topic, event.Data, event.Time)
}
```

### 1.5.3 管道模式 (Pipeline)

```go
// 管道模式
type Pipeline struct {
    stages []Stage
}

type Stage func(input interface{}) interface{}

func NewPipeline() *Pipeline {
    return &Pipeline{
        stages: make([]Stage, 0),
    }
}

func (p *Pipeline) AddStage(stage Stage) *Pipeline {
    p.stages = append(p.stages, stage)
    return p
}

func (p *Pipeline) Execute(input interface{}) interface{} {
    result := input
    
    for _, stage := range p.stages {
        result = stage(result)
    }
    
    return result
}

// 并发管道
type ConcurrentPipeline struct {
    stages []Stage
}

func (cp *ConcurrentPipeline) Execute(input interface{}) interface{} {
    if len(cp.stages) == 0 {
        return input
    }
    
    // 创建通道链
    channels := make([]chan interface{}, len(cp.stages)+1)
    for i := range channels {
        channels[i] = make(chan interface{}, 1)
    }
    
    // 启动所有阶段
    for i, stage := range cp.stages {
        go func(stage Stage, in, out chan interface{}) {
            for data := range in {
                result := stage(data)
                out <- result
            }
            close(out)
        }(stage, channels[i], channels[i+1])
    }
    
    // 发送初始数据
    channels[0] <- input
    close(channels[0])
    
    // 获取最终结果
    return <-channels[len(channels)-1]
}
```

## 1.6 云原生模式

### 1.6.1 健康检查模式 (Health Check)

```go
// 健康检查模式
type HealthChecker interface {
    Check() HealthStatus
}

type HealthStatus struct {
    Status    string `json:"status"`
    Message   string `json:"message"`
    Timestamp time.Time `json:"timestamp"`
}

type HealthCheckRegistry struct {
    checkers map[string]HealthChecker
    mu       sync.RWMutex
}

func NewHealthCheckRegistry() *HealthCheckRegistry {
    return &HealthCheckRegistry{
        checkers: make(map[string]HealthChecker),
    }
}

func (hcr *HealthCheckRegistry) Register(name string, checker HealthChecker) {
    hcr.mu.Lock()
    defer hcr.mu.Unlock()
    hcr.checkers[name] = checker
}

func (hcr *HealthCheckRegistry) CheckAll() map[string]HealthStatus {
    hcr.mu.RLock()
    defer hcr.mu.RUnlock()
    
    results := make(map[string]HealthStatus)
    
    for name, checker := range hcr.checkers {
        results[name] = checker.Check()
    }
    
    return results
}

// 数据库健康检查
type DatabaseHealthChecker struct {
    db *sql.DB
}

func (dhc *DatabaseHealthChecker) Check() HealthStatus {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := dhc.db.PingContext(ctx); err != nil {
        return HealthStatus{
            Status:    "unhealthy",
            Message:   err.Error(),
            Timestamp: time.Now(),
        }
    }
    
    return HealthStatus{
        Status:    "healthy",
        Message:   "database connection ok",
        Timestamp: time.Now(),
    }
}
```

### 1.6.2 配置管理模式 (Configuration Management)

```go
// 配置管理模式
type Configuration struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func NewConfiguration() *Configuration {
    return &Configuration{
        data: make(map[string]interface{}),
    }
}

func (c *Configuration) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func (c *Configuration) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}

func (c *Configuration) GetString(key string) (string, error) {
    value, exists := c.Get(key)
    if !exists {
        return "", fmt.Errorf("key %s not found", key)
    }
    
    if str, ok := value.(string); ok {
        return str, nil
    }
    
    return "", fmt.Errorf("value for key %s is not a string", key)
}

// 配置热重载
type ConfigWatcher struct {
    config     *Configuration
    filePath   string
    stopChan   chan struct{}
}

func NewConfigWatcher(config *Configuration, filePath string) *ConfigWatcher {
    return &ConfigWatcher{
        config:   config,
        filePath: filePath,
        stopChan: make(chan struct{}),
    }
}

func (cw *ConfigWatcher) Start() {
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                cw.reloadConfig()
            case <-cw.stopChan:
                return
            }
        }
    }()
}

func (cw *ConfigWatcher) Stop() {
    close(cw.stopChan)
}

func (cw *ConfigWatcher) reloadConfig() {
    // 重新加载配置文件
    data, err := os.ReadFile(cw.filePath)
    if err != nil {
        fmt.Printf("Failed to reload config: %v\n", err)
        return
    }
    
    var configData map[string]interface{}
    if err := json.Unmarshal(data, &configData); err != nil {
        fmt.Printf("Failed to parse config: %v\n", err)
        return
    }
    
    // 更新配置
    for key, value := range configData {
        cw.config.Set(key, value)
    }
    
    fmt.Println("Configuration reloaded successfully")
}
```

## 1.7 性能优化模式

### 1.7.1 对象池模式 (Object Pool)

```go
// 对象池模式
type ObjectPool[T any] struct {
    pool sync.Pool
    new  func() T
}

func NewObjectPool[T any](newFunc func() T) *ObjectPool[T] {
    return &ObjectPool[T]{
        pool: sync.Pool{
            New: func() interface{} {
                return newFunc()
            },
        },
        new: newFunc,
    }
}

func (op *ObjectPool[T]) Get() T {
    return op.pool.Get().(T)
}

func (op *ObjectPool[T]) Put(obj T) {
    op.pool.Put(obj)
}

// 连接池
type ConnectionPool struct {
    connections chan *Connection
    factory     func() *Connection
    maxSize     int
}

type Connection struct {
    id   string
    data interface{}
}

func NewConnectionPool(factory func() *Connection, maxSize int) *ConnectionPool {
    return &ConnectionPool{
        connections: make(chan *Connection, maxSize),
        factory:     factory,
        maxSize:     maxSize,
    }
}

func (cp *ConnectionPool) Get() *Connection {
    select {
    case conn := <-cp.connections:
        return conn
    default:
        return cp.factory()
    }
}

func (cp *ConnectionPool) Put(conn *Connection) {
    select {
    case cp.connections <- conn:
    default:
        // 池已满，丢弃连接
    }
}
```

### 1.7.2 缓存模式 (Cache)

```go
// LRU 缓存
type LRUCache[K comparable, V any] struct {
    capacity int
    cache    map[K]*Node[K, V]
    head     *Node[K, V]
    tail     *Node[K, V]
    mu       sync.RWMutex
}

type Node[K comparable, V any] struct {
    key   K
    value V
    prev  *Node[K, V]
    next  *Node[K, V]
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
    cache := &LRUCache[K, V]{
        capacity: capacity,
        cache:    make(map[K]*Node[K, V]),
    }
    
    cache.head = &Node[K, V]{}
    cache.tail = &Node[K, V]{}
    cache.head.next = cache.tail
    cache.tail.prev = cache.head
    
    return cache
}

func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if node, exists := lru.cache[key]; exists {
        lru.moveToFront(node)
        return node.value, true
    }
    
    var zero V
    return zero, false
}

func (lru *LRUCache[K, V]) Put(key K, value V) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if node, exists := lru.cache[key]; exists {
        node.value = value
        lru.moveToFront(node)
        return
    }
    
    node := &Node[K, V]{
        key:   key,
        value: value,
    }
    
    lru.cache[key] = node
    lru.addToFront(node)
    
    if len(lru.cache) > lru.capacity {
        lru.removeLRU()
    }
}

func (lru *LRUCache[K, V]) moveToFront(node *Node[K, V]) {
    lru.removeNode(node)
    lru.addToFront(node)
}

func (lru *LRUCache[K, V]) addToFront(node *Node[K, V]) {
    node.next = lru.head.next
    node.prev = lru.head
    lru.head.next.prev = node
    lru.head.next = node
}

func (lru *LRUCache[K, V]) removeNode(node *Node[K, V]) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (lru *LRUCache[K, V]) removeLRU() {
    lruNode := lru.tail.prev
    lru.removeNode(lruNode)
    delete(lru.cache, lruNode.key)
}
```

## 1.8 总结

本文档全面介绍了 Go 1.25 中的各种设计模式，包括：

1. **创建型模式**：单例、工厂、建造者等
2. **结构型模式**：适配器、装饰器、代理等
3. **行为型模式**：观察者、策略、命令等
4. **并发型模式**：工作池、发布订阅、管道等
5. **云原生模式**：健康检查、配置管理等
6. **性能优化模式**：对象池、缓存等

这些模式充分利用了 Go 1.25 的新特性，如泛型、改进的并发模型等，为构建高性能、可维护的 Go 应用程序提供了最佳实践。
