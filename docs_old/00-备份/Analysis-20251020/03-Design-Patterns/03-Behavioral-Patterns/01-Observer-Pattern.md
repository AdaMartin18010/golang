# 3.3.1 观察者模式 (Observer Pattern)

## 3.3.1.1 目录

## 3.3.1.2 1. 概述

### 3.3.1.2.1 定义

观察者模式定义对象间的一种一对多的依赖关系，当一个对象的状态发生改变时，所有依赖于它的对象都得到通知并被自动更新。

**形式化定义**:
$$Observer = (Subject, Observer, ConcreteSubject, ConcreteObserver, Notification)$$

其中：

- $Subject$ 是主题接口
- $Observer$ 是观察者接口
- $ConcreteSubject$ 是具体主题
- $ConcreteObserver$ 是具体观察者
- $Notification$ 是通知机制

### 3.3.1.2.2 核心特征

- **一对多关系**: 一个主题可以有多个观察者
- **松耦合**: 主题和观察者之间松耦合
- **自动通知**: 状态变化时自动通知观察者
- **动态订阅**: 观察者可以动态订阅和取消订阅

## 3.3.1.3 2. 理论基础

### 3.3.1.3.1 数学形式化

**定义 2.1** (观察者模式): 观察者模式是一个六元组 $O = (S, O, N, A, F, V)$

其中：

- $S$ 是主题集合
- $O$ 是观察者集合
- $N$ 是通知函数，$N: S \times O \rightarrow Notification$
- $A$ 是附加函数，$A: S \times O \rightarrow S$
- $F$ 是分离函数，$F: S \times O \rightarrow S$
- $V$ 是验证规则

**定理 2.1** (通知传递性): 对于任意主题 $s \in S$ 和观察者 $o \in O$，如果 $o$ 订阅了 $s$，则 $s$ 的状态变化会通知 $o$

**证明**: 由通知函数的实现保证。

### 3.3.1.3.2 范畴论视角

在范畴论中，观察者模式可以表示为：

$$Observer : Subject \times Observer \rightarrow Notification$$

其中 $Subject$ 和 $Observer$ 是对象范畴，$Notification$ 是通知范畴。

## 3.3.1.4 3. Go语言实现

### 3.3.1.4.1 基础观察者模式

```go
package observer

import (
    "fmt"
    "sync"
)

// Observer 观察者接口
type Observer interface {
    Update(subject Subject)
    GetID() string
}

// Subject 主题接口
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
    GetState() interface{}
    SetState(state interface{})
}

// ConcreteSubject 具体主题
type ConcreteSubject struct {
    observers []Observer
    state     interface{}
    mu        sync.RWMutex
}

func NewConcreteSubject() *ConcreteSubject {
    return &ConcreteSubject{
        observers: make([]Observer, 0),
    }
}

func (s *ConcreteSubject) Attach(observer Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 检查是否已经存在
    for _, obs := range s.observers {
        if obs.GetID() == observer.GetID() {
            return
        }
    }
    
    s.observers = append(s.observers, observer)
    fmt.Printf("Observer %s attached to subject\n", observer.GetID())
}

func (s *ConcreteSubject) Detach(observer Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, obs := range s.observers {
        if obs.GetID() == observer.GetID() {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            fmt.Printf("Observer %s detached from subject\n", observer.GetID())
            return
        }
    }
}

func (s *ConcreteSubject) Notify() {
    s.mu.RLock()
    observers := make([]Observer, len(s.observers))
    copy(observers, s.observers)
    s.mu.RUnlock()
    
    fmt.Printf("Subject notifying %d observers\n", len(observers))
    for _, observer := range observers {
        observer.Update(s)
    }
}

func (s *ConcreteSubject) GetState() interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state
}

func (s *ConcreteSubject) SetState(state interface{}) {
    s.mu.Lock()
    s.state = state
    s.mu.Unlock()
    
    fmt.Printf("Subject state changed to: %v\n", state)
    s.Notify()
}

func (s *ConcreteSubject) GetObserverCount() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.observers)
}

// ConcreteObserver 具体观察者
type ConcreteObserver struct {
    id   string
    name string
}

func NewConcreteObserver(id, name string) *ConcreteObserver {
    return &ConcreteObserver{
        id:   id,
        name: name,
    }
}

func (o *ConcreteObserver) Update(subject Subject) {
    state := subject.GetState()
    fmt.Printf("Observer %s (%s) received update: %v\n", o.name, o.id, state)
}

func (o *ConcreteObserver) GetID() string {
    return o.id
}

func (o *ConcreteObserver) GetName() string {
    return o.name
}

```

