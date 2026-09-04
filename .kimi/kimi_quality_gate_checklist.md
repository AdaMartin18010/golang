# Kimi 质量门 Checklist

> **EN**: Kimi Quality Gate Checklist
> **Summary**: Scenario-based command checklist for Kimi (and human maintainers) to verify content changes before committing.
> **Scope**: `E:/_src/rust-lang`

---

## 使用说明

根据你本次修改的类型，勾选并执行对应命令。所有命令必须在仓库根目录执行。**只有 `run_quality_gates.sh` exit 0 才能推送。**

---

## 场景 A：新建/修改 `concept/` 权威页

### A1 内容自检（生成后立即）

- [ ] EN 标题与 Summary 已填写
- [ ] Bloom 层级、A/S/P 标记、双维定位一致
- [ ] 前置/后置概念链接有效且含低层链接
- [ ] 包含至少 1 个 `rust` 可运行块和 1 个 `rust,compile_fail` 反例
- [ ] 包含 Mermaid mindmap
- [ ] References 覆盖 P0/P1/P2
- [ ] 形成至少一对双向链接

### A2 必跑命令

```bash
# 1. 死链检查
python scripts/kb_auditor.py --link-check

# 2. 代码块实测（std-only 默认抽样）
python scripts/check_concept_code_blocks.py --strict

# 3. 若新增依赖 crate 的代码块，先 build workspace 再测
# cargo build --workspace
# python scripts/check_concept_code_blocks.py --strict --with-deps

# 4. 权威覆盖率
python scripts/check_concept_authority_coverage.py --strict --include-crates

# 5. 元数据一致性
python scripts/check_metadata_consistency.py --strict

# 6. 命名规范
python scripts/check_naming_convention.py --strict

# 7. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 场景 B：版本补丁响应（如 1.98.1）

### B1 内容自检

- [ ] 补丁页头部包含 EN、Summary、发布日期、Bloom 层级
- [ ] 正文包含对称差分析（A∩B、B\A、A\B）
- [ ] 反命题表 ≥6 条
- [ ] MSRV 未提升（`Cargo.toml` 仍为 `1.98.0`）
- [ ] 受影响的 concept 页已注入补丁提示并双向链接

### B2 必跑命令

```bash
# 1. MSRV 一致性
python scripts/check_msrv_consistency.py --strict

# 2. 死链检查
python scripts/kb_auditor.py --link-check

# 3. 版本语义注入双向链接
python scripts/check_version_semantic_injection.py --strict

# 4. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 场景 C：新建/修改 `docs/` / `content/` / `crates/*/docs/` 工程页

### C1 内容自检

- [ ] 头部声明 canonical 来源链接
- [ ] 正文不重复 `concept/` 已有概念推导
- [ ] 保留应用场景、决策树、操作步骤、链接

### C2 必跑命令

```bash
# 1. 死链检查（覆盖 concept/ + docs/ + content/）
python scripts/kb_auditor.py --link-check

# 2. Stub 纯净度（若文件声明为 stub/redirect）
python scripts/check_stub_purity.py --strict

# 3. 权威覆盖率
python scripts/check_concept_authority_coverage.py --strict --include-crates

# 4. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 场景 D：只改 Kimi 配置/模板/Prompt（如 C6）

### D1 内容自检

- [ ] 新文件路径符合 `.kimi/` 命名规范
- [ ] `AGENTS.md` 已同步引用新资产
- [ ] 未在 `book/`、`tmp/`、构建产物目录中新增内容

### D2 必跑命令

```bash
# 1. 死链检查（确保新增 markdown 内部链接有效）
python scripts/kb_auditor.py --link-check

# 2. 命名规范
python scripts/check_naming_convention.py --strict

# 3. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 场景 E：新增/修改 `crates/` 代码或示例

### E1 内容自检

- [ ] `Cargo.toml` 继承 workspace 元数据
- [ ] `[lints] workspace = true` 已声明
- [ ] 新增示例在对应 `crates/xxx/examples/` 中可编译

### E2 必跑命令

```bash
# 1. 构建
 cargo check --workspace

# 2. 测试
 cargo test --workspace --quiet

# 3. Clippy
 cargo clippy --workspace -- -D warnings

# 4. 游离示例编译
python scripts/check_examples_compile.py --strict

# 5. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 场景 F：新增/修改 Quiz

### F1 内容自检

- [ ] quiz 文件位于 `concept/XX_quizzes/` 目录
- [ ] 题型 ≥3 种
- [ ] 每题标注难度与 Bloom 层级
- [ ] quiz→concept 与 concept→quiz 双向链接

### F2 必跑命令

```bash
# 1. Quiz 体系一致性
python scripts/check_quiz_system.py --strict

# 2. 死链检查
python scripts/kb_auditor.py --link-check

# 3. 全部门（推送前）
bash scripts/run_quality_gates.sh
```

---

## 通用红线

- 禁止在 `book/`、`tmp/`、独立 workspace 的构建产物目录中直接写内容。
- 禁止复制已有权威页正文到非权威位置。
- 禁止未经验证声明“已完成”或“全部通过”。
- 只有 `bash scripts/run_quality_gates.sh` 输出 `All 23 quality gates passed` 且 exit 0 后，才可执行 `git push origin main`。

---

## 快速命令速查

```bash
# 一键全部门
bash scripts/run_quality_gates.sh

# 单项快速检查
python scripts/kb_auditor.py --link-check
python scripts/check_concept_code_blocks.py --strict
python scripts/check_concept_authority_coverage.py --strict --include-crates
python scripts/check_metadata_consistency.py --strict
python scripts/check_msrv_consistency.py --strict
python scripts/check_quiz_system.py --strict
python scripts/check_naming_convention.py --strict

# 语义观察门快速检查
python scripts/check_stub_purity.py --strict
python scripts/check_cross_domain_coverage.py --strict
python scripts/check_kg_relation_precision.py --strict
python scripts/check_decision_trees.py --strict
python scripts/check_version_semantic_injection.py --strict
python scripts/check_dep_centralization.py --strict
```
