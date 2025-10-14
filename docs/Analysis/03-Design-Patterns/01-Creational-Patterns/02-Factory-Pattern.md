# 3.1.1 工厂模式 (Factory Pattern)

## 3.1.1.1 目录

## 3.1.1.2 1. 概述

### 3.1.1.2.1 定义

工厂模式定义一个创建对象的接口，让子类决定实例化哪一个类。

**形式化定义**:
$$Factory = (Creator, Product, ConcreteCreator, ConcreteProduct, FactoryMethod)$$

其中：

- $Creator$ 是创建者接口
- $Product$ 是产品接口
- $ConcreteCreator$ 是具体创建者
- $ConcreteProduct$ 是具体产品
- $FactoryMethod$ 是工厂方法

### 3.1.1.2.2 核心特征

- **封装创建**: 封装对象创建逻辑
- **多态性**: 支持多种产品类型
- **扩展性**: 易于添加新产品
- **解耦**: 客户端与具体产品解耦

## 3.1.1.3 2. 理论基础

### 3.1.1.3.1 数学形式化

**定义 2.1** (工厂模式): 工厂模式是一个五元组 $F = (C, P, M, R, V)$

其中：

- $C$ 是创建者集合
- $P$ 是产品集合
- $M$ 是工厂方法映射，$M: C \rightarrow P$
- $R$ 是关系约束
- $V$ 是验证规则

**定理 2.1** (创建映射): 对于任意创建者 $c \in C$，存在唯一产品 $p \in P$ 使得 $M(c) = p$

**证明**: 由工厂方法的确定性保证。

### 3.1.1.3.2 范畴论视角

在范畴论中，工厂模式可以表示为：

$$Factory : Creator \rightarrow Product$$

其中 $Creator$ 和 $Product$ 是对象范畴。

## 3.1.1.4 3. Go语言实现

### 3.1.1.4.1 简单工厂模式

```go
package factory

import "fmt"

// Product 产品接口
type Product interface {
    Use() string
    GetName() string
}

// ConcreteProductA 具体产品A
type ConcreteProductA struct {
    name string
}

func (p *ConcreteProductA) Use() string {
    return "Using Product A"
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

// ConcreteProductB 具体产品B
type ConcreteProductB struct {
    name string
}

func (p *ConcreteProductB) Use() string {
    return "Using Product B"
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

// SimpleFactory 简单工厂
type SimpleFactory struct{}

// CreateProduct 创建产品
func (f *SimpleFactory) CreateProduct(productType string) (Product, error) {
    switch productType {
    case "A":
        return &ConcreteProductA{name: "ProductA"}, nil
    case "B":
        return &ConcreteProductB{name: "ProductB"}, nil
    default:
        return nil, fmt.Errorf("unknown product type: %s", productType)
    }
}

// 工厂函数
func NewProduct(productType string) (Product, error) {
    factory := &SimpleFactory{}
    return factory.CreateProduct(productType)
}

```

### 3.1.1.4.2 工厂方法模式

```go
package factorymethod

import "fmt"

// Product 产品接口
type Product interface {
    Use() string
    GetName() string
}

// Creator 创建者接口
type Creator interface {
    CreateProduct() Product
    GetCreatorName() string
}

// ConcreteCreatorA 具体创建者A
type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) CreateProduct() Product {
    return &ConcreteProductA{name: "ProductA"}
}

func (c *ConcreteCreatorA) GetCreatorName() string {
    return "CreatorA"
}

// ConcreteCreatorB 具体创建者B
type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) CreateProduct() Product {
    return &ConcreteProductB{name: "ProductB"}
}

func (c *ConcreteCreatorB) GetCreatorName() string {
    return "CreatorB"
}

// ConcreteProductA 具体产品A
type ConcreteProductA struct {
    name string
}

func (p *ConcreteProductA) Use() string {
    return "Using Product A"
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

// ConcreteProductB 具体产品B
type ConcreteProductB struct {
    name string
}

func (p *ConcreteProductB) Use() string {
    return "Using Product B"
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

```

### 3.1.1.4.3 抽象工厂模式

