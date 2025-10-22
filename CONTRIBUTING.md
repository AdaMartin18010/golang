# 贡献指南

感谢您对本项目的关注！我们欢迎各种形式的贡献。

---

## 📋 目录

- [贡献指南](#贡献指南)
  - [📋 目录](#-目录)
  - [🤝 行为准则](#-行为准则)
    - [我们的承诺](#我们的承诺)
    - [我们的标准](#我们的标准)
  - [💡 如何贡献](#-如何贡献)
    - [报告Bug](#报告bug)
    - [建议新功能](#建议新功能)
    - [改进文档](#改进文档)
  - [🔄 开发流程](#-开发流程)
    - [1. Fork和Clone](#1-fork和clone)
    - [2. 创建分支](#2-创建分支)
    - [3. 开发](#3-开发)
    - [4. 提交更改](#4-提交更改)
    - [5. 创建Pull Request](#5-创建pull-request)
  - [📝 代码规范](#-代码规范)
    - [Go代码风格](#go代码风格)
    - [命名规范](#命名规范)
    - [注释规范](#注释规范)
    - [错误处理](#错误处理)
  - [🎯 提交规范](#-提交规范)
    - [格式](#格式)
    - [Type类型](#type类型)
    - [示例](#示例)
  - [🧪 测试要求](#-测试要求)
    - [单元测试](#单元测试)
    - [基准测试](#基准测试)
    - [覆盖率要求](#覆盖率要求)
    - [运行测试](#运行测试)
  - [🔍 代码审查](#-代码审查)
    - [审查清单](#审查清单)
    - [审查过程](#审查过程)
  - [🎓 学习资源](#-学习资源)
    - [Go语言](#go语言)
    - [并发编程](#并发编程)
    - [测试](#测试)
  - [📞 获取帮助](#-获取帮助)
  - [🏆 贡献者](#-贡献者)

---

## 🤝 行为准则

### 我们的承诺

为了营造一个开放和友好的环境，我们作为贡献者和维护者承诺：无论年龄、体型、残疾、种族、性别认同和表达、经验水平、国籍、个人外貌、种族、宗教或性认同和性取向如何，参与我们的项目和社区对每个人来说都是无骚扰的体验。

### 我们的标准

积极行为的例子：

- ✅ 使用友好和包容的语言
- ✅ 尊重不同的观点和经验
- ✅ 优雅地接受建设性批评
- ✅ 关注对社区最有利的事情
- ✅ 对其他社区成员表示同情

不可接受行为的例子：

- ❌ 使用性化的语言或图像
- ❌ 挑衅、侮辱/贬损性评论，以及人身或政治攻击
- ❌ 公开或私下骚扰
- ❌ 未经明确许可，发布他人的私人信息
- ❌ 在专业环境中可能被认为不适当的其他行为

---

## 💡 如何贡献

### 报告Bug

如果您发现了bug，请：

1. **检查已有Issue**: 确保该bug尚未被报告
2. **创建新Issue**: 使用Bug报告模板
3. **提供详细信息**:
   - 清晰的标题和描述
   - 重现步骤
   - 预期行为和实际行为
   - 环境信息（Go版本、OS等）
   - 代码示例或错误日志

### 建议新功能

如果您有新功能的想法：

1. **检查已有Issue**: 确保功能尚未被建议
2. **创建Feature Request**: 描述功能和使用场景
3. **讨论设计**: 等待维护者反馈
4. **实现功能**: 获得批准后开始开发

### 改进文档

文档改进总是受欢迎的：

- 修复拼写错误
- 改进示例
- 添加缺失的说明
- 翻译文档

---

## 🔄 开发流程

### 1. Fork和Clone

```bash
# Fork仓库
# 在GitHub上点击Fork按钮

# Clone你的fork
git clone https://github.com/your-username/golang.git
cd golang

# 添加upstream远程仓库
git remote add upstream https://github.com/original-owner/golang.git
```

### 2. 创建分支

```bash
# 从main创建新分支
git checkout -b feature/your-feature-name

# 或者
git checkout -b fix/your-bug-fix
```

分支命名规范：

- `feature/` - 新功能
- `fix/` - Bug修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 3. 开发

```bash
# 安装依赖
go mod download

# 运行测试
go test ./...

# 运行质量检查
go fmt ./...
go vet ./...
golangci-lint run
```

### 4. 提交更改

```bash
# 暂存更改
git add .

# 提交（遵循提交规范）
git commit -m "feat: add new concurrency pattern"

# 推送到你的fork
git push origin feature/your-feature-name
```

### 5. 创建Pull Request

1. 访问GitHub上的原始仓库
2. 点击"New Pull Request"
3. 选择你的分支
4. 填写PR模板
5. 等待审查

---

## 📝 代码规范

### Go代码风格

遵循官方Go代码风格：

```go
// ✅ 好的示例
func CalculateSum(numbers []int) int {
    sum := 0
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// ❌ 不好的示例
func calculate_sum(numbers []int) int {
    Sum := 0
    for i := 0; i < len(numbers); i++ {
        Sum = Sum + numbers[i]
    }
    return Sum
}
```

### 命名规范

- **包名**: 小写，简短，不使用下划线
- **导出函数**: 大写开头，驼峰命名
- **私有函数**: 小写开头，驼峰命名
- **常量**: 驼峰命名（不是全大写）
- **接口**: 以`-er`结尾（如`Reader`, `Writer`）

### 注释规范

```go
// Package patterns provides common concurrency patterns.
package patterns

// RateLimiter implements a token bucket rate limiter.
// It allows controlling the rate of operations.
type RateLimiter struct {
    rate     int
    capacity int
}

// NewRateLimiter creates a new rate limiter with the given rate and capacity.
// rate is the number of tokens added per second.
// capacity is the maximum number of tokens the bucket can hold.
func NewRateLimiter(rate, capacity int) *RateLimiter {
    return &RateLimiter{
        rate:     rate,
        capacity: capacity,
    }
}
```

### 错误处理

```go
// ✅ 好的示例
func ProcessData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("empty data")
    }
    
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("parse data: %w", err)
    }
    
    return saveResult(result)
}

// ❌ 不好的示例
func ProcessData(data []byte) error {
    result, _ := parseData(data) // 忽略错误
    saveResult(result)
    return nil
}
```

---

## 🎯 提交规范

使用[Conventional Commits](https://www.conventionalcommits.org/)规范：

### 格式

```text
<type>(<scope>): <subject>

<body>

<footer>
```

### Type类型

- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `perf`: 性能优化

### 示例

```text
feat(concurrency): add rate limiter pattern

Implement a token bucket rate limiter with the following features:
- Configurable rate and capacity
- Thread-safe operations
- Context support for cancellation

Closes #123
```

---

## 🧪 测试要求

### 单元测试

所有新代码必须包含测试：

```go
func TestRateLimiter(t *testing.T) {
    rl := NewRateLimiter(10, 20)
    
    // 测试基本功能
    if !rl.Allow() {
        t.Error("First request should be allowed")
    }
    
    // 测试边界条件
    for i := 0; i < 20; i++ {
        rl.Allow()
    }
    
    if rl.Allow() {
        t.Error("Should reject when bucket is empty")
    }
}
```

### 基准测试

性能关键代码需要基准测试：

```go
func BenchmarkRateLimiter(b *testing.B) {
    rl := NewRateLimiter(10000, 10000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rl.Allow()
    }
}
```

### 覆盖率要求

- 新代码覆盖率 > 80%
- 核心包覆盖率 > 70%
- 整体项目覆盖率 > 60%

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包
go test ./pkg/concurrency/...

# 生成覆盖率报告
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=. -benchmem ./...
```

---

## 🔍 代码审查

### 审查清单

在提交PR前，请确保：

- [ ] 代码遵循Go规范
- [ ] 所有测试通过
- [ ] 新代码有测试覆盖
- [ ] 文档已更新
- [ ] 提交消息规范
- [ ] 无linter警告
- [ ] 性能无退化

### 审查过程

1. **自动检查**: CI会自动运行测试和linter
2. **人工审查**: 维护者会审查代码
3. **反馈修改**: 根据反馈进行修改
4. **合并**: 审查通过后合并

---

## 🎓 学习资源

### Go语言

- [Go官方文档](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 并发编程

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)

### 测试

- [Testing in Go](https://golang.org/pkg/testing/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

---

## 📞 获取帮助

遇到问题？可以：

- 📖 查看[文档](docs/)
- 💬 在[Discussions](https://github.com/yourusername/golang/discussions)提问
- 🐛 提交[Issue](https://github.com/yourusername/golang/issues)
- 📧 联系维护者

---

## 🏆 贡献者

感谢所有贡献者！

[贡献者列表](https://github.com/yourusername/golang/graphs/contributors)

---

**感谢您的贡献！** 🎉
