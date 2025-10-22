# Workspace 迁移方案

> **制定日期**: 2025-10-22  
> **执行阶段**: Phase 2 - Day 1  
> **目标**: 建立标准Go项目结构

---

## 📋 迁移概述

### 目标

将当前的examples导向结构迁移到标准的Go项目布局：
- `pkg/` - 可复用的公共库
- `internal/` - 内部实现包
- `cmd/` - 可执行程序
- `examples/` - 示例代码

### 原则

```text
✅ 保持向后兼容 - examples/保留供学习使用
✅ 逐步迁移 - 一次一个模块，充分验证
✅ 清晰职责 - pkg/是库，examples/是示例
✅ 标准布局 - 遵循Go社区最佳实践
```

---

## 🔍 现状分析

### 当前结构

```text
examples/
├── advanced/
│   ├── ai-agent/          ← 适合迁移到 pkg/agent/
│   ├── http3/             ← 适合迁移到 pkg/http3/
│   └── weak-pointer-cache/ ← 适合迁移到 pkg/memory/
├── concurrency/
│   └── patterns/          ← 适合迁移到 pkg/concurrency/
├── modern-features/       ← 保留为示例
├── observability/         ← 适合迁移到 pkg/observability/
└── testing/               ← 保留为示例
```

### 模块评估

| 模块路径 | 是否迁移 | 目标位置 | 理由 |
|----------|---------|----------|------|
| `examples/advanced/ai-agent/` | ✅ | `pkg/agent/` | 完整的库，可复用 |
| `examples/advanced/http3/` | ✅ | `pkg/http3/` | HTTP/3实现库 |
| `examples/advanced/weak-pointer-cache/` | ✅ | `pkg/memory/` | 内存管理库 |
| `examples/concurrency/patterns/` | ✅ | `pkg/concurrency/` | 并发模式库 |
| `examples/observability/` | ✅ | `pkg/observability/` | 可观测性库 |
| `examples/modern-features/` | ❌ | 保留 | 教学示例 |
| `examples/testing/` | ❌ | 保留 | 测试示例 |
| `examples/basic/` | ❌ | 保留 | 基础示例 |

---

## 📦 目标结构设计

### 完整目录树

```text
golang/
│
├── cmd/                          # 可执行程序
│   └── gox/                      # CLI工具 (已有)
│       ├── main.go
│       ├── go.mod
│       └── README.md
│
├── pkg/                          # 公共库 (新增)
│   │
│   ├── agent/                    # AI Agent库
│   │   ├── core/                 # 核心实现
│   │   │   ├── agent.go         # BaseAgent
│   │   │   ├── decision_engine.go
│   │   │   ├── learning_engine.go
│   │   │   └── multimodal_interface.go
│   │   ├── examples/             # 库的使用示例
│   │   ├── docs/                 # 库文档
│   │   │   ├── ARCHITECTURE.md
│   │   │   └── QUICK_START.md
│   │   ├── agent_test.go         # 测试
│   │   ├── go.mod
│   │   └── README.md
│   │
│   ├── concurrency/              # 并发库
│   │   ├── patterns/             # 并发模式
│   │   │   ├── pipeline.go
│   │   │   ├── worker_pool.go
│   │   │   └── fan_out_fan_in.go
│   │   ├── go.mod
│   │   └── README.md
│   │
│   ├── http3/                    # HTTP/3库
│   │   ├── client.go
│   │   ├── server.go
│   │   ├── go.mod
│   │   └── README.md
│   │
│   ├── memory/                   # 内存管理库
│   │   ├── weak_pointer.go
│   │   ├── arena.go
│   │   ├── go.mod
│   │   └── README.md
│   │
│   └── observability/            # 可观测性库
│       ├── metrics.go
│       ├── tracing.go
│       ├── go.mod
│       └── README.md
│
├── internal/                     # 内部包 (新增)
│   ├── types/                    # 内部类型
│   │   └── common.go
│   └── utils/                    # 内部工具
│       └── helpers.go
│
├── examples/                     # 示例代码 (保留)
│   ├── basic/                    # 基础示例
│   ├── concurrency/              # 并发示例
│   ├── modern-features/          # 新特性示例
│   ├── testing/                  # 测试示例
│   ├── go.mod
│   └── README.md
│
├── docs/                         # 项目文档
│   ├── guides/
│   ├── tutorials/
│   └── INDEX.md
│
├── scripts/                      # 工具脚本
│
├── .github/                      # GitHub配置 (新增)
│   └── workflows/
│       ├── ci.yml               # CI工作流
│       └── docs.yml             # 文档部署
│
├── go.work                       # Workspace配置
├── README.md
└── 📖-README-项目导航.md
```

### import路径规划

```text
迁移后的import路径:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
github.com/yourusername/golang/pkg/agent
github.com/yourusername/golang/pkg/concurrency
github.com/yourusername/golang/pkg/http3
github.com/yourusername/golang/pkg/memory
github.com/yourusername/golang/pkg/observability

github.com/yourusername/golang/internal/types
github.com/yourusername/golang/internal/utils

github.com/yourusername/golang/cmd/gox
```

---

## 🚀 迁移步骤

### Phase 1: 准备阶段 (Day 1)

**Step 1: 创建目录结构**

```bash
# 创建pkg/目录
mkdir pkg

# 创建各模块目录
mkdir -p pkg/agent/core
mkdir -p pkg/agent/docs
mkdir -p pkg/agent/examples

mkdir -p pkg/concurrency/patterns
mkdir -p pkg/http3
mkdir -p pkg/memory
mkdir -p pkg/observability

# 创建internal/目录
mkdir -p internal/types
mkdir -p internal/utils

# 创建GitHub配置目录
mkdir -p .github/workflows
```

