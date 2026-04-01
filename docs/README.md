# Go 1.26 Clean Architecture 项目文档

> **版本**: v2.0
> **更新日期**: 2026-04-02
> **Go版本**: 1.26.1
> **文档数**: 15篇核心文档

---

## 📚 核心文档索引

### 🏗️ 架构设计 (5篇)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [clean-architecture.md](architecture/clean-architecture.md) | 整洁架构设计 | ⭐⭐⭐⭐⭐ |
| [clean-architecture-2026-best-practices.md](architecture/clean-architecture-2026-best-practices.md) | 2026最佳实践 | ⭐⭐⭐⭐ |
| [opentelemetry.md](architecture/tech-stack/observability/opentelemetry.md) | 可观测性架构 | ⭐⭐⭐⭐ |
| [OPENTELEMETRY-2026-UPDATE.md](architecture/tech-stack/observability/OPENTELEMETRY-2026-UPDATE.md) | OTel 2026更新 | ⭐⭐⭐⭐ |
| [ent-orm.md](architecture/tech-stack/data/ent-orm.md) | 数据持久化 | ⭐⭐⭐ |

### 📖 Go 1.26 核心知识 (5篇)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [05-csp-formal-model.md](go126-comprehensive-guide/05-csp-formal-model.md) | CSP形式化模型 | ⭐⭐⭐⭐⭐ |
| [01-language-features.md](go126-comprehensive-guide/01-language-features.md) | 语言特性 | ⭐⭐⭐⭐ |
| [03-type-system.md](go126-comprehensive-guide/03-type-system.md) | 类型系统 | ⭐⭐⭐⭐ |
| [09-concurrency-patterns.md](go126-comprehensive-guide/09-concurrency-patterns.md) | 并发模式 | ⭐⭐⭐⭐ |
| [26-memory-management.md](go126-comprehensive-guide/26-memory-management.md) | 内存管理 | ⭐⭐⭐ |

### 📋 参考文档 (3篇)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [00-Go-1.26完整知识体系总览-2026.md](00-Go-1.26完整知识体系总览-2026.md) | 知识总览 | ⭐⭐⭐⭐ |
| [go126-package-management.md](go126-package-management.md) | 包管理 | ⭐⭐⭐ |
| [CORE-DOCUMENTS.md](CORE-DOCUMENTS.md) | 文档索引 | ⭐⭐⭐ |

---

## 🚀 快速开始

1. **新用户**: 从 [clean-architecture.md](architecture/clean-architecture.md) 开始
2. **Go 1.26 特性**: 阅读 [01-language-features.md](go126-comprehensive-guide/01-language-features.md)
3. **并发编程**: 深入学习 [05-csp-formal-model.md](go126-comprehensive-guide/05-csp-formal-model.md)
4. **可观测性**: 查看 [opentelemetry.md](architecture/tech-stack/observability/opentelemetry.md)

---

## 📁 文档结构

```
docs/
├── README.md                                    # 本文档
├── CORE-DOCUMENTS.md                            # 高质量文档索引
├── 00-Go-1.26完整知识体系总览-2026.md          # 知识总览
├── go126-package-management.md                  # 包管理
├── architecture/                                # 架构文档
│   ├── clean-architecture.md                   # 整洁架构
│   ├── clean-architecture-2026-best-practices.md
│   └── tech-stack/                             # 技术栈
│       ├── observability/
│       └── data/
└── go126-comprehensive-guide/                  # Go 1.26核心
    ├── 01-language-features.md
    ├── 03-type-system.md
    ├── 05-csp-formal-model.md
    ├── 09-concurrency-patterns.md
    └── 26-memory-management.md
```

---

## 📝 文档说明

- **总文档数**: 15篇（从1002篇精简）
- **文档质量**: 核心高质量文档
- **更新频率**: 按重要程度维护
- **归档文档**: 备份在 `docs-backup-20260402-062435.zip`

---

## 🔍 查找文档

使用以下关键词快速定位：

- **架构**: `architecture/`
- **Go特性**: `go126-comprehensive-guide/`
- **并发**: `05-csp-formal-model.md`, `09-concurrency-patterns.md`
- **可观测性**: `opentelemetry.md`, `OPENTELEMETRY-2026-UPDATE.md`

---

*最后更新: 2026-04-02*
*文档清理完成 ✅*
