# Kimi 内容生成要求：Go 分层概念知识体系

> **EN**: Kimi Content Generation Requirements for the Go Layered Concept Knowledge Base
> **Summary**: Reusable content-creation contract for Kimi when editing `go-knowledge-base/`, `docs/`, `view/`, or `examples/` in this repository.
> **Scope**: `E:/_src/golang` and all subdirectories.
> **Canonical companion**: [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) (project-level governance) and [`AGENTS.md`](../AGENTS.md) (repository rules).

---

## 1. 何时使用本要求

当你（Kimi）被要求：

- 新建或补全 `go-knowledge-base/01..05-*` 五维权威页（FT/LD/EC/TS/AD 编号）；
- 在 `docs/` / `view/` / `examples/` 写指南、版本分析或工程页；
- 针对 Go patch release（如 1.27.1）创建补丁跟踪页并注入双向链接；
- 对网络/开源库/新惯用法进行语义对齐与回填；

必须先阅读本文件，再阅读主题相关的五维权威页，最后才生成内容。

---

## 2. 元数据：每页头部必须项

### 2.1 权威页模板（`go-knowledge-base/0X-维度/`）

```markdown
# LD-NNN: 中文标题 (English Title)

> **维度**: Language Design
> **级别**: S (16+ KB) / A / B   <!-- 按实测大小：S >16KB，A >8KB，其余 B -->
> **标签**: #tag1 #go127
> **Go 版本**: 1.27+              <!-- 以 Go 1.27 工具链验证为基线 -->
> **Bloom 层级**: Lx               <!-- L0 元层 … L7 未来/研究 -->
> **前置概念**: `[..](../LD-0xx.md)` · **后置概念**: `[..]`
> **定理链**: Input → Operation → Output / Invariant
```

### 2.2 非权威页模板（`docs/` / `view/` / `examples/`）

```markdown
# 中文标题

> **维度**: {维度} | **Go 版本**: 1.27+
> **权威来源**: 通用概念解释见五维权威页 `{FT|LD|EC|TS|AD}-NNN-....md`（相对路径示例：`../go-knowledge-base/0X-维度/XX-NNN-....md`，替换为真实链接）。
> 本文仅保留应用场景、决策、操作步骤与链接，不重复概念推导。
```

### 2.3 关键字段解释

| 字段 | 要求 |
| --- | --- |
| `维度` | FT（Formal-Theory）/ LD（Language-Design）/ EC（Engineering-CloudNative）/ TS（Technology-Stack）/ AD（Application-Domains），编号在维度内唯一。 |
| `级别` | 按实测大小：S >16 KB、A >8 KB、其余 B；stub 一律标 B。 |
| `Go 版本` | 以 Go 1.27 工具链验证为基线；版本特性页标注具体版本区间。 |
| `Bloom 层级` | L0（元层/框架）/ L1（基础语法与语义）/ L2（类型系统、接口、组合、错误处理）/ L3（泛型、反射、并发原语 channel/select/context）/ L4（调度器 GMP、GC、内存模型、cgo/汇编、运行时）/ L5（跨语言/范式对比、工程架构）/ L6（生态、设计模式、分布式系统设计）/ L7（未来特性、研究、提案预览）。 |
| `前置概念` | 至少一个低层链接（高层 L6/L7 页必须含 L4 或以下链接）。禁止循环自引用。 |
| `后置概念` | 自然延伸，可含 quiz、版本页、形式化页、示例链接。 |
| `定理链` | 必须存在且贴合本页主题；禁止模板化，需对应正文推理。 |

---

## 3. 正文结构：章节六件套（每页必备）

### 3.1 权威定义

- 一句话给出**精确语义定义**，而非比喻。
- 紧接着列出**核心约束/不变式**，用表格或项目符号。

### 3.2 核心机制

