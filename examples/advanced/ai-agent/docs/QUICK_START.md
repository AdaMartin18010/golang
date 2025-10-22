# AI-Agent 快速开始

> 5分钟快速上手AI-Agent系统

---

## 📦 安装

### 前提条件

```bash
# 需要 Go 1.23+
go version

# 确保在项目根目录
cd golang
```

### 获取代码

```bash
# 代码已在 examples/advanced/ai-agent/
cd examples/advanced/ai-agent
```

---

## 🚀 第一个Agent

### 1. 创建简单Agent

```go
package main

import (
    "context"
    "fmt"
    "github.com/yourusername/golang/examples/advanced/ai-agent/core"
)

func main() {
    // 1. 创建Agent配置
    config := core.AgentConfig{
        Name: "MyFirstAgent",
        Type: "assistant",
    }
    
    // 2. 创建Agent实例
    agent := core.NewBaseAgent("agent-001", config)
    
    // 3. 初始化组件
    agent.SetLearningEngine(core.NewLearningEngine(nil))
    agent.SetDecisionEngine(core.NewDecisionEngine(nil))
    
    // 4. 启动Agent
    ctx := context.Background()
    if err := agent.Start(ctx); err != nil {
        panic(err)
    }
    defer agent.Stop()
    
    // 5. 处理任务
    input := core.Input{
        ID:   "task-1",
        Type: "text",
        Data: map[string]interface{}{
            "message": "Hello, Agent!",
        },
    }
    
    output, err := agent.Process(input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Agent Response: %+v\n", output)
}
```

### 2. 运行程序

```bash
# 编译运行
go run main.go

# 或者运行测试
go test -v ./...
```

---

## 💡 核心功能示例

### 决策引擎使用

```go
// 创建决策引擎
decisionEngine := core.NewDecisionEngine(&core.DecisionConfig{
    Strategy: "rule-based",
    Timeout:  time.Second * 5,
})

// 添加决策规则
decisionEngine.AddRule(core.Rule{
    Condition: func(input core.Input) bool {
        return input.Type == "urgent"
    },
    Action: func(input core.Input) core.Decision {
        return core.Decision{
            Action:   "immediate",
            Priority: core.HighPriority,
        }
    },
})

// 执行决策
decision := decisionEngine.Decide(input)
fmt.Printf("Decision: %s (Priority: %d)\n", decision.Action, decision.Priority)
```

### 学习引擎使用

```go
// 创建学习引擎
learningEngine := core.NewLearningEngine(&core.LearningConfig{
    LearningRate: 0.01,
    BufferSize:   1000,
})

// 记录经验
learningEngine.Learn(core.Experience{
    State:  currentState,
    Action: takenAction,
    Reward: 1.0,
    NextState: nextState,
})

// 评估策略
score := learningEngine.EvaluatePolicy()
fmt.Printf("Policy Score: %.2f\n", score)
```

### 多模态处理

```go
// 创建多模态接口
multimodal := core.NewMultimodalInterface(&core.MultimodalConfig{
    EnableText:  true,
    EnableAudio: true,
    EnableImage: true,
})

// 处理文本
textFeatures := multimodal.ProcessText("Hello, world!")

// 处理音频
audioFeatures := multimodal.ProcessAudio(audioData)

// 融合特征
fusedFeatures := multimodal.FuseFeatures(textFeatures, audioFeatures)
```

---

## 🎯 实战示例

### 示例1: 智能客服Agent

```go
// 创建客服Agent
customerServiceAgent := core.NewBaseAgent("cs-agent", core.AgentConfig{
    Name: "CustomerService",
    Type: "service",
})

// 配置决策引擎 - 规则化响应
decisionEngine := core.NewDecisionEngine(&core.DecisionConfig{
    Strategy: "rule-based",
})

// 添加客服规则
decisionEngine.AddRule(core.Rule{
    Condition: func(input core.Input) bool {
        msg, _ := input.Data["message"].(string)
        return strings.Contains(msg, "退款")
    },
    Action: func(input core.Input) core.Decision {
        return core.Decision{
            Action: "handle_refund",
            Data: map[string]interface{}{
                "response": "我会帮您处理退款申请...",
            },
        }
    },
})

customerServiceAgent.SetDecisionEngine(decisionEngine)

// 处理客户请求
input := core.Input{
    ID:   "req-001",
    Type: "customer_query",
    Data: map[string]interface{}{
        "message": "我要退款",
    },
}

response, _ := customerServiceAgent.Process(input)
fmt.Println(response.Data["response"])
```

### 示例2: 任务自动化Agent

