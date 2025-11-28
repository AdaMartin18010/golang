# Concurrency Pattern Generator (CPG)

**Version**: v1.0.0  
**Go Version**: 1.25.3  
**Status**: 🚀 Active Development  
**Theory**: 文档02 CSP并发模型 + 文档16 并发模式

---

## 📚 简介

Concurrency Pattern Generator（CPG）是一个基于**CSP形式化验证**的Go并发模式代码生成工具。它能够生成30+种经过形式化验证的并发模式代码，每个模式都附带：

- ✅ **CSP进程定义**
- ✅ **Happens-Before关系分析**
- ✅ **死锁自由证明**
- ✅ **数据竞争分析**
- ✅ **形式化注释**

---

## 🎯 核心功能

### 30+并发模式

#### 1. 经典模式 (5个)

- Worker Pool
- Fan-In
- Fan-Out
- Pipeline
- Generator

#### 2. 同步模式 (8个)

- Mutex Pattern
- RWMutex Pattern
- WaitGroup Pattern
- Once Pattern
- Cond Pattern
- Semaphore
- Barrier
- CountDownLatch

#### 3. 控制流模式 (5个)

- Context Cancellation
- Context Timeout
- Context WithValue
- Graceful Shutdown
- Rate Limiting

#### 4. 数据流模式 (7个)

- Producer-Consumer
- Buffered Channel
- Unbuffered Channel
- Select Pattern
- For-Select Loop
- Done Channel
- Error Channel

#### 5. 高级模式 (5个)

- Actor Model
- Session Types
- Future/Promise
- Map-Reduce
- Pub-Sub

---

## 🚀 快速开始

### 安装

```bash
go install github.com/your-org/concurrency-pattern-generator/cmd/cpg@latest
```

### 生成Worker Pool模式

```bash
cpg generate --pattern worker-pool --output worker_pool.go
```

### 交互式模式

```bash
cpg interactive
```

### 批量生成

```bash
cpg batch --config patterns.yaml
```

---

## 💡 使用示例

### 示例1: 生成Worker Pool

```bash
$ cpg generate --pattern worker-pool --workers 10 --output pool.go

✅ 生成成功: pool.go
📊 统计:
   - 代码行数: 85
   - CSP验证: ✓ 通过
   - 安全性: ✓ 死锁自由, ✓ 竞争自由
```

生成的代码：

```go
// Pattern: Worker Pool
// CSP Model: Pool = (worker₁ || worker₂ || ... || worker₁₀)
// Safety Properties:
//   - Deadlock-free: ✓ (All workers can terminate)
//   - Race-free: ✓ (Channel synchronization)
// Theory: 文档02 第3.2节, 文档16 第1.1节
//
// Happens-Before Relations:
//   1. job sent → job received
//   2. result computed → result sent
//   3. done closed → all workers exit

package main

import (
 "context"
 "sync"
)

type Job struct {
 ID   int
 Data interface{}
}

type Result struct {
 JobID int
 Data  interface{}
 Error error
}

func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
 results := make(chan Result)
 var wg sync.WaitGroup

 // Start workers
 for i := 0; i < numWorkers; i++ {
  wg.Add(1)
  go func(workerID int) {
   defer wg.Done()
   for {
    select {
    case <-ctx.Done():
     return
    case job, ok := <-jobs:
     if !ok {
      return
     }
     // Process job
     result := processJob(job)
     select {
     case results <- result:
     case <-ctx.Done():
      return
     }
    }
   }
  }(i)
 }

 // Close results when all workers done
 go func() {
  wg.Wait()
  close(results)
 }()

 return results
}

func processJob(job Job) Result {
 // User implements this
 return Result{JobID: job.ID, Data: job.Data}
}
```

---

### 示例2: Fan-Out + Fan-In

```bash
cpg generate --pattern fan-out-in --fanout 5 --output fanout.go
```

---

## 📐 形式化理论基础

### CSP模型

