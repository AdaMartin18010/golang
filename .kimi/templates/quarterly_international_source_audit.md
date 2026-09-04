# Quarterly International Source Semantic Audit

**Audit quarter**: QX YYYY
**Auditor**: @username
**Scope**: 抽样 5–8 篇五维核心权威页，与对应国际权威来源（Go Spec、Effective Go、Go Blog、pkg.go.dev、Proposals、Go 内存模型）对照。目标：发现语义漂移、缺失边界条件、过期版本标注。审计记录归档 `docs/tracking/`。

---

## 1. Sample Selection

Select 5–8 pages from the priority list below, ensuring coverage across L1–L4 and at least one cross-domain boundary page:

| Priority | Page path | Authority source to compare | Reason for selection |
| ---: | --- | --- | --- |
| P0 | `go-knowledge-base/02-Language-Design/LD-001-Go-Memory-Model-Formal.md` | [Go Memory Model](https://go.dev/ref/mem) | Normative source; anchors must stay valid |
| P0 | `go-knowledge-base/02-Language-Design/LD-037-Go-1.27-Generic-Methods.md` | [Go Spec — Method declarations](https://go.dev/ref/spec#Method_declarations), [Go 1.27 Release Notes](https://go.dev/doc/go1.27) | 当前版本核心特性 |
| P0 | `go-knowledge-base/02-Language-Design/LD-003-Go-Garbage-Collector-Formal.md` | [Go Blog — GC 系列](https://go.dev/blog), [runtime 源码](https://go.dev/src/runtime/) | 跨域（GC × 调度） |
| P1 | `go-knowledge-base/02-Language-Design/LD-008-Go-Error-Handling-Patterns.md` | [Effective Go](https://go.dev/doc/effective_go), [errors 包文档](https://pkg.go.dev/errors) | L2 核心；高频漂移风险 |
| P1 | `go-knowledge-base/01-Formal-Theory/FT-006-Paxos-Formal.md` | [Raft paper](https://raft.github.io/raft.pdf), [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) | L4 形式化基线 |
| P1 | `go-knowledge-base/02-Language-Design/LD-002-Go-Concurrency-CSP-Formal.md` | [CSP 原始论文](https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf), [Go Spec — Channel types](https://go.dev/ref/spec#Channel_types) | 并发语义边界 |
| P2 | `go-knowledge-base/04-Technology-Stack/TS-008-NATS-Messaging-Patterns.md` | [NATS 官方文档](https://docs.nats.io/), [nats.go](https://github.com/nats-io/nats.go) | L6 生态 |

---

## 2. Comparison Dimensions

For each sampled page, check the following against the authority source:

| Dimension | Question to answer | Pass criterion |
| --- | --- | --- |
| **Definition drift** | 页面定义是否与权威来源一致？ | 无矛盾；项目特定扩展显式标注 `project-specific` |
| **Boundary completeness** | 所有例外、前置条件、未定义行为触发条件是否列全？ | 覆盖不少于权威来源 |
| **Version alignment** | `> **Go 版本**` 标注与当前稳定版（1.27.x）一致？ | ≤ 当前稳定 patch 版本 |
| **Anchor validity** | 外链到 Spec/pkg.go.dev 的锚点仍可达？ | 全部 HTTP 200（go.googlesource.com 不可达时注明并跳过） |
| **Code block freshness** | ```go 示例在当前稳定版下 vet 通过？反例确实失败？ | 抽样实测通过 |
| **Cross-link health** | 页面是否链接到相关权威页与 indices？ | 活跃区无死链 |

---

## 3. Audit Findings

| # | Page path | Authority source | Dimension | Finding | Severity (P0/P1/P2) | Action item |
| ---: | --- | --- | --- | --- | --- | --- |
| 1 | | | | | | |
| 2 | | | | | | |
| 3 | | | | | | |
| 4 | | | | | | |
| 5 | | | | | | |

---

## 4. Mandatory Commands

Run before sign-off:

```bash
# 死链复扫（活跃区 0；跳过 archive/、代码围栏；http(s) 外链不做本地判定）
python scripts/tmp/rescan_deadlinks.py

# 头部四字段全量校验（维度/级别/标签/Go 版本）
python scripts/tmp/normalize_headers.py   # 二次运行应为幂等（fields_added=0 或接近）

# 编号唯一性
python scripts/tmp/dup_id_scan.py

# Go 验证抽验（示例模块）
cd examples/go127-features && GOWORK=off go vet ./... && GOWORK=off go test ./...
```

---

## 5. Summary and Action Items

| Priority | Action | Owner | Due date |
|---|---|---|---|
| | | | |

**Overall audit conclusion**: [ ] No drift / [ ] Minor drift remediated / [ ] Major drift requires follow-up sprint

**Sign-off**: _________________
