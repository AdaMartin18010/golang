# Wire 知识矩阵与思维导图

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📊 知识矩阵

### 1. Wire 核心概念矩阵

| 概念 | 定义 | 重要性 | 难度 | 使用频率 |
|------|------|--------|------|---------|
| **Provider 函数** | 创建依赖的函数 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **wire.Build** | 声明依赖关系 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **构建标签** | `//go:build wireinject` | ⭐⭐⭐⭐ | ⭐ | ⭐⭐⭐⭐ |
| **生成代码** | `wire_gen.go` | ⭐⭐⭐⭐ | ⭐ | ⭐⭐⭐⭐ |
| **Provider 集合** | `wire.NewSet` | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **接口绑定** | `wire.Bind` | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **值绑定** | `wire.Value` | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **结构体 Provider** | `wire.Struct` | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **字段 Provider** | `wire.FieldsOf` | ⭐⭐ | ⭐⭐⭐ | ⭐ |

### 2. 依赖注入模式矩阵

| 模式 | 适用场景 | 复杂度 | 灵活性 | 推荐度 |
|------|---------|--------|--------|--------|
| **构造函数注入** | 大多数场景 | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **接口注入** | 需要多态 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **值注入** | 配置、常量 | ⭐ | ⭐⭐ | ⭐⭐⭐ |
| **结构体注入** | 复杂对象 | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |

### 3. 架构层次矩阵

| 层次 | Provider 类型 | 依赖来源 | 被依赖 | 复杂度 |
|------|--------------|---------|--------|--------|
| **配置层** | Config | 无 | Infrastructure | ⭐ |
| **基础设施层** | Database, Cache, MQ | Config | Domain | ⭐⭐ |
| **领域层** | Repository | Infrastructure | Application | ⭐⭐⭐ |
| **应用层** | Service | Domain | Interface | ⭐⭐⭐ |
| **接口层** | Router, Server | Application | App | ⭐⭐ |
| **应用组装** | App | Interface | 无 | ⭐ |

### 4. 错误处理矩阵

| 场景 | 处理方式 | 重要性 | 复杂度 |
|------|---------|--------|--------|
| **Provider 错误** | 返回 error | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **依赖创建失败** | 立即返回 | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **错误包装** | fmt.Errorf | ⭐⭐⭐⭐ | ⭐ |
| **错误上下文** | 提供详细信息 | ⭐⭐⭐⭐ | ⭐⭐ |

### 5. 测试策略矩阵

| 策略 | 适用场景 | 复杂度 | 效果 |
|------|---------|--------|------|
| **Mock Provider** | 单元测试 | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **测试 Provider** | 集成测试 | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **测试配置** | 环境隔离 | ⭐⭐ | ⭐⭐⭐⭐ |

---

## 🗺️ 思维导图

### 1. Wire 核心概念思维导图

```text
Wire 依赖注入
│
├── 基础概念
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
│   │   └── func NewXxx(...) (*Xxx, error)
│   │
│   ├── 2. 声明依赖（wire.Build）
│   │   └── wire.Build(NewXxx, NewYyy, ...)
│   │
│   ├── 3. 运行 Wire 生成代码
│   │   └── $ wire ./scripts/wire
│   │
│   ├── 4. 使用生成的代码
│   │   └── app, err := wire.InitializeApp(cfg)
│   │
│   └── 5. 编译运行
│       └── $ go build
│
├── 高级特性
│   ├── Provider 集合（wire.NewSet）
│   │   ├── 组织相关 Provider
│   │   └── 提高可维护性
│   │
│   ├── 接口绑定（wire.Bind）
│   │   ├── 绑定接口和实现
│   │   └── 支持多态
│   │
│   ├── 值绑定（wire.Value）
│   │   ├── 注入常量值
│   │   └── 注入配置
│   │
│   ├── 结构体 Provider（wire.Struct）
│   │   ├── 自动注入字段
│   │   └── 简化代码
│   │
│   └── 字段 Provider（wire.FieldsOf）
│       ├── 从结构体提取字段
│       └── 支持嵌套结构
│
└── 最佳实践
    ├── 命名规范
    │   └── NewXxx 格式
    │
    ├── 层次组织
    │   └── 按架构层次组织 Provider
    │
    ├── 单一职责
    │   └── 每个 Provider 只创建一个依赖
    │
    ├── 错误处理
    │   └── 返回 error 并提供上下文
    │
    └── 避免循环依赖
        └── 设计单向依赖
```

