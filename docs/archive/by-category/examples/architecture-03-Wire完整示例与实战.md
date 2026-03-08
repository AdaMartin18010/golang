# Wire 完整示例与实战

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录

- [Wire 完整示例与实战](#wire-完整示例与实战)
  - [📋 目录](#-目录)
  - [1. 基础示例](#1-基础示例)
    - [1.1 最简单的示例](#11-最简单的示例)
    - [1.2 带错误的示例](#12-带错误的示例)
  - [2. 进阶示例](#2-进阶示例)
    - [2.1 使用 Provider 集合](#21-使用-provider-集合)
    - [2.2 使用接口绑定](#22-使用接口绑定)
    - [2.3 使用值绑定](#23-使用值绑定)
    - [2.4 使用结构体 Provider](#24-使用结构体-provider)
  - [3. 实战案例](#3-实战案例)
    - [3.1 完整的 HTTP 服务](#31-完整的-http-服务)
    - [3.2 多环境配置](#32-多环境配置)
    - [3.3 条件依赖](#33-条件依赖)
  - [4. 常见模式](#4-常见模式)
    - [4.1 单例模式](#41-单例模式)
    - [4.2 工厂模式](#42-工厂模式)
    - [4.3 装饰器模式](#43-装饰器模式)
  - [5. 故障排查](#5-故障排查)
    - [5.1 常见错误](#51-常见错误)
      - [错误 1: 缺少 Provider](#错误-1-缺少-provider)
      - [错误 2: 循环依赖](#错误-2-循环依赖)
      - [错误 3: 类型不匹配](#错误-3-类型不匹配)
    - [5.2 调试技巧](#52-调试技巧)
      - [技巧 1: 查看生成的代码](#技巧-1-查看生成的代码)
      - [技巧 2: 使用 Wire 的调试选项](#技巧-2-使用-wire-的调试选项)
      - [技巧 3: 逐步添加 Provider](#技巧-3-逐步添加-provider)
  - [6. 性能优化](#6-性能优化)
    - [6.1 减少 Provider 函数调用](#61-减少-provider-函数调用)
    - [6.2 延迟初始化](#62-延迟初始化)

---

## 1. 基础示例

### 1.1 最简单的示例

```go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
)

// Config 配置
type Config struct {
    DatabaseURL string
    Port        int
}

// Database 数据库
type Database struct {
    URL string
}

// Service 服务
type Service struct {
    DB *Database
}

// App 应用
type App struct {
    Service *Service
}

// Provider 函数
func NewConfig() *Config {
    return &Config{
        DatabaseURL: "postgres://localhost/db",
        Port:        8080,
    }
}

func NewDatabase(cfg *Config) (*Database, error) {
    return &Database{URL: cfg.DatabaseURL}, nil
}

func NewService(db *Database) *Service {
    return &Service{DB: db}
}

func NewApp(service *Service) *App {
    return &App{Service: service}
}

// Wire 配置
func InitializeApp() (*App, error) {
    wire.Build(
        NewConfig,
        NewDatabase,
        NewService,
        NewApp,
    )
    return nil, nil
}
```

**生成的代码** (`wire_gen.go`):

```go
//go:build !wireinject
// +build !wireinject

package wire

// InitializeApp 初始化应用（自动生成）
func InitializeApp() (*App, error) {
    config := NewConfig()
    database, err := NewDatabase(config)
    if err != nil {
        return nil, err
    }
    service := NewService(database)
    app := NewApp(service)
    return app, nil
}
```

### 1.2 带错误的示例

```go
// Provider 函数返回错误
func NewDatabase(cfg *Config) (*Database, error) {
    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("database URL is required")
    }
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    return &Database{conn: db}, nil
}

// Wire 自动处理错误传播
func InitializeApp() (*App, error) {
    wire.Build(
        NewConfig,
        NewDatabase,  // 如果这里失败，整个初始化失败
        NewService,
        NewApp,
    )
    return nil, nil
}
```

---

## 2. 进阶示例

### 2.1 使用 Provider 集合

```go
// 定义 Provider 集合
var (
    // 数据库相关 Provider
    DatabaseProviderSet = wire.NewSet(
        NewDatabase,
        NewUserRepository,
        NewOrderRepository,
    )

    // 服务相关 Provider
    ServiceProviderSet = wire.NewSet(
        NewUserService,
        NewOrderService,
    )

    // 接口相关 Provider
    InterfaceProviderSet = wire.NewSet(
        NewHTTPRouter,
        NewGRPCServer,
    )
)

// 使用 Provider 集合
func InitializeApp(cfg *Config) (*App, error) {
    wire.Build(
        DatabaseProviderSet,
        ServiceProviderSet,
        InterfaceProviderSet,
        NewApp,
    )
    return nil, nil
}
```

### 2.2 使用接口绑定

```go
// 定义接口
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Get(ctx context.Context, id string) (*User, error)
}

// 实现接口
type UserRepositoryImpl struct {
    db *Database
}

func NewUserRepository(db *Database) *UserRepositoryImpl {
    return &UserRepositoryImpl{db: db}
}

// 使用接口绑定
var RepositoryProviderSet = wire.NewSet(
    NewUserRepository,
    wire.Bind(new(UserRepository), new(*UserRepositoryImpl)),
)

// Service 依赖接口
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### 2.3 使用值绑定

```go
// 绑定配置值
func InitializeApp() (*App, error) {
    wire.Build(
        wire.Value(&Config{
            DatabaseURL: "postgres://localhost/db",
            Port:        8080,
        }),
        NewDatabase,
        NewApp,
    )
    return nil, nil
}

// 或者绑定多个值
func InitializeApp() (*App, error) {
    wire.Build(
        wire.Values(
            &Config{DatabaseURL: "postgres://localhost/db"},
            &LoggerConfig{Level: "info"},
        ),
        NewDatabase,
        NewLogger,
        NewApp,
    )
    return nil, nil
}
```

### 2.4 使用结构体 Provider

```go
// 定义结构体
type App struct {
    Database *Database
    Service  *Service
    Router   *Router
}

// 使用结构体 Provider
func InitializeApp() (*App, error) {
    wire.Build(
        NewDatabase,
        NewService,
        NewRouter,
        wire.Struct(new(App), "*"),  // 注入所有字段
    )
    return nil, nil
}

// 或者只注入特定字段
func InitializeApp() (*App, error) {
    wire.Build(
        NewDatabase,
        NewService,
        NewRouter,
        wire.Struct(new(App), "Database", "Service"),  // 只注入指定字段
    )
    return nil, nil
}
```

---

## 3. 实战案例

### 3.1 完整的 HTTP 服务

```go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "github.com/yourusername/golang/internal/config"
    "github.com/yourusername/golang/internal/infrastructure/database"
    "github.com/yourusername/golang/internal/infrastructure/cache"
    "github.com/yourusername/golang/internal/infrastructure/messaging"
    "github.com/yourusername/golang/internal/application/user"
    "github.com/yourusername/golang/internal/application/order"
    "github.com/yourusername/golang/internal/interfaces/http"
    "github.com/yourusername/golang/internal/interfaces/grpc"
)

// Provider 集合
var (
    // 配置 Provider
    ConfigProviderSet = wire.NewSet(
        config.NewConfig,
    )

    // 基础设施 Provider
    InfrastructureProviderSet = wire.NewSet(
        database.NewDatabase,
        cache.NewCache,
        messaging.NewMessageQueue,
    )

    // 仓储 Provider
    RepositoryProviderSet = wire.NewSet(
        database.NewUserRepository,
        database.NewOrderRepository,
        wire.Bind(new(user.Repository), new(*database.UserRepository)),
        wire.Bind(new(order.Repository), new(*database.OrderRepository)),
    )

    // 服务 Provider
    ServiceProviderSet = wire.NewSet(
        user.NewService,
        order.NewService,
    )

    // 接口 Provider
    InterfaceProviderSet = wire.NewSet(
        http.NewRouter,
        grpc.NewServer,
    )
)

// 初始化应用
func InitializeApp(cfgPath string) (*App, error) {
    wire.Build(
        ConfigProviderSet,
        InfrastructureProviderSet,
        RepositoryProviderSet,
        ServiceProviderSet,
        InterfaceProviderSet,
        NewApp,
    )
    return nil, nil
}

// App 结构
type App struct {
    Config    *config.Config
    Database  *database.Database
    Cache     *cache.Cache
    MQ        *messaging.MessageQueue
    Router    *http.Router
    GRPCServer *grpc.Server
}

// NewApp 创建应用
func NewApp(
    cfg *config.Config,
    db *database.Database,
    c *cache.Cache,
    mq *messaging.MessageQueue,
    router *http.Router,
    grpcServer *grpc.Server,
) *App {
    return &App{
        Config:     cfg,
        Database:   db,
        Cache:      c,
        MQ:         mq,
        Router:     router,
        GRPCServer: grpcServer,
    }
}
```

### 3.2 多环境配置

```go
// 环境类型
type Environment string

const (
    EnvDevelopment Environment = "development"
    EnvStaging     Environment = "staging"
    EnvProduction  Environment = "production"
)

// 根据环境创建不同的 Provider
func NewConfig(env Environment) (*Config, error) {
    switch env {
    case EnvDevelopment:
        return &Config{
            DatabaseURL: "postgres://localhost/dev_db",
            LogLevel:    "debug",
        }, nil
    case EnvStaging:
        return &Config{
            DatabaseURL: os.Getenv("DATABASE_URL"),
            LogLevel:    "info",
        }, nil
    case EnvProduction:
        return &Config{
            DatabaseURL: os.Getenv("DATABASE_URL"),
            LogLevel:    "warn",
        }, nil
    default:
        return nil, fmt.Errorf("unknown environment: %s", env)
    }
}

// 初始化应用（传入环境）
func InitializeApp(env Environment) (*App, error) {
    wire.Build(
        NewConfig,
        NewDatabase,
        NewApp,
    )
    return nil, nil
}
```

### 3.3 条件依赖

```go
// 使用接口支持条件依赖
type Logger interface {
    Log(msg string)
}

// 实现 1：控制台日志
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
    return &ConsoleLogger{}
}

func (l *ConsoleLogger) Log(msg string) {
    fmt.Println(msg)
}

// 实现 2：文件日志
type FileLogger struct {
    file *os.File
}

func NewFileLogger(path string) (*FileLogger, error) {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &FileLogger{file: file}, nil
}

func (l *FileLogger) Log(msg string) {
    l.file.WriteString(msg + "\n")
}

// 根据配置选择不同的实现
func NewLogger(cfg *Config) (Logger, error) {
    if cfg.LogToFile {
        return NewFileLogger(cfg.LogPath)
    }
    return NewConsoleLogger(), nil
}

// 使用接口绑定
var LoggerProviderSet = wire.NewSet(
    NewLogger,
    wire.Bind(new(Logger), new(*ConsoleLogger)),
    wire.Bind(new(Logger), new(*FileLogger)),
)
```

---

## 4. 常见模式

### 4.1 单例模式

```go
// Wire 默认创建单例（每个类型只创建一次）
func NewDatabase(cfg *Config) (*Database, error) {
    // 这个函数只会被调用一次
    return &Database{URL: cfg.DatabaseURL}, nil
}

// 多个 Provider 依赖同一个 Database
func NewUserRepository(db *Database) UserRepository {
    return &UserRepositoryImpl{db: db}  // 使用同一个 db 实例
}

func NewOrderRepository(db *Database) OrderRepository {
    return &OrderRepositoryImpl{db: db}  // 使用同一个 db 实例
}
```

### 4.2 工厂模式

```go
// 使用 Provider 函数作为工厂
func NewRepositoryFactory(db *Database) *RepositoryFactory {
    return &RepositoryFactory{db: db}
}

type RepositoryFactory struct {
    db *Database
}

func (f *RepositoryFactory) CreateUserRepository() UserRepository {
    return &UserRepositoryImpl{db: f.db}
}

func (f *RepositoryFactory) CreateOrderRepository() OrderRepository {
    return &OrderRepositoryImpl{db: f.db}
}
```

### 4.3 装饰器模式

```go
// 基础 Repository
func NewUserRepository(db *Database) UserRepository {
    return &UserRepositoryImpl{db: db}
}

// 带缓存的 Repository（装饰器）
func NewCachedUserRepository(
    repo UserRepository,
    cache *Cache,
) UserRepository {
    return &CachedUserRepository{
        repo:  repo,
        cache: cache,
    }
}

// 使用接口绑定
var RepositoryProviderSet = wire.NewSet(
    NewUserRepository,
    NewCachedUserRepository,
    wire.Bind(new(UserRepository), new(*CachedUserRepository)),
)
```

---

## 5. 故障排查

### 5.1 常见错误

#### 错误 1: 缺少 Provider

```
wire: no provider found for *database.Database
```

**原因**：没有为 `*database.Database` 类型提供 Provider。

**解决方案**：

```go
// 添加 Provider
func NewDatabase(cfg *Config) (*Database, error) {
    return &Database{}, nil
}

// 在 wire.Build 中包含
wire.Build(
    NewConfig,
    NewDatabase,  // 添加这个
    NewApp,
)
```

#### 错误 2: 循环依赖

```
wire: cycle detected in dependency graph
```

**原因**：Provider 之间存在循环依赖。

**解决方案**：

1. 重新设计依赖关系
2. 使用接口解耦
3. 延迟初始化

#### 错误 3: 类型不匹配

```
wire: *database.UserRepositoryImpl is not assignable to user.Repository
```

**原因**：类型不匹配，需要接口绑定。

**解决方案**：

```go
var RepositoryProviderSet = wire.NewSet(
    NewUserRepository,
    wire.Bind(new(user.Repository), new(*database.UserRepositoryImpl)),
)
```

### 5.2 调试技巧

#### 技巧 1: 查看生成的代码

```bash
# 生成代码后查看
cat wire_gen.go
```

#### 技巧 2: 使用 Wire 的调试选项

```bash
# 显示详细的依赖图
wire -v ./scripts/wire

# 显示 Provider 信息
wire -show-providers ./scripts/wire
```

#### 技巧 3: 逐步添加 Provider

```go
// 先添加基础 Provider
wire.Build(
    NewConfig,
    NewDatabase,
)

// 逐步添加更多 Provider
wire.Build(
    NewConfig,
    NewDatabase,
    NewUserRepository,
    NewUserService,
)
```

---

## 6. 性能优化

### 6.1 减少 Provider 函数调用

```go
// ❌ 错误：每次调用都创建新对象
func NewConfig() *Config {
    return &Config{...}  // 每次都创建新对象
}

// ✅ 正确：使用单例或缓存
var configInstance *Config
var configOnce sync.Once

func NewConfig() *Config {
    configOnce.Do(func() {
        configInstance = &Config{...}
    })
    return configInstance
}
```

### 6.2 延迟初始化

```go
// 使用 lazy initialization
type Service struct {
    repo UserRepository
    cache *Cache
}

func NewService(repo UserRepository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetCache() *Cache {
    if s.cache == nil {
        s.cache = NewCache()  // 延迟初始化
    }
    return s.cache
}
```

---

**最后更新**: 2025-01-XX
