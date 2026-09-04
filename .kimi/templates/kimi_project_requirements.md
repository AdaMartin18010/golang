# Kimi 项目协作要求模板（Go 版，可复用）

> **用途**：当你要用 Kimi 创建与 `golang` 知识库类似的分层、可验证、可搜索 Go 知识体系 + 可编译代码库时，可直接复用本要求。
> **来源**：从 `E:/_src/golang` 的 `AGENTS.md`、`.gimi/` 治理模板与 2026-09 内容合规审计实践中提炼。
> **当前基线**：Go 1.27.1（2026-08 发布）；42 个 go module，统一 `go 1.27`。

---

## 1. 项目定位与架构

### 1.1 核心定位

项目应为**分层、可验证、可搜索**的 Go 知识体系 + 可编译代码库，目标是为每个主题维护**单一、权威、可演进**的解释来源。

- 不要堆叠文档，要维护权威来源。
- 每个概念/主题必须只有一个权威页（canonical page）。
- 新增内容前先查重；发现重复时按 canonical 规则合并或 stub 化（stub ≤ 25 行 / 2000 字节，只留一句话 + canonical 链接）。

### 1.2 推荐目录职责

| 目录 | 职责 | 是否可作为权威来源 |
| --- | --- | --- |
| `go-knowledge-base/01..05-*` | 权威概念层（五维：Formal-Theory / Language-Design / Engineering-CloudNative / Technology-Stack / Application-Domains），FT/LD/EC/TS/AD 编号 | ✅ 唯一权威层 |
| `go-knowledge-base/indices/` | 多维索引（by-date / by-topic / complete-map 等） | ❌ 新增权威页必须同步登记 |
| `go-knowledge-base/learning-paths/` | 学习路径 | ❌ 引用权威页 |
| `docs/` `view/` | 指南、版本分析、形式化专题 | ⚠️ 概念解释必须回链五维权威页 |
| `examples/` | 可运行示例（go125/go126-features/go127-features 等） | ❌ |
| `pkg/` `internal/` `cmd/` | 生产代码 | ❌ |
| `pknowledge/` | 个人笔记 | ❌ |
| `archive/` | 只读历史 | ❌ 不是权威来源 |
| `test/` | 测试（unit / integration / e2e） | ❌ 不能替代概念解释 |

### 1.3 认知分层（Bloom / L0–L7）

为每个权威页标注 Bloom 层级：

- **L0**: 元层 / 框架 / 方法论
- **L1**: 基础语法与语义（变量、控制流、函数）
- **L2**: 类型系统、接口、组合、错误处理
- **L3**: 泛型、反射、并发原语（channel/select/context）
- **L4**: 调度器（GMP）、GC、内存模型、cgo/汇编、运行时
- **L5**: 跨语言/范式对比、工程架构
- **L6**: 生态、设计模式、分布式系统设计
- **L7**: 未来特性、研究、预览（如 roadmap 提案）

跨层引用规则：

- 高层（L6/L7）页面前置概念中必须包含至少一个低层（如 L4）链接。
- 禁止循环自引用或自引用作为前置/后置概念。

---

## 2. 文件命名与格式

### 2.1 命名规范

- 权威页命名：`{维度前缀}-NNN-{Kebab-Title}.md`（如 `LD-037-Go-1.27-Generic-Methods.md`）。
- 子目录导览用 `README.md`；`00_` 保留给导览/框架页。
- 禁止中文文件名、空格（历史/豁免目录如 `view/Go1.26.1语法/` 已存档登记）。
- 编号在维度内唯一；冲突组保留入链最多者，其余重编为维度最大号递增。
- 专题系列可集中同一目录并配 README 索引。

### 2.2 权威页头部元数据模板

每个权威页顶部必须包含：

```markdown
# LD-NNN: 中文标题 (English Title)

> **维度**: Language Design
> **级别**: S (16+ KB) / A / B   <!-- 按实测大小：S >16KB，A >8KB，其余 B -->
> **标签**: #tag1 #go127
> **Go 版本**: 1.27+              <!-- 以 Go 1.27 工具链验证为基线 -->
> **Bloom 层级**: Lx               <!-- L0 元层 … L7 未来/研究 -->
> **前置概念**: [..](../LD-0xx.md) · **后置概念**: [..]
> **定理链**: Input → Operation → Output / Invariant
```

### 2.3 Stub / 重定向模板

非权威位置或合并后的文件必须改为 stub，正文不超过 25 行 / 2000 字节：

