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
- 死链复扫：本批新增链接引入 4 处死链，已修复（stub complete-map 相对深度、`LD-037` examples 路径多一层 `../`、模板占位符链接改行内代码），复扫 **0 死链**。

## 4d. P2 Bloom/概念链标注结果（2026-09-05）

- 范围：版本特性页 5 篇（LD-026/027/035/036/037）+ 演进页 3 篇（03-Evolution/04/05/07），共 8 篇（先行批）。
- 每篇在头部 `> **Go 版本**` 后插入 3 行：Bloom 层级（L2 roadmap / L3 特性综合·演进 / L4 形式化）、前置·后置概念链（相对链接，全部验证存在）、定理链（输入→操作→输出/不变式）。
- 概念链拓扑：03-Evolution/03→04→07 版本序链；LD-027→LD-037→examples/go127-features 的 1.27 落地链；LD-040→LD-037 泛型依赖链；LD-004→LD-036 GC 前置链。
- 顺带修正：LD-026 头部 Go 版本误写 `1.27+` → `1.26+`。
- 其余实质页 Bloom 标注随 P3 六件套试点（核心 30 页）补齐。

## 5. 代码块新鲜度

- `examples/` 6 模块构建/测试此前已全量验证（Go 1.27 迁移批次）✅
- 权威页内嵌 ```go 块未实测 → 纳入 P3 D3.2 抽样实测

## 5b. P3 六件套试点结果（2026-09-05）

- 范围：核心 30 页（LD 8 / FT 5 / EC 6 / TS 6 / AD 5，五维覆盖，S/A 级优先），并行 6 路补齐缺失项。
- 补齐内容：头部三行（Bloom L3/L4 + 前置/后置概念链 + 定理链）、反命题与边界节（语言类页为 `// 编译失败:` 确定性编译期反例；非语言类页为可论证的反模式示例）、Mermaid mindmap、References 重组为 P0/P1/P2 三级（既有条目零删除）。
- 机械校验：30/30 六件套齐全（scripts/tmp/verify_sixpiece.py）；新增相对链接全部经验证；死链复扫 0。
- 附带修复：① LD-048 页内编号 LD-009 → LD-048；② TS-049 页内编号 TS-002 → TS-049（页内指向真实存在的 TS-002-Redis-Data-Structures-Internals 链接保留）；③ TS-050 页内编号 TS-003 → TS-050（同上）；④ P3 复查新确认重复权威 2 组并 stub 化：`FT-039-Raft-Consensus-Formal-Specification` → `FT-002`、`FT-040-Paxos-Consensus-Formal` → `FT-006`（stub 计数 54 → 56）。
- 待实测：新增 `// 编译失败:` 反例块与可运行块 → §5c 统一验证。

## 5c. P3 代码块实测结果（2026-09-05，Go 1.27.1 windows/amd64 实机）