- 分小节讲原理，每节配 **1 个可运行 ```go 代码块**。
- 关键推理步骤使用 `⟹`（推出）和 `⟸`（反推）标记。
- 涉及并发顺序（happens-before）、`unsafe`/cgo 边界、GC 与内存可见性、接口方法集等主题时必须有**形式化语义提示**（即使不展开证明，也要指出其依赖的公理/规则，如 Go 内存模型的 channel 同步规则）。

### 3.3 工程实践

- 使用场景、权衡、最佳实践、常用标准库/知名第三方库引用。
- 必须引用 **P0 官方**、**P1 学术/形式化**、**P2 生态/社区** 至少各一个来源（详见 [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) §4）。

### 3.4 反命题与边界分析

- 节标题固定为 **「反命题与边界分析」** 或 **「X. 反命题与边界分析」**。
- 至少包含一个 ```go 编译失败反例块（首行 `// 编译失败: <原因>`），并解释**失败原因**（引用编译器错误类别文本，如 `declared and not used`、`cannot use … (type X) as type Y`、`invalid operation`）与**正确写法**。
- 使用表格列出命题、真假值、说明。

### 3.5 思维导图

- 使用 Mermaid `mindmap`。
- 覆盖：定义、机制、边界、对比/权衡、迁移/实践。
- 禁止纯文本堆砌，节点需有层次。

### 3.6 参考来源 / References

```markdown
## References

- **P0 官方**: [The Go Programming Language Specification](https://go.dev/ref/spec) · [Go 1.27 Release Notes](https://go.dev/doc/go1.27) · [Proposal 仓库](https://go.googlesource.com/proposal)
- **P1 学术**: [Go 内存模型](https://go.dev/ref/mem) · [Dijkstra / Lamport 等经典文献] · [论文 DOI](...)
- **P2 生态**: [Go Blog](https://go.dev/blog) · [知名 Go 项目](https://github.com/etcd-io/raft)
```

### 3.7 思维表征方式

除上述章节外，每页还应根据主题选择并正确使用以下表征方式：

- **思维导图**：Mermaid `mindmap`，覆盖定义/机制/边界/实践/关联（详见 `.kimi/kimi_thinking_representation_requirements.md` §2）。
- **多维矩阵对比表**：用于概念/版本/写法对比（详见 §3 与本文件 §6）。
- **概念五元组**：定义-属性-关系-示例-反例。
- **决策树**：Mermaid flowchart，用于编译失败类别诊断与迁移判定（Go 无错误码体系，改为「编译失败类别映射」；属可选扩展机制）。
- **语义关联 / KG 关系**：使用具体谓词 `dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`（详见 §6 与 `kimi_kg_topology_requirements.md`）。
- **故障树 / 边界扩展树**：用于根因分析与边界突破场景。
- **定理推理链**：使用 `⟹`/`⟸` 与编号定理（详见 §10）。
- **表征空间分析**：能/不能/痛苦表达三维边界、等价表达谱系、跨语言对比矩阵（属可选扩展机制，详见 [`.kimi/kimi_semantic_space_analysis_requirements.md`](./kimi_semantic_space_analysis_requirements.md)）。

---

## 4. 代码块规范：10 桶分类

`go-knowledge-base/` 五维目录中所有 ```go 代码块必须可归入以下类别之一，并接受 §8.2 的实测检查：

| 类别 | Markdown 标记 | 用途 |
| --- | --- | --- |
| 可运行示例 | `go` | 标准库自足（`package main` + `func main()` 或 `TestXxx`），所在模块 `GOWORK=off go vet` 通过。 |
| 故意失败反例 | `go` + 首行 `// 编译失败: <确定性编译期错误原因>` | 必须确实编译失败，并注明期望的错误类别文本（Go 编译器无公开错误码体系）与正确写法。 |
| 依赖外部模块 | `go` + 首行 `// 不可编译: 依赖外部模块 <path>` | 展示第三方库 API，不在正文编译；必须在 `examples/<module>/` 提供 go.mod 与可运行版本。 |
| 需要 main 包装 | `go` 内嵌 `func main()` | 所有独立示例必须包含入口函数或 `TestXxx`。 |
| 片段/伪代码 | `text` | 仅用于说明算法或流程，不编译。 |
| 交互式输出 | `bash` / `text` | 命令行、编译器输出、日志。 |
| Mermaid 图 | `mermaid` | 思维导图、流程图、状态图。 |
| 版本/环境标记 | 元数据 `> **Go 版本**` + 正文说明 build tags | Go 的版本要求写在头部元数据，平台/实验特性要求用 build constraints 或 GOEXPERIMENT 说明。 |
| 测试用例 | `go` 内嵌 `func TestXxx` | 表驱动测试等验证性质的代码。 |
| 反模式归档 | `go` + 首行 `// 不可编译: 历史反模式` | 历史上存在但已不推荐的写法，明确标注。 |

