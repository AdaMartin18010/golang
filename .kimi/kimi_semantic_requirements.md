# Kimi 语义内容与一致性要求

> **EN**: Kimi Semantic Content and Consistency Requirements
> **Summary**: Reusable contract for semantic depth, cross-file consistency, KG/semantic-network alignment, frontier-proposal tracking pages, and observation-gate discipline when Kimi edits `go-knowledge-base/`, `docs/`, `view/`, `examples/` or生产代码 in this repository.
> **Scope**: `E:/_src/golang`
> **Canonical companion**: [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md) (page-level format), [`.kimi/kimi_thinking_representation_requirements.md`](kimi_thinking_representation_requirements.md) (representations), [`.kimi/kimi_kg_topology_requirements.md`](kimi_kg_topology_requirements.md) (KG/topology).

---

## 1. 适用范围

本文件规定 Kimi 在生成本仓库内容时必须遵守的**语义层规则**：

1. 语义一致性（定义唯一、跨文件一致、glossary 对齐、go 指令单一事实源、语义漂移审计）。
2. 知识图谱 / 语义网络（KG 谓词实例化、关系塌缩治理、交叉/边界语义域）。
3. Go 前沿提案跟踪页专用模板（含 GOEXPERIMENT 机制）。
4. 语义健康与语义观察门纪律。

> 形式与结构要求见 [`.kimi/kimi_content_requirements.md`](kimi_content_requirements.md)；思维表征格式见 [`.kimi/kimi_thinking_representation_requirements.md`](kimi_thinking_representation_requirements.md)；KG/拓扑技术要求见 [`.kimi/kimi_kg_topology_requirements.md`](kimi_kg_topology_requirements.md)。

---

## 2. 语义一致性机制

### 2.1 定义唯一性（Canonical 规则）

- 每个 Go 概念/主题在五维权威层（`go-knowledge-base/01-Formal-Theory/` … `05-Application-Domains/`）中**只能有一个**权威页。
- 若发现同主题双权威页，必须按 `AGENTS.md` §6 红线合并：保留较完整版本，另一版本改为重定向 stub（≤ 25 行 / 2000 字节，只留一句话 + canonical 链接）。
- 新建页面前必须先查重：搜索 `go-knowledge-base/indices/complete-map.md` 与五维目录确认无同主题权威页；可用 `python scripts/tmp/dup_canon_scan.py` 辅助扫描重复 canonical。

### 2.2 跨文件语义一致性审计

- 关键术语在全库必须保持一致定义，包括但不限于：
  - 并发语义：`happens-before`、`channel` 关闭/接收规则、`select` 公平性、`context` 取消与超时传播。
  - 类型系统：`interface` 方法集判定规则、类型嵌入（embedding）与方法提升、接口满足（satisfaction）判定时点。
  - 内存与 unsafe：`unsafe.Pointer` 转换规则、cgo 指针传递限制（Go 指针不得跨 C 边界持有）、`string`/`[]byte` 零拷贝转换边界。
  - 泛型：GCShape/stenciling 实例化语义、类型参数约束满足、`comparable` 与 `any`。
  - 错误处理：`errors.Is`/`errors.As` 匹配规则、`panic`/`recover` 语义边界。
- 修改涉及这些术语的页面后，应复扫相关权威页与引用页，确认定义表述一致。
- 若发现跨文件定义矛盾，必须将改动收敛到五维权威页，并在其他位置改为链接或引用。

### 2.3 Glossary 对齐

- 全库 glossary 以 `go-knowledge-base/GLOSSARY.md` 为权威表。
- 新建/修改概念定义时，检查权威表是否已收录该术语；若未收录，应在权威页首次出现时用 `**术语**` 加粗，并在后续季度审计中补入 glossary。
- 新增权威页必须同步登记 `go-knowledge-base/indices/`（by-date / by-topic / complete-index）。

### 2.4 go 指令单一事实源

- 42 个 go module 的 `go.mod` 中 `go 1.27` 指令是唯一事实源；**go 指令即最低支持版本**。
- patch release（如工具链 1.27.1）**不提升** `go` 指令；文档中可写 `go 1.27` 或注明验证工具链 `1.27.1`，但 go 指令保持 `1.27`。
- 版本升级（如 1.28）时必须全仓批量对齐所有 `go.mod`、`go.work`、CI（`.github/workflows/`）与 `deployments/docker/Dockerfile`。
- 新增/修改后可用 `grep -rn "^go " --include=go.mod .` 抽检一致性。

