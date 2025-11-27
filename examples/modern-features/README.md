# Go 现代化特性示例集

> **Go 1.21-1.25 现代化特性完整示例**

本目录包含从 `docs/` 迁移出来的所有可运行代码示例，展示 Go 语言的现代化特性和最佳实践。

---

## 📋 目录结构

### 01-new-features/ - 新特性深度解析

- **泛型类型别名** - 泛型的高级用法
  - 基础示例
  - 高级模式
  - 性能分析
  - 实际应用案例
- **Swiss Table优化** - Map性能优化
- **测试增强** - Go 1.22+测试新特性
- **WASM与WASI** - WebAssembly支持
- **for循环变量语义** - Go 1.22循环变量改进

### 02-concurrency-2.0/ - 并发编程2.0

- Worker Pool 高级实现
- 并发模式最佳实践
- 性能优化技巧

### 03-stdlib-enhancements/ - 标准库增强

- **结构化日志 slog** - 现代化日志库
- **HTTP路由 ServeMux** - 增强的路由功能
- **并发原语与模式** - 新的并发工具

### 05-performance-toolchain/ - 性能与工具链

- **PGO (Profile-Guided Optimization)** - 配置引导优化
- **CGO与互操作性** - C语言互操作
  - 基础用法
  - 内存管理
  - 性能分析
- **链接器与编译器优化** - 构建优化

### 06-architecture-patterns/ - 架构模式现代化

- **Clean Architecture** - 完整的Clean架构实现
  - Domain层
  - Repository层
  - Use Case层
  - Delivery层

### 07-performance-optimization/ - 性能优化2.0

- **Zero-Copy优化**
  - Sendfile实现
  - 网络缓冲区优化
  - 内存池管理
- **SIMD优化**
  - 向量运算
  - 矩阵计算
  - 图像处理
  - 加密优化

### 08-cloud-native/ - 云原生集成

- **Kubernetes Operator** - K8s操作器实现
  - 事件记录
  - 指标收集
  - 资源管理

### 09-cloud-native-2.0/ - 云原生2.0

- **Kubernetes Operator** - 高级操作器
- **Service Mesh集成** - 服务网格
- **GitOps流水线** - 自动化部署

---

## 🚀 快速开始

### 环境要求

```text
Go版本: 1.21+ (部分示例需要 1.22-1.25)
操作系统: Windows/Linux/macOS
```

### 运行示例

```bash
# 泛型示例
cd examples/modern-features/01-new-features/01-generic-type-alias
go test -v ./...

# PGO优化
cd examples/modern-features/05-performance-toolchain/01-Profile-Guided-Optimization-PGO
go run main.go

# Clean Architecture
cd examples/modern-features/06-architecture-patterns/01-Clean-Architecture
go run cmd/main.go

# Zero-Copy优化
cd examples/modern-features/07-performance-optimization/01-zero-copy/sendfile
go run server.go
```

---

## 📊 特性覆盖

| 分类 | Go版本 | 示例数 | 状态 |
|------|--------|--------|------|
| 泛型特性 | 1.21+ | 5 | ✅ |
| 标准库增强 | 1.21-1.22 | 3 | ✅ |
| 性能优化 | 1.21+ | 8 | ✅ |
| 架构模式 | 通用 | 2 | ✅ |
| 云原生 | 通用 | 3 | ✅ |

---

## 📖 学习路径

### 初级 (1-2小时)

1. 结构化日志 slog
2. HTTP路由 ServeMux
3. For循环变量语义

### 中级 (3-5小时)

1. 泛型类型别名
2. PGO性能优化
3. WASM与WASI
4. Worker Pool实现

### 高级 (1-2天)

1. Clean Architecture完整实现
2. Zero-Copy性能优化
3. SIMD向量计算
4. Kubernetes Operator

---

## 💡 最佳实践

### 代码组织

- 每个示例独立可运行
- 包含完整的测试用例
- 配备详细的注释说明

### 性能考虑

- PGO优化示例
- 基准测试完整
- 内存使用分析

### 测试标准

- 单元测试覆盖
- 集成测试完善
- 性能基准测试

---

## 🔗 相关文档

- [主文档](../../docs/02-Go语言现代化/README.md) - 详细理论文档
- [示例索引](../README.md) - 所有示例总览
- [项目README](../../README.md) - 项目主页

---

## 📝 注意事项

1. **版本要求**: 部分示例需要特定Go版本
2. **依赖管理**: 每个子目录可能有独立的 `go.mod`
3. **环境配置**: WASM示例需要特殊的构建标签
4. **性能测试**: SIMD示例需要特定CPU特性支持

---

**迁移日期**: 2025年10月19日
**来源**: docs/02-Go语言现代化/
**状态**: ✅ 完全可运行