### 3.3.1.4.2 事件驱动观察者模式

```go
package eventobserver

import (
    "fmt"
    "sync"
    "time"
)

// Event 事件
type Event struct {
    Type      string
    Data      interface{}
    Timestamp time.Time
    Source    string
}

func NewEvent(eventType string, data interface{}, source string) *Event {
    return &Event{
        Type:      eventType,
        Data:      data,
        Timestamp: time.Now(),
        Source:    source,
    }
}

func (e *Event) String() string {
    return fmt.Sprintf("Event[%s] from %s at %s: %v", 
        e.Type, e.Source, e.Timestamp.Format("15:04:05"), e.Data)
}

// EventHandler 事件处理器接口
type EventHandler interface {
    HandleEvent(event *Event)
    GetID() string
    GetName() string
}

// EventBus 事件总线
type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    if eb.handlers[eventType] == nil {
        eb.handlers[eventType] = make([]EventHandler, 0)
    }
    
    // 检查是否已经存在
    for _, h := range eb.handlers[eventType] {
        if h.GetID() == handler.GetID() {
            return
        }
    }
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
    fmt.Printf("Handler %s subscribed to event type: %s\n", handler.GetName(), eventType)
}

func (eb *EventBus) Unsubscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    if handlers, exists := eb.handlers[eventType]; exists {
        for i, h := range handlers {
            if h.GetID() == handler.GetID() {
                eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
                fmt.Printf("Handler %s unsubscribed from event type: %s\n", handler.GetName(), eventType)
                return
            }
        }
    }
}

func (eb *EventBus) Publish(event *Event) {
    eb.mu.RLock()
    handlers := make([]EventHandler, 0)
    if h, exists := eb.handlers[event.Type]; exists {
        handlers = make([]EventHandler, len(h))
        copy(handlers, h)
    }
    eb.mu.RUnlock()
    
    fmt.Printf("Publishing event: %s\n", event)
    for _, handler := range handlers {
        go func(h EventHandler) {
            h.HandleEvent(event)
        }(handler)
    }
}

func (eb *EventBus) GetHandlerCount(eventType string) int {
    eb.mu.RLock()
    defer eb.mu.RUnlock()
    
    if handlers, exists := eb.handlers[eventType]; exists {
        return len(handlers)
    }
    return 0
}

// LoggingHandler 日志处理器
type LoggingHandler struct {
    id   string
    name string
}

func NewLoggingHandler(id, name string) *LoggingHandler {
    return &LoggingHandler{
        id:   id,
        name: name,
    }
}

func (l *LoggingHandler) HandleEvent(event *Event) {
    fmt.Printf("[%s] %s: %s\n", l.name, event.Type, event.Data)
}

func (l *LoggingHandler) GetID() string {
    return l.id
}

func (l *LoggingHandler) GetName() string {
    return l.name
}

// EmailHandler 邮件处理器
type EmailHandler struct {
    id   string
    name string
}

func NewEmailHandler(id, name string) *EmailHandler {
    return &EmailHandler{
        id:   id,
        name: name,
    }
}

func (e *EmailHandler) HandleEvent(event *Event) {
    if event.Type == "user_registered" {
        userData := event.Data.(map[string]interface{})
        email := userData["email"].(string)
        fmt.Printf("[%s] Sending welcome email to: %s\n", e.name, email)
    }
}

func (e *EmailHandler) GetID() string {
    return e.id
}

func (e *EmailHandler) GetName() string {
    return e.name
}

// NotificationHandler 通知处理器
type NotificationHandler struct {
    id   string
    name string
}

func NewNotificationHandler(id, name string) *NotificationHandler {
    return &NotificationHandler{
        id:   id,
        name: name,
    }
}

func (n *NotificationHandler) HandleEvent(event *Event) {
    if event.Type == "order_created" {
        orderData := event.Data.(map[string]interface{})
        orderID := orderData["order_id"].(string)
        fmt.Printf("[%s] Sending push notification for order: %s\n", n.name, orderID)
    }
}

func (n *NotificationHandler) GetID() string {
    return n.id
}

func (n *NotificationHandler) GetName() string {
    return n.name
}

```

