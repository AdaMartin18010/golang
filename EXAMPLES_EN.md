# 📚 Go Examples Showcase

> **Complete Go 1.23-1.25 Features and Architecture Patterns Collection**  
> **100+ Test Cases | Production-Ready Code | Best Practices**

---

## 🎯 Quick Navigation

<table>
<tr>
<td width="50%">

### 🔥 Most Popular

- [AI-Agent Architecture](#ai-agent-architecture) ⭐⭐⭐⭐⭐
- [Go 1.25 Runtime Optimizations](#go-125-runtime-optimizations) ⭐⭐⭐⭐⭐  
- [Concurrency Patterns](#concurrency-patterns) ⭐⭐⭐⭐
- [Modern Features Collection](#modern-features) ⭐⭐⭐⭐

</td>
<td width="50%">

### 📖 Browse by Topic

- [Advanced Features](#advanced-features)
- [Performance Optimization](#performance-optimization)
- [Cloud Native](#cloud-native)
- [Testing Framework](#testing-framework)

</td>
</tr>
</table>

---

## ⭐ Featured Examples

### AI-Agent Architecture

> **Complete Production-Grade AI Agent Implementation** | 18 Test Cases | 100% Coverage

**Location**: [`examples/advanced/ai-agent/`](./examples/advanced/ai-agent/)

```bash
cd examples/advanced/ai-agent
go test -v ./...
```

**Core Features**:

- ✅ Decision Engine - Intelligent decision processing
- ✅ Learning Engine - Experience learning & optimization
- ✅ Base Agent - Complete lifecycle management
- ✅ Multi-Agent Coordination - Distributed collaboration
- ✅ Multimodal Interface - Unified interaction abstraction

**Test Coverage**:

- Decision Engine: 7 test cases
- Learning Engine: 9 test cases  
- Base Agent: 2 test cases

---

### Go 1.25 Runtime Optimizations

> **Latest Go Runtime Features** | Greentea GC | Container-Aware

**Location**: [`examples/go125/runtime/`](./examples/go125/runtime/)

```bash
# Greentea GC Optimization
cd examples/go125/runtime/gc_optimization
go test -v

# Container-Aware Scheduling
cd examples/go125/runtime/container_scheduling
go run main.go

# Memory Allocator Enhancements
cd examples/go125/runtime/memory_allocator
go test -bench=.
```

**Core Features**:

- ✅ Greentea GC - Next-generation garbage collector
- ✅ Container Scheduling - Adaptive container environment
- ✅ Memory Optimization - High-performance allocation

---

### Concurrency Patterns

> **Classic Concurrency Patterns & Best Practices** | 15+ Test Cases

**Location**: [`examples/concurrency/`](./examples/concurrency/)

```bash
cd examples/concurrency
go test -v
```

**Core Patterns**:

- ✅ Pipeline Pattern - Stream processing
- ✅ Worker Pool Pattern - Concurrent worker pool
- ✅ Fan-Out/Fan-In - Distribution and aggregation
- ✅ Context Pattern - Context propagation & cancellation

**Example Programs**:

```bash
# Run Pipeline example
cd examples/concurrency/pipeline_example
go run main.go

# Run Worker Pool example
cd examples/concurrency/worker_pool_example
go run main.go
```

---

### Modern Features

> **Go 1.23-1.25 Modern Features Collection** | 95 Code Files

**Location**: [`examples/modern-features/`](./examples/modern-features/)

**Complete Directory**:

#### 1. New Features Deep Dive

📁 `01-new-features/`

- Generic Type Aliases - Advanced generics
- Swiss Table Optimization - High-performance Map
- Testing Enhancements - Loop variable semantics
- WASM & WASI - WebAssembly support

#### 2. Concurrency 2.0

📁 `02-concurrency-2.0/`

- Advanced Worker Pool patterns
- Concurrency safety best practices

#### 3. Standard Library Enhancements

📁 `03-stdlib-enhancements/`

- `slog` - Structured logging
- `ServeMux` - New HTTP router
- Concurrency primitives enhancements

#### 4. Performance & Toolchain

📁 `05-performance-toolchain/`

- PGO - Profile-Guided Optimization
- CGO & Interoperability
- Compiler & linker optimizations

#### 5. Modern Architecture Patterns

📁 `06-architecture-patterns/`

- Clean Architecture
- Hexagonal Architecture

#### 6. Performance Optimization 2.0

📁 `07-performance-optimization/`

- Zero-Copy techniques
- SIMD optimization

#### 7. Cloud Native 2.0

📁 `09-cloud-native-2.0/`

- Kubernetes Operator
- Service Mesh integration
- GitOps pipeline

```bash
# View complete documentation
cat examples/modern-features/README.md

# Run all tests
cd examples/modern-features
go test ./... -v
```

---

## 📋 Complete Examples Catalog

### 🔥 Advanced Features

| Example | Description | Location | Tests |
|---------|-------------|----------|-------|
| AI-Agent | Intelligent agent architecture | `advanced/ai-agent/` | ✅ 18 |
| HTTP/3 Server | HTTP/3 and QUIC | `advanced/http3-server/` | ✅ |
| Weak Pointer Cache | Memory optimization | `advanced/cache-weak-pointer/` | ✅ |
| Arena Allocator | Custom allocation | `advanced/arena-allocator/` | ✅ |
| Worker Pool | Concurrent worker pool | `advanced/worker-pool/` | ✅ |

### 🚀 Go 1.25 Features

| Example | Description | Location | Tests |
|---------|-------------|----------|-------|
| Greentea GC | New garbage collector | `go125/runtime/gc_optimization/` | ✅ |
| Container Scheduling | Container awareness | `go125/runtime/container_scheduling/` | ✅ |
| Memory Optimization | Allocator enhancements | `go125/runtime/memory_allocator/` | ✅ |
| ASan Detection | Memory safety | `go125/toolchain/asan_memory_leak/` | ✅ |

### 🆕 Modern Features (95 Files)

| Category | Examples | Location |
|----------|----------|----------|
| New Features Deep Dive | 6 topics | `modern-features/01-new-features/` |
| Concurrency 2.0 | 1 topic | `modern-features/02-concurrency-2.0/` |
| Stdlib Enhancements | 3 topics | `modern-features/03-stdlib-enhancements/` |
| Performance & Toolchain | 3 topics | `modern-features/05-performance-toolchain/` |
| Architecture Patterns | 2 topics | `modern-features/06-architecture-patterns/` |
| Performance Optimization | 2 topics | `modern-features/07-performance-optimization/` |
| Cloud Native 2.0 | 3 topics | `modern-features/09-cloud-native-2.0/` |

### 🔄 Concurrency Programming

| Example | Description | Tests |
|---------|-------------|-------|
| Pipeline | Pipeline pattern | ✅ 6 |
| Worker Pool | Worker pool pattern | ✅ 7 |
| Fan-Out/Fan-In | Distribution pattern | ✅ |
| Context | Context management | ✅ |

### 🧪 Testing Framework

| Example | Description | Location |
|---------|-------------|----------|
| Integration Test Framework | Complete test system | `testing-framework/` |
| Performance Regression | Automated monitoring | `testing-framework/` |

### 📊 Observability

| Example | Description | Location |
|---------|-------------|----------|
| OpenTelemetry | Complete observability | `observability/` |
| Prometheus | Metrics monitoring | `observability/prometheus/` |
| Grafana | Visualization dashboard | `observability/grafana/` |

---

## 🚀 Quick Start

### 1️⃣ Run All Tests

```bash
# From project root
go test ./examples/... -v
```

### 2️⃣ Recommended Learning Path

#### Beginner (1-2 hours)

```bash
# Step 1: Concurrency Basics
cd examples/concurrency
go test -v

# Step 2: View example programs
cd pipeline_example
go run main.go
```

#### Intermediate (3-5 hours)

```bash
# Step 1: AI-Agent Architecture
cd examples/advanced/ai-agent
go test -v ./...

# Step 2: Modern Features
cd examples/modern-features/03-stdlib-enhancements
go test -v ./...
```

#### Advanced (5-10 hours)

```bash
# Step 1: Go 1.25 Runtime
cd examples/go125/runtime
go test -v ./...

# Step 2: Performance Optimization
cd examples/modern-features/07-performance-optimization
go test -bench=. -benchmem
```

### 3️⃣ Use by Scenario

#### Scenario 1: Learning Concurrency

```bash
# 1. Read documentation
cat docs/03-并发编程/README.md

# 2. Run examples
cd examples/concurrency
go test -v

# 3. View real-world cases
cd examples/advanced/ai-agent/coordination
go test -v
```

#### Scenario 2: Learning Go 1.25 Features

```bash
# 1. Read documentation
cat docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md

# 2. Run runtime examples
cd examples/go125/runtime/gc_optimization
go test -v

# 3. Test toolchain enhancements
cd examples/go125/toolchain/asan_memory_leak
go test -v
```

#### Scenario 3: Performance Optimization

```bash
# 1. View PGO examples
cd examples/modern-features/05-performance-toolchain/01-Profile-Guided-Optimization-PGO
go test -bench=.

# 2. Learn Zero-Copy
cd examples/modern-features/07-performance-optimization/01-zero-copy
go test -bench=.

# 3. SIMD Optimization
cd examples/modern-features/07-performance-optimization/02-simd-optimization
go test -bench=.
```

---

## 📊 Examples Statistics

### Overall Stats

| Metric | Count |
|--------|-------|
| Total Examples | 50+ |
| Code Files | 150+ |
| Test Cases | 100+ |
| Topics Covered | 10+ |

### By Category

```text
🔥 Advanced Features    : 5 examples
🚀 Go 1.25 Features     : 4 examples
🆕 Modern Features      : 95 files, 20+ subtopics
🔄 Concurrency          : 4 core patterns
🧪 Testing Framework    : 1 complete system
📊 Observability        : 1 integrated solution
```

### Test Coverage

```text
✅ AI-Agent          : 18 test cases
✅ Concurrency       : 15+ test cases
✅ Runtime           : 10+ test cases
✅ Toolchain         : 5+ test cases
```

---

## 📖 Related Documentation

| Document | Description | Location |
|----------|-------------|----------|
| **Complete Examples Index** | Detailed examples list | [examples/README.md](./examples/README.md) |
| **Modern Features Guide** | 95 files explained | [examples/modern-features/README.md](./examples/modern-features/README.md) |
| **Concurrency Documentation** | Theory + Practice | [docs/03-并发编程/](./docs/03-并发编程/) |
| **Go 1.25 Documentation** | New features explained | [docs/02-Go语言现代化/](./docs/02-Go语言现代化/) |
| **Quick Reference** | One-page cheatsheet | [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) |

---

## 🔧 Development Tools

### Verify Project Structure

```bash
# Windows
powershell -ExecutionPolicy Bypass -File scripts/verify_structure.ps1

# Linux/macOS
bash scripts/verify_structure.sh
```

### Code Quality Check

```bash
# Quality scan
powershell -ExecutionPolicy Bypass -File scripts/scan_code_quality.ps1

# Test statistics
powershell -ExecutionPolicy Bypass -File scripts/test_summary.ps1
```

### Project Statistics

```bash
cd scripts/project_stats
go run main.go
```

---

## ❓ FAQ

**Q: How to run all examples?**

```bash
go test ./examples/... -v
```

**Q: Can examples be used in production?**
A: Yes, AI-Agent, concurrency patterns, and performance optimization examples are production-grade code.

**Q: Can't find an old example?**
A: All code has been reorganized to `examples/` directory. See [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)

**Q: How to contribute new examples?**
A: See [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## 🎯 Learning Recommendations

### By Difficulty

1. **Beginner** ⭐: Concurrency basic patterns
2. **Intermediate** ⭐⭐: Stdlib enhancements, ServeMux
3. **Advanced** ⭐⭐⭐: AI-Agent, performance optimization
4. **Expert** ⭐⭐⭐⭐: Go 1.25 runtime, SIMD

### By Time

- **1 hour**: Quick overview → Run `concurrency` examples
- **Half day**: Deep dive → AI-Agent complete tests
- **1-2 days**: Comprehensive → All modern-features
- **1 week**: Expert level → Go 1.25 + performance

---

## 📞 Get Help

- 📖 See [FAQ.md](./FAQ.md)
- 📝 Read [CONTRIBUTING.md](./CONTRIBUTING.md)  
- 🐛 Submit [GitHub Issues](../../issues)
- 💬 Join [Discussions](../../discussions)

---

<div align="center">

**Examples Version**: 2.0.0  
**Last Updated**: October 19, 2025  
**Go Version Required**: 1.23+ (Recommended 1.25+)

---

[Complete Documentation](./docs/) | [Quick Start](./QUICK_START.md) | [Contributing](./CONTRIBUTING.md) | [Back to Home](./README.md)

</div>
