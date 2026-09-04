# Kimi 内容生成要求：Rust 分层概念知识体系

> **EN**: Kimi Content Generation Requirements for the Rust Layered Concept Knowledge Base
> **Summary**: Reusable content-creation contract for Kimi when editing `concept/`, `docs/`, `content/`, or `crates/*/docs/` in this repository.
> **Scope**: `E:/_src/rust-lang` and all subdirectories.
> **Canonical companion**: [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) (project-level governance) and [`AGENTS.md`](../AGENTS.md) (repository rules).

---

## 1. 何时使用本要求

当你（Kimi）被要求：

- 新建或补全 `concept/` 权威页；
- 在 `docs/` / `content/` / `crates/*/docs/` 写摘要、指南或工程页；
- 针对版本更新（如 Rust 1.98.1 patch）创建补丁跟踪页并注入双向链接；
- 对网络/开源库/新惯用法进行语义对齐与回填；

必须先阅读本文件，再阅读主题相关的 `concept/` 权威页，最后才生成内容。

---

## 2. 元数据：每页头部必须项

### 2.1 权威页模板（`concept/`）

```markdown
# 中文标题

> **EN**: English Title
> **Summary**: One-sentence English abstract that states the *semantic guarantee* or *core problem* addressed.
>
> **Rust 版本**: 1.98.0+ (Edition 2024)
> **Bloom 层级**: Lx
> **权威来源**: 本文件为 `concept/` 权威页。
> **受众**: [初学者 / 进阶 / 专家 / 研究者]
> **内容分级**: [入门级 / 进阶级 / 专家级 / 综述级]
> **A/S/P 标记**: **S+A+P** — Structure + Application + Procedure
> **双维定位**: C×Ana   <!-- Concept × Analysis/Application/... -->
> **前置概念**: [概念A](../../xx/xx/xx.md) · [概念B](../../yy/yy/yy.md)
> **后置概念**: [概念C](../zz/zz.md) · [概念D](NN_file.md)
> **定理链**: 前提 → 不变式 → 结论 / 见「X. 反命题与边界分析」节
```

### 2.2 摘要/工程页模板（`docs/` / `content/` / `crates/*/docs/`）

```markdown
# 中文标题

> **EN**: English Title
> **Summary**: One-sentence English abstract scoped to this guide/use-case.
>
> **权威来源**: 通用 Rust 概念解释见 [`concept/xxx/xxx.md`](../../concept/xxx/xxx.md)。
> 本文仅保留应用场景、决策树、操作步骤与链接，不重复概念推导。
```

### 2.3 关键字段解释

| 字段 | 要求 |
|---|---|
| `EN` | 标题必须是英文，首字母大写（介词/连词小写），不使用缩写。 |
| `Summary` | 一句完整英文，说明**语义保证**或**核心问题**，避免空泛形容词。 |
| `Bloom 层级` | L0（元）/ L1（基础语法）/ L2（类型与控制）/ L3（泛型/trait/生命周期）/ L4（async/unsafe/并发/FFI）/ L5（跨语言/范式）/ L6（生态/设计模式/算法/系统）/ L7（未来/研究）。 |
| `权威来源` | `concept/` 页声明为权威；非 `concept/` 页必须给出 canonical 链接。 |
| `前置概念` | 至少一个低层链接（L5 页必须含 L4 或以下链接）。禁止循环自引用。 |
| `后置概念` | 自然延伸，可含 quiz、版本页、形式化页、crate 示例。 |
| `定理链` | 必须存在；禁止模板化，需贴合本页具体主题。 |

---

## 3. 正文结构：每页必备章节

### 3.1 权威定义

- 一句话给出**精确语义定义**，而非比喻。
- 紧接着列出**核心约束/不变式**，用表格或项目符号。

### 3.2 核心机制

- 分小节讲原理，每节配 **1 个可运行 `rust` 代码块**。
- 关键推理步骤使用 `⟹`（推出）和 `⟸`（反推）标记。
- 涉及生命周期、unsafe 边界、并发顺序、Pin 契约等必须有**形式化语义提示**（即使不展开证明，也要指出其依赖的公理/定理）。