### 3.3.1.4.3 股票价格观察者模式

```go
package stockobserver

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// StockPrice 股票价格
type StockPrice struct {
    Symbol    string
    Price     float64
    Change    float64
    ChangePct float64
    Volume    int64
    Timestamp time.Time
}

func NewStockPrice(symbol string, price float64) *StockPrice {
    return &StockPrice{
        Symbol:    symbol,
        Price:     price,
        Timestamp: time.Now(),
    }
}

func (s *StockPrice) UpdatePrice(newPrice float64) {
    s.Change = newPrice - s.Price
    s.ChangePct = (s.Change / s.Price) * 100
    s.Price = newPrice
    s.Timestamp = time.Now()
}

func (s *StockPrice) String() string {
    changeSymbol := "+"
    if s.Change < 0 {
        changeSymbol = ""
    }
    return fmt.Sprintf("%s: $%.2f %s%.2f (%.2f%%)", 
        s.Symbol, s.Price, changeSymbol, s.Change, s.ChangePct)
}

// StockSubject 股票主题
type StockSubject struct {
    symbol    string
    price     *StockPrice
    observers []StockObserver
    mu        sync.RWMutex
    stopChan  chan bool
}

func NewStockSubject(symbol string, initialPrice float64) *StockSubject {
    return &StockSubject{
        symbol:    symbol,
        price:     NewStockPrice(symbol, initialPrice),
        observers: make([]StockObserver, 0),
        stopChan:  make(chan bool),
    }
}

func (s *StockSubject) Attach(observer StockObserver) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for _, obs := range s.observers {
        if obs.GetID() == observer.GetID() {
            return
        }
    }
    
    s.observers = append(s.observers, observer)
    fmt.Printf("Observer %s attached to %s\n", observer.GetName(), s.symbol)
}

func (s *StockSubject) Detach(observer StockObserver) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, obs := range s.observers {
        if obs.GetID() == observer.GetID() {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            fmt.Printf("Observer %s detached from %s\n", observer.GetName(), s.symbol)
            return
        }
    }
}

func (s *StockSubject) Notify() {
    s.mu.RLock()
    observers := make([]StockObserver, len(s.observers))
    copy(observers, s.observers)
    s.mu.RUnlock()
    
    for _, observer := range observers {
        observer.Update(s.price)
    }
}

func (s *StockSubject) UpdatePrice(newPrice float64) {
    s.price.UpdatePrice(newPrice)
    fmt.Printf("Stock %s price updated: %s\n", s.symbol, s.price)
    s.Notify()
}

func (s *StockSubject) StartSimulation() {
    go func() {
        ticker := time.NewTicker(2 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                // 模拟价格变化
                change := (rand.Float64() - 0.5) * 10 // -5 到 +5 的变化
                newPrice := s.price.Price + change
                if newPrice < 0 {
                    newPrice = 0.01
                }
                s.UpdatePrice(newPrice)
            case <-s.stopChan:
                return
            }
        }
    }()
}

func (s *StockSubject) StopSimulation() {
    close(s.stopChan)
}

func (s *StockSubject) GetCurrentPrice() *StockPrice {
    return s.price
}

// StockObserver 股票观察者接口
type StockObserver interface {
    Update(price *StockPrice)
    GetID() string
    GetName() string
}

// PriceAlertObserver 价格提醒观察者
type PriceAlertObserver struct {
    id       string
    name     string
    symbol   string
    minPrice float64
    maxPrice float64
}

func NewPriceAlertObserver(id, name, symbol string, minPrice, maxPrice float64) *PriceAlertObserver {
    return &PriceAlertObserver{
        id:       id,
        name:     name,
        symbol:   symbol,
        minPrice: minPrice,
        maxPrice: maxPrice,
    }
}

func (p *PriceAlertObserver) Update(price *StockPrice) {
    if price.Symbol != p.symbol {
        return
    }
    
    if price.Price <= p.minPrice {
        fmt.Printf("[%s] ALERT: %s price dropped to $%.2f (below minimum $%.2f)\n", 
            p.name, p.symbol, price.Price, p.minPrice)
    } else if price.Price >= p.maxPrice {
        fmt.Printf("[%s] ALERT: %s price rose to $%.2f (above maximum $%.2f)\n", 
            p.name, p.symbol, price.Price, p.maxPrice)
    }
}

func (p *PriceAlertObserver) GetID() string {
    return p.id
}

func (p *PriceAlertObserver) GetName() string {
    return p.name
}

// PortfolioObserver 投资组合观察者
type PortfolioObserver struct {
    id        string
    name      string
    holdings  map[string]int
    totalValue float64
}

func NewPortfolioObserver(id, name string) *PortfolioObserver {
    return &PortfolioObserver{
        id:       id,
        name:     name,
        holdings: make(map[string]int),
    }
}

func (p *PortfolioObserver) AddHolding(symbol string, shares int) {
    p.holdings[symbol] = shares
}

func (p *PortfolioObserver) Update(price *StockPrice) {
    if shares, exists := p.holdings[price.Symbol]; exists {
        holdingValue := float64(shares) * price.Price
        fmt.Printf("[%s] Portfolio update: %s holding value = $%.2f (%d shares @ $%.2f)\n", 
            p.name, price.Symbol, holdingValue, shares, price.Price)
    }
}

func (p *PortfolioObserver) GetID() string {
    return p.id
}

func (p *PortfolioObserver) GetName() string {
    return p.name
}

// ChartObserver 图表观察者
type ChartObserver struct {
    id     string
    name   string
    symbol string
    prices []float64
    maxLen int
}

func NewChartObserver(id, name, symbol string, maxLen int) *ChartObserver {
    return &ChartObserver{
        id:     id,
        name:   name,
        symbol: symbol,
        prices: make([]float64, 0),
        maxLen: maxLen,
    }
}

func (c *ChartObserver) Update(price *StockPrice) {
    if price.Symbol != c.symbol {
        return
    }
    
    c.prices = append(c.prices, price.Price)
    if len(c.prices) > c.maxLen {
        c.prices = c.prices[1:]
    }
    
    fmt.Printf("[%s] Chart update for %s: %d data points, latest: $%.2f\n", 
        c.name, c.symbol, len(c.prices), price.Price)
}

func (c *ChartObserver) GetID() string {
    return c.id
}

func (c *ChartObserver) GetName() string {
    return c.name
}

```

