# 内容合规审计基线（2026-09-04）

> 按 `AGENTS.md` + `.kimi/templates/kimi_project_requirements.md` 要求的全项目基线盘点。
> 治理决策：**分层推进**（机器可检项全量；六件套只做核心 30 页 + 新增页强制）；根目录历史文档归档 archive/。
> 详细死链清单见 [dead-links-2026-09.md](dead-links-2026-09.md)。

---

## 1. 规模与登记缺口

| 区域 | 文件数 | 说明 |
| --- | --- | --- |
| go-knowledge-base 五维 | 668 | FT-59 / LD-37 / EC-200 / TS-65 / AD-44 编号页 + ~263 篇 NN- 前缀/未编号 |
| indices | 11 个索引 | **登记滞后**：complete-index 仅 120 篇；by-topic 实际引用 403 篇；LD-013~034 等大量未登记 |
| docs | 38 | 根级 4 篇 TS-034~037 错位（应属 04-Technology-Stack） |
| view | 60 | formal/ 子目录死链重灾区 |
| examples | 55（6 个 go module） | — |
| 仓库根 | ~20 篇历史状态/计划文档 | 待归档 archive/ |

## 2. 元数据齐全度（逐页扫描，头部 3KB 检测）

| 维度 | 总数 | 标签 | Go 版本 | Bloom | 定理链 | mindmap | 反例 | P0/P1/P2 | References |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 01-Formal-Theory | 85 | 58 | 0 | 0 | 0 | 0 | 0 | 9 | 72 |
| 02-Language-Design | 113 | 48 | 11 | 0 | 0 | 0 | 1 | 8 | 52 |
| 03-Engineering-CloudNative | 261 | 190 | 0 | 0 | 0 | 0 | 0 | 14 | 103 |
| 04-Technology-Stack | 121 | 79 | 3 | 0 | 0 | 0 | 0 | 18 | 107 |
| 05-Application-Domains | 88 | 27 | 0 | 0 | 0 | 0 | 0 | 7 | 23 |
| **合计** | **668** | **402** | **14** | **0** | **0** | **0** | **1** | **56** | **357** |

结论：标签覆盖率 60%、References 53%，其余字段（Go 版本/Bloom/定理链/mindmap/反例/P0P1P2）几乎为零 → P2 头部统一 + P3 核心试点。

## 3. 死链基线

