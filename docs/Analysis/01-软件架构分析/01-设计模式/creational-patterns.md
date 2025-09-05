# Go语言创建型模式详解

<!-- TOC START -->
- [Go语言创建型模式详解](#go语言创建型模式详解)
  - [1.1 📋 概述](#11--概述)
  - [1.2 🏗️ 单例模式 (Singleton)](#12-️-单例模式-singleton)
    - [1.2.1 概念定义](#121-概念定义)
    - [1.2.2 形式化描述](#122-形式化描述)
    - [1.2.3 Go语言实现](#123-go语言实现)
    - [1.2.4 线程安全实现](#124-线程安全实现)
    - [1.2.5 性能分析](#125-性能分析)
  - [1.3 🏭 工厂方法模式 (Factory Method)](#13--工厂方法模式-factory-method)
    - [1.3.1 概念定义](#131-概念定义)
    - [1.3.2 形式化描述](#132-形式化描述)
    - [1.3.3 Go语言实现](#133-go语言实现)
    - [1.3.4 函数式实现](#134-函数式实现)
  - [1.4 🏢 抽象工厂模式 (Abstract Factory)](#14--抽象工厂模式-abstract-factory)
    - [1.4.1 概念定义](#141-概念定义)
    - [1.4.2 形式化描述](#142-形式化描述)
    - [1.4.3 Go语言实现](#143-go语言实现)
  - [1.5 🔨 建造者模式 (Builder)](#15--建造者模式-builder)
    - [1.5.1 概念定义](#151-概念定义)
    - [1.5.2 形式化描述](#152-形式化描述)
    - [1.5.3 Go语言实现](#153-go语言实现)
    - [1.5.4 函数式建造者](#154-函数式建造者)
  - [1.6 🧬 原型模式 (Prototype)](#16--原型模式-prototype)
    - [1.6.1 概念定义](#161-概念定义)
    - [1.6.2 形式化描述](#162-形式化描述)
    - [1.6.3 Go语言实现](#163-go语言实现)
  - [1.7 📊 性能对比分析](#17--性能对比分析)
  - [1.8 🎯 最佳实践](#18--最佳实践)
<!-- TOC END -->

## 1.1 📋 概述

创建型模式关注对象的创建过程，在Go语言中，这些模式通过接口、结构体、函数和并发原语来实现。Go语言的简洁语法和强大的类型系统为创建型模式提供了优雅的实现方式。

## 1.2 🏗️ 单例模式 (Singleton)

### 1.2.1 概念定义

单例模式确保一个类只有一个实例，并提供一个全局访问点。

**数学定义**:
设 $S$ 为单例类，$I$ 为实例集合，则：
$$|I| = 1 \land \forall i \in I : i \in S$$

### 1.2.2 形式化描述

```go
type Singleton interface {
    GetInstance() *Singleton
}

// 单例约束
type SingletonConstraint struct {
    instance *Singleton
    once     sync.Once
}
```

### 1.2.3 Go语言实现

```go
package main

import (
    "fmt"
    "sync"
)

// 单例结构体
type Database struct {
    connection string
}

// 全局变量方式（简单但不推荐）
var (
    instance *Database
    once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{
            connection: "database_connection",
        }
    })
    return instance
}

// 方法实现
func (db *Database) Connect() {
    fmt.Printf("Connected to: %s\n", db.connection)
}

func main() {
    db1 := GetInstance()
    db2 := GetInstance()
    
    fmt.Printf("db1 == db2: %t\n", db1 == db2) // true
    db1.Connect()
}
```

### 1.2.4 线程安全实现

```go
// 使用sync.Once确保线程安全
type ThreadSafeSingleton struct {
    data string
}

var (
    instance *ThreadSafeSingleton
    once     sync.Once
)

func GetThreadSafeInstance() *ThreadSafeSingleton {
    once.Do(func() {
        instance = &ThreadSafeSingleton{
            data: "thread_safe_data",
        }
    })
    return instance
}

// 使用互斥锁的替代实现
type MutexSingleton struct {
    data string
    mu   sync.RWMutex
}

var (
    mutexInstance *MutexSingleton
    mutexOnce     sync.Once
)

func GetMutexInstance() *MutexSingleton {
    mutexOnce.Do(func() {
        mutexInstance = &MutexSingleton{
            data: "mutex_data",
        }
    })
    return mutexInstance
}

func (s *MutexSingleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}
```

### 1.2.5 性能分析

| 实现方式 | 内存开销 | 性能 | 线程安全 | 推荐度 |
|---------|---------|------|---------|--------|
| sync.Once | 低 | 高 | ✅ | ⭐⭐⭐⭐⭐ |
| sync.Mutex | 中 | 中 | ✅ | ⭐⭐⭐ |
| 全局变量 | 低 | 高 | ❌ | ⭐⭐ |

## 1.3 🏭 工厂方法模式 (Factory Method)

### 1.3.1 概念定义

工厂方法模式定义一个创建对象的接口，但让子类决定实例化哪个类。

**数学定义**:
设 $F$ 为工厂函数，$T$ 为类型集合，$O$ 为对象集合，则：
$$F: T \rightarrow O$$

### 1.3.2 形式化描述

```go
// 产品接口
type Product interface {
    Use() string
}

// 工厂接口
type Factory interface {
    CreateProduct() Product
}
```

### 1.3.3 Go语言实现

```go
package main

import "fmt"

// 产品接口
type Vehicle interface {
    Drive() string
    GetType() string
}

// 具体产品
type Car struct {
    brand string
}

func (c Car) Drive() string {
    return fmt.Sprintf("Driving %s car", c.brand)
}

func (c Car) GetType() string {
    return "Car"
}

type Motorcycle struct {
    brand string
}

func (m Motorcycle) Drive() string {
    return fmt.Sprintf("Riding %s motorcycle", m.brand)
}

func (m Motorcycle) GetType() string {
    return "Motorcycle"
}

// 工厂接口
type VehicleFactory interface {
    CreateVehicle(brand string) Vehicle
}

// 具体工厂
type CarFactory struct{}

func (cf CarFactory) CreateVehicle(brand string) Vehicle {
    return Car{brand: brand}
}

type MotorcycleFactory struct{}

func (mf MotorcycleFactory) CreateVehicle(brand string) Vehicle {
    return Motorcycle{brand: brand}
}

// 工厂注册表
type VehicleFactoryRegistry struct {
    factories map[string]VehicleFactory
}

func NewVehicleFactoryRegistry() *VehicleFactoryRegistry {
    return &VehicleFactoryRegistry{
        factories: make(map[string]VehicleFactory),
    }
}

func (vfr *VehicleFactoryRegistry) RegisterFactory(vehicleType string, factory VehicleFactory) {
    vfr.factories[vehicleType] = factory
}

func (vfr *VehicleFactoryRegistry) CreateVehicle(vehicleType, brand string) (Vehicle, error) {
    factory, exists := vfr.factories[vehicleType]
    if !exists {
        return nil, fmt.Errorf("unknown vehicle type: %s", vehicleType)
    }
    return factory.CreateVehicle(brand), nil
}

func main() {
    registry := NewVehicleFactoryRegistry()
    registry.RegisterFactory("car", CarFactory{})
    registry.RegisterFactory("motorcycle", MotorcycleFactory{})
    
    car, _ := registry.CreateVehicle("car", "Toyota")
    motorcycle, _ := registry.CreateVehicle("motorcycle", "Honda")
    
    fmt.Println(car.Drive())
    fmt.Println(motorcycle.Drive())
}
```

### 1.3.4 函数式实现

```go
// 函数式工厂
type VehicleCreator func(brand string) Vehicle

var vehicleCreators = map[string]VehicleCreator{
    "car": func(brand string) Vehicle {
        return Car{brand: brand}
    },
    "motorcycle": func(brand string) Vehicle {
        return Motorcycle{brand: brand}
    },
}

func CreateVehicle(vehicleType, brand string) (Vehicle, error) {
    creator, exists := vehicleCreators[vehicleType]
    if !exists {
        return nil, fmt.Errorf("unknown vehicle type: %s", vehicleType)
    }
    return creator(brand), nil
}
```

## 1.4 🏢 抽象工厂模式 (Abstract Factory)

### 1.4.1 概念定义

抽象工厂模式提供一个接口，用于创建相关或依赖对象的家族，而不需要指定它们的具体类。

**数学定义**:
设 $AF$ 为抽象工厂，$P$ 为产品族，$F$ 为具体工厂，则：
$$AF: P \rightarrow F$$

### 1.4.2 形式化描述

```go
// 抽象产品族
type AbstractProductA interface {
    OperationA() string
}

type AbstractProductB interface {
    OperationB() string
}

// 抽象工厂
type AbstractFactory interface {
    CreateProductA() AbstractProductA
    CreateProductB() AbstractProductB
}
```

### 1.4.3 Go语言实现

```go
package main

import "fmt"

// 抽象产品
type Button interface {
    Render() string
}

type Dialog interface {
    Show() string
}

// 具体产品 - Windows系列
type WindowsButton struct{}

func (wb WindowsButton) Render() string {
    return "Windows Button"
}

type WindowsDialog struct{}

func (wd WindowsDialog) Show() string {
    return "Windows Dialog"
}

// 具体产品 - Mac系列
type MacButton struct{}

func (mb MacButton) Render() string {
    return "Mac Button"
}

type MacDialog struct{}

func (md MacDialog) Show() string {
    return "Mac Dialog"
}

// 抽象工厂
type UIFactory interface {
    CreateButton() Button
    CreateDialog() Dialog
}

// 具体工厂
type WindowsUIFactory struct{}

func (wuf WindowsUIFactory) CreateButton() Button {
    return WindowsButton{}
}

func (wuf WindowsUIFactory) CreateDialog() Dialog {
    return WindowsDialog{}
}

type MacUIFactory struct{}

func (muf MacUIFactory) CreateButton() Button {
    return MacButton{}
}

func (muf MacUIFactory) CreateDialog() Dialog {
    return MacDialog{}
}

// 客户端代码
type Application struct {
    factory UIFactory
}

func NewApplication(factory UIFactory) *Application {
    return &Application{factory: factory}
}

func (app *Application) CreateUI() {
    button := app.factory.CreateButton()
    dialog := app.factory.CreateDialog()
    
    fmt.Println(button.Render())
    fmt.Println(dialog.Show())
}

func main() {
    // Windows应用
    windowsApp := NewApplication(WindowsUIFactory{})
    windowsApp.CreateUI()
    
    // Mac应用
    macApp := NewApplication(MacUIFactory{})
    macApp.CreateUI()
}
```

## 1.5 🔨 建造者模式 (Builder)

### 1.5.1 概念定义

建造者模式将一个复杂对象的构建与它的表示分离，使得同样的构建过程可以创建不同的表示。

**数学定义**:
设 $B$ 为建造者，$P$ 为产品，$S$ 为构建步骤，则：
$$B: S_1 \times S_2 \times ... \times S_n \rightarrow P$$

### 1.5.2 形式化描述

```go
// 产品
type Product struct {
    PartA string
    PartB string
    PartC string
}

// 建造者接口
type Builder interface {
    BuildPartA() Builder
    BuildPartB() Builder
    BuildPartC() Builder
    GetResult() Product
}
```

### 1.5.3 Go语言实现

```go
package main

import "fmt"

// 产品
type Computer struct {
    CPU    string
    Memory string
    Storage string
    GPU     string
}

func (c Computer) String() string {
    return fmt.Sprintf("Computer{CPU: %s, Memory: %s, Storage: %s, GPU: %s}",
        c.CPU, c.Memory, c.Storage, c.GPU)
}

// 建造者接口
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetStorage(storage string) ComputerBuilder
    SetGPU(gpu string) ComputerBuilder
    Build() Computer
}

// 具体建造者
type GamingComputerBuilder struct {
    computer Computer
}

func NewGamingComputerBuilder() *GamingComputerBuilder {
    return &GamingComputerBuilder{}
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

func (gcb *GamingComputerBuilder) Build() Computer {
    return gcb.computer
}

// 办公电脑建造者
type OfficeComputerBuilder struct {
    computer Computer
}

func NewOfficeComputerBuilder() *OfficeComputerBuilder {
    return &OfficeComputerBuilder{}
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

func (ocb *OfficeComputerBuilder) SetGPU(gpu string) ComputerBuilder {
    ocb.computer.GPU = gpu
    return ocb
}

func (ocb *OfficeComputerBuilder) Build() Computer {
    return ocb.computer
}

// 导演类
type ComputerDirector struct{}

func (cd *ComputerDirector) BuildGamingComputer(builder ComputerBuilder) Computer {
    return builder.
        SetCPU("Intel i9-12900K").
        SetMemory("32GB DDR5").
        SetStorage("1TB NVMe SSD").
        SetGPU("RTX 4080").
        Build()
}

func (cd *ComputerDirector) BuildOfficeComputer(builder ComputerBuilder) Computer {
    return builder.
        SetCPU("Intel i5-12400").
        SetMemory("16GB DDR4").
        SetStorage("512GB SSD").
        SetGPU("Integrated").
        Build()
}

func main() {
    director := &ComputerDirector{}
    
    // 构建游戏电脑
    gamingBuilder := NewGamingComputerBuilder()
    gamingComputer := director.BuildGamingComputer(gamingBuilder)
    fmt.Println("Gaming Computer:", gamingComputer)
    
    // 构建办公电脑
    officeBuilder := NewOfficeComputerBuilder()
    officeComputer := director.BuildOfficeComputer(officeBuilder)
    fmt.Println("Office Computer:", officeComputer)
}
```

### 1.5.4 函数式建造者

```go
// 函数式建造者
type ComputerConfig struct {
    CPU     string
    Memory  string
    Storage string
    GPU     string
}

type ComputerOption func(*ComputerConfig)

func WithCPU(cpu string) ComputerOption {
    return func(c *ComputerConfig) {
        c.CPU = cpu
    }
}

func WithMemory(memory string) ComputerOption {
    return func(c *ComputerConfig) {
        c.Memory = memory
    }
}

func WithStorage(storage string) ComputerOption {
    return func(c *ComputerConfig) {
        c.Storage = storage
    }
}

func WithGPU(gpu string) ComputerOption {
    return func(c *ComputerConfig) {
        c.GPU = gpu
    }
}

func NewComputer(options ...ComputerOption) Computer {
    config := &ComputerConfig{}
    for _, option := range options {
        option(config)
    }
    
    return Computer{
        CPU:     config.CPU,
        Memory:  config.Memory,
        Storage: config.Storage,
        GPU:     config.GPU,
    }
}

// 使用示例
func main() {
    computer := NewComputer(
        WithCPU("Intel i7-12700K"),
        WithMemory("32GB DDR5"),
        WithStorage("1TB NVMe SSD"),
        WithGPU("RTX 4070"),
    )
    
    fmt.Println("Custom Computer:", computer)
}
```

## 1.6 🧬 原型模式 (Prototype)

### 1.6.1 概念定义

原型模式用原型实例指定创建对象的种类，并且通过拷贝这些原型创建新的对象。

**数学定义**:
设 $P$ 为原型，$C$ 为克隆函数，$O$ 为对象，则：
$$C: P \rightarrow O \text{ where } O \cong P$$

### 1.6.2 形式化描述

```go
// 原型接口
type Prototype interface {
    Clone() Prototype
    GetID() string
}
```

### 1.6.3 Go语言实现

```go
package main

import (
    "fmt"
    "time"
)

// 原型接口
type Document interface {
    Clone() Document
    GetTitle() string
    GetContent() string
    SetTitle(title string)
    SetContent(content string)
}

// 具体原型
type Report struct {
    title   string
    content string
    author  string
    date    time.Time
}

func NewReport(title, content, author string) *Report {
    return &Report{
        title:   title,
        content: content,
        author:  author,
        date:    time.Now(),
    }
}

func (r *Report) Clone() Document {
    // 深拷贝
    return &Report{
        title:   r.title,
        content: r.content,
        author:  r.author,
        date:    r.date,
    }
}

func (r *Report) GetTitle() string {
    return r.title
}

func (r *Report) GetContent() string {
    return r.content
}

func (r *Report) SetTitle(title string) {
    r.title = title
}

func (r *Report) SetContent(content string) {
    r.content = content
}

// 原型管理器
type DocumentManager struct {
    prototypes map[string]Document
}

func NewDocumentManager() *DocumentManager {
    return &DocumentManager{
        prototypes: make(map[string]Document),
    }
}

func (dm *DocumentManager) RegisterPrototype(name string, prototype Document) {
    dm.prototypes[name] = prototype
}

func (dm *DocumentManager) CreateDocument(name string) (Document, error) {
    prototype, exists := dm.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    return prototype.Clone(), nil
}

func main() {
    // 创建原型管理器
    manager := NewDocumentManager()
    
    // 注册原型
    reportPrototype := NewReport("Monthly Report", "This is a monthly report", "John Doe")
    manager.RegisterPrototype("report", reportPrototype)
    
    // 克隆文档
    report1, _ := manager.CreateDocument("report")
    report2, _ := manager.CreateDocument("report")
    
    // 修改克隆的文档
    report1.SetTitle("Q1 Report")
    report1.SetContent("This is Q1 report")
    
    report2.SetTitle("Q2 Report")
    report2.SetContent("This is Q2 report")
    
    fmt.Printf("Report 1: %s - %s\n", report1.GetTitle(), report1.GetContent())
    fmt.Printf("Report 2: %s - %s\n", report2.GetTitle(), report2.GetContent())
}
```

## 1.7 📊 性能对比分析

| 模式 | 内存使用 | 创建速度 | 灵活性 | 复杂度 | 适用场景 |
|------|---------|---------|--------|--------|----------|
| 单例 | 低 | 高 | 低 | 低 | 全局资源 |
| 工厂方法 | 中 | 中 | 高 | 中 | 多态创建 |
| 抽象工厂 | 中 | 中 | 高 | 高 | 产品族 |
| 建造者 | 中 | 低 | 高 | 中 | 复杂对象 |
| 原型 | 中 | 高 | 中 | 低 | 相似对象 |

## 1.8 🎯 最佳实践

### 1.8.1 选择原则

1. **单例模式**: 适用于全局唯一资源（数据库连接、配置等）
2. **工厂方法**: 适用于需要多态创建的场景
3. **抽象工厂**: 适用于创建相关产品族的场景
4. **建造者**: 适用于创建复杂对象的场景
5. **原型**: 适用于创建相似对象的场景

### 1.8.2 Go语言特定建议

1. **使用接口**: 充分利用Go语言的接口系统
2. **函数式编程**: 结合函数式编程特性简化实现
3. **并发安全**: 考虑并发环境下的线程安全
4. **性能优化**: 利用Go语言的性能特性
5. **简洁性**: 保持代码简洁，避免过度设计

---

**注意**: 本文档基于`/model/Software/DesignPattern/`目录中的创建型模式内容，结合Go语言特性进行了重新整理和实现，确保内容的准确性和实用性。