### 2. 依赖关系思维导图

```
依赖关系图
│
├── 配置层（Config Layer）
│   └── NewConfig() → *Config
│       │
│       └── 被依赖：Infrastructure Layer
│
├── 基础设施层（Infrastructure Layer）
│   ├── NewDatabase(cfg) → *Database
│   │   └── 被依赖：Domain Layer
│   │
│   ├── NewCache(cfg) → *Cache
│   │   └── 被依赖：Domain Layer
│   │
│   └── NewMQ(cfg) → *MessageQueue
│       └── 被依赖：Application Layer
│
├── 领域层（Domain Layer）
│   ├── NewUserRepository(db) → UserRepository
│   │   └── 被依赖：Application Layer
│   │
│   ├── NewOrderRepository(db) → OrderRepository
│   │   └── 被依赖：Application Layer
│   │
│   └── NewProductRepository(db) → ProductRepository
│       └── 被依赖：Application Layer
│
├── 应用层（Application Layer）
│   ├── NewUserService(repo) → *UserService
│   │   └── 被依赖：Interface Layer
│   │
│   ├── NewOrderService(repo) → *OrderService
│   │   └── 被依赖：Interface Layer
│   │
│   └── NewProductService(repo) → *ProductService
│       └── 被依赖：Interface Layer
│
├── 接口层（Interface Layer）
│   ├── NewHTTPRouter(services) → *Router
│   │   └── 被依赖：App
│   │
│   ├── NewGRPCServer(services) → *Server
│   │   └── 被依赖：App
│   │
│   └── NewGraphQLServer(services) → *Server
│       └── 被依赖：App
│
└── 应用组装（App Assembly）
    └── NewApp(router, servers) → *App
        └── 无被依赖（顶层）
```

### 3. 错误处理思维导图

```
错误处理
│
├── Provider 函数错误
│   ├── 返回 error
│   │   └── func NewXxx(...) (*Xxx, error)
│   │
│   ├── 错误传播
│   │   └── Wire 自动传播错误
│   │
│   └── 错误包装
│       └── fmt.Errorf("context: %w", err)
│
├── 依赖创建失败
│   ├── 立即返回错误
│   │   └── 不创建后续依赖
│   │
│   ├── 不创建后续依赖
│   │   └── 避免部分初始化
│   │
│   └── 清理已创建的资源
│       └── defer 清理（如果需要）
│
└── 错误处理最佳实践
    ├── 使用 fmt.Errorf 包装错误
    │   └── 提供上下文信息
    │
    ├── 提供上下文信息
    │   └── 说明失败的原因和位置
    │
    └── 避免静默失败
        └── 始终返回错误
```

### 4. 测试策略思维导图

```
测试策略
│
├── Mock Provider
│   ├── 创建 Mock 对象
│   │   └── type MockRepository struct { ... }
│   │
│   ├── 手动注入 Mock
│   │   └── service := NewService(mockRepo)
│   │
│   └── 适用场景
│       └── 单元测试
│
├── 测试 Provider
│   ├── 创建测试 Provider
│   │   └── func NewTestDatabase() (*Database, error)
│   │
│   ├── 使用测试配置
│   │   └── 内存数据库、测试配置等
│   │
│   └── 适用场景
│       └── 集成测试
│
└── 测试配置
    ├── 环境隔离
    │   └── 测试环境独立配置
    │
    ├── 数据隔离
    │   └── 测试数据独立
    │
    └── 适用场景
        └── 端到端测试
```

---

## 📈 学习路径

### 初学者路径

```
1. 理解依赖注入概念
   ↓
2. 学习 Provider 函数
   ↓
3. 学习 wire.Build
   ↓
4. 运行第一个 Wire 示例
   ↓
5. 理解生成的代码
```

### 进阶路径

```
1. 学习 Provider 集合
   ↓
2. 学习接口绑定
   ↓
3. 学习高级特性
   ↓
4. 实践复杂场景
   ↓
5. 优化依赖关系
```

### 专家路径

```
1. 深入理解 Wire 原理
   ↓
2. 设计复杂依赖关系
   ↓
3. 优化性能
   ↓
4. 贡献代码
   ↓
5. 分享经验
```

---

**最后更新**: 2025-01-XX
