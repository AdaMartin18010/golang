# crypto-hpke

Go 1.26 `crypto/hpke` 标准库用法示例（HPKE 混合公钥加密，定义于 RFC 9180）。

演示 API 组合：KEM / KDF / AEAD 套件选择、接收方密钥对生成等。生产环境请使用完善的密钥管理，每次操作用临时密钥。

## 运行

```bash
cd examples/crypto-hpke
GOWORK=off go run .
```

要求 Go 1.26+（`crypto/hpke` 自 Go 1.26 引入）。本目录无独立 go.mod，使用工具链默认模块模式即可运行。
