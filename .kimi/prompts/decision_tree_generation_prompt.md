# 可复用 Prompt：生成/补全决策树

> **状态**: 可选扩展 — 本仓库 AGENTS.md 体系尚未启用该机制，从 rust-lang 项目适配保留以便未来复用；不进入质量门。
> **用途**：根据 Go 编译失败类别（编译器/运行时错误文本）或迁移/判定逻辑，在 `go-knowledge-base/indices/knowledge_topology/decision_trees.yaml` 中生成或补全决策树。
> **配套文件**：[`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) §5；[`.kimi/templates/decision_tree_template.md`](../templates/decision_tree_template.md)。
> **说明**：Go 编译器无公开的「错误码 → 判定节点」对应物，判定目标以「编译失败类别映射」表达；所有错误类别须在 Go 1.27.1 工具链下实测可复现。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中为以下主题生成/补全决策树。

**主题**：{{topic}}
**树 ID**：{{tree_id}}  <!-- 如 J-ASYNC-12, DF-UNSAFE-08 -->
**目标编译失败类别**：{{go_fail_categories}}  <!-- 如 "declared and not used", "cannot use ... as ...", "all goroutines are asleep - deadlock", "WARNING: DATA RACE" -->
**判定目标**：{{decision_goal}}  <!-- 如：诊断 channel / select 相关的编译错误与运行时竞争 -->

要求：
1. 读取现有 `go-knowledge-base/indices/knowledge_topology/decision_trees.yaml`，确认 tree_id 不冲突。
2. 按 `.kimi/templates/decision_tree_template.md` 的 YAML 结构书写。
3. 每个 decision 节点必须有明确的 yes/no 分支或固定选项；每个 action 节点必须关联 `go_fail_categories: [<编译失败类别描述>]`。
4. 引用的每个编译失败类别用 GOWORK=off go build ./... / go vet ./... / go build -race ./... 在 Go 1.27.1 下实测可复现，禁止凭记忆编造错误文本。
5. 树必须无死端（dead_end = 0）。
6. 定量节点（decision/action 带明确判据）占比尽量 ≥50%。
7. 在相关权威页（go-knowledge-base/0X-维度/）中嵌入 Mermaid flowchart 可视化该树，并形成双向链接。
8. 生成后验证：
   python scripts/tmp/rescan_deadlinks.py   # 死链 0
9. 不要声明"已完成"或"全部通过"，除非检查返回 exit 0。
```

---

## 示例填充

```text
**topic**: channel 泄漏与数据竞争判定
**tree_id**: J-ASYNC-12
**go_fail_categories**: "all goroutines are asleep - deadlock", "WARNING: DATA RACE", "select statement discards result of receive expression"
**decision_goal**: 诊断 channel / select 混用导致的死锁与数据竞争
```

---

## 注意

- 本 prompt 属于可选扩展，产出不进入 AGENTS.md 质量门；若对应机制未来启用，需同步更新 AGENTS.md §5。
- 编译失败类别按错误类别描述记录（Go 无 E0xxx 错误码）；文本以工具链实测为准。
- 提交信息惯例为 `update`；push 由用户决定。