- 全仓扫描相对 md 链接 **5981 条**，死链目标 **326 个**（完整清单见 dead-links-2026-09.md）
- 重灾区分类：
  1. `VISUAL-ATLAS.md`（repo 根缺失）被 view/formal/** 引用 ~40 次
  2. `view/formal/Go/**` 内部 `./Go/...`、`../02-language-analysis/...` 等错误相对路径 ~50 次
  3. `pkg/*/README.md` → `../../docs/00-框架拓展计划.md`（不存在）~7 次
  4. `deployments/*/README.md`、`scripts/deploy/README.md` → `docs/deployment/*.md`（目录已重构）~10 次
  5. `examples/advanced/**` → `docs/02-Go语言现代化/**`（已重构）~10 次
  6. `docs/architecture/tech-stack/**` → `../00-技术栈概览.md` 等同级缺失文件 ~6 次
  7. `go-knowledge-base/templates/template-*.md` 占位链接（`LD-XXX-Name.md` 等）~23 次——模板占位符，改为 `<占位说明>` 风格即合规
  8. `archive/**` 内部链接（只读区，低优先，随归档政策不修复内容，仅登记）
  9. `.github/README.md:268`、`docs-deploy.yml:211` → `docs/DOCUMENT_STANDARD.md`（不存在）
- **设施问题**：`scripts/check-unfixed-links.ps1` 实为硬编码旧 docs 路径的静态检查器（非死链扫描器），需重写为真实扫描器或修正 AGENTS.md 对其的描述

## 4. 命名与结构

- 五维内无中文/空格文件名 ✅；docs/view/根目录待 P1 B1.4 扫描
- 编号连续性：LD-001~037 连续 ✅；FT/EC/TS/AD 编号存在但含未编号页 ~263 篇（P2 头部规范时统一处理，不强制重新编号）
- stub 文件 0 篇；重复权威页待 B1.6 抽样核查

## 4b. B1.6 重复权威核查结果（2026-09-05）

**方法**：641 篇权威页标题词元两两比对 + 编号唯一性全量扫描（原为"抽查 10-15 组"，实际升级为全量，因为发现了系统性问题）。

**发现一：FT 维度 43 篇占位垃圾页**。`01-Formal-Theory/` 下 43 个文件（含 8 个 README）约 98% 内容为自动生成模板残留（`## 1. 主题 1`、`E_1 = mc^2 + x^2 + …`、`定义 1.1 (核心概念 1)` 等），无任何实质内容。
**处置**：全部 stub 化（≤15 行：一句话主题 + canonical 链接或"待编写"声明 + complete-map 指针）。其中有 26 篇指向真实 canonical（如 FT-001 内存模型→LD-001、FT-011 Gossip→FT-027），17 篇主题暂无权威页，标注待编写。

**发现二：7 篇真实内容重复页**。LD-023（→LD-008）、EC-097（→EC-041）、EC-099-Kubernetes-134（→EC-099）、TS-001-PostgreSQL-18（→TS-001）、TS-004-Elasticsearch-90（→TS-012）、TS-012-Elasticsearch-Internals-Formal（→TS-012）、03-Network/06-NATS（→TS-008）。全部 stub 化指向 canonical。

**发现三：编号冲突 111 组（去 stub 后）**。EC/FT/LD/TS/AD 五维均存在多系列编号撞车。处置：冲突组保留入链最多者（平票取大文件），其余 143 篇重编为维度最大号递增（EC→EC-219、LD→LD-052、TS→TS-062、AD→AD-044、FT→FT-041），全仓引用（明文+URL编码）同步更新，索引经引用替换自动一致。

**残留口径**：19 个编号仍有一名 stub 别名持有者（stub 头部已声明占位与 canonical，非权威页，不再冲突）。FT 维度实质内容页 = 59 − 43 = 16 篇 + 跨维 stub 指向。

**基线修正**：五维 641 篇中，stub/占位 51 篇（43 篇内容页 + 8 个 README，合计 8.0%），实质权威页 590 篇。FT 维度内容最薄（16 篇实质页），后续内容建设优先级见月审记录。

## 4c. P2 头部元数据规范化结果（2026-09-05）

- 范围：五维 641 篇（590 实质页 + 51 stub；stub 仅保留维度/级别/状态三行，不强制四字段）。
- 处置：批量规范器补齐/标准化 1612 字段——级别一律按实测大小重算（S >16KB / A >8KB / B），缺失标签给维度默认标签，缺失 Go 版本给 `1.27+`（工具链基线）；3 篇英文头部（AD-010/AD-028/02-Authentication）转中文字段。
- 终验：**590 实质页四字段覆盖率 100%**（全量扫描 + 随机抽检 20 页双通道）。
- 附带发现：`01-Semantics/01-Operational-Semantics.md` 为单行字面 `\n` 垃圾文件（逃避了行式占位扫描），已 stub 化（累计垃圾页 44 篇）。
- dup 复查第二批（2026-09-05）：score≥6 剩余候选复查，新确认 3 组真实重复并 stub 化：`LD-029-Go-126-Features` → `LD-026-Go-126-New-Features`（浅层 Release 概览并入深度页）；`LD-030-Go-127-Roadmap` → `LD-027-Go-1-27-Roadmap-Features`（同上）；`LD-049-Go-Generics-Formal`（20KB）→ `LD-040-Go-Generics-Formal`（35KB，章节已覆盖实现/基准/模式）。stub 计数 51 → 54。

## 5. 代码块新鲜度

- `examples/` 6 模块构建/测试此前已全量验证（Go 1.27 迁移批次）✅
- 权威页内嵌 ```go 块未实测 → 纳入 P3 D3.2 抽样实测

## 6. 阶段验收对照表

| 阶段 | 任务 | 验收标准 | 状态 |
| --- | --- | --- | --- |
| P0 | 基线盘点 | 本文档 + 死链清单 | ✅ |
| P1 | indices 重建 | complete-index 登记 == 实际篇数 | ✅ |
| P1 | 错位归位 | docs/ 根部无维度前缀文件；TS-034~037 入 04-Technology-Stack | ✅ |
| P1 | 根目录归档 | 根目录仅留活跃文件 | ✅ |
| P1 | 命名核查 | 0 中文/空格文件名 | ✅ |
| P1 | 死链修复 | 326 → 0（archive/ 与模板占位符除外并注明） | ✅ |
| P1 | 重复权威 | 同主题唯一权威页（50 页 stub 化、143 页重编号，见 §4b） | ✅ |
| P2 | 头部统一 | 590 实质页维度/级别/标签/Go 版本四字段齐全（100%，见 §4c） | ✅ |
| P2 | Bloom/概念链 | 版本页 + 核心 30 页标注 | ⬜ |
| P3 | 六件套试点 | 核心 30 页反例+mindmap+P0P1P2 | ⬜ |
| P3 | 代码块实测 | 抽样页 go 块 build/test 通过 | ⬜ |
| P3 | 双向链接 | 版本特性 ↔ 权威页双向可达 | ⬜ |
| P4 | 外围治理 | examples/view/docs/pknowledge 合规 | ⬜ |
| P5 | 持续机制 | 首期月审记录 + CI 对齐 | ⬜ |
| P6 | 终验 | 死链 0、索引一致、vet 抽样过、分阶段提交 | ⬜ |