## 3.3.1.5 4. 工程案例

### 3.3.1.5.1 配置管理观察者模式

```go
package configobserver

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "sync"
    "time"
)

// ConfigChange 配置变更
type ConfigChange struct {
    Key       string
    OldValue  interface{}
    NewValue  interface{}
    Timestamp time.Time
    Source    string
}

func NewConfigChange(key string, oldValue, newValue interface{}, source string) *ConfigChange {
    return &ConfigChange{
        Key:       key,
        OldValue:  oldValue,
        NewValue:  newValue,
        Timestamp: time.Now(),
        Source:    source,
    }
}

func (c *ConfigChange) String() string {
    return fmt.Sprintf("Config change: %s = %v -> %v (from %s)", 
        c.Key, c.OldValue, c.NewValue, c.Source)
}

// ConfigObserver 配置观察者接口
type ConfigObserver interface {
    OnConfigChanged(change *ConfigChange)
    GetID() string
    GetName() string
}

// ConfigManager 配置管理器
type ConfigManager struct {
    config   map[string]interface{}
    observers []ConfigObserver
    mu       sync.RWMutex
    filePath string
}

func NewConfigManager(filePath string) *ConfigManager {
    cm := &ConfigManager{
        config:    make(map[string]interface{}),
        observers: make([]ConfigObserver, 0),
        filePath:  filePath,
    }
    
    // 加载初始配置
    cm.loadConfig()
    
    return cm
}

func (cm *ConfigManager) loadConfig() error {
    data, err := ioutil.ReadFile(cm.filePath)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config map[string]interface{}
    if err := json.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse config file: %w", err)
    }
    
    cm.mu.Lock()
    cm.config = config
    cm.mu.Unlock()
    
    return nil
}

func (cm *ConfigManager) saveConfig() error {
    cm.mu.RLock()
    data, err := json.MarshalIndent(cm.config, "", "  ")
    cm.mu.RUnlock()
    
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }
    
    if err := ioutil.WriteFile(cm.filePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }
    
    return nil
}

func (cm *ConfigManager) Get(key string) (interface{}, bool) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    value, exists := cm.config[key]
    return value, exists
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.mu.Lock()
    oldValue, exists := cm.config[key]
    cm.config[key] = value
    cm.mu.Unlock()
    
    if !exists {
        oldValue = nil
    }
    
    change := NewConfigChange(key, oldValue, value, "ConfigManager")
    cm.notifyObservers(change)
    
    // 保存到文件
    if err := cm.saveConfig(); err != nil {
        fmt.Printf("Failed to save config: %v\n", err)
    }
}

func (cm *ConfigManager) AddObserver(observer ConfigObserver) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    for _, obs := range cm.observers {
        if obs.GetID() == observer.GetID() {
            return
        }
    }
    
    cm.observers = append(cm.observers, observer)
    fmt.Printf("Config observer %s added\n", observer.GetName())
}

func (cm *ConfigManager) RemoveObserver(observer ConfigObserver) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    for i, obs := range cm.observers {
        if obs.GetID() == observer.GetID() {
            cm.observers = append(cm.observers[:i], cm.observers[i+1:]...)
            fmt.Printf("Config observer %s removed\n", observer.GetName())
            return
        }
    }
}

func (cm *ConfigManager) notifyObservers(change *ConfigChange) {
    cm.mu.RLock()
    observers := make([]ConfigObserver, len(cm.observers))
    copy(observers, cm.observers)
    cm.mu.RUnlock()
    
    fmt.Printf("Notifying %d observers of config change: %s\n", len(observers), change)
    for _, observer := range observers {
        go observer.OnConfigChanged(change)
    }
}

// LoggingConfigObserver 日志配置观察者
type LoggingConfigObserver struct {
    id   string
    name string
}

func NewLoggingConfigObserver(id, name string) *LoggingConfigObserver {
    return &LoggingConfigObserver{
        id:   id,
        name: name,
    }
}

func (l *LoggingConfigObserver) OnConfigChanged(change *ConfigChange) {
    if change.Key == "log_level" {
        fmt.Printf("[%s] Log level changed from %v to %v\n", 
            l.name, change.OldValue, change.NewValue)
    } else if change.Key == "log_file" {
        fmt.Printf("[%s] Log file changed from %v to %v\n", 
            l.name, change.OldValue, change.NewValue)
    }
}

func (l *LoggingConfigObserver) GetID() string {
    return l.id
}

func (l *LoggingConfigObserver) GetName() string {
    return l.name
}

// DatabaseConfigObserver 数据库配置观察者
type DatabaseConfigObserver struct {
    id   string
    name string
}

func NewDatabaseConfigObserver(id, name string) *DatabaseConfigObserver {
    return &DatabaseConfigObserver{
        id:   id,
        name: name,
    }
}

func (d *DatabaseConfigObserver) OnConfigChanged(change *ConfigChange) {
    if change.Key == "database_url" {
        fmt.Printf("[%s] Database URL changed, reconnecting...\n", d.name)
        // 这里可以实现数据库重连逻辑
    } else if change.Key == "database_pool_size" {
        fmt.Printf("[%s] Database pool size changed from %v to %v\n", 
            d.name, change.OldValue, change.NewValue)
    }
}

func (d *DatabaseConfigObserver) GetID() string {
    return d.id
}

func (d *DatabaseConfigObserver) GetName() string {
    return d.name
}

// CacheConfigObserver 缓存配置观察者
type CacheConfigObserver struct {
    id   string
    name string
}

func NewCacheConfigObserver(id, name string) *CacheConfigObserver {
    return &CacheConfigObserver{
        id:   id,
        name: name,
    }
}

func (c *CacheConfigObserver) OnConfigChanged(change *ConfigChange) {
    if change.Key == "cache_ttl" {
        fmt.Printf("[%s] Cache TTL changed from %v to %v seconds\n", 
            c.name, change.OldValue, change.NewValue)
    } else if change.Key == "cache_max_size" {
        fmt.Printf("[%s] Cache max size changed from %v to %v\n", 
            c.name, change.OldValue, change.NewValue)
    }
}

func (c *CacheConfigObserver) GetID() string {
    return c.id
}

func (c *CacheConfigObserver) GetName() string {
    return c.name
}

```

