# 3.3.1 模板方法模式 (Template Method Pattern)

## 3.3.1.1 目录

## 3.3.1.2 1. 概述

### 3.3.1.2.1 定义

模板方法模式定义了一个算法的骨架，将一些步骤延迟到子类中实现。模板方法使得子类可以在不改变算法结构的情况下，重新定义算法的某些特定步骤。

**形式化定义**:
$$TemplateMethod = (AbstractClass, ConcreteClass_1, ConcreteClass_2, ..., ConcreteClass_n)$$

其中：

- $AbstractClass$ 是抽象类
- $ConcreteClass_i$ 是具体实现类

### 3.3.1.2.2 核心特征

- **算法骨架**: 定义算法的基本结构
- **步骤抽象**: 将可变步骤抽象化
- **子类实现**: 子类实现具体步骤
- **结构稳定**: 算法结构不可改变

## 3.3.1.3 2. 理论基础

### 3.3.1.3.1 数学形式化

**定义 2.1** (模板方法模式): 模板方法模式是一个三元组 $T = (A, S, I)$

其中：

- $A$ 是算法骨架
- $S$ 是步骤集合
- $I$ 是实现函数，$I: S \times ConcreteClass \rightarrow Result$

**定理 2.1** (算法稳定性): 对于任意具体类 $c \in ConcreteClass$，算法骨架 $A$ 保持不变。

### 3.3.1.3.2 范畴论视角

在范畴论中，模板方法模式可以表示为：

$$TemplateMethod : Algorithm \times Implementation \rightarrow Result$$

## 3.3.1.4 3. Go语言实现

### 3.3.1.4.1 基础模板方法模式

```go
package template

import "fmt"

// AbstractClass 抽象类
type AbstractClass interface {
    TemplateMethod()
    PrimitiveOperation1()
    PrimitiveOperation2()
    Hook() bool
}

// BaseClass 基础类
type BaseClass struct{}

func (b *BaseClass) TemplateMethod() {
    fmt.Println("Template method started")
    b.PrimitiveOperation1()
    
    if b.Hook() {
        b.PrimitiveOperation2()
    }
    
    fmt.Println("Template method finished")
}

func (b *BaseClass) PrimitiveOperation1() {
    fmt.Println("Default implementation of operation 1")
}

func (b *BaseClass) PrimitiveOperation2() {
    fmt.Println("Default implementation of operation 2")
}

func (b *BaseClass) Hook() bool {
    return true
}

// ConcreteClassA 具体类A
type ConcreteClassA struct {
    BaseClass
}

func NewConcreteClassA() *ConcreteClassA {
    return &ConcreteClassA{}
}

func (c *ConcreteClassA) PrimitiveOperation1() {
    fmt.Println("ConcreteClassA: Operation 1")
}

func (c *ConcreteClassA) PrimitiveOperation2() {
    fmt.Println("ConcreteClassA: Operation 2")
}

func (c *ConcreteClassA) Hook() bool {
    return true
}

// ConcreteClassB 具体类B
type ConcreteClassB struct {
    BaseClass
}

func NewConcreteClassB() *ConcreteClassB {
    return &ConcreteClassB{}
}

func (c *ConcreteClassB) PrimitiveOperation1() {
    fmt.Println("ConcreteClassB: Operation 1")
}

func (c *ConcreteClassB) PrimitiveOperation2() {
    fmt.Println("ConcreteClassB: Operation 2")
}

func (c *ConcreteClassB) Hook() bool {
    return false
}

```

### 3.3.1.4.2 数据处理器模板方法

