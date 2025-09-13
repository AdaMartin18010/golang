# Go语言设计模式体系

<!-- TOC START -->
- [Go语言设计模式体系](#go语言设计模式体系)
  - [1.1 🏭 创建型模式](#11--创建型模式)
    - [1.1.1 工厂模式](#111-工厂模式)
    - [1.1.2 建造者模式](#112-建造者模式)
    - [1.1.3 单例模式](#113-单例模式)
  - [1.2 🏗️ 结构型模式](#12-️-结构型模式)
    - [1.2.1 适配器模式](#121-适配器模式)
    - [1.2.2 装饰器模式](#122-装饰器模式)
    - [1.2.3 代理模式](#123-代理模式)
  - [1.3 🎭 行为型模式](#13--行为型模式)
    - [1.3.1 观察者模式](#131-观察者模式)
    - [1.3.2 策略模式](#132-策略模式)
    - [1.3.3 命令模式](#133-命令模式)
  - [1.4 ⚡ 并发模式](#14--并发模式)
    - [1.4.1 生产者消费者模式](#141-生产者消费者模式)
    - [1.4.2 工作池模式](#142-工作池模式)
    - [1.4.3 管道模式](#143-管道模式)
  - [1.5 🎯 模式应用指南](#15--模式应用指南)
    - [1.5.1 模式选择决策树](#151-模式选择决策树)
    - [1.5.2 模式组合策略](#152-模式组合策略)
<!-- TOC END -->

## 1.1 🏭 创建型模式

### 1.1.1 工厂模式

**简单工厂**:

```go
// 产品接口
type PaymentProcessor interface {
    ProcessPayment(amount float64) error
}

// 具体产品
type CreditCardProcessor struct{}

func (ccp *CreditCardProcessor) ProcessPayment(amount float64) error {
    fmt.Printf("Processing credit card payment: $%.2f\n", amount)
    return nil
}

type PayPalProcessor struct{}

func (ppp *PayPalProcessor) ProcessPayment(amount float64) error {
    fmt.Printf("Processing PayPal payment: $%.2f\n", amount)
    return nil
}

// 简单工厂
type PaymentFactory struct{}

func (pf *PaymentFactory) CreateProcessor(paymentType string) (PaymentProcessor, error) {
    switch paymentType {
    case "creditcard":
        return &CreditCardProcessor{}, nil
    case "paypal":
        return &PayPalProcessor{}, nil
    default:
        return nil, fmt.Errorf("unsupported payment type: %s", paymentType)
    }
}
```

**抽象工厂**:

```go
// 抽象工厂接口
type DatabaseFactory interface {
    CreateConnection() DatabaseConnection
    CreateQuery() DatabaseQuery
}

// 具体工厂
type MySQLFactory struct{}

func (mf *MySQLFactory) CreateConnection() DatabaseConnection {
    return &MySQLConnection{}
}

func (mf *MySQLFactory) CreateQuery() DatabaseQuery {
    return &MySQLQuery{}
}

type PostgreSQLFactory struct{}

func (pf *PostgreSQLFactory) CreateConnection() DatabaseConnection {
    return &PostgreSQLConnection{}
}

func (pf *PostgreSQLFactory) CreateQuery() DatabaseQuery {
    return &PostgreSQLQuery{}
}

// 产品接口
type DatabaseConnection interface {
    Connect() error
    Disconnect() error
}

type DatabaseQuery interface {
    Execute(query string) (interface{}, error)
}

// 具体产品
type MySQLConnection struct{}

func (mc *MySQLConnection) Connect() error {
    fmt.Println("Connecting to MySQL")
    return nil
}

func (mc *MySQLConnection) Disconnect() error {
    fmt.Println("Disconnecting from MySQL")
    return nil
}

type MySQLQuery struct{}

func (mq *MySQLQuery) Execute(query string) (interface{}, error) {
    fmt.Printf("Executing MySQL query: %s\n", query)
    return nil, nil
}
```

### 1.1.2 建造者模式

```go
// 产品
type Computer struct {
    CPU    string
    Memory string
    Storage string
    GPU    string
}

// 建造者接口
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetStorage(storage string) ComputerBuilder
    SetGPU(gpu string) ComputerBuilder
    Build() *Computer
}

// 具体建造者
type GamingComputerBuilder struct {
    computer *Computer
}

func NewGamingComputerBuilder() *GamingComputerBuilder {
    return &GamingComputerBuilder{
        computer: &Computer{},
    }
}

func (gcb *GamingComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    gcb.computer.CPU = cpu
    return gcb
}

func (gcb *GamingComputerBuilder) SetMemory(memory string) ComputerBuilder {
    gcb.computer.Memory = memory
    return gcb
}

func (gcb *GamingComputerBuilder) SetStorage(storage string) ComputerBuilder {
    gcb.computer.Storage = storage
    return gcb
}

func (gcb *GamingComputerBuilder) SetGPU(gpu string) ComputerBuilder {
    gcb.computer.GPU = gpu
    return gcb
}

func (gcb *GamingComputerBuilder) Build() *Computer {
    return gcb.computer
}

// 指挥者
type ComputerDirector struct {
    builder ComputerBuilder
}

func (cd *ComputerDirector) SetBuilder(builder ComputerBuilder) {
    cd.builder = builder
}

func (cd *ComputerDirector) BuildGamingComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i9-12900K").
        SetMemory("32GB DDR5").
        SetStorage("1TB NVMe SSD").
        SetGPU("RTX 4080").
        Build()
}

func (cd *ComputerDirector) BuildOfficeComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i5-12400").
        SetMemory("16GB DDR4").
        SetStorage("512GB SSD").
        SetGPU("Integrated").
        Build()
}
```

### 1.1.3 单例模式

```go
// 线程安全的单例
type DatabaseManager struct {
    connection *sql.DB
}

var (
    instance *DatabaseManager
    once     sync.Once
)

func GetDatabaseManager() *DatabaseManager {
    once.Do(func() {
        instance = &DatabaseManager{}
        instance.initialize()
    })
    return instance
}

func (dm *DatabaseManager) initialize() {
    // 初始化数据库连接
    db, err := sql.Open("postgres", "connection_string")
    if err != nil {
        panic(err)
    }
    dm.connection = db
}

func (dm *DatabaseManager) GetConnection() *sql.DB {
    return dm.connection
}

// 使用sync.Mutex的单例
type ConfigManager struct {
    config map[string]interface{}
    mu     sync.RWMutex
}

var (
    configInstance *ConfigManager
    configOnce     sync.Once
)

func GetConfigManager() *ConfigManager {
    configOnce.Do(func() {
        configInstance = &ConfigManager{
            config: make(map[string]interface{}),
        }
        configInstance.loadConfig()
    })
    return configInstance
}

func (cm *ConfigManager) Get(key string) (interface{}, bool) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    value, exists := cm.config[key]
    return value, exists
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.config[key] = value
}

func (cm *ConfigManager) loadConfig() {
    // 加载配置
    cm.config["database_url"] = "postgres://localhost/mydb"
    cm.config["redis_url"] = "redis://localhost:6379"
}
```

## 1.2 🏗️ 结构型模式

### 1.2.1 适配器模式

```go
// 目标接口
type MediaPlayer interface {
    Play(audioType string, fileName string)
}

// 被适配者
type AdvancedMediaPlayer interface {
    PlayVlc(fileName string)
    PlayMp4(fileName string)
}

type VlcPlayer struct{}

func (vp *VlcPlayer) PlayVlc(fileName string) {
    fmt.Printf("Playing vlc file: %s\n", fileName)
}

func (vp *VlcPlayer) PlayMp4(fileName string) {
    // VLC播放器不支持MP4
}

type Mp4Player struct{}

func (mp *Mp4Player) PlayVlc(fileName string) {
    // MP4播放器不支持VLC
}

func (mp *Mp4Player) PlayMp4(fileName string) {
    fmt.Printf("Playing mp4 file: %s\n", fileName)
}

// 适配器
type MediaAdapter struct {
    advancedPlayer AdvancedMediaPlayer
}

func NewMediaAdapter(audioType string) *MediaAdapter {
    switch audioType {
    case "vlc":
        return &MediaAdapter{advancedPlayer: &VlcPlayer{}}
    case "mp4":
        return &MediaAdapter{advancedPlayer: &Mp4Player{}}
    default:
        return nil
    }
}

func (ma *MediaAdapter) Play(audioType string, fileName string) {
    switch audioType {
    case "vlc":
        ma.advancedPlayer.PlayVlc(fileName)
    case "mp4":
        ma.advancedPlayer.PlayMp4(fileName)
    }
}

// 客户端
type AudioPlayer struct {
    mediaAdapter *MediaAdapter
}

func (ap *AudioPlayer) Play(audioType string, fileName string) {
    switch audioType {
    case "mp3":
        fmt.Printf("Playing mp3 file: %s\n", fileName)
    case "vlc", "mp4":
        ap.mediaAdapter = NewMediaAdapter(audioType)
        ap.mediaAdapter.Play(audioType, fileName)
    default:
        fmt.Printf("Invalid media type: %s\n", audioType)
    }
}
```

### 1.2.2 装饰器模式

```go
// 组件接口
type Coffee interface {
    GetDescription() string
    GetCost() float64
}

// 具体组件
type SimpleCoffee struct{}

func (sc *SimpleCoffee) GetDescription() string {
    return "Simple coffee"
}

func (sc *SimpleCoffee) GetCost() float64 {
    return 2.0
}

// 装饰器基类
type CoffeeDecorator struct {
    coffee Coffee
}

func (cd *CoffeeDecorator) GetDescription() string {
    return cd.coffee.GetDescription()
}

func (cd *CoffeeDecorator) GetCost() float64 {
    return cd.coffee.GetCost()
}

// 具体装饰器
type MilkDecorator struct {
    CoffeeDecorator
}

func NewMilkDecorator(coffee Coffee) *MilkDecorator {
    return &MilkDecorator{
        CoffeeDecorator: CoffeeDecorator{coffee: coffee},
    }
}

func (md *MilkDecorator) GetDescription() string {
    return md.coffee.GetDescription() + ", milk"
}

func (md *MilkDecorator) GetCost() float64 {
    return md.coffee.GetCost() + 0.5
}

type SugarDecorator struct {
    CoffeeDecorator
}

func NewSugarDecorator(coffee Coffee) *SugarDecorator {
    return &SugarDecorator{
        CoffeeDecorator: CoffeeDecorator{coffee: coffee},
    }
}

func (sd *SugarDecorator) GetDescription() string {
    return sd.coffee.GetDescription() + ", sugar"
}

func (sd *SugarDecorator) GetCost() float64 {
    return sd.coffee.GetCost() + 0.2
}

// 使用示例
func ExampleDecorator() {
    coffee := &SimpleCoffee{}
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffee.GetCost(), coffee.GetDescription())
    
    coffeeWithMilk := NewMilkDecorator(coffee)
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffeeWithMilk.GetCost(), coffeeWithMilk.GetDescription())
    
    coffeeWithMilkAndSugar := NewSugarDecorator(coffeeWithMilk)
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffeeWithMilkAndSugar.GetCost(), coffeeWithMilkAndSugar.GetDescription())
}
```

### 1.2.3 代理模式

```go
// 主题接口
type Image interface {
    Display()
}

// 真实主题
type RealImage struct {
    fileName string
}

func NewRealImage(fileName string) *RealImage {
    ri := &RealImage{fileName: fileName}
    ri.loadFromDisk()
    return ri
}

func (ri *RealImage) loadFromDisk() {
    fmt.Printf("Loading image from disk: %s\n", ri.fileName)
}

func (ri *RealImage) Display() {
    fmt.Printf("Displaying image: %s\n", ri.fileName)
}

// 代理
type ProxyImage struct {
    realImage *RealImage
    fileName  string
}

func NewProxyImage(fileName string) *ProxyImage {
    return &ProxyImage{fileName: fileName}
}

func (pi *ProxyImage) Display() {
    if pi.realImage == nil {
        pi.realImage = NewRealImage(pi.fileName)
    }
    pi.realImage.Display()
}

// 虚拟代理
type VirtualProxy struct {
    realSubject RealSubject
    loaded      bool
}

type RealSubject struct {
    data string
}

func (rs *RealSubject) Request() string {
    return rs.data
}

func (vp *VirtualProxy) Request() string {
    if !vp.loaded {
        vp.realSubject = RealSubject{data: "Heavy data loaded"}
        vp.loaded = true
        fmt.Println("Real subject loaded")
    }
    return vp.realSubject.Request()
}
```

## 1.3 🎭 行为型模式

### 1.3.1 观察者模式

```go
// 观察者接口
type Observer interface {
    Update(message string)
}

// 主题接口
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify(message string)
}

// 具体主题
type NewsAgency struct {
    observers []Observer
    news      string
}

func (na *NewsAgency) Attach(observer Observer) {
    na.observers = append(na.observers, observer)
}

func (na *NewsAgency) Detach(observer Observer) {
    for i, obs := range na.observers {
        if obs == observer {
            na.observers = append(na.observers[:i], na.observers[i+1:]...)
            break
        }
    }
}

func (na *NewsAgency) Notify(message string) {
    for _, observer := range na.observers {
        observer.Update(message)
    }
}

func (na *NewsAgency) SetNews(news string) {
    na.news = news
    na.Notify(news)
}

// 具体观察者
type NewsChannel struct {
    name string
}

func (nc *NewsChannel) Update(message string) {
    fmt.Printf("%s received news: %s\n", nc.name, message)
}

// 使用示例
func ExampleObserver() {
    agency := &NewsAgency{}
    
    channel1 := &NewsChannel{name: "CNN"}
    channel2 := &NewsChannel{name: "BBC"}
    
    agency.Attach(channel1)
    agency.Attach(channel2)
    
    agency.SetNews("Breaking: Important news!")
}
```

### 1.3.2 策略模式

```go
// 策略接口
type PaymentStrategy interface {
    Pay(amount float64) error
}

// 具体策略
type CreditCardStrategy struct {
    cardNumber string
    cvv        string
}

func (ccs *CreditCardStrategy) Pay(amount float64) error {
    fmt.Printf("Paid $%.2f using credit card ending in %s\n", amount, ccs.cardNumber[len(ccs.cardNumber)-4:])
    return nil
}

type PayPalStrategy struct {
    email string
}

func (pps *PayPalStrategy) Pay(amount float64) error {
    fmt.Printf("Paid $%.2f using PayPal account: %s\n", amount, pps.email)
    return nil
}

type BankTransferStrategy struct {
    accountNumber string
}

func (bts *BankTransferStrategy) Pay(amount float64) error {
    fmt.Printf("Paid $%.2f using bank transfer to account: %s\n", amount, bts.accountNumber)
    return nil
}

// 上下文
type PaymentContext struct {
    strategy PaymentStrategy
}

func (pc *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    pc.strategy = strategy
}

func (pc *PaymentContext) ExecutePayment(amount float64) error {
    if pc.strategy == nil {
        return fmt.Errorf("no payment strategy set")
    }
    return pc.strategy.Pay(amount)
}

// 使用示例
func ExampleStrategy() {
    context := &PaymentContext{}
    
    // 使用信用卡支付
    context.SetStrategy(&CreditCardStrategy{
        cardNumber: "1234567890123456",
        cvv:        "123",
    })
    context.ExecutePayment(100.0)
    
    // 切换到PayPal支付
    context.SetStrategy(&PayPalStrategy{
        email: "user@example.com",
    })
    context.ExecutePayment(50.0)
}
```

### 1.3.3 命令模式

```go
// 命令接口
type Command interface {
    Execute()
    Undo()
}

// 具体命令
type LightOnCommand struct {
    light *Light
}

func (loc *LightOnCommand) Execute() {
    loc.light.On()
}

func (loc *LightOnCommand) Undo() {
    loc.light.Off()
}

type LightOffCommand struct {
    light *Light
}

func (lofc *LightOffCommand) Execute() {
    lofc.light.Off()
}

func (lofc *LightOffCommand) Undo() {
    lofc.light.On()
}

// 接收者
type Light struct {
    isOn bool
}

func (l *Light) On() {
    l.isOn = true
    fmt.Println("Light is on")
}

func (l *Light) Off() {
    l.isOn = false
    fmt.Println("Light is off")
}

// 调用者
type RemoteControl struct {
    commands map[string]Command
    lastCommand Command
}

func NewRemoteControl() *RemoteControl {
    return &RemoteControl{
        commands: make(map[string]Command),
    }
}

func (rc *RemoteControl) SetCommand(slot string, command Command) {
    rc.commands[slot] = command
}

func (rc *RemoteControl) PressButton(slot string) {
    if command, exists := rc.commands[slot]; exists {
        command.Execute()
        rc.lastCommand = command
    }
}

func (rc *RemoteControl) PressUndo() {
    if rc.lastCommand != nil {
        rc.lastCommand.Undo()
    }
}
```

## 1.4 ⚡ 并发模式

### 1.4.1 生产者消费者模式

```go
// 生产者
type Producer struct {
    id       int
    dataChan chan<- Data
    done     <-chan struct{}
}

type Data struct {
    ID   int
    Data string
}

func (p *Producer) Start() {
    go func() {
        defer close(p.dataChan)
        
        for i := 0; ; i++ {
            select {
            case p.dataChan <- Data{ID: i, Data: fmt.Sprintf("Producer %d: Data %d", p.id, i)}:
                time.Sleep(100 * time.Millisecond)
            case <-p.done:
                return
            }
        }
    }()
}

// 消费者
type Consumer struct {
    id       int
    dataChan <-chan Data
    done     <-chan struct{}
}

func (c *Consumer) Start() {
    go func() {
        for {
            select {
            case data, ok := <-c.dataChan:
                if !ok {
                    return
                }
                c.processData(data)
            case <-c.done:
                return
            }
        }
    }()
}

func (c *Consumer) processData(data Data) {
    fmt.Printf("Consumer %d processing: %s\n", c.id, data.Data)
    time.Sleep(200 * time.Millisecond)
}

// 协调器
type ProducerConsumerSystem struct {
    producers []*Producer
    consumers []*Consumer
    dataChan  chan Data
    done      chan struct{}
}

func NewProducerConsumerSystem(producerCount, consumerCount, bufferSize int) *ProducerConsumerSystem {
    dataChan := make(chan Data, bufferSize)
    done := make(chan struct{})
    
    pcs := &ProducerConsumerSystem{
        dataChan: dataChan,
        done:     done,
    }
    
    // 创建生产者
    for i := 0; i < producerCount; i++ {
        producer := &Producer{
            id:       i,
            dataChan: dataChan,
            done:     done,
        }
        pcs.producers = append(pcs.producers, producer)
    }
    
    // 创建消费者
    for i := 0; i < consumerCount; i++ {
        consumer := &Consumer{
            id:       i,
            dataChan: dataChan,
            done:     done,
        }
        pcs.consumers = append(pcs.consumers, consumer)
    }
    
    return pcs
}

func (pcs *ProducerConsumerSystem) Start() {
    // 启动生产者
    for _, producer := range pcs.producers {
        producer.Start()
    }
    
    // 启动消费者
    for _, consumer := range pcs.consumers {
        consumer.Start()
    }
}

func (pcs *ProducerConsumerSystem) Stop() {
    close(pcs.done)
}
```

### 1.4.2 工作池模式

```go
// 工作池
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultChan chan Result
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Data  interface{}
    Error error
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, queueSize),
        resultChan: make(chan Result, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case job := <-wp.jobQueue:
            result := wp.processJob(job)
            select {
            case wp.resultChan <- result:
            case <-wp.ctx.Done():
                return
            }
        case <-wp.ctx.Done():
            return
        }
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // 模拟工作处理
    time.Sleep(100 * time.Millisecond)
    
    return Result{
        JobID: job.ID,
        Data:  fmt.Sprintf("Processed job %d", job.ID),
        Error: nil,
    }
}

func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobQueue <- job:
        return nil
    case <-wp.ctx.Done():
        return wp.ctx.Err()
    default:
        return fmt.Errorf("job queue is full")
    }
}

func (wp *WorkerPool) GetResult() <-chan Result {
    return wp.resultChan
}

func (wp *WorkerPool) Stop() {
    wp.cancel()
    wp.wg.Wait()
    close(wp.jobQueue)
    close(wp.resultChan)
}
```

### 1.4.3 管道模式

```go
// 管道阶段
type PipelineStage[T, U any] interface {
    Process(input <-chan T) <-chan U
}

// Map阶段
type MapStage[T, U any] struct {
    mapper func(T) U
}

func (ms *MapStage[T, U]) Process(input <-chan T) <-chan U {
    output := make(chan U, 100)
    
    go func() {
        defer close(output)
        for item := range input {
            output <- ms.mapper(item)
        }
    }()
    
    return output
}

// Filter阶段
type FilterStage[T any] struct {
    predicate func(T) bool
}

func (fs *FilterStage[T]) Process(input <-chan T) <-chan T {
    output := make(chan T, 100)
    
    go func() {
        defer close(output)
        for item := range input {
            if fs.predicate(item) {
                output <- item
            }
        }
    }()
    
    return output
}

// 管道
type Pipeline[T any] struct {
    stages []PipelineStage[T, T]
}

func NewPipeline[T any](stages ...PipelineStage[T, T]) *Pipeline[T] {
    return &Pipeline[T]{stages: stages}
}

func (p *Pipeline[T]) Process(input <-chan T) <-chan T {
    current := input
    
    for _, stage := range p.stages {
        current = stage.Process(current)
    }
    
    return current
}

// 使用示例
func ExamplePipeline() {
    input := make(chan int, 100)
    
    pipeline := NewPipeline(
        &FilterStage[int]{predicate: func(x int) bool { return x%2 == 0 }},
        &MapStage[int, int]{mapper: func(x int) int { return x * 2 }},
    )
    
    output := pipeline.Process(input)
    
    // 发送数据
    go func() {
        for i := 1; i <= 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // 接收结果
    for result := range output {
        fmt.Println(result) // 输出: 4, 8, 12, 16, 20
    }
}
```

## 1.5 🎯 模式应用指南

### 1.5.1 模式选择决策树

```go
// 模式选择决策
type PatternDecision struct {
    Problem     string
    Constraints []string
    Recommended []string
    Avoid       []string
}

var PatternDecisions = []PatternDecision{
    {
        Problem:     "需要创建复杂对象",
        Constraints: []string{"对象构建步骤复杂", "需要不同配置"},
        Recommended: []string{"建造者模式", "抽象工厂模式"},
        Avoid:       []string{"简单工厂模式"},
    },
    {
        Problem:     "需要适配不兼容接口",
        Constraints: []string{"接口不匹配", "不能修改现有代码"},
        Recommended: []string{"适配器模式"},
        Avoid:       []string{"修改现有接口"},
    },
    {
        Problem:     "需要动态添加功能",
        Constraints: []string{"功能组合", "运行时决定"},
        Recommended: []string{"装饰器模式", "策略模式"},
        Avoid:       []string{"继承"},
    },
    {
        Problem:     "需要解耦发送者和接收者",
        Constraints: []string{"请求需要排队", "支持撤销操作"},
        Recommended: []string{"命令模式"},
        Avoid:       []string{"直接调用"},
    },
}
```

### 1.5.2 模式组合策略

```go
// 模式组合示例
type CompositePattern struct {
    factory    *AbstractFactory
    decorators []Decorator
    observers  []Observer
}

func (cp *CompositePattern) CreateDecoratedObservableService() Service {
    // 使用工厂创建服务
    service := cp.factory.CreateService()
    
    // 使用装饰器添加功能
    for _, decorator := range cp.decorators {
        service = decorator.Decorate(service)
    }
    
    // 添加观察者
    for _, observer := range cp.observers {
        service.Attach(observer)
    }
    
    return service
}

// 模式最佳实践
type PatternBestPractices struct {
    patterns map[string][]string
}

func NewPatternBestPractices() *PatternBestPractices {
    return &PatternBestPractices{
        patterns: map[string][]string{
            "单例模式": {
                "使用sync.Once确保线程安全",
                "考虑是否需要单例",
                "避免全局状态",
            },
            "工厂模式": {
                "使用接口定义产品",
                "考虑使用依赖注入",
                "避免过度设计",
            },
            "观察者模式": {
                "注意内存泄漏",
                "考虑使用事件总线",
                "处理异步通知",
            },
        },
    }
}
```

---

**设计模式体系**: 2025年1月  
**模块状态**: ✅ **已完成**  
**质量等级**: 🏆 **企业级**
