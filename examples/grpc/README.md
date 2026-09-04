# grpc

gRPC 服务端/客户端最小示例。代码不依赖 internal 包，服务端内联定义了简化的 gRPC 服务；客户端中生成代码的 import 已注释，接入真实 proto 后取消注释即可。

## 结构

- `cmd/server/main.go` — gRPC 服务端示例
- `cmd/client/main.go` — gRPC 客户端示例（TODO：生成 gRPC 代码后取消注释 proto import）

## 运行

本目录无独立 go.mod，属于上级 `examples` 模块（`example.com/golang-examples`，go 指令 `go 1.27`）。在 `examples/` 目录下执行：

```bash
cd examples
GOWORK=off go run ./grpc/cmd/server
GOWORK=off go run ./grpc/cmd/client
```

> 注意：当前 `examples/go.sum` 缺少 `google.golang.org/grpc` 等依赖的校验条目，`GOWORK=off` 构建会报 `missing go.sum entry`（既有问题，未在本期修复）。可先执行 `GOWORK=off go mod tidy` 补齐后再运行；或使用 workspace 模式（`go work` 已包含 `./examples`）不加 `GOWORK=off` 运行。