**硬性要求**：

- 每个新增权威页至少包含 **1 个可运行 `go` 块** 和 **1 个编译失败反例块**。
- 代码块必须自包含，避免依赖未声明的变量/函数。
- 编译失败反例块必须说明"为何失败"（引用编译器错误类别文本）与"如何修复"。
- 在仓库根跑任何 go 命令必须加 `GOWORK=off`（根 `go.work` 不覆盖全部 42 个模块），或进入具体模块目录。

---

## 5. 链接与交叉引用规范

### 5.1 双向链接

- 版本特性页 ↔ 相关五维概念权威页（红线：版本特性必须映射回概念权威页）。
- quiz 页 ↔ 权威页。
- `examples/goXXX-features/` 示例 ↔ 权威页。
- 工程页/指南 ↔ 权威页。

### 5.2 链接格式

- 同一目录内使用相对路径 `(NNN-Kebab-Title.md)` 或 `(./NNN-Kebab-Title.md)`。
- 跨目录使用 `../../xx/xx/xx.md` 或 `../../../xx/xx/xx.md`。
- 禁止绝对路径或 URL 指向本地仓库内的 markdown 文件。
- 外部链接必须是 HTTPS，优先官方（P0）/学术（P1）域名。

### 5.3 死链检查

新增/修改后必须运行：

```bash
powershell scripts/check-unfixed-links.ps1
```

（口径：活跃区 0 死链；`archive/` 只读不修并注明；扫描跳过代码围栏与 http(s) 外链。）

---

## 6. 语义深度与推理要求

### 6.1 定理链

- 每页元数据中的 `定理链` 必须对应正文中的推理。
- 正文推理使用 `⟹` / `⟸` 标记，禁止无意义的模板化三段论。
- 涉及并发/GC/cgo/unsafe 时，必须引用具体规则（Go 内存模型、调度器文档）或论文。

### 6.2 对称差分析

当比较两个版本、两个概念或两种写法时，使用**集合对称差**视角：

```text
A = X 的语义/语法/API/行为
B = Y 的语义/语法/API/行为
A ∩ B：交集
B \ A：仅 Y 有
A \ B：仅 X 有
```

典型场景：Go 1.26 ↔ 1.27 特性差异、泛型写法 ↔ 反射写法、`sync.Mutex` ↔ channel 同步。

### 6.3 反命题表

| 命题 | 真假 | 说明 |
|---|---|---|
| ... | ✅/❌ | 引用代码或定理 |

---

## 7. 空父章节与导航回填

### 7.1 禁止空壳父章节

- 如果一个维度子目录（如 `go-knowledge-base/03-Engineering-CloudNative/01_concurrency/`）包含子文件，其 README 或 `00_` 导览页必须：
  1. 说明本目录主题；
  2. 列出所有子文件并给出 one-line 摘要；
  3. 链接到前置/后置目录。

### 7.2 indices 与索引同步

- 新建五维权威页后，必须同步登记 `go-knowledge-base/indices/`（by-date / by-topic / complete-index）。
- 跨领域主题更新 `go-knowledge-base/indices/` 下的交叉引用矩阵。
- 学习路径更新见 `go-knowledge-base/learning-paths/`。

---

## 8. Kimi 调用格式与自检

### 8.1 生成内容前的自检问题