### 3.3 工程实践

- 使用场景、权衡、最佳实践、常见 crate/库引用。
- 必须引用 **P0 官方**、**P1 学术/形式化**、**P2 生态/社区** 至少各一个来源。

### 3.4 反命题与边界分析

- 节标题固定为 **「反命题与边界分析」** 或 **「X. 反命题与边界分析」**。
- 至少包含一个 `rust,compile_fail` 块，并解释**失败原因**与**正确写法**。
- 使用表格列出命题、真假值、说明。

### 3.5 思维导图

- 使用 Mermaid `mindmap`。
- 覆盖：定义、机制、边界、对比/权衡、迁移/实践。
- 禁止纯文本堆砌，节点需有层次。

### 3.6 参考来源 / References

```markdown
## 参考来源 / References

- **P0 官方**: [The Rust Reference — ...](https://doc.rust-lang.org/reference/...) · [RFC XXXX](https://rust-lang.github.io/rfcs/XXXX-...)
- **P1 学术**: [RustBelt: Securing the Foundations of Rust](https://plv.mpi-sws.org/rustbelt/popl18/) · [Paper DOI](...)
- **P2 生态**: [crate docs](https://docs.rs/...) · [Rust Blog](https://blog.rust-lang.org/...)
```

### 3.7 思维表征方式

除上述章节外，每页还应根据主题选择并正确使用以下表征方式：

- **思维导图**：Mermaid `mindmap`，覆盖定义/机制/边界/实践/关联（详见 `.kimi/kimi_thinking_representation_requirements.md` §2）。
- **多维矩阵对比表**：用于概念/版本/写法对比（详见 §3）。
- **概念五元组**：定义-属性-关系-示例-反例（详见 §4）。
- **决策树**：Mermaid flowchart / YAML，用于错误诊断与迁移判定（详见 §5）。
- **语义关联 / KG 关系**：使用具体谓词 `dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`（详见 §6）。
- **故障树 / 边界扩展树**：用于根因分析与边界突破场景（详见 §7）。
- **定理推理链**：使用 `⟹`/`⟸` 与编号定理（详见 §8）。
- **表征空间分析**：能/不能/痛苦表达三维边界、等价表达谱系、机制组合代数、跨语言对比矩阵（详见 [`.kimi/kimi_semantic_space_analysis_requirements.md`](../kimi_semantic_space_analysis_requirements.md)）。

---

## 4. 代码块规范：10 桶分类

`concept/` 中所有 ```rust 代码块必须可归入以下类别之一，并接受 `check_concept_code_blocks.py --strict` 实测：

| 类别 | Markdown 标记 | 用途 |
|---|---|---|
| 可运行示例 | `rust` | std-only，可直接 `rustc --edition 2024` 编译通过。 |
| 故意失败反例 | `rust,compile_fail` | 必须确实编译失败，并标注期望的错误码（`E0xxx`）或失败原因。 |
| 依赖外部 crate | `rust,ignore` | 展示 crate API，但不直接编译；尽量在同一 crate 的 `examples/` 中提供可运行版本。 |
| 需要 main 包装 | `rust` 内嵌 `fn main()` | 所有独立示例必须包含 `fn main()` 或 `#[test]`。 |
| 片段/伪代码 | `text` / `pseudo` | 仅用于说明算法或流程，不编译。 |
| 交互式输出 | `bash` / `text` | 命令行、编译器输出、日志。 |
| Mermaid 图 | `mermaid` | 思维导图、流程图、状态图。 |
| 版本/环境标记 | `rust,edition2024` / `rust,ignore` | 明确 edition 或平台要求。 |
| 测试用例 | `rust` 内嵌 `#[test]` | 验证性质的代码。 |
| 反模式归档 | `rust,ignore` + 标题“反模式” | 历史上存在但已不推荐的写法。 |

**硬性要求**：

- 每个新增 `concept/` 页至少包含 **1 个 `rust` 块** 和 **1 个 `rust,compile_fail` 块**。
- 代码块必须自包含，避免依赖未声明的变量/函数。
- `compile_fail` 块必须说明“为何失败”与“如何修复”。

