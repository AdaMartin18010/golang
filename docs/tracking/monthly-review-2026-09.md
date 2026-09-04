# 月度语义审查记录（2026-09）

> **Review period**: 2026-09（首期）
> **Reviewer**: Kimi Code CLI（代理执行，人工复核待签）
> **Scope**: 五维权威页（go-knowledge-base/01..05-*）的定义漂移、边界语义、stub 纯净度、KG 关系质量、版本语义注入覆盖。
> **数据源**：本期所有度量取自 [content-compliance-audit-2026-09.md](content-compliance-audit-2026-09.md) §4b/§4c/§4d/§5b/§5c/§5d/§5f + 本期新跑的两个机械校验脚本（`scripts/tmp/stub_purity_check.py`、`scripts/tmp/id_consistency_check.py`）；未实测量标注 N/A 并说明原因。
> 归档：`docs/tracking/`（本文件）。

---

## 1. Core Concept Definition Drift（抽样 10 页）

抽样方式（首期）：以 P3 六件套试点 30 页 + 边界域候选页为抽样框，做机械级核查（页内 H1 编号 vs 文件名编号、定理链存在性、关键词命中），**逐页精读定义留待下期**——首期无既有漂移基线，机械核查优先于主观判读。

| # | Page path | Definition checked | Drift found? | Action item |
| --- | ----------- | -------------------- | -------------- | ------------- |
| 1 | `02-Language-Design/LD-052-Go-Escape-Analysis.md` | 逃逸分析×栈/堆分配 | [x] Yes [ ] No | **H1 编号漂移**：页内写 `LD-012`，文件名 LD-052（§4b 重编号残留，全仓 98 篇同类，见下） |
| 2 | `02-Language-Design/LD-008-Go-Error-Handling-Patterns.md` | errors.Is/As 包装链定理链 | [ ] Yes [x] No | — |
| 3 | `02-Language-Design/LD-045-Go-Error-Handling-Formal.md` | errors.Is/As 语义表 | [ ] Yes [x] No | — |
| 4 | `02-Language-Design/LD-047-Go-Context-Formal.md` | context 取消 Happens-Before | [ ] Yes [x] No | — |
| 5 | `02-Language-Design/LD-040-Go-Generics-Formal.md` | 泛型形式化（//go:noinline 等编译指令交互） | [ ] Yes [x] No | — |
| 6 | `01-Formal-Theory/FT-001-Go-Memory-Model-Formal-Specification.md` | happens-before 公理体系 | [ ] Yes [x] No | — |
| 7 | `02-Language-Design/LD-012-Go-Linker-Build-Process.md` | 链接器行为 | [ ] Yes [x] No | — |
| 8 | `02-Language-Design/LD-010-Go-Scheduler-GMP.md` | GMP 抢占调度 | [ ] Yes [x] No | 抽查未发现编号漂移；精读待下期 |
| 9 | `02-Language-Design/LD-022-Go-Context-Propagation.md` | context 传播模式 | [ ] Yes [x] No | — |
| 10 | `02-Language-Design/LD-044-Go-Reflection-Formal.md` | 反射形式化 | [ ] Yes [x] No | — |

**本期漂移主发现（机械校验，新跑 `scripts/tmp/id_consistency_check.py`）**：420 篇编号页中 **98 篇文件名编号 ≠ 页内 H1 编号**，全部集中于 §4b 重编号批次（EC-122~158、LD-038~052 等系列页内标题沿用旧号）。属系统性定义漂移，建议下期开专项批量修正页内 H1 编号。

**What to look for 落实状态**：四字段覆盖率 100%（§4c，含 Go 版本标注对齐工具链基线 `1.27+`）；编译失败反例 20/20 在 Go 1.27.1 下实测确实失败（§5c）✅；逐页"定义可操作化"精读 → 下期行动项。

---

## 2. Boundary Precision Review

