# 3.2.1 装饰器模式 (Decorator Pattern)

<!-- TOC START -->
- [3.2.1 装饰器模式 (Decorator Pattern)](#321-装饰器模式-decorator-pattern)
  - [3.2.1.1 目录](#3211-目录)
  - [3.2.1.2 1. 概述](#3212-1-概述)
    - [3.2.1.2.1 定义](#32121-定义)
    - [3.2.1.2.2 核心特征](#32122-核心特征)
  - [3.2.1.3 2. 理论基础](#3213-2-理论基础)
    - [3.2.1.3.1 数学形式化](#32131-数学形式化)
    - [3.2.1.3.2 范畴论视角](#32132-范畴论视角)
  - [3.2.1.4 3. Go语言实现](#3214-3-go语言实现)
    - [3.2.1.4.1 基础装饰器模式](#32141-基础装饰器模式)
    - [3.2.1.4.2 函数式装饰器模式](#32142-函数式装饰器模式)
    - [3.2.1.4.3 中间件装饰器模式](#32143-中间件装饰器模式)
  - [3.2.1.5 4. 工程案例](#3215-4-工程案例)
    - [3.2.1.5.1 HTTP客户端装饰器](#32151-http客户端装饰器)
    - [3.2.1.5.2 数据库连接装饰器](#32152-数据库连接装饰器)
    - [3.2.1.5.3 日志装饰器](#32153-日志装饰器)
  - [3.2.1.6 5. 批判性分析](#3216-5-批判性分析)
    - [3.2.1.6.1 优势](#32161-优势)
    - [3.2.1.6.2 劣势](#32162-劣势)
    - [3.2.1.6.3 行业对比](#32163-行业对比)
    - [3.2.1.6.4 最新趋势](#32164-最新趋势)
  - [3.2.1.7 6. 面试题与考点](#3217-6-面试题与考点)
    - [3.2.1.7.1 基础考点](#32171-基础考点)
    - [3.2.1.7.2 进阶考点](#32172-进阶考点)
  - [3.2.1.8 7. 术语表](#3218-7-术语表)
  - [3.2.1.9 8. 常见陷阱](#3219-8-常见陷阱)
  - [3.2.1.10 9. 相关主题](#32110-9-相关主题)
  - [3.2.1.11 10. 学习路径](#32111-10-学习路径)
    - [3.2.1.11.1 新手路径](#321111-新手路径)
    - [3.2.1.11.2 进阶路径](#321112-进阶路径)
    - [3.2.1.11.3 高阶路径](#321113-高阶路径)
<!-- TOC END -->

## 3.2.1.1 目录

## 3.2.1.2 1. 概述

### 3.2.1.2.1 定义

装饰器模式动态地给对象添加一些额外的职责，就增加功能来说，装饰器模式比生成子类更为灵活。

**形式化定义**:
$$Decorator = (Component, ConcreteComponent, Decorator, ConcreteDecorator, DecoratorChain)$$

其中：

- $Component$ 是组件接口
- $ConcreteComponent$ 是具体组件
- $Decorator$ 是装饰器抽象类
- $ConcreteDecorator$ 是具体装饰器
- $DecoratorChain$ 是装饰器链

### 3.2.1.2.2 核心特征

- **动态扩展**: 运行时动态添加功能
- **组合优先**: 使用组合而非继承
- **透明性**: 装饰器与组件接口一致
- **灵活性**: 支持任意组合装饰器

## 3.2.1.3 2. 理论基础

### 3.2.1.3.1 数学形式化

**定义 2.1** (装饰器模式): 装饰器模式是一个六元组 $D = (C, D, M, R, F, V)$

其中：

- $C$ 是组件集合
- $D$ 是装饰器集合
- $M$ 是装饰映射，$M: C \times D \rightarrow C$
- $R$ 是关系约束
- $F$ 是功能集合
- $V$ 是验证规则

**定理 2.1** (功能叠加性): 对于任意组件 $c \in C$ 和装饰器 $d \in D$，$M(c, d)$ 包含 $c$ 的所有功能加上 $d$ 的额外功能

**证明**: 由装饰器的组合性质保证。

### 3.2.1.3.2 范畴论视角

在范畴论中，装饰器模式可以表示为：

$$Decorator : Component \times Decorator \rightarrow Component$$

其中 $Component$ 和 $Decorator$ 是对象范畴。

## 3.2.1.4 3. Go语言实现

### 3.2.1.4.1 基础装饰器模式

```go
package decorator

import "fmt"

// Component 组件接口
type Component interface {
    Operation() string
}

// ConcreteComponent 具体组件
type ConcreteComponent struct {
    name string
}

func (c *ConcreteComponent) Operation() string {
    return fmt.Sprintf("ConcreteComponent(%s)", c.name)
}

// Decorator 装饰器抽象类
type Decorator struct {
    component Component
}

func (d *Decorator) Operation() string {
    return d.component.Operation()
}

// ConcreteDecoratorA 具体装饰器A
type ConcreteDecoratorA struct {
    Decorator
    addedState string
}

func NewConcreteDecoratorA(component Component) *ConcreteDecoratorA {
    return &ConcreteDecoratorA{
        Decorator:  Decorator{component: component},
        addedState: "DecoratorA",
    }
}

func (d *ConcreteDecoratorA) Operation() string {
    return fmt.Sprintf("%s + %s", d.Decorator.Operation(), d.addedState)
}

// ConcreteDecoratorB 具体装饰器B
type ConcreteDecoratorB struct {
    Decorator
    addedBehavior func() string
}

func NewConcreteDecoratorB(component Component) *ConcreteDecoratorB {
    return &ConcreteDecoratorB{
        Decorator: Decorator{component: component},
        addedBehavior: func() string {
            return "DecoratorB"
        },
    }
}

func (d *ConcreteDecoratorB) Operation() string {
    return fmt.Sprintf("%s + %s", d.Decorator.Operation(), d.addedBehavior())
}
```

### 3.2.1.4.2 函数式装饰器模式

```go
package functionaldecorator

import "fmt"

// ComponentFunc 组件函数类型
type ComponentFunc func() string

// DecoratorFunc 装饰器函数类型
type DecoratorFunc func(ComponentFunc) ComponentFunc

// LoggingDecorator 日志装饰器
func LoggingDecorator(prefix string) DecoratorFunc {
    return func(component ComponentFunc) ComponentFunc {
        return func() string {
            result := component()
            fmt.Printf("[%s] %s\n", prefix, result)
            return result
        }
    }
}

// TimingDecorator 计时装饰器
func TimingDecorator() DecoratorFunc {
    return func(component ComponentFunc) ComponentFunc {
        return func() string {
            start := time.Now()
            result := component()
            duration := time.Since(start)
            fmt.Printf("Operation took %v\n", duration)
            return result
        }
    }
}

// CachingDecorator 缓存装饰器
func CachingDecorator() DecoratorFunc {
    cache := make(map[string]string)
    return func(component ComponentFunc) ComponentFunc {
        return func() string {
            key := "cache_key" // 实际应用中需要更复杂的键生成
            if cached, exists := cache[key]; exists {
                return cached
            }
            result := component()
            cache[key] = result
            return result
        }
    }
}

// 使用示例
func Example() {
    // 基础组件
    component := func() string {
        return "Hello, World!"
    }
    
    // 应用装饰器
    decorated := LoggingDecorator("INFO")(
        TimingDecorator()(
            CachingDecorator()(component),
        ),
    )
    
    // 执行
    result := decorated()
    fmt.Println("Result:", result)
}
```

### 3.2.1.4.3 中间件装饰器模式

```go
package middleware

import (
    "fmt"
    "net/http"
    "time"
)

// Handler 处理器接口
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HandlerFunc 处理器函数类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}

// Middleware 中间件类型
type Middleware func(Handler) Handler

// LoggingMiddleware 日志中间件
func LoggingMiddleware() Middleware {
    return func(next Handler) Handler {
        return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)
            
            next.ServeHTTP(w, r)
            
            fmt.Printf("Completed %s %s in %v\n", r.Method, r.URL.Path, time.Since(start))
        })
    }
}

// AuthMiddleware 认证中间件
func AuthMiddleware(token string) Middleware {
    return func(next Handler) Handler {
        return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authToken := r.Header.Get("Authorization")
            if authToken != token {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// CORSMiddleware CORS中间件
func CORSMiddleware() Middleware {
    return func(next Handler) Handler {
        return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Compose 组合多个中间件
func Compose(middlewares ...Middleware) Middleware {
    return func(next Handler) Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}
```

## 3.2.1.5 4. 工程案例

### 3.2.1.5.1 HTTP客户端装饰器

```go
package httpdecorator

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// HTTPClient HTTP客户端接口
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

// BaseHTTPClient 基础HTTP客户端
type BaseHTTPClient struct {
    client *http.Client
}

func NewBaseHTTPClient() *BaseHTTPClient {
    return &BaseHTTPClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (b *BaseHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return b.client.Do(req)
}

// RetryDecorator 重试装饰器
type RetryDecorator struct {
    client HTTPClient
    maxRetries int
    backoff    time.Duration
}

func NewRetryDecorator(client HTTPClient, maxRetries int, backoff time.Duration) *RetryDecorator {
    return &RetryDecorator{
        client:     client,
        maxRetries: maxRetries,
        backoff:    backoff,
    }
}

func (r *RetryDecorator) Do(req *http.Request) (*http.Response, error) {
    var lastErr error
    
    for i := 0; i <= r.maxRetries; i++ {
        resp, err := r.client.Do(req)
        if err == nil {
            return resp, nil
        }
        
        lastErr = err
        if i < r.maxRetries {
            time.Sleep(r.backoff * time.Duration(i+1))
        }
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", r.maxRetries, lastErr)
}

// LoggingDecorator 日志装饰器
type LoggingDecorator struct {
    client HTTPClient
}

func NewLoggingDecorator(client HTTPClient) *LoggingDecorator {
    return &LoggingDecorator{client: client}
}

func (l *LoggingDecorator) Do(req *http.Request) (*http.Response, error) {
    start := time.Now()
    fmt.Printf("HTTP Request: %s %s\n", req.Method, req.URL.String())
    
    resp, err := l.client.Do(req)
    
    duration := time.Since(start)
    if err != nil {
        fmt.Printf("HTTP Error: %s %s - %v (took %v)\n", req.Method, req.URL.String(), err, duration)
    } else {
        fmt.Printf("HTTP Response: %s %s - %d (took %v)\n", req.Method, req.URL.String(), resp.StatusCode, duration)
    }
    
    return resp, err
}

// CachingDecorator 缓存装饰器
type CachingDecorator struct {
    client HTTPClient
    cache  map[string]*CachedResponse
}

type CachedResponse struct {
    Response *http.Response
    Expiry   time.Time
}

func NewCachingDecorator(client HTTPClient) *CachingDecorator {
    return &CachingDecorator{
        client: client,
        cache:  make(map[string]*CachedResponse),
    }
}

func (c *CachingDecorator) Do(req *http.Request) (*http.Response, error) {
    // 只缓存GET请求
    if req.Method != "GET" {
        return c.client.Do(req)
    }
    
    key := fmt.Sprintf("%s %s", req.Method, req.URL.String())
    
    // 检查缓存
    if cached, exists := c.cache[key]; exists && time.Now().Before(cached.Expiry) {
        return cached.Response, nil
    }
    
    // 执行请求
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    
    // 缓存响应（只缓存成功的响应）
    if resp.StatusCode == http.StatusOK {
        c.cache[key] = &CachedResponse{
            Response: resp,
            Expiry:   time.Now().Add(5 * time.Minute),
        }
    }
    
    return resp, nil
}
```

### 3.2.1.5.2 数据库连接装饰器

```go
package dbdecorator

import (
    "database/sql"
    "fmt"
    "time"
)

// DBInterface 数据库接口
type DBInterface interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    Exec(query string, args ...interface{}) (sql.Result, error)
    Close() error
}

// BaseDB 基础数据库
type BaseDB struct {
    db *sql.DB
}

func NewBaseDB(driver, dsn string) (*BaseDB, error) {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, err
    }
    
    return &BaseDB{db: db}, nil
}

func (b *BaseDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return b.db.Query(query, args...)
}

func (b *BaseDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    return b.db.Exec(query, args...)
}

func (b *BaseDB) Close() error {
    return b.db.Close()
}

// MetricsDecorator 指标装饰器
type MetricsDecorator struct {
    db      DBInterface
    metrics map[string]int64
}

func NewMetricsDecorator(db DBInterface) *MetricsDecorator {
    return &MetricsDecorator{
        db:      db,
        metrics: make(map[string]int64),
    }
}

func (m *MetricsDecorator) Query(query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    rows, err := m.db.Query(query, args...)
    duration := time.Since(start)
    
    m.metrics["query_count"]++
    m.metrics["query_duration_ns"] += duration.Nanoseconds()
    
    return rows, err
}

func (m *MetricsDecorator) Exec(query string, args ...interface{}) (sql.Result, error) {
    start := time.Now()
    result, err := m.db.Exec(query, args...)
    duration := time.Since(start)
    
    m.metrics["exec_count"]++
    m.metrics["exec_duration_ns"] += duration.Nanoseconds()
    
    return result, err
}

func (m *MetricsDecorator) Close() error {
    return m.db.Close()
}

// RetryDecorator 重试装饰器
type RetryDecorator struct {
    db         DBInterface
    maxRetries int
    backoff    time.Duration
}

func NewRetryDecorator(db DBInterface, maxRetries int, backoff time.Duration) *RetryDecorator {
    return &RetryDecorator{
        db:         db,
        maxRetries: maxRetries,
        backoff:    backoff,
    }
}

func (r *RetryDecorator) Query(query string, args ...interface{}) (*sql.Rows, error) {
    var lastErr error
    
    for i := 0; i <= r.maxRetries; i++ {
        rows, err := r.db.Query(query, args...)
        if err == nil {
            return rows, nil
        }
        
        lastErr = err
        if i < r.maxRetries {
            time.Sleep(r.backoff * time.Duration(i+1))
        }
    }
    
    return nil, fmt.Errorf("query failed after %d retries: %w", r.maxRetries, lastErr)
}

func (r *RetryDecorator) Exec(query string, args ...interface{}) (sql.Result, error) {
    var lastErr error
    
    for i := 0; i <= r.maxRetries; i++ {
        result, err := r.db.Exec(query, args...)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        if i < r.maxRetries {
            time.Sleep(r.backoff * time.Duration(i+1))
        }
    }
    
    return nil, fmt.Errorf("exec failed after %d retries: %w", r.maxRetries, lastErr)
}

func (r *RetryDecorator) Close() error {
    return r.db.Close()
}
```

### 3.2.1.5.3 日志装饰器

```go
package logdecorator

import (
    "fmt"
    "log"
    "os"
    "time"
)

// Logger 日志接口
type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
    Warn(message string)
}

// BaseLogger 基础日志器
type BaseLogger struct {
    logger *log.Logger
}

func NewBaseLogger() *BaseLogger {
    return &BaseLogger{
        logger: log.New(os.Stdout, "", log.LstdFlags),
    }
}

func (b *BaseLogger) Info(message string) {
    b.logger.Printf("[INFO] %s", message)
}

func (b *BaseLogger) Error(message string) {
    b.logger.Printf("[ERROR] %s", message)
}

func (b *BaseLogger) Debug(message string) {
    b.logger.Printf("[DEBUG] %s", message)
}

func (b *BaseLogger) Warn(message string) {
    b.logger.Printf("[WARN] %s", message)
}

// TimestampDecorator 时间戳装饰器
type TimestampDecorator struct {
    logger Logger
}

func NewTimestampDecorator(logger Logger) *TimestampDecorator {
    return &TimestampDecorator{logger: logger}
}

func (t *TimestampDecorator) Info(message string) {
    timestamp := time.Now().Format(time.RFC3339)
    t.logger.Info(fmt.Sprintf("[%s] %s", timestamp, message))
}

func (t *TimestampDecorator) Error(message string) {
    timestamp := time.Now().Format(time.RFC3339)
    t.logger.Error(fmt.Sprintf("[%s] %s", timestamp, message))
}

func (t *TimestampDecorator) Debug(message string) {
    timestamp := time.Now().Format(time.RFC3339)
    t.logger.Debug(fmt.Sprintf("[%s] %s", timestamp, message))
}

func (t *TimestampDecorator) Warn(message string) {
    timestamp := time.Now().Format(time.RFC3339)
    t.logger.Warn(fmt.Sprintf("[%s] %s", timestamp, message))
}

// FileDecorator 文件装饰器
type FileDecorator struct {
    logger Logger
    file   *os.File
}

func NewFileDecorator(logger Logger, filename string) (*FileDecorator, error) {
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }
    
    return &FileDecorator{
        logger: logger,
        file:   file,
    }, nil
}

func (f *FileDecorator) Info(message string) {
    f.logger.Info(message)
    f.file.WriteString(fmt.Sprintf("[INFO] %s\n", message))
}

func (f *FileDecorator) Error(message string) {
    f.logger.Error(message)
    f.file.WriteString(fmt.Sprintf("[ERROR] %s\n", message))
}

func (f *FileDecorator) Debug(message string) {
    f.logger.Debug(message)
    f.file.WriteString(fmt.Sprintf("[DEBUG] %s\n", message))
}

func (f *FileDecorator) Warn(message string) {
    f.logger.Warn(message)
    f.file.WriteString(fmt.Sprintf("[WARN] %s\n", message))
}

func (f *FileDecorator) Close() error {
    return f.file.Close()
}
```

## 3.2.1.6 5. 批判性分析

### 3.2.1.6.1 优势

1. **动态扩展**: 运行时动态添加功能
2. **组合优先**: 使用组合而非继承
3. **单一职责**: 每个装饰器只负责一个功能
4. **开闭原则**: 对扩展开放，对修改关闭

### 3.2.1.6.2 劣势

1. **复杂性**: 装饰器链可能变得复杂
2. **性能开销**: 多层装饰器可能影响性能
3. **调试困难**: 装饰器链调试复杂
4. **过度使用**: 可能导致过度设计

### 3.2.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 结构体 | 高 | 中 |
| Java | 抽象类 + 接口 | 中 | 中 |
| C++ | 虚函数 | 中 | 中 |
| Python | 装饰器语法 | 高 | 低 |

### 3.2.1.6.4 最新趋势

1. **函数式装饰器**: 使用高阶函数
2. **中间件模式**: 在Web框架中广泛应用
3. **AOP装饰器**: 面向切面编程
4. **配置驱动**: 基于配置的装饰器

## 3.2.1.7 6. 面试题与考点

### 3.2.1.7.1 基础考点

1. **Q**: 装饰器模式与继承的区别？
   **A**: 装饰器使用组合，继承使用类层次结构

2. **Q**: 什么时候使用装饰器模式？
   **A**: 需要动态添加功能、避免类爆炸时

3. **Q**: 装饰器模式的优缺点？
   **A**: 优点：灵活性、组合优先；缺点：复杂性、性能开销

### 3.2.1.7.2 进阶考点

1. **Q**: 如何设计装饰器链？
   **A**: 使用组合模式、中间件模式

2. **Q**: 装饰器模式在Web框架中的应用？
   **A**: 中间件、拦截器、过滤器

3. **Q**: 如何处理装饰器的性能问题？
   **A**: 缓存、异步、批量处理

## 3.2.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 装饰器模式 | 动态添加功能的设计模式 | Decorator Pattern |
| 组件 | 被装饰的对象 | Component |
| 装饰器 | 添加功能的包装器 | Decorator |
| 装饰器链 | 多个装饰器的组合 | Decorator Chain |
| 中间件 | Web框架中的装饰器 | Middleware |

## 3.2.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 装饰器链过长 | 性能下降 | 优化装饰器顺序 |
| 循环依赖 | 装饰器相互依赖 | 重新设计架构 |
| 状态共享 | 装饰器间状态混乱 | 明确状态边界 |
| 过度装饰 | 功能过度复杂 | 简化装饰器设计 |

## 3.2.1.10 9. 相关主题

- [适配器模式](./01-Adapter-Pattern.md)
- [代理模式](./03-Proxy-Pattern.md)
- [外观模式](./04-Facade-Pattern.md)
- [桥接模式](./05-Bridge-Pattern.md)
- [组合模式](./06-Composite-Pattern.md)

## 3.2.1.11 10. 学习路径

### 3.2.1.11.1 新手路径

1. 理解装饰器模式的基本概念
2. 学习基础装饰器实现
3. 实现简单的装饰器链
4. 理解组合优于继承的原则

### 3.2.1.11.2 进阶路径

1. 学习函数式装饰器
2. 理解中间件模式
3. 掌握装饰器的应用场景
4. 学习装饰器的最佳实践

### 3.2.1.11.3 高阶路径

1. 分析装饰器在大型项目中的应用
2. 理解装饰器与架构设计的关系
3. 掌握装饰器的性能优化
4. 学习装饰器的替代方案

---

**相关文档**: [结构型模式总览](./README.md) | [设计模式总览](../README.md)