```go
package dataprocessor

import (
    "fmt"
    "time"
)

// DataProcessor 数据处理器接口
type DataProcessor interface {
    Process()
    LoadData() []interface{}
    TransformData(data []interface{}) []interface{}
    SaveData(data []interface{})
    ValidateData(data []interface{}) bool
    LogProcess(start, end time.Time, recordCount int)
}

// BaseDataProcessor 基础数据处理器
type BaseDataProcessor struct{}

func (b *BaseDataProcessor) Process() {
    start := time.Now()
    
    fmt.Println("Starting data processing...")
    
    // 加载数据
    data := b.LoadData()
    fmt.Printf("Loaded %d records\n", len(data))
    
    // 验证数据
    if !b.ValidateData(data) {
        fmt.Println("Data validation failed")
        return
    }
    
    // 转换数据
    transformedData := b.TransformData(data)
    fmt.Printf("Transformed %d records\n", len(transformedData))
    
    // 保存数据
    b.SaveData(transformedData)
    fmt.Println("Data saved successfully")
    
    end := time.Now()
    b.LogProcess(start, end, len(transformedData))
}

func (b *BaseDataProcessor) LoadData() []interface{} {
    fmt.Println("Default: Loading data")
    return make([]interface{}, 0)
}

func (b *BaseDataProcessor) TransformData(data []interface{}) []interface{} {
    fmt.Println("Default: Transforming data")
    return data
}

func (b *BaseDataProcessor) SaveData(data []interface{}) {
    fmt.Println("Default: Saving data")
}

func (b *BaseDataProcessor) ValidateData(data []interface{}) bool {
    fmt.Println("Default: Validating data")
    return true
}

func (b *BaseDataProcessor) LogProcess(start, end time.Time, recordCount int) {
    duration := end.Sub(start)
    fmt.Printf("Process completed in %v, processed %d records\n", duration, recordCount)
}

// CSVProcessor CSV数据处理器
type CSVProcessor struct {
    BaseDataProcessor
    filePath string
}

func NewCSVProcessor(filePath string) *CSVProcessor {
    return &CSVProcessor{
        filePath: filePath,
    }
}

func (c *CSVProcessor) LoadData() []interface{} {
    fmt.Printf("Loading CSV data from %s\n", c.filePath)
    // 模拟加载CSV数据
    return []interface{}{
        map[string]string{"name": "John", "age": "30"},
        map[string]string{"name": "Jane", "age": "25"},
        map[string]string{"name": "Bob", "age": "35"},
    }
}

func (c *CSVProcessor) TransformData(data []interface{}) []interface{} {
    fmt.Println("Transforming CSV data")
    transformed := make([]interface{}, 0)
    
    for _, item := range data {
        if record, ok := item.(map[string]string); ok {
            // 转换年龄为整数
            transformed = append(transformed, map[string]interface{}{
                "name": record["name"],
                "age":  30, // 简化处理
            })
        }
    }
    
    return transformed
}

func (c *CSVProcessor) SaveData(data []interface{}) {
    fmt.Printf("Saving %d records to database\n", len(data))
    // 模拟保存到数据库
}

func (c *CSVProcessor) ValidateData(data []interface{}) bool {
    fmt.Println("Validating CSV data")
    return len(data) > 0
}

// JSONProcessor JSON数据处理器
type JSONProcessor struct {
    BaseDataProcessor
    filePath string
}

func NewJSONProcessor(filePath string) *JSONProcessor {
    return &JSONProcessor{
        filePath: filePath,
    }
}

func (j *JSONProcessor) LoadData() []interface{} {
    fmt.Printf("Loading JSON data from %s\n", j.filePath)
    // 模拟加载JSON数据
    return []interface{}{
        map[string]interface{}{"id": 1, "name": "Product1", "price": 100.0},
        map[string]interface{}{"id": 2, "name": "Product2", "price": 200.0},
        map[string]interface{}{"id": 3, "name": "Product3", "price": 150.0},
    }
}

func (j *JSONProcessor) TransformData(data []interface{}) []interface{} {
    fmt.Println("Transforming JSON data")
    transformed := make([]interface{}, 0)
    
    for _, item := range data {
        if record, ok := item.(map[string]interface{}); ok {
            // 计算折扣价格
            price := record["price"].(float64)
            discountedPrice := price * 0.9 // 10% 折扣
            transformed = append(transformed, map[string]interface{}{
                "id":    record["id"],
                "name":  record["name"],
                "price": discountedPrice,
            })
        }
    }
    
    return transformed
}

func (j *JSONProcessor) SaveData(data []interface{}) {
    fmt.Printf("Saving %d records to cache\n", len(data))
    // 模拟保存到缓存
}

func (j *JSONProcessor) ValidateData(data []interface{}) bool {
    fmt.Println("Validating JSON data")
    return len(data) > 0
}

```

