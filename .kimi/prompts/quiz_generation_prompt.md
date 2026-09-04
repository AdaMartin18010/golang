# 可复用 Prompt：Quiz 生成

> **用途**：为 `concept/` 权威页生成或补全 quiz，并注册到 `concept/XX_quizzes/` 与 `quiz_registry.yaml`。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) §5.3、AGENTS.md §5.1 门 19。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中为以下概念页生成 quiz。

**目标概念页**：{{concept_path}}
**概念名称**：{{concept_name}}
**Bloom 层级**：{{bloom_level}}
**目标 quiz 目录**：{{quiz_dir}}  <!-- 如 concept/05_quizzes/ -->

要求：
1. 先读取目标概念页与 `quiz_registry.yaml`，确认该主题是否已注册 quiz。
2. 若已存在，补全题目；否则新建 quiz 文件（命名：`NN_quiz_{{topic_snake}}.md`）。
3. 每份 quiz 必须包含 ≥3 种题型，例如：
   - 单选题（考察定义）
   - 多选题（考察边界/反例）
   - 判断题（考察反命题）
   - 代码阅读题（给出 rust 片段，问输出/错误码）
4. 每题必须标注：
   - 难度（easy / medium / hard）
   - 对应 Bloom 层级
   - 答案解析（引用概念页具体节或定理编号）
5. 在 quiz 文件头部添加指向概念页的链接；在概念页头部或“后置概念”中添加指向 quiz 文件的链接，形成双向链接。
6. 更新 `quiz_registry.yaml`，确保注册表与文件一致。
7. 运行：
   python scripts/check_quiz_system.py --strict
   python scripts/kb_auditor.py --link-check
8. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**concept_path**: concept/03_advanced/01_async/08_pin_unpin.md
**concept_name**: Pin and Unpin
**bloom_level**: L4-L5
**quiz_dir**: concept/05_quizzes/
```