```go
package abstractfactory

import "fmt"

// AbstractProductA 抽象产品A
type AbstractProductA interface {
    UseA() string
}

// AbstractProductB 抽象产品B
type AbstractProductB interface {
    UseB() string
}

// AbstractFactory 抽象工厂
type AbstractFactory interface {
    CreateProductA() AbstractProductA
    CreateProductB() AbstractProductB
}

// ConcreteProductA1 具体产品A1
type ConcreteProductA1 struct{}

func (p *ConcreteProductA1) UseA() string {
    return "Using Product A1"
}

// ConcreteProductA2 具体产品A2
type ConcreteProductA2 struct{}

func (p *ConcreteProductA2) UseA() string {
    return "Using Product A2"
}

// ConcreteProductB1 具体产品B1
type ConcreteProductB1 struct{}

func (p *ConcreteProductB1) UseB() string {
    return "Using Product B1"
}

// ConcreteProductB2 具体产品B2
type ConcreteProductB2 struct{}

func (p *ConcreteProductB2) UseB() string {
    return "Using Product B2"
}

// ConcreteFactory1 具体工厂1
type ConcreteFactory1 struct{}

func (f *ConcreteFactory1) CreateProductA() AbstractProductA {
    return &ConcreteProductA1{}
}

func (f *ConcreteFactory1) CreateProductB() AbstractProductB {
    return &ConcreteProductB1{}
}

// ConcreteFactory2 具体工厂2
type ConcreteFactory2 struct{}

func (f *ConcreteFactory2) CreateProductA() AbstractProductA {
    return &ConcreteProductA2{}
}

func (f *ConcreteFactory2) CreateProductB() AbstractProductB {
    return &ConcreteProductB2{}
}

// NewFactory 创建工厂
func NewFactory(factoryType string) (AbstractFactory, error) {
    switch factoryType {
    case "1":
        return &ConcreteFactory1{}, nil
    case "2":
        return &ConcreteFactory2{}, nil
    default:
        return nil, fmt.Errorf("unknown factory type: %s", factoryType)
    }
}

```

## 3.1.1.5 4. 工程案例

### 3.1.1.5.1 数据库连接工厂

```go
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
)

// Database 数据库接口
type Database interface {
    Connect() error
    Query(sql string) (*sql.Rows, error)
    Close() error
    GetName() string
}

// PostgreSQL 数据库
type PostgreSQL struct {
    db   *sql.DB
    dsn  string
    name string
}

func (p *PostgreSQL) Connect() error {
    db, err := sql.Open("postgres", p.dsn)
    if err != nil {
        return err
    }
    p.db = db
    return nil
}

func (p *PostgreSQL) Query(sql string) (*sql.Rows, error) {
    return p.db.Query(sql)
}

func (p *PostgreSQL) Close() error {
    return p.db.Close()
}

func (p *PostgreSQL) GetName() string {
    return p.name
}

// MySQL 数据库
type MySQL struct {
    db   *sql.DB
    dsn  string
    name string
}

func (m *MySQL) Connect() error {
    db, err := sql.Open("mysql", m.dsn)
    if err != nil {
        return err
    }
    m.db = db
    return nil
}

func (m *MySQL) Query(sql string) (*sql.Rows, error) {
    return m.db.Query(sql)
}

func (m *MySQL) Close() error {
    return m.db.Close()
}

func (m *MySQL) GetName() string {
    return m.name
}

// DatabaseFactory 数据库工厂
type DatabaseFactory struct{}

// CreateDatabase 创建数据库连接
func (f *DatabaseFactory) CreateDatabase(dbType, dsn string) (Database, error) {
    switch dbType {
    case "postgres":
        return &PostgreSQL{dsn: dsn, name: "PostgreSQL"}, nil
    case "mysql":
        return &MySQL{dsn: dsn, name: "MySQL"}, nil
    default:
        return nil, fmt.Errorf("unsupported database type: %s", dbType)
    }
}

```

### 3.1.1.5.2 日志工厂

```go
package logger

import (
    "fmt"
    "log"
    "os"
)

// Logger 日志接口
type Logger interface {
    Info(format string, v ...interface{})
    Error(format string, v ...interface{})
    Debug(format string, v ...interface{})
    GetType() string
}

// ConsoleLogger 控制台日志
type ConsoleLogger struct {
    logger *log.Logger
}

func (c *ConsoleLogger) Info(format string, v ...interface{}) {
    c.logger.Printf("[INFO] "+format, v...)
}

func (c *ConsoleLogger) Error(format string, v ...interface{}) {
    c.logger.Printf("[ERROR] "+format, v...)
}

func (c *ConsoleLogger) Debug(format string, v ...interface{}) {
    c.logger.Printf("[DEBUG] "+format, v...)
}

func (c *ConsoleLogger) GetType() string {
    return "Console"
}

// FileLogger 文件日志
type FileLogger struct {
    logger *log.Logger
    file   *os.File
}

func (f *FileLogger) Info(format string, v ...interface{}) {
    f.logger.Printf("[INFO] "+format, v...)
}

func (f *FileLogger) Error(format string, v ...interface{}) {
    f.logger.Printf("[ERROR] "+format, v...)
}

func (f *FileLogger) Debug(format string, v ...interface{}) {
    f.logger.Printf("[DEBUG] "+format, v...)
}

func (f *FileLogger) GetType() string {
    return "File"
}

// LoggerFactory 日志工厂
type LoggerFactory struct{}

// CreateLogger 创建日志器
func (f *LoggerFactory) CreateLogger(loggerType, filename string) (Logger, error) {
    switch loggerType {
    case "console":
        return &ConsoleLogger{
            logger: log.New(os.Stdout, "", log.LstdFlags),
        }, nil
    case "file":
        file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            return nil, fmt.Errorf("failed to open log file: %w", err)
        }
        return &FileLogger{
            logger: log.New(file, "", log.LstdFlags),
            file:   file,
        }, nil
    default:
        return nil, fmt.Errorf("unsupported logger type: %s", loggerType)
    }
}

```