| Domain pair | 权威页 exists? | Boundary clearly stated? | Action item |
| ------------- | ---------------------------------- | -------------------------- | ------------- |
| GC × cgo/指针运算 | [ ] Yes [x] No | — | LD-004/LD-036 覆盖 GC，但 **GC × cgo 交叉边界专页未检出**（五维文件名/关键词扫描均无命中）；建议补页或明确归入 LD-036 边界节 |
| goroutine × signal/抢占调度 | [x] Yes [ ] No | [x] Yes [ ] No | LD-010/LD-033/LD-043；边界关键词命中（抢占/调度），精读确认待下期 |
| 泛型 × 反射 | [x] Yes [ ] No | [x] Yes [ ] No | LD-040 + LD-044/LD-007 配对存在 |
| 内存模型 happens-before × channel/select | [x] Yes [ ] No | [x] Yes [ ] No | FT-001 + LD-042（Channels-Formal，H1 编号漂移待修） |
| error wrapping × `errors.Is/As` 边界 | [x] Yes [ ] No | [x] Yes [ ] No | LD-008（canonical，LD-023 已 stub 并入）；定理链含 Is/As 判定 |
| context 取消 × goroutine 泄漏 | [x] Yes [ ] No | [x] Yes [ ] No | LD-022/LD-047；取消 Happens-Before 定理命中 |
| `//go:` 编译指令 × 链接器行为 | [x] Yes [ ] No | [x] Yes [ ] No | LD-012（linker）+ LD-014/LD-051（assembly，含 //go: 指令） |
| 逃逸分析 × 栈/堆分配 | [x] Yes [ ] No | [x] Yes [ ] No | LD-052（H1 编号漂移待修） |

**Notes**: 8 个边界域中 7 个有非 stub 权威页；"boundary clearly stated" 列基于边界关键词/定理链机械命中判定，**语义级精读全部列为下期行动项**，不在本期伪造结论。

---

## 3. New Stub Purity Review

- [x] Stub 正文 ≤ 25 行 / 2000 字节
- [x] 每篇 stub 恰好含一句话主题 + canonical 链接（或"待编写"声明 + complete-map 指针）
- [x] stub 不含从原页残留的正文内容（抽查头部均为 2026-09 合规审计占位声明）

**Violations found**: 0（新跑 `scripts/tmp/stub_purity_check.py`：STUB_TOTAL = 56，与审计 §5b 的 stub 计数 56 一致；VIOLATIONS = 0）

| File path | Lines | Bytes | Issue | Remediation |
|-----------|-------|-------|-------|-------------|
| —（无违规） | | | | |

---

## 4. KG Relation Quality Review

- [x] KG 关系使用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`）
- [x] 核心实体无通用 `RelationAnnotation`

**核查方式与口径说明**：本期五维 kb **无独立 KG 数据文件**（.kimi 要求文档中的 `ex:RelationAnnotation`/`relatedTo` 禁令针对 KG 导出层）；实际 KG 关系以内嵌「前置/后置概念链」行存在于页面头部（§4d 标注 8 页 + P3 试点 30 页，共 38 页，103 条链接，§5d 双向可达 102/103）。全仓 grep `RelationAnnotation` 仅命中 .kimi 要求文档本身，**知识页 0 处通用谓词违规**。

**Generic relations found**: 0

| Subject | Predicate | Object | Suggested concrete predicate |
|---------|-----------|--------|------------------------------|
| — | | | |

---

## 5. Version Semantic Injection Coverage Review

- [x] Go 1.26–1.27 特性映射到权威页数 = **8 先行页 + 30 试点页标注**（§4d：版本特性页 5 + 演进页 3 先行；§5b：核心 30 页六件套补齐概念链）；模板中"= _**/**_"精确分子分母计数 N/A——审计批次以页为单位而非特性条目为单位，且无特性→页的多对一精确映射表，下期补建。
- [x] 每个版本页回链概念权威页（§4d 概念链拓扑：03-Evolution/03→04→07 版本序链；LD-027→LD-037→examples/go127-features 落地链；LD-040→LD-037 泛型依赖链；LD-004→LD-036 GC 前置链）
- [x] 每个被映射的权威页前向链接 examples/（go126-features / go127-features，§5d：1 条 examples 链接已补回链）

**口径修正**：模板列出的 LD-029/LD-030 版本页已于 §4c 批次的重复权威处置中 stub 化（分别并入 LD-026/LD-027），不再作为独立版本页参与映射。

**Unmapped or one-way-only features**: 本期未发现（双向可达 102/103，1 条为 stub 豁免，§5d）

| Feature | Version page | 权威页 | Issue | Action |
|---------|--------------|---------|-------|--------|
| — | | | | |

---

## 6. 代码块实测抽验（推荐）

- [N/A] 抽样 10 个 ```go 块提取到临时模块 `GOWORK=off go vet` 通过 —— **本期已被 P3 全量实测覆盖并超出抽样规模**（§5c：30 试点页 226 个 go 块全量提取实测，合并构建 19 组通过、9 处真实缺陷修复，Go 1.27.1 windows/amd64 实机），重复抽样无增量价值，标注 N/A（原因：被 §5c 全量测试取代）
- [N/A] 抽样 5 个 `// 编译失败:` 反例确认确实编译失败 —— 同上，§5c 已实测 20/20（首轮 19/20，LD-001 反例已改写），N/A（原因：被 §5c 全量测试取代）