对每个新建/修改的页面，生成前自问：

- [ ] 该主题在五维目录与 `indices/complete-index` 是否已存在权威页？若存在，只建 stub 或补充，不重复。
- [ ] 双语标题 + 元数据是否齐全（维度/级别/标签/Go 版本/Bloom/定理链）？
- [ ] 前置/后置概念链接是否有效且含低层链接？
- [ ] 是否包含可运行 `go` 块与 `// 编译失败:` 反例？
- [ ] 是否有 Mermaid mindmap？
- [ ] References 是否覆盖 P0/P1/P2？
- [ ] 是否已登记 indices（by-date / by-topic / complete-index）？
- [ ] 是否形成至少一对双向链接？

### 8.2 生成内容后的必跑命令

```bash
# 1. 死链（活跃区 0 死链；archive/ 只读跳过）
powershell scripts/check-unfixed-links.ps1

# 2. 代码块实测：抽样提取权威页 ```go 块写入临时模块验证
#    正向块：GOWORK=off go vet 通过
#    反向块：GOWORK=off go build 必须失败，且错误与首行注释一致

# 3. 文档格式
powershell scripts/check-markdown-format.ps1

# 4. 综合质量检查
powershell scripts/check_quality.ps1

# 5. 全部门（预提交钩子已内置：gofmt → go vet → golangci-lint → 相关包测试）
bash .githooks/pre-commit
```

### 8.3 禁止声明

- 禁止说"已完成"或"全部通过"，除非引用上述检查的退出码（exit 0）。
- 禁止在 `tmp/`、构建产物目录中直接写内容。
- 禁止复制已有权威页的正文到非权威位置。

---

## 9. 版本补丁页专用要求

当响应 Go patch release（如 1.27.1）时，新建/更新 `docs/tracking/`（或对应维度下）的版本补丁页必须包含：

### 9.1 头部字段

```markdown
# Go 1.27.1 稳定补丁

> **Go 版本**: **1.27.1 stable**（YYYY-MM-DD）
> **Bloom 层级**: L2-L3
> **权威来源**: 本文件为版本补丁跟踪页；语义解释回链五维权威页。
> **前置概念**: [Go 版本跟踪](../...) · [Go 1.27.0 稳定特性](../...)
> **后置概念**: [Go 1.28 前沿特性预览](../...)
```

### 9.2 正文必备章节

1. **补丁要点**：发布日期、修复类型（编译器/runtime/标准库/安全）、影响范围、是否建议升级。
2. **对称差分析**：用集合记号明确 $A \cap B$、$B \setminus A$、$A \setminus B$（相对上一版本）。
3. **技术细节**：触发路径、运行时表现、最小风险代码模式（可运行或说明性）。
4. **迁移建议**：工具链获取方式（`go install golang.org/dl/go1.27.1@latest` 或包管理器）、go 指令是否变更、CI/安全关键项目注意事项。
5. **反命题与边界分析**：至少 6 条命题，标明真假与理由。
6. **双向链接**：链接到受影响的概念权威页（GC、调度器、channel、标准库包等），并确保这些页反向链接回本补丁页。

### 9.3 go 指令（最低版本）规则

- patch release **不提升** `go.mod` 的 `go` 指令（Go 的最低版本口径：**go 指令即最低版本**）。
- 文档中可写 `1.27.0+` 或 `1.27.1 stable`，但最低版本事实源保持各模块 `go.mod` 的统一 `go 1.27`。
- 必须确认无模块在 patch 响应中被意外提升 go 指令（`grep -r "^go " --include=go.mod .` 全仓一致性核对）。

---

## 10. 形式化与定理链要求

### 10.1 定理链格式

元数据中的 `定理链` 必须对应正文中的具体推理，禁止空泛模板：

```markdown
> **定理链**: T-081 [Tier 2] channel send happens-before receive → T-082 [Tier 2] 同步点前写入可见 → T-083 [Tier 3] 生产者-消费者无需额外锁
```

正文引用：

```markdown
**T-081** 对无缓冲 channel 的 send happens-before 对应 receive 完成（Go 内存模型）。
⟹ **T-082** goroutine A 在 send 之前写入的变量，对 receive 完成之后的 goroutine B 可见。
⟹ **T-083** 以 channel 作为唯一同步点的生产者-消费者结构，不需要额外互斥锁。
```

### 10.2 形式化内容最小要求

L4-L5 页（并发、GC、cgo/unsafe、内存模型、运行时、跨语言形式化对比）必须至少包含以下之一：

- happens-before 规则表（channel、sync 原语、go 语句的启动顺序）。
- 不变式（invariant）表格（如 GC 三色抽象不变式、互斥锁临界区保护不变式）。
- 与形式化来源（Go 内存模型文档、调度器/GC 设计文档、经典并发文献）的显式对齐。

### 10.3 禁止

- 禁止把"定义 → 示例 → 结论"包装成伪定理链。
- 禁止引用不存在的定理编号。

---

## 11. 空父章节回填要求

### 11.1 触发条件

当目录满足以下任一条件时，其父章节（README 或 `00_` 导览页）必须非空：

- 目录下存在 ≥2 个子 markdown 文件。
- 目录被 `go-knowledge-base/indices/` 或学习路径索引引用。
- 目录名暗示它是一个主题域（如 `03-Engineering-CloudNative/` 下的 `01_concurrency/`）。

### 11.2 导览页最小内容

```markdown
# 目录中文名