### 3.3.1.4.3 构建器模板方法

```go
package builder

import "fmt"

// BuildProcess 构建过程接口
type BuildProcess interface {
    Build()
    SetWalls()
    SetRoof()
    SetWindows()
    SetDoors()
    SetFoundation()
    Validate() bool
}

// House 房屋
type House struct {
    Walls      string
    Roof       string
    Windows    string
    Doors      string
    Foundation string
}

func (h *House) String() string {
    return fmt.Sprintf("House with %s foundation, %s walls, %s roof, %s windows, %s doors",
        h.Foundation, h.Walls, h.Roof, h.Windows, h.Doors)
}

// BaseBuilder 基础构建器
type BaseBuilder struct {
    house *House
}

func NewBaseBuilder() *BaseBuilder {
    return &BaseBuilder{
        house: &House{},
    }
}

func (b *BaseBuilder) Build() {
    fmt.Println("Starting house construction...")
    
    b.SetFoundation()
    b.SetWalls()
    b.SetRoof()
    b.SetWindows()
    b.SetDoors()
    
    if b.Validate() {
        fmt.Printf("House built successfully: %s\n", b.house)
    } else {
        fmt.Println("House construction failed validation")
    }
}

func (b *BaseBuilder) SetFoundation() {
    fmt.Println("Setting default foundation")
    b.house.Foundation = "Concrete"
}

func (b *BaseBuilder) SetWalls() {
    fmt.Println("Setting default walls")
    b.house.Walls = "Brick"
}

func (b *BaseBuilder) SetRoof() {
    fmt.Println("Setting default roof")
    b.house.Roof = "Shingle"
}

func (b *BaseBuilder) SetWindows() {
    fmt.Println("Setting default windows")
    b.house.Windows = "Standard"
}

func (b *BaseBuilder) SetDoors() {
    fmt.Println("Setting default doors")
    b.house.Doors = "Wood"
}

func (b *BaseBuilder) Validate() bool {
    return b.house.Foundation != "" && b.house.Walls != "" && 
           b.house.Roof != "" && b.house.Windows != "" && b.house.Doors != ""
}

func (b *BaseBuilder) GetHouse() *House {
    return b.house
}

// WoodenHouseBuilder 木屋构建器
type WoodenHouseBuilder struct {
    BaseBuilder
}

func NewWoodenHouseBuilder() *WoodenHouseBuilder {
    return &WoodenHouseBuilder{
        BaseBuilder: *NewBaseBuilder(),
    }
}

func (w *WoodenHouseBuilder) SetFoundation() {
    fmt.Println("Setting wooden foundation")
    w.house.Foundation = "Wooden"
}

func (w *WoodenHouseBuilder) SetWalls() {
    fmt.Println("Setting wooden walls")
    w.house.Walls = "Wood"
}

func (w *WoodenHouseBuilder) SetRoof() {
    fmt.Println("Setting wooden roof")
    w.house.Roof = "Wood"
}

func (w *WoodenHouseBuilder) SetWindows() {
    fmt.Println("Setting wooden windows")
    w.house.Windows = "Wooden"
}

func (w *WoodenHouseBuilder) SetDoors() {
    fmt.Println("Setting wooden doors")
    w.house.Doors = "Wooden"
}

// StoneHouseBuilder 石屋构建器
type StoneHouseBuilder struct {
    BaseBuilder
}

func NewStoneHouseBuilder() *StoneHouseBuilder {
    return &StoneHouseBuilder{
        BaseBuilder: *NewBaseBuilder(),
    }
}

func (s *StoneHouseBuilder) SetFoundation() {
    fmt.Println("Setting stone foundation")
    s.house.Foundation = "Stone"
}

func (s *StoneHouseBuilder) SetWalls() {
    fmt.Println("Setting stone walls")
    s.house.Walls = "Stone"
}

func (s *StoneHouseBuilder) SetRoof() {
    fmt.Println("Setting stone roof")
    s.house.Roof = "Stone"
}

func (s *StoneHouseBuilder) SetWindows() {
    fmt.Println("Setting stone windows")
    s.house.Windows = "Stone"
}

func (s *StoneHouseBuilder) SetDoors() {
    fmt.Println("Setting stone doors")
    s.house.Doors = "Stone"
}

```

