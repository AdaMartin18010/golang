# 可观测性完整示例

**版本**: 2024最新
**技术栈**: OpenTelemetry + Prometheus + Tempo + Loki + Grafana

---

## 🎯 功能展示

本示例展示完整的可观测性栈：

1. **分布式追踪** - Tempo (OpenTelemetry)
2. **指标监控** - Prometheus
3. **日志聚合** - Loki
4. **可视化** - Grafana
5. **eBPF 监控** - 系统级可观测性

---

## 🏗️ 架构

```text
┌─────────────┐
│   Go App    │ (OTLP gRPC)
└──────┬──────┘
       │
       ↓
┌──────────────────────┐
│  OTEL Collector      │ (0.114.0 最新)
│  - Receivers         │
│  - Processors        │
│  - Exporters         │
└──────┬───────────────┘
       │
       ├─→ Tempo (Traces)
       ├─→ Prometheus (Metrics)
       └─→ Loki (Logs)
              │
              ↓
       ┌─────────────┐
       │   Grafana   │ (可视化)
       └─────────────┘
```

---

## 🚀 快速开始

### 1. 启动可观测性栈

```bash
# 启动所有服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f collector
```

### 2. 访问服务

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Tempo**: http://localhost:3200
- **Loki**: http://localhost:3100
- **OTEL Collector Health**: http://localhost:13133/health
- **OTEL Collector zpages**: http://localhost:55679/debug/servicez

### 3. 运行 Go 应用

```bash
# 运行基础示例
cd app
go run main.go

# 运行 eBPF 监控示例 (需要 Linux)
cd ebpf-monitoring
sudo go run main.go

# 运行系统监控示例
cd system-monitoring
go run main.go
```

---

## 📊 验证可观测性

### 1. 查看追踪

1. 打开 Grafana: http://localhost:3000
2. 进入 Explore
3. 选择 Tempo 数据源
4. 搜索追踪

### 2. 查看指标

1. 打开 Grafana: http://localhost:3000
2. 进入 Explore
3. 选择 Prometheus 数据源
4. 查询: `otel_*` 或 `ebpf_*`

### 3. 查看日志

1. 打开 Grafana: http://localhost:3000
2. 进入 Explore
3. 选择 Loki 数据源
4. 查询: `{service_name="my-service"}`

---

## 🔧 配置说明

### OpenTelemetry Collector

**文件**: `otelcol.yaml`

**特性**:
- ✅ OTLP gRPC/HTTP receivers
- ✅ 批处理优化
- ✅ 内存限制保护
- ✅ 资源自动检测
- ✅ 智能采样
- ✅ 多后端导出

### Tempo

**文件**: `tempo.yaml`

**特性**:
- ✅ OTLP 接收
- ✅ 本地存储
- ✅ 7天数据保留
- ✅ Metrics 生成

### Prometheus

**文件**: `prometheus.yaml`

**特性**:
- ✅ 收集 OTEL Collector 指标
- ✅ 收集应用指标
- ✅ 15秒抓取间隔

---

## 📈 监控指标

### 应用指标

```
# HTTP 请求
http_requests_total
http_request_duration_seconds

# 系统资源
system_cpu_usage
system_memory_usage
system_disk_usage

# eBPF 指标
ebpf_syscall_count
ebpf_syscall_duration
ebpf_tcp_connections
ebpf_tcp_bytes
```

### Collector 指标

```
# 接收器
otelcol_receiver_accepted_spans
otelcol_receiver_refused_spans

# 处理器
otelcol_processor_batch_batch_send_size

# 导出器
otelcol_exporter_sent_spans
otelcol_exporter_send_failed_spans
```

---

## 🎯 最佳实践

### 1. 生产环境配置

```yaml
# otelcol.yaml 生产配置
processors:
  batch:
    timeout: 10s
    send_batch_size: 1024

  memory_limiter:
    limit_mib: 2048  # 根据实际调整
    spike_limit_mib: 512

exporters:
  otlp:
    tls:
      insecure: false  # 启用 TLS
      cert_file: /path/to/cert.pem
      key_file: /path/to/key.pem
```

### 2. 采样策略

```yaml
# 智能采样
processors:
  probabilistic_sampler:
    sampling_percentage: 10  # 10% 采样

  # 或使用尾部采样（更智能）
  tail_sampling:
    policies:
      - name: error-traces
        type: status_code
        status_code: {status_codes: [ERROR]}
      - name: slow-traces
        type: latency
        latency: {threshold_ms: 1000}
```

### 3. 资源限制

```yaml
# Docker Compose
services:
  collector:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

---

## 🔍 故障排查

### Collector 无法启动

```bash
# 检查配置
docker-compose logs collector

# 验证配置文件
docker-compose exec collector otelcol validate --config=/etc/otelcol/config.yaml
```

### 应用无法连接 Collector

```bash
# 检查端口
netstat -tlnp | grep 4317

# 测试连接
telnet localhost 4317
```

### 数据未显示在 Grafana

1. 检查 Collector 日志
2. 检查 Tempo/Prometheus 日志
3. 验证数据源配置
4. 检查时间范围

---

## 📚 相关资源

- [OpenTelemetry Collector 文档](https://opentelemetry.io/docs/collector/)
- [Tempo 文档](https://grafana.com/docs/tempo/latest/)
- [Prometheus 文档](https://prometheus.io/docs/)
- [Grafana 文档](https://grafana.com/docs/grafana/latest/)

---

## 🎯 扩展示例

### 添加 Jaeger

```yaml
services:
  jaeger:
    image: jaegertracing/all-in-one:1.61
    ports:
      - "16686:16686"  # UI
      - "14268:14268"  # HTTP
```

### 添加 Zipkin

```yaml
exporters:
  zipkin:
    endpoint: http://zipkin:9411/api/v2/spans
```

---

**状态**: ✅ 生产就绪
**版本**: 2024最新
**更新**: 2025-12-03
