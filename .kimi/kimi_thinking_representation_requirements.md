# Kimi 思维表征与语义内容要求

> **EN**: Kimi Thinking Representation and Semantic Content Requirements
> **Summary**: Reusable contract for the visual, structural, and semantic representations that Kimi must produce or verify when editing `concept/` and related pages.
> **Scope**: `E:/_src/rust-lang`
> **Companion**: [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md) (page-level format), [`.kimi/kimi_quality_gate_checklist.md`](kimi_quality_gate_checklist.md) (commands).

---

## 1. 适用范围

本文件规定 Kimi 在生成本仓库内容时必须使用的**思维表征方式**：

1. 思维导图（Mermaid mindmap）
2. 多维矩阵对比表
3. 概念定义-属性-关系-示例-反例 五元组
4. 决策树（Mermaid flowchart / YAML）
5. 语义关联 / KG 关系
6. 故障树 / 边界扩展树
7. 定理推理链（⟹ / ⟸）

---

## 2. 思维导图（Mermaid mindmap）

### 2.1 使用场景

- 每个 `concept/` 权威页顶部或第 1.5 节必须包含一个 mindmap。
- 用于概述概念的定义、机制、边界、实践与关联。

### 2.2 强制结构

```mermaid
mindmap
  root((主题))
    定义
      一句话语义
      核心不变式
    机制
      机制 A
      机制 B
    边界
      反例 1
      反例 2
    实践
      最佳实践
      常见陷阱
    关联
      前置概念
      后置概念
```

### 2.3 质量要求

- 节点层级 ≥3 层（root → 一级分支 → 二级分支）。
- 每个一级分支下至少 2 个二级节点。
- 禁止使用纯文本长句作为节点；节点应为名词或短语。
- 同一页的 mindmap 必须与正文章节结构一致。

### 2.4 反模式

- ❌ 只列出 2–3 个一级分支，无细节。
- ❌ 节点是完整段落而非关键词。
- ❌ mindmap 与正文内容脱节。

---

## 3. 多维矩阵对比表

### 3.1 使用场景

- 比较两个及以上概念、版本、写法、运行时、crate 时。
- 用于展示多维度下的差异与权衡。

### 3.2 强制结构

```markdown
| 维度 | 对象 A | 对象 B | 对象 C |
|---|---|---|---|
| **语义** | ... | ... | ... |
| **语法** | ... | ... | ... |
| **性能** | ... | ... | ... |
| **安全保证** | ... | ... | ... |
| **典型场景** | ... | ... | ... |
```

### 3.3 质量要求

- 至少 4 个维度；每个维度必须可验证（语义/语法/性能/安全/生态/平台）。
- 单元格内容避免空泛形容词，必须给出具体规则或代码片段。
- 若用于版本对比，必须包含“对称差”视角：交集、仅 A 有、仅 B 有。

---

## 4. 概念五元组：定义-属性-关系-示例-反例

### 4.1 使用场景

- 每个 `concept/` 权威页在“权威定义”节后，必须用表格或结构化块呈现概念五元组。

### 4.2 强制结构

```markdown
| 元素 | 内容 | 要求 |
|---|---|---|
| **定义** | 一句话精确语义 | 可转化为判定过程 |
| **属性** | 核心约束/不变式列表 | 不可撤销、不可破坏 |
| **关系** | 前置/后置/互斥/等价概念 | 使用具体语义谓词 |
| **示例** | `rust` 可运行代码 | 展示正确用法 |
| **反例** | `rust,compile_fail` 代码 | 展示错误原因与修复 |
```

### 4.3 关系谓词

必须使用具体语义谓词，禁止通用 `related to`：

- `dependsOn`：A 依赖 B
- `entails`：A 语义上蕴含 B
- `mutexWith`：A 与 B 不能同时成立
- `refines`：A 细化 B
- `equivalentTo`：A 与 B 等价
- `counterExample`：A 是 B 的反例

---

## 5. 决策树（Decision Tree）

### 5.1 使用场景

- 诊断 rustc error code 的判定路径。
- 版本迁移、兼容性、工具链选择的决策流程。
- 存放在 `concept/00_meta/knowledge_topology/decision_trees.yaml` 并在 concept 页引用。

### 5.2 YAML 结构

```yaml
trees:
  - id: J-EXAMPLE-01
    title: 示例判定树
    description: 判定某个编译错误的原因
    root: start
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
      type_check:
        type: action
        text: 检查类型推导与 trait bound
        rustc_codes: [E0277, E0308]
```

