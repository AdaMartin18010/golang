# 📖 项目重组迁移指南

> **适用对象**: 项目贡献者、使用者  
> **更新日期**: 2025年10月19日

---

## 🎯 概述

如果您之前使用过本项目，请阅读本指南了解最新的项目结构变化。

---

## 🔄 主要变更

### 1. 文档位置

#### ❌ 旧路径
```
docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/
├── core/agent.go
├── go.mod
└── README.md
```

#### ✅ 新路径
```
docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/
└── README.md  (指向 examples/advanced/ai-agent/)

examples/advanced/ai-agent/
├── core/agent.go
├── go.mod
└── README.md
```

### 2. 示例代码位置

#### 所有代码现在在 `examples/`

| 旧位置 | 新位置 | 说明 |
|-------|--------|------|
| `docs/02-Go语言现代化/01-新特性深度解析/` | `examples/modern-features/01-new-features/` | 新特性示例 |
| `docs/02-Go语言现代化/05-性能与工具链/` | `examples/modern-features/05-performance-toolchain/` | 性能工具 |
| `docs/02-Go语言现代化/06-架构模式现代化/` | `examples/modern-features/06-architecture-patterns/` | 架构模式 |
| `docs/02-Go语言现代化/07-性能优化2.0/` | `examples/modern-features/07-performance-optimization/` | 性能优化 |
| `docs/02-Go语言现代化/12-Go-1.25运行时优化/examples/` | `examples/go125/runtime/` | Go 1.25运行时 |
| `docs/02-Go语言现代化/13-Go-1.25工具链增强/examples/` | `examples/go125/toolchain/` | Go 1.25工具链 |

### 3. 报告位置

#### ❌ 旧路径
```
根目录/
├── Phase-1-阶段性总结-2025-10-19.md
├── 📢持续推进完成简报-2025-10-19.md
└── ...
```

#### ✅ 新路径
```
reports/phase-reports/
├── Phase-1-阶段性总结-2025-10-19.md
├── 📢持续推进完成简报-2025-10-19.md
└── ...
```

### 4. 历史内容归档

#### 已归档内容

```
archive/
├── model/          # 900+个旧文档
└── old-docs/       # 旧的计划文档
```

这些内容已归档，不影响当前使用。

---

## 🚀 快速查找指南

### 我想找...

#### 📚 学习文档
```bash
# 所有文档都在 docs/ 目录
cd docs/

# 查看文档索引
cat INDEX.md

# 具体主题
cd 03-并发编程/          # 并发编程深度解析
cd 02-Go语言现代化/      # Go 1.23-1.25新特性
cd 05-性能优化/          # 性能优化指南
```

#### 💻 可运行代码
```bash
# 所有代码都在 examples/ 目录
cd examples/

# 查看示例索引
cat README.md

# 热门示例
cd advanced/ai-agent/           # AI-Agent完整实现
cd go125/runtime/               # Go 1.25运行时优化
cd modern-features/             # 现代化特性集合
cd concurrency/                 # 并发模式示例
```

#### 📊 项目报告
```bash
# 所有报告都在 reports/ 目录
cd reports/

# 查看报告索引
cat README.md

# 阶段报告
cd phase-reports/

# 历史报告
cd archive/
```

#### 🗄️ 历史内容
```bash
# 历史内容在 archive/ 目录
cd archive/

# 旧 model 目录
cd model/

# 旧文档
cd old-docs/
```

---

## 📖 使用示例

### 场景1: 学习AI-Agent架构

**步骤**:

1. **阅读理论文档**
   ```bash
   # 查看架构设计文档
   cat docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/README.md
   ```

2. **运行示例代码**
   ```bash
   # 进入代码目录
   cd examples/advanced/ai-agent/
   
   # 运行测试
   go test -v ./...
   
   # 查看实现
   cat core/agent.go
   ```

### 场景2: 学习Go 1.25新特性

**步骤**:

