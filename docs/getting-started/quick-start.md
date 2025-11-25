# 🚀 快速开始指南

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [🚀 快速开始指南](#快速开始指南)
  - [安装Go](#安装go)
- [输出: go version go1.25.x windows/amd64](#输出-go-version-go125x-windowsamd64)
- [设置GOPATH (可选，Go 1.11+已启用Go Modules)](#设置gopath-可选go-111已启用go-modules)
- [设置Go代理 (国内用户推荐)](#设置go代理-国内用户推荐)
  - [第一个Go程序](#第一个go程序)
- [输出: Hello, World!](#输出-hello-world)
- [输出: Hello, World!](#输出-hello-world)
  - [如何使用本文档库](#如何使用本文档库)
  - [推荐学习路径](#推荐学习路径)
  - [快速参考](#快速参考)
- [初始化模块](#初始化模块)
- [下载依赖](#下载依赖)
- [整理依赖](#整理依赖)
- [查看依赖图](#查看依赖图)
- [运行代码](#运行代码)
- [编译](#编译)
- [交叉编译](#交叉编译)
- [运行测试](#运行测试)
- [运行测试(详细输出)](#运行测试详细输出)
- [运行测试(覆盖率)](#运行测试覆盖率)
- [基准测试](#基准测试)
- [格式化代码](#格式化代码)
- [检查代码](#检查代码)
- [静态分析](#静态分析)
  - [获取帮助](#获取帮助)
  - [💡 学习建议](#学习建议)
  - [下一步](#下一步)
  - [🎁 附加资源](#附加资源)
  - [🆘 遇到问题？](#遇到问题)
  - [🚀 准备好了吗？开始您的Go语言之旅](#准备好了吗开始您的go语言之旅)

---

## 安装Go

### Windows

1. **下载安装包**

   ```text
   https://go.dev/dl/
   ```

   下载 `go1.25.x.windows-amd64.msi`

2. **安装**
   - 双击运行安装包
   - 默认安装到 `C:\Program Files\Go`

3. **验证安装**

   ```powershell
   go version
   # 输出: go version go1.25.x windows/amd64
   ```

### macOS

1. **使用Homebrew** (推荐)

   ```bash
   brew install go
   ```

2. **或下载安装包**

   ```text
   https://go.dev/dl/
   ```

   下载 `go1.25.x.darwin-amd64.pkg` (Intel) 或 `go1.25.x.darwin-arm64.pkg` (M1/M2)

3. **验证安装**

   ```bash
   go version
   ```

### Linux

1. **下载并解压**

   ```bash
   wget https://go.dev/dl/go1.25.x.linux-amd64.tar.gz
   sudo rm -rf /usr/local/go
   sudo tar -C /usr/local -xzf go1.25.x.linux-amd64.tar.gz
   ```

2. **配置环境变量**

   ```bash
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **验证安装**

   ```bash
   go version
   ```

### 配置Go环境

```bash
# 设置GOPATH (可选，Go 1.11+已启用Go Modules)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# 设置Go代理 (国内用户推荐)
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

---

## 第一个Go程序

### 2. 初始化Go模块

```bash
go mod init hello-world
```

### 3. 创建main.go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### 4. 运行程序

```bash
go run main.go
# 输出: Hello, World!
```

### 5. 编译程序

```bash
go build -o hello
./hello  # Windows: hello.exe
# 输出: Hello, World!
```

### 🎉 恭喜

您已经完成了第一个Go程序！

**下一步**: 学习[变量和常量](01-语言基础/01-语法基础/02-变量和常量.md)

---

## 如何使用本文档库

### 📚 文档结构

本文档库按照**从基础到高级**的顺序组织：

```text
docs-new/
├── 01-语言基础/          ⭐ 从这里开始
├── 02-数据结构与算法/
├── 03-Web开发/           🌐 Web应用
├── 04-数据库编程/
├── 05-微服务架构/        🏗️ 分布式系统
├── 06-云原生与容器/      ☁️ 容器化部署
├── 07-性能优化/          ⚡ 性能提升
├── 08-架构设计/          🏛️ 架构模式
├── 09-工程实践/
├── 10-Go版本特性/        🆕 版本特性
├── 11-高级专题/          🎓 高级主题
├── 12-行业应用/          🏢 实战案例
└── 13-参考资料/          📖 技术报告
```

### 🔍 查找内容的方法

#### 方法1: 使用索引 (推荐)

查看 [INDEX.md](INDEX.md) - 全部177个文档的完整索引

- 按技术主题索引
- 按难度等级索引
- 按应用场景索引
- 按Go版本索引

#### 方法3: 使用搜索

- 在GitHub中按 `/` 键搜索
- 或使用IDE的全局搜索功能

### 📖 阅读文档

每个文档包含：

- ✅ **简介**: 概述主题
- ✅ **代码示例**: 可运行的代码
- ✅ **最佳实践**: 实战经验
- ✅ **常见问题**: FAQ
- ✅ **扩展阅读**: 相关文档链接

---

## 推荐学习路径

### 🎯 根据您的目标选择

#### 零基础学习 (2-4周)

```text
01-语言基础 → 03-Web开发基础 → 小项目实战
```

**详细路径**: [零基础入门路径](LEARNING_PATHS.md#%e9%9b%b6%e5%9f%ba%e7%a1%80%e5%85%a5%e9%97%a8%e8%b7%af%e5%be%84)

#### Web开发 (4-8周)

```text
语言基础 → Web开发 → 数据库编程 → 项目实战
```

**详细路径**: [Web开发路径](LEARNING_PATHS.md#web%e5%bc%80%e5%8f%91%e8%b7%af%e5%be%84)

#### 微服务开发 (8-12周)

```text
语言基础 → Web开发 → 微服务架构 → 云原生 → 大型项目
```

**详细路径**: [微服务开发路径](LEARNING_PATHS.md#%e5%be%ae%e6%9c%8d%e5%8a%a1%e5%bc%80%e5%8f%91%e8%b7%af%e5%be%84)

#### 算法面试 (8-12周)

```text
语言基础 → 数据结构与算法 → LeetCode刷题
```

**详细路径**: [算法面试路径](LEARNING_PATHS.md#%e7%ae%97%e6%b3%95%e9%9d%a2%e8%af%95%e8%b7%af%e5%be%84)

### 📅 按时间规划

| 时间 | 目标 | 路径 |
|------|------|------|
| 4周 | 快速入门 | [4周计划](LEARNING_PATHS.md#%ef%bf%bd%ef%bf%bd-4%e5%91%a8%e5%bf%ab%e9%80%9f%e5%85%a5%e9%97%a8) |
| 3个月 | 成为Go开发者 | [3个月计划](LEARNING_PATHS.md#%ef%bf%bd%ef%bf%bd-3%e4%b8%aa%e6%9c%88%e6%88%90%e4%b8%bago%e5%bc%80%e5%8f%91%e8%80%85) |
| 6个月 | 进阶工程师 | [6个月计划](LEARNING_PATHS.md#%ef%bf%bd%ef%bf%bd-6%e4%b8%aa%e6%9c%88%e8%bf%9b%e9%98%b6) |
| 1年 | 成为专家 | [1年计划](LEARNING_PATHS.md#%ef%bf%bd%ef%bf%bd-1%e5%b9%b4%e6%88%90%e4%b8%ba%e4%b8%93%e5%ae%b6) |

---

## 快速参考

### 常用命令

#### 项目管理

```bash
# 初始化模块
go mod init <module-name>

# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 查看依赖图
go mod graph
```

#### 运行与编译

```bash
# 运行代码
go run main.go

# 编译
go build
go build -o myapp

# 交叉编译
GOOS=linux GOARCH=amd64 go build
```

#### 测试

```bash
# 运行测试
go test

# 运行测试(详细输出)
go test -v

# 运行测试(覆盖率)
go test -cover

# 基准测试
go test -bench=.
```

#### 格式化与检查

```bash
# 格式化代码
go fmt ./...

# 检查代码
go vet ./...

# 静态分析
golangci-lint run
```

### 核心概念一览

| 概念 | 简介 | 文档链接 |
|------|------|---------|
| Goroutine | 轻量级并发 | [详情](01-语言基础/02-并发编程/02-Goroutine基础.md) |
| Channel | 通信机制 | [详情](01-语言基础/02-并发编程/03-Channel基础.md) |
| Interface | 接口 | [详情](01-语言基础/01-语法基础/03-基本数据类型.md) |
| Defer | 延迟执行 | [详情](01-语言基础/01-语法基础/04-流程控制.md) |
| Slice | 切片 | [详情](02-数据结构与算法/01-基础数据结构.md) |
| Map | 字典/哈希表 | [详情](02-数据结构与算法/01-基础数据结构.md) |

### 常用标准库

| 包 | 用途 | 文档链接 |
|------|------|---------|
| `fmt` | 格式化I/O | [官方文档](https://pkg.go.dev/fmt) |
| `net/http` | HTTP客户端和服务器 | [详情](03-Web开发/02-net-http包.md) |
| `encoding/json` | JSON编解码 | [官方文档](https://pkg.go.dev/encoding/json) |
| `database/sql` | SQL数据库接口 | [详情](04-数据库编程/01-MySQL编程.md) |
| `Context` | 上下文控制 | [详情](01-语言基础/02-并发编程/05-select与context.md) |
| `sync` | 同步原语 | [详情](01-语言基础/02-并发编程/06-sync包.md) |
| `testing` | 测试框架 | [详情](09-工程实践/00-Go测试深度实战指南.md) |

---

## 获取帮助

### 📚 本文档库

1. **技术索引**: [INDEX.md](INDEX.md) - 查找特定主题
2. **学习路径**: [LEARNING_PATHS.md](LEARNING_PATHS.md) - 系统学习计划
3. **常见问题**: [FAQ.md](FAQ.md) - 常见问题解答
4. **术语表**: [GLOSSARY.md](GLOSSARY.md) - 技术术语

### 🌐 官方资源

- **Go官网**: <https://go.dev/>
- **Go文档**: <https://pkg.go.dev/>
- **Go博客**: <https://go.dev/blog/>
- **Go Playground**: <https://go.dev/play/>

### 💬 社区支持

- **Go官方论坛**: <https://forum.golangbridge.org/>
- **Stack Overflow**: 搜索 `[go]` 标签
- **Reddit**: r/golang
- **Discord**: Gophers Slack

### 📖 推荐书籍

1. **《Go程序设计语言》** - 入门经典
2. **《Go语言实战》** - 实战指南
3. **《Go并发编程实战》** - 并发深入

### 🎓 在线课程

- **Tour of Go**: <https://go.dev/tour/>
- **Go by Example**: <https://gobyexample.com/>
- **Effective Go**: <https://go.dev/doc/effective_go>

---

## 💡 学习建议

### 1. 动手实践 🔨

- ✅ **每天写代码**: 至少1小时
- ✅ **完成练习**: 每章节的练习题
- ✅ **做小项目**: 巩固知识

### 2. 循序渐进 📈

- ✅ **不要跳步**: 打好基础
- ✅ **理论+实践**: 看完就写
- ✅ **反复练习**: 熟能生巧

### 3. 阅读代码 📖

- ✅ **读标准库**: 学习最佳实践
- ✅ **读优秀项目**: Docker, Kubernetes, Prometheus
- ✅ **参与开源**: 贡献代码

### 4. 持续学习 🚀

- ✅ **关注更新**: Go版本特性
- ✅ **技术分享**: 写博客、做分享
- ✅ **社区参与**: 加入Go社区

---

## 下一步

### 🎯 立即开始

选择一条路径，开始学习：

1. **零基础**: → [Hello World](01-语言基础/01-语法基础/01-Hello-World.md)
2. **有编程经验**: → [并发编程](01-语言基础/02-并发编程/README.md)
3. **想做Web**: → [Web开发](03-Web开发/README.md)
4. **想做微服务**: → [微服务架构](05-微服务架构/README.md)
5. **刷算法题**: → [数据结构与算法](02-数据结构与算法/README.md)

### 📚 深入学习

- 查看 [学习路径图](LEARNING_PATHS.md)
- 浏览 [技术索引](INDEX.md)
- 加入 Go社区

---

## 🎁 附加资源

### IDE推荐

1. **VS Code** + Go扩展 (推荐)
2. **GoLand** (JetBrains)
3. **Vim/Neovim** + vim-go

### 工具推荐

- **golangci-lint**: 代码检查
- **air**: 热重载
- **delve**: 调试器
- **mockgen**: Mock生成器

### 有用的网站

- **Go官方文档**: <https://pkg.go.dev/>
- **Go by Example**: <https://gobyexample.com/>
- **Go Playground**: <https://go.dev/play/>
- **Awesome Go**: <https://awesome-go.com/>

---

## 🆘 遇到问题？

1. **查看** [常见问题](FAQ.md)
2. **搜索** [技术索引](INDEX.md)
3. **提问** 在社区或Stack Overflow
4. **提交** GitHub Issue

---

<div align="center">

## 🚀 准备好了吗？开始您的Go语言之旅
