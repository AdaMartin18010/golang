# EC-M02: Design Patterns in Go (S-Level)

> **维度**: Engineering-CloudNative / Methodology
> **级别**: S (18+ KB)
> **标签**: #design-patterns #go #creational #structural #behavioral #concurrency
> **权威来源**:
>
> - [Design Patterns: Elements of Reusable Object-Oriented Software](https://en.wikipedia.org/wiki/Design_Patterns) - Gang of Four (1994)
> - [Go Design Patterns](https://www.packtpub.com/product/go-design-patterns/9781786466204) - Mario Castro Contreras (2017)
> - [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/) - Katherine Cox-Buday (2017)
> - [Cloud Native Go](https://www.oreilly.com/library/view/cloud-native-go/9781492076322/) - Matthew A. Titmus (2021)

---

## 1. 设计模式的形式化分类

### 1.1 模式分类体系

```
Design Patterns in Go
├── Creational (创建型)
│   ├── Singleton
│   ├── Factory
│   ├── Abstract Factory
│   ├── Builder
│   └── Prototype
├── Structural (结构型)
│   ├── Adapter
│   ├── Bridge
│   ├── Composite
│   ├── Decorator
│   ├── Facade
│   ├── Flyweight
│   └── Proxy
├── Behavioral (行为型)
│   ├── Chain of Responsibility
│   ├── Command
│   ├── Iterator
│   ├── Mediator
│   ├── Memento
│   ├── Observer
│   ├── State
│   ├── Strategy
│   ├── Template Method
│   └── Visitor
└── Concurrency (并发型)
    ├── Barrier
    ├── Future/Promise
    ├── Pipeline
    ├── Worker Pool
    ├── Fan-Out/Fan-In
    └── Cancellation
```

### 1.2 形式化定义

**定义 1.1 (设计模式)**
设计模式 $P$ 是在特定上下文 $C$ 中针对重复出现问题 $O$ 的可重用解决方案 $S$：

```
P = ⟨C, O, S, F, T⟩
```

其中：

- $C$: 适用上下文
- $O$: 解决的问题
- $S$: 解决方案结构
- $F$: 参与的角色（ forces）
- $T$: 权衡与后果

---

## 2. 创建型模式 (Creational Patterns)

### 2.1 单例模式 (Singleton)

**意图**: 确保一个类只有一个实例，并提供一个全局访问点。

**Go 实现**:

```go
package singleton

import (
    "sync"
    "sync/atomic"
)

// 使用 sync.Once 实现线程安全单例
type Database struct {
    conn   string
    status atomic.Int32
}

var (
    instance *Database
    once     sync.Once
)

// GetInstance 返回单例实例
func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{conn: "connected"}
        instance.status.Store(1)
    })
    return instance
}

// 替代方案: 依赖注入（推荐用于测试）
type Service struct {
    db DatabaseInterface
}

func NewService(db DatabaseInterface) *Service {
    return &Service{db: db}
}

// 接口便于 mock
type DatabaseInterface interface {
    Query(sql string) ([]Row, error)
    Close() error
}
```

**生产考虑**:

- 单例难以测试，优先使用依赖注入
- `sync.Once` 保证线程安全且高效
- 考虑延迟初始化的替代方案

### 2.2 工厂模式 (Factory)

**意图**: 定义创建对象的接口，让子类决定实例化哪个类。

**简单工厂**:

```go
package factory

import "fmt"

// Transport 运输接口
type Transport interface {
    Deliver() string
    Cost() float64
}

// Truck 卡车
type Truck struct {
    capacity int
}

func (t *Truck) Deliver() string {
    return "Deliver by land in a box"
}

func (t *Truck) Cost() float64 {
    return 100.0
}

// Ship 轮船
type Ship struct {
    capacity int
}

func (s *Ship) Deliver() string {
    return "Deliver by sea in a container"
}

func (s *Ship) Cost() float64 {
    return 50.0
}

// Plane 飞机
type Plane struct {
    capacity int
}

func (p *Plane) Deliver() string {
    return "Deliver by air"
}

func (p *Plane) Cost() float64 {
    return 500.0
}

// TransportType 运输类型
type TransportType string

const (
    TransportTruck TransportType = "truck"
    TransportShip  TransportType = "ship"
    TransportPlane TransportType = "plane"
)

// CreateTransport 工厂函数
func CreateTransport(typ TransportType) (Transport, error) {
    switch typ {
    case TransportTruck:
        return &Truck{capacity: 1000}, nil
    case TransportShip:
        return &Ship{capacity: 10000}, nil
    case TransportPlane:
        return &Plane{capacity: 100}, nil
    default:
        return nil, fmt.Errorf("unknown transport type: %s", typ)
    }
}

// 使用
func main() {
    transport, err := factory.CreateTransport(factory.TransportTruck)
    if err != nil {
        panic(err)
    }
    fmt.Println(transport.Deliver())
}
```

**抽象工厂**:

```go
// 抽象工厂 - UI 组件族
type UIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
    CreateTextField() TextField
}

type Button interface {
    Render()
    OnClick()
}

type Checkbox interface {
    Render()
    Toggle()
}

type TextField interface {
    Render()
    SetText(text string)
}

// Windows 工厂
type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (w *WindowsFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{}
}

func (w *WindowsFactory) CreateTextField() TextField {
    return &WindowsTextField{}
}

// Mac 工厂
type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
    return &MacButton{}
}

func (m *MacFactory) CreateCheckbox() Checkbox {
    return &MacCheckbox{}
}

func (m *MacFactory) CreateTextField() TextField {
    return &MacTextField{}
}

// 应用程序不关心具体平台
func CreateUI(factory UIFactory) {
    button := factory.CreateButton()
    checkbox := factory.CreateCheckbox()

    button.Render()
    checkbox.Render()
}
```

### 2.3 构建者模式 (Builder)

**意图**: 将复杂对象的构建与其表示分离。

**Go 惯用写法 - 函数选项模式**:

```go
package builder

import (
    "crypto/tls"
    "fmt"
    "time"
)

// Server 配置
type Server struct {
    host        string
    port        int
    timeout     time.Duration
    maxConn     int
    tls         *tls.Config
    middlewares []Middleware
}

// ServerOption 函数选项类型
type ServerOption func(*Server) error

// WithHost 设置主机
func WithHost(host string) ServerOption {
    return func(s *Server) error {
        if host == "" {
            return fmt.Errorf("host cannot be empty")
        }
        s.host = host
        return nil
    }
}

// WithPort 设置端口
func WithPort(port int) ServerOption {
    return func(s *Server) error {
        if port <= 0 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
        s.port = port
        return nil
    }
}

// WithTimeout 设置超时
func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = d
    }
}

// WithMaxConnections 设置最大连接数
func WithMaxConnections(n int) ServerOption {
    return func(s *Server) {
        s.maxConn = n
    }
}

// WithTLS 设置 TLS
func WithTLS(config *tls.Config) ServerOption {
    return func(s *Server) {
        s.tls = config
    }
}

// WithMiddleware 添加中间件
func WithMiddleware(mw ...Middleware) ServerOption {
    return func(s *Server) {
        s.middlewares = append(s.middlewares, mw...)
    }
}

// NewServer 创建服务器
func NewServer(opts ...ServerOption) (*Server, error) {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
        maxConn: 100,
    }

    for _, opt := range opts {
        if err := opt(s); err != nil {
            return nil, err
        }
    }

    return s, nil
}

type Middleware func(Handler) Handler
type Handler func(Request) Response

// 使用示例
func main() {
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS12,
    }

    srv, err := NewServer(
        WithHost("0.0.0.0"),
        WithPort(9090),
        WithTimeout(60*time.Second),
        WithMaxConnections(1000),
        WithTLS(tlsConfig),
        WithMiddleware(LoggingMiddleware, AuthMiddleware),
    )
    if err != nil {
        panic(err)
    }

    _ = srv
}

func LoggingMiddleware(next Handler) Handler {
    return func(req Request) Response {
        fmt.Println("Logging request")
        return next(req)
    }
}

func AuthMiddleware(next Handler) Handler {
    return func(req Request) Response {
        fmt.Println("Authenticating")
        return next(req)
    }
}

type Request struct{}
type Response struct{}
```

---

## 3. 结构型模式 (Structural Patterns)

### 3.1 适配器模式 (Adapter)

**意图**: 将类的接口转换成客户希望的另一个接口。

```go
package adapter

import "fmt"

// Target 目标接口 - 我们期望的接口
type JSONProcessor interface {
    ProcessJSON(data []byte) error
    GetResult() string
}

// Adaptee 被适配者 - 已存在的接口
type XMLService struct {
    xmlData string
}

func (x *XMLService) ProcessXML(data []byte) error {
    x.xmlData = string(data)
    fmt.Println("Processing XML:", x.xmlData)
    return nil
}

func (x *XMLService) GetXMLResult() string {
    return x.xmlData
}

// Adapter 适配器
type XMLToJSONAdapter struct {
    xmlService *XMLService
    result     string
}

func NewXMLToJSONAdapter(service *XMLService) *XMLToJSONAdapter {
    return &XMLToJSONAdapter{xmlService: service}
}

func (a *XMLToJSONAdapter) ProcessJSON(data []byte) error {
    // 转换 JSON 到 XML
    xmlData := convertJSONToXML(data)

    // 使用 XML 服务处理
    if err := a.xmlService.ProcessXML(xmlData); err != nil {
        return err
    }

    // 转换结果回 JSON
    a.result = convertXMLToJSON(a.xmlService.GetXMLResult())
    return nil
}

func (a *XMLToJSONAdapter) GetResult() string {
    return a.result
}

func convertJSONToXML(data []byte) []byte {
    // 转换逻辑
    return data
}

func convertXMLToJSON(data string) string {
    // 转换逻辑
    return data
}

// 使用
func ProcessWithAdapter(adapter JSONProcessor, jsonData []byte) error {
    return adapter.ProcessJSON(jsonData)
}
```

### 3.2 装饰器模式 (Decorator)

**意图**: 动态地给对象添加额外职责。

```go
package decorator

import (
    "context"
    "fmt"
    "log/slog"
    "time"
)

// Handler 处理器类型
type Handler func(ctx context.Context, req *Request) (*Response, error)

// Request/Response 类型
type Request struct {
    ID      string
    Payload interface{}
}

type Response struct {
    Data  interface{}
    Error error
}

// 日志装饰器
func Logging(logger *slog.Logger) func(Handler) Handler {
    return func(next Handler) Handler {
        return func(ctx context.Context, req *Request) (*Response, error) {
            start := time.Now()
            logger.Info("request started",
                "id", req.ID,
                "time", start,
            )

            resp, err := next(ctx, req)

            logger.Info("request completed",
                "id", req.ID,
                "duration", time.Since(start),
                "error", err,
            )
            return resp, err
        }
    }
}

// 重试装饰器
func Retry(maxAttempts int, delay time.Duration) func(Handler) Handler {
    return func(next Handler) Handler {
        return func(ctx context.Context, req *Request) (*Response, error) {
            var resp *Response
            var err error

            for i := 0; i < maxAttempts; i++ {
                resp, err = next(ctx, req)
                if err == nil {
                    return resp, nil
                }

                // 指数退避
                backoff := delay * time.Duration(1<<i)
                select {
                case <-time.After(backoff):
                case <-ctx.Done():
                    return nil, ctx.Err()
                }
            }

            return nil, fmt.Errorf("failed after %d attempts: %w", maxAttempts, err)
        }
    }
}

// 超时装饰器
func Timeout(d time.Duration) func(Handler) Handler {
    return func(next Handler) Handler {
        return func(ctx context.Context, req *Request) (*Response, error) {
            ctx, cancel := context.WithTimeout(ctx, d)
            defer cancel()
            return next(ctx, req)
        }
    }
}

// 限流装饰器
func RateLimiter(limit int) func(Handler) Handler {
    semaphore := make(chan struct{}, limit)

    return func(next Handler) Handler {
        return func(ctx context.Context, req *Request) (*Response, error) {
            select {
            case semaphore <- struct{}{}:
                defer func() { <-semaphore }()
                return next(ctx, req)
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }
    }
}

// 业务逻辑
func BusinessLogic(ctx context.Context, req *Request) (*Response, error) {
    // 处理业务
    return &Response{Data: "success"}, nil
}

// 使用 - 链式装饰
func main() {
    logger := slog.Default()

    handler := RateLimiter(10)(
        Timeout(30*time.Second)(
            Retry(3, time.Second)(
                Logging(logger)(
                    BusinessLogic,
                ),
            ),
        ),
    )

    req := &Request{ID: "123", Payload: "data"}
    resp, err := handler(context.Background(), req)
    _ = resp
    _ = err
}
```

### 3.3 外观模式 (Facade)

**意图**: 为子系统中的一组接口提供一个统一的高层接口。

```go
package facade

import (
    "errors"
    "fmt"
)

// 复杂子系统
type InventoryService struct {
    stock map[string]int
}

func NewInventoryService() *InventoryService {
    return &InventoryService{
        stock: map[string]int{
            "product1": 100,
            "product2": 50,
        },
    }
}

func (i *InventoryService) Check(productID string, quantity int) bool {
    return i.stock[productID] >= quantity
}

func (i *InventoryService) Reserve(productID string, quantity int) error {
    if !i.Check(productID, quantity) {
        return errors.New("insufficient stock")
    }
    i.stock[productID] -= quantity
    return nil
}

func (i *InventoryService) Release(productID string, quantity int) {
    i.stock[productID] += quantity
}

type PaymentService struct {
    transactions []Transaction
}

type Transaction struct {
    OrderID string
    Amount  float64
    Status  string
}

func (p *PaymentService) Process(orderID string, amount float64) (string, error) {
    tx := Transaction{
        OrderID: orderID,
        Amount:  amount,
        Status:  "completed",
    }
    p.transactions = append(p.transactions, tx)
    return tx.OrderID, nil
}

func (p *PaymentService) Refund(transactionID string) error {
    // 退款逻辑
    return nil
}

type ShippingService struct {
    shipments []Shipment
}

type Shipment struct {
    OrderID string
    Status  string
}

func (s *ShippingService) Schedule(orderID string, address string) (string, error) {
    shipment := Shipment{
        OrderID: orderID,
        Status:  "scheduled",
    }
    s.shipments = append(s.shipments, shipment)
    return shipment.OrderID, nil
}

type NotificationService struct{}

func (n *NotificationService) SendEmail(to, subject, body string) error {
    fmt.Printf("Sending email to %s: %s\n", to, subject)
    return nil
}

func (n *NotificationService) SendSMS(phone, message string) error {
    fmt.Printf("Sending SMS to %s: %s\n", phone, message)
    return nil
}

// OrderFacade 外观
type OrderFacade struct {
    inventory *InventoryService
    payment   *PaymentService
    shipping  *ShippingService
    notify    *NotificationService
}

func NewOrderFacade() *OrderFacade {
    return &OrderFacade{
        inventory: NewInventoryService(),
        payment:   &PaymentService{},
        shipping:  &ShippingService{},
        notify:    &NotificationService{},
    }
}

type OrderRequest struct {
    UserID      string
    ProductID   string
    Quantity    int
    Amount      float64
    Email       string
    Phone       string
    Address     string
}

type OrderResult struct {
    OrderID       string
    TransactionID string
    ShipmentID    string
    Status        string
}

func (f *OrderFacade) PlaceOrder(req OrderRequest) (*OrderResult, error) {
    // 1. 检查库存
    if !f.inventory.Check(req.ProductID, req.Quantity) {
        return nil, errors.New("product out of stock")
    }

    // 2. 预留库存
    if err := f.inventory.Reserve(req.ProductID, req.Quantity); err != nil {
        return nil, err
    }

    // 3. 处理支付
    txID, err := f.payment.Process(req.UserID, req.Amount)
    if err != nil {
        f.inventory.Release(req.ProductID, req.Quantity)
        return nil, fmt.Errorf("payment failed: %w", err)
    }

    // 4. 安排配送
    shipmentID, err := f.shipping.Schedule(req.UserID, req.Address)
    if err != nil {
        f.inventory.Release(req.ProductID, req.Quantity)
        f.payment.Refund(txID)
        return nil, fmt.Errorf("shipping failed: %w", err)
    }

    // 5. 发送通知
    f.notify.SendEmail(req.Email, "Order Confirmed", "Your order has been placed")
    f.notify.SendSMS(req.Phone, "Order confirmed!")

    return &OrderResult{
        OrderID:       req.UserID,
        TransactionID: txID,
        ShipmentID:    shipmentID,
        Status:        "confirmed",
    }, nil
}
```

---

## 4. 行为型模式 (Behavioral Patterns)

### 4.1 策略模式 (Strategy)

**意图**: 定义算法族，分别封装起来，让它们可以互相替换。

```go
package strategy

import "fmt"

// PaymentStrategy 支付策略接口
type PaymentStrategy interface {
    Pay(amount float64) error
    GetName() string
}

// CreditCard 信用卡策略
type CreditCard struct {
    number     string
    cvv        string
    expiryDate string
}

func NewCreditCard(number, cvv, expiry string) *CreditCard {
    return &CreditCard{
        number:     number,
        cvv:        cvv,
        expiryDate: expiry,
    }
}

func (c *CreditCard) Pay(amount float64) error {
    fmt.Printf("Paying %.2f using Credit Card ending with %s\n",
        amount, c.number[len(c.number)-4:])
    return nil
}

func (c *CreditCard) GetName() string {
    return "CreditCard"
}

// PayPal PayPal策略
type PayPal struct {
    email    string
    password string
}

func NewPayPal(email, password string) *PayPal {
    return &PayPal{email: email, password: password}
}

func (p *PayPal) Pay(amount float64) error {
    fmt.Printf("Paying %.2f using PayPal account %s\n", amount, p.email)
    return nil
}

func (p *PayPal) GetName() string {
    return "PayPal"
}

// Crypto 加密货币策略
type Crypto struct {
    walletAddress string
    currency      string
}

func NewCrypto(address, currency string) *Crypto {
    return &Crypto{walletAddress: address, currency: currency}
}

func (c *Crypto) Pay(amount float64) error {
    fmt.Printf("Paying %.2f %s to wallet %s\n",
        amount, c.currency, c.walletAddress[:8]+"...")
    return nil
}

func (c *Crypto) GetName() string {
    return "Crypto(" + c.currency + ")"
}

// PaymentContext 支付上下文
type PaymentContext struct {
    strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentContext) ExecutePayment(amount float64) error {
    if p.strategy == nil {
        return fmt.Errorf("payment strategy not set")
    }
    return p.strategy.Pay(amount)
}

func (p *PaymentContext) GetStrategyName() string {
    if p.strategy == nil {
        return "none"
    }
    return p.strategy.GetName()
}

// 使用
func main() {
    context := &PaymentContext{}

    // 使用信用卡
    context.SetStrategy(NewCreditCard("1234567890123456", "123", "12/25"))
    context.ExecutePayment(100.50)

    // 切换到 PayPal
    context.SetStrategy(NewPayPal("user@example.com", "password"))
    context.ExecutePayment(50.25)

    // 切换到加密货币
    context.SetStrategy(NewCrypto("0x1234567890abcdef", "ETH"))
    context.ExecutePayment(0.5)
}
```

### 4.2 观察者模式 (Observer)

**意图**: 定义对象间的一对多依赖，当一个对象改变状态，所有依赖者都会收到通知。

```go
package observer

import (
    "fmt"
    "sync"
)

// Event 事件
type Event struct {
    Type    string
    Payload interface{}
}

// Observer 观察者接口
type Observer interface {
    Update(event Event)
    GetID() string
}

// Subject 主题接口
type Subject interface {
    Attach(o Observer)
    Detach(o Observer)
    Notify(event Event)
}

// ConcreteSubject 具体主题
type ConcreteSubject struct {
    observers []Observer
    mu        sync.RWMutex
}

func (s *ConcreteSubject) Attach(o Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers = append(s.observers, o)
}

func (s *ConcreteSubject) Detach(o Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for i, obs := range s.observers {
        if obs.GetID() == o.GetID() {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            break
        }
    }
}

func (s *ConcreteSubject) Notify(event Event) {
    s.mu.RLock()
    observers := make([]Observer, len(s.observers))
    copy(observers, s.observers)
    s.mu.RUnlock()

    for _, o := range observers {
        go o.Update(event) // 异步通知
    }
}

func (s *ConcreteSubject) GetObserverCount() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.observers)
}

// EmailNotifier 邮件通知者
type EmailNotifier struct {
    id    string
    email string
}

func NewEmailNotifier(id, email string) *EmailNotifier {
    return &EmailNotifier{id: id, email: email}
}

func (e *EmailNotifier) Update(event Event) {
    fmt.Printf("[EmailNotifier %s] Sending email to %s for event: %s\n",
        e.id, e.email, event.Type)
}

func (e *EmailNotifier) GetID() string {
    return e.id
}

// SMSNotifier 短信通知者
type SMSNotifier struct {
    id    string
    phone string
}

func NewSMSNotifier(id, phone string) *SMSNotifier {
    return &SMSNotifier{id: id, phone: phone}
}

func (s *SMSNotifier) Update(event Event) {
    fmt.Printf("[SMSNotifier %s] Sending SMS to %s for event: %s\n",
        s.id, s.phone, event.Type)
}

func (s *SMSNotifier) GetID() string {
    return s.id
}

// LogNotifier 日志通知者
type LogNotifier struct {
    id string
}

func NewLogNotifier(id string) *LogNotifier {
    return &LogNotifier{id: id}
}

func (l *LogNotifier) Update(event Event) {
    fmt.Printf("[LogNotifier %s] Logging event: %s, payload: %v\n",
        l.id, event.Type, event.Payload)
}

func (l *LogNotifier) GetID() string {
    return l.id
}

// 使用
func main() {
    subject := &ConcreteSubject{}

    emailObs := NewEmailNotifier("email1", "user@example.com")
    smsObs := NewSMSNotifier("sms1", "+1234567890")
    logObs := NewLogNotifier("log1")

    subject.Attach(emailObs)
    subject.Attach(smsObs)
    subject.Attach(logObs)

    // 发布事件
    subject.Notify(Event{Type: "order_created", Payload: map[string]string{"order_id": "123"}})

    // 移除观察者
    subject.Detach(smsObs)

    subject.Notify(Event{Type: "order_shipped", Payload: map[string]string{"order_id": "123"}})
}
```

### 4.3 模板方法模式 (Template Method)

**意图**: 定义算法骨架，将某些步骤延迟到子类。

```go
package template

import "fmt"

// DataMiner 数据挖掘器接口
type DataMiner interface {
    Mine(path string)
}

// BaseDataMiner 基础数据挖掘器
type BaseDataMiner struct {
    file string
}

func (b *BaseDataMiner) open(file string) {
    b.file = file
    fmt.Println("Opening file:", file)
}

func (b *BaseDataMiner) close() {
    fmt.Println("Closing file:", b.file)
}

// 必须由子类实现的方法
func (b *BaseDataMiner) extract() []byte {
    panic("must implement")
}

func (b *BaseDataMiner) parse(data []byte) string {
    panic("must implement")
}

func (b *BaseDataMiner) analyze(data string) string {
    fmt.Println("Analyzing data...")
    return data
}

func (b *BaseDataMiner) sendReport(data string) {
    fmt.Println("Sending report:", data)
}

// Mine 模板方法
func (b *BaseDataMiner) Mine(path string, extractFunc func() []byte, parseFunc func([]byte) string) {
    b.open(path)

    rawData := extractFunc()
    parsedData := parseFunc(rawData)
    analyzedData := b.analyze(parsedData)
    b.sendReport(analyzedData)

    b.close()
}

// PDFDataMiner PDF 数据挖掘器
type PDFDataMiner struct {
    BaseDataMiner
}

func NewPDFDataMiner() *PDFDataMiner {
    return &PDFDataMiner{}
}

func (p *PDFDataMiner) MinePDF(path string) {
    p.Mine(path, p.extractPDF, p.parsePDF)
}

func (p *PDFDataMiner) extractPDF() []byte {
    fmt.Println("Extracting PDF data")
    return []byte("raw pdf data")
}

func (p *PDFDataMiner) parsePDF(data []byte) string {
    fmt.Println("Parsing PDF data")
    return "parsed pdf: " + string(data)
}

// CSVDataMiner CSV 数据挖掘器
type CSVDataMiner struct {
    BaseDataMiner
}

func NewCSVDataMiner() *CSVDataMiner {
    return &CSVDataMiner{}
}

func (c *CSVDataMiner) MineCSV(path string) {
    c.Mine(path, c.extractCSV, c.parseCSV)
}

func (c *CSVDataMiner) extractCSV() []byte {
    fmt.Println("Extracting CSV data")
    return []byte("raw csv data")
}

func (c *CSVDataMiner) parseCSV(data []byte) string {
    fmt.Println("Parsing CSV data")
    return "parsed csv: " + string(data)
}
```

---

## 5. 并发模式 (Concurrency Patterns)

### 5.1 工作池模式 (Worker Pool)

**意图**: 限制同时执行的任务数量，复用 goroutine。

```go
package workerpool

import (
    "context"
    "fmt"
    "sync"
)

// Job 任务接口
type Job interface {
    Process() error
    GetID() string
}

// Pool 工作池
type Pool struct {
    workers   int
    jobQueue  chan Job
    resultQueue chan Result
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

// Result 处理结果
type Result struct {
    JobID string
    Error error
}

// NewPool 创建工作池
func NewPool(workers int, queueSize int) *Pool {
    ctx, cancel := context.WithCancel(context.Background())
    return &Pool{
        workers:     workers,
        jobQueue:    make(chan Job, queueSize),
        resultQueue: make(chan Result, queueSize),
        ctx:         ctx,
        cancel:      cancel,
    }
}

// Start 启动工作池
func (p *Pool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *Pool) worker(id int) {
    defer p.wg.Done()

    for {
        select {
        case job, ok := <-p.jobQueue:
            if !ok {
                return
            }
            err := job.Process()
            p.resultQueue <- Result{JobID: job.GetID(), Error: err}

        case <-p.ctx.Done():
            return
        }
    }
}

// Submit 提交任务
func (p *Pool) Submit(job Job) bool {
    select {
    case p.jobQueue <- job:
        return true
    case <-p.ctx.Done():
        return false
    }
}

// GetResults 获取结果通道
func (p *Pool) GetResults() <-chan Result {
    return p.resultQueue
}

// Stop 停止工作池
func (p *Pool) Stop() {
    p.cancel()
    close(p.jobQueue)
    p.wg.Wait()
    close(p.resultQueue)
}

// 示例任务
type SimpleJob struct {
    id   string
    data string
}

func (j *SimpleJob) Process() error {
    fmt.Printf("Processing job %s with data: %s\n", j.id, j.data)
    return nil
}

func (j *SimpleJob) GetID() string {
    return j.id
}
```

### 5.2 Pipeline 模式

**意图**: 通过通道连接处理阶段，构建数据流。

```go
package pipeline

import "context"

// Item 数据项
type Item struct {
    ID    string
    Value interface{}
}

// Stage 处理阶段
type Stage func(ctx context.Context, in <-chan Item) <-chan Item

// Pipeline 构建管道
func Pipeline(stages ...Stage) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := in
        for _, stage := range stages {
            out = stage(ctx, out)
        }
        return out
    }
}

// Generator 生成器阶段
func Generator(items ...Item) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item)
        go func() {
            defer close(out)
            for _, item := range items {
                select {
                case out <- item:
                case <-ctx.Done():
                    return
                }
            }
        }()
        return out
    }
}

// Filter 过滤阶段
func Filter(predicate func(Item) bool) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item)
        go func() {
            defer close(out)
            for item := range in {
                if predicate(item) {
                    select {
                    case out <- item:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
        return out
    }
}

// Map 映射阶段
func Map(transform func(Item) Item) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item)
        go func() {
            defer close(out)
            for item := range in {
                select {
                case out <- transform(item):
                case <-ctx.Done():
                    return
                }
            }
        }()
        return out
    }
}

// Buffer 缓冲阶段
func Buffer(size int) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item, size)
        go func() {
            defer close(out)
            for item := range in {
                select {
                case out <- item:
                case <-ctx.Done():
                    return
                }
            }
        }()
        return out
    }
}

// FanOut 扇出阶段
func FanOut(workers int, stage Stage) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item)

        // 创建多个 worker
        inputs := make([]chan Item, workers)
        for i := 0; i < workers; i++ {
            inputs[i] = make(chan Item)
            go func(input chan Item) {
                result := stage(ctx, input)
                for item := range result {
                    select {
                    case out <- item:
                    case <-ctx.Done():
                        return
                    }
                }
            }(inputs[i])
        }

        // 分发输入
        go func() {
            defer close(out)
            i := 0
            for item := range in {
                select {
                case inputs[i%workers] <- item:
                    i++
                case <-ctx.Done():
                    return
                }
            }
            for _, ch := range inputs {
                close(ch)
            }
        }()

        return out
    }
}

// FanIn 扇入阶段
func FanIn(inputs ...<-chan Item) Stage {
    return func(ctx context.Context, in <-chan Item) <-chan Item {
        out := make(chan Item)
        var wg sync.WaitGroup

        output := func(c <-chan Item) {
            defer wg.Done()
            for item := range c {
                select {
                case out <- item:
                case <-ctx.Done():
                    return
                }
            }
        }

        wg.Add(len(inputs))
        for _, input := range inputs {
            go output(input)
        }

        go func() {
            wg.Wait()
            close(out)
        }()

        return out
    }
}
```

---

## 6. 架构决策记录

### ADR-001: 优先使用组合而非继承

**状态**: 已接受

**背景**: Go 不支持传统面向对象的继承，使用结构体嵌入和接口组合。

**决策**: 使用组合实现代码复用，通过接口定义行为契约。

**后果**:

- 更灵活的设计
- 避免继承层次过深
- 更好的可测试性

### ADR-002: 使用函数选项模式替代 Builder

**状态**: 已接受

**背景**: 构造函数参数过多时，需要更好的方式来传递可选配置。

**决策**: 对于可选参数，使用函数选项模式（Functional Options）。

**后果**:

- API 向后兼容
- 自文档化
- 编译时类型安全
- 可以验证参数

### ADR-003: 优先使用接口定义依赖

**状态**: 已接受

**背景**: 需要解耦模块，便于测试和替换实现。

**决策**: 依赖接口而非具体类型，接口定义在消费者侧。

**后果**:

- 松耦合设计
- 易于 mock 测试
- 清晰的服务边界

---

## 7. 生产检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Design Patterns Production Checklist                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  单例模式:                                                                   │
│  □ 使用 sync.Once 保证线程安全                                               │
│  □ 考虑依赖注入替代全局状态                                                   │
│  □ 单例可序列化时需要特殊处理                                                 │
│                                                                              │
│  工厂模式:                                                                   │
│  □ 工厂返回接口而非具体类型                                                   │
│  □ 错误处理完整（未知类型返回错误）                                            │
│  □ 支持配置参数传递                                                          │
│                                                                              │
│  构建者模式:                                                                 │
│  □ 提供合理的默认值                                                          │
│  □ 验证必需参数                                                              │
│  □ 支持错误返回                                                              │
│                                                                              │
│  装饰器模式:                                                                 │
│  □ 保持被装饰接口的一致性                                                     │
│  □ 支持装饰器链式组合                                                         │
│  □ 注意性能开销（嵌套调用）                                                   │
│                                                                              │
│  观察者模式:                                                                 │
│  □ 异步通知防止阻塞                                                          │
│  □ 支持取消订阅防止内存泄漏                                                    │
│  □ 处理 panic 防止单个观察者影响其他                                          │
│                                                                              │
│  工作池模式:                                                                 │
│  □ 设置合理的队列大小防止 OOM                                                 │
│  □ 支持优雅关闭                                                              │
│  □ 监控队列深度和工作线程利用率                                               │
│                                                                              │
│  Pipeline:                                                                   │
│  □ 支持 context 取消                                                         │
│  □ 处理通道关闭                                                              │
│  □ 背压控制防止内存无限增长                                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. 参考文献

1. **Gamma, E., Helm, R., Johnson, R., & Vlissides, J. (1994)**. Design Patterns: Elements of Reusable Object-Oriented Software. *Addison-Wesley*.
2. **Contreras, M. C. (2017)**. Go Design Patterns. *Packt Publishing*.
3. **Cox-Buday, K. (2017)**. Concurrency in Go. *O'Reilly Media*.
4. **Titmus, M. A. (2021)**. Cloud Native Go. *O'Reilly Media*.
5. **Hoare, C. A. R. (1978)**. Communicating Sequential Processes. *CACM*.

---

**质量评级**: S (18+ KB, 完整形式化 + 生产实践)