### 5.3 质量要求

- 每个决策节点必须有明确的 `yes` / `no` 分支或固定选项。
- 行动节点必须关联具体 `rustc_codes`（格式 `E\d{4}`）。
- 树必须无死端（dead_end = 0）。
- 定量节点占比 ≥50%。
- Top 30 常见 rustc error code 覆盖率 ≥80%。

### 5.4 Mermaid 可视化

在 concept 页中可用 Mermaid flowchart 展示树的简化视图：

```mermaid
flowchart TD
    A[错误是否涉及生命周期？] -- 是 --> B[检查 lifetime elision]
    A -- 否 --> C[检查类型推导]
    B --> D[E0106 / E0597]
    C --> E[E0277 / E0308]
```

---

## 6. 语义关联 / KG 关系

### 6.1 使用场景

- `concept_kb.json` / `kg/` 中的知识图谱关系。
- concept 页正文中的前置/后置/互斥/等价声明。

### 6.2 强制谓词

核心概念周边禁止使用通用 `ex:RelationAnnotation`：

| 谓词 | 含义 | 示例 |
|---|---|---|
| `dependsOn` | 依赖 | `Pin` dependsOn `Unpin` |
| `entails` | 蕴含 | `Send + Sync` entails `Sync` |
| `mutexWith` | 互斥 | `unsafe` raw pointer deref mutexWith safe guarantee |
| `refines` | 细化 | `TreeBorrows` refines `StackedBorrows` |
| `equivalentTo` | 等价 | `&T` and `&mut T` immutable borrows equivalent in read-only context |
| `counterExample` | 反例 | `dangling pointer` counterExample `valid for read` |
| `hasPart` / `partOf` | 组成 | `async fn` hasPart `Future` state machine |

### 6.3 质量要求

- 每个 `concept/` 页至少通过前置/后置链接形成 2 条具体语义关系。
- KG 中核心 50 实体周边的 `generic_ratio` 必须为 0%。

---

## 7. 故障树 / 边界扩展树

### 7.1 使用场景

- 分析某个 UB、segfault、编译错误的根因。
- 展示概念边界如何被突破。

### 7.2 强制结构

```mermaid
mindmap
  root((运行时崩溃))
    编译期保证被破坏
      vtable 零槽位
      async trait object
     unsafe 误用
      悬垂指针
      数据竞争
     FFI 边界
      ABI 不匹配
      生命周期桥接错误
```

或使用 fault_tree 表格：

```markdown
| 顶层事件 | 中间事件 | 基本事件 | 触发条件 | 预防措施 |
|---|---|---|---|---|
| segfault | vtable 零槽位 | rustc miscompilation | 复杂 where 子句 + async trait object | 升级 1.98.1 |
```

### 7.3 质量要求

- 必须区分“顶层事件”“中间事件”“基本事件”。
- 每个基本事件必须对应具体代码模式或编译器行为。
- 必须给出预防措施或修复建议。

---

## 8. 定理推理链（⟹ / ⟸）

### 8.1 使用场景

- L3-L5 概念页的核心推理。
- 形式化页（L4）的公理/定理/证明骨架。

### 8.2 强制格式

```markdown
**T-081** 前提：...
⟹ **T-082** 不变式：...
⟹ **T-083** 结论：...
```

### 8.3 质量要求

- 每个定理编号在元数据和正文中一致。
- 推理步骤必须可验证，禁止跳跃。
- 涉及 unsafe/并发/形式化时，必须引用 P1 来源。
- 反命题节必须引用至少一个定理编号。

---

## 9. 各表征方式与质量门对应

| 表征方式 | 检查脚本 | 通过标准 |
|---|---|---|
| mindmap | `check_mindmap_coverage.py --strict` | 覆盖率 ≥10%，反例节存在率 ≥40% |
| 决策树 | `check_decision_trees.py --strict` | dead_end=0，Top30 覆盖率 ≥80% |
| KG 关系 | `check_kg_relation_precision.py --strict` | 核心 generic_ratio=0% |
| 定理链 | `check_metadata_consistency.py --strict` + 人工 | 编号一致、正文有引用 |
| 多维矩阵/五元组/故障树 | `kb_auditor.py --link-check` + 人工 | 表格完整、链接有效 |

---

## 10. Kimi 生成前自检

对每页内容，生成前确认：

