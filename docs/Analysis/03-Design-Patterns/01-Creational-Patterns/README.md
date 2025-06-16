# 创建型设计模式分析

## 目录

1. [概念定义](#1-概念定义)
2. [单例模式 (Singleton Pattern)](#2-单例模式-singleton-pattern)
3. [工厂方法模式 (Factory Method Pattern)](#3-工厂方法模式-factory-method-pattern)
4. [抽象工厂模式 (Abstract Factory Pattern)](#4-抽象工厂模式-abstract-factory-pattern)
5. [建造者模式 (Builder Pattern)](#5-建造者模式-builder-pattern)
6. [原型模式 (Prototype Pattern)](#6-原型模式-prototype-pattern)
7. [模式比较](#7-模式比较)
8. [最佳实践](#8-最佳实践)

## 1. 概念定义

### 定义 1.1 (创建型模式)
创建型模式处理对象创建机制，试图在适合特定情况的场景下创建对象：
$$\mathcal{C}_{Creational} = \{Singleton, FactoryMethod, AbstractFactory, Builder, Prototype\}$$

### 定义 1.2 (对象创建)
对象创建函数：
$$Create: Class \times Parameters \rightarrow Object$$

### 定义 1.3 (创建模式分类)
创建型模式按创建方式分类：
1. **单例模式**: 确保一个类只有一个实例
2. **工厂模式**: 通过工厂方法创建对象
3. **建造者模式**: 分步骤构建复杂对象
4. **原型模式**: 通过克隆创建对象

## 2. 单例模式 (Singleton Pattern)

### 2.1 形式化定义

#### 定义 2.1 (单例模式)
单例模式确保一个类只有一个实例，并提供全局访问点：
$$Singleton(C) = \{instance \in C | \forall c \in C : c = instance\}$$

#### 定理 2.1 (单例唯一性)
对于任意类 $C$，单例模式保证：
$$\exists! instance \in C : \forall c \in C : c = instance$$

**证明**:
1. 假设存在两个不同的实例 $instance_1$ 和 $instance_2$
2. 根据单例模式定义，$instance_1 = instance_2$
3. 这与假设矛盾
4. 因此，单例实例是唯一的

### 2.2 Golang实现

#### 2.2.1 线程安全单例
```go
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "initialized"}
    })
    return instance
}

func (s *Singleton) GetData() string {
    return s.data
}

func (s *Singleton) SetData(data string) {
    s.data = data
}
```

#### 2.2.2 延迟初始化单例
```go
type LazySingleton struct {
    data string
}

var (
    lazyInstance *LazySingleton
    lazyMutex    sync.Mutex
)

func GetLazyInstance() *LazySingleton {
    if lazyInstance == nil {
        lazyMutex.Lock()
        defer lazyMutex.Unlock()
        
        if lazyInstance == nil {
            lazyInstance = &LazySingleton{data: "lazy initialized"}
        }
    }
    return lazyInstance
}
```

#### 2.2.3 配置单例
```go
type Config struct {
    DatabaseURL string
    Port        int
    Debug       bool
    mutex       sync.RWMutex
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

func (c *Config) GetDatabaseURL() string {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.DatabaseURL
}

func (c *Config) SetDatabaseURL(url string) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.DatabaseURL = url
}
```

### 2.3 性能分析

#### 定义 2.2 (单例性能)
单例模式的性能特征：
$$Performance(Singleton) = O(1) \text{ for access}$$

#### 内存使用分析
```go
type MemoryAnalyzer struct {
    instances map[string]int
    mutex     sync.RWMutex
}

func (ma *MemoryAnalyzer) AnalyzeMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Heap Alloc: %d bytes\n", m.HeapAlloc)
    fmt.Printf("Heap Sys: %d bytes\n", m.HeapSys)
    fmt.Printf("Num Goroutines: %d\n", runtime.NumGoroutine())
}
```

## 3. 工厂方法模式 (Factory Method Pattern)

### 3.1 形式化定义

#### 定义 3.1 (工厂方法)
工厂方法定义一个创建对象的接口，让子类决定实例化哪个类：
$$FactoryMethod(C, F) = \{f \in F | f: \emptyset \rightarrow C\}$$

#### 定义 3.2 (产品族)
产品族定义为：
$$\mathcal{P}_{Family} = \{Product_1, Product_2, ..., Product_n\}$$

#### 定理 3.1 (工厂方法可扩展性)
工厂方法模式支持开闭原则：
$$\forall Product_{new} \in \mathcal{P}_{Family} : \exists Factory_{new} \in \mathcal{F}$$

### 3.2 Golang实现

#### 3.2.1 基础工厂方法
```go
type Product interface {
    Operation() string
}

type Creator interface {
    FactoryMethod() Product
    SomeOperation() string
}

type ConcreteCreator struct{}

func (c *ConcreteCreator) FactoryMethod() Product {
    return &ConcreteProduct{}
}

func (c *ConcreteCreator) SomeOperation() string {
    product := c.FactoryMethod()
    return "Creator: " + product.Operation()
}

type ConcreteProduct struct{}

func (cp *ConcreteProduct) Operation() string {
    return "ConcreteProduct operation"
}
```

#### 3.2.2 数据库连接工厂
```go
type Database interface {
    Connect() error
    Query(sql string) ([]map[string]interface{}, error)
    Close() error
}

type DatabaseCreator interface {
    CreateDatabase(config DatabaseConfig) Database
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type MySQLCreator struct{}

func (mc *MySQLCreator) CreateDatabase(config DatabaseConfig) Database {
    return &MySQLDatabase{
        config: config,
    }
}

type PostgreSQLCreator struct{}

func (pc *PostgreSQLCreator) CreateDatabase(config DatabaseConfig) Database {
    return &PostgreSQLDatabase{
        config: config,
    }
}

type MySQLDatabase struct {
    config DatabaseConfig
    conn   *sql.DB
}

func (md *MySQLDatabase) Connect() error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
        md.config.Username, md.config.Password, 
        md.config.Host, md.config.Port, md.config.Database)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    
    md.conn = db
    return nil
}

func (md *MySQLDatabase) Query(sql string) ([]map[string]interface{}, error) {
    rows, err := md.conn.Query(sql)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    return md.scanRows(rows)
}

func (md *MySQLDatabase) Close() error {
    return md.conn.Close()
}

func (md *MySQLDatabase) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var result []map[string]interface{}
    
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        
        result = append(result, row)
    }
    
    return result, nil
}
```

#### 3.2.3 日志工厂
```go
type Logger interface {
    Log(level LogLevel, message string) error
    SetLevel(level LogLevel)
}

type LogLevel int

const (
    Debug LogLevel = iota
    Info
    Warn
    Error
)

type LoggerCreator interface {
    CreateLogger(config LoggerConfig) Logger
}

type LoggerConfig struct {
    Level      LogLevel
    OutputPath string
    Format     string
}

type FileLoggerCreator struct{}

func (flc *FileLoggerCreator) CreateLogger(config LoggerConfig) Logger {
    return &FileLogger{
        config: config,
        file:   nil,
    }
}

type ConsoleLoggerCreator struct{}

func (clc *ConsoleLoggerCreator) CreateLogger(config LoggerConfig) Logger {
    return &ConsoleLogger{
        config: config,
    }
}

type FileLogger struct {
    config LoggerConfig
    file   *os.File
    mutex  sync.Mutex
}

func (fl *FileLogger) Log(level LogLevel, message string) error {
    if level < fl.config.Level {
        return nil
    }
    
    fl.mutex.Lock()
    defer fl.mutex.Unlock()
    
    if fl.file == nil {
        file, err := os.OpenFile(fl.config.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return err
        }
        fl.file = file
    }
    
    logEntry := fmt.Sprintf("[%s] %s: %s\n", 
        time.Now().Format("2006-01-02 15:04:05"),
        level.String(),
        message)
    
    _, err := fl.file.WriteString(logEntry)
    return err
}

func (fl *FileLogger) SetLevel(level LogLevel) {
    fl.config.Level = level
}

func (l LogLevel) String() string {
    switch l {
    case Debug:
        return "DEBUG"
    case Info:
        return "INFO"
    case Warn:
        return "WARN"
    case Error:
        return "ERROR"
    default:
        return "UNKNOWN"
    }
}
```

## 4. 抽象工厂模式 (Abstract Factory Pattern)

### 4.1 形式化定义

#### 定义 4.1 (抽象工厂)
抽象工厂提供一个创建一系列相关或相互依赖对象的接口：
$$AbstractFactory(F_1, F_2, ..., F_n) = \prod_{i=1}^{n} F_i$$

#### 定义 4.2 (产品系列)
产品系列定义为：
$$\mathcal{S}_{Product} = \{Series_1, Series_2, ..., Series_n\}$$

#### 定理 4.1 (抽象工厂一致性)
抽象工厂确保产品系列的一致性：
$$\forall Series_i \in \mathcal{S}_{Product} : Compatible(Series_i)$$

### 4.2 Golang实现

#### 4.2.1 基础抽象工厂
```go
type AbstractFactory interface {
    CreateProductA() ProductA
    CreateProductB() ProductB
}

type ProductA interface {
    OperationA() string
}

type ProductB interface {
    OperationB() string
}

type ConcreteFactory1 struct{}

func (cf1 *ConcreteFactory1) CreateProductA() ProductA {
    return &ConcreteProductA1{}
}

func (cf1 *ConcreteFactory1) CreateProductB() ProductB {
    return &ConcreteProductB1{}
}

type ConcreteFactory2 struct{}

func (cf2 *ConcreteFactory2) CreateProductA() ProductA {
    return &ConcreteProductA2{}
}

func (cf2 *ConcreteFactory2) CreateProductB() ProductB {
    return &ConcreteProductB2{}
}

type ConcreteProductA1 struct{}

func (cpa1 *ConcreteProductA1) OperationA() string {
    return "ConcreteProductA1 operation"
}

type ConcreteProductA2 struct{}

func (cpa2 *ConcreteProductA2) OperationA() string {
    return "ConcreteProductA2 operation"
}

type ConcreteProductB1 struct{}

func (cpb1 *ConcreteProductB1) OperationB() string {
    return "ConcreteProductB1 operation"
}

type ConcreteProductB2 struct{}

func (cpb2 *ConcreteProductB2) OperationB() string {
    return "ConcreteProductB2 operation"
}
```

#### 4.2.2 UI组件工厂
```go
type UIComponent interface {
    Render() string
    SetStyle(style Style)
}

type Button interface {
    UIComponent
    Click() string
}

type TextField interface {
    UIComponent
    SetText(text string)
    GetText() string
}

type UIFactory interface {
    CreateButton() Button
    CreateTextField() TextField
}

type MaterialUIFactory struct{}

func (muf *MaterialUIFactory) CreateButton() Button {
    return &MaterialButton{
        text:  "Button",
        style: MaterialStyle{},
    }
}

func (muf *MaterialUIFactory) CreateTextField() TextField {
    return &MaterialTextField{
        text:  "",
        style: MaterialStyle{},
    }
}

type BootstrapUIFactory struct{}

func (buf *BootstrapUIFactory) CreateButton() Button {
    return &BootstrapButton{
        text:  "Button",
        style: BootstrapStyle{},
    }
}

func (buf *BootstrapUIFactory) CreateTextField() TextField {
    return &BootstrapTextField{
        text:  "",
        style: BootstrapStyle{},
    }
}

type Style interface {
    GetCSS() string
}

type MaterialStyle struct{}

func (ms MaterialStyle) GetCSS() string {
    return "material-design-style"
}

type BootstrapStyle struct{}

func (bs BootstrapStyle) GetCSS() string {
    return "bootstrap-style"
}

type MaterialButton struct {
    text  string
    style Style
}

func (mb *MaterialButton) Render() string {
    return fmt.Sprintf("<button class='%s'>%s</button>", mb.style.GetCSS(), mb.text)
}

func (mb *MaterialButton) SetStyle(style Style) {
    mb.style = style
}

func (mb *MaterialButton) Click() string {
    return "Material button clicked"
}

type MaterialTextField struct {
    text  string
    style Style
}

func (mtf *MaterialTextField) Render() string {
    return fmt.Sprintf("<input class='%s' value='%s' />", mtf.style.GetCSS(), mtf.text)
}

func (mtf *MaterialTextField) SetStyle(style Style) {
    mtf.style = style
}

func (mtf *MaterialTextField) SetText(text string) {
    mtf.text = text
}

func (mtf *MaterialTextField) GetText() string {
    return mtf.text
}
```

#### 4.2.3 数据库抽象工厂
```go
type DatabaseFactory interface {
    CreateConnection() Connection
    CreateQueryBuilder() QueryBuilder
    CreateTransaction() Transaction
}

type Connection interface {
    Connect() error
    Disconnect() error
    IsConnected() bool
}

type QueryBuilder interface {
    Select(columns ...string) QueryBuilder
    From(table string) QueryBuilder
    Where(condition string) QueryBuilder
    Build() string
}

type Transaction interface {
    Begin() error
    Commit() error
    Rollback() error
}

type MySQLFactory struct{}

func (mf *MySQLFactory) CreateConnection() Connection {
    return &MySQLConnection{}
}

func (mf *MySQLFactory) CreateQueryBuilder() QueryBuilder {
    return &MySQLQueryBuilder{}
}

func (mf *MySQLFactory) CreateTransaction() Transaction {
    return &MySQLTransaction{}
}

type PostgreSQLFactory struct{}

func (pf *PostgreSQLFactory) CreateConnection() Connection {
    return &PostgreSQLConnection{}
}

func (pf *PostgreSQLFactory) CreateQueryBuilder() QueryBuilder {
    return &PostgreSQLQueryBuilder{}
}

func (pf *PostgreSQLFactory) CreateTransaction() Transaction {
    return &PostgreSQLTransaction{}
}

type MySQLConnection struct {
    connected bool
}

func (mc *MySQLConnection) Connect() error {
    mc.connected = true
    return nil
}

func (mc *MySQLConnection) Disconnect() error {
    mc.connected = false
    return nil
}

func (mc *MySQLConnection) IsConnected() bool {
    return mc.connected
}

type MySQLQueryBuilder struct {
    columns  []string
    table    string
    where    string
}

func (mqb *MySQLQueryBuilder) Select(columns ...string) QueryBuilder {
    mqb.columns = columns
    return mqb
}

func (mqb *MySQLQueryBuilder) From(table string) QueryBuilder {
    mqb.table = table
    return mqb
}

func (mqb *MySQLQueryBuilder) Where(condition string) QueryBuilder {
    mqb.where = condition
    return mqb
}

func (mqb *MySQLQueryBuilder) Build() string {
    query := "SELECT "
    if len(mqb.columns) == 0 {
        query += "*"
    } else {
        query += strings.Join(mqb.columns, ", ")
    }
    
    query += " FROM " + mqb.table
    
    if mqb.where != "" {
        query += " WHERE " + mqb.where
    }
    
    return query
}
```

## 5. 建造者模式 (Builder Pattern)

### 5.1 形式化定义

#### 定义 5.1 (建造者模式)
建造者模式分步骤构建复杂对象：
$$Builder(Product) = \{Step_1, Step_2, ..., Step_n\} \rightarrow Product$$

#### 定义 5.2 (构建步骤)
构建步骤序列：
$$BuildSteps = [Step_1, Step_2, ..., Step_n]$$

#### 定理 5.1 (建造者完整性)
建造者模式确保对象构建的完整性：
$$\forall Product : Complete(Product) \iff \forall Step_i \in BuildSteps : Step_i(Product)$$

### 5.2 Golang实现

#### 5.2.1 基础建造者
```go
type Product struct {
    PartA string
    PartB string
    PartC string
}

type Builder interface {
    BuildPartA() Builder
    BuildPartB() Builder
    BuildPartC() Builder
    GetResult() *Product
}

type ConcreteBuilder struct {
    product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
    return &ConcreteBuilder{
        product: &Product{},
    }
}

func (cb *ConcreteBuilder) BuildPartA() Builder {
    cb.product.PartA = "Part A"
    return cb
}

func (cb *ConcreteBuilder) BuildPartB() Builder {
    cb.product.PartB = "Part B"
    return cb
}

func (cb *ConcreteBuilder) BuildPartC() Builder {
    cb.product.PartC = "Part C"
    return cb
}

func (cb *ConcreteBuilder) GetResult() *Product {
    return cb.product
}

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
        BuildPartA().
        BuildPartB().
        BuildPartC().
        GetResult()
}
```

#### 5.2.2 HTTP请求建造者
```go
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
    Timeout time.Duration
}

type HTTPRequestBuilder struct {
    request *HTTPRequest
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
    return &HTTPRequestBuilder{
        request: &HTTPRequest{
            Headers: make(map[string]string),
            Timeout: 30 * time.Second,
        },
    }
}

func (hrb *HTTPRequestBuilder) SetMethod(method string) *HTTPRequestBuilder {
    hrb.request.Method = method
    return hrb
}

func (hrb *HTTPRequestBuilder) SetURL(url string) *HTTPRequestBuilder {
    hrb.request.URL = url
    return hrb
}

func (hrb *HTTPRequestBuilder) AddHeader(key, value string) *HTTPRequestBuilder {
    hrb.request.Headers[key] = value
    return hrb
}

func (hrb *HTTPRequestBuilder) SetBody(body []byte) *HTTPRequestBuilder {
    hrb.request.Body = body
    return hrb
}

func (hrb *HTTPRequestBuilder) SetTimeout(timeout time.Duration) *HTTPRequestBuilder {
    hrb.request.Timeout = timeout
    return hrb
}

func (hrb *HTTPRequestBuilder) Build() *HTTPRequest {
    return hrb.request
}

func (hrb *HTTPRequestBuilder) Execute() (*http.Response, error) {
    req, err := http.NewRequest(hrb.request.Method, hrb.request.URL, bytes.NewReader(hrb.request.Body))
    if err != nil {
        return nil, err
    }
    
    for key, value := range hrb.request.Headers {
        req.Header.Set(key, value)
    }
    
    client := &http.Client{
        Timeout: hrb.request.Timeout,
    }
    
    return client.Do(req)
}

// 使用示例
func ExampleHTTPRequestBuilder() {
    response, err := NewHTTPRequestBuilder().
        SetMethod("POST").
        SetURL("https://api.example.com/data").
        AddHeader("Content-Type", "application/json").
        AddHeader("Authorization", "Bearer token").
        SetBody([]byte(`{"key": "value"}`)).
        SetTimeout(10 * time.Second).
        Execute()
    
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    
    fmt.Printf("Status: %s\n", response.Status)
}
```

#### 5.2.3 配置建造者
```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Logging  LoggingConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type ServerConfig struct {
    Host string
    Port int
    SSL  bool
}

type LoggingConfig struct {
    Level      string
    OutputPath string
    Format     string
}

type ConfigBuilder struct {
    config *Config
}

func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{
        config: &Config{},
    }
}

func (cb *ConfigBuilder) SetDatabase(host string, port int, username, password, database string) *ConfigBuilder {
    cb.config.Database = DatabaseConfig{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        Database: database,
    }
    return cb
}

func (cb *ConfigBuilder) SetServer(host string, port int, ssl bool) *ConfigBuilder {
    cb.config.Server = ServerConfig{
        Host: host,
        Port: port,
        SSL:  ssl,
    }
    return cb
}

func (cb *ConfigBuilder) SetLogging(level, outputPath, format string) *ConfigBuilder {
    cb.config.Logging = LoggingConfig{
        Level:      level,
        OutputPath: outputPath,
        Format:     format,
    }
    return cb
}

func (cb *ConfigBuilder) Build() *Config {
    return cb.config
}

func (cb *ConfigBuilder) BuildFromFile(filepath string) (*Config, error) {
    data, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    cb.config = &config
    return cb.config, nil
}

func (cb *ConfigBuilder) Validate() error {
    if cb.config.Database.Host == "" {
        return errors.New("database host is required")
    }
    
    if cb.config.Server.Port <= 0 {
        return errors.New("server port must be positive")
    }
    
    if cb.config.Logging.Level == "" {
        return errors.New("logging level is required")
    }
    
    return nil
}
```

## 6. 原型模式 (Prototype Pattern)

### 6.1 形式化定义

#### 定义 6.1 (原型模式)
原型模式通过克隆现有对象来创建新对象：
$$Prototype(Object) = Clone(Object) \rightarrow NewObject$$

#### 定义 6.2 (克隆深度)
克隆深度定义为：
$$CloneDepth = \{Shallow, Deep\}$$

#### 定理 6.1 (原型等价性)
原型克隆的对象与原对象在结构上等价：
$$Clone(Object) \equiv Object$$

### 6.2 Golang实现

#### 6.2.1 基础原型
```go
type Prototype interface {
    Clone() Prototype
}

type ConcretePrototype struct {
    Field1 string
    Field2 int
    Field3 []string
}

func (cp *ConcretePrototype) Clone() Prototype {
    // 深拷贝
    cloned := &ConcretePrototype{
        Field1: cp.Field1,
        Field2: cp.Field2,
        Field3: make([]string, len(cp.Field3)),
    }
    
    copy(cloned.Field3, cp.Field3)
    return cloned
}

func (cp *ConcretePrototype) ShallowClone() Prototype {
    // 浅拷贝
    return &ConcretePrototype{
        Field1: cp.Field1,
        Field2: cp.Field2,
        Field3: cp.Field3, // 共享切片
    }
}
```

#### 6.2.2 文档原型
```go
type Document interface {
    Clone() Document
    GetContent() string
    SetContent(content string)
}

type TextDocument struct {
    content string
    style   DocumentStyle
}

type DocumentStyle struct {
    FontSize int
    FontName string
    Color    string
}

func (td *TextDocument) Clone() Document {
    return &TextDocument{
        content: td.content,
        style:   td.style, // 结构体是值类型，自动深拷贝
    }
}

func (td *TextDocument) GetContent() string {
    return td.content
}

func (td *TextDocument) SetContent(content string) {
    td.content = content
}

type DocumentManager struct {
    templates map[string]Document
}

func NewDocumentManager() *DocumentManager {
    return &DocumentManager{
        templates: make(map[string]Document),
    }
}

func (dm *DocumentManager) RegisterTemplate(name string, template Document) {
    dm.templates[name] = template
}

func (dm *DocumentManager) CreateDocument(templateName string) (Document, error) {
    template, exists := dm.templates[templateName]
    if !exists {
        return nil, fmt.Errorf("template %s not found", templateName)
    }
    
    return template.Clone(), nil
}

// 使用示例
func ExampleDocumentPrototype() {
    manager := NewDocumentManager()
    
    // 注册模板
    template := &TextDocument{
        content: "This is a template",
        style: DocumentStyle{
            FontSize: 12,
            FontName: "Arial",
            Color:    "black",
        },
    }
    
    manager.RegisterTemplate("default", template)
    
    // 创建新文档
    doc1, _ := manager.CreateDocument("default")
    doc1.SetContent("Document 1 content")
    
    doc2, _ := manager.CreateDocument("default")
    doc2.SetContent("Document 2 content")
    
    fmt.Printf("Doc1: %s\n", doc1.GetContent())
    fmt.Printf("Doc2: %s\n", doc2.GetContent())
}
```

#### 6.2.3 游戏对象原型
```go
type GameObject interface {
    Clone() GameObject
    GetPosition() Position
    SetPosition(position Position)
    Update()
}

type Position struct {
    X, Y, Z float64
}

type Enemy struct {
    position Position
    health   int
    speed    float64
    behavior string
}

func (e *Enemy) Clone() GameObject {
    return &Enemy{
        position: e.position,
        health:   e.health,
        speed:    e.speed,
        behavior: e.behavior,
    }
}

func (e *Enemy) GetPosition() Position {
    return e.position
}

func (e *Enemy) SetPosition(position Position) {
    e.position = position
}

func (e *Enemy) Update() {
    // 更新敌人行为
    switch e.behavior {
    case "patrol":
        e.patrol()
    case "chase":
        e.chase()
    case "attack":
        e.attack()
    }
}

func (e *Enemy) patrol() {
    // 巡逻逻辑
}

func (e *Enemy) chase() {
    // 追击逻辑
}

func (e *Enemy) attack() {
    // 攻击逻辑
}

type GameObjectFactory struct {
    prototypes map[string]GameObject
}

func NewGameObjectFactory() *GameObjectFactory {
    return &GameObjectFactory{
        prototypes: make(map[string]GameObject),
    }
}

func (gof *GameObjectFactory) RegisterPrototype(name string, prototype GameObject) {
    gof.prototypes[name] = prototype
}

func (gof *GameObjectFactory) CreateObject(name string, position Position) (GameObject, error) {
    prototype, exists := gof.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    
    obj := prototype.Clone()
    obj.SetPosition(position)
    return obj, nil
}

// 使用示例
func ExampleGameObjectPrototype() {
    factory := NewGameObjectFactory()
    
    // 注册敌人原型
    enemyPrototype := &Enemy{
        position: Position{0, 0, 0},
        health:   100,
        speed:    5.0,
        behavior: "patrol",
    }
    
    factory.RegisterPrototype("enemy", enemyPrototype)
    
    // 创建多个敌人实例
    for i := 0; i < 5; i++ {
        enemy, _ := factory.CreateObject("enemy", Position{
            X: float64(i * 10),
            Y: 0,
            Z: 0,
        })
        
        // 每个敌人可以有不同的行为
        if e, ok := enemy.(*Enemy); ok {
            if i%2 == 0 {
                e.behavior = "chase"
            } else {
                e.behavior = "attack"
            }
        }
    }
}
```

## 7. 模式比较

### 7.1 性能比较

#### 定义 7.1 (模式性能)
各创建型模式的性能特征：
- **单例模式**: $O(1)$ 访问时间
- **工厂方法**: $O(1)$ 创建时间
- **抽象工厂**: $O(n)$ 创建时间，其中 $n$ 是产品数量
- **建造者模式**: $O(m)$ 构建时间，其中 $m$ 是步骤数量
- **原型模式**: $O(k)$ 克隆时间，其中 $k$ 是对象复杂度

### 7.2 内存使用比较

```go
type PatternMemoryAnalyzer struct {
    patterns map[string]func() interface{}
}

func NewPatternMemoryAnalyzer() *PatternMemoryAnalyzer {
    return &PatternMemoryAnalyzer{
        patterns: make(map[string]func() interface{}),
    }
}

func (pma *PatternMemoryAnalyzer) RegisterPattern(name string, creator func() interface{}) {
    pma.patterns[name] = creator
}

func (pma *PatternMemoryAnalyzer) AnalyzeMemoryUsage(patternName string, iterations int) {
    creator, exists := pma.patterns[patternName]
    if !exists {
        return
    }
    
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    var objects []interface{}
    for i := 0; i < iterations; i++ {
        obj := creator()
        objects = append(objects, obj)
    }
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    memoryUsed := m2.HeapAlloc - m1.HeapAlloc
    fmt.Printf("Pattern: %s, Objects: %d, Memory: %d bytes\n", 
        patternName, iterations, memoryUsed)
}
```

### 7.3 适用场景比较

| 模式 | 适用场景 | 优点 | 缺点 |
|------|----------|------|------|
| 单例 | 全局状态管理 | 内存效率高 | 测试困难 |
| 工厂方法 | 对象创建复杂 | 扩展性好 | 类数量增加 |
| 抽象工厂 | 产品族创建 | 一致性保证 | 复杂度高 |
| 建造者 | 复杂对象构建 | 步骤清晰 | 代码冗长 |
| 原型 | 对象克隆 | 性能好 | 深拷贝复杂 |

## 8. 最佳实践

### 8.1 选择指南

#### 8.1.1 何时使用单例模式
- 需要全局唯一实例
- 资源管理（配置、连接池）
- 日志记录器

#### 8.1.2 何时使用工厂方法
- 对象创建逻辑复杂
- 需要支持扩展
- 子类决定创建对象

#### 8.1.3 何时使用抽象工厂
- 需要创建产品族
- 产品间有依赖关系
- 需要保证一致性

#### 8.1.4 何时使用建造者模式
- 对象构建步骤多
- 构建过程复杂
- 需要不同构建方式

#### 8.1.5 何时使用原型模式
- 对象创建成本高
- 需要对象副本
- 避免重复初始化

### 8.2 实现建议

#### 8.2.1 线程安全
```go
type ThreadSafeSingleton struct {
    data string
    mutex sync.RWMutex
}

func (tss *ThreadSafeSingleton) GetData() string {
    tss.mutex.RLock()
    defer tss.mutex.RUnlock()
    return tss.data
}

func (tss *ThreadSafeSingleton) SetData(data string) {
    tss.mutex.Lock()
    defer tss.mutex.Unlock()
    tss.data = data
}
```

#### 8.2.2 错误处理
```go
type SafeBuilder struct {
    product *Product
    errors  []error
}

func (sb *SafeBuilder) BuildPartA() *SafeBuilder {
    if sb.product == nil {
        sb.product = &Product{}
    }
    
    // 构建逻辑
    if err := sb.validatePartA(); err != nil {
        sb.errors = append(sb.errors, err)
    }
    
    return sb
}

func (sb *SafeBuilder) Build() (*Product, error) {
    if len(sb.errors) > 0 {
        return nil, fmt.Errorf("build failed: %v", sb.errors)
    }
    
    return sb.product, nil
}
```

#### 8.2.3 性能优化
```go
type ObjectPool struct {
    pool chan interface{}
    new  func() interface{}
}

func NewObjectPool(size int, newFunc func() interface{}) *ObjectPool {
    pool := make(chan interface{}, size)
    for i := 0; i < size; i++ {
        pool <- newFunc()
    }
    
    return &ObjectPool{
        pool: pool,
        new:  newFunc,
    }
}

func (op *ObjectPool) Get() interface{} {
    select {
    case obj := <-op.pool:
        return obj
    default:
        return op.new()
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

### 8.3 测试策略

#### 8.3.1 单元测试
```go
func TestSingleton(t *testing.T) {
    instance1 := GetInstance()
    instance2 := GetInstance()
    
    if instance1 != instance2 {
        t.Error("Singleton instances are not the same")
    }
}

func TestFactoryMethod(t *testing.T) {
    creator := &ConcreteCreator{}
    product := creator.FactoryMethod()
    
    if product == nil {
        t.Error("Factory method returned nil")
    }
    
    result := product.Operation()
    expected := "ConcreteProduct operation"
    
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

#### 8.3.2 性能测试
```go
func BenchmarkSingleton(b *testing.B) {
    for i := 0; i < b.N; i++ {
        GetInstance()
    }
}

func BenchmarkFactoryMethod(b *testing.B) {
    creator := &ConcreteCreator{}
    for i := 0; i < b.N; i++ {
        creator.FactoryMethod()
    }
}
```

---

*最后更新时间: 2024-01-XX*
*版本: 1.0.0* 