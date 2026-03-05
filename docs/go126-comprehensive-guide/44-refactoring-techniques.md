# Go代码重构技术

> 改善代码质量、可维护性和性能的重构实践

---

## 一、重构基础原则

### 1.1 为什么需要重构

```text
技术债务：
────────────────────────────────────────

什么是技术债务：
- 为了快速交付而做出的妥协
- 设计不完善但"能用"的代码
- 缺少测试的代码

技术债务成本：
┌─────────────────────────────────────┐
│  时间   │ 债务成本 / 初始开发成本   │
├─────────────────────────────────────┤
│  1个月  │        1.2x              │
│  6个月  │        2.0x              │
│  1年    │        4.0x              │
│  2年    │       10.0x+             │
└─────────────────────────────────────┘

重构 vs 重写：
────────────────────────────────────────

重构：
- 小步修改
- 保持功能不变
- 逐步改善设计
- 风险低

重写：
- 从零开始
- 功能可能改变
- 高风险
- 时间长

选择重构的条件：
- 代码基础大部分可用
- 有测试覆盖
- 问题局部化

选择重写的条件：
- 架构严重错误
- 技术栈过时
- 维护成本过高

重构安全网：
────────────────────────────────────────

1. 测试覆盖：
   - 重构前确保有测试
   - 测试通过 = 功能正确
   - 没有测试，先写测试

2. 版本控制：
   - 小步提交
   - 清晰的提交信息
   - 可以随时回滚

3. 代码审查：
   - 重构后让同事review
   - 发现潜在问题
   - 知识共享
```

### 1.2 代码异味识别

```text
常见代码异味：
────────────────────────────────────────

1. 过长函数：

// 异味
func ProcessOrder(order Order) error {
    // 100+ 行代码
    // 验证、计算、数据库、通知...
}

// 重构后
func ProcessOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    total := calculateTotal(order)
    if err := saveOrder(order, total); err != nil {
        return err
    }
    return notifyCustomer(order)
}

2. 重复代码：

// 异味
func ProcessA(data Data) {
    // 10行处理逻辑
}

func ProcessB(data Data) {
    // 同样的10行处理逻辑
}

// 重构后
func processCommon(data Data) {
    // 10行处理逻辑
}

func ProcessA(data Data) {
    processCommon(data)
    // A特有逻辑
}

3. 过长参数列表：

// 异味
func CreateUser(name, email, phone, address, city, country string, age int, active bool)

// 重构后
type UserProfile struct {
    Name, Email, Phone string
    Address            Address
    Age                int
    Active             bool
}

func CreateUser(profile UserProfile)

4. 上帝对象：

// 异味
type OrderManager struct {
    // 包含订单、库存、支付、物流等所有功能
}

// 重构后
type OrderService struct {
    inventory *InventoryService
    payment   *PaymentService
    shipping  *ShippingService
}
```

---

## 二、具体重构技术

### 2.1 提取函数

```text
何时提取：
────────────────────────────────────────

- 代码块可以被命名描述
- 代码块被多处使用
- 函数过长需要拆分

示例：
────────────────────────────────────────

// 重构前
func CalculateTotal(items []Item) float64 {
    var subtotal float64
    for _, item := range items {
        price := item.Price
        if item.Discount > 0 {
            price = price * (1 - item.Discount)
        }
        subtotal += price * float64(item.Quantity)
    }

    tax := subtotal * 0.08
    if subtotal > 100 {
        tax = tax * 0.95  // 税收折扣
    }

    return subtotal + tax
}

// 重构后
func CalculateTotal(items []Item) float64 {
    subtotal := calculateSubtotal(items)
    tax := calculateTax(subtotal)
    return subtotal + tax
}

func calculateSubtotal(items []Item) float64 {
    var total float64
    for _, item := range items {
        total += calculateItemPrice(item)
    }
    return total
}

func calculateItemPrice(item Item) float64 {
    price := item.Price
    if item.Discount > 0 {
        price = price * (1 - item.Discount)
    }
    return price * float64(item.Quantity)
}

func calculateTax(subtotal float64) float64 {
    tax := subtotal * 0.08
    if subtotal > 100 {
        tax = tax * 0.95
    }
    return tax
}

好处：
- 每个函数职责单一
- 函数名即文档
- 易于测试
- 易于复用
```

### 2.2 引入接口

