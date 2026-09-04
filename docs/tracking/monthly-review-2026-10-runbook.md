# 2026-10 月审运行手册 + 季度审计排期

> **用途**：2026-10 月度语义审查（模板 `.kimi/templates/monthly_semantic_review.md`）的全量复跑手册，数据全部来自机器校验脚本，零手工统计。季度审计排期见 §3。
> **前置**：Go 1.27.1（`GOWORK=off` 进模块目录）、python 3、pwsh（质量门）。

## 1. 度量复跑清单（按序执行，全部命令在仓库根）

| 指标 | 命令 | 本期基线（2026-09-05） |
| --- | --- | --- |
| 死链 | `python scripts/tmp/rescan_deadlinks.py` | TOTAL_DEAD=0 |
| 六件套 | `python scripts/tmp/verify_sixpiece.py` | 30/30 |
| 编号一致性 | `python scripts/tmp/id_consistency_check.py` | MISMATCH=0（CHECKED=421） |
| 双向链接 | `python scripts/tmp/check_bidir.py` | 103/103 |
| stub 纯净度 | `python scripts/tmp/stub_purity_check.py` | 56 stub，VIOLATIONS=0 |
| 代码块全量 | `python scripts/tmp/extract_go_blocks.py && python scripts/tmp/test_go_blocks.py && python scripts/tmp/merge_test.py` | 226 块；compile_fail 20/20；runnable BAD=0；merge BAD=0（OK=33 / skip_dep=45 / exempt=3 / known-sibling=11） |
| 索引一致性 | `python scripts/tmp/diff_map_actual.py` | 编号文档双向一致（README 除外） |
| CI 质量门 | `pwsh -File scripts/check-unfixed-links.ps1` | TOTAL_DEAD=0，Exit 0 |
| gofmt 存量 | `gofmt -l archive/ examples/ scripts/ go-knowledge-base/ test/` | 289 个未格式化（CI 不强制，仅记录趋势） |

## 2. 本期已解项（下期月审不应再出现）

- 编号漂移 98 篇（Q0 已清零）；GC × cgo 专页（LD-053 已建）；examples go.sum（已补齐）；片段豁免 24 组（已收敛为 3 块标注豁免）；FT-009 ID 复用（实质页 → FT-044）；RUN-BAD 14（已清零，11 项入豁免表 + 4 页真缺陷修复）；CI 死引用与 `\|\| true`（已修）；go.work 覆盖 8→42 模块，go-test 矩阵化 38 模块；质量门脚本重写并挂载 quality-gates.yml。

## 3. 季度国际源审计排期（2026-Q4）

| 批次 | 时间 | 范围 | 对照源 |
| --- | --- | --- | --- |
| Q4-1 | 2026-10 月审同期 | LD 维度抽样 20 页定义漂移精读 | Go Spec + Proposals |
| Q4-2 | 2026-11 月审同期 | EC/TS 维度边界语义 + 版本特性→权威页映射表 | Go Blog + pkg.go.dev |
| Q4-3 | 2026-12 月审同期 | 全仓 References 存活率抽验 10% | 全部 P0 源 |

- 模板：`.kimi/templates/quarterly_international_source_audit.md`。
- 2026-10 月审精读重点（上期移交）：定义漂移逐页精读（本期仅机械核查）、边界域语义精读、版本特性→权威页精确映射表。
- 已知待办：release.yml 第 81 行 YAML 语法错误（既有）；docs-deploy 部署脚本不复制 go-knowledge-base（页脚链接站点内不可达，需改复制逻辑）；AD-038 blocks 3/7 为 package main 无 main 的库式片段（同页惯例，非缺陷）；AD-037 block 5 领域符号未定义（zap 依赖组内，未修）；check_quality.ps1 加 `-FailOnHigh` 后可从报告门升强制门。
