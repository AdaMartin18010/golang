# Kimi 项目协作要求模板（Go 版，可复用）

> **用途**：当你要用 Kimi 维护本 Go 项目（或创建类似的分层、可验证、可搜索 Go 知识体系）时，遵守本要求。
> **来源**：从 rust-lang 项目（E:/_src/rust-lang）的 `AGENTS.md`、质量门脚本与目录治理实践提炼，并已全面适配本仓库（E:/_src/golang，Go 1.27.1）。
> **核心原则**：**单一权威来源、机器可验证、语义可追溯**。

---

## 1. 项目定位与架构

### 1.1 核心定位

项目是**分层、可验证、可搜索**的 Go 知识体系 + 可编译代码库，为每个主题维护**单一、权威、可演进**的解释来源。

- 不要堆叠文档，要维护权威来源。
- 每个概念/主题只能有一个权威页（canonical page）。
- 新增内容前先查重（搜索 `go-knowledge-base/indices/`）；发现重复时按 canonical 规则合并或 stub 化。

### 1.2 目录职责（本仓实际）

| 目录 | 职责 | 是否可作为权威来源 |
| --- | --- | --- |
| `go-knowledge-base/01-Formal-Theory/` ~ `05-Application-Domains/` | 权威概念层（五维），每个主题的唯一深度解释（LD-NNN / FT-NNN / EC-NNN / TS-* / AD-NNN 编号体系） | ✅ 是 |
| `go-knowledge-base/indices/` | 多维索引（by-date / by-topic / by-difficulty / complete-index 等），新增权威页必须同步登记 | ❌ 索引，不是内容 |
| `go-knowledge-base/learning-paths/` | 学习路径编排 | ❌ 引用权威页，不替代 |
| `go-knowledge-base/.authoritative-sources/` | 权威外部源存档（论文/课件/代码片段/TLA+） | ❌ 素材库 |
| `docs/` | 指南、版本分析、研究报告 | ❌ 概念解释必须链接到五维权威页 |
| `examples/` | 可运行代码示例（go125 / go126-features / go127-features 等） | ❌ 概念解释不能放在这里 |
| `view/` | 专题视角与形式化分析 | ⚠️ 仅当五维未覆盖时，并回链 |
| `pknowledge/` | 个人学习笔记与卡片盒 | ❌ 私有层 |
| `pkg/` + `internal/` + `cmd/` | 生产代码与框架（42 个 go module） | ❌ 实现，不是概念权威 |
| `archive/` | 只读历史归档 | ❌ 不是权威来源 |
| `.kimi/` | 协作治理模板（本目录） | ❌ 治理文件 |

### 1.3 认知分层（Bloom / L0–L7）

为每个内容页标注 Bloom 层级：

- **L0**: 元层 / 框架 / 方法论
- **L1**: 基础语法与语义
- **L2**: 类型系统、控制流、函数
- **L3**: 泛型、接口、类型约束
- **L4**: 并发、channel、unsafe、CGO、运行时
- **L5**: 跨语言/范式对比、工程架构
- **L6**: 生态、设计模式、算法、系统设计
- **L7**: 未来特性、研究、预览（如 GOEXPERIMENT）

跨层引用规则：

- 高层（L6）页面前置概念中必须包含至少一个低层（如 L5）链接。
- 禁止循环自引用或自引用作为前置/后置概念。

---

## 2. 文件命名与格式

### 2.1 命名规范（本仓实际约定）

- 权威页命名：`{维度前缀}-NNN-{Kebab-Title}.md`，如 `LD-037-Go-1.27-Generic-Methods.md`。
- 目录内文件使用两位连续序号 `NN-`（从 01 起）；`00_` 保留给导览/README。
- 禁止中文文件名、空格；标题用中文 + 英文双语（见元数据模板）。
- 禁止双前缀（如 `06_20_`）与异形前缀（如 `1_2_`）。
- 新增权威页后，必须同步更新 `go-knowledge-base/indices/` 中的 by-date / by-topic / complete-index 等索引。

### 2.2 Markdown 文件元数据模板

每个权威页顶部必须包含（对齐本仓 LD 实际格式）：

```markdown
# LD-NNN: 中文标题 (English Title)

> **维度**: Language Design
> **级别**: S (16+ KB) / A / B
> **标签**: #tag1 #tag2 #go127
> **Go 版本**: 1.27+   <!-- 或项目对应版本 -->
> **Bloom 层级**: Lx
> **权威来源**: 本文件为 `02-Language-Design/` 权威页。
> **前置概念**: [概念A](../LD-0xx/xx.md) · [概念B](../yy.md)
> **后置概念**: [概念C](../LD-0zz/zz.md)
> **定理链**: Input → Operation → Output / Invariant
>
> - [Go Spec §xxx](https://go.dev/ref/spec#xxx) - P0
> - [相关 Proposal](https://go.googlesource.com/proposal/+/HEAD/design/xxxxx.md) - P0
```

