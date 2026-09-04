# 可复用 Prompt：生成/补全 Go 概念权威页

> **用途**：当你需要 Kimi 在 `golang` 知识库中新建或补全五维权威页、工程指南页、摘要 stub 时使用。
> **配套文件**：执行前必须已读取 [`AGENTS.md`](../../AGENTS.md) 与 [`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)。
> **专用 prompt**：版本补丁、语义对称差审计、形式化/定理链、代码示例请分别使用同目录下的专用 prompt。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中，针对以下主题进行内容生成/补全。

**主题**：{{topic}}
**目标 Bloom 层级**：{{bloom_level}}
**目标位置**：{{target_path}}   <!-- go-knowledge-base/0X-维度/ 下，{FT|LD|EC|TS|AD}-NNN-{Kebab-Title}.md -->
**类型**：{{type}}  <!-- 权威页 / 版本跟踪页 / 工程指南 / 摘要 stub -->
**前置已知**：{{context}}

要求：
1. 先搜索 go-knowledge-base/ 与 indices/complete-map.md 确认是否已有同主题权威页；若存在，只补充不重复，并建立双向链接。
2. 头部按 AGENTS.md §3 元数据模板：维度 / 级别（按实测大小 S>16KB、A>8KB、B）/ 标签 / Go 版本 / Bloom 层级 / 前置/后置概念 / 定理链。
3. 正文必须包含章节六件套：权威定义、核心机制（配可运行 ```go）、工程实践、反命题与边界（含 `// 编译失败:` 反例 + 原因）、Mermaid mindmap、P0/P1/P2 References。
4. 代码块规范：可运行示例自包含；反例确实编译失败；外部依赖示例放 examples/<module>/。
5. 所有本地 markdown 链接使用相对路径；新增权威页后同步登记 indices（by-date / by-topic / complete-index）。
6. 如果主题涉及 Go 版本差异，使用对称差视角（A∩B、B\A、A\B）分析，并映射回概念权威页（红线）。
7. 生成后验证：
   python scripts/tmp/rescan_deadlinks.py   # 死链 0
   cd <涉及模块> && GOWORK=off go vet ./...   # 代码块所在模块
8. 不要声明"已完成"或"全部通过"，除非检查返回 exit 0。
```

---

## 示例填充

```text
**主题**: Go 1.27 json/v2 迁移要点
**目标 Bloom 层级**: L3
**目标位置**: go-knowledge-base/02-Language-Design/LD-047-Go-Context-Formal.md 旁边新建 LD-0xx-Go-127-JSON-v2.md
**类型**: 版本跟踪页
**前置已知**: encoding/json/v2 在 1.27 作为实验包提供，encoding/json/v2/experimental 下。
```

---

## 注意

- 本 prompt 不会自动运行质量门，生成后必须手动执行上述检查。
- 若发现同主题多权威页，必须按 canonical 规则合并或 stub 化，禁止双权威页。
- 提交信息惯例为 `update`；push 由用户决定。
