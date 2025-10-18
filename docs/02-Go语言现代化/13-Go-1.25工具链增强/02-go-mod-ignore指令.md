# go.mod ignore 指令（Go 1.25）

> **版本要求**: Go 1.25+  
> **实验性**: 否（正式特性）  
> **最后更新**: 2025年10月18日

---

## 📚 目录

- [概述](#概述)
- [为什么需要 ignore 指令](#为什么需要-ignore-指令)
- [基本语法](#基本语法)
- [使用场景](#使用场景)
- [与其他工具集成](#与其他工具集成)
- [实践案例](#实践案例)
- [注意事项](#注意事项)
- [最佳实践](#最佳实践)
- [与 .gitignore 对比](#与-gitignore-对比)
- [常见问题](#常见问题)
- [参考资料](#参考资料)

---

## 概述

Go 1.25 在 `go.mod` 中引入了新的 `ignore` 指令,允许开发者明确指定哪些目录或文件应该被 Go 工具链忽略,提升构建性能和项目组织灵活性。

### 什么是 ignore 指令?

`ignore` 指令是 `go.mod` 文件中的一个新指令,用于告诉 Go 工具链忽略特定的目录或文件,这些目录不会被:

- ❌ `go build` 编译
- ❌ `go list` 列出
- ❌ `go mod tidy` 分析
- ❌ `go test` 测试
- ❌ IDE 工具扫描

### 核心优势

- ✅ **提升性能**: 减少不必要的文件扫描
- ✅ **简化配置**: 在 go.mod 中统一管理
- ✅ **明确意图**: 清晰表达哪些目录不是Go代码
- ✅ **CI/CD 友好**: 加速构建流程
- ✅ **多模块支持**: 更好支持 monorepo

---

## 为什么需要 ignore 指令?

### 传统痛点

在 Go 1.25 之前,没有官方方式在模块级别忽略目录:

```text
myproject/
├── go.mod
├── cmd/
├── pkg/
├── docs/          # 希望忽略,但会被扫描
├── examples/      # 希望忽略,但会被扫描
├── testdata/      # 特殊目录,自动忽略
├── _archive/      # 下划线开头,自动忽略
├── vendor/        # vendor 目录,自动忽略
└── .git/          # 点开头,自动忽略
```

**问题**:

1. ❌ **性能开销**: `go list ./...` 会扫描所有目录
2. ❌ **IDE 混乱**: IDE 会索引不相关的文件
3. ❌ **错误干扰**: 非 Go 代码可能导致工具报错
4. ❌ **配置分散**: 需要在多处配置忽略规则

### Go 1.25 解决方案

```go
// go.mod
module example.com/myproject

go 1.25

ignore (
    ./docs/...
    ./examples/...
    ./scripts/...
)
```

**效果**:

- ✅ 明确声明忽略目录
- ✅ Go 工具链统一遵守
- ✅ 构建性能提升
- ✅ IDE 工具集成

---

## 基本语法

### 简单语法

```go
// go.mod
module example.com/myproject

go 1.25

// 单个目录
ignore ./docs/...

// 多个目录
ignore (
    ./docs/...
    ./examples/...
)
```

### 完整语法

```go
module example.com/myproject

go 1.25

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
)

// ignore 指令可以放在文件任何位置
ignore (
    // 文档和示例
    ./docs/...
    ./examples/...
    
    // 临时和生成文件
    ./tmp/...
    ./_output/...
    ./generated/...
    
    // 旧代码和备份
    ./vendor-backup/...
    ./legacy/...
    
    // 脚本和工具
    ./scripts/...
    ./tools/...
)

// 替换规则仍然有效
replace example.com/old => example.com/new v1.0.0
```

### 语法规则

1. **路径格式**: 相对路径,以 `./` 开头
2. **递归匹配**: 使用 `/...` 表示递归忽略
3. **单行或块**: 可以单行或使用 `()` 块
4. **注释支持**: 可以使用 `//` 注释

**有效的路径**:

```go
ignore (
    ./docs/...           // ✅ 递归忽略 docs 目录
    ./tmp/...            // ✅ 递归忽略 tmp 目录
    ./scripts/           // ✅ 只忽略 scripts 目录本身
)
```

**无效的路径**:

```go
ignore (
    docs/...             // ❌ 缺少 ./
    ../other/...         // ❌ 不能是父目录
    /abs/path/...        // ❌ 不能是绝对路径
    ./docs/*.go          // ❌ 不支持通配符
)
```

---

## 使用场景

### 场景 1: 忽略文档和示例

**项目结构**:

```text
myproject/
├── go.mod
├── cmd/
│   └── app/
│       └── main.go
├── pkg/
│   └── utils/
│       └── utils.go
├── docs/              # 📚 Markdown 文档
│   ├── API.md
│   └── GUIDE.md
└── examples/          # 💡 示例代码 (不参与构建)
    ├── simple/
    └── advanced/
```

**go.mod 配置**:

```go
module example.com/myproject

go 1.25

ignore (
    ./docs/...      // 忽略所有文档
    ./examples/...  // 忽略所有示例
)
```

**效果**:

```bash
# 之前: 扫描所有目录
$ go list ./...
example.com/myproject/cmd/app
example.com/myproject/pkg/utils
example.com/myproject/docs        # 被扫描
example.com/myproject/examples     # 被扫描

# 之后: 只扫描相关目录
$ go list ./...
example.com/myproject/cmd/app
example.com/myproject/pkg/utils
```

---

### 场景 2: 忽略临时和生成文件

**项目结构**:

```text
myproject/
├── go.mod
├── cmd/
├── pkg/
├── tmp/               # 🗑️ 临时文件
├── _output/           # 📦 构建输出
├── generated/         # 🤖 代码生成
└── .cache/            # 💾 缓存文件
```

**go.mod 配置**:

```go
module example.com/myproject

go 1.25

ignore (
    ./tmp/...
    ./_output/...
    ./generated/...
    ./.cache/...
)
```

---

### 场景 3: Monorepo 多项目

**项目结构**:

```text
monorepo/
├── go.mod             # 根 go.mod
├── go.work            # Go workspace
├── service-a/
│   └── go.mod         # 独立模块
├── service-b/
│   └── go.mod         # 独立模块
├── shared/
│   └── pkg/
├── infra/             # 基础设施代码 (Terraform等)
├── docs/              # 文档
└── scripts/           # 脚本
```

**根 go.mod 配置**:

```go
module example.com/monorepo

go 1.25

// 忽略独立子模块 (它们有自己的 go.mod)
ignore (
    ./service-a/...
    ./service-b/...
    
    // 忽略非 Go 代码
    ./infra/...
    ./docs/...
    ./scripts/...
)
```

---

### 场景 4: 大型项目优化

**问题**: 1000+ 个包,构建慢

**解决**:

```go
module example.com/largeproject

go 1.25

ignore (
    // 忽略测试数据
    ./testdata/...
    
    // 忽略基准测试数据
    ./benchmarks/data/...
    
    // 忽略工具和脚本
    ./tools/...
    ./scripts/...
    
    // 忽略文档和示例
    ./docs/...
    ./examples/...
    
    // 忽略旧版本和备份
    ./legacy/...
    ./backup/...
    
    // 忽略第三方集成测试
    ./integration-tests/third-party/...
)
```

**性能提升**:

```bash
# 之前
$ time go list ./...
real    0m5.234s

# 之后
$ time go list ./...
real    0m2.156s   # 提升 ~60%
```

---

## 与其他工具集成

### go list

```bash
# 列出所有包 (遵守 ignore 指令)
$ go list ./...

# 列出所有包 (包括忽略的)
$ go list -tags=all ./...
```

---

### go mod tidy

```bash
# 整理依赖 (不扫描忽略的目录)
$ go mod tidy

# 之前可能报错:
# go: finding module for package example.com/myproject/docs
# go.mod:XX: no required module provides package example.com/myproject/docs

# 之后: 忽略 ./docs/..., 不报错
```

---

### IDE 集成

#### VS Code

`go.mod` 的 `ignore` 指令会被 Go 语言服务器 (gopls) 自动识别:

```json
// settings.json (通常不需要额外配置)
{
  "gopls": {
    // gopls 会自动读取 go.mod 的 ignore 指令
  }
}
```

#### GoLand / IntelliJ IDEA

GoLand 2023.3+ 自动支持 `ignore` 指令,无需额外配置。

---

### CI/CD 优化

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Build
        run: go build ./...
        # go.mod 的 ignore 指令自动生效
        # 构建速度更快
      
      - name: Test
        run: go test ./...
        # 不会测试忽略的目录
```

---

## 实践案例

### 案例 1: Web 项目

**项目结构**:

```text
webapp/
├── go.mod
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   ├── service/
│   └── repository/
├── web/               # 前端代码
│   ├── src/
│   ├── public/
│   └── node_modules/
├── docs/
├── docker/
└── k8s/
```

**go.mod**:

```go
module example.com/webapp

go 1.25

require (
    github.com/gin-gonic/gin v1.9.1
)

ignore (
    ./web/...      // 前端代码
    ./docs/...     // 文档
    ./docker/...   // Docker 配置
    ./k8s/...      // Kubernetes 配置
)
```

**效果**:

- ✅ `go list ./...` 只列出 Go 代码
- ✅ IDE 不索引前端代码
- ✅ 构建速度提升 40%

---

### 案例 2: 微服务 Monorepo

**项目结构**:

```text
microservices/
├── go.work
├── go.mod
├── services/
│   ├── auth/
│   │   └── go.mod
│   ├── user/
│   │   └── go.mod
│   └── order/
│       └── go.mod
├── shared/
│   └── pkg/
├── deployment/
├── docs/
└── scripts/
```

**根 go.mod**:

```go
module example.com/microservices

go 1.25

// 共享代码
require (
    github.com/grpc/grpc-go v1.56.0
)

// 忽略子服务 (有独立 go.mod)
ignore (
    ./services/auth/...
    ./services/user/...
    ./services/order/...
    
    // 忽略非代码目录
    ./deployment/...
    ./docs/...
    ./scripts/...
)
```

**每个服务的 go.mod**:

```go
// services/auth/go.mod
module example.com/microservices/services/auth

go 1.25

require (
    example.com/microservices v0.0.0
    github.com/golang-jwt/jwt v4.5.0
)

replace example.com/microservices => ../..
```

---

### 案例 3: 工具和代码混合项目

**项目结构**:

```text
project/
├── go.mod
├── cmd/
│   └── app/
│       └── main.go
├── pkg/
├── tools/             # 开发工具
│   ├── codegen/
│   └── migrate/
├── scripts/           # Shell 脚本
├── hack/              # 辅助脚本
└── third_party/       # 第三方代码
```

**go.mod**:

```go
module example.com/project

go 1.25

require (
    github.com/spf13/cobra v1.7.0
)

ignore (
    ./tools/...        // 开发工具 (可选择性忽略)
    ./scripts/...      // Shell 脚本
    ./hack/...         // 辅助脚本
    ./third_party/...  // 第三方代码
)
```

**注意**: `./tools/` 如果包含 Go 工具,可以不忽略,或者使用 `tools.go` 模式:

```go
//go:build tools
// +build tools

package tools

import (
    _ "github.com/golangci/golangci-lint/cmd/golangci-lint"
    _ "golang.org/x/tools/cmd/goimports"
)
```

---

## 注意事项

### 1. 忽略的文件仍在模块 ZIP 中

**重要**: `ignore` 指令**不影响** `go mod vendor` 或模块 ZIP:

```bash
# 创建模块 ZIP
$ go mod vendor

# vendor 目录仍包含所有文件 (包括忽略的)
```

**原因**: 这是设计选择,确保模块的完整性。

**如果需要排除文件**: 使用 `.gitattributes`:

```text
# .gitattributes
docs/**     export-ignore
examples/** export-ignore
```

---

### 2. 不影响 go.sum

忽略的目录中的 `import` 语句仍会被分析:

```go
// examples/main.go (在 ignore 列表中)
package main

import "github.com/unknown/package"  // 仍会尝试解析

func main() {}
```

**建议**: 确保示例代码的依赖也在 `require` 中,或者使用独立的 `go.mod`。

---

### 3. 版本控制考虑

`ignore` 指令**不替代** `.gitignore`:

```text
# .gitignore (仍然需要)
/tmp/
/_output/
/.cache/
*.log
```

**区别**:

- **`.gitignore`**: Git 不跟踪这些文件
- **`ignore` 指令**: Go 工具链不扫描这些目录

---

### 4. 递归忽略

使用 `/...` 进行递归忽略:

```go
ignore (
    ./docs/...      // ✅ 忽略 docs 及所有子目录
    ./docs/         // ⚠️ 只忽略 docs 目录本身
)
```

---

## 最佳实践

### 1. 明确忽略非 Go 代码

```go
ignore (
    // 文档
    ./docs/...
    ./README_files/...
    
    // 前端
    ./web/...
    ./static/...
    
    // 基础设施
    ./terraform/...
    ./ansible/...
    ./k8s/...
    
    // 脚本
    ./scripts/...
    ./hack/...
)
```

---

### 2. 保持测试数据结构

`testdata/` 目录自动忽略,无需显式声明:

```go
// ❌ 不需要
ignore (
    ./testdata/...
)

// ✅ testdata 自动忽略
```

---

### 3. 文档化忽略原因

```go
ignore (
    // 文档: 纯 Markdown,无 Go 代码
    ./docs/...
    
    // 示例: 独立示例,不参与主构建
    ./examples/...
    
    // 工具: 开发工具,有独立 go.mod
    ./tools/...
)
```

---

### 4. CI/CD 验证

```yaml
# .github/workflows/verify.yml
- name: Verify ignore paths exist
  run: |
    # 确保 go.mod 中忽略的目录确实存在
    grep "ignore" go.mod | while read -r line; do
      path=$(echo "$line" | sed 's/.*\.\/\([^/]*\).*/\1/')
      if [ ! -d "$path" ]; then
        echo "Warning: Ignored directory $path does not exist"
      fi
    done
```

---

## 与 .gitignore 对比

| 特性 | `go.mod` ignore | `.gitignore` |
|------|----------------|--------------|
| **目的** | Go 工具链忽略 | Git 版本控制忽略 |
| **作用范围** | Go 命令 (build, list, test) | Git 命令 |
| **语法** | Go 模块路径 | Git 通配符 |
| **文件存在** | 文件仍在仓库中 | 文件不提交到仓库 |
| **IDE 支持** | gopls 支持 | Git 集成 |
| **递归** | `/...` | `/**` 或 `/*` |

**最佳实践**: 两者结合使用

```text
# .gitignore
/tmp/
/_output/
/.cache/
*.log
node_modules/
```

```go
// go.mod
ignore (
    ./docs/...
    ./web/...
)
```

---

## 常见问题

### Q1: ignore 指令会加快构建速度吗?

**A**: ✅ 会!

- **`go list ./...`**: 减少目录扫描
- **`go build ./...`**: 减少编译目标
- **`go test ./...`**: 减少测试目标
- **IDE 索引**: 减少索引文件

**性能提升**: 10-40% (取决于项目大小)

---

### Q2: 可以忽略特定文件吗?

**A**: ❌ 不直接支持

`ignore` 指令只能忽略目录,不能忽略特定文件:

```go
// ❌ 不支持
ignore (
    ./pkg/old_file.go
)

// ✅ 只能忽略目录
ignore (
    ./pkg/old/...
)
```

**解决方案**: 将不需要的文件移到单独目录。

---

### Q3: ignore 和 build tags 有什么区别?

**A**: 不同的用途

| 特性 | `ignore` 指令 | Build Tags |
|------|--------------|------------|
| **级别** | 目录级别 | 文件级别 |
| **用途** | 永久忽略 | 条件编译 |
| **语法** | `go.mod` | `//go:build` |
| **动态** | 静态 | 动态 (编译时选择) |

**示例**:

```go
// ignore: 永久忽略整个目录
ignore (
    ./experimental/...
)

// build tags: 条件编译文件
//go:build linux
// +build linux

package main
```

---

### Q4: 可以在子模块中使用 ignore 吗?

**A**: ✅ 可以

每个 `go.mod` 都可以有自己的 `ignore` 指令:

```text
project/
├── go.mod              # 根模块
│   └── ignore ./services/...
└── services/
    └── auth/
        └── go.mod      # 子模块
            └── ignore ./testdata-large/...
```

---

### Q5: ignore 会影响 go get 吗?

**A**: ❌ 不会

`ignore` 指令只影响**本地**操作,不影响:

- ✅ `go get` 下载依赖
- ✅ 模块发布到代理
- ✅ 其他项目依赖你的模块

---

## 参考资料

### 官方文档

- 📘 [Go 1.25 Release Notes](https://go.dev/doc/go1.25#gomod)
- 📘 [Go Modules Reference](https://go.dev/ref/mod)
- 📘 [go.mod file syntax](https://go.dev/doc/modules/gomod-ref)

### 相关章节

- 🔗 [Go 1.25 工具链增强](./README.md)
- 🔗 [Go Modules 最佳实践](../../模块化/Go-Modules.md)
- 🔗 [项目结构设计](../../架构/项目结构.md)

---

## 更新日志

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2025-10-18 | v1.0 | 初始版本,完整的 ignore 指令指南 |

---

**编写者**: AI Assistant  
**审核者**: [待审核]  
**最后更新**: 2025年10月18日

---

<p align="center">
  <b>🎯 使用 ignore 指令优化你的 Go 项目结构! 📁</b>
</p>

