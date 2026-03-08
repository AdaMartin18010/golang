# 文档归档报告

**日期**: 2026-03-08
**归档操作**: 分类整理（移动而非删除）

---

## 归档统计

| 类别 | 数量 | 说明 |
|------|------|------|
| knowledge-maps | 59 | 知识图谱、概念定义、总览、导图 |
| navigation | 36 | README、INDEX、索引、导航 |
| reports | 22 | 总结、报告、计划、roadmap、changelog |
| matrices | 18 | 对比矩阵 |
| getting-started | 10 | 快速开始、入门、FAQ、cheatsheet |
| templates | 3 | 模板文件 |
| other | 0 | 其他 |
| **总计** | **148** | |

---

## 目录结构

```
docs/archive/
├── by-category/
│   ├── knowledge-maps/    # 59个 - 知识类文档
│   ├── navigation/         # 36个 - 导航类文档
│   ├── reports/            # 22个 - 报告类文档
│   ├── matrices/           # 18个 - 对比矩阵
│   ├── getting-started/    # 10个 - 入门类文档
│   ├── templates/          # 3个 - 模板
│   └── other/              # 0个 - 其他
├── low-quality/            # 历史归档的低质量文档
└── ARCHIVE-REPORT.md       # 本报告
```

---

## 核心文档保留

**保留在原地的高质量文档**: 476个

**主要保留目录**:

- `docs/architecture/` - 架构设计文档
- `docs/go126-comprehensive-guide/` - 深度技术指南
- `docs/development/` - 开发指南
- `docs/advanced/` - 高级主题
- `docs/security/` - 安全文档
- `docs/deployment/` - 部署文档
- `docs/formal-specs/` - 形式化规格

---

## 归档原则

1. **模板化内容** → 归档至 `templates/`
2. **临时性报告** → 归档至 `reports/`
3. **索引导航** → 归档至 `navigation/`
4. **重复知识图谱** → 归档至 `knowledge-maps/`
5. **对比矩阵** → 归档至 `matrices/`
6. **入门FAQ** → 归档至 `getting-started/`

---

## 后续建议

1. **定期审查**: 每季度审查archive目录，清理过时内容
2. **恢复机制**: 如需恢复，直接从相应类别目录移回
3. **文档合并**: 考虑将重复的knowledge-maps合并
4. **核心精简**: 继续精简核心文档至<200篇

---

**归档操作完成** ✓