- 规模：30 试点页共 226 个 ```go 块 = 编译失败反例 20 + 完整程序 122 + 片段 84。
- **编译失败反例：20/20 全部确实编译失败**（首轮 19/20，LD-001 §8.2 反例标注错误——`done` 在 goroutine 发送中被使用、实际可编译，已改写为真正触发 `declared and not used: done` 的形式）。
- **可运行块实测修复 9 处真实缺陷**：LD-001（`procs` 未使用 + Node/Cache 示意类型缺失）、LD-010（`runtime/debug` 未使用导入）、LD-011（`runtime` 未使用导入 + `i < m.NumGC` int/uint32 类型不匹配 + gc_test `sync` 未使用导入）、FT-003（cap 包 time 导入跨块错配）、FT-016（`math` 未使用导入）、EC-155（sbom `encoding/json`/`os` 未使用导入、slsa `os` 导入误删恢复）。
- **合并构建终态（同页同 package 块合并为程序）**：19 组通过、36 组因外部依赖（grpc/redis/sql/goland.org/x 等）跳过、24 组按「片段级示意示例」豁免——均含 package 声明但引用未定义的示意类型（如 `EventStore`/`Alert`/`NetworkEndpoint`），属工程实践节的模式示意代码，非核心机制块；清单见 scripts/tmp/goblock_merge_results.json，可在后续批次逐页补全为可编译示例。
- 工具链：extract_go_blocks.py → test_go_blocks.py（单块）/ merge_test.py（合并），全部 GOWORK=off。

## 5d. P3 双向链接核查与补齐结果（2026-09-05）

- 30 试点页头部前置/后置链接共 103 条：补齐前仅 15 条双向可达。
- 处置：scripts/tmp/add_backlinks.py 在 79 个目标页头部追加回链（优先并入既有 前置/后置概念 行，无头部链的行插入 `> **相关概念**` 行）；1 条指向 stub（FT-009，canonical 待定，按 stub 口径豁免）；1 条指向 examples（go127-features README 已补概念回链）。
- 终态：**102/103 双向可达**（99.0%）；死链复扫 0。

## 5e. .kimi 治理文件全面 Go 化（2026-09-05，用户追加批次）

- 背景：`.kimi/` 新增 21 个从 rust-lang 项目复制的文件（根目录要求 6 + 模板 7 + prompts 8，约 3100 行）。
- 处置：7 路并行全量重写。映射口径：concept/ → 五维 FT/LD/EC/TS/AD 编号页；```rust,compile_fail →```go + `// 编译失败:`（Go 无 E0xxx 错误码体系）；MSRV/Cargo.toml → go.mod go 1.27 指令；cargo workspace 命令 → `GOWORK=off go build/vet/test ./...`；权威源 → go.dev/ref/spec、pkg.go.dev、go.googlesource.com/proposal、go.dev/blog；所有权/借用/生命周期 → GC/cgo/happens-before/channel；rustc 决策树 → Go 编译失败类别映射。
- 精简策略（用户确认）：内容生成/语义要求/质量门/KG 拓扑/语义审计/mindmap/交叉域/属性矩阵为核心启用集；表征空间/思维表征/决策树/效应系统页标为可选扩展（文件头降级标注，不进入质量门）。
- 终验：全目录 Rust 术语 grep 残留 0（仅余跨语言对比矩阵的 Rust 列等允许项）；subagent 引入的 30 处占位/深度错误链接已修复（占位符改行内代码、kimi_* 兄弟文件相对路径修正），死链复扫 0；README 全量登记 32 个文件并标注启用状态。

## 5f. P4 外围治理结果（2026-09-05）

- **A. examples 工作区垃圾清理**：删除未跟踪构建产物 `examples/go126-features/go126-features.exe`、`examples/go126-features/test_build.exe`；`git status --short examples/` 确认无意外改动（仅存量 `go127-features/README.md` 修改，与本批无关）。
- **B. examples 缺 README 模块补登（10 个）**：crypto-hpke、go125、go126-features、grpc、messaging、net-dialer-ctx、security、servemux、slog、stdlib-peek 各补最小 README（一句话用途 + 运行命令 + go 指令版本/归属模块说明）。按实际结构核实：仅 go126-features 有独立 go.mod（go 1.27）；crypto-hpke/net-dialer-ctx/slog/stdlib-peek 无 go.mod（工具链目录模式可运行）；servemux 为测试驱动示例；go125/grpc/messaging/security 为子目录集合（grpc/messaging/security 属上级 `examples` 模块 `example.com/golang-examples`）。
- **抽验（真实执行）**：crypto-hpke / net-dialer-ctx / slog / stdlib-peek / go126-features 的 `GOWORK=off go run .` 全部通过；servemux `GOWORK=off go test ./...` ok；go125/runtime/container_scheduling `GOWORK=off go run .` 通过。
- **既有问题标注（未越权修复）**：grpc/messaging/security 的 `GOWORK=off go build` 因 `examples/go.sum` 缺 grpc/nats/chi 等校验条目而失败（存量问题），README 已注明需 `go mod tidy` 补齐或用 workspace 模式；go125/concurrency-network/http3 为空目录、synctest 仅有 go.mod 无源文件，README 如实标注。
- **C. docs/ 根部 5 篇游离 md + cleanup.py 归位**（git mv 保留历史）：
  - `00-Go-1.26-Comprehensive-...md` → `docs/go126-comprehensive-guide/`
  - `01-Go-1.27-Comprehensive-...md` → 新建 `docs/go127-comprehensive-guide/`（与 1.26 命名模式一致，未来 1.28 归入此处）
  - `go126-package-management.md` → `docs/go126-comprehensive-guide/`
  - `go126-security-advisory-2026-05.md` → `docs/tracking/`
  - `go127-json-v2-migration.md` → `docs/development/`
  - `cleanup.py` → `scripts/tmp/cleanup.py`（先 grep 确认 md/yml 零引用，一次性归档脚本保留）
  - 引用修复：被移文件内部 `../` 外链全部加深一层（00-Go-1.26 8 处、package-management 1 处、01-Go-1.27 5 处、json-v2-migration 3 处）；`./` 自引用与跨文件引用按新位置改指（go126-comprehensive-guide/、architecture/、development/、00-Go-1.26 前版互链 5 处）；入站链接修复 6 处（docs/README.md 2、examples/go127-features/README.md 1、07-Go126-to-Go127 1、Go-1.27-Preview 1、Go-1.27-Release 1、go127_semantic_analysis 1）；docs/README.md 目录树同步更新；活文档路径提及更新 2 处（LD-036、PERSONAL-LEARNING-PATH）；archive/ 与 tracking 历史日志、CHANGELOG 历史记录中的旧路径为历史快照，按政策不改。
