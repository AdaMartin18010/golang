# Go模块管理

Go模块管理完整指南，涵盖go.mod、go.sum、Go Workspace和版本管理。

---

## 📋 目录

- [📚 文档列表](#文档列表)
- [🚀 快速示例](#快速示例)
  - [初始化模块](#初始化模块)
  - [go.mod示例](#go-mod示例)
  - [Go Workspace](#go-workspace)
  - [go.work示例](#go-work示例)
  - [依赖管理](#依赖管理)
- [📖 系统文档](#系统文档)
- [🔗 相关资源](#相关资源)

## 📚 文档列表

1. **[Go Modules基础](./01-Go-Modules基础.md)** ⭐⭐⭐⭐⭐
   - go.mod文件结构
   - module, require, replace
   - 版本选择算法

2. **[依赖管理](./02-依赖管理.md)** ⭐⭐⭐⭐⭐
   - go get, go mod tidy
   - 私有仓库
   - 版本约束

3. **[Go Workspace](./03-Go-Workspace.md)** ⭐⭐⭐⭐⭐
   - go.work文件
   - 多模块开发
   - Monorepo支持

4. **[版本管理](./04-版本管理.md)** ⭐⭐⭐⭐
   - 语义化版本(SemVer)
   - 版本标签
   - 发布流程

---

## 🚀 快速示例

### 初始化模块

```bash
go mod init github.com/username/myproject
```

### go.mod示例

```go
module github.com/username/myproject

go 1.25.3

require (
    github.com/gin-gonic/gin v1.10.0
    gorm.io/gorm v1.25.5
)

replace github.com/old/module => github.com/new/module v1.0.0
```

### Go Workspace

```bash
go work init ./module1 ./module2
```

### go.work示例

```go
go 1.25.3

use (
    ./backend
    ./frontend
    ./shared
)
```

### 依赖管理

```bash
# 添加依赖
go get github.com/gin-gonic/gin@latest

# 更新依赖
go get -u ./...

# 清理依赖
go mod tidy

# 下载依赖
go mod download

# 查看依赖树
go mod graph
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

## 🔗 相关资源

- [Go Workspace完整指南](../../00-Go-Workspace完整指南-Go1.25.3.md)
- [Go Modules与Workspace完整对比](../../00-Go-Modules与Workspace完整对比-2025.md)

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
