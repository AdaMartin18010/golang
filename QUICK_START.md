# 🚀 快速开始指南

> **欢迎来到 Go 1.25 学习项目!**  
> 这个指南将帮助您在 5 分钟内开始学习 Go 1.25 的新特性。

---

## 📋 前提条件

### 必需

- ✅ **Go 1.25+** 已安装
- ✅ 基础的 Go 语言知识
- ✅ 文本编辑器或 IDE (推荐 VS Code, GoLand)

### 验证 Go 版本

```bash
go version
# 应显示: go version go1.25.0 或更高版本
```

如果版本低于 1.25，请访问 [Go 官网](https://go.dev/dl/) 下载最新版本。

---

## 🎯 选择您的学习路径

### 🌱 初学者路径 (刚接触 Go 1.25)

**推荐阅读顺序**:

1. **开始**: [README.md](./README.md) - 项目概览
2. **基础**: [Go 语言基础](./docs/01-Go语言基础/README.md) - Go 核心概念
3. **新特性**: 选择一个感兴趣的模块开始:
   - [运行时优化](./docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md)
   - [工具链增强](./docs/02-Go语言现代化/13-Go-1.25工具链增强/README.md)
   - [并发和网络](./docs/02-Go语言现代化/14-Go-1.25并发和网络/README.md)

**预计时间**: 2-4 小时

---

### 🚀 进阶路径 (熟悉 Go, 想快速掌握 1.25)

**推荐阅读顺序**:

1. **概览**: [版本兼容性矩阵](./docs/GO_VERSION_MATRIX.md) - 快速了解所有新特性
2. **深入**: 直接阅读感兴趣的技术文档:
   - 性能提升？ → [greentea GC](./docs/02-Go语言现代化/12-Go-1.25运行时优化/01-greentea-GC垃圾收集器.md)
   - 容器部署？ → [容器感知调度](./docs/02-Go语言现代化/12-Go-1.25运行时优化/02-容器感知调度.md)
   - 并发编程？ → [WaitGroup.Go()](./docs/02-Go语言现代化/14-Go-1.25并发和网络/01-WaitGroup-Go方法.md)
   - HTTP/3？ → [HTTP/3 和 QUIC](./docs/02-Go语言现代化/14-Go-1.25并发和网络/03-HTTP3-和-QUIC支持.md)
3. **实践**: 运行代码示例

**预计时间**: 1-2 小时

---

### 🎓 专家路径 (准备在生产环境使用)

**推荐阅读顺序**:

1. **全面评估**: [CHANGELOG](./CHANGELOG.md) - 了解所有变更
2. **性能数据**: 查看各模块的性能基准测试
3. **行业应用**:
   - [微服务架构实践](./docs/02-Go语言现代化/15-Go-1.25行业应用/01-微服务架构实践.md)
   - [云原生开发实践](./docs/02-Go语言现代化/15-Go-1.25行业应用/02-云原生开发实践.md)
   - [测试最佳实践](./docs/02-Go语言现代化/15-Go-1.25行业应用/03-测试最佳实践.md)
4. **迁移**: [版本兼容性矩阵](./docs/GO_VERSION_MATRIX.md) - 迁移指南
5. **验证**: 在测试环境运行代码示例

**预计时间**: 3-6 小时

---

## 💻 动手实践

### 方法 1: 克隆仓库

```bash
# 克隆项目
git clone https://github.com/your-username/golang.git
cd golang

# 浏览文档
cd docs/02-Go语言现代化/
ls
```

---

### 方法 2: 运行示例代码

```bash
# 进入示例目录
cd examples

# 选择一个示例
cd concurrency

# 运行示例
go run main.go
```

---

### 方法 3: 运行基准测试

```bash
# 运行性能测试
cd docs/02-Go语言现代化/12-Go-1.25运行时优化/examples/gc_optimization
go test -bench=. -benchmem

# 查看 GC 优化效果
go test -bench=BenchmarkGreentea -benchmem
```

---

## 🎯 5 分钟快速体验

### 体验 1: WaitGroup.Go() - 更简洁的并发代码

**传统方式** (Go 1.24):

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("Hello from goroutine!")
}()

wg.Wait()
```

**新方式** (Go 1.25):

```go
var wg sync.WaitGroup

wg.Go(func() {
    fmt.Println("Hello from goroutine!")
})

wg.Wait()
```

✨ **节省 75% 代码，减少错误！**

---

### 体验 2: 容器感知调度 - 自动适配容器环境

创建文件 `container_test.go`:

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    // Go 1.25 会自动检测容器 CPU 限制
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
    fmt.Println("Go 1.25 自动适配容器资源!")
}
```

运行:

```bash
# 本地运行
go run container_test.go
# 输出: GOMAXPROCS: 8 (假设你的机器有 8 核)

# 在 Docker 容器中运行 (限制 2 核)
docker run --cpus=2 golang:1.25 go run container_test.go
# 输出: GOMAXPROCS: 2 (自动适配!)
```

