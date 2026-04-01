# Wire 依赖注入替代方案评估

## 背景

Google Wire 已于 2025 年 8 月 25 日归档（archived），当前版本 v0.7.0 将是最终版本。需要评估替代方案以确保项目长期可维护性。

## 当前 Wire 使用情况

### 使用范围
- **位置**: `scripts/wire/wire.go`
- **用途**: 初始化应用，注入 Router、Service、Repository 等依赖
- **复杂度**: 低（仅初始化 App 结构体）

### 当前依赖关系
```
Config
  ↓
postgres.NewUserRepository (Infrastructure)
  ↓
appuser.NewService (Application)
  ↓
chirouter.NewRouter (Interfaces)
  ↓
NewApp (组装)
```

## 可选方案对比

| 方案 | 优点 | 缺点 | 工作量 | 推荐度 |
|------|------|------|--------|--------|
| **保留 Wire** | 无需改动；代码生成性能优秀 | 项目已归档，无未来更新 | 无 | ⭐⭐⭐ |
| **手动注入** | 零依赖；最透明；Go 惯用法 | 需要编写工厂代码 | 中等 | ⭐⭐⭐⭐⭐ |
| **uber-go/dig** | 运行时 DI；API 简洁 | 运行时反射；依赖第三方 | 中等 | ⭐⭐⭐ |
| **samber/do** | 支持泛型；轻量级 | 运行时 DI；社区较小 | 中等 | ⭐⭐⭐ |

## 推荐方案：手动依赖注入

### 理由
1. **简单场景**：当前依赖关系简单，只有 4 个组件需要组装
2. **Go 惯用法**：Go 社区推荐在简单场景下使用手动注入
3. **零依赖**：移除 Wire 依赖，减少外部依赖风险
4. **透明可调试**：代码直观，易于理解和调试

### 实施计划

#### Phase 1: 创建手动注入代码
创建 `cmd/server/main.go` 的手动注入版本：

```go
// 手动依赖注入示例
func initializeApp(cfg *config.Config) (*app.App, error) {
    // Infrastructure Layer
    repo := postgres.NewUserRepository(cfg.Database)
    
    // Application Layer
    userService := appuser.NewService(repo)
    
    // Interfaces Layer
    router := chirouter.NewRouter(userService)
    
    // App Assembly
    return &app.App{Router: router}, nil
}
```

#### Phase 2: 逐步替换
1. 在 `cmd/server/main.go` 中实现手动注入
2. 验证功能正常
3. 删除 `scripts/wire/` 目录
4. 移除 `go.mod` 中的 Wire 依赖

#### Phase 3: 清理
- 更新 Makefile（移除 wire 生成命令）
- 更新文档

## 保留 Wire 的理由

如果选择**保留 Wire**：
1. 项目功能完整，无已知 Bug
2. 编译时 DI 性能优于运行时 DI
3. 依赖关系不复杂，代码生成结果稳定
4. 归档不代表不可用，类似 `github.com/google/gops`

### 风险缓解
- 固定版本 `v0.7.0`，不自动升级
- 监控社区 fork（如 `github.com/go-wire/wire`）

## 决策建议

| 场景 | 建议 |
|------|------|
| **长期维护（>2年）** | 迁移到手动注入 |
| **短期项目（<1年）** | 保留 Wire v0.7.0 |
| **团队熟悉 Wire** | 保留 Wire，制定迁移计划 |
| **团队不熟悉 Wire** | 迁移到手动注入 |

## 当前项目建议

**建议：保留 Wire，但制定迁移计划**

理由：
1. 当前项目依赖关系简单，Wire 工作稳定
2. 迁移到手动注入需要修改所有 cmd 下的 main.go
3. 可以等到下次架构调整时一并处理
4. 在文档中标记 Wire 已归档，提醒团队

### 立即执行的行动
1. 固定 Wire 版本（已固定 v0.7.0）
2. 在 README 中添加说明
3. 创建迁移计划文档（未来执行）
