# 可复用 Prompt：为 concept 页生成 Mermaid mindmap

> **用途**：根据已有 `concept/` 页内容，生成或优化 Mermaid `mindmap`。
> **配套文件**：[`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) §2。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中为以下 concept 页生成一个 Mermaid mindmap。

**目标页面**：{{target_path}}
**主题**：{{topic}}
**Bloom 层级**：{{bloom_level}}

要求：
1. 先读取目标页面内容，提取：定义、核心机制、边界/反例、工程实践、前后置概念。
2. 使用 Mermaid `mindmap` 语法，root 节点为主题名称。
3. 必须包含 5 个一级分支：定义、机制、边界、实践、关联。
4. 每个一级分支下至少 2 个二级节点；节点为名词/短语，禁止完整句子。
5. 若页面涉及版本差异或对比，可增加“对比/权衡”一级分支。
6. 生成的 mindmap 必须与正文章节结构一致。
7. 将 mindmap 插入页面合适位置（通常是第 1.5 节或第 5 节）。
8. 运行 `python scripts/check_mindmap_coverage.py --strict` 确认无回归。
9. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**target_path**: concept/03_advanced/01_async/08_pin_unpin.md
**topic**: Pin and Unpin
**bloom_level**: L4-L5
```
