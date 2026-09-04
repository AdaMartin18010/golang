# 可复用 Prompt：生成 Go 交叉 / 边界语义域权威页

> **用途**：为两个或多个核心概念交叉形成的边界语义域（如 channel × context、cgo × GC、泛型 × 反射）新建 `go-knowledge-base/` 权威页。
> **配套文件**：执行前必须已读取 [`AGENTS.md`](../../AGENTS.md) 与 [`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)；页面骨架使用 `.kimi/templates/kimi_cross_domain_concept_template.md`。
> **归集**：核心启用集。

---

## Prompt 模板

```text
你是 Go 语义深度文档作者。请为以下交叉域在 E:/_src/golang 仓库中新建权威页。

**交叉域名称**：{{name}}
**参与概念**：{{domain A}}, {{domain B}}, {{optional domain C}}
**判定目标**：{{decision_goal}}   <!-- 如：诊断 channel × context 混用导致的数据竞争与泄漏 -->
**相关权威页路径**：{{paths}}      <!-- 各参与概念的五维页相对路径 -->

要求：
1. 先搜索 go-knowledge-base/ 与 indices/complete-map.md 确认是否已有同主题权威页；若存在，只补充不重复。
2. 页面头部按 AGENTS.md §3 元数据模板：维度 / 级别 / 标签 / Go 版本 / Bloom 层级（交叉域通常 L4–L5）/ 前置/后置概念（必须包含所有参与概念）/ 定理链。
3. 正文必须包含：
   - 一句话权威定义与组合不变式
   - 涉及概念矩阵（≥4 个维度）
   - 核心机制与组合规则（使用 `⟹` / `⟸` 标记关键推理）
   - 形式化契约或不变式表格
   - 反命题与边界分析（≥6 条命题）
   - 至少一个 ```go 编译失败反例，首行注释 `// 编译失败: <确定性编译期错误原因>`（如 "declared and not used"、"cannot use … as …"、"invalid operation"），并给出正确写法
   - 决策树 / 判定表（如启用），并链接到 `go-knowledge-base/indices/knowledge_topology/decision_trees.yaml`
   - Mermaid mindmap（定义/机制/边界/实践/关联）
   - References：P0（go.dev/ref/spec、pkg.go.dev、proposal）+ P1（学术）+ P2（生态）至少各一
4. 确保与所有参与概念页形成双向链接（前置/后置概念 + 「相关概念」节）。
5. 若涉及 KG 实体，生成语义谓词关系（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`），禁止使用通用 RelationAnnotation。
6. 代码块自包含；可运行示例用 GOWORK=off go vet 验证；反例确实编译失败。
7. 生成后同步登记 indices（by-date / by-topic / complete-index），并运行：
   python scripts/tmp/rescan_deadlinks.py   # 死链 0
8. 不要声明"已完成"或"全部通过"，除非检查返回 exit 0。

## 输出格式

- 输出完整 markdown 文件内容（可直接写入 `go-knowledge-base/0X-维度/`）。
- 使用中文；EN 标题与 Summary 必须英文。
- stdlib-only 示例必须可在 go 1.27 下直接 `GOWORK=off go vet` / `go test` 通过。
```

---

## 约束

- 禁止将参与概念的通用解释复制到本页；只保留交叉后产生的新语义（组合不变式）。
- 每个编译失败反例必须解释失败原因（引用编译器错误文本）与正确写法。
- 必须引用至少一个 P1 形式化/学术来源（如 Go 内存模型、调度器论文、GC 相关论文）。
- 交叉域页同样遵守 canonical 红线：每个主题只有一个权威页。

---

## 示例填充

```text
**name**: channel × context 组合边界
**domain A**: channel happens-before
**domain B**: context.Context 生命周期
**decision_goal**: 诊断 select-ctx.Done() 混用导致的泄漏与数据竞争
**paths**: go-knowledge-base/02-Language-Design/LD-0xx-Channel-Happens-Before.md, go-knowledge-base/03-Engineering-CloudNative/EC-0xx-Context-Lifecycle.md
```