## 3.3.1.5 4. 工程案例

### 3.3.1.5.1 HTTP请求处理器模板方法

```go
package httpprocessor

import (
    "fmt"
    "net/http"
    "time"
)

// RequestProcessor HTTP请求处理器接口
type RequestProcessor interface {
    Process()
    Authenticate() bool
    Validate() bool
    ProcessRequest() interface{}
    LogRequest(start, end time.Time)
}

// BaseRequestProcessor 基础请求处理器
type BaseRequestProcessor struct {
    request  *http.Request
    response http.ResponseWriter
}

func NewBaseRequestProcessor(w http.ResponseWriter, r *http.Request) *BaseRequestProcessor {
    return &BaseRequestProcessor{
        request:  r,
        response: w,
    }
}

func (b *BaseRequestProcessor) Process() {
    start := time.Now()
    
    fmt.Printf("Processing %s request to %s\n", b.request.Method, b.request.URL.Path)
    
    // 认证
    if !b.Authenticate() {
        http.Error(b.response, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // 验证
    if !b.Validate() {
        http.Error(b.response, "Bad Request", http.StatusBadRequest)
        return
    }
    
    // 处理请求
    result := b.ProcessRequest()
    
    // 返回结果
    b.response.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(b.response, "%v", result)
    
    end := time.Now()
    b.LogRequest(start, end)
}

func (b *BaseRequestProcessor) Authenticate() bool {
    fmt.Println("Default: Basic authentication")
    return true
}

func (b *BaseRequestProcessor) Validate() bool {
    fmt.Println("Default: Basic validation")
    return true
}

func (b *BaseRequestProcessor) ProcessRequest() interface{} {
    fmt.Println("Default: Processing request")
    return map[string]string{"message": "Default response"}
}

func (b *BaseRequestProcessor) LogRequest(start, end time.Time) {
    duration := end.Sub(start)
    fmt.Printf("Request processed in %v\n", duration)
}

// UserRequestProcessor 用户请求处理器
type UserRequestProcessor struct {
    BaseRequestProcessor
}

func NewUserRequestProcessor(w http.ResponseWriter, r *http.Request) *UserRequestProcessor {
    return &UserRequestProcessor{
        BaseRequestProcessor: *NewBaseRequestProcessor(w, r),
    }
}

func (u *UserRequestProcessor) Authenticate() bool {
    fmt.Println("User: Token-based authentication")
    // 检查Authorization头
    auth := u.request.Header.Get("Authorization")
    return auth != ""
}

func (u *UserRequestProcessor) Validate() bool {
    fmt.Println("User: Validating user request")
    return u.request.Method == "GET" || u.request.Method == "POST"
}

func (u *UserRequestProcessor) ProcessRequest() interface{} {
    fmt.Println("User: Processing user request")
    return map[string]interface{}{
        "user_id": 123,
        "name":    "John Doe",
        "email":   "john@example.com",
    }
}

// AdminRequestProcessor 管理员请求处理器
type AdminRequestProcessor struct {
    BaseRequestProcessor
}

func NewAdminRequestProcessor(w http.ResponseWriter, r *http.Request) *AdminRequestProcessor {
    return &AdminRequestProcessor{
        BaseRequestProcessor: *NewBaseRequestProcessor(w, r),
    }
}

func (a *AdminRequestProcessor) Authenticate() bool {
    fmt.Println("Admin: Admin authentication")
    // 检查管理员权限
    auth := a.request.Header.Get("Authorization")
    return auth == "admin-token"
}

func (a *AdminRequestProcessor) Validate() bool {
    fmt.Println("Admin: Validating admin request")
    return a.request.Method == "GET" || a.request.Method == "POST" || 
           a.request.Method == "PUT" || a.request.Method == "DELETE"
}

func (a *AdminRequestProcessor) ProcessRequest() interface{} {
    fmt.Println("Admin: Processing admin request")
    return map[string]interface{}{
        "admin_id": 1,
        "role":     "super_admin",
        "permissions": []string{"read", "write", "delete"},
    }
}

```

