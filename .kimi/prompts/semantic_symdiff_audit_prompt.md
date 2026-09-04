# 可复用 Prompt：语义对称差审计

> **用途**：比较两个 Rust 版本、两个概念或两种写法之间的语义/语法/API/行为差异，输出集合论形式的对称差分析。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) §9.2。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中对以下主题进行语义对称差审计。

**对比对象 A**：{{object_a}}
**对比对象 B**：{{object_b}}
**审计维度**：{{dimensions}}  <!-- 例如：语法 / 语义 / API / 行为 / 平台支持 / 工具链 -->
**目标位置**：{{target_path}}  <!-- 已有 concept/ 页或新建页 -->

要求：
1. 先读取 A、B 在 `concept/` 中的权威页（若存在），避免重复定义。
2. 用集合记号明确写出：
   - A ∩ B：两者完全相同的部分
   - B \ A：仅 B 有的部分
   - A \ B：仅 A 有的部分
3. 每个维度用表格呈现，避免大段文字堆砌。
4. 对 B \ A 和 A \ B 中的每一项给出：
   - 具体变更/差异描述
   - 触发条件或代码模式
   - 迁移/修复建议
   - 相关 rustc error code（如适用）
5. 至少包含一个 `rust,compile_fail` 反例展示“在 A 上可行、在 B 上失败”或反之。
6. 生成 Mermaid mindmap 总结对称差。
7. 建立与相关 concept/ 页的双向链接。
8. 运行：
   python scripts/kb_auditor.py --link-check
   python scripts/check_concept_code_blocks.py --strict
9. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**object_a**: Rust 1.98.0
**object_b**: Rust 1.98.1
**dimensions**: 语法、语义、标准库 API、平台支持、Cargo/工具链
**target_path**: concept/07_future/00_version_tracking/rust_1_98_1.md
```
