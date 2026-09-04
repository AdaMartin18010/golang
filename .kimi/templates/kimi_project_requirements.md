# Kimi 项目协作要求模板（可复用）

> **用途**：当你要用 Kimi 创建与 `rust-lang` 知识库类似的分层、可验证、可搜索项目时，可直接复用本要求。
> **来源**：从 `E:/_src/rust-lang` 的 `AGENTS.md`、质量门脚本与目录治理实践中提炼。

---

## 1. 项目定位与架构

### 1.1 核心定位

项目应为**分层、可验证、可搜索**的知识体系/代码库，目标是为每个主题维护**单一、权威、可演进**的解释来源。

- 不要堆叠文档，要维护权威来源。
- 每个概念/主题必须只有一个权威页（canonical page）。
- 新增内容前先查重；发现重复时按 canonical 规则合并或 stub 化。

### 1.2 推荐目录职责

| 目录 | 职责 | 是否可作为权威来源 |
|---|---|---|
| `concept/` | 权威概念层，每个主题的唯一深度解释 | ✅ 是 |
| `docs/` | 指南、参考、实践、研究报告 | ❌ 否；概念解释必须链接到 `concept/` |
| `content/` | 专题深度内容套件 | ⚠️ 仅当 `concept/` 未覆盖时 |
| `crates/` | 可编译代码示例与 workspace | ❌ 概念解释不能放在这里 |
| `exercises/` | 练习题与答案 | ❌ 不能替代概念解释 |
| `archive/` | 只读历史归档 | ❌ 不是权威来源 |
| `book/` | 构建产物输出目录 | ❌ 构建产物 |
| `tmp/` | 临时文件与缓存 | ❌ 临时目录 |

### 1.3 认知分层（Bloom / L0–L7）

为每个内容页标注 Bloom 层级：

- **L0**: 元层 / 框架 / 方法论
- **L1**: 基础语法与语义
- **L2**: 类型系统、控制流、函数
- **L3**: 泛型、trait、生命周期
- **L4**: 异步、并发、unsafe、FFI
- **L5**: 跨语言/范式对比、工程架构
- **L6**: 生态、设计模式、算法、系统设计
- **L7**: 未来特性、研究、预览

跨层引用规则：

- 高层（L6）页面前置概念中必须包含至少一个低层（如 L5）链接。
- 禁止循环自引用或自引用作为前置/后置概念。

---

## 2. 文件命名与格式

### 2.1 命名规范

- 使用 `snake_case` 或 `number_prefix_snake_case`。
- 目录内文件使用两位连续序号 `NN_`（从 01 起）；`00_` 保留给导览/README。
- 禁止中文文件名、空格、混合大小写（历史/豁免目录除外）。
- 禁止双前缀（如 `06_20_`）与异形前缀（如 `1_2_`）。
- 专题系列可集中同一目录并配 README 索引。

### 2.2 Markdown 文件元数据模板

每个权威页顶部必须包含：

```markdown
# 中文标题

**EN**: English Title
**Summary**: One-sentence English abstract.

> **Rust 版本**: 1.98.0+ (Edition 2024)   <!-- 或项目对应版本 -->
> **Bloom 层级**: Lx
> **权威来源**: 本文件为 `concept/` 权威页。
> **受众**: [初学者/进阶/专家/研究者]
> **内容分级**: [入门级/进阶级/专家级/综述级]
> **A/S/P 标记**: **S+A+P** — Structure + Application + Procedure
> **双维定位**: C×App   <!-- Concept × Application/Analysis/... -->
> **前置概念**: [概念A](../../xx/xx/xx.md) · [概念B](../../yy/yy/yy.md)
> **后置概念**: [概念C](../zz/zz.md) · [概念D](NN_file.md)
> **定理链**: Input → Operation → Output / Invariant
```

### 2.3 Stub / 重定向模板

非权威位置或合并后的文件必须改为 stub，正文不超过 25 行 / 2000 字节：

```markdown
# 中文标题

**EN**: English Title
**Summary**: One-sentence English abstract.

> **权威来源**: [concept/xxx/xxx.md](../../../concept/xxx/xxx.md)
> 本文件为重定向 stub：完整解释请见上述权威页。
```

---

## 3. 内容质量要求

### 3.1 章节结构（每页必备）

