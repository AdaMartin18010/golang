# 3.1.1 原型模式 (Prototype Pattern)

<!-- TOC START -->
- [3.1.1 原型模式 (Prototype Pattern)](#311-原型模式-prototype-pattern)
  - [3.1.1.1 目录](#3111-目录)
  - [3.1.1.2 1. 概述](#3112-1-概述)
    - [3.1.1.2.1 定义](#31121-定义)
    - [3.1.1.2.2 核心特征](#31122-核心特征)
  - [3.1.1.3 2. 理论基础](#3113-2-理论基础)
    - [3.1.1.3.1 数学形式化](#31131-数学形式化)
    - [3.1.1.3.2 范畴论视角](#31132-范畴论视角)
  - [3.1.1.4 3. Go语言实现](#3114-3-go语言实现)
    - [3.1.1.4.1 基础原型模式](#31141-基础原型模式)
    - [3.1.1.4.2 原型注册表模式](#31142-原型注册表模式)
    - [3.1.1.4.3 函数式原型模式](#31143-函数式原型模式)
  - [3.1.1.5 4. 工程案例](#3115-4-工程案例)
    - [3.1.1.5.1 配置对象原型](#31151-配置对象原型)
    - [3.1.1.5.2 文档模板原型](#31152-文档模板原型)
  - [3.1.1.6 5. 批判性分析](#3116-5-批判性分析)
    - [3.1.1.6.1 优势](#31161-优势)
    - [3.1.1.6.2 劣势](#31162-劣势)
    - [3.1.1.6.3 行业对比](#31163-行业对比)
    - [3.1.1.6.4 最新趋势](#31164-最新趋势)
  - [3.1.1.7 6. 面试题与考点](#3117-6-面试题与考点)
    - [3.1.1.7.1 基础考点](#31171-基础考点)
    - [3.1.1.7.2 进阶考点](#31172-进阶考点)
  - [3.1.1.8 7. 术语表](#3118-7-术语表)
  - [3.1.1.9 8. 常见陷阱](#3119-8-常见陷阱)
  - [3.1.1.10 9. 相关主题](#31110-9-相关主题)
  - [3.1.1.11 10. 学习路径](#31111-10-学习路径)
    - [3.1.1.11.1 新手路径](#311111-新手路径)
    - [3.1.1.11.2 进阶路径](#311112-进阶路径)
    - [3.1.1.11.3 高阶路径](#311113-高阶路径)
<!-- TOC END -->

## 3.1.1.1 目录

## 3.1.1.2 1. 概述

### 3.1.1.2.1 定义

原型模式用原型实例指定创建对象的种类，并且通过复制这些原型创建新的对象。

**形式化定义**:
$$Prototype = (Prototype, Clone, ConcretePrototype, Registry, CloneMethod)$$

其中：

- $Prototype$ 是原型接口
- $Clone$ 是克隆方法
- $ConcretePrototype$ 是具体原型
- $Registry$ 是原型注册表
- $CloneMethod$ 是克隆方法集合

### 3.1.1.2.2 核心特征

- **对象复制**: 通过复制现有对象创建新对象
- **性能优化**: 避免昂贵的初始化过程
- **动态创建**: 运行时动态创建对象
- **配置复用**: 复用复杂对象的配置

## 3.1.1.3 2. 理论基础

### 3.1.1.3.1 数学形式化

**定义 2.1** (原型模式): 原型模式是一个五元组 $P = (P, C, R, M, V)$

其中：

- $P$ 是原型集合
- $C$ 是克隆方法集合
- $R$ 是注册表
- $M$ 是克隆映射，$M: P \rightarrow P$
- $V$ 是验证规则

**定理 2.1** (克隆一致性): 对于任意原型 $p \in P$，$M(p)$ 与 $p$ 在结构上等价

**证明**: 由克隆方法的实现保证。

### 3.1.1.3.2 范畴论视角

在范畴论中，原型模式可以表示为：

$$Clone : Prototype \rightarrow Prototype$$

其中 $Prototype$ 是对象范畴。

## 3.1.1.4 3. Go语言实现

### 3.1.1.4.1 基础原型模式

```go
package prototype

import (
    "encoding/json"
    "fmt"
)

// Prototype 原型接口
type Prototype interface {
    Clone() Prototype
    GetName() string
    SetName(name string)
}

// ConcretePrototype 具体原型
type ConcretePrototype struct {
    Name   string
    Data   map[string]interface{}
    Config []string
}

func (p *ConcretePrototype) Clone() Prototype {
    // 深拷贝实现
    clone := &ConcretePrototype{
        Name:   p.Name,
        Data:   make(map[string]interface{}),
        Config: make([]string, len(p.Config)),
    }
    
    // 复制Data
    for key, value := range p.Data {
        clone.Data[key] = value
    }
    
    // 复制Config
    copy(clone.Config, p.Config)
    
    return clone
}

func (p *ConcretePrototype) GetName() string {
    return p.Name
}

func (p *ConcretePrototype) SetName(name string) {
    p.Name = name
}

// 使用JSON序列化的深拷贝
func (p *ConcretePrototype) DeepClone() Prototype {
    data, err := json.Marshal(p)
    if err != nil {
        panic(fmt.Sprintf("failed to marshal: %v", err))
    }
    
    var clone ConcretePrototype
    if err := json.Unmarshal(data, &clone); err != nil {
        panic(fmt.Sprintf("failed to unmarshal: %v", err))
    }
    
    return &clone
}

```

### 3.1.1.4.2 原型注册表模式

```go
package registry

import (
    "fmt"
    "sync"
)

// Prototype 原型接口
type Prototype interface {
    Clone() Prototype
    GetName() string
}

// PrototypeRegistry 原型注册表
type PrototypeRegistry struct {
    prototypes map[string]Prototype
    mu         sync.RWMutex
}

func NewPrototypeRegistry() *PrototypeRegistry {
    return &PrototypeRegistry{
        prototypes: make(map[string]Prototype),
    }
}

// Register 注册原型
func (r *PrototypeRegistry) Register(name string, prototype Prototype) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.prototypes[name] = prototype
}

// Get 获取原型
func (r *PrototypeRegistry) Get(name string) (Prototype, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    prototype, exists := r.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    
    return prototype.Clone(), nil
}

// List 列出所有原型名称
func (r *PrototypeRegistry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    names := make([]string, 0, len(r.prototypes))
    for name := range r.prototypes {
        names = append(names, name)
    }
    
    return names
}

// ConcretePrototypeA 具体原型A
type ConcretePrototypeA struct {
    Name string
    Data map[string]string
}

func (p *ConcretePrototypeA) Clone() Prototype {
    clone := &ConcretePrototypeA{
        Name: p.Name,
        Data: make(map[string]string),
    }
    
    for key, value := range p.Data {
        clone.Data[key] = value
    }
    
    return clone
}

func (p *ConcretePrototypeA) GetName() string {
    return p.Name
}

// ConcretePrototypeB 具体原型B
type ConcretePrototypeB struct {
    Name  string
    Value int
    Tags  []string
}

func (p *ConcretePrototypeB) Clone() Prototype {
    clone := &ConcretePrototypeB{
        Name:  p.Name,
        Value: p.Value,
        Tags:  make([]string, len(p.Tags)),
    }
    
    copy(clone.Tags, p.Tags)
    
    return clone
}

func (p *ConcretePrototypeB) GetName() string {
    return p.Name
}

```

### 3.1.1.4.3 函数式原型模式

```go
package functionalprototype

import (
    "encoding/json"
    "fmt"
)

// Cloneable 可克隆接口
type Cloneable interface {
    Clone() Cloneable
}

// CloneFunc 克隆函数类型
type CloneFunc func() Cloneable

// PrototypeManager 原型管理器
type PrototypeManager struct {
    prototypes map[string]CloneFunc
}

func NewPrototypeManager() *PrototypeManager {
    return &PrototypeManager{
        prototypes: make(map[string]CloneFunc),
    }
}

// Register 注册原型
func (pm *PrototypeManager) Register(name string, cloneFunc CloneFunc) {
    pm.prototypes[name] = cloneFunc
}

// Create 创建对象
func (pm *PrototypeManager) Create(name string) (Cloneable, error) {
    cloneFunc, exists := pm.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    
    return cloneFunc(), nil
}

// Product 产品
type Product struct {
    Name     string
    Price    float64
    Category string
    Tags     []string
}

func (p *Product) Clone() Cloneable {
    clone := &Product{
        Name:     p.Name,
        Price:    p.Price,
        Category: p.Category,
        Tags:     make([]string, len(p.Tags)),
    }
    
    copy(clone.Tags, p.Tags)
    
    return clone
}

// 使用JSON的深拷贝
func (p *Product) DeepClone() *Product {
    data, err := json.Marshal(p)
    if err != nil {
        panic(fmt.Sprintf("failed to marshal: %v", err))
    }
    
    var clone Product
    if err := json.Unmarshal(data, &clone); err != nil {
        panic(fmt.Sprintf("failed to unmarshal: %v", err))
    }
    
    return &clone
}

```

## 3.1.1.5 4. 工程案例

### 3.1.1.5.1 配置对象原型

```go
package configprototype

import (
    "encoding/json"
    "fmt"
    "sync"
)

// Config 配置对象
type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Cache    CacheConfig    `json:"cache"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type ServerConfig struct {
    Port    int    `json:"port"`
    Timeout int    `json:"timeout"`
    Mode    string `json:"mode"`
}

type CacheConfig struct {
    RedisURL string `json:"redis_url"`
    TTL      int    `json:"ttl"`
}

// ConfigPrototype 配置原型
type ConfigPrototype struct {
    configs map[string]*Config
    mu      sync.RWMutex
}

func NewConfigPrototype() *ConfigPrototype {
    return &ConfigPrototype{
        configs: make(map[string]*Config),
    }
}

// Register 注册配置原型
func (cp *ConfigPrototype) Register(name string, config *Config) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    cp.configs[name] = config
}

// Clone 克隆配置
func (cp *ConfigPrototype) Clone(name string) (*Config, error) {
    cp.mu.RLock()
    defer cp.mu.RUnlock()
    
    config, exists := cp.configs[name]
    if !exists {
        return nil, fmt.Errorf("config %s not found", name)
    }
    
    // 使用JSON进行深拷贝
    data, err := json.Marshal(config)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal config: %w", err)
    }
    
    var clone Config
    if err := json.Unmarshal(data, &clone); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &clone, nil
}

// 使用示例
func Example() {
    // 创建原型管理器
    prototype := NewConfigPrototype()
    
    // 注册开发环境配置
    devConfig := &Config{
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "dev",
            Password: "dev123",
        },
        Server: ServerConfig{
            Port:    8080,
            Timeout: 30,
            Mode:    "debug",
        },
        Cache: CacheConfig{
            RedisURL: "redis://localhost:6379",
            TTL:      300,
        },
    }
    prototype.Register("dev", devConfig)
    
    // 克隆配置并修改
    prodConfig, err := prototype.Clone("dev")
    if err != nil {
        panic(err)
    }
    
    // 修改为生产环境配置
    prodConfig.Database.Host = "prod-db.example.com"
    prodConfig.Database.Password = "prod123"
    prodConfig.Server.Port = 80
    prodConfig.Server.Mode = "release"
    prodConfig.Cache.RedisURL = "redis://prod-redis.example.com:6379"
    
    // 注册生产环境配置
    prototype.Register("prod", prodConfig)
}

```

### 3.1.1.5.2 文档模板原型

```go
package documentprototype

import (
    "encoding/json"
    "fmt"
    "time"
)

// Document 文档
type Document struct {
    Title       string            `json:"title"`
    Content     string            `json:"content"`
    Author      string            `json:"author"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    Tags        []string          `json:"tags"`
    Metadata    map[string]string `json:"metadata"`
    Version     int               `json:"version"`
}

// DocumentPrototype 文档原型
type DocumentPrototype struct {
    templates map[string]*Document
}

func NewDocumentPrototype() *DocumentPrototype {
    return &DocumentPrototype{
        templates: make(map[string]*Document),
    }
}

// RegisterTemplate 注册模板
func (dp *DocumentPrototype) RegisterTemplate(name string, template *Document) {
    dp.templates[name] = template
}

// CreateDocument 创建文档
func (dp *DocumentPrototype) CreateDocument(templateName, title, author string) (*Document, error) {
    template, exists := dp.templates[templateName]
    if !exists {
        return nil, fmt.Errorf("template %s not found", templateName)
    }
    
    // 深拷贝模板
    data, err := json.Marshal(template)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal template: %w", err)
    }
    
    var document Document
    if err := json.Unmarshal(data, &document); err != nil {
        return nil, fmt.Errorf("failed to unmarshal template: %w", err)
    }
    
    // 设置新文档的属性
    document.Title = title
    document.Author = author
    document.CreatedAt = time.Now()
    document.UpdatedAt = time.Now()
    document.Version = 1
    
    return &document, nil
}

// 使用示例
func Example() {
    // 创建原型管理器
    prototype := NewDocumentPrototype()
    
    // 注册技术文档模板
    techTemplate := &Document{
        Title:     "技术文档模板",
        Content:   "# {title}\n\n## 概述\n\n## 技术细节\n\n## 总结",
        Author:    "系统",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Tags:      []string{"技术", "文档", "模板"},
        Metadata: map[string]string{
            "category": "技术",
            "format":   "markdown",
        },
        Version: 1,
    }
    prototype.RegisterTemplate("tech", techTemplate)
    
    // 注册用户手册模板
    userTemplate := &Document{
        Title:     "用户手册模板",
        Content:   "# {title}\n\n## 简介\n\n## 使用说明\n\n## 常见问题",
        Author:    "系统",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Tags:      []string{"用户", "手册", "模板"},
        Metadata: map[string]string{
            "category": "用户",
            "format":   "markdown",
        },
        Version: 1,
    }
    prototype.RegisterTemplate("user", userTemplate)
    
    // 创建技术文档
    techDoc, err := prototype.CreateDocument("tech", "Go语言并发编程指南", "张三")
    if err != nil {
        panic(err)
    }
    
    // 创建用户手册
    userDoc, err := prototype.CreateDocument("user", "系统使用手册", "李四")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("技术文档: %+v\n", techDoc)
    fmt.Printf("用户手册: %+v\n", userDoc)
}

```

## 3.1.1.6 5. 批判性分析

### 3.1.1.6.1 优势

1. **性能优化**: 避免昂贵的初始化过程
2. **动态创建**: 运行时动态创建对象
3. **配置复用**: 复用复杂对象的配置
4. **减少耦合**: 客户端与具体类解耦

### 3.1.1.6.2 劣势

1. **深拷贝复杂性**: 复杂对象的深拷贝实现困难
2. **内存开销**: 复制大对象可能消耗大量内存
3. **循环引用**: 处理循环引用问题复杂
4. **性能权衡**: 在某些场景下可能不如直接创建

### 3.1.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 手动实现 | 高 | 中 |
| Java | Cloneable接口 | 中 | 中 |
| C++ | 拷贝构造函数 | 高 | 中 |
| Python | copy模块 | 中 | 低 |

### 3.1.1.6.4 最新趋势

1. **序列化克隆**: 使用JSON/Protocol Buffers
2. **浅拷贝优化**: 结合浅拷贝和懒加载
3. **不可变对象**: 利用不可变性简化克隆
4. **代码生成**: 自动生成克隆代码

## 3.1.1.7 6. 面试题与考点

### 3.1.1.7.1 基础考点

1. **Q**: 原型模式与工厂模式的区别？
   **A**: 原型模式通过复制创建，工厂模式通过实例化创建

2. **Q**: 深拷贝与浅拷贝的区别？
   **A**: 深拷贝复制所有嵌套对象，浅拷贝只复制引用

3. **Q**: 什么时候使用原型模式？
   **A**: 对象创建成本高、需要动态创建、配置复用时

### 3.1.1.7.2 进阶考点

1. **Q**: 如何实现深拷贝？
   **A**: 使用JSON序列化、手动递归复制、反射

2. **Q**: 原型模式在缓存中的应用？
   **A**: 缓存模板对象，快速创建新实例

3. **Q**: 如何处理循环引用？
   **A**: 使用引用映射、序列化方式、标记访问

## 3.1.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 原型模式 | 通过复制原型创建对象的设计模式 | Prototype Pattern |
| 深拷贝 | 复制对象及其所有嵌套对象 | Deep Copy |
| 浅拷贝 | 只复制对象的引用 | Shallow Copy |
| 原型注册表 | 管理原型对象的容器 | Prototype Registry |
| 克隆方法 | 复制对象的方法 | Clone Method |

## 3.1.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 浅拷贝问题 | 共享可变状态 | 实现深拷贝 |
| 循环引用 | 无限递归 | 使用引用映射 |
| 性能问题 | 大对象复制慢 | 使用懒加载 |
| 内存泄漏 | 复制过多对象 | 使用对象池 |

## 3.1.1.10 9. 相关主题

- [单例模式](./01-Singleton-Pattern.md)
- [工厂模式](./02-Factory-Pattern.md)
- [建造者模式](./03-Builder-Pattern.md)
- [抽象工厂模式](./05-Abstract-Factory-Pattern.md)
- [对象池模式](../04-Concurrent-Patterns/02-Object-Pool-Pattern.md)

## 3.1.1.11 10. 学习路径

### 3.1.1.11.1 新手路径

1. 理解原型模式的基本概念
2. 学习浅拷贝和深拷贝的区别
3. 实现简单的原型模式
4. 理解性能优化的意义

### 3.1.1.11.2 进阶路径

1. 学习原型注册表模式
2. 理解深拷贝的实现方法
3. 掌握原型模式的应用场景
4. 学习原型模式的最佳实践

### 3.1.1.11.3 高阶路径

1. 分析原型模式在大型项目中的应用
2. 理解原型模式与性能优化的关系
3. 掌握原型模式的性能调优
4. 学习原型模式的替代方案

---

**相关文档**: [创建型模式总览](./README.md) | [设计模式总览](../README.md)
