# 可复用 Prompt：为 Go 权威页生成 Mermaid mindmap

> **用途**：根据已有 `go-knowledge-base/` 权威页内容，生成或优化 Mermaid `mindmap`。
> **配套文件**：执行前必须已读取 [`AGENTS.md`](../../AGENTS.md)（§3 章节六件套 / §6 红线）；`.kimi/templates/mindmap_template.md`。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中为以下权威页生成一个 Mermaid mindmap。

**目标页面**：{{target_path}}  <!-- go-knowledge-base/0X-维度/ 下的 {FT|LD|EC|TS|AD}-NNN-Kebab-Title.md -->
**主题**：{{topic}}
**Bloom 层级**：{{bloom_level}}

要求：
1. 先读取目标页面内容，提取：权威定义、核心机制、边界/编译失败反例、工程实践、前置/后置概念。
2. 使用 Mermaid `mindmap` 语法，root 节点为主题名称。
3. 必须包含 5 个一级分支：定义、机制、边界、实践、关联。
4. 每个一级分支下至少 2 个二级节点；节点为名词/短语，禁止完整句子。
5. 若页面涉及版本差异或对比（如 go 指令 1.26/1.27），可增加“对比/权衡”一级分支。
6. 生成的 mindmap 必须与正文章节六件套结构一致。
7. 将 mindmap 插入页面的「思维导图」节（六件套第 5 项）；若页面已有 mindmap，只优化不重复插入。
8. 不要改变页面元数据与编号；插入后同步检查前置/后置概念相对路径无死链。
9. 运行：
   python scripts/tmp/rescan_deadlinks.py   # 死链 0
10. 不要声明"已完成"或"全部通过"，除非检查返回 exit 0。
```

---

## 示例填充

```text
**target_path**: go-knowledge-base/02-Language-Design/LD-037-Go-1.27-Generic-Methods.md
**topic**: Go 1.27 泛型方法
**bloom_level**: L3
```

---

## 注意

- 本 prompt 不会自动运行质量门，生成后必须手动执行上述检查。
- mindmap 是权威页的必备节，缺该节的页面视为六件套不完整，应优先补齐。
- 禁止生成纯文本堆砌的 mindmap，必须有层次结构且与正文一致。
