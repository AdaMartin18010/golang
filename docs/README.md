# Go 1.26 Clean Architecture 项目文档

> **版本**: v2.1  
> **更新日期**: 2026-04-02  
> **Go版本**: 1.26.1  
> **文档数**: 18篇核心文档

---

## 📚 核心文档索引

### 🏗️ 架构设计 (6篇)

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

### 🚀 操作文档 (4篇)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [API 文档](api/README.md) | HTTP/gRPC/GraphQL | ⭐⭐⭐⭐ |
| [部署指南](deployment/README.md) | Docker/K8s 部署 | ⭐⭐⭐⭐ |
| [开发环境搭建](development/setup.md) | 本地开发 | ⭐⭐⭐⭐ |
| [项目总览](../README.md) | 项目根文档 | ⭐⭐⭐⭐⭐ |

---

## 🎯 快速导航

### 新用户
1. [开发环境搭建](development/setup.md) - 30分钟上手
2. [clean-architecture.md](architecture/clean-architecture.md) - 了解架构

### 开发者
1. [API 文档](api/README.md) - 接口参考
2. [部署指南](deployment/README.md) - 部署应用

### 架构师
1. [05-csp-formal-model.md](go126-comprehensive-guide/05-csp-formal-model.md) - 并发理论
2. [clean-architecture-2026-best-practices.md](architecture/clean-architecture-2026-best-practices.md) - 最佳实践

---

## 📁 文档结构

```
docs/
├── README.md                                    # 本文档
├── 00-Go-1.26完整知识体系总览-2026.md          # 知识总览
├── go126-package-management.md                  # 包管理
├── architecture/                                # 架构文档
│   ├── clean-architecture.md                   # 整洁架构
│   ├── clean-architecture-2026-best-practices.md
│   └── tech-stack/                             # 技术栈
│       ├── observability/
│       └── data/
├── go126-comprehensive-guide/                  # Go 1.26核心
│   ├── 01-language-features.md
│   ├── 03-type-system.md
│   ├── 05-csp-formal-model.md
│   ├── 09-concurrency-patterns.md
│   └── 26-memory-management.md
├── api/                                         # API 文档
├── deployment/                                  # 部署指南
└── development/                                 # 开发指南
```

---

## 📈 文档统计

- **总文档数**: 18篇
- **架构文档**: 6篇
- **Go 核心**: 5篇
- **操作文档**: 4篇
- **参考文档**: 3篇

---

## 🆘 获取帮助

- 查看 [开发环境搭建](development/setup.md)
- 阅读 [API 文档](api/README.md)
- 参考 [部署指南](deployment/README.md)
- 创建 [GitHub Issue](../issues)

---

## 📝 贡献

查看 [CONTRIBUTING.md](../CONTRIBUTING.md) 了解如何贡献。

---

*最后更新: 2026-04-02*  
*文档清理与优化完成 ✅*