> **EN**: English Title
> **摘要**: 本目录涵盖 ...

## 子主题

| 序号 | 文件 | 内容摘要 | Bloom 层级 |
|---|---|---|---|
| 01 | `LD-041-....md` | 一句话摘要 | L3 |
| 02 | `LD-042-....md` | 一句话摘要 | L4 |

## 前置/后置目录

- 前置：`../go-knowledge-base/0X-维度/README.md`（占位，替换为真实相对路径）
- 后置：`../go-knowledge-base/0Y-维度/README.md`（占位，替换为真实相对路径）
```

### 11.3 回填步骤

1. 定位空章节：遍历五维目录，找 README/`00_` 为空或被索引引用但无摘要的目录。
2. 按上表补全导览页。
3. 更新 `go-knowledge-base/indices/`（by-date / by-topic / complete-index）。
4. 运行 `powershell scripts/check-unfixed-links.ps1` 验证链接。

---

## 12. 代码块 10 桶与实测检查的对应关系

| 桶 # | Markdown 标记 | 检查方式 | 要求 |
| --- | --- | --- | --- |
| 1 | `go`（含 `func main()`） | 抽样提取至临时模块 `GOWORK=off go vet` | 必须通过 |
| 2 | `go` + `// 编译失败:` | 抽样提取 `GOWORK=off go build` | 必须编译失败，错误类别与注释一致 |
| 3 | `go` + `// 不可编译: 依赖外部模块` | 对应 `examples/<module>/` 可运行版本 | go.mod 齐全，`go vet ./...` 通过 |
| 4 | `go` + `// 不可编译: <原因>` | 正文必须解释为何忽略 | 平台相关/伪代码需说明 |
| 5 | 版本要求写元数据 + build tags | 人工核对 | `go` 指令 ≥ 声称的最低版本 |
| 6 | `go` 内嵌 `func TestXxx` | 所在模块 `GOWORK=off go test` | 必须通过 |
| 7 | `text` | 不编译 | 仅说明算法/流程 |
| 8 | `bash` | 命令需真实可执行 | 仓库根 go 命令必须带 `GOWORK=off` |
| 9 | `mermaid` | 语法检查 | mindmap 需有层次 |
| 10 | `go` + `// 不可编译: 历史反模式` | 标注明确 | 反模式归档 |

**执行顺序**：先跑 §8.2 的死链与格式门；新增代码块先抽样 `go vet` 实测，再视影响面跑相关模块的 `GOWORK=off go test ./...`。

---

## 13. 与现有模板的协作关系

