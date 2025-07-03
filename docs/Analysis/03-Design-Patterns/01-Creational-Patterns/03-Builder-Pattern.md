# 建造者模式 (Builder Pattern)

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

建造者模式将一个复杂对象的构建与它的表示分离，使得同样的构建过程可以创建不同的表示。

**形式化定义**:
$$Builder = (Director, Builder, Product, ConcreteBuilder, BuildSteps)$$

其中：

- $Director$ 是指导者，控制构建过程
- $Builder$ 是建造者接口
- $Product$ 是产品
- $ConcreteBuilder$ 是具体建造者
- $BuildSteps$ 是构建步骤集合

### 1.2 核心特征

- **分步构建**: 将复杂对象构建过程分解为多个步骤
- **链式调用**: 支持方法链式调用
- **参数验证**: 在构建过程中进行参数验证
- **不可变对象**: 构建完成后对象不可变

## 2. 理论基础

### 2.1 数学形式化

**定义 2.1** (建造者模式): 建造者模式是一个六元组 $B = (D, B, P, S, M, V)$

其中：

- $D$ 是指导者集合
- $B$ 是建造者集合
- $P$ 是产品集合
- $S$ 是构建步骤集合
- $M$ 是构建方法映射，$M: B \times S \rightarrow B$
- $V$ 是验证规则

**定理 2.1** (构建完整性): 对于任意建造者 $b \in B$，存在构建序列 $s_1, s_2, ..., s_n$ 使得 $M(...M(M(b, s_1), s_2)..., s_n) \in P$

**证明**: 由建造者模式的完整性保证。

### 2.2 范畴论视角

在范畴论中，建造者模式可以表示为：

$$Builder : Builder \times Steps \rightarrow Builder$$

其中 $Builder$ 和 $Steps$ 是对象范畴。

## 3. Go语言实现

### 3.1 基础建造者模式

```go
package builder

import "fmt"

// Product 产品
type Product struct {
    PartA string
    PartB string
    PartC string
}

func (p *Product) Show() string {
    return fmt.Sprintf("Product: %s, %s, %s", p.PartA, p.PartB, p.PartC)
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
    return &ConcreteBuilder{
        product: &Product{},
    }
}

func (b *ConcreteBuilder) BuildPartA() {
    b.product.PartA = "PartA"
}

func (b *ConcreteBuilder) BuildPartB() {
    b.product.PartB = "PartB"
}

func (b *ConcreteBuilder) BuildPartC() {
    b.product.PartC = "PartC"
}

func (b *ConcreteBuilder) GetResult() *Product {
    return b.product
}

// Director 指导者
type Director struct {
    builder Builder
}

func NewDirector(builder Builder) *Director {
    return &Director{builder: builder}
}

func (d *Director) Construct() *Product {
    d.builder.BuildPartA()
    d.builder.BuildPartB()
    d.builder.BuildPartC()
    return d.builder.GetResult()
}
```

### 3.2 链式建造者模式

```go
package chainbuilder

import "fmt"

// Product 产品
type Product struct {
    Name     string
    Price    float64
    Category string
    Tags     []string
}

// ProductBuilder 产品建造者
type ProductBuilder struct {
    product *Product
}

func NewProductBuilder() *ProductBuilder {
    return &ProductBuilder{
        product: &Product{},
    }
}

func (b *ProductBuilder) SetName(name string) *ProductBuilder {
    b.product.Name = name
    return b
}

func (b *ProductBuilder) SetPrice(price float64) *ProductBuilder {
    if price < 0 {
        panic("price cannot be negative")
    }
    b.product.Price = price
    return b
}

func (b *ProductBuilder) SetCategory(category string) *ProductBuilder {
    b.product.Category = category
    return b
}

func (b *ProductBuilder) AddTag(tag string) *ProductBuilder {
    b.product.Tags = append(b.product.Tags, tag)
    return b
}

func (b *ProductBuilder) Build() *Product {
    if b.product.Name == "" {
        panic("product name is required")
    }
    if b.product.Price == 0 {
        panic("product price is required")
    }
    return b.product
}

// 使用示例
func Example() {
    product := NewProductBuilder().
        SetName("Laptop").
        SetPrice(999.99).
        SetCategory("Electronics").
        AddTag("portable").
        AddTag("computer").
        Build()
    
    fmt.Printf("Product: %+v\n", product)
}
```

### 3.3 函数式建造者模式