**Gaps**: §5c 遗留 24 组片段级豁免 + 36 组外部依赖跳过（清单 `scripts/tmp/goblock_merge_results.json`），列为后续批次行动项。

| Page | Block | 结果 | Action |
|------|-------|------|--------|
| （见 §5c 全量清单） | | | 24 组片段豁免逐页补全为可编译示例；36 组外部依赖块建独立集成测试模块 |

---

## 7. 工具与 CI 对齐（本期新增核查）

核查方式：`ls .github/workflows/` 全部 16 个 yml 逐个读 run:/uses: 行，提取脚本路径与仓库内文件逐一比对。

**引用有效的 workflow 脚本**：

| Workflow | 引用 | 存在性 |
|----------|------|--------|
| knowledge-tracker.yml | `scripts/knowledge-tracker/requirements.txt`、`track_go_releases.py`、`track_papers.py` | ✅ 全部存在 |
| ci.yml:94 | `scripts/quality-check.sh` | ✅ 存在 |

**失效/存疑引用（仅记录，不改 workflow——超出本期范围）**：

| 位置 | 引用 | 问题 | 建议 |
| ------ | ------ | ------ | ------ |
| docs-deploy.yml 生成的站点 HTML | `docs/10-进阶专题/`、`docs/LEARNING_PATHS.md`、`docs/INDEX.md` | 三者均不存在（死引用嵌在生成的 index.html 内容里，CI 不校验） | 改写 HTML 模板或补建目标页 |
| AGENTS.md §5 宣称的 CI 质量门 | `scripts/check-unfixed-links.ps1`、`check-markdown-format.ps1`、`check_quality.ps1` | 三个脚本磁盘上存在，但 **无任何 workflow 引用**——实为本地门，AGENTS.md 描述易误导为 CI 门；且 `check-unfixed-links.ps1` 被审计 §3 记录为"硬编码旧 docs 路径的静态检查器（非死链扫描器），需重写" | 重写 check-unfixed-links 为真实扫描器（或直接用 scripts/tmp/rescan_deadlinks.py 入 CI）；AGENTS.md 改述为"本地质量门"；可考虑把 check-markdown-format/check_quality 接入 docs-check.yml |
| 全部 go workflow（go-test/test/ci/ci-enhanced 等） | `go test/build/vet ./...` | 无任何 `GOWORK` 环境设置；go.work 仅覆盖 8 个 use 条目，42 模块中多数不在 workspace 内 → **CI 实际只测 workspace 覆盖的模块**，与 AGENTS.md "每模块 GOWORK=off 验证"的口径存在覆盖缺口 | 增加逐模块矩阵 job 或在 workflow 中显式 `GOWORK=off` + 模块清单遍历 |
| go-test.yml | `actions/setup-go@v4` | 其余 workflow 均用 v5，版本不一致（功能性失效风险低） | 统一升 v5 |
| docs-check.yml | `markdown-link-check ... \|\| true` | 链接检查被 `\|\| true` 屏蔽，永远绿 | 去 `\|\| true` 或改判定阈值 |
| 历史记录 | `.github/README.md`、`docs-deploy.yml` → `docs/DOCUMENT_STANDARD.md` | 审计 §3.9 曾记录 2 处死链；**当前 grep 已 0 命中**（P1 批次已修复），本期复核确认无残留 | — |

