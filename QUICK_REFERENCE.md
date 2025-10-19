# 📋 快速参考卡片

> **新项目结构快速查询** | 一页搞定所有查找

---

## 🎯 我想找

### 📚 学习某个主题

```bash
cd docs/
ls                    # 查看所有主题
cat INDEX.md          # 查看文档索引
```

**主要主题**:

- `01-Go语言基础/` → 基础语法
- `02-Go语言现代化/` → Go 1.23-1.25新特性  
- `03-并发编程/` → 并发编程深度解析
- `04-设计模式/` → 设计模式实践
- `05-性能优化/` → 性能优化指南
- `06-微服务架构/` → 微服务架构
- `07-云原生与部署/` → 云原生实践

### 💻 运行代码示例

```bash
cd examples/
ls                    # 查看所有示例
cat README.md         # 查看示例索引
```

**热门示例**:

- `advanced/ai-agent/` → AI-Agent完整实现 ⭐
- `go125/` → Go 1.25特性 ⭐
- `modern-features/` → 现代化特性集合 ⭐
- `concurrency/` → 并发模式
- `testing-framework/` → 测试框架

### 📊 查看项目报告

```bash
cd reports/
ls                    # 查看报告分类
cat README.md         # 查看报告索引
```

**报告分类**:

- `phase-reports/` → 阶段报告
- `archive/` → 历史归档
- `code-quality/` → 代码质量
- `daily-summaries/` → 日常总结

### 🗄️ 查看历史内容

```bash
cd archive/
ls                    # 查看归档内容
```

**归档内容**:

- `model/` → 900+旧文档
- `old-docs/` → 杂散文档

---

## 🔧 常用命令

### 验证项目结构

```bash
# Windows
powershell -ExecutionPolicy Bypass -File scripts/verify_structure.ps1

# Linux/macOS  
bash scripts/verify_structure.sh
```

### 代码质量检查

```bash
# Windows
powershell -ExecutionPolicy Bypass -File scripts/scan_code_quality.ps1

# Linux/macOS
bash scripts/scan_code_quality.sh
```

### 运行测试

```bash
powershell -ExecutionPolicy Bypass -File scripts/test_summary.ps1
```

### 项目统计

```bash
cd scripts/project_stats
go run main.go
```

---

## 📖 重要文档

| 文档 | 用途 | 位置 |
|------|------|------|
| **README.md** | 项目主页 | 根目录 |
| **RESTRUCTURE.md** | 重组说明 | 根目录 |
| **MIGRATION_GUIDE.md** | 迁移指南 | 根目录 |
| **QUICK_REFERENCE.md** | 本文档 | 根目录 |
| **docs/INDEX.md** | 文档索引 | docs/ |
| **examples/README.md** | 示例索引 | examples/ |
| **reports/README.md** | 报告索引 | reports/ |

---

## 🎯 按场景查找

### 场景1: 学习AI-Agent

```bash
# 1. 阅读文档
cat docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/README.md

# 2. 运行代码
cd examples/advanced/ai-agent/
go test -v ./...
```

### 场景2: 学习Go 1.25新特性

```bash
# 1. 阅读文档
cat docs/02-Go语言现代化/12-Go-1.25运行时优化/README.md

# 2. 运行示例
cd examples/go125/runtime/gc_optimization/
go test -v
```

### 场景3: 学习并发编程

```bash
# 1. 阅读文档
cat docs/03-并发编程/README.md

# 2. 运行示例
cd examples/concurrency/
go test -v
```

### 场景4: 贡献代码

```bash
# 1. 阅读贡献指南
cat CONTRIBUTING.md

# 2. 验证环境
scripts/verify_structure.ps1

# 3. 代码检查
scripts/scan_code_quality.ps1

# 4. 运行测试
scripts/test_summary.ps1
```

---

## 🗂️ 目录职责速查

| 目录 | 职责 | 规则 |
|------|------|------|
| `docs/` | 📚 学习文档 | ❌ 不放代码 |
| `examples/` | 💻 可运行代码 | ✅ 独立可运行 |
| `reports/` | 📊 项目报告 | ✅ 按类型分类 |
| `archive/` | 🗄️ 历史归档 | ⚠️ 只读 |
| `scripts/` | 🔧 开发工具 | ✅ 可执行脚本 |

---

## 🚀 快速开始

### 新手入门 (5分钟)

```bash
# 1. 查看项目结构
cat README.md

# 2. 运行第一个示例
cd examples/concurrency
go test -v

# 3. 查看文档索引
cat docs/INDEX.md
```

### 开发者工作流

```bash
# 每日流程
1. scripts/verify_structure.ps1    # 验证结构
2. scripts/scan_code_quality.ps1   # 质量检查
3. scripts/test_summary.ps1        # 运行测试
```

---

## ❓ 常见问题

**Q: 找不到以前的代码？**
A: 所有代码已移至 `examples/`，查看 `MIGRATION_GUIDE.md`

**Q: 文档中的代码示例呢？**
A: 完整代码在 `examples/`，文档中有链接指向

**Q: 如何验证项目结构？**
A: 运行 `scripts/verify_structure.ps1`

**Q: 如何查找特定功能？**
A: 查看对应的 INDEX 或 README 文件

---

## 📞 获取帮助

| 需求 | 资源 |
|------|------|
| 详细说明 | [RESTRUCTURE.md](RESTRUCTURE.md) |
| 迁移指南 | [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) |
| 常见问题 | [FAQ.md](FAQ.md) |
| 贡献指南 | [CONTRIBUTING.md](CONTRIBUTING.md) |
| 提问题 | [GitHub Issues](../../issues) |

---

## 🎯 记住这些

### ✅ 现在这样做

```bash
✅ 在 docs/ 找文档
✅ 在 examples/ 运行代码
✅ 在 reports/ 看报告
✅ 用 scripts/ 验证
```

### ❌ 不要这样做

```bash
❌ 在 docs/ 找代码（已移走）
❌ 在根目录找报告（已移走）
❌ 在 model/ 找文档（已归档）
```

---

**快速参考版本**: 1.0  
**最后更新**: 2025年10月19日  
**基于项目版本**: 2.0.0+

---

<div align="center">

💡 **提示**: 收藏本页，随时快速查找！

[详细文档](RESTRUCTURE.md) | [返回主页](README.md)

</div>
