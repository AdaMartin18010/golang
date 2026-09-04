# 可复用 Prompt：形式化与定理链补全

> **用途**：为 L4-L5 概念权威页补充或修正定理链、形式化语义、不变式与学术来源对齐。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) §10。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中补全/修正以下概念页的形式化与定理链。

**目标页面**：{{target_path}}
**主题**：{{topic}}
**当前 Bloom 层级**：{{bloom_level}}
**相关形式化来源**：{{formal_sources}}  <!-- 如 RustBelt、Stacked Borrows、Tree Borrows、MiniRust -->

要求：
1. 先读取目标页与所有前置/后置概念页，确保定理链不自相矛盾、不循环。
2. 在元数据中添加或修正 `定理链` 字段，格式如：
   T-081 [Tier 2] 前提 → T-082 [Tier 2] 不变式 → T-083 [Tier 3] 结论
3. 在正文中为每个定理编号写出至少一句话解释，并用 `⟹` / `⟸` 标记推理方向。
4. 补充以下至少一种形式化内容：
   - 小步/大步操作语义规则
   - 霍尔逻辑/分离逻辑断言
   - 不变式表格
   - UB 触发条件表
5. 反命题与边界分析节必须引用至少一个定理编号。
6. 所有形式化声明必须对齐 P0/P1 权威来源（Rust Reference、Rustonomicon、论文）。
7. 不要引入无法验证的“伪定理”。
8. 运行：
   python scripts/check_metadata_consistency.py --strict
   python scripts/kb_auditor.py --link-check
   python scripts/check_concept_code_blocks.py --strict
9. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**target_path**: concept/03_advanced/01_async/08_pin_unpin.md
**topic**: Pin/Unpin 不动性保证
**bloom_level**: L4-L5
**formal_sources**: RustBelt (Jung et al. 2018), std::pin module docs, RFC 2349
```