1. **权威定义**：一句话定义 + 核心约束。
2. **核心机制**：分小节讲原理，配代码示例。
3. **工程实践**：使用场景、权衡、最佳实践。
4. **反命题与边界分析**：至少一个 `compile_fail` 或反例。
5. **思维导图**：Mermaid `mindmap`。
6. **参考来源 / References**：P0 官方 + P1 学术 + P2 生态。

### 3.2 代码块规范

- 可运行示例使用 ```rust。
- 故意编译失败的反例使用 ```rust,compile_fail，并确保确实失败。
- 需要外部 crate 的示例使用 ```rust,ignore 或实现 std-only 版本。
- 伪代码/片段使用 ```text 或```pseudo。
- 每个新增概念页至少包含一个 `rust` 块和一个 `rust,compile_fail` 块。

### 3.3 反例要求

- 每页必须包含“反命题与边界分析”节。
- 反例要解释错误原因，而不是只贴代码。
- 反例覆盖率应达到项目设定的基线（如 ≥40%）。

### 3.4 思维导图

- 每页使用 Mermaid `mindmap`。
- 覆盖定义、机制、边界、对比/权衡。
- 避免纯文本堆砌，要有层次结构。

### 3.5 定理链与推理

- 每页在元数据中添加 `定理链`。
- 正文中使用 `⟹`（推出）和 `⟸`（反推）标记关键推理。
- 避免模板化定理链，要贴合本页主题。

---

## 4. 权威来源分级（P0 / P1 / P2）

每个内容页必须在 References 中覆盖至少一个 P0、P1、P2 来源：

| 级别 | 含义 | 典型域名 |
|---|---|---|
| **P0 官方** | 官方语言/框架文档 | `doc.rust-lang.org`, `rust-lang.github.io`, `github.com/rust-lang`, `rustc-dev-guide`, `ferrocene.dev` |
| **P1 学术/形式化** | 论文、形式化验证、经典教材 | `plv.mpi-sws.org`, `arxiv.org`, `acm.org`, `dl.acm.org`, `ieee.org`, `springer.com` |
| **P2 生态/社区** | crate、博客、知名开源项目 | `docs.rs`, `crates.io`, `blog.rust-lang.org`, `tokio.rs`, `github.com/verus-lang`, `github.com/creusot-rs` |

示例 References 节：

```markdown
## 参考来源 / References