### 3.3.1.5.2 用户状态观察者模式

```go
package userobserver

import (
    "fmt"
    "sync"
    "time"
)

// UserEvent 用户事件
type UserEvent struct {
    UserID    string
    EventType string
    Data      map[string]interface{}
    Timestamp time.Time
}

func NewUserEvent(userID, eventType string, data map[string]interface{}) *UserEvent {
    return &UserEvent{
        UserID:    userID,
        EventType: eventType,
        Data:      data,
        Timestamp: time.Now(),
    }
}

func (u *UserEvent) String() string {
    return fmt.Sprintf("UserEvent[%s] for user %s: %v", 
        u.EventType, u.UserID, u.Data)
}

// UserObserver 用户观察者接口
type UserObserver interface {
    OnUserEvent(event *UserEvent)
    GetID() string
    GetName() string
}

// UserManager 用户管理器
type UserManager struct {
    users     map[string]*User
    observers []UserObserver
    mu        sync.RWMutex
}

type User struct {
    ID       string
    Username string
    Email    string
    Status   string
    LastSeen time.Time
}

func NewUserManager() *UserManager {
    return &UserManager{
        users:     make(map[string]*User),
        observers: make([]UserObserver, 0),
    }
}

func (um *UserManager) CreateUser(id, username, email string) {
    user := &User{
        ID:       id,
        Username: username,
        Email:    email,
        Status:   "active",
        LastSeen: time.Now(),
    }
    
    um.mu.Lock()
    um.users[id] = user
    um.mu.Unlock()
    
    event := NewUserEvent(id, "user_created", map[string]interface{}{
        "username": username,
        "email":    email,
    })
    
    um.notifyObservers(event)
}

func (um *UserManager) UpdateUserStatus(userID, status string) {
    um.mu.Lock()
    user, exists := um.users[userID]
    if exists {
        oldStatus := user.Status
        user.Status = status
        user.LastSeen = time.Now()
        um.mu.Unlock()
        
        event := NewUserEvent(userID, "status_changed", map[string]interface{}{
            "old_status": oldStatus,
            "new_status": status,
        })
        
        um.notifyObservers(event)
    } else {
        um.mu.Unlock()
    }
}

func (um *UserManager) UserLogin(userID string) {
    um.mu.Lock()
    user, exists := um.users[userID]
    if exists {
        user.LastSeen = time.Now()
        um.mu.Unlock()
        
        event := NewUserEvent(userID, "user_login", map[string]interface{}{
            "login_time": user.LastSeen,
        })
        
        um.notifyObservers(event)
    } else {
        um.mu.Unlock()
    }
}

func (um *UserManager) UserLogout(userID string) {
    um.mu.Lock()
    user, exists := um.users[userID]
    if exists {
        user.LastSeen = time.Now()
        um.mu.Unlock()
        
        event := NewUserEvent(userID, "user_logout", map[string]interface{}{
            "logout_time": user.LastSeen,
        })
        
        um.notifyObservers(event)
    } else {
        um.mu.Unlock()
    }
}

func (um *UserManager) AddObserver(observer UserObserver) {
    um.mu.Lock()
    defer um.mu.Unlock()
    
    for _, obs := range um.observers {
        if obs.GetID() == observer.GetID() {
            return
        }
    }
    
    um.observers = append(um.observers, observer)
    fmt.Printf("User observer %s added\n", observer.GetName())
}

func (um *UserManager) RemoveObserver(observer UserObserver) {
    um.mu.Lock()
    defer um.mu.Unlock()
    
    for i, obs := range um.observers {
        if obs.GetID() == observer.GetID() {
            um.observers = append(um.observers[:i], um.observers[i+1:]...)
            fmt.Printf("User observer %s removed\n", observer.GetName())
            return
        }
    }
}

func (um *UserManager) notifyObservers(event *UserEvent) {
    um.mu.RLock()
    observers := make([]UserObserver, len(um.observers))
    copy(observers, um.observers)
    um.mu.RUnlock()
    
    fmt.Printf("Notifying %d observers of user event: %s\n", len(observers), event)
    for _, observer := range observers {
        go observer.OnUserEvent(event)
    }
}

// EmailNotificationObserver 邮件通知观察者
type EmailNotificationObserver struct {
    id   string
    name string
}

func NewEmailNotificationObserver(id, name string) *EmailNotificationObserver {
    return &EmailNotificationObserver{
        id:   id,
        name: name,
    }
}

func (e *EmailNotificationObserver) OnUserEvent(event *UserEvent) {
    switch event.EventType {
    case "user_created":
        email := event.Data["email"].(string)
        fmt.Printf("[%s] Sending welcome email to: %s\n", e.name, email)
    case "status_changed":
        newStatus := event.Data["new_status"].(string)
        if newStatus == "suspended" {
            fmt.Printf("[%s] Sending suspension notification to user: %s\n", e.name, event.UserID)
        }
    }
}

func (e *EmailNotificationObserver) GetID() string {
    return e.id
}

func (e *EmailNotificationObserver) GetName() string {
    return e.name
}

// AnalyticsObserver 分析观察者
type AnalyticsObserver struct {
    id   string
    name string
}

func NewAnalyticsObserver(id, name string) *AnalyticsObserver {
    return &AnalyticsObserver{
        id:   id,
        name: name,
    }
}

func (a *AnalyticsObserver) OnUserEvent(event *UserEvent) {
    fmt.Printf("[%s] Recording analytics event: %s for user %s\n", 
        a.name, event.EventType, event.UserID)
    
    // 这里可以实现分析数据记录逻辑
    switch event.EventType {
    case "user_login":
        fmt.Printf("[%s] User login analytics recorded\n", a.name)
    case "user_logout":
        fmt.Printf("[%s] User logout analytics recorded\n", a.name)
    case "status_changed":
        fmt.Printf("[%s] Status change analytics recorded\n", a.name)
    }
}

func (a *AnalyticsObserver) GetID() string {
    return a.id
}

func (a *AnalyticsObserver) GetName() string {
    return a.name
}

// SecurityObserver 安全观察者
type SecurityObserver struct {
    id   string
    name string
}

func NewSecurityObserver(id, name string) *SecurityObserver {
    return &SecurityObserver{
        id:   id,
        name: name,
    }
}

func (s *SecurityObserver) OnUserEvent(event *UserEvent) {
    switch event.EventType {
    case "user_login":
        fmt.Printf("[%s] Security check for user login: %s\n", s.name, event.UserID)
        // 这里可以实现安全检查逻辑
    case "status_changed":
        newStatus := event.Data["new_status"].(string)
        if newStatus == "suspended" {
            fmt.Printf("[%s] Security alert: User %s suspended\n", s.name, event.UserID)
        }
    }
}

func (s *SecurityObserver) GetID() string {
    return s.id
}

func (s *SecurityObserver) GetName() string {
    return s.name
}

```