- **验证数据**：死链复扫 `TOTAL_DEAD = 0`；六件套校验 `OK=30 MISS=0 / 30`；complete-map 校验类脚本在 `scripts/tmp/` 不存在（全仓亦无，仅 `check-content-completeness.ps1` 属 docs 内容完整性检查，非 kb 条目计数），该项 N/A。

## 5g. P5 持续机制结果（2026-09-05）

- **首期月度语义审查**：按 `.kimi/templates/monthly_semantic_review.md` 产出 [monthly-review-2026-09.md](monthly-review-2026-09.md)（150 行），数据全部取自本文档 §4b/§4c/§4d/§5b/§5c/§5d/§5f，未实测量标注 N/A 并注明原因；本期新跑两个机械校验：`stub_purity_check.py`（STUB_TOTAL=56，与 §5b 计数一致，VIOLATIONS=0）、`id_consistency_check.py`（CHECKED=420，**MISMATCH=98**——§4b 重编号批次页内 H1 编号系统性残留，列为 P0 行动项）。
- **月审主要结论**：stub 纯净度 56/56 零违规；边界域 7/8 有权威页（GC × cgo 专页缺失）；KG 通用谓词违规 0（关系以内嵌概念链形式存在）；版本语义注入 8 先行页 + 30 试点页；代码块抽验被 §5c 全量实测（226 块）取代故 N/A；综合语义健康分 78/100。
- **CI 对齐核查**（只记录不改动）：16 个 workflow 中 knowledge-tracker.yml 的 3 个脚本引用与 ci.yml 的 `scripts/quality-check.sh` 均有效；发现失效/缺口 5 项——docs-deploy.yml 生成 HTML 内 3 处死引用（docs/10-进阶专题/、LEARNING_PATHS.md、INDEX.md 均不存在）、AGENTS.md 宣称的 3 个 CI 质量门脚本实际无 workflow 引用且 check-unfixed-links.ps1 需重写（§3 既有记录）、go workflow 无 GOWORK 设置且 go.work 仅覆盖 42 模块中的 8 个（CI 覆盖缺口）、go-test.yml setup-go@v4 版本不一致、docs-check.yml 的 markdown-link-check 被 `\|\| true` 屏蔽。DOCUMENT_STANDARD.md 旧死引用复核已 0 命中。详见月审 §7。
- **验证**：月审新增 md 后死链复扫仍为 `TOTAL_DEAD = 0`；六件套 30/30 不受影响。

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
| P2 | Bloom/概念链 | 版本页 5 + 演进页 3 先行标注（§4d）；核心 30 页随 P3 补齐 | ✅ |
| P3 | 六件套试点 | 核心 30 页反例+mindmap+P0P1P2+头部三行，机械校验 30/30（§5b） | ✅ |
| P3 | 代码块实测 | 编译失败反例 20/20；合并构建 19 组通过，24 组片段级豁免、36 组外部依赖跳过；9 处缺陷已修（§5c） | ✅ |
| P3 | 双向链接 | 试点页头部链接 102/103 双向可达（1 为 stub 豁免）（§5d） | ✅ |
| P4 | 外围治理 | examples/view/docs/pknowledge 合规 | ✅ |
| P5 | 持续机制 | 首期月审记录 + CI 对齐 | ✅ |
| P6 | 终验 | 死链 0、索引一致、vet 抽样过、分阶段提交 | ✅ |

