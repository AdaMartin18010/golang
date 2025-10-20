# 3.3.1 状态模式 (State Pattern)

## 3.3.1.1 目录

## 3.3.1.2 1. 概述

### 3.3.1.2.1 定义

状态模式允许对象在内部状态改变时改变其行为，对象看起来好像修改了其类。

**形式化定义**:
$$State = (Context, State, ConcreteState_1, ConcreteState_2, ..., ConcreteState_n)$$

其中：

- $Context$ 是上下文类
- $State$ 是状态接口
- $ConcreteState_i$ 是具体状态实现

### 3.3.1.2.2 核心特征

- **状态封装**: 每个状态都被封装在独立的类中
- **状态转换**: 对象可以在不同状态间转换
- **行为变化**: 不同状态下对象行为不同
- **状态管理**: 集中管理状态转换逻辑

## 3.3.1.3 2. 理论基础

### 3.3.1.3.1 数学形式化

**定义 2.1** (状态模式): 状态模式是一个四元组 $S = (C, \Sigma, T, B)$

其中：

- $C$ 是上下文集合
- $\Sigma$ 是状态集合
- $T$ 是状态转换函数，$T: C \times \Sigma \times Event \rightarrow \Sigma$
- $B$ 是行为函数，$B: C \times \Sigma \times Action \rightarrow Result$

**定理 2.1** (状态转换一致性): 对于任意上下文 $c \in C$ 和状态 $\sigma \in \Sigma$，状态转换必须保持一致性。

### 3.3.1.3.2 范畴论视角

在范畴论中，状态模式可以表示为：

$$State : Context \times Event \rightarrow NewState$$

## 3.3.1.4 3. Go语言实现

### 3.3.1.4.1 基础状态模式

```go
package state

import "fmt"

// State 状态接口
type State interface {
    Handle(context *Context)
    GetName() string
}

// Context 上下文
type Context struct {
    state State
}

func NewContext() *Context {
    return &Context{
        state: NewStateA(),
    }
}

func (c *Context) SetState(state State) {
    c.state = state
    fmt.Printf("Context state changed to: %s\n", state.GetName())
}

func (c *Context) Request() {
    c.state.Handle(c)
}

func (c *Context) GetCurrentState() string {
    return c.state.GetName()
}

// StateA 状态A
type StateA struct{}

func NewStateA() *StateA {
    return &StateA{}
}

func (s *StateA) Handle(context *Context) {
    fmt.Println("StateA is handling the request")
    context.SetState(NewStateB())
}

func (s *StateA) GetName() string {
    return "StateA"
}

// StateB 状态B
type StateB struct{}

func NewStateB() *StateB {
    return &StateB{}
}

func (s *StateB) Handle(context *Context) {
    fmt.Println("StateB is handling the request")
    context.SetState(NewStateC())
}

func (s *StateB) GetName() string {
    return "StateB"
}

// StateC 状态C
type StateC struct{}

func NewStateC() *StateC {
    return &StateC{}
}

func (s *StateC) Handle(context *Context) {
    fmt.Println("StateC is handling the request")
    context.SetState(NewStateA())
}

func (s *StateC) GetName() string {
    return "StateC"
}

```

### 3.3.1.4.2 自动售货机状态模式

