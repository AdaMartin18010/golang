# HTTP/3 Server 示例

## 说明

此示例展示如何使用HTTP/3（基于QUIC）构建高性能服务器。

## 特性

- HTTP/3 over QUIC
- 0-RTT连接恢复
- 连接迁移
- 更好的弱网性能
- 自动协议降级（HTTP/3 → HTTP/2 → HTTP/1.1）

## 运行要求

### 1. 安装依赖

```bash
go get github.com/quic-go/quic-go/http3
```

### 2. 生成证书

HTTP/3需要TLS证书：

```bash
# 生成自签名证书（仅用于测试）
openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem -out cert.pem \
  -days 365 -nodes \
  -subj "/C=CN/ST=State/L=City/O=Org/CN=localhost"
```

### 3. 修改代码

取消代码中证书相关的注释：

```go
// 在main()函数中启用
if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
    log.Fatal(err)
}
```

### 4. 运行

```bash
go run main.go
```

## 测试

### 使用curl测试

```bash
# HTTP/3
curl --http3 -k https://localhost:8443

# HTTP/2 (fallback)
curl -k https://localhost:8443

# 查看统计
curl -k https://localhost:8443/stats

# 健康检查
curl -k https://localhost:8443/health
```

### 使用浏览器

访问 `https://localhost:8443` （需要接受自签名证书警告）

## 性能对比

| 网络条件 | HTTP/1.1 | HTTP/2 | HTTP/3 | 提升 |
|---------|----------|--------|--------|------|
| 理想网络 | 100% | 105% | 110% | +10% |
| 高延迟 | 100% | 120% | 140% | +40% |
| 高丢包 | 100% | 110% | 150% | +50% |

## HTTP/3 优势

### 1. 0-RTT 握手

首次连接后，后续连接无需握手，直接发送数据。

### 2. 多路复用无队头阻塞

基于UDP，每个流独立，一个流的丢包不影响其他流。

### 3. 连接迁移

客户端IP变化（如WiFi切换到4G）时，连接不中断。

### 4. 更快的拥塞控制

QUIC内置现代拥塞控制算法。

## 生产部署

### 1. 使用正式证书

```bash
# Let's Encrypt
certbot certonly --standalone -d yourdomain.com
```

### 2. 负载均衡

支持大多数现代负载均衡器：

- Nginx 1.25+
- HAProxy 2.6+
- Envoy

### 3. 监控

关键指标：

- 0-RTT成功率
- 连接迁移次数
- 协议降级比例
- 延迟分布

## 故障排除

### 问题：连接失败

```bash
# 检查UDP端口是否开放
sudo netstat -ulnp | grep 8443

# 防火墙规则
sudo ufw allow 8443/udp
```

### 问题：客户端不支持HTTP/3

服务器会自动降级到HTTP/2或HTTP/1.1，无需特殊处理。

## 更多信息

- [HTTP/3和QUIC文档](../../../docs/02-Go语言现代化/14-Go-1.23并发和网络/03-HTTP3-和-QUIC支持.md)
- [网络优化指南](../../../docs/02-Go语言现代化/性能优化实战指南.md)