```markdown
# {原标题}

> **维度**: {维度} | **级别**: B (stub)
> **状态**: 占位 — 原内容与 canonical 页重复，已按红线合并。

**主题**: 一句话说明。

**权威页（canonical）**: [..](相对路径.md)
```

---

## 3. 内容质量要求

### 3.1 章节六件套（每页必备）

1. **权威定义**：一句话定义 + 核心约束。
2. **核心机制**：分小节讲原理，配可运行 ```go 代码示例。
3. **工程实践**：使用场景、权衡、最佳实践。
4. **反命题与边界**：至少一个含 `// 编译失败:` 注释的反例，并解释失败原因与正确写法。
5. **Mermaid mindmap**：覆盖定义、机制、边界、实践。
6. **References**：P0 官方 + P1 学术 + P2 生态，三者至少各一。

### 3.2 代码块规范

- 可运行示例使用 ```go，必须自包含（`package main` + `func main()` 或可独立测试）。
- 故意编译失败的反例使用 ```go，首行注释 `// 编译失败: <原因>`，并给出期望的编译器错误。
- 需要外部依赖的示例放入对应 `examples/<module>/` 并提供 go.mod，不在正文内假装自包含。
- 伪代码/输出片段使用 ```text。
- 每个新增权威页至少包含一个可运行 `go` 块和一个编译失败反例。

### 3.3 反例要求

- 每页必须包含「反命题与边界」节。
- 反例要解释错误原因（引用编译器错误文本），而不是只贴代码。
- 反例必须真实可复现：用 `GOWORK=off go vet` 或 `go build` 验证过。

### 3.4 思维导图

- 每页使用 Mermaid `mindmap`。
- 覆盖定义、机制、边界、对比/权衡。
- 避免纯文本堆砌，要有层次结构。

### 3.5 定理链与推理

- 每页在元数据中添加 `定理链`。
- 正文中使用 `⟹`（推出）和 `⟸`（反推）标记关键推理。
- 避免模板化定理链，要贴合本页主题。

---

## 4. 权威来源分级（P0 / P1 / P2）

每个权威页必须在 References 中覆盖至少一个 P0、P1、P2 来源：

| 级别 | 含义 | 典型域名 |
| --- | --- | --- |
| **P0 官方** | Go 语言规范、标准库、官方提案 | `go.dev/ref/spec`, `pkg.go.dev`, `go.googlesource.com/proposal`, `github.com/golang`, `go.dev/doc` |
| **P1 学术/形式化** | 论文、形式化验证、经典教材、内存模型 | `arxiv.org`, `acm.org`, `dl.acm.org`, `ieee.org`, `usenix.org`, `lamport.azurewebsites.net`, `raft.github.io` |
| **P2 生态/社区** | 知名 Go 项目、Go Blog | `go.dev/blog`, `github.com/<知名 Go 项目>`（如 etcd、grpc-go、prometheus、kubernetes client-go） |

示例 References 节：

```markdown
## References