---

## 5. 链接与交叉引用规范

### 5.1 双向链接

- 版本特性页 ↔ 相关 `concept/` 权威页。
- quiz 页 ↔ concept 页。
- crate 示例 ↔ concept 页。
- 工程页 ↔ concept 权威页。

### 5.2 链接格式

- 同一目录内使用相对路径 `(NN_file.md)` 或 `(./NN_file.md)`。
- 跨目录使用 `../../xx/xx/xx.md` 或 `../../../xx/xx/xx.md`。
- 禁止绝对路径或 URL 指向本地仓库内的 markdown 文件。
- 外部链接必须是 HTTPS，优先官方/学术域名。

### 5.3 死链检查

新增/修改后必须运行：

```bash
python scripts/kb_auditor.py --link-check
```

---

## 6. 语义深度与推理要求

### 6.1 定理链

- 每页元数据中的 `定理链` 必须对应正文中的推理。
- 正文推理使用 `⟹` / `⟸` 标记，禁止无意义的模板化三段论。
- 涉及 unsafe/并发/形式化时，必须引用具体定理或论文。

### 6.2 对称差分析

当比较两个版本、两个概念或两种写法时，使用**集合对称差**视角：

```text
A = X 的语义/语法/API/行为
B = Y 的语义/语法/API/行为
A ∩ B：交集
B \ A：仅 Y 有
A \ B：仅 X 有
```

### 6.3 反命题表

| 命题 | 真假 | 说明 |
|---|---|---|
| ... | ✅/❌ | 引用代码或定理 |

---

## 7. 空父章节与导航回填

### 7.1 禁止空壳父章节

- 如果一个目录（如 `concept/03_advanced/01_async/`）包含子文件，其 README 或 `00_` 导览页必须：
  1. 说明本目录主题；
  2. 列出所有子文件并给出 one-line 摘要；
  3. 链接到前置/后置目录。

### 7.2 SUMMARY 与索引同步

- 新建 `concept/` 文件后，检查 `concept/SUMMARY.md` 是否包含该文件。
- 更新 `concept/00_meta/04_navigation/` 下的索引（如 `03_concept_index.md`、`04_inter_layer_map.md`）。
- 跨领域主题更新 `concept/00_meta/04_navigation/01_cross_reference_matrix.md`。

---

## 8. Kimi 调用格式与自检

### 8.1 生成内容前的自检问题

对每个新建/修改的页面，生成前自问：

- [ ] 该主题在 `concept/` 是否已存在？若存在，只建 stub 或补充，不重复。
- [ ] EN 标题与 Summary 是否完整？
- [ ] Bloom 层级、A/S/P 标记、双维定位是否一致？
- [ ] 前置/后置概念链接是否有效且含低层链接？
- [ ] 是否包含 `rust` 可运行块与 `rust,compile_fail` 反例？
- [ ] 是否有 Mermaid mindmap？
- [ ] References 是否覆盖 P0/P1/P2？
- [ ] 是否形成至少一对双向链接？

### 8.2 生成内容后的必跑命令

```bash
# 1. 死链
python scripts/kb_auditor.py --link-check

# 2. 代码块（默认抽样；依赖块加 --with-deps）
python scripts/check_concept_code_blocks.py --strict

# 3. 权威覆盖（含 crates docs）
python scripts/check_concept_authority_coverage.py --strict --include-crates

# 4. 元数据一致性
python scripts/check_metadata_consistency.py --strict

# 5. 命名规范
python scripts/check_naming_convention.py --strict

# 6. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

### 8.3 禁止声明

- 禁止说“已完成”或“全部通过”，除非引用 `run_quality_gates.sh` 的退出码。
- 禁止在 `book/`、`tmp/`、独立 workspace 的构建产物目录中直接写内容。
- 禁止复制已有权威页的正文到非权威位置。

---

## 9. 版本补丁页专用要求

当响应 Rust patch release（如 1.98.1）时，新建/更新 `concept/07_future/00_version_tracking/rust_1_XX_Y.md` 必须包含：

### 9.1 头部字段

```markdown
# Rust 1.98.1 稳定补丁

