# 🔄 迁移指南 - v1.x to v2.0

> **目标版本**: v2.0.0  
> **源版本**: v1.x  
> **预计迁移时间**: 30-60分钟  
> **难度**: 中等

---

## 📋 目录

- [🔄 迁移指南 - v1.x to v2.0](#-迁移指南---v1x-to-v20)
  - [📋 目录](#-目录)
  - [🎯 快速概览](#-快速概览)
    - [主要变更](#主要变更)
    - [兼容性矩阵](#兼容性矩阵)
  - [📝 迁移前准备](#-迁移前准备)
    - [1. 备份当前代码](#1-备份当前代码)
    - [2. 检查当前版本](#2-检查当前版本)
    - [3. 记录当前依赖](#3-记录当前依赖)
    - [4. 准备迁移清单](#4-准备迁移清单)
  - [🚀 分步迁移指南](#-分步迁移指南)
    - [步骤1: 更新Go模块](#步骤1-更新go模块)
      - [1.1 更新go.mod](#11-更新gomod)
      - [1.2 验证模块](#12-验证模块)
    - [步骤2: 更新导入路径](#步骤2-更新导入路径)
      - [2.1 创建替换脚本](#21-创建替换脚本)
      - [2.2 手动检查特殊情况](#22-手动检查特殊情况)
    - [步骤3: 更新API调用](#步骤3-更新api调用)
      - [3.1 Observability模块](#31-observability模块)
      - [3.2 Agent模块](#32-agent模块)
    - [步骤4: 更新配置](#步骤4-更新配置)
      - [4.1 日志配置](#41-日志配置)
      - [4.2 文件权限配置](#42-文件权限配置)
    - [步骤5: 运行测试](#步骤5-运行测试)
    - [步骤6: 更新文档](#步骤6-更新文档)
  - [💥 Breaking Changes详解](#-breaking-changes详解)
    - [1. 目录结构变更](#1-目录结构变更)
      - [变更对照表](#变更对照表)
      - [迁移影响](#迁移影响)
    - [2. API签名变更](#2-api签名变更)
      - [Observability - Register函数](#observability---register函数)
    - [3. 文件权限变更](#3-文件权限变更)
      - [默认权限更加严格](#默认权限更加严格)
      - [迁移检查清单](#迁移检查清单)
  - [📝 代码迁移示例](#-代码迁移示例)
    - [示例1: 完整的HTTP服务迁移](#示例1-完整的http服务迁移)
    - [示例2: 并发模式迁移](#示例2-并发模式迁移)
  - [❓ 常见迁移问题](#-常见迁移问题)
    - [Q1: 导入路径找不到](#q1-导入路径找不到)
    - [Q2: 类型不匹配](#q2-类型不匹配)
    - [Q3: 测试失败](#q3-测试失败)
    - [Q4: 性能下降](#q4-性能下降)
  - [🔙 回滚方案](#-回滚方案)
    - [方案1: Git回滚](#方案1-git回滚)
    - [方案2: 使用v1.x版本](#方案2-使用v1x版本)
    - [方案3: 渐进式迁移](#方案3-渐进式迁移)
  - [✅ 迁移验证清单](#-迁移验证清单)
  - [📞 获取帮助](#-获取帮助)
  - [📚 相关资源](#-相关资源)

---

## 🎯 快速概览

### 主要变更

| 变更类型 | 影响范围 | 严重程度 |
|---------|---------|---------|
| 目录结构重组 | 所有导入路径 | ⚠️ 高 |
| API签名变更 | 部分API | ℹ️ 中 |
| 配置格式 | 配置文件 | ℹ️ 低 |
| 文件权限 | 文件操作 | ℹ️ 低 |

### 兼容性矩阵

| 组件 | v1.x | v2.0 | 向后兼容 |
|------|------|------|---------|
| 核心API | ✅ | ✅ | 部分 |
| 导入路径 | ✅ | ⚠️ | ❌ |
| 配置文件 | ✅ | ✅ | ✅ |
| 数据格式 | ✅ | ✅ | ✅ |

---

## 📝 迁移前准备

### 1. 备份当前代码

```bash
# 创建分支
git checkout -b pre-v2-migration

# 提交当前状态
git add .
git commit -m "Backup before v2.0 migration"

# 推送到远程
git push origin pre-v2-migration
```

### 2. 检查当前版本

```bash
# 检查Go版本
go version

# 必须是1.25.3+
# 如果不是，先升级Go
```

### 3. 记录当前依赖

```bash
# 保存当前依赖列表
go list -m all > dependencies-before-migration.txt

# 运行测试确保一切正常
go test ./...
```

### 4. 准备迁移清单

- [ ] 更新Go版本到1.25.3+
- [ ] 备份代码
- [ ] 记录当前配置
- [ ] 通知团队
- [ ] 计划回滚方案

---

## 🚀 分步迁移指南

### 步骤1: 更新Go模块

#### 1.1 更新go.mod

```bash
# 编辑go.mod，更新版本
go get github.com/yourusername/golang@v2.0.0

# 清理依赖
go mod tidy
```

#### 1.2 验证模块

```bash
# 验证下载的模块
go mod verify

# 查看依赖图
go mod graph | grep github.com/yourusername/golang
```

---

### 步骤2: 更新导入路径

这是最主要的变更。需要批量替换所有导入路径。

#### 2.1 创建替换脚本

**Linux/macOS** - `migrate-imports.sh`:

```bash
#!/bin/bash

echo "开始迁移导入路径..."

# AI-Agent
find . -type f -name "*.go" -exec sed -i '' \
    's|examples/advanced/ai-agent/core|github.com/yourusername/golang/pkg/agent/core|g' {} +

# Concurrency
find . -type f -name "*.go" -exec sed -i '' \
    's|examples/concurrency|github.com/yourusername/golang/pkg/concurrency/patterns|g' {} +

# HTTP/3
find . -type f -name "*.go" -exec sed -i '' \
    's|examples/advanced/http3|github.com/yourusername/golang/pkg/http3|g' {} +

# Memory
find . -type f -name "*.go" -exec sed -i '' \
    's|examples/advanced/memory|github.com/yourusername/golang/pkg/memory|g' {} +

echo "迁移完成！请检查结果。"
```

**Windows** - `migrate-imports.ps1`:

```powershell
Write-Host "开始迁移导入路径..."

# AI-Agent
Get-ChildItem -Path . -Filter *.go -Recurse | ForEach-Object {
    (Get-Content $_.FullName) `
        -replace 'examples/advanced/ai-agent/core', 'github.com/yourusername/golang/pkg/agent/core' |
    Set-Content $_.FullName
}

# Concurrency
Get-ChildItem -Path . -Filter *.go -Recurse | ForEach-Object {
    (Get-Content $_.FullName) `
        -replace 'examples/concurrency', 'github.com/yourusername/golang/pkg/concurrency/patterns' |
    Set-Content $_.FullName
}

# HTTP/3
Get-ChildItem -Path . -Filter *.go -Recurse | ForEach-Object {
    (Get-Content $_.FullName) `
        -replace 'examples/advanced/http3', 'github.com/yourusername/golang/pkg/http3' |
    Set-Content $_.FullName
}

Write-Host "迁移完成！"
```

#### 2.2 手动检查特殊情况

某些导入可能需要手动调整：

```go
// ❌ v1.x
import (
    agent "path/to/examples/advanced/ai-agent/core"
    "local/custom/wrapper"
)

// ✅ v2.0
import (
    "github.com/yourusername/golang/pkg/agent/core"
    // 自定义包装器可能需要更新
)
```

---

### 步骤3: 更新API调用

#### 3.1 Observability模块

**Metrics注册**:

```go
// ❌ v1.x
counter := NewCounter("requests", "Total requests", nil)
Register(counter)  // 错误未处理

// ✅ v2.0 - 选项1: 显式忽略
counter := NewCounter("requests", "Total requests", nil)
_ = Register(counter)  // #nosec G104 - 显式忽略

// ✅ v2.0 - 选项2: 处理错误
counter := NewCounter("requests", "Total requests", nil)
if err := Register(counter); err != nil {
    log.Printf("Failed to register counter: %v", err)
}

// ✅ v2.0 - 选项3: 使用便捷函数（推荐）
counter := RegisterCounter("requests", "Total requests", nil)
```

**文件操作权限**:

```go
// ❌ v1.x
fileHook, err := NewFileHook("app.log", ErrorLevel, 0666)

// ✅ v2.0 - 更安全的权限
fileHook, err := NewFileHook("app.log", ErrorLevel, 0600)
```

#### 3.2 Agent模块

**随机数生成**:

```go
// ❌ v1.x - 可能触发安全警告
import "math/rand"
value := rand.Float64()

// ✅ v2.0 - 添加注释说明
import "math/rand"
// #nosec G404 - 用于非安全相关的探索策略
value := rand.Float64()

// 或者使用crypto/rand（如果需要更高安全性）
import "crypto/rand"
import "math/big"

func secureRandFloat() (float64, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(1000000))
    if err != nil {
        return 0, err
    }
    return float64(n.Int64()) / 1000000.0, nil
}
```

**文件路径验证**:

```go
// ❌ v1.x
data, err := os.ReadFile(filename)

// ✅ v2.0 - 添加路径验证
import "path/filepath"

func safeReadFile(baseDir, filename string) ([]byte, error) {
    cleanPath := filepath.Clean(filename)
    fullPath := filepath.Join(baseDir, cleanPath)
    
    // 验证路径在安全范围内
    if !strings.HasPrefix(fullPath, baseDir) {
        return nil, fmt.Errorf("invalid path: outside base directory")
    }
    
    return os.ReadFile(fullPath)
}
```

---

### 步骤4: 更新配置

#### 4.1 日志配置

```go
// ❌ v1.x
logger := NewLogger(InfoLevel, os.Stdout)

// ✅ v2.0 - 添加推荐的钩子
logger := NewLogger(InfoLevel, os.Stdout)
logger.AddHook(NewMetricsHook())  // 自动记录日志指标
```

#### 4.2 文件权限配置

如果你的代码中有硬编码的文件权限：

```go
// ❌ v1.x
const (
    LogFileMode = 0666  // 过于宽松
    ConfigMode  = 0644
)

// ✅ v2.0 - 更安全的权限
const (
    LogFileMode = 0600  // 只有所有者可读写
    ConfigMode  = 0600  // 配置文件也应该限制权限
    DirMode     = 0700  // 目录权限
)
```

---

### 步骤5: 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test github.com/yourusername/golang/pkg/agent/...

# 运行带竞态检测的测试
go test -race ./...

# 运行覆盖率测试
go test -cover ./...
```

---

### 步骤6: 更新文档

更新项目中的文档引用：

```markdown
<!-- ❌ 旧文档 -->
See [Agent Documentation](examples/advanced/ai-agent/README.md)

<!-- ✅ 新文档 -->
See [Agent Documentation](pkg/agent/README.md)
```

---

## 💥 Breaking Changes详解

### 1. 目录结构变更

#### 变更对照表

| v1.x | v2.0 | 说明 |
|------|------|------|
| `examples/advanced/ai-agent/` | `pkg/agent/` | Agent框架模块化 |
| `examples/concurrency/` | `pkg/concurrency/` | 并发模式库 |
| `examples/advanced/http3/` | `pkg/http3/` | HTTP/3服务器 |
| `examples/advanced/memory/` | `pkg/memory/` | 内存管理 |
| (新增) | `pkg/observability/` | 可观测性 |
| (新增) | `cmd/gox/` | CLI工具 |

#### 迁移影响

- **高影响**: 所有导入路径需要更新
- **迁移工具**: 提供自动化脚本
- **测试建议**: 完整的回归测试

### 2. API签名变更

#### Observability - Register函数

```go
// v1.x
func Register(metric Metric)

// v2.0
func Register(metric Metric) error

// 迁移策略
// 1. 处理错误（推荐）
if err := Register(metric); err != nil {
    // 处理错误
}

// 2. 显式忽略
_ = Register(metric)  // #nosec G104

// 3. 使用新的便捷函数
counter := RegisterCounter(name, help, labels)  // 自动忽略错误
```

### 3. 文件权限变更

#### 默认权限更加严格

```go
// v1.x
os.OpenFile(file, flags, 0666)  // rw-rw-rw-
os.MkdirAll(dir, 0755)          // rwxr-xr-x

// v2.0 (推荐)
os.OpenFile(file, flags, 0600)  // rw-------
os.MkdirAll(dir, 0700)          // rwx------
```

#### 迁移检查清单

- [ ] 检查所有`os.OpenFile`调用
- [ ] 检查所有`os.MkdirAll`调用
- [ ] 检查所有`os.WriteFile`调用
- [ ] 确保权限符合安全要求

---

## 📝 代码迁移示例

### 示例1: 完整的HTTP服务迁移

**v1.x代码**:

```go
package main

import (
    "log"
    agent "path/to/examples/advanced/ai-agent/core"
)

func main() {
    // 创建Agent
    myAgent := agent.NewBaseAgent("server")
    
    // 简单使用
    log.Println("Server started")
}
```

**v2.0代码**:

```go
package main

import (
    "context"
    "github.com/yourusername/golang/pkg/agent/core"
    "github.com/yourusername/golang/pkg/observability"
)

func main() {
    // 初始化可观测性
    logger := observability.NewLogger(observability.InfoLevel, nil)
    logger.AddHook(observability.NewMetricsHook())
    observability.SetDefaultLogger(logger)
    
    // 创建Agent
    myAgent := core.NewBaseAgent("server")
    
    // 使用追踪
    ctx := context.Background()
    span, ctx := observability.StartSpan(ctx, "server-start")
    defer span.Finish()
    
    observability.Info("Server started")
}
```

### 示例2: 并发模式迁移

**v1.x代码**:

```go
import "path/to/examples/concurrency"

func processJobs() {
    // 使用worker pool
    pool := concurrency.NewWorkerPool(10)
    // ...
}
```

**v2.0代码**:

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/concurrency/patterns"
    "github.com/yourusername/golang/pkg/observability"
)

func processJobs() {
    ctx := context.Background()
    
    // 使用追踪
    span, ctx := observability.StartSpan(ctx, "process-jobs")
    defer span.Finish()
    
    // 使用rate limiter
    limiter := patterns.NewTokenBucket(100, time.Second)
    
    // 使用worker pool with context
    jobs := make(chan Job, 100)
    results := patterns.WorkerPool(ctx, 10, jobs)
    
    // 处理结果
    for result := range results {
        // ...
    }
}
```

---

## ❓ 常见迁移问题

### Q1: 导入路径找不到

**问题**:

```text
cannot find package "examples/advanced/ai-agent/core"
```

**解决**:

```bash
# 1. 更新go.mod
go get github.com/yourusername/golang@v2.0.0

# 2. 更新导入
# 使用提供的迁移脚本

# 3. 清理缓存
go clean -modcache
go mod tidy
```

### Q2: 类型不匹配

**问题**:

```text
cannot use Register(counter) (type error) as type () in assignment
```

**解决**:

```go
// 方案1: 处理错误
if err := Register(counter); err != nil {
    log.Printf("Error: %v", err)
}

// 方案2: 忽略错误
_ = Register(counter)

// 方案3: 使用新API
counter := RegisterCounter("name", "help", nil)
```

### Q3: 测试失败

**问题**:

```text
permission denied when creating log file
```

**解决**:

```go
// 检查文件权限
// 从 0666 改为 0600

// 测试时使用临时目录
tmpDir := t.TempDir()
logFile := filepath.Join(tmpDir, "test.log")
```

### Q4: 性能下降

**问题**: 迁移后性能测试显示下降

**排查**:

```bash
# 运行基准测试
go test -bench=. -benchmem ./...

# 对比v1.x和v2.0的结果
# 检查是否有不必要的额外检查

# 使用pprof分析
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

---

## 🔙 回滚方案

如果迁移遇到严重问题，可以回滚：

### 方案1: Git回滚

```bash
# 回滚到迁移前的提交
git checkout pre-v2-migration

# 或创建新分支
git checkout -b v1-x-stable
```

### 方案2: 使用v1.x版本

```bash
# 在go.mod中固定v1.x版本
go get github.com/yourusername/golang@v1.9.0

# 恢复旧的导入路径
# ...

# 重新构建
go build ./...
```

### 方案3: 渐进式迁移

如果全量迁移风险大，可以：

1. **保持v1.x运行**
2. **创建新的v2.0服务**
3. **逐步切换流量**
4. **最终完全迁移**

---

## ✅ 迁移验证清单

完成迁移后，使用此清单验证：

- [ ] 所有导入路径已更新
- [ ] go.mod版本正确 (v2.0.0)
- [ ] 所有测试通过 (`go test ./...`)
- [ ] 基准测试性能符合预期
- [ ] 文档已更新
- [ ] CI/CD管道正常
- [ ] 生产环境验证通过
- [ ] 监控指标正常
- [ ] 无安全告警

---

## 📞 获取帮助

如果迁移过程中遇到问题：

1. **查看文档**: [Release Notes](RELEASE_NOTES_v2.0.0.md)
2. **搜索Issues**: [GitHub Issues](https://github.com/yourusername/golang/issues)
3. **提问**: [GitHub Discussions](https://github.com/yourusername/golang/discussions)
4. **紧急支持**: <your-email@example.com>

---

## 📚 相关资源

- [Release Notes v2.0.0](RELEASE_NOTES_v2.0.0.md)
- [安装指南](INSTALLATION.md)
- [快速开始](QUICK_START.md)
- [完整文档](docs/README.md)

---

**祝迁移顺利！** 🚀

如果成功完成迁移，欢迎分享你的经验！