```go
// 创建自动化Agent
automationAgent := core.NewBaseAgent("auto-agent", core.AgentConfig{
    Name: "Automation",
    Type: "automation",
})

// 配置学习引擎 - 优化执行策略
learningEngine := core.NewLearningEngine(&core.LearningConfig{
    LearningRate: 0.01,
})

automationAgent.SetLearningEngine(learningEngine)

// 执行任务并学习
for i := 0; i < 100; i++ {
    task := generateTask(i)
    result := automationAgent.Process(task)
    
    // 根据结果学习
    reward := evaluateResult(result)
    learningEngine.Learn(core.Experience{
        State:  task,
        Action: result.Action,
        Reward: reward,
    })
}

// 评估学习效果
finalScore := learningEngine.EvaluatePolicy()
fmt.Printf("Final Policy Score: %.2f\n", finalScore)
```

### 示例3: 多Agent协作

```go
// 创建多个Agent
agents := []core.Agent{
    core.NewBaseAgent("agent-1", core.AgentConfig{Name: "Analyzer"}),
    core.NewBaseAgent("agent-2", core.AgentConfig{Name: "Executor"}),
    core.NewBaseAgent("agent-3", core.AgentConfig{Name: "Monitor"}),
}

// 启动所有Agent
for _, agent := range agents {
    agent.Start(context.Background())
}

// 创建协调决策引擎
consensusEngine := core.NewDecisionEngine(&core.DecisionConfig{
    Strategy: "consensus",
})

// 执行共识决策
task := core.Input{
    ID:   "collaborative-task",
    Type: "complex",
    Data: map[string]interface{}{
        "description": "需要多个Agent协作",
    },
}

decision := consensusEngine.ConsensusDecision(agents, task)
fmt.Printf("Consensus Decision: %s\n", decision.Action)
```

---

## 🧪 运行测试

### 运行所有测试

```bash
cd examples/advanced/ai-agent
go test -v ./...
```

### 运行特定测试

```bash
# 测试决策引擎
go test -v ./core -run TestDecisionEngine

# 测试学习引擎
go test -v ./core -run TestLearningEngine

# 测试多模态接口
go test -v ./core -run TestMultimodalInterface
```

### 查看覆盖率

```bash
go test -cover ./...

# 生成详细覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📊 性能测试

### 基准测试

```bash
# 运行基准测试
go test -bench=. -benchmem ./...

# 决策性能
go test -bench=BenchmarkDecision ./core

# 学习性能
go test -bench=BenchmarkLearning ./core
```

### 性能分析

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

---

## 🔧 配置选项

### Agent配置

```go
type AgentConfig struct {
    Name     string        // Agent名称
    Type     string        // Agent类型
    Timeout  time.Duration // 超时时间
    MaxTasks int           // 最大任务数
}
```

### 决策引擎配置

```go
type DecisionConfig struct {
    Strategy   string        // 决策策略: rule-based, probabilistic, consensus
    Timeout    time.Duration // 决策超时
    MaxRetries int           // 最大重试次数
}
```

### 学习引擎配置

```go
type LearningConfig struct {
    LearningRate float64 // 学习率
    BufferSize   int     // 经验缓冲区大小
    BatchSize    int     // 批量学习大小
}
```

---

## 📚 下一步

### 深入学习

1. **[架构文档](ARCHITECTURE.md)** - 了解系统架构
2. **[API文档](API.md)** - 查看完整API
3. **[示例集合](EXAMPLES.md)** - 更多实战示例

### 扩展功能

- 自定义决策算法
- 实现新的学习策略
- 添加新的模态处理
- 构建多Agent系统

### 最佳实践

- 合理设置超时时间
- 监控Agent状态
- 定期评估性能
- 优化并发处理

---

## ❓ 常见问题

### Q: 如何处理超时？

A: 使用Context控制:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

agent.Start(ctx)
```

### Q: 如何实现持久化？

A: 实现StateStore接口:

```go
type StateStore interface {
    Save(state AgentState) error
    Load() (AgentState, error)
}
```

### Q: 如何监控Agent性能？

A: 使用Metrics:

```go
metrics := agent.GetMetrics()
fmt.Printf("Processed: %d, Success: %d\n", 
    metrics.TotalProcessed, 
    metrics.SuccessCount)
```

---

## 📞 获取帮助

- 📖 [完整文档](../README.md)
- 🐛 [提交Issue](https://github.com/yourusername/golang/issues)
- 💬 [讨论区](https://github.com/yourusername/golang/discussions)

---

**快速开始版本**: v1.0  
**最后更新**: 2025-10-22