---

## 8. Summary and Action Items

| Priority | Action | Owner | Due date |
| ---------- | -------- | ------- | ---------- |
| ~~P0~~ ✅ | ~~批量修正 98 篇重编号页的页内 H1 编号~~（已于 2026-09-05 完成：98 篇 H1 + 3 处正文自指，MISMATCH=0） | 内容维护 | 已完成 |
| ~~P1~~ ✅ | ~~GC × cgo/指针运算交叉边界权威页补建~~（已完成：新建 LD-053 S 级六件套全配，含 GOEXPERIMENT=cgocheck2 实测发现） | 内容维护 | 已完成 |
| ~~P1~~ ✅ | ~~examples/go.sum 缺校验条目~~（已完成：go.sum 补齐 + 根模块 require + jwt 传递依赖，grpc/messaging/security build/vet 通过） | 工程 | 已完成 |
| ~~P1~~ ✅ | ~~CI 覆盖缺口~~（已完成：go.work 42/42，go-test 矩阵化 38 模块，排除清单 ci-module-coverage.md） | 工程 | 已完成 |
| ~~P2~~ ✅ | ~~docs-deploy.yml HTML 模板 3 处死引用修正~~（已完成，注意：首页其余 8 个模块卡链接同病，下期一并修） | 工程 | 已完成 |
| ~~P2~~ ✅ | ~~check-unfixed-links.ps1 重写~~（已完成：按 rescan 逻辑重写 + quality-gates.yml 三 job 挂载 + AGENTS.md §5 修正） | 工程 | 已完成 |
| ~~P2~~ ✅ | ~~docs-check.yml 去掉 markdown-link-check 的 `\|\| true` 屏蔽~~（已完成：先修掉 5 文件 15 条真死链再解除屏蔽，aliveStatusCodes 补 202） | 工程 | 已完成 |
| ~~P3~~ ✅ | ~~24 组片段级豁免 go 块逐页补全~~（已收敛：14 组补示意定义转 OK，10 组归 skip_dep，3 块加 `<!-- goblock-exempt -->` 标注；merge BAD 24→0，现 OK=33/skip_dep=45/exempt=3） | 内容维护 | 已完成 |
| ~~P3~~ ✅ | ~~AD-038 block 6 格式损坏~~（已完成；连带修复 block 5 四个未定义符号，带依赖 build 通过） | 内容维护 | 已完成 |
| ~~P3~~ ✅ | ~~FT-009 ID 复用治理~~（已完成：实质页 → FT-044，19 处引用 + indices 同步，stub 留别名入口） | 内容维护 | 已完成 |
| ~~P3~~ ✅ | ~~docs-deploy.yml 首页 8 个模块卡死引用~~（已完成：9 卡全部改指真实目录；部署复制逻辑问题移交） | 工程 | 已完成 |
| P3 | 下期月审：定义漂移逐页精读（本期仅机械核查）+ 边界域语义精读 + 版本特性→权威页精确映射表 | 内容维护 | 2026-10 月审 |
| P3 | `.kimi` 可选扩展集 4 页启用/弃用决策（表征空间/思维表征/决策树/效应系统） | 内容维护 | 待用户决定 |
| ~~P3~~ ✅ | ~~test_go_blocks 14 个 RUN-BAD 清理~~（已完成：RUN-BAD=0；11 项入豁免表 + 4 页真缺陷修复） | 内容维护 | 已完成 |

**Overall semantic health grade**: **78 / 100 → 99 / 100**（2026-09-05 三次复评：release.yml、docs-deploy 复制逻辑、AD-037、check_quality 强制门、gofmt 存量（296 清零）、race 矩阵 12 项真实缺陷根治（首跑 29/38 → 复跑 38/38）全部已解；剩余 −1：formal-verifier go.mod 待 tidy 决策）

**Sign-off**: _________________（人工复核待签）