- **P0 官方**: [The Go Programming Language Specification](https://go.dev/ref/spec) · [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
- **P1 学术**: [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) · [Dijkstra / Lamport 等经典文献]
- **P2 生态**: [Go Blog](https://go.dev/blog) · [etcd raft](https://github.com/etcd-io/raft)
```

---

## 5. 链接与交叉引用

### 5.1 死链检查

- 所有本地 markdown 链接必须有效（相对路径，区分大小写）。
- 新增页面前必须运行死链扫描（活跃区 0 死链；archive/ 只读不修并注明）。
- stub 的 canonical 链接不能失效。

### 5.2 跨层引用

- 高层页面必须引用低层权威页。
- 禁止死端页面（无出链/入链的孤立页）。
- 前置/后置概念中的相对路径必须正确。

### 5.3 双向链接

- 版本特性、示例代码等必须形成双向链接：版本页 → 概念权威页；权威页 → `examples/goXXX-features/`。
- 新增权威页后，必须同步登记 `go-knowledge-base/indices/`（by-date / by-topic / complete-index）。

---

## 6. 代码与构建规范

### 6.1 Go module 规范

- 42 个 go module 统一 `go 1.27`；根 `go.work` 不覆盖全部模块。
- **必须 `GOWORK=off`**：在仓库根跑 go 命令一律加 `GOWORK=off`（或进入具体模块目录）。
- 禁止在仓库根裸跑 go 命令。

### 6.2 构建验证

每模块：`GOWORK=off go build ./... && go vet ./... && go test ./...`；`golangci-lint run`（配置 `.golangci.yml`）。
Docker：`DOCKER_BUILDKIT=0 docker build -f deployments/docker/Dockerfile -t <tag> .`（经典构建器为既定方案）。
版本升级（如 1.28）时全仓批量对齐：所有 `go.mod`、`go.work`、CI、Dockerfile。

---

## 7. 质量门（Quality Gates）

### 7.1 阻断门（必须全部通过）

1. `gofmt -l`（无输出）
2. `GOWORK=off go vet ./...`
3. `GOWORK=off go test ./...`（相关包）
4. `golangci-lint run`
5. 死链扫描：活跃区 0 死链（扫描器跳过 archive/ 与代码围栏）
6. 头部四字段齐全：维度 / 级别 / 标签 / Go 版本（实质权威页 100%）
7. 权威页唯一性：同主题仅一个 canonical，重复已 stub 化
8. 编号唯一性：维度内 FT/LD/EC/TS/AD 编号无冲突（stub 别名除外并注明）
9. indices 一致：complete-map 登记数 == 实际篇数
10. 代码块实测：抽样提取权威页 ```go 块 `go vet` 通过

### 7.2 周期性审查门

- 月度语义审查：`.kimi/templates/monthly_semantic_review.md`（定义漂移、stub 纯净度、KG 关系、版本语义注入）。
- 季度权威源审计：`.kimi/templates/quarterly_international_source_audit.md`（对照 Go Spec / Effective Go / Go Blog / pkg.go.dev / Proposals）。
- 内容合规审计（按需）：命名、归档、错位归位、元数据、六件套覆盖率。

### 7.3 运行方式

```bash
# 本地（预提交钩子已内置）
bash .githooks/pre-commit

# 死链复扫（Python 扫描器口径：跳过 archive/、代码围栏，排除 http(s) 外链）
python scripts/tmp/rescan_deadlinks.py

# Go 验证（示例模块）
cd examples/go127-features && GOWORK=off go vet ./... && GOWORK=off go test ./...
```

---

## 8. 新增内容工作流

1. **查重**：搜索五维目录与 `indices/complete-map.md` 确认主题是否已存在权威页。
2. **定位层级**：确定 Bloom 层级与维度归属（FT/LD/EC/TS/AD）。
3. **创建权威页**：按 §2.2 元数据模板 + §3 章节六件套书写。
4. **添加 References**：确保 P0/P1/P2 全覆盖。
5. **链接前置/后置**：至少包含一个低层链接，无死链。
6. **登记索引**：by-date / by-topic / complete-index。
7. **跑质量门**：§7.1 全过。
8. **提交**：提交信息惯例为 `update`；push 由用户决定。

---

## 9. 红线与禁止事项

- 每个主题只有一个权威页；重复内容合并或 stub 化。
- 新增权威页必须同步登记 `go-knowledge-base/indices/`。
- 禁止未经验证的"完成"声明——以构建/测试/链接检查结果为准。
- KG 关系必须用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`），避免通用 `RelationAnnotation`。
- 版本特性必须映射回概念权威页（如 1.27 特性 → `LD-037` + `examples/go127-features/`）。
- 禁止在仓库根裸跑 go 命令；禁止把临时文件、构建产物提交入库。
- 不要在 `examples/` 的 README 里复制通用概念解释；应链接到五维权威页。

---

## 10. 推荐工具链

- **构建**：Go 1.27.1 toolchain（`GOWORK=off`）、golangci-lint、gofmt
- **文档**：Mermaid（mindmap）、Markdown 死链扫描器（Python）
- **审计**：`.kimi/templates/` 周期性审查模板 + `docs/tracking/` 跟踪记录
- **预提交**：`.githooks/pre-commit`

---

## 11. 快速检查清单

新增/修改权威页前自问：

- [ ] 主题是否已存在权威页？（查 complete-map）
- [ ] 双语标题 + 元数据齐全（维度/级别/标签/Go 版本/Bloom/定理链）？
- [ ] 前置含低层链接，无死链？
- [ ] 可运行 `go` 块 + `// 编译失败:` 反例 + Mermaid mindmap？
- [ ] References 覆盖 P0/P1/P2？
- [ ] 已登记 indices（by-date / by-topic / complete-index）？
- [ ] gofmt / go vet / 死链扫描通过？

---

> **提示**：本模板应根据具体项目调整（如语言版本、目录结构、特殊质量门）。核心原则是：**单一权威来源、机器可验证、语义可追溯**。
