
## 目标与范围

- 覆盖端到端可观测性四要素：Trace、Metrics、Logs、Profile，统一以 OpenTelemetry（以下简称 OTel）为采集与标准层；Collector 作为汇聚与处理中心；下游可插拔（Prometheus/Tempo/Loki/Mimir/Jaeger/OTel-Collector/云厂商）。
- 实施优先级：Trace → Metrics → Logs → Profile，分阶段灰度接入，控制采样与成本。
- 输出：最小可运行方案、环境化部署参考（本地/容器/K8s）、Go 接入规范、SLO 与告警基线。

## 总体架构

```text
应用（Go SDK: otel/sdk + instrumentation）
    ├─ Trace: OTLP/gRPC →
    ├─ Metrics: OTLP/gRPC →  OTel Collector（管道：接收 → 处理器 → 导出）
    └─ Logs: OTLP/gRPC   →
                                ├─ Tempo/Jaeger（Trace）
                                ├─ Prometheus/Mimir（Metrics，经 Prometheus Remote Write 或 OTLP → Prometheus 接收）
                                ├─ Loki/ELK（Logs）
                                └─ Pyroscope/Parca/云厂商（Profile）
```

- 采样与隐私：在 SDK 与 Collector 双层可控（先头部采样，后端再采样/尾部采样）；敏感字段在 Processor 阶段脱敏/丢弃。
- 统一上下文：W3C TraceContext（traceparent/tracestate）+ Baggage；跨服务一致传播。
- 稳态与容量：Collector 侧启用队列、重试、批处理；按 QPS 与基数(Cardinality)治理 Metrics 标签和值域。

## 分阶段落地计划

1) P0（2 周）最小闭环

   - 应用嵌入 OTel SDK（HTTP 客户端/服务端、数据库驱动基础埋点）
   - Collector 单副本 + Tempo + Prometheus + Loki（可使用 Grafana All-in-One 或简化 Compose）
   - 统一追踪 ID 注入日志；基础仪表盘（延迟、错误率、吞吐量）

2) P1（2-4 周）稳定与标准化

   - 采样策略：基于流量/端点/错误率的 head sampling；关键链路 100% 采样
   - 指标治理：RED/USE 指标基线，限制高基数标签
   - 日志治理：结构化日志（slog）统一字段；敏感字段脱敏
   - Collector 横向扩展；引入队列/重试/批处理与限流

3) P2（4-6 周）成本与可靠性

   - 尾部采样（tail sampling）按异常/延迟阈值提升采样
   - 异地/多环境路由；归档与保存期策略
   - SLO/告警与演练（看板驱动运维 SRE）

4) P3（长期）Profiling 与智能关联

   - 持续分析（Continuous Profiling）接入 Pyroscope/Parca
   - Trace ↔ Profile ↔ Logs ↔ Metrics 关联跳转与根因辅助定位

## 数据与处理规范

- 追踪命名：`service.name`、`service.version`、`deployment.environment` 必填；Span 名使用语义化（HTTP METHOD + 路径模板）
- 指标命名：采用 OTel 语义约定；限制标签数量与取值基数；长尾指标下沉为日志或事件
- 日志：优先结构化（Go 1.21+ `log/slog`），包含 `trace_id`、`span_id`、`severity`、`source`
- 脱敏与合规：在 Collector Processor 层做字段删除/哈希；遵循数据最小化原则

## 参考部署（本地 Docker Compose）

```yaml
version: '3.9'
services:
  collector:
    image: otel/opentelemetry-collector:0.104.0
    command: ["--config=/etc/otelcol/config.yaml"]
    volumes:
      - ./otelcol.yaml:/etc/otelcol/config.yaml:ro
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP

  tempo:
    image: grafana/tempo:2.5.0
    ports: ["3200:3200"]

  prometheus:
    image: prom/prometheus:v2.54.1
    ports: ["9090:9090"]

  loki:
    image: grafana/loki:3.1.0
    ports: ["3100:3100"]

  grafana:
    image: grafana/grafana:11.2.0
    ports: ["3000:3000"]
```

Collector 最小配置（`otelcol.yaml`）示例：

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:
processors:
  batch: {}
  memory_limiter: { check_interval: 5s, limit_percentage: 80, spike_limit_percentage: 25 }
exporters:
  otlp/tempo:
    endpoint: tempo:4317
    tls: { insecure: true }
  prometheus:
    endpoint: 0.0.0.0:8889
  loki:
    endpoint: http://loki:3100/loki/api/v1/push
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/tempo]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [prometheus]
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [loki]
```

## Go 应用最小接入规范

```go
// go.mod 需包含 go.opentelemetry.io/otel 相关依赖
import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitOTel(serviceName, env string) (func(context.Context) error, error) {
    res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName(serviceName),
        semconv.DeploymentEnvironment(env),
    ))
    exp, err := otlptracegrpc.New(context.Background())
    if err != nil { return nil, err }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
    return tp.Shutdown, nil
}
```

日志关联（Go 1.21+ `slog` 与 Trace 上下文）：

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
ctx, span := otel.Tracer("svc").Start(context.Background(), "op")
defer span.End()
traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
logger.Info("handling op", "trace_id", traceID)
```

## 采样与成本策略

- Head Sampling：按端点/租户/比例采样，默认 1-5%；关键路径提高至 100%
- Tail Sampling（可选）：对高延迟/错误的 Trace 提高采样（Collector Tail Sampling Processor）
- 指标：以 RED/USE 为核心，限制高基数标签；业务维度用 exemplar 关联 Trace
- 日志：仅在需要时上报详细参数；启用日志等级与采样
- 存储与保留：Trace/Loki 短期高保真（3-7 天），中长期降采样/归档

## SLO 与告警基线

- SLI：
  - 可用性：成功请求比例（2xx/3xx）
  - 延迟：P50/P90/P99
  - 错误率：5xx 比例
- SLO：
  - API：30 天内 99.9% 可用性；P99 < 500ms
- 告警：
  - 错误率瞬时 > 基线阈值（滑窗）
  - 延迟 P99 > 阈值 且 持续 5 分钟
  - 无数据/Exporter 失败/队列溢出

## 风险与缓解

- 高基数/高吞吐导致的存储与查询成本：预防性标签治理与采样
- 采样导致的排障遗漏：关键路径 100% + Tail 提升
- 多后端耦合：通过 Collector 实现解耦与路由

## 成功标准（验收）

- 本地/容器环境可在 30 分钟内部署成功并看到 Trace/Metrics/Logs
- 关键服务具备统一追踪 ID 串联日志与指标
- 发生错误或 P99 异常时可 5 分钟内定位到可疑代码段
