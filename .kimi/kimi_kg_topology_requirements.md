# Kimi 知识图谱与拓扑要求

> **EN**: Kimi Knowledge Graph and Topology Requirements
> **Summary**: Reusable contract for maintaining the knowledge graph (KG), taxonomy, atlas relation semantics, and decision-tree mappings in this repository.
> **Scope**: `E:/_src/rust-lang`
> **Canonical companion**: [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md) (page format), [`.kimi/kimi_semantic_requirements.md`](kimi_semantic_requirements.md) (semantic consistency).

---

## 1. 适用范围

本文件规定 Kimi 在生成/修改 `concept/` 内容、KG 数据、atlas 页面或决策树时必须遵守的**知识图谱与拓扑规则**：

1. Taxonomy 与机器可读领域模型。
2. KG 实体属性与语义谓词。
3. Atlas 关系语义与 KG 同步。
4. 决策树与 KG 实体 / rustc error code 的映射。
5. KG 刷新流水线。

---

## 2. Taxonomy 与领域模型

### 2.1 `taxonomy.yaml` 单一事实源

- `concept/00_meta/taxonomy.yaml` 定义全库领域分层与主题域。
- 每个 `concept/` 文件对应一个主题，必须能在 `taxonomy.yaml` 中找到其 `layer`（L0–L7）与 `domain`（如 `types`、`ownership`、`async`、`concurrency`、`unsafe`、`ffi`、`formal`）。
- 新增主题域前，先检查 `taxonomy.yaml` 是否已存在；若需新增，应同步更新 `concept/00_meta/04_navigation/` 索引。

### 2.2 KG 实体属性

每个 KG 实体应至少包含：

```yaml
{
  "id": "concept_pin",
  "label": { "en": "Pin", "zh": "Pin 类型" },
  "layer": "L4",
  "domain": ["async", "unsafe"],
  "canonical_path": "concept/03_advanced/01_async/08_pin_unpin.md"
}
```

- `layer` 与 Bloom 层级一致。
- `domain` 可多个，覆盖主题交叉场景。
- `canonical_path` 指向 `concept/` 权威页；若为预览特性，可指向 `concept/07_future/...`。

---

## 3. KG 语义谓词使用规范

### 3.1 谓词目录

| 谓词 | 方向 | 使用场景 |
|---|---|---|
| `dependsOn` | 单向 | A 的语义/实现依赖 B。例如 `Pin` dependsOn `Unpin`。 |
| `entails` | 单向 | A 成立则 B 必然成立。例如 `Send + Sync` entails `Sync`。 |
| `mutexWith` | 双向 | A 与 B 在同一上下文中互斥。例如 `unsafe` raw deref mutexWith safe borrow guarantee。 |
| `refines` | 单向 | A 是 B 的更精确版本。例如 `TreeBorrows` refines `StackedBorrows`。 |
| `equivalentTo` | 双向 | A 与 B 在特定语境下语义等价。使用须附 `scope` 说明。 |
| `counterExample` | 单向 | A 是 B 的反例/破坏场景。例如 `dangling pointer` counterExample `valid for read`。 |
| `hasPart` / `partOf` | 单向 | 组成关系。例如 `async fn` hasPart `Future` state machine。 |

### 3.2 禁止使用通用谓词

- 核心 50 实体周边**禁止**使用 `ex:RelationAnnotation`。
- 非核心实体若暂时无法判定具体谓词，可先用 `ex:relatedTo`，但必须在注释或审计模板中说明原因，并在下一轮 KG 刷新中迁移。

### 3.3 谓词实例化证据

- 每个具体谓词边应附带 `evidence` 字段，指向 `concept/` 页中的定理编号、章节锚点或 rustc error code。
- 示例：

  ```json
  {
    "subject": "concept_pin",
    "predicate": "dependsOn",
    "object": "concept_unpin",
    "evidence": "concept/03_advanced/01_async/08_pin_unpin.md#pin-contract"
  }
  ```

---

## 4. Atlas 关系语义与 KG 同步

### 4.1 Atlas 符号映射

`concept/00_meta/knowledge_topology/` 各 atlas 页面中使用的符号必须映射到 KG 谓词：

| Atlas 符号 | 含义 | KG 谓词 |
|---|---|---|
| `⟹` / `=>` | 推出 / 依赖 | `dependsOn` / `entails` |
| `⟸` / `<=` | 反推 / 被依赖 | `dependsOn`（反向） |
| `⊣` / `⊥` | 互斥 | `mutexWith` |
| `⟺` / `⇔` | 等价（特定 scope） | `equivalentTo` |
| `→` | 前置到后置学习路径 | `dependsOn` |
| `↛` | 不成立路径 | `mutexWith` / `counterExample` |

### 4.2 关系塌缩治理

