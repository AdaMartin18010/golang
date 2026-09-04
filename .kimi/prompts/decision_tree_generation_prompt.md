# 可复用 Prompt：生成/补全决策树

> **用途**：根据 rustc error code 或迁移/判定逻辑，在 `concept/00_meta/knowledge_topology/decision_trees.yaml` 中生成或补全决策树。
> **配套文件**：[`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) §5。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中为以下主题生成/补全决策树。

**主题**：{{topic}}
**树 ID**：{{tree_id}}  <!-- 如 J-ASYNC-12, DF-UNSAFE-08 -->
**目标 error codes**：{{rustc_codes}}  <!-- 如 E0373, E0700 -->
**判定目标**：{{decision_goal}}  <!-- 如：诊断 async 相关编译错误 -->

要求：
1. 读取现有 `concept/00_meta/knowledge_topology/decision_trees.yaml`，确认 tree_id 不冲突。
2. 按 `.kimi/templates/decision_tree_template.md` 的 YAML 结构书写。
3. 每个 decision 节点必须有明确的 yes/no 分支或固定选项；每个 action 节点必须关联 `rustc_codes: [E0xxx]`。
4. 树必须无死端（dead_end = 0）。
5. 定量节点（decision/action 带明确判据）占比尽量 ≥50%。
6. 在相关 concept 页中嵌入 Mermaid flowchart 可视化该树。
7. 运行：
   python scripts/check_decision_trees.py --strict
   python scripts/kb_auditor.py --link-check
8. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**topic**: async trait object vtable miscompilation
**tree_id**: J-ASYNC-12
**rustc_codes**: E0373, E0597, E0700, E0729, E0733
**decision_goal**: 诊断 async trait object 相关编译错误与运行时崩溃
```