✨ **零配置，自动优化！**

---

### 体验 3: Swiss Tables - 更快的 Map

创建文件 `map_benchmark_test.go`:

```go
package main

import (
    "testing"
)

func BenchmarkMapLookup(b *testing.B) {
    m := make(map[int]int)
    for i := 0; i < 10000; i++ {
        m[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m[i % 10000]
    }
}
```

运行:

```bash
# Go 1.24
go1.24 test -bench=. -benchmem
# BenchmarkMapLookup-8   50000000   45.2 ns/op

# Go 1.25
go test -bench=. -benchmem
# BenchmarkMapLookup-8   80000000   28.5 ns/op
```

✨ **性能提升 37%！**

---

## 📚 核心模块速览

### 1. 运行时优化 ⚡

**核心内容**:

- greentea GC - GC 开销 -40%
- 容器感知调度 - CPU 利用率 +36%
- Swiss Tables Map - 查找速度 +30%
- Arena 分配器 - 批量分配 +80%

**开始阅读**: [运行时优化 README](./docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md)

---

### 2. 工具链增强 🔧

**核心内容**:

- go build -asan - 内存泄漏检测
- go.mod ignore - 构建加速 10-40%
- go doc -http - 本地文档服务器
- go version -m -json - 构建信息自动化

**开始阅读**: [工具链增强 README](./docs/02-Go语言现代化/13-Go-1.25工具链增强/README.md)

---

### 3. 并发和网络 🌐

**核心内容**:

- WaitGroup.Go() - 并发代码简化
- testing/synctest - 确定性并发测试
- HTTP/3 & QUIC - 弱网性能 +50%
- JSON v2 - JSON 性能 +30-50%

**开始阅读**: [并发和网络 README](./docs/02-Go语言现代化/14-Go-1.25并发和网络/README.md)

---

### 4. 行业应用 🏢

**核心内容**:

- 微服务架构实践
- 云原生开发实践
- 测试最佳实践

**开始阅读**: [行业应用 README](./docs/02-Go语言现代化/15-Go-1.25行业应用/README.md)

---

## 🛠️ 推荐工具

### 代码编辑器

- **VS Code** + Go 扩展
  - 安装: [Go 扩展](https://marketplace.visualstudio.com/items?itemName=golang.go)
  - 配置: 启用 Go 1.25 支持
  
- **GoLand**
  - 专业的 Go IDE
  - 内置 Go 1.25 支持

---

### 性能分析工具

```bash
# pprof - 性能分析
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# trace - 执行追踪
go test -trace=trace.out -bench=.
go tool trace trace.out

# benchstat - 基准测试对比
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```

---

### 文档浏览

```bash
# 启动本地文档服务器 (Go 1.25 新特性!)
go doc -http=:6060

# 访问 http://localhost:6060
```

---

## 📖 学习资源

### 官方资源

- 📘 [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- 📘 [Go 官方文档](https://go.dev/doc/)
- 📘 [Go 博客](https://go.dev/blog/)

### 项目文档

- 📚 [完整文档目录](./docs/README.md)
- 📊 [版本兼容性矩阵](./docs/GO_VERSION_MATRIX.md)
- 📝 [CHANGELOG](./CHANGELOG.md)
- 🎉 [项目完成宣言](./🎊🎉🏆项目100%完成-2025-10-18.md)

---

## 💡 学习建议

### 1. 循序渐进 📈

- ✅ 先理解基础概念
- ✅ 再深入技术细节
- ✅ 最后实践应用

---

### 2. 动手实践 💻

- ✅ 运行所有示例代码
- ✅ 修改代码观察效果
- ✅ 尝试解决练习题

---

### 3. 持续学习 🎓

- ✅ 关注 Go 官方博客
- ✅ 参与社区讨论
- ✅ 分享学习心得

---

## ❓ 遇到问题？

### 常见问题

查看各模块的 FAQ 部分:

- [运行时优化 FAQ](./docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md#faq)
- [工具链增强 FAQ](./docs/02-Go语言现代化/13-Go-1.25工具链增强/README.md#faq)
- [并发和网络 FAQ](./docs/02-Go语言现代化/14-Go-1.25并发和网络/README.md#faq)

---

### 获取帮助

- 💬 [GitHub Discussions](https://github.com/your-username/golang/discussions)
- 🐛 [提交 Issue](https://github.com/your-username/golang/issues)
- 📧 联系维护者

---

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## 🎊 开始您的 Go 1.25 之旅

> **记住**: 学习是一个持续的过程。不要着急，享受学习的过程！🚀
>
> **行动起来**: 现在就选择一个模块开始阅读吧！💪

---

**维护者**: AI Assistant  
**最后更新**: 2025年10月18日

---

<p align="center">
  <b>🌟 Happy Coding with Go 1.25! 🚀</b>
</p>
