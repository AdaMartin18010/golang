# messaging

消息队列示例，目前仅包含 NATS 示例。

## 结构

- `nats/cmd/publish-subscribe/main.go` — NATS 发布/订阅示例
- `nats/cmd/request-reply/main.go` — NATS 请求/应答示例

## 运行

本目录无独立 go.mod，属于上级 `examples` 模块（`example.com/golang-examples`，go 指令 `go 1.27`）。运行需要本地或远端的 NATS Server。在 `examples/` 目录下执行：

```bash
cd examples
GOWORK=off go run ./messaging/nats/cmd/publish-subscribe
GOWORK=off go run ./messaging/nats/cmd/request-reply
```

> 注意：当前 `examples/go.sum` 缺少 `github.com/nats-io/nats.go` 的校验条目，`GOWORK=off` 构建会报 `missing go.sum entry`（既有问题，未在本期修复）。可先执行 `GOWORK=off go mod tidy` 补齐后再运行；或使用 workspace 模式（`go work` 已包含 `./examples`）不加 `GOWORK=off` 运行。