## 3.3.1.6 5. 批判性分析

### 3.3.1.6.1 优势

1. **松耦合**: 主题和观察者之间松耦合
2. **一对多关系**: 支持一个主题多个观察者
3. **动态订阅**: 观察者可以动态订阅和取消订阅
4. **自动通知**: 状态变化时自动通知观察者

### 3.3.1.6.2 劣势

1. **内存泄漏**: 观察者可能无法正确释放
2. **通知顺序**: 观察者通知顺序不确定
3. **循环依赖**: 可能产生循环依赖
4. **性能问题**: 大量观察者可能影响性能

### 3.3.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 切片 | 高 | 中 |
| Java | 接口 + 列表 | 中 | 中 |
| C++ | 虚函数 | 中 | 中 |
| Python | 回调函数 | 高 | 低 |

### 3.3.1.6.4 最新趋势

1. **响应式编程**: RxGo、响应式流
2. **事件驱动**: 事件总线、消息队列
3. **异步观察者**: 异步通知机制
4. **微服务观察者**: 服务间事件通知

## 3.3.1.7 6. 面试题与考点

### 3.3.1.7.1 基础考点

1. **Q**: 观察者模式与发布订阅模式的区别？
   **A**: 观察者直接通信，发布订阅通过中介

2. **Q**: 什么时候使用观察者模式？
   **A**: 需要一对多通知、状态变化通知时

