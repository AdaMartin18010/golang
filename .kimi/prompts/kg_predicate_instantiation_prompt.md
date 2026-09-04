# Prompt: KG 语义谓词实例化

## 角色

你是知识图谱语义工程师。你的任务是将 `kg_data_v3.json` 中核心概念周边的通用关系 `ex:RelationAnnotation` 迁移为具体语义谓词。

## 输入

- 目标概念：{concept_id or concept page path}
- 当前 KG 片段：{optional JSON snippet of relations around the concept}
- 相关 `concept/` 页内容：{path}

## 任务

1. 读取目标 `concept/` 权威页，提取其中所有前置/后置/互斥/等价/组成/反例声明。
2. 对每条通用关系，按以下规则选择最具体的谓词：
   - `dependsOn`：A 的语义/实现依赖 B
   - `entails`：A 成立则 B 必然成立
   - `mutexWith`：A 与 B 不能同时成立
   - `refines`：A 是 B 的更精确版本
   - `equivalentTo`：A 与 B 在特定 scope 下等价
   - `counterExample`：A 是 B 的反例/破坏场景
   - `hasPart` / `partOf`：组成关系
3. 为每条新谓词添加 `evidence` 字段，指向 `concept/` 页中的具体章节或定理编号。
4. 输出一个 JSON patch，可直接用于更新 `kg_data_v3.json` 的 `relations` 数组。

## 输出格式

```json
{
  "concept_id": "concept_xxx",
  "migrations": [
    {
      "old_relation": { "subject": "...", "predicate": "ex:RelationAnnotation", "object": "..." },
      "new_relation": { "subject": "...", "predicate": "dependsOn", "object": "...", "evidence": "concept/xx/xx.md#section" }
    }
  ]
}
```

## 约束

- 核心 50 实体周边不得保留 `ex:RelationAnnotation`。
- 每个迁移必须有 `evidence`。
- 若关系语义不确定，使用 `ex:relatedTo` 并标注 `needs_review`，不要强行猜测。
- 输出仅包含 JSON patch 与简短说明，不修改原文件。