> **EN**: Rust 1.98.1 Stable Patch
> **Summary**: One-sentence abstract covering fix type, symmetric difference, and upgrade guidance.
>
> **Rust 版本**: **1.98.1 stable**（YYYY-MM-DD）
> **Bloom 层级**: L2-L3
> **权威来源**: 本文件为 `concept/` 权威页（Rust 1.98.1 patch 跟踪页）。
> **前置概念**: [Rust 版本跟踪](01_rust_version_tracking.md) · [Rust 1.98.0 稳定特性](rust_1_98_stabilized.md)
> **后置概念**: [Rust 1.99+ 前沿特性预览](rust_1_99_preview.md)
```

### 9.2 正文必备章节

1. **补丁要点**：发布日期、修复类型、影响范围、是否建议升级。
2. **对称差分析**：用集合记号明确 $A \cap B$、$B \setminus A$、$A \setminus B$。
3. **技术细节**：触发路径、运行时表现、最小风险代码模式（可运行或说明性）。
4. **迁移建议**：`rustup update stable`、MSRV 是否变更、CI/安全关键项目注意事项。
5. **反命题与边界分析**：至少 6 条命题，标明真假与理由。
6. **双向链接**：链接到受影响的 concept 页（Trait、Unsafe、Memory Model、Toolchain 等），并确保这些页反向链接回本补丁页。

### 9.3 MSRV 规则

- patch release **不提升** `Cargo.toml` 的 `rust-version`。
- 文档中可写 `1.98.0+` 或 `1.98.1 stable`，但 MSRV 事实源保持 `1.98.0`。
- 必须运行 `python scripts/check_msrv_consistency.py --strict` 确认无回归。

---

## 10. 形式化与定理链要求

### 10.1 定理链格式

元数据中的 `定理链` 必须对应正文中的具体推理，禁止空泛模板：

```markdown
> **定理链**: T-081 [Tier 2] `Pin::new_unchecked` 前提 → T-082 [Tier 2] 投影保持不动性 → T-083 [Tier 3] async 状态机安全
```

正文引用：

```markdown
**T-081** `Pin::new_unchecked` 要求调用方承诺 `T: !Unpin` 或 `T` 在 `Pin<&mut T>` 存活期间不被移动。
⟹ **T-082** 对 `Pin<&mut T>` 的字段投影必须保持“若移动该字段会破坏自引用，则禁止”。
⟹ **T-083** async/await 状态机依赖 `Pin` 保证唤醒时内部自引用仍然有效。
```

### 10.2 形式化内容最小要求

L4-L5 页（async、unsafe、并发、FFI、形式化）必须至少包含以下之一：

- 操作语义规则（小步/大步）。
- 霍尔逻辑/分离逻辑断言。
- 不变式（invariant）表格。
- 与形式化来源（RustBelt、Stacked Borrows、Tree Borrows、MiniRust）的显式对齐。

### 10.3 禁止

- 禁止把“定义 → 示例 → 结论”包装成伪定理链。
- 禁止引用不存在的定理编号。

---

## 11. 空父章节回填要求

### 11.1 触发条件

当目录满足以下任一条件时，其父章节（README 或 `00_` 导览页）必须非空：

- 目录下存在 ≥2 个子 markdown 文件。
- 目录被 `concept/SUMMARY.md` 或导航索引引用。
- 目录名暗示它是一个主题域（如 `01_async/`、`02_unsafe/`）。

### 11.2 导览页最小内容

```markdown
# 目录中文名

> **EN**: English Title
> **Summary**: 本目录涵盖 ...

## 子主题

| 序号 | 文件 | 内容摘要 | Bloom 层级 |
|---|---|---|---|
| 01 | [01_xxx.md](01_xxx.md) | 一句话摘要 | L3 |
| 02 | [02_yyy.md](02_yyy.md) | 一句话摘要 | L4 |

## 前置/后置目录

