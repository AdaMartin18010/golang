# Go 1.23 设计模式全面指南

> 本文档系统梳理 **Go 1.23** 语言实现的23种经典设计模式，涵盖创建型、结构型、行为型三大类别。
>
> **版本**: 2.0
> **适用Go版本**: Go 1.23+
> **更新内容**:
>
> - 新增迭代器模式（利用Go 1.23 range-over-func特性）
> - 使用 `iter.Seq` / `iter.Seq2` 优化集合遍历
> - 结合 `slices` / `maps` 包的迭代器函数简化代码

---

## 目录

- [Go 1.23 设计模式全面指南](#go-123-设计模式全面指南)
  - [目录](#目录)
  - [前言](#前言)
  - [创建型模式](#创建型模式)
    - [1. 单例模式（Singleton）](#1-单例模式singleton)
      - [概念定义](#概念定义)
      - [意图与动机](#意图与动机)
      - [UML结构图](#uml结构图)
      - [Go语言实现](#go语言实现)
      - [完整示例：配置管理器](#完整示例配置管理器)
      - [反例说明](#反例说明)
      - [优缺点分析](#优缺点分析)
    - [2. 工厂方法模式（Factory Method）](#2-工厂方法模式factory-method)
      - [概念定义](#概念定义-1)
      - [意图与动机](#意图与动机-1)
      - [UML结构图](#uml结构图-1)
      - [Go语言实现](#go语言实现-1)
      - [完整示例：日志记录器工厂](#完整示例日志记录器工厂)
      - [反例说明](#反例说明-1)
      - [优缺点分析](#优缺点分析-1)
    - [3. 抽象工厂模式（Abstract Factory）](#3-抽象工厂模式abstract-factory)
      - [概念定义](#概念定义-2)
      - [意图与动机](#意图与动机-2)
      - [UML结构图](#uml结构图-2)
      - [Go语言实现](#go语言实现-2)
      - [完整示例：跨平台UI组件](#完整示例跨平台ui组件)
      - [反例说明](#反例说明-2)
      - [优缺点分析](#优缺点分析-2)
    - [4. 建造者模式（Builder）](#4-建造者模式builder)
      - [概念定义](#概念定义-3)
      - [意图与动机](#意图与动机-3)
      - [UML结构图](#uml结构图-3)
      - [Go语言实现](#go语言实现-3)
      - [完整示例：HTTP请求构建器](#完整示例http请求构建器)
      - [反例说明](#反例说明-3)
      - [优缺点分析](#优缺点分析-3)
    - [5. 原型模式（Prototype）](#5-原型模式prototype)
      - [概念定义](#概念定义-4)
      - [意图与动机](#意图与动机-4)
      - [UML结构图](#uml结构图-4)
      - [Go语言实现](#go语言实现-4)
      - [完整示例：文档模板系统](#完整示例文档模板系统)
      - [反例说明](#反例说明-4)
      - [优缺点分析](#优缺点分析-4)
  - [结构型模式](#结构型模式)
    - [6. 适配器模式（Adapter）](#6-适配器模式adapter)
      - [概念定义](#概念定义-5)
      - [意图与动机](#意图与动机-5)
      - [UML结构图](#uml结构图-5)
      - [Go语言实现](#go语言实现-5)
      - [完整示例：支付系统适配器](#完整示例支付系统适配器)
      - [反例说明](#反例说明-5)
      - [优缺点分析](#优缺点分析-5)
    - [7. 桥接模式（Bridge）](#7-桥接模式bridge)
      - [概念定义](#概念定义-6)
      - [意图与动机](#意图与动机-6)
      - [UML结构图](#uml结构图-6)
      - [Go语言实现](#go语言实现-6)
      - [完整示例：图形渲染系统](#完整示例图形渲染系统)
      - [反例说明](#反例说明-6)
      - [优缺点分析](#优缺点分析-6)
    - [8. 组合模式（Composite）](#8-组合模式composite)
      - [概念定义](#概念定义-7)
      - [意图与动机](#意图与动机-7)
      - [UML结构图](#uml结构图-7)
      - [Go语言实现](#go语言实现-7)
      - [完整示例：文件系统](#完整示例文件系统)
      - [反例说明](#反例说明-7)
      - [优缺点分析](#优缺点分析-7)
    - [9. 装饰器模式（Decorator）](#9-装饰器模式decorator)
      - [概念定义](#概念定义-8)
      - [意图与动机](#意图与动机-8)
      - [UML结构图](#uml结构图-8)
      - [Go语言实现](#go语言实现-8)
      - [完整示例：HTTP中间件链](#完整示例http中间件链)
      - [反例说明](#反例说明-8)
      - [优缺点分析](#优缺点分析-8)
    - [10. 外观模式（Facade）](#10-外观模式facade)
      - [概念定义](#概念定义-9)
      - [意图与动机](#意图与动机-9)
      - [UML结构图](#uml结构图-9)
      - [Go语言实现](#go语言实现-9)
      - [完整示例：电商订单系统外观](#完整示例电商订单系统外观)
      - [反例说明](#反例说明-9)
      - [优缺点分析](#优缺点分析-9)
    - [11. 享元模式（Flyweight）](#11-享元模式flyweight)
      - [概念定义](#概念定义-10)
      - [意图与动机](#意图与动机-10)
      - [UML结构图](#uml结构图-10)
      - [Go语言实现](#go语言实现-10)
      - [完整示例：游戏角色渲染系统](#完整示例游戏角色渲染系统)
      - [反例说明](#反例说明-10)
      - [优缺点分析](#优缺点分析-10)
    - [12. 代理模式（Proxy）](#12-代理模式proxy)
      - [概念定义](#概念定义-11)
      - [意图与动机](#意图与动机-11)
      - [UML结构图](#uml结构图-11)
      - [Go语言实现](#go语言实现-11)
      - [完整示例：图片加载代理](#完整示例图片加载代理)
      - [反例说明](#反例说明-11)
      - [优缺点分析](#优缺点分析-11)
  - [行为型模式](#行为型模式)
    - [13. 责任链模式（Chain of Responsibility）](#13-责任链模式chain-of-responsibility)
      - [概念定义](#概念定义-12)
      - [意图与动机](#意图与动机-12)
      - [UML结构图](#uml结构图-12)
      - [Go语言实现](#go语言实现-12)
      - [完整示例：请求处理链](#完整示例请求处理链)
      - [反例说明](#反例说明-12)
      - [优缺点分析](#优缺点分析-12)
    - [14. 命令模式（Command）](#14-命令模式command)
      - [概念定义](#概念定义-13)
      - [意图与动机](#意图与动机-13)
      - [UML结构图](#uml结构图-13)
      - [Go语言实现](#go语言实现-13)
      - [完整示例：文本编辑器](#完整示例文本编辑器)
      - [反例说明](#反例说明-13)
      - [优缺点分析](#优缺点分析-13)
    - [15. 解释器模式（Interpreter）](#15-解释器模式interpreter)
      - [概念定义](#概念定义-14)
      - [意图与动机](#意图与动机-14)
      - [UML结构图](#uml结构图-14)
      - [Go语言实现](#go语言实现-14)
      - [完整示例：规则引擎](#完整示例规则引擎)
      - [反例说明](#反例说明-14)
      - [优缺点分析](#优缺点分析-14)
    - [16. 迭代器模式（Iterator）](#16-迭代器模式iterator)
      - [概念定义](#概念定义-15)
      - [意图与动机](#意图与动机-15)
      - [UML结构图](#uml结构图-15)
      - [Go语言实现](#go语言实现-15)
      - [完整示例：自定义集合迭代器](#完整示例自定义集合迭代器)
      - [反例说明](#反例说明-15)
      - [优缺点分析](#优缺点分析-15)
    - [17. 中介者模式（Mediator）](#17-中介者模式mediator)
      - [概念定义](#概念定义-16)
      - [意图与动机](#意图与动机-16)
      - [UML结构图](#uml结构图-16)
      - [Go语言实现](#go语言实现-16)
      - [完整示例：聊天室系统](#完整示例聊天室系统)
      - [反例说明](#反例说明-16)
      - [优缺点分析](#优缺点分析-16)
    - [18. 备忘录模式（Memento）](#18-备忘录模式memento)
      - [概念定义](#概念定义-17)
      - [意图与动机](#意图与动机-17)
      - [UML结构图](#uml结构图-17)
      - [Go语言实现](#go语言实现-17)
      - [完整示例：文档编辑器撤销系统](#完整示例文档编辑器撤销系统)
      - [反例说明](#反例说明-17)
      - [优缺点分析](#优缺点分析-17)
    - [19. 观察者模式（Observer）](#19-观察者模式observer)
      - [概念定义](#概念定义-18)
      - [意图与动机](#意图与动机-18)
      - [UML结构图](#uml结构图-18)
      - [Go语言实现](#go语言实现-18)
      - [完整示例：事件总线系统](#完整示例事件总线系统)
      - [反例说明](#反例说明-18)
      - [优缺点分析](#优缺点分析-18)
    - [20. 状态模式（State）](#20-状态模式state)
      - [概念定义](#概念定义-19)
      - [意图与动机](#意图与动机-19)
      - [UML结构图](#uml结构图-19)
      - [Go语言实现](#go语言实现-19)
      - [完整示例：订单状态机](#完整示例订单状态机)
      - [反例说明](#反例说明-19)
      - [优缺点分析](#优缺点分析-19)
    - [21. 策略模式（Strategy）](#21-策略模式strategy)
      - [概念定义](#概念定义-20)
      - [意图与动机](#意图与动机-20)
      - [UML结构图](#uml结构图-20)
      - [Go语言实现](#go语言实现-20)
      - [完整示例：支付策略系统](#完整示例支付策略系统)
      - [反例说明](#反例说明-20)
      - [优缺点分析](#优缺点分析-20)
    - [22. 模板方法模式（Template Method）](#22-模板方法模式template-method)
      - [概念定义](#概念定义-21)
      - [意图与动机](#意图与动机-21)
      - [UML结构图](#uml结构图-21)
      - [Go语言实现](#go语言实现-21)
      - [完整示例：数据导入框架](#完整示例数据导入框架)
      - [反例说明](#反例说明-21)
      - [优缺点分析](#优缺点分析-21)
    - [23. 访问者模式（Visitor）](#23-访问者模式visitor)
      - [概念定义](#概念定义-22)
      - [意图与动机](#意图与动机-22)
      - [UML结构图](#uml结构图-22)
      - [Go语言实现](#go语言实现-22)
      - [完整示例：文档导出系统](#完整示例文档导出系统)
      - [反例说明](#反例说明-22)
      - [优缺点分析](#优缺点分析-22)
  - [总结](#总结)
    - [设计模式对比表](#设计模式对比表)
    - [Go语言设计模式最佳实践](#go语言设计模式最佳实践)
    - [参考资源](#参考资源)

---

## 前言

设计模式是软件工程中对常见问题的可复用解决方案。Go语言作为一门现代编程语言，具有独特的特性：

- **接口隐式实现**: 无需显式声明实现关系
- **组合优于继承**: 通过嵌入实现代码复用
- **一等函数**: 函数可作为参数和返回值
- **Goroutine与Channel**: 原生并发支持
- **零值初始化**: 所有类型都有默认值

这些特性深刻影响了设计模式在Go中的实现方式，使得某些模式变得更简单，而另一些则需要不同的实现策略。

---

## 创建型模式

创建型模式关注对象的创建机制，提供在不同场景下创建对象的灵活方式。

---

### 1. 单例模式（Singleton）

#### 概念定义

单例模式确保一个类只有一个实例，并提供一个全局访问点。该模式控制实例化过程，确保系统中只存在一个实例对象。

#### 意图与动机

- **资源控制**: 数据库连接池、线程池等需要限制实例数量
- **全局状态**: 配置管理器、日志记录器等需要统一状态
- **协调访问**: 文件系统、设备驱动等需要串行访问的资源

#### UML结构图

```text
┌─────────────────┐
│   Singleton     │
├─────────────────┤
│ - instance      │◄──── 唯一实例
│ - mutex         │      同步锁
├─────────────────┤
│ + GetInstance() │◄──── 全局访问点
│ + BusinessMethod()│
└─────────────────┘
```

#### Go语言实现

Go语言实现单例模式有多种方式，从简单到复杂：

**方式一：懒汉式（线程不安全）**

```go
package singleton

// 不推荐：非线程安全
type Singleton struct {
    data string
}

var instance *Singleton

func GetInstance() *Singleton {
    if instance == nil {
        instance = &Singleton{data: "initialized"}
    }
    return instance
}
```

**方式二：饿汉式（推荐简单场景）**

```go
package singleton

// 包初始化时创建，线程安全
var instance = &Singleton{data: "initialized"}

type Singleton struct {
    data string
}

func GetInstance() *Singleton {
    return instance
}
```

**方式三：双重检查锁定（推荐复杂场景）**

```go
package singleton

import (
    "sync"
    "sync/atomic"
)

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
```

#### 完整示例：配置管理器

```go
package main

import (
    "fmt"
    "sync"
)

// ConfigManager 配置管理器单例
type ConfigManager struct {
    mu     sync.RWMutex
    config map[string]string
}

var (
    configInstance *ConfigManager
    configOnce     sync.Once
)

// GetConfigManager 获取配置管理器实例
func GetConfigManager() *ConfigManager {
    configOnce.Do(func() {
        configInstance = &ConfigManager{
            config: make(map[string]string),
        }
        // 加载默认配置
        configInstance.config["app.name"] = "MyApp"
        configInstance.config["app.version"] = "1.0.0"
        configInstance.config["db.host"] = "localhost"
        configInstance.config["db.port"] = "5432"
    })
    return configInstance
}

// Get 获取配置值
func (c *ConfigManager) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.config[key]
}

// Set 设置配置值
func (c *ConfigManager) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.config[key] = value
}

// GetAll 获取所有配置
func (c *ConfigManager) GetAll() map[string]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    // 返回副本避免外部修改
    copy := make(map[string]string)
    for k, v := range c.config {
        copy[k] = v
    }
    return copy
}

func main() {
    // 获取单例实例
    cm1 := GetConfigManager()
    cm2 := GetConfigManager()

    // 验证是同一实例
    fmt.Printf("同一实例: %v\n", cm1 == cm2)

    // 使用配置
    fmt.Printf("App Name: %s\n", cm1.Get("app.name"))
    fmt.Printf("App Version: %s\n", cm1.Get("app.version"))

    // 修改配置
    cm1.Set("app.version", "2.0.0")
    fmt.Printf("Updated Version: %s\n", cm2.Get("app.version"))

    // 并发测试
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            cm := GetConfigManager()
            cm.Set(fmt.Sprintf("key_%d", n), fmt.Sprintf("value_%d", n))
        }(i)
    }
    wg.Wait()

    fmt.Printf("配置数量: %d\n", len(cm1.GetAll()))
}
```

#### 反例说明

**错误1：非线程安全**

```go
// 错误：并发环境下可能创建多个实例
func GetInstance() *Singleton {
    if instance == nil {          // 线程A和B同时到达这里
        instance = &Singleton{}   // 可能都执行创建
    }
    return instance
}
```

**错误2：过度使用锁**

```go
// 错误：每次访问都加锁，性能差
var mutex sync.Mutex

func GetInstance() *Singleton {
    mutex.Lock()          // 不必要的锁
    defer mutex.Unlock()
    if instance == nil {
        instance = &Singleton{}
    }
    return instance
}
```

**错误3：使用全局变量直接暴露**

```go
// 错误：破坏了封装，外部可直接修改
var Instance *Singleton  // 不要这样做
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 控制实例数量，节约资源 | 违反单一职责原则 |
| 提供全局访问点 | 隐藏依赖关系，难以测试 |
| 延迟初始化可能 | 并发环境下需要额外处理 |
| 避免重复创建开销 | 扩展困难（难以有多个实例） |

**适用场景**:

- 数据库连接池
- 配置管理器
- 日志记录器
- 缓存管理器
- 线程池

**Go语言特化**:

- 使用`sync.Once`实现线程安全的延迟初始化
- 包级变量实现饿汉式单例
- 避免使用`init()`函数进行复杂初始化

---

### 2. 工厂方法模式（Factory Method）

#### 概念定义

工厂方法模式定义了一个创建对象的接口，但让子类决定实例化哪个类。工厂方法将类的实例化延迟到子类。

#### 意图与动机

- **解耦**: 将对象创建与使用分离
- **扩展性**: 新增产品类型无需修改现有代码
- **多态**: 通过统一接口创建不同类型的对象

#### UML结构图

```
                    ┌─────────────────┐
                    │     Product     │◄────── 产品接口
                    │   (interface)   │
                    └────────┬────────┘
                             │
           ┌─────────────────┼─────────────────┐
           │                 │                 │
    ┌──────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐
    │ ConcreteP1  │   │ ConcreteP2  │   │ ConcreteP3  │
    └─────────────┘   └─────────────┘   └─────────────┘
           ▲                 ▲                 ▲
           │                 │                 │
    ┌──────┴─────────────────┴─────────────────┴──────┐
    │                   Creator                       │◄──── 创建者接口
    │                 (interface)                     │
    │  + FactoryMethod() Product                      │
    └────────────────────┬────────────────────────────┘
                         │
              ┌──────────▼──────────┐
              │   ConcreteCreator   │◄──── 具体创建者
              │  + FactoryMethod()  │
              └─────────────────────┘
```

#### Go语言实现

Go语言通过接口和函数类型实现工厂方法：

```go
package factorymethod

// Product 产品接口
type Product interface {
    Use() string
    GetName() string
}

// Creator 创建者接口
type Creator interface {
    FactoryMethod() Product
}

// ConcreteProductA 具体产品A
type ConcreteProductA struct {
    name string
}

func (p *ConcreteProductA) Use() string {
    return "使用产品A"
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

// ConcreteProductB 具体产品B
type ConcreteProductB struct {
    name string
}

func (p *ConcreteProductB) Use() string {
    return "使用产品B"
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

// ConcreteCreatorA 具体创建者A
type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) FactoryMethod() Product {
    return &ConcreteProductA{name: "产品A"}
}

// ConcreteCreatorB 具体创建者B
type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) FactoryMethod() Product {
    return &ConcreteProductB{name: "产品B"}
}
```

#### 完整示例：日志记录器工厂

```go
package main

import (
    "fmt"
    "time"
)

// Logger 日志记录器接口
type Logger interface {
    Log(message string)
    Logf(format string, args ...interface{})
}

// LoggerFactory 日志工厂接口
type LoggerFactory interface {
    CreateLogger() Logger
}

// ConsoleLogger 控制台日志
type ConsoleLogger struct {
    prefix string
}

func (c *ConsoleLogger) Log(message string) {
    fmt.Printf("[%s] [CONSOLE] %s: %s\n",
        time.Now().Format("2006-01-02 15:04:05"), c.prefix, message)
}

func (c *ConsoleLogger) Logf(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    c.Log(message)
}

// FileLogger 文件日志
type FileLogger struct {
    filename string
    prefix   string
}

func (f *FileLogger) Log(message string) {
    // 实际实现会写入文件
    fmt.Printf("[%s] [FILE:%s] %s: %s\n",
        time.Now().Format("2006-01-02 15:04:05"), f.filename, f.prefix, message)
}

func (f *FileLogger) Logf(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    f.Log(message)
}

// NetworkLogger 网络日志
type NetworkLogger struct {
    endpoint string
    prefix   string
}

func (n *NetworkLogger) Log(message string) {
    // 实际实现会发送到网络端点
    fmt.Printf("[%s] [NETWORK:%s] %s: %s\n",
        time.Now().Format("2006-01-02 15:04:05"), n.endpoint, n.prefix, message)
}

func (n *NetworkLogger) Logf(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    n.Log(message)
}

// ConsoleLoggerFactory 控制台日志工厂
type ConsoleLoggerFactory struct {
    prefix string
}

func (f *ConsoleLoggerFactory) CreateLogger() Logger {
    return &ConsoleLogger{prefix: f.prefix}
}

// FileLoggerFactory 文件日志工厂
type FileLoggerFactory struct {
    filename string
    prefix   string
}

func (f *FileLoggerFactory) CreateLogger() Logger {
    return &FileLogger{filename: f.filename, prefix: f.prefix}
}

// NetworkLoggerFactory 网络日志工厂
type NetworkLoggerFactory struct {
    endpoint string
    prefix   string
}

func (f *NetworkLoggerFactory) CreateLogger() Logger {
    return &NetworkLogger{endpoint: f.endpoint, prefix: f.prefix}
}

// Client 客户端代码
func Client(factory LoggerFactory) {
    logger := factory.CreateLogger()
    logger.Log("系统启动")
    logger.Logf("当前时间: %s", time.Now().Format("15:04:05"))
    logger.Log("操作完成")
}

func main() {
    // 使用控制台日志
    fmt.Println("=== 控制台日志 ===")
    consoleFactory := &ConsoleLoggerFactory{prefix: "APP"}
    Client(consoleFactory)

    // 使用文件日志
    fmt.Println("\n=== 文件日志 ===")
    fileFactory := &FileLoggerFactory{filename: "/var/log/app.log", prefix: "APP"}
    Client(fileFactory)

    // 使用网络日志
    fmt.Println("\n=== 网络日志 ===")
    networkFactory := &NetworkLoggerFactory{endpoint: "http://logs.example.com", prefix: "APP"}
    Client(networkFactory)

    // 运行时动态选择
    fmt.Println("\n=== 动态选择 ===")
    var factory LoggerFactory
    config := "file" // 从配置读取

    switch config {
    case "console":
        factory = &ConsoleLoggerFactory{prefix: "DYNAMIC"}
    case "file":
        factory = &FileLoggerFactory{filename: "dynamic.log", prefix: "DYNAMIC"}
    case "network":
        factory = &NetworkLoggerFactory{endpoint: "http://dynamic.example.com", prefix: "DYNAMIC"}
    default:
        factory = &ConsoleLoggerFactory{prefix: "DYNAMIC"}
    }

    Client(factory)
}
```

#### 反例说明

**错误1：直接使用new创建**

```go
// 错误：客户端直接依赖具体类
func main() {
    logger := &ConsoleLogger{}  // 紧耦合
    logger.Log("message")
}
```

**错误2：工厂返回具体类型**

```go
// 错误：返回具体类型而非接口
func (f *ConsoleLoggerFactory) CreateLogger() *ConsoleLogger {
    return &ConsoleLogger{}  // 应该返回 Logger 接口
}
```

**错误3：工厂包含业务逻辑**

```go
// 错误：工厂不应该包含产品使用逻辑
type LoggerFactory struct{}

func (f *LoggerFactory) CreateAndUse() {
    logger := &ConsoleLogger{}
    logger.Log("message")  // 工厂只应该负责创建
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 符合开闭原则 | 类数量增加 |
| 解耦创建与使用 | 代码复杂度增加 |
| 易于扩展新产品 | 每个产品都需要对应工厂 |
| 单一职责原则 | 简单场景可能过度设计 |

**适用场景**:

- 日志记录系统
- 数据库连接创建
- UI组件创建
- 文档解析器

**Go语言特化**:

- 使用函数类型简化工厂：`type Factory func() Product`
- 结合闭包创建参数化工厂
- 利用接口隐式实现减少样板代码

---

### 3. 抽象工厂模式（Abstract Factory）

#### 概念定义

抽象工厂模式提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们具体的类。该模式创建的是"产品族"。

#### 意图与动机

- **产品族一致性**: 确保同一风格/主题的产品一起使用
- **跨平台**: 不同平台（Windows/Mac）的UI组件
- **主题切换**: 深色/浅色主题的整体切换

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                    AbstractFactory                          │
│  + CreateButton() Button                                    │
│  + CreateCheckbox() Checkbox                                │
│  + CreateTextField() TextField                              │
└────────────┬────────────────────────────────────────────────┘
             │
    ┌────────┴────────┐
    │                 │
┌───▼────┐      ┌────▼────┐
│Windows  │      │  Mac    │
│Factory  │      │ Factory │
└───┬────┘      └────┬────┘
    │                 │
    ▼                 ▼
┌────────┐      ┌────────┐
│WinBtn  │      │MacBtn  │
│WinChk  │      │MacChk  │
│WinTxt  │      │MacTxt  │
└────────┘      └────────┘
```

#### Go语言实现

```go
package abstractfactory

// Button 按钮接口
type Button interface {
    Render() string
    OnClick()
}

// Checkbox 复选框接口
type Checkbox interface {
    Render() string
    Toggle()
}

// TextField 文本框接口
type TextField interface {
    Render() string
    SetText(text string)
}

// GUIFactory GUI工厂接口
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
    CreateTextField() TextField
}
```

#### 完整示例：跨平台UI组件

```go
package main

import "fmt"

// ============ 产品接口 ============

type Button interface {
    Render() string
    OnClick()
}

type Checkbox interface {
    Render() string
    Toggle()
}

type TextField interface {
    Render() string
    SetText(text string)
}

// ============ Windows 产品族 ============

type WindowsButton struct {
    label string
}

func (b *WindowsButton) Render() string {
    return fmt.Sprintf("[Windows按钮] %s", b.label)
}

func (b *WindowsButton) OnClick() {
    fmt.Println("Windows按钮被点击")
}

type WindowsCheckbox struct {
    checked bool
    label   string
}

func (c *WindowsCheckbox) Render() string {
    status := "[ ]"
    if c.checked {
        status = "[X]"
    }
    return fmt.Sprintf("[Windows复选框] %s %s", status, c.label)
}

func (c *WindowsCheckbox) Toggle() {
    c.checked = !c.checked
    fmt.Printf("Windows复选框状态: %v\n", c.checked)
}

type WindowsTextField struct {
    text string
}

func (t *WindowsTextField) Render() string {
    return fmt.Sprintf("[Windows文本框] |%s|", t.text)
}

func (t *WindowsTextField) SetText(text string) {
    t.text = text
    fmt.Printf("Windows文本框内容: %s\n", text)
}

// ============ Mac 产品族 ============

type MacButton struct {
    label string
}

func (b *MacButton) Render() string {
    return fmt.Sprintf("⟨Mac按钮⟩ %s", b.label)
}

func (b *MacButton) OnClick() {
    fmt.Println("Mac按钮被点击")
}

type MacCheckbox struct {
    checked bool
    label   string
}

func (c *MacCheckbox) Render() string {
    status := "○"
    if c.checked {
        status = "◉"
    }
    return fmt.Sprintf("⟨Mac复选框⟩ %s %s", status, c.label)
}

func (c *MacCheckbox) Toggle() {
    c.checked = !c.checked
    fmt.Printf("Mac复选框状态: %v\n", c.checked)
}

type MacTextField struct {
    text string
}

func (t *MacTextField) Render() string {
    return fmt.Sprintf("⟨Mac文本框⟩ [%s]", t.text)
}

func (t *MacTextField) SetText(text string) {
    t.text = text
    fmt.Printf("Mac文本框内容: %s\n", text)
}

// ============ 具体工厂 ============

type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
    return &WindowsButton{label: "Windows风格"}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{label: "Windows风格", checked: false}
}

func (f *WindowsFactory) CreateTextField() TextField {
    return &WindowsTextField{}
}

type MacFactory struct{}

func (f *MacFactory) CreateButton() Button {
    return &MacButton{label: "Mac风格"}
}

func (f *MacFactory) CreateCheckbox() Checkbox {
    return &MacCheckbox{label: "Mac风格", checked: false}
}

func (f *MacFactory) CreateTextField() TextField {
    return &MacTextField{}
}

// ============ 应用程序 ============

type Application struct {
    factory GUIFactory
    button  Button
    checkbox Checkbox
    textfield TextField
}

func NewApplication(factory GUIFactory) *Application {
    return &Application{
        factory: factory,
    }
}

func (a *Application) CreateUI() {
    a.button = a.factory.CreateButton()
    a.checkbox = a.factory.CreateCheckbox()
    a.textfield = a.factory.CreateTextField()
}

func (a *Application) Render() {
    fmt.Println(a.button.Render())
    fmt.Println(a.checkbox.Render())
    fmt.Println(a.textfield.Render())
}

type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
    CreateTextField() TextField
}

// GetFactory 根据平台获取工厂
func GetFactory(platform string) GUIFactory {
    switch platform {
    case "windows":
        return &WindowsFactory{}
    case "mac":
        return &MacFactory{}
    default:
        return &WindowsFactory{}
    }
}

func main() {
    // Windows 应用
    fmt.Println("========== Windows 应用 ==========")
    winFactory := GetFactory("windows")
    winApp := NewApplication(winFactory)
    winApp.CreateUI()
    winApp.Render()

    // Mac 应用
    fmt.Println("\n========== Mac 应用 ==========")
    macFactory := GetFactory("mac")
    macApp := NewApplication(macFactory)
    macApp.CreateUI()
    macApp.Render()

    // 运行时切换
    fmt.Println("\n========== 运行时切换 ==========")
    platforms := []string{"windows", "mac"}
    for _, p := range platforms {
        fmt.Printf("\n--- 切换到 %s ---\n", p)
        factory := GetFactory(p)
        app := NewApplication(factory)
        app.CreateUI()
        app.Render()
    }
}
```

#### 反例说明

**错误1：混用不同产品族**

```go
// 错误：混用Windows按钮和Mac复选框
button := &WindowsButton{}
checkbox := &MacCheckbox{}  // 风格不一致！
```

**错误2：客户端直接创建产品**

```go
// 错误：客户端应该知道工厂而非具体产品
func main() {
    button := &WindowsButton{}  // 应该通过工厂创建
}
```

**错误3：工厂方法返回错误类型**

```go
// 错误：返回不兼容的类型
func (f *WindowsFactory) CreateButton() *WindowsButton {
    // 应该返回 Button 接口
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 产品族一致性保证 | 类数量大幅增加 |
| 易于切换产品族 | 扩展新产品类型困难 |
| 符合开闭原则（产品族） | 代码复杂度较高 |
| 隔离具体类 | 接口变更影响所有工厂 |

**适用场景**:

- 跨平台UI框架
- 主题系统（深色/浅色）
- 不同数据库访问层
- 游戏资源管理（不同画质级别）

**Go语言特化**:

- 使用map注册工厂实现动态加载
- 结合配置文件选择工厂
- 利用类型断言进行运行时类型检查

---

### 4. 建造者模式（Builder）

#### 概念定义

建造者模式将一个复杂对象的构建与其表示分离，使得同样的构建过程可以创建不同的表示。该模式适用于创建具有多个组成部分的复杂对象。

#### 意图与动机

- **复杂对象创建**: 对象有多个可选/必选参数
- **不可变对象**: 创建后状态不可修改
- **可读性**: 链式调用提高代码可读性

#### UML结构图

```
┌──────────────────────────────────────────────────────────┐
│                      Director                            │
│  - builder: Builder                                      │
│  + Construct()                                           │
└──────────────────────────┬───────────────────────────────┘
                           │ uses
                           ▼
┌──────────────────────────────────────────────────────────┐
│                      Builder                             │◄──── 建造者接口
│  + BuildPartA()                                          │
│  + BuildPartB()                                          │
│  + BuildPartC()                                          │
│  + GetResult() Product                                   │
└──────────────────────────┬───────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│    ConcreteBuilder1     │   │    ConcreteBuilder2     │
│  - product: Product     │   │  - product: Product     │
├─────────────────────────┤   ├─────────────────────────┤
│  + BuildPartA()         │   │  + BuildPartA()         │
│  + BuildPartB()         │   │  + BuildPartB()         │
│  + BuildPartC()         │   │  + BuildPartC()         │
│  + GetResult()          │   │  + GetResult()          │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

Go语言有两种实现方式：传统方式和函数式选项模式。

**方式一：传统建造者**

```go
package builder

// Product 产品
type Product struct {
    partA string
    partB string
    partC string
}

func (p *Product) String() string {
    return fmt.Sprintf("Product{A=%s, B=%s, C=%s}", p.partA, p.partB, p.partC)
}

// Builder 建造者接口
type Builder interface {
    BuildPartA()
    BuildPartB()
    BuildPartC()
    GetResult() *Product
}

// ConcreteBuilder 具体建造者
type ConcreteBuilder struct {
    product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
    return &ConcreteBuilder{product: &Product{}}
}

func (b *ConcreteBuilder) BuildPartA() {
    b.product.partA = "PartA"
}

func (b *ConcreteBuilder) BuildPartB() {
    b.product.partB = "PartB"
}

func (b *ConcreteBuilder) BuildPartC() {
    b.product.partC = "PartC"
}

func (b *ConcreteBuilder) GetResult() *Product {
    return b.product
}

// Director 指挥者
type Director struct {
    builder Builder
}

func NewDirector(builder Builder) *Director {
    return &Director{builder: builder}
}

func (d *Director) Construct() {
    d.builder.BuildPartA()
    d.builder.BuildPartB()
    d.builder.BuildPartC()
}
```

**方式二：函数式选项模式（Go惯用法）**

```go
package builder

// Config 配置对象
type Config struct {
    Host     string
    Port     int
    Username string
    Password string
    Timeout  int
    Retries  int
    Debug    bool
}

// Option 配置选项函数类型
type Option func(*Config)

// WithHost 设置主机
func WithHost(host string) Option {
    return func(c *Config) {
        c.Host = host
    }
}

// WithPort 设置端口
func WithPort(port int) Option {
    return func(c *Config) {
        c.Port = port
    }
}

// WithAuth 设置认证信息
func WithAuth(username, password string) Option {
    return func(c *Config) {
        c.Username = username
        c.Password = password
    }
}

// WithTimeout 设置超时
func WithTimeout(timeout int) Option {
    return func(c *Config) {
        c.Timeout = timeout
    }
}

// WithRetries 设置重试次数
func WithRetries(retries int) Option {
    return func(c *Config) {
        c.Retries = retries
    }
}

// WithDebug 启用调试模式
func WithDebug() Option {
    return func(c *Config) {
        c.Debug = true
    }
}

// NewConfig 创建配置（函数式选项模式）
func NewConfig(opts ...Option) *Config {
    // 默认值
    cfg := &Config{
        Host:    "localhost",
        Port:    8080,
        Timeout: 30,
        Retries: 3,
        Debug:   false,
    }

    // 应用选项
    for _, opt := range opts {
        opt(cfg)
    }

    return cfg
}
```

#### 完整示例：HTTP请求构建器

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// HTTPRequest HTTP请求
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    io.Reader
    Timeout time.Duration
}

func (r *HTTPRequest) String() string {
    return fmt.Sprintf("%s %s (timeout: %v)", r.Method, r.URL, r.Timeout)
}

// RequestBuilder 请求构建器
type RequestBuilder struct {
    method  string
    url     string
    headers map[string]string
    body    interface{}
    timeout time.Duration
}

// NewRequestBuilder 创建请求构建器
func NewRequestBuilder() *RequestBuilder {
    return &RequestBuilder{
        headers: make(map[string]string),
        timeout: 30 * time.Second,
    }
}

// Method 设置请求方法
func (b *RequestBuilder) Method(method string) *RequestBuilder {
    b.method = method
    return b
}

// URL 设置URL
func (b *RequestBuilder) URL(url string) *RequestBuilder {
    b.url = url
    return b
}

// Get 快捷方法：GET请求
func (b *RequestBuilder) Get(url string) *RequestBuilder {
    b.method = http.MethodGet
    b.url = url
    return b
}

// Post 快捷方法：POST请求
func (b *RequestBuilder) Post(url string) *RequestBuilder {
    b.method = http.MethodPost
    b.url = url
    return b
}

// Put 快捷方法：PUT请求
func (b *RequestBuilder) Put(url string) *RequestBuilder {
    b.method = http.MethodPut
    b.url = url
    return b
}

// Delete 快捷方法：DELETE请求
func (b *RequestBuilder) Delete(url string) *RequestBuilder {
    b.method = http.MethodDelete
    b.url = url
    return b
}

// Header 添加请求头
func (b *RequestBuilder) Header(key, value string) *RequestBuilder {
    b.headers[key] = value
    return b
}

// ContentType 设置Content-Type
func (b *RequestBuilder) ContentType(contentType string) *RequestBuilder {
    b.headers["Content-Type"] = contentType
    return b
}

// JSON 设置JSON内容类型和请求体
func (b *RequestBuilder) JSON(body interface{}) *RequestBuilder {
    b.headers["Content-Type"] = "application/json"
    b.body = body
    return b
}

// Body 设置请求体
func (b *RequestBuilder) Body(body interface{}) *RequestBuilder {
    b.body = body
    return b
}

// Timeout 设置超时
func (b *RequestBuilder) Timeout(timeout time.Duration) *RequestBuilder {
    b.timeout = timeout
    return b
}

// Build 构建请求
func (b *RequestBuilder) Build() (*http.Request, error) {
    var bodyReader io.Reader

    if b.body != nil {
        switch v := b.body.(type) {
        case string:
            bodyReader = bytes.NewBufferString(v)
        case []byte:
            bodyReader = bytes.NewBuffer(v)
        default:
            jsonData, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            bodyReader = bytes.NewBuffer(jsonData)
        }
    }

    req, err := http.NewRequest(b.method, b.url, bodyReader)
    if err != nil {
        return nil, err
    }

    for key, value := range b.headers {
        req.Header.Set(key, value)
    }

    return req, nil
}

// Execute 执行请求
func (b *RequestBuilder) Execute() (*http.Response, error) {
    req, err := b.Build()
    if err != nil {
        return nil, err
    }

    client := &http.Client{Timeout: b.timeout}
    return client.Do(req)
}

// 函数式选项模式版本
type RequestOption func(*HTTPRequest)

func NewHTTPRequest(url string, opts ...RequestOption) *HTTPRequest {
    req := &HTTPRequest{
        Method:  http.MethodGet,
        URL:     url,
        Headers: make(map[string]string),
        Timeout: 30 * time.Second,
    }

    for _, opt := range opts {
        opt(req)
    }

    return req
}

func WithMethod(method string) RequestOption {
    return func(r *HTTPRequest) {
        r.Method = method
    }
}

func WithHeader(key, value string) RequestOption {
    return func(r *HTTPRequest) {
        r.Headers[key] = value
    }
}

func WithBody(body io.Reader) RequestOption {
    return func(r *HTTPRequest) {
        r.Body = body
    }
}

func WithRequestTimeout(timeout time.Duration) RequestOption {
    return func(r *HTTPRequest) {
        r.Timeout = timeout
    }
}

func main() {
    // 方式一：链式调用建造者
    fmt.Println("=== 链式调用建造者 ===")

    req1, err := NewRequestBuilder().
        Get("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Header("X-Request-ID", "abc-123").
        Timeout(10 * time.Second).
        Build()

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Request: %s %s\n", req1.Method, req1.URL)
        fmt.Printf("Headers: %v\n", req1.Header)
    }

    // POST请求
    req2, err := NewRequestBuilder().
        Post("https://api.example.com/users").
        JSON(map[string]interface{}{
            "name":  "John",
            "email": "john@example.com",
        }).
        Header("Authorization", "Bearer token123").
        Build()

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("\nRequest: %s %s\n", req2.Method, req2.URL)
        fmt.Printf("Content-Type: %s\n", req2.Header.Get("Content-Type"))
    }

    // 方式二：函数式选项模式
    fmt.Println("\n=== 函数式选项模式 ===")

    req3 := NewHTTPRequest("https://api.example.com/users",
        WithMethod(http.MethodPost),
        WithHeader("Authorization", "Bearer token456"),
        WithHeader("X-API-Key", "secret-key"),
        WithRequestTimeout(5*time.Second),
    )

    fmt.Printf("Request: %s %s\n", req3.Method, req3.URL)
    fmt.Printf("Headers: %v\n", req3.Headers)
    fmt.Printf("Timeout: %v\n", req3.Timeout)
}
```

#### 反例说明

**错误1：过多参数的构造函数**

```go
// 错误：参数过多，难以维护
func NewConfig(host string, port int, username string, password string,
               timeout int, retries int, debug bool) *Config {
    // 调用时：NewConfig("host", 8080, "user", "pass", 30, 3, false)
    // 难以知道每个参数的含义
}
```

**错误2：可变的建造者状态**

```go
// 错误：建造者应该每次创建新实例
type Builder struct {
    product *Product
}

func (b *Builder) Build() *Product {
    return b.product  // 每次都返回同一实例！
}
```

**错误3：缺少验证**

```go
// 错误：没有验证必需参数
func (b *RequestBuilder) Build() *http.Request {
    // 应该检查 method 和 url 是否设置
    return &http.Request{}
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 可读性强（链式调用） | 代码量增加 |
| 可选参数灵活 | 简单对象可能过度设计 |
| 不可变对象支持 | 需要额外内存 |
| 参数含义清晰 | 学习成本 |

**适用场景**:

- 复杂对象创建（多个参数）
- SQL查询构建
- HTTP请求构建
- 配置对象创建
- 测试数据构建

**Go语言特化**:

- 函数式选项模式是Go社区标准做法
- 广泛用于标准库（如`grpc.Dial`）
- 结合默认值提供良好体验

---

### 5. 原型模式（Prototype）

#### 概念定义

原型模式通过复制现有对象来创建新对象，而不是通过实例化类。该模式允许对象在不了解创建细节的情况下创建其他对象。

#### 意图与动机

- **避免子类化**: 不需要为每个类创建工厂
- **运行时灵活性**: 动态添加/删除产品
- **复杂对象**: 对象创建成本高时，复制更高效

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Prototype                              │◄──── 原型接口
│  + Clone() Prototype                                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│   ConcretePrototype1    │   │   ConcretePrototype2    │
│  - field1               │   │  - fieldA               │
│  - field2               │   │  - fieldB               │
├─────────────────────────┤   ├─────────────────────────┤
│  + Clone() Prototype    │   │  + Clone() Prototype    │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

Go语言实现原型模式有两种方式：实现Clone方法或使用序列化。

```go
package prototype

// Prototype 原型接口
type Prototype interface {
    Clone() Prototype
}

// ConcretePrototype 具体原型
type ConcretePrototype struct {
    Name  string
    Value int
    Items []string
}

// Clone 浅拷贝
func (c *ConcretePrototype) Clone() Prototype {
    // 创建新实例
    clone := &ConcretePrototype{
        Name:  c.Name,
        Value: c.Value,
        Items: c.Items, // 浅拷贝：共享切片
    }
    return clone
}

// DeepClone 深拷贝
func (c *ConcretePrototype) DeepClone() Prototype {
    // 创建新实例并复制所有字段
    clone := &ConcretePrototype{
        Name:  c.Name,
        Value: c.Value,
        Items: make([]string, len(c.Items)),
    }
    copy(clone.Items, c.Items) // 深拷贝：复制切片内容
    return clone
}
```

#### 完整示例：文档模板系统

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// Document 文档接口
type Document interface {
    Clone() Document
    Display() string
}

// BaseDocument 基础文档
type BaseDocument struct {
    ID        string
    Title     string
    Content   string
    Author    string
    CreatedAt time.Time
    Tags      []string
    Metadata  map[string]string
}

// Clone 浅拷贝
func (d *BaseDocument) Clone() Document {
    // 浅拷贝（共享切片和map）
    return &BaseDocument{
        ID:        d.ID + "_copy",
        Title:     d.Title,
        Content:   d.Content,
        Author:    d.Author,
        CreatedAt: time.Now(),
        Tags:      d.Tags,     // 共享
        Metadata:  d.Metadata, // 共享
    }
}

// DeepClone 深拷贝
func (d *BaseDocument) DeepClone() Document {
    // 深拷贝（复制所有内容）
    clone := &BaseDocument{
        ID:        d.ID + "_copy",
        Title:     d.Title,
        Content:   d.Content,
        Author:    d.Author,
        CreatedAt: time.Now(),
        Tags:      make([]string, len(d.Tags)),
        Metadata:  make(map[string]string),
    }

    copy(clone.Tags, d.Tags)
    for k, v := range d.Metadata {
        clone.Metadata[k] = v
    }

    return clone
}

func (d *BaseDocument) Display() string {
    return fmt.Sprintf("[%s] %s by %s (Tags: %v)", d.ID, d.Title, d.Author, d.Tags)
}

// Report 报告文档
type Report struct {
    BaseDocument
    Department string
    Period     string
    Data       []float64
}

func (r *Report) Clone() Document {
    baseClone := r.BaseDocument.DeepClone().(*BaseDocument)
    clone := &Report{
        BaseDocument: *baseClone,
        Department:   r.Department,
        Period:       r.Period,
        Data:         make([]float64, len(r.Data)),
    }
    copy(clone.Data, r.Data)
    return clone
}

func (r *Report) Display() string {
    return fmt.Sprintf("[Report:%s] %s - %s (%s)",
        r.ID, r.Title, r.Department, r.Period)
}

// Contract 合同文档
type Contract struct {
    BaseDocument
    Client     string
    Amount     float64
    Signatures []string
}

func (c *Contract) Clone() Document {
    baseClone := c.BaseDocument.DeepClone().(*BaseDocument)
    clone := &Contract{
        BaseDocument: *baseClone,
        Client:       c.Client,
        Amount:       c.Amount,
        Signatures:   make([]string, len(c.Signatures)),
    }
    copy(clone.Signatures, c.Signatures)
    return clone
}

func (c *Contract) Display() string {
    return fmt.Sprintf("[Contract:%s] %s - Client: %s ($%.2f)",
        c.ID, c.Title, c.Client, c.Amount)
}

// DocumentRegistry 文档注册表（原型管理器）
type DocumentRegistry struct {
    prototypes map[string]Document
}

func NewDocumentRegistry() *DocumentRegistry {
    return &DocumentRegistry{
        prototypes: make(map[string]Document),
    }
}

func (r *DocumentRegistry) Register(name string, prototype Document) {
    r.prototypes[name] = prototype
}

func (r *DocumentRegistry) Create(name string) (Document, error) {
    prototype, ok := r.prototypes[name]
    if !ok {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    return prototype.Clone(), nil
}

// JSONClone 使用JSON序列化实现深拷贝
func JSONClone(src, dst interface{}) error {
    data, err := json.Marshal(src)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dst)
}

func main() {
    // 创建文档注册表
    registry := NewDocumentRegistry()

    // 注册报告模板
    reportTemplate := &Report{
        BaseDocument: BaseDocument{
            ID:        "template_report",
            Title:     "月度销售报告",
            Content:   "报告内容模板...",
            Author:    "系统自动生成",
            Tags:      []string{"report", "sales", "monthly"},
            Metadata:  map[string]string{"department": "sales"},
        },
        Department: "销售部",
        Period:     "2024-01",
        Data:       []float64{1000, 2000, 3000},
    }
    registry.Register("monthly_report", reportTemplate)

    // 注册合同模板
    contractTemplate := &Contract{
        BaseDocument: BaseDocument{
            ID:        "template_contract",
            Title:     "标准服务合同",
            Content:   "合同条款模板...",
            Author:    "法务部",
            Tags:      []string{"contract", "service"},
            Metadata:  map[string]string{"type": "service"},
        },
        Client:     "客户名称",
        Amount:     0,
        Signatures: []string{},
    }
    registry.Register("service_contract", contractTemplate)

    // 从模板创建新文档
    fmt.Println("=== 从模板创建文档 ===")

    doc1, err := registry.Create("monthly_report")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println(doc1.Display())
    }

    doc2, err := registry.Create("service_contract")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println(doc2.Display())
    }

    // 深拷贝 vs 浅拷贝对比
    fmt.Println("\n=== 深拷贝 vs 浅拷贝 ===")

    original := &BaseDocument{
        ID:    "original",
        Title: "原始文档",
        Tags:  []string{"tag1", "tag2"},
        Metadata: map[string]string{"key": "value"},
    }

    // 浅拷贝
    shallowClone := original.Clone().(*BaseDocument)
    // 深拷贝
    deepClone := original.DeepClone().(*BaseDocument)

    // 修改原始文档
    original.Tags[0] = "modified"
    original.Metadata["key"] = "modified"

    fmt.Printf("Original Tags[0]: %s\n", original.Tags[0])
    fmt.Printf("Shallow Clone Tags[0]: %s (共享)\n", shallowClone.Tags[0])
    fmt.Printf("Deep Clone Tags[0]: %s (独立)\n", deepClone.Tags[0])

    // 使用JSON深拷贝
    fmt.Println("\n=== JSON深拷贝 ===")
    jsonClone := &BaseDocument{}
    if err := JSONClone(original, jsonClone); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("JSON Clone Title: %s\n", jsonClone.Title)
    }
}
```

#### 反例说明

**错误1：浅拷贝导致意外修改**

```go
// 错误：没有深拷贝引用类型
func (c *Config) Clone() *Config {
    return &Config{
        Items: c.Items,  // 共享切片，修改会影响原对象
    }
}
```

**错误2：复制未导出字段**

```go
// 错误：无法访问未导出字段
type Document struct {
    name string  // 未导出
}

func (d *Document) Clone() *Document {
    return &Document{
        name: d.name,  // 同一包内可以，跨包不行
    }
}
```

**错误3：循环引用导致无限递归**

```go
// 错误：循环引用
type Node struct {
    Value int
    Next  *Node  // 循环引用会导致Clone无限递归
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 隐藏创建复杂性 | 深拷贝实现复杂 |
| 运行时添加产品 | 循环引用处理困难 |
| 减少子类数量 | 需要实现Clone方法 |
| 性能优于new | 克隆复杂对象成本高 |

**适用场景**:

- 文档模板系统
- 游戏对象复制
- 配置对象复制
- 复杂对象的快速创建

**Go语言特化**:

- 使用JSON序列化实现通用深拷贝
- 注意切片和map的引用语义
- 利用接口实现多态克隆

---

## 结构型模式

结构型模式关注类和对象的组合，通过继承或组合创建更大的结构。

---

### 6. 适配器模式（Adapter）

#### 概念定义

适配器模式将一个类的接口转换成客户希望的另一个接口。适配器让原本接口不兼容的类可以一起工作。

#### 意图与动机

- **接口兼容**: 使用现有类但接口不匹配
- **复用**: 复用已有的类而不修改其代码
- **集成**: 集成第三方库或遗留系统

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                       Client                                │
│                    (使用Target)                             │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                       Target                                │◄──── 目标接口
│  + Request()                                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│      Adaptee            │   │       Adapter           │
│  + SpecificRequest()    │◄──│  - adaptee: Adaptee     │
│                         │   │  + Request()            │
└─────────────────────────┘   │    adaptee.SpecificRequest()
                              └─────────────────────────┘
```

#### Go语言实现

Go语言通过嵌入和接口实现适配器模式：

```go
package adapter

// Target 目标接口
type Target interface {
    Request() string
}

// Adaptee 被适配者
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Adaptee specific request"
}

// Adapter 适配器（对象适配器）
type Adapter struct {
    adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}

// 类适配器（Go通过嵌入模拟）
type ClassAdapter struct {
    Adaptee  // 嵌入
}

func (c *ClassAdapter) Request() string {
    return c.SpecificRequest()
}
```

#### 完整示例：支付系统适配器

```go
package main

import (
    "fmt"
    "time"
)

// PaymentProcessor 标准支付接口（Target）
type PaymentProcessor interface {
    ProcessPayment(amount float64, currency string) (string, error)
    RefundPayment(transactionID string) error
    GetTransactionStatus(transactionID string) (string, error)
}

// ============ 旧版支付系统（需要适配） ============

// LegacyPaymentSystem 旧版支付系统
type LegacyPaymentSystem struct {
    apiKey string
}

func NewLegacyPaymentSystem(apiKey string) *LegacyPaymentSystem {
    return &LegacyPaymentSystem{apiKey: apiKey}
}

// MakePayment 旧版支付方法
func (l *LegacyPaymentSystem) MakePayment(amount int, currencyCode string) (int, error) {
    fmt.Printf("[Legacy] Processing payment: %d %s\n", amount, currencyCode)
    // 模拟处理
    return 1000 + int(time.Now().Unix())%1000, nil
}

// CancelPayment 旧版取消方法
func (l *LegacyPaymentSystem) CancelPayment(transactionID int) error {
    fmt.Printf("[Legacy] Cancelling transaction: %d\n", transactionID)
    return nil
}

// CheckStatus 旧版查询方法
func (l *LegacyPaymentSystem) CheckStatus(transactionID int) string {
    return "completed"
}

// ============ 第三方支付系统（需要适配） ============

// ThirdPartyPayment 第三方支付
type ThirdPartyPayment struct {
    merchantID string
    secretKey  string
}

func NewThirdPartyPayment(merchantID, secretKey string) *ThirdPartyPayment {
    return &ThirdPartyPayment{merchantID: merchantID, secretKey: secretKey}
}

// CreateTransaction 第三方创建交易
func (t *ThirdPartyPayment) CreateTransaction(params map[string]interface{}) map[string]interface{} {
    fmt.Printf("[ThirdParty] Creating transaction: %v\n", params)
    return map[string]interface{}{
        "transaction_id": fmt.Sprintf("TP%d", time.Now().Unix()),
        "status":         "success",
    }
}

// ReverseTransaction 第三方退款
func (t *ThirdPartyPayment) ReverseTransaction(txID string) map[string]interface{} {
    fmt.Printf("[ThirdParty] Reversing transaction: %s\n", txID)
    return map[string]interface{}{"status": "reversed"}
}

// QueryTransaction 第三方查询
func (t *ThirdPartyPayment) QueryTransaction(txID string) map[string]interface{} {
    return map[string]interface{}{
        "transaction_id": txID,
        "status":         "completed",
    }
}

// ============ 适配器实现 ============

// LegacyPaymentAdapter 旧版支付适配器
type LegacyPaymentAdapter struct {
    legacy *LegacyPaymentSystem
}

func NewLegacyPaymentAdapter(apiKey string) *LegacyPaymentAdapter {
    return &LegacyPaymentAdapter{
        legacy: NewLegacyPaymentSystem(apiKey),
    }
}

func (a *LegacyPaymentAdapter) ProcessPayment(amount float64, currency string) (string, error) {
    // 转换金额（元转分）
    amountInCents := int(amount * 100)
    txID, err := a.legacy.MakePayment(amountInCents, currency)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("LEGACY_%d", txID), nil
}

func (a *LegacyPaymentAdapter) RefundPayment(transactionID string) error {
    // 解析ID
    var legacyID int
    fmt.Sscanf(transactionID, "LEGACY_%d", &legacyID)
    return a.legacy.CancelPayment(legacyID)
}

func (a *LegacyPaymentAdapter) GetTransactionStatus(transactionID string) (string, error) {
    var legacyID int
    fmt.Sscanf(transactionID, "LEGACY_%d", &legacyID)
    return a.legacy.CheckStatus(legacyID), nil
}

// ThirdPartyPaymentAdapter 第三方支付适配器
type ThirdPartyPaymentAdapter struct {
    thirdParty *ThirdPartyPayment
}

func NewThirdPartyPaymentAdapter(merchantID, secretKey string) *ThirdPartyPaymentAdapter {
    return &ThirdPartyPaymentAdapter{
        thirdParty: NewThirdPartyPayment(merchantID, secretKey),
    }
}

func (a *ThirdPartyPaymentAdapter) ProcessPayment(amount float64, currency string) (string, error) {
    params := map[string]interface{}{
        "amount":   amount,
        "currency": currency,
    }
    result := a.thirdParty.CreateTransaction(params)
    if txID, ok := result["transaction_id"].(string); ok {
        return txID, nil
    }
    return "", fmt.Errorf("payment failed")
}

func (a *ThirdPartyPaymentAdapter) RefundPayment(transactionID string) error {
    result := a.thirdParty.ReverseTransaction(transactionID)
    if status, ok := result["status"].(string); ok && status == "reversed" {
        return nil
    }
    return fmt.Errorf("refund failed")
}

func (a *ThirdPartyPaymentAdapter) GetTransactionStatus(transactionID string) (string, error) {
    result := a.thirdParty.QueryTransaction(transactionID)
    if status, ok := result["status"].(string); ok {
        return status, nil
    }
    return "", fmt.Errorf("status not found")
}

// ============ 支付服务 ============

type PaymentService struct {
    processor PaymentProcessor
}

func NewPaymentService(processor PaymentProcessor) *PaymentService {
    return &PaymentService{processor: processor}
}

func (s *PaymentService) Checkout(amount float64, currency string) error {
    txID, err := s.processor.ProcessPayment(amount, currency)
    if err != nil {
        return err
    }
    fmt.Printf("Payment successful! Transaction ID: %s\n", txID)

    status, _ := s.processor.GetTransactionStatus(txID)
    fmt.Printf("Transaction status: %s\n", status)

    return nil
}

func main() {
    fmt.Println("========== 使用旧版支付系统 ==========")
    legacyAdapter := NewLegacyPaymentAdapter("legacy-api-key-123")
    service1 := NewPaymentService(legacyAdapter)
    service1.Checkout(99.99, "USD")

    fmt.Println("\n========== 使用第三方支付系统 ==========")
    thirdPartyAdapter := NewThirdPartyPaymentAdapter("merchant-456", "secret-key")
    service2 := NewPaymentService(thirdPartyAdapter)
    service2.Checkout(149.99, "EUR")

    // 运行时切换
    fmt.Println("\n========== 运行时切换支付系统 ==========")
    processors := []PaymentProcessor{
        NewLegacyPaymentAdapter("key1"),
        NewThirdPartyPaymentAdapter("m1", "s1"),
    }

    for i, processor := range processors {
        fmt.Printf("\n--- Payment Method %d ---\n", i+1)
        service := NewPaymentService(processor)
        service.Checkout(50.0, "CNY")
    }
}
```

#### 反例说明

**错误1：修改被适配者代码**

```go
// 错误：直接修改第三方库
func (t *ThirdPartyPayment) ProcessPayment() {  // 不要修改！
    // ...
}
```

**错误2：适配器暴露被适配者**

```go
// 错误：暴露内部实现
type Adapter struct {
    Adaptee  // 公开嵌入，破坏了封装
}
```

**错误3：适配器包含业务逻辑**

```go
// 错误：适配器不应该有业务逻辑
func (a *Adapter) Request() {
    a.adaptee.SpecificRequest()
    // 业务逻辑不应该在这里
    doBusinessLogic()
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 复用现有类 | 增加系统复杂度 |
| 解耦客户端和被适配者 | 过多适配器难以维护 |
| 符合开闭原则 | 调用链路变长 |
| 灵活性高 | 需要额外内存 |

**适用场景**:

- 集成第三方库
- 遗留系统迁移
- 统一不同接口
- 数据格式转换

**Go语言特化**:

- 通过嵌入实现类适配器
- 接口隐式实现简化适配
- 函数类型适配函数签名

---

### 7. 桥接模式（Bridge）

#### 概念定义

桥接模式将抽象部分与实现部分分离，使它们都可以独立变化。该模式通过组合而非继承来连接抽象和实现。

#### 意图与动机

- **维度分离**: 多个独立变化的维度
- **避免类爆炸**: 避免m×n个子类
- **运行时切换**: 动态改变实现

#### UML结构图

```
        Abstraction                    Implementation
        (抽象部分)                      (实现部分)
┌─────────────────────┐        ┌─────────────────────┐
│   Shape             │        │   Renderer          │◄──── 实现接口
│ - renderer:Renderer │◄───────┤  + RenderCircle()   │
│ + Draw()            │  组合  │  + RenderSquare()   │
└─────────┬───────────┘        └──────────┬──────────┘
          │                               │
    ┌─────┴─────┐               ┌─────────┴─────────┐
    │           │               │                   │
┌───▼───┐   ┌───▼───┐      ┌────▼────┐       ┌─────▼────┐
│Circle │   │Square │      │  SVG    │       │  Canvas  │
└───────┘   └───────┘      │Renderer │       │ Renderer │
                           └─────────┘       └──────────┘
```

#### Go语言实现

```go
package bridge

// Implementation 实现接口
type Implementation interface {
    OperationImpl() string
}

// Abstraction 抽象
type Abstraction struct {
    implementation Implementation
}

func NewAbstraction(impl Implementation) *Abstraction {
    return &Abstraction{implementation: impl}
}

func (a *Abstraction) Operation() string {
    return a.implementation.OperationImpl()
}

// RefinedAbstraction 扩展抽象
type RefinedAbstraction struct {
    Abstraction
}

func (r *RefinedAbstraction) ExtendedOperation() string {
    return "Extended: " + r.Operation()
}

// ConcreteImplementationA 具体实现A
type ConcreteImplementationA struct{}

func (c *ConcreteImplementationA) OperationImpl() string {
    return "ConcreteImplementationA"
}

// ConcreteImplementationB 具体实现B
type ConcreteImplementationB struct{}

func (c *ConcreteImplementationB) OperationImpl() string {
    return "ConcreteImplementationB"
}
```

#### 完整示例：图形渲染系统

```go
package main

import "fmt"

// Renderer 渲染器接口（实现部分）
type Renderer interface {
    RenderCircle(radius float64, x, y float64) string
    RenderSquare(side float64, x, y float64) string
    RenderTriangle(base, height float64, x, y float64) string
}

// ============ 具体渲染器 ============

// SVGRenderer SVG渲染器
type SVGRenderer struct{}

func (s *SVGRenderer) RenderCircle(radius float64, x, y float64) string {
    return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" />`, x, y, radius)
}

func (s *SVGRenderer) RenderSquare(side float64, x, y float64) string {
    return fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" />`,
        x, y, side, side)
}

func (s *SVGRenderer) RenderTriangle(base, height float64, x, y float64) string {
    return fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" />`,
        x, y+height, x+base/2, y, x+base, y+height)
}

// CanvasRenderer Canvas渲染器
type CanvasRenderer struct{}

func (c *CanvasRenderer) RenderCircle(radius float64, x, y float64) string {
    return fmt.Sprintf("ctx.arc(%.1f, %.1f, %.1f, 0, 2 * Math.PI);", x, y, radius)
}

func (c *CanvasRenderer) RenderSquare(side float64, x, y float64) string {
    return fmt.Sprintf("ctx.rect(%.1f, %.1f, %.1f, %.1f);", x, y, side, side)
}

func (c *CanvasRenderer) RenderTriangle(base, height float64, x, y float64) string {
    return fmt.Sprintf("ctx.moveTo(%.1f, %.1f); ctx.lineTo(%.1f, %.1f); ctx.lineTo(%.1f, %.1f);",
        x, y+height, x+base/2, y, x+base, y+height)
}

// OpenGLRenderer OpenGL渲染器
type OpenGLRenderer struct{}

func (o *OpenGLRenderer) RenderCircle(radius float64, x, y float64) string {
    return fmt.Sprintf("glDrawCircle(%.1f, %.1f, %.1f);", x, y, radius)
}

func (o *OpenGLRenderer) RenderSquare(side float64, x, y float64) string {
    return fmt.Sprintf("glDrawRect(%.1f, %.1f, %.1f, %.1f);", x, y, side, side)
}

func (o *OpenGLRenderer) RenderTriangle(base, height float64, x, y float64) string {
    return fmt.Sprintf("glDrawTriangle(%.1f, %.1f, %.1f, %.1f, %.1f, %.1f);",
        x, y+height, x+base/2, y, x+base, y+height)
}

// ============ 形状（抽象部分） ============

// Shape 形状接口
type Shape interface {
    Draw() string
    GetArea() float64
}

// BaseShape 基础形状
type BaseShape struct {
    renderer Renderer
    x, y     float64
}

func NewBaseShape(renderer Renderer, x, y float64) *BaseShape {
    return &BaseShape{renderer: renderer, x: x, y: y}
}

// Circle 圆形
type Circle struct {
    *BaseShape
    radius float64
}

func NewCircle(renderer Renderer, radius, x, y float64) *Circle {
    return &Circle{
        BaseShape: NewBaseShape(renderer, x, y),
        radius:    radius,
    }
}

func (c *Circle) Draw() string {
    return c.renderer.RenderCircle(c.radius, c.x, c.y)
}

func (c *Circle) GetArea() float64 {
    return 3.14159 * c.radius * c.radius
}

// Square 正方形
type Square struct {
    *BaseShape
    side float64
}

func NewSquare(renderer Renderer, side, x, y float64) *Square {
    return &Square{
        BaseShape: NewBaseShape(renderer, x, y),
        side:      side,
    }
}

func (s *Square) Draw() string {
    return s.renderer.RenderSquare(s.side, s.x, s.y)
}

func (s *Square) GetArea() float64 {
    return s.side * s.side
}

// Triangle 三角形
type Triangle struct {
    *BaseShape
    base, height float64
}

func NewTriangle(renderer Renderer, base, height, x, y float64) *Triangle {
    return &Triangle{
        BaseShape: NewBaseShape(renderer, x, y),
        base:      base,
        height:    height,
    }
}

func (t *Triangle) Draw() string {
    return t.renderer.RenderTriangle(t.base, t.height, t.x, t.y)
}

func (t *Triangle) GetArea() float64 {
    return 0.5 * t.base * t.height
}

// ============ 高级形状 ============

// ColoredShape 带颜色的形状（装饰器风格，但使用桥接）
type ColoredShape struct {
    shape Shape
    color string
}

func NewColoredShape(shape Shape, color string) *ColoredShape {
    return &ColoredShape{shape: shape, color: color}
}

func (c *ColoredShape) Draw() string {
    return fmt.Sprintf("<!-- Color: %s -->\n%s", c.color, c.shape.Draw())
}

func (c *ColoredShape) GetArea() float64 {
    return c.shape.GetArea()
}

func main() {
    // 创建不同渲染器
    svg := &SVGRenderer{}
    canvas := &CanvasRenderer{}
    opengl := &OpenGLRenderer{}

    fmt.Println("========== SVG 渲染 ==========")
    circle1 := NewCircle(svg, 50, 100, 100)
    square1 := NewSquare(svg, 80, 200, 200)
    fmt.Println(circle1.Draw())
    fmt.Println(square1.Draw())

    fmt.Println("\n========== Canvas 渲染 ==========")
    circle2 := NewCircle(canvas, 50, 100, 100)
    square2 := NewSquare(canvas, 80, 200, 200)
    fmt.Println(circle2.Draw())
    fmt.Println(square2.Draw())

    fmt.Println("\n========== OpenGL 渲染 ==========")
    circle3 := NewCircle(opengl, 50, 100, 100)
    triangle := NewTriangle(opengl, 100, 80, 300, 300)
    fmt.Println(circle3.Draw())
    fmt.Println(triangle.Draw())

    // 运行时切换渲染器
    fmt.Println("\n========== 运行时切换 ==========")
    shapes := []Shape{
        NewCircle(svg, 30, 50, 50),
        NewSquare(canvas, 60, 150, 150),
        NewTriangle(opengl, 90, 70, 250, 250),
    }

    for _, shape := range shapes {
        fmt.Printf("Area: %.2f\n", shape.GetArea())
        fmt.Println(shape.Draw())
        fmt.Println()
    }

    // 组合使用
    fmt.Println("========== 带颜色的形状 ==========")
    coloredCircle := NewColoredShape(NewCircle(svg, 40, 100, 100), "red")
    fmt.Println(coloredCircle.Draw())
}
```

#### 反例说明

**错误1：使用继承而非组合**

```go
// 错误：类爆炸
type CircleSVG struct{}      // 2D × 3Renderer = 6类
type CircleCanvas struct{}
type SquareSVG struct{}
type SquareCanvas struct{}
// ...
```

**错误2：抽象与实现耦合**

```go
// 错误：直接依赖具体实现
type Circle struct {
    renderer *SVGRenderer  // 应该依赖接口
}
```

**错误3：接口过于庞大**

```go
// 错误：接口包含不相关方法
type Renderer interface {
    RenderCircle()
    RenderSquare()
    RenderTriangle()
    SaveToFile()    // 不属于渲染职责
    Print()         // 不属于渲染职责
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 分离抽象和实现 | 增加设计复杂度 |
| 支持多维度扩展 | 需要正确识别维度 |
| 符合开闭原则 | 对简单系统过度设计 |
| 运行时切换实现 | 增加间接层 |

**适用场景**:

- 图形渲染系统
- 跨平台UI框架
- 数据库驱动
- 消息队列实现

**Go语言特化**:

- 接口作为桥接的天然选择
- 嵌入简化抽象层实现
- 函数类型作为轻量级实现

---

### 8. 组合模式（Composite）

#### 概念定义

组合模式将对象组合成树形结构以表示"部分-整体"的层次结构。组合模式使得用户对单个对象和组合对象的使用具有一致性。

#### 意图与动机

- **树形结构**: 表示层次结构数据
- **统一接口**: 一致处理叶子和容器
- **递归组合**: 支持任意深度的嵌套

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                     Component                               │◄──── 组件接口
│  + Operation()                                              │
│  + Add(Component)                                           │
│  + Remove(Component)                                        │
│  + GetChild(int)                                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│         Leaf            │   │       Composite         │
│  + Operation()          │   │  - children: []Component│
│                         │   ├─────────────────────────┤
│                         │   │  + Operation()          │
│                         │   │  + Add(Component)       │
│                         │   │  + Remove(Component)    │
│                         │   │  + GetChild(int)        │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

```go
package composite

// Component 组件接口
type Component interface {
    Operation() string
    GetName() string
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

func (c *Composite) Operation() string {
    result := c.name + "\n"
    for _, child := range c.children {
        result += "  " + child.Operation() + "\n"
    }
    return result
}

func (c *Composite) GetName() string {
    return c.name
}

func (c *Composite) Add(child Component) {
    c.children = append(c.children, child)
}

func (c *Composite) Remove(child Component) {
    for i, ch := range c.children {
        if ch == child {
            c.children = append(c.children[:i], c.children[i+1:]...)
            return
        }
    }
}

func (c *Composite) GetChild(index int) Component {
    if index >= 0 && index < len(c.children) {
        return c.children[index]
    }
    return nil
}

// Leaf 叶子节点
type Leaf struct {
    name string
}

func NewLeaf(name string) *Leaf {
    return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
    return l.name
}

func (l *Leaf) GetName() string {
    return l.name
}
```

#### 完整示例：文件系统

```go
package main

import (
    "fmt"
    "strings"
)

// FileSystemComponent 文件系统组件接口
type FileSystemComponent interface {
    GetName() string
    GetSize() int
    Display(indent string) string
    Search(keyword string) []FileSystemComponent
}

// File 文件（叶子节点）
type File struct {
    name string
    size int
}

func NewFile(name string, size int) *File {
    return &File{name: name, size: size}
}

func (f *File) GetName() string {
    return f.name
}

func (f *File) GetSize() int {
    return f.size
}

func (f *File) Display(indent string) string {
    return fmt.Sprintf("%s📄 %s (%d bytes)", indent, f.name, f.size)
}

func (f *File) Search(keyword string) []FileSystemComponent {
    if strings.Contains(strings.ToLower(f.name), strings.ToLower(keyword)) {
        return []FileSystemComponent{f}
    }
    return nil
}

// Directory 目录（组合节点）
type Directory struct {
    name     string
    children []FileSystemComponent
}

func NewDirectory(name string) *Directory {
    return &Directory{
        name:     name,
        children: make([]FileSystemComponent, 0),
    }
}

func (d *Directory) GetName() string {
    return d.name
}

func (d *Directory) GetSize() int {
    total := 0
    for _, child := range d.children {
        total += child.GetSize()
    }
    return total
}

func (d *Directory) Display(indent string) string {
    result := fmt.Sprintf("%s📁 %s/ (%d bytes)", indent, d.name, d.GetSize())
    for _, child := range d.children {
        result += "\n" + child.Display(indent+"  ")
    }
    return result
}

func (d *Directory) Search(keyword string) []FileSystemComponent {
    var results []FileSystemComponent

    if strings.Contains(strings.ToLower(d.name), strings.ToLower(keyword)) {
        results = append(results, d)
    }

    for _, child := range d.children {
        results = append(results, child.Search(keyword)...)
    }

    return results
}

func (d *Directory) Add(component FileSystemComponent) {
    d.children = append(d.children, component)
}

func (d *Directory) Remove(component FileSystemComponent) {
    for i, child := range d.children {
        if child == component {
            d.children = append(d.children[:i], d.children[i+1:]...)
            return
        }
    }
}

func (d *Directory) GetChildren() []FileSystemComponent {
    return d.children
}

// ============ 高级功能 ============

// SizeCalculator 大小计算器
type SizeCalculator struct {
    totalSize int
    fileCount int
    dirCount  int
}

func (sc *SizeCalculator) Calculate(component FileSystemComponent) {
    switch c := component.(type) {
    case *File:
        sc.totalSize += c.GetSize()
        sc.fileCount++
    case *Directory:
        sc.dirCount++
        for _, child := range c.GetChildren() {
            sc.Calculate(child)
        }
    }
}

func (sc *SizeCalculator) Report() string {
    return fmt.Sprintf("Total Size: %d bytes, Files: %d, Directories: %d",
        sc.totalSize, sc.fileCount, sc.dirCount)
}

func main() {
    // 构建文件系统
    root := NewDirectory("root")

    // 创建文档目录
    documents := NewDirectory("documents")
    documents.Add(NewFile("resume.pdf", 1024))
    documents.Add(NewFile("cover_letter.docx", 512))

    // 创建图片目录
    pictures := NewDirectory("pictures")
    pictures.Add(NewFile("vacation.jpg", 2048))
    pictures.Add(NewFile("profile.png", 1024))

    // 创建子目录
    screenshots := NewDirectory("screenshots")
    screenshots.Add(NewFile("error.png", 512))
    screenshots.Add(NewFile("success.png", 512))
    pictures.Add(screenshots)

    // 创建代码目录
    code := NewDirectory("code")
    code.Add(NewFile("main.go", 256))
    code.Add(NewFile("utils.go", 128))

    // 组装目录结构
    root.Add(documents)
    root.Add(pictures)
    root.Add(code)
    root.Add(NewFile("README.md", 256))

    // 显示文件系统
    fmt.Println("========== 文件系统结构 ==========")
    fmt.Println(root.Display(""))

    // 搜索文件
    fmt.Println("\n========== 搜索 'pic' ==========")
    results := root.Search("pic")
    for _, r := range results {
        fmt.Printf("Found: %s\n", r.GetName())
    }

    // 搜索文件
    fmt.Println("\n========== 搜索 '.go' ==========")
    results = root.Search(".go")
    for _, r := range results {
        fmt.Printf("Found: %s\n", r.GetName())
    }

    // 统计信息
    fmt.Println("\n========== 统计信息 ==========")
    calculator := &SizeCalculator{}
    calculator.Calculate(root)
    fmt.Println(calculator.Report())

    // 统一处理
    fmt.Println("\n========== 统一处理所有组件 ==========")
    components := []FileSystemComponent{
        NewFile("single.txt", 100),
        NewDirectory("empty_dir"),
        root,
    }

    for _, c := range components {
        fmt.Printf("%s: %d bytes\n", c.GetName(), c.GetSize())
    }
}
```

#### 反例说明

**错误1：叶子节点实现Add/Remove**

```go
// 错误：叶子不应该有这些方法
type Leaf struct{}

func (l *Leaf) Add(c Component) {
    panic("Cannot add to leaf")  // 不应该实现
}
```

**错误2：类型判断破坏透明性**

```go
// 错误：不应该需要类型判断
func Process(c Component) {
    if composite, ok := c.(*Composite); ok {
        // 特殊处理组合节点
    } else {
        // 处理叶子节点
    }
}
```

**错误3：循环引用**

```go
// 错误：循环引用会导致无限递归
parent := NewComposite("parent")
child := NewComposite("child")
parent.Add(child)
child.Add(parent)  // 循环引用！
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 统一处理单个和组合对象 | 限制组件类型（必须实现所有接口方法） |
| 易于添加新组件类型 | 设计复杂 |
| 支持递归结构 | 类型安全降低 |
| 符合开闭原则 | 调试困难 |

**适用场景**:

- 文件系统
- UI组件树
- 组织架构
- 菜单系统
- 表达式解析

**Go语言特化**:

- 类型断言处理特定类型
- 接口嵌入简化实现
- 空接口作为通用容器

---

### 9. 装饰器模式（Decorator）

#### 概念定义

装饰器模式动态地给一个对象添加额外的职责。装饰器模式提供了一种比继承更灵活的替代方案来扩展功能。

#### 意图与动机

- **动态扩展**: 运行时添加功能
- **单一职责**: 每个装饰器只负责一个功能
- **组合功能**: 通过组合实现多种功能叠加

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                     Component                               │◄──── 组件接口
│  + Operation()                                              │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│    ConcreteComponent    │   │      Decorator          │
│  + Operation()          │   │  - component: Component │
│                         │   ├─────────────────────────┤
│                         │   │  + Operation()          │
│                         │   │    component.Operation()│
└─────────────────────────┘   └───────────┬─────────────┘
                                          │
                              ┌───────────┴───────────┐
                              │                       │
                              ▼                       ▼
                    ┌─────────────────┐   ┌─────────────────┐
                    │ConcreteDecoratorA│   │ConcreteDecoratorB│
                    │ + Operation()   │   │ + Operation()   │
                    │   addedBehavior()│   │   addedState    │
                    └─────────────────┘   └─────────────────┘
```

#### Go语言实现

```go
package decorator

// Component 组件接口
type Component interface {
    Operation() string
}

// ConcreteComponent 具体组件
type ConcreteComponent struct{}

func (c *ConcreteComponent) Operation() string {
    return "ConcreteComponent"
}

// Decorator 装饰器基类
type Decorator struct {
    component Component
}

func NewDecorator(component Component) *Decorator {
    return &Decorator{component: component}
}

func (d *Decorator) Operation() string {
    if d.component != nil {
        return d.component.Operation()
    }
    return ""
}

// ConcreteDecoratorA 具体装饰器A
type ConcreteDecoratorA struct {
    Decorator
}

func NewConcreteDecoratorA(component Component) *ConcreteDecoratorA {
    return &ConcreteDecoratorA{Decorator: *NewDecorator(component)}
}

func (d *ConcreteDecoratorA) Operation() string {
    return "DecoratorA(" + d.Decorator.Operation() + ")"
}

// ConcreteDecoratorB 具体装饰器B
type ConcreteDecoratorB struct {
    Decorator
    addedState string
}

func NewConcreteDecoratorB(component Component) *ConcreteDecoratorB {
    return &ConcreteDecoratorB{
        Decorator:  *NewDecorator(component),
        addedState: "B_State",
    }
}

func (d *ConcreteDecoratorB) Operation() string {
    return "DecoratorB[" + d.addedState + "](" + d.Decorator.Operation() + ")"
}
```

#### 完整示例：HTTP中间件链

```go
package main

import (
    "fmt"
    "time"
)

// Handler HTTP处理器接口
type Handler interface {
    ServeHTTP(request string) string
}

// HandlerFunc 函数类型适配器
type HandlerFunc func(request string) string

func (f HandlerFunc) ServeHTTP(request string) string {
    return f(request)
}

// BaseHandler 基础处理器
type BaseHandler struct{}

func (h *BaseHandler) ServeHTTP(request string) string {
    return fmt.Sprintf("BaseHandler processed: %s", request)
}

// Middleware 中间件类型
type Middleware func(Handler) Handler

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        start := time.Now()
        fmt.Printf("[LOG] Request: %s\n", request)

        response := next.ServeHTTP(request)

        duration := time.Since(start)
        fmt.Printf("[LOG] Response: %s (took %v)\n", response, duration)

        return response
    })
}

// AuthMiddleware 认证中间件
func AuthMiddleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        // 检查认证
        if !isAuthenticated(request) {
            return "Error: Unauthorized"
        }
        fmt.Println("[AUTH] Authentication successful")
        return next.ServeHTTP(request)
    })
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
    requests map[string]int
    limit    int
}

func NewRateLimitMiddleware(limit int) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        requests: make(map[string]int),
        limit:    limit,
    }
}

func (r *RateLimitMiddleware) Middleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        clientID := extractClientID(request)

        if r.requests[clientID] >= r.limit {
            return "Error: Rate limit exceeded"
        }

        r.requests[clientID]++
        fmt.Printf("[RATE] Request count for %s: %d/%d\n",
            clientID, r.requests[clientID], r.limit)

        return next.ServeHTTP(request)
    })
}

// CacheMiddleware 缓存中间件
type CacheMiddleware struct {
    cache map[string]string
}

func NewCacheMiddleware() *CacheMiddleware {
    return &CacheMiddleware{
        cache: make(map[string]string),
    }
}

func (c *CacheMiddleware) Middleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        // 检查缓存
        if cached, ok := c.cache[request]; ok {
            fmt.Printf("[CACHE] Cache hit for: %s\n", request)
            return cached
        }

        // 执行请求
        response := next.ServeHTTP(request)

        // 缓存结果
        c.cache[request] = response
        fmt.Printf("[CACHE] Cached response for: %s\n", request)

        return response
    })
}

// MetricsMiddleware 指标中间件
func MetricsMiddleware(next Handler) Handler {
    requestCount := 0

    return HandlerFunc(func(request string) string {
        requestCount++
        fmt.Printf("[METRICS] Total requests: %d\n", requestCount)
        return next.ServeHTTP(request)
    })
}

// Chain 中间件链
type Chain struct {
    handler     Handler
    middlewares []Middleware
}

func NewChain(handler Handler) *Chain {
    return &Chain{handler: handler}
}

func (c *Chain) Use(middlewares ...Middleware) *Chain {
    c.middlewares = append(c.middlewares, middlewares...)
    return c
}

func (c *Chain) Then() Handler {
    handler := c.handler
    // 倒序应用中间件（先添加的先执行）
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        handler = c.middlewares[i](handler)
    }
    return handler
}

// 辅助函数
func isAuthenticated(request string) bool {
    return len(request) > 10  // 简单模拟
}

func extractClientID(request string) string {
    if len(request) > 5 {
        return request[:5]
    }
    return "unknown"
}

func main() {
    // 基础处理器
    baseHandler := &BaseHandler{}

    fmt.Println("========== 单个中间件 ==========")
    loggedHandler := LoggingMiddleware(baseHandler)
    result := loggedHandler.ServeHTTP("GET /api/users")
    fmt.Printf("Result: %s\n\n", result)

    fmt.Println("========== 多个中间件组合 ==========")
    // 手动组合
    handler := LoggingMiddleware(
        AuthMiddleware(
            baseHandler,
        ),
    )
    result = handler.ServeHTTP("GET /api/users authenticated")
    fmt.Printf("Result: %s\n\n", result)

    fmt.Println("========== 使用Chain ==========")
    rateLimiter := NewRateLimitMiddleware(3)
    cache := NewCacheMiddleware()

    chain := NewChain(baseHandler).
        Use(MetricsMiddleware).
        Use(LoggingMiddleware).
        Use(rateLimiter.Middleware).
        Use(cache.Middleware).
        Use(AuthMiddleware)

    finalHandler := chain.Then()

    // 发送多个请求
    requests := []string{
        "GET /api/users authenticated",
        "GET /api/users authenticated",  // 缓存命中
        "GET /api/products authenticated",
        "GET /api/orders authenticated",
        "GET /api/users authenticated",  // 可能触发限流
    }

    for _, req := range requests {
        fmt.Printf("\n--- Request: %s ---\n", req)
        response := finalHandler.ServeHTTP(req)
        fmt.Printf("Final Response: %s\n", response)
    }

    // 展示装饰器叠加效果
    fmt.Println("\n========== 装饰器叠加效果 ==========")
    simpleHandler := HandlerFunc(func(r string) string {
        return "Response"
    })

    decorated := LoggingMiddleware(
        MetricsMiddleware(
            AuthMiddleware(
                simpleHandler,
            ),
        ),
    )

    decorated.ServeHTTP("test request")
}
```

#### 反例说明

**错误1：过多装饰器嵌套**

```go
// 错误：过度嵌套难以阅读
result := A(B(C(D(E(F(base))))))
```

**错误2：装饰器修改接口**

```go
// 错误：装饰器应该保持接口一致
type Decorator struct {
    component Component
}

func (d *Decorator) NewMethod() {  // 不要添加新方法
    // ...
}
```

**错误3：装饰器依赖具体类型**

```go
// 错误：应该依赖接口
type Decorator struct {
    component *ConcreteComponent  // 应该使用 Component 接口
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 运行时扩展功能 | 可能产生大量小对象 |
| 单一职责原则 | 调试困难（多层包装） |
| 灵活组合功能 | 过度使用导致复杂 |
| 符合开闭原则 | 性能开销 |

**适用场景**:

- HTTP中间件
- I/O流包装
- GUI组件装饰
- 数据验证链
- 缓存层

**Go语言特化**:

- 函数类型作为轻量级装饰器
- `http.Handler`接口标准实践
- 中间件链式组合模式

---

### 10. 外观模式（Facade）

#### 概念定义

外观模式为子系统中的一组接口提供一个统一的接口。外观定义了一个高层接口，使子系统更容易使用。

#### 意图与动机

- **简化接口**: 提供简单的统一接口
- **解耦**: 隔离客户端和子系统
- **分层**: 构建清晰的系统层次

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Client                                 │
│                    (使用Facade)                             │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                       Facade                                │◄──── 外观
│  - subsystem1: Subsystem1                                   │
│  - subsystem2: Subsystem2                                   │
│  - subsystem3: Subsystem3                                   │
├─────────────────────────────────────────────────────────────┤
│  + SimpleOperation()                                        │
│    subsystem1.OperationA()                                  │
│    subsystem2.OperationB()                                  │
│    subsystem3.OperationC()                                  │
└─────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  Subsystem1   │  │  Subsystem2   │  │  Subsystem3   │
│ + OperationA()│  │ + OperationB()│  │ + OperationC()│
│ + OperationX()│  │ + OperationY()│  │ + OperationZ()│
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

```go
package facade

import "fmt"

// Subsystem1 子系统1
type Subsystem1 struct{}

func (s *Subsystem1) Operation1() string {
    return "Subsystem1: Ready!\n"
}

func (s *Subsystem1) OperationN() string {
    return "Subsystem1: Go!\n"
}

// Subsystem2 子系统2
type Subsystem2 struct{}

func (s *Subsystem2) Operation1() string {
    return "Subsystem2: Get ready!\n"
}

func (s *Subsystem2) OperationZ() string {
    return "Subsystem2: Fire!\n"
}

// Facade 外观
type Facade struct {
    subsystem1 *Subsystem1
    subsystem2 *Subsystem2
}

func NewFacade() *Facade {
    return &Facade{
        subsystem1: &Subsystem1{},
        subsystem2: &Subsystem2{},
    }
}

func (f *Facade) Operation() string {
    result := "Facade initializes subsystems:\n"
    result += f.subsystem1.Operation1()
    result += f.subsystem2.Operation1()
    result += "Facade orders subsystems to perform the action:\n"
    result += f.subsystem1.OperationN()
    result += f.subsystem2.OperationZ()
    return result
}
```

#### 完整示例：电商订单系统外观

```go
package main

import (
    "fmt"
    "time"
)

// ============ 子系统 ============

// InventorySystem 库存系统
type InventorySystem struct{}

func (i *InventorySystem) CheckStock(productID string, quantity int) bool {
    fmt.Printf("[库存] 检查产品 %s 的库存，数量: %d\n", productID, quantity)
    return true  // 模拟有库存
}

func (i *InventorySystem) ReserveStock(productID string, quantity int) error {
    fmt.Printf("[库存] 预留产品 %s，数量: %d\n", productID, quantity)
    return nil
}

func (i *InventorySystem) ReleaseStock(productID string, quantity int) {
    fmt.Printf("[库存] 释放产品 %s，数量: %d\n", productID, quantity)
}

// PaymentSystem 支付系统
type PaymentSystem struct{}

func (p *PaymentSystem) ProcessPayment(userID string, amount float64) (string, error) {
    fmt.Printf("[支付] 处理用户 %s 的支付，金额: %.2f\n", userID, amount)
    return fmt.Sprintf("PAY%d", time.Now().Unix()), nil
}

func (p *PaymentSystem) RefundPayment(paymentID string) error {
    fmt.Printf("[支付] 退款: %s\n", paymentID)
    return nil
}

// ShippingSystem 物流系统
type ShippingSystem struct{}

func (s *ShippingSystem) CreateShipment(orderID string, address string) (string, error) {
    fmt.Printf("[物流] 创建订单 %s 的配送，地址: %s\n", orderID, address)
    return fmt.Sprintf("SHIP%d", time.Now().Unix()), nil
}

func (s *ShippingSystem) CancelShipment(shipmentID string) error {
    fmt.Printf("[物流] 取消配送: %s\n", shipmentID)
    return nil
}

// NotificationSystem 通知系统
type NotificationSystem struct{}

func (n *NotificationSystem) SendOrderConfirmation(userID, orderID string) {
    fmt.Printf("[通知] 发送订单确认给用户 %s，订单: %s\n", userID, orderID)
}

func (n *NotificationSystem) SendShippingNotification(userID, shipmentID string) {
    fmt.Printf("[通知] 发送配送通知给用户 %s，配送: %s\n", userID, shipmentID)
}

// UserSystem 用户系统
type UserSystem struct{}

func (u *UserSystem) GetUserAddress(userID string) string {
    fmt.Printf("[用户] 获取用户 %s 的地址\n", userID)
    return "123 Main St, City, Country"
}

func (u *UserSystem) ValidateUser(userID string) bool {
    fmt.Printf("[用户] 验证用户 %s\n", userID)
    return true
}

// ============ 外观 ============

// OrderFacade 订单外观
type OrderFacade struct {
    inventory     *InventorySystem
    payment       *PaymentSystem
    shipping      *ShippingSystem
    notification  *NotificationSystem
    user          *UserSystem
}

func NewOrderFacade() *OrderFacade {
    return &OrderFacade{
        inventory:    &InventorySystem{},
        payment:      &PaymentSystem{},
        shipping:     &ShippingSystem{},
        notification: &NotificationSystem{},
        user:         &UserSystem{},
    }
}

// PlaceOrder 下单（简化复杂流程）
func (o *OrderFacade) PlaceOrder(userID string, items map[string]int, totalAmount float64) (*OrderResult, error) {
    // 1. 验证用户
    if !o.user.ValidateUser(userID) {
        return nil, fmt.Errorf("用户验证失败")
    }

    // 2. 检查库存
    for productID, quantity := range items {
        if !o.inventory.CheckStock(productID, quantity) {
            return nil, fmt.Errorf("产品 %s 库存不足", productID)
        }
    }

    // 3. 预留库存
    for productID, quantity := range items {
        if err := o.inventory.ReserveStock(productID, quantity); err != nil {
            return nil, err
        }
    }

    // 4. 处理支付
    paymentID, err := o.payment.ProcessPayment(userID, totalAmount)
    if err != nil {
        // 回滚：释放库存
        for productID, quantity := range items {
            o.inventory.ReleaseStock(productID, quantity)
        }
        return nil, fmt.Errorf("支付失败: %v", err)
    }

    // 5. 创建订单ID
    orderID := fmt.Sprintf("ORD%d", time.Now().Unix())

    // 6. 创建配送
    address := o.user.GetUserAddress(userID)
    shipmentID, err := o.shipping.CreateShipment(orderID, address)
    if err != nil {
        // 回滚：退款、释放库存
        o.payment.RefundPayment(paymentID)
        for productID, quantity := range items {
            o.inventory.ReleaseStock(productID, quantity)
        }
        return nil, fmt.Errorf("创建配送失败: %v", err)
    }

    // 7. 发送通知
    o.notification.SendOrderConfirmation(userID, orderID)
    o.notification.SendShippingNotification(userID, shipmentID)

    return &OrderResult{
        OrderID:     orderID,
        PaymentID:   paymentID,
        ShipmentID:  shipmentID,
        TotalAmount: totalAmount,
    }, nil
}

// CancelOrder 取消订单
type OrderResult struct {
    OrderID     string
    PaymentID   string
    ShipmentID  string
    TotalAmount float64
}

func (o *OrderFacade) CancelOrder(orderID, paymentID, shipmentID string, items map[string]int) error {
    fmt.Printf("\n[外观] 取消订单 %s\n", orderID)

    // 1. 取消配送
    if err := o.shipping.CancelShipment(shipmentID); err != nil {
        return err
    }

    // 2. 退款
    if err := o.payment.RefundPayment(paymentID); err != nil {
        return err
    }

    // 3. 释放库存
    for productID, quantity := range items {
        o.inventory.ReleaseStock(productID, quantity)
    }

    fmt.Println("[外观] 订单取消成功")
    return nil
}

// GetOrderStatus 获取订单状态（简化查询）
func (o *OrderFacade) GetOrderStatus(orderID, shipmentID string) map[string]string {
    return map[string]string{
        "order_id":    orderID,
        "shipment_id": shipmentID,
        "status":      "processing",
    }
}

func main() {
    // 创建外观
    orderFacade := NewOrderFacade()

    fmt.Println("========== 下单流程 ==========")
    items := map[string]int{
        "PROD001": 2,
        "PROD002": 1,
    }

    result, err := orderFacade.PlaceOrder("USER123", items, 299.99)
    if err != nil {
        fmt.Printf("下单失败: %v\n", err)
    } else {
        fmt.Printf("\n下单成功!\n")
        fmt.Printf("订单ID: %s\n", result.OrderID)
        fmt.Printf("支付ID: %s\n", result.PaymentID)
        fmt.Printf("配送ID: %s\n", result.ShipmentID)
        fmt.Printf("总金额: %.2f\n", result.TotalAmount)
    }

    // 取消订单
    fmt.Println("\n========== 取消订单 ==========")
    if result != nil {
        err := orderFacade.CancelOrder(result.OrderID, result.PaymentID, result.ShipmentID, items)
        if err != nil {
            fmt.Printf("取消失败: %v\n", err)
        }
    }

    // 对比：不使用外观
    fmt.Println("\n========== 不使用外观（复杂代码） ==========")
    fmt.Println("需要手动调用:")
    fmt.Println("1. user.ValidateUser()")
    fmt.Println("2. inventory.CheckStock() for each item")
    fmt.Println("3. inventory.ReserveStock() for each item")
    fmt.Println("4. payment.ProcessPayment()")
    fmt.Println("5. shipping.CreateShipment()")
    fmt.Println("6. notification.SendOrderConfirmation()")
    fmt.Println("7. 处理各种错误和回滚...")
}
```

#### 反例说明

**错误1：外观过于复杂**

```go
// 错误：外观不应该包含所有子系统方法
type Facade struct {
    // ...
}

func (f *Facade) Subsystem1Method1() {  // 不要暴露所有方法
    f.subsystem1.Method1()
}
```

**错误2：外观成为上帝对象**

```go
// 错误：外观不应该包含业务逻辑
type Facade struct{}

func (f *Facade) ComplexBusinessLogic() {  // 应该放在领域层
    // 大量业务逻辑
}
```

**错误3：阻止直接访问子系统**

```go
// 错误：不应该强制只能通过外观访问
type Subsystem struct{}  // 不要设为私有

func (s *Subsystem) Method() {}  // 应该允许直接访问
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 简化接口 | 可能成为上帝对象 |
| 解耦子系统 | 过度简化可能丢失灵活性 |
| 分层清晰 | 额外抽象层 |
| 易于使用 | 需要维护外观 |

**适用场景**:

- 复杂子系统封装
- 分层架构
- 遗留系统包装
- API网关

**Go语言特化**:

- 包级别外观（公开函数）
- 接口组合简化外观
- 上下文传递配置

---

### 11. 享元模式（Flyweight）

#### 概念定义

享元模式运用共享技术有效地支持大量细粒度的对象。享元模式通过共享对象来减少内存使用。

#### 意图与动机

- **内存优化**: 减少大量相似对象的内存占用
- **性能提升**: 减少对象创建开销
- **状态分离**: 区分内部状态和外部状态

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                    FlyweightFactory                         │
│  - flyweights: map[string]Flyweight                         │
├─────────────────────────────────────────────────────────────┤
│  + GetFlyweight(key) Flyweight                              │
│    if not exists:                                           │
│      create new Flyweight                                   │
│    return flyweights[key]                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │ manages
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Flyweight                              │◄──── 享元接口
│  + Operation(extrinsicState)                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│   ConcreteFlyweight     │   │  UnsharedConcrete       │
│  - intrinsicState       │   │    Flyweight            │
│  + Operation()          │   │  - allState             │
│    uses intrinsicState  │   │  + Operation()          │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

```go
package flyweight

// Flyweight 享元接口
type Flyweight interface {
    Operation(extrinsicState string) string
}

// ConcreteFlyweight 具体享元
type ConcreteFlyweight struct {
    intrinsicState string  // 内部状态（共享）
}

func NewConcreteFlyweight(state string) *ConcreteFlyweight {
    return &ConcreteFlyweight{intrinsicState: state}
}

func (f *ConcreteFlyweight) Operation(extrinsicState string) string {
    return fmt.Sprintf("Intrinsic: %s, Extrinsic: %s",
        f.intrinsicState, extrinsicState)
}

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
    flyweights map[string]Flyweight
}

func NewFlyweightFactory() *FlyweightFactory {
    return &FlyweightFactory{
        flyweights: make(map[string]Flyweight),
    }
}

func (f *FlyweightFactory) GetFlyweight(key string) Flyweight {
    if flyweight, exists := f.flyweights[key]; exists {
        return flyweight
    }

    // 创建新的享元
    flyweight := NewConcreteFlyweight(key)
    f.flyweights[key] = flyweight
    return flyweight
}

func (f *FlyweightFactory) Count() int {
    return len(f.flyweights)
}
```

#### 完整示例：游戏角色渲染系统

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// ============ 享元部分（共享数据） ============

// TreeType 树的类型（享元）
type TreeType struct {
    name       string
    color      string
    texture    string
    meshData   []byte  // 大量几何数据
}

func NewTreeType(name, color, texture string, meshSize int) *TreeType {
    return &TreeType{
        name:     name,
        color:    color,
        texture:  texture,
        meshData: make([]byte, meshSize), // 模拟大量数据
    }
}

func (t *TreeType) GetName() string {
    return t.name
}

func (t *TreeType) GetColor() string {
    return t.color
}

func (t *TreeType) GetMemorySize() int {
    return len(t.meshData) + len(t.name) + len(t.color) + len(t.texture)
}

// TreeTypeFactory 树类型工厂
type TreeTypeFactory struct {
    treeTypes map[string]*TreeType
}

func NewTreeTypeFactory() *TreeTypeFactory {
    return &TreeTypeFactory{
        treeTypes: make(map[string]*TreeType),
    }
}

func (f *TreeTypeFactory) GetTreeType(name, color, texture string, meshSize int) *TreeType {
    key := fmt.Sprintf("%s_%s_%s", name, color, texture)

    if treeType, exists := f.treeTypes[key]; exists {
        fmt.Printf("[享元] 复用树类型: %s\n", key)
        return treeType
    }

    fmt.Printf("[享元] 创建新树类型: %s\n", key)
    treeType := NewTreeType(name, color, texture, meshSize)
    f.treeTypes[key] = treeType
    return treeType
}

func (f *TreeTypeFactory) GetTypeCount() int {
    return len(f.treeTypes)
}

func (f *TreeTypeFactory) GetTotalMemory() int {
    total := 0
    for _, treeType := range f.treeTypes {
        total += treeType.GetMemorySize()
    }
    return total
}

// ============ 上下文部分（外部状态） ============

// Tree 树实例（包含外部状态）
type Tree struct {
    x, y    int           // 位置（外部状态）
    treeType *TreeType    // 引用享元（内部状态）
}

func NewTree(x, y int, treeType *TreeType) *Tree {
    return &Tree{
        x:        x,
        y:        y,
        treeType: treeType,
    }
}

func (t *Tree) Draw() string {
    return fmt.Sprintf("在(%d, %d)绘制%s色的%s",
        t.x, t.y, t.treeType.GetColor(), t.treeType.GetName())
}

func (t *Tree) GetPosition() (int, int) {
    return t.x, t.y
}

// ============ 森林管理器 ============

// Forest 森林
type Forest struct {
    trees     []*Tree
    factory   *TreeTypeFactory
}

func NewForest() *Forest {
    return &Forest{
        trees:   make([]*Tree, 0),
        factory: NewTreeTypeFactory(),
    }
}

func (f *Forest) PlantTree(x, y int, name, color, texture string, meshSize int) {
    treeType := f.factory.GetTreeType(name, color, texture, meshSize)
    tree := NewTree(x, y, treeType)
    f.trees = append(f.trees, tree)
}

func (f *Forest) Draw() {
    for _, tree := range f.trees {
        fmt.Println(tree.Draw())
    }
}

func (f *Forest) GetTreeCount() int {
    return len(f.trees)
}

func (f *Forest) GetMemoryReport() string {
    uniqueTypes := f.factory.GetTypeCount()
    sharedMemory := f.factory.GetTotalMemory()
    totalTrees := len(f.trees)

    // 假设每个树实例的位置信息占用16字节
    contextMemory := totalTrees * 16

    // 如果不使用享元，需要的内存
    nonSharedMemory := totalTrees * (1024 + 16) // mesh + position

    savedMemory := nonSharedMemory - (sharedMemory + contextMemory)
    savedPercent := float64(savedMemory) / float64(nonSharedMemory) * 100

    return fmt.Sprintf(
        "内存报告:\n"+
        "  树实例数量: %d\n"+
        "  唯一树类型: %d\n"+
        "  共享数据内存: %d bytes\n"+
        "  上下文内存: %d bytes\n"+
        "  总内存使用: %d bytes\n"+
        "  不使用享元需要: %d bytes\n"+
        "  节省内存: %d bytes (%.1f%%)",
        totalTrees, uniqueTypes, sharedMemory, contextMemory,
        sharedMemory+contextMemory, nonSharedMemory,
        savedMemory, savedPercent,
    )
}

// ============ 其他享元示例：字符渲染 ============

// CharacterStyle 字符样式（享元）
type CharacterStyle struct {
    font     string
    size     int
    color    string
    bold     bool
    italic   bool
}

func NewCharacterStyle(font string, size int, color string, bold, italic bool) *CharacterStyle {
    return &CharacterStyle{
        font:   font,
        size:   size,
        color:  color,
        bold:   bold,
        italic: italic,
    }
}

func (c *CharacterStyle) String() string {
    return fmt.Sprintf("%s-%d-%s", c.font, c.size, c.color)
}

// CharacterStyleFactory 字符样式工厂
type CharacterStyleFactory struct {
    styles map[string]*CharacterStyle
}

func NewCharacterStyleFactory() *CharacterStyleFactory {
    return &CharacterStyleFactory{
        styles: make(map[string]*CharacterStyle),
    }
}

func (f *CharacterStyleFactory) GetStyle(font string, size int, color string, bold, italic bool) *CharacterStyle {
    key := fmt.Sprintf("%s_%d_%s_%v_%v", font, size, color, bold, italic)

    if style, exists := f.styles[key]; exists {
        return style
    }

    style := NewCharacterStyle(font, size, color, bold, italic)
    f.styles[key] = style
    return style
}

// Character 字符（包含外部状态）
type Character struct {
    char   rune
    x, y   int
    style  *CharacterStyle
}

func NewCharacter(char rune, x, y int, style *CharacterStyle) *Character {
    return &Character{
        char:  char,
        x:     x,
        y:     y,
        style: style,
    }
}

func (c *Character) Render() string {
    style := ""
    if c.style.bold {
        style += "Bold "
    }
    if c.style.italic {
        style += "Italic "
    }
    return fmt.Sprintf("'%c' at (%d,%d) with %s%dpt %s %s",
        c.char, c.x, c.y, style, c.style.size, c.style.font, c.style.color)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    fmt.Println("========== 森林示例 ==========")
    forest := NewForest()

    // 种植大量树，但只有几种类型
    treeTypes := []struct {
        name    string
        color   string
        texture string
    }{
        {"橡树", "绿色", "粗糙"},
        {"松树", "深绿", "光滑"},
        {"枫树", "红色", "中等"},
        {"桦树", "白色", "光滑"},
    }

    // 种植100棵树
    for i := 0; i < 100; i++ {
        treeType := treeTypes[rand.Intn(len(treeTypes))]
        x := rand.Intn(1000)
        y := rand.Intn(1000)
        forest.PlantTree(x, y, treeType.name, treeType.color, treeType.texture, 1024)
    }

    // 显示部分树
    fmt.Println("\n部分树的绘制：")
    for i := 0; i < 5; i++ {
        fmt.Println(forest.trees[i].Draw())
    }

    // 内存报告
    fmt.Println("\n" + forest.GetMemoryReport())

    // 字符渲染示例
    fmt.Println("\n========== 字符渲染示例 ==========")
    styleFactory := NewCharacterStyleFactory()

    // 创建文档
    var characters []*Character
    text := "Hello, World! This is a test."
    x, y := 0, 0

    // 普通文本
    normalStyle := styleFactory.GetStyle("Arial", 12, "black", false, false)
    for _, char := range text {
        if char == ' ' {
            x += 8
            continue
        }
        characters = append(characters, NewCharacter(char, x, y, normalStyle))
        x += 8
    }

    // 标题（粗体）
    x, y = 0, 20
    titleStyle := styleFactory.GetStyle("Arial", 16, "black", true, false)
    for _, char := range "Title" {
        characters = append(characters, NewCharacter(char, x, y, titleStyle))
        x += 10
    }

    // 强调（斜体红色）
    x, y = 0, 40
    emphasisStyle := styleFactory.GetStyle("Arial", 12, "red", false, true)
    for _, char := range "Important!" {
        characters = append(characters, NewCharacter(char, x, y, emphasisStyle))
        x += 8
    }

    // 显示字符
    fmt.Println("文档内容：")
    for _, char := range characters {
        fmt.Println(char.Render())
    }

    fmt.Printf("\n字符总数: %d\n", len(characters))
    fmt.Printf("唯一样式数: %d\n", len(styleFactory.styles))
}
```

#### 反例说明

**错误1：混淆内部和外部状态**

```go
// 错误：位置应该是外部状态
type TreeType struct {
    x, y int  // 错误：位置不应该在享元中
}
```

**错误2：享元可变**

```go
// 错误：享元应该是不可变的
type TreeType struct {
    color string
}

func (t *TreeType) SetColor(c string) {  // 不要提供修改方法
    t.color = c
}
```

**错误3：过度使用享元**

```go
// 错误：对象数量不多时不需要享元
if objectCount < 100 {  // 少量对象不需要享元
    // 直接使用普通对象
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 大幅减少内存使用 | 代码复杂度增加 |
| 提升性能 | 需要分离内外状态 |
| 集中管理对象 | 运行时开销 |
| 符合单一职责 | 线程安全需要额外处理 |

**适用场景**:

- 游戏对象渲染
- 文本编辑器字符样式
- 地图瓦片
- 数据库连接池
- 字符串常量池

**Go语言特化**:

- `string`类型本身就是享元
- 使用`sync.Pool`实现对象池
- map作为享元工厂的自然选择

---

### 12. 代理模式（Proxy）

#### 概念定义

代理模式为其他对象提供一种代理以控制对这个对象的访问。代理对象在客户端和目标对象之间起到中介作用。

#### 意图与动机

- **访问控制**: 控制对对象的访问权限
- **延迟加载**: 延迟创建开销大的对象
- **远程代理**: 访问远程对象
- **缓存**: 缓存对象结果

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                       Client                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Subject                                │◄──── 主题接口
│  + Request()                                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│     RealSubject         │   │        Proxy            │
│  + Request()            │   │  - realSubject:Subject  │
│                         │   ├─────────────────────────┤
│                         │   │  + Request()            │
│                         │   │    checkAccess()        │
│                         │   │    realSubject.Request()│
│                         │   │    logAccess()          │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

```go
package proxy

// Subject 主题接口
type Subject interface {
    Request() string
}

// RealSubject 真实主题
type RealSubject struct{}

func (r *RealSubject) Request() string {
    return "RealSubject: Handling request."
}

// Proxy 代理
type Proxy struct {
    realSubject *RealSubject
}

func NewProxy() *Proxy {
    return &Proxy{}
}

func (p *Proxy) Request() string {
    if p.realSubject == nil {
        p.realSubject = &RealSubject{}
    }

    // 访问控制
    if !p.checkAccess() {
        return "Proxy: Access denied."
    }

    result := p.realSubject.Request()
    p.logAccess()

    return result
}

func (p *Proxy) checkAccess() bool {
    // 检查访问权限
    return true
}

func (p *Proxy) logAccess() {
    // 记录访问日志
}
```

#### 完整示例：图片加载代理

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Image 图片接口
type Image interface {
    Display() string
    GetFileName() string
    GetSize() int
}

// ============ 真实图片 ============

// RealImage 真实图片
type RealImage struct {
    filename string
    size     int
    data     []byte
}

func NewRealImage(filename string) *RealImage {
    fmt.Printf("[RealImage] 从磁盘加载: %s\n", filename)

    // 模拟加载大图片
    time.Sleep(1 * time.Second)

    // 模拟图片数据
    size := 1024 * 1024 * 5  // 5MB
    data := make([]byte, size)

    return &RealImage{
        filename: filename,
        size:     size,
        data:     data,
    }
}

func (r *RealImage) Display() string {
    return fmt.Sprintf("显示图片: %s (%d bytes)", r.filename, r.size)
}

func (r *RealImage) GetFileName() string {
    return r.filename
}

func (r *RealImage) GetSize() int {
    return r.size
}

// ============ 代理实现 ============

// ProxyImage 图片代理（虚拟代理）
type ProxyImage struct {
    filename  string
    realImage *RealImage
    mu        sync.Mutex
}

func NewProxyImage(filename string) *ProxyImage {
    return &ProxyImage{filename: filename}
}

func (p *ProxyImage) Display() string {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.realImage == nil {
        p.realImage = NewRealImage(p.filename)
    }
    return p.realImage.Display()
}

func (p *ProxyImage) GetFileName() string {
    return p.filename
}

func (p *ProxyImage) GetSize() int {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.realImage != nil {
        return p.realImage.GetSize()
    }
    return 0  // 尚未加载
}

// ============ 保护代理 ============

// ProtectedImage 保护代理
type ProtectedImage struct {
    image       Image
    userRole    string
    allowedRoles []string
}

func NewProtectedImage(image Image, userRole string, allowedRoles []string) *ProtectedImage {
    return &ProtectedImage{
        image:        image,
        userRole:     userRole,
        allowedRoles: allowedRoles,
    }
}

func (p *ProtectedImage) hasAccess() bool {
    for _, role := range p.allowedRoles {
        if role == p.userRole {
            return true
        }
    }
    return false
}

func (p *ProtectedImage) Display() string {
    if !p.hasAccess() {
        return fmt.Sprintf("[拒绝] 用户角色 '%s' 无权访问此图片", p.userRole)
    }
    return p.image.Display()
}

func (p *ProtectedImage) GetFileName() string {
    return p.image.GetFileName()
}

func (p *ProtectedImage) GetSize() int {
    if !p.hasAccess() {
        return 0
    }
    return p.image.GetSize()
}

// ============ 缓存代理 ============

// CachedImage 缓存代理
type CachedImage struct {
    image     Image
    cache     string
    cached    bool
    cacheTime time.Time
    ttl       time.Duration
}

func NewCachedImage(image Image, ttl time.Duration) *CachedImage {
    return &CachedImage{
        image: image,
        ttl:   ttl,
    }
}

func (c *CachedImage) Display() string {
    if c.cached && time.Since(c.cacheTime) < c.ttl {
        fmt.Println("[缓存] 返回缓存结果")
        return c.cache
    }

    result := c.image.Display()
    c.cache = result
    c.cached = true
    c.cacheTime = time.Now()
    return result
}

func (c *CachedImage) GetFileName() string {
    return c.image.GetFileName()
}

func (c *CachedImage) GetSize() int {
    return c.image.GetSize()
}

// ============ 日志代理 ============

// LoggingImage 日志代理
type LoggingImage struct {
    image Image
}

func NewLoggingImage(image Image) *LoggingImage {
    return &LoggingImage{image: image}
}

func (l *LoggingImage) Display() string {
    start := time.Now()
    fmt.Printf("[日志] 开始显示: %s\n", l.image.GetFileName())

    result := l.image.Display()

    duration := time.Since(start)
    fmt.Printf("[日志] 显示完成，耗时: %v\n", duration)
    return result
}

func (l *LoggingImage) GetFileName() string {
    return l.image.GetFileName()
}

func (l *LoggingImage) GetSize() int {
    return l.image.GetSize()
}

// ============ 图片管理器 ============

// ImageGallery 图片画廊
type ImageGallery struct {
    images []Image
}

func NewImageGallery() *ImageGallery {
    return &ImageGallery{
        images: make([]Image, 0),
    }
}

func (g *ImageGallery) AddImage(image Image) {
    g.images = append(g.images, image)
}

func (g *ImageGallery) DisplayAll() {
    for _, img := range g.images {
        fmt.Println(img.Display())
        fmt.Println()
    }
}

func (g *ImageGallery) GetTotalSize() int {
    total := 0
    for _, img := range g.images {
        total += img.GetSize()
    }
    return total
}

func main() {
    fmt.Println("========== 虚拟代理（延迟加载） ==========")

    // 创建代理图片列表（不立即加载）
    images := []Image{
        NewProxyImage("photo1.jpg"),
        NewProxyImage("photo2.jpg"),
        NewProxyImage("photo3.jpg"),
    }

    fmt.Println("图片列表创建完成（尚未加载）")
    fmt.Println()

    // 只显示第一张图片
    fmt.Println("显示第一张图片：")
    fmt.Println(images[0].Display())

    fmt.Println("\n再次显示第一张图片（已缓存）：")
    fmt.Println(images[0].Display())

    fmt.Println("\n========== 保护代理 ==========")

    adminImage := NewProtectedImage(
        NewProxyImage("secret.jpg"),
        "admin",
        []string{"admin", "manager"},
    )

    guestImage := NewProtectedImage(
        NewProxyImage("secret.jpg"),
        "guest",
        []string{"admin", "manager"},
    )

    fmt.Println("Admin访问：")
    fmt.Println(adminImage.Display())

    fmt.Println("\nGuest访问：")
    fmt.Println(guestImage.Display())

    fmt.Println("\n========== 缓存代理 ==========")

    cachedImage := NewCachedImage(
        NewProxyImage("cached.jpg"),
        5*time.Second,
    )

    fmt.Println("第一次调用：")
    fmt.Println(cachedImage.Display())

    fmt.Println("\n第二次调用（缓存）：")
    fmt.Println(cachedImage.Display())

    fmt.Println("\n========== 日志代理 ==========")

    loggedImage := NewLoggingImage(NewProxyImage("logged.jpg"))
    fmt.Println(loggedImage.Display())

    fmt.Println("\n========== 代理组合 ==========")

    // 组合多个代理
    complexImage := NewLoggingImage(
        NewCachedImage(
            NewProtectedImage(
                NewProxyImage("complex.jpg"),
                "admin",
                []string{"admin"},
            ),
            10*time.Second,
        ),
    )

    fmt.Println(complexImage.Display())
}
```

#### 反例说明

**错误1：代理暴露真实对象**

```go
// 错误：不应该暴露内部实现
type Proxy struct {
    RealSubject *RealSubject  // 公开字段，破坏了封装
}
```

**错误2：代理修改接口**

```go
// 错误：代理应该保持接口一致
type Proxy struct{}

func (p *Proxy) Request() string { return "" }
func (p *Proxy) ExtraMethod() {}  // 不要添加额外方法
```

**错误3：过度代理**

```go
// 错误：简单对象不需要代理
func NewProxyForSimpleObject() *Proxy {  // 过度设计
    // 简单对象直接使用即可
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 控制对象访问 | 增加响应时间 |
| 延迟加载 | 代码复杂度增加 |
| 透明性（客户端无感知） | 过多代理难以维护 |
| 多种代理可组合 | 调试困难 |

**适用场景**:

- 虚拟代理（延迟加载）
- 保护代理（访问控制）
- 缓存代理（结果缓存）
- 远程代理（RPC）
- 智能引用（引用计数）

**Go语言特化**:

- `net/http/httputil.ReverseProxy`标准实现
- 结合Context实现超时控制
- 函数闭包实现轻量级代理

---

## 行为型模式

行为型模式关注对象之间的通信和职责分配，定义对象如何交互和分配职责。

---

### 13. 责任链模式（Chain of Responsibility）

#### 概念定义

责任链模式使多个对象都有机会处理请求，从而避免请求的发送者和接收者之间的耦合关系。将这些对象连成一条链，并沿着这条链传递请求，直到有一个对象处理它为止。

#### 意图与动机

- **解耦**: 发送者和接收者解耦
- **动态组合**: 运行时改变处理链
- **单一职责**: 每个处理器只负责一种处理

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                       Client                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler                                │◄──── 处理器接口
│  - next: Handler                                            │
├─────────────────────────────────────────────────────────────┤
│  + SetNext(Handler)                                         │
│  + Handle(request)                                          │
│    if canHandle:                                            │
│      return handle(request)                                 │
│    else if next != nil:                                     │
│      return next.Handle(request)                            │
│    return nil                                               │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│  ConcreteHandler1       │   │  ConcreteHandler2       │
│  + canHandle()          │   │  + canHandle()          │
│  + handle()             │   │  + handle()             │
└─────────────────────────┘   └─────────────────────────┘
```

#### Go语言实现

```go
package chain

// Handler 处理器接口
type Handler interface {
    SetNext(handler Handler) Handler
    Handle(request interface{}) interface{}
}

// BaseHandler 基础处理器
type BaseHandler struct {
    next Handler
}

func (b *BaseHandler) SetNext(handler Handler) Handler {
    b.next = handler
    return handler
}

func (b *BaseHandler) HandleNext(request interface{}) interface{} {
    if b.next != nil {
        return b.next.Handle(request)
    }
    return nil
}
```

#### 完整示例：请求处理链

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

// Request 请求
type Request struct {
    ID       string
    Type     string
    Priority int
    Content  string
    User     string
}

// Response 响应
type Response struct {
    Success bool
    Message string
    Handler string
}

// Handler 处理器接口
type Handler interface {
    SetNext(next Handler) Handler
    Handle(request *Request) *Response
    GetName() string
}

// BaseHandler 基础处理器
type BaseHandler struct {
    next Handler
    name string
}

func (b *BaseHandler) SetNext(next Handler) Handler {
    b.next = next
    return next
}

func (b *BaseHandler) HandleNext(request *Request) *Response {
    if b.next != nil {
        return b.next.Handle(request)
    }
    return &Response{
        Success: false,
        Message: "没有处理器能处理此请求",
        Handler: b.name,
    }
}

func (b *BaseHandler) GetName() string {
    return b.name
}

// ============ 具体处理器 ============

// AuthHandler 认证处理器
type AuthHandler struct {
    BaseHandler
    validUsers []string
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        BaseHandler: BaseHandler{name: "AuthHandler"},
        validUsers:  []string{"admin", "user1", "user2"},
    }
}

func (a *AuthHandler) Handle(request *Request) *Response {
    fmt.Printf("[%s] 检查用户认证: %s\n", a.name, request.User)

    valid := false
    for _, user := range a.validUsers {
        if user == request.User {
            valid = true
            break
        }
    }

    if !valid {
        return &Response{
            Success: false,
            Message: "用户未认证",
            Handler: a.name,
        }
    }

    fmt.Printf("[%s] 认证通过\n", a.name)
    return a.HandleNext(request)
}

// RateLimitHandler 限流处理器
type RateLimitHandler struct {
    BaseHandler
    requests map[string][]time.Time
    limit    int
    window   time.Duration
}

func NewRateLimitHandler(limit int, window time.Duration) *RateLimitHandler {
    return &RateLimitHandler{
        BaseHandler: BaseHandler{name: "RateLimitHandler"},
        requests:    make(map[string][]time.Time),
        limit:       limit,
        window:      window,
    }
}

func (r *RateLimitHandler) Handle(request *Request) *Response {
    fmt.Printf("[%s] 检查限流: %s\n", r.name, request.User)

    now := time.Now()
    userRequests := r.requests[request.User]

    // 清理过期请求
    var validRequests []time.Time
    for _, t := range userRequests {
        if now.Sub(t) < r.window {
            validRequests = append(validRequests, t)
        }
    }

    if len(validRequests) >= r.limit {
        return &Response{
            Success: false,
            Message: "请求过于频繁，请稍后再试",
            Handler: r.name,
        }
    }

    validRequests = append(validRequests, now)
    r.requests[request.User] = validRequests

    fmt.Printf("[%s] 限流检查通过 (%d/%d)\n", r.name, len(validRequests), r.limit)
    return r.HandleNext(request)
}

// ValidationHandler 验证处理器
type ValidationHandler struct {
    BaseHandler
}

func NewValidationHandler() *ValidationHandler {
    return &ValidationHandler{
        BaseHandler: BaseHandler{name: "ValidationHandler"},
    }
}

func (v *ValidationHandler) Handle(request *Request) *Response {
    fmt.Printf("[%s] 验证请求: %s\n", v.name, request.ID)

    if request.Content == "" {
        return &Response{
            Success: false,
            Message: "请求内容不能为空",
            Handler: v.name,
        }
    }

    if len(request.Content) > 1000 {
        return &Response{
            Success: false,
            Message: "请求内容过长",
            Handler: v.name,
        }
    }

    fmt.Printf("[%s] 验证通过\n", v.name)
    return v.HandleNext(request)
}

// PermissionHandler 权限处理器
type PermissionHandler struct {
    BaseHandler
    permissions map[string][]string
}

func NewPermissionHandler() *PermissionHandler {
    return &PermissionHandler{
        BaseHandler: BaseHandler{name: "PermissionHandler"},
        permissions: map[string][]string{
            "admin": {"read", "write", "delete"},
            "user1": {"read", "write"},
            "user2": {"read"},
        },
    }
}

func (p *PermissionHandler) Handle(request *Request) *Response {
    fmt.Printf("[%s] 检查权限: %s -> %s\n", p.name, request.User, request.Type)

    userPerms, exists := p.permissions[request.User]
    if !exists {
        return &Response{
            Success: false,
            Message: "用户无权限",
            Handler: p.name,
        }
    }

    requiredPerm := ""
    switch request.Type {
    case "CREATE", "UPDATE":
        requiredPerm = "write"
    case "DELETE":
        requiredPerm = "delete"
    default:
        requiredPerm = "read"
    }

    hasPerm := false
    for _, perm := range userPerms {
        if perm == requiredPerm {
            hasPerm = true
            break
        }
    }

    if !hasPerm {
        return &Response{
            Success: false,
            Message: fmt.Sprintf("用户无%s权限", requiredPerm),
            Handler: p.name,
        }
    }

    fmt.Printf("[%s] 权限检查通过\n", p.name)
    return p.HandleNext(request)
}

// BusinessHandler 业务处理器
type BusinessHandler struct {
    BaseHandler
}

func NewBusinessHandler() *BusinessHandler {
    return &BusinessHandler{
        BaseHandler: BaseHandler{name: "BusinessHandler"},
    }
}

func (b *BusinessHandler) Handle(request *Request) *Response {
    fmt.Printf("[%s] 处理业务逻辑: %s\n", b.name, request.Type)

    // 模拟业务处理
    time.Sleep(100 * time.Millisecond)

    return &Response{
        Success: true,
        Message: fmt.Sprintf("请求 %s 处理成功", request.ID),
        Handler: b.name,
    }
}

// ============ 链式构建器 ============

// ChainBuilder 链式构建器
type ChainBuilder struct {
    first Handler
    last  Handler
}

func NewChainBuilder() *ChainBuilder {
    return &ChainBuilder{}
}

func (b *ChainBuilder) Add(handler Handler) *ChainBuilder {
    if b.first == nil {
        b.first = handler
        b.last = handler
    } else {
        b.last.SetNext(handler)
        b.last = handler
    }
    return b
}

func (b *ChainBuilder) Build() Handler {
    return b.first
}

// ============ 处理链管理器 ============

type RequestProcessor struct {
    chain Handler
}

func NewRequestProcessor() *RequestProcessor {
    // 构建处理链
    chain := NewChainBuilder().
        Add(NewAuthHandler()).
        Add(NewRateLimitHandler(5, time.Minute)).
        Add(NewValidationHandler()).
        Add(NewPermissionHandler()).
        Add(NewBusinessHandler()).
        Build()

    return &RequestProcessor{chain: chain}
}

func (p *RequestProcessor) Process(request *Request) *Response {
    fmt.Printf("\n处理请求: %s (Type: %s, User: %s)\n",
        request.ID, request.Type, request.User)
    fmt.Println(strings.Repeat("-", 50))

    response := p.chain.Handle(request)

    fmt.Println(strings.Repeat("-", 50))
    fmt.Printf("结果: Success=%v, Message=%s, Handler=%s\n",
        response.Success, response.Message, response.Handler)

    return response
}

func main() {
    processor := NewRequestProcessor()

    // 测试用例1：正常请求
    processor.Process(&Request{
        ID:      "REQ001",
        Type:    "CREATE",
        Content: "创建新订单",
        User:    "admin",
    })

    // 测试用例2：未认证用户
    processor.Process(&Request{
        ID:      "REQ002",
        Type:    "READ",
        Content: "查询数据",
        User:    "unknown",
    })

    // 测试用例3：权限不足
    processor.Process(&Request{
        ID:      "REQ003",
        Type:    "DELETE",
        Content: "删除数据",
        User:    "user2",
    })

    // 测试用例4：验证失败
    processor.Process(&Request{
        ID:      "REQ004",
        Type:    "CREATE",
        Content: "",  // 空内容
        User:    "admin",
    })

    // 测试用例5：触发限流
    fmt.Println("\n========== 限流测试 ==========")
    for i := 0; i < 7; i++ {
        processor.Process(&Request{
            ID:      fmt.Sprintf("REQ100%d", i),
            Type:    "READ",
            Content: "查询",
            User:    "user1",
        })
    }

    // 动态构建不同处理链
    fmt.Println("\n========== 自定义处理链 ==========")
    customChain := NewChainBuilder().
        Add(NewValidationHandler()).
        Add(NewBusinessHandler()).
        Build()

    response := customChain.Handle(&Request{
        ID:      "CUSTOM001",
        Type:    "TEST",
        Content: "测试内容",
        User:    "anyone",
    })

    fmt.Printf("自定义链结果: %v - %s\n", response.Success, response.Message)
}
```

#### 反例说明

**错误1：处理器不调用下一个**

```go
// 错误：处理完成后应该继续传递
func (h *Handler) Handle(request *Request) *Response {
    if h.canHandle(request) {
        return h.doHandle(request)  // 忘记调用下一个
    }
    return nil
}
```

**错误2：循环链**

```go
// 错误：循环引用会导致无限循环
handler1.SetNext(handler2)
handler2.SetNext(handler1)  // 循环！
```

**错误3：处理器包含过多逻辑**

```go
// 错误：每个处理器应该只负责单一职责
func (h *Handler) Handle(request *Request) *Response {
    // 认证
    // 限流
    // 验证
    // 权限检查
    // 业务处理
    // 太多职责！
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 解耦发送者和接收者 | 请求可能未被处理 |
| 动态组合处理器 | 调试困难 |
| 单一职责原则 | 性能开销（遍历链） |
| 符合开闭原则 | 链配置复杂 |

**适用场景**:

- 请求过滤器链
- 日志处理器链
- 中间件链
- 审批流程
- 异常处理

**Go语言特化**:

- `http.Handler`中间件链
- 函数闭包实现轻量级处理器
- 结合Context传递状态

---

### 14. 命令模式（Command）

#### 概念定义

命令模式将请求封装为对象，从而可以用不同的请求、队列或日志来参数化客户端。命令模式也支持可撤销的操作。

#### 意图与动机

- **解耦**: 解耦调用者和接收者
- **可撤销**: 支持操作撤销
- **队列**: 支持请求排队
- **日志**: 支持操作日志

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                       Client                                │
└──────────────────────────┬──────────────────────────────────┘
                           │ creates
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Command                                │◄──── 命令接口
│  + Execute()                                                │
│  + Undo()                                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│   ConcreteCommand       │   │   ConcreteCommandB      │
│  - receiver: Receiver   │   │  - receiver: Receiver   │
├─────────────────────────┤   ├─────────────────────────┤
│  + Execute()            │   │  + Execute()            │
│    receiver.Action()    │   │    receiver.Action()    │
│  + Undo()               │   │  + Undo()               │
│    receiver.UndoAction()│   │    receiver.UndoAction()│
└─────────────────────────┘   └─────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Invoker                                │
│  - commands: []Command                                      │
│  - history: []Command                                       │
├─────────────────────────────────────────────────────────────┤
│  + ExecuteCommand(Command)                                  │
│  + Undo()                                                   │
│  + Redo()                                                   │
└─────────────────────────────────────────────────────────────┘
```

#### Go语言实现

```go
package command

// Command 命令接口
type Command interface {
    Execute() error
    Undo() error
    GetName() string
}

// Receiver 接收者
type Receiver struct {
    state string
}

func (r *Receiver) Action() {
    r.state = "modified"
}

func (r *Receiver) UndoAction() {
    r.state = "original"
}

// ConcreteCommand 具体命令
type ConcreteCommand struct {
    receiver *Receiver
    name     string
}

func NewConcreteCommand(receiver *Receiver, name string) *ConcreteCommand {
    return &ConcreteCommand{
        receiver: receiver,
        name:     name,
    }
}

func (c *ConcreteCommand) Execute() error {
    c.receiver.Action()
    return nil
}

func (c *ConcreteCommand) Undo() error {
    c.receiver.UndoAction()
    return nil
}

func (c *ConcreteCommand) GetName() string {
    return c.name
}
```

#### 完整示例：文本编辑器

```go
package main

import (
    "fmt"
    "strings"
)

// ============ 接收者：编辑器 ============

// TextEditor 文本编辑器
type TextEditor struct {
    content   string
    clipboard string
    selection struct {
        start int
        end   int
    }
}

func NewTextEditor() *TextEditor {
    return &TextEditor{
        content: "",
    }
}

func (e *TextEditor) Insert(text string, position int) {
    if position < 0 || position > len(e.content) {
        position = len(e.content)
    }
    e.content = e.content[:position] + text + e.content[position:]
}

func (e *TextEditor) Delete(start, end int) string {
    if start < 0 {
        start = 0
    }
    if end > len(e.content) {
        end = len(e.content)
    }
    deleted := e.content[start:end]
    e.content = e.content[:start] + e.content[end:]
    return deleted
}

func (e *TextEditor) Copy(start, end int) {
    if start < 0 {
        start = 0
    }
    if end > len(e.content) {
        end = len(e.content)
    }
    e.clipboard = e.content[start:end]
}

func (e *TextEditor) Paste(position int) {
    e.Insert(e.clipboard, position)
}

func (e *TextEditor) Select(start, end int) {
    e.selection.start = start
    e.selection.end = end
}

func (e *TextEditor) GetSelection() (int, int) {
    return e.selection.start, e.selection.end
}

func (e *TextEditor) GetContent() string {
    return e.content
}

func (e *TextEditor) SetContent(content string) {
    e.content = content
}

func (e *TextEditor) GetClipboard() string {
    return e.clipboard
}

// ============ 命令接口和具体命令 ============

// Command 命令接口
type Command interface {
    Execute() error
    Undo() error
    GetName() string
}

// InsertCommand 插入命令
type InsertCommand struct {
    editor   *TextEditor
    text     string
    position int
    prevContent string
    name     string
}

func NewInsertCommand(editor *TextEditor, text string, position int) *InsertCommand {
    return &InsertCommand{
        editor:   editor,
        text:     text,
        position: position,
        name:     fmt.Sprintf("Insert '%s' at %d", text, position),
    }
}

func (c *InsertCommand) Execute() error {
    c.prevContent = c.editor.GetContent()
    c.editor.Insert(c.text, c.position)
    return nil
}

func (c *InsertCommand) Undo() error {
    c.editor.SetContent(c.prevContent)
    return nil
}

func (c *InsertCommand) GetName() string {
    return c.name
}

// DeleteCommand 删除命令
type DeleteCommand struct {
    editor      *TextEditor
    start       int
    end         int
    deletedText string
    prevContent string
    name        string
}

func NewDeleteCommand(editor *TextEditor, start, end int) *DeleteCommand {
    return &DeleteCommand{
        editor: editor,
        start:  start,
        end:    end,
        name:   fmt.Sprintf("Delete from %d to %d", start, end),
    }
}

func (c *DeleteCommand) Execute() error {
    c.prevContent = c.editor.GetContent()
    c.deletedText = c.editor.Delete(c.start, c.end)
    return nil
}

func (c *DeleteCommand) Undo() error {
    c.editor.SetContent(c.prevContent)
    return nil
}

func (c *DeleteCommand) GetName() string {
    return c.name
}

// CopyCommand 复制命令
type CopyCommand struct {
    editor *TextEditor
    start  int
    end    int
    name   string
}

func NewCopyCommand(editor *TextEditor, start, end int) *CopyCommand {
    return &CopyCommand{
        editor: editor,
        start:  start,
        end:    end,
        name:   fmt.Sprintf("Copy from %d to %d", start, end),
    }
}

func (c *CopyCommand) Execute() error {
    c.editor.Copy(c.start, c.end)
    return nil
}

func (c *CopyCommand) Undo() error {
    // 复制操作不需要撤销
    return nil
}

func (c *CopyCommand) GetName() string {
    return c.name
}

// PasteCommand 粘贴命令
type PasteCommand struct {
    editor      *TextEditor
    position    int
    prevContent string
    name        string
}

func NewPasteCommand(editor *TextEditor, position int) *PasteCommand {
    return &PasteCommand{
        editor:   editor,
        position: position,
        name:     fmt.Sprintf("Paste at %d", position),
    }
}

func (c *PasteCommand) Execute() error {
    c.prevContent = c.editor.GetContent()
    c.editor.Paste(c.position)
    return nil
}

func (c *PasteCommand) Undo() error {
    c.editor.SetContent(c.prevContent)
    return nil
}

func (c *PasteCommand) GetName() string {
    return c.name
}

// MacroCommand 宏命令（组合命令）
type MacroCommand struct {
    commands []Command
    name     string
}

func NewMacroCommand(name string) *MacroCommand {
    return &MacroCommand{
        commands: make([]Command, 0),
        name:     name,
    }
}

func (m *MacroCommand) Add(command Command) {
    m.commands = append(m.commands, command)
}

func (m *MacroCommand) Execute() error {
    for _, cmd := range m.commands {
        if err := cmd.Execute(); err != nil {
            return err
        }
    }
    return nil
}

func (m *MacroCommand) Undo() error {
    // 逆序撤销
    for i := len(m.commands) - 1; i >= 0; i-- {
        if err := m.commands[i].Undo(); err != nil {
            return err
        }
    }
    return nil
}

func (m *MacroCommand) GetName() string {
    return m.name
}

// ============ 调用者：命令管理器 ============

// CommandManager 命令管理器
type CommandManager struct {
    history     []Command
    currentIdx  int
    editor      *TextEditor
}

func NewCommandManager(editor *TextEditor) *CommandManager {
    return &CommandManager{
        history:    make([]Command, 0),
        currentIdx: -1,
        editor:     editor,
    }
}

func (m *CommandManager) Execute(command Command) error {
    // 执行命令
    if err := command.Execute(); err != nil {
        return err
    }

    // 清除撤销后的历史
    if m.currentIdx < len(m.history)-1 {
        m.history = m.history[:m.currentIdx+1]
    }

    // 添加到历史
    m.history = append(m.history, command)
    m.currentIdx++

    fmt.Printf("执行: %s\n", command.GetName())
    return nil
}

func (m *CommandManager) Undo() error {
    if m.currentIdx < 0 {
        return fmt.Errorf("没有可撤销的操作")
    }

    command := m.history[m.currentIdx]
    if err := command.Undo(); err != nil {
        return err
    }

    fmt.Printf("撤销: %s\n", command.GetName())
    m.currentIdx--
    return nil
}

func (m *CommandManager) Redo() error {
    if m.currentIdx >= len(m.history)-1 {
        return fmt.Errorf("没有可重做的操作")
    }

    m.currentIdx++
    command := m.history[m.currentIdx]
    if err := command.Execute(); err != nil {
        return err
    }

    fmt.Printf("重做: %s\n", command.GetName())
    return nil
}

func (m *CommandManager) GetHistory() []string {
    var names []string
    for i, cmd := range m.history {
        prefix := "  "
        if i == m.currentIdx {
            prefix = "► "
        }
        names = append(names, prefix+cmd.GetName())
    }
    return names
}

// ============ UI ============

func printState(editor *TextEditor, manager *CommandManager) {
    fmt.Println()
    fmt.Println("=" + strings.Repeat("=", 50))
    fmt.Printf("内容: \"%s\"\n", editor.GetContent())
    fmt.Printf("剪贴板: \"%s\"\n", editor.GetClipboard())
    fmt.Println("-" + strings.Repeat("-", 50))
    fmt.Println("历史:")
    for _, h := range manager.GetHistory() {
        fmt.Println(h)
    }
    fmt.Println("=" + strings.Repeat("=", 50))
}

func main() {
    editor := NewTextEditor()
    manager := NewCommandManager(editor)

    fmt.Println("========== 文本编辑器演示 ==========")

    // 插入文本
    manager.Execute(NewInsertCommand(editor, "Hello, ", 0))
    printState(editor, manager)

    manager.Execute(NewInsertCommand(editor, "World!", 7))
    printState(editor, manager)

    // 选择并复制
    manager.Execute(NewCopyCommand(editor, 0, 5))
    printState(editor, manager)

    // 粘贴
    manager.Execute(NewPasteCommand(editor, 12))
    printState(editor, manager)

    // 删除
    manager.Execute(NewDeleteCommand(editor, 5, 7))
    printState(editor, manager)

    // 撤销操作
    fmt.Println("\n========== 撤销操作 ==========")
    manager.Undo()
    printState(editor, manager)

    manager.Undo()
    printState(editor, manager)

    // 重做
    fmt.Println("\n========== 重做操作 ==========")
    manager.Redo()
    printState(editor, manager)

    // 宏命令
    fmt.Println("\n========== 宏命令 ==========")
    macro := NewMacroCommand("格式化文本")
    macro.Add(NewInsertCommand(editor, "[", 0))
    macro.Add(NewInsertCommand(editor, "]", len(editor.GetContent())+1))

    manager.Execute(macro)
    printState(editor, manager)

    // 撤销宏
    manager.Undo()
    printState(editor, manager)
}
```

#### 反例说明

**错误1：命令直接执行而非封装**

```go
// 错误：没有封装请求
func main() {
    editor := &TextEditor{}
    editor.Insert("text", 0)  // 直接调用，无法撤销
}
```

**错误2：命令不保存状态**

```go
// 错误：没有保存之前的状态
type DeleteCommand struct {
    editor *TextEditor
}

func (c *DeleteCommand) Execute() {
    c.editor.Delete(0, 5)
    // 没有保存删除的内容，无法撤销
}
```

**错误3：命令包含业务逻辑**

```go
// 错误：命令应该只负责调用接收者
type Command struct {
    // ...
}

func (c *Command) Execute() {
    // 大量业务逻辑不应该在这里
    validate()
    process()
    save()
    notify()
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 解耦调用者和接收者 | 类数量增加 |
| 支持撤销/重做 | 代码复杂度增加 |
| 支持命令队列 | 可能产生大量命令对象 |
| 支持宏命令 | 维护历史记录需要内存 |

**适用场景**:

- 文本编辑器（撤销/重做）
- 事务系统
- 任务队列
- 宏录制
- 远程调用

**Go语言特化**:

- 函数类型作为轻量级命令
- 结合Channel实现命令队列
- Context传递取消信号

---

### 15. 解释器模式（Interpreter）

#### 概念定义

解释器模式给定一个语言，定义它的文法的一种表示，并定义一个解释器，这个解释器使用该表示来解释语言中的句子。

#### 意图与动机

- **领域特定语言**: 实现简单的DSL
- **语法解析**: 解析和执行特定语法
- **规则引擎**: 实现业务规则

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                    Context                                  │
│  - variables: map[string]int                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Expression                                │◄──── 表达式接口
│  + Interpret(Context) int                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ Terminal      │  │ NonTerminal   │  │ NonTerminal   │
│ Expression    │  │  (Add)        │  │  (Subtract)   │
│ (Number)      │  │ - left: Exp   │  │ - left: Exp   │
│ - value: int  │  │ - right: Exp  │  │ - right: Exp  │
├───────────────┤  ├───────────────┤  ├───────────────┤
│ + Interpret() │  │ + Interpret() │  │ + Interpret() │
│   return val  │  │   return l+r  │  │   return l-r  │
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

```go
package interpreter

// Expression 表达式接口
type Expression interface {
    Interpret() int
}

// NumberExpression 数字表达式（终结符）
type NumberExpression struct {
    value int
}

func NewNumberExpression(value int) *NumberExpression {
    return &NumberExpression{value: value}
}

func (n *NumberExpression) Interpret() int {
    return n.value
}

// AddExpression 加法表达式（非终结符）
type AddExpression struct {
    left  Expression
    right Expression
}

func NewAddExpression(left, right Expression) *AddExpression {
    return &AddExpression{left: left, right: right}
}

func (a *AddExpression) Interpret() int {
    return a.left.Interpret() + a.right.Interpret()
}

// SubtractExpression 减法表达式
type SubtractExpression struct {
    left  Expression
    right Expression
}

func NewSubtractExpression(left, right Expression) *SubtractExpression {
    return &SubtractExpression{left: left, right: right}
}

func (s *SubtractExpression) Interpret() int {
    return s.left.Interpret() - s.right.Interpret()
}
```

#### 完整示例：规则引擎

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
)

// ============ 表达式接口 ============

type Expression interface {
    Evaluate(context map[string]interface{}) (interface{}, error)
    String() string
}

// ============ 基本表达式 ============

// NumberExpression 数字表达式
type NumberExpression struct {
    value float64
}

func NewNumberExpression(value float64) *NumberExpression {
    return &NumberExpression{value: value}
}

func (n *NumberExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    return n.value, nil
}

func (n *NumberExpression) String() string {
    return fmt.Sprintf("%g", n.value)
}

// StringExpression 字符串表达式
type StringExpression struct {
    value string
}

func NewStringExpression(value string) *StringExpression {
    return &StringExpression{value: value}
}

func (s *StringExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    return s.value, nil
}

func (s *StringExpression) String() string {
    return fmt.Sprintf("\"%s\"", s.value)
}

// VariableExpression 变量表达式
type VariableExpression struct {
    name string
}

func NewVariableExpression(name string) *VariableExpression {
    return &VariableExpression{name: name}
}

func (v *VariableExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    if val, ok := context[v.name]; ok {
        return val, nil
    }
    return nil, fmt.Errorf("变量 %s 未定义", v.name)
}

func (v *VariableExpression) String() string {
    return v.name
}

// ============ 比较表达式 ============

// EqualsExpression 等于表达式
type EqualsExpression struct {
    left  Expression
    right Expression
}

func NewEqualsExpression(left, right Expression) *EqualsExpression {
    return &EqualsExpression{left: left, right: right}
}

func (e *EqualsExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    leftVal, err := e.left.Evaluate(context)
    if err != nil {
        return nil, err
    }
    rightVal, err := e.right.Evaluate(context)
    if err != nil {
        return nil, err
    }
    return leftVal == rightVal, nil
}

func (e *EqualsExpression) String() string {
    return fmt.Sprintf("(%s == %s)", e.left.String(), e.right.String())
}

// GreaterThanExpression 大于表达式
type GreaterThanExpression struct {
    left  Expression
    right Expression
}

func NewGreaterThanExpression(left, right Expression) *GreaterThanExpression {
    return &GreaterThanExpression{left: left, right: right}
}

func (g *GreaterThanExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    leftVal, err := g.left.Evaluate(context)
    if err != nil {
        return nil, err
    }
    rightVal, err := g.right.Evaluate(context)
    if err != nil {
        return nil, err
    }

    leftNum, ok1 := toFloat64(leftVal)
    rightNum, ok2 := toFloat64(rightVal)
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("比较操作需要数字")
    }

    return leftNum > rightNum, nil
}

func (g *GreaterThanExpression) String() string {
    return fmt.Sprintf("(%s > %s)", g.left.String(), g.right.String())
}

// ============ 逻辑表达式 ============

// AndExpression 与表达式
type AndExpression struct {
    left  Expression
    right Expression
}

func NewAndExpression(left, right Expression) *AndExpression {
    return &AndExpression{left: left, right: right}
}

func (a *AndExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    leftVal, err := a.left.Evaluate(context)
    if err != nil {
        return nil, err
    }
    if !toBool(leftVal) {
        return false, nil
    }

    rightVal, err := a.right.Evaluate(context)
    if err != nil {
        return nil, err
    }
    return toBool(rightVal), nil
}

func (a *AndExpression) String() string {
    return fmt.Sprintf("(%s AND %s)", a.left.String(), a.right.String())
}

// OrExpression 或表达式
type OrExpression struct {
    left  Expression
    right Expression
}

func NewOrExpression(left, right Expression) *OrExpression {
    return &OrExpression{left: left, right: right}
}

func (o *OrExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    leftVal, err := o.left.Evaluate(context)
    if err != nil {
        return nil, err
    }
    if toBool(leftVal) {
        return true, nil
    }

    rightVal, err := o.right.Evaluate(context)
    if err != nil {
        return nil, err
    }
    return toBool(rightVal), nil
}

func (o *OrExpression) String() string {
    return fmt.Sprintf("(%s OR %s)", o.left.String(), o.right.String())
}

// NotExpression 非表达式
type NotExpression struct {
    operand Expression
}

func NewNotExpression(operand Expression) *NotExpression {
    return &NotExpression{operand: operand}
}

func (n *NotExpression) Evaluate(context map[string]interface{}) (interface{}, error) {
    val, err := n.operand.Evaluate(context)
    if err != nil {
        return nil, err
    }
    return !toBool(val), nil
}

func (n *NotExpression) String() string {
    return fmt.Sprintf("(NOT %s)", n.operand.String())
}

// ============ 辅助函数 ============

func toFloat64(v interface{}) (float64, bool) {
    switch val := v.(type) {
    case float64:
        return val, true
    case int:
        return float64(val), true
    case int64:
        return float64(val), true
    default:
        return 0, false
    }
}

func toBool(v interface{}) bool {
    switch val := v.(type) {
    case bool:
        return val
    case float64:
        return val != 0
    case int:
        return val != 0
    case string:
        return val != ""
    default:
        return v != nil
    }
}

// ============ 解析器 ============

type Parser struct {
    tokens []string
    pos    int
}

func NewParser(input string) *Parser {
    // 简单分词
    input = strings.ReplaceAll(input, "(", " ( ")
    input = strings.ReplaceAll(input, ")", " ) ")
    tokens := strings.Fields(input)
    return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (Expression, error) {
    return p.parseExpression()
}

func (p *Parser) parseExpression() (Expression, error) {
    return p.parseOr()
}

func (p *Parser) parseOr() (Expression, error) {
    left, err := p.parseAnd()
    if err != nil {
        return nil, err
    }

    for p.match("OR") {
        right, err := p.parseAnd()
        if err != nil {
            return nil, err
        }
        left = NewOrExpression(left, right)
    }

    return left, nil
}

func (p *Parser) parseAnd() (Expression, error) {
    left, err := p.parseEquality()
    if err != nil {
        return nil, err
    }

    for p.match("AND") {
        right, err := p.parseEquality()
        if err != nil {
            return nil, err
        }
        left = NewAndExpression(left, right)
    }

    return left, nil
}

func (p *Parser) parseEquality() (Expression, error) {
    left, err := p.parseComparison()
    if err != nil {
        return nil, err
    }

    if p.match("==") {
        right, err := p.parseComparison()
        if err != nil {
            return nil, err
        }
        left = NewEqualsExpression(left, right)
    }

    return left, nil
}

func (p *Parser) parseComparison() (Expression, error) {
    left, err := p.parseUnary()
    if err != nil {
        return nil, err
    }

    if p.match(">") {
        right, err := p.parseUnary()
        if err != nil {
            return nil, err
        }
        left = NewGreaterThanExpression(left, right)
    }

    return left, nil
}

func (p *Parser) parseUnary() (Expression, error) {
    if p.match("NOT") {
        operand, err := p.parseUnary()
        if err != nil {
            return nil, err
        }
        return NewNotExpression(operand), nil
    }

    return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
    if p.match("(") {
        expr, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        if !p.match(")") {
            return nil, fmt.Errorf("期望 )")
        }
        return expr, nil
    }

    token := p.consume()

    // 尝试解析为数字
    if num, err := strconv.ParseFloat(token, 64); err == nil {
        return NewNumberExpression(num), nil
    }

    // 尝试解析为字符串
    if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
        return NewStringExpression(token[1 : len(token)-1]), nil
    }

    // 否则作为变量
    return NewVariableExpression(token), nil
}

func (p *Parser) match(expected string) bool {
    if p.pos >= len(p.tokens) {
        return false
    }
    if p.tokens[p.pos] == expected {
        p.pos++
        return true
    }
    return false
}

func (p *Parser) consume() string {
    if p.pos >= len(p.tokens) {
        return ""
    }
    token := p.tokens[p.pos]
    p.pos++
    return token
}

// ============ 规则引擎 ============

type RuleEngine struct {
    rules map[string]Expression
}

func NewRuleEngine() *RuleEngine {
    return &RuleEngine{
        rules: make(map[string]Expression),
    }
}

func (e *RuleEngine) AddRule(name, rule string) error {
    parser := NewParser(rule)
    expr, err := parser.Parse()
    if err != nil {
        return err
    }
    e.rules[name] = expr
    return nil
}

func (e *RuleEngine) Evaluate(ruleName string, context map[string]interface{}) (bool, error) {
    expr, ok := e.rules[ruleName]
    if !ok {
        return false, fmt.Errorf("规则 %s 不存在", ruleName)
    }

    result, err := expr.Evaluate(context)
    if err != nil {
        return false, err
    }

    return toBool(result), nil
}

func (e *RuleEngine) GetRuleExpression(name string) string {
    if expr, ok := e.rules[name]; ok {
        return expr.String()
    }
    return ""
}

func main() {
    engine := NewRuleEngine()

    fmt.Println("========== 规则引擎演示 ==========")

    // 添加规则
    rules := map[string]string{
        "adult_check":   "age >= 18",
        "vip_discount":  "is_vip == true AND total > 100",
        "free_shipping": "total > 200 OR is_vip == true",
        "complex_rule":  "(age >= 18 AND country == \"US\") OR (age >= 21 AND country == \"EU\")",
    }

    for name, rule := range rules {
        if err := engine.AddRule(name, rule); err != nil {
            fmt.Printf("添加规则 %s 失败: %v\n", name, err)
        } else {
            fmt.Printf("规则 '%s': %s\n", name, engine.GetRuleExpression(name))
        }
    }

    // 测试规则
    fmt.Println("\n========== 规则测试 ==========")

    testCases := []struct {
        ruleName string
        context  map[string]interface{}
    }{
        {
            "adult_check",
            map[string]interface{}{"age": 25},
        },
        {
            "adult_check",
            map[string]interface{}{"age": 16},
        },
        {
            "vip_discount",
            map[string]interface{}{"is_vip": true, "total": 150},
        },
        {
            "vip_discount",
            map[string]interface{}{"is_vip": false, "total": 150},
        },
        {
            "free_shipping",
            map[string]interface{}{"total": 250, "is_vip": false},
        },
        {
            "complex_rule",
            map[string]interface{}{"age": 20, "country": "US"},
        },
        {
            "complex_rule",
            map[string]interface{}{"age": 20, "country": "EU"},
        },
    }

    for _, tc := range testCases {
        result, err := engine.Evaluate(tc.ruleName, tc.context)
        if err != nil {
            fmt.Printf("规则 '%s' 评估失败: %v\n", tc.ruleName, err)
        } else {
            fmt.Printf("规则 '%s' with %v => %v\n", tc.ruleName, tc.context, result)
        }
    }
}
```

#### 反例说明

**错误1：过于复杂的文法**

```go
// 错误：解释器不适合复杂文法
// 应该使用成熟的解析器生成器（如yacc、ANTLR）
```

**错误2：没有错误处理**

```go
// 错误：应该返回错误而非panic
func (e *Expression) Interpret() int {
    return e.left.Interpret() / e.right.Interpret()  // 可能除零
}
```

**错误3：混合解析和执行**

```go
// 错误：解析和执行应该分离
func Evaluate(input string) int {
    // 解析和执行混在一起
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 易于改变和扩展文法 | 复杂文法难以维护 |
| 实现简单文法容易 | 性能问题（递归解释） |
| 增加新的解释表达式容易 | 不适合复杂文法 |

**适用场景**:

- 简单DSL
- 配置文件解析
- 规则引擎
- 查询语言

**Go语言特化**:

- 使用`text/scanner`包辅助解析
- 接口实现表达式多态
- 递归下降解析器

---

### 16. 迭代器模式（Iterator）

#### 概念定义

迭代器模式提供一种方法顺序访问一个聚合对象中各个元素，而又无需暴露该对象的内部表示。

#### 意图与动机

- **统一接口**: 统一遍历不同集合的方式
- **隐藏实现**: 不暴露集合内部结构
- **多遍历**: 支持同时多个遍历

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Iterator                               │◄──── 迭代器接口
│  + HasNext() bool                                           │
│  + Next() interface{}                                       │
│  + Reset()                                                  │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│  ConcreteIterator       │   │  ConcreteIteratorB      │
│  - collection: []Item   │   │  - collection: Map      │
│  - index: int           │   │  - keys: []Key          │
├─────────────────────────┤   │  - index: int           │
│  + HasNext()            │   ├─────────────────────────┤
│  + Next()               │   │  + HasNext()            │
│  + Reset()              │   │  + Next()               │
└─────────────────────────┘   │  + Reset()              │
                              └─────────────────────────┘
                                         ▲
                                         │ creates
┌─────────────────────────────────────────────────────────────┐
│                     Aggregate                               │◄──── 聚合接口
│  + CreateIterator() Iterator                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  ConcreteAggregate                          │
│  - items: []Item                                            │
├─────────────────────────────────────────────────────────────┤
│  + CreateIterator() Iterator                                │
└─────────────────────────────────────────────────────────────┘
```

#### Go语言实现

Go语言内置了迭代器支持：`range`关键字和`container/iterator`包。

```go
package iterator

// Iterator 迭代器接口
type Iterator interface {
    HasNext() bool
    Next() interface{}
    Reset()
}

// Aggregate 聚合接口
type Aggregate interface {
    CreateIterator() Iterator
}

// ConcreteIterator 具体迭代器
type ConcreteIterator struct {
    collection []interface{}
    index      int
}

func NewConcreteIterator(collection []interface{}) *ConcreteIterator {
    return &ConcreteIterator{
        collection: collection,
        index:      0,
    }
}

func (c *ConcreteIterator) HasNext() bool {
    return c.index < len(c.collection)
}

func (c *ConcreteIterator) Next() interface{} {
    if c.HasNext() {
        item := c.collection[c.index]
        c.index++
        return item
    }
    return nil
}

func (c *ConcreteIterator) Reset() {
    c.index = 0
}

// ConcreteAggregate 具体聚合
type ConcreteAggregate struct {
    items []interface{}
}

func (c *ConcreteAggregate) Add(item interface{}) {
    c.items = append(c.items, item)
}

func (c *ConcreteAggregate) CreateIterator() Iterator {
    return NewConcreteIterator(c.items)
}
```

#### 完整示例：自定义集合迭代器

```go
package main

import (
    "fmt"
    "strings"
)

// ============ 迭代器接口 ============

type Iterator interface {
    HasNext() bool
    Next() interface{}
    Reset()
}

// ============ 图书和书架 ============

// Book 图书
type Book struct {
    Title  string
    Author string
    ISBN   string
    Price  float64
}

func (b *Book) String() string {
    return fmt.Sprintf("《%s》 by %s ($%.2f)", b.Title, b.Author, b.Price)
}

// BookShelf 书架
type BookShelf struct {
    name  string
    books []*Book
}

func NewBookShelf(name string) *BookShelf {
    return &BookShelf{
        name:  name,
        books: make([]*Book, 0),
    }
}

func (s *BookShelf) AddBook(book *Book) {
    s.books = append(s.books, book)
}

func (s *BookShelf) GetName() string {
    return s.name
}

// BookIterator 图书迭代器
type BookIterator struct {
    shelf *BookShelf
    index int
}

func NewBookIterator(shelf *BookShelf) *BookIterator {
    return &BookIterator{
        shelf: shelf,
        index: 0,
    }
}

func (i *BookIterator) HasNext() bool {
    return i.index < len(i.shelf.books)
}

func (i *BookIterator) Next() interface{} {
    if i.HasNext() {
        book := i.shelf.books[i.index]
        i.index++
        return book
    }
    return nil
}

func (i *BookIterator) Reset() {
    i.index = 0
}

// ReverseBookIterator 反向迭代器
type ReverseBookIterator struct {
    shelf *BookShelf
    index int
}

func NewReverseBookIterator(shelf *BookShelf) *ReverseBookIterator {
    return &ReverseBookIterator{
        shelf: shelf,
        index: len(shelf.books) - 1,
    }
}

func (i *ReverseBookIterator) HasNext() bool {
    return i.index >= 0
}

func (i *ReverseBookIterator) Next() interface{} {
    if i.HasNext() {
        book := i.shelf.books[i.index]
        i.index--
        return book
    }
    return nil
}

func (i *ReverseBookIterator) Reset() {
    i.index = len(i.shelf.books) - 1
}

// AuthorFilterIterator 作者过滤迭代器
type AuthorFilterIterator struct {
    shelf      *BookShelf
    author     string
    index      int
}

func NewAuthorFilterIterator(shelf *BookShelf, author string) *AuthorFilterIterator {
    return &AuthorFilterIterator{
        shelf:  shelf,
        author: author,
        index:  0,
    }
}

func (i *AuthorFilterIterator) HasNext() bool {
    for i.index < len(i.shelf.books) {
        if i.shelf.books[i.index].Author == i.author {
            return true
        }
        i.index++
    }
    return false
}

func (i *AuthorFilterIterator) Next() interface{} {
    for i.index < len(i.shelf.books) {
        book := i.shelf.books[i.index]
        i.index++
        if book.Author == i.author {
            return book
        }
    }
    return nil
}

func (i *AuthorFilterIterator) Reset() {
    i.index = 0
}

// ============ 图书馆 ============

// Library 图书馆
type Library struct {
    shelves map[string]*BookShelf
}

func NewLibrary() *Library {
    return &Library{
        shelves: make(map[string]*BookShelf),
    }
}

func (l *Library) AddShelf(shelf *BookShelf) {
    l.shelves[shelf.GetName()] = shelf
}

func (l *Library) GetShelf(name string) *BookShelf {
    return l.shelves[name]
}

// LibraryIterator 图书馆迭代器（遍历所有书架）
type LibraryIterator struct {
    library     *Library
    shelfNames  []string
    currentShelf int
    shelfIterator Iterator
}

func NewLibraryIterator(library *Library) *LibraryIterator {
    names := make([]string, 0, len(library.shelves))
    for name := range library.shelves {
        names = append(names, name)
    }

    iter := &LibraryIterator{
        library:     library,
        shelfNames:  names,
        currentShelf: 0,
    }
    iter.initShelfIterator()
    return iter
}

func (i *LibraryIterator) initShelfIterator() {
    if i.currentShelf < len(i.shelfNames) {
        shelf := i.library.shelves[i.shelfNames[i.currentShelf]]
        i.shelfIterator = NewBookIterator(shelf)
    }
}

func (i *LibraryIterator) HasNext() bool {
    for i.currentShelf < len(i.shelfNames) {
        if i.shelfIterator.HasNext() {
            return true
        }
        i.currentShelf++
        if i.currentShelf < len(i.shelfNames) {
            i.initShelfIterator()
        }
    }
    return false
}

func (i *LibraryIterator) Next() interface{} {
    for i.currentShelf < len(i.shelfNames) {
        if i.shelfIterator.HasNext() {
            return i.shelfIterator.Next()
        }
        i.currentShelf++
        if i.currentShelf < len(i.shelfNames) {
            i.initShelfIterator()
        }
    }
    return nil
}

func (i *LibraryIterator) Reset() {
    i.currentShelf = 0
    i.initShelfIterator()
}

// ============ Go风格迭代器 ============

// IteratorFunc 迭代器函数类型
type IteratorFunc func(yield func(interface{}) bool)

// BookShelf.Range Go风格的Range方法
func (s *BookShelf) Range(yield func(*Book) bool) {
    for _, book := range s.books {
        if !yield(book) {
            break
        }
    }
}

// BookShelf.Filter 过滤方法
func (s *BookShelf) Filter(predicate func(*Book) bool) []*Book {
    var result []*Book
    for _, book := range s.books {
        if predicate(book) {
            result = append(result, book)
        }
    }
    return result
}

// BookShelf.Map 映射方法
func (s *BookShelf) Map(transform func(*Book) string) []string {
    var result []string
    for _, book := range s.books {
        result = append(result, transform(book))
    }
    return result
}

func main() {
    // 创建书架
    fictionShelf := NewBookShelf("小说")
    fictionShelf.AddBook(&Book{"三体", "刘慈欣", "9787536692930", 39.99})
    fictionShelf.AddBook(&Book{"流浪地球", "刘慈欣", "9787536692931", 29.99})
    fictionShelf.AddBook(&Book{"哈利波特", "J.K.罗琳", "9780747532699", 49.99})

    techShelf := NewBookShelf("技术")
    techShelf.AddBook(&Book{"Go语言编程", "许式伟", "9787115547017", 79.99})
    techShelf.AddBook(&Book{"深入理解计算机系统", "Randal E. Bryant", "9787111544935", 139.99})

    // 创建图书馆
    library := NewLibrary()
    library.AddShelf(fictionShelf)
    library.AddShelf(techShelf)

    fmt.Println("========== 正向迭代 ==========")
    iter := NewBookIterator(fictionShelf)
    for iter.HasNext() {
        book := iter.Next().(*Book)
        fmt.Println(book)
    }

    fmt.Println("\n========== 反向迭代 ==========")
    reverseIter := NewReverseBookIterator(fictionShelf)
    for reverseIter.HasNext() {
        book := reverseIter.Next().(*Book)
        fmt.Println(book)
    }

    fmt.Println("\n========== 过滤迭代（刘慈欣的书） ==========")
    filterIter := NewAuthorFilterIterator(fictionShelf, "刘慈欣")
    for filterIter.HasNext() {
        book := filterIter.Next().(*Book)
        fmt.Println(book)
    }

    fmt.Println("\n========== 图书馆整体迭代 ==========")
    libIter := NewLibraryIterator(library)
    for libIter.HasNext() {
        book := libIter.Next().(*Book)
        fmt.Println(book)
    }

    fmt.Println("\n========== Go风格迭代 ==========")
    fmt.Println("遍历小说书架:")
    fictionShelf.Range(func(book *Book) bool {
        fmt.Println(book)
        return true
    })

    fmt.Println("\n过滤价格大于40的书:")
    expensiveBooks := fictionShelf.Filter(func(b *Book) bool {
        return b.Price > 40
    })
    for _, book := range expensiveBooks {
        fmt.Println(book)
    }

    fmt.Println("\n获取所有书名:")
    titles := fictionShelf.Map(func(b *Book) string {
        return b.Title
    })
    fmt.Println(strings.Join(titles, ", "))
}
```

#### 反例说明

**错误1：暴露内部实现**

```go
// 错误：暴露了内部切片
type Collection struct {
    Items []Item  // 公开字段
}
```

**错误2：迭代器修改集合**

```go
// 错误：迭代时不应该修改集合
func (i *Iterator) Remove() {
    // 可能导致索引错误
}
```

**错误3：不支持并发**

```go
// 错误：并发访问需要同步
type Iterator struct {
    collection []Item
    index      int  // 并发修改不安全
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 统一遍历接口 | 增加类数量 |
| 隐藏集合实现 | 简单集合可能过度设计 |
| 支持多种遍历 | 性能开销 |
| 支持并发遍历 | 需要额外同步 |

**适用场景**:

- 复杂集合遍历
- 多种遍历方式
- 需要隐藏实现
- 树形结构遍历

**Go语言特化**:

- `range`关键字内置支持
- Channel作为迭代器
- 函数式编程风格（Filter/Map/Reduce）

---

### 17. 中介者模式（Mediator）

#### 概念定义

中介者模式用一个中介对象来封装一系列的对象交互。中介者使各对象不需要显式地相互引用，从而使其耦合松散，而且可以独立地改变它们之间的交互。

#### 意图与动机

- **解耦**: 减少对象之间的直接引用
- **集中控制**: 集中管理对象交互
- **简化通信**: 简化多对多通信

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Mediator                               │◄──── 中介者接口
│  + Send(message, colleague)                                 │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  ConcreteMediator                           │
│  - colleague1: Colleague                                    │
│  - colleague2: Colleague                                    │
├─────────────────────────────────────────────────────────────┤
│  + Send(message, colleague)                                 │
│    if colleague == colleague1:                              │
│      colleague2.Receive(message)                            │
│    else:                                                    │
│      colleague1.Receive(message)                            │
└──────────────────────────┬──────────────────────────────────┘
                           │ manages
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│   Colleague   │  │   Colleague1  │  │   Colleague2  │
│  - mediator   │  │  + Send()     │  │  + Send()     │
├───────────────┤  │    mediator.Send()              │
│  + Send()     │  │  + Receive()  │  │  + Receive()  │
│  + Receive()  │  └───────────────┘  └───────────────┘
└───────────────┘
```

#### Go语言实现

```go
package mediator

// Mediator 中介者接口
type Mediator interface {
    Send(message string, colleague Colleague)
}

// Colleague 同事接口
type Colleague interface {
    Send(message string)
    Receive(message string)
    GetName() string
}

// ConcreteMediator 具体中介者
type ConcreteMediator struct {
    colleagues []Colleague
}

func NewConcreteMediator() *ConcreteMediator {
    return &ConcreteMediator{
        colleagues: make([]Colleague, 0),
    }
}

func (m *ConcreteMediator) Register(colleague Colleague) {
    m.colleagues = append(m.colleagues, colleague)
}

func (m *ConcreteMediator) Send(message string, sender Colleague) {
    for _, colleague := range m.colleagues {
        if colleague != sender {
            colleague.Receive(fmt.Sprintf("From %s: %s", sender.GetName(), message))
        }
    }
}

// ConcreteColleague 具体同事
type ConcreteColleague struct {
    name     string
    mediator Mediator
}

func NewConcreteColleague(name string, mediator Mediator) *ConcreteColleague {
    return &ConcreteColleague{
        name:     name,
        mediator: mediator,
    }
}

func (c *ConcreteColleague) Send(message string) {
    c.mediator.Send(message, c)
}

func (c *ConcreteColleague) Receive(message string) {
    fmt.Printf("%s received: %s\n", c.name, message)
}

func (c *ConcreteColleague) GetName() string {
    return c.name
}
```

#### 完整示例：聊天室系统

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

// ============ 消息 ============

type Message struct {
    From      string
    To        string
    Content   string
    Timestamp time.Time
    Type      string // broadcast, private, system
}

func (m *Message) String() string {
    timeStr := m.Timestamp.Format("15:04:05")
    switch m.Type {
    case "broadcast":
        return fmt.Sprintf("[%s] %s: %s", timeStr, m.From, m.Content)
    case "private":
        return fmt.Sprintf("[%s] %s -> %s: %s", timeStr, m.From, m.To, m.Content)
    case "system":
        return fmt.Sprintf("[%s] [系统] %s", timeStr, m.Content)
    default:
        return fmt.Sprintf("[%s] %s: %s", timeStr, m.From, m.Content)
    }
}

// ============ 用户 ============

type User struct {
    name     string
    room     ChatRoom
    isOnline bool
}

func NewUser(name string) *User {
    return &User{
        name:     name,
        isOnline: false,
    }
}

func (u *User) GetName() string {
    return u.name
}

func (u *User) Join(room ChatRoom) {
    u.room = room
    u.isOnline = true
    room.Register(u)
}

func (u *User) Leave() {
    if u.room != nil {
        u.room.Unregister(u)
        u.isOnline = false
        u.room = nil
    }
}

func (u *User) SendBroadcast(message string) {
    if u.room != nil {
        u.room.Broadcast(&Message{
            From:      u.name,
            Content:   message,
            Timestamp: time.Now(),
            Type:      "broadcast",
        })
    }
}

func (u *User) SendPrivate(to, message string) {
    if u.room != nil {
        u.room.PrivateMessage(&Message{
            From:      u.name,
            To:        to,
            Content:   message,
            Timestamp: time.Now(),
            Type:      "private",
        })
    }
}

func (u *User) Receive(message *Message) {
    fmt.Printf("[%s] %s\n", u.name, message)
}

func (u *User) IsOnline() bool {
    return u.isOnline
}

// ============ 聊天室 ============

type ChatRoom interface {
    Register(user *User)
    Unregister(user *User)
    Broadcast(message *Message)
    PrivateMessage(message *Message)
    GetUsers() []*User
    GetUser(name string) *User
}

// ChatRoomImpl 聊天室实现
type ChatRoomImpl struct {
    name  string
    users map[string]*User
}

func NewChatRoom(name string) ChatRoom {
    return &ChatRoomImpl{
        name:  name,
        users: make(map[string]*User),
    }
}

func (c *ChatRoomImpl) Register(user *User) {
    c.users[user.GetName()] = user
    c.Broadcast(&Message{
        From:      "系统",
        Content:   fmt.Sprintf("%s 加入了聊天室", user.GetName()),
        Timestamp: time.Now(),
        Type:      "system",
    })
}

func (c *ChatRoomImpl) Unregister(user *User) {
    delete(c.users, user.GetName())
    c.Broadcast(&Message{
        From:      "系统",
        Content:   fmt.Sprintf("%s 离开了聊天室", user.GetName()),
        Timestamp: time.Now(),
        Type:      "system",
    })
}

func (c *ChatRoomImpl) Broadcast(message *Message) {
    fmt.Printf("\n=== %s ===\n", c.name)
    for _, user := range c.users {
        user.Receive(message)
    }
}

func (c *ChatRoomImpl) PrivateMessage(message *Message) {
    if toUser, ok := c.users[message.To]; ok {
        toUser.Receive(message)
        // 发送者也收到确认
        if fromUser, ok := c.users[message.From]; ok {
            fromUser.Receive(&Message{
                From:      "系统",
                Content:   fmt.Sprintf("私信已发送给 %s", message.To),
                Timestamp: time.Now(),
                Type:      "system",
            })
        }
    } else {
        if fromUser, ok := c.users[message.From]; ok {
            fromUser.Receive(&Message{
                From:      "系统",
                Content:   fmt.Sprintf("用户 %s 不存在", message.To),
                Timestamp: time.Now(),
                Type:      "system",
            })
        }
    }
}

func (c *ChatRoomImpl) GetUsers() []*User {
    var users []*User
    for _, user := range c.users {
        users = append(users, user)
    }
    return users
}

func (c *ChatRoomImpl) GetUser(name string) *User {
    return c.users[name]
}

// ============ 高级聊天室（带功能） ============

type AdvancedChatRoom struct {
    ChatRoom
    bannedWords []string
    history     []*Message
    maxHistory  int
}

func NewAdvancedChatRoom(name string) *AdvancedChatRoom {
    return &AdvancedChatRoom{
        ChatRoom:    NewChatRoom(name),
        bannedWords: []string{"脏话", "spam"},
        history:     make([]*Message, 0),
        maxHistory:  100,
    }
}

func (a *AdvancedChatRoom) filterContent(content string) string {
    filtered := content
    for _, word := range a.bannedWords {
        filtered = strings.ReplaceAll(filtered, word, "***")
    }
    return filtered
}

func (a *AdvancedChatRoom) Broadcast(message *Message) {
    // 过滤内容
    message.Content = a.filterContent(message.Content)

    // 保存历史
    a.history = append(a.history, message)
    if len(a.history) > a.maxHistory {
        a.history = a.history[1:]
    }

    // 调用父类方法
    a.ChatRoom.Broadcast(message)
}

func (a *AdvancedChatRoom) GetHistory() []*Message {
    return a.history
}

func main() {
    // 创建聊天室
    room := NewChatRoom("Go程序员交流群")

    // 创建用户
    alice := NewUser("Alice")
    bob := NewUser("Bob")
    charlie := NewUser("Charlie")

    // 用户加入
    fmt.Println("========== 用户加入 ==========")
    alice.Join(room)
    bob.Join(room)
    charlie.Join(room)

    // 广播消息
    fmt.Println("\n========== 广播消息 ==========")
    alice.SendBroadcast("大家好！")
    bob.SendBroadcast("欢迎Alice！")

    // 私信
    fmt.Println("\n========== 私信 ==========")
    alice.SendPrivate("Bob", "你好Bob，有个问题想请教你")
    bob.SendPrivate("Alice", "没问题，请问")

    // 私信给不存在的用户
    alice.SendPrivate("David", "你好")

    // 用户离开
    fmt.Println("\n========== 用户离开 ==========")
    charlie.Leave()

    // 高级聊天室
    fmt.Println("\n========== 高级聊天室 ==========")
    advancedRoom := NewAdvancedChatRoom("高级群")

    dave := NewUser("Dave")
    eve := NewUser("Eve")

    dave.Join(advancedRoom)
    eve.Join(advancedRoom)

    dave.SendBroadcast("这是一条正常消息")
    eve.SendBroadcast("这条消息包含脏话内容")

    fmt.Println("\n历史记录:")
    for _, msg := range advancedRoom.GetHistory() {
        fmt.Println(msg)
    }
}
```

#### 反例说明

**错误1：中介者成为上帝对象**

```go
// 错误：中介者不应该包含所有业务逻辑
type Mediator struct {
    // 包含大量业务逻辑
}
```

**错误2：同事直接通信**

```go
// 错误：应该通过中介者通信
type Colleague1 struct {
    colleague2 *Colleague2  // 直接引用
}
```

**错误3：中介者过于复杂**

```go
// 错误：中介者不应该知道太多细节
type Mediator struct {
    // 不应该包含具体业务逻辑
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 减少对象间耦合 | 中介者可能变得复杂 |
| 集中控制交互 | 中介者成为瓶颈 |
| 简化对象协议 | 过度集中 |
| 易于修改交互 | 难以维护 |

**适用场景**:

- 聊天系统
- 事件总线
- MVC控制器
- 协调多个对象

**Go语言特化**:

- Channel作为中介
- 事件驱动架构
- 发布-订阅模式

---

### 18. 备忘录模式（Memento）

#### 概念定义

备忘录模式在不破坏封装性的前提下，捕获一个对象的内部状态，并在该对象之外保存这个状态。这样以后就可将该对象恢复到原先保存的状态。

#### 意图与动机

- **状态保存**: 保存和恢复对象状态
- **撤销功能**: 支持撤销操作
- **快照**: 创建对象快照

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Originator                             │
│  - state: State                                             │
├─────────────────────────────────────────────────────────────┤
│  + SetState(state)                                          │
│  + CreateMemento() Memento                                  │
│    return Memento{state}                                    │
│  + RestoreMemento(memento)                                  │
│    state = memento.GetState()                               │
└──────────────────────────┬──────────────────────────────────┘
                           │ creates
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Memento                                │◄──── 备忘录
│  - state: State                                             │
├─────────────────────────────────────────────────────────────┤
│  + GetState() State                                         │
└──────────────────────────┬──────────────────────────────────┘
                           │ managed by
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     Caretaker                               │
│  - mementos: []Memento                                      │
│  - current: int                                             │
├─────────────────────────────────────────────────────────────┤
│  + AddMemento(memento)                                      │
│  + GetMemento(index) Memento                                │
│  + Undo() Memento                                           │
│  + Redo() Memento                                           │
└─────────────────────────────────────────────────────────────┘
```

#### Go语言实现

```go
package memento

// Memento 备忘录
type Memento struct {
    state string
}

func NewMemento(state string) *Memento {
    return &Memento{state: state}
}

func (m *Memento) GetState() string {
    return m.state
}

// Originator 原发器
type Originator struct {
    state string
}

func (o *Originator) SetState(state string) {
    o.state = state
}

func (o *Originator) GetState() string {
    return o.state
}

func (o *Originator) CreateMemento() *Memento {
    return NewMemento(o.state)
}

func (o *Originator) RestoreMemento(memento *Memento) {
    o.state = memento.GetState()
}

// Caretaker 负责人
type Caretaker struct {
    mementos []*Memento
}

func NewCaretaker() *Caretaker {
    return &Caretaker{
        mementos: make([]*Memento, 0),
    }
}

func (c *Caretaker) AddMemento(memento *Memento) {
    c.mementos = append(c.mementos, memento)
}

func (c *Caretaker) GetMemento(index int) *Memento {
    if index >= 0 && index < len(c.mementos) {
        return c.mementos[index]
    }
    return nil
}
```

#### 完整示例：文档编辑器撤销系统

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// ============ 文档状态 ============

type DocumentState struct {
    Title      string
    Content    string
    CursorPos  int
    ModifiedAt time.Time
}

func (d *DocumentState) Clone() *DocumentState {
    return &DocumentState{
        Title:      d.Title,
        Content:    d.Content,
        CursorPos:  d.CursorPos,
        ModifiedAt: d.ModifiedAt,
    }
}

func (d *DocumentState) String() string {
    return fmt.Sprintf("Document{Title: %s, Content: %s..., Cursor: %d}",
        d.Title, d.Content[:min(20, len(d.Content))], d.CursorPos)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// ============ 备忘录 ============

type DocumentMemento struct {
    state      *DocumentState
    savedAt    time.Time
    description string
}

func NewDocumentMemento(state *DocumentState, description string) *DocumentMemento {
    return &DocumentMemento{
        state:       state.Clone(),
        savedAt:     time.Now(),
        description: description,
    }
}

func (m *DocumentMemento) GetState() *DocumentState {
    return m.state.Clone()
}

func (m *DocumentMemento) GetDescription() string {
    return m.description
}

func (m *DocumentMemento) GetSavedAt() time.Time {
    return m.savedAt
}

// ============ 文档（原发器） ============

type Document struct {
    state *DocumentState
}

func NewDocument(title string) *Document {
    return &Document{
        state: &DocumentState{
            Title:      title,
            Content:    "",
            CursorPos:  0,
            ModifiedAt: time.Now(),
        },
    }
}

func (d *Document) SetTitle(title string) {
    d.state.Title = title
    d.state.ModifiedAt = time.Now()
}

func (d *Document) GetTitle() string {
    return d.state.Title
}

func (d *Document) Insert(text string, position int) {
    if position < 0 || position > len(d.state.Content) {
        position = len(d.state.Content)
    }
    d.state.Content = d.state.Content[:position] + text + d.state.Content[position:]
    d.state.CursorPos = position + len(text)
    d.state.ModifiedAt = time.Now()
}

func (d *Document) Delete(start, end int) string {
    if start < 0 {
        start = 0
    }
    if end > len(d.state.Content) {
        end = len(d.state.Content)
    }
    deleted := d.state.Content[start:end]
    d.state.Content = d.state.Content[:start] + d.state.Content[end:]
    d.state.CursorPos = start
    d.state.ModifiedAt = time.Now()
    return deleted
}

func (d *Document) GetContent() string {
    return d.state.Content
}

func (d *Document) SetContent(content string) {
    d.state.Content = content
    d.state.ModifiedAt = time.Now()
}

func (d *Document) GetCursorPos() int {
    return d.state.CursorPos
}

func (d *Document) SetCursorPos(pos int) {
    d.state.CursorPos = pos
}

func (d *Document) CreateMemento(description string) *DocumentMemento {
    return NewDocumentMemento(d.state, description)
}

func (d *Document) RestoreMemento(memento *DocumentMemento) {
    d.state = memento.GetState()
}

func (d *Document) String() string {
    return d.state.String()
}

// ============ 历史管理器（负责人） ============

type HistoryManager struct {
    history    []*DocumentMemento
    currentIdx int
    maxSize    int
}

func NewHistoryManager(maxSize int) *HistoryManager {
    return &HistoryManager{
        history:    make([]*DocumentMemento, 0),
        currentIdx: -1,
        maxSize:    maxSize,
    }
}

func (h *HistoryManager) Save(memento *DocumentMemento) {
    // 删除当前位置之后的历史
    if h.currentIdx < len(h.history)-1 {
        h.history = h.history[:h.currentIdx+1]
    }

    // 添加新状态
    h.history = append(h.history, memento)
    h.currentIdx++

    // 限制历史大小
    if len(h.history) > h.maxSize {
        h.history = h.history[1:]
        h.currentIdx--
    }
}

func (h *HistoryManager) Undo() *DocumentMemento {
    if h.currentIdx > 0 {
        h.currentIdx--
        return h.history[h.currentIdx]
    }
    return nil
}

func (h *HistoryManager) Redo() *DocumentMemento {
    if h.currentIdx < len(h.history)-1 {
    h.currentIdx++
        return h.history[h.currentIdx]
    }
    return nil
}

func (h *HistoryManager) CanUndo() bool {
    return h.currentIdx > 0
}

func (h *HistoryManager) CanRedo() bool {
    return h.currentIdx < len(h.history)-1
}

func (h *HistoryManager) GetHistory() []string {
    var result []string
    for i, m := range h.history {
        prefix := "  "
        if i == h.currentIdx {
            prefix = "► "
        }
        result = append(result, fmt.Sprintf("%s%s (%s)",
            prefix, m.GetDescription(), m.GetSavedAt().Format("15:04:05")))
    }
    return result
}

func (h *HistoryManager) Clear() {
    h.history = make([]*DocumentMemento, 0)
    h.currentIdx = -1
}

// ============ 编辑器 ============

type Editor struct {
    document *Document
    history  *HistoryManager
}

func NewEditor(title string) *Editor {
    doc := NewDocument(title)
    editor := &Editor{
        document: doc,
        history:  NewHistoryManager(50),
    }
    // 保存初始状态
    editor.Save("初始状态")
    return editor
}

func (e *Editor) Save(description string) {
    memento := e.document.CreateMemento(description)
    e.history.Save(memento)
}

func (e *Editor) Undo() bool {
    if !e.history.CanUndo() {
        fmt.Println("没有可撤销的操作")
        return false
    }

    memento := e.history.Undo()
    if memento != nil {
        e.document.RestoreMemento(memento)
        fmt.Printf("撤销: %s\n", memento.GetDescription())
        return true
    }
    return false
}

func (e *Editor) Redo() bool {
    if !e.history.CanRedo() {
        fmt.Println("没有可重做的操作")
        return false
    }

    memento := e.history.Redo()
    if memento != nil {
        e.document.RestoreMemento(memento)
        fmt.Printf("重做: %s\n", memento.GetDescription())
        return true
    }
    return false
}

func (e *Editor) Insert(text string) {
    e.document.Insert(text, e.document.GetCursorPos())
    e.Save(fmt.Sprintf("插入: %s", text))
}

func (e *Editor) Delete(length int) {
    pos := e.document.GetCursorPos()
    deleted := e.document.Delete(pos-length, pos)
    e.Save(fmt.Sprintf("删除: %s", deleted))
}

func (e *Editor) SetTitle(title string) {
    oldTitle := e.document.GetTitle()
    e.document.SetTitle(title)
    e.Save(fmt.Sprintf("修改标题: %s -> %s", oldTitle, title))
}

func (e *Editor) PrintState() {
    fmt.Println()
    fmt.Println("=" + "="*50)
    fmt.Printf("标题: %s\n", e.document.GetTitle())
    fmt.Printf("内容: %s\n", e.document.GetContent())
    fmt.Printf("光标位置: %d\n", e.document.GetCursorPos())
    fmt.Println("-" + "-"*50)
    fmt.Println("历史:")
    for _, h := range e.history.GetHistory() {
        fmt.Println(h)
    }
    fmt.Println("=" + "="*50)
}

// ============ 序列化备忘录（持久化） ============

type SerializableMemento struct {
    State       *DocumentState `json:"state"`
    SavedAt     time.Time      `json:"saved_at"`
    Description string         `json:"description"`
}

func (m *DocumentMemento) ToJSON() ([]byte, error) {
    sm := &SerializableMemento{
        State:       m.state,
        SavedAt:     m.savedAt,
        Description: m.description,
    }
    return json.Marshal(sm)
}

func MementoFromJSON(data []byte) (*DocumentMemento, error) {
    var sm SerializableMemento
    if err := json.Unmarshal(data, &sm); err != nil {
        return nil, err
    }
    return &DocumentMemento{
        state:       sm.State,
        savedAt:     sm.SavedAt,
        description: sm.Description,
    }, nil
}

func main() {
    editor := NewEditor("未命名文档")

    fmt.Println("========== 文档编辑器演示 ==========")
    editor.PrintState()

    // 编辑操作
    fmt.Println("\n========== 编辑操作 ==========")

    editor.Insert("Hello")
    editor.PrintState()

    editor.Insert(" World")
    editor.PrintState()

    editor.SetTitle("我的文档")
    editor.PrintState()

    editor.Insert("!")
    editor.PrintState()

    // 撤销操作
    fmt.Println("\n========== 撤销操作 ==========")
    editor.Undo()
    editor.PrintState()

    editor.Undo()
    editor.PrintState()

    // 重做操作
    fmt.Println("\n========== 重做操作 ==========")
    editor.Redo()
    editor.PrintState()

    // 新操作（会清除重做历史）
    fmt.Println("\n========== 新操作（清除重做历史） ==========")
    editor.Insert("!!!")
    editor.PrintState()

    // 此时无法重做
    editor.Redo()

    // 序列化示例
    fmt.Println("\n========== 序列化备忘录 ==========")
    memento := editor.document.CreateMemento("手动保存")
    jsonData, err := memento.ToJSON()
    if err != nil {
        fmt.Printf("序列化失败: %v\n", err)
    } else {
        fmt.Printf("序列化成功: %s\n", string(jsonData))
    }
}
```

#### 反例说明

**错误1：备忘录暴露原发器内部**

```go
// 错误：备忘录不应该暴露内部实现
type Memento struct {
    State *State  // 公开字段，破坏了封装
}
```

**错误2：备忘录过大**

```go
// 错误：不应该保存整个对象
type Memento struct {
    Document *Document  // 保存整个文档，内存占用大
}
```

**错误3：没有限制历史大小**

```go
// 错误：应该限制历史记录大小
type Caretaker struct {
    mementos []Memento  // 无限增长会导致内存溢出
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 保持封装性 | 内存开销 |
| 简化原发器 | 管理成本高 |
| 支持撤销/重做 | 历史记录维护 |
| 状态快照 | 深拷贝开销 |

**适用场景**:

- 编辑器撤销功能
- 游戏存档
- 事务回滚
- 配置快照

**Go语言特化**:

- 使用JSON序列化持久化
- 结构体复制实现深拷贝
- 接口隐藏实现细节

---

### 19. 观察者模式（Observer）

#### 概念定义

观察者模式定义对象间的一种一对多的依赖关系，当一个对象的状态发生改变时，所有依赖于它的对象都得到通知并被自动更新。

#### 意图与动机

- **松耦合**: 主题和观察者松耦合
- **事件通知**: 状态变化自动通知
- **广播通信**: 一对多通信

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Subject                                │◄──── 主题接口
│  - observers: []Observer                                    │
├─────────────────────────────────────────────────────────────┤
│  + Attach(observer)                                         │
│  + Detach(observer)                                         │
│  + Notify()                                                 │
│    for each observer:                                       │
│      observer.Update()                                      │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  ConcreteSubject                            │
│  - state: State                                             │
├─────────────────────────────────────────────────────────────┤
│  + GetState()                                               │
│  + SetState(state)                                          │
│    this.state = state                                       │
│    Notify()                                                 │
└──────────────────────────┬──────────────────────────────────┘
                           │ notifies
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│   Observer    │  │  ObserverA    │  │  ObserverB    │
│  + Update()   │  │  + Update()   │  │  + Update()   │
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

Go语言推荐使用Channel实现观察者模式。

```go
package observer

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
    GetState() string
}

// ConcreteSubject 具体主题
type ConcreteSubject struct {
    state     string
    observers []Observer
}

func NewConcreteSubject() *ConcreteSubject {
    return &ConcreteSubject{
        observers: make([]Observer, 0),
    }
}

func (s *ConcreteSubject) Attach(observer Observer) {
    s.observers = append(s.observers, observer)
}

func (s *ConcreteSubject) Detach(observer Observer) {
    for i, obs := range s.observers {
        if obs.GetID() == observer.GetID() {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            return
        }
    }
}

func (s *ConcreteSubject) Notify() {
    for _, observer := range s.observers {
        observer.Update(s)
    }
}

func (s *ConcreteSubject) GetState() string {
    return s.state
}

func (s *ConcreteSubject) SetState(state string) {
    s.state = state
    s.Notify()
}

// ConcreteObserver 具体观察者
type ConcreteObserver struct {
    id   string
    name string
}

func NewConcreteObserver(id, name string) *ConcreteObserver {
    return &ConcreteObserver{id: id, name: name}
}

func (o *ConcreteObserver) Update(subject Subject) {
    fmt.Printf("%s received update: %s\n", o.name, subject.GetState())
}

func (o *ConcreteObserver) GetID() string {
    return o.id
}
```

#### 完整示例：事件总线系统

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// ============ 事件 ============

type EventType string

const (
    EventUserLogin    EventType = "user:login"
    EventUserLogout   EventType = "user:logout"
    EventOrderCreated EventType = "order:created"
    EventOrderPaid    EventType = "order:paid"
    EventDataChanged  EventType = "data:changed"
)

type Event struct {
    Type      EventType
    Data      interface{}
    Timestamp time.Time
    Source    string
}

func NewEvent(eventType EventType, data interface{}, source string) *Event {
    return &Event{
        Type:      eventType,
        Data:      data,
        Timestamp: time.Now(),
        Source:    source,
    }
}

func (e *Event) String() string {
    return fmt.Sprintf("[%s] %s from %s: %v",
        e.Timestamp.Format("15:04:05"), e.Type, e.Source, e.Data)
}

// ============ 观察者接口 ============

type EventListener interface {
    OnEvent(event *Event)
    GetListenerID() string
    Supports(eventType EventType) bool
}

// ============ 事件总线 ============

type EventBus struct {
    listeners map[EventType][]EventListener
    mu        sync.RWMutex
    history   []*Event
    maxHistory int
}

func NewEventBus() *EventBus {
    return &EventBus{
        listeners:  make(map[EventType][]EventListener),
        history:    make([]*Event, 0),
        maxHistory: 100,
    }
}

func (b *EventBus) Subscribe(listener EventListener, eventTypes ...EventType) {
    b.mu.Lock()
    defer b.mu.Unlock()

    for _, eventType := range eventTypes {
        b.listeners[eventType] = append(b.listeners[eventType], listener)
    }

    fmt.Printf("[EventBus] %s 订阅了 %v\n", listener.GetListenerID(), eventTypes)
}

func (b *EventBus) Unsubscribe(listener EventListener, eventTypes ...EventType) {
    b.mu.Lock()
    defer b.mu.Unlock()

    for _, eventType := range eventTypes {
        listeners := b.listeners[eventType]
        for i, l := range listeners {
            if l.GetListenerID() == listener.GetListenerID() {
                b.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
                break
            }
        }
    }

    fmt.Printf("[EventBus] %s 取消订阅了 %v\n", listener.GetListenerID(), eventTypes)
}

func (b *EventBus) Publish(event *Event) {
    b.mu.RLock()
    listeners := b.listeners[event.Type]
    b.mu.RUnlock()

    // 保存历史
    b.addToHistory(event)

    fmt.Printf("[EventBus] 发布事件: %s\n", event)

    // 异步通知监听器
    var wg sync.WaitGroup
    for _, listener := range listeners {
        if listener.Supports(event.Type) {
            wg.Add(1)
            go func(l EventListener) {
                defer wg.Done()
                l.OnEvent(event)
            }(listener)
        }
    }
    wg.Wait()
}

func (b *EventBus) addToHistory(event *Event) {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.history = append(b.history, event)
    if len(b.history) > b.maxHistory {
        b.history = b.history[1:]
    }
}

func (b *EventBus) GetHistory(eventType EventType) []*Event {
    b.mu.RLock()
    defer b.mu.RUnlock()

    var result []*Event
    for _, event := range b.history {
        if event.Type == eventType {
            result = append(result, event)
        }
    }
    return result
}

// ============ 具体监听器 ============

// LoggerListener 日志监听器
type LoggerListener struct {
    id string
}

func NewLoggerListener() *LoggerListener {
    return &LoggerListener{id: "logger"}
}

func (l *LoggerListener) OnEvent(event *Event) {
    fmt.Printf("[Logger] %s\n", event)
}

func (l *LoggerListener) GetListenerID() string {
    return l.id
}

func (l *LoggerListener) Supports(eventType EventType) bool {
    return true  // 监听所有事件
}

// EmailListener 邮件通知监听器
type EmailListener struct {
    id string
}

func NewEmailListener() *EmailListener {
    return &EmailListener{id: "email"}
}

func (e *EmailListener) OnEvent(event *Event) {
    switch event.Type {
    case EventUserLogin:
        fmt.Printf("[Email] 发送登录通知: %v\n", event.Data)
    case EventOrderPaid:
        fmt.Printf("[Email] 发送订单确认: %v\n", event.Data)
    }
}

func (e *EmailListener) GetListenerID() string {
    return e.id
}

func (e *EmailListener) Supports(eventType EventType) bool {
    return eventType == EventUserLogin || eventType == EventOrderPaid
}

// AnalyticsListener 分析监听器
type AnalyticsListener struct {
    id       string
    counters map[EventType]int
    mu       sync.Mutex
}

func NewAnalyticsListener() *AnalyticsListener {
    return &AnalyticsListener{
        id:       "analytics",
        counters: make(map[EventType]int),
    }
}

func (a *AnalyticsListener) OnEvent(event *Event) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.counters[event.Type]++
    fmt.Printf("[Analytics] %s 计数: %d\n", event.Type, a.counters[event.Type])
}

func (a *AnalyticsListener) GetListenerID() string {
    return a.id
}

func (a *AnalyticsListener) Supports(eventType EventType) bool {
    return true
}

func (a *AnalyticsListener) GetStats() map[EventType]int {
    a.mu.Lock()
    defer a.mu.Unlock()

    stats := make(map[EventType]int)
    for k, v := range a.counters {
        stats[k] = v
    }
    return stats
}

// CacheInvalidationListener 缓存失效监听器
type CacheInvalidationListener struct {
    id    string
    cache map[string]interface{}
}

func NewCacheInvalidationListener() *CacheInvalidationListener {
    return &CacheInvalidationListener{
        id:    "cache",
        cache: make(map[string]interface{}),
    }
}

func (c *CacheInvalidationListener) OnEvent(event *Event) {
    if event.Type == EventDataChanged {
        key := event.Data.(string)
        delete(c.cache, key)
        fmt.Printf("[Cache] 缓存失效: %s\n", key)
    }
}

func (c *CacheInvalidationListener) GetListenerID() string {
    return c.id
}

func (c *CacheInvalidationListener) Supports(eventType EventType) bool {
    return eventType == EventDataChanged
}

func (c *CacheInvalidationListener) Set(key string, value interface{}) {
    c.cache[key] = value
}

// ============ 使用示例 ============

func main() {
    // 创建事件总线
    eventBus := NewEventBus()

    // 创建监听器
    logger := NewLoggerListener()
    email := NewEmailListener()
    analytics := NewAnalyticsListener()
    cache := NewCacheInvalidationListener()

    // 订阅事件
    eventBus.Subscribe(logger, EventUserLogin, EventUserLogout, EventOrderCreated, EventOrderPaid, EventDataChanged)
    eventBus.Subscribe(email, EventUserLogin, EventOrderPaid)
    eventBus.Subscribe(analytics, EventUserLogin, EventOrderCreated, EventOrderPaid)
    eventBus.Subscribe(cache, EventDataChanged)

    fmt.Println("\n========== 发布事件 ==========")

    // 发布用户登录事件
    eventBus.Publish(NewEvent(EventUserLogin, map[string]string{
        "user_id":  "123",
        "username": "alice",
    }, "auth-service"))

    // 发布订单创建事件
    eventBus.Publish(NewEvent(EventOrderCreated, map[string]interface{}{
        "order_id": "ORD-001",
        "amount":   99.99,
    }, "order-service"))

    // 发布订单支付事件
    eventBus.Publish(NewEvent(EventOrderPaid, map[string]interface{}{
        "order_id": "ORD-001",
        "amount":   99.99,
    }, "payment-service"))

    // 发布数据变更事件
    eventBus.Publish(NewEvent(EventDataChanged, "user:123", "user-service"))

    // 查看统计
    fmt.Println("\n========== 统计信息 ==========")
    stats := analytics.GetStats()
    for eventType, count := range stats {
        fmt.Printf("%s: %d\n", eventType, count)
    }

    // 取消订阅
    fmt.Println("\n========== 取消订阅 ==========")
    eventBus.Unsubscribe(email, EventUserLogin)

    // 再次发布事件
    fmt.Println("\n========== 再次发布事件 ==========")
    eventBus.Publish(NewEvent(EventUserLogin, map[string]string{
        "user_id":  "456",
        "username": "bob",
    }, "auth-service"))

    // 查看历史
    fmt.Println("\n========== 事件历史 ==========")
    loginHistory := eventBus.GetHistory(EventUserLogin)
    for _, event := range loginHistory {
        fmt.Println(event)
    }
}
```

#### 反例说明

**错误1：循环依赖**

```go
// 错误：A观察B，B观察A会导致循环通知
subjectA.Attach(observerB)
subjectB.Attach(observerA)
```

**错误2：通知时修改观察者列表**

```go
// 错误：通知时修改列表会导致问题
func (s *Subject) Notify() {
    for _, observer := range s.observers {
        observer.Update(s)
        s.Detach(observer)  // 不要在遍历时修改
    }
}
```

**错误3：同步通知阻塞**

```go
// 错误：同步通知可能阻塞
func (s *Subject) Notify() {
    for _, observer := range s.observers {
        observer.Update(s)  // 如果某个观察者慢，会阻塞其他观察者
    }
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 松耦合 | 可能导致循环依赖 |
| 支持广播 | 通知顺序不确定 |
| 动态添加观察者 | 性能问题（大量观察者） |
| 符合开闭原则 | 调试困难 |

**适用场景**:

- 事件驱动系统
- MVC模式
- 消息队列
- 数据绑定

**Go语言特化**:

- Channel作为观察者通信
- `sync.RWMutex`保护观察者列表
- 协程异步通知

---

### 20. 状态模式（State）

#### 概念定义

状态模式允许对象在内部状态改变时改变它的行为。对象看起来似乎修改了它的类。

#### 意图与动机

- **消除条件语句**: 用多态替代复杂的条件判断
- **状态封装**: 将状态行为封装到独立类
- **状态转换**: 显式定义状态转换规则

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Context                                │
│  - state: State                                             │
├─────────────────────────────────────────────────────────────┤
│  + Request()                                                │
│    state.Handle()                                           │
│  + SetState(state)                                          │
│  + GetState() State                                         │
└──────────────────────────┬──────────────────────────────────┘
                           │ delegates to
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                       State                                 │◄──── 状态接口
│  + Handle(context)                                          │
│  + Enter()                                                  │
│  + Exit()                                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  StateA       │  │  StateB       │  │  StateC       │
│  + Handle()   │  │  + Handle()   │  │  + Handle()   │
│    context.Set│  │    context.Set│  │    context.Set│
│    State(B)   │  │    State(C)   │  │    State(A)   │
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

```go
package state

// State 状态接口
type State interface {
    Handle(context *Context)
    GetName() string
}

// Context 上下文
type Context struct {
    state State
}

func NewContext(initialState State) *Context {
    return &Context{state: initialState}
}

func (c *Context) SetState(state State) {
    c.state = state
}

func (c *Context) GetState() State {
    return c.state
}

func (c *Context) Request() {
    c.state.Handle(c)
}

// ConcreteStateA 具体状态A
type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) {
    fmt.Println("StateA handling request, transitioning to StateB")
    context.SetState(&ConcreteStateB{})
}

func (s *ConcreteStateA) GetName() string {
    return "StateA"
}

// ConcreteStateB 具体状态B
type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) {
    fmt.Println("StateB handling request, transitioning to StateA")
    context.SetState(&ConcreteStateA{})
}

func (s *ConcreteStateB) GetName() string {
    return "StateB"
}
```

#### 完整示例：订单状态机

```go
package main

import (
    "fmt"
    "time"
)

// ============ 订单状态接口 ============

type OrderState interface {
    Pay(order *Order, amount float64) error
    Ship(order *Order, trackingNumber string) error
    Deliver(order *Order) error
    Cancel(order *Order) error
    GetName() string
}

// ============ 订单 ============

type Order struct {
    ID             string
    Amount         float64
    State          OrderState
    TrackingNumber string
    CreatedAt      time.Time
    PaidAt         *time.Time
    ShippedAt      *time.Time
    DeliveredAt    *time.Time
    CancelledAt    *time.Time
}

func NewOrder(id string, amount float64) *Order {
    order := &Order{
        ID:        id,
        Amount:    amount,
        CreatedAt: time.Now(),
    }
    order.SetState(&PendingState{})
    return order
}

func (o *Order) SetState(state OrderState) {
    fmt.Printf("订单 %s: %s -> %s\n", o.ID, o.State.GetName(), state.GetName())
    o.State = state
}

func (o *Order) Pay(amount float64) error {
    return o.State.Pay(o, amount)
}

func (o *Order) Ship(trackingNumber string) error {
    return o.State.Ship(o, trackingNumber)
}

func (o *Order) Deliver() error {
    return o.State.Deliver(o)
}

func (o *Order) Cancel() error {
    return o.State.Cancel(o)
}

func (o *Order) GetStatus() string {
    return o.State.GetName()
}

func (o *Order) String() string {
    return fmt.Sprintf("订单[%s] 金额: %.2f 状态: %s", o.ID, o.Amount, o.GetStatus())
}

// ============ 具体状态 ============

// PendingState 待支付状态
type PendingState struct{}

func (s *PendingState) Pay(order *Order, amount float64) error {
    if amount < order.Amount {
        return fmt.Errorf("支付金额不足")
    }

    now := time.Now()
    order.PaidAt = &now
    order.SetState(&PaidState{})
    fmt.Printf("订单 %s 支付成功\n", order.ID)
    return nil
}

func (s *PendingState) Ship(order *Order, trackingNumber string) error {
    return fmt.Errorf("订单未支付，无法发货")
}

func (s *PendingState) Deliver(order *Order) error {
    return fmt.Errorf("订单未支付，无法配送")
}

func (s *PendingState) Cancel(order *Order) error {
    now := time.Now()
    order.CancelledAt = &now
    order.SetState(&CancelledState{})
    fmt.Printf("订单 %s 已取消\n", order.ID)
    return nil
}

func (s *PendingState) GetName() string {
    return "待支付"
}

// PaidState 已支付状态
type PaidState struct{}

func (s *PaidState) Pay(order *Order, amount float64) error {
    return fmt.Errorf("订单已支付")
}

func (s *PaidState) Ship(order *Order, trackingNumber string) error {
    order.TrackingNumber = trackingNumber
    now := time.Now()
    order.ShippedAt = &now
    order.SetState(&ShippedState{})
    fmt.Printf("订单 %s 已发货，物流单号: %s\n", order.ID, trackingNumber)
    return nil
}

func (s *PaidState) Deliver(order *Order) error {
    return fmt.Errorf("订单未发货，无法配送")
}

func (s *PaidState) Cancel(order *Order) error {
    now := time.Now()
    order.CancelledAt = &now
    order.SetState(&CancelledState{})
    fmt.Printf("订单 %s 已取消，将退款\n", order.ID)
    return nil
}

func (s *PaidState) GetName() string {
    return "已支付"
}

// ShippedState 已发货状态
type ShippedState struct{}

func (s *ShippedState) Pay(order *Order, amount float64) error {
    return fmt.Errorf("订单已支付")
}

func (s *ShippedState) Ship(order *Order, trackingNumber string) error {
    return fmt.Errorf("订单已发货")
}

func (s *ShippedState) Deliver(order *Order) error {
    now := time.Now()
    order.DeliveredAt = &now
    order.SetState(&DeliveredState{})
    fmt.Printf("订单 %s 已送达\n", order.ID)
    return nil
}

func (s *ShippedState) Cancel(order *Order) error {
    return fmt.Errorf("订单已发货，无法取消")
}

func (s *ShippedState) GetName() string {
    return "已发货"
}

// DeliveredState 已送达状态
type DeliveredState struct{}

func (s *DeliveredState) Pay(order *Order, amount float64) error {
    return fmt.Errorf("订单已完成")
}

func (s *DeliveredState) Ship(order *Order, trackingNumber string) error {
    return fmt.Errorf("订单已完成")
}

func (s *DeliveredState) Deliver(order *Order) error {
    return fmt.Errorf("订单已送达")
}

func (s *DeliveredState) Cancel(order *Order) error {
    return fmt.Errorf("订单已完成，无法取消")
}

func (s *DeliveredState) GetName() string {
    return "已送达"
}

// CancelledState 已取消状态
type CancelledState struct{}

func (s *CancelledState) Pay(order *Order, amount float64) error {
    return fmt.Errorf("订单已取消")
}

func (s *CancelledState) Ship(order *Order, trackingNumber string) error {
    return fmt.Errorf("订单已取消")
}

func (s *CancelledState) Deliver(order *Order) error {
    return fmt.Errorf("订单已取消")
}

func (s *CancelledState) Cancel(order *Order) error {
    return fmt.Errorf("订单已取消")
}

func (s *CancelledState) GetName() string {
    return "已取消"
}

// ============ 状态机 ============

type OrderStateMachine struct {
    transitions map[string]map[string]bool
}

func NewOrderStateMachine() *OrderStateMachine {
    return &OrderStateMachine{
        transitions: map[string]map[string]bool{
            "待支付": {"已支付": true, "已取消": true},
            "已支付": {"已发货": true, "已取消": true},
            "已发货": {"已送达": true},
            "已送达": {},
            "已取消": {},
        },
    }
}

func (sm *OrderStateMachine) CanTransition(from, to string) bool {
    if transitions, ok := sm.transitions[from]; ok {
        return transitions[to]
    }
    return false
}

func (sm *OrderStateMachine) GetAllowedTransitions(state string) []string {
    var allowed []string
    if transitions, ok := sm.transitions[state]; ok {
        for to := range transitions {
            allowed = append(allowed, to)
        }
    }
    return allowed
}

// ============ 订单服务 ============

type OrderService struct {
    orders map[string]*Order
    stateMachine *OrderStateMachine
}

func NewOrderService() *OrderService {
    return &OrderService{
        orders:       make(map[string]*Order),
        stateMachine: NewOrderStateMachine(),
    }
}

func (s *OrderService) CreateOrder(id string, amount float64) *Order {
    order := NewOrder(id, amount)
    s.orders[id] = order
    fmt.Printf("创建订单: %s\n", order)
    return order
}

func (s *OrderService) GetOrder(id string) *Order {
    return s.orders[id]
}

func (s *OrderService) GetAllowedActions(orderID string) []string {
    order := s.orders[orderID]
    if order == nil {
        return nil
    }
    return s.stateMachine.GetAllowedTransitions(order.GetStatus())
}

func main() {
    service := NewOrderService()

    fmt.Println("========== 订单状态机演示 ==========")

    // 创建订单
    order := service.CreateOrder("ORD-001", 299.99)
    fmt.Println()

    // 正常流程
    fmt.Println("========== 正常流程 ==========")
    order.Pay(299.99)
    order.Ship("SF123456789")
    order.Deliver()

    fmt.Println("\n========== 创建新订单测试取消 ==========")
    order2 := service.CreateOrder("ORD-002", 199.99)
    order2.Cancel()

    // 尝试对已取消订单操作
    fmt.Println("\n========== 对已取消订单操作 ==========")
    err := order2.Pay(199.99)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    }

    fmt.Println("\n========== 创建新订单测试非法操作 ==========")
    order3 := service.CreateOrder("ORD-003", 399.99)

    // 尝试未支付就发货
    err = order3.Ship("SF987654321")
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    }

    // 支付后尝试再次支付
    order3.Pay(399.99)
    err = order3.Pay(399.99)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    }

    fmt.Println("\n========== 查看允许的操作 ==========")
    order4 := service.CreateOrder("ORD-004", 99.99)
    fmt.Printf("当前状态: %s\n", order4.GetStatus())
    fmt.Printf("允许的操作: %v\n", service.GetAllowedActions("ORD-004"))

    order4.Pay(99.99)
    fmt.Printf("当前状态: %s\n", order4.GetStatus())
    fmt.Printf("允许的操作: %v\n", service.GetAllowedActions("ORD-004"))
}
```

#### 反例说明

**错误1：使用大量条件语句**

```go
// 错误：应该使用状态模式
func (o *Order) Handle(action string) {
    if o.State == "pending" {
        if action == "pay" {
            // ...
        } else if action == "cancel" {
            // ...
        }
    } else if o.State == "paid" {
        // ... 更多条件
    }
    // 难以维护
}
```

**错误2：状态转换不明确**

```go
// 错误：状态转换应该由状态控制
type Order struct {
    State State
}

func (o *Order) Pay() {
    o.State = &PaidState{}  // 外部直接修改状态
}
```

**错误3：状态共享数据**

```go
// 错误：状态不应该共享可变数据
type StateA struct {
    sharedData *Data  // 多个状态共享会导致问题
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 消除复杂条件语句 | 类数量增加 |
| 状态行为集中 | 状态转换复杂 |
| 易于添加新状态 | 可能过度设计 |
| 符合单一职责 | 状态分散 |

**适用场景**:

- 订单状态机
- 工作流引擎
- 游戏角色状态
- 网络连接状态

**Go语言特化**:

- 接口实现状态多态
- 嵌入共享状态数据
- 函数类型作为轻量级状态

---

### 21. 策略模式（Strategy）

#### 概念定义

策略模式定义一系列算法，把它们一个个封装起来，并且使它们可以互相替换。策略模式让算法的变化独立于使用算法的客户。

#### 意图与动机

- **算法封装**: 封装不同的算法
- **运行时切换**: 动态选择算法
- **消除条件语句**: 用多态替代条件判断

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Context                                │
│  - strategy: Strategy                                       │
├─────────────────────────────────────────────────────────────┤
│  + SetStrategy(strategy)                                    │
│  + ExecuteStrategy(data)                                    │
│    strategy.Execute(data)                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │ uses
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     Strategy                                │◄──── 策略接口
│  + Execute(data)                                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  StrategyA    │  │  StrategyB    │  │  StrategyC    │
│  + Execute()  │  │  + Execute()  │  │  + Execute()  │
│    // Alg A   │  │    // Alg B   │  │    // Alg C   │
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

```go
package strategy

// Strategy 策略接口
type Strategy interface {
    Execute(data []int) []int
    GetName() string
}

// Context 上下文
type Context struct {
    strategy Strategy
}

func NewContext(strategy Strategy) *Context {
    return &Context{strategy: strategy}
}

func (c *Context) SetStrategy(strategy Strategy) {
    c.strategy = strategy
}

func (c *Context) ExecuteStrategy(data []int) []int {
    return c.strategy.Execute(data)
}

// BubbleSortStrategy 冒泡排序
type BubbleSortStrategy struct{}

func (b *BubbleSortStrategy) Execute(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)

    n := len(result)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    return result
}

func (b *BubbleSortStrategy) GetName() string {
    return "BubbleSort"
}

// QuickSortStrategy 快速排序
type QuickSortStrategy struct{}

func (q *QuickSortStrategy) Execute(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    quickSort(result, 0, len(result)-1)
    return result
}

func quickSort(arr []int, low, high int) {
    if low < high {
        pi := partition(arr, low, high)
        quickSort(arr, low, pi-1)
        quickSort(arr, pi+1, high)
    }
}

func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    for j := low; j < high; j++ {
        if arr[j] < pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

func (q *QuickSortStrategy) GetName() string {
    return "QuickSort"
}
```

#### 完整示例：支付策略系统

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// ============ 支付策略接口 ============

type PaymentStrategy interface {
    Pay(amount float64) (*PaymentResult, error)
    GetName() string
    GetFeeRate() float64
}

type PaymentResult struct {
    Success       bool
    TransactionID string
    Amount        float64
    Fee           float64
    Message       string
}

func (r *PaymentResult) String() string {
    return fmt.Sprintf("支付%s [ID: %s] 金额: %.2f 手续费: %.2f - %s",
        map[bool]string{true: "成功", false: "失败"}[r.Success],
        r.TransactionID, r.Amount, r.Fee, r.Message)
}

// ============ 具体支付策略 ============

// CreditCardStrategy 信用卡支付
type CreditCardStrategy struct {
    cardNumber string
    cvv        string
    expiryDate string
    name       string
}

func NewCreditCardStrategy(cardNumber, cvv, expiryDate, name string) *CreditCardStrategy {
    return &CreditCardStrategy{
        cardNumber: cardNumber,
        cvv:        cvv,
        expiryDate: expiryDate,
        name:       name,
    }
}

func (c *CreditCardStrategy) Pay(amount float64) (*PaymentResult, error) {
    // 模拟处理
    time.Sleep(100 * time.Millisecond)

    fee := amount * c.GetFeeRate()

    return &PaymentResult{
        Success:       true,
        TransactionID: generateTransactionID("CC"),
        Amount:        amount,
        Fee:           fee,
        Message:       fmt.Sprintf("信用卡 %s 支付", maskCardNumber(c.cardNumber)),
    }, nil
}

func (c *CreditCardStrategy) GetName() string {
    return "信用卡"
}

func (c *CreditCardStrategy) GetFeeRate() float64 {
    return 0.015  // 1.5%
}

func maskCardNumber(cardNumber string) string {
    if len(cardNumber) > 4 {
        return "****" + cardNumber[len(cardNumber)-4:]
    }
    return cardNumber
}

// PayPalStrategy PayPal支付
type PayPalStrategy struct {
    email    string
    password string
}

func NewPayPalStrategy(email, password string) *PayPalStrategy {
    return &PayPalStrategy{
        email:    email,
        password: password,
    }
}

func (p *PayPalStrategy) Pay(amount float64) (*PaymentResult, error) {
    time.Sleep(150 * time.Millisecond)

    fee := amount * p.GetFeeRate()

    return &PaymentResult{
        Success:       true,
        TransactionID: generateTransactionID("PP"),
        Amount:        amount,
        Fee:           fee,
        Message:       fmt.Sprintf("PayPal账户 %s 支付", p.email),
    }, nil
}

func (p *PayPalStrategy) GetName() string {
    return "PayPal"
}

func (p *PayPalStrategy) GetFeeRate() float64 {
    return 0.029 + 0.30/amount  // 2.9% + $0.30
}

// CryptoStrategy 加密货币支付
type CryptoStrategy struct {
    walletAddress string
    cryptoType    string
}

func NewCryptoStrategy(walletAddress, cryptoType string) *CryptoStrategy {
    return &CryptoStrategy{
        walletAddress: walletAddress,
        cryptoType:    cryptoType,
    }
}

func (c *CryptoStrategy) Pay(amount float64) (*PaymentResult, error) {
    time.Sleep(500 * time.Millisecond)  // 区块链确认较慢

    fee := amount * c.GetFeeRate()

    return &PaymentResult{
        Success:       true,
        TransactionID: generateTransactionID("CR"),
        Amount:        amount,
        Fee:           fee,
        Message:       fmt.Sprintf("%s 钱包 %s 支付", c.cryptoType, maskAddress(c.walletAddress)),
    }, nil
}

func (c *CryptoStrategy) GetName() string {
    return "加密货币(" + c.cryptoType + ")"
}

func (c *CryptoStrategy) GetFeeRate() float64 {
    return 0.001  // 0.1%
}

func maskAddress(address string) string {
    if len(address) > 8 {
        return address[:4] + "..." + address[len(address)-4:]
    }
    return address
}

// BankTransferStrategy 银行转账
type BankTransferStrategy struct {
    accountNumber string
    bankCode      string
    accountName   string
}

func NewBankTransferStrategy(accountNumber, bankCode, accountName string) *BankTransferStrategy {
    return &BankTransferStrategy{
        accountNumber: accountNumber,
        bankCode:      bankCode,
        accountName:   accountName,
    }
}

func (b *BankTransferStrategy) Pay(amount float64) (*PaymentResult, error) {
    time.Sleep(200 * time.Millisecond)

    if amount < 10 {
        return &PaymentResult{
            Success: false,
            Message: "银行转账最低金额为10元",
        }, nil
    }

    fee := amount * b.GetFeeRate()

    return &PaymentResult{
        Success:       true,
        TransactionID: generateTransactionID("BT"),
        Amount:        amount,
        Fee:           fee,
        Message:       fmt.Sprintf("银行转账 %s", b.accountName),
    }, nil
}

func (b *BankTransferStrategy) GetName() string {
    return "银行转账"
}

func (b *BankTransferStrategy) GetFeeRate() float64 {
    return 0.005  // 0.5%
}

// ============ 支付上下文 ============

type PaymentContext struct {
    strategy PaymentStrategy
}

func NewPaymentContext(strategy PaymentStrategy) *PaymentContext {
    return &PaymentContext{strategy: strategy}
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentContext) Pay(amount float64) (*PaymentResult, error) {
    fmt.Printf("使用 %s 支付 %.2f\n", p.strategy.GetName(), amount)
    return p.strategy.Pay(amount)
}

func (p *PaymentContext) GetFeeEstimate(amount float64) float64 {
    return amount * p.strategy.GetFeeRate()
}

// ============ 支付服务 ============

type PaymentService struct {
    strategies map[string]PaymentStrategy
}

func NewPaymentService() *PaymentService {
    return &PaymentService{
        strategies: make(map[string]PaymentStrategy),
    }
}

func (s *PaymentService) RegisterStrategy(name string, strategy PaymentStrategy) {
    s.strategies[name] = strategy
}

func (s *PaymentService) GetStrategy(name string) PaymentStrategy {
    return s.strategies[name]
}

func (s *PaymentService) GetAllStrategies() map[string]PaymentStrategy {
    return s.strategies
}

func (s *PaymentService) CompareFees(amount float64) []struct {
    Name string
    Fee  float64
} {
    var comparisons []struct {
        Name string
        Fee  float64
    }

    for name, strategy := range s.strategies {
        comparisons = append(comparisons, struct {
            Name string
            Fee  float64
        }{
            Name: name,
            Fee:  amount * strategy.GetFeeRate(),
        })
    }

    return comparisons
}

// 辅助函数
func generateTransactionID(prefix string) string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%s%d%d", prefix, time.Now().Unix(), rand.Intn(1000))
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // 创建支付服务
    service := NewPaymentService()

    // 注册支付策略
    service.RegisterStrategy("creditcard", NewCreditCardStrategy("1234567890123456", "123", "12/25", "张三"))
    service.RegisterStrategy("paypal", NewPayPalStrategy("zhangsan@example.com", "password"))
    service.RegisterStrategy("bitcoin", NewCryptoStrategy("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "BTC"))
    service.RegisterStrategy("bank", NewBankTransferStrategy("6222021234567890123", "ICBC", "张三"))

    fmt.Println("========== 支付方式手续费比较 ==========")
    amount := 1000.0
    comparisons := service.CompareFees(amount)
    for _, cmp := range comparisons {
        fmt.Printf("%s: %.2f (%.2f%%)\n", cmp.Name, cmp.Fee, cmp.Fee/amount*100)
    }

    fmt.Println("\n========== 使用不同支付方式 ==========")

    // 使用信用卡
    context := NewPaymentContext(service.GetStrategy("creditcard"))
    result, _ := context.Pay(100)
    fmt.Println(result)

    fmt.Println()

    // 切换到PayPal
    context.SetStrategy(service.GetStrategy("paypal"))
    result, _ = context.Pay(100)
    fmt.Println(result)

    fmt.Println()

    // 切换到加密货币
    context.SetStrategy(service.GetStrategy("bitcoin"))
    result, _ = context.Pay(100)
    fmt.Println(result)

    fmt.Println()

    // 切换到银行转账（小额会失败）
    context.SetStrategy(service.GetStrategy("bank"))
    result, _ = context.Pay(5)
    fmt.Println(result)

    fmt.Println()

    result, _ = context.Pay(100)
    fmt.Println(result)

    fmt.Println("\n========== 运行时选择策略 ==========")

    // 根据金额选择最优策略
    largeAmount := 10000.0
    fmt.Printf("大额支付 %.2f，选择最优策略:\n", largeAmount)

    // 加密货币手续费最低
    context.SetStrategy(service.GetStrategy("bitcoin"))
    result, _ = context.Pay(largeAmount)
    fmt.Println(result)
}
```

#### 反例说明

**错误1：使用大量条件语句**

```go
// 错误：应该使用策略模式
func Pay(method string, amount float64) {
    if method == "creditcard" {
        // 信用卡支付逻辑
    } else if method == "paypal" {
        // PayPal支付逻辑
    } else if method == "bitcoin" {
        // 比特币支付逻辑
    }
    // 难以扩展
}
```

**错误2：策略包含上下文数据**

```go
// 错误：策略不应该依赖上下文的具体实现
type Strategy struct {
    context *ConcreteContext  // 应该依赖接口
}
```

**错误3：策略创建开销大**

```go
// 错误：频繁创建策略对象
for _, item := range items {
    strategy := &ConcreteStrategy{}  // 每次循环都创建
    strategy.Execute(item)
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 消除条件语句 | 类数量增加 |
| 运行时切换算法 | 客户端需要知道策略 |
| 易于扩展 | 策略选择开销 |
| 符合开闭原则 | 简单场景过度设计 |

**适用场景**:

- 支付方式选择
- 排序算法选择
- 压缩算法选择
- 路由策略

**Go语言特化**:

- 函数类型作为轻量级策略
- 闭包捕获策略状态
- 接口实现策略多态

---

### 22. 模板方法模式（Template Method）

#### 概念定义

模板方法模式定义一个操作中的算法的骨架，而将一些步骤延迟到子类中。模板方法使得子类可以不改变一个算法的结构即可重定义该算法的某些特定步骤。

#### 意图与动机

- **代码复用**: 复用算法骨架
- **扩展点**: 定义可扩展的算法步骤
- **控制反转**: 父类控制流程，子类实现细节

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                  AbstractClass                              │
├─────────────────────────────────────────────────────────────┤
│  + TemplateMethod()                                         │
│    PrimitiveOperation1()                                    │
│    PrimitiveOperation2()                                    │
│    Hook()                                                   │
├─────────────────────────────────────────────────────────────┤
│  # PrimitiveOperation1()   ◄──── 抽象方法（必须实现）       │
│  # PrimitiveOperation2()                                    │
│  # Hook() {}               ◄──── 钩子方法（可选重写）       │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  ConcreteClass                              │
├─────────────────────────────────────────────────────────────┤
│  # PrimitiveOperation1()   ◄──── 具体实现                   │
│  # PrimitiveOperation2()                                    │
│  # Hook()                  ◄──── 可选重写                   │
└─────────────────────────────────────────────────────────────┘
```

#### Go语言实现

Go语言没有继承，通过嵌入和接口实现模板方法模式。

```go
package template

// Algorithm 算法接口
type Algorithm interface {
    Step1()
    Step2()
    Step3()
}

// Template 模板
type Template struct {
    algorithm Algorithm
}

func NewTemplate(alg Algorithm) *Template {
    return &Template{algorithm: alg}
}

func (t *Template) Execute() {
    fmt.Println("Template: 开始执行")
    t.algorithm.Step1()
    t.algorithm.Step2()
    t.algorithm.Step3()
    fmt.Println("Template: 执行完成")
}

// ConcreteAlgorithm 具体算法
type ConcreteAlgorithm struct{}

func (c *ConcreteAlgorithm) Step1() {
    fmt.Println("Concrete: Step1")
}

func (c *ConcreteAlgorithm) Step2() {
    fmt.Println("Concrete: Step2")
}

func (c *ConcreteAlgorithm) Step3() {
    fmt.Println("Concrete: Step3")
}
```

#### 完整示例：数据导入框架

```go
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "strings"
)

// ============ 数据记录 ============

type DataRecord struct {
    Fields map[string]string
}

func NewDataRecord() *DataRecord {
    return &DataRecord{
        Fields: make(map[string]string),
    }
}

func (r *DataRecord) Set(key, value string) {
    r.Fields[key] = value
}

func (r *DataRecord) Get(key string) string {
    return r.Fields[key]
}

// ============ 数据导入器接口 ============

type DataImporter interface {
    // 必须实现的方法
    ReadSource(source string) ([]DataRecord, error)
    Validate(record *DataRecord) error
    Transform(record *DataRecord) (*DataRecord, error)
    Save(records []DataRecord) error

    // 可选的钩子方法
    BeforeImport()
    AfterImport(records []DataRecord)
    OnError(err error)
    GetName() string
}

// ============ 模板基类（通过嵌入实现） ============

type BaseImporter struct {
    Name string
}

func (b *BaseImporter) BeforeImport() {
    fmt.Printf("[%s] 开始导入...\n", b.Name)
}

func (b *BaseImporter) AfterImport(records []DataRecord) {
    fmt.Printf("[%s] 导入完成，共 %d 条记录\n", b.Name, len(records))
}

func (b *BaseImporter) OnError(err error) {
    fmt.Printf("[%s] 错误: %v\n", b.Name, err)
}

func (b *BaseImporter) GetName() string {
    return b.Name
}

// ImportTemplate 导入模板
type ImportTemplate struct {
    importer DataImporter
}

func NewImportTemplate(importer DataImporter) *ImportTemplate {
    return &ImportTemplate{importer: importer}
}

func (t *ImportTemplate) Import(source string) error {
    // 1. 前置处理
    t.importer.BeforeImport()

    // 2. 读取数据源
    records, err := t.importer.ReadSource(source)
    if err != nil {
        t.importer.OnError(err)
        return err
    }

    // 3. 验证和转换
    var validRecords []DataRecord
    for i, record := range records {
        if err := t.importer.Validate(&record); err != nil {
            fmt.Printf("记录 %d 验证失败: %v\n", i, err)
            continue
        }

        transformed, err := t.importer.Transform(&record)
        if err != nil {
            fmt.Printf("记录 %d 转换失败: %v\n", i, err)
            continue
        }

        validRecords = append(validRecords, *transformed)
    }

    // 4. 保存
    if err := t.importer.Save(validRecords); err != nil {
        t.importer.OnError(err)
        return err
    }

    // 5. 后置处理
    t.importer.AfterImport(validRecords)

    return nil
}

// ============ CSV导入器 ============

type CSVImporter struct {
    BaseImporter
    delimiter rune
}

func NewCSVImporter(delimiter rune) *CSVImporter {
    return &CSVImporter{
        BaseImporter: BaseImporter{Name: "CSV导入器"},
        delimiter:    delimiter,
    }
}

func (c *CSVImporter) ReadSource(source string) ([]DataRecord, error) {
    fmt.Println("[CSV] 读取CSV数据...")

    reader := csv.NewReader(strings.NewReader(source))
    reader.Comma = c.delimiter

    // 读取表头
    headers, err := reader.Read()
    if err != nil {
        return nil, err
    }

    var records []DataRecord
    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        record := NewDataRecord()
        for i, value := range row {
            if i < len(headers) {
                record.Set(headers[i], value)
            }
        }
        records = append(records, *record)
    }

    return records, nil
}

func (c *CSVImporter) Validate(record *DataRecord) error {
    if record.Get("id") == "" {
        return fmt.Errorf("ID不能为空")
    }
    return nil
}

func (c *CSVImporter) Transform(record *DataRecord) (*DataRecord, error) {
    // 转换数据格式
    transformed := NewDataRecord()
    for key, value := range record.Fields {
        // 统一转换为大写
        transformed.Set(strings.ToUpper(key), strings.TrimSpace(value))
    }
    return transformed, nil
}

func (c *CSVImporter) Save(records []DataRecord) error {
    fmt.Printf("[CSV] 保存 %d 条记录到数据库\n", len(records))
    return nil
}

// ============ JSON导入器 ============

type JSONImporter struct {
    BaseImporter
}

func NewJSONImporter() *JSONImporter {
    return &JSONImporter{
        BaseImporter: BaseImporter{Name: "JSON导入器"},
    }
}

func (j *JSONImporter) ReadSource(source string) ([]DataRecord, error) {
    fmt.Println("[JSON] 读取JSON数据...")

    var data []map[string]interface{}
    if err := json.Unmarshal([]byte(source), &data); err != nil {
        return nil, err
    }

    var records []DataRecord
    for _, item := range data {
        record := NewDataRecord()
        for key, value := range item {
            record.Set(key, fmt.Sprintf("%v", value))
        }
        records = append(records, *record)
    }

    return records, nil
}

func (j *JSONImporter) Validate(record *DataRecord) error {
    if record.Get("id") == "" {
        return fmt.Errorf("ID不能为空")
    }
    if record.Get("name") == "" {
        return fmt.Errorf("名称不能为空")
    }
    return nil
}

func (j *JSONImporter) Transform(record *DataRecord) (*DataRecord, error) {
    // JSON特有的转换逻辑
    transformed := NewDataRecord()
    for key, value := range record.Fields {
        // 添加前缀
        transformed.Set("data."+key, value)
    }
    return transformed, nil
}

func (j *JSONImporter) Save(records []DataRecord) error {
    fmt.Printf("[JSON] 保存 %d 条记录到数据库\n", len(records))
    return nil
}

// 重写钩子方法
func (j *JSONImporter) BeforeImport() {
    fmt.Println("[JSON] 准备JSON导入环境...")
}

// ============ XML导入器 ============

type XMLImporter struct {
    BaseImporter
}

func NewXMLImporter() *XMLImporter {
    return &XMLImporter{
        BaseImporter: BaseImporter{Name: "XML导入器"},
    }
}

func (x *XMLImporter) ReadSource(source string) ([]DataRecord, error) {
    fmt.Println("[XML] 读取XML数据...")

    // 简化XML解析
    var records []DataRecord
    // 实际实现会使用xml.Unmarshal
    record := NewDataRecord()
    record.Set("id", "1")
    record.Set("name", "Test")
    records = append(records, *record)

    return records, nil
}

func (x *XMLImporter) Validate(record *DataRecord) error {
    return nil
}

func (x *XMLImporter) Transform(record *DataRecord) (*DataRecord, error) {
    return record, nil
}

func (x *XMLImporter) Save(records []DataRecord) error {
    fmt.Printf("[XML] 保存 %d 条记录到数据库\n", len(records))
    return nil
}

// ============ 报告生成器（另一个模板方法示例） ============

type ReportGenerator interface {
    PrepareData() ([]map[string]interface{}, error)
    FormatData(data []map[string]interface{}) string
    SaveReport(content string, filename string) error
    GetExtension() string
}

type BaseReportGenerator struct {
    Name string
}

func (b *BaseReportGenerator) Generate(filename string, generator ReportGenerator) error {
    // 1. 准备数据
    data, err := generator.PrepareData()
    if err != nil {
        return err
    }

    // 2. 格式化数据
    content := generator.FormatData(data)

    // 3. 保存报告
    fullFilename := filename + generator.GetExtension()
    return generator.SaveReport(content, fullFilename)
}

// PDF报告生成器
type PDFReportGenerator struct {
    BaseReportGenerator
}

func NewPDFReportGenerator() *PDFReportGenerator {
    return &PDFReportGenerator{
        BaseReportGenerator: BaseReportGenerator{Name: "PDF"},
    }
}

func (p *PDFReportGenerator) PrepareData() ([]map[string]interface{}, error) {
    fmt.Println("[PDF] 准备数据...")
    return []map[string]interface{}{
        {"name": "Sales", "value": 1000},
        {"name": "Profit", "value": 200},
    }, nil
}

func (p *PDFReportGenerator) FormatData(data []map[string]interface{}) string {
    fmt.Println("[PDF] 格式化为PDF...")
    return "PDF Content"
}

func (p *PDFReportGenerator) SaveReport(content string, filename string) error {
    fmt.Printf("[PDF] 保存报告: %s\n", filename)
    return nil
}

func (p *PDFReportGenerator) GetExtension() string {
    return ".pdf"
}

// Excel报告生成器
type ExcelReportGenerator struct {
    BaseReportGenerator
}

func NewExcelReportGenerator() *ExcelReportGenerator {
    return &ExcelReportGenerator{
        BaseReportGenerator: BaseReportGenerator{Name: "Excel"},
    }
}

func (e *ExcelReportGenerator) PrepareData() ([]map[string]interface{}, error) {
    fmt.Println("[Excel] 准备数据...")
    return []map[string]interface{}{
        {"name": "Sales", "value": 1000},
        {"name": "Profit", "value": 200},
    }, nil
}

func (e *ExcelReportGenerator) FormatData(data []map[string]interface{}) string {
    fmt.Println("[Excel] 格式化为Excel...")
    return "Excel Content"
}

func (e *ExcelReportGenerator) SaveReport(content string, filename string) error {
    fmt.Printf("[Excel] 保存报告: %s\n", filename)
    return nil
}

func (e *ExcelReportGenerator) GetExtension() string {
    return ".xlsx"
}

func main() {
    fmt.Println("========== CSV导入 ==========")
    csvData := `id,name,age
1,Alice,30
2,Bob,25
3,Charlie,35`

    csvImporter := NewCSVImporter(',')
    csvTemplate := NewImportTemplate(csvImporter)
    csvTemplate.Import(csvData)

    fmt.Println("\n========== JSON导入 ==========")
    jsonData := `[
        {"id": "1", "name": "Alice", "age": 30},
        {"id": "2", "name": "Bob", "age": 25}
    ]`

    jsonImporter := NewJSONImporter()
    jsonTemplate := NewImportTemplate(jsonImporter)
    jsonTemplate.Import(jsonData)

    fmt.Println("\n========== XML导入 ==========")
    xmlImporter := NewXMLImporter()
    xmlTemplate := NewImportTemplate(xmlImporter)
    xmlTemplate.Import("<data><item id='1' name='Test'/></data>")

    fmt.Println("\n========== 报告生成 ==========")

    pdfGenerator := NewPDFReportGenerator()
    pdfGenerator.Generate("sales_report", pdfGenerator)

    fmt.Println()

    excelGenerator := NewExcelReportGenerator()
    excelGenerator.Generate("sales_report", excelGenerator)
}
```

#### 反例说明

**错误1：模板方法过于复杂**

```go
// 错误：模板方法应该只定义骨架
func TemplateMethod() {
    // 太多步骤，应该拆分
    Step1()
    Step2()
    Step3()
    Step4()
    Step5()
    Step6()
    // ...
}
```

**错误2：子类破坏算法结构**

```go
// 错误：子类不应该覆盖模板方法
func (c *Concrete) TemplateMethod() {
    // 完全重写，破坏了算法结构
}
```

**错误3：过多的抽象方法**

```go
// 错误：过多的抽象方法增加子类负担
type Template interface {
    Step1()  // 必须实现
    Step2()  // 必须实现
    Step3()  // 必须实现
    Step4()  // 必须实现
    Step5()  // 必须实现
    // 太多！
}
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 代码复用 | 类数量增加 |
| 控制算法结构 | 继承限制灵活性 |
| 易于扩展 | 骨架变更影响所有子类 |
| 符合开闭原则 | 需要理解父类结构 |

**适用场景**:

- 数据导入/导出
- 报告生成
- 测试框架
- 构建流程

**Go语言特化**:

- 嵌入实现代码复用
- 接口定义扩展点
- 函数类型作为钩子

---

### 23. 访问者模式（Visitor）

#### 概念定义

访问者模式表示一个作用于某对象结构中的各元素的操作。它使你可以在不改变各元素的类的前提下定义作用于这些元素的新操作。

#### 意图与动机

- **添加操作**: 不修改类而添加新操作
- **分离关注点**: 将算法与对象结构分离
- **复杂遍历**: 遍历复杂结构执行操作

#### UML结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Visitor                                │◄──── 访问者接口
│  + VisitConcreteElementA(element)                           │
│  + VisitConcreteElementB(element)                           │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  ConcreteVisitor                            │
│  + VisitConcreteElementA(element)                           │
│    // 对A的操作                                             │
│  + VisitConcreteElementB(element)                           │
│    // 对B的操作                                             │
└──────────────────────────┬──────────────────────────────────┘
                           │ visits
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Element                                │◄──── 元素接口
│  + Accept(visitor)                                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  ElementA     │  │  ElementB     │  │  ElementC     │
│  + Accept(v)  │  │  + Accept(v)  │  │  + Accept(v)  │
│    v.VisitA() │  │    v.VisitB() │  │    v.VisitC() │
└───────────────┘  └───────────────┘  └───────────────┘
```

#### Go语言实现

```go
package visitor

// Visitor 访问者接口
type Visitor interface {
    VisitConcreteElementA(element *ConcreteElementA)
    VisitConcreteElementB(element *ConcreteElementB)
}

// Element 元素接口
type Element interface {
    Accept(visitor Visitor)
}

// ConcreteElementA 具体元素A
type ConcreteElementA struct {
    Name string
}

func (e *ConcreteElementA) Accept(visitor Visitor) {
    visitor.VisitConcreteElementA(e)
}

// ConcreteElementB 具体元素B
type ConcreteElementB struct {
    Value int
}

func (e *ConcreteElementB) Accept(visitor Visitor) {
    visitor.VisitConcreteElementB(e)
}

// ConcreteVisitor 具体访问者
type ConcreteVisitor struct{}

func (v *ConcreteVisitor) VisitConcreteElementA(element *ConcreteElementA) {
    fmt.Printf("Visiting ElementA: %s\n", element.Name)
}

func (v *ConcreteVisitor) VisitConcreteElementB(element *ConcreteElementB) {
    fmt.Printf("Visiting ElementB: %d\n", element.Value)
}

// ObjectStructure 对象结构
type ObjectStructure struct {
    elements []Element
}

func (o *ObjectStructure) Attach(element Element) {
    o.elements = append(o.elements, element)
}

func (o *ObjectStructure) Accept(visitor Visitor) {
    for _, element := range o.elements {
        element.Accept(visitor)
    }
}
```

#### 完整示例：文档导出系统

```go
package main

import (
    "fmt"
    "strings"
)

// ============ 文档元素接口 ============

type DocumentElement interface {
    Accept(visitor DocumentVisitor)
    GetText() string
}

// ============ 具体文档元素 ============

// Paragraph 段落
type Paragraph struct {
    Text string
}

func NewParagraph(text string) *Paragraph {
    return &Paragraph{Text: text}
}

func (p *Paragraph) Accept(visitor DocumentVisitor) {
    visitor.VisitParagraph(p)
}

func (p *Paragraph) GetText() string {
    return p.Text
}

// Heading 标题
type Heading struct {
    Level int
    Text  string
}

func NewHeading(level int, text string) *Heading {
    return &Heading{Level: level, Text: text}
}

func (h *Heading) Accept(visitor DocumentVisitor) {
    visitor.VisitHeading(h)
}

func (h *Heading) GetText() string {
    return h.Text
}

// Link 链接
type Link struct {
    Text string
    URL  string
}

func NewLink(text, url string) *Link {
    return &Link{Text: text, URL: url}
}

func (l *Link) Accept(visitor DocumentVisitor) {
    visitor.VisitLink(l)
}

func (l *Link) GetText() string {
    return l.Text
}

// Image 图片
type Image struct {
    AltText string
    Source  string
}

func NewImage(altText, source string) *Image {
    return &Image{AltText: altText, Source: source}
}

func (i *Image) Accept(visitor DocumentVisitor) {
    visitor.VisitImage(i)
}

func (i *Image) GetText() string {
    return i.AltText
}

// List 列表
type List struct {
    Items []string
}

func NewList(items []string) *List {
    return &List{Items: items}
}

func (l *List) Accept(visitor DocumentVisitor) {
    visitor.VisitList(l)
}

func (l *List) GetText() string {
    return strings.Join(l.Items, ", ")
}

// ============ 访问者接口 ============

type DocumentVisitor interface {
    VisitParagraph(p *Paragraph)
    VisitHeading(h *Heading)
    VisitLink(l *Link)
    VisitImage(i *Image)
    VisitList(l *List)
    GetResult() string
}

// ============ HTML导出访问者 ============

type HTMLExportVisitor struct {
    builder strings.Builder
}

func NewHTMLExportVisitor() *HTMLExportVisitor {
    return &HTMLExportVisitor{}
}

func (v *HTMLExportVisitor) VisitParagraph(p *Paragraph) {
    v.builder.WriteString(fmt.Sprintf("<p>%s</p>\n", p.Text))
}

func (v *HTMLExportVisitor) VisitHeading(h *Heading) {
    tag := fmt.Sprintf("h%d", h.Level)
    v.builder.WriteString(fmt.Sprintf("<%s>%s</%s>\n", tag, h.Text, tag))
}

func (v *HTMLExportVisitor) VisitLink(l *Link) {
    v.builder.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, l.URL, l.Text))
}

func (v *HTMLExportVisitor) VisitImage(i *Image) {
    v.builder.WriteString(fmt.Sprintf(`<img src="%s" alt="%s" />`, i.Source, i.AltText))
}

func (v *HTMLExportVisitor) VisitList(l *List) {
    v.builder.WriteString("<ul>\n")
    for _, item := range l.Items {
        v.builder.WriteString(fmt.Sprintf("  <li>%s</li>\n", item))
    }
    v.builder.WriteString("</ul>\n")
}

func (v *HTMLExportVisitor) GetResult() string {
    return v.builder.String()
}

// ============ Markdown导出访问者 ============

type MarkdownExportVisitor struct {
    builder strings.Builder
}

func NewMarkdownExportVisitor() *MarkdownExportVisitor {
    return &MarkdownExportVisitor{}
}

func (v *MarkdownExportVisitor) VisitParagraph(p *Paragraph) {
    v.builder.WriteString(p.Text + "\n\n")
}

func (v *MarkdownExportVisitor) VisitHeading(h *Heading) {
    prefix := strings.Repeat("#", h.Level)
    v.builder.WriteString(fmt.Sprintf("%s %s\n\n", prefix, h.Text))
}

func (v *MarkdownExportVisitor) VisitLink(l *Link) {
    v.builder.WriteString(fmt.Sprintf("[%s](%s)", l.Text, l.URL))
}

func (v *MarkdownExportVisitor) VisitImage(i *Image) {
    v.builder.WriteString(fmt.Sprintf("![%s](%s)", i.AltText, i.Source))
}

func (v *MarkdownExportVisitor) VisitList(l *List) {
    for _, item := range l.Items {
        v.builder.WriteString(fmt.Sprintf("- %s\n", item))
    }
    v.builder.WriteString("\n")
}

func (v *MarkdownExportVisitor) GetResult() string {
    return v.builder.String()
}

// ============ 纯文本导出访问者 ============

type PlainTextExportVisitor struct {
    builder strings.Builder
}

func NewPlainTextExportVisitor() *PlainTextExportVisitor {
    return &PlainTextExportVisitor{}
}

func (v *PlainTextExportVisitor) VisitParagraph(p *Paragraph) {
    v.builder.WriteString(p.Text + "\n\n")
}

func (v *PlainTextExportVisitor) VisitHeading(h *Heading) {
    v.builder.WriteString(h.Text + "\n")
    v.builder.WriteString(strings.Repeat("=", len(h.Text)) + "\n\n")
}

func (v *PlainTextExportVisitor) VisitLink(l *Link) {
    v.builder.WriteString(fmt.Sprintf("%s (%s)", l.Text, l.URL))
}

func (v *PlainTextExportVisitor) VisitImage(i *Image) {
    v.builder.WriteString(fmt.Sprintf("[图片: %s]", i.AltText))
}

func (v *PlainTextExportVisitor) VisitList(l *List) {
    for i, item := range l.Items {
        v.builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
    }
    v.builder.WriteString("\n")
}

func (v *PlainTextExportVisitor) GetResult() string {
    return v.builder.String()
}

// ============ 字数统计访问者 ============

type WordCountVisitor struct {
    totalWords int
    counts     map[string]int
}

func NewWordCountVisitor() *WordCountVisitor {
    return &WordCountVisitor{
        counts: make(map[string]int),
    }
}

func (v *WordCountVisitor) VisitParagraph(p *Paragraph) {
    words := len(strings.Fields(p.Text))
    v.totalWords += words
    v.counts["paragraph"] += words
}

func (v *WordCountVisitor) VisitHeading(h *Heading) {
    words := len(strings.Fields(h.Text))
    v.totalWords += words
    v.counts["heading"] += words
}

func (v *WordCountVisitor) VisitLink(l *Link) {
    words := len(strings.Fields(l.Text))
    v.totalWords += words
    v.counts["link"] += words
}

func (v *WordCountVisitor) VisitImage(i *Image) {
    words := len(strings.Fields(i.AltText))
    v.totalWords += words
    v.counts["image"] += words
}

func (v *WordCountVisitor) VisitList(l *List) {
    for _, item := range l.Items {
        words := len(strings.Fields(item))
        v.totalWords += words
        v.counts["list"] += words
    }
}

func (v *WordCountVisitor) GetResult() string {
    result := fmt.Sprintf("总字数: %d\n", v.totalWords)
    for elementType, count := range v.counts {
        result += fmt.Sprintf("  %s: %d\n", elementType, count)
    }
    return result
}

// ============ 文档结构 ============

type Document struct {
    Title    string
    Elements []DocumentElement
}

func NewDocument(title string) *Document {
    return &Document{
        Title:    title,
        Elements: make([]DocumentElement, 0),
    }
}

func (d *Document) Add(element DocumentElement) {
    d.Elements = append(d.Elements, element)
}

func (d *Document) Accept(visitor DocumentVisitor) {
    for _, element := range d.Elements {
        element.Accept(visitor)
    }
}

// ============ 导出管理器 ============

type ExportManager struct{}

func NewExportManager() *ExportManager {
    return &ExportManager{}
}

func (m *ExportManager) Export(document *Document, format string) string {
    var visitor DocumentVisitor

    switch format {
    case "html":
        visitor = NewHTMLExportVisitor()
    case "markdown":
        visitor = NewMarkdownExportVisitor()
    case "text":
        visitor = NewPlainTextExportVisitor()
    case "wordcount":
        visitor = NewWordCountVisitor()
    default:
        return "不支持的格式"
    }

    document.Accept(visitor)
    return visitor.GetResult()
}

func main() {
    // 创建文档
    doc := NewDocument("示例文档")
    doc.Add(NewHeading(1, "欢迎来到Go设计模式"))
    doc.Add(NewParagraph("这是访问者模式的示例文档。访问者模式允许我们在不改变元素类的情况下定义新的操作。"))
    doc.Add(NewHeading(2, "主要特点"))
    doc.Add(NewList([]string{
        "分离算法与对象结构",
        "易于添加新操作",
        "集中相关操作",
    }))
    doc.Add(NewParagraph("了解更多信息，请访问"))
    doc.Add(NewLink("Go官网", "https://golang.org"))
    doc.Add(NewImage("Go Logo", "https://golang.org/logo.png"))

    manager := NewExportManager()

    fmt.Println("========== HTML导出 ==========")
    fmt.Println(manager.Export(doc, "html"))

    fmt.Println("========== Markdown导出 ==========")
    fmt.Println(manager.Export(doc, "markdown"))

    fmt.Println("========== 纯文本导出 ==========")
    fmt.Println(manager.Export(doc, "text"))

    fmt.Println("========== 字数统计 ==========")
    fmt.Println(manager.Export(doc, "wordcount"))
}
```

#### 反例说明

**错误1：元素类型频繁变化**

```go
// 错误：如果经常添加新元素类型，访问者模式不合适
// 每次添加新元素都需要修改所有访问者
type NewElement struct{}

// 需要修改所有Visitor接口和实现
```

**错误2：访问者依赖具体元素**

```go
// 错误：访问者不应该依赖元素的内部实现
func (v *Visitor) VisitElement(e *Element) {
    // 访问元素的私有字段
    _ = e.privateField  // 破坏了封装
}
```

**错误3：过度使用访问者**

```go
// 错误：简单场景不需要访问者
// 直接使用方法即可
```

#### 优缺点分析

| 优点 | 缺点 |
|------|------|
| 易于添加新操作 | 难以添加新元素类型 |
| 集中相关操作 | 破坏封装 |
| 遍历灵活 | 代码复杂度高 |
| 符合单一职责 | 学习曲线陡峭 |

**适用场景**:

- 文档导出
- AST遍历
- 编译器
- 复杂对象结构操作

**Go语言特化**:

- 接口实现双重分发
- 类型断言处理未知类型
- 函数类型作为轻量级访问者

---

## 总结

### 设计模式对比表

| 模式 | 类型 | 主要目的 | Go特性应用 |
|------|------|----------|------------|
| 单例 | 创建型 | 唯一实例 | sync.Once |
| 工厂方法 | 创建型 | 解耦创建 | 接口隐式实现 |
| 抽象工厂 | 创建型 | 产品族 | 接口组合 |
| 建造者 | 创建型 | 复杂对象构建 | 函数式选项 |
| 原型 | 创建型 | 对象复制 | JSON序列化 |
| 适配器 | 结构型 | 接口兼容 | 嵌入 |
| 桥接 | 结构型 | 维度分离 | 接口组合 |
| 组合 | 结构型 | 树形结构 | 接口递归 |
| 装饰器 | 结构型 | 动态扩展 | 函数闭包 |
| 外观 | 结构型 | 简化接口 | 包级别函数 |
| 享元 | 结构型 | 共享对象 | map工厂 |
| 代理 | 结构型 | 访问控制 | Channel |
| 责任链 | 行为型 | 请求传递 | 接口链 |
| 命令 | 行为型 | 请求封装 | 函数类型 |
| 解释器 | 行为型 | 语言解析 | 递归下降 |
| 迭代器 | 行为型 | 统一遍历 | range关键字 |
| 中介者 | 行为型 | 解耦通信 | Channel |
| 备忘录 | 行为型 | 状态保存 | 结构体复制 |
| 观察者 | 行为型 | 事件通知 | Channel |
| 状态 | 行为型 | 状态机 | 接口多态 |
| 策略 | 行为型 | 算法替换 | 函数类型 |
| 模板方法 | 行为型 | 算法骨架 | 嵌入 |
| 访问者 | 行为型 | 操作扩展 | 双重分发 |

### Go语言设计模式最佳实践

1. **优先使用组合而非继承**: Go没有继承，通过嵌入实现代码复用
2. **利用接口隐式实现**: 减少样板代码
3. **函数类型作为轻量级模式**: 策略、命令等模式可用函数实现
4. **Channel用于通信**: 观察者、中介者等模式可用Channel实现
5. **Context传递状态**: 跨层传递请求上下文
6. **sync包处理并发**: Once、Mutex、RWMutex等

### 参考资源

- 《设计模式：可复用面向对象软件的基础》(GoF)
- 《Go语言设计模式》
- Go官方文档: <https://golang.org/doc/>

---

*文档版本: 1.0*
*最后更新: 2024年*