### 2.5 语义漂移审计

- **季度国际来源抽样审计**：每季度抽样 5–8 个核心五维权威页，与 The Go Programming Language Specification、Effective Go、Go Blog、pkg.go.dev、go.googlesource.com/proposal 等权威来源对比，检查定义漂移（模板见 [`.kimi/templates/quarterly_international_source_audit.md`](templates/quarterly_international_source_audit.md)）。
- **月度语义深度评审**：每月检查新增/修改页是否仍满足本文件 §2–§4 要求（模板见 [`.kimi/templates/monthly_semantic_review.md`](templates/monthly_semantic_review.md)）。
- 审计输出按 [`.kimi/templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) 填写，发现漂移时必须：
  1. 定位五维权威页；
  2. 按权威来源修正定义或明确标注“实现定义/平台相关/依赖具体编译器版本”；
  3. 更新 References 与 Go 版本声明。

---

## 3. 知识图谱与语义网络

### 3.1 KG 谓词目录

核心概念周边**禁止**使用通用 `ex:RelationAnnotation` 或 `ex:relatedTo`，必须使用以下语义谓词之一：

| 谓词 | 含义 | 示例 |
| --- | --- | --- |
| `dependsOn` | A 依赖 B | `select` dependsOn `channel` |
| `entails` | A 语义上蕴含 B | `close(ch)` entails `接收端立即获得零值与 ok=false` |
| `mutexWith` | A 与 B 不能同时成立 | `unsafe.Pointer` 算术 mutexWith GC 安全保证 |
| `refines` | A 细化 B | `atomic.Int64` refines `sync.Mutex + int64` |
| `equivalentTo` | A 与 B 等价 | `for range ch` equivalentTo `recv 循环直至 channel 关闭` |
| `counterExample` | A 是 B 的反例 | `data race` counterExample `happens-before` |
| `hasPart` / `partOf` | 组成 | `goroutine` hasPart `GMP 调度模型` |

### 3.2 谓词实例化流程

- 新建/修改五维权威页后，若涉及核心实体周边关系，KG 关系必须使用 §3.1 的具体谓词表达，并写入该页元数据/KG 节。
- 核心实体周边禁止出现无差别的通用关系边；谓词误用由季度语义审计复核。

### 3.3 关系塌缩治理

- 文档中凡是语义连接（依赖、蕴含、互斥、等价、细化）必须落到具体 KG 谓词，禁止用“相关”“参见”式无差别表述替代语义关系。
- 新增内容不得引入新的关系塌缩；月审/季审抽样检查塌缩率。

### 3.4 Taxonomy 与机器可读领域模型

- `go-knowledge-base/indices/`（complete-map / by-topic / prerequisite-graph 等）是领域模型与登记事实源。
- 每个 KG 实体应携带 `layer`（L0–L7，对应权威页 Bloom 层级）和 `domain`（如 concurrency、gc、unsafe-cgo、types、generics、errors 等）属性。
- 新建权威页时，根据其 Bloom 层级与主题域同步登记 `indices/`，并保持 complete-map 登记数 == 实际篇数。

---

## 4. 交叉 / 边界语义域

### 4.1 关键交叉域清单

以下主题必须在五维权威层中存在**非 stub 权威页**，并链接相关单点概念：

- `cgo + unsafe`（Go 指针跨 C 边界规则）
- `defer + 返回值 / 闭包变量捕获`
- `goroutine + panic / recover`
- `channel + select + context 取消`
- `interface 方法集 + 类型嵌入`
- `泛型 + 反射`（GCShape × `reflect.Type` 的边界）
- `sync.Pool + GC`（对象复用与回收时机）
- `内存模型 happens-before × channel / sync 原语`
- `error wrapping（errors.Is/As）× 自定义错误类型`
- `string ↔ []byte 零拷贝转换 × unsafe`
- `finalizer（runtime.SetFinalizer）× 对象复活`
- `go:generate + 构建约束（build tags）`

> 缺口由月度语义评审抽样捕获，并按 §4.2 评估是否新建独立权威页。

### 4.2 何时新建独立权威页

当某个主题同时满足以下条件时，必须新建独立五维权威页，而不是仅在相关页中分散讨论：

1. 涉及两个及以上 L2–L4 核心概念（如 cgo + unsafe）。
2. 存在独立的编译器判定规则或错误类别（如 “declared and not used”、“invalid operation”、“cannot use … as …”）。
3. 有独立的版本语义（如某 GOEXPERIMENT 在 1.27 的行为变化）。
4. 有独立的安全边界或正确性风险（data race、GC 悬挂、cgo 指针违规、内存模型违例）。

### 4.3 交叉域页必备章节

1. **定义与涉及概念矩阵**：列出参与概念及它们之间的语义关系（用 §3.1 谓词）。
2. **形式化契约 / 不变式**：明确组合后的保证与限制。
3. **判定表 / 决策矩阵**：用表格或 Mermaid 给出组合场景 → 语义结论。
4. **反命题与边界分析**：至少 6 条命题，含 `// 编译失败:` 或运行时违例反例。
5. **双向链接**：与所有参与概念页互指，并确保 KG 中新增对应关系边。

---

## 5. Go 前沿提案跟踪页

> **状态**: 可选扩展 — 本仓库 AGENTS.md 体系尚未启用该机制，从 rust-lang 项目适配保留以便未来复用；不进入质量门。

### 5.1 适用场景

- 已被接受但尚未发布的 Go 提案（go.googlesource.com/proposal 中 accepted 状态，如未来版本语言变更）。
- 需要通过 `GOEXPERIMENT` 启用的实验特性（实验包、工具链行为开关）。
- 需要跟踪学术 lineage、设计空间、语法提案、版本语义注入的主题。

### 5.2 页面头部模板

```markdown
# 中文标题

> **EN**: English Title
> **Summary**: One-sentence abstract covering the design space, current proposal status, and experiment gating.
>
> **Go 版本**: 1.27 stable / GOEXPERIMENT=xxx / 设计提案（accepted, unreleased）
> **Bloom 层级**: L7
> **权威来源**: 本文件为五维权威层跟踪页（Go 前沿提案）。
> **前置概念**: [..](../LD-0xx.md) · [..](../EC-0xx.md)
> **后置概念**: [Go 1.28 版本跟踪](go_1_28_preview.md)
```

### 5.3 正文必备章节

1. **学术谱系 / 设计动机**：列出关键提案（proposal issue）、设计文档、讨论时间线。
2. **设计空间分类矩阵**：开放 vs 封闭、静态 vs 动态、显式 vs 隐式等维度。
3. **当前提案语法 / API**：使用 ```go 展示可能的语法或 API，首行注释 `// 不可编译: <原因>` 或标注“提案，非稳定决策”。
4. **组合规则**：与现有机制的交互（如与泛型、与 arena、与 json/v2 实验包的组合约束）。
5. **与现有概念的交叉分析**：如 Proposal × 内存模型、Proposal × 泛型、Proposal × cgo。
6. **反命题与边界分析**：标注“当前不可用”“可能变化”等。
7. **版本语义注入**：链接到受影响的五维权威页，并确保这些页反向链接回本页。

### 5.4 版本语义注入

- 前沿提案页必须在相关权威页的“版本兼容性”小节中留下反向链接。
- 稳定版本特性必须映射回概念权威页（红线：如 1.27 特性 → `LD-037` + `examples/go127-features/`）。

---

## 6. 语义健康与观察门纪律

### 6.1 质量门结构

- **阻断门**：以 [`.kimi/templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) §7.1 与 `AGENTS.md` §5 为准（gofmt、`GOWORK=off go vet`、`GOWORK=off go test`、golangci-lint、死链 0、头部四字段齐全、canonical 唯一、编号唯一、indices 一致、```go 块实测）。必须全部通过。
- **语义观察门**：默认 `continue-on-error`，用于持续度量，达标后按 §6.2 规则转阻断：
  1. stub 纯净度（stub 无正文残留，canonical 链接有效）
  2. 交叉域覆盖（§4.1 清单逐项有非 stub 权威页）
  3. KG 谓词精度（核心实体周边无通用关系边）
  4. 双向链接覆盖（版本特性↔概念页、`scripts/tmp/check_bidir.py` 辅助）
  5. 版本语义注入（`examples/goXXX-features/` ↔ 概念权威页双向链接）
  6. go 指令一致性（42 个 `go.mod` 的 `go` 指令无漂移）

### 6.2 观察门转正规则

- 任一观察门必须**连续 4 周（或连续 10 次 CI 运行）达标**，且本地检查当前 exit 0，才可转阻断。
- **禁止**以口头或一次性指示绕过该规则。
- 转正后若指标退化，按阻断门流程处置（修复或经评估后回调为观察门）。

### 6.3 “完成”声明纪律

- 禁止未经验证声明“已完成”“全部通过”“100%”。
- 任何完成声明必须引用可机器复核的证据：`bash .githooks/pre-commit` 输出、`scripts/check-unfixed-links.ps1` 结果、`scripts/tmp/verify_sixpiece.py` 结果、CI 记录。
- 声明必须覆盖全部阻断门与语义观察门（若观察门状态变化，须一并核对）。

---

## 7. 与质量门对应关系

| 语义要求 | 检查方式 | 通过标准 |
| --- | --- | --- |
| 定义唯一性 / canonical | `python scripts/tmp/dup_canon_scan.py` + complete-map 查重 | 0 双权威页 |
| 跨文件一致性 | 关键术语复扫 + 季审抽样（对照 Go Spec / Effective Go） | 0 错误级发现 |
| Glossary 对齐 | 对照 `go-knowledge-base/GLOSSARY.md` | 差异清单可解释或清零 |
| go 指令单一事实源 | `grep -rn "^go " --include=go.mod .` | 全部 `go 1.27` |
| 交叉域覆盖 | §4.1 清单逐项核对 | 12/12 覆盖（或经评估豁免并注明） |
| KG 谓词精度 | 核心实体周边 KG 边审计 | 无通用关系边 |
| 版本语义注入 | `python scripts/tmp/check_bidir.py` + examples/goXXX-features 双向链接 | 覆盖率达标 |
| 六件套完整性 | `python scripts/tmp/verify_sixpiece.py` | 权威页 100% 覆盖 |
| ```go 块实测 | `python scripts/tmp/extract_go_blocks.py` 提取后 `go vet` | 抽样块全部通过 |
| 综合语义健康 | 月审（`monthly_semantic_review.md`）+ 死链复扫 `python scripts/tmp/rescan_deadlinks.py` | grade = OK，死链 0 |

---

## 8. Kimi 生成前自检

对每页新增/修改内容，生成前确认：

- [ ] 该主题在五维权威层是否已存在权威页？若存在，只建 stub 或补充，不重复。
- [ ] 是否涉及关键术语（happens-before、方法集、unsafe、cgo 指针、泛型实例化等）？若涉及，检查跨文件定义一致性。
- [ ] 是否属于 §4.1 交叉/边界语义域之一？若是，按 §4 新建独立权威页。
- [ ] 是否涉及前沿提案 / GOEXPERIMENT？若是，按 §5 模板书写（该机制为可选扩展）。
- [ ] KG 关系是否使用具体谓词（dependsOn/entails/mutexWith/refines/equivalentTo/counterExample）？
- [ ] 是否形成至少一对双向链接（版本特性↔概念、工程页↔概念、learning-path↔权威页）？
- [ ] 声明“完成”前是否已跑 `bash .githooks/pre-commit` 并保存输出？

---

## 9. 修订历史

- 2026-09-05: 初版，整合语义审计、KG 谓词、交叉域覆盖、前沿提案跟踪计划与 `AGENTS.md` 观察门纪律。
- 2026-09: Go 版重写（从 rust-lang 项目适配）— scope 迁移至 `E:/_src/golang`；概念树映射为五维 `go-knowledge-base/`；最低版本口径改为 `go` 指令即最低支持版本；原交叉域清单替换为 Go 交叉域（cgo×unsafe、defer×闭包、泛型×反射等）；原效应系统页模板机制改为 Go 前沿提案跟踪页（标可选扩展）；质量门映射到本仓库真实脚本（`scripts/tmp/dup_canon_scan.py`、`check_bidir.py`、`verify_sixpiece.py`、`extract_go_blocks.py`、`rescan_deadlinks.py` 等）。
