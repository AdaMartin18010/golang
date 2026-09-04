# 可复用 Prompt：生成/补全 Rust 概念权威页

> **用途**：当你需要 Kimi 在 `rust-lang` 知识库中新建或补全通用 `concept/` 页、工程指南页、摘要 stub 时使用。
> **配套文件**：执行前必须已读取 [`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) 与 [`AGENTS.md`](../../AGENTS.md)。
> **专用 prompt**：版本补丁、语义对称差审计、形式化/定理链、代码示例、quiz 生成请分别使用同目录下的专用 prompt。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中，针对以下主题进行内容生成/补全。

**主题**：{{topic}}
**目标层级**：{{bloom_level}}
**目标位置**：{{target_path}}
**类型**：{{type}}  <!-- 权威页 / 补丁跟踪页 / 工程指南 / 摘要 stub -->
**前置已知**：{{context}}

要求：
1. 先搜索 `concept/` 中是否已有相关权威页；若存在，只补充不重复，并建立双向链接。
2. 按 `.kimi/templates/kimi_project_requirements.md` 的元数据模板书写头部（EN、Summary、Bloom、权威来源、前置/后置概念、定理链）。
3. 正文必须包含：权威定义、核心机制（配 rust 代码块）、工程实践、反命题与边界分析（含 rust,compile_fail）、Mermaid mindmap、P0/P1/P2 References。
4. 代码块遵循 10 桶规则：可运行示例用 `rust`，反例用 `rust,compile_fail`，外部 crate 用 `rust,ignore`。
5. 所有本地 markdown 链接使用相对路径；新增页面后同步更新 `concept/SUMMARY.md` 与相关导航索引。
6. 如果主题涉及 Rust 版本差异，使用对称差视角（A∩B、B\\A、A\\B）分析。
7. 生成后运行：
   python scripts/kb_auditor.py --link-check
   python scripts/check_concept_code_blocks.py --strict
   python scripts/check_concept_authority_coverage.py --strict --include-crates
   python scripts/check_metadata_consistency.py --strict
8. 不要声明“已完成”或“全部通过”，除非质量门脚本返回 exit 0。
```

---

## 示例填充

```text
**主题**：Rust 1.98.1 patch 响应与 vtable miscompilation 修复
**目标层级**：L2-L3
**目标位置**：concept/07_future/00_version_tracking/rust_1_98_1.md
**类型**：补丁跟踪页
**前置已知**：Rust 1.98.0 存在 rustc vtable 生成错误编译，可导致 dyn Trait 方法槽位置零。
```

---

## 注意

- 本 prompt 不会自动运行质量门，生成后必须手动执行 §1 中列出的命令。
- 若 `detect_content_overlap.py` 报告重复，必须按 canonical 规则合并，禁止双权威页。