```go
package vendingmachine

import "fmt"

// VendingMachineState 自动售货机状态接口
type VendingMachineState interface {
    InsertCoin(machine *VendingMachine)
    EjectCoin(machine *VendingMachine)
    SelectProduct(machine *VendingMachine, product string)
    Dispense(machine *VendingMachine)
    GetName() string
}

// VendingMachine 自动售货机
type VendingMachine struct {
    state           VendingMachineState
    coins           int
    products        map[string]int
    selectedProduct string
}

func NewVendingMachine() *VendingMachine {
    vm := &VendingMachine{
        coins:    0,
        products: make(map[string]int),
    }
    vm.state = NewIdleState()
    return vm
}

func (v *VendingMachine) SetState(state VendingMachineState) {
    v.state = state
    fmt.Printf("Vending machine state changed to: %s\n", state.GetName())
}

func (v *VendingMachine) InsertCoin() {
    v.state.InsertCoin(v)
}

func (v *VendingMachine) EjectCoin() {
    v.state.EjectCoin(v)
}

func (v *VendingMachine) SelectProduct(product string) {
    v.state.SelectProduct(v, product)
}

func (v *VendingMachine) Dispense() {
    v.state.Dispense(v)
}

func (v *VendingMachine) AddProduct(name string, price int) {
    v.products[name] = price
    fmt.Printf("Added product: %s (price: %d coins)\n", name, price)
}

func (v *VendingMachine) GetCurrentState() string {
    return v.state.GetName()
}

// IdleState 空闲状态
type IdleState struct{}

func NewIdleState() *IdleState {
    return &IdleState{}
}

func (i *IdleState) InsertCoin(machine *VendingMachine) {
    machine.coins++
    fmt.Printf("Coin inserted. Total coins: %d\n", machine.coins)
    machine.SetState(NewHasCoinState())
}

func (i *IdleState) EjectCoin(machine *VendingMachine) {
    fmt.Println("No coin to eject")
}

func (i *IdleState) SelectProduct(machine *VendingMachine, product string) {
    fmt.Println("Please insert a coin first")
}

func (i *IdleState) Dispense(machine *VendingMachine) {
    fmt.Println("Please insert a coin first")
}

func (i *IdleState) GetName() string {
    return "Idle"
}

// HasCoinState 有硬币状态
type HasCoinState struct{}

func NewHasCoinState() *HasCoinState {
    return &HasCoinState{}
}

func (h *HasCoinState) InsertCoin(machine *VendingMachine) {
    machine.coins++
    fmt.Printf("Coin inserted. Total coins: %d\n", machine.coins)
}

func (h *HasCoinState) EjectCoin(machine *VendingMachine) {
    fmt.Printf("Ejecting %d coins\n", machine.coins)
    machine.coins = 0
    machine.SetState(NewIdleState())
}

func (h *HasCoinState) SelectProduct(machine *VendingMachine, product string) {
    price, exists := machine.products[product]
    if !exists {
        fmt.Printf("Product %s not available\n", product)
        return
    }
    
    if machine.coins >= price {
        machine.selectedProduct = product
        machine.coins -= price
        fmt.Printf("Product %s selected. Remaining coins: %d\n", product, machine.coins)
        machine.SetState(NewSoldState())
    } else {
        fmt.Printf("Insufficient coins. Need %d, have %d\n", price, machine.coins)
    }
}

func (h *HasCoinState) Dispense(machine *VendingMachine) {
    fmt.Println("Please select a product first")
}

func (h *HasCoinState) GetName() string {
    return "Has Coin"
}

// SoldState 售出状态
type SoldState struct{}

func NewSoldState() *SoldState {
    return &SoldState{}
}

func (s *SoldState) InsertCoin(machine *VendingMachine) {
    fmt.Println("Please wait, dispensing product")
}

func (s *SoldState) EjectCoin(machine *VendingMachine) {
    fmt.Println("Cannot eject coin while dispensing")
}

func (s *SoldState) SelectProduct(machine *VendingMachine, product string) {
    fmt.Println("Please wait, dispensing product")
}

func (s *SoldState) Dispense(machine *VendingMachine) {
    fmt.Printf("Dispensing %s\n", machine.selectedProduct)
    machine.selectedProduct = ""
    
    if machine.coins > 0 {
        machine.SetState(NewHasCoinState())
    } else {
        machine.SetState(NewIdleState())
    }
}

func (s *SoldState) GetName() string {
    return "Sold"
}

```

### 3.3.1.4.3 订单状态模式

