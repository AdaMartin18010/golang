# Kimi 语义内容与一致性要求

> **EN**: Kimi Semantic Content and Consistency Requirements
> **Summary**: Reusable contract for semantic depth, cross-file consistency, KG/semantic-network alignment, effect-system pages, and observation-gate discipline when Kimi edits `concept/`, `docs/`, `content/`, or `crates/*/docs/` in this repository.
> **Scope**: `E:/_src/rust-lang`
> **Canonical companion**: [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md) (page-level format), [`.kimi/kimi_thinking_representation_requirements.md`](kimi_thinking_representation_requirements.md) (representations), [`.kimi/kimi_kg_topology_requirements.md`](kimi_kg_topology_requirements.md) (KG/topology).

---

## 1. 适用范围

本文件规定 Kimi 在生成本仓库内容时必须遵守的**语义层规则**：

1. 语义一致性（定义唯一、跨文件一致、glossary 对齐、MSRV 单一事实源、语义漂移审计）。
2. 知识图谱 / 语义网络（KG 谓词实例化、关系塌缩治理、交叉/边界语义域）。
3. 效应系统 / 预览特性页专用模板。
4. 语义健康与 6 个语义观察门纪律。

> 形式与结构要求见 [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md)；思维表征格式见 [`.kimi/kimi_thinking_representation_requirements.md`](kimi_thinking_representation_requirements.md)；KG/拓扑/决策树技术要求见 [`.kimi/kimi_kg_topology_requirements.md`](kimi_kg_topology_requirements.md)。

---

## 2. 语义一致性机制

### 2.1 定义唯一性（Canonical 规则）

- 每个 Rust 概念/主题在 `concept/` 中**只能有一个**权威页。
- 若发现同主题双权威页，必须按 `AGENTS.md` §3.3 合并：保留较完整版本，另一版本改为重定向 stub。
- 新建页面前必须先查重：运行 `python scripts/detect_content_overlap.py` 并搜索 `concept/` 同主题文件。

### 2.2 跨文件语义一致性审计

- 关键术语（如 `Send`/`Sync`、`Pin`、`unsafe` five superpowers、`async fn`/`Future` 等价表述、`Pin` 投影规则、variance、let chains 等）在全库必须保持一致定义。
- 修改涉及这些术语的页面后，应运行：

  ```bash
  python scripts/concept_consistency_auditor.py --strict
  ```

- 若脚本报告跨文件定义矛盾，必须将改动收敛到 `concept/` 权威页，并在其他位置改为链接或引用。

### 2.3 Glossary 对齐

- 全库 glossary 以 `concept/00_meta/01_terminology/01_terminology_glossary.md` 为权威表。
- 新建/修改概念定义时，检查权威表是否已收录该术语；若未收录，应在权威页首次出现时用 `**术语**` 加粗，并在后续季度审计中补入 glossary。
- 运行：

  ```bash
  python scripts/check_glossary_alignment.py --strict
  ```

### 2.4 MSRV 单一事实源

- 根 `Cargo.toml` 的 `rust-version = "1.98.0"` 是唯一事实源。
- patch release（如 1.98.1）**不提升** `rust-version`；文档中可写 `1.98.0+` 或 `1.98.1 stable`，但 MSRV 保持 `1.98.0`。
- 新增/修改后运行：

  ```bash
  python scripts/check_msrv_consistency.py --strict
  ```

### 2.5 语义漂移审计

- **季度国际来源抽样审计**：每季度抽样 5–8 个核心 `concept/` 页，与 The Rust Reference、The Rustonomicon、TRPL、RFC、RustBelt/Stacked Borrows/Tree Borrows 等权威来源对比，检查定义漂移。
- **月度语义深度评审**：每月检查新增/修改页是否仍满足本文件 §2–§4 要求。
- 审计输出按 [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) 填写，发现漂移时必须：
  1. 定位 `concept/` 权威页；
  2. 按权威来源修正定义或明确标注“实现定义/平台相关”；
  3. 更新 References 与版本声明。

---

## 3. 知识图谱与语义网络

### 3.1 KG 谓词目录

核心概念周边**禁止**使用通用 `ex:RelationAnnotation`，必须使用以下语义谓词之一：

