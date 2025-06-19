# 创建型模式详细分析

## 目录

1. [概述](#概述)
2. [工厂模式](#工厂模式)
3. [抽象工厂模式](#抽象工厂模式)
4. [单例模式](#单例模式)
5. [建造者模式](#建造者模式)
6. [原型模式](#原型模式)
7. [最佳实践](#最佳实践)

## 概述

创建型模式专注于对象的创建机制，通过将对象的创建与使用分离，提高系统的灵活性和可维护性。本文档深入分析五种主要的创建型模式在Golang中的实现方案。

### 核心特征

- **封装创建逻辑**: 隐藏对象创建的复杂性
- **提高灵活性**: 支持不同的创建策略
- **降低耦合度**: 客户端与具体产品解耦
- **支持扩展**: 易于添加新的产品类型
- **控制实例**: 管理对象的生命周期

## 工厂模式

### 形式化定义

**定义 18.1** (工厂模式)
工厂模式是一个五元组 $\mathcal{FP} = (F, P, C, I, R)$，其中：

- $F$ 是工厂 (Factory)
- $P$ 是产品集合 (Products)
- $C$ 是创建条件 (Creation Conditions)
- $I$ 是接口 (Interface)
- $R$ 是创建规则 (Creation Rules)

### 简单工厂模式

```go
// 产品接口
type Product interface {
    Use() string
    GetName() string
}

// 具体产品A
type ConcreteProductA struct {
    name string
}

func (cpa *ConcreteProductA) Use() string {
    return "Using ConcreteProductA"
}

func (cpa *ConcreteProductA) GetName() string {
    return cpa.name
}

// 具体产品B
type ConcreteProductB struct {
    name string
}

func (cpb *ConcreteProductB) Use() string {
    return "Using ConcreteProductB"
}

func (cpb *ConcreteProductB) GetName() string {
    return cpb.name
}

// 简单工厂
type SimpleFactory struct{}

func (sf *SimpleFactory) CreateProduct(productType string) (Product, error) {
    switch productType {
    case "A":
        return &ConcreteProductA{name: "ProductA"}, nil
    case "B":
        return &ConcreteProductB{name: "ProductB"}, nil
    default:
        return nil, fmt.Errorf("unknown product type: %s", productType)
    }
}

// 工厂方法模式
type FactoryMethod interface {
    CreateProduct() Product
}

// 具体工厂A
type ConcreteFactoryA struct{}

func (cfa *ConcreteFactoryA) CreateProduct() Product {
    return &ConcreteProductA{name: "ProductA"}
}

// 具体工厂B
type ConcreteFactoryB struct{}

func (cfb *ConcreteFactoryB) CreateProduct() Product {
    return &ConcreteProductB{name: "ProductB"}
}

// 工厂管理器
type FactoryManager struct {
    factories map[string]FactoryMethod
    mu        sync.RWMutex
}

func NewFactoryManager() *FactoryManager {
    return &FactoryManager{
        factories: make(map[string]FactoryMethod),
    }
}

func (fm *FactoryManager) RegisterFactory(name string, factory FactoryMethod) {
    fm.mu.Lock()
    defer fm.mu.Unlock()
    fm.factories[name] = factory
}

func (fm *FactoryManager) CreateProduct(factoryName string) (Product, error) {
    fm.mu.RLock()
    defer fm.mu.RUnlock()
    
    factory, exists := fm.factories[factoryName]
    if !exists {
        return nil, fmt.Errorf("factory %s not found", factoryName)
    }
    
    return factory.CreateProduct(), nil
}
```

### 高级工厂模式

```go
// 产品族接口
type ProductFamily interface {
    CreateChair() Chair
    CreateTable() Table
    CreateSofa() Sofa
}

// 家具接口
type Chair interface {
    Sit() string
}

type Table interface {
    Put() string
}

type Sofa interface {
    Lie() string
}

// 现代风格产品族
type ModernChair struct{}

func (mc *ModernChair) Sit() string {
    return "Sitting on modern chair"
}

type ModernTable struct{}

func (mt *ModernTable) Put() string {
    return "Putting on modern table"
}

type ModernSofa struct{}

func (ms *ModernSofa) Lie() string {
    return "Lying on modern sofa"
}

type ModernProductFamily struct{}

func (mpf *ModernProductFamily) CreateChair() Chair {
    return &ModernChair{}
}

func (mpf *ModernProductFamily) CreateTable() Table {
    return &ModernTable{}
}

func (mpf *ModernProductFamily) CreateSofa() Sofa {
    return &ModernSofa{}
}

// 古典风格产品族
type ClassicChair struct{}

func (cc *ClassicChair) Sit() string {
    return "Sitting on classic chair"
}

type ClassicTable struct{}

func (ct *ClassicTable) Put() string {
    return "Putting on classic table"
}

type ClassicSofa struct{}

func (cs *ClassicSofa) Lie() string {
    return "Lying on classic sofa"
}

type ClassicProductFamily struct{}

func (cpf *ClassicProductFamily) CreateChair() Chair {
    return &ClassicChair{}
}

func (cpf *ClassicProductFamily) CreateTable() Table {
    return &ClassicTable{}
}

func (cpf *ClassicProductFamily) CreateSofa() Sofa {
    return &ClassicSofa{}
}

// 产品族工厂
type ProductFamilyFactory struct {
    families map[string]ProductFamily
    mu       sync.RWMutex
}

func NewProductFamilyFactory() *ProductFamilyFactory {
    factory := &ProductFamilyFactory{
        families: make(map[string]ProductFamily),
    }
    
    // 注册默认产品族
    factory.RegisterFamily("modern", &ModernProductFamily{})
    factory.RegisterFamily("classic", &ClassicProductFamily{})
    
    return factory
}

func (pff *ProductFamilyFactory) RegisterFamily(name string, family ProductFamily) {
    pff.mu.Lock()
    defer pff.mu.Unlock()
    pff.families[name] = family
}

func (pff *ProductFamilyFactory) GetFamily(name string) (ProductFamily, error) {
    pff.mu.RLock()
    defer pff.mu.RUnlock()
    
    family, exists := pff.families[name]
    if !exists {
        return nil, fmt.Errorf("product family %s not found", name)
    }
    
    return family, nil
}
```

## 抽象工厂模式

### 形式化定义

**定义 18.2** (抽象工厂模式)
抽象工厂模式是一个六元组 $\mathcal{AFP} = (AF, PF, CF, I, R, C)$，其中：

- $AF$ 是抽象工厂 (Abstract Factory)
- $PF$ 是产品族 (Product Families)
- $CF$ 是具体工厂 (Concrete Factories)
- $I$ 是接口集合 (Interface Set)
- $R$ 是创建规则 (Creation Rules)
- $C$ 是约束条件 (Constraints)

### 实现示例

```go
// 抽象工厂接口
type AbstractFactory interface {
    CreateDatabase() Database
    CreateCache() Cache
    CreateQueue() Queue
}

// 数据存储接口
type Database interface {
    Connect() error
    Query(sql string) ([]map[string]interface{}, error)
    Close() error
}

type Cache interface {
    Set(key string, value interface{}) error
    Get(key string) (interface{}, error)
    Delete(key string) error
}

type Queue interface {
    Push(message interface{}) error
    Pop() (interface{}, error)
    Size() int
}

// MySQL产品族
type MySQLDatabase struct {
    connectionString string
}

func (md *MySQLDatabase) Connect() error {
    log.Printf("Connecting to MySQL: %s", md.connectionString)
    return nil
}

func (md *MySQLDatabase) Query(sql string) ([]map[string]interface{}, error) {
    log.Printf("Executing MySQL query: %s", sql)
    return []map[string]interface{}{{"result": "mysql_data"}}, nil
}

func (md *MySQLDatabase) Close() error {
    log.Printf("Closing MySQL connection")
    return nil
}

type RedisCache struct {
    address string
}

func (rc *RedisCache) Set(key string, value interface{}) error {
    log.Printf("Setting Redis key %s: %v", key, value)
    return nil
}

func (rc *RedisCache) Get(key string) (interface{}, error) {
    log.Printf("Getting Redis key %s", key)
    return "redis_value", nil
}

func (rc *RedisCache) Delete(key string) error {
    log.Printf("Deleting Redis key %s", key)
    return nil
}

type RabbitMQQueue struct {
    url string
}

func (rmq *RabbitMQQueue) Push(message interface{}) error {
    log.Printf("Pushing message to RabbitMQ: %v", message)
    return nil
}

func (rmq *RabbitMQQueue) Pop() (interface{}, error) {
    log.Printf("Popping message from RabbitMQ")
    return "rabbitmq_message", nil
}

func (rmq *RabbitMQQueue) Size() int {
    return 10
}

type MySQLFactory struct {
    config map[string]string
}

func NewMySQLFactory(config map[string]string) *MySQLFactory {
    return &MySQLFactory{config: config}
}

func (mf *MySQLFactory) CreateDatabase() Database {
    return &MySQLDatabase{
        connectionString: mf.config["mysql_connection"],
    }
}

func (mf *MySQLFactory) CreateCache() Cache {
    return &RedisCache{
        address: mf.config["redis_address"],
    }
}

func (mf *MySQLFactory) CreateQueue() Queue {
    return &RabbitMQQueue{
        url: mf.config["rabbitmq_url"],
    }
}

// PostgreSQL产品族
type PostgreSQLDatabase struct {
    connectionString string
}

func (pd *PostgreSQLDatabase) Connect() error {
    log.Printf("Connecting to PostgreSQL: %s", pd.connectionString)
    return nil
}

func (pd *PostgreSQLDatabase) Query(sql string) ([]map[string]interface{}, error) {
    log.Printf("Executing PostgreSQL query: %s", sql)
    return []map[string]interface{}{{"result": "postgresql_data"}}, nil
}

func (pd *PostgreSQLDatabase) Close() error {
    log.Printf("Closing PostgreSQL connection")
    return nil
}

type MemcachedCache struct {
    servers []string
}

func (mc *MemcachedCache) Set(key string, value interface{}) error {
    log.Printf("Setting Memcached key %s: %v", key, value)
    return nil
}

func (mc *MemcachedCache) Get(key string) (interface{}, error) {
    log.Printf("Getting Memcached key %s", key)
    return "memcached_value", nil
}

func (mc *MemcachedCache) Delete(key string) error {
    log.Printf("Deleting Memcached key %s", key)
    return nil
}

type KafkaQueue struct {
    brokers []string
}

func (kq *KafkaQueue) Push(message interface{}) error {
    log.Printf("Pushing message to Kafka: %v", message)
    return nil
}

func (kq *KafkaQueue) Pop() (interface{}, error) {
    log.Printf("Popping message from Kafka")
    return "kafka_message", nil
}

func (kq *KafkaQueue) Size() int {
    return 20
}

type PostgreSQLFactory struct {
    config map[string]string
}

func NewPostgreSQLFactory(config map[string]string) *PostgreSQLFactory {
    return &PostgreSQLFactory{config: config}
}

func (pf *PostgreSQLFactory) CreateDatabase() Database {
    return &PostgreSQLDatabase{
        connectionString: pf.config["postgresql_connection"],
    }
}

func (pf *PostgreSQLFactory) CreateCache() Cache {
    return &MemcachedCache{
        servers: strings.Split(pf.config["memcached_servers"], ","),
    }
}

func (pf *PostgreSQLFactory) CreateQueue() Queue {
    return &KafkaQueue{
        brokers: strings.Split(pf.config["kafka_brokers"], ","),
    }
}

// 抽象工厂管理器
type AbstractFactoryManager struct {
    factories map[string]AbstractFactory
    mu        sync.RWMutex
}

func NewAbstractFactoryManager() *AbstractFactoryManager {
    return &AbstractFactoryManager{
        factories: make(map[string]AbstractFactory),
    }
}

func (afm *AbstractFactoryManager) RegisterFactory(name string, factory AbstractFactory) {
    afm.mu.Lock()
    defer afm.mu.Unlock()
    afm.factories[name] = factory
}

func (afm *AbstractFactoryManager) GetFactory(name string) (AbstractFactory, error) {
    afm.mu.RLock()
    defer afm.mu.RUnlock()
    
    factory, exists := afm.factories[name]
    if !exists {
        return nil, fmt.Errorf("factory %s not found", name)
    }
    
    return factory, nil
}

func (afm *AbstractFactoryManager) CreateProductFamily(factoryName string) (*ProductFamily, error) {
    factory, err := afm.GetFactory(factoryName)
    if err != nil {
        return nil, err
    }
    
    return &ProductFamily{
        Database: factory.CreateDatabase(),
        Cache:    factory.CreateCache(),
        Queue:    factory.CreateQueue(),
    }, nil
}

// 产品族
type ProductFamily struct {
    Database Database
    Cache    Cache
    Queue    Queue
}
```

## 单例模式

### 形式化定义

**定义 18.3** (单例模式)
单例模式是一个四元组 $\mathcal{SP} = (I, C, L, T)$，其中：

- $I$ 是实例 (Instance)
- $C$ 是创建条件 (Creation Condition)
- $L$ 是生命周期 (Lifecycle)
- $T$ 是线程安全 (Thread Safety)

### 实现示例

```go
// 基础单例
type Singleton struct {
    data string
    mu   sync.RWMutex
}

func (s *Singleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

func (s *Singleton) SetData(data string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data = data
}

// 饿汉式单例
type EagerSingleton struct {
    data string
}

var eagerInstance = &EagerSingleton{
    data: "Eager Singleton Data",
}

func GetEagerSingleton() *EagerSingleton {
    return eagerInstance
}

// 懒汉式单例
type LazySingleton struct {
    data string
}

var (
    lazyInstance *LazySingleton
    lazyOnce     sync.Once
)

func GetLazySingleton() *LazySingleton {
    lazyOnce.Do(func() {
        lazyInstance = &LazySingleton{
            data: "Lazy Singleton Data",
        }
    })
    return lazyInstance
}

// 双重检查锁定单例
type DoubleCheckSingleton struct {
    data string
}

var (
    doubleCheckInstance *DoubleCheckSingleton
    doubleCheckMu       sync.Mutex
)

func GetDoubleCheckSingleton() *DoubleCheckSingleton {
    if doubleCheckInstance == nil {
        doubleCheckMu.Lock()
        defer doubleCheckMu.Unlock()
        
        if doubleCheckInstance == nil {
            doubleCheckInstance = &DoubleCheckSingleton{
                data: "Double Check Singleton Data",
            }
        }
    }
    return doubleCheckInstance
}

// 单例管理器
type SingletonManager struct {
    singletons map[string]interface{}
    mu         sync.RWMutex
}

func NewSingletonManager() *SingletonManager {
    return &SingletonManager{
        singletons: make(map[string]interface{}),
    }
}

func (sm *SingletonManager) GetSingleton(name string, creator func() interface{}) interface{} {
    sm.mu.RLock()
    if instance, exists := sm.singletons[name]; exists {
        sm.mu.RUnlock()
        return instance
    }
    sm.mu.RUnlock()
    
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // 双重检查
    if instance, exists := sm.singletons[name]; exists {
        return instance
    }
    
    instance := creator()
    sm.singletons[name] = instance
    return instance
}

// 配置单例
type Config struct {
    DatabaseURL string
    RedisURL    string
    Port        int
    mu          sync.RWMutex
}

func (c *Config) GetDatabaseURL() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.DatabaseURL
}

func (c *Config) SetDatabaseURL(url string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.DatabaseURL = url
}

func (c *Config) GetRedisURL() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.RedisURL
}

func (c *Config) SetRedisURL(url string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.RedisURL = url
}

func (c *Config) GetPort() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.Port
}

func (c *Config) SetPort(port int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Port = port
}

var (
    configInstance *Config
    configOnce     sync.Once
)

func GetConfig() *Config {
    configOnce.Do(func() {
        configInstance = &Config{
            DatabaseURL: "localhost:5432",
            RedisURL:    "localhost:6379",
            Port:        8080,
        }
    })
    return configInstance
}
```

## 建造者模式

### 形式化定义

**定义 18.4** (建造者模式)
建造者模式是一个五元组 $\mathcal{BP} = (B, D, P, S, R)$，其中：

- $B$ 是建造者 (Builder)
- $D$ 是导演 (Director)
- $P$ 是产品 (Product)
- $S$ 是步骤 (Steps)
- $R$ 是结果 (Result)

### 实现示例

```go
// 产品
type Computer struct {
    CPU       string
    Memory    string
    Storage   string
    Graphics  string
    Network   string
    Power     string
    Case      string
}

func (c *Computer) String() string {
    return fmt.Sprintf("Computer{CPU: %s, Memory: %s, Storage: %s, Graphics: %s, Network: %s, Power: %s, Case: %s}",
        c.CPU, c.Memory, c.Storage, c.Graphics, c.Network, c.Power, c.Case)
}

// 抽象建造者
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetStorage(storage string) ComputerBuilder
    SetGraphics(graphics string) ComputerBuilder
    SetNetwork(network string) ComputerBuilder
    SetPower(power string) ComputerBuilder
    SetCase(caseType string) ComputerBuilder
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

func (gcb *GamingComputerBuilder) SetGraphics(graphics string) ComputerBuilder {
    gcb.computer.Graphics = graphics
    return gcb
}

func (gcb *GamingComputerBuilder) SetNetwork(network string) ComputerBuilder {
    gcb.computer.Network = network
    return gcb
}

func (gcb *GamingComputerBuilder) SetPower(power string) ComputerBuilder {
    gcb.computer.Power = power
    return gcb
}

func (gcb *GamingComputerBuilder) SetCase(caseType string) ComputerBuilder {
    gcb.computer.Case = caseType
    return gcb
}

func (gcb *GamingComputerBuilder) Build() *Computer {
    return gcb.computer
}

// 办公电脑建造者
type OfficeComputerBuilder struct {
    computer *Computer
}

func NewOfficeComputerBuilder() *OfficeComputerBuilder {
    return &OfficeComputerBuilder{
        computer: &Computer{},
    }
}

func (ocb *OfficeComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    ocb.computer.CPU = cpu
    return ocb
}

func (ocb *OfficeComputerBuilder) SetMemory(memory string) ComputerBuilder {
    ocb.computer.Memory = memory
    return ocb
}

func (ocb *OfficeComputerBuilder) SetStorage(storage string) ComputerBuilder {
    ocb.computer.Storage = storage
    return ocb
}

func (ocb *OfficeComputerBuilder) SetGraphics(graphics string) ComputerBuilder {
    ocb.computer.Graphics = graphics
    return ocb
}

func (ocb *OfficeComputerBuilder) SetNetwork(network string) ComputerBuilder {
    ocb.computer.Network = network
    return ocb
}

func (ocb *OfficeComputerBuilder) SetPower(power string) ComputerBuilder {
    ocb.computer.Power = power
    return ocb
}

func (ocb *OfficeComputerBuilder) SetCase(caseType string) ComputerBuilder {
    ocb.computer.Case = caseType
    return ocb
}

func (ocb *OfficeComputerBuilder) Build() *Computer {
    return ocb.computer
}

// 导演
type ComputerDirector struct {
    builder ComputerBuilder
}

func NewComputerDirector(builder ComputerBuilder) *ComputerDirector {
    return &ComputerDirector{builder: builder}
}

func (cd *ComputerDirector) SetBuilder(builder ComputerBuilder) {
    cd.builder = builder
}

func (cd *ComputerDirector) BuildGamingComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i9-12900K").
        SetMemory("32GB DDR5").
        SetStorage("2TB NVMe SSD").
        SetGraphics("RTX 4090").
        SetNetwork("WiFi 6E").
        SetPower("850W Gold").
        SetCase("ATX Gaming Case").
        Build()
}

func (cd *ComputerDirector) BuildOfficeComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i5-12400").
        SetMemory("16GB DDR4").
        SetStorage("512GB SSD").
        SetGraphics("Integrated").
        SetNetwork("Ethernet").
        SetPower("450W Bronze").
        SetCase("Micro ATX Case").
        Build()
}

// 建造者管理器
type BuilderManager struct {
    builders map[string]ComputerBuilder
    director *ComputerDirector
    mu       sync.RWMutex
}

func NewBuilderManager() *BuilderManager {
    return &BuilderManager{
        builders: make(map[string]ComputerBuilder),
        director: NewComputerDirector(nil),
    }
}

func (bm *BuilderManager) RegisterBuilder(name string, builder ComputerBuilder) {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    bm.builders[name] = builder
}

func (bm *BuilderManager) GetBuilder(name string) (ComputerBuilder, error) {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    
    builder, exists := bm.builders[name]
    if !exists {
        return nil, fmt.Errorf("builder %s not found", name)
    }
    
    return builder, nil
}

func (bm *BuilderManager) BuildComputer(builderName, computerType string) (*Computer, error) {
    builder, err := bm.GetBuilder(builderName)
    if err != nil {
        return nil, err
    }
    
    bm.director.SetBuilder(builder)
    
    switch computerType {
    case "gaming":
        return bm.director.BuildGamingComputer(), nil
    case "office":
        return bm.director.BuildOfficeComputer(), nil
    default:
        return nil, fmt.Errorf("unknown computer type: %s", computerType)
    }
}
```

## 原型模式

### 形式化定义

**定义 18.5** (原型模式)
原型模式是一个四元组 $\mathcal{PP} = (P, C, R, M)$，其中：

- $P$ 是原型 (Prototype)
- $C$ 是克隆方法 (Clone Method)
- $R$ 是注册表 (Registry)
- $M$ 是管理器 (Manager)

### 实现示例

```go
// 原型接口
type Prototype interface {
    Clone() Prototype
    GetName() string
}

// 具体原型
type Document struct {
    title    string
    content  string
    author   string
    created  time.Time
    modified time.Time
}

func (d *Document) Clone() Prototype {
    return &Document{
        title:    d.title,
        content:  d.content,
        author:   d.author,
        created:  d.created,
        modified: time.Now(),
    }
}

func (d *Document) GetName() string {
    return d.title
}

func (d *Document) SetTitle(title string) {
    d.title = title
    d.modified = time.Now()
}

func (d *Document) SetContent(content string) {
    d.content = content
    d.modified = time.Now()
}

func (d *Document) String() string {
    return fmt.Sprintf("Document{Title: %s, Author: %s, Created: %s, Modified: %s}",
        d.title, d.author, d.created.Format("2006-01-02 15:04:05"), d.modified.Format("2006-01-02 15:04:05"))
}

// 深度克隆
type DeepCloneable interface {
    DeepClone() DeepCloneable
}

type ComplexObject struct {
    ID       string
    Data     map[string]interface{}
    Children []*ComplexObject
    mu       sync.RWMutex
}

func (co *ComplexObject) DeepClone() DeepCloneable {
    co.mu.RLock()
    defer co.mu.RUnlock()
    
    cloned := &ComplexObject{
        ID:   co.ID,
        Data: make(map[string]interface{}),
    }
    
    // 深度复制数据
    for key, value := range co.Data {
        if mapValue, ok := value.(map[string]interface{}); ok {
            clonedMap := make(map[string]interface{})
            for k, v := range mapValue {
                clonedMap[k] = v
            }
            cloned.Data[key] = clonedMap
        } else {
            cloned.Data[key] = value
        }
    }
    
    // 深度复制子对象
    for _, child := range co.Children {
        clonedChild := child.DeepClone().(*ComplexObject)
        cloned.Children = append(cloned.Children, clonedChild)
    }
    
    return cloned
}

// 原型注册表
type PrototypeRegistry struct {
    prototypes map[string]Prototype
    mu         sync.RWMutex
}

func NewPrototypeRegistry() *PrototypeRegistry {
    return &PrototypeRegistry{
        prototypes: make(map[string]Prototype),
    }
}

func (pr *PrototypeRegistry) Register(name string, prototype Prototype) {
    pr.mu.Lock()
    defer pr.mu.Unlock()
    pr.prototypes[name] = prototype
}

func (pr *PrototypeRegistry) Get(name string) (Prototype, error) {
    pr.mu.RLock()
    defer pr.mu.RUnlock()
    
    prototype, exists := pr.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    
    return prototype.Clone(), nil
}

func (pr *PrototypeRegistry) List() []string {
    pr.mu.RLock()
    defer pr.mu.RUnlock()
    
    names := make([]string, 0, len(pr.prototypes))
    for name := range pr.prototypes {
        names = append(names, name)
    }
    
    return names
}

// 原型管理器
type PrototypeManager struct {
    registry *PrototypeRegistry
}

func NewPrototypeManager() *PrototypeManager {
    return &PrototypeManager{
        registry: NewPrototypeRegistry(),
    }
}

func (pm *PrototypeManager) InitializePrototypes() {
    // 注册默认原型
    defaultDoc := &Document{
        title:    "Default Document",
        content:  "This is a default document template.",
        author:   "System",
        created:  time.Now(),
        modified: time.Now(),
    }
    pm.registry.Register("default", defaultDoc)
    
    // 注册报告原型
    reportDoc := &Document{
        title:    "Report Template",
        content:  "This is a report template with predefined structure.",
        author:   "System",
        created:  time.Now(),
        modified: time.Now(),
    }
    pm.registry.Register("report", reportDoc)
    
    // 注册信件原型
    letterDoc := &Document{
        title:    "Letter Template",
        content:  "This is a letter template with formal structure.",
        author:   "System",
        created:  time.Now(),
        modified: time.Now(),
    }
    pm.registry.Register("letter", letterDoc)
}

func (pm *PrototypeManager) CreateDocument(templateName string) (*Document, error) {
    prototype, err := pm.registry.Get(templateName)
    if err != nil {
        return nil, err
    }
    
    if doc, ok := prototype.(*Document); ok {
        return doc, nil
    }
    
    return nil, fmt.Errorf("prototype is not a document")
}

func (pm *PrototypeManager) CreateCustomDocument(templateName, title, content, author string) (*Document, error) {
    doc, err := pm.CreateDocument(templateName)
    if err != nil {
        return nil, err
    }
    
    doc.SetTitle(title)
    doc.SetContent(content)
    doc.author = author
    
    return doc, nil
}
```

## 最佳实践

### 1. 错误处理

```go
// 创建型模式错误类型
type CreationalPatternError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Pattern string `json:"pattern,omitempty"`
    Details string `json:"details,omitempty"`
}

func (e *CreationalPatternError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// 错误码常量
const (
    ErrCodeFactoryNotFound = "FACTORY_NOT_FOUND"
    ErrCodeProductNotFound = "PRODUCT_NOT_FOUND"
    ErrCodeSingletonError  = "SINGLETON_ERROR"
    ErrCodeBuilderError    = "BUILDER_ERROR"
    ErrCodePrototypeError  = "PROTOTYPE_ERROR"
)

// 统一错误处理
func HandleCreationalPatternError(err error, pattern string) *CreationalPatternError {
    switch {
    case errors.Is(err, ErrFactoryNotFound):
        return &CreationalPatternError{
            Code:    ErrCodeFactoryNotFound,
            Message: "Factory not found",
            Pattern: pattern,
        }
    case errors.Is(err, ErrSingletonError):
        return &CreationalPatternError{
            Code:    ErrCodeSingletonError,
            Message: "Singleton error",
            Pattern: pattern,
        }
    default:
        return &CreationalPatternError{
            Code: ErrCodeProductNotFound,
            Message: "Product not found",
        }
    }
}
```

### 2. 监控和日志

```go
// 创建型模式指标
type CreationalPatternMetrics struct {
    factoryUsage    prometheus.CounterVec
    singletonUsage  prometheus.CounterVec
    builderUsage    prometheus.CounterVec
    prototypeUsage  prometheus.CounterVec
    patternErrors   prometheus.CounterVec
}

func NewCreationalPatternMetrics() *CreationalPatternMetrics {
    return &CreationalPatternMetrics{
        factoryUsage: *prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "creational_factory_usage_total",
            Help: "Total number of factory pattern usage",
        }, []string{"factory_type", "product_type"}),
        singletonUsage: *prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "creational_singleton_usage_total",
            Help: "Total number of singleton pattern usage",
        }, []string{"singleton_type"}),
        builderUsage: *prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "creational_builder_usage_total",
            Help: "Total number of builder pattern usage",
        }, []string{"builder_type", "product_type"}),
        prototypeUsage: *prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "creational_prototype_usage_total",
            Help: "Total number of prototype pattern usage",
        }, []string{"prototype_type"}),
        patternErrors: *prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "creational_pattern_errors_total",
            Help: "Total number of creational pattern errors",
        }, []string{"pattern_type", "error_type"}),
    }
}

// 创建型模式日志
type CreationalPatternLogger struct {
    logger *zap.Logger
}

func (l *CreationalPatternLogger) LogFactoryUsage(factoryType, productType string) {
    l.logger.Info("factory usage",
        zap.String("factory_type", factoryType),
        zap.String("product_type", productType),
    )
}

func (l *CreationalPatternLogger) LogSingletonUsage(singletonType string) {
    l.logger.Info("singleton usage",
        zap.String("singleton_type", singletonType),
    )
}

func (l *CreationalPatternLogger) LogBuilderUsage(builderType, productType string) {
    l.logger.Info("builder usage",
        zap.String("builder_type", builderType),
        zap.String("product_type", productType),
    )
}
```

### 3. 测试策略

```go
// 单元测试
func TestFactoryPattern_CreateProduct(t *testing.T) {
    factory := &SimpleFactory{}
    
    // 测试创建产品A
    productA, err := factory.CreateProduct("A")
    if err != nil {
        t.Errorf("Failed to create product A: %v", err)
    }
    
    if productA.GetName() != "ProductA" {
        t.Errorf("Expected ProductA, got %s", productA.GetName())
    }
    
    // 测试创建产品B
    productB, err := factory.CreateProduct("B")
    if err != nil {
        t.Errorf("Failed to create product B: %v", err)
    }
    
    if productB.GetName() != "ProductB" {
        t.Errorf("Expected ProductB, got %s", productB.GetName())
    }
    
    // 测试未知产品类型
    _, err = factory.CreateProduct("C")
    if err == nil {
        t.Error("Expected error for unknown product type")
    }
}

// 集成测试
func TestAbstractFactory_CreateProductFamily(t *testing.T) {
    // 创建工厂管理器
    manager := NewAbstractFactoryManager()
    
    // 注册MySQL工厂
    mysqlConfig := map[string]string{
        "mysql_connection": "localhost:5432",
        "redis_address":    "localhost:6379",
        "rabbitmq_url":     "localhost:5672",
    }
    manager.RegisterFactory("mysql", NewMySQLFactory(mysqlConfig))
    
    // 创建产品族
    family, err := manager.CreateProductFamily("mysql")
    if err != nil {
        t.Errorf("Failed to create product family: %v", err)
    }
    
    if family.Database == nil {
        t.Error("Expected database to be created")
    }
    
    if family.Cache == nil {
        t.Error("Expected cache to be created")
    }
    
    if family.Queue == nil {
        t.Error("Expected queue to be created")
    }
}

// 性能测试
func BenchmarkSingleton_GetInstance(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = GetLazySingleton()
    }
}

func BenchmarkBuilder_BuildComputer(b *testing.B) {
    builder := NewGamingComputerBuilder()
    director := NewComputerDirector(builder)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = director.BuildGamingComputer()
    }
}
```

---

## 总结

本文档深入分析了五种主要的创建型模式：

1. **工厂模式**: 封装对象创建逻辑，支持不同的创建策略
2. **抽象工厂模式**: 创建相关产品族，确保产品兼容性
3. **单例模式**: 确保类只有一个实例，提供全局访问点
4. **建造者模式**: 分步骤构建复杂对象，支持链式调用
5. **原型模式**: 通过克隆创建对象，避免重复初始化

每种模式都包含形式化定义、Golang实现、高级用法和最佳实践，为实际项目开发提供了完整的参考方案。

---

**最后更新**: 2024-12-19  
**当前状态**: ✅ 创建型模式分析完成  
**下一步**: 结构型模式详细分析