```go
package functionalbuilder

import "fmt"

// Product 产品
type Product struct {
    Name     string
    Price    float64
    Category string
    Tags     []string
}

// ProductOption 产品选项函数
type ProductOption func(*Product)

// WithName 设置名称
func WithName(name string) ProductOption {
    return func(p *Product) {
        p.Name = name
    }
}

// WithPrice 设置价格
func WithPrice(price float64) ProductOption {
    return func(p *Product) {
        if price >= 0 {
            p.Price = price
        }
    }
}

// WithCategory 设置分类
func WithCategory(category string) ProductOption {
    return func(p *Product) {
        p.Category = category
    }
}

// WithTags 设置标签
func WithTags(tags ...string) ProductOption {
    return func(p *Product) {
        p.Tags = append(p.Tags, tags...)
    }
}

// NewProduct 创建产品
func NewProduct(options ...ProductOption) *Product {
    product := &Product{}
    
    for _, option := range options {
        option(product)
    }
    
    return product
}

// 使用示例
func Example() {
    product := NewProduct(
        WithName("Smartphone"),
        WithPrice(599.99),
        WithCategory("Electronics"),
        WithTags("mobile", "communication"),
    )
    
    fmt.Printf("Product: %+v\n", product)
}
```

## 4. 工程案例

### 4.1 HTTP请求建造者

```go
package httpbuilder

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// HTTPRequest HTTP请求
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
    Timeout time.Duration
}

// HTTPRequestBuilder HTTP请求建造者
type HTTPRequestBuilder struct {
    request *HTTPRequest
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
    return &HTTPRequestBuilder{
        request: &HTTPRequest{
            Headers: make(map[string]string),
            Timeout: 30 * time.Second,
        },
    }
}

func (b *HTTPRequestBuilder) SetMethod(method string) *HTTPRequestBuilder {
    b.request.Method = method
    return b
}

func (b *HTTPRequestBuilder) SetURL(url string) *HTTPRequestBuilder {
    b.request.URL = url
    return b
}

func (b *HTTPRequestBuilder) SetHeader(key, value string) *HTTPRequestBuilder {
    b.request.Headers[key] = value
    return b
}

func (b *HTTPRequestBuilder) SetJSONBody(data interface{}) *HTTPRequestBuilder {
    jsonData, err := json.Marshal(data)
    if err != nil {
        panic(fmt.Sprintf("failed to marshal JSON: %v", err))
    }
    b.request.Body = jsonData
    b.SetHeader("Content-Type", "application/json")
    return b
}

func (b *HTTPRequestBuilder) SetTimeout(timeout time.Duration) *HTTPRequestBuilder {
    b.request.Timeout = timeout
    return b
}

func (b *HTTPRequestBuilder) Build() *HTTPRequest {
    if b.request.Method == "" {
        b.request.Method = "GET"
    }
    if b.request.URL == "" {
        panic("URL is required")
    }
    return b.request
}

// Execute 执行请求
func (r *HTTPRequest) Execute() (*http.Response, error) {
    req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
    if err != nil {
        return nil, err
    }
    
    for key, value := range r.Headers {
        req.Header.Set(key, value)
    }
    
    client := &http.Client{Timeout: r.Timeout}
    return client.Do(req)
}
```

### 4.2 数据库查询建造者

```go
package querybuilder

import (
    "fmt"
    "strings"
)

// Query 查询对象
type Query struct {
    Table     string
    Fields    []string
    Where     []string
    OrderBy   []string
    Limit     int
    Offset    int
    SQL       string
}

// QueryBuilder 查询建造者
type QueryBuilder struct {
    query *Query
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        query: &Query{
            Fields:  []string{"*"},
            Where:   []string{},
            OrderBy: []string{},
        },
    }
}

func (b *QueryBuilder) From(table string) *QueryBuilder {
    b.query.Table = table
    return b
}

func (b *QueryBuilder) Select(fields ...string) *QueryBuilder {
    if len(fields) > 0 {
        b.query.Fields = fields
    }
    return b
}

func (b *QueryBuilder) Where(condition string) *QueryBuilder {
    b.query.Where = append(b.query.Where, condition)
    return b
}

func (b *QueryBuilder) OrderBy(field string, direction string) *QueryBuilder {
    if direction == "" {
        direction = "ASC"
    }
    b.query.OrderBy = append(b.query.OrderBy, fmt.Sprintf("%s %s", field, direction))
    return b
}

func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
    b.query.Limit = limit
    return b
}

func (b *QueryBuilder) Offset(offset int) *QueryBuilder {
    b.query.Offset = offset
    return b
}

func (b *QueryBuilder) Build() *Query {
    if b.query.Table == "" {
        panic("table name is required")
    }
    
    // 构建SQL
    sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(b.query.Fields, ", "), b.query.Table)
    
    if len(b.query.Where) > 0 {
        sql += " WHERE " + strings.Join(b.query.Where, " AND ")
    }
    
    if len(b.query.OrderBy) > 0 {
        sql += " ORDER BY " + strings.Join(b.query.OrderBy, ", ")
    }
    
    if b.query.Limit > 0 {
        sql += fmt.Sprintf(" LIMIT %d", b.query.Limit)
    }
    
    if b.query.Offset > 0 {
        sql += fmt.Sprintf(" OFFSET %d", b.query.Offset)
    }
    
    b.query.SQL = sql
    return b.query
}
```