### 2.3 Stub / 重定向模板

非权威位置或合并后的文件必须改为 stub，正文不超过 25 行 / 2000 字节：

```markdown
# 中文标题

> **权威来源**: [LD-0xx 中文标题](../02-Language-Design/LD-0xx-Xxx.md)
> 本文件为重定向 stub：完整解释请见上述权威页。
```

---

## 3. 内容质量要求

### 3.1 章节结构（每页必备）

1. **权威定义**：一句话定义 + 核心约束（可用"定义 X.Y"编号，如 LD-010 形式化理论页）。
2. **核心机制**：分小节讲原理，配可运行的 ```go 代码示例。
3. **工程实践**：使用场景、权衡、最佳实践。
4. **反命题与边界分析**：至少一个反例（Go 无 compile_fail 属性，用标注 + 实测说明，见 3.2）。
5. **思维导图**：Mermaid `mindmap`。
6. **参考来源 / References**：P0 官方 + P1 学术 + P2 生态。

### 3.2 代码块规范

- 可运行示例使用 ```go，且必须实际通过 `go build` / `go test` 验证（本仓 `examples/` 即此用途）。
- 故意编译失败的反例使用 ```go 并紧跟 `// 编译失败:` 标注，说明报错原因；保持与结论一致的实测行为。
- 需要特定 GOEXPERIMENT 的示例（如 `simd`）需注明环境要求与降级路径。
- 伪代码/片段使用 ```text。
- 每个新增概念页至少包含一个可运行 `go` 块和一个编译失败反例块。

### 3.3 反例要求

- 每页必须包含"反命题与边界分析"节。
- 反例要解释错误原因，而不是只贴代码。
- 反例覆盖率应达到项目设定的基线（如 ≥40%）。

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

每个内容页必须在 References 中覆盖至少一个 P0、P1、P2 来源：

| 级别 | 含义 | 典型域名 |
| --- | --- | --- |
| **P0 官方** | Go 语言/工具链官方文档与提案 | `go.dev/ref/spec`, `go.dev/doc`, `go.dev/blog`, `pkg.go.dev`, `go.googlesource.com/proposal`, `github.com/golang/go` |
| **P1 学术/形式化** | 论文、形式化验证、经典教材 | `arxiv.org`, `acm.org`, `dl.acm.org`, `ieee.org`, `springer.com`（如 *The Go Memory Model*、CSP 原始论文） |
| **P2 生态/社区** | 知名 Go 项目、技术博客 | `github.com`（gin/ent/grpc-go 等）, `go.dev/wiki`, `research.swtch.com` |

示例 References 节：

```markdown
## 参考来源 / References

