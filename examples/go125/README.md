# go125

Go 1.25 特性示例集合，按主题分为三个子目录。多数子目录是独立 module（各有 go.mod，go 指令均为 `go 1.27`），运行/测试命令需在对应子目录内执行。

## 子目录清单

| 子目录 | 内容 | 运行方式 |
| --- | --- | --- |
| `concurrency-network/http3` | （空目录，占位） | — |
| `concurrency-network/synctest` | `testing/synctest` 实验性测试框架 | 目前只有 go.mod，无 `.go` 源文件，暂无可运行内容 |
| `runtime/container_scheduling` | 容器感知调度验证工具（Go 1.27+ container-aware GOMAXPROCS） | `cd examples/go125/runtime/container_scheduling && GOWORK=off go run .` |
| `runtime/gc_optimization` | GC 优化测试（greentea 测试） | `cd examples/go125/runtime/gc_optimization && GOWORK=off go test ./...` |
| `runtime/memory_allocator` | 内存分配器基准测试 | `cd examples/go125/runtime/memory_allocator && GOWORK=off go test -bench=. ./...` |
| `toolchain/asan_memory_leak` | ASan 内存泄漏检测示例（mock 实现，无需 CGO） | `cd examples/go125/toolchain/asan_memory_leak && GOWORK=off go run .` |

> 注：每个子目录需先 `cd` 进入后再执行 go 命令（独立 module）。
