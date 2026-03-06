# Wire 依赖注入全面解析

> **版本**: v1.0
> **日期**: 2025-01-XX
> **状态**: ✅ 完整版

---

## 📋 目录

- [Wire 依赖注入全面解析](#wire-依赖注入全面解析)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
    - [1.1 什么是 Wire？](#11-什么是-wire)
    - [1.2 核心特性](#12-核心特性)
    - [1.3 为什么选择 Wire？](#13-为什么选择-wire)
      - [1.3.1 与其他 DI 工具对比](#131-与其他-di-工具对比)
      - [1.3.2 选择 Wire 的理由](#132-选择-wire-的理由)
  - [2. 核心原理](#2-核心原理)
    - [2.1 工作原理](#21-工作原理)
    - [2.2 依赖解析算法](#22-依赖解析算法)
    - [2.3 类型匹配机制](#23-类型匹配机制)
  - [3. 知识矩阵](#3-知识矩阵)
    - [3.1 Wire 知识矩阵](#31-wire-知识矩阵)
    - [3.2 依赖注入模式矩阵](#32-依赖注入模式矩阵)
    - [3.3 架构层次矩阵](#33-架构层次矩阵)
  - [4. 思维导图](#4-思维导图)
    - [4.1 Wire 核心概念思维导图](#41-wire-核心概念思维导图)
    - [4.2 依赖关系思维导图](#42-依赖关系思维导图)
    - [4.3 错误处理思维导图](#43-错误处理思维导图)
  - [5. 深度论证](#5-深度论证)
    - [5.1 为什么需要依赖注入？](#51-为什么需要依赖注入)
      - [5.1.1 问题：紧耦合](#511-问题紧耦合)
      - [5.1.2 解决方案：依赖注入](#512-解决方案依赖注入)
      - [5.1.3 新问题：手动管理依赖](#513-新问题手动管理依赖)
      - [5.1.4 最终解决方案：Wire](#514-最终解决方案wire)
    - [5.2 Wire vs 其他 DI 工具](#52-wire-vs-其他-di-工具)
      - [5.2.1 Wire vs Dig](#521-wire-vs-dig)
      - [5.2.2 Wire vs 手动注入](#522-wire-vs-手动注入)
    - [5.3 Wire 的适用场景](#53-wire-的适用场景)
      - [5.3.1 适合使用 Wire 的场景](#531-适合使用-wire-的场景)
      - [5.3.2 不适合使用 Wire 的场景](#532-不适合使用-wire-的场景)
  - [6. 完整示例](#6-完整示例)
    - [6.1 基础示例](#61-基础示例)
    - [6.2 使用 Provider 集合](#62-使用-provider-集合)
    - [6.3 使用接口绑定](#63-使用接口绑定)
    - [6.4 使用值绑定](#64-使用值绑定)
    - [6.5 使用结构体 Provider](#65-使用结构体-provider)
  - [7. 最佳实践](#7-最佳实践)
    - [7.1 Provider 函数设计](#71-provider-函数设计)
      - [7.1.1 命名规范](#711-命名规范)
      - [7.1.2 单一职责](#712-单一职责)
      - [7.1.3 错误处理](#713-错误处理)
    - [7.2 依赖关系设计](#72-依赖关系设计)
      - [7.2.1 避免循环依赖](#721-避免循环依赖)
      - [7.2.2 按层次组织](#722-按层次组织)
    - [7.3 测试策略](#73-测试策略)
      - [7.3.1 Mock Provider](#731-mock-provider)
      - [7.3.2 测试 Provider](#732-测试-provider)
  - [8. 常见问题](#8-常见问题)
    - [8.1 如何调试 Wire 生成的代码？](#81-如何调试-wire-生成的代码)
    - [8.2 如何处理循环依赖？](#82-如何处理循环依赖)
    - [8.3 如何在运行时动态配置依赖？](#83-如何在运行时动态配置依赖)
    - [8.4 如何处理可选依赖？](#84-如何处理可选依赖)
  - [9. 总结](#9-总结)

---

## 1. 概述

### 1.1 什么是 Wire？

**Wire** 是 Google 开源的 Go 语言依赖注入工具，它通过**编译时代码生成**实现依赖注入，而不是运行时反射。

### 1.2 核心特性

| 特性 | 说明 | 优势 |
|------|------|------|
| **编译时生成** | 在编译时生成依赖注入代码 | 类型安全、性能优秀 |
| **零反射** | 不使用运行时反射 | 性能优秀、易于调试 |
| **类型安全** | 使用 Go 的类型系统 | 编译时检查、IDE 支持 |
| **易于调试** | 生成的代码可查看 | 易于理解和调试 |
| **IDE 支持** | IDE 可以理解依赖关系 | 代码补全、重构支持 |

### 1.3 为什么选择 Wire？

#### 1.3.1 与其他 DI 工具对比

| 工具 | 实现方式 | 性能 | 类型安全 | 调试难度 |
|------|---------|------|---------|---------|
| **Wire** | 编译时代码生成 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Dig** | 运行时反射 | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **Fx** | 运行时反射 | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **手动注入** | 手动编写 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |

#### 1.3.2 选择 Wire 的理由

1. **性能优势**：编译时生成，零反射，性能最优
2. **类型安全**：编译时检查，避免运行时错误
3. **易于调试**：生成的代码可查看，易于理解
4. **IDE 支持**：IDE 可以理解依赖关系，提供代码补全
5. **Google 支持**：Google 官方维护，稳定可靠

---

## 2. 核心原理

### 2.1 工作原理

```text
┌─────────────────────────────────────────────────────────────┐
│                    Wire 工作流程                              │
└─────────────────────────────────────────────────────────────┘

1. 定义 Provider 函数
   ↓
   func NewDatabase(cfg *Config) (*Database, error) { ... }
   func NewRepository(db *Database) Repository { ... }
   func NewService(repo Repository) *Service { ... }

2. 声明依赖关系（wire.Build）
   ↓
   func InitializeApp(cfg *Config) (*App, error) {
       wire.Build(
           NewDatabase,
           NewRepository,
           NewService,
           NewApp,
       )
       return nil, nil
   }

3. 运行 Wire 生成代码
   ↓
   $ wire ./scripts/wire

4. 生成 wire_gen.go
   ↓
   func InitializeApp(cfg *Config) (*App, error) {
       // 自动生成的依赖注入代码
       database, err := NewDatabase(cfg)
       if err != nil {
           return nil, err
       }
       repository := NewRepository(database)
       service := NewService(repository)
       app := NewApp(service)
       return app, nil
   }

5. 使用生成的代码
   ↓
   app, err := wire.InitializeApp(cfg)
   if err != nil {
       log.Fatal(err)
   }
```

### 2.2 依赖解析算法

Wire 使用**拓扑排序**算法解析依赖关系：

```text
1. 构建依赖图（Dependency Graph）
   - 节点：Provider 函数的返回值类型
   - 边：Provider 函数的参数依赖

2. 拓扑排序（Topological Sort）
   - 找到所有没有依赖的节点（入度为 0）
   - 依次处理这些节点
   - 更新依赖图，继续处理

3. 生成代码
   - 按照拓扑排序的顺序生成代码
   - 确保依赖在依赖者之前创建
```

### 2.3 类型匹配机制

Wire 通过**类型匹配**确定依赖关系：

```go
// Provider 函数签名
func NewDatabase(cfg *Config) (*Database, error)

// 依赖注入时
func NewRepository(db *Database) Repository  // 匹配 *Database 类型
func NewService(repo Repository) *Service    // 匹配 Repository 接口类型
```

**匹配规则**：

1. **精确匹配**：类型完全一致
2. **接口匹配**：实现类型匹配接口类型
3. **指针匹配**：`*T` 和 `T` 可以相互匹配（通过解引用/取地址）

---

## 3. 知识矩阵

### 3.1 Wire 知识矩阵

| 维度 | 知识点 | 说明 | 重要性 |
|------|--------|------|--------|
| **基础概念** | Provider 函数 | 创建依赖的函数 | ⭐⭐⭐⭐⭐ |
| | wire.Build | 声明依赖关系 | ⭐⭐⭐⭐⭐ |
| | 构建标签 | `//go:build wireinject` | ⭐⭐⭐⭐ |
| | 生成代码 | `wire_gen.go` | ⭐⭐⭐⭐ |
| **高级特性** | Provider 集合 | `wire.NewSet` | ⭐⭐⭐⭐ |
| | 接口绑定 | `wire.Bind` | ⭐⭐⭐ |
| | 值绑定 | `wire.Value` | ⭐⭐⭐ |
| | 结构体 Provider | `wire.Struct` | ⭐⭐⭐ |
| | 字段 Provider | `wire.FieldsOf` | ⭐⭐⭐ |
| **错误处理** | 错误返回 | Provider 返回 error | ⭐⭐⭐⭐⭐ |
| | 错误传播 | 自动传播错误 | ⭐⭐⭐⭐ |
| **最佳实践** | 命名规范 | `NewXxx` 格式 | ⭐⭐⭐⭐ |
| | 层次组织 | 按架构层次组织 | ⭐⭐⭐⭐⭐ |
| | 单一职责 | 每个 Provider 只创建一个依赖 | ⭐⭐⭐⭐⭐ |
| | 避免循环依赖 | 设计单向依赖 | ⭐⭐⭐⭐⭐ |

### 3.2 依赖注入模式矩阵

| 模式 | 适用场景 | 示例 | 优缺点 |
|------|---------|------|--------|
| **构造函数注入** | 大多数场景 | `NewService(repo Repository)` | ✅ 类型安全<br>✅ 易于测试<br>❌ 参数较多时复杂 |
| **接口注入** | 需要多态 | `NewService(repo UserRepository)` | ✅ 灵活<br>✅ 可替换实现<br>❌ 需要定义接口 |
| **值注入** | 配置、常量 | `wire.Value(cfg)` | ✅ 简单<br>❌ 不够灵活 |
| **结构体注入** | 复杂对象 | `wire.Struct(new(App), "*")` | ✅ 自动注入字段<br>❌ 不够明确 |

### 3.3 架构层次矩阵

| 层次 | Provider 类型 | 依赖关系 | 示例 |
|------|--------------|---------|------|
| **配置层** | Config Provider | 无依赖 | `NewConfig()` |
| **基础设施层** | Infrastructure Provider | 依赖 Config | `NewDatabase(cfg)` |
| **领域层** | Domain Provider | 依赖 Infrastructure | `NewRepository(db)` |
| **应用层** | Application Provider | 依赖 Domain | `NewService(repo)` |
| **接口层** | Interface Provider | 依赖 Application | `NewRouter(service)` |
| **应用组装** | App Provider | 依赖所有层 | `NewApp(router)` |

---

## 4. 思维导图

### 4.1 Wire 核心概念思维导图

```text
Wire 依赖注入
│
├── 核心概念
│   ├── Provider 函数
│   │   ├── 定义：创建依赖的函数
│   │   ├── 命名：NewXxx 格式
│   │   ├── 返回值：依赖对象 + error
│   │   └── 参数：声明依赖关系
│   │
│   ├── wire.Build
│   │   ├── 作用：声明依赖关系
│   │   ├── 参数：Provider 函数列表
│   │   └── 返回：生成的代码
│   │
│   └── 构建标签
│       ├── wireinject：标记需要生成的函数
│       └── !wireinject：标记生成的代码
│
├── 工作流程
│   ├── 1. 定义 Provider
│   ├── 2. 声明依赖（wire.Build）
│   ├── 3. 运行 Wire 生成代码
│   ├── 4. 使用生成的代码
│   └── 5. 编译运行
│
├── 高级特性
│   ├── Provider 集合（wire.NewSet）
│   ├── 接口绑定（wire.Bind）
│   ├── 值绑定（wire.Value）
│   ├── 结构体 Provider（wire.Struct）
│   └── 字段 Provider（wire.FieldsOf）
│
└── 最佳实践
    ├── 命名规范
    ├── 层次组织
    ├── 单一职责
    ├── 错误处理
    └── 避免循环依赖
```

### 4.2 依赖关系思维导图

```text
依赖关系图
│
├── 配置层（Config）
│   └── NewConfig() → *Config
│
├── 基础设施层（Infrastructure）
│   ├── NewDatabase(cfg) → *Database
│   ├── NewCache(cfg) → *Cache
│   └── NewMQ(cfg) → *MessageQueue
│
├── 领域层（Domain）
│   ├── NewUserRepository(db) → UserRepository
│   ├── NewOrderRepository(db) → OrderRepository
│   └── NewProductRepository(db) → ProductRepository
│
├── 应用层（Application）
│   ├── NewUserService(repo) → *UserService
│   ├── NewOrderService(repo) → *OrderService
│   └── NewProductService(repo) → *ProductService
│
├── 接口层（Interface）
│   ├── NewHTTPRouter(services) → *Router
│   ├── NewGRPCServer(services) → *Server
│   └── NewGraphQLServer(services) → *Server
│
└── 应用组装（App）
    └── NewApp(router, servers) → *App
```

### 4.3 错误处理思维导图

```text
错误处理
│
├── Provider 函数错误
│   ├── 返回 error
│   ├── 错误传播
│   └── 错误包装
│
├── 依赖创建失败
│   ├── 立即返回错误
│   ├── 不创建后续依赖
│   └── 清理已创建的资源
│
└── 错误处理最佳实践
    ├── 使用 fmt.Errorf 包装错误
    ├── 提供上下文信息
    └── 避免静默失败
```

---

## 5. 深度论证

### 5.1 为什么需要依赖注入？

#### 5.1.1 问题：紧耦合

```go
// ❌ 紧耦合：直接创建依赖
type UserService struct {
    repo *UserRepository
}

func NewUserService() *UserService {
    // 直接创建依赖，难以测试和替换
    db := sql.Open("postgres", "connection string")
    repo := NewUserRepository(db)
    return &UserService{repo: repo}
}
```

**问题**：

- 难以测试（无法 mock 依赖）
- 难以替换实现（硬编码依赖）
- 违反依赖倒置原则（依赖具体实现）

#### 5.1.2 解决方案：依赖注入

```go
// ✅ 依赖注入：通过参数注入依赖
type UserService struct {
    repo UserRepository  // 依赖接口，不依赖具体实现
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

**优势**：

- 易于测试（可以注入 mock 对象）
- 易于替换实现（注入不同的实现）
- 符合依赖倒置原则（依赖抽象）

#### 5.1.3 新问题：手动管理依赖

```go
// ❌ 手动管理依赖：复杂且容易出错
func main() {
    cfg := NewConfig()
    db, err := NewDatabase(cfg)
    if err != nil {
        log.Fatal(err)
    }
    repo := NewUserRepository(db)
    service := NewUserService(repo)
    router := NewRouter(service)
    app := NewApp(router)
    // ... 更多依赖管理代码
}
```

**问题**：

- 代码冗长
- 容易出错（依赖顺序错误）
- 难以维护（依赖关系复杂时）

#### 5.1.4 最终解决方案：Wire

```go
// ✅ Wire：自动管理依赖
func InitializeApp(cfg *Config) (*App, error) {
    wire.Build(
        NewDatabase,
        NewUserRepository,
        NewUserService,
        NewRouter,
        NewApp,
    )
    return nil, nil
}

// 使用
app, err := wire.InitializeApp(cfg)
```

**优势**：

- 自动管理依赖关系
- 编译时检查依赖
- 生成的代码可查看
- 易于维护

### 5.2 Wire vs 其他 DI 工具

#### 5.2.1 Wire vs Dig

| 维度 | Wire | Dig |
|------|------|-----|
| **实现方式** | 编译时代码生成 | 运行时反射 |
| **性能** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **类型安全** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **调试难度** | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **学习曲线** | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **IDE 支持** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |

**结论**：Wire 在性能、类型安全、调试方面优于 Dig。

#### 5.2.2 Wire vs 手动注入

| 维度 | Wire | 手动注入 |
|------|------|---------|
| **代码量** | 少（自动生成） | 多（手动编写） |
| **维护成本** | 低 | 高 |
| **错误率** | 低（编译时检查） | 高（容易出错） |
| **灵活性** | 中 | 高 |

**结论**：Wire 在代码量、维护成本、错误率方面优于手动注入。

### 5.3 Wire 的适用场景

#### 5.3.1 适合使用 Wire 的场景

1. **大型项目**：依赖关系复杂，需要统一管理
2. **多环境部署**：需要不同的依赖配置
3. **测试驱动开发**：需要频繁 mock 依赖
4. **团队协作**：需要统一的依赖注入规范

#### 5.3.2 不适合使用 Wire 的场景

1. **小型项目**：依赖关系简单，手动注入即可
2. **原型开发**：快速迭代，不需要复杂的依赖管理
3. **学习项目**：需要理解依赖注入原理

---

## 6. 完整示例

### 6.1 基础示例

```go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "github.com/yourusername/golang/internal/config"
    "github.com/yourusername/golang/internal/infrastructure/database"
    "github.com/yourusername/golang/internal/application/user"
    "github.com/yourusername/golang/internal/interfaces/http"
)

// InitializeApp 初始化应用
func InitializeApp(cfg *config.Config) (*App, error) {
    wire.Build(
        // Infrastructure
        database.NewDatabase,
        database.NewUserRepository,

        // Application
        user.NewService,

        // Interface
        http.NewRouter,

        // App
        NewApp,
    )
    return nil, nil
}
```

### 6.2 使用 Provider 集合

```go
// Provider 集合
var (
    // 数据库相关 Provider
    DatabaseProviderSet = wire.NewSet(
        database.NewDatabase,
        database.NewUserRepository,
        database.NewOrderRepository,
    )

    // 应用层 Provider
    ApplicationProviderSet = wire.NewSet(
        user.NewService,
        order.NewService,
    )

    // 接口层 Provider
    InterfaceProviderSet = wire.NewSet(
        http.NewRouter,
        grpc.NewServer,
    )
)

// 使用 Provider 集合
func InitializeApp(cfg *config.Config) (*App, error) {
    wire.Build(
        DatabaseProviderSet,
        ApplicationProviderSet,
        InterfaceProviderSet,
        NewApp,
    )
    return nil, nil
}
```

### 6.3 使用接口绑定

```go
// 定义接口
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Get(ctx context.Context, id string) (*User, error)
}

// Provider 返回具体实现
func NewUserRepository(db *Database) *UserRepositoryImpl {
    return &UserRepositoryImpl{db: db}
}

// 使用接口绑定
var RepositoryProviderSet = wire.NewSet(
    NewUserRepository,
    wire.Bind(new(UserRepository), new(*UserRepositoryImpl)),
)
```

### 6.4 使用值绑定

```go
// 绑定配置值
func InitializeApp() (*App, error) {
    wire.Build(
        wire.Value(&Config{
            DatabaseURL: "postgres://...",
            Port:        8080,
        }),
        NewDatabase,
        NewApp,
    )
    return nil, nil
}
```

### 6.5 使用结构体 Provider

```go
// 自动注入结构体字段
type App struct {
    Database *Database
    Service  *Service
    Router   *Router
}

func InitializeApp() (*App, error) {
    wire.Build(
        NewDatabase,
        NewService,
        NewRouter,
        wire.Struct(new(App), "*"),  // 注入所有字段
    )
    return nil, nil
}
```

---

## 7. 最佳实践

### 7.1 Provider 函数设计

#### 7.1.1 命名规范

```go
// ✅ 正确：使用 NewXxx 格式
func NewDatabase(cfg *Config) (*Database, error)
func NewUserRepository(db *Database) UserRepository
func NewUserService(repo UserRepository) *UserService

// ❌ 错误：不使用 NewXxx 格式
func CreateDatabase(cfg *Config) (*Database, error)
func MakeUserRepository(db *Database) UserRepository
```

#### 7.1.2 单一职责

```go
// ✅ 正确：每个 Provider 只创建一个依赖
func NewDatabase(cfg *Config) (*Database, error) {
    return sql.Open("postgres", cfg.DSN)
}

func NewUserRepository(db *Database) UserRepository {
    return &UserRepositoryImpl{db: db}
}

// ❌ 错误：一个 Provider 创建多个依赖
func NewDatabaseAndRepository(cfg *Config) (*Database, UserRepository, error) {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil {
        return nil, nil, err
    }
    repo := NewUserRepository(db)
    return db, repo, nil
}
```

#### 7.1.3 错误处理

```go
// ✅ 正确：返回错误并提供上下文
func NewDatabase(cfg *Config) (*Database, error) {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    return db, nil
}

// ❌ 错误：静默失败或不提供上下文
func NewDatabase(cfg *Config) (*Database, error) {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil {
        return nil, err  // 缺少上下文信息
    }
    return db, nil
}
```

### 7.2 依赖关系设计

#### 7.2.1 避免循环依赖

```go
// ❌ 错误：循环依赖
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func NewRepository(service *Service) Repository {
    return &RepositoryImpl{service: service}
}

// ✅ 正确：单向依赖
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func NewRepository(db *Database) Repository {
    return &RepositoryImpl{db: db}
}
```

#### 7.2.2 按层次组织

```go
// ✅ 正确：按架构层次组织 Provider
var (
    // Infrastructure Layer
    InfrastructureProviderSet = wire.NewSet(
        NewDatabase,
        NewCache,
        NewMQ,
    )

    // Domain Layer
    DomainProviderSet = wire.NewSet(
        NewUserRepository,
        NewOrderRepository,
    )

    // Application Layer
    ApplicationProviderSet = wire.NewSet(
        NewUserService,
        NewOrderService,
    )

    // Interface Layer
    InterfaceProviderSet = wire.NewSet(
        NewRouter,
        NewGRPCServer,
    )
)
```

### 7.3 测试策略

#### 7.3.1 Mock Provider

```go
// 测试时使用 Mock Provider
func TestUserService(t *testing.T) {
    // 创建 Mock Repository
    mockRepo := &MockUserRepository{}

    // 手动注入 Mock
    service := NewUserService(mockRepo)

    // 测试...
}
```

#### 7.3.2 测试 Provider

```go
// 为测试创建专门的 Provider
func NewTestDatabase() (*Database, error) {
    // 使用测试数据库
    return sql.Open("sqlite3", ":memory:")
}

// 测试时使用测试 Provider
func TestApp(t *testing.T) {
    app, err := InitializeTestApp()
    if err != nil {
        t.Fatal(err)
    }
    // 测试...
}
```

---

## 8. 常见问题

### 8.1 如何调试 Wire 生成的代码？

**答案**：查看 `wire_gen.go` 文件，Wire 会生成完整的依赖注入代码，可以直接查看和调试。

### 8.2 如何处理循环依赖？

**答案**：

1. **重新设计依赖关系**：避免循环依赖
2. **使用接口**：通过接口解耦
3. **延迟初始化**：使用 `wire.Bind` 和接口

### 8.3 如何在运行时动态配置依赖？

**答案**：

1. **使用配置对象**：通过 Config Provider 注入配置
2. **使用环境变量**：在 Config Provider 中读取环境变量
3. **使用配置文件**：在 Config Provider 中读取配置文件

### 8.4 如何处理可选依赖？

**答案**：

1. **使用接口**：定义可选接口，实现为空操作
2. **使用指针**：使用 `*Type` 表示可选依赖
3. **使用 Provider 集合**：为不同场景创建不同的 Provider 集合

---

## 9. 总结

Wire 是一个强大的依赖注入工具，通过编译时代码生成实现类型安全的依赖注入。它适合大型项目，能够显著提高代码的可维护性和可测试性。

**核心优势**：

- ✅ 编译时检查，类型安全
- ✅ 零反射，性能优秀
- ✅ 生成的代码可查看，易于调试
- ✅ IDE 支持良好

**适用场景**：

- ✅ 大型项目
- ✅ 多环境部署
- ✅ 测试驱动开发
- ✅ 团队协作

---

**最后更新**: 2025-01-XX