每个模式都有对应的CSP进程定义：

```text
Worker Pool:
  Pool = worker₁ || worker₂ || ... || workerₙ
  workerᵢ = jobs?job → process(job) → results!result → workerᵢ
          □ done → SKIP

Pipeline:
  Pipeline = stage₁ >> stage₂ >> ... >> stageₙ
  stageᵢ = input?x → process(x) → output!y → stageᵢ
```

### 安全性验证

1. **Deadlock Freedom**
   - 证明所有进程可以终止
   - 无循环依赖

2. **Race Freedom**
   - Channel同步保证
   - Happens-Before关系

3. **Liveness**
   - 最终所有消息被处理
   - 无饿死

---

## 🎨 高级特性

### 1. 自定义模板

```yaml
# custom.yaml
pattern: custom-pool
workers: 20
buffer_size: 100
timeout: 30s
error_handling: retry
max_retries: 3
```

```bash
cpg generate --config custom.yaml
```

### 2. 性能优化建议

CPG会分析并提供：

- Buffer大小建议
- Worker数量优化
- 性能瓶颈警告

### 3. 代码分析

```bash
cpg verify --file worker_pool.go

📊 分析结果:
   ✅ CSP模型: 符合Worker Pool定义
   ✅ 死锁检查: 通过
   ✅ 竞争检查: 通过
   ⚠️  性能建议: 考虑增加buffer大小到50
```

---

## 📖 文档

### 理论文档

- [文档02: CSP并发模型与形式化证明](../../docs/01-语言基础/00-Go-1.25.3形式化理论体系/02-CSP并发模型与形式化证明.md)
- [文档16: Go并发模式完整形式化分析](../../docs/01-语言基础/00-Go-1.25.3形式化理论体系/16-Go并发模式完整形式化分析-2025.md)

### 使用手册

- [快速开始](docs/quick-start.md)
- [模式参考](docs/patterns.md)
- [配置指南](docs/configuration.md)
- [最佳实践](docs/best-practices.md)

---

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定模式测试
go test ./pkg/patterns -run TestWorkerPool

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🏗️ 项目结构

```text
concurrency-pattern-generator/
├── cmd/cpg/
│   └── main.go              # CLI主程序
├── pkg/
│   ├── generator/
│   │   ├── generator.go     # 生成器核心
│   │   └── template.go      # 模板引擎
│   ├── patterns/
│   │   ├── classic.go       # 经典模式
│   │   ├── sync.go          # 同步模式
│   │   ├── control.go       # 控制流模式
│   │   ├── dataflow.go      # 数据流模式
│   │   └── advanced.go      # 高级模式
│   └── verifier/
│       ├── csp.go           # CSP验证
│       └── safety.go        # 安全性检查
├── templates/               # 模式模板
├── testdata/               # 测试数据
├── docs/                   # 文档
├── go.mod
└── README.md
```

---

## 🎓 学习资源

### 对初学者

1. 从Worker Pool开始
2. 理解CSP模型
3. 阅读生成的注释
4. 运行测试用例

### 对工程师

1. 自定义模板
2. 性能优化
3. 错误处理策略
4. 生产环境部署

### 对研究者

1. CSP形式化方法
2. 并发模式抽象
3. 安全性验证
4. 代码生成技术

---

## 💬 贡献

欢迎贡献！

- 报告Bug
- 添加新模式
- 改进文档
- 提供反馈

---

## 📜 许可证

MIT License

---

## 🙏 致谢

基于以下理论研究：

- CSP (Communicating Sequential Processes) - Tony Hoare
- Go Memory Model - Go Team
- Go 1.25.3形式化理论体系

---

<div align="center">

**Concurrency Pattern Generator**-

基于形式化验证的Go并发代码生成工具

**[开始使用](#🚀-快速开始)** | **[查看模式](#30并发模式)** | **[理论基础](#📐-形式化理论基础)**

Made with ❤️ for Go Concurrency

</div>