| 文件 | 层级 | 用途 |
| --- | --- | --- |
| [`AGENTS.md`](../AGENTS.md) | 仓库规则 | 人+Agent 都必须遵守的硬性约束、质量门、红线。 |
| [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) | 项目级 | 创建"类似 golang 的新项目"时可复用的整体架构与质量门清单。 |
| [`.kimi/kimi_thinking_representation_requirements.md`](./kimi_thinking_representation_requirements.md) | 表征级 | mindmap、矩阵、五元组、决策树、KG 关系、故障树、定理链的格式与质量要求。 |
| [`.kimi/templates/concept_page_template.md`](templates/concept_page_template.md) | 页级 | 新建五维权威页时的 copy-paste 骨架。 |
| [`.kimi/kimi_quality_gate_checklist.md`](./kimi_quality_gate_checklist.md) | 操作级 | 按场景列出必须跑的质量门命令。 |
| `.kimi/kimi_content_requirements.md`（本文件） | 内容级 | Kimi 在**本仓库**内生成/修改内容时必须遵循的格式、语义、链接、代码块细则。 |
| [`.kimi/kimi_semantic_requirements.md`](./kimi_semantic_requirements.md) | 语义级 | 语义一致性、交叉域覆盖、前沿提案跟踪页、观察门纪律。 |
| [`.kimi/kimi_kg_topology_requirements.md`](./kimi_kg_topology_requirements.md) | 拓扑级 | KG 谓词、taxonomy、atlas 关系、决策树映射与刷新流水线。 |
| [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) | 审计模板 | 季度/月度语义审计报告骨架。 |
| [`.kimi/templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md) | 页模板 | 前沿提案跟踪页骨架（可选扩展：Go 无效应系统对应提案，已改造为 Go 前沿提案跟踪）。 |
| [`.kimi/templates/kimi_cross_domain_concept_template.md`](templates/kimi_cross_domain_concept_template.md) | 页模板 | 交叉/边界语义域权威页骨架（如 cgo+并发、unsafe+GC）。 |
| [`.kimi/templates/mindmap_template.md`](templates/mindmap_template.md) | 页模板 | Mermaid mindmap 骨架。 |
| [`.kimi/templates/concept_attribute_matrix_template.md`](templates/concept_attribute_matrix_template.md) | 页模板 | 概念属性矩阵骨架。 |
| [`.kimi/templates/decision_tree_template.md`](templates/decision_tree_template.md) | 页模板 | 编译失败类别映射决策树骨架（可选扩展）。 |
| [`.kimi/prompts/content_generation_prompt.md`](prompts/content_generation_prompt.md) | Prompt | 生成/补全五维权威页。 |
| [`.kimi/prompts/semantic_audit_prompt.md`](prompts/semantic_audit_prompt.md) | Prompt | 让 Kimi 执行语义审计。 |
| [`.kimi/prompts/kg_predicate_instantiation_prompt.md`](prompts/kg_predicate_instantiation_prompt.md) | Prompt | 让 Kimi 将 KG 通用关系实例化为语义谓词。 |
| [`.kimi/prompts/effect_system_page_prompt.md`](prompts/effect_system_page_prompt.md) | Prompt | 让 Kimi 生成前沿提案跟踪页（可选扩展）。 |
| [`.kimi/prompts/cross_domain_concept_prompt.md`](prompts/cross_domain_concept_prompt.md) | Prompt | 让 Kimi 生成交叉/边界语义域页。 |
| [`.kimi/prompts/semantic_drift_review_prompt.md`](prompts/semantic_drift_review_prompt.md) | Prompt | 让 Kimi 对比权威来源评审语义漂移。 |
| [`.kimi/prompts/mindmap_generation_prompt.md`](prompts/mindmap_generation_prompt.md) | Prompt | 让 Kimi 生成 Mermaid mindmap。 |
| [`.kimi/kimi_semantic_space_analysis_requirements.md`](./kimi_semantic_space_analysis_requirements.md) | 语义边界级 | 表征空间分析（可选扩展）。 |
| [`.kimi/templates/kimi_semantic_space_page_template.md`](templates/kimi_semantic_space_page_template.md) | 页模板 | 表征空间总论页骨架（可选扩展）。 |
| [`.kimi/prompts/semantic_space_analysis_prompt.md`](prompts/semantic_space_analysis_prompt.md) | Prompt | 让 Kimi 生成/审计表征空间页（可选扩展）。 |

---

## 14. 语义层要求索引

生成或修改内容时，除本文件外，还必须在以下场景阅读对应语义层资产：

| 场景 | 必读资产 |
| --- | --- |
| 新建/修改五维权威页，涉及关键术语定义 | [`.kimi/kimi_semantic_requirements.md`](./kimi_semantic_requirements.md) §2 |
| 新建/修改权威页，涉及 KG 实体/关系 | [`.kimi/kimi_kg_topology_requirements.md`](./kimi_kg_topology_requirements.md) §3–§6 |
| 新建 atlas 页面或编译失败类别决策树（可选扩展） | [`.kimi/kimi_kg_topology_requirements.md`](./kimi_kg_topology_requirements.md) §4–§5 |
| 新建 Go 前沿提案 / 预览特性跟踪页（可选扩展） | [`.kimi/templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md) + [`.kimi/kimi_semantic_requirements.md`](./kimi_semantic_requirements.md) §5 |
| 新建 cgo+并发 / unsafe+GC 等交叉域页 | [`.kimi/templates/kimi_cross_domain_concept_template.md`](templates/kimi_cross_domain_concept_template.md) + [`.kimi/kimi_semantic_requirements.md`](./kimi_semantic_requirements.md) §4 |
| 声明质量门"全部通过" | [`.kimi/kimi_semantic_requirements.md`](./kimi_semantic_requirements.md) §6.3 |
| 运行季度/月度语义审计 | [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) + [`.kimi/prompts/semantic_audit_prompt.md`](prompts/semantic_audit_prompt.md) |
| 新建/改写表征空间总论，或补充表征空间映射标注（可选扩展） | [`.kimi/kimi_semantic_space_analysis_requirements.md`](./kimi_semantic_space_analysis_requirements.md) + [`.kimi/templates/kimi_semantic_space_page_template.md`](templates/kimi_semantic_space_page_template.md) + [`.kimi/prompts/semantic_space_analysis_prompt.md`](prompts/semantic_space_analysis_prompt.md) |