```go
package orderstate

import (
    "fmt"
    "time"
)

// OrderState 订单状态接口
type OrderState interface {
    Confirm(order *Order)
    Ship(order *Order)
    Deliver(order *Order)
    Cancel(order *Order)
    GetName() string
}

// Order 订单
type Order struct {
    ID          string
    CustomerID  string
    Items       []string
    TotalAmount float64
    State       OrderState
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewOrder(id, customerID string, items []string, totalAmount float64) *Order {
    order := &Order{
        ID:          id,
        CustomerID:  customerID,
        Items:       items,
        TotalAmount: totalAmount,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    order.State = NewPendingState()
    return order
}

func (o *Order) SetState(state OrderState) {
    o.State = state
    o.UpdatedAt = time.Now()
    fmt.Printf("Order %s state changed to: %s\n", o.ID, state.GetName())
}

func (o *Order) Confirm() {
    o.State.Confirm(o)
}

func (o *Order) Ship() {
    o.State.Ship(o)
}

func (o *Order) Deliver() {
    o.State.Deliver(o)
}

func (o *Order) Cancel() {
    o.State.Cancel(o)
}

func (o *Order) GetCurrentState() string {
    return o.State.GetName()
}

// PendingState 待确认状态
type PendingState struct{}

func NewPendingState() *PendingState {
    return &PendingState{}
}

func (p *PendingState) Confirm(order *Order) {
    fmt.Printf("Order %s confirmed\n", order.ID)
    order.SetState(NewConfirmedState())
}

func (p *PendingState) Ship(order *Order) {
    fmt.Println("Cannot ship unconfirmed order")
}

func (p *PendingState) Deliver(order *Order) {
    fmt.Println("Cannot deliver unconfirmed order")
}

func (p *PendingState) Cancel(order *Order) {
    fmt.Printf("Order %s cancelled\n", order.ID)
    order.SetState(NewCancelledState())
}

func (p *PendingState) GetName() string {
    return "Pending"
}

// ConfirmedState 已确认状态
type ConfirmedState struct{}

func NewConfirmedState() *ConfirmedState {
    return &ConfirmedState{}
}

func (c *ConfirmedState) Confirm(order *Order) {
    fmt.Println("Order is already confirmed")
}

func (c *ConfirmedState) Ship(order *Order) {
    fmt.Printf("Order %s shipped\n", order.ID)
    order.SetState(NewShippedState())
}

func (c *ConfirmedState) Deliver(order *Order) {
    fmt.Println("Cannot deliver unshipped order")
}

func (c *ConfirmedState) Cancel(order *Order) {
    fmt.Printf("Order %s cancelled\n", order.ID)
    order.SetState(NewCancelledState())
}

func (c *ConfirmedState) GetName() string {
    return "Confirmed"
}

// ShippedState 已发货状态
type ShippedState struct{}

func NewShippedState() *ShippedState {
    return &ShippedState{}
}

func (s *ShippedState) Confirm(order *Order) {
    fmt.Println("Order is already confirmed")
}

func (s *ShippedState) Ship(order *Order) {
    fmt.Println("Order is already shipped")
}

func (s *ShippedState) Deliver(order *Order) {
    fmt.Printf("Order %s delivered\n", order.ID)
    order.SetState(NewDeliveredState())
}

func (s *ShippedState) Cancel(order *Order) {
    fmt.Println("Cannot cancel shipped order")
}

func (s *ShippedState) GetName() string {
    return "Shipped"
}

// DeliveredState 已送达状态
type DeliveredState struct{}

func NewDeliveredState() *DeliveredState {
    return &DeliveredState{}
}

func (d *DeliveredState) Confirm(order *Order) {
    fmt.Println("Order is already delivered")
}

func (d *DeliveredState) Ship(order *Order) {
    fmt.Println("Order is already delivered")
}

func (d *DeliveredState) Deliver(order *Order) {
    fmt.Println("Order is already delivered")
}

func (d *DeliveredState) Cancel(order *Order) {
    fmt.Println("Cannot cancel delivered order")
}

func (d *DeliveredState) GetName() string {
    return "Delivered"
}

// CancelledState 已取消状态
type CancelledState struct{}

func NewCancelledState() *CancelledState {
    return &CancelledState{}
}

func (c *CancelledState) Confirm(order *Order) {
    fmt.Println("Cannot confirm cancelled order")
}

func (c *CancelledState) Ship(order *Order) {
    fmt.Println("Cannot ship cancelled order")
}

func (c *CancelledState) Deliver(order *Order) {
    fmt.Println("Cannot deliver cancelled order")
}

func (c *CancelledState) Cancel(order *Order) {
    fmt.Println("Order is already cancelled")
}

func (c *CancelledState) GetName() string {
    return "Cancelled"
}

```

## 3.3.1.5 4. 工程案例

### 3.3.1.5.1 网络连接状态模式

