# examples

可运行示例与测试集合。按专题分目录：

- `concurrency/` 并发与并发测试规范
- `servemux/` 新路由示例（Go 1.22+）
- `pgo/` PGO 示例（Go 1.21+）
- `slog/` 结构化日志示例（Go 1.21+）
- `observability/` OTel 最小可运行示例（Trace/Metrics/Logs）

运行：

```bash
go test ./...
```

运行 OTel 示例：

```bash
cd examples/observability
docker compose up -d
go run .
```
