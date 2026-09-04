# Prompt: KG 语义谓词实例化（Go 版）

> **状态**: 可选扩展 — 本仓库 AGENTS.md 体系尚未启用该机制，从 rust-lang 项目适配保留以便未来复用；不进入质量门。

## 角色

你是知识图谱语义工程师。你的任务是将本仓库 KG 关系存储（如 `go-knowledge-base/indices/prerequisite-graph.md` 中的关系条目，或从五维权威页元数据导出的 KG JSON）中核心概念周边的通用关系 `ex:RelationAnnotation` 迁移为具体语义谓词。

## 输入

- 目标概念：`{FT|LD|EC|TS|AD}-NNN` 编号或权威页路径（`go-knowledge-base/0X-维度/FT|LD|EC|TS|AD-NNN-{Kebab-Title}.md`）
- 当前 KG 片段：`{optional JSON snippet of relations around the concept}`
- 相关权威页内容：`{path}`

## 任务

1. 读取目标权威页，提取其中所有前置/后置/互斥/等价/组成/反例声明（元数据的「前置概念 / 后置概念」字段 + 正文章节）。
2. 对每条通用关系，按以下规则选择最具体的谓词：
   - `dependsOn`：A 的语义/实现依赖 B
   - `entails`：A 成立则 B 必然成立
   - `mutexWith`：A 与 B 不能同时成立
   - `refines`：A 是 B 的更精确版本
   - `equivalentTo`：A 与 B 在特定 scope 下等价
   - `counterExample`：A 是 B 的反例/破坏场景
   - `hasPart` / `partOf`：组成关系
3. 为每条新谓词添加 `evidence` 字段，指向权威页中的具体章节、定理编号或元数据字段。
4. 输出一个 JSON patch，可直接用于更新 KG 关系存储的 `relations` 数组。

## 输出格式

```json
{
  "concept_id": "LD-037",
  "migrations": [
    {
      "old_relation": { "subject": "...", "predicate": "ex:RelationAnnotation", "object": "..." },
      "new_relation": { "subject": "...", "predicate": "dependsOn", "object": "...", "evidence": "go-knowledge-base/02-Language-Design/LD-0xx-...md#反命题与边界" }
    }
  ]
}
```

## 约束

- 核心概念周边不得保留 `ex:RelationAnnotation`。
- 每个迁移必须有 `evidence`。
- 若关系语义不确定，使用 `ex:relatedTo` 并标注 `needs_review`，不要强行猜测。
- 输出仅包含 JSON patch 与简短说明，不修改原文件。
- 迁移完成后，权威页元数据中的前置/后置概念链接保持相对路径不变，不做重写。