3. **Q**: 观察者模式的优缺点？
   **A**: 优点：松耦合、动态订阅；缺点：内存泄漏、通知顺序

### 3.3.1.7.2 进阶考点

1. **Q**: 如何避免观察者模式的内存泄漏？
   **A**: 正确管理生命周期、使用弱引用

2. **Q**: 观察者模式在微服务中的应用？
   **A**: 事件驱动架构、服务间通信

3. **Q**: 如何处理观察者的异常？
   **A**: 异常隔离、重试机制、死信队列

## 3.3.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 观察者模式 | 一对多依赖关系的设计模式 | Observer Pattern |
| 主题 | 被观察的对象 | Subject |
| 观察者 | 观察主题的对象 | Observer |
| 通知 | 主题向观察者发送消息 | Notification |
| 订阅 | 观察者注册到主题 | Subscription |

## 3.3.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 内存泄漏 | 观察者无法释放 | 正确管理生命周期 |
| 通知顺序 | 观察者通知顺序不确定 | 定义通知顺序 |
| 循环依赖 | 主题和观察者相互依赖 | 避免循环引用 |
| 性能问题 | 大量观察者影响性能 | 异步通知、批量处理 |

## 3.3.1.10 9. 相关主题

- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)

## 3.3.1.11 10. 学习路径

### 3.3.1.11.1 新手路径

1. 理解观察者模式的基本概念
2. 学习主题和观察者的关系
3. 实现简单的观察者模式
4. 理解通知机制

### 3.3.1.11.2 进阶路径

1. 学习复杂的观察者实现
2. 理解观察者的性能优化
3. 掌握观察者的应用场景
4. 学习观察者的最佳实践

### 3.3.1.11.3 高阶路径

1. 分析观察者在大型项目中的应用
2. 理解观察者与架构设计的关系
3. 掌握观察者的性能调优
4. 学习观察者的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)
