# 中介者模式 (Mediator Pattern)

## 目录

1. [概述](#1-概述)
2. [理论基础](#2-理论基础)
3. [Go语言实现](#3-go语言实现)
4. [工程案例](#4-工程案例)
5. [批判性分析](#5-批判性分析)
6. [面试题与考点](#6-面试题与考点)
7. [术语表](#7-术语表)
8. [常见陷阱](#8-常见陷阱)
9. [相关主题](#9-相关主题)
10. [学习路径](#10-学习路径)

## 1. 概述

### 1.1 定义

中介者模式用一个中介对象来封装一系列的对象交互，中介者使各对象不需要显式地相互引用，从而使其耦合松散，而且可以独立地改变它们之间的交互。

**形式化定义**:
$$Mediator = (Mediator, Colleague, ConcreteMediator, ConcreteColleague)$$

其中：

- $Mediator$ 是中介者接口
- $Colleague$ 是同事接口
- $ConcreteMediator$ 是具体中介者
- $ConcreteColleague$ 是具体同事

### 1.2 核心特征

- **集中控制**: 中介者集中控制对象交互
- **松耦合**: 对象间不直接相互引用
- **简化交互**: 简化复杂的对象交互
- **易于维护**: 交互逻辑集中管理

## 2. 理论基础

### 2.1 数学形式化

**定义 2.1** (中介者模式): 中介者模式是一个四元组 $M = (Med, Col, I, C)$

其中：

- $Med$ 是中介者集合
- $Col$ 是同事集合
- $I$ 是交互函数，$I: Med \times Col \times Col \rightarrow Result$
- $C$ 是协调函数，$C: Med \times Event \rightarrow Action$

**定理 2.1** (交互简化): 对于任意同事 $c_1, c_2 \in Col$，通过中介者 $m \in Med$ 的交互比直接交互更简单。

### 2.2 范畴论视角

在范畴论中，中介者模式可以表示为：

$$Mediator : Colleague \times Colleague \rightarrow Interaction$$

## 3. Go语言实现

### 3.1 基础中介者模式

```go
package mediator

import "fmt"

// Mediator 中介者接口
type Mediator interface {
    Send(message string, colleague Colleague)
    RegisterColleague(colleague Colleague)
}

// Colleague 同事接口
type Colleague interface {
    Send(message string)
    Receive(message string)
    GetName() string
    SetMediator(mediator Mediator)
}

// ConcreteMediator 具体中介者
type ConcreteMediator struct {
    colleagues map[string]Colleague
}

func NewConcreteMediator() *ConcreteMediator {
    return &ConcreteMediator{
        colleagues: make(map[string]Colleague),
    }
}

func (m *ConcreteMediator) Send(message string, sender Colleague) {
    fmt.Printf("Mediator: %s sends message: %s\n", sender.GetName(), message)
    
    for name, colleague := range m.colleagues {
        if colleague != sender {
            fmt.Printf("Mediator: Forwarding message to %s\n", name)
            colleague.Receive(message)
        }
    }
}

func (m *ConcreteMediator) RegisterColleague(colleague Colleague) {
    m.colleagues[colleague.GetName()] = colleague
    colleague.SetMediator(m)
    fmt.Printf("Mediator: Registered colleague %s\n", colleague.GetName())
}

// ConcreteColleagueA 具体同事A
type ConcreteColleagueA struct {
    name     string
    mediator Mediator
}

func NewConcreteColleagueA(name string) *ConcreteColleagueA {
    return &ConcreteColleagueA{
        name: name,
    }
}

func (c *ConcreteColleagueA) Send(message string) {
    fmt.Printf("%s: Sending message: %s\n", c.name, message)
    c.mediator.Send(message, c)
}

func (c *ConcreteColleagueA) Receive(message string) {
    fmt.Printf("%s: Received message: %s\n", c.name, message)
}

func (c *ConcreteColleagueA) GetName() string {
    return c.name
}

func (c *ConcreteColleagueA) SetMediator(mediator Mediator) {
    c.mediator = mediator
}

// ConcreteColleagueB 具体同事B
type ConcreteColleagueB struct {
    name     string
    mediator Mediator
}

func NewConcreteColleagueB(name string) *ConcreteColleagueB {
    return &ConcreteColleagueB{
        name: name,
    }
}

func (c *ConcreteColleagueB) Send(message string) {
    fmt.Printf("%s: Sending message: %s\n", c.name, message)
    c.mediator.Send(message, c)
}

func (c *ConcreteColleagueB) Receive(message string) {
    fmt.Printf("%s: Received message: %s\n", c.name, message)
}

func (c *ConcreteColleagueB) GetName() string {
    return c.name
}

func (c *ConcreteColleagueB) SetMediator(mediator Mediator) {
    c.mediator = mediator
}
```

### 3.2 聊天室中介者模式

```go
package chatroom

import (
    "fmt"
    "time"
)

// Message 消息
type Message struct {
    From      string
    Content   string
    Timestamp time.Time
    Type      string
}

func NewMessage(from, content, msgType string) *Message {
    return &Message{
        From:      from,
        Content:   content,
        Timestamp: time.Now(),
        Type:      msgType,
    }
}

func (m *Message) String() string {
    return fmt.Sprintf("[%s] %s: %s", 
        m.Timestamp.Format("15:04:05"), m.From, m.Content)
}

// ChatRoomMediator 聊天室中介者
type ChatRoomMediator struct {
    users    map[string]*User
    messages []*Message
}

func NewChatRoomMediator() *ChatRoomMediator {
    return &ChatRoomMediator{
        users:    make(map[string]*User),
        messages: make([]*Message, 0),
    }
}

func (c *ChatRoomMediator) SendMessage(message *Message) {
    c.messages = append(c.messages, message)
    
    switch message.Type {
    case "broadcast":
        c.broadcastMessage(message)
    case "private":
        c.sendPrivateMessage(message)
    case "system":
        c.sendSystemMessage(message)
    }
}

func (c *ChatRoomMediator) broadcastMessage(message *Message) {
    fmt.Printf("Broadcasting: %s\n", message)
    for _, user := range c.users {
        if user.GetName() != message.From {
            user.ReceiveMessage(message)
        }
    }
}

func (c *ChatRoomMediator) sendPrivateMessage(message *Message) {
    // 简化实现：假设私聊消息格式为 "@username: content"
    fmt.Printf("Private message: %s\n", message)
    // 这里可以实现私聊逻辑
}

func (c *ChatRoomMediator) sendSystemMessage(message *Message) {
    fmt.Printf("System message: %s\n", message)
    for _, user := range c.users {
        user.ReceiveMessage(message)
    }
}

func (c *ChatRoomMediator) AddUser(user *User) {
    c.users[user.GetName()] = user
    user.SetMediator(c)
    
    // 发送系统消息通知新用户加入
    joinMessage := NewMessage("System", 
        fmt.Sprintf("%s joined the chat room", user.GetName()), "system")
    c.SendMessage(joinMessage)
}

func (c *ChatRoomMediator) RemoveUser(user *User) {
    delete(c.users, user.GetName())
    
    // 发送系统消息通知用户离开
    leaveMessage := NewMessage("System", 
        fmt.Sprintf("%s left the chat room", user.GetName()), "system")
    c.SendMessage(leaveMessage)
}

func (c *ChatRoomMediator) GetUserCount() int {
    return len(c.users)
}

func (c *ChatRoomMediator) GetMessageHistory() []*Message {
    return c.messages
}

// User 用户
type User struct {
    name     string
    mediator *ChatRoomMediator
    online   bool
}

func NewUser(name string) *User {
    return &User{
        name:   name,
        online: true,
    }
}

func (u *User) SendMessage(content, msgType string) {
    if !u.online {
        fmt.Printf("%s is offline, cannot send message\n", u.name)
        return
    }
    
    message := NewMessage(u.name, content, msgType)
    u.mediator.SendMessage(message)
}

func (u *User) ReceiveMessage(message *Message) {
    if u.online {
        fmt.Printf("%s received: %s\n", u.name, message)
    }
}

func (u *User) GetName() string {
    return u.name
}

func (u *User) SetMediator(mediator *ChatRoomMediator) {
    u.mediator = mediator
}

func (u *User) SetOnline(online bool) {
    u.online = online
    if !online {
        fmt.Printf("%s went offline\n", u.name)
    } else {
        fmt.Printf("%s came online\n", u.name)
    }
}

func (u *User) IsOnline() bool {
    return u.online
}
```

### 3.3 航空交通管制中介者模式

```go
package airtrafficcontrol

import (
    "fmt"
    "time"
)

// Aircraft 飞机
type Aircraft struct {
    id       string
    callsign string
    altitude int
    speed    int
    position Position
    mediator *AirTrafficController
}

type Position struct {
    latitude  float64
    longitude float64
}

func NewAircraft(id, callsign string) *Aircraft {
    return &Aircraft{
        id:       id,
        callsign: callsign,
        altitude: 0,
        speed:    0,
        position: Position{0, 0},
    }
}

func (a *Aircraft) RequestTakeoff() {
    fmt.Printf("%s requesting takeoff\n", a.callsign)
    a.mediator.HandleTakeoffRequest(a)
}

func (a *Aircraft) RequestLanding() {
    fmt.Printf("%s requesting landing\n", a.callsign)
    a.mediator.HandleLandingRequest(a)
}

func (a *Aircraft) RequestAltitudeChange(newAltitude int) {
    fmt.Printf("%s requesting altitude change to %d feet\n", a.callsign, newAltitude)
    a.mediator.HandleAltitudeChange(a, newAltitude)
}

func (a *Aircraft) UpdatePosition(lat, lng float64) {
    a.position.latitude = lat
    a.position.longitude = lng
    a.mediator.HandlePositionUpdate(a)
}

func (a *Aircraft) ReceiveClearance(clearance string) {
    fmt.Printf("%s received clearance: %s\n", a.callsign, clearance)
}

func (a *Aircraft) SetMediator(mediator *AirTrafficController) {
    a.mediator = mediator
}

func (a *Aircraft) GetID() string {
    return a.id
}

func (a *Aircraft) GetCallsign() string {
    return a.callsign
}

func (a *Aircraft) GetAltitude() int {
    return a.altitude
}

func (a *Aircraft) GetPosition() Position {
    return a.position
}

// AirTrafficController 空中交通管制员
type AirTrafficController struct {
    aircraft map[string]*Aircraft
    runway   *Runway
    airspace *Airspace
}

func NewAirTrafficController() *AirTrafficController {
    return &AirTrafficController{
        aircraft: make(map[string]*Aircraft),
        runway:   NewRunway(),
        airspace: NewAirspace(),
    }
}

func (a *AirTrafficController) RegisterAircraft(aircraft *Aircraft) {
    a.aircraft[aircraft.GetID()] = aircraft
    aircraft.SetMediator(a)
    fmt.Printf("ATC: Registered aircraft %s\n", aircraft.GetCallsign())
}

func (a *AirTrafficController) HandleTakeoffRequest(aircraft *Aircraft) {
    if a.runway.IsAvailable() {
        a.runway.Occupy(aircraft.GetCallsign())
        aircraft.ReceiveClearance("Cleared for takeoff")
        fmt.Printf("ATC: %s cleared for takeoff\n", aircraft.GetCallsign())
    } else {
        aircraft.ReceiveClearance("Hold position, runway occupied")
        fmt.Printf("ATC: %s hold position\n", aircraft.GetCallsign())
    }
}

func (a *AirTrafficController) HandleLandingRequest(aircraft *Aircraft) {
    if a.runway.IsAvailable() {
        a.runway.Occupy(aircraft.GetCallsign())
        aircraft.ReceiveClearance("Cleared to land")
        fmt.Printf("ATC: %s cleared to land\n", aircraft.GetCallsign())
    } else {
        aircraft.ReceiveClearance("Go around, runway occupied")
        fmt.Printf("ATC: %s go around\n", aircraft.GetCallsign())
    }
}

func (a *AirTrafficController) HandleAltitudeChange(aircraft *Aircraft, newAltitude int) {
    if a.airspace.IsAltitudeAvailable(newAltitude, aircraft.GetPosition()) {
        aircraft.altitude = newAltitude
        aircraft.ReceiveClearance(fmt.Sprintf("Cleared to %d feet", newAltitude))
        fmt.Printf("ATC: %s cleared to %d feet\n", aircraft.GetCallsign(), newAltitude)
    } else {
        aircraft.ReceiveClearance("Altitude not available, maintain current altitude")
        fmt.Printf("ATC: %s maintain current altitude\n", aircraft.GetCallsign())
    }
}

func (a *AirTrafficController) HandlePositionUpdate(aircraft *Aircraft) {
    // 检查与其他飞机的冲突
    for _, other := range a.aircraft {
        if other != aircraft && a.isConflict(aircraft, other) {
            aircraft.ReceiveClearance("Traffic alert, adjust course")
            fmt.Printf("ATC: Traffic alert for %s\n", aircraft.GetCallsign())
            return
        }
    }
}

func (a *AirTrafficController) isConflict(aircraft1, aircraft2 *Aircraft) bool {
    // 简化的冲突检测
    pos1 := aircraft1.GetPosition()
    pos2 := aircraft2.GetPosition()
    
    distance := a.calculateDistance(pos1, pos2)
    altitudeDiff := abs(aircraft1.GetAltitude() - aircraft2.GetAltitude())
    
    return distance < 5.0 && altitudeDiff < 1000 // 5海里内且高度差小于1000英尺
}

func (a *AirTrafficController) calculateDistance(pos1, pos2 Position) float64 {
    // 简化的距离计算
    latDiff := pos1.latitude - pos2.latitude
    lngDiff := pos1.longitude - pos2.longitude
    return latDiff*latDiff + lngDiff*lngDiff
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

// Runway 跑道
type Runway struct {
    available bool
    occupiedBy string
}

func NewRunway() *Runway {
    return &Runway{
        available: true,
    }
}

func (r *Runway) IsAvailable() bool {
    return r.available
}

func (r *Runway) Occupy(aircraft string) {
    r.available = false
    r.occupiedBy = aircraft
    fmt.Printf("Runway occupied by %s\n", aircraft)
}

func (r *Runway) Release() {
    r.available = true
    r.occupiedBy = ""
    fmt.Println("Runway released")
}

// Airspace 空域
type Airspace struct {
    maxAltitude int
    minAltitude int
}

func NewAirspace() *Airspace {
    return &Airspace{
        maxAltitude: 45000,
        minAltitude: 0,
    }
}

func (a *Airspace) IsAltitudeAvailable(altitude int, position Position) bool {
    return altitude >= a.minAltitude && altitude <= a.maxAltitude
}
```

## 4. 工程案例

### 4.1 事件总线中介者模式

```go
package eventbus

import (
    "fmt"
    "sync"
    "time"
)

// Event 事件
type Event struct {
    Type      string
    Data      interface{}
    Source    string
    Timestamp time.Time
}

func NewEvent(eventType string, data interface{}, source string) *Event {
    return &Event{
        Type:      eventType,
        Data:      data,
        Source:    source,
        Timestamp: time.Now(),
    }
}

func (e *Event) String() string {
    return fmt.Sprintf("Event[%s] from %s at %s: %v", 
        e.Type, e.Source, e.Timestamp.Format("15:04:05"), e.Data)
}

// EventBus 事件总线中介者
type EventBus struct {
    subscribers map[string][]Subscriber
    mu          sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]Subscriber),
    }
}

func (e *EventBus) Subscribe(eventType string, subscriber Subscriber) {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if e.subscribers[eventType] == nil {
        e.subscribers[eventType] = make([]Subscriber, 0)
    }
    
    e.subscribers[eventType] = append(e.subscribers[eventType], subscriber)
    fmt.Printf("EventBus: %s subscribed to %s events\n", 
        subscriber.GetName(), eventType)
}

func (e *EventBus) Unsubscribe(eventType string, subscriber Subscriber) {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if subscribers, exists := e.subscribers[eventType]; exists {
        for i, sub := range subscribers {
            if sub == subscriber {
                e.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
                fmt.Printf("EventBus: %s unsubscribed from %s events\n", 
                    subscriber.GetName(), eventType)
                return
            }
        }
    }
}

func (e *EventBus) Publish(event *Event) {
    e.mu.RLock()
    subscribers := make([]Subscriber, 0)
    if subs, exists := e.subscribers[event.Type]; exists {
        subscribers = make([]Subscriber, len(subs))
        copy(subscribers, subs)
    }
    e.mu.RUnlock()
    
    fmt.Printf("EventBus: Publishing event %s\n", event)
    for _, subscriber := range subscribers {
        go subscriber.HandleEvent(event)
    }
}

// Subscriber 订阅者接口
type Subscriber interface {
    HandleEvent(event *Event)
    GetName() string
}

// Publisher 发布者接口
type Publisher interface {
    PublishEvent(eventType string, data interface{})
    SetEventBus(eventBus *EventBus)
}

// ConcretePublisher 具体发布者
type ConcretePublisher struct {
    name     string
    eventBus *EventBus
}

func NewConcretePublisher(name string) *ConcretePublisher {
    return &ConcretePublisher{
        name: name,
    }
}

func (c *ConcretePublisher) PublishEvent(eventType string, data interface{}) {
    event := NewEvent(eventType, data, c.name)
    c.eventBus.Publish(event)
}

func (c *ConcretePublisher) SetEventBus(eventBus *EventBus) {
    c.eventBus = eventBus
}

// ConcreteSubscriber 具体订阅者
type ConcreteSubscriber struct {
    name string
}

func NewConcreteSubscriber(name string) *ConcreteSubscriber {
    return &ConcreteSubscriber{
        name: name,
    }
}

func (c *ConcreteSubscriber) HandleEvent(event *Event) {
    fmt.Printf("%s handling event: %s\n", c.name, event)
}

func (c *ConcreteSubscriber) GetName() string {
    return c.name
}
```

## 5. 批判性分析

### 5.1 优势

1. **集中控制**: 中介者集中控制对象交互
2. **松耦合**: 对象间不直接相互引用
3. **简化交互**: 简化复杂的对象交互
4. **易于维护**: 交互逻辑集中管理

### 5.2 劣势

1. **中介者复杂**: 中介者可能变得复杂
2. **性能问题**: 中介者可能成为性能瓶颈
3. **单点故障**: 中介者故障影响整个系统
4. **扩展困难**: 添加新对象类型困难

### 5.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口+通道 | 高 | 中 |
| Java | 接口 | 中 | 中 |
| C++ | 虚函数 | 高 | 中 |
| Python | 对象 | 中 | 低 |

### 5.4 最新趋势

1. **事件驱动**: 使用事件驱动架构
2. **消息队列**: 使用消息队列作为中介者
3. **微服务**: 使用API网关作为中介者
4. **响应式**: 使用响应式编程

## 6. 面试题与考点

### 6.1 基础考点

1. **Q**: 中介者模式与观察者模式的区别？
   **A**: 中介者关注对象交互，观察者关注状态变化

2. **Q**: 什么时候使用中介者模式？
   **A**: 对象间交互复杂、需要集中控制时

3. **Q**: 中介者模式的优缺点？
   **A**: 优点：集中控制、松耦合；缺点：中介者复杂、单点故障

### 6.2 进阶考点

1. **Q**: 如何避免中介者的复杂性？
   **A**: 分层中介者、职责分离、事件驱动

2. **Q**: 中介者模式在微服务中的应用？
   **A**: API网关、服务网格、消息总线

3. **Q**: 如何处理中介者的性能问题？
   **A**: 异步处理、缓存、负载均衡

## 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 中介者模式 | 封装对象交互的设计模式 | Mediator Pattern |
| 中介者 | 协调对象交互的对象 | Mediator |
| 同事 | 相互交互的对象 | Colleague |
| 交互 | 对象间的通信 | Interaction |

## 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 中介者复杂 | 中介者承担过多职责 | 职责分离、分层设计 |
| 性能瓶颈 | 中介者成为性能瓶颈 | 异步处理、缓存 |
| 单点故障 | 中介者故障影响系统 | 冗余、故障转移 |
| 扩展困难 | 添加新对象类型困难 | 接口设计、插件架构 |

## 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)
- [访问者模式](./07-Visitor-Pattern.md)

## 10. 学习路径

### 10.1 新手路径

1. 理解中介者模式的基本概念
2. 学习中介者和同事的关系
3. 实现简单的中介者模式
4. 理解交互协调机制

### 10.2 进阶路径

1. 学习复杂的中介者实现
2. 理解中介者的性能优化
3. 掌握中介者的应用场景
4. 学习中介者的最佳实践

### 10.3 高阶路径

1. 分析中介者在大型项目中的应用
2. 理解中介者与架构设计的关系
3. 掌握中介者的性能调优
4. 学习中介者的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)