| 谓词 | 含义 | 示例 |
|---|---|---|
| `dependsOn` | A 依赖 B | `Pin` dependsOn `Unpin` |
| `entails` | A 语义上蕴含 B | `Send + Sync` entails `Sync` |
| `mutexWith` | A 与 B 不能同时成立 | `unsafe` raw pointer deref mutexWith safe guarantee |
| `refines` | A 细化 B | `TreeBorrows` refines `Stacked Borrows` |
| `equivalentTo` | A 与 B 等价 | `&T` and `&mut T` immutable borrows equivalent in read-only context |
| `counterExample` | A 是 B 的反例 | `dangling pointer` counterExample `valid for read` |
| `hasPart` / `partOf` | 组成 | `async fn` hasPart `Future` state machine |

### 3.2 谓词实例化流程

- 新建/修改 `concept/` 权威页后，若涉及核心 50 实体周边关系，应运行 KG 刷新流水线（见 `AGENTS.md` §7 KG 刷新与谓词实例化）。
- 核心 50 实体周边 `generic_ratio` 必须为 0%，由 `check_kg_relation_precision.py --strict` 捕获。

### 3.3 关系塌缩治理

- atlas 层间/层内关系符号（⟹/⊣/⟺/↔）必须对应到具体 KG 谓词，禁止无差别使用 `ex:relatedTo`。
- 关系塌缩率由 `check_topology_quality.py --strict` 监控；新增内容不得引入新的高频关系塌缩。

### 3.4 Taxonomy 与机器可读领域模型

- `concept/00_meta/taxonomy.yaml` 是领域模型单一事实源。
- 每个 KG 实体应携带 `layer`（L0–L7）和 `domain`（如 concurrency、async、unsafe、ffi、types 等）属性。
- 新建 `concept/` 页时，根据其 Bloom 层级与主题域更新或校验 `taxonomy.yaml`。

---

## 4. 交叉 / 边界语义域

### 4.1 关键交叉域清单

以下主题必须在 `concept/` 中存在**非 stub 权威页**，并链接相关单点概念：

- `let chains`
- `unsafe extern blocks`
- `async + unsafe`
- `FFI + async`
- `Send / Sync boundaries`
- `Pin + lifetimes`
- `no_std async`
- `const generics + trait objects`
- `GAT + async`
- `drop × concurrency / async`
- `async cancellation safety`
- `allocator_api` 稳定边界
- `match ergonomics` / default binding mode
- `tail expression drop`
- `const trait impl`
- `unsafe extern blocks`（Edition 2024）

> 缺口由 `scripts/check_cross_domain_coverage.py --strict` 捕获。

### 4.2 何时新建独立权威页

当某个主题同时满足以下条件时，必须新建独立 `concept/` 权威页，而不是仅在相关页中分散讨论：

1. 涉及两个及以上 L2–L4 核心概念（如 async + unsafe）。
2. 存在独立的 rustc error code 或编译器判定规则。
3. 有独立的版本语义（如 Edition 2024 `unsafe extern`）。
4. 有独立的安全边界或 UB 风险。

### 4.3 交叉域页必备章节

1. **定义与涉及概念矩阵**：列出参与概念及它们之间的语义关系。
2. **形式化契约 / 不变式**：明确组合后的保证与限制。
3. **决策树 / 判定表**：链接到 `decision_trees.yaml` 中对应树，并给出 Mermaid 可视化。
4. **反命题与边界分析**：至少 6 条命题，含 `rust,compile_fail` 反例。
5. **双向链接**：与所有参与概念页互指，并确保 KG 中新增对应关系边。

---

## 5. 效应系统 / 预览特性页

### 5.1 适用场景

- Rust 预览特性（如 effect system、`gen` blocks、`throws`、`with`-clauses、`const trait impl` 新语法等）。
- 需要跟踪学术 lineage、设计空间、语法提案、版本语义注入的主题。

### 5.2 页面头部模板

```markdown
# 中文标题

> **EN**: English Title
> **Summary**: One-sentence abstract covering the design space, current syntax proposal, and stability status.
>
> **Rust 版本**: 1.99 beta / nightly / 设计提案
> **Bloom 层级**: L7
> **权威来源**: 本文件为 `concept/` 权威页（Rust 1.99+ 预览特性跟踪页）。
> **前置概念**: [Effects and Purity](../xx/xx.md) · [Async](../xx/xx.md)
> **后置概念**: [Rust 版本跟踪](rust_1_99_preview.md)
```

