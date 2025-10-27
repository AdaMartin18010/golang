---
title: 发布说明 v2025.09-P1
slug: release-2025-09-p1
---

## 亮点

- 新增站点索引与版本矩阵，确立 docs 为唯一发布源；
- 建立 examples 模块并提供并发/ServeMux/PGO/slog 可运行示例；
- 文档元数据模板落地（10+篇），Markdown 规范与链接检查上线；
- 三域（并发/测试/架构）交叉链接与首版去重报告发布。

## 变更详情

- 新增：`docs/SITE_INDEX.md`、`docs/VERSION_MATRIX.md`、`docs/_META_TEMPLATE.md`；
- 新增：`examples/` 及并发、ServeMux、PGO、slog 示例；
- 新增：`CROSSLINK_PLAN.md`、`DEDUPE_REPORT.md`；
- 更新：多处 README 增加元数据头与相关链接。

## 升级指引

- 文档引用请以 `docs/` 为准；
- 示例运行：`cd examples && go test ./...` 或按子目录说明运行；
- 历史页已保留跳转提示，建议更新外部链接至主线入口。
