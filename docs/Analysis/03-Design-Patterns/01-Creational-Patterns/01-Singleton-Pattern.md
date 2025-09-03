# 3.1.1 单例模式 (Singleton Pattern)

<!-- TOC START -->
- [3.1.1 单例模式 (Singleton Pattern)](#311-单例模式-singleton-pattern)
  - [3.1.1.1 目录](#3111-目录)
  - [3.1.1.2 1. 概述](#3112-1-概述)
    - [3.1.1.2.1 定义](#31121-定义)
    - [3.1.1.2.2 核心特征](#31122-核心特征)
  - [3.1.1.3 2. 理论基础](#3113-2-理论基础)
    - [3.1.1.3.1 数学形式化](#31131-数学形式化)
    - [3.1.1.3.2 范畴论视角](#31132-范畴论视角)
  - [3.1.1.4 3. Go语言实现](#3114-3-go语言实现)
    - [3.1.1.4.1 基础实现](#31141-基础实现)
    - [3.1.1.4.2 配置管理单例](#31142-配置管理单例)
    - [3.1.1.4.3 连接池单例](#31143-连接池单例)
  - [3.1.1.5 4. 工程案例](#3115-4-工程案例)
    - [3.1.1.5.1 日志系统单例](#31151-日志系统单例)
    - [3.1.1.5.2 缓存管理器单例](#31152-缓存管理器单例)
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

单例模式确保一个类只有一个实例，并提供全局访问点。

**形式化定义**:
$$Singleton = (Instance, GetInstance, Constructor, State)$$

其中：

- $Instance$ 是唯一实例
- $GetInstance()$ 是获取实例的方法
- $Constructor$ 是私有构造函数
- $State$ 是实例状态

### 3.1.1.2.2 核心特征

- **唯一性**: 确保全局只有一个实例
- **全局访问**: 提供统一的访问接口
- **延迟初始化**: 支持懒加载
- **线程安全**: 保证并发环境下的正确性

## 3.1.1.3 2. 理论基础

### 3.1.1.3.1 数学形式化

**定义 2.1** (单例模式): 单例模式是一个四元组 $S = (C, I, M, V)$

其中：

- $C$ 是类定义
- $I$ 是实例集合，$|I| = 1$
- $M$ 是方法集合
- $V$ 是验证规则

**定理 2.1** (唯一性保证): 对于任意时刻 $t$，$|I_t| \leq 1$

**证明**: 通过构造函数私有化和同步机制保证。

### 3.1.1.3.2 范畴论视角

在范畴论中，单例模式可以表示为：

$$Singleton : 1 \rightarrow C$$

其中 $1$ 是单元素集合，$C$ 是对象类。

## 3.1.1.4 3. Go语言实现

### 3.1.1.4.1 基础实现

```go
package singleton

import (
    "fmt"
    "sync"
)

// 单例结构体
type Singleton struct {
    data string
    mu   sync.RWMutex
}

var (
    instance *Singleton
    once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "initialized",
        }
    })
    return instance
}

// GetData 获取数据
func (s *Singleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

// SetData 设置数据
func (s *Singleton) SetData(data string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data = data
}
```

### 3.1.1.4.2 配置管理单例

```go
package config

import (
    "encoding/json"
    "os"
    "sync"
)

// Config 配置结构
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

// ConfigManager 配置管理器单例
type ConfigManager struct {
    config *Config
    mu     sync.RWMutex
}

var (
    configInstance *ConfigManager
    configOnce     sync.Once
)

// GetConfigInstance 获取配置实例
func GetConfigInstance() *ConfigManager {
    configOnce.Do(func() {
        configInstance = &ConfigManager{}
    })
    return configInstance
}

// LoadConfig 加载配置
func (cm *ConfigManager) LoadConfig(filename string) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()
    
    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return fmt.Errorf("failed to decode config: %w", err)
    }
    
    cm.config = &config
    return nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *Config {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.config
}
```

### 3.1.1.4.3 连接池单例

```go
package connection

import (
    "database/sql"
    "fmt"
    "sync"
    _ "github.com/lib/pq"
)

// ConnectionPool 连接池单例
type ConnectionPool struct {
    db   *sql.DB
    mu   sync.RWMutex
    size int
}

var (
    poolInstance *ConnectionPool
    poolOnce     sync.Once
)

// GetPoolInstance 获取连接池实例
func GetPoolInstance() *ConnectionPool {
    poolOnce.Do(func() {
        poolInstance = &ConnectionPool{
            size: 10,
        }
    })
    return poolInstance
}

// Initialize 初始化连接池
func (cp *ConnectionPool) Initialize(dsn string) error {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    db.SetMaxOpenConns(cp.size)
    db.SetMaxIdleConns(cp.size / 2)
    
    cp.db = db
    return nil
}

// GetDB 获取数据库连接
func (cp *ConnectionPool) GetDB() *sql.DB {
    cp.mu.RLock()
    defer cp.mu.RUnlock()
    return cp.db
}
```

## 3.1.1.5 4. 工程案例

### 3.1.1.5.1 日志系统单例

```go
package logger

import (
    "log"
    "os"
    "sync"
)

// Logger 日志单例
type Logger struct {
    logger *log.Logger
    mu     sync.Mutex
}

var (
    loggerInstance *Logger
    loggerOnce     sync.Once
)

// GetLogger 获取日志实例
func GetLogger() *Logger {
    loggerOnce.Do(func() {
        loggerInstance = &Logger{
            logger: log.New(os.Stdout, "[APP] ", log.LstdFlags),
        }
    })
    return loggerInstance
}

// Info 信息日志
func (l *Logger) Info(format string, v ...interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Printf("[INFO] "+format, v...)
}

// Error 错误日志
func (l *Logger) Error(format string, v ...interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Printf("[ERROR] "+format, v...)
}
```

### 3.1.1.5.2 缓存管理器单例

```go
package cache

import (
    "sync"
    "time"
)

// CacheItem 缓存项
type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}

// CacheManager 缓存管理器单例
type CacheManager struct {
    cache map[string]CacheItem
    mu    sync.RWMutex
}

var (
    cacheInstance *CacheManager
    cacheOnce     sync.Once
)

// GetCacheInstance 获取缓存实例
func GetCacheInstance() *CacheManager {
    cacheOnce.Do(func() {
        cacheInstance = &CacheManager{
            cache: make(map[string]CacheItem),
        }
    })
    return cacheInstance
}

// Set 设置缓存
func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.cache[key] = CacheItem{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

// Get 获取缓存
func (cm *CacheManager) Get(key string) (interface{}, bool) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    item, exists := cm.cache[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(item.Expiration) {
        delete(cm.cache, key)
        return nil, false
    }
    
    return item.Value, true
}
```

## 3.1.1.6 5. 批判性分析

### 3.1.1.6.1 优势

1. **资源管理**: 有效管理全局资源
2. **状态一致性**: 保证全局状态一致
3. **性能优化**: 避免重复创建对象
4. **配置集中**: 统一配置管理

### 3.1.1.6.2 劣势

1. **全局状态**: 引入全局状态，增加复杂性
2. **测试困难**: 单例状态影响测试
3. **并发问题**: 需要额外的同步机制
4. **扩展性差**: 难以扩展和修改

### 3.1.1.6.3 行业对比

| 语言 | 实现方式 | 线程安全 | 性能 |
|------|----------|----------|------|
| Go | sync.Once | 内置保证 | 高 |
| Java | volatile + synchronized | 需要手动保证 | 中 |
| C++ | std::call_once | 内置保证 | 高 |
| Python | 装饰器 | 需要手动保证 | 中 |

### 3.1.1.6.4 最新趋势

1. **依赖注入**: 替代单例模式
2. **函数式编程**: 避免全局状态
3. **微服务架构**: 减少单例使用
4. **容器化**: 通过容器管理状态

## 3.1.1.7 6. 面试题与考点

### 3.1.1.7.1 基础考点

1. **Q**: 如何实现线程安全的单例模式？
   **A**: 使用 `sync.Once` 或 `sync.Mutex`

2. **Q**: 单例模式的优缺点是什么？
   **A**: 优点：资源管理、状态一致；缺点：全局状态、测试困难

3. **Q**: 如何避免单例模式的全局状态问题？
   **A**: 使用依赖注入、函数式编程、容器化

### 3.1.1.7.2 进阶考点

1. **Q**: 单例模式在微服务架构中的应用？
   **A**: 主要用于配置管理、连接池、日志系统

2. **Q**: 如何测试单例模式？
   **A**: 使用接口抽象、依赖注入、测试容器

3. **Q**: 单例模式与设计原则的冲突？
   **A**: 违反单一职责、开闭原则，但符合DRY原则

## 3.1.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 单例模式 | 确保一个类只有一个实例的设计模式 | Singleton Pattern |
| 全局访问点 | 提供统一访问实例的接口 | Global Access Point |
| 延迟初始化 | 在首次使用时才创建实例 | Lazy Initialization |
| 线程安全 | 在并发环境下正确工作的特性 | Thread Safety |
| 同步机制 | 保证并发安全的机制 | Synchronization Mechanism |

## 3.1.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 双重检查锁定 | 在Go中不必要 | 使用 `sync.Once` |
| 全局状态污染 | 影响测试和调试 | 使用依赖注入 |
| 内存泄漏 | 单例持有大量资源 | 实现清理机制 |
| 循环依赖 | 单例间相互依赖 | 重新设计架构 |

## 3.1.1.10 9. 相关主题

- [工厂模式](./02-Factory-Pattern.md)
- [建造者模式](./03-Builder-Pattern.md)
- [原型模式](./04-Prototype-Pattern.md)
- [抽象工厂模式](./05-Abstract-Factory-Pattern.md)
- [依赖注入模式](../04-Concurrent-Patterns/01-Dependency-Injection.md)

## 3.1.1.11 10. 学习路径

### 3.1.1.11.1 新手路径

1. 理解单例模式的基本概念
2. 学习Go的 `sync.Once` 机制
3. 实现简单的单例模式
4. 理解线程安全的重要性

### 3.1.1.11.2 进阶路径

1. 学习不同场景下的单例应用
2. 理解单例模式的优缺点
3. 学习替代方案（依赖注入）
4. 掌握测试单例模式的方法

### 3.1.1.11.3 高阶路径

1. 分析单例模式在大型项目中的应用
2. 理解单例模式与架构设计的关系
3. 掌握单例模式的性能优化
4. 学习单例模式的最佳实践

---

**相关文档**: [创建型模式总览](./README.md) | [设计模式总览](../README.md)
