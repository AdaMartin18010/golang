# Monthly Semantic Review Checklist

**Review period**: YYYY-MM
**Reviewer**: @username
**Scope**: 五维权威页（go-knowledge-base/01..05-*）的定义漂移、边界语义、stub 纯净度、KG 关系质量、版本语义注入覆盖。审计记录归档 `docs/tracking/`。

---

## 1. Core Concept Definition Drift (sample 10 pages)

Sample 10 篇实质权威页 across L1–L4，检查定义是否仍然精确、无歧义、与当前 Go 稳定版（1.27.x）一致。

| # | Page path | Definition checked | Drift found? | Action item |
|---|-----------|--------------------|--------------|-------------|
| 1 | `go-knowledge-base/02-Language-Design/LD-0xx-....md` |  | [ ] Yes [ ] No |  |
| 2 | `go-knowledge-base/01-Formal-Theory/FT-0xx-....md` |  | [ ] Yes [ ] No |  |
| 3 | `go-knowledge-base/03-Engineering-CloudNative/EC-0xx-....md` |  | [ ] Yes [ ] No |  |
| 4 | `go-knowledge-base/04-Technology-Stack/TS-0xx-....md` |  | [ ] Yes [ ] No |  |
| 5 | `go-knowledge-base/05-Application-Domains/AD-0xx-....md` |  | [ ] Yes [ ] No |  |
| 6 | `go-knowledge-base/02-Language-Design/02-Language-Features/xx.md` |  | [ ] Yes [ ] No |  |
| 7 | `go-knowledge-base/01-Formal-Theory/FT-0xx-....md` |  | [ ] Yes [ ] No |  |
| 8 | `go-knowledge-base/04-Technology-Stack/01-Core-Library/xx.md` |  | [ ] Yes [ ] No |  |
| 9 | `go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/xx.md` |  | [ ] Yes [ ] No |  |
| 10 | `go-knowledge-base/02-Language-Design/LD-0xx-....md` |  | [ ] Yes [ ] No |  |

**What to look for**:

- Definitions are operational（可转化为判定程序，而非口号）。
- No circular definitions.
- `> **Go 版本**` 标注与当前稳定版一致，或显式标记为实验/预览（`GOEXPERIMENT`）。
- 反例存在且在最新 Go 工具链下仍然编译失败（用 `GOWORK=off go vet` 抽验）。

---

## 2. Boundary Precision Review

Check key cross-domain / boundary semantics。每个边界域应有一个非 stub 五维权威页。

| Domain pair | 权威页 exists? | Boundary clearly stated? | Action item |
| ------------- | ---------------------------------- | -------------------------- | ------------- |
| GC × cgo/指针运算 | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| goroutine × signal/抢占调度 | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| 泛型 × 反射（type switches on type parameters） | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| 内存模型 happens-before × channel/select | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| error wrapping × `errors.Is/As` 边界 | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| context 取消 × goroutine 泄漏 | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| `//go:` 编译指令 × 链接器行为 | [ ] Yes [ ] No | [ ] Yes [ ] No | |
| 逃逸分析 × 栈/堆分配 | [ ] Yes [ ] No | [ ] Yes [ ] No | |

**Notes / additional boundary topics**:

---

## 3. New Stub Purity Review

- [ ] Stub 正文 ≤ 25 行 / 2000 字节
- [ ] 每篇 stub 恰好含一句话主题 + canonical 链接（或"待编写"声明 + complete-map 指针）
- [ ] stub 不含从原页残留的正文内容

**Violations found**:

| File path | Lines | Bytes | Issue | Remediation |
|-----------|-------|-------|-------|-------------|
|           |       |       |       |             |

**Policy reminder**: stub 必须纯净。超出"一句话 + canonical 链接"的内容必须移动到 canonical 权威页。

---

## 4. KG Relation Quality Review

- [ ] KG 关系使用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`）
- [ ] 核心实体无通用 `RelationAnnotation`

**Generic relations found**:

| Subject | Predicate | Object | Suggested concrete predicate |
|---------|-----------|--------|------------------------------|
|         |           |        |                              |

---

## 5. Version Semantic Injection Coverage Review

- [ ] Go 1.26–1.27 特性映射到权威页数 = _**/**_
- [ ] 每个版本页（LD-026/027/029/030/035/036/037 + 03-Evolution）回链概念权威页
- [ ] 每个被映射的权威页前向链接 `examples/go126-features/` 或 `examples/go127-features/`

**Unmapped or one-way-only features**:

| Feature | Version page | 权威页 | Issue | Action |
|---------|--------------|---------|-------|--------|
|         |              |         |       |        |

**Policy reminder**: 版本特性必须映射回概念权威页；孤立的 release notes 不充分。

---

## 6. 代码块实测抽验（推荐）

- [ ] 抽样 10 个 ```go 块提取到临时模块 `GOWORK=off go vet` 通过
- [ ] 抽样 5 个 `// 编译失败:` 反例确认确实编译失败

**Gaps**:

| Page | Block | 结果 | Action |
|------|-------|------|--------|
|      |       |      |        |

---

## 7. Summary and Action Items

| Priority | Action | Owner | Due date |
|----------|--------|-------|----------|
|          |        |       |          |

**Overall semantic health grade**: _**/ 100
**Sign-off**: _________________