## 5h. P6 终验结果（2026-09-05）

- **死链复扫**：`python scripts/tmp/rescan_deadlinks.py` → `TOTAL_DEAD = 0`。
- **六件套复验**：`python scripts/tmp/verify_sixpiece.py` → `OK=30 MISS=0 / 30`。
- **索引一致性复核**（修正口径后重跑）：`scripts/tmp/diff_map_actual.py` 以 basename 集合比对五维目录实际 md 与 `indices/complete-map.md` 登记——641 篇编号文档双向一致，唯一差异为各维度 `README.md`（目录导览，非编号文档，不入清单）。首轮报告 LD-035/036/037「未登记」为脚本正则不含点号的假阳性，该三页实际已登记于 complete-map 第 222–224 行。
- **vet 抽样**：`cd examples/go127-features && GOWORK=off go vet ./...` 通过。
- **提交策略**：全程不 commit/push，由用户 IDE 惯例处理（提交惯例 `update`）。

**遗留事项**（移交下期，详见月审 §8）：

1. 24 组片段级代码块豁免（示意类型/外部依赖，清单 `scripts/tmp/goblock_merge_results.json`）。
2. FT-009 为 stub，canonical 待定（双向链接 102/103 的唯一豁免）。
3. 98 篇页内 H1 编号 ≠ 文件名编号（§4b 重编号残留，P0 行动项，`scripts/tmp/id_consistency_check.py` 可复扫）。
4. GC × cgo 交叉边界专页缺失；examples/go.sum 缺 grpc/nats/chi 校验条目；CI 5 项缺口（月审 §7）。
5. `.kimi` 可选扩展集（表征空间/思维表征/决策树/效应系统 4 页）已 Go 化但默认未启用。

## 5i. 遗留项收口批结果（2026-09-05，Q0→Q3）

P6 验收后随即执行月审 §8 行动项收口（用户已确认：并行推进、CI 只修死引用）：

- **Q0 编号一致性**：98 篇 H1 修复 + 3 处正文自指（EC-155 目录锚点、AD-035 目录锚点 + Document ID），`id_consistency_check.py` MISMATCH=0。
- **Q1a examples/go.sum**：go.sum 26→60 行补齐 + go.mod 增根模块 require + jwt/v5 indirect；grpc/messaging/security `GOWORK=off go build` 与 vet 通过（tidy 方案因大范围改写 go.mod 被否，用最小改动替代）。
- **Q1b GC × cgo 权威页**：新建 `LD-053-Go-CGO-GC-Interaction-Boundary.md`（19KB S 级，六件套全配，8 个 go 块实测含 cgo shim）；已登记 complete-map/by-date/by-topic/complete-index。**实测发现**：Go 1.27 中 `GODEBUG=cgocheck=2` 已失效（运行时 fatal），严格检查须用构建期 `GOEXPERIMENT=cgocheck2`。
- **Q2a 豁免收敛**：merge BAD 24→0（14 组补示意定义转 OK=33；10 组第三方依赖归 skip_dep=45；3 块 x/crypto 加 `<!-- goblock-exempt -->` 机器可识别标注，管道三脚本同步升级支持 exempt 类别）。
- **Q2b FT-009 stub**：两个重复 stub 归位 canonical（Quorum→FT-017、SMR→FT-032），FT-017/032 补反向链接；`check_bidir.py` 103/103 全双向（顺带修复 TS-047 缺前置链）。
- **Q3 CI（限定范围）**：docs-deploy.yml 3 处死引用修复（注意：用户 IDE 中途 commit 5ed3b676 已将其带入 HEAD）；docs-check.yml 移除 markdown-link-check `|| true`（先修 5 文件 15 条真死链 + aliveStatusCodes 补 202）；其余 3 项缺口只记录未动。