## 3.3.1.6 5. 批判性分析

### 3.3.1.6.1 优势

1. **算法复用**: 算法骨架可以复用
2. **扩展性**: 子类可以扩展算法步骤
3. **结构稳定**: 算法结构不可改变
4. **代码复用**: 减少重复代码

### 3.3.1.6.2 劣势

1. **继承限制**: 只能通过继承扩展
2. **步骤固定**: 算法步骤顺序固定
3. **违反开闭**: 修改算法需要修改基类
4. **理解困难**: 算法流程可能难以理解

### 3.3.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口+嵌入 | 高 | 中 |
| Java | 抽象类 | 中 | 中 |
| C++ | 虚函数 | 高 | 中 |
| Python | 抽象基类 | 中 | 低 |

### 3.3.1.6.4 最新趋势

1. **函数式模板**: 使用函数作为模板
2. **组合模式**: 使用组合替代继承
3. **策略模板**: 结合策略模式
4. **微服务模板**: 服务模板方法

## 3.3.1.7 6. 面试题与考点

### 3.3.1.7.1 基础考点

1. **Q**: 模板方法模式与策略模式的区别？
   **A**: 模板方法关注算法结构，策略模式关注算法选择

2. **Q**: 什么时候使用模板方法模式？
   **A**: 算法步骤固定、需要复用算法骨架时

3. **Q**: 模板方法模式的优缺点？
   **A**: 优点：算法复用、结构稳定；缺点：继承限制、步骤固定

### 3.3.1.7.2 进阶考点

1. **Q**: 如何避免继承的限制？
   **A**: 使用组合模式、函数式编程

2. **Q**: 模板方法模式在微服务中的应用？
   **A**: 服务模板、API模板、部署模板

3. **Q**: 如何处理算法的变化？
   **A**: 钩子方法、策略模式、组合模式

## 3.3.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 模板方法模式 | 定义算法骨架的设计模式 | Template Method Pattern |
| 抽象类 | 定义算法骨架的类 | Abstract Class |
| 具体类 | 实现具体步骤的类 | Concrete Class |
| 钩子方法 | 可选的算法步骤 | Hook Method |

## 3.3.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 继承限制 | 只能通过继承扩展 | 使用组合模式 |
| 步骤固定 | 算法步骤顺序固定 | 使用策略模式 |
| 违反开闭 | 修改算法需要修改基类 | 使用钩子方法 |
| 理解困难 | 算法流程难以理解 | 文档化、简化设计 |

## 3.3.1.10 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)

## 3.3.1.11 10. 学习路径

### 3.3.1.11.1 新手路径

1. 理解模板方法模式的基本概念
2. 学习算法骨架的定义
3. 实现简单的模板方法模式
4. 理解钩子方法的作用

### 3.3.1.11.2 进阶路径

1. 学习复杂的模板方法实现
2. 理解模板方法的性能优化
3. 掌握模板方法的应用场景
4. 学习模板方法的最佳实践

### 3.3.1.11.3 高阶路径

1. 分析模板方法在大型项目中的应用
2. 理解模板方法与架构设计的关系
3. 掌握模板方法的性能调优
4. 学习模板方法的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)
