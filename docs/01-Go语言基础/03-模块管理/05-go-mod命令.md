# go mod 命令详解

> **简介**: 全面介绍 `go mod` 命令的各种子命令、选项和使用场景，帮助开发者熟练掌握 Go Modules 的日常操作和高级技巧。

<!-- TOC START -->
- [go mod 命令详解](#go-mod-命令详解)
  - [1. 📚 理论分析](#1-理论分析)
    - [1.1 命令概述](#11-命令概述)
    - [1.2 命令分类](#12-命令分类)
  - [2. 💻 核心命令详解](#2-核心命令详解)
    - [2.1 go mod init](#21-go-mod-init)
    - [2.2 go mod tidy](#22-go-mod-tidy)
    - [2.3 go mod download](#23-go-mod-download)
    - [2.4 go mod verify](#24-go-mod-verify)
    - [2.5 go mod graph](#25-go-mod-graph)
    - [2.6 go mod why](#26-go-mod-why)
    - [2.7 go mod edit](#27-go-mod-edit)
    - [2.8 go mod vendor](#28-go-mod-vendor)
  - [3. 🔧 实践应用](#3-实践应用)
    - [3.1 初始化项目](#31-初始化项目)
    - [3.2 依赖管理](#32-依赖管理)
    - [3.3 故障排查](#33-故障排查)
  - [4. 📊 命令速查表](#4-命令速查表)
  - [5. 🎯 最佳实践](#5-最佳实践)
  - [6. ⚠️ 常见陷阱](#6-️-常见陷阱)
  - [7. 🔍 常见问题](#7-常见问题)
  - [8. 📚 扩展阅读](#8-扩展阅读)
<!-- TOC END -->

---

## 1. 📚 理论分析

### 1.1 命令概述

`go mod` 是 Go Modules 的核心命令工具，提供了模块管理的各种功能：

```bash
go mod <command> [arguments]
```

**设计理念**:

- **自动化**: 尽可能自动处理依赖关系
- **透明**: 明确显示操作结果
- **安全**: 验证依赖完整性
- **高效**: 使用缓存加速操作

### 1.2 命令分类

| 分类 | 命令 | 用途 |
|------|------|------|
| **初始化** | `init` | 初始化新模块 |
| **维护** | `tidy`, `download` | 整理和下载依赖 |
| **验证** | `verify` | 验证依赖完整性 |
| **查询** | `graph`, `why` | 分析依赖关系 |
| **编辑** | `edit` | 修改 go.mod 文件 |
| **特殊** | `vendor` | 创建vendor目录 |

---

## 2. 💻 核心命令详解

### 2.1 go mod init

**功能**: 初始化一个新的模块，创建 go.mod 文件

#### 2.1.1 基本用法

```bash
# 语法
go mod init [module-path]

# 示例
go mod init github.com/username/project
```

#### 2.1.2 自动推断模块路径

```bash
# 在 Git 仓库中自动推断
cd my-git-repo
go mod init  # 自动使用 git remote 路径
```

#### 2.1.3 本地开发

```bash
# 使用简短名称（仅本地开发）
go mod init myapp

# 注意：发布时应使用完整路径
go mod edit -module=github.com/username/myapp
```

#### 2.1.4 生成的 go.mod

```go
module github.com/username/project

go 1.25  // 自动检测当前 Go 版本
```

**选项**:

- 无特殊选项

**错误处理**:

```bash
# 如果已存在 go.mod
go mod init  # Error: go.mod already exists

# 解决: 删除后重新初始化
rm go.mod go.sum
go mod init github.com/username/project
```

---

### 2.2 go mod tidy

**功能**: 整理 go.mod 和 go.sum，添加缺失的依赖，删除未使用的依赖

#### 2.2.1 基本用法

```bash
go mod tidy
```

#### 2.2.2 详细输出

```bash
# 显示详细信息
go mod tidy -v

# 输出示例
go: downloading github.com/gin-gonic/gin v1.9.1
go: downloading golang.org/x/net v0.17.0
```

#### 2.2.3 指定 Go 版本兼容

```bash
# 兼容旧版本 (Go 1.16)
go mod tidy -go=1.16

# 兼容最新版本
go mod tidy -go=1.25
```

#### 2.2.4 典型场景

**场景1: 添加新依赖后**

```bash
# 1. 在代码中导入新包
import "github.com/gin-gonic/gin"

# 2. 运行 tidy
go mod tidy

# 结果: go.mod 自动添加依赖
```

**场景2: 删除未使用的依赖**

```bash
# 1. 从代码中删除 import
# 2. 运行 tidy
go mod tidy

# 结果: go.mod 自动删除未使用的依赖
```

**选项**:

- `-v`: 显示详细输出
- `-go=version`: 更新 go 版本
- `-compat=version`: 保持与指定版本兼容

---

### 2.3 go mod download

**功能**: 下载依赖到本地缓存，但不修改 go.mod

#### 2.3.1 基本用法

```bash
# 下载所有依赖
go mod download

# 下载指定依赖
go mod download github.com/gin-gonic/gin
go mod download github.com/gin-gonic/gin@v1.9.1
```

#### 2.3.2 JSON 输出

```bash
# 输出下载信息为 JSON
go mod download -json

# 输出示例
{
    "Path": "github.com/gin-gonic/gin",
    "Version": "v1.9.1",
    "Info": "/path/to/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.info",
    "GoMod": "/path/to/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.mod",
    "Zip": "/path/to/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.zip",
    "Dir": "/path/to/cache/github.com/gin-gonic/gin@v1.9.1"
}
```

#### 2.3.3 CI/CD 使用

```bash
# Dockerfile 示例
FROM golang:1.25 AS builder
WORKDIR /app

# 先下载依赖（利用 Docker 层缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制代码
COPY . .
RUN go build -o myapp
```

```yaml
# GitHub Actions 示例
- name: Download dependencies
  run: go mod download
  
- name: Cache dependencies
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

**选项**:

- `-json`: JSON 格式输出
- `-x`: 显示执行的命令

---

### 2.4 go mod verify

**功能**: 验证依赖的完整性，检查缓存中的依赖是否被篡改

#### 2.4.1 基本用法

```bash
go mod verify
```

#### 2.4.2 输出结果

**成功**:

```bash
$ go mod verify
all modules verified
```

**失败**:

```bash
$ go mod verify
github.com/gin-gonic/gin v1.9.1: dir has been modified
```

#### 2.4.3 修复损坏的缓存

```bash
# 1. 清除缓存
go clean -modcache

# 2. 重新下载
go mod download

# 3. 再次验证
go mod verify
```

#### 2.4.4 CI/CD 中的使用

```bash
# 在构建前验证
go mod download
go mod verify
go build
```

**用途**:

- 检测依赖篡改
- 确保构建安全
- CI/CD 流程验证

---

### 2.5 go mod graph

**功能**: 打印模块依赖图

#### 2.5.1 基本用法

```bash
go mod graph
```

#### 2.5.2 输出格式

```text
github.com/user/project github.com/gin-gonic/gin@v1.9.1
github.com/gin-gonic/gin@v1.9.1 github.com/gin-contrib/sse@v0.1.0
github.com/gin-gonic/gin@v1.9.1 github.com/go-playground/validator/v10@v10.14.0
...
```

格式: `module_A module_B@version` 表示 A 依赖 B

#### 2.5.3 可视化依赖图

```bash
# 安装可视化工具
go install golang.org/x/exp/cmd/modgraphviz@latest

# 生成可视化图
go mod graph | modgraphviz | dot -Tsvg -o graph.svg
```

#### 2.5.4 过滤特定依赖

```bash
# 查找特定包的依赖链
go mod graph | grep gin-gonic

# 查找谁依赖某个包
go mod graph | grep "github.com/specific/package"
```

**用途**:

- 分析依赖关系
- 查找传递依赖
- 调试依赖问题

---

### 2.6 go mod why

**功能**: 解释为什么需要某个依赖

#### 2.6.1 基本用法

```bash
# 查看为什么需要某个包
go mod why github.com/gin-gonic/gin

# 输出示例
# github.com/username/project
github.com/username/project
github.com/gin-gonic/gin
```

#### 2.6.2 查询多个包

```bash
# 同时查询多个包
go mod why github.com/gin-gonic/gin golang.org/x/sync
```

#### 2.6.3 查询所有依赖

```bash
# 查看所有包的依赖原因
go mod why -m all
```

#### 2.6.4 查询供应商模式

```bash
# 在 vendor 模式下查询
go mod why -vendor github.com/gin-gonic/gin
```

**选项**:

- `-m`: 按模块而非包查询
- `-vendor`: 在 vendor 目录中查询

**用途**:

- 理解间接依赖
- 优化依赖树
- 排查无用依赖

---

### 2.7 go mod edit

**功能**: 编辑 go.mod 文件（脚本友好）

#### 2.7.1 修改模块路径

```bash
# 修改模块路径
go mod edit -module=github.com/new/path
```

#### 2.7.2 修改 Go 版本

```bash
# 设置 Go 版本
go mod edit -go=1.25
```

#### 2.7.3 添加依赖

```bash
# 添加或更新依赖
go mod edit -require=github.com/gin-gonic/gin@v1.9.1

# 批量添加
go mod edit -require=github.com/pkg1@v1.0.0 -require=github.com/pkg2@v2.0.0
```

#### 2.7.4 删除依赖

```bash
# 删除依赖
go mod edit -droprequire=github.com/unused/package
```

#### 2.7.5 替换依赖

```bash
# 替换为其他版本
go mod edit -replace=github.com/old/pkg@v1.0.0=github.com/new/pkg@v2.0.0

# 替换为本地路径
go mod edit -replace=github.com/some/pkg=../local/pkg

# 删除替换
go mod edit -dropreplace=github.com/old/pkg@v1.0.0
```

#### 2.7.6 排除版本

```bash
# 排除特定版本
go mod edit -exclude=github.com/broken/pkg@v1.2.3

# 删除排除
go mod edit -dropexclude=github.com/broken/pkg@v1.2.3
```

#### 2.7.7 JSON 输出

```bash
# 以 JSON 格式打印 go.mod
go mod edit -json
```

#### 2.7.8 格式化

```bash
# 格式化 go.mod
go mod edit -fmt
```

**常用脚本**:

```bash
#!/bin/bash
# 批量更新依赖到本地开发版本

go mod edit -replace=github.com/proj/auth=../auth
go mod edit -replace=github.com/proj/api=../api
go mod edit -replace=github.com/proj/db=../db

go mod tidy
```

**选项**:

- `-module`: 修改模块路径
- `-go`: 修改 Go 版本
- `-require`: 添加依赖
- `-droprequire`: 删除依赖
- `-replace`: 添加替换
- `-dropreplace`: 删除替换
- `-exclude`: 排除版本
- `-dropexclude`: 删除排除
- `-json`: JSON 输出
- `-fmt`: 格式化
- `-print`: 打印而不写入文件

---

### 2.8 go mod vendor

**功能**: 将依赖复制到 vendor 目录

#### 2.8.1 基本用法

```bash
# 创建 vendor 目录
go mod vendor
```

#### 2.8.2 生成的目录结构

```text
project/
├── go.mod
├── go.sum
├── main.go
└── vendor/
    ├── modules.txt          # 记录vendor的模块
    └── github.com/
        └── gin-gonic/
            └── gin/
                └── ...
```

#### 2.8.3 使用 vendor 构建

```bash
# 显式使用 vendor
go build -mod=vendor

# 或设置环境变量
export GOFLAGS="-mod=vendor"
go build
```

#### 2.8.4 验证 vendor

```bash
# 验证 vendor 是否与 go.mod 一致
go mod vendor
git diff --exit-code vendor/
```

#### 2.8.5 清理 vendor

```bash
# 删除 vendor 目录
rm -rf vendor/

# 恢复使用模块缓存
go build
```

**使用场景**:

- 离线构建
- 企业内网环境
- 确保构建一致性
- CI/CD 加速

**选项**:

- `-v`: 显示详细输出
- `-e`: 验证包

---

## 3. 🔧 实践应用

### 3.1 初始化项目

```bash
# 完整流程
mkdir myproject
cd myproject

# 初始化模块
go mod init github.com/username/myproject

# 创建代码
cat > main.go << 'EOF'
package main
import "fmt"
func main() {
    fmt.Println("Hello, Modules!")
}
EOF

# 整理依赖
go mod tidy

# 验证
go build
./myproject
```

### 3.2 依赖管理

#### 3.2.1 添加依赖的完整流程

```bash
# 1. 在代码中使用
# import "github.com/gin-gonic/gin"

# 2. 自动下载并添加到 go.mod
go mod tidy

# 3. 验证依赖
go mod verify

# 4. 查看依赖树
go mod graph

# 5. 提交更改
git add go.mod go.sum
git commit -m "Add gin dependency"
```

#### 3.2.2 更新依赖

```bash
# 1. 列出可更新的依赖
go list -u -m all

# 2. 更新所有依赖
go get -u ./...

# 3. 整理
go mod tidy

# 4. 测试
go test ./...

# 5. 提交
git add go.mod go.sum
git commit -m "Update dependencies"
```

#### 3.2.3 锁定依赖版本

```bash
# 使用具体版本
go get github.com/gin-gonic/gin@v1.9.1

# 或编辑 go.mod 后
go mod download
```

### 3.3 故障排查

#### 3.3.1 依赖下载失败

```bash
# 1. 检查代理设置
go env GOPROXY

# 2. 更换代理
go env -w GOPROXY=https://goproxy.cn,direct

# 3. 清除缓存重试
go clean -modcache
go mod download
```

#### 3.3.2 依赖版本冲突

```bash
# 1. 查看依赖图
go mod graph | grep conflicting-package

# 2. 查看为什么需要
go mod why github.com/conflicting/package

# 3. 使用 replace 统一版本
go mod edit -replace=old@v1.0.0=new@v2.0.0
go mod tidy
```

#### 3.3.3 go.sum 校验失败

```bash
# 1. 验证依赖
go mod verify

# 2. 如果失败，清除并重新下载
go clean -modcache
rm go.sum
go mod download
go mod tidy
```

---

## 4. 📊 命令速查表

| 命令 | 功能 | 常用场景 | 频率 |
|------|------|---------|------|
| `go mod init` | 初始化模块 | 创建新项目 | 低 |
| `go mod tidy` | 整理依赖 | 添加/删除依赖后 | 高 |
| `go mod download` | 下载依赖 | CI/CD 构建 | 中 |
| `go mod verify` | 验证依赖 | 安全检查 | 中 |
| `go mod graph` | 打印依赖图 | 分析依赖 | 低 |
| `go mod why` | 解释依赖原因 | 调试依赖 | 低 |
| `go mod edit` | 编辑 go.mod | 脚本操作 | 中 |
| `go mod vendor` | 创建 vendor | 离线构建 | 低 |

---

## 5. 🎯 最佳实践

### 5.1 日常开发

```bash
# 1. 每次添加依赖后
go mod tidy

# 2. 定期验证
go mod verify

# 3. 提交前检查
go mod tidy
git diff go.mod go.sum
```

### 5.2 团队协作

```bash
# 1. 克隆项目后
git clone ...
cd project
go mod download  # 下载依赖
go mod verify    # 验证完整性

# 2. 更新依赖前
git pull
go mod tidy
go test ./...    # 确保测试通过

# 3. 提交依赖变更
git add go.mod go.sum
git commit -m "Update/Add dependencies"
```

### 5.3 CI/CD

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      # 缓存依赖
      - uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      # 下载并验证
      - run: go mod download
      - run: go mod verify
      
      # 构建和测试
      - run: go build -v ./...
      - run: go test -v ./...
```

---

## 6. ⚠️ 常见陷阱

### 6.1 忘记运行 tidy

❌ **错误做法**:

```bash
# 修改代码后直接提交
git add .
git commit -m "Add feature"
```

✅ **正确做法**:

```bash
go mod tidy
git add .
git commit -m "Add feature"
```

### 6.2 手动编辑 go.sum

❌ **错误做法**:

```bash
# 手动编辑 go.sum
vim go.sum
```

✅ **正确做法**:

```bash
# 让 Go 工具管理
go mod tidy
```

### 6.3 不提交 go.sum

❌ **错误做法**:

```gitignore
# .gitignore
go.sum
```

✅ **正确做法**:

```bash
# 提交 go.sum
git add go.mod go.sum
```

### 6.4 过度使用 vendor

❌ **过度使用**:

```bash
# 每个项目都用 vendor
go mod vendor
```

✅ **按需使用**:

```bash
# 仅在必要时使用 vendor
# - 离线构建
# - 企业内网
# - 特殊要求
```

---

## 7. 🔍 常见问题

### 7.1 Q: go mod init 后可以更改模块路径吗？

**A**: 可以，使用 `go mod edit`

```bash
go mod edit -module=github.com/new/path
```

### 7.2 Q: go mod tidy 会删除哪些依赖？

**A**: 删除以下依赖：

- 代码中未使用的直接依赖
- 不再需要的间接依赖
- go.mod 中多余的 require

### 7.3 Q: 如何强制使用特定版本？

**A**: 使用 `replace` 指令

```bash
go mod edit -replace=pkg@v1.0.0=pkg@v2.0.0
```

### 7.4 Q: vendor 和模块缓存有什么区别？

**A**:

| 特性 | Vendor | 模块缓存 |
|------|--------|---------|
| 位置 | 项目内 `vendor/` | 全局 `$GOPATH/pkg/mod` |
| 提交 | 提交到版本控制 | 不提交 |
| 共享 | 不共享 | 所有项目共享 |
| 大小 | 较大（每个项目都有） | 较小（全局一份） |

### 7.5 Q: 如何查看模块的所有可用版本？

**A**: 使用 `go list`

```bash
go list -m -versions github.com/gin-gonic/gin
```

---

## 8. 📚 扩展阅读

### 8.1 官方文档

- [go mod Command](https://go.dev/ref/mod#go-mod-init) - 官方命令参考
- [Module Commands](https://go.dev/cmd/go/#hdr-Module_maintenance) - 模块维护命令
- [go Command](https://go.dev/cmd/go/) - Go 命令行工具

### 8.2 相关文档

- [Go Modules简介](./01-Go-Modules简介.md)
- [go.mod文件详解](./02-go-mod文件详解.md)
- [依赖管理](./06-依赖管理.md)

### 8.3 工具和脚本

- [modgraphviz](https://pkg.go.dev/golang.org/x/exp/cmd/modgraphviz) - 依赖图可视化
- [go-mod-upgrade](https://github.com/oligot/go-mod-upgrade) - 依赖更新工具
- [gomods](https://github.com/Helcaraxan/gomod) - 依赖分析工具

---

**文档维护者**: Go Modules Team  
**最后更新**: 2025年10月19日  
**文档状态**: 完成  
**适用版本**: Go 1.11+，推荐 Go 1.25.3+