### 5.3 正文必备章节

1. **学术谱系 / 设计动机**：列出关键文章/RFC/论文与时间线。
2. **设计空间分类矩阵**：开放 vs 封闭、静态 vs 动态、显式 vs 隐式等维度。
3. **当前语法提案**：使用 `rust,ignore` 展示可能语法，明确标注“提案，非稳定决策”。
4. **效应代数 / 组合规则**：效果并集、排除、互斥、别名。
5. **与现有概念的交叉分析**：如 Effect × Pin、Effect × async、Effect × const。
6. **反命题与边界分析**：标注“当前不可编译”“可能变化”等。
7. **版本语义注入**：链接到受影响的 `concept/` 权威页，并确保这些页反向链接回本页。

### 5.4 版本语义注入

- 预览特性页必须在 `concept/` 相关权威页的“版本兼容性”小节中留下反向链接。
- 运行：

  ```bash
  python scripts/check_version_semantic_injection.py --strict
  ```

---

## 6. 语义健康与观察门纪律

### 6.1 质量门结构

- **23 个阻断门**：必须全部通过才能推送。
- **6 个语义观察门**：默认 `continue-on-error`，用于持续度量，达标后按规则转阻断：
  1. `check_stub_purity.py --strict`
  2. `check_cross_domain_coverage.py --strict`
  3. `check_kg_relation_precision.py --strict`
  4. `check_decision_trees.py --strict`
  5. `check_version_semantic_injection.py --strict`
  6. `check_dep_centralization.py --strict`

### 6.2 观察门转正规则

- 任一观察门必须**连续 4 周（或连续 10 次 CI 运行）达标**，且本地 `--strict` 当前 `exit 0`，才可转阻断。
- **禁止**以口头或一次性指示绕过该规则。
- 转正后若指标退化，按阻断门流程处置（修复或经评估后回调为观察门）。

### 6.3 “完成”声明纪律

- 禁止未经验证声明“已完成”“全部通过”“100%”。
- 任何完成声明必须引用可机器复核的证据：`bash scripts/run_quality_gates.sh` 输出、`reports/` 基线、CI 记录。
- 声明必须覆盖全部 23 个阻断门与 6 个语义观察门（若观察门状态变化，须一并核对）。

---

## 7. 与质量门对应关系

| 语义要求 | 检查脚本 | 通过标准 |
|---|---|---|
| 定义唯一性 / canonical | `check_canonical_uniqueness.py --strict` | 0 双权威页 |
| 跨文件一致性 | `concept_consistency_auditor.py --strict` | 0 错误级发现 |
| Glossary 对齐 | `check_glossary_alignment.py --strict` | 差异清单可解释或清零 |
| MSRV 单一事实源 | `check_msrv_consistency.py --strict` | 无不一致声明 |
| 交叉域覆盖 | `check_cross_domain_coverage.py --strict` | 16/16 覆盖 |
| KG 谓词精度 | `check_kg_relation_precision.py --strict` | 核心 generic_ratio = 0% |
| 版本语义注入 | `check_version_semantic_injection.py --strict` | 1.90–1.99 特性双向链接覆盖率达标 |
| 综合语义健康 | `semantic_health.py --strict` | grade = OK |

---

## 8. Kimi 生成前自检

对每页新增/修改内容，生成前确认：

- [ ] 该主题在 `concept/` 是否已存在权威页？若存在，只建 stub 或补充，不重复。
- [ ] 是否涉及关键术语？若涉及，检查跨文件定义一致性。
- [ ] 是否属于 16 个交叉/边界语义域之一？若是，按 §4 新建独立权威页。
- [ ] 是否涉及预览特性 / 效应系统？若是，按 §5 模板书写。
- [ ] KG 关系是否使用具体谓词（dependsOn/entails/mutexWith/refines/equivalentTo/counterExample）？
- [ ] 是否形成至少一对双向链接（版本特性↔概念、quiz↔概念、工程页↔概念）？
- [ ] 声明“完成”前是否已跑 `bash scripts/run_quality_gates.sh` 并保存输出？

---

## 9. 修订历史

- 2026-09-05: 初版，整合 `.kimi/archive/` 语义审计、KG 谓词、交叉域覆盖、效应系统计划与 `AGENTS.md` 观察门纪律。
