# 🚀 Go 1.25.3 Workspace 快速开始

> 5分钟快速了解和使用 Go Workspace 模式

## 📋 目录

- [什么是 Workspace](#什么是-workspace)
- [为什么要用 Workspace](#为什么要用-workspace)
- [快速开始](#快速开始)
- [常用命令](#常用命令)
- [常见问题](#常见问题)

---

## 什么是 Workspace？

Go Workspace（工作区）是 **Go 1.18+** 引入的特性，允许在一个项目中同时开发多个模块。

```text
传统模式（单模块）          Workspace 模式（多模块）
───────────────            ─────────────────────
go.mod                     go.work          ← 工作区配置
main.go                    ├── go.work.sum
service/                   ├── pkg/
                           │   ├── agent/
                           │   │   └── go.mod
                           │   └── concurrency/
                           │       └── go.mod
                           ├── cmd/
                           │   └── app/
                           │       └── go.mod
                           └── examples/
                               └── go.mod
```

---

## 为什么要用 Workspace？

### 1. 本地开发便利 🎯

**传统方式**：

```bash
# 修改 pkg/agent 后，需要手动替换
cd examples
go mod edit -replace example.com/agent=../pkg/agent
go mod tidy
```

**Workspace 方式**：

```bash
# 修改 pkg/agent 后，自动使用最新代码
cd examples
go run .    # 自动使用 workspace 中的本地 pkg/agent
```

### 2. 统一依赖管理 📦

```bash
# 一个命令同步所有模块的依赖
go work sync

# 测试所有模块
go work test ./...

# 构建所有程序
go work build ./cmd/...
```

### 3. 清晰的模块边界 🏗️

```text
每个功能都是独立模块：
- 可以单独测试
- 可以独立发布
- 依赖关系清晰
```

---

## 快速开始

### Step 1: 检查 Go 版本

```bash
go version
# 需要 Go 1.18 或更高版本
# 推荐 Go 1.25.3
```

### Step 2: 创建 `go.work` 文件

**方法 A: 手动创建**:

```go
// go.work
go 1.25.3

use (
    ./pkg/agent
    ./pkg/concurrency
    ./examples
    ./cmd/app
)
```

**方法 B: 使用命令创建**:

```bash
# 初始化 workspace
go work init

# 添加模块
go work use ./pkg/agent
go work use ./pkg/concurrency
go work use ./examples

# 或自动发现所有模块
go work use ./...
```

### Step 3: 同步依赖

```bash
go work sync
```

### Step 4: 验证

```bash
# 测试所有模块
go work test ./...

# 查看 workspace 配置
cat go.work
```

---

## 常用命令

### 🔧 Workspace 管理

| 命令 | 说明 | 示例 |
|-----|------|------|
| `go work init` | 初始化 workspace | `go work init ./examples` |
| `go work use` | 添加模块 | `go work use ./pkg/agent` |
| `go work use ./...` | 自动发现所有模块 | - |
| `go work edit` | 编辑 go.work | `go work edit -replace a=b` |
| `go work sync` | 同步所有模块的依赖 | - |

### 📦 模块操作

| 命令 | 说明 | 示例 |
|-----|------|------|
| `go work test ./...` | 测试所有模块 | - |
| `go work build ./...` | 构建所有模块 | - |
| `go work vet ./...` | 静态检查所有模块 | - |
| `go work fmt ./...` | 格式化所有模块 | - |

### 🔍 调试

```bash
# 查看某个包的实际路径
go list -f '{{.Dir}}' example.com/agent

# 查看所有模块
go list -m all

# 查看依赖图
go mod graph
```

---

## 本项目的 Workspace 结构

```text
golang/
├── go.work                    # ← Workspace 配置文件
├── go.work.sum               # ← 校验和
│
├── pkg/                      # 可复用库（多个模块）
│   ├── agent/
│   │   ├── go.mod           # 独立模块
│   │   └── core/
│   ├── concurrency/
│   │   ├── go.mod           # 独立模块
│   │   └── pipeline/
│   └── http3/
│       ├── go.mod           # 独立模块
│       └── server/
│
├── cmd/                      # 可执行程序（多个模块）
│   ├── ai-agent/
│   │   ├── go.mod           # 独立模块
│   │   └── main.go
│   └── http3-server/
│       ├── go.mod           # 独立模块
│       └── main.go
│
└── examples/                 # 示例代码（单个模块）
    ├── go.mod               # 统一的 examples 模块
    ├── 01-basic/
    ├── 02-concurrency/
    └── 03-web-development/
```

---

## 实际使用示例

### 场景 1: 开发新功能

```bash
# 1. 修改库代码
cd pkg/agent
# 编辑 agent.go

# 2. 在示例中测试（自动使用本地代码）
cd ../../examples/05-ai-agent
go run .    # ← 自动使用 workspace 中的最新代码

# 3. 运行测试
cd ../../
go work test ./pkg/agent ./examples/05-ai-agent
```

### 场景 2: 更新依赖

```bash
# 更新特定模块的依赖
cd pkg/agent
go get -u github.com/gin-gonic/gin
go mod tidy

# 同步到其他模块
cd ../..
go work sync

# 验证
go work test ./...
```

### 场景 3: 添加新模块

```bash
# 1. 创建新模块
mkdir -p pkg/newfeature
cd pkg/newfeature
go mod init github.com/yourusername/newfeature

# 2. 添加到 workspace
cd ../..
go work use ./pkg/newfeature

# 3. 验证
go work sync
```

---

## 常见问题

### Q1: go.work 应该提交到 Git 吗？

**A**: 取决于项目类型：

| 项目类型 | 是否提交 | 原因 |
|---------|---------|------|
| **应用程序** | ✅ 是 | 团队共享统一的开发环境 |
| **库/SDK** | ❌ 否 | 用户可能有不同的 workspace 配置 |
| **Monorepo** | ✅ 是 | 统一管理多个相关项目 |

**建议**：

- 提交 `go.work`
- 在 `.gitignore` 中添加 `go.work.sum`

### Q2: Workspace vs go.mod replace？

| 特性 | Workspace | go.mod replace |
|-----|-----------|----------------|
| **用途** | 本地开发 | 临时替换或私有仓库 |
| **作用域** | 所有模块 | 单个模块 |
| **提交** | 可选 | 通常提交 |
| **优先级** | 更高 | 较低 |

**建议**：优先使用 Workspace，replace 用于特殊情况

### Q3: 如何禁用 Workspace？

```bash
# 临时禁用
GOWORK=off go test ./...

# 或重命名文件
mv go.work go.work.bak

# 恢复
mv go.work.bak go.work
```

### Q4: Workspace 影响性能吗？

**A**: 不会。Workspace 只影响模块解析，不影响运行时性能。

### Q5: 多人协作时如何使用？

**团队规范**：

```bash
# 每个人在自己的分支开发
git checkout -b feature/my-feature

# 定期同步主分支
git pull origin main
go work sync

# 提交前测试所有模块
go work test ./...
```

---

## 🎯 最佳实践

### 1. 模块粒度

```text
✅ 推荐：按功能划分
pkg/
  ├── agent/        # AI Agent 功能
  ├── concurrency/  # 并发工具
  └── http3/        # HTTP/3 服务

❌ 避免：过度拆分
pkg/
  ├── agent-core/
  ├── agent-types/
  ├── agent-utils/  # 太细碎了
```

### 2. 依赖方向

```text
✅ 正确的依赖方向
cmd/app → pkg/agent → pkg/concurrency
         ↓
      examples

❌ 避免循环依赖
pkg/agent ← → pkg/concurrency
```

### 3. 版本管理

```go
// pkg/agent/go.mod
module github.com/yourusername/agent

go 1.25.3  // ← 明确指定版本

require (
    github.com/gin-gonic/gin v1.11.0  // ← 锁定大版本
)
```

### 4. 测试策略

```bash
# 单元测试：在模块目录下
cd pkg/agent
go test ./...

# 集成测试：在 workspace 根目录
go work test ./pkg/agent ./examples/05-ai-agent

# 全量测试：测试所有模块
go work test ./...
```

---

## 📚 进阶阅读

### 官方文档

- [Go Workspace Tutorial](https://go.dev/doc/tutorial/workspaces)
- [Go Modules Reference](https://go.dev/ref/mod)

### 相关文章

- [Workspace 设计文档](https://go.googlesource.com/proposal/+/master/design/45713-workspace.md)
- [Go 1.18 Release Notes](https://go.dev/doc/go1.18)

---

## 🎊 总结

**Workspace 模式的核心优势**：

1. ✅ **开发便利** - 本地修改立即生效
2. ✅ **统一管理** - 一个命令管理所有模块
3. ✅ **模块独立** - 清晰的边界和职责
4. ✅ **团队协作** - 共享一致的开发环境

**何时使用 Workspace**：

- ✅ 单仓库多模块（Monorepo）
- ✅ 同时开发相关的多个库
- ✅ 需要频繁修改依赖库
- ✅ 团队开发标准化

**何时不用 Workspace**：

- ❌ 单模块项目（不需要）
- ❌ 只使用外部依赖（不需要）
- ❌ 发布独立的库（用户不需要）

---

## 🚀 下一步

1. 阅读 [完整重构方案](RESTRUCTURE_PROPOSAL_GO1.25.3.md)
2. 运行 [迁移脚本](scripts/migrate-to-workspace.ps1)
3. 查看 [项目新结构](#本项目的-workspace-结构)
4. 开始开发！🎉

---

**Last Updated**: 2025-10-22  
**Go Version**: 1.25.3  
**Document Version**: 1.0