---

## 15. 修订历史

- 2026-09-04: 初版治理模板，从 AGENTS.md、质量门脚本与近期 patch 响应实践中提炼。
- 2026-09-04 (C6): 新增 §9 版本补丁页、§10 形式化与定理链、§11 空父章节回填、§12 代码块 10 桶与检查对应关系；调整 §13 协作关系表。
- 2026-09-05 (T1): 新增 §3.7 思维表征方式，并在 §13 协作关系表中引入 `kimi_thinking_representation_requirements.md`。
- 2026-09-05 (S1): 新增 §14 语义层要求索引；扩展 §13 协作关系表，引入语义/拓扑要求文件及模板与 prompt。
- 2026-09-06 (SS1): 新增 §3.7 表征空间分析条目；在 §13、§14 引入表征空间要求文件、对应模板与 prompt。
- 2026-09-XX: 从 rust-lang 项目治理模板全量 Go 化：作用域映射到五维 `go-knowledge-base/`（FT/LD/EC/TS/AD 编号）；代码块改为 ```go + `// 编译失败:`/`// 不可编译:` 首行注释口径（删除错误码与版本方言标记体系）；构建命令改为 `GOWORK=off go build ./... && go vet ./... && go test ./...`；权威源分级映射到 go.dev/pkg.go.dev/proposal/Go Blog；版本基线 Go 1.27.1；检查脚本替换为本仓库真实存在的 `scripts/check-unfixed-links.ps1`、`scripts/check-markdown-format.ps1`、`scripts/check_quality.ps1` 与 `.githooks/pre-commit`。
