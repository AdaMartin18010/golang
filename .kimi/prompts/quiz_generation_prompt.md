# 可复用 Prompt：Quiz 生成（可选扩展）

> **用途**：为五维权威页生成理解度自测题。
> **⚠️ 本项目状态**：当前仓库**未启用 quiz 体系**（无 `quiz_registry.yaml`、无 `concept/XX_quizzes/`、质量门不含 quiz 检查）。本 prompt 仅作可选扩展：生成结果需人工评审后决定是否附在权威页尾部（折叠节）或 `go-knowledge-base/learning-paths/` 下，**不进阻断质量门**。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中为以下权威页生成理解度自测题。

**目标概念页**：{{concept_path}}
**概念名称**：{{concept_name}}
**Bloom 层级**：{{bloom_level}}
**输出位置**：{{output_mode}}  <!-- 附在页尾折叠节 / learning-paths 下独立文件 -->

要求：
1. 先读取目标概念页，确保题目只考察页内已声明的内容（定义、不变式、反例）。
2. 题型 ≥3 种：
   - 单选题（考察定义）
   - 多选题（考察边界/反例）
   - 判断题（考察反命题）
   - 代码阅读题（给出 ```go 片段，问输出/编译错误）
3. 每题标注：难度（easy / medium / hard）、对应 Bloom 层级、答案解析（引用页内具体节）。
4. 在题目文件头部添加指向概念页的链接；若附在页尾，与「相关概念」节并列，不破坏六件套结构。
5. 验证：
   python scripts/tmp/rescan_deadlinks.py
6. 不要声明"已完成"，除非检查通过。
```

---

## 示例填充

```text
**concept_path**: go-knowledge-base/02-Language-Design/LD-037-Go-1.27-Generic-Methods.md
**concept_name**: Go 1.27 泛型方法
**bloom_level**: L3
**output_mode**: 附在页尾折叠节
```
