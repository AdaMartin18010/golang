# 🚀 Quick Start Guide

> **Welcome to the Go 1.23+ Learning Project!**  
> This guide will help you get started with Go 1.23+ new features in 5 minutes.

**Languages**: [中文](QUICK_START.md) | **English**

---

## 📋 Prerequisites

### Required

- ✅ **Go 1.23+** installed (1.24 or 1.25 recommended)
- ✅ Basic Go knowledge
- ✅ Text editor or IDE (VS Code, GoLand recommended)

### Verify Go Version

```bash
go version
# Should show: go version go1.23.0 or higher
```

If your version is lower, visit [Go official website](https://go.dev/dl/) to download the latest version.

---

## 🎯 Choose Your Learning Path

### 🌱 Beginner Path (New to Go 1.23+)

**Recommended Reading Order**:

1. **Start**: [README_EN.md](./README_EN.md) - Project overview
2. **Basics**: Review Go fundamentals if needed
3. **New Features**: Choose an interesting module:
   - [Concurrency Patterns](./examples/concurrency/)
   - [AI-Agent Architecture](./examples/advanced/ai-agent/)
   - [Go 1.25 Features](./examples/go125/)

**Estimated Time**: 2-4 hours

---

### 🚀 Intermediate Path (Familiar with Go, Want to Master 1.25)

**Recommended Reading Order**:

1. **Overview**: [Examples Showcase](./EXAMPLES_EN.md) - Quick overview of all features
2. **Deep Dive**: Read technical documentation you're interested in:
   - Performance boost? → Greentea GC
   - Container deployment? → Container-aware scheduling
   - Concurrency? → WaitGroup.Go()
   - HTTP/3? → HTTP/3 and QUIC support
3. **Practice**: Run code examples

**Estimated Time**: 1-2 hours

---

### 🎓 Expert Path (Ready for Production)

**Recommended Reading Order**:

1. **Comprehensive**: [Documentation Index](./docs/INDEX.md) - Complete knowledge system
2. **Performance**: Review benchmarks in each module
3. **Industry**: Best practices:
   - Microservices architecture
   - Cloud-native development
   - Testing strategies
4. **Migration**: Plan your upgrade path
5. **Validation**: Test examples in your environment

**Estimated Time**: 3-6 hours

---

## 💻 Hands-On Practice

### Method 1: Clone Repository

```bash
# Clone project
git clone https://github.com/AdaMartin18010/golang.git
cd golang

# Browse documentation
cd docs/02-Go语言现代化/
ls
```

---

### Method 2: Run Example Code

#### Example 1: WaitGroup.Go() - Go 1.23+ Feature

```bash
# Navigate to example directory
cd examples/concurrency

# Run all tests
go test -v .

# Run pipeline example
cd pipeline_example
go run main.go

# Expected output:
# === RUN   TestBasicWaitGroup
# --- PASS: TestBasicWaitGroup (0.00s)
# ...
# PASS
# ok   waitgroup_go 0.123s
```

**What you'll learn**:

- ✅ How to use WaitGroup for concurrency
- ✅ Panic recovery in goroutines
- ✅ Concurrent slice processing
- ✅ Concurrency limits

---

#### Example 2: Pipeline Concurrency Pattern

```bash
# Navigate to concurrency examples
cd examples/concurrency

# Run pipeline tests
go test -v . -run Pipeline

# Expected output:
# === RUN   TestSimplePipeline
# --- PASS: TestSimplePipeline (0.00s)
# === RUN   TestMultiStagePipeline
# --- PASS: TestMultiStagePipeline (0.00s)
# ...
# PASS
```

**What you'll learn**:

- ✅ Pipeline pattern implementation
- ✅ Channel-based data flow
- ✅ Timeout handling
- ✅ Fan-out/fan-in patterns

---

#### Example 3: Worker Pool Pattern

```bash
# Still in examples/concurrency
go test -v . -run WorkerPool

# Expected output:
# === RUN   TestBasicWorkerPool
# --- PASS: TestBasicWorkerPool (0.01s)
# === RUN   TestContextAwareWorkerPool
# --- PASS: TestContextAwareWorkerPool (0.01s)
# ...
# PASS
```

**What you'll learn**:

- ✅ Worker pool implementation
- ✅ Task queue management
- ✅ Load balancing
- ✅ Graceful shutdown

---

#### Example 4: AI-Agent Architecture

```bash
# Navigate to AI-Agent directory
cd docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构

# Run all tests
go test -v ./...

# Run specific engine tests
go test -v ./core -run Decision  # DecisionEngine
go test -v ./core -run Learning  # LearningEngine
go test -v . -run Agent          # BaseAgent

# Expected output:
# === RUN   TestDecisionEngine
# --- PASS: TestDecisionEngine (0.00s)
# ...
# PASS
```

**What you'll learn**:

- ✅ Multi-agent collaboration
- ✅ Decision consensus mechanism
- ✅ Reinforcement learning
- ✅ Knowledge base management

---

### Method 3: Run All Tests

```bash
# Use test summary script (recommended)
cd /path/to/golang
powershell -ExecutionPolicy Bypass -File scripts/test_summary.ps1

# Or run all tests manually
go test -v ./...

# With race detection
go test -v -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Expected Output**:

```text
📊 Test Summary:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Concurrency     14 tests  PASS
✅ WaitGroup.Go    13 tests  PASS  
✅ AI-Agent        18 tests  PASS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📈 Total:          45 tests  100% pass
```

---

## 📚 Key Documentation

### Must-Read Documents

| Document | Description | Time |
|----------|-------------|------|
| [README_EN.md](./README_EN.md) | Project overview | 5 min |
| [EXAMPLES_EN.md](./EXAMPLES_EN.md) | 45 complete examples | 10 min |
| [Doc Index](./docs/INDEX.md) | Systematic learning path | 15 min |

### By Topic

**Concurrency**:

- [Concurrency Examples](./examples/concurrency/)
- [Pipeline Pattern](./examples/concurrency/pipeline_test.go)
- [Worker Pool Pattern](./examples/concurrency/worker_pool_test.go)

**AI-Agent**:

- [DecisionEngine](./docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/core/decision_engine.go)
- [LearningEngine](./docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/core/learning_engine.go)
- [BaseAgent](./docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/core/agent.go)

**Performance**:

- Greentea GC documentation
- Container-aware scheduling
- Arena allocator

---

## 🎯 Learning Milestones

### Week 1: Foundations ✓

- [ ] Understand Go 1.23+ new features overview
- [ ] Run WaitGroup.Go() examples
- [ ] Complete basic concurrency examples
- [ ] Read 3-5 documentation files

### Week 2: Concurrency Mastery ✓

- [ ] Master Pipeline pattern
- [ ] Implement Worker Pool
- [ ] Understand fan-out/fan-in
- [ ] Complete all concurrency tests

### Week 3: Advanced Topics ✓

- [ ] Study AI-Agent architecture
- [ ] Understand DecisionEngine
- [ ] Learn LearningEngine
- [ ] Run AI-Agent examples

### Week 4: Production Ready ✓

- [ ] Review best practices
- [ ] Understand performance optimization
- [ ] Plan migration strategy
- [ ] Complete final project

---

## 💡 Tips & Tricks

### Learning Strategies

1. **Start Simple**: Begin with WaitGroup.Go() - easiest entry point
2. **Practice Incrementally**: Run examples after reading each doc
3. **Use Tests**: Tests are excellent learning resources
4. **Take Notes**: Document your understanding and questions

### Common Pitfalls

❌ **Don't**:

- Skip basic concurrency concepts
- Ignore error handling
- Optimize prematurely
- Skip writing tests

✅ **Do**:

- Build solid foundations
- Handle errors properly
- Write correct code first, optimize later
- Write tests first (TDD)

### IDE Setup

**VS Code**:

```bash
# Install Go extension
code --install-extension golang.go

# Configure settings.json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"]
}
```

**GoLand**:

- Enable Go Modules (Preferences → Go → Go Modules)
- Configure test runner (Run → Edit Configurations)

---

## 🔧 Troubleshooting

### Issue: "package not found"

```bash
# Solution: Download dependencies
go mod download
go mod tidy
```

### Issue: "CGO required"

**For ASan example**: Use the mock version

```bash
cd examples/go125/toolchain/asan_memory_leak
go run main_mock.go  # Use mock version instead
```

### Issue: Tests failing

```bash
# Check Go version
go version  # Should be 1.23+