- 前置：[../00_xxx](../00_xxx/README.md)
- 后置：[../02_yyy](../02_yyy/README.md)
```

### 11.3 回填步骤

1. 运行 `python scripts/audit_content_completeness.py --json tmp/completeness.json` 定位空章节。
2. 按上表补全导览页。
3. 更新 `concept/SUMMARY.md` 与 `concept/00_meta/04_navigation/` 索引。
4. 运行 `python scripts/kb_auditor.py --link-check` 验证链接。

---

## 12. 代码块 10 桶与 `check_concept_code_blocks.py` 的对应关系

| 桶 # | Markdown 标记 | 检查器分类 | 要求 |
|---|---|---|---|
| 1 | `rust`（含 `fn main()`） | `std-only runnable` | 必须可通过 `rustc --edition 2024` 编译 |
| 2 | `rust,compile_fail` | `compile_fail` | 必须确实编译失败；建议标注 `E0xxx` |
| 3 | `rust,ignore`（依赖 crate） | `dependency` | 需 `--with-deps` 实测；或在 crate examples 中提供可运行版本 |
| 4 | `rust,ignore`（平台/伪代码） | `ignored` | 必须在正文解释为何忽略 |
| 5 | `rust,edition2021` / `edition2024` | `edition-tagged` | 标签必须与 MSRV 兼容 |
| 6 | `rust` 内嵌 `#[test]` | `test` | 必须可通过 `cargo test` 风格编译 |
| 7 | `text` / `pseudo` | `pseudo` | 不编译，仅说明算法/流程 |
| 8 | `bash` / `sh` | `shell` | 命令示例 |
| 9 | `mermaid` | `diagram` | 语法需通过 mermaid-cli 检查 |
| 10 | `rust,ignore`（历史反模式） | `anti-pattern` | 明确标注为“反模式/已废弃” |

**执行顺序**：先跑 `python scripts/check_concept_code_blocks.py --strict` 默认抽样；若新增依赖块，先 `cargo build --workspace` 再跑 `--with-deps`。

---

## 13. 与现有模板的协作关系

| 文件 | 层级 | 用途 |
|---|---|---|
| [`AGENTS.md`](../AGENTS.md) | 仓库规则 | 人+Agent 都必须遵守的硬性约束、质量门、红线。 |
| [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) | 项目级 | 创建“类似 rust-lang 的新项目”时可复用的整体架构与质量门清单。 |
| [`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) | 表征级 | mindmap、矩阵、五元组、决策树、KG 关系、故障树、定理链的格式与质量要求。 |
| [`.kimi/templates/concept_page_template.md`](templates/concept_page_template.md) | 页级 | 新建 `concept/` 权威页时的 copy-paste 骨架。 |
| [`.kimi/kimi_quality_gate_checklist.md`](../kimi_quality_gate_checklist.md) | 操作级 | 按场景列出必须跑的质量门命令。 |
| `.kimi/kimi_content_requirements.md`（本文件） | 内容级 | Kimi 在**本仓库**内生成/修改内容时必须遵循的格式、语义、链接、代码块细则。 |
| [`.kimi/kimi_semantic_requirements.md`](../kimi_semantic_requirements.md) | 语义级 | 语义一致性、交叉域覆盖、效应系统页、观察门纪律。 |
| [`.kimi/kimi_kg_topology_requirements.md`](../kimi_kg_topology_requirements.md) | 拓扑级 | KG 谓词、taxonomy、atlas 关系、决策树映射与刷新流水线。 |
| [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) | 审计模板 | 季度/月度语义审计报告骨架。 |
| [`.kimi/templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md) | 页模板 | 效应系统 / 预览特性权威页骨架。 |
| [`.kimi/templates/kimi_cross_domain_concept_template.md`](templates/kimi_cross_domain_concept_template.md) | 页模板 | 交叉/边界语义域权威页骨架。 |
| [`.kimi/prompts/semantic_audit_prompt.md`](prompts/semantic_audit_prompt.md) | Prompt | 让 Kimi 执行语义审计。 |
| [`.kimi/prompts/kg_predicate_instantiation_prompt.md`](prompts/kg_predicate_instantiation_prompt.md) | Prompt | 让 Kimi 将 KG 通用关系实例化为语义谓词。 |
| [`.kimi/prompts/effect_system_page_prompt.md`](prompts/effect_system_page_prompt.md) | Prompt | 让 Kimi 生成效应系统/预览特性页。 |
| [`.kimi/prompts/cross_domain_concept_prompt.md`](prompts/cross_domain_concept_prompt.md) | Prompt | 让 Kimi 生成交叉/边界语义域页。 |
| [`.kimi/prompts/semantic_drift_review_prompt.md`](prompts/semantic_drift_review_prompt.md) | Prompt | 让 Kimi 对比权威来源评审语义漂移。 |
| [`.kimi/kimi_semantic_space_analysis_requirements.md`](../kimi_semantic_space_analysis_requirements.md) | 语义边界级 | 表征空间分析：表达力、封闭性、等价表达、跨语言对比。 |
| [`.kimi/templates/kimi_semantic_space_page_template.md`](templates/kimi_semantic_space_page_template.md) | 页模板 | 新建 `concept/00_meta/00_framework/semantic_space.md` 风格总论页骨架。 |
| [`.kimi/prompts/semantic_space_analysis_prompt.md`](prompts/semantic_space_analysis_prompt.md) | Prompt | 让 Kimi 生成/审计表征空间页与映射标注。 |