1. **阅读文档**
   ```bash
   cat docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md
   ```

2. **运行示例**
   ```bash
   cd examples/go125/runtime/gc_optimization/
   go test -v
   ```

### 场景3: 学习并发编程

**步骤**:

1. **阅读理论**
   ```bash
   cat docs/03-并发编程/README.md
   ```

2. **实践代码**
   ```bash
   cd examples/concurrency/
   go test -v
   ```

---

## 🔧 开发者工作流

### 贡献新内容

#### 添加文档
```bash
# 1. 在 docs/ 相应目录创建 .md 文件
docs/03-并发编程/14-新主题.md

# 2. 更新 docs/INDEX.md 索引
```

#### 添加代码示例
```bash
# 1. 在 examples/ 相应目录创建代码
examples/concurrency/new_pattern_test.go

# 2. 确保代码可独立运行
go test -v ./examples/...

# 3. 更新 examples/README.md
```

#### 添加报告
```bash
# 放入 reports/phase-reports/
reports/phase-reports/Phase-4-*.md
```

### 维护规范

| 目录 | 职责 | 规则 |
|------|------|------|
| `docs/` | 📚 纯文档 | ❌ 不放代码文件 |
| `examples/` | 💻 可运行代码 | ✅ 每个示例独立可运行 |
| `reports/` | 📊 项目报告 | ✅ 按类型分类 |
| `archive/` | 🗄️ 历史归档 | ⚠️ 只读，不再更新 |
| 根目录 | 📄 核心文档 | ❌ 不放临时文件 |

---

## ❓ 常见问题

### Q: 找不到以前的代码文件？

**A**: 所有代码已移至 `examples/` 目录，请查看：

1. `examples/advanced/` - 高级特性
2. `examples/go125/` - Go 1.25特性
3. `examples/modern-features/` - 现代化特性集合

### Q: 文档中代码示例去哪了？

**A**: 
- 小的代码片段仍在文档中作为示例
- 完整的可运行项目已移至 `examples/`
- 文档中有指向代码的链接

### Q: 旧的 model/ 目录去哪了？

**A**: 已归档至 `archive/model/`，不影响当前工作。

### Q: 如何运行示例？

**A**: 
```bash
# 查看示例索引
cat examples/README.md

# 进入具体示例目录
cd examples/advanced/ai-agent/

# 运行测试
go test -v ./...
```

### Q: 如何查找特定主题？

**A**: 
```bash
# 文档
cat docs/INDEX.md

# 示例
cat examples/README.md

# 报告
cat reports/README.md
```

---

## 📚 相关文档

- [RESTRUCTURE.md](RESTRUCTURE.md) - 完整重组说明
- [README.md](README.md) - 项目主页
- [examples/README.md](examples/README.md) - 示例索引
- [docs/INDEX.md](docs/INDEX.md) - 文档索引

---

## 🎯 重要提示

### ✅ 现在这样做

```bash
# ✅ 在 docs/ 中查找文档
cd docs/03-并发编程/

# ✅ 在 examples/ 中运行代码
cd examples/concurrency/
go test -v

# ✅ 在 reports/ 中查看报告
cd reports/phase-reports/
```

### ❌ 不要这样做

```bash
# ❌ 不要在 docs/ 中找代码
cd docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/
go test  # 没有代码文件了

# ❌ 不要在根目录找报告
ls Phase-*.md  # 已移至 reports/phase-reports/

# ❌ 不要在 model/ 找文档
cd model/  # 已移至 archive/model/
```

---

## 🔜 需要帮助？

如果您在迁移过程中遇到问题：

1. 查看 [RESTRUCTURE.md](RESTRUCTURE.md) 了解完整变更
2. 查看 [FAQ.md](FAQ.md) 常见问题
3. 提交 [Issue](../../issues)

---

**迁移指南版本**: 1.0  
**最后更新**: 2025年10月19日  
**适用项目版本**: 2.0.0+

