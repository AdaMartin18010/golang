# Kimi 知识图谱与拓扑要求

> **EN**: Kimi Knowledge Graph and Topology Requirements
> **Summary**: Reusable contract for maintaining the knowledge graph (KG), taxonomy, relation semantics, and decision-tree mappings in this repository.
> **Scope**: `E:/_src/golang`
> **Canonical companion**: [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md) (page format), [`.kimi/kimi_semantic_requirements.md`](kimi_semantic_requirements.md) (semantic consistency).

---

## 1. 适用范围

本文件规定 Kimi 在生成/修改 `go-knowledge-base/` 内容、KG 数据、索引关系页或决策树时必须遵守的**知识图谱与拓扑规则**：

1. 五维目录与领域模型（taxonomy）。
2. KG 实体属性与语义谓词。
3. 索引/关系页（prerequisite-graph、cross-reference、complete-map）的关系语义与 KG 同步。
4. 决策树与 KG 实体 / Go 编译失败类别的映射。
5. KG 刷新与校验流水线。

---

## 2. Taxonomy 与领域模型

### 2.1 五维目录 + `indices/complete-map.md` 单一事实源

- 权威概念层固定为五维：`go-knowledge-base/01-Formal-Theory/`、`02-Language-Design/`、`03-Engineering-CloudNative/`、`04-Technology-Stack/`、`05-Application-Domains/`，页命名 `{FT|LD|EC|TS|AD}-NNN-{Kebab-Title}.md`。
- `go-knowledge-base/indices/complete-map.md` 与 `complete-index.md` 是全库主题登记的唯一权威索引；每个权威页必须能在其中找到其维度前缀、编号、`layer`（L0–L7）与主题域。
- 主题域（domain）示例：`types`、`concurrency`、`gc`、`unsafe`、`cgo`、`generics`、`formal`、`networking`、`observability`。
- 新增主题页前，先查 `complete-map.md` 是否已有同主题权威页；若新增维度内编号，须同步更新 `indices/`（by-date / by-topic / complete-index）。

### 2.2 KG 实体属性

每个 KG 实体应至少包含：

```yaml
{
  "id": "concept_channel",
  "label": { "en": "Channel", "zh": "Channel 通道" },
  "layer": "L3",
  "domain": ["concurrency", "formal"],
  "canonical_path": "go-knowledge-base/02-Language-Design/LD-002-Go-Concurrency-CSP-Formal.md"
}
```

- `layer` 与页头 Bloom 层级一致（L0 元层 … L7 未来/研究）。
- `domain` 可多个，覆盖主题交叉场景。
- `canonical_path` 指向五维权威页；预览/roadmap 特性可指向对应的版本跟踪页，但必须回链概念权威页。

---

## 3. KG 语义谓词使用规范

### 3.1 谓词目录

| 谓词 | 方向 | 使用场景 |
| --- | --- | --- |
| `dependsOn` | 单向 | A 的语义/实现依赖 B。例如 `select` dependsOn `channel`。 |
| `entails` | 单向 | A 成立则 B 必然成立。例如 `unbuffered_channel_send` entails `recv_happens_before`（Go 内存模型：无缓冲通道接收先行于发送完成）。 |
| `mutexWith` | 双向 | A 与 B 在同一上下文中互斥。例如 `data_race` mutexWith `happens_before` 保证。 |
| `refines` | 单向 | A 是 B 的更精确版本。例如 `atomic.Int64`（类型化原子操作）refines `atomic.AddInt64`（函数式 API）。 |
| `equivalentTo` | 双向 | A 与 B 在特定语境下语义等价。使用须附 `scope` 说明。例如 `sync.Once.Do` equivalentTo `double-checked locking`（scope: 单实例初始化）。 |
| `counterExample` | 单向 | A 是 B 的反例/破坏场景。例如 `racy_counter` counterExample `happens_before`。 |

### 3.2 禁止使用通用谓词

- 核心实体周边**禁止**使用 `ex:RelationAnnotation` 或无信息谓词（如裸 `relatedTo`）。
- 非核心实体若暂时无法判定具体谓词，可先用 `relatedTo` 并在注释或审计模板中说明原因，在下一轮 KG 刷新（月度语义审查）中迁移为具体谓词。

### 3.3 谓词实例化证据

- 每个具体谓词边应附带 `evidence` 字段，指向权威页中的定理编号、章节锚点或 Go 编译失败类别（Go 编译器无公开错误码体系，证据用错误类别 + 编译器原文描述，如 `declared and not used`）。
- 示例：

  ```json
  {
    "subject": "concept_select",
    "predicate": "dependsOn",
    "object": "concept_channel",
    "evidence": "go-knowledge-base/02-Language-Design/LD-002-Go-Concurrency-CSP-Formal.md#select-desugaring"
  }
  ```

---

## 4. 索引关系页语义与 KG 同步

### 4.1 关系符号映射

`go-knowledge-base/indices/`（`prerequisite-graph.md`、`cross-reference.md` 及各权威页定理链）中使用的符号必须映射到 KG 谓词：

| 符号 | 含义 | KG 谓词 |
| --- | --- | --- |
| `⟹` / `=>` | 推出 / 依赖 | `dependsOn` / `entails` |
| `⟸` / `<=` | 反推 / 被依赖 | `dependsOn`（反向） |
| `⊣` / `⊥` | 互斥 | `mutexWith` |
| `⟺` / `⇔` | 等价（特定 scope） | `equivalentTo` |
| `→` | 前置到后置学习路径 | `dependsOn` |
| `↛` | 不成立路径 | `mutexWith` / `counterExample` |

