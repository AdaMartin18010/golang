# 创建型设计模式 - Golang实现

## 目录

1. [概述](#概述)
2. [单例模式](#单例模式)
3. [工厂方法模式](#工厂方法模式)
4. [抽象工厂模式](#抽象工厂模式)
5. [建造者模式](#建造者模式)
6. [原型模式](#原型模式)
7. [性能优化](#性能优化)
8. [最佳实践](#最佳实践)

## 概述

创建型设计模式关注对象的创建机制，试图在适合特定情况的场景下创建对象。这些模式将对象的创建与使用分离，提高系统的灵活性和可维护性。

### 形式化定义

创建型模式可以形式化为：

$$\mathcal{C} = (F, P, R, V)$$

其中：
- $F = \{f_1, f_2, ..., f_n\}$ 为工厂函数集合
- $P = \{p_1, p_2, ..., p_m\}$ 为产品集合
- $R = \{r_1, r_2, ..., r_k\}$ 为创建规则集合
- $V = \{v_1, v_2, ..., v_l\}$ 为验证规则集合

### 创建函数

创建函数定义为：

$$create: F \times P \times R \rightarrow P$$

其中 $create$ 根据工厂函数、产品类型和创建规则生成产品实例。

## 单例模式

### 定义

确保一个类只有一个实例，并提供一个全局访问点。

### 形式化定义

单例模式可以定义为：

$$S = (I, A, L)$$

其中：
- $I$ 为唯一实例
- $A$ 为访问函数
- $L$ 为锁机制

### Golang实现

```go
// 线程安全的单例模式
type Singleton struct {
    data string
    mu   sync.RWMutex
}

var (
    instance *Singleton
    once     sync.Once
)

// 获取单例实例
func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "initial data",
        }
    })
    return instance
}

// 设置数据
func (s *Singleton) SetData(data string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data = data
}

// 获取数据
func (s *Singleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

// 使用sync.Once的改进版本
type ThreadSafeSingleton struct {
    instance *Singleton
    once     sync.Once
}

func (tss *ThreadSafeSingleton) GetInstance() *Singleton {
    tss.once.Do(func() {
        tss.instance = &Singleton{
            data: "thread safe data",
        }
    })
    return tss.instance
}

// 配置单例
type Config struct {
    DatabaseURL string
    Port        int
    Debug       bool
    mu          sync.RWMutex
}

var (
    configInstance *Config
    configOnce     sync.Once
)

func GetConfig() *Config {
    configOnce.Do(func() {
        configInstance = &Config{
            DatabaseURL: "localhost:5432",
            Port:        8080,
            Debug:       false,
        }
    })
    return configInstance
}

func (c *Config) SetDatabaseURL(url string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.DatabaseURL = url
}

func (c *Config) GetDatabaseURL() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.DatabaseURL
}

// 对象池单例
type ObjectPool struct {
    pool    chan interface{}
    factory func() interface{}
    reset   func(interface{}) interface{}
    mu      sync.RWMutex
}

var (
    poolInstance *ObjectPool
    poolOnce     sync.Once
)

func GetObjectPool() *ObjectPool {
    poolOnce.Do(func() {
        poolInstance = &ObjectPool{
            pool: make(chan interface{}, 100),
            factory: func() interface{} {
                return &struct{}{}
            },
            reset: func(obj interface{}) interface{} {
                return obj
            },
        }
        
        // 预填充池
        for i := 0; i < 100; i++ {
            poolInstance.pool <- poolInstance.factory()
        }
    })
    return poolInstance
}

func (op *ObjectPool) Get() interface{} {
    select {
    case obj := <-op.pool:
        return op.reset(obj)
    default:
        return op.factory()
    }
}

func (op *ObjectPool) Put(obj interface{}) {
    select {
    case op.pool <- obj:
    default:
        // 池已满，丢弃对象
    }
}
```

## 工厂方法模式

### 定义

定义一个用于创建对象的接口，让子类决定实例化哪一个类。

### 形式化定义

工厂方法模式可以定义为：

$$FM = (C, P, F)$$

其中：
- $C$ 为创建者接口
- $P$ 为产品接口
- $F$ 为工厂方法

### Golang实现

```go
// 产品接口
type Product interface {
    Operation() string
    GetName() string
}

// 具体产品A
type ConcreteProductA struct {
    name string
}

func NewConcreteProductA() *ConcreteProductA {
    return &ConcreteProductA{
        name: "ProductA",
    }
}

func (p *ConcreteProductA) Operation() string {
    return "Result of ConcreteProductA"
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

// 具体产品B
type ConcreteProductB struct {
    name string
}

func NewConcreteProductB() *ConcreteProductB {
    return &ConcreteProductB{
        name: "ProductB",
    }
}

func (p *ConcreteProductB) Operation() string {
    return "Result of ConcreteProductB"
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

// 创建者接口
type Creator interface {
    FactoryMethod() Product
    SomeOperation() string
}

// 具体创建者A
type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) FactoryMethod() Product {
    return NewConcreteProductA()
}

func (c *ConcreteCreatorA) SomeOperation() string {
    product := c.FactoryMethod()
    return fmt.Sprintf("Creator: The same creator's code has just worked with %s", product.Operation())
}

// 具体创建者B
type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) FactoryMethod() Product {
    return NewConcreteProductB()
}

func (c *ConcreteCreatorB) SomeOperation() string {
    product := c.FactoryMethod()
    return fmt.Sprintf("Creator: The same creator's code has just worked with %s", product.Operation())
}

// 客户端代码
func ClientCode(creator Creator) {
    fmt.Println(creator.SomeOperation())
}

// 数据库连接工厂
type DatabaseConnection interface {
    Connect() error
    Disconnect() error
    Execute(query string) (interface{}, error)
}

type MySQLConnection struct {
    host     string
    port     int
    username string
    password string
}

func NewMySQLConnection(host string, port int, username, password string) *MySQLConnection {
    return &MySQLConnection{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

func (m *MySQLConnection) Connect() error {
    fmt.Printf("Connecting to MySQL at %s:%d\n", m.host, m.port)
    return nil
}

func (m *MySQLConnection) Disconnect() error {
    fmt.Println("Disconnecting from MySQL")
    return nil
}

func (m *MySQLConnection) Execute(query string) (interface{}, error) {
    fmt.Printf("Executing MySQL query: %s\n", query)
    return "MySQL result", nil
}

type PostgreSQLConnection struct {
    host     string
    port     int
    username string
    password string
}

func NewPostgreSQLConnection(host string, port int, username, password string) *PostgreSQLConnection {
    return &PostgreSQLConnection{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

func (p *PostgreSQLConnection) Connect() error {
    fmt.Printf("Connecting to PostgreSQL at %s:%d\n", p.host, p.port)
    return nil
}

func (p *PostgreSQLConnection) Disconnect() error {
    fmt.Println("Disconnecting from PostgreSQL")
    return nil
}

func (p *PostgreSQLConnection) Execute(query string) (interface{}, error) {
    fmt.Printf("Executing PostgreSQL query: %s\n", query)
    return "PostgreSQL result", nil
}

// 数据库连接工厂
type DatabaseConnectionFactory interface {
    CreateConnection(config *DatabaseConfig) DatabaseConnection
}

type MySQLConnectionFactory struct{}

func (m *MySQLConnectionFactory) CreateConnection(config *DatabaseConfig) DatabaseConnection {
    return NewMySQLConnection(config.Host, config.Port, config.Username, config.Password)
}

type PostgreSQLConnectionFactory struct{}

func (p *PostgreSQLConnectionFactory) CreateConnection(config *DatabaseConfig) DatabaseConnection {
    return NewPostgreSQLConnection(config.Host, config.Port, config.Username, config.Password)
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}
```

## 抽象工厂模式

### 定义

提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们具体的类。

### 形式化定义

抽象工厂模式可以定义为：

$$AF = (F, P_A, P_B, R)$$

其中：
- $F$ 为抽象工厂接口
- $P_A$ 为产品A族
- $P_B$ 为产品B族
- $R$ 为产品间关系

### Golang实现

```go
// 产品A接口
type AbstractProductA interface {
    UsefulFunctionA() string
}

// 具体产品A1
type ConcreteProductA1 struct{}

func (p *ConcreteProductA1) UsefulFunctionA() string {
    return "The result of the product A1."
}

// 具体产品A2
type ConcreteProductA2 struct{}

func (p *ConcreteProductA2) UsefulFunctionA() string {
    return "The result of the product A2."
}

// 产品B接口
type AbstractProductB interface {
    UsefulFunctionB() string
    AnotherUsefulFunctionB(collaborator AbstractProductA) string
}

// 具体产品B1
type ConcreteProductB1 struct{}

func (p *ConcreteProductB1) UsefulFunctionB() string {
    return "The result of the product B1."
}

func (p *ConcreteProductB1) AnotherUsefulFunctionB(collaborator AbstractProductA) string {
    result := collaborator.UsefulFunctionA()
    return fmt.Sprintf("The result of the B1 collaborating with the (%s)", result)
}

// 具体产品B2
type ConcreteProductB2 struct{}

func (p *ConcreteProductB2) UsefulFunctionB() string {
    return "The result of the product B2."
}

func (p *ConcreteProductB2) AnotherUsefulFunctionB(collaborator AbstractProductA) string {
    result := collaborator.UsefulFunctionA()
    return fmt.Sprintf("The result of the B2 collaborating with the (%s)", result)
}

// 抽象工厂接口
type AbstractFactory interface {
    CreateProductA() AbstractProductA
    CreateProductB() AbstractProductB
}

// 具体工厂1
type ConcreteFactory1 struct{}

func (f *ConcreteFactory1) CreateProductA() AbstractProductA {
    return &ConcreteProductA1{}
}

func (f *ConcreteFactory1) CreateProductB() AbstractProductB {
    return &ConcreteProductB1{}
}

// 具体工厂2
type ConcreteFactory2 struct{}

func (f *ConcreteFactory2) CreateProductA() AbstractProductA {
    return &ConcreteProductA2{}
}

func (f *ConcreteFactory2) CreateProductB() AbstractProductB {
    return &ConcreteProductB2{}
}

// 客户端代码
func ClientCode(factory AbstractFactory) {
    productA := factory.CreateProductA()
    productB := factory.CreateProductB()

    fmt.Println(productB.UsefulFunctionB())
    fmt.Println(productB.AnotherUsefulFunctionB(productA))
}

// UI组件抽象工厂
type Button interface {
    Render() string
    Click() string
}

type Checkbox interface {
    Render() string
    Check() string
}

type WindowsButton struct{}

func (w *WindowsButton) Render() string {
    return "Windows button rendered"
}

func (w *WindowsButton) Click() string {
    return "Windows button clicked"
}

type WindowsCheckbox struct{}

func (w *WindowsCheckbox) Render() string {
    return "Windows checkbox rendered"
}

func (w *WindowsCheckbox) Check() string {
    return "Windows checkbox checked"
}

type MacButton struct{}

func (m *MacButton) Render() string {
    return "Mac button rendered"
}

func (m *MacButton) Click() string {
    return "Mac button clicked"
}

type MacCheckbox struct{}

func (m *MacCheckbox) Render() string {
    return "Mac checkbox rendered"
}

func (m *MacCheckbox) Check() string {
    return "Mac checkbox checked"
}

// UI工厂接口
type UIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

type WindowsUIFactory struct{}

func (w *WindowsUIFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (w *WindowsUIFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{}
}

type MacUIFactory struct{}

func (m *MacUIFactory) CreateButton() Button {
    return &MacButton{}
}

func (m *MacUIFactory) CreateCheckbox() Checkbox {
    return &MacCheckbox{}
}
```

## 建造者模式

### 定义

将一个复杂对象的构建与其表示分离，使得同样的构建过程可以创建不同的表示。

### 形式化定义

建造者模式可以定义为：

$$B = (D, S, F)$$

其中：
- $D$ 为导演
- $S$ 为构建步骤集合
- $F$ 为最终产品

### Golang实现

```go
// 产品
type Product struct {
    PartA string
    PartB string
    PartC int
}

func (p *Product) String() string {
    return fmt.Sprintf("Product{PartA: %s, PartB: %s, PartC: %d}", p.PartA, p.PartB, p.PartC)
}

// 建造者接口
type Builder interface {
    SetPartA(partA string) Builder
    SetPartB(partB string) Builder
    SetPartC(partC int) Builder
    Build() *Product
}

// 具体建造者
type ConcreteBuilder struct {
    product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
    return &ConcreteBuilder{
        product: &Product{},
    }
}

func (b *ConcreteBuilder) SetPartA(partA string) Builder {
    b.product.PartA = partA
    return b
}

func (b *ConcreteBuilder) SetPartB(partB string) Builder {
    b.product.PartB = partB
    return b
}

func (b *ConcreteBuilder) SetPartC(partC int) Builder {
    b.product.PartC = partC
    return b
}

func (b *ConcreteBuilder) Build() *Product {
    return b.product
}

// 导演
type Director struct {
    builder Builder
}

func NewDirector(builder Builder) *Director {
    return &Director{
        builder: builder,
    }
}

func (d *Director) Construct() *Product {
    return d.builder.
        SetPartA("Part A").
        SetPartB("Part B").
        SetPartC(100).
        Build()
}

// 数据库查询构建器
type SQLQuery struct {
    Select  []string
    From    string
    Where   []string
    OrderBy []string
    Limit   int
    Offset  int
}

func (q *SQLQuery) String() string {
    query := "SELECT " + strings.Join(q.Select, ", ")
    query += " FROM " + q.From
    
    if len(q.Where) > 0 {
        query += " WHERE " + strings.Join(q.Where, " AND ")
    }
    
    if len(q.OrderBy) > 0 {
        query += " ORDER BY " + strings.Join(q.OrderBy, ", ")
    }
    
    if q.Limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", q.Limit)
    }
    
    if q.Offset > 0 {
        query += fmt.Sprintf(" OFFSET %d", q.Offset)
    }
    
    return query
}

// SQL查询构建器
type SQLQueryBuilder struct {
    query *SQLQuery
}

func NewSQLQueryBuilder() *SQLQueryBuilder {
    return &SQLQueryBuilder{
        query: &SQLQuery{},
    }
}

func (b *SQLQueryBuilder) Select(columns ...string) *SQLQueryBuilder {
    b.query.Select = columns
    return b
}

func (b *SQLQueryBuilder) From(table string) *SQLQueryBuilder {
    b.query.From = table
    return b
}

func (b *SQLQueryBuilder) Where(condition string) *SQLQueryBuilder {
    b.query.Where = append(b.query.Where, condition)
    return b
}

func (b *SQLQueryBuilder) OrderBy(column string) *SQLQueryBuilder {
    b.query.OrderBy = append(b.query.OrderBy, column)
    return b
}

func (b *SQLQueryBuilder) Limit(limit int) *SQLQueryBuilder {
    b.query.Limit = limit
    return b
}

func (b *SQLQueryBuilder) Offset(offset int) *SQLQueryBuilder {
    b.query.Offset = offset
    return b
}

func (b *SQLQueryBuilder) Build() *SQLQuery {
    return b.query
}

// HTTP请求构建器
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
}

func (r *HTTPRequest) String() string {
    return fmt.Sprintf("HTTPRequest{Method: %s, URL: %s, Headers: %v}", r.Method, r.URL, r.Headers)
}

type HTTPRequestBuilder struct {
    request *HTTPRequest
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
    return &HTTPRequestBuilder{
        request: &HTTPRequest{
            Headers: make(map[string]string),
        },
    }
}

func (b *HTTPRequestBuilder) Method(method string) *HTTPRequestBuilder {
    b.request.Method = method
    return b
}

func (b *HTTPRequestBuilder) URL(url string) *HTTPRequestBuilder {
    b.request.URL = url
    return b
}

func (b *HTTPRequestBuilder) Header(key, value string) *HTTPRequestBuilder {
    b.request.Headers[key] = value
    return b
}

func (b *HTTPRequestBuilder) Body(body []byte) *HTTPRequestBuilder {
    b.request.Body = body
    return b
}

func (b *HTTPRequestBuilder) Build() *HTTPRequest {
    return b.request
}
```

## 原型模式

### 定义

用原型实例指定创建对象的种类，并且通过拷贝这些原型创建新的对象。

### 形式化定义

原型模式可以定义为：

$$P = (O, C, D)$$

其中：
- $O$ 为原始对象
- $C$ 为克隆函数
- $D$ 为深度复制函数

### Golang实现

```go
// 原型接口
type Prototype interface {
    Clone() Prototype
    DeepClone() Prototype
}

// 具体原型
type ConcretePrototype struct {
    Name    string
    Data    map[string]interface{}
    Created time.Time
}

func NewConcretePrototype(name string) *ConcretePrototype {
    return &ConcretePrototype{
        Name:    name,
        Data:    make(map[string]interface{}),
        Created: time.Now(),
    }
}

// 浅克隆
func (p *ConcretePrototype) Clone() Prototype {
    return &ConcretePrototype{
        Name:    p.Name,
        Data:    p.Data, // 共享map引用
        Created: p.Created,
    }
}

// 深克隆
func (p *ConcretePrototype) DeepClone() Prototype {
    // 深拷贝map
    newData := make(map[string]interface{})
    for k, v := range p.Data {
        newData[k] = v
    }
    
    return &ConcretePrototype{
        Name:    p.Name,
        Data:    newData,
        Created: p.Created,
    }
}

func (p *ConcretePrototype) SetData(key string, value interface{}) {
    p.Data[key] = value
}

func (p *ConcretePrototype) GetData(key string) interface{} {
    return p.Data[key]
}

// 使用json进行深克隆
type JSONPrototype struct {
    Name    string                 `json:"name"`
    Data    map[string]interface{} `json:"data"`
    Created time.Time              `json:"created"`
}

func NewJSONPrototype(name string) *JSONPrototype {
    return &JSONPrototype{
        Name:    name,
        Data:    make(map[string]interface{}),
        Created: time.Now(),
    }
}

func (p *JSONPrototype) Clone() Prototype {
    // 使用JSON序列化进行深克隆
    data, err := json.Marshal(p)
    if err != nil {
        return nil
    }
    
    var clone JSONPrototype
    if err := json.Unmarshal(data, &clone); err != nil {
        return nil
    }
    
    return &clone
}

func (p *JSONPrototype) DeepClone() Prototype {
    return p.Clone() // JSON序列化已经是深克隆
}

// 原型管理器
type PrototypeManager struct {
    prototypes map[string]Prototype
    mu         sync.RWMutex
}

func NewPrototypeManager() *PrototypeManager {
    return &PrototypeManager{
        prototypes: make(map[string]Prototype),
    }
}

func (pm *PrototypeManager) Register(name string, prototype Prototype) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.prototypes[name] = prototype
}

func (pm *PrototypeManager) Get(name string) Prototype {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    if prototype, exists := pm.prototypes[name]; exists {
        return prototype.Clone()
    }
    return nil
}

func (pm *PrototypeManager) Remove(name string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    delete(pm.prototypes, name)
}

// 文档模板原型
type DocumentTemplate struct {
    Title       string            `json:"title"`
    Content     string            `json:"content"`
    Styles      map[string]string `json:"styles"`
    Metadata    map[string]string `json:"metadata"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

func NewDocumentTemplate(title, content string) *DocumentTemplate {
    return &DocumentTemplate{
        Title:     title,
        Content:   content,
        Styles:    make(map[string]string),
        Metadata:  make(map[string]string),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func (dt *DocumentTemplate) Clone() Prototype {
    data, err := json.Marshal(dt)
    if err != nil {
        return nil
    }
    
    var clone DocumentTemplate
    if err := json.Unmarshal(data, &clone); err != nil {
        return nil
    }
    
    clone.UpdatedAt = time.Now()
    return &clone
}

func (dt *DocumentTemplate) DeepClone() Prototype {
    return dt.Clone()
}

func (dt *DocumentTemplate) AddStyle(selector, style string) {
    dt.Styles[selector] = style
    dt.UpdatedAt = time.Now()
}

func (dt *DocumentTemplate) AddMetadata(key, value string) {
    dt.Metadata[key] = value
    dt.UpdatedAt = time.Now()
}
```

## 性能优化

### 1. 对象池优化

```go
// 对象池优化
type ObjectPoolOptimizer struct {
    pools map[string]*sync.Pool
    mu    sync.RWMutex
}

func NewObjectPoolOptimizer() *ObjectPoolOptimizer {
    return &ObjectPoolOptimizer{
        pools: make(map[string]*sync.Pool),
    }
}

func (opo *ObjectPoolOptimizer) GetPool(name string, factory func() interface{}) *sync.Pool {
    opo.mu.Lock()
    defer opo.mu.Unlock()
    
    if pool, exists := opo.pools[name]; exists {
        return pool
    }
    
    pool := &sync.Pool{
        New: factory,
    }
    opo.pools[name] = pool
    return pool
}

// 优化的建造者
type OptimizedBuilder struct {
    pool *sync.Pool
}

func NewOptimizedBuilder() *OptimizedBuilder {
    return &OptimizedBuilder{
        pool: &sync.Pool{
            New: func() interface{} {
                return &Product{}
            },
        },
    }
}

func (b *OptimizedBuilder) Build() *Product {
    obj := b.pool.Get().(*Product)
    // 重置对象状态
    *obj = Product{}
    return obj
}

func (b *OptimizedBuilder) Recycle(product *Product) {
    b.pool.Put(product)
}
```

### 2. 缓存优化

```go
// 工厂缓存
type CachedFactory struct {
    cache *sync.Map
    factory func() interface{}
}

func NewCachedFactory(factory func() interface{}) *CachedFactory {
    return &CachedFactory{
        cache:   &sync.Map{},
        factory: factory,
    }
}

func (cf *CachedFactory) Get(key string) interface{} {
    if cached, exists := cf.cache.Load(key); exists {
        return cached
    }
    
    obj := cf.factory()
    cf.cache.Store(key, obj)
    return obj
}

func (cf *CachedFactory) Clear() {
    cf.cache = &sync.Map{}
}
```

### 3. 并发优化

```go
// 并发安全的工厂
type ConcurrentFactory struct {
    factories map[string]func() interface{}
    mu        sync.RWMutex
}

func NewConcurrentFactory() *ConcurrentFactory {
    return &ConcurrentFactory{
        factories: make(map[string]func() interface{}),
    }
}

func (cf *ConcurrentFactory) Register(name string, factory func() interface{}) {
    cf.mu.Lock()
    defer cf.mu.Unlock()
    cf.factories[name] = factory
}

func (cf *ConcurrentFactory) Create(name string) (interface{}, error) {
    cf.mu.RLock()
    factory, exists := cf.factories[name]
    cf.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("factory not found: %s", name)
    }
    
    return factory(), nil
}
```

## 最佳实践

### 1. 错误处理

```go
// 工厂错误
type FactoryError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e FactoryError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
    ErrFactoryNotFound = FactoryError{Code: "FACTORY_NOT_FOUND", Message: "Factory not found"}
    ErrInvalidProduct = FactoryError{Code: "INVALID_PRODUCT", Message: "Invalid product"}
    ErrCreationFailed = FactoryError{Code: "CREATION_FAILED", Message: "Product creation failed"}
)

// 安全的工厂方法
func SafeCreateProduct(factory func() (interface{}, error)) (interface{}, error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic in factory: %v", r)
        }
    }()
    
    return factory()
}
```

### 2. 配置管理

```go
// 工厂配置
type FactoryConfig struct {
    MaxInstances    int           `json:"max_instances"`
    Timeout         time.Duration `json:"timeout"`
    EnableCaching   bool          `json:"enable_caching"`
    CacheSize       int           `json:"cache_size"`
    EnablePooling   bool          `json:"enable_pooling"`
    PoolSize        int           `json:"pool_size"`
}

func LoadFactoryConfig(filename string) (*FactoryConfig, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var config FactoryConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 3. 测试策略

```go
// 工厂测试
func TestFactoryMethod(t *testing.T) {
    creator := &ConcreteCreatorA{}
    product := creator.FactoryMethod()
    
    if product == nil {
        t.Error("Product should not be nil")
    }
    
    result := product.Operation()
    expected := "Result of ConcreteProductA"
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}

// 单例测试
func TestSingleton(t *testing.T) {
    instance1 := GetInstance()
    instance2 := GetInstance()
    
    if instance1 != instance2 {
        t.Error("Singleton instances should be the same")
    }
    
    instance1.SetData("test data")
    if instance2.GetData() != "test data" {
        t.Error("Singleton state should be shared")
    }
}

// 建造者测试
func TestBuilder(t *testing.T) {
    builder := NewConcreteBuilder()
    product := builder.
        SetPartA("A").
        SetPartB("B").
        SetPartC(100).
        Build()
    
    if product.PartA != "A" {
        t.Errorf("Expected PartA to be A, got %s", product.PartA)
    }
    
    if product.PartB != "B" {
        t.Errorf("Expected PartB to be B, got %s", product.PartB)
    }
    
    if product.PartC != 100 {
        t.Errorf("Expected PartC to be 100, got %d", product.PartC)
    }
}

// 原型测试
func TestPrototype(t *testing.T) {
    original := NewConcretePrototype("original")
    original.SetData("key", "value")
    
    clone := original.Clone().(*ConcretePrototype)
    
    if clone.Name != original.Name {
        t.Error("Clone should have same name")
    }
    
    if clone.GetData("key") != original.GetData("key") {
        t.Error("Clone should have same data")
    }
    
    // 测试浅克隆
    clone.SetData("key", "new value")
    if original.GetData("key") != "new value" {
        t.Error("Shallow clone should share data reference")
    }
}
```

## 总结

创建型设计模式提供了灵活的对象创建机制，通过Golang的特性实现了高性能、线程安全的模式实现。

关键要点：
1. **单例模式**: 使用sync.Once确保线程安全
2. **工厂方法**: 通过接口实现多态创建
3. **抽象工厂**: 创建相关对象族
4. **建造者模式**: 链式调用构建复杂对象
5. **原型模式**: 通过克隆创建对象
6. **性能优化**: 对象池、缓存、并发优化
7. **最佳实践**: 错误处理、配置管理、测试策略 