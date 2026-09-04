# net-dialer-ctx

Go 1.26 `net.Dialer` 上下文感知拨号示例。

演示新增的 `DialIP` / `DialTCP` / `DialUDP` / `DialUnix` 方法，均接受 `context.Context` 以支持取消与超时控制。

## 运行

```bash
cd examples/net-dialer-ctx
GOWORK=off go run .
```

要求 Go 1.26+。本目录无独立 go.mod，使用工具链默认模块模式即可运行。