**Step 2: 迁移第一个模块 (agent)**

```bash
# 复制agent代码
cp -r examples/advanced/ai-agent/core/* pkg/agent/core/
cp -r examples/advanced/ai-agent/docs/* pkg/agent/docs/
cp -r examples/advanced/ai-agent/*.go pkg/agent/

# 复制测试文件
cp examples/advanced/ai-agent/*_test.go pkg/agent/
```

**Step 3: 创建go.mod**

```bash
cd pkg/agent
go mod init github.com/yourusername/golang/pkg/agent

# 设置Go版本
go mod edit -go=1.25

# 添加依赖
go mod tidy
```

**Step 4: 更新package声明**

```go
// 在所有.go文件中
// 从: package main 或 package ai-agent
// 改为: package agent (对于core/外的文件)
//      package core (对于core/内的文件)
```

**Step 5: 编译验证**

```bash
cd pkg/agent
go build ./...
go test ./...
```

### Phase 2: 更新Workspace配置 (Day 1-2)

**Step 6: 更新go.work**

```go
go 1.25.3

use (
    ./cmd/gox
    ./examples
    
    // 新增的pkg模块
    ./pkg/agent
    ./pkg/concurrency
    ./pkg/http3
    ./pkg/memory
    ./pkg/observability
)
```

**Step 7: 同步Workspace**

```bash
go work sync
```

### Phase 3: 迁移其他模块 (Day 2-3)

**Step 8: 迁移concurrency模块**

```bash
# 创建模块
mkdir -p pkg/concurrency/patterns
cp -r examples/concurrency/patterns/*.go pkg/concurrency/patterns/

cd pkg/concurrency
go mod init github.com/yourusername/golang/pkg/concurrency
go mod edit -go=1.25
go build ./...
go test ./...
```

**Step 9: 迁移http3模块**

```bash
mkdir -p pkg/http3
cp -r examples/advanced/http3/*.go pkg/http3/

cd pkg/http3
go mod init github.com/yourusername/golang/pkg/http3
go mod edit -go=1.25
go build ./...
go test ./...
```

**Step 10: 迁移memory模块**

```bash
mkdir -p pkg/memory
cp -r examples/advanced/weak-pointer-cache/*.go pkg/memory/
cp -r examples/modern-features/memory/*.go pkg/memory/

cd pkg/memory
go mod init github.com/yourusername/golang/pkg/memory
go mod edit -go=1.25
go build ./...
go test ./...
```

**Step 11: 迁移observability模块**

```bash
mkdir -p pkg/observability
cp -r examples/observability/*.go pkg/observability/

cd pkg/observability
go mod init github.com/yourusername/golang/pkg/observability
go mod edit -go=1.25
go build ./...
go test ./...
```

### Phase 4: 全面验证 (Day 3-4)

**Step 12: 更新examples引用**

```go
// examples中使用pkg的示例
import (
    "github.com/yourusername/golang/pkg/agent"
    "github.com/yourusername/golang/pkg/concurrency"
)
```

**Step 13: 全局编译测试**

```bash
# 在根目录
go work sync
go build ./...
go test ./...
gox quality
```

**Step 14: 更新文档**

```bash
# 更新README.md
# 更新各模块的README
# 更新import示例
# 更新📖-README-项目导航.md
```

---

## ✅ 验证清单

### 每个模块迁移后

```text
□ go.mod文件已创建
□ package声明正确
□ 代码可编译 (go build ./...)
□ 测试通过 (go test ./...)
□ 添加到go.work
□ README.md已创建
```

### 全部迁移完成后

```text
□ 所有pkg模块可编译
□ 所有tests通过
□ go work sync成功
□ gox quality通过
□ examples中的引用已更新
□ 文档已更新
□ 旧代码标记为deprecated (保留)
```

---

## 📝 每个模块的go.mod模板

```go
module github.com/yourusername/golang/pkg/MODULENAME

go 1.25

require (
    // 添加实际依赖
)
```

---

## 🎯 成功标准

### 必须达成

```text
✅ 所有模块可独立编译
✅ 所有测试通过
✅ go.work配置正确
✅ import路径更新完成
✅ 文档已更新
```

### 质量标准

```text
✅ 0编译错误
✅ 0测试失败
✅ go vet无警告
✅ 符合Go标准布局
```

---

## ⚠️ 注意事项

### 不要做的事

```text
❌ 不要删除examples/中的原始代码
❌ 不要在一次提交中迁移所有模块
❌ 不要在迁移时修改业务逻辑
❌ 不要忘记更新import路径
```

### 应该做的事

```text
✅ 每迁移一个模块就验证一次
✅ 保持examples/作为教学示例
✅ 在pkg/中添加完整文档
✅ 每个模块独立的go.mod
✅ 记录所有改动
```

---

## 📊 预期时间

```text
Day 1 (4h):  准备 + agent模块迁移
Day 2 (6h):  其他4个模块迁移
Day 3 (4h):  更新引用 + 验证
Day 4 (2h):  文档更新 + 最终验证
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
总计: 16小时 (2个工作日)
```

---

## 🔄 回退计划

如果迁移出现问题：

```bash
# 1. 使用git恢复
git checkout -- .

# 2. 清理新创建的目录
rm -rf pkg/
rm -rf internal/

# 3. 恢复go.work
git checkout go.work

# 4. 重新规划策略
```

---

<div align="center">

## ✅ 迁移方案制定完成

**下一步**: 开始执行迁移  
**预计完成**: Day 4

</div>