## 3.1.1.6 5. 批判性分析

### 3.1.1.6.1 优势

1. **封装创建逻辑**: 隐藏对象创建复杂性
2. **支持扩展**: 易于添加新产品类型
3. **解耦**: 客户端与具体产品解耦
4. **多态性**: 支持运行时多态

### 3.1.1.6.2 劣势

1. **类爆炸**: 可能导致类数量增加
2. **复杂性**: 增加系统复杂性
3. **性能开销**: 额外的抽象层开销
4. **过度设计**: 简单场景可能过度设计

### 3.1.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 函数 | 高 | 低 |
| Java | 抽象类 + 接口 | 中 | 中 |
| C++ | 虚函数 | 中 | 中 |
| Python | 类方法 | 中 | 低 |

### 3.1.1.6.4 最新趋势

1. **函数式工厂**: 使用函数而非类
2. **依赖注入**: 结合DI容器
3. **泛型支持**: 利用Go 1.18+泛型
4. **配置驱动**: 基于配置的工厂

## 3.1.1.7 6. 面试题与考点

### 3.1.1.7.1 基础考点

1. **Q**: 工厂模式的三种类型是什么？
   **A**: 简单工厂、工厂方法、抽象工厂

2. **Q**: 工厂模式与new操作符的区别？
   **A**: 工厂模式封装创建逻辑，支持多态和扩展

3. **Q**: 什么时候使用工厂模式？
   **A**: 对象创建复杂、需要扩展、需要解耦时

### 3.1.1.7.2 进阶考点

1. **Q**: 工厂模式在微服务中的应用？
   **A**: 服务发现、客户端工厂、配置工厂

2. **Q**: 如何避免工厂模式的类爆炸？
   **A**: 使用参数化工厂、配置驱动、泛型

3. **Q**: 工厂模式与依赖注入的关系？
   **A**: DI容器本质上是高级工厂模式

## 3.1.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 工厂模式 | 封装对象创建逻辑的设计模式 | Factory Pattern |
| 简单工厂 | 使用静态方法创建对象的工厂 | Simple Factory |
| 工厂方法 | 使用虚函数创建对象的工厂 | Factory Method |
| 抽象工厂 | 创建相关对象族的工厂 | Abstract Factory |
| 产品族 | 相关产品的集合 | Product Family |

## 3.1.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 类爆炸 | 产品类过多 | 使用参数化工厂 |
| 过度设计 | 简单场景使用复杂工厂 | 评估实际需求 |
| 性能问题 | 抽象层过多 | 使用函数式工厂 |
| 维护困难 | 工厂逻辑复杂 | 简化工厂逻辑 |

## 3.1.1.10 9. 相关主题

- [单例模式](./01-Singleton-Pattern.md)
- [建造者模式](./03-Builder-Pattern.md)
- [原型模式](./04-Prototype-Pattern.md)
- [抽象工厂模式](./05-Abstract-Factory-Pattern.md)
- [依赖注入模式](../04-Concurrent-Patterns/01-Dependency-Injection.md)

## 3.1.1.11 10. 学习路径

### 3.1.1.11.1 新手路径

1. 理解工厂模式的基本概念
2. 学习简单工厂模式
3. 实现基础的工厂方法
4. 理解多态性的重要性

### 3.1.1.11.2 进阶路径

1. 学习工厂方法模式
2. 理解抽象工厂模式
3. 掌握工厂模式的应用场景
4. 学习工厂模式的最佳实践

### 3.1.1.11.3 高阶路径

1. 分析工厂模式在大型项目中的应用
2. 理解工厂模式与架构设计的关系
3. 掌握工厂模式的性能优化
4. 学习工厂模式的替代方案

---

**相关文档**: [创建型模式总览](./README.md) | [设计模式总览](../README.md)
