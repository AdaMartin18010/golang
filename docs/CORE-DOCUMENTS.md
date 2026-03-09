# 高质量核心文档索引

> **文档质量评估标准**: 满足≥3条以下标准
>
> - ✅ 完整的形式化语义/推导
> - ✅ 生产环境实战验证
> - ✅ 深度性能基准测试
> - ✅ 架构选型决策树
> - ✅ 量化收益分析

**版本**: v1.0
**更新日期**: 2026-03-09
**文档总数**: 6篇（目标<50篇）

---

## 📋 核心高质量文档列表

### 形式化语义/推导类 (2篇)

| # | 文档路径 | 质量评级 | 核心亮点 |
|---|---------|---------|---------|
| 1 | [`go126-comprehensive-guide/05-csp-formal-model.md`](./go126-comprehensive-guide/05-csp-formal-model.md) | ⭐⭐⭐⭐⭐ (5/5) | **形式化语义完整** - 基于Hoare CSP理论的Go并发完整形式化模型 |
| 2 | [`reference/versions/06-Go-1.26特性-v3/derivation/D1-formal-semantics.md`](./reference/versions/06-Go-1.26特性-v3/derivation/D1-formal-semantics.md) | ⭐⭐⭐⭐⭐ (5/5) | **形式化推导** - Go 1.26新特性的操作语义与类型系统推导 |

### 架构设计类 (1篇)

| # | 文档路径 | 质量评级 | 核心亮点 |
|---|---------|---------|---------|
| 3 | [`architecture/clean-architecture.md`](./architecture/clean-architecture.md) | ⭐⭐⭐⭐ (4/5) | **架构选型决策树** - 四层架构设计、量化收益分析 |

### 性能优化类 (2篇)

| # | 文档路径 | 质量评级 | 核心亮点 |
|---|---------|---------|---------|
| 4 | [`go126-comprehensive-guide/37-performance-profiling-complete.md`](./go126-comprehensive-guide/37-performance-profiling-complete.md) | ⭐⭐⭐⭐ (4/5) | **pprof实战案例** - 从理论到实践的完整性能分析手册 |
| 5 | [`advanced/performance/04-PGO深度实践指南.md`](./advanced/performance/04-PGO深度实践指南.md) | ⭐⭐⭐⭐ (4/5) | **生产环境收集策略** - Profile-Guided Optimization深度实践 |

### 数据技术栈类 (1篇)

| # | 文档路径 | 质量评级 | 核心亮点 |
|---|---------|---------|---------|
| 6 | [`architecture/tech-stack/data/ent-orm.md`](./architecture/tech-stack/data/ent-orm.md) | ⭐⭐⭐⭐ (4/5) | **性能基准测试** - Ent ORM选型论证与优化技巧 |

---

## 🏷️ 质量标签说明

| 标签 | 含义 | 标准数量 |
|-----|------|---------|
| ⭐⭐⭐⭐⭐ | 卓越 | 5条标准全满足 |
| ⭐⭐⭐⭐ | 优秀 | 满足4条标准 |
| ⭐⭐⭐ | 良好 | 满足3条标准 |

---

## 📂 文档组织结构

```text
docs/
├── core/                    # 核心高质量文档（本索引所列）
│   └── （待迁移整理）
├── archive/                 # 归档文档（保留但不维护）
│   ├── by-category/         # 按主题分类的归档
│   ├── low-quality/         # 低质量文档
│   └── ...
├── architecture/            # 架构文档
├── advanced/                # 高级主题
├── go126-comprehensive-guide/  # Go 1.26综合指南
├── reference/versions/      # 版本特性参考
└── ...                      # 其他功能性文档
```

---

## 🔍 快速导航

### 按主题查找

- **学习Go并发**: ➡️ [`05-csp-formal-model.md`](./go126-comprehensive-guide/05-csp-formal-model.md)
- **性能优化**: ➡️ [`37-performance-profiling-complete.md`](./go126-comprehensive-guide/37-performance-profiling-complete.md) + [`04-PGO深度实践指南.md`](./advanced/performance/04-PGO深度实践指南.md)
- **架构设计**: ➡️ [`clean-architecture.md`](./architecture/clean-architecture.md)
- **数据库ORM**: ➡️ [`ent-orm.md`](./architecture/tech-stack/data/ent-orm.md)
- **形式化语义**: ➡️ [`D1-formal-semantics.md`](./reference/versions/06-Go-1.26特性-v3/derivation/D1-formal-semantics.md)

---

## 📈 质量统计

| 分类 | 数量 | 占比 |
|-----|------|-----|
| 形式化语义 | 2篇 | 33% |
| 架构设计 | 1篇 | 17% |
| 性能优化 | 2篇 | 33% |
| 数据技术 | 1篇 | 17% |
| **总计** | **6篇** | **100%** |

---

## 📝 维护说明

- **核心文档**: 持续维护，确保内容最新
- **归档文档**: 保留历史价值，但不再更新
- **纳入标准**: 新文档需满足≥3条质量标准方可加入核心列表

---

*最后更新: 2026-03-09*
