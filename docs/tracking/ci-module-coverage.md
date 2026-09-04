# CI 模块覆盖矩阵（go.work × go-test）

> 维护者：R1a 任务（go.work 覆盖矩阵 + go-test CI 矩阵化）
> 生成日期：2026-08-25
> 依据：仓库 42 个已跟踪 go module（`git ls-files '*go.mod'`，排除 gitignored 的 `scripts/tmp/` 临时产物）

## 背景

- 根 `go.work` 原先仅覆盖 8 个模块；R1a 将其余 34 个模块逐一 `go work use` 纳入。
- 纳入校验：每加入一个模块执行 `go list -m all`（workspace 模式）验证模块图无版本冲突。
- **结果：42 个模块全部纳入 go.work，无一因版本冲突回退。**
- 注意：`go build ./...` 在仓库根只构建根模块包；workspace 模式下逐模块独立构建仍应使用 `GOWORK=off`。

## go-test CI 矩阵（38 个模块）

`.github/workflows/go-test.yml` 以 `fail-fast: false` 矩阵对每个模块执行
`GOWORK=off go build ./... && go vet ./... && go test -race ./... && go test -bench=. ./...`
（`03-testing-loop` 仅含测试文件，build 步骤按 GoFiles 计数自动跳过；coverage 上传仅根模块）。

| # | 模块 | go.work | CI 矩阵 | 备注 |
|---|------|---------|---------|------|
| 1 | `.`（根模块） | ✅ | ✅ | |
| 2 | `archive/ai-agent/examples/ai-agent` | ✅ | ✅ | |
| 3 | `archive/ai-agent/pkg/agent` | ✅ | ✅ | |
| 4 | `archive/formal-verification/tools/concurrency-pattern-generator` | ✅ | ✅ | |
| 5 | `archive/formal-verification/tools/formal-verifier` | ✅ | ✅ | |
| 6 | `examples` | ✅ | ✅ | R1a 修复：`go mod tidy` 补齐声明（此前依赖 workspace 遮蔽） |
| 7 | `examples/concurrency/pipeline_example` | ✅ | ✅ | |
| 8 | `examples/concurrency/worker_pool_example` | ✅ | ✅ | |
| 9 | `examples/go125/runtime/container_scheduling` | ✅ | ✅ | |
| 10 | `examples/go125/runtime/gc_optimization` | ✅ | ✅ | |
| 11 | `examples/go125/runtime/memory_allocator` | ✅ | ✅ | |
| 12 | `examples/go125/toolchain/asan_memory_leak` | ✅ | ✅ | CI 含 `-race`，与 asan 示例并用存在风险，见遗留问题 |
| 13 | `examples/go126-features` | ✅ | ✅ | |
| 14 | `examples/go127-features` | ✅ | ✅ | |
| 15 | `examples/modern-features/01-new-features/01-generic-type-alias` | ✅ | ✅ | |
| 16 | `examples/modern-features/01-new-features/03-testing-enhancement` | ✅ | ✅ | |
| 17 | `examples/modern-features/01-new-features/04-go125-new-features/01-iter-demo` | ✅ | ✅ | |
| 18 | `examples/modern-features/01-new-features/04-go125-new-features/02-unique-demo` | ✅ | ✅ | |
| 19 | `examples/modern-features/01-new-features/04-go125-new-features/03-testing-loop` | ✅ | ✅ | 仅含 `main_test.go`，build 步骤自动跳过 |
| 20 | `examples/modern-features/06-architecture-patterns/01-Clean-Architecture` | ✅ | ✅ | |
| 21 | `examples/modern-features/07-performance-optimization` | ✅ | ✅ | R1a 修复：`Buffer.ReadFrom/WriteTo` 改为标准 `(int64, error)` 签名，`BufferChain.ReadFrom` 更名 `FillFrom`（vet stdmethods） |
| 22 | `examples/modern-features/08-cloud-native/kubernetes-operator` | ✅ | ✅ | |
| 23 | `examples/modern-features/09-cloud-native-2.0/01-Kubernetes-Operator` | ✅ | ✅ | |
| 24 | `examples/testing-framework` | ✅ | ✅ | |
| 25 | `examples/web-crawler` | ✅ | ✅ | |
| 26 | `go-knowledge-base/examples/distributed-cache` | ✅ | ✅ | |
| 27 | `go-knowledge-base/examples/event-driven-system` | ✅ | ✅ | |
| 28 | `go-knowledge-base/examples/microservices-platform/services/user-service` | ✅ | ✅ | R1a 修复：删除未使用的 `fmt` 导入 |
| 29 | `go-knowledge-base/examples/rate-limiter` | ✅ | ✅ | |
| 30 | `go-knowledge-base/examples/task-scheduler` | ✅ | ✅ | |
| 31 | `pkg/concurrency` | ✅ | ✅ | R1a 修复：`context_test.go` 补 `defer cancel()`（vet lostcancel） |
| 32 | `pkg/http3` | ✅ | ✅ | |
| 33 | `pkg/memory` | ✅ | ✅ | |
| 34 | `pkg/observability` | ✅ | ✅ | |
| 35 | `scripts/fix_code_format` | ✅ | ✅ | |
| 36 | `scripts/format_docs` | ✅ | ✅ | |
| 37 | `scripts/gen_changelog` | ✅ | ✅ | |
| 38 | `scripts/project_stats` | ✅ | ✅ | |

## CI 排除清单（4 个模块，其中 3 个在 go.work 内）

| 模块 | go.work | CI 矩阵 | 排除理由 |
|------|---------|---------|----------|
| `examples/go125/concurrency-network/synctest` | ✅ | ❌ | 空模块：仅含 `go.mod`，无任何 Go 源文件，`go build/vet/test ./...` 均为 "no packages"（预存状态，非 R1a 引入） |
| `examples/modern-features/05-performance-toolchain/02-cgo-interop` | ✅ | ❌ | 本地独立构建失败：`performance_test.go` 报 `use of cgo in test ... not supported`，`basic/main.go` 的 C 代码在本机 llvm-mingw gcc 下报 `-Wimplicit-int` 错误。示例代码与 cgo 工具链深度耦合，修复超出 R1a 范围 |
| `examples/observability` | ✅ | ❌ | 本地独立构建失败：示例代码与 `pkg/observability` 当前 API 漂移（`systemMonitor.Start(ctx)` 签名不符、`Config` 结构体字段缺失、`status.Degraded` 未定义等 10+ 处）。R1a 曾尝试以 require+replace 修复依赖解析，编译错误属于示例源码层面，已回退 go.mod 改动 |
| `scripts`（根） | ✅ | ❌ | 本地独立构建失败：`wire/providers.go` 引用 `github.com/google/wire` 未声明于 go.mod，且导入 `github.com/yourusername/golang/internal/...`（internal 包跨模块越界，compiler 强制拒绝，replace 无法解决） |

## 不计入统计的临时模块

`scripts/tmp/`（`goblock_merge/`、`goblock_test/` 等 100+ 个 go.mod）为 gitignored 的临时测试产物，不属于 42 个模块统计口径，不纳入 go.work 与 CI。
