# AGENTS.md — Go 项目协作约定

> 本文件是 AI 代理与人类协作者在本仓库工作的执行摘要。完整治理模板见 `.kimi/templates/`。
> 核心原则：**单一权威来源、机器可验证、语义可追溯**。当前工具链：**Go 1.27.1**（2026-08 发布）。

## 1. 项目定位

分层、可验证、可搜索的 Go 知识体系 + 可编译代码库（42 个 go module，统一 `go 1.27`）。每个概念/主题只有一个权威页（canonical）。

## 2. 目录职责

| 目录 | 职责 | 权威来源？ |
| --- | --- | --- |
| `go-knowledge-base/01..05-*` | 权威概念层（五维：Formal-Theory / Language-Design / Engineering-CloudNative / Technology-Stack / Application-Domains），LD/FT/EC/TS/AD 编号 | ✅ 唯一权威层 |
| `go-knowledge-base/indices/` | 多维索引（by-date / by-topic / complete-index 等） | ❌ 新增权威页必须同步登记 |
| `go-knowledge-base/learning-paths/` | 学习路径 | ❌ 引用权威页 |
| `docs/` `view/` | 指南、版本分析、形式化专题 | ⚠️ 概念解释回链五维 |
| `examples/` | 可运行示例（go125/go126-features/go127-features 等） | ❌ |
| `pkg/` `internal/` `cmd/` | 生产代码 | ❌ |
| `pknowledge/` | 个人笔记 | ❌ |
| `archive/` | 只读历史 | ❌ |

## 3. 文档规范

- **命名**：`{维度前缀}-NNN-{Kebab-Title}.md`（如 `LD-037-Go-1.27-Generic-Methods.md`）；`00_` 保留给导览。
- **权威页头部**（对齐 LD 实际格式）：

```markdown
# LD-NNN: 中文标题 (English Title)

> **维度**: Language Design
> **级别**: S (16+ KB) / A / B
> **标签**: #tag1 #go127
> **Go 版本**: 1.27+
> **Bloom 层级**: Lx   <!-- L0 元层 … L7 未来/研究 -->
> **前置概念**: [..](../LD-0xx.md) · **后置概念**: [..]
> **定理链**: Input → Operation → Output / Invariant
```

- **章节六件套**：权威定义 → 核心机制（可运行 ```go）→ 工程实践 → 反命题与边界（含 `// 编译失败:` 反例并解释原因）→ Mermaid mindmap → References。
- **References 分级**：P0 官方（go.dev/ref/spec、pkg.go.dev、go.googlesource.com/proposal、github.com/golang）＋ P1 学术（arxiv/acm/ieee、Go 内存模型等）＋ P2 生态（知名 Go 项目、go.dev/blog）。三者至少各一。

## 4. 构建与验证

- **必须 `GOWORK=off`**：CI 与发布流程一律 `GOWORK=off` 逐模块构建（go-test.yml 已矩阵化 38 模块，排除清单见 `docs/tracking/ci-module-coverage.md`）；本地根 `go.work` 已覆盖全部 42 个已跟踪模块，根 workspace 命令可不带，但进入具体模块目录 standalone 构建仍建议显式 `GOWORK=off`。
- 每模块：`GOWORK=off go build ./... && go vet ./... && go test ./...`；`golangci-lint run`（配置 `.golangci.yml`）。
- Docker：`DOCKER_BUILDKIT=0 docker build -f deployments/docker/Dockerfile -t <tag> .`（本机 buildkit 拉取 metadata 网络不稳，经典构建器为既定方案；builder 阶段已内置 `GOWORK=off`）。
- 版本升级（如 1.28）时全仓批量对齐：所有 `go.mod`、`go.work`、CI（`.github/workflows/`）、`deployments/docker/Dockerfile`。

## 5. 质量门

- 本地门：`.githooks/pre-commit`（gofmt → go vet → golangci-lint → 相关包测试）。
- CI 门：`.github/workflows/quality-gates.yml`（死链检查 `scripts/check-unfixed-links.ps1`（Exit 1 强制）+ 文档质量报告 `scripts/check_quality.ps1`（仅报告）+ gofmt（强制，范围 `cmd/` `internal/` `pkg/`））；其余 workflow 见 go-test / go-lint / go-fix / docs-check / knowledge-tracker 等。
- 文档辅助（本地手动跑）：`scripts/check-markdown-format.ps1`（依赖 markdownlint-cli，未装则回退内置基础检查）。
- gofmt 存量：`archive/` `examples/` `scripts/` `go-knowledge-base/` `test/` 尚有历史未格式化文件（289 个，其中 archive 占 105），不纳入 CI 强制门，存量清理另行安排。
- 周期性审查模板：`.kimi/templates/monthly_semantic_review.md`（月度）、`quarterly_international_source_audit.md`（季度，对照 Go Spec / Effective Go / Go Blog / pkg.go.dev / Proposals）。

## 6. 红线

- 每个主题只有一个权威页；重复内容合并或 stub 化（stub ≤ 25 行 / 2000 字节，只留一句话 + canonical 链接）。
- 新增权威页必须同步登记 `go-knowledge-base/indices/`。
- 禁止未经验证的"完成"声明——以构建/测试/链接检查结果为准。
- KG 关系必须用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`）。
- 版本特性必须映射回概念权威页（如 1.27 特性 → `LD-037` + `examples/go127-features/`）。
- 禁止在仓库根裸跑 go 命令；禁止把临时文件、构建产物提交入库。
- 提交信息惯例为 `update`；push 由用户决定。

## 7. 快速清单（新增/修改权威页）

- [ ] 查重（indices 中无同主题权威页）
- [ ] 双语标题 + 元数据齐全（维度/级别/标签/Go 版本/Bloom/定理链）
- [ ] 前置含低层链接，无死链
- [ ] 可运行 `go` 块 + 编译失败反例 + Mermaid mindmap
- [ ] References 覆盖 P0/P1/P2
- [ ] 已登记 indices（by-date / by-topic / complete-index）
- [ ] gofmt / go vet / 链接检查通过
