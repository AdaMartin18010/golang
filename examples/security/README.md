# security

安全相关示例，目前仅包含一个 HTTP 认证/授权示例。

## 结构

- `auth-example/main.go` — 基于 `go-chi/chi/v5` 的 HTTP 服务，演示 JWT 认证（`github.com/yourusername/golang/pkg/security/jwt`）与 RBAC 授权（`pkg/security/rbac`）的接入方式

## 运行

本目录无独立 go.mod，属于上级 `examples` 模块（`example.com/golang-examples`，go 指令 `go 1.27`），并依赖仓库根模块的 `pkg/security`。在 `examples/` 目录下执行：

```bash
cd examples
GOWORK=off go run ./security/auth-example
```

> 注意：当前 `examples/go.sum` 缺少根模块（`github.com/yourusername/golang`）与 `go-chi/chi` 的校验条目，`GOWORK=off` 构建会报 `missing go.sum entry`（既有问题，未在本期修复）。可先执行 `GOWORK=off go mod tidy` 补齐后再运行；或使用 workspace 模式（`go work` 同时包含 `.` 与 `./examples`）不加 `GOWORK=off` 运行。
