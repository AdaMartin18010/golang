# 3.2.1 适配器模式 (Adapter Pattern)

<!-- TOC START -->
- [3.2.1 适配器模式 (Adapter Pattern)](#321-适配器模式-adapter-pattern)
  - [3.2.1.1 目录](#3211-目录)
  - [3.2.1.2 1. 概述](#3212-1-概述)
    - [3.2.1.2.1 定义](#32121-定义)
    - [3.2.1.2.2 核心特征](#32122-核心特征)
  - [3.2.1.3 2. 理论基础](#3213-2-理论基础)
    - [3.2.1.3.1 数学形式化](#32131-数学形式化)
    - [3.2.1.3.2 范畴论视角](#32132-范畴论视角)
  - [3.2.1.4 3. Go语言实现](#3214-3-go语言实现)
    - [3.2.1.4.1 类适配器模式](#32141-类适配器模式)
    - [3.2.1.4.2 对象适配器模式](#32142-对象适配器模式)
    - [3.2.1.4.3 函数适配器模式](#32143-函数适配器模式)
  - [3.2.1.5 4. 工程案例](#3215-4-工程案例)
    - [3.2.1.5.1 第三方API适配器](#32151-第三方api适配器)
    - [3.2.1.5.2 数据库适配器](#32152-数据库适配器)
    - [3.2.1.5.3 日志适配器](#32153-日志适配器)
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

适配器模式将一个类的接口转换成客户希望的另外一个接口，使得原本由于接口不兼容而不能一起工作的那些类可以一起工作。

**形式化定义**:
$$Adapter = (Target, Adaptee, Adapter, Client, InterfaceMapping)$$

其中：

- $Target$ 是目标接口
- $Adaptee$ 是适配者类
- $Adapter$ 是适配器类
- $Client$ 是客户端
- $InterfaceMapping$ 是接口映射关系

### 3.2.1.2.2 核心特征

- **接口转换**: 将不兼容的接口转换为兼容接口
- **透明性**: 客户端无需知道适配器的存在
- **复用性**: 复用现有类而不修改其代码
- **扩展性**: 支持多种适配方式

## 3.2.1.3 2. 理论基础

### 3.2.1.3.1 数学形式化

**定义 2.1** (适配器模式): 适配器模式是一个五元组 $A = (T, A, M, R, V)$

其中：

- $T$ 是目标接口集合
- $A$ 是适配者集合
- $M$ 是映射函数，$M: A \rightarrow T$
- $R$ 是关系约束
- $V$ 是验证规则

**定理 2.1** (接口兼容性): 对于任意适配者 $a \in A$，存在目标接口 $t \in T$ 使得 $M(a) = t$

**证明**: 由适配器的实现保证。

### 3.2.1.3.2 范畴论视角

在范畴论中，适配器模式可以表示为：

$$Adapter : Adaptee \rightarrow Target$$

其中 $Adaptee$ 和 $Target$ 是对象范畴。

## 3.2.1.4 3. Go语言实现

### 3.2.1.4.1 类适配器模式

```go
package adapter

import "fmt"

// Target 目标接口
type Target interface {
    Request() string
}

// Adaptee 适配者
type Adaptee struct {
    data string
}

func (a *Adaptee) SpecificRequest() string {
    return fmt.Sprintf("Adaptee: %s", a.data)
}

// Adapter 适配器
type Adapter struct {
    *Adaptee // 嵌入适配者
}

func NewAdapter(data string) *Adapter {
    return &Adapter{
        Adaptee: &Adaptee{data: data},
    }
}

func (a *Adapter) Request() string {
    // 调用适配者的方法并转换结果
    result := a.SpecificRequest()
    return fmt.Sprintf("Adapter: %s", result)
}

```

### 3.2.1.4.2 对象适配器模式

```go
package objectadapter

import "fmt"

// Target 目标接口
type Target interface {
    Request() string
}

// Adaptee 适配者
type Adaptee struct {
    data string
}

func (a *Adaptee) SpecificRequest() string {
    return fmt.Sprintf("Adaptee: %s", a.data)
}

// Adapter 适配器
type Adapter struct {
    adaptee *Adaptee // 组合适配者
}

func NewAdapter(adaptee *Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    // 调用适配者的方法并转换结果
    result := a.adaptee.SpecificRequest()
    return fmt.Sprintf("Adapter: %s", result)
}

```

### 3.2.1.4.3 函数适配器模式

```go
package functionadapter

import (
    "fmt"
    "strconv"
)

// TargetFunc 目标函数类型
type TargetFunc func(int) string

// AdapteeFunc 适配者函数类型
type AdapteeFunc func(string) int

// AdapterFunc 适配器函数
func AdapterFunc(adaptee AdapteeFunc) TargetFunc {
    return func(value int) string {
        // 将int转换为string，调用适配者函数，再将结果转换回string
        strValue := strconv.Itoa(value)
        result := adaptee(strValue)
        return fmt.Sprintf("Adapted: %d", result)
    }
}

// 使用示例
func Example() {
    // 适配者函数：将字符串转换为数字并返回
    adapteeFunc := func(s string) int {
        if num, err := strconv.Atoi(s); err == nil {
            return num * 2
        }
        return 0
    }
    
    // 创建适配器
    targetFunc := AdapterFunc(adapteeFunc)
    
    // 使用目标接口
    result := targetFunc(5)
    fmt.Println(result) // 输出: Adapted: 10
}

```

## 3.2.1.5 4. 工程案例

### 3.2.1.5.1 第三方API适配器

```go
package apiadapter

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// PaymentService 支付服务接口
type PaymentService interface {
    ProcessPayment(amount float64, currency string) (string, error)
    GetPaymentStatus(paymentID string) (string, error)
}

// StripeAPI Stripe API（第三方服务）
type StripeAPI struct {
    apiKey string
    client *http.Client
}

func NewStripeAPI(apiKey string) *StripeAPI {
    return &StripeAPI{
        apiKey: apiKey,
        client: &http.Client{Timeout: 30 * time.Second},
    }
}

func (s *StripeAPI) CreatePaymentIntent(amount int64, currency string) (*StripePaymentIntent, error) {
    // 模拟Stripe API调用
    return &StripePaymentIntent{
        ID:     "pi_stripe_123",
        Amount: amount,
        Status: "succeeded",
    }, nil
}

func (s *StripeAPI) RetrievePaymentIntent(paymentIntentID string) (*StripePaymentIntent, error) {
    // 模拟Stripe API调用
    return &StripePaymentIntent{
        ID:     paymentIntentID,
        Amount: 1000,
        Status: "succeeded",
    }, nil
}

type StripePaymentIntent struct {
    ID     string `json:"id"`
    Amount int64  `json:"amount"`
    Status string `json:"status"`
}

// StripeAdapter Stripe适配器
type StripeAdapter struct {
    stripe *StripeAPI
}

func NewStripeAdapter(apiKey string) *StripeAdapter {
    return &StripeAdapter{
        stripe: NewStripeAPI(apiKey),
    }
}

func (s *StripeAdapter) ProcessPayment(amount float64, currency string) (string, error) {
    // 将金额转换为分（Stripe使用最小货币单位）
    amountInCents := int64(amount * 100)
    
    // 调用Stripe API
    paymentIntent, err := s.stripe.CreatePaymentIntent(amountInCents, currency)
    if err != nil {
        return "", fmt.Errorf("failed to create payment intent: %w", err)
    }
    
    return paymentIntent.ID, nil
}

func (s *StripeAdapter) GetPaymentStatus(paymentID string) (string, error) {
    // 调用Stripe API
    paymentIntent, err := s.stripe.RetrievePaymentIntent(paymentID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve payment intent: %w", err)
    }
    
    // 转换状态
    switch paymentIntent.Status {
    case "succeeded":
        return "completed", nil
    case "processing":
        return "pending", nil
    case "requires_payment_method":
        return "failed", nil
    default:
        return "unknown", nil
    }
}

```

### 3.2.1.5.2 数据库适配器

```go
package databaseadapter

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
)

// DatabaseInterface 数据库接口
type DatabaseInterface interface {
    Connect(dsn string) error
    Query(sql string, args ...interface{}) ([]map[string]interface{}, error)
    Execute(sql string, args ...interface{}) (int64, error)
    Close() error
}

// MySQLDatabase MySQL数据库
type MySQLDatabase struct {
    db *sql.DB
}

func (m *MySQLDatabase) Connect(dsn string) error {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    m.db = db
    return nil
}

func (m *MySQLDatabase) Query(sql string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := m.db.Query(sql, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var results []map[string]interface{}
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        results = append(results, row)
    }
    
    return results, nil
}

func (m *MySQLDatabase) Execute(sql string, args ...interface{}) (int64, error) {
    result, err := m.db.Exec(sql, args...)
    if err != nil {
        return 0, err
    }
    
    return result.RowsAffected()
}

func (m *MySQLDatabase) Close() error {
    return m.db.Close()
}

// PostgreSQLDatabase PostgreSQL数据库
type PostgreSQLDatabase struct {
    db *sql.DB
}

func (p *PostgreSQLDatabase) Connect(dsn string) error {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }
    p.db = db
    return nil
}

func (p *PostgreSQLDatabase) Query(sql string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := p.db.Query(sql, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var results []map[string]interface{}
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        results = append(results, row)
    }
    
    return results, nil
}

func (p *PostgreSQLDatabase) Execute(sql string, args ...interface{}) (int64, error) {
    result, err := p.db.Exec(sql, args...)
    if err != nil {
        return 0, err
    }
    
    return result.RowsAffected()
}

func (p *PostgreSQLDatabase) Close() error {
    return p.db.Close()
}

// DatabaseFactory 数据库工厂
type DatabaseFactory struct{}

// CreateDatabase 创建数据库实例
func (df *DatabaseFactory) CreateDatabase(dbType string) (DatabaseInterface, error) {
    switch dbType {
    case "mysql":
        return &MySQLDatabase{}, nil
    case "postgres":
        return &PostgreSQLDatabase{}, nil
    default:
        return nil, fmt.Errorf("unsupported database type: %s", dbType)
    }
}

```

### 3.2.1.5.3 日志适配器

```go
package logadapter

import (
    "fmt"
    "log"
    "os"
    "time"
)

// LoggerInterface 日志接口
type LoggerInterface interface {
    Info(message string)
    Error(message string)
    Debug(message string)
    Warn(message string)
}

// SimpleLogger 简单日志器
type SimpleLogger struct {
    logger *log.Logger
}

func NewSimpleLogger() *SimpleLogger {
    return &SimpleLogger{
        logger: log.New(os.Stdout, "", log.LstdFlags),
    }
}

func (s *SimpleLogger) Info(message string) {
    s.logger.Printf("[INFO] %s", message)
}

func (s *SimpleLogger) Error(message string) {
    s.logger.Printf("[ERROR] %s", message)
}

func (s *SimpleLogger) Debug(message string) {
    s.logger.Printf("[DEBUG] %s", message)
}

func (s *SimpleLogger) Warn(message string) {
    s.logger.Printf("[WARN] %s", message)
}

// StructuredLogger 结构化日志器
type StructuredLogger struct {
    logger *log.Logger
}

func NewStructuredLogger() *StructuredLogger {
    return &StructuredLogger{
        logger: log.New(os.Stdout, "", log.LstdFlags),
    }
}

func (s *StructuredLogger) Log(level, message string) {
    timestamp := time.Now().Format(time.RFC3339)
    s.logger.Printf(`{"timestamp":"%s","level":"%s","message":"%s"}`, timestamp, level, message)
}

// StructuredLoggerAdapter 结构化日志适配器
type StructuredLoggerAdapter struct {
    logger *StructuredLogger
}

func NewStructuredLoggerAdapter() *StructuredLoggerAdapter {
    return &StructuredLoggerAdapter{
        logger: NewStructuredLogger(),
    }
}

func (s *StructuredLoggerAdapter) Info(message string) {
    s.logger.Log("INFO", message)
}

func (s *StructuredLoggerAdapter) Error(message string) {
    s.logger.Log("ERROR", message)
}

func (s *StructuredLoggerAdapter) Debug(message string) {
    s.logger.Log("DEBUG", message)
}

func (s *StructuredLoggerAdapter) Warn(message string) {
    s.logger.Log("WARN", message)
}

```

## 3.2.1.6 5. 批判性分析

### 3.2.1.6.1 优势

1. **接口兼容**: 解决接口不兼容问题
2. **代码复用**: 复用现有类而不修改
3. **透明性**: 客户端无需知道适配器存在
4. **扩展性**: 支持多种适配方式

### 3.2.1.6.2 劣势

1. **复杂性**: 增加系统复杂性
2. **性能开销**: 额外的适配层开销
3. **维护困难**: 适配器代码需要维护
4. **过度使用**: 可能导致过度设计

### 3.2.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 结构体 | 高 | 低 |
| Java | 类继承 + 接口 | 中 | 中 |
| C++ | 多重继承 | 中 | 中 |
| Python | 鸭子类型 | 高 | 低 |

### 3.2.1.6.4 最新趋势

1. **函数式适配器**: 使用函数作为适配器
2. **泛型适配器**: 利用Go 1.18+泛型
3. **自动适配**: 代码生成适配器
4. **配置驱动**: 基于配置的适配

## 3.2.1.7 6. 面试题与考点

### 3.2.1.7.1 基础考点

1. **Q**: 适配器模式与装饰器模式的区别？
   **A**: 适配器改变接口，装饰器增强功能

2. **Q**: 什么时候使用适配器模式？
   **A**: 需要集成第三方库、接口不兼容时

3. **Q**: 类适配器与对象适配器的区别？
   **A**: 类适配器使用继承，对象适配器使用组合

### 3.2.1.7.2 进阶考点

1. **Q**: 如何设计通用的适配器？
   **A**: 使用接口、泛型、配置驱动

2. **Q**: 适配器模式在微服务中的应用？
   **A**: 服务间通信、API版本兼容

3. **Q**: 如何处理适配器的性能问题？
   **A**: 缓存、连接池、异步处理

## 3.2.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 适配器模式 | 转换接口的设计模式 | Adapter Pattern |
| 目标接口 | 客户端期望的接口 | Target Interface |
| 适配者 | 需要适配的类 | Adaptee |
| 适配器 | 转换接口的类 | Adapter |
| 接口映射 | 接口间的转换关系 | Interface Mapping |

## 3.2.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 过度适配 | 适配器过多 | 评估实际需求 |
| 性能问题 | 适配层开销大 | 优化适配逻辑 |
| 维护困难 | 适配器代码复杂 | 简化适配逻辑 |
| 接口污染 | 目标接口过于复杂 | 设计简洁接口 |

## 3.2.1.10 9. 相关主题

- [装饰器模式](./02-Decorator-Pattern.md)
- [代理模式](./03-Proxy-Pattern.md)
- [外观模式](./04-Facade-Pattern.md)
- [桥接模式](./05-Bridge-Pattern.md)
- [组合模式](./06-Composite-Pattern.md)

## 3.2.1.11 10. 学习路径

### 3.2.1.11.1 新手路径

1. 理解适配器模式的基本概念
2. 学习类适配器和对象适配器
3. 实现简单的适配器
4. 理解接口兼容的重要性

### 3.2.1.11.2 进阶路径

1. 学习函数式适配器
2. 理解适配器的性能优化
3. 掌握适配器的应用场景
4. 学习适配器的最佳实践

### 3.2.1.11.3 高阶路径

1. 分析适配器在大型项目中的应用
2. 理解适配器与架构设计的关系
3. 掌握适配器的性能调优
4. 学习适配器的替代方案

---

**相关文档**: [结构型模式总览](./README.md) | [设计模式总览](../README.md)
