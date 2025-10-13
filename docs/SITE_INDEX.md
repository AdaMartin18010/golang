---
title: 站点总索引（唯一发布源）
slug: site-index
topic: navigation
level: overview
goVersion: 
lastReviewed: 2025-09-15
owner: core-team
status: active
---

## 站点总索引（docs 为唯一发布源）

> 本索引作为全站导航入口。`docs/` 为唯一发布源；`model/` 作为素材库；`00-备份/` 已归档，不再直接面向发布。

### 导航地图

- 基础 → `01-Go语言基础/`
- 并发主线 → `12-并发编程/`、`02-Go语言现代化/03-标准库增强/03-并发原语与模式/`
- 现代化与新特性 → `02-Go语言现代化/`
- 测试与质量 → `18-最佳实践/`、`04-质量保证体系/`
- 架构与模式 → `03-软件体系架构/`、`14-设计模式/`、`02-Go语言现代化/06-架构模式现代化/`
- 性能与工具链 → `15-性能优化/`、`02-Go语言现代化/05-性能与工具链/`
- 云原生与部署 → `19-云原生与部署/`、`02-Go语言现代化/09-云原生2.0实现/`
- 可观测性 → `docs/observability/OTel-方案.md`
- 行业与项目实践 → `16-行业应用/`、`17-项目实践/`

### 快速入口

- 版本矩阵与兼容性 → `docs/VERSION_MATRIX.md`
- 架构专题总览 → `docs/architecture_README.md`
- 质量报告与技术报告 → `docs/QUALITY_REPORT.md`、`docs/TECHNICAL_REPORT.md`
- 发布说明 → `RELEASE_NOTES.md`
- 讨论与提案 → `DISCUSSIONS.md`
- FAQ 汇总 → `FAQ.md`
- 报表与度量 → `reports/README.md`
- 可观测性（OTel 方案） → `docs/observability/OTel-方案.md`
  - Go 接入规范 → `docs/observability/Go-接入规范.md`
  - SLO 与成本策略 → `docs/observability/SLO-与成本策略.md`
- Examples（可运行示例） → `examples/README.md`
  - 并发 → `examples/concurrency/`
  - ServeMux（1.22+） → `examples/servemux/`
  - PGO（1.21+） → `examples/pgo/`
  - slog（1.21+） → `examples/slog/`

### 学习路径（建议）

1. 基础 → 并发 → 测试 → 现代化新特性（slog/ServeMux/PGO）
2. 架构与模式（Clean/Hex）→ 微服务与云原生
3. 性能优化与工程化 → 发布与度量

### 维护约定

- 所有文档需包含元数据头（见 `docs/_META_TEMPLATE.md`），并在更新后刷新 `lastReviewed`。
- 新专题必须提供：对比表、可运行示例、基准/结论、版本适配说明。
- 链接使用相对路径，避免指向 `00-备份/` 与 `model/`。
- 文档质量要求：内容完整、代码可运行、示例丰富、结构清晰。

### 文档质量标准

- **内容完整性**: 涵盖核心概念、实践案例、最佳实践
- **代码质量**: 可运行、有注释、遵循Go规范
- **结构清晰**: 逻辑清晰、层次分明、易于导航
- **更新及时**: 保持与Go版本同步，定期更新

### 贡献指南

- 遵循统一的文档模板和格式
- 提供完整的代码示例和测试
- 包含必要的图表和说明
- 定期更新和维护内容