```go
package connectionstate

import (
    "fmt"
    "time"
)

// ConnectionState 连接状态接口
type ConnectionState interface {
    Connect(conn *Connection)
    Disconnect(conn *Connection)
    Send(conn *Connection, data []byte)
    Receive(conn *Connection) []byte
    GetName() string
}

// Connection 网络连接
type Connection struct {
    ID       string
    State    ConnectionState
    Data     []byte
    IsActive bool
}

func NewConnection(id string) *Connection {
    conn := &Connection{
        ID:       id,
        IsActive: false,
    }
    conn.State = NewDisconnectedState()
    return conn
}

func (c *Connection) SetState(state ConnectionState) {
    c.State = state
    fmt.Printf("Connection %s state changed to: %s\n", c.ID, state.GetName())
}

func (c *Connection) Connect() {
    c.State.Connect(c)
}

func (c *Connection) Disconnect() {
    c.State.Disconnect(c)
}

func (c *Connection) Send(data []byte) {
    c.State.Send(c, data)
}

func (c *Connection) Receive() []byte {
    return c.State.Receive(c)
}

func (c *Connection) GetCurrentState() string {
    return c.State.GetName()
}

// DisconnectedState 断开连接状态
type DisconnectedState struct{}

func NewDisconnectedState() *DisconnectedState {
    return &DisconnectedState{}
}

func (d *DisconnectedState) Connect(conn *Connection) {
    fmt.Printf("Connecting to %s...\n", conn.ID)
    time.Sleep(100 * time.Millisecond) // 模拟连接延迟
    conn.IsActive = true
    conn.SetState(NewConnectedState())
}

func (d *DisconnectedState) Disconnect(conn *Connection) {
    fmt.Println("Already disconnected")
}

func (d *DisconnectedState) Send(conn *Connection, data []byte) {
    fmt.Println("Cannot send data: not connected")
}

func (d *DisconnectedState) Receive(conn *Connection) []byte {
    fmt.Println("Cannot receive data: not connected")
    return nil
}

func (d *DisconnectedState) GetName() string {
    return "Disconnected"
}

// ConnectedState 已连接状态
type ConnectedState struct{}

func NewConnectedState() *ConnectedState {
    return &ConnectedState{}
}

func (c *ConnectedState) Connect(conn *Connection) {
    fmt.Println("Already connected")
}

func (c *ConnectedState) Disconnect(conn *Connection) {
    fmt.Printf("Disconnecting from %s...\n", conn.ID)
    conn.IsActive = false
    conn.SetState(NewDisconnectedState())
}

func (c *ConnectedState) Send(conn *Connection, data []byte) {
    fmt.Printf("Sending %d bytes to %s\n", len(data), conn.ID)
    conn.Data = append(conn.Data, data...)
}

func (c *ConnectedState) Receive(conn *Connection) []byte {
    if len(conn.Data) > 0 {
        data := conn.Data[0]
        conn.Data = conn.Data[1:]
        fmt.Printf("Received data from %s\n", conn.ID)
        return data
    }
    fmt.Println("No data to receive")
    return nil
}

func (c *ConnectedState) GetName() string {
    return "Connected"
}

```

## 3.3.1.6 5. 批判性分析

### 3.3.1.6.1 优势

1. **状态封装**: 每个状态独立封装
2. **状态转换**: 清晰的状态转换逻辑
3. **行为变化**: 不同状态下行为不同
4. **扩展性**: 新增状态不影响现有代码

### 3.3.1.6.2 劣势

1. **状态数量**: 状态过多会增加复杂度
2. **转换复杂**: 状态转换逻辑可能复杂
3. **内存开销**: 每个状态都需要对象
4. **调试困难**: 状态转换难以调试

### 3.3.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 | 高 | 中 |
| Java | 接口 | 中 | 中 |
| C++ | 虚函数 | 高 | 中 |
| Python | 类 | 中 | 低 |

### 3.3.1.6.4 最新趋势

1. **状态机**: 使用状态机库
2. **事件驱动**: 基于事件的状态转换
3. **持久化状态**: 状态的持久化存储
4. **分布式状态**: 分布式状态管理

## 3.3.1.7 6. 面试题与考点

### 3.3.1.7.1 基础考点

1. **Q**: 状态模式与策略模式的区别？
   **A**: 状态模式关注状态转换，策略模式关注算法选择

2. **Q**: 什么时候使用状态模式？
   **A**: 对象行为依赖于状态、状态转换复杂时

3. **Q**: 状态模式的优缺点？
   **A**: 优点：状态封装、转换清晰；缺点：状态数量多、转换复杂

### 3.3.1.7.2 进阶考点

1. **Q**: 如何避免状态爆炸？
   **A**: 状态合并、层次化状态、状态机

2. **Q**: 状态模式在微服务中的应用？
   **A**: 服务状态管理、工作流引擎

3. **Q**: 如何处理状态转换的并发？
   **A**: 状态锁、原子操作、事件溯源

## 3.3.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 状态模式 | 管理对象状态的设计模式 | State Pattern |
| 状态 | 对象的当前状态 | State |
| 上下文 | 包含状态的对象 | Context |
| 状态转换 | 从一个状态到另一个状态 | State Transition |

## 3.3.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 状态爆炸 | 状态数量过多 | 状态合并、层次化 |
| 转换复杂 | 状态转换逻辑复杂 | 状态机、事件驱动 |
| 内存开销 | 每个状态都需要对象 | 状态池、单例状态 |
| 调试困难 | 状态转换难以调试 | 日志记录、可视化 |

## 3.3.1.10 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)

## 3.3.1.11 10. 学习路径

### 3.3.1.11.1 新手路径

1. 理解状态模式的基本概念
2. 学习状态和上下文的关系
3. 实现简单的状态模式
4. 理解状态转换机制

### 3.3.1.11.2 进阶路径

1. 学习复杂的状态实现
2. 理解状态的性能优化
3. 掌握状态的应用场景
4. 学习状态的最佳实践

### 3.3.1.11.3 高阶路径

1. 分析状态在大型项目中的应用
2. 理解状态与架构设计的关系
3. 掌握状态的性能调优
4. 学习状态的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)
