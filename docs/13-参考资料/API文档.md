# API 文档

> **项目**: Golang Learning Project  
> **版本**: v2.0  
> **最后更新**: 2025-10-22

---

## 📚 模块概览

本项目提供以下可复用的Go包：

1. **pkg/agent** - AI代理框架
2. **pkg/concurrency** - 并发模式
3. **pkg/http3** - HTTP/3服务器
4. **pkg/memory** - 内存管理
5. **pkg/observability** - 可观测性

---

## 📦 pkg/agent - AI代理框架

### 导入

```go
import "github.com/yourusername/golang/pkg/agent/core"
```

### 核心类型

#### Agent 接口

```go
type Agent interface {
    ID() string
    Start(ctx context.Context) error
    Stop() error
    Process(input Input) (Output, error)
    Learn(experience Experience) error
    GetStatus() Status
}
```

#### BaseAgent

基础代理实现。

**创建**:

```go
config := core.AgentConfig{
    Name:         "MyAgent",
    Type:         "worker",
    MaxLoad:      0.8,
    Timeout:      5 * time.Second,
    Retries:      3,
    Capabilities: []string{"processing"},
}

agent := core.NewBaseAgent("agent-1", config)
```

**使用示例**:

```go
ctx := context.Background()

// 启动代理
if err := agent.Start(ctx); err != nil {
    log.Fatal(err)
}
defer agent.Stop()

// 处理输入
input := core.Input{
    ID:   "task-1",
    Type: "process",
    Data: map[string]interface{}{"value": 42},
}

output, err := agent.Process(input)
if err != nil {
    log.Fatal(err)
}
```

### DecisionEngine

决策引擎用于选择合适的代理处理任务。

```go
engine := core.NewDecisionEngine(nil)

// 注册代理
engine.RegisterAgent(&myAgent)

// 做出决策
decision := engine.Decide(input)
```

### LearningEngine

学习引擎用于从经验中学习。

```go
learner := core.NewLearningEngine(nil)

experience := core.Experience{
    Input:  input,
    Output: output,
    Reward: 1.0,
}

learner.Learn(experience)
```

---

## 📦 pkg/concurrency - 并发模式

### 导入1

```go
import "github.com/yourusername/golang/pkg/concurrency/patterns"
```

### Pipeline模式

```go
// 生成器
gen := func(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// 处理器
sq := func(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// 使用
nums := gen(1, 2, 3, 4)
squared := sq(nums)

for n := range squared {
    fmt.Println(n) // 1, 4, 9, 16
}
```

### Worker Pool模式

```go
const numWorkers = 5
jobs := make(chan int, 100)
results := make(chan int, 100)

// 启动workers
for w := 1; w <= numWorkers; w++ {
    go worker(jobs, results)
}

// 发送任务
for j := 1; j <= 10; j++ {
    jobs <- j
}
close(jobs)

// 收集结果
for r := 1; r <= 10; r++ {
    result := <-results
    fmt.Println(result)
}
```

---

## 📦 pkg/http3 - HTTP/3服务器

### 导入2

```go
import "github.com/yourusername/golang/pkg/http3"
```

### 基本使用

```go
mux := http.NewServeMux()
mux.HandleFunc("/", handleRoot)
mux.HandleFunc("/health", handleHealth)

server := &http.Server{
    Addr:    ":8443",
    Handler: mux,
}

log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
```

### 响应结构

```go
type Response struct {
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Protocol  string    `json:"protocol"`
    Server    string    `json:"server"`
}
```

---

## 📦 pkg/memory - 内存管理

### Arena Allocator

用于批量短生命周期对象的内存管理。

```go
records := []Record{
    {ID: 1, Name: "A", Value: 10.0},
    {ID: 2, Name: "B", Value: 20.0},
}

results := processWithArena(records)
```

### Weak Pointer Cache

使用弱引用避免内存泄漏的缓存。

```go
cache := NewWeakCache()

// 设置值
value := &Value{
    Data: "cached data",
    Size: 1024,
}
cache.Set("key1", value)

// 获取值
if val, found := cache.Get("key1"); found {
    fmt.Println(val.Data)
}

// 清理过期条目
cache.Cleanup()
```

---

## 📦 pkg/observability - 可观测性

### OTel配置

项目使用OpenTelemetry进行可观测性。

**环境变量**:

- `OTEL_SERVICE_NAME`: 服务名称
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OTLP端点
- `OTEL_EXPORTER_OTLP_INSECURE`: 是否使用不安全连接

**示例**:

```bash
export OTEL_SERVICE_NAME=my-service
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
export OTEL_EXPORTER_OTLP_INSECURE=true
```

---

## 🛠️ CLI工具 - gox

### 安装

```bash
go build -o gox ./cmd/gox
```

### 命令

```bash
# 运行测试
gox test

# 构建项目
gox build

# 生成覆盖率报告
gox coverage

# 项目统计
gox stats

# 代码质量检查
gox quality

# 版本信息
gox version
```

---

## 📖 完整文档

### 在线文档

使用godoc查看完整API文档：

```bash
# 安装pkgsite
go install golang.org/x/pkgsite/cmd/pkgsite@latest

# 启动文档服务器
pkgsite -http=:8080

# 访问 http://localhost:8080
```

### 生成文档

```bash
# 生成HTML文档
godoc -http=:6060

# 访问 http://localhost:6060/pkg/
```

---

## 🔗 相关链接

- **项目首页**: [README.md](README.md)
- **快速开始**: [docs/QUICK_START.md](docs/QUICK_START.md)
- **学习路径**: [docs/LEARNING_PATHS.md](docs/LEARNING_PATHS.md)
- **文档索引**: [docs/INDEX.md](docs/INDEX.md)

---

## 📝 示例代码

完整的示例代码位于 `examples/` 目录：

- `examples/advanced/ai-agent/` - AI代理示例
- `examples/concurrency/` - 并发模式示例
- `examples/advanced/http3-server/` - HTTP/3服务器示例
- `examples/observability/` - 可观测性示例

---

## 🤝 贡献

欢迎贡献代码和文档！请查看 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 📞 支持

- **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-repo/discussions)
- **文档**: [docs/](docs/)

---

**生成时间**: 2025-10-22  
**API版本**: v2.0  
**Go版本**: 1.25.3+