```text
解耦依赖：
────────────────────────────────────────

// 重构前：强耦合
type OrderService struct {
    db *sql.DB
}

func (s *OrderService) CreateOrder(o Order) error {
    _, err := s.db.Exec("INSERT INTO orders ...", o.ID)
    return err
}

// 重构后：依赖接口
type OrderRepository interface {
    Save(Order) error
    GetByID(int) (Order, error)
}

type OrderService struct {
    repo OrderRepository
}

// 实现1：数据库
type DBRepository struct {
    db *sql.DB
}

func (r *DBRepository) Save(o Order) error {
    _, err := r.db.Exec("INSERT INTO orders ...", o.ID)
    return err
}

// 实现2：内存（用于测试）
type MemoryRepository struct {
    orders map[int]Order
}

func (r *MemoryRepository) Save(o Order) error {
    r.orders[o.ID] = o
    return nil
}

// 测试变得容易
func TestOrderService(t *testing.T) {
    repo := &MemoryRepository{orders: make(map[int]Order)}
    service := &OrderService{repo: repo}
    // 测试...
}

接口隔离：
────────────────────────────────────────

// 不良：大接口
type DataStore interface {
    Get(key string) (Value, error)
    Set(key string, value Value) error
    Delete(key string) error
    List(prefix string) ([]Value, error)
    Watch(key string) (<-chan Event, error)
    // ...更多方法
}

// 良好：小接口组合
type Reader interface {
    Get(key string) (Value, error)
}

type Writer interface {
    Set(key string, value Value) error
    Delete(key string) error
}

type Watcher interface {
    Watch(key string) (<-chan Event, error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

### 2.3 消除重复

```text
模板方法模式：
────────────────────────────────────────

// 重构前：重复的处理流程
type EmailNotifier struct{}

func (n *EmailNotifier) Notify(user User, message string) error {
    if !user.Active {
        return nil
    }
    // 发送邮件...
}

type SMSNotifier struct{}

func (n *SMSNotifier) Notify(user User, message string) error {
    if !user.Active {
        return nil
    }
    // 发送短信...
}

// 重构后：提取公共流程
type Notifier interface {
    Notify(User, string) error
}

type baseNotifier struct{}

func (n *baseNotifier) notifyIfActive(user User, fn func() error) error {
    if !user.Active {
        return nil
    }
    return fn()
}

type EmailNotifier struct {
    baseNotifier
}

func (n *EmailNotifier) Notify(user User, message string) error {
    return n.notifyIfActive(user, func() error {
        // 发送邮件...
        return nil
    })
}

策略模式：
────────────────────────────────────────

// 重构前：大量if-else
func CalculateShipping(order Order) float64 {
    switch order.Method {
    case "standard":
        return 5.0
    case "express":
        return 15.0
    case "free":
        if order.Total > 50 {
            return 0
        }
        return 5.0
    }
    return 0
}

// 重构后：策略模式
type ShippingStrategy interface {
    Calculate(order Order) float64
}

type StandardShipping struct{}

func (s *StandardShipping) Calculate(order Order) float64 {
    return 5.0
}

type ExpressShipping struct{}

func (s *ExpressShipping) Calculate(order Order) float64 {
    return 15.0
}

type FreeShipping struct{}

func (s *FreeShipping) Calculate(order Order) float64 {
    if order.Total > 50 {
        return 0
    }
    return 5.0
}

var strategies = map[string]ShippingStrategy{
    "standard": &StandardShipping{},
    "express":  &ExpressShipping{},
    "free":     &FreeShipping{},
}

func CalculateShipping(order Order) float64 {
    strategy, ok := strategies[order.Method]
    if !ok {
        return 0
    }
    return strategy.Calculate(order)
}
```

---

## 三、性能重构

### 3.1 内存优化

```text
减少分配：
────────────────────────────────────────

// 重构前：每次调用都分配
func process(data []byte) []byte {
    result := make([]byte, 0)
    for _, b := range data {
        result = append(result, transform(b))
    }
    return result
}

// 重构后：预分配容量
func process(data []byte) []byte {
    result := make([]byte, 0, len(data))
    for _, b := range data {
        result = append(result, transform(b))
    }
    return result
}

对象池：
────────────────────────────────────────

// 重构前：频繁创建大对象
func handleRequest(w http.ResponseWriter, r *http.Request) {
    buf := make([]byte, 64*1024)  // 64KB
    // 使用buf...
}

// 重构后：使用sync.Pool
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 64*1024)
    },
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // 使用buf...
}
```

### 3.2 并发优化

```text
串行改并行：
────────────────────────────────────────

// 重构前：串行处理
func processAll(items []Item) []Result {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = process(item)
    }
    return results
}

// 重构后：并行处理
func processAll(items []Item) []Result {
    results := make([]Result, len(items))
    var wg sync.WaitGroup
    wg.Add(len(items))

    for i, item := range items {
        go func(i int, item Item) {
            defer wg.Done()
            results[i] = process(item)
        }(i, item)
    }

    wg.Wait()
    return results
}

批量处理：
────────────────────────────────────────

// 重构前：逐个处理（大量数据库往返）
func saveUsers(users []User) error {
    for _, u := range users {
        if err := db.Save(&u); err != nil {
            return err
        }
    }
    return nil
}

// 重构后：批量插入（一次数据库操作）
func saveUsers(users []User) error {
    return db.CreateInBatches(users, 100).Error
}
```

---

*本章提供了代码重构的完整技术指南。*