- **P0 官方**: [The Rust Programming Language](https://doc.rust-lang.org/book/title-page.html) · [The Rust Reference](https://doc.rust-lang.org/reference/title-page.html)
- **P1 学术**: [RustBelt: Securing the Foundations of the Rust Programming Language](https://plv.mpi-sws.org/rustbelt/popl18/)
- **P2 生态**: [docs.rs](https://docs.rs) · [crates.io](https://crates.io)
```

---

## 5. 链接与交叉引用

### 5.1 死链检查

- 所有本地 markdown 链接必须有效。
- 新增页面前必须运行链接检查。
- 重定向 stub 的 canonical 链接不能失效。

### 5.2 跨层引用

- 高层页面必须引用低层权威页。
- 禁止死端页面（无出链/入链的孤立页）。
- 前置/后置概念中的相对路径必须正确。

### 5.3 双向链接

- 版本特性、quiz、crate 示例等必须形成双向链接。
- 新增 `concept/` 页后，应同步更新相关索引/SUMMARY。

---

## 6. 代码与构建规范

### 6.1 Workspace 规范（如适用）

- 所有 workspace member 继承 workspace 元数据，禁止硬编码重复值。
- 声明 `[lints] workspace = true`。
- `rust-version` 保持项目统一 MSRV。
- feature 名使用 `kebab-case`。

### 6.2 构建验证

- `cargo check --workspace` 必须通过。
- `cargo test --workspace --quiet` 必须通过。
- `cargo clippy --workspace -- -D warnings` 必须通过。

---

## 7. 质量门（Quality Gates）

### 7.1 阻断门（必须全部通过）

项目应至少包含以下阻断质量门：

1. `cargo check --workspace`
2. `cargo test --workspace --quiet`
3. `cargo clippy --workspace -- -D warnings`
4. `cargo audit --no-fetch`
5. `cargo vet --locked`
6. `mdbook build`（如使用 mdbook）
7. 死链检查：`kb_auditor.py --link-check`
8. 内容重叠检测：`detect_content_overlap.py`
9. 双语标注检查：`add_bilingual_annotations.py --mode check-only`
10. Mermaid 语法检查
11. 拓扑质量：`check_topology_quality.py --strict`
12. KG 形态检查：`check_kg_shapes.py --strict`
13. 权威页唯一性：`check_canonical_uniqueness.py --strict`
14. 概念一致性：`concept_consistency_auditor.py --strict`
15. 段落级去重：`detect_content_overlap_v2.py` + `triage_overlap.py`
16. 权威覆盖率：`check_concept_authority_coverage.py --strict --include-crates`
17. 游离示例编译：`check_examples_compile.py --strict`
18. 命名规范：`check_naming_convention.py --strict`
19. 测验体系：`check_quiz_system.py --strict`
20. 元数据一致性：`check_metadata_consistency.py --strict`
21. 概念代码块实测：`check_concept_code_blocks.py --strict`
22. 思维导图覆盖：`check_mindmap_coverage.py --strict`
23. 综合语义健康：`semantic_health.py --strict`

### 7.2 语义观察门（达标后转阻断）

- stub 纯净度：`check_stub_purity.py --strict`
- 交叉/边界语义覆盖：`check_cross_domain_coverage.py --strict`
- KG 谓词精度：`check_kg_relation_precision.py --strict`
- 决策树错误码映射：`check_decision_trees.py --strict`
- 版本语义注入：`check_version_semantic_injection.py --strict`
- 依赖集中化：`check_dep_centralization.py --strict`

### 7.3 运行方式

```bash
# 一键运行全部门
bash scripts/run_quality_gates.sh

# 关键单项检查
python scripts/kb_auditor.py --link-check
python scripts/check_concept_code_blocks.py --strict
python scripts/check_concept_authority_coverage.py --strict --include-crates
python scripts/check_canonical_uniqueness.py --strict
```

---

## 8. 新增内容工作流

1. **查重**：运行重叠检测或搜索目标主题是否已存在。
2. **定位层级**：确定 Bloom 层级与目录位置。
3. **创建权威页**：按元数据模板、章节结构、代码块规范书写。
4. **添加 References**：确保 P0/P1/P2 全覆盖。
5. **链接前置/后置**：至少包含一个低层（如 L5）链接，无死链。
6. **跑单项检查**：死链、代码块、权威覆盖率。
7. **跑完整质量门**：`run_quality_gates.sh`。
8. **提交并推送**：`git add -A && git commit -m "feat: ..." && git push origin main`。

---

## 9. 红线与禁止事项

- 不要在构建产物目录中直接修改内容。
- 不要把临时文件提交到版本控制。
- 不要在 `crates/*/docs/` 中复制通用概念解释；应链接到 `concept/`。
- 禁止未经验证的“完成”声明；必须通过机器可复核的质量门。
- Stub/redirect 文件正文不得超过 25 行 / 2000 字节。
- KG 关系必须使用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`），避免通用 `RelationAnnotation`。
- 交叉/边界语义域必须有 `concept/` 权威页。
- 版本特性必须映射回 `concept/` 权威页。

---

## 10. 推荐工具链

- **构建**：cargo, clippy, rustfmt
- **文档**：mdbook, mermaid-cli
- **审计**：自定义 Python 脚本（`scripts/` 目录）
- **安全**：cargo-audit, cargo-vet
- **预提交**：`scripts/git_hooks/pre-commit`

---

## 11. 快速检查清单

新增/修改 concept 页前自问：

- [ ] 主题是否已存在权威页？
- [ ] EN 标题与 Summary 是否填写？
- [ ] Bloom 层级是否正确？
- [ ] 前置/后置概念链接是否有效且含低层链接？
- [ ] 是否包含 `rust` 代码块与 `rust,compile_fail` 反例？
- [ ] 是否有 Mermaid mindmap？
- [ ] References 是否覆盖 P0/P1/P2？
- [ ] 是否通过 `kb_auditor.py --link-check`？
- [ ] 是否通过 `check_concept_code_blocks.py --strict`？
- [ ] 是否通过 `check_concept_authority_coverage.py --strict --include-crates`？
- [ ] 是否通过完整 `run_quality_gates.sh`？

---

> **提示**：本模板应根据具体项目调整（如语言版本、目录结构、特殊质量门）。核心原则是：**单一权威来源、机器可验证、语义可追溯**。
