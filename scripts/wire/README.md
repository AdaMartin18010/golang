# Wire 依赖注入

**版本**: v2.0
**更新日期**: 2025-12-03
**Wire 版本**: v0.6.0

---

## 🎯 Wire 配置结构

### 文件组织

```text
scripts/wire/
├── wire.go           # Wire 注入配置 ✅
├── providers.go      # Provider 函数集合 ✅
├── wire_gen.go       # 生成的代码（自动生成）
└── README.md         # 本文档
```

---

## 🏗️ Provider 组织

### Provider Sets 按层次组织

```go
// 可观测性 Provider
ObservabilityProviderSet = wire.NewSet(
    NewOTLPIntegration,
    NewSystemMonitor,
    NewPlatformMonitor,
)

// 安全 Provider
SecurityProviderSet = wire.NewSet(
    NewJWTTokenManager,
    NewRBACSystem,
)

// 数据库 Provider
DatabaseProviderSet = wire.NewSet(
    NewEntClient,
    NewUserRepository,
)

// 应用层 Provider
ApplicationProviderSet = wire.NewSet(
    NewUserService,
)

// 接口层 Provider
InterfaceProviderSet = wire.NewSet(
    NewRouter,
)
```

### 依赖关系

```text
Config
  ↓
ObservabilityProviderSet + SecurityProviderSet
  ↓
DatabaseProviderSet
  ↓
ApplicationProviderSet
  ↓
InterfaceProviderSet
  ↓
App
```

---

## 🚀 使用方法

### 1. 生成代码

```bash
# 方法1: 使用 Makefile
make generate-wire

# 方法2: 直接使用 go generate
cd scripts/wire && go generate

# 方法3: 直接运行 wire
wire ./scripts/wire
```

### 2. 在应用中使用

```go
package main

import (
    "log"
    "github.com/yourusername/golang/scripts/wire"
    "github.com/yourusername/golang/internal/config"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 使用 Wire 初始化应用
    app, err := wire.InitializeApp(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 运行应用
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 📚 Provider 函数规范

### 1. 命名规范

```go
// ✅ 正确：New + 类型名
func NewUserRepository(client *ent.Client) *UserRepository

// ❌ 错误：不清晰的命名
func GetUserRepo(c *ent.Client) *UserRepository
```

### 2. 返回值规范

```go
// ✅ 推荐：返回具体类型 + error
func NewService(repo Repository) (*Service, error)

// ✅ 可接受：返回具体类型
func NewConfig() *Config

// ⚠️ 慎用：返回接口（需要 wire.Bind）
func NewService(repo Repository) Service
```

### 3. 接口绑定

```go
// 当 Provider 返回具体类型，但需要接口时
var ProviderSet = wire.NewSet(
    NewUserRepository,  // 返回 *UserRepositoryImpl
    wire.Bind(new(UserRepository), new(*UserRepositoryImpl)),
)
```

---

## 🎯 最佳实践

### 1. 按层次组织

- ✅ 每层一个 ProviderSet
- ✅ 清晰的依赖关系
- ✅ 易于维护和测试

### 2. 错误处理

```go
// ✅ Provider 应该返回 error
func NewDatabase(cfg *Config) (*Database, error) {
    db, err := sql.Open(...)
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

### 3. 清理函数

```go
// ✅ 返回清理函数
func NewDatabase(cfg *Config) (*Database, func(), error) {
    db, err := sql.Open(...)
    if err != nil {
        return nil, nil, err
    }

    cleanup := func() {
        db.Close()
    }

    return db, cleanup, nil
}
```

---

## 🔧 常见问题

### Q1: Wire 生成失败？

```bash
# 检查语法
go build ./scripts/wire

# 查看详细错误
wire -v ./scripts/wire
```

### Q2: 循环依赖？

- 检查 Provider 函数的参数
- 使用接口打破循环
- 重新设计依赖关系

### Q3: 如何添加新的 Provider？

1. 在对应包中定义 Provider 函数
2. 在 `providers.go` 中添加到对应的 ProviderSet
3. 运行 `make generate-wire`

---

**状态**: ✅ 完整实现
**下一步**: 根据业务需求添加更多 Provider