- `check_topology_quality.py --strict` 监控关系塌缩率；新增 atlas 内容不得让高频关系符号占比超过基线。
- 禁止在 atlas 中连续使用同一符号表达不同语义（例如用 `⟹` 同时表示依赖、蕴含、学习顺序）。

### 4.3 Atlas ↔ KG 双向同步

- 新增/修改 atlas 页面后，必须重新生成 KG，确保 atlas 中的符号关系被实例化为具体谓词边。
- 标准刷新流水线见 `AGENTS.md` §7。

---

## 5. 决策树与 KG / rustc 映射

### 5.1 决策树注册表

- 决策树元数据存储于 `concept/00_meta/knowledge_topology/decision_trees.yaml`。
- 每棵树必须有全局唯一 `id`（如 `J-EXAMPLE-01`）、`title`、`description`、`root` 节点。

### 5.2 节点规范

```yaml
nodes:
  start:
    type: decision
    text: 错误是否涉及生命周期？
    yes: lifetime_check
    no: type_check
  lifetime_check:
    type: action
    text: 检查 lifetime elision 与显式标注
    rustc_codes: [E0106, E0597]
    concepts: [concept_lifetime_elision, concept_lifetime_annotation]
```

- `type`: `decision` / `action` / `info`。
- `decision` 节点必须提供明确分支（`yes`/`no` 或 `options` 列表）。
- `action` 节点必须关联 `rustc_codes`（格式 `E\d{4}`）和/或 `concepts`（KG 实体 ID）。

### 5.3 与 KG 实体的链接

- 每个 `action` 节点的 `concepts` 字段应指向 KG 实体 ID。
- 在 `concept/` 权威页中引用决策树时，使用锚点链接：

  ```markdown
  对应判定路径见 [`J-LIFETIME-01`](../../00_meta/knowledge_topology/decision_trees.yaml#Lxx)。
  ```

### 5.4 覆盖率要求

- Top 30 常见 rustc error code 覆盖率 ≥80%。
- 每棵决策树无死端（dead_end = 0）。
- 定量节点占比 ≥50%。

---

## 6. KG 刷新标准流水线

新增/重命名 `concept/` 页或调整 atlas 后，按以下顺序刷新 KG：

```bash
# 1. 生成/刷新索引
python scripts/generate_kg_index.py

# 2. 注入稳定版本特性实体（当前 1.98.0）
python scripts/inject_rust198_kg_features.py

# 3. 注入 beta 版本特性实体（当前 1.99 beta）
python scripts/inject_rust199_kg_features.py

# 4. 生成 KG v3
python scripts/generate_kg_v3.py

# 5. 修正版本特性关系
python scripts/patch_rust198_kg_relations.py

# 6. 实例化语义谓词
python scripts/apply_kg_semantic_predicates.py --all-batches --apply

# 7. 抽取反例/互斥边
python scripts/apply_kg_counterexample_predicates.py --apply

# 8. 回退剩余通用关系到 relatedTo
python scripts/fallback_kg_generic_to_related.py --apply

# 9. 压缩 relatedTo 为具体谓词
python scripts/compress_kg_relatedto.py --apply

# 10. 校验
python scripts/check_kg_shapes.py --strict
python scripts/check_kg_relation_precision.py --strict
```

> 若本次修改未涉及概念重命名或版本特性，可跳过步骤 2–5，但步骤 1、6–10 必须在新增概念关系后执行。

---

## 7. 质量门对应关系

| 规则 | 检查脚本 | 通过标准 |
|---|---|---|
| KG 形态合规 | `check_kg_shapes.py --strict` | K1–K7 全 0 |
| 核心谓词精度 | `check_kg_relation_precision.py --strict` | 核心 50 实体 generic_ratio = 0% |
| 拓扑质量 | `check_topology_quality.py --strict` | T1–T6 达标 |
| 决策树结构 | `check_decision_trees.py --strict` | dead_end = 0，Top30 ≥80% |
| atlas 死链 | `kb_auditor.py --link-check` | 0 死链 |

---

## 8. Kimi 生成前自检

- [ ] 新建 `concept/` 页是否已在 `taxonomy.yaml` 中确定 layer/domain？
- [ ] 前置/后置概念关系是否使用了具体 KG 谓词？
- [ ] 是否避免在核心 50 实体周边使用 `ex:RelationAnnotation`？
- [ ] 若新增/修改 atlas，符号（⟹/⊣/⟺）是否映射到具体谓词？
- [ ] 若新增决策树节点，是否关联了 `rustc_codes` 和/或 KG 实体？
- [ ] 修改完成后是否按 §6 流水线刷新 KG 并校验？

---

## 9. 修订历史

- 2026-09-05: 初版，从 `AGENTS.md` KG 刷新链、历史语义审计与拓扑整改实践中提炼，覆盖 taxonomy、KG 谓词、atlas 关系、决策树映射与刷新流水线。
