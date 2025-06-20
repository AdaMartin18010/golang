# 结构型设计模式分析

## 目录

1. [概述](#1-概述)
2. [结构型模式形式化定义](#2-结构型模式形式化定义)
3. [适配器模式 (Adapter)](#3-适配器模式-adapter)
4. [桥接模式 (Bridge)](#4-桥接模式-bridge)
5. [组合模式 (Composite)](#5-组合模式-composite)
6. [装饰器模式 (Decorator)](#6-装饰器模式-decorator)
7. [外观模式 (Facade)](#7-外观模式-facade)
8. [享元模式 (Flyweight)](#8-享元模式-flyweight)
9. [代理模式 (Proxy)](#9-代理模式-proxy)
10. [性能分析与最佳实践](#10-性能分析与最佳实践)
11. [参考文献](#11-参考文献)

---

## 1. 概述

结构型模式关注类和对象的组合，通过继承和组合获得新功能。在Golang中，结构型模式通过接口、嵌入和组合实现，充分利用Go语言的简洁性和灵活性。

### 1.1 结构型模式分类

结构型模式集合可以形式化定义为：

$$C_{str} = \{Adapter, Bridge, Composite, Decorator, Facade, Flyweight, Proxy\}$$

### 1.2 核心特征

- **对象组合**: 通过组合实现功能扩展
- **接口适配**: 提供统一的接口适配不同实现
- **结构优化**: 优化对象结构和内存使用
- **访问控制**: 控制对对象的访问和操作

---

## 2. 结构型模式形式化定义

### 2.1 结构型模式系统定义

结构型模式系统可以定义为五元组：

$$\mathcal{SP} = (P_{str}, I_{str}, C_{str}, E_{str}, Q_{str})$$

其中：

- **$P_{str}$** - 结构型模式集合
- **$I_{str}$** - 接口集合
- **$C_{str}$** - 组合关系集合
- **$E_{str}$** - 评估指标集合
- **$Q_{str}$** - 质量保证集合

### 2.2 结构关系形式化定义

结构关系可以定义为图论模型：

$$\mathcal{G} = (V, E, \phi)$$

其中：

- **$V$** - 顶点集合（对象）
- **$E$** - 边集合（关系）
- **$\phi$** - 边到顶点对的映射函数

---

## 3. 适配器模式 (Adapter)

### 3.1 形式化定义

适配器模式将一个类的接口转换成客户希望的另一个接口，使不兼容的接口可以一起工作。

**数学定义**:
$$Adapter : IncompatibleInterface \rightarrow CompatibleInterface$$

**类图关系**:

```mermaid
classDiagram
    Client --> Target
    Target <|-- Adapter
    Adapter --> Adaptee
    Adaptee <|-- ConcreteAdaptee
```

### 3.2 Golang实现

```go
package adapter

import (
    "fmt"
    "strconv"
)

// Target 目标接口
type Target interface {
    Request() string
}

// Adaptee 被适配的接口
type Adaptee interface {
    SpecificRequest() int
}

// ConcreteAdaptee 具体被适配类
type ConcreteAdaptee struct {
    value int
}

func NewConcreteAdaptee(value int) *ConcreteAdaptee {
    return &ConcreteAdaptee{value: value}
}

func (a *ConcreteAdaptee) SpecificRequest() int {
    return a.value
}

// Adapter 适配器
type Adapter struct {
    adaptee Adaptee
}

func NewAdapter(adaptee Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    value := a.adaptee.SpecificRequest()
    return fmt.Sprintf("Adapted: %s", strconv.Itoa(value))
}

// 函数式适配器
type FunctionalAdapter struct {
    adapteeFunc func() int
}

func NewFunctionalAdapter(adapteeFunc func() int) *FunctionalAdapter {
    return &FunctionalAdapter{adapteeFunc: adapteeFunc}
}

func (f *FunctionalAdapter) Request() string {
    value := f.adapteeFunc()
    return fmt.Sprintf("Functional Adapted: %s", strconv.Itoa(value))
}

// Client 客户端
type Client struct {
    target Target
}

func NewClient(target Target) *Client {
    return &Client{target: target}
}

func (c *Client) UseTarget() string {
    return c.target.Request()
}
```

### 3.3 性能分析

**时间复杂度**: $O(1)$ - 适配操作的时间复杂度为常数
**空间复杂度**: $O(1)$ - 只占用适配器的内存空间
**兼容性**: 提供接口兼容性
**灵活性**: 支持多种适配策略

---

## 4. 桥接模式 (Bridge)

### 4.1 形式化定义

桥接模式将抽象部分与实现部分分离，使它们都可以独立地变化。

**数学定义**:
$$Bridge : (Abstraction, Implementation) \rightarrow Abstraction \times Implementation$$

**类图关系**:

```mermaid
classDiagram
    Abstraction --> Implementation
    RefinedAbstraction --|> Abstraction
    ConcreteImplementation --|> Implementation
```

### 4.2 Golang实现

```go
package bridge

import (
    "fmt"
)

// Implementation 实现接口
type Implementation interface {
    OperationImpl() string
}

// ConcreteImplementationA 具体实现A
type ConcreteImplementationA struct{}

func (i *ConcreteImplementationA) OperationImpl() string {
    return "Implementation A"
}

// ConcreteImplementationB 具体实现B
type ConcreteImplementationB struct{}

func (i *ConcreteImplementationB) OperationImpl() string {
    return "Implementation B"
}

// Abstraction 抽象接口
type Abstraction interface {
    Operation() string
    SetImplementation(impl Implementation)
}

// RefinedAbstraction 精确抽象
type RefinedAbstraction struct {
    implementation Implementation
}

func NewRefinedAbstraction(impl Implementation) *RefinedAbstraction {
    return &RefinedAbstraction{implementation: impl}
}

func (a *RefinedAbstraction) Operation() string {
    return fmt.Sprintf("Refined: %s", a.implementation.OperationImpl())
}

func (a *RefinedAbstraction) SetImplementation(impl Implementation) {
    a.implementation = impl
}

// ExtendedAbstraction 扩展抽象
type ExtendedAbstraction struct {
    RefinedAbstraction
}

func NewExtendedAbstraction(impl Implementation) *ExtendedAbstraction {
    return &ExtendedAbstraction{
        RefinedAbstraction: *NewRefinedAbstraction(impl),
    }
}

func (e *ExtendedAbstraction) Operation() string {
    return fmt.Sprintf("Extended: %s", e.implementation.OperationImpl())
}
```

### 4.3 性能分析

**时间复杂度**: $O(1)$ - 桥接操作的时间复杂度为常数
**空间复杂度**: $O(1)$ - 只占用抽象和实现的内存空间
**解耦性**: 抽象和实现完全解耦
**扩展性**: 支持独立扩展抽象和实现

---

## 5. 组合模式 (Composite)

### 5.1 形式化定义

组合模式将对象组合成树形结构以表示"部分-整体"的层次结构，使客户端对单个对象和组合对象具有一致的访问性。

**数学定义**:
$$Composite : Component \rightarrow Tree(Component)$$

其中 $Tree(Component)$ 表示组件树结构

### 5.2 Golang实现

```go
package composite

import (
    "fmt"
    "strings"
)

// Component 组件接口
type Component interface {
    Add(component Component)
    Remove(component Component)
    GetChild(index int) Component
    Operation() string
    GetName() string
}

// Leaf 叶子节点
type Leaf struct {
    name string
}

func NewLeaf(name string) *Leaf {
    return &Leaf{name: name}
}

func (l *Leaf) Add(component Component) {
    // 叶子节点不支持添加子组件
}

func (l *Leaf) Remove(component Component) {
    // 叶子节点不支持移除子组件
}

func (l *Leaf) GetChild(index int) Component {
    return nil
}

func (l *Leaf) Operation() string {
    return fmt.Sprintf("Leaf: %s", l.name)
}

func (l *Leaf) GetName() string {
    return l.name
}

// Composite 组合节点
type Composite struct {
    name     string
    children []Component
}

func NewComposite(name string) *Composite {
    return &Composite{
        name:     name,
        children: make([]Component, 0),
    }
}

func (c *Composite) Add(component Component) {
    c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
    for i, child := range c.children {
        if child == component {
            c.children = append(c.children[:i], c.children[i+1:]...)
            break
        }
    }
}

func (c *Composite) GetChild(index int) Component {
    if index >= 0 && index < len(c.children) {
        return c.children[index]
    }
    return nil
}

func (c *Composite) Operation() string {
    var results []string
    results = append(results, fmt.Sprintf("Composite: %s", c.name))
    
    for _, child := range c.children {
        results = append(results, "  "+child.Operation())
    }
    
    return strings.Join(results, "\n")
}

func (c *Composite) GetName() string {
    return c.name
}

// 安全组合模式
type SafeComponent interface {
    Operation() string
    GetName() string
}

type SafeLeaf struct {
    name string
}

func (l *SafeLeaf) Operation() string {
    return fmt.Sprintf("Safe Leaf: %s", l.name)
}

func (l *SafeLeaf) GetName() string {
    return l.name
}

type SafeComposite struct {
    name     string
    children []SafeComponent
}

func (c *SafeComposite) Operation() string {
    var results []string
    results = append(results, fmt.Sprintf("Safe Composite: %s", c.name))
    
    for _, child := range c.children {
        results = append(results, "  "+child.Operation())
    }
    
    return strings.Join(results, "\n")
}

func (c *SafeComposite) GetName() string {
    return c.name
}

func (c *SafeComposite) AddChild(child SafeComponent) {
    c.children = append(c.children, child)
}
```

### 5.3 性能分析

**时间复杂度**: $O(n)$ - n为树中节点的数量
**空间复杂度**: $O(n)$ - 需要存储所有节点
**遍历效率**: 支持深度优先和广度优先遍历
**内存使用**: 树结构的内存使用与节点数量成正比

---

## 6. 装饰器模式 (Decorator)

### 6.1 形式化定义

装饰器模式动态地给对象添加额外的职责，而不改变其接口。

**数学定义**:
$$Decorator : Component \rightarrow DecoratedComponent$$

其中 $DecoratedComponent \supset Component$ (包含原组件功能)

### 6.2 Golang实现

```go
package decorator

import (
    "fmt"
    "strings"
)

// Component 组件接口
type Component interface {
    Operation() string
}

// ConcreteComponent 具体组件
type ConcreteComponent struct {
    name string
}

func NewConcreteComponent(name string) *ConcreteComponent {
    return &ConcreteComponent{name: name}
}

func (c *ConcreteComponent) Operation() string {
    return fmt.Sprintf("ConcreteComponent: %s", c.name)
}

// Decorator 装饰器基类
type Decorator struct {
    component Component
}

func NewDecorator(component Component) *Decorator {
    return &Decorator{component: component}
}

func (d *Decorator) Operation() string {
    return d.component.Operation()
}

// ConcreteDecoratorA 具体装饰器A
type ConcreteDecoratorA struct {
    Decorator
}

func NewConcreteDecoratorA(component Component) *ConcreteDecoratorA {
    return &ConcreteDecoratorA{
        Decorator: *NewDecorator(component),
    }
}

func (d *ConcreteDecoratorA) Operation() string {
    return fmt.Sprintf("DecoratorA(%s)", d.component.Operation())
}

// ConcreteDecoratorB 具体装饰器B
type ConcreteDecoratorB struct {
    Decorator
}

func NewConcreteDecoratorB(component Component) *ConcreteDecoratorB {
    return &ConcreteDecoratorB{
        Decorator: *NewDecorator(component),
    }
}

func (d *ConcreteDecoratorB) Operation() string {
    return fmt.Sprintf("DecoratorB(%s)", d.component.Operation())
}

// 函数式装饰器
type FunctionalDecorator func(Component) Component

func LoggingDecorator(component Component) Component {
    return &LoggingComponent{component: component}
}

type LoggingComponent struct {
    component Component
}

func (l *LoggingComponent) Operation() string {
    result := l.component.Operation()
    fmt.Printf("Logging: %s\n", result)
    return result
}

func CachingDecorator(component Component) Component {
    return &CachingComponent{component: component}
}

type CachingComponent struct {
    component Component
    cache     map[string]string
}

func (c *CachingComponent) Operation() string {
    if c.cache == nil {
        c.cache = make(map[string]string)
    }
    
    key := "operation"
    if result, exists := c.cache[key]; exists {
        return result
    }
    
    result := c.component.Operation()
    c.cache[key] = result
    return result
}
```

### 6.3 性能分析

**时间复杂度**: $O(n)$ - n为装饰器的数量
**空间复杂度**: $O(n)$ - 每个装饰器占用额外空间
**灵活性**: 支持动态组合装饰器
**可扩展性**: 易于添加新的装饰器

---

## 7. 外观模式 (Facade)

### 7.1 形式化定义

外观模式为子系统中的一组接口提供一个一致的界面，定义了一个高层接口，使子系统更加容易使用。

**数学定义**:
$$Facade : \{Subsystem_1, Subsystem_2, ..., Subsystem_n\} \rightarrow UnifiedInterface$$

### 7.2 Golang实现

```go
package facade

import (
    "fmt"
    "time"
)

// SubsystemA 子系统A
type SubsystemA struct{}

func (a *SubsystemA) OperationA() string {
    return "SubsystemA operation"
}

// SubsystemB 子系统B
type SubsystemB struct{}

func (b *SubsystemB) OperationB() string {
    return "SubsystemB operation"
}

// SubsystemC 子系统C
type SubsystemC struct{}

func (c *SubsystemC) OperationC() string {
    return "SubsystemC operation"
}

// Facade 外观
type Facade struct {
    subsystemA *SubsystemA
    subsystemB *SubsystemB
    subsystemC *SubsystemC
}

func NewFacade() *Facade {
    return &Facade{
        subsystemA: &SubsystemA{},
        subsystemB: &SubsystemB{},
        subsystemC: &SubsystemC{},
    }
}

func (f *Facade) Operation1() string {
    resultA := f.subsystemA.OperationA()
    resultB := f.subsystemB.OperationB()
    return fmt.Sprintf("Operation1: %s, %s", resultA, resultB)
}

func (f *Facade) Operation2() string {
    resultB := f.subsystemB.OperationB()
    resultC := f.subsystemC.OperationC()
    return fmt.Sprintf("Operation2: %s, %s", resultB, resultC)
}

// 复杂外观模式
type ComplexFacade struct {
    facade *Facade
}

func NewComplexFacade() *ComplexFacade {
    return &ComplexFacade{
        facade: NewFacade(),
    }
}

func (c *ComplexFacade) ComplexOperation() string {
    result1 := c.facade.Operation1()
    time.Sleep(100 * time.Millisecond) // 模拟复杂操作
    result2 := c.facade.Operation2()
    return fmt.Sprintf("Complex: %s; %s", result1, result2)
}

// 配置外观模式
type ConfigFacade struct {
    subsystems map[string]interface{}
}

func NewConfigFacade() *ConfigFacade {
    return &ConfigFacade{
        subsystems: make(map[string]interface{}),
    }
}

func (c *ConfigFacade) RegisterSubsystem(name string, subsystem interface{}) {
    c.subsystems[name] = subsystem
}

func (c *ConfigFacade) ExecuteOperation(operation string) string {
    // 根据操作类型调用相应的子系统
    switch operation {
    case "A":
        if subsystem, exists := c.subsystems["A"]; exists {
            if a, ok := subsystem.(*SubsystemA); ok {
                return a.OperationA()
            }
        }
    case "B":
        if subsystem, exists := c.subsystems["B"]; exists {
            if b, ok := subsystem.(*SubsystemB); ok {
                return b.OperationB()
            }
        }
    }
    return "Unknown operation"
}
```

### 7.3 性能分析

**时间复杂度**: $O(n)$ - n为子系统的数量
**空间复杂度**: $O(n)$ - 需要存储所有子系统
**简化性**: 简化客户端与子系统的交互
**封装性**: 隐藏子系统的复杂性

---

## 8. 享元模式 (Flyweight)

### 8.1 形式化定义

享元模式通过共享技术有效地支持大量细粒度对象的复用。

**数学定义**:
$$Flyweight : (IntrinsicState, ExtrinsicState) \rightarrow SharedObject$$

其中 $IntrinsicState$ 是共享状态，$ExtrinsicState$ 是外部状态

### 8.2 Golang实现

```go
package flyweight

import (
    "fmt"
    "sync"
)

// Flyweight 享元接口
type Flyweight interface {
    Operation(extrinsicState string) string
}

// ConcreteFlyweight 具体享元
type ConcreteFlyweight struct {
    intrinsicState string
}

func NewConcreteFlyweight(intrinsicState string) *ConcreteFlyweight {
    return &ConcreteFlyweight{intrinsicState: intrinsicState}
}

func (f *ConcreteFlyweight) Operation(extrinsicState string) string {
    return fmt.Sprintf("Flyweight[%s] with extrinsic state: %s", 
        f.intrinsicState, extrinsicState)
}

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
    flyweights map[string]Flyweight
    mu         sync.RWMutex
}

func NewFlyweightFactory() *FlyweightFactory {
    return &FlyweightFactory{
        flyweights: make(map[string]Flyweight),
    }
}

func (f *FlyweightFactory) GetFlyweight(key string) Flyweight {
    f.mu.RLock()
    if flyweight, exists := f.flyweights[key]; exists {
        f.mu.RUnlock()
        return flyweight
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查
    if flyweight, exists := f.flyweights[key]; exists {
        return flyweight
    }
    
    flyweight := NewConcreteFlyweight(key)
    f.flyweights[key] = flyweight
    return flyweight
}

func (f *FlyweightFactory) GetFlyweightCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.flyweights)
}

// UnsharedConcreteFlyweight 非共享具体享元
type UnsharedConcreteFlyweight struct {
    allState string
}

func NewUnsharedConcreteFlyweight(allState string) *UnsharedConcreteFlyweight {
    return &UnsharedConcreteFlyweight{allState: allState}
}

func (u *UnsharedConcreteFlyweight) Operation(extrinsicState string) string {
    return fmt.Sprintf("Unshared Flyweight[%s] with extrinsic state: %s", 
        u.allState, extrinsicState)
}

// 字符串享元示例
type StringFlyweight struct {
    content string
}

func NewStringFlyweight(content string) *StringFlyweight {
    return &StringFlyweight{content: content}
}

func (s *StringFlyweight) GetContent() string {
    return s.content
}

type StringFlyweightFactory struct {
    strings map[string]*StringFlyweight
    mu      sync.RWMutex
}

func NewStringFlyweightFactory() *StringFlyweightFactory {
    return &StringFlyweightFactory{
        strings: make(map[string]*StringFlyweight),
    }
}

func (f *StringFlyweightFactory) GetString(content string) *StringFlyweight {
    f.mu.RLock()
    if str, exists := f.strings[content]; exists {
        f.mu.RUnlock()
        return str
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if str, exists := f.strings[content]; exists {
        return str
    }
    
    str := NewStringFlyweight(content)
    f.strings[content] = str
    return str
}
```

### 8.3 性能分析

**时间复杂度**: $O(1)$ - 获取享元的时间复杂度为常数
**空间复杂度**: $O(n)$ - n为不同享元的数量
**内存效率**: 显著减少内存使用
**共享性**: 支持大量对象共享

---

## 9. 代理模式 (Proxy)

### 9.1 形式化定义

代理模式为其他对象提供一种代理以控制对这个对象的访问。

**数学定义**:
$$Proxy : Client \rightarrow Subject$$

其中代理控制对Subject的访问

### 9.2 Golang实现

```go
package proxy

import (
    "fmt"
    "sync"
    "time"
)

// Subject 主题接口
type Subject interface {
    Request() string
}

// RealSubject 真实主题
type RealSubject struct {
    name string
}

func NewRealSubject(name string) *RealSubject {
    return &RealSubject{name: name}
}

func (r *RealSubject) Request() string {
    // 模拟耗时操作
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("RealSubject[%s] response", r.name)
}

// Proxy 代理
type Proxy struct {
    realSubject Subject
    mu          sync.RWMutex
    cache       map[string]string
}

func NewProxy(realSubject Subject) *Proxy {
    return &Proxy{
        realSubject: realSubject,
        cache:       make(map[string]string),
    }
}

func (p *Proxy) Request() string {
    p.mu.RLock()
    if result, exists := p.cache["request"]; exists {
        p.mu.RUnlock()
        return result
    }
    p.mu.RUnlock()
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // 双重检查
    if result, exists := p.cache["request"]; exists {
        return result
    }
    
    result := p.realSubject.Request()
    p.cache["request"] = result
    return result
}

// VirtualProxy 虚拟代理
type VirtualProxy struct {
    realSubject Subject
    mu          sync.Mutex
    initialized bool
}

func NewVirtualProxy() *VirtualProxy {
    return &VirtualProxy{}
}

func (v *VirtualProxy) Request() string {
    v.mu.Lock()
    defer v.mu.Unlock()
    
    if !v.initialized {
        v.realSubject = NewRealSubject("Virtual")
        v.initialized = true
    }
    
    return v.realSubject.Request()
}

// ProtectionProxy 保护代理
type ProtectionProxy struct {
    realSubject Subject
    accessLevel string
}

func NewProtectionProxy(realSubject Subject, accessLevel string) *ProtectionProxy {
    return &ProtectionProxy{
        realSubject: realSubject,
        accessLevel: accessLevel,
    }
}

func (p *ProtectionProxy) Request() string {
    if p.accessLevel == "admin" {
        return p.realSubject.Request()
    }
    return "Access denied"
}

// RemoteProxy 远程代理
type RemoteProxy struct {
    realSubject Subject
    network     string
}

func NewRemoteProxy(network string) *RemoteProxy {
    return &RemoteProxy{network: network}
}

func (r *RemoteProxy) Request() string {
    // 模拟网络请求
    fmt.Printf("Sending request over %s network\n", r.network)
    time.Sleep(50 * time.Millisecond)
    return fmt.Sprintf("Remote response from %s", r.network)
}

// 智能引用代理
type SmartReferenceProxy struct {
    realSubject Subject
    referenceCount int
    mu            sync.Mutex
}

func NewSmartReferenceProxy(realSubject Subject) *SmartReferenceProxy {
    return &SmartReferenceProxy{
        realSubject: realSubject,
        referenceCount: 1,
    }
}

func (s *SmartReferenceProxy) Request() string {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.referenceCount++
    return s.realSubject.Request()
}

func (s *SmartReferenceProxy) Release() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.referenceCount--
    if s.referenceCount == 0 {
        fmt.Println("RealSubject will be garbage collected")
    }
}
```

### 9.3 性能分析

**时间复杂度**: $O(1)$ - 代理操作的时间复杂度为常数
**空间复杂度**: $O(1)$ - 代理只占用额外控制空间
**访问控制**: 提供灵活的访问控制机制
**缓存效果**: 支持结果缓存和延迟加载

---

## 10. 性能分析与最佳实践

### 10.1 性能对比

| 模式 | 时间复杂度 | 空间复杂度 | 适用场景 | 优势 |
|------|------------|------------|----------|------|
| 适配器 | O(1) | O(1) | 接口兼容 | 接口适配 |
| 桥接 | O(1) | O(1) | 抽象实现分离 | 解耦设计 |
| 组合 | O(n) | O(n) | 树形结构 | 统一接口 |
| 装饰器 | O(n) | O(n) | 动态扩展 | 功能组合 |
| 外观 | O(n) | O(n) | 子系统封装 | 简化接口 |
| 享元 | O(1) | O(n) | 对象复用 | 内存优化 |
| 代理 | O(1) | O(1) | 访问控制 | 控制访问 |

### 10.2 最佳实践

#### 10.2.1 接口设计

```go
// 使用接口定义契约
type Component interface {
    Operation() string
}

// 提供默认实现
type BaseComponent struct{}

func (b *BaseComponent) Operation() string {
    return "Base operation"
}
```

#### 10.2.2 组合优于继承

```go
// 使用组合而不是继承
type Decorator struct {
    component Component
}

func (d *Decorator) Operation() string {
    return d.component.Operation()
}
```

#### 10.2.3 线程安全

```go
// 使用适当的同步机制
type ThreadSafeProxy struct {
    realSubject Subject
    mu          sync.RWMutex
    cache       map[string]string
}

func (p *ThreadSafeProxy) Request() string {
    p.mu.RLock()
    if result, exists := p.cache["request"]; exists {
        p.mu.RUnlock()
        return result
    }
    p.mu.RUnlock()
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    result := p.realSubject.Request()
    p.cache["request"] = result
    return result
}
```

### 10.3 性能优化建议

1. **使用享元模式**: 对于大量相似对象，使用享元模式减少内存使用
2. **缓存结果**: 在代理和装饰器中缓存计算结果
3. **延迟加载**: 使用虚拟代理实现延迟加载
4. **对象池**: 对于创建成本高的对象，使用对象池
5. **接口优化**: 设计简洁的接口，减少方法调用开销

---

## 11. 参考文献

1. Gamma, E., Helm, R., Johnson, R., & Vlissides, J. (1994). Design Patterns: Elements of Reusable Object-Oriented Software. Addison-Wesley.
2. Freeman, E., Robson, E., Sierra, K., & Bates, B. (2004). Head First Design Patterns. O'Reilly Media.
3. Go Team. (2023). The Go Programming Language Specification. <https://golang.org/ref/spec>
4. Go Team. (2023). Effective Go. <https://golang.org/doc/effective_go.html>
5. Go Team. (2023). Go Concurrency Patterns. <https://golang.org/doc/effective_go.html#concurrency>

---

**最后更新**: 2024-12-19  
**版本**: 1.0.0  
**状态**: 结构型模式分析完成