---

## 14. 语义层要求索引

生成或修改内容时，除本文件外，还必须在以下场景阅读对应语义层资产：

| 场景 | 必读资产 |
|---|---|
| 新建/修改 `concept/` 权威页，涉及关键术语定义 | [`.kimi/kimi_semantic_requirements.md`](../kimi_semantic_requirements.md) §2 |
| 新建/修改 `concept/` 页，涉及 KG 实体/关系 | [`.kimi/kimi_kg_topology_requirements.md`](../kimi_kg_topology_requirements.md) §3–§6 |
| 新建/修改 atlas 页面或决策树 | [`.kimi/kimi_kg_topology_requirements.md`](../kimi_kg_topology_requirements.md) §4–§5 |
| 新建 Rust 预览特性 / 效应系统页 | [`.kimi/templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md) + [`.kimi/kimi_semantic_requirements.md`](../kimi_semantic_requirements.md) §5 |
| 新建 async+unsafe / FFI+async 等交叉域页 | [`.kimi/templates/kimi_cross_domain_concept_template.md`](templates/kimi_cross_domain_concept_template.md) + [`.kimi/kimi_semantic_requirements.md`](../kimi_semantic_requirements.md) §4 |
| 声明质量门“全部通过” | [`.kimi/kimi_semantic_requirements.md`](../kimi_semantic_requirements.md) §6.3 |
| 运行季度/月度语义审计 | [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) + [`.kimi/prompts/semantic_audit_prompt.md`](prompts/semantic_audit_prompt.md) |
| 新建/改写 `concept/00_meta/00_framework/semantic_space.md` 风格总论，或补充表征空间映射标注 | [`.kimi/kimi_semantic_space_analysis_requirements.md`](../kimi_semantic_space_analysis_requirements.md) + [`.kimi/templates/kimi_semantic_space_page_template.md`](templates/kimi_semantic_space_page_template.md) + [`.kimi/prompts/semantic_space_analysis_prompt.md`](prompts/semantic_space_analysis_prompt.md) |

---

## 15. 修订历史

- 2026-09-04: 初版，从 AGENTS.md、质量门脚本与近期 1.98.1 patch 响应实践中提炼。
- 2026-09-04 (C6): 新增 §9 版本补丁页、§10 形式化与定理链、§11 空父章节回填、§12 代码块 10 桶与检查器对应关系；调整 §13 协作关系表。
- 2026-09-05 (T1): 新增 §3.7 思维表征方式，并在 §13 协作关系表中引入 `kimi_thinking_representation_requirements.md`。
- 2026-09-05 (S1): 新增 §14 语义层要求索引；扩展 §13 协作关系表，引入 `kimi_semantic_requirements.md`、`kimi_kg_topology_requirements.md` 及 3 模板/5 prompt。
- 2026-09-06 (SS1): 新增 §3.7 表征空间分析条目；在 §13、§14 引入 `kimi_semantic_space_analysis_requirements.md`、对应模板与 prompt。
