# 后量子 TLS 部署指南

> **文档编号**: DEP-2026-PQTLS
> **适用版本**: Go 1.26+
> **最后更新**: 2026-05-06

---

## 一、背景

Go 1.26 在 `crypto/tls` 中默认启用了**混合后量子密钥交换**：

```text
默认启用的密钥交换算法:
────────────────────────────────────────
SecP256r1MLKEM768   ← 默认首选
SecP384r1MLKEM1024  ← 高安全需求

特性:
├─ 结合经典椭圆曲线 (ECDHE) 和 ML-KEM (NIST 后量子标准)
├─ 即使未来量子计算机破解 ECDHE，ML-KEM 仍提供安全
├─ 零配置：升级 Go 1.26 即自动启用
└─ 与现有 TLS 1.3 客户端/服务器完全兼容
```

---

## 二、服务器配置

```go
// 零配置即可启用后量子 TLS
// Go 1.26 默认已包含后量子密钥交换

// 如需显式控制（高级场景）
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS13,
    // Go 1.26 默认已启用 MLKEM，无需额外配置
}

server := &http.Server{
    Addr:      ":443",
    TLSConfig: tlsConfig,
}
```

---

## 三、客户端配置

```go
// 客户端同样自动使用后量子密钥交换
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS13,
}

tr := &http.Transport{
    TLSClientConfig: tlsConfig,
}
client := &http.Client{Transport: tr}
```

---

## 四、验证后量子连接

```go
// 检查实际使用的密钥交换算法
conn, err := tls.Dial("tcp", "example.com:443", &tls.Config{
    MinVersion: tls.VersionTLS13,
})
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

state := conn.ConnectionState()
fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(state.CipherSuite))
fmt.Printf("TLS Version: %s\n", versionName(state.Version))

// Go 1.26 新增: 检查是否使用后量子密钥交换
// 如果 server 支持，state 会反映 ML-KEM 协商结果
```

---

## 五、兼容性矩阵

| 客户端 | 服务器 | 结果 |
| -------- | -------- | ------ |
| Go 1.26 | Go 1.26 | ✅ 后量子 + 经典混合 |
| Go 1.26 | Go 1.25 | ✅ 经典 ECDHE（回退） |
| Go 1.25 | Go 1.26 | ✅ 经典 ECDHE（回退） |
| 旧浏览器 | Go 1.26 | ✅ 经典算法（回退） |

---

## 六、性能影响

| 指标 | 经典 TLS | 后量子 TLS | 变化 |
| ------ | --------- | ----------- | ------ |
| 握手延迟 | 1 RTT | 1 RTT | 无变化 |
| 公钥操作 | ECDHE | ECDHE + ML-KEM | +~20% CPU |
| 密钥大小 | ~32 bytes | ~1184 bytes (MLKEM768) | 增大 |

**建议**: 对于高并发场景，可考虑 session resumption (0-RTT) 减少握手频率。

---

## 七、部署检查清单

- [ ] Go 版本升级到 1.26.2+
- [ ] TLS 1.3 已启用（Go 1.26 默认）
- [ ] 监控握手失败率（确认无兼容性回归）
- [ ] 测试与旧客户端的连接回退
