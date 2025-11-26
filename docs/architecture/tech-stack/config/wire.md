# 1. ⚙️ Wire 依赖注入深度解析

> **简介**: 本文档详细阐述了 Wire 依赖注入的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. ⚙️ Wire 依赖注入深度解析](#1-️-wire-依赖注入深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 Provider 函数编写](#131-provider-函数编写)
    - [1.3.2 代码生成](#132-代码生成)
    - [1.3.3 使用生成的代码](#133-使用生成的代码)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 Provider 设计最佳实践](#141-provider-设计最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Wire 是什么？**

Wire 是 Google 开源的 Go 语言依赖注入工具。

**核心特性**:

- ✅ **编译时注入**: 编译时生成代码
- ✅ **类型安全**: 编译时检查
- ✅ **零反射**: 运行时零反射
- ✅ **性能**: 性能优秀

---

## 1.2 选型论证

**为什么选择 Wire？**

**论证矩阵**:

| 评估维度 | 权重 | Wire | Fx | Dig | 手动注入 | 说明 |
|---------|------|------|-----|-----|---------|------|
| **编译时检查** | 35% | 10 | 5 | 3 | 8 | Wire 编译时检查 |
| **性能** | 25% | 10 | 7 | 6 | 10 | Wire 零反射，性能优秀 |
| **易用性** | 20% | 8 | 9 | 7 | 6 | Wire 使用简单 |
| **类型安全** | 15% | 10 | 8 | 6 | 9 | Wire 类型安全 |
| **生态支持** | 5% | 8 | 7 | 6 | 10 | Wire 生态良好 |
| **加权总分** | - | **9.30** | 7.20 | 5.55 | 8.50 | Wire 得分最高 |

**核心优势**:

1. **编译时检查（权重 35%）**:
   - 编译时生成代码，编译时检查依赖
   - 减少运行时错误，提高可靠性
   - IDE 支持好，错误提示清晰

2. **性能（权重 25%）**:
   - 零反射，性能优秀
   - 生成的代码高效
   - 适合高性能场景

3. **类型安全（权重 15%）**:
   - 编译时类型检查
   - 避免运行时类型错误
   - 与 Go 类型系统完美集成

**为什么不选择其他依赖注入方案？**

1. **Fx**:
   - ✅ 功能丰富，支持生命周期管理
   - ❌ 运行时注入，性能不如 Wire
   - ❌ 编译时检查不如 Wire
   - ❌ 学习成本较高

2. **Dig**:
   - ✅ 功能完整，支持复杂场景
   - ❌ 运行时注入，性能不如 Wire
   - ❌ 编译时检查不如 Wire
   - ❌ API 较复杂

3. **手动注入**:
   - ✅ 完全控制，性能最优
   - ❌ 代码量大，维护成本高
   - ❌ 容易出错，无编译时检查
   - ❌ 不适合大型项目

---

## 1.3 实际应用

### 1.3.1 Provider 函数编写

**Provider 函数示例**:

```go
// internal/wire/wire.go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "github.com/yourusername/golang/internal/application/user"
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
    entrepo "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
)

// ProviderSet 定义 Provider 集合
var ProviderSet = wire.NewSet(
    // 数据库
    ent.NewClient,

    // 仓储
    entrepo.NewUserRepository,

    // 服务
    user.NewService,
)

// InitializeApp 初始化应用
func InitializeApp() (*App, error) {
    wire.Build(
        ProviderSet,
        NewApp,
    )
    return &App{}, nil
}
```

### 1.3.2 代码生成

**生成代码**:

```bash
# 生成 Wire 代码
go generate ./internal/wire
```

### 1.3.3 使用生成的代码

**使用示例**:

```go
// cmd/server/main.go
package main

import (
    "github.com/yourusername/golang/internal/wire"
)

func main() {
    app, err := wire.InitializeApp()
    if err != nil {
        log.Fatal(err)
    }

    app.Run()
}
```

---

## 1.4 最佳实践

### 1.4.1 Provider 设计最佳实践

**为什么需要良好的 Provider 设计？**

良好的 Provider 设计可以提高依赖注入的可维护性和可测试性。

**Provider 设计原则**:

1. **单一职责**: 每个 Provider 只负责一个依赖
2. **接口绑定**: 使用接口绑定，提高灵活性
3. **错误处理**: 正确处理 Provider 错误
4. **测试支持**: 支持测试场景

**实际应用示例**:

```go
// Provider 设计最佳实践
// 定义接口
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Get(ctx context.Context, id string) (*User, error)
}

// Provider 函数
func NewUserRepository(client *ent.Client) UserRepository {
    return entrepo.NewUserRepository(client)
}

// 使用接口绑定
var ProviderSet = wire.NewSet(
    ent.NewClient,
    wire.Bind(new(UserRepository), new(*entrepo.UserRepository)),
    user.NewService,
)
```

**最佳实践要点**:

1. **单一职责**: 每个 Provider 只负责一个依赖
2. **接口绑定**: 使用接口绑定，提高灵活性
3. **错误处理**: 正确处理 Provider 错误
4. **测试支持**: 支持测试场景，便于单元测试

---

## 📚 扩展阅读

- [Wire 官方文档](https://github.com/google/wire)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Wire 依赖注入的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