## 5. 批判性分析

### 5.1 优势

1. **分步构建**: 复杂对象构建过程清晰
2. **参数验证**: 构建过程中进行验证
3. **不可变性**: 构建完成后对象不可变
4. **链式调用**: 提供流畅的API

### 5.2 劣势

1. **复杂性**: 增加系统复杂性
2. **性能开销**: 额外的构建步骤开销
3. **过度设计**: 简单对象可能过度设计
4. **学习成本**: 需要理解构建流程

### 5.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 链式方法 | 高 | 中 |
| Java | Builder类 | 中 | 中 |
| C++ | 流式接口 | 中 | 中 |
| Python | 字典参数 | 高 | 低 |

### 5.4 最新趋势

1. **函数式建造者**: 使用函数选项模式
2. **类型安全**: 利用Go的类型系统
3. **验证集成**: 内置参数验证
4. **代码生成**: 自动生成建造者代码

## 6. 面试题与考点

### 6.1 基础考点

1. **Q**: 建造者模式与工厂模式的区别？
   **A**: 建造者关注构建过程，工厂关注创建结果

2. **Q**: 什么时候使用建造者模式？
   **A**: 对象构建复杂、需要参数验证、需要不可变对象时

3. **Q**: 链式调用的优势是什么？
   **A**: 提供流畅的API、减少临时变量、提高可读性

### 6.2 进阶考点

1. **Q**: 如何实现参数验证？
   **A**: 在Build方法中进行验证，使用panic或返回error

2. **Q**: 建造者模式与函数式编程的关系？
   **A**: 函数选项模式是函数式建造者的体现

3. **Q**: 如何保证对象不可变性？
   **A**: 构建完成后返回副本或使用只读接口

## 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 建造者模式 | 分步构建复杂对象的设计模式 | Builder Pattern |
| 指导者 | 控制构建过程的对象 | Director |
| 建造者 | 负责构建产品的对象 | Builder |
| 链式调用 | 方法连续调用的编程风格 | Method Chaining |
| 函数选项 | 使用函数作为选项的模式 | Functional Options |

## 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 过度复杂 | 简单对象使用建造者 | 评估实际需求 |
| 参数验证缺失 | 构建无效对象 | 在Build方法中验证 |
| 性能问题 | 过多中间对象 | 使用对象池或复用 |
| API设计不当 | 链式调用不流畅 | 精心设计API |

## 9. 相关主题

- [单例模式](./01-Singleton-Pattern.md)
- [工厂模式](./02-Factory-Pattern.md)
- [原型模式](./04-Prototype-Pattern.md)
- [抽象工厂模式](./05-Abstract-Factory-Pattern.md)
- [函数式编程模式](../03-Behavioral-Patterns/01-Functional-Patterns.md)

## 10. 学习路径

### 10.1 新手路径

1. 理解建造者模式的基本概念
2. 学习基础的建造者实现
3. 实现简单的链式调用
4. 理解参数验证的重要性

### 10.2 进阶路径

1. 学习函数式建造者模式
2. 理解不可变对象的设计
3. 掌握复杂对象的构建流程
4. 学习建造者模式的最佳实践

### 10.3 高阶路径

1. 分析建造者模式在大型项目中的应用
2. 理解建造者模式与架构设计的关系
3. 掌握建造者模式的性能优化
4. 学习建造者模式的替代方案

---

**相关文档**: [创建型模式总览](./README.md) | [设计模式总览](../README.md)