### 4.2 关系塌缩治理

- 禁止在关系页中连续使用同一符号表达不同语义（例如用 `⟹` 同时表示依赖、蕴含、学习顺序）；新增内容应保持谓词语义分布，不得让高频关系符号占比显著超过基线。
- 关系塌缩率由月度语义审查（`.kimi/templates/monthly_semantic_review.md`）人工核查并记录于 `docs/tracking/`。

### 4.3 索引 ↔ KG 双向同步

- 新增/修改 `prerequisite-graph.md`、`cross-reference.md` 或权威页定理链后，必须同步更新 `complete-map.md` 登记，确保符号关系被实例化为具体谓词边。
- 质量门与校验流水线见 §6 与 `AGENTS.md` §5。

---

## 5. 决策树与 KG / Go 编译失败类别映射

### 5.1 决策树注册表

- 决策树模板见 `.kimi/templates/decision_tree_template.md`，生成规范见 `.kimi/prompts/decision_tree_generation_prompt.md`；注册表（如启用）集中登记每棵树的元数据。
- 每棵树必须有全局唯一 `id`（如 `J-EXAMPLE-01`）、`title`、`description`、`root` 节点。

### 5.2 节点规范

```yaml
nodes:
  start:
    type: decision
    text: 错误是否为编译期类型错误？
    yes: type_check
    no: runtime_check
  type_check:
    type: action
    text: 检查接口方法集与类型嵌入是否匹配
    failure_categories: ["cannot use ... as ...", "missing method ..."]
    concepts: [concept_interface_method_set, concept_type_embedding]
```

- `type`: `decision` / `action` / `info`。
- `decision` 节点必须提供明确分支（`yes`/`no` 或 `options` 列表）。
- `action` 节点必须关联 `failure_categories`（Go 编译失败类别，如 `declared and not used`、`cannot use … as …`、`invalid operation`、`too many arguments`）和/或 `concepts`（KG 实体 ID）。

### 5.3 与 KG 实体的链接

- 每个 `action` 节点的 `concepts` 字段应指向 KG 实体 ID。
- 在权威页中引用决策树时，使用锚点链接：

  ```markdown
  对应判定路径见 [`J-TYPE-01`](../../indices/decision-trees.yaml#Lxx)。
  ```

  （链接路径以决策树注册表实际位置为准。）

### 5.4 覆盖率要求

- Top 常见 Go 编译失败类别覆盖率 ≥80%（Go 无错误码体系，类别按编译器错误文本归并统计）。
- 每棵决策树无死端（dead_end = 0）。
- 定量节点占比 ≥50%。

---

## 6. KG 刷新与校验标准流水线

新增/重命名权威页、调整定理链或索引关系页后，按以下顺序刷新与校验：

```bash
# 1. 权威页唯一性 / 同主题撞车扫描
python scripts/tmp/dup_canon_scan.py

# 2. 死链复扫（活跃区 0 死链；跳过 archive/ 与代码围栏）
python scripts/tmp/rescan_deadlinks.py

# 3. 六件套机械校验（定义/机制/实践/反例/mindmap/References）
python scripts/tmp/verify_sixpiece.py

# 4. go 代码块实测（可运行必过、编译失败反例必失败）
python scripts/tmp/extract_go_blocks.py
python scripts/tmp/test_go_blocks.py

# 5. 双向链接核查（前置/后置回链）
python scripts/tmp/check_bidir.py

# 6. 语义谓词审计（KG 谓词精度、关系塌缩）
#    按 .kimi/templates/monthly_semantic_review.md 人工核查并记录
```

> 若本次修改未涉及概念重命名或版本特性，可只做步骤 1–2、5；新增概念关系后 3–6 必须执行。版本特性变更须额外映射回概念权威页（如 1.27 特性 → `LD-037` + `examples/go127-features/`）。

---

## 7. 质量门对应关系

| 规则 | 检查方式 | 通过标准 |
| --- | --- | --- |
| KG 形态合规（编号唯一、元数据齐全） | `python scripts/tmp/annotate_bloom.py` + 人工核对 AGENTS.md §3 模板 | 维度内 FT/LD/EC/TS/AD 编号无冲突 |
| 核心谓词精度 | 月度语义审查（`.kimi/templates/monthly_semantic_review.md`） | 核心实体周边 generic_ratio = 0% |
| 拓扑质量（关系塌缩） | 月审人工核查 + `docs/tracking/` 记录 | 高频符号占比不超基线 |
| 决策树结构 | 按 §5.2/§5.4 人工校验 | dead_end = 0，Top 类别 ≥80% |
| 索引/权威页死链 | `python scripts/tmp/rescan_deadlinks.py` | 活跃区 0 死链 |

---

## 8. Kimi 生成前自检

- [ ] 新建权威页是否已在 `indices/complete-map.md` 中确定维度/编号/layer/domain？
- [ ] 前置/后置概念关系是否使用了具体 KG 谓词？
- [ ] 是否避免在核心实体周边使用 `ex:RelationAnnotation` / 裸 `relatedTo`？
- [ ] 若新增/修改索引关系页，符号（⟹/⊣/⟺）是否映射到具体谓词？
- [ ] 若新增决策树节点，是否关联了 `failure_categories` 和/或 KG 实体？
- [ ] 修改完成后是否按 §6 流水线刷新并校验？

---

## 9. 修订历史

- 2026-09: Go 版，从本仓库 `AGENTS.md` KG 语义谓词红线、五维目录约定与语义审计实践中提炼，覆盖 taxonomy、KG 谓词、索引关系语义、决策树映射与刷新校验流水线。