- **P0 官方**: [Go Spec](https://go.dev/ref/spec) · [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
- **P1 学术**: [The Go Memory Model](https://go.dev/ref/mem)
- **P2 生态**: [research.swtch.com](https://research.swtch.com) · [Go Wiki](https://go.dev/wiki)
```

---

## 5. 链接与交叉引用

### 5.1 死链检查

- 所有本地 markdown 链接必须有效。
- 新增页面前必须运行链接检查（`scripts/check-unfixed-links.ps1` 等）。
- 重定向 stub 的 canonical 链接不能失效。

### 5.2 跨层引用

- 高层页面必须引用低层权威页。
- 禁止死端页面（无出链/入链的孤立页）。
- 前置/后置概念中的相对路径必须正确。

### 5.3 双向链接

- 版本特性、示例、测验等必须与权威页形成双向链接。
- 新增权威页后，应同步更新 `go-knowledge-base/indices/` 相关索引。

---

## 6. 代码与构建规范

### 6.1 多模块规范（本仓 42 个 go.mod）

- 所有模块保持统一的 `go 1.27` 版本声明，升级时全仓批量对齐（含 CI、Dockerfile、go.work）。
- 根 `go.work` 不覆盖全部模块；**在仓库根目录运行 go 命令必须加 `GOWORK=off`**，或进入具体模块目录操作。
- 本地 `replace` 指向的子模块（如 `pkg/observability`），构建上下文必须包含其 go.mod/go.sum（见 `deployments/docker/Dockerfile` deps 阶段）。

### 6.2 构建验证

在每个模块内：

```bash
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off go test ./...
```

- `golangci-lint run`（配置见 `.golangci.yml`）必须通过。
- 容器构建：`DOCKER_BUILDKIT=0 docker build -f deployments/docker/Dockerfile -t <tag> .`（本机网络对 buildkit 拉取 metadata 不稳，经典构建器为既定绕行方案）。

---

## 7. 质量门（Quality Gates）

### 7.1 阻断门（必须全部通过）

本仓质量设施（rust-lang 项目的 23 条 Python 脚本门已映射为以下实际设施）：

1. `gofmt -l`（无输出，见 `.githooks/pre-commit` 第 1 步）
2. `go vet ./...`（每模块，GOWORK=off）
3. `golangci-lint run --timeout=5m`（若已安装）
4. `go test -short`（改动包，见 `.githooks/pre-commit` 第 4 步）
5. `.github/workflows/go-test.yml` — CI 全量测试
6. `.github/workflows/go-lint.yml` — CI 静态检查
7. `.github/workflows/go-fix.yml` — go fix modernizer 检查
8. `.github/workflows/docs-check.yml` — 文档格式与链接检查
9. `.github/workflows/knowledge-tracker.yml` — 知识库跟踪
10. `scripts/check-unfixed-links.ps1` — 死链检查
11. `scripts/check-markdown-format.ps1` — Markdown 格式
12. `scripts/check_quality.ps1` — 综合质量检查

### 7.2 语义观察门（达标后转阻断）

- 权威页唯一性：同一主题仅一个 canonical 页（查 `go-knowledge-base/indices/` 登记）。
- stub 纯净度：stub 正文 ≤ 25 行 / 2000 字节。
- 交叉/边界语义覆盖：跨维度主题（如并发形式模型）必须有权威页。
- 版本语义注入：Go 1.26/1.27 特性页必须映射回对应概念权威页（如 LD-037 ↔ examples/go127-features）。
- 双语完整性：权威页标题与 EN/Summary 标注齐全。

---

## 8. 新增内容工作流

1. **查重**：搜索 `go-knowledge-base/indices/` 确认主题是否已有权威页。
2. **定位层级**：确定 Bloom 层级与维度目录（五维之一）。
3. **创建权威页**：按 §2.2 元数据模板、§3.1 章节结构、§3.2 代码块规范书写。
4. **添加 References**：确保 P0/P1/P2 全覆盖（§4）。
5. **链接前置/后置**：至少包含一个低层链接，无死链。
6. **同步索引**：在 `go-knowledge-base/indices/` 的 by-date / by-topic / complete-index 登记新页。
7. **跑检查**：gofmt、go vet、链接检查。
8. **提交**：`git add -A && git commit`（本仓提交信息惯例为 `update`；push 由用户决定）。

---

## 9. 红线与禁止事项

- 不要在 `go-knowledge-base/indices/` 中手造与权威页冲突的内容；索引只登记、不解释。
- 不要把临时文件提交到版本控制。
- 不要在 `examples/` 或 `pkg/` 中复制通用概念解释；应链接到五维权威页。
- 禁止未经验证的"完成"声明；必须通过机器可复核的检查（构建、测试、链接）。
- Stub/redirect 文件正文不得超过 25 行 / 2000 字节。
- KG 关系必须使用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`），避免通用 `Relation`。
- 交叉/边界语义域必须有权威页。
- 版本特性必须映射回概念权威页（如 go127 特性 → LD-037 / go-knowledge-base 相应页）。
- 在仓库根目录禁止裸跑 go 命令（go.work 不含全部模块），必须 `GOWORK=off`。

---

## 10. 推荐工具链

- **构建**：go 1.27.1 toolchain、gofmt、golangci-lint（`.golangci.yml`）
- **文档**：Mermaid（mindmap/flowchart）、`docs/templates/TECHNICAL-ARTICLE-TEMPLATE.md`
- **审计**：`scripts/` 下 PowerShell / Python 检查脚本
- **CI**：`.github/workflows/`（go-test / go-lint / go-fix / docs-check / knowledge-tracker）
- **预提交**：`.githooks/pre-commit`

---

## 11. 快速检查清单

新增/修改权威页前自问：

- [ ] 主题是否已存在权威页？（查 indices）
- [ ] 双语标题与 EN/Summary 是否填写？
- [ ] Bloom 层级与维度目录是否正确？
- [ ] 前置/后置概念链接是否有效且含低层链接？
- [ ] 是否包含可运行 `go` 代码块与编译失败反例？
- [ ] 是否有 Mermaid mindmap？
- [ ] References 是否覆盖 P0/P1/P2？
- [ ] 是否已同步登记 `go-knowledge-base/indices/`？
- [ ] `gofmt` / `go vet` 是否通过？
- [ ] 死链检查是否通过？

---

> **提示**：本模板已按本 Go 项目实际（Go 1.27.1、五维知识库、42 模块、GOWORK=off 约定）校准。后续版本升级（如 Go 1.28）时，同步更新 §2.2、§4、§6、§7 中的版本引用。