- [ ] 是否已包含 mindmap？
- [ ] 若涉及对比，是否使用多维矩阵？
- [ ] 概念五元组是否完整（定义/属性/关系/示例/反例）？
- [ ] 关系是否使用具体谓词（dependsOn/entails/mutexWith/refines/equivalentTo/counterExample）？
- [ ] 若涉及错误诊断，是否可转化为决策树节点？
- [ ] 定理链是否编号一致、推理可验证？

---

## 11. 表征空间分析专用表征

当主题为 Rust 表征空间 / 语义边界 / 表达力 / 跨语言对比时，除通用表征方式外，还必须包含以下专用表征：

### 11.1 三维边界表

| 类别 | 概念 | Rust 表达 | 语义保持 | 成本/痛点 | 权威来源 |
|---|---|---|---|---|---|
| 能且高效 | ... | ... | ... | ... | ... |
| 能但痛苦 | ... | ... | ... | ... | ... |
| 不能表达 | ... | ... | — | ... | ... |

要求：

- 「不能表达」必须引用官方决策或设计哲学（RFC、Reference、Release Notes）。
- 「能但痛苦」必须说明结构性痛点，而非简单抱怨语法啰嗦。

### 11.2 等价表达谱系

使用层级缩进文本或 Mermaid flowchart 展示同一语义的不同 Rust 表达，并标注：

- 观察等价
- 语义等价
- 性能等价

### 11.3 机制组合代数

将 Rust 机制抽象为算子（`Own`、`Borrow`、`Lifetime`、`Trait`、`Generic`、`Const`、`Unsafe`、`Async`），给出：

- 合法组合规则；
- 非法组合规则（配 `rust,compile_fail`）；
- 组合选择决策树。

### 11.4 跨语言对比矩阵

至少 5 个维度对比 Rust / C++ / Haskell / Go / Java（或其他目标语言），并附包含关系图。

### 11.5 反命题决策树与定理矩阵

- 至少 3 个反命题，每个用 Mermaid graph TD，反例节点红、修正节点绿。
- 至少 8 行定理一致性矩阵：断言、前提 ⟹ 结论、反例/边界、典型场景、失效条件。

详细要求见 [`.kimi/kimi_semantic_space_analysis_requirements.md`](../kimi_semantic_space_analysis_requirements.md)。

---

## 12. 语义网络与 KG 谓词表征

思维表征不仅服务于人类阅读，也必须能进入机器可消费的 KG。具体要求见 [`.kimi/kimi_kg_topology_requirements.md`](../kimi_kg_topology_requirements.md)。本节只列出与表征格式直接相关的要点：

### 12.1 Atlas 符号 → KG 谓词

| Atlas 符号 | 含义 | 对应 KG 谓词 |
|---|---|---|
| `⟹` | 推出 / 依赖 | `dependsOn` / `entails` |
| `⟸` | 被依赖 | `dependsOn`（反向） |
| `⊣` / `⊥` | 互斥 | `mutexWith` |
| `⟺` / `⇔` | 等价（特定 scope） | `equivalentTo` |
| `→` | 学习路径 / 前置 | `dependsOn` |
| `↛` | 不成立路径 | `mutexWith` / `counterExample` |

### 12.2 概念五元组中的“关系”列

概念五元组表格的“关系”列必须使用具体谓词，禁止写“相关”“有关”等空泛词。

### 12.3 决策树的语义锚定

每个 Mermaid 决策树节点应能在 `decision_trees.yaml` 中找到对应节点，并关联 KG 实体 ID 与 `rustc_codes`。

---

## 13. 效应系统与版本预览特性表征

当主题属于 Rust 预览特性（尤其是效应系统、`throws`、`with`-clauses、`gen` blocks 等）时，表征方式需额外满足：

1. **学术谱系图**：使用时间线或 mindmap 展示关键文章/RFC/论文。
2. **设计空间矩阵**：≥3 个维度对比开放/封闭、显式/隐式、静态/动态等。
3. **效果代数表**：并集、排除、互斥、别名。
4. **语法提案代码块**：使用 `rust,ignore` 并加 `⚠️ 设计提案` 声明。
5. **版本语义注入标记**：在页面末尾明确列出需反向链接的 `concept/` 权威页。

专用模板见 [`.kimi/templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md)。

---

## 14. 修订历史

- 2026-09-05: 初版，覆盖 7 类思维表征方式的格式、质量要求与质量门对应关系。
- 2026-09-05 (S1): 新增 §12 语义网络与 KG 谓词表征、§13 效应系统与版本预览特性表征。
