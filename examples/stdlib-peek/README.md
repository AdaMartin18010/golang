# stdlib-peek

Go 1.26 `bytes.Buffer.Peek` 示例。

`Peek(n)` 返回后续 n 字节但不推进缓冲区，适用于解析协议时需要先检查数据再决定消费的场景。

## 运行

```bash
cd examples/stdlib-peek
GOWORK=off go run .
```

要求 Go 1.26+。本目录无独立 go.mod，使用工具链默认模块模式即可运行。