# Clean cache
go clean -cache -testcache -modcache

# Re-run tests
go test -v ./...
```

### Issue: Import errors

```bash
# Verify module path
cat go.mod  # Check module name

# Update imports
go mod tidy
```

---

## 📊 Progress Tracking

### Checklist

**Getting Started** ✓

- [ ] Installed Go 1.23+
- [ ] Cloned repository
- [ ] Verified environment
- [ ] Read README

**Basic Examples** ✓

- [ ] WaitGroup.Go() example
- [ ] Pipeline pattern
- [ ] Worker Pool pattern
- [ ] All tests passing

**Advanced Features** ✓

- [ ] AI-Agent basics
- [ ] DecisionEngine
- [ ] LearningEngine
- [ ] Complete examples

**Production Ready** ✓

- [ ] Best practices understood
- [ ] Performance optimization
- [ ] Testing strategies
- [ ] Ready to deploy

---

## 🎓 Next Steps

### Completed Quick Start?

1. **Explore More**: [Complete Examples](./EXAMPLES_EN.md) - 45 examples
2. **Deep Dive**: [Documentation Index](./docs/INDEX.md) - Systematic learning
3. **Contribute**: [Contributing Guide](./CONTRIBUTING_EN.md) - Join the project
4. **Stay Updated**: Star and watch the repository

### Join Community

- 🐛 [Report Issues](https://github.com/AdaMartin18010/golang/issues)
- 💡 [Request Features](https://github.com/AdaMartin18010/golang/issues)
- 💬 [Discussions](https://github.com/AdaMartin18010/golang/discussions)
- 🤝 [Contribute](CONTRIBUTING_EN.md)

---

## 📞 Need Help?

### Resources

- **Documentation**: [docs/INDEX.md](docs/INDEX.md)
- **Examples**: [EXAMPLES_EN.md](EXAMPLES_EN.md)
- **FAQ**: [FAQ.md](FAQ.md)
- **Issues**: [GitHub Issues](https://github.com/AdaMartin18010/golang/issues)

### Quick Commands

```bash
# View examples
cat EXAMPLES_EN.md

# Run specific test
go test -v ./path/to/package -run TestName

# Check coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

---

<div align="center">

## 🎉 You're Ready to Go

**Choose your path and start exploring!**

[📚 Examples](EXAMPLES_EN.md) • [📖 Documentation](docs/INDEX.md) • [🤝 Contributing](CONTRIBUTING_EN.md)

---

**Last Updated**: October 19, 2025  
**Languages**: [中文](QUICK_START.md) | **English**

Happy Learning! 🚀

</div>