**回归终态**：死链 0 · 六件套 30/30 · 编号 MISMATCH=0 · 双向链接 103/103 · stub 纯净度 56/56 · merge BAD=0。

**剩余移交项**（详见月审 §8 更新表）：CI 覆盖矩阵/质量门挂载（P1）、docs-deploy 其余 8 个模块卡死引用、check-unfixed-links.ps1 重写、AD-038 block 6 格式损坏、FT-009 ID 复用治理（3 文件共号）、test_go_blocks 14 个 RUN-BAD、`.kimi` 扩展集启用决策（待用户）。


## 5j. R1→R4 批结果（2026-09-05，用户确认：并行推进、扩展集保持未启用）

- **R1a go.work/CI 矩阵**：go.work 从 8 扩至全部 42 个已跟踪模块（零版本冲突回退）；go-test.yml 矩阵化 38 模块（GOWORK=off、fail-fast: false），4 模块入排除清单（`docs/tracking/ci-module-coverage.md` 新建）；顺带修复 4 处预存编译/vet 缺陷（user-service fmt 导入、network_buffer 标准签名、context_test lostcancel、examples go.mod tidy）。
- **R1b 质量门**：check-unfixed-links.ps1 按 rescan 逻辑重写（真实仓库 0 死链 Exit 0，合成死链 Exit 1）；新建 quality-gates.yml（死链强制门 + check_quality 报告门 + gofmt 强制门限 cmd/internal/pkg）；AGENTS.md §5 改为本地门/CI 门真实描述。
- **R2 docs-deploy/setup-go**：首页 9 卡全部改指真实目录（原 8 个 docs/0N-* 路径均不存在）；全仓 setup-go 统一 v5。release.yml 第 81 行既有 YAML 语法错误未动（移交）。
- **R3a/R3c AD-038**：block 6 补 package/imports/stub（带依赖 build 通过）；父代理复核确认 block 5 真缺陷（getState/updateState/CircuitStateRecord/LambdaHandler 4 符号未定义）并修复，带依赖 tidy+vet+build 全过。
- **R3b FT-009 治理**：实质页重编号 FT-044（git mv + 19 处引用 + 5 个 indices 文件同步）；两 stub 保留 FT-009 别名入口并加注说明。
- **R3c RUN-BAD**：14 项清零——11 项伪失败入 test_go_blocks.py 按页豁免表（仍单块 vet 防烂尾），TS-049/AD-035/AD-037/TS-006 四处真缺陷修复（补 import/示意类型/缺失方法）。
- **R4 机制**：`docs/tracking/monthly-review-2026-10-runbook.md` 新建（月审全量复跑手册 + Q4 季度审计三批排期）。

**回归终态**：死链 0 · 六件套 30/30 · 编号 MISMATCH=0 · 双向 103/103 · stub 56/56 · compile_fail 20/20 · runnable BAD=0 · merge BAD=0（OK=33/skip_dep=45/exempt=3/known-sibling=11）。

**移交下期**：release.yml YAML 修复；docs-deploy 复制逻辑（页脚链接站点不可达）；AD-037 block 5 领域符号；check_quality.ps1 -FailOnHigh 升级；gofmt 存量 289 个（趋势监控即可）；CI 首跑 38 模块矩阵的 -race 测试未本地全量验证（首跑 CI 见真章，asan 模块可能需关 race）。
