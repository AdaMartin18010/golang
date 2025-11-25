# 📦 安装指南

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: v2.0.0

---

## 📋 目录

- [📦 安装指南](#安装指南)
  - [🎯 系统要求](#系统要求)
  - [🚀 安装方式](#安装方式)
- [创建项目目录](#创建项目目录)
- [初始化Go模块](#初始化go模块)
- [安装主包](#安装主包)
- [或安装特定模块](#或安装特定模块)
- [克隆完整仓库](#克隆完整仓库)
- [或克隆特定版本](#或克隆特定版本)
- [下载所有依赖](#下载所有依赖)
- [验证依赖](#验证依赖)
- [构建所有模块](#构建所有模块)
- [或使用CLI工具](#或使用cli工具)
- [运行测试](#运行测试)
- [安装所有包](#安装所有包)
- [或安装特定命令](#或安装特定命令)
- [拉取官方镜像（如果有）](#拉取官方镜像如果有)
- [复制go.mod和go.sum](#复制gomod和gosum)
- [复制源码](#复制源码)
- [构建](#构建)
- [运行阶段](#运行阶段)
- [复制构建产物](#复制构建产物)
- [运行](#运行)
- [构建镜像](#构建镜像)
- [运行容器](#运行容器)
  - [🛠️ 安装CLI工具](#️-安装cli工具)
- [克隆仓库](#克隆仓库)
- [构建并安装](#构建并安装)
- [或构建到当前目录](#或构建到当前目录)
- [检查版本](#检查版本)
- [查看帮助](#查看帮助)
- [列出所有命令](#列出所有命令)
- [Linux/macOS](#linuxmacos)
- [或添加到 ~/.zshrc](#或添加到-zshrc)
- [Windows (PowerShell)](#windows-powershell)
- [永久添加（Windows）](#永久添加windows)
  - [✅ 验证安装](#验证安装)
- [检查Go版本](#检查go版本)
- [输出: go version go1.25.3 ...](#输出-go-version-go1253)
- [检查Go环境](#检查go环境)
- [列出已安装的包](#列出已安装的包)
- [查看包信息](#查看包信息)
- [进入项目目录](#进入项目目录)
- [运行测试](#运行测试)
- [克隆示例代码](#克隆示例代码)
- [运行示例](#运行示例)
  - [⚙️ 环境配置](#️-环境配置)
- [使用Go代理](#使用go代理)
- [或使用阿里云代理](#或使用阿里云代理)
- [启用Go模块](#启用go模块)
- [配置私有仓库](#配置私有仓库)
- [配置不进行checksum验证的模块](#配置不进行checksum验证的模块)
- [创建工作区](#创建工作区)
- [初始化工作区](#初始化工作区)
- [添加模块](#添加模块)
- [同步](#同步)
  - [🔧 IDE配置](#ide配置)
  - [❓ 常见问题](#常见问题)
- [Linux/macOS](#linuxmacos)
- [或使用用户目录](#或使用用户目录)
- [使用代理](#使用代理)
- [清理缓存](#清理缓存)
- [重新下载](#重新下载)
- [清理模块缓存](#清理模块缓存)
- [删除go.sum](#删除gosum)
- [重新整理](#重新整理)
- [检查gox路径](#检查gox路径)
- [如果找不到，检查GOPATH/bin是否在PATH中](#如果找不到检查gopathbin是否在path中)
- [手动添加](#手动添加)
- [检查Go版本](#检查go版本)
- [必须是1.25.3+](#必须是1253)
- [如果版本过低，升级Go](#如果版本过低升级go)
- [验证依赖](#验证依赖)
- [下载缺失的依赖](#下载缺失的依赖)
  - [🆘 获取帮助](#获取帮助)
  - [📚 下一步](#下一步)

---

## 🎯 系统要求

### 必需条件

| 组件 | 最低版本 | 推荐版本 |
|------|---------|----------|
| Go | 1.25.3 | 1.25.3+ |
| 操作系统 | - | Windows 10+, Linux 4.0+, macOS 10.15+ |
| 内存 | 512MB | 2GB+ |
| 磁盘空间 | 100MB | 500MB+ |

### 可选依赖

- **Git** - 从源码安装时需要
- **Make** - 使用Makefile时需要（可选）
- **Docker** - 使用容器化部署时需要

---

## 🚀 安装方式

### 使用Go Modules（推荐）

这是最简单和推荐的安装方式。

#### 1. 创建新项目

```bash
# 创建项目目录
mkdir my-golang-project
cd my-golang-project

# 初始化Go模块
go mod init my-golang-project
```

#### 2. 安装包

```bash
# 安装主包
go get github.com/yourusername/golang@v2.0.0

# 或安装特定模块
go get github.com/yourusername/golang/pkg/agent@v2.0.0
go get github.com/yourusername/golang/pkg/concurrency@v2.0.0
go get github.com/yourusername/golang/pkg/http3@v2.0.0
go get github.com/yourusername/golang/pkg/memory@v2.0.0
go get github.com/yourusername/golang/pkg/observability@v2.0.0
```

#### 3. 同步依赖

```bash
go mod tidy
```

#### 4. 在代码中使用

```go
package main

import (
    "Context"
    "fmt"
    "github.com/yourusername/golang/pkg/observability"
)

func main() {
    // 使用observability
    observability.Info("Application started")

    ctx := Context.Background()
    span, ctx := observability.StartSpan(ctx, "main-operation")
    defer span.Finish()

    fmt.Println("Hello, Golang v2.0!")
}
```

---

### 从源码安装

适合需要修改源码或贡献代码的用户。

#### 1. 克隆仓库

```bash
# 克隆完整仓库
git clone https://github.com/yourusername/golang.git
cd golang

# 或克隆特定版本
git clone -b v2.0.0 https://github.com/yourusername/golang.git
cd golang
```

#### 2. 安装依赖

```bash
# 下载所有依赖
go mod download

# 验证依赖
go mod verify
```

#### 3. 构建项目

```bash
# 构建所有模块
go build ./...

# 或使用CLI工具
cd cmd/gox
go build -o gox

# 运行测试
go test ./...
```

#### 4. 安装到GOPATH

```bash
# 安装所有包
go install ./...

# 或安装特定命令
go install ./cmd/gox@latest
```

---

### 使用Docker

适合容器化部署的场景。

#### 1. 拉取镜像

```bash
# 拉取官方镜像（如果有）
docker pull yourusername/golang:v2.0.0
```

#### 2. 或构建自己的镜像

创建 `Dockerfile`:

```dockerfile
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建
RUN go build -o /app/main ./cmd/your-app

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/main .

# 运行
ENTRYPOINT ["./main"]
```

构建和运行:

```bash
# 构建镜像
docker build -t my-golang-app:v2.0.0 .

# 运行容器
docker run -d -p 8080:8080 my-golang-app:v2.0.0
```

---

## 🛠️ 安装CLI工具

CLI工具(`gox`)提供了便捷的项目管理功能。

### 方式1: 直接安装

```bash
go install github.com/yourusername/golang/cmd/gox@v2.0.0
```

### 方式2: 从源码安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/golang.git
cd golang/cmd/gox

# 构建并安装
go install

# 或构建到当前目录
go build -o gox
```

### 验证安装

```bash
# 检查版本
gox version

# 查看帮助
gox help

# 列出所有命令
gox
```

### 配置环境变量

确保 `$GOPATH/bin` 在你的 `PATH` 中：

```bash
# Linux/macOS
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# 或添加到 ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc

# Windows (PowerShell)
$env:Path += ";$env:GOPATH\bin"

# 永久添加（Windows）
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:GOPATH\bin", "User")
```

---

## ✅ 验证安装

### 1. 检查Go环境

```bash
# 检查Go版本
go version
# 输出: go version go1.25.3 ...

# 检查Go环境
go env
```

### 2. 验证包安装

```bash
# 列出已安装的包
go list -m github.com/yourusername/golang

# 查看包信息
go list -m -versions github.com/yourusername/golang
```

### 3. 运行测试

```bash
# 进入项目目录
cd $GOPATH/pkg/mod/github.com/yourusername/golang@v2.0.0

# 运行测试
go test ./...
```

### 4. 运行示例

```bash
# 克隆示例代码
git clone https://github.com/yourusername/golang.git
cd golang/examples

# 运行示例
cd modern-features/observability
go run main.go
```

---

## ⚙️ 环境配置

### Go代理配置（中国用户）

```bash
# 使用Go代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 或使用阿里云代理
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

### Go模块配置

```bash
# 启用Go模块
go env -w GO111MODULE=on

# 配置私有仓库
go env -w GOPRIVATE=github.com/your-org/*

# 配置不进行checksum验证的模块
go env -w GONOSUMDB=github.com/your-org/*
```

### 工作区配置（可选）

如果你需要同时开发多个模块：

```bash
# 创建工作区
mkdir my-workspace
cd my-workspace

# 初始化工作区
go work init

# 添加模块
go work use ./module1
go work use ./module2

# 同步
go work sync
```

---

## 🔧 IDE配置

### VS Code

1. 安装Go扩展

    ```bash
    code --install-extension golang.go
    ```

2. 配置 `settings.json`:

    ```json
    {
        "go.useLanguageServer": true,
        "go.toolsManagement.autoUpdate": true,
        "go.lintTool": "golangci-lint",
        "go.lintOnSave": "workspace",
        "go.testFlags": ["-v", "-race"],
        "go.coverOnSave": true
    }
    ```

### GoLand

1. 打开设置: `File` > `Settings` > `Go`
2. 配置GOROOT指向Go 1.25.3
3. 配置GOPATH
4. 启用Go Modules: `Go` > `Go Modules` > `Enable Go Modules integration`

### Vim/Neovim

使用 `vim-go`:

```vim
" 在 .vimrc 或 init.vim 中添加
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

" Go配置
let g:go_fmt_command = "goimports"
let g:go_auto_type_info = 1
let g:go_def_mode='gopls'
let g:go_info_mode='gopls'
```

---

## ❓ 常见问题

### Q1: 安装时提示"permission denied"

**A**: 需要管理员权限或修改GOPATH权限

```bash
# Linux/macOS
sudo chown -R $USER:$USER $GOPATH

# 或使用用户目录
export GOPATH=$HOME/go
```

### Q2: 找不到包

**A**: 检查代理设置和网络连接

```bash
# 使用代理
go env -w GOPROXY=https://goproxy.cn,direct

# 清理缓存
go clean -modcache

# 重新下载
go mod download
```

### Q3: 版本冲突

**A**: 清理并重新安装

```bash
# 清理模块缓存
go clean -modcache

# 删除go.sum
rm go.sum

# 重新整理
go mod tidy
```

### Q4: CLI工具无法找到

**A**: 检查PATH配置

```bash
# 检查gox路径
which gox

# 如果找不到，检查GOPATH/bin是否在PATH中
echo $PATH | grep $GOPATH/bin

# 手动添加
export PATH=$PATH:$GOPATH/bin
```

### Q5: 构建失败

**A**: 检查Go版本和依赖

```bash
# 检查Go版本
go version

# 必须是1.25.3+
# 如果版本过低，升级Go

# 验证依赖
go mod verify

# 下载缺失的依赖
go mod download
```

### Q6: 导入路径错误

**A**: 检查模块路径

```go
// ❌ 错误
import "golang/pkg/agent"

// ✅ 正确
import "github.com/yourusername/golang/pkg/agent/core"
```

---

## 🆘 获取帮助

如果遇到其他问题：

1. **查看文档**: [完整文档](docs/README.md)
2. **搜索Issues**: [GitHub Issues](https://github.com/yourusername/golang/issues)
3. **提问**: [GitHub Discussions](https://github.com/yourusername/golang/discussions)
4. **联系**: <your-email@example.com>

---

## 📚 下一步
