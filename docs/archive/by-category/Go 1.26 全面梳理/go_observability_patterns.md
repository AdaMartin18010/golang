# Go 1.23 可观测性基础设计模型全面梳理

> 本文档系统性地梳理了 **Go 1.23** 语言中OTLP、eBPF等可观测性基础设计模型，涵盖链路追踪、指标收集、日志聚合、性能剖析等核心领域。
>
> **Go 1.23 更新**：
>
> - `runtime/pprof` 最大栈深度从32提升至128帧
> - `runtime/trace` 崩溃时自动刷新追踪数据
> - 利用 `iter` 包优化日志和追踪数据的迭代处理
> - 使用 `unique` 包优化高频日志字段的内存使用

---

## 目录

- [Go 1.23 可观测性基础设计模型全面梳理](#go-123-可观测性基础设计模型全面梳理)
  - [目录](#目录)
  - [1. OpenTelemetry基础](#1-opentelemetry基础)
    - [1.1 概念定义](#11-概念定义)
    - [1.2 架构图](#12-架构图)
    - [1.3 三大支柱详解](#13-三大支柱详解)
      - [1.3.1 Trace（链路追踪）](#131-trace链路追踪)
      - [1.3.2 Metrics（指标）](#132-metrics指标)
      - [1.3.3 Logs（日志）](#133-logs日志)
    - [1.4 OTLP协议详解](#14-otlp协议详解)
    - [1.5 Collector架构](#15-collector架构)
    - [1.6 Go SDK使用](#16-go-sdk使用)
      - [1.6.1 完整示例](#161-完整示例)
      - [1.6.2 反例说明](#162-反例说明)
      - [1.6.3 最佳实践](#163-最佳实践)
  - [2. 链路追踪（Tracing）](#2-链路追踪tracing)
    - [2.1 概念定义](#21-概念定义)
    - [2.2 Span与Trace详细模型](#22-span与trace详细模型)
    - [2.3 上下文传播机制](#23-上下文传播机制)
    - [2.4 采样策略](#24-采样策略)
    - [2.5 Jaeger/Zipkin集成](#25-jaegerzipkin集成)
  - [3. 指标收集（Metrics）](#3-指标收集metrics)
    - [3.1 概念定义](#31-概念定义)
    - [3.2 指标类型详解](#32-指标类型详解)
    - [3.3 Prometheus集成](#33-prometheus集成)
    - [3.4 OpenTelemetry Metrics](#34-opentelemetry-metrics)
    - [3.5 反例说明](#35-反例说明)
    - [3.6 最佳实践](#36-最佳实践)
  - [4. 日志聚合（Logging）](#4-日志聚合logging)
    - [4.1 概念定义](#41-概念定义)
    - [4.2 结构化日志（zap/logrus）](#42-结构化日志zaplogrus)
    - [4.3 日志与Trace关联](#43-日志与trace关联)
    - [4.4 反例说明](#44-反例说明)
    - [4.5 最佳实践](#45-最佳实践)
  - [5. eBPF基础](#5-ebpf基础)
    - [5.1 概念定义](#51-概念定义)
    - [5.2 eBPF架构](#52-ebpf架构)
    - [5.3 eBPF程序生命周期](#53-ebpf程序生命周期)
    - [5.4 Go与eBPF交互（cilium/ebpf库）](#54-go与ebpf交互ciliumebpf库)
    - [5.5 eBPF C程序示例](#55-ebpf-c程序示例)
    - [5.6 性能剖析应用](#56-性能剖析应用)
    - [5.7 反例说明](#57-反例说明)
    - [5.8 最佳实践](#58-最佳实践)
  - [6. 性能剖析（Profiling）](#6-性能剖析profiling)
    - [6.1 概念定义](#61-概念定义)
    - [6.2 pprof架构](#62-pprof架构)
    - [6.3 CPU Profiling](#63-cpu-profiling)
    - [6.4 Memory Profiling](#64-memory-profiling)
    - [6.5 Goroutine Profiling](#65-goroutine-profiling)
    - [6.6 Block和Mutex Profiling](#66-block和mutex-profiling)
    - [6.7 HTTP pprof端点](#67-http-pprof端点)
    - [6.8 反例说明](#68-反例说明)
    - [6.9 最佳实践](#69-最佳实践)
  - [7. 健康检查](#7-健康检查)
    - [7.1 概念定义](#71-概念定义)
    - [7.2 健康检查架构](#72-健康检查架构)
    - [7.3 Go实现示例](#73-go实现示例)
    - [7.4 Kubernetes配置示例](#74-kubernetes配置示例)
    - [7.5 高级健康检查模式](#75-高级健康检查模式)
    - [7.6 反例说明](#76-反例说明)
    - [7.7 最佳实践](#77-最佳实践)
  - [8. 总结与最佳实践](#8-总结与最佳实践)
    - [8.1 可观测性三支柱整合](#81-可观测性三支柱整合)
    - [8.2 Go可观测性最佳实践清单](#82-go可观测性最佳实践清单)
      - [Trace最佳实践](#trace最佳实践)
      - [Metrics最佳实践](#metrics最佳实践)
      - [Logging最佳实践](#logging最佳实践)
      - [eBPF最佳实践](#ebpf最佳实践)
      - [Profiling最佳实践](#profiling最佳实践)
      - [Health Check最佳实践](#health-check最佳实践)
    - [8.3 推荐工具链](#83-推荐工具链)

---

## 1. OpenTelemetry基础

### 1.1 概念定义

**OpenTelemetry（OTel）** 是一个开源的可观测性框架，由Cloud Native Computing Foundation（CNCF）孵化，旨在提供标准化的遥测数据收集、处理和导出机制。
它是OpenTracing和OpenCensus的合并产物，现已成为云原生可观测性的事实标准。

**核心设计目标**：

- **vendor-neutral**：与特定后端解耦，支持多种导出目标
- **统一标准**：提供跨语言的统一API和SDK
- **自动 instrumentation**：最小化代码侵入性
- **高性能**：低开销的数据收集机制

### 1.2 架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         OpenTelemetry 架构                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                      │
│  │ Application │    │ Application │    │ Application │                      │
│  │   (Go)      │    │  (Python)   │    │  (Java)     │                      │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                      │
│         │                  │                  │                             │
│         ▼                  ▼                  ▼                             │
│  ┌─────────────────────────────────────────────────────┐                   │
│  │              OpenTelemetry API/SDK                   │                   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────┐  │                   │
│  │  │  Trace  │  │ Metrics │  │  Logs   │  │ Baggage│  │                   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └───┬────┘  │                   │
│  └───────┼────────────┼────────────┼───────────┼───────┘                   │
│          │            │            │           │                            │
│          ▼            ▼            ▼           ▼                            │
│  ┌─────────────────────────────────────────────────────┐                   │
│  │              OpenTelemetry Collector                 │                   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────┐  │                   │
│  │  │Receiver │→ │ Processor│→ │  Batch  │→ │Exporter│  │                   │
│  │  │ (OTLP)  │  │(Filter)  │  │ (Queue) │  │(Various)│  │                   │
│  │  └─────────┘  └─────────┘  └─────────┘  └────────┘  │                   │
│  └─────────────────────────────────────────────────────┘                   │
│          │            │            │           │                            │
│          ▼            ▼            ▼           ▼                            │
│  ┌─────────────┐  ┌──────────┐  ┌─────────┐  ┌──────────┐                  │
│  │   Jaeger    │  │Prometheus│  │  Loki   │  │  Tempo   │                  │
│  │  (Tracing)  │  │(Metrics) │  │ (Logs)  │  │(Tracing) │                  │
│  └─────────────┘  └──────────┘  └─────────┘  └──────────┘                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.3 三大支柱详解

#### 1.3.1 Trace（链路追踪）

**定义**：Trace记录请求在分布式系统中的完整调用路径，由多个Span组成的有向无环图（DAG）。

**核心概念**：

- **Trace**：一次完整的请求调用链，具有唯一的TraceID
- **Span**：Trace中的单个操作单元，包含操作名称、起止时间、属性、事件
- **SpanContext**：携带TraceID、SpanID、采样标志等传播信息

```
Trace (TraceID: abc123)
│
├── Span A (Root) [0ms - 100ms]
│   ├── Span B (Child) [10ms - 50ms]
│   │   └── Span D (Child) [15ms - 40ms]
│   └── Span C (Child) [60ms - 90ms]
│
└── Span E (FollowsFrom) [100ms - 150ms]
```

#### 1.3.2 Metrics（指标）

**定义**：Metrics是在一段时间内聚合的数值测量，用于监控系统的状态和性能。

**指标类型**：

| 类型 | 描述 | 适用场景 |
|------|------|----------|
| Counter | 单调递增的累计值 | 请求总数、错误数 |
| UpDownCounter | 可增可减的累计值 | 队列长度、连接数 |
| Gauge | 瞬时值 | 温度、内存使用量 |
| Histogram | 数值分布统计 | 请求延迟分布 |
| ObservableCounter | 异步观测的累计值 | CPU时间、网络IO |

#### 1.3.3 Logs（日志）

**定义**：Logs是离散的事件记录，包含时间戳和结构化/非结构化消息。

**OpenTelemetry日志模型**：

```
LogRecord {
    Timestamp: 2024-01-01T00:00:00Z
    Severity: INFO
    Body: "User login successful"
    Attributes: {
        user_id: "12345"
        ip_address: "192.168.1.1"
    }
    TraceID: abc123
    SpanID: def456
}
```

### 1.4 OTLP协议详解

**定义**：OpenTelemetry Protocol（OTLP）是OpenTelemetry定义的标准传输协议，用于遥测数据（Trace、Metrics、Logs）的传输。

**协议特性**：

- **传输层**：支持gRPC（默认）和HTTP/1.1（Protobuf/JSON）
- **数据编码**：Protocol Buffers（高效二进制）或JSON（可读性）
- **压缩**：支持GZIP等压缩算法
- **安全**：支持TLS/mTLS加密传输

**OTLP/gRPC消息结构**：

```protobuf
// trace.proto
message ExportTraceServiceRequest {
    repeated ResourceSpans resource_spans = 1;
}

message ResourceSpans {
    Resource resource = 1;
    repeated ScopeSpans scope_spans = 2;
}

message ScopeSpans {
    InstrumentationScope scope = 1;
    repeated Span spans = 2;
}
```

**OTLP架构**：

```
┌─────────────────────────────────────────────────────────────┐
│                      OTLP Protocol                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Client (SDK)                    Server (Collector)         │
│  ┌─────────────┐                ┌─────────────┐             │
│  │  Exporter   │ ──gRPC/HTTP──→ │  Receiver   │             │
│  │             │   Protobuf     │             │             │
│  │ - Batch     │   JSON         │ - Decode    │             │
│  │ - Compress  │   TLS          │ - Validate  │             │
│  │ - Retry     │                │ - Route     │             │
│  └─────────────┘                └─────────────┘             │
│                                                              │
│  关键特性：                                                   │
│  • 请求批处理（Batching）                                      │
│  • 超时与重试机制                                              │
│  • 背压处理（Backpressure）                                    │
│  • 元数据传递（Headers）                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.5 Collector架构

**定义**：OpenTelemetry Collector是一个可部署的代理/服务，负责接收、处理、导出遥测数据。

**Collector组件**：

```
┌─────────────────────────────────────────────────────────────────┐
│                  OpenTelemetry Collector                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  Receivers  │ →  │  Processors │ →  │  Exporters  │          │
│  │             │    │             │    │             │          │
│  │ • OTLP      │    │ • Batch     │    │ • OTLP      │          │
│  │ • Jaeger    │    │ • Memory    │    │ • Jaeger    │          │
│  │ • Prometheus│    │   Limiter   │    │ • Prometheus│          │
│  │ • Zipkin    │    │ • Filter    │    │ • Zipkin    │          │
│  │ • Kafka     │    │ • Resource  │    │ • Kafka     │          │
│  │             │    │   Processor │    │ • File      │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│         ↑                                    ↓                   │
│         └──────────── Connectors ────────────┘                   │
│                                                                  │
│  部署模式：                                                       │
│  • Agent模式：与应用同机部署，收集本机数据                          │
│  • Gateway模式：独立集群部署，集中处理数据                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.6 Go SDK使用

#### 1.6.1 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
 "go.opentelemetry.io/otel/sdk/resource"
 sdktrace "go.opentelemetry.io/otel/sdk/trace"
 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
 "go.opentelemetry.io/otel/trace"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials/insecure"
)

// 全局Tracer提供者
var tracer trace.Tracer

// initTracer 初始化Tracer提供者
func initTracer() (*sdktrace.TracerProvider, error) {
 ctx := context.Background()

 // 创建OTLP gRPC导出器
 conn, err := grpc.DialContext(ctx, "localhost:4317",
  grpc.WithTransportCredentials(insecure.NewCredentials()),
  grpc.WithBlock(),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
 }

 exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
 if err != nil {
  return nil, fmt.Errorf("failed to create trace exporter: %w", err)
 }

 // 创建资源（描述产生遥测数据的实体）
 res, err := resource.New(ctx,
  resource.WithAttributes(
   semconv.ServiceName("my-go-service"),
   semconv.ServiceVersion("1.0.0"),
   attribute.String("deployment.environment", "production"),
  ),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create resource: %w", err)
 }

 // 创建TracerProvider
 tp := sdktrace.NewTracerProvider(
  sdktrace.WithBatcher(exporter),           // 批处理导出
  sdktrace.WithResource(res),               // 关联资源
  sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样策略
 )

 // 设置为全局TracerProvider
 otel.SetTracerProvider(tp)
 tracer = tp.Tracer("my-go-service")

 return tp, nil
}

// processOrder 模拟订单处理
func processOrder(ctx context.Context, orderID string) error {
 // 创建子Span
 ctx, span := tracer.Start(ctx, "processOrder",
  trace.WithAttributes(
   attribute.String("order.id", orderID),
  ),
 )
 defer span.End()

 // 添加事件
 span.AddEvent("validating_order", trace.WithAttributes(
  attribute.String("validation.type", "inventory"),
 ))

 // 模拟处理
 time.Sleep(50 * time.Millisecond)

 // 验证库存
 if err := validateInventory(ctx, orderID); err != nil {
  span.RecordError(err)
  span.SetStatus(trace.Status{Code: trace.StatusCodeError, Description: err.Error()})
  return err
 }

 // 处理支付
 if err := processPayment(ctx, orderID); err != nil {
  span.RecordError(err)
  return err
 }

 span.SetAttributes(attribute.Bool("order.success", true))
 return nil
}

// validateInventory 验证库存
func validateInventory(ctx context.Context, orderID string) error {
 _, span := tracer.Start(ctx, "validateInventory")
 defer span.End()

 time.Sleep(20 * time.Millisecond)
 return nil
}

// processPayment 处理支付
func processPayment(ctx context.Context, orderID string) error {
 _, span := tracer.Start(ctx, "processPayment")
 defer span.End()

 time.Sleep(30 * time.Millisecond)
 return nil
}

func main() {
 // 初始化Tracer
 tp, err := initTracer()
 if err != nil {
  log.Fatalf("Failed to initialize tracer: %v", err)
 }
 defer func() {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if err := tp.Shutdown(ctx); err != nil {
   log.Printf("Error shutting down tracer provider: %v", err)
  }
 }()

 // 创建根Span
 ctx, span := tracer.Start(context.Background(), "main")
 defer span.End()

 // 处理订单
 if err := processOrder(ctx, "ORDER-12345"); err != nil {
  log.Printf("Failed to process order: %v", err)
 }

 fmt.Println("Trace completed")
}
```

#### 1.6.2 反例说明

```go
// ❌ 错误示例1：未正确传播上下文
func badProcessOrder(orderID string) {  // 缺少context.Context参数
 _, span := tracer.Start(context.Background(), "processOrder")  // 使用Background，丢失父Span关联
 defer span.End()
 // ...
}

// ❌ 错误示例2：忘记结束Span
func badProcessOrder2(ctx context.Context) {
 _, span := tracer.Start(ctx, "processOrder")
 // 忘记defer span.End() - 导致Span永不结束，内存泄漏
}

// ❌ 错误示例3：在循环中创建过多Span
func badProcessBatch(ctx context.Context, items []string) {
 for _, item := range items {
  _, span := tracer.Start(ctx, "processItem")  // 每个item都创建Span，开销过大
  // 应该使用单个Span记录批量操作
  span.End()
 }
}

// ❌ 错误示例4：未处理错误状态
func badProcessPayment(ctx context.Context) error {
 _, span := tracer.Start(ctx, "processPayment")
 defer span.End()

 err := doPayment()
 if err != nil {
  // 未记录错误到Span
  return err
 }
 return nil
}
```

#### 1.6.3 最佳实践

```go
// ✅ 正确示例1：上下文传播
func goodProcessOrder(ctx context.Context, orderID string) error {
 ctx, span := tracer.Start(ctx, "processOrder")  // 正确传播上下文
 defer span.End()
 // ...
}

// ✅ 正确示例2：批量操作使用单个Span
func goodProcessBatch(ctx context.Context, items []string) {
 ctx, span := tracer.Start(ctx, "processBatch",
  trace.WithAttributes(attribute.Int("batch.size", len(items))),
 )
 defer span.End()

 for _, item := range items {
  // 批量处理，不创建子Span
  processItem(item)
 }
}

// ✅ 正确示例3：正确记录错误
func goodProcessPayment(ctx context.Context) error {
 ctx, span := tracer.Start(ctx, "processPayment")
 defer span.End()

 err := doPayment()
 if err != nil {
  span.RecordError(err)  // 记录错误
  span.SetStatus(trace.Status{Code: trace.StatusCodeError, Description: err.Error()})
  return err
 }
 span.SetAttributes(attribute.Bool("payment.success", true))
 return nil
}

// ✅ 正确示例4：使用属性而非日志
func goodRecordMetrics(ctx context.Context) {
 ctx, span := tracer.Start(ctx, "recordMetrics")
 defer span.End()

 // 使用属性记录结构化数据
 span.SetAttributes(
  attribute.Int64("db.query_time_ms", 150),
  attribute.String("db.statement", "SELECT * FROM users"),
  attribute.Int("db.rows_affected", 100),
 )
}
```

---

## 2. 链路追踪（Tracing）

### 2.1 概念定义

**链路追踪（Distributed Tracing）** 是一种监控和诊断分布式系统中请求流的技术。它通过记录请求在各个服务间的调用路径，帮助开发者理解系统行为、定位性能瓶颈和排查故障。

**核心术语**：

| 术语 | 定义 | 类比 |
|------|------|------|
| Trace | 一次完整请求的调用链，由多个Span组成 | 一次旅行的完整行程 |
| Span | Trace中的单个操作单元，包含开始/结束时间、操作名称 | 旅行中的单个站点 |
| SpanContext | 携带TraceID、SpanID等传播信息 | 旅行护照 |
| Parent Span | 创建其他Span的Span | 父节点 |
| Child Span | 被其他Span创建的Span | 子节点 |
| Baggage | 跨Span传递的键值对数据 | 随身携带的行李 |

### 2.2 Span与Trace详细模型

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Span 数据结构                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Span {                                                                      │
│      // 标识信息                                                            │
│      TraceID:      16字节唯一标识符                                          │
│      SpanID:       8字节唯一标识符                                           │
│      ParentSpanID: 8字节（可选）                                             │
│      Name:         "HTTP GET /api/users"                                     │
│                                                                              │
│      // 时间信息                                                            │
│      StartTime:    2024-01-01T00:00:00.000Z                                  │
│      EndTime:      2024-01-01T00:00:00.150Z                                  │
│      Duration:     150ms                                                     │
│                                                                              │
│      // 属性（键值对）                                                       │
│      Attributes: {                                                           │
│          "http.method": "GET"                                                │
│          "http.url": "/api/users"                                            │
│          "http.status_code": 200                                             │
│          "user.id": "12345"                                                  │
│      }                                                                       │
│                                                                              │
│      // 事件（时间线标记）                                                    │
│      Events: [                                                               │
│          { Time: T+10ms,  Name: "cache_miss" }                               │
│          { Time: T+50ms,  Name: "db_query_start" }                           │
│          { Time: T+120ms, Name: "db_query_end" }                             │
│      ]                                                                       │
│                                                                              │
│      // 状态                                                                 │
│      Status: { Code: OK, Description: "" }                                   │
│                                                                              │
│      // 链接（关联其他Trace）                                                 │
│      Links: [                                                                │
│          { TraceID: XXX, SpanID: YYY, Attributes: {...} }                    │
│      ]                                                                       │
│  }                                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 上下文传播机制

**定义**：上下文传播（Context Propagation）是在分布式系统中传递Trace信息的技术，确保跨服务调用的Span能够正确关联。

**传播方式**：

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      上下文传播机制                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Service A                              Service B                           │
│  ┌─────────────┐                        ┌─────────────┐                     │
│  │ Span A      │ ──HTTP Request───────→ │ Span B      │                     │
│  │             │   Headers:             │ (Child of A)│                     │
│  │ traceparent:│   traceparent:         │             │                     │
│  │ 00-abc123-  │   00-abc123-def456-01  │             │                     │
│  │ def456-01   │                        │             │                     │
│  │             │   tracestate:          │             │                     │
│  │ tracestate: │   vendor=value         │             │                     │
│  │ vendor=val  │                        │             │                     │
│  │ baggage:    │   baggage:             │             │                     │
│  │ user=john   │   user=john            │             │                     │
│  └─────────────┘                        └─────────────┘                     │
│                                                                              │
│  W3C Trace Context标准：                                                     │
│  • traceparent: 00-<trace-id>-<parent-id>-<flags>                           │
│  • tracestate:  厂商特定的扩展信息                                           │
│  • baggage:     用户自定义的键值对                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Go实现示例**：

```go
package main

import (
 "context"
 "fmt"
 "net/http"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/baggage"
 "go.opentelemetry.io/otel/propagation"
 "go.opentelemetry.io/otel/trace"
)

// 初始化传播器
var propagator = propagation.NewCompositeTextMapPropagator(
 propagation.TraceContext{},  // W3C traceparent/tracestate
 propagation.Baggage{},       // W3C baggage
)

// HTTP客户端：注入上下文到请求头
func makeRequestWithContext(ctx context.Context, url string) (*http.Response, error) {
 req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
 if err != nil {
  return nil, err
 }

 // 注入Trace上下文到HTTP头
 propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

 client := &http.Client{}
 return client.Do(req)
}

// HTTP服务端：从请求头提取上下文
func handler(w http.ResponseWriter, r *http.Request) {
 // 提取上下文
 ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

 // 获取Baggage
 bag := baggage.FromContext(ctx)
 userMember := bag.Member("user.id")
 fmt.Printf("User ID from baggage: %s\n", userMember.Value())

 // 创建子Span
 tracer := otel.Tracer("http-server")
 ctx, span := tracer.Start(ctx, "handleRequest")
 defer span.End()

 // 处理请求...
 w.Write([]byte("OK"))
}

// 设置Baggage
func setBaggage(ctx context.Context, key, value string) context.Context {
 member, err := baggage.NewMember(key, value)
 if err != nil {
  return ctx
 }
 bag, err := baggage.New(member)
 if err != nil {
  return ctx
 }
 return baggage.ContextWithBaggage(ctx, bag)
}

func main() {
 // 示例：设置和传播Baggage
 ctx := context.Background()
 ctx = setBaggage(ctx, "user.id", "12345")
 ctx = setBaggage(ctx, "tenant.id", "acme-corp")

 // 创建Span，Baggage会自动传播
 tracer := otel.Tracer("example")
 ctx, span := tracer.Start(ctx, "main")
 defer span.End()

 // 模拟HTTP调用
 // resp, err := makeRequestWithContext(ctx, "http://service-b/api")
}
```

### 2.4 采样策略

**定义**：采样（Sampling）决定哪些Trace应该被记录和导出，用于控制数据量和成本。

**采样类型**：

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        采样策略对比                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. AlwaysOnSampler (全采样)                                                 │
│     ├── 特点：记录所有Trace                                                   │
│     ├── 适用：开发环境、低流量系统                                            │
│     └── 缺点：高流量时数据量过大                                              │
│                                                                              │
│  2. AlwaysOffSampler (不采样)                                                │
│     ├── 特点：不记录任何Trace                                                 │
│     ├── 适用：完全禁用追踪                                                    │
│     └── 缺点：无可见性                                                        │
│                                                                              │
│  3. TraceIDRatioBased (概率采样)                                             │
│     ├── 特点：按TraceID哈希值比例采样                                         │
│     ├── 示例：采样率0.1表示10%的Trace被记录                                   │
│     ├── 优点：确定性采样（相同TraceID总是相同结果）                            │
│     └── 适用：高流量生产环境                                                  │
│                                                                              │
│  4. ParentBased (基于父Span)                                                 │
│     ├── 特点：根据父Span的采样决策                                            │
│     ├── 配置：                                                                │
│     │   • root: 根Span采样策略                                                │
│     │   • remoteParentSampled: 远程父Span已采样                               │
│     │   • remoteParentNotSampled: 远程父Span未采样                            │
│     │   • localParentSampled: 本地父Span已采样                                │
│     │   • localParentNotSampled: 本地父Span未采样                             │
│     └── 适用：保持Trace完整性                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Go实现示例**：

```go
package main

import (
 "context"
 "fmt"

 sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// 配置采样策略
func configureSamplers() {
 // 1. 全采样 - 开发环境
 alwaysOn := sdktrace.AlwaysSample()

 // 2. 概率采样 - 10%采样率，生产环境
 ratioBased := sdktrace.TraceIDRatioBased(0.1)

 // 3. 基于父Span的采样 - 复杂配置
 parentBased := sdktrace.ParentBased(
  sdktrace.TraceIDRatioBased(0.1), // 根Span使用10%采样
  sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
  sdktrace.WithRemoteParentNotSampled(sdktrace.NeverSample()),
  sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
  sdktrace.WithLocalParentNotSampled(sdktrace.NeverSample()),
 )

 // 4. 自定义采样器 - 基于属性过滤
 customSampler := &attributeBasedSampler{
  attributeKey:   "http.url",
  attributeValue: "/health",
  delegate:       sdktrace.NeverSample(), // /health端点不采样
 }

 // 创建TracerProvider
 _ = sdktrace.NewTracerProvider(
  sdktrace.WithSampler(parentBased),
 )

 fmt.Printf("Samplers configured: %+v\n", []sdktrace.Sampler{
  alwaysOn, ratioBased, parentBased, customSampler,
 })
}

// 自定义采样器：基于属性过滤
type attributeBasedSampler struct {
 attributeKey   string
 attributeValue string
 delegate       sdktrace.Sampler
}

func (s *attributeBasedSampler) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
 // 检查属性
 for _, attr := range params.Attributes {
  if string(attr.Key) == s.attributeKey && attr.Value.AsString() == s.attributeValue {
   return s.delegate.ShouldSample(params)
  }
 }
 // 默认采样
 return sdktrace.AlwaysSample().ShouldSample(params)
}

func (s *attributeBasedSampler) Description() string {
 return fmt.Sprintf("AttributeBasedSampler{%s=%s}", s.attributeKey, s.attributeValue)
}

func main() {
 configureSamplers()
}
```

### 2.5 Jaeger/Zipkin集成

**Jaeger集成示例**：

```go
package main

import (
 "context"
 "fmt"
 "io"
 "log"
 "time"

 "github.com/jaegertracing/jaeger-client-go"
 jaegercfg "github.com/jaegertracing/jaeger-client-go/config"
 "github.com/opentracing/opentracing-go"
 "github.com/opentracing/opentracing-go/ext"
)

// initJaeger 初始化Jaeger Tracer
func initJaeger(serviceName string) (opentracing.Tracer, io.Closer, error) {
 cfg := &jaegercfg.Configuration{
  ServiceName: serviceName,
  Sampler: &jaegercfg.SamplerConfig{
   Type:  jaeger.SamplerTypeConst,
   Param: 1, // 全采样
  },
  Reporter: &jaegercfg.ReporterConfig{
   LogSpans:           true,
   LocalAgentHostPort: "localhost:6831",
  },
 }

 tracer, closer, err := cfg.NewTracer()
 if err != nil {
  return nil, nil, fmt.Errorf("failed to create Jaeger tracer: %w", err)
 }

 opentracing.SetGlobalTracer(tracer)
 return tracer, closer, nil
}

// 模拟数据库操作
func queryDatabase(ctx context.Context, userID string) error {
 span, ctx := opentracing.StartSpanFromContext(ctx, "queryDatabase")
 defer span.Finish()

 span.SetTag("db.user_id", userID)
 span.SetTag("db.table", "users")

 // 模拟查询延迟
 time.Sleep(20 * time.Millisecond)

 span.LogKV("event", "query_completed")
 return nil
}

// HTTP处理函数
func handleUserRequest(ctx context.Context, userID string) error {
 span, ctx := opentracing.StartSpanFromContext(ctx, "handleUserRequest")
 defer span.Finish()

 ext.HTTPMethod.Set(span, "GET")
 ext.HTTPUrl.Set(span, "/api/users/"+userID)

 // 查询数据库
 if err := queryDatabase(ctx, userID); err != nil {
  ext.Error.Set(span, true)
  span.SetTag("error.message", err.Error())
  return err
 }

 ext.HTTPStatusCode.Set(span, 200)
 return nil
}

func main() {
 tracer, closer, err := initJaeger("my-service")
 if err != nil {
  log.Fatalf("Failed to initialize Jaeger: %v", err)
 }
 defer closer.Close()

 // 创建根Span
 span := tracer.StartSpan("main-operation")
 ctx := opentracing.ContextWithSpan(context.Background(), span)

 // 处理请求
 if err := handleUserRequest(ctx, "12345"); err != nil {
  log.Printf("Request failed: %v", err)
 }

 span.Finish()

 // 等待上报
 time.Sleep(2 * time.Second)
 fmt.Println("Trace sent to Jaeger")
}
```

---

## 3. 指标收集（Metrics）

### 3.1 概念定义

**指标（Metrics）** 是在一段时间内聚合的数值测量，用于监控系统的状态和性能。与Trace记录单个请求不同，Metrics关注的是系统整体行为的统计特征。

**核心特性**：

- **聚合性**：多个测量值合并为统计值
- **时间序列**：按时间顺序记录的数据点
- **低开销**：相比Trace，Metrics收集成本更低
- **可查询**：支持复杂的查询和聚合操作

### 3.2 指标类型详解

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         指标类型对比                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Counter（计数器）                                                        │
│     ┌─────────────────────────────────────────┐                             │
│     │ 值                                      │                             │
│     │ 10 ┤                                    │ 单调递增                    │
│     │  8 ┤          ┌───┐                     │ 只能增加                    │
│     │  6 ┤    ┌───┐ │   │                     │ 适用：请求数、错误数         │
│     │  4 ┤    │   │ │   │                     │                             │
│     │  2 ┤┌───┤   │ │   │                     │                             │
│     │  0 ┼────┴───┴─┴───┴────→ 时间           │                             │
│     └─────────────────────────────────────────┘                             │
│                                                                              │
│  2. Gauge（仪表盘）                                                          │
│     ┌─────────────────────────────────────────┐                             │
│     │ 值                                      │                             │
│     │ 10 ┤      ┌───┐                         │ 可增可减                    │
│     │  8 ┤      │   │    ┌───┐                │ 瞬时值                      │
│     │  6 ┤  ┌───┤   │    │   │                │ 适用：内存、温度、队列长度   │
│     │  4 ┤  │   │   │┌───┤   │                │                             │
│     │  2 ┤┌─┤   │   ││   │   │                │                             │
│     │  0 ┼─┴───┴───┴┴───┴───┴────→ 时间       │                             │
│     └─────────────────────────────────────────┘                             │
│                                                                              │
│  3. Histogram（直方图）                                                      │
│     ┌─────────────────────────────────────────┐                             │
│     │ 频率                                    │                             │
│     │    ┤  ██                               │ 数值分布                    │
│     │    ┤ ████                              │ 预定义桶边界                 │
│     │    ┤██████  ██                          │ 适用：延迟分布               │
│     │    ┤████████████                        │                             │
│     │    ┼────┬────┬────┬────→ 数值范围      │                             │
│     │      0-10  10-50  50-100                │                             │
│     └─────────────────────────────────────────┘                             │
│                                                                              │
│  4. Summary（摘要）                                                          │
│     ┌─────────────────────────────────────────┐                             │
│     │ 值                                      │                             │
│     │    ┤──────────P99──────────┐            │ 分位数计算                  │
│     │    ┤────────P95────────┐   │            │ 客户端计算                  │
│     │    ┤──────P90──────┐   │   │            │ 适用：SLA监控               │
│     │    ┤              │   │   │            │                             │
│     │    ┼──────────────┴───┴───┴────→ 时间  │                             │
│     └─────────────────────────────────────────┘                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Prometheus集成

**Prometheus** 是最流行的开源监控和告警工具，采用拉取（Pull）模型收集指标。

**数据模型**：

```
# 指标名称 + 标签 = 时间序列
http_requests_total{method="GET", endpoint="/api/users", status="200"}

# 四种指标类型
# 1. Counter - 累计值
http_requests_total 1027

# 2. Gauge - 瞬时值
memory_usage_bytes 536870912

# 3. Histogram - 分布统计
http_request_duration_seconds_bucket{le="0.1"} 240
http_request_duration_seconds_bucket{le="0.5"} 983
http_request_duration_seconds_bucket{le="1.0"} 1027
http_request_duration_seconds_sum 534.2
http_request_duration_seconds_count 1027

# 4. Summary - 分位数
http_request_duration_seconds{quantile="0.5"} 0.0234
http_request_duration_seconds{quantile="0.9"} 0.1456
http_request_duration_seconds{quantile="0.99"} 0.5678
```

**Go实现示例**：

```go
package main

import (
 "fmt"
 "math/rand"
 "net/http"
 "time"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 定义指标变量
var (
 // Counter: 请求总数
 httpRequestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total number of HTTP requests",
  },
  []string{"method", "endpoint", "status"},
 )

 // Counter: 错误总数
 httpErrorsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_errors_total",
   Help: "Total number of HTTP errors",
  },
  []string{"method", "endpoint", "error_type"},
 )

 // Gauge: 活跃连接数
 activeConnections = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "active_connections",
   Help: "Number of active connections",
  },
 )

 // Gauge: 内存使用量
 memoryUsage = promauto.NewGaugeFunc(
  prometheus.GaugeOpts{
   Name: "memory_usage_bytes",
   Help: "Current memory usage in bytes",
  },
  func() float64 {
   // 模拟内存使用
   return float64(rand.Intn(1000000000))
  },
 )

 // Histogram: 请求延迟分布
 httpRequestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "HTTP request duration in seconds",
   Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
  },
  []string{"method", "endpoint"},
 )

 // Summary: 请求延迟分位数
 httpRequestDurationSummary = promauto.NewSummaryVec(
  prometheus.SummaryOpts{
   Name:       "http_request_duration_summary_seconds",
   Help:       "HTTP request duration summary in seconds",
   Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
  },
  []string{"method", "endpoint"},
 )
)

// 记录HTTP请求
func recordHTTPRequest(method, endpoint string, duration time.Duration, statusCode int) {
 status := fmt.Sprintf("%d", statusCode)
 httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
 httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
 httpRequestDurationSummary.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// 记录HTTP错误
func recordHTTPError(method, endpoint, errorType string) {
 httpErrorsTotal.WithLabelValues(method, endpoint, errorType).Inc()
}

// 模拟HTTP处理
func handleRequest(w http.ResponseWriter, r *http.Request) {
 start := time.Now()

 // 增加活跃连接数
 activeConnections.Inc()
 defer activeConnections.Dec()

 // 模拟处理时间
 duration := time.Duration(rand.Intn(100)) * time.Millisecond
 time.Sleep(duration)

 // 随机返回状态码
 statusCode := 200
 if rand.Float32() < 0.1 {
  statusCode = 500
  recordHTTPError(r.Method, r.URL.Path, "internal_error")
 }

 w.WriteHeader(statusCode)
 w.Write([]byte("OK"))

 recordHTTPRequest(r.Method, r.URL.Path, time.Since(start), statusCode)
}

func main() {
 // 注册指标端点
 http.Handle("/metrics", promhttp.Handler())
 http.HandleFunc("/api/users", handleRequest)
 http.HandleFunc("/api/orders", handleRequest)

 fmt.Println("Server starting on :8080")
 fmt.Println("Metrics available at http://localhost:8080/metrics")
 http.ListenAndServe(":8080", nil)
}
```

### 3.4 OpenTelemetry Metrics

**OpenTelemetry Metrics** 提供统一的指标API，支持多种导出后端。

**Go实现示例**：

```go
package main

import (
 "context"
 "fmt"
 "log"
 "math/rand"
 "time"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
 "go.opentelemetry.io/otel/metric"
 sdkmetric "go.opentelemetry.io/otel/sdk/metric"
 "go.opentelemetry.io/otel/sdk/resource"
 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials/insecure"
)

// 全局Meter
var meter metric.Meter

// initMetrics 初始化Metrics提供者
func initMetrics() (*sdkmetric.MeterProvider, error) {
 ctx := context.Background()

 // 创建OTLP导出器
 conn, err := grpc.DialContext(ctx, "localhost:4317",
  grpc.WithTransportCredentials(insecure.NewCredentials()),
  grpc.WithBlock(),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
 }

 exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
 if err != nil {
  return nil, fmt.Errorf("failed to create metric exporter: %w", err)
 }

 // 创建资源
 res, err := resource.New(ctx,
  resource.WithAttributes(
   semconv.ServiceName("metrics-service"),
   semconv.ServiceVersion("1.0.0"),
  ),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create resource: %w", err)
 }

 // 创建MeterProvider
 mp := sdkmetric.NewMeterProvider(
  sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
  sdkmetric.WithResource(res),
 )

 otel.SetMeterProvider(mp)
 meter = mp.Meter("metrics-service")

 return mp, nil
}

// MetricsDemo 演示各种指标类型
func MetricsDemo() {
 ctx := context.Background()

 // 1. Counter - 请求计数
 requestCounter, err := meter.Int64Counter(
  "http_requests_total",
  metric.WithDescription("Total HTTP requests"),
 )
 if err != nil {
  log.Printf("Failed to create counter: %v", err)
  return
 }

 // 2. UpDownCounter - 活跃连接数
 activeConnections, err := meter.Int64UpDownCounter(
  "active_connections",
  metric.WithDescription("Number of active connections"),
 )
 if err != nil {
  log.Printf("Failed to create updown counter: %v", err)
  return
 }

 // 3. Histogram - 请求延迟
 requestDuration, err := meter.Float64Histogram(
  "http_request_duration_seconds",
  metric.WithDescription("HTTP request duration"),
  metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
 )
 if err != nil {
  log.Printf("Failed to create histogram: %v", err)
  return
 }

 // 4. ObservableGauge - 内存使用量
 _, err = meter.Float64ObservableGauge(
  "memory_usage_bytes",
  metric.WithDescription("Memory usage in bytes"),
  metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
   // 模拟获取内存使用
   memory := float64(rand.Intn(1000000000))
   o.Observe(memory)
   return nil
  }),
 )
 if err != nil {
  log.Printf("Failed to create observable gauge: %v", err)
  return
 }

 // 模拟请求处理
 for i := 0; i < 100; i++ {
  // 增加活跃连接
  activeConnections.Add(ctx, 1)

  // 模拟处理时间
  duration := time.Duration(rand.Intn(100)) * time.Millisecond
  time.Sleep(duration)

  // 记录指标
  attrs := attribute.NewSet(
   attribute.String("method", "GET"),
   attribute.String("endpoint", "/api/users"),
   attribute.Int("status", 200),
  )
  requestCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
  requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributeSet(attrs))

  // 减少活跃连接
  activeConnections.Add(ctx, -1)
 }
}

func main() {
 mp, err := initMetrics()
 if err != nil {
  log.Fatalf("Failed to initialize metrics: %v", err)
 }
 defer func() {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if err := mp.Shutdown(ctx); err != nil {
   log.Printf("Error shutting down meter provider: %v", err)
  }
 }()

 // 运行演示
 MetricsDemo()

 // 等待导出
 time.Sleep(5 * time.Second)
 fmt.Println("Metrics demo completed")
}
```

### 3.5 反例说明

```go
// ❌ 错误示例1：标签值过多导致基数爆炸
func badHighCardinality() {
 // 使用用户ID作为标签 - 会导致时间序列数量爆炸
 httpRequestsTotal.WithLabelValues("GET", "/api/users", userID).Inc()
 // 正确做法：使用有限的、预定义的标签值
}

// ❌ 错误示例2：Counter递减
func badCounterDecrement() {
 // Counter只能递增，递减是错误的
 httpRequestsTotal.WithLabelValues("GET", "/api/users", "200").Add(-1)
 // 正确做法：使用Gauge或UpDownCounter
}

// ❌ 错误示例3：忘记注册指标
func badUnregisteredMetric() {
 counter := prometheus.NewCounter(prometheus.CounterOpts{
  Name: "unregistered_counter",
 })
 // 忘记调用prometheus.MustRegister(counter)
 counter.Inc() // 不会生效
}

// ❌ 错误示例4：Histogram桶边界不合理
func badHistogramBuckets() {
 // 桶边界设置不合理
 prometheus.NewHistogram(prometheus.HistogramOpts{
  Name:    "bad_histogram",
  Buckets: []float64{1, 2, 3, 4, 5}, // 对于秒级延迟来说太小
 })
 // 正确做法：根据实际值范围设置合理的桶边界
}

// ❌ 错误示例5：在热路径创建指标
func badCreateInHotPath() {
 for i := 0; i < 1000000; i++ {
  // 每次迭代都创建新的指标 - 性能灾难
  counter := prometheus.NewCounter(...)
  counter.Inc()
 }
}
```

### 3.6 最佳实践

```go
// ✅ 正确示例1：预定义标签值
var validStatuses = []string{"200", "400", "500"}

func goodRecordRequest(status string) {
 // 验证标签值
 for _, valid := range validStatuses {
  if status == valid {
   httpRequestsTotal.WithLabelValues("GET", "/api/users", status).Inc()
   return
  }
 }
 // 未知状态归类为"other"
 httpRequestsTotal.WithLabelValues("GET", "/api/users", "other").Inc()
}

// ✅ 正确示例2：使用常量标签
var (
 // 在初始化时定义标签
 requestLabels = prometheus.Labels{
  "service": "user-service",
  "version": "1.0.0",
 }
)

// ✅ 正确示例3：合理的Histogram桶
var goodHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
 Name:    "http_request_duration_seconds",
 Help:    "HTTP request duration",
 Buckets: prometheus.DefBuckets, // 使用默认桶: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
})

// ✅ 正确示例4：指标缓存
var (
 // 预创建带标签的指标
 userAPIRequests = httpRequestsTotal.MustCurryWith(prometheus.Labels{
  "endpoint": "/api/users",
 })
)

func goodRecordUserAPI(method, status string) {
 // 直接使用缓存的指标
 userAPIRequests.WithLabelValues(method, status).Inc()
}

// ✅ 正确示例5：使用中间件模式
func metricsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  start := time.Now()

  // 包装ResponseWriter以捕获状态码
  wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

  next.ServeHTTP(wrapped, r)

  duration := time.Since(start)
  recordHTTPRequest(r.Method, r.URL.Path, duration, wrapped.statusCode)
 })
}

type responseWriter struct {
 http.ResponseWriter
 statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
 rw.statusCode = code
 rw.ResponseWriter.WriteHeader(code)
}
```

---

## 4. 日志聚合（Logging）

### 4.1 概念定义

**日志聚合（Log Aggregation）** 是将分散在各处的日志数据收集、存储、索引和分析的过程。在分布式系统中，日志聚合是故障排查和审计追踪的关键手段。

**日志演进**：

```
传统日志 → 结构化日志 → 与Trace关联 → 可观测性统一
   │           │            │              │
   ▼           ▼            ▼              ▼
纯文本      JSON格式     TraceID注入    OTLP统一导出
难解析      可查询       跨系统关联      与Metrics/Trace统一
```

### 4.2 结构化日志（zap/logrus）

**zap** 是Uber开发的高性能结构化日志库，采用零分配设计。

**Go实现示例**：

```go
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "time"

 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
 "go.uber.org/zap/zapcore"
)

// Logger 封装zap logger
type Logger struct {
 *zap.Logger
}

// NewLogger 创建新的Logger
func NewLogger() (*Logger, error) {
 config := zap.NewProductionConfig()
 config.EncoderConfig.TimeKey = "timestamp"
 config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
 config.EncoderConfig.StacktraceKey = "stacktrace"

 logger, err := config.Build()
 if err != nil {
  return nil, err
 }

 return &Logger{logger}, nil
}

// WithContext 添加上下文信息
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
 fields := []zap.Field{}

 // 提取Trace信息
 span := trace.SpanFromContext(ctx)
 if span != nil {
  spanContext := span.SpanContext()
  if spanContext.IsValid() {
   fields = append(fields,
    zap.String("trace_id", spanContext.TraceID().String()),
    zap.String("span_id", spanContext.SpanID().String()),
    zap.Bool("trace_sampled", spanContext.IsSampled()),
   )
  }
 }

 return l.With(fields...)
}

// Info 记录Info级别日志
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
 l.WithContext(ctx).Info(msg, fields...)
}

// Error 记录Error级别日志
func (l *Logger) Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
 allFields := append([]zap.Field{zap.Error(err)}, fields...)
 l.WithContext(ctx).Error(msg, allFields...)
}

// Warn 记录Warn级别日志
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
 l.WithContext(ctx).Warn(msg, fields...)
}

// Debug 记录Debug级别日志
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
 l.WithContext(ctx).Debug(msg, fields...)
}

// 示例：业务逻辑中使用日志
func processOrder(ctx context.Context, logger *Logger, orderID string) error {
 logger.Info(ctx, "开始处理订单",
  zap.String("order_id", orderID),
  zap.String("event", "order_processing_started"),
 )

 // 模拟处理
 time.Sleep(100 * time.Millisecond)

 // 记录成功
 logger.Info(ctx, "订单处理完成",
  zap.String("order_id", orderID),
  zap.String("event", "order_processing_completed"),
  zap.Duration("processing_time", 100*time.Millisecond),
 )

 return nil
}

// 示例：记录错误
func processPayment(ctx context.Context, logger *Logger, paymentID string) error {
 logger.Info(ctx, "开始处理支付",
  zap.String("payment_id", paymentID),
 )

 // 模拟错误
 err := fmt.Errorf("payment gateway timeout")
 logger.Error(ctx, "支付处理失败", err,
  zap.String("payment_id", paymentID),
  zap.String("error_type", "gateway_timeout"),
  zap.Int("retry_count", 3),
 )

 return err
}

// 自定义Encoder示例
func customEncoderExample() {
 encoderConfig := zapcore.EncoderConfig{
  TimeKey:        "ts",
  LevelKey:       "level",
  NameKey:        "logger",
  CallerKey:      "caller",
  FunctionKey:    zapcore.OmitKey,
  MessageKey:     "msg",
  StacktraceKey:  "stacktrace",
  LineEnding:     zapcore.DefaultLineEnding,
  EncodeLevel:    zapcore.LowercaseLevelEncoder,
  EncodeTime:     zapcore.EpochMillisTimeEncoder,
  EncodeDuration: zapcore.SecondsDurationEncoder,
  EncodeCaller:   zapcore.ShortCallerEncoder,
 }

 core := zapcore.NewCore(
  zapcore.NewJSONEncoder(encoderConfig),
  zapcore.AddSync(&customSyncer{}),
  zapcore.InfoLevel,
 )

 logger := zap.New(core)
 logger.Info("custom encoder example", zap.String("key", "value"))
}

// 自定义Syncer
type customSyncer struct{}

func (c *customSyncer) Write(p []byte) (n int, err error) {
 // 自定义处理逻辑，如发送到远程日志服务
 fmt.Printf("Custom syncer received: %s", string(p))
 return len(p), nil
}

func (c *customSyncer) Sync() error {
 return nil
}

func main() {
 logger, err := NewLogger()
 if err != nil {
  panic(err)
 }
 defer logger.Sync()

 ctx := context.Background()

 // 记录业务日志
 if err := processOrder(ctx, logger, "ORDER-12345"); err != nil {
  fmt.Printf("Process order failed: %v\n", err)
 }

 if err := processPayment(ctx, logger, "PAY-67890"); err != nil {
  fmt.Printf("Process payment failed: %v\n", err)
 }

 // 结构化日志输出示例
 logEntry := map[string]interface{}{
  "timestamp": time.Now().UTC().Format(time.RFC3339),
  "level":     "INFO",
  "message":   "Application started",
  "service":   "user-service",
  "version":   "1.0.0",
  "metadata": map[string]interface{}{
   "pid":     12345,
   "host":    "server-01",
   "region":  "us-east-1",
  },
 }
 jsonBytes, _ := json.MarshalIndent(logEntry, "", "  ")
 fmt.Printf("\nStructured log example:\n%s\n", string(jsonBytes))
}
```

### 4.3 日志与Trace关联

**核心概念**：将日志与Trace关联，可以在排查问题时快速定位相关调用链。

```go
package main

import (
 "context"
 "fmt"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
 "go.opentelemetry.io/otel/sdk/resource"
 sdktrace "go.opentelemetry.io/otel/sdk/trace"
 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
)

// TraceLogger 集成Trace和日志
type TraceLogger struct {
 logger *zap.Logger
 tracer trace.Tracer
}

// NewTraceLogger 创建TraceLogger
func NewTraceLogger() (*TraceLogger, error) {
 // 初始化Tracer
 exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
 if err != nil {
  return nil, err
 }

 res, _ := resource.New(context.Background(),
  resource.WithAttributes(
   semconv.ServiceName("trace-logger-service"),
  ),
 )

 tp := sdktrace.NewTracerProvider(
  sdktrace.WithBatcher(exporter),
  sdktrace.WithResource(res),
 )
 otel.SetTracerProvider(tp)

 // 初始化Logger
 logger, err := zap.NewProduction()
 if err != nil {
  return nil, err
 }

 return &TraceLogger{
  logger: logger,
  tracer: tp.Tracer("trace-logger-service"),
 }, nil
}

// LogWithSpan 记录带Trace信息的日志
func (tl *TraceLogger) LogWithSpan(ctx context.Context, level, msg string, fields ...zap.Field) {
 span := trace.SpanFromContext(ctx)
 if span != nil {
  spanContext := span.SpanContext()
  if spanContext.IsValid() {
   fields = append(fields,
    zap.String("trace_id", spanContext.TraceID().String()),
    zap.String("span_id", spanContext.SpanID().String()),
   )
  }
 }

 switch level {
 case "info":
  tl.logger.Info(msg, fields...)
 case "error":
  tl.logger.Error(msg, fields...)
 case "warn":
  tl.logger.Warn(msg, fields...)
 case "debug":
  tl.logger.Debug(msg, fields...)
 }
}

// StartSpan 开始新的Span并记录日志
func (tl *TraceLogger) StartSpan(ctx context.Context, name string, fields ...zap.Field) (context.Context, trace.Span) {
 ctx, span := tl.tracer.Start(ctx, name)

 // 记录Span开始日志
 tlc.LogWithSpan(ctx, "info", fmt.Sprintf("Span started: %s", name),
  append(fields, zap.String("span_name", name))...,
 )

 return ctx, span
}

// EndSpan 结束Span并记录日志
func (tl *TraceLogger) EndSpan(ctx context.Context, span trace.Span, fields ...zap.Field) {
 // 记录Span结束日志
 tlc.LogWithSpan(ctx, "info", fmt.Sprintf("Span ended: %s", span.SpanContext().SpanID().String()),
  append(fields, zap.String("span_id", span.SpanContext().SpanID().String()))...,
 )
 span.End()
}

// LogToSpan 将日志同时记录到Span和日志系统
func (tl *TraceLogger) LogToSpan(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
 span := trace.SpanFromContext(ctx)
 if span != nil {
  // 添加到Span事件
  span.AddEvent(msg, trace.WithAttributes(attrs...))
 }

 // 同时记录到日志
 fields := make([]zap.Field, 0, len(attrs))
 for _, attr := range attrs {
  fields = append(fields, zap.Any(string(attr.Key), attr.Value))
 }
 tlc.LogWithSpan(ctx, "info", msg, fields...)
}

var tlc *TraceLogger

func init() {
 var err error
 tlc, err = NewTraceLogger()
 if err != nil {
  panic(err)
 }
}

// 业务函数示例
func processRequest(ctx context.Context, requestID string) error {
 ctx, span := tlc.StartSpan(ctx, "processRequest",
  zap.String("request_id", requestID),
 )
 defer tlc.EndSpan(ctx, span)

 // 记录业务事件到Span和日志
 tlc.LogToSpan(ctx, "validating_request",
  attribute.String("request.id", requestID),
  attribute.String("validation.type", "auth"),
 )

 // 模拟验证
 // ...

 tlc.LogToSpan(ctx, "request_validated",
  attribute.Bool("validation.success", true),
 )

 // 记录数据库查询
 tlc.LogToSpan(ctx, "database_query",
  attribute.String("db.statement", "SELECT * FROM users WHERE id = ?"),
  attribute.String("db.system", "postgresql"),
 )

 return nil
}

func main() {
 ctx := context.Background()

 // 创建根Span
 ctx, span := tlc.StartSpan(ctx, "main-operation")
 defer tlc.EndSpan(ctx, span)

 // 处理请求
 if err := processRequest(ctx, "REQ-12345"); err != nil {
  tlc.LogWithSpan(ctx, "error", "Request processing failed",
   zap.Error(err),
  )
 }

 fmt.Println("\n日志和Trace关联完成")
}
```

### 4.4 反例说明

```go
// ❌ 错误示例1：使用fmt.Printf记录日志
func badPrintfLogging() {
 fmt.Printf("User %s logged in at %s\n", userID, time.Now())
 // 问题：非结构化，难以解析和查询
}

// ❌ 错误示例2：日志级别混乱
func badLogLevel() {
 // 错误级别使用不当
 logger.Info("Database connection failed") // 应该是Error级别
 logger.Error("Application started")       // 应该是Info级别
}

// ❌ 错误示例3：敏感信息泄露
func badSensitiveData() {
 logger.Info("User login",
  zap.String("password", userPassword), // ❌ 泄露密码
  zap.String("credit_card", cardNumber), // ❌ 泄露信用卡号
 )
}

// ❌ 错误示例4：在热路径创建logger
func badCreateLoggerInHotPath() {
 for i := 0; i < 1000000; i++ {
  logger, _ := zap.NewProduction() // 每次迭代都创建新logger
  logger.Info("processing item")
 }
}

// ❌ 错误示例5：忽略日志错误
func badIgnoreError() {
 logger, err := zap.NewProduction()
 _ = err // 忽略错误
 logger.Info("this might panic if logger is nil")
}

// ❌ 错误示例6：不关联Trace信息
func badNoTraceContext() {
 // 日志中没有TraceID，无法与分布式追踪关联
 logger.Info("processing request", zap.String("request_id", reqID))
}
```

### 4.5 最佳实践

```go
// ✅ 正确示例1：使用结构化日志
func goodStructuredLogging() {
 logger.Info("user_logged_in",
  zap.String("user_id", userID),
  zap.String("ip_address", clientIP),
  zap.Time("login_time", time.Now()),
  zap.String("user_agent", userAgent),
 )
}

// ✅ 正确示例2：正确的日志级别
func goodLogLevel() {
 logger.Debug("entering function", zap.String("function", "processOrder"))
 logger.Info("order processed", zap.String("order_id", orderID))
 logger.Warn("slow query detected", zap.Duration("query_time", 5*time.Second))
 logger.Error("database connection failed", zap.Error(err))
}

// ✅ 正确示例3：脱敏处理
func goodSanitizeData() {
 // 脱敏敏感信息
 maskedCard := maskCreditCard(cardNumber)
 logger.Info("payment_processed",
  zap.String("card_last4", maskedCard),
  zap.String("transaction_id", txnID),
 )
}

func maskCreditCard(card string) string {
 if len(card) < 4 {
  return "****"
 }
 return "****-****-****-" + card[len(card)-4:]
}

// ✅ 正确示例4：单例Logger
var (
 logger *zap.Logger
 once   sync.Once
)

func GetLogger() *zap.Logger {
 once.Do(func() {
  var err error
  logger, err = zap.NewProduction()
  if err != nil {
   panic(err)
  }
 })
 return logger
}

// ✅ 正确示例5：关联Trace上下文
func goodTraceContext(ctx context.Context) {
 span := trace.SpanFromContext(ctx)
 if span != nil {
  spanContext := span.SpanContext()
  logger.Info("processing",
   zap.String("trace_id", spanContext.TraceID().String()),
   zap.String("span_id", spanContext.SpanID().String()),
  )
 }
}

// ✅ 正确示例6：使用字段组
var commonFields = []zap.Field{
 zap.String("service", "user-service"),
 zap.String("version", "1.0.0"),
 zap.String("environment", "production"),
}

func goodWithCommonFields() {
 logger.With(commonFields...).Info("request processed")
}

// ✅ 正确示例7：采样日志（高频日志）
var sampledLogger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
 return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
}))

func goodSampledLogging() {
 // 每秒最多记录100条相同日志
 sampledLogger.Info("heartbeat")
}
```

---

## 5. eBPF基础

### 5.1 概念定义

**eBPF（Extended Berkeley Packet Filter）** 是一种革命性的内核技术，允许在内核空间安全地执行用户定义的程序，无需修改内核源码或加载内核模块。

**核心特性**：

- **安全性**：通过Verifier确保程序不会导致内核崩溃
- **高性能**：JIT编译为本地机器码，接近原生性能
- **灵活性**：可动态加载和卸载程序
- **可观测性**：可访问内核内部数据结构

**eBPF演进**：

```
1992: BPF - 网络包过滤
    │
    ▼
2014: eBPF - 扩展功能，支持更多内核钩子
    │
    ▼
2020+: eBPF生态成熟，成为云原生可观测性基石
```

### 5.2 eBPF架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         eBPF 架构                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  User Space                              Kernel Space                       │
│  ┌─────────────┐                         ┌─────────────────────────────┐    │
│  │  Go Program │                         │      eBPF Subsystem         │    │
│  │             │                         │  ┌─────────────────────┐    │    │
│  │ ┌─────────┐ │  ┌─────────────┐        │  │    eBPF Verifier    │    │    │
│  │ │ Cilium  │ │  │  eBPF Byte  │───────→│  │  (安全检查)          │    │    │
│  │ │ eBPF    │ │  │    Code     │        │  └─────────────────────┘    │    │
│  │ │ Library │ │  │  (.o文件)   │        │  ┌─────────────────────┐    │    │
│  │ └────┬────┘ │  └─────────────┘        │  │   eBPF JIT Compiler │    │    │
│  │      │      │                         │  │  (编译为机器码)       │    │    │
│  │  ┌───┴───┐  │  ┌─────────────┐        │  └─────────────────────┘    │    │
│  │  │Maps   │←─┼──│  eBPF Maps  │←───────│  ┌─────────────────────┐    │    │
│  │  │(数据) │  │  │  (共享内存)  │        │  │   eBPF Programs     │    │    │
│  │  └───┬───┘  │  └─────────────┘        │  │  ┌───────────────┐  │    │    │
│  │      │      │                         │  │  │  Kprobe       │  │    │    │
│  │  ┌───┴───┐  │                         │  │  │  Tracepoint   │  │    │    │
│  │  │Events │  │                         │  │  │  XDP          │  │    │    │
│  │  │(事件) │  │                         │  │  │  Socket Filter│  │    │    │
│  │  └───────┘  │                         │  │  │  TC           │  │    │    │
│  └─────────────┘                         │  │  │  LSM          │  │    │    │
│                                          │  │  └───────────────┘  │    │    │
│                                          │  └─────────────────────┘    │    │
│                                          └─────────────────────────────┘    │
│                                                                              │
│  eBPF程序类型：                                                               │
│  • Kprobe/Kretprobe: 内核函数入口/返回钩子                                    │
│  • Uprobe/Uretprobe: 用户空间函数钩子                                         │
│  • Tracepoint: 内核预定义跟踪点                                               │
│  • XDP: 网络包处理（最早阶段）                                                │
│  • TC: 流量控制                                                               │
│  • LSM: Linux安全模块                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 eBPF程序生命周期

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      eBPF 程序生命周期                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 编写              2. 编译               3. 加载              4. 验证     │
│  ┌─────────┐         ┌─────────┐          ┌─────────┐         ┌─────────┐   │
│  │ C/Rust  │────────→│  LLVM   │─────────→│  bpf()  │────────→│Verifier │   │
│  │ 源码    │         │ 编译器   │          │ 系统调用 │         │ 检查    │   │
│  └─────────┘         └─────────┘          └─────────┘         └────┬────┘   │
│       │                   │                  │                     │        │
│       ▼                   ▼                  ▼                     ▼        │
│  程序逻辑            eBPF字节码           加载到内核            安全检查     │
│                                                                              │
│  5. JIT编译           6. 附加               7. 执行              8. 清理     │
│  ┌─────────┐         ┌─────────┐          ┌─────────┐         ┌─────────┐   │
│  │  JIT    │────────→│ Attach  │─────────→│  运行   │────────→│  Detach │   │
│  │ 编译器   │         │ 到钩子  │          │  处理   │         │  卸载   │   │
│  └─────────┘         └─────────┘          └─────────┘         └─────────┘   │
│       │                   │                  │                     │        │
│       ▼                   ▼                  ▼                     ▼        │
│  机器码              绑定到事件            触发执行              资源释放     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.4 Go与eBPF交互（cilium/ebpf库）

**cilium/ebpf** 是Go语言中最流行的eBPF库，提供了纯Go的eBPF程序加载和管理功能。

**安装**：

```bash
go get github.com/cilium/ebpf
```

**Go实现示例 - 系统调用追踪**：

```go
package main

import (
 "encoding/binary"
 "fmt"
 "log"
 "os"
 "os/signal"
 "syscall"

 "github.com/cilium/ebpf"
 "github.com/cilium/ebpf/link"
 "github.com/cilium/ebpf/ringbuf"
 "github.com/cilium/ebpf/rlimit"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf sys_enter.bpf.c -- -I../headers

// Event 定义eBPF程序发送的事件结构
type Event struct {
 PID  uint32
 UID  uint32
 Comm [16]byte
 SyscallID uint32
}

// 嵌入式eBPF程序（简化版）
// 实际使用时需要编译C代码为eBPF字节码
const bpfProgram = `
#include <linux/bpf.h>
#include <linux/ptrace.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct event {
    u32 pid;
    u32 uid;
    char comm[16];
    u32 syscall_id;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} events SEC(".maps");

SEC("tracepoint/raw_syscalls/sys_enter")
int trace_sys_enter(struct trace_event_raw_sys_enter *ctx)
{
    struct event *e;

    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e)
        return 0;

    e->pid = bpf_get_current_pid_tgid() >> 32;
    e->uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    e->syscall_id = ctx->id;

    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
`

// loadEBPFProgram 加载eBPF程序
func loadEBPFProgram() (*ebpf.Collection, error) {
 // 解除内存限制
 if err := rlimit.RemoveMemlock(); err != nil {
  return nil, fmt.Errorf("failed to remove memlock: %w", err)
 }

 // 加载eBPF程序
 // 实际使用时需要预编译的eBPF对象文件
 spec, err := ebpf.LoadCollectionSpec("sys_enter.bpf.o")
 if err != nil {
  return nil, fmt.Errorf("failed to load collection spec: %w", err)
 }

 coll, err := ebpf.NewCollection(spec)
 if err != nil {
  return nil, fmt.Errorf("failed to create collection: %w", err)
 }

 return coll, nil
}

// traceSyscalls 追踪系统调用
func traceSyscalls() error {
 // 加载eBPF程序
 coll, err := loadEBPFProgram()
 if err != nil {
  return err
 }
 defer coll.Close()

 // 获取eBPF程序
 prog := coll.Programs["trace_sys_enter"]
 if prog == nil {
  return fmt.Errorf("program not found")
 }

 // 获取Ring Buffer Map
 eventsMap := coll.Maps["events"]
 if eventsMap == nil {
  return fmt.Errorf("events map not found")
 }

 // 附加到Tracepoint
 tp, err := link.Tracepoint("raw_syscalls", "sys_enter", prog, nil)
 if err != nil {
  return fmt.Errorf("failed to attach tracepoint: %w", err)
 }
 defer tp.Close()

 // 创建Ring Buffer读取器
 rd, err := ringbuf.NewReader(eventsMap)
 if err != nil {
  return fmt.Errorf("failed to create ringbuf reader: %w", err)
 }
 defer rd.Close()

 // 处理信号
 sig := make(chan os.Signal, 1)
 signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

 fmt.Println("Tracing syscalls... Press Ctrl+C to stop")

 // 读取事件
 go func() {
  for {
   record, err := rd.Read()
   if err != nil {
    if err == ringbuf.ErrClosed {
     return
    }
    log.Printf("Failed to read from ringbuf: %v", err)
    continue
   }

   var event Event
   if err := binary.Read(record.RawSample, binary.LittleEndian, &event); err != nil {
    log.Printf("Failed to parse event: %v", err)
    continue
   }

   fmt.Printf("PID: %d, UID: %d, Comm: %s, Syscall: %d\n",
    event.PID, event.UID, string(event.Comm[:]), event.SyscallID)
  }
 }()

 <-sig
 fmt.Println("\nStopping...")
 return nil
}

func main() {
 if err := traceSyscalls(); err != nil {
  log.Fatalf("Error: %v", err)
 }
}
```

### 5.5 eBPF C程序示例

**sys_enter.bpf.c** - 系统调用入口追踪：

```c
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

#define TASK_COMM_LEN 16

// 事件结构
tstruct event {
    u32 pid;
    u32 uid;
    char comm[TASK_COMM_LEN];
    u32 syscall_id;
    u64 timestamp;
};

// Ring Buffer Map定义
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} events SEC(".maps");

// 系统调用入口追踪程序
SEC("tracepoint/raw_syscalls/sys_enter")
int trace_sys_enter(struct trace_event_raw_sys_enter *ctx)
{
    struct event *e;
    u64 id;

    // 预留Ring Buffer空间
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e)
        return 0;

    // 获取进程ID
    id = bpf_get_current_pid_tgid();
    e->pid = id >> 32;

    // 获取用户ID
    id = bpf_get_current_uid_gid();
    e->uid = id & 0xFFFFFFFF;

    // 获取进程名
    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    // 获取系统调用号
    e->syscall_id = ctx->id;

    // 获取时间戳
    e->timestamp = bpf_ktime_get_ns();

    // 提交事件
    bpf_ringbuf_submit(e, 0);

    return 0;
}

// 许可证声明
char LICENSE[] SEC("license") = "GPL";
```

### 5.6 性能剖析应用

**eBPF在性能剖析中的应用**：

```go
package main

import (
 "bytes"
 "encoding/binary"
 "fmt"
 "log"
 "os"
 "os/signal"
 "sort"
 "syscall"
 "time"

 "github.com/cilium/ebpf"
 "github.com/cilium/ebpf/link"
 "github.com/cilium/ebpf/perf"
 "github.com/cilium/ebpf/rlimit"
)

// StackTrace 栈跟踪信息
type StackTrace struct {
 PID       uint32
 Comm      [16]byte
 KStackID  int32
 UStackID  int32
}

// StackCounts 栈计数器
type StackCounts struct {
 Counts map[string]uint64
}

// CPUProfiler CPU性能剖析器
type CPUProfiler struct {
 coll       *ebpf.Collection
 profileLink link.Link
 reader     *perf.Reader
 stackTraces *ebpf.Map
 sampleRate int
}

// NewCPUProfiler 创建CPU剖析器
func NewCPUProfiler(sampleRate int) (*CPUProfiler, error) {
 // 解除内存限制
 if err := rlimit.RemoveMemlock(); err != nil {
  return nil, err
 }

 // 加载eBPF程序
 spec, err := ebpf.LoadCollectionSpec("cpu_profile.bpf.o")
 if err != nil {
  return nil, fmt.Errorf("failed to load spec: %w", err)
 }

 // 设置采样率
 spec.Maps["sample_rate"].Contents = []ebpf.MapKV{{
  Key:   uint32(0),
  Value: uint64(sampleRate),
 }}

 coll, err := ebpf.NewCollection(spec)
 if err != nil {
  return nil, fmt.Errorf("failed to create collection: %w", err)
 }

 return &CPUProfiler{
  coll:       coll,
  stackTraces: coll.Maps["stack_traces"],
  sampleRate: sampleRate,
 }, nil
}

// Start 开始剖析
func (p *CPUProfiler) Start() error {
 prog := p.coll.Programs["do_sample"]
 if prog == nil {
  return fmt.Errorf("program not found")
 }

 // 附加到perf事件
 profileLink, err := link.AttachPerfEvent(link.PerfEventOptions{
  Type:   link.PerfEventTypeSoftware,
  Config: link.PerfEventConfigCPUClock,
  SampleRate: uint64(p.sampleRate),
 }, prog)
 if err != nil {
  return fmt.Errorf("failed to attach perf event: %w", err)
 }
 p.profileLink = profileLink

 // 创建perf reader
 countsMap := p.coll.Maps["counts"]
 reader, err := perf.NewReader(countsMap, 4096)
 if err != nil {
  return fmt.Errorf("failed to create perf reader: %w", err)
 }
 p.reader = reader

 return nil
}

// Stop 停止剖析
func (p *CPUProfiler) Stop() error {
 if p.profileLink != nil {
  p.profileLink.Close()
 }
 if p.reader != nil {
  p.reader.Close()
 }
 p.coll.Close()
 return nil
}

// ReadSamples 读取采样数据
func (p *CPUProfiler) ReadSamples(duration time.Duration) (*StackCounts, error) {
 counts := &StackCounts{
  Counts: make(map[string]uint64),
 }

 timeout := time.After(duration)

 for {
  select {
  case <-timeout:
   return counts, nil
  default:
   record, err := p.reader.Read()
   if err != nil {
    if err == perf.ErrClosed {
     return counts, nil
    }
    continue
   }

   var stack StackTrace
   if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &stack); err != nil {
    continue
   }

   // 解析栈跟踪
   stackKey := fmt.Sprintf("%s[%d]", string(stack.Comm[:]), stack.PID)
   counts.Counts[stackKey]++
  }
 }
}

// PrintFlameGraph 打印火焰图格式数据
func (sc *StackCounts) PrintFlameGraph() {
 // 按计数排序
 type kv struct {
  Key   string
  Value uint64
 }

 var sorted []kv
 for k, v := range sc.Counts {
  sorted = append(sorted, kv{k, v})
 }

 sort.Slice(sorted, func(i, j int) bool {
  return sorted[i].Value > sorted[j].Value
 })

 // 打印火焰图格式
 fmt.Println("\n=== Flame Graph Data ===")
 for _, kv := range sorted {
  fmt.Printf("%s %d\n", kv.Key, kv.Value)
 }
}

func main() {
 // 创建CPU剖析器 (99Hz采样率)
 profiler, err := NewCPUProfiler(99)
 if err != nil {
  log.Fatalf("Failed to create profiler: %v", err)
 }
 defer profiler.Stop()

 // 开始剖析
 if err := profiler.Start(); err != nil {
  log.Fatalf("Failed to start profiler: %v", err)
 }

 fmt.Println("CPU profiling started... Press Ctrl+C to stop")

 // 处理信号
 sig := make(chan os.Signal, 1)
 signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

 // 在另一个goroutine中收集数据
 done := make(chan *StackCounts)
 go func() {
  counts, err := profiler.ReadSamples(30 * time.Second)
  if err != nil {
   log.Printf("Error reading samples: %v", err)
  }
  done <- counts
 }()

 select {
 case <-sig:
  fmt.Println("\nStopping profiler...")
 case counts := <-done:
  counts.PrintFlameGraph()
 }
}
```

### 5.7 反例说明

```go
// ❌ 错误示例1：未解除内存限制
func badNoMemlock() {
 // 忘记调用rlimit.RemoveMemlock()
 coll, err := ebpf.NewCollection(spec) // 可能失败
}

// ❌ 错误示例2：未正确关闭资源
func badResourceLeak() {
 coll, _ := ebpf.NewCollection(spec)
 // 忘记defer coll.Close() - 资源泄漏
}

// ❌ 错误示例3：在eBPF程序中使用循环
// C代码中的错误
/*
SEC("kprobe/do_sys_open")
int bad_loop(struct pt_regs *ctx) {
    for (int i = 0; i < 1000000; i++) { // ❌ 无限循环会被Verifier拒绝
        // ...
    }
    return 0;
}
*/

// ❌ 错误示例4：访问无效内存
// C代码中的错误
/*
SEC("kprobe/do_sys_open")
int bad_memory_access(struct pt_regs *ctx) {
    char *ptr = NULL;
    char c = *ptr; // ❌ 空指针解引用
    return 0;
}
*/

// ❌ 错误示例5：Map操作不当
func badMapOperation() {
 // 未检查Map是否存在
 m := coll.Maps["non_existent_map"] // 返回nil
 m.Update(key, value, 0)            // panic
}
```

### 5.8 最佳实践

```go
// ✅ 正确示例1：解除内存限制
func goodMemlock() {
 if err := rlimit.RemoveMemlock(); err != nil {
  log.Fatalf("Failed to remove memlock: %v", err)
 }
 // ...
}

// ✅ 正确示例2：资源清理
func goodResourceCleanup() {
 coll, err := ebpf.NewCollection(spec)
 if err != nil {
  log.Fatalf("Failed to create collection: %v", err)
 }
 defer coll.Close() // 确保资源释放
 // ...
}

// ✅ 正确示例3：错误处理
func goodErrorHandling() {
 prog := coll.Programs["my_program"]
 if prog == nil {
  log.Fatal("Program not found")
 }

 link, err := link.AttachTracepoint("syscalls", "sys_enter_openat", prog, nil)
 if err != nil {
  log.Fatalf("Failed to attach: %v", err)
 }
 defer link.Close()
}

// ✅ 正确示例4：使用CO-RE（Compile Once, Run Everywhere）
// C代码中的正确做法
/*
#include "vmlinux.h"  // 使用BTF信息
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>

SEC("kprobe/do_sys_open")
int good_core_read(struct pt_regs *ctx) {
    struct task_struct *task = (struct task_struct *)bpf_get_current_task();

    // 使用BPF_CORE_READ宏进行安全读取
    pid_t pid = BPF_CORE_READ(task, pid);

    return 0;
}
*/

// ✅ 正确示例5：Ring Buffer使用
func goodRingBuffer() {
 // 创建Ring Buffer Reader
 rd, err := ringbuf.NewReader(eventsMap)
 if err != nil {
  log.Fatalf("Failed to create reader: %v", err)
 }
 defer rd.Close()

 // 使用超时读取
 record, err := rd.Read()
 if err != nil {
  if err == ringbuf.ErrClosed {
   return
  }
  log.Printf("Read error: %v", err)
 }
 // 处理record...
}

// ✅ 正确示例6：批量处理事件
func goodBatchProcessing() {
 batch := make([]Event, 0, 100)

 for i := 0; i < 100; i++ {
  record, err := rd.Read()
  if err != nil {
   break
  }

  var event Event
  if err := binary.Read(record.RawSample, binary.LittleEndian, &event); err != nil {
   continue
  }
  batch = append(batch, event)
 }

 // 批量处理
 processBatch(batch)
}
```

---

## 6. 性能剖析（Profiling）

### 6.1 概念定义

**性能剖析（Profiling）** 是一种动态程序分析技术，用于测量程序的空间（内存）或时间复杂度、特定指令的使用频率、函数调用频率和持续时间。

**剖析类型**：

| 类型 | 描述 | 用途 |
|------|------|------|
| CPU Profile | CPU时间消耗分析 | 找出CPU热点函数 |
| Memory Profile | 内存分配分析 | 检测内存泄漏、优化分配 |
| Goroutine Profile | Goroutine状态分析 | 检测死锁、Goroutine泄漏 |
| Block Profile | 阻塞分析 | 找出同步阻塞点 |
| Mutex Profile | 互斥锁竞争分析 | 优化锁粒度 |
| ThreadCreate Profile | 线程创建分析 | 检测线程泄漏 |

### 6.2 pprof架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         pprof 架构                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go Application                    pprof Tool                              │
│  ┌─────────────────┐               ┌─────────────────┐                      │
│  │ runtime/pprof   │               │  go tool pprof  │                      │
│  │                 │               │                 │                      │
│  │ ┌─────────────┐ │               │ ┌─────────────┐ │                      │
│  │ │ CPU Profile │─┼──────┐        │ │ Interactive │ │                      │
│  │ │ (采样)      │ │      │        │ │ CLI         │ │                      │
│  │ └─────────────┘ │      │        │ └─────────────┘ │                      │
│  │ ┌─────────────┐ │      │        │ ┌─────────────┐ │                      │
│  │ │ Mem Profile │─┼──────┼────────│→│ Web UI      │ │                      │
│  │ │ (堆快照)    │ │      │        │ └─────────────┘ │                      │
│  │ └─────────────┘ │      │        │ ┌─────────────┐ │                      │
│  │ ┌─────────────┐ │      │        │ │ Flame Graph │ │                      │
│  │ │ Goroutine   │─┼──────┘        │ │ SVG/PNG     │ │                      │
│  │ │ Profile     │ │    Profile    │ └─────────────┘ │                      │
│  │ └─────────────┘ │    Data       └─────────────────┘                      │
│  │ ┌─────────────┐ │                                                        │
│  │ │ Block/Mutex │ │                                                        │
│  │ │ Profile     │ │                                                        │
│  │ └─────────────┘ │                                                        │
│  └─────────────────┘                                                        │
│                                                                              │
│  采样机制：                                                                   │
│  • CPU: SIGPROF信号，默认100Hz采样                                           │
│  • Memory: 每次分配时采样（按大小或频率）                                     │
│  • Goroutine: 全量栈跟踪                                                     │
│  • Block/Mutex: 事件驱动采样                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 CPU Profiling

**概念**：CPU Profiling通过定期中断程序执行，记录当前调用栈，从而统计各函数消耗的CPU时间。

**Go实现示例**：

```go
package main

import (
 "fmt"
 "log"
 "os"
 "runtime"
 "runtime/pprof"
 "time"
)

// CPU密集型函数
func cpuIntensiveTask(n int) int {
 if n <= 1 {
  return n
 }
 return cpuIntensiveTask(n-1) + cpuIntensiveTask(n-2)
}

// 模拟工作负载
func simulateWork() {
 for i := 0; i < 10; i++ {
  go func() {
   for {
    _ = cpuIntensiveTask(35)
   }
  }()
 }
}

// startCPUProfile 开始CPU剖析
func startCPUProfile(filename string) (*os.File, error) {
 f, err := os.Create(filename)
 if err != nil {
  return nil, fmt.Errorf("could not create CPU profile: %v", err)
 }

 if err := pprof.StartCPUProfile(f); err != nil {
  f.Close()
  return nil, fmt.Errorf("could not start CPU profile: %v", err)
 }

 return f, nil
}

// stopCPUProfile 停止CPU剖析
func stopCPUProfile(f *os.File) {
 pprof.StopCPUProfile()
 f.Close()
}

func main() {
 // 设置最大CPU核心数
 runtime.GOMAXPROCS(runtime.NumCPU())

 // 开始CPU剖析
 f, err := startCPUProfile("cpu.prof")
 if err != nil {
  log.Fatal(err)
 }
 defer stopCPUProfile(f)

 fmt.Println("CPU profiling started...")

 // 启动工作负载
 simulateWork()

 // 运行一段时间
 time.Sleep(30 * time.Second)

 fmt.Println("CPU profiling completed. Profile saved to cpu.prof")
 fmt.Println("Analyze with: go tool pprof cpu.prof")
}
```

**分析CPU Profile**：

```bash
# 交互式分析
go tool pprof cpu.prof

# 常用命令：
# (pprof) top          - 显示最耗时的函数
# (pprof) top 20       - 显示前20个函数
# (pprof) list main    - 显示main包的源码级分析
# (pprof) web          - 生成SVG图形
# (pprof) png          - 生成PNG图像
# (pprof) pdf          - 生成PDF文档
# (pprof) flamegraph   - 生成火焰图

# 直接生成火焰图
go tool pprof -http=:8080 cpu.prof
```

### 6.4 Memory Profiling

**概念**：Memory Profiling记录堆内存分配信息，帮助识别内存泄漏和高内存消耗点。

**Go实现示例**：

```go
package main

import (
 "fmt"
 "log"
 "os"
 "runtime"
 "runtime/pprof"
 "time"
)

// 模拟内存泄漏
type Data struct {
 buffer []byte
}

var leakySlice []*Data

func memoryLeak() {
 // 不断分配内存但不释放
 for {
  data := &Data{
   buffer: make([]byte, 1024*1024), // 1MB
  }
  leakySlice = append(leakySlice, data)
  time.Sleep(100 * time.Millisecond)
 }
}

// 正常内存使用
func normalMemoryUsage() {
 for {
  data := make([]byte, 1024*1024) // 1MB
  _ = data
  time.Sleep(100 * time.Millisecond)
  runtime.GC() // 强制GC
 }
}

// writeHeapProfile 写入堆内存剖析
func writeHeapProfile(filename string) error {
 f, err := os.Create(filename)
 if err != nil {
  return fmt.Errorf("could not create memory profile: %v", err)
 }
 defer f.Close()

 runtime.GC() // 先执行GC获取准确数据
 if err := pprof.WriteHeapProfile(f); err != nil {
  return fmt.Errorf("could not write memory profile: %v", err)
 }

 return nil
}

func main() {
 // 启动内存泄漏模拟
 go memoryLeak()
 go normalMemoryUsage()

 // 定期保存内存剖析
 ticker := time.NewTicker(10 * time.Second)
 defer ticker.Stop()

 count := 0
 for range ticker.C {
  count++
  filename := fmt.Sprintf("heap_%d.prof", count)
  if err := writeHeapProfile(filename); err != nil {
   log.Printf("Failed to write heap profile: %v", err)
   continue
  }
  fmt.Printf("Heap profile saved to %s\n", filename)

  if count >= 3 {
   break
  }
 }

 fmt.Println("\nMemory profiling completed.")
 fmt.Println("Compare profiles with: go tool pprof -diff_base heap_1.prof heap_3.prof")
}
```

**分析Memory Profile**：

```bash
# 查看内存分配
go tool pprof heap.prof

# 常用命令：
# (pprof) top          - 显示分配最多的函数
# (pprof) list main    - 显示源码级分配详情
# (pprof) alloc_space  - 按分配空间排序（默认）
# (pprof) inuse_space  - 按使用空间排序
# (pprof) alloc_objects - 按分配对象数排序
# (pprof) inuse_objects - 按使用对象数排序

# 比较两个时间点的内存
go tool pprof -diff_base heap_1.prof heap_3.prof
```

### 6.5 Goroutine Profiling

**概念**：Goroutine Profiling记录所有Goroutine的栈跟踪，用于检测死锁、Goroutine泄漏和分析并发模式。

**Go实现示例**：

```go
package main

import (
 "fmt"
 "log"
 "os"
 "runtime/pprof"
 "sync"
 "time"
)

// 模拟Goroutine泄漏
func goroutineLeak() {
 var wg sync.WaitGroup

 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   // 模拟长时间运行的任务
   ch := make(chan int)
   go func() {
    // 这个Goroutine永远不会退出
    for {
     select {
     case <-ch:
      return
     default:
      time.Sleep(1 * time.Second)
     }
    }
   }()

   // 任务完成后不关闭channel
   time.Sleep(100 * time.Millisecond)
  }(i)
 }

 wg.Wait()
}

// 模拟死锁
func deadlockExample() {
 ch1 := make(chan int)
 ch2 := make(chan int)

 go func() {
  <-ch1
  ch2 <- 1
 }()

 go func() {
  <-ch2
  ch1 <- 1
 }()
}

// writeGoroutineProfile 写入Goroutine剖析
func writeGoroutineProfile(filename string) error {
 f, err := os.Create(filename)
 if err != nil {
  return fmt.Errorf("could not create goroutine profile: %v", err)
 }
 defer f.Close()

 if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
  return fmt.Errorf("could not write goroutine profile: %v", err)
 }

 return nil
}

func main() {
 // 启动Goroutine泄漏模拟
 go goroutineLeak()
 go deadlockExample()

 // 等待Goroutine创建
 time.Sleep(5 * time.Second)

 // 保存Goroutine剖析
 if err := writeGoroutineProfile("goroutine.prof"); err != nil {
  log.Fatal(err)
 }

 fmt.Println("Goroutine profile saved to goroutine.prof")
 fmt.Println("Analyze with: go tool pprof goroutine.prof")

 // 保持运行以便分析
 time.Sleep(1 * time.Hour)
}
```

**分析Goroutine Profile**：

```bash
# 查看Goroutine状态
go tool pprof goroutine.prof

# 常用命令：
# (pprof) top          - 显示Goroutine数量最多的栈
# (pprof) list main    - 显示源码级Goroutine创建点
# (pprof) traces       - 显示所有Goroutine的栈跟踪

# 直接查看文本格式的Goroutine dump
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 查看所有Goroutine的完整栈
curl http://localhost:8080/debug/pprof/goroutine?debug=2
```

### 6.6 Block和Mutex Profiling

**概念**：Block Profiling记录Goroutine阻塞事件，Mutex Profiling记录互斥锁竞争。

**Go实现示例**：

```go
package main

import (
 "fmt"
 "log"
 "os"
 "runtime"
 "runtime/pprof"
 "sync"
 "time"
)

// 模拟锁竞争
func mutexContention() {
 var mu sync.Mutex
 var counter int

 for i := 0; i < 100; i++ {
  go func() {
   for {
    mu.Lock()
    counter++
    mu.Unlock()
   }
  }()
 }
}

// 模拟Channel阻塞
func channelBlocking() {
 ch := make(chan int)

 for i := 0; i < 100; i++ {
  go func() {
   for {
    ch <- 1 // 阻塞发送
   }
  }()
 }

 // 缓慢消费
 go func() {
  for {
   <-ch
   time.Sleep(10 * time.Millisecond)
  }
 }()
}

// 启用Block和Mutex剖析
func enableProfiling() {
 // 启用Block剖析（采样率1表示100%采样）
 runtime.SetBlockProfileRate(1)

 // 启用Mutex剖析（采样率5表示5%采样）
 runtime.SetMutexProfileFraction(5)
}

// writeBlockProfile 写入Block剖析
func writeBlockProfile(filename string) error {
 f, err := os.Create(filename)
 if err != nil {
  return fmt.Errorf("could not create block profile: %v", err)
 }
 defer f.Close()

 if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
  return fmt.Errorf("could not write block profile: %v", err)
 }

 return nil
}

// writeMutexProfile 写入Mutex剖析
func writeMutexProfile(filename string) error {
 f, err := os.Create(filename)
 if err != nil {
  return fmt.Errorf("could not create mutex profile: %v", err)
 }
 defer f.Close()

 if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
  return fmt.Errorf("could not write mutex profile: %v", err)
 }

 return nil
}

func main() {
 // 启用剖析
 enableProfiling()

 // 启动竞争模拟
 go mutexContention()
 go channelBlocking()

 // 运行一段时间
 time.Sleep(10 * time.Second)

 // 保存剖析数据
 if err := writeBlockProfile("block.prof"); err != nil {
  log.Printf("Failed to write block profile: %v", err)
 } else {
  fmt.Println("Block profile saved to block.prof")
 }

 if err := writeMutexProfile("mutex.prof"); err != nil {
  log.Printf("Failed to write mutex profile: %v", err)
 } else {
  fmt.Println("Mutex profile saved to mutex.prof")
 }
}
```

### 6.7 HTTP pprof端点

**net/http/pprof** 包提供了便捷的HTTP端点用于实时剖析。

**Go实现示例**：

```go
package main

import (
 "fmt"
 "net/http"
 _ "net/http/pprof" // 自动注册pprof端点
 "runtime"
 "time"
)

func init() {
 // 启用所有剖析类型
 runtime.SetBlockProfileRate(1)
 runtime.SetMutexProfileFraction(5)
}

func main() {
 // 启动HTTP服务器
 go func() {
  fmt.Println("pprof server starting on :6060")
  fmt.Println("Available endpoints:")
  fmt.Println("  http://localhost:6060/debug/pprof/")
  fmt.Println("  http://localhost:6060/debug/pprof/profile")
  fmt.Println("  http://localhost:6060/debug/pprof/heap")
  fmt.Println("  http://localhost:6060/debug/pprof/goroutine")
  fmt.Println("  http://localhost:6060/debug/pprof/block")
  fmt.Println("  http://localhost:6060/debug/pprof/mutex")
  http.ListenAndServe("localhost:6060", nil)
 }()

 // 模拟工作负载
 for {
  simulateWork()
  time.Sleep(1 * time.Second)
 }
}

func simulateWork() {
 // CPU工作
 for i := 0; i < 1000000; i++ {
  _ = i * i
 }

 // 内存分配
 _ = make([]byte, 1024*1024)
}
```

**pprof端点列表**：

| 端点 | 描述 | 参数 |
|------|------|------|
| /debug/pprof/ | 索引页面 | - |
| /debug/pprof/profile | CPU Profile | seconds=30 |
| /debug/pprof/heap | 堆内存 | gc=1 |
| /debug/pprof/goroutine | Goroutine | debug=1/2 |
| /debug/pprof/block | 阻塞分析 | - |
| /debug/pprof/mutex | 锁竞争 | - |
| /debug/pprof/threadcreate | 线程创建 | - |
| /debug/pprof/allocs | 分配分析 | - |
| /debug/pprof/cmdline | 命令行参数 | - |

### 6.8 反例说明

```go
// ❌ 错误示例1：在生产环境过度采样
func badHighSamplingRate() {
 // 10000Hz采样会严重影响性能
 pprof.StartCPUProfile(f) // 默认100Hz已经足够
}

// ❌ 错误示例2：忘记关闭profile文件
func badNoClose() {
 f, _ := os.Create("cpu.prof")
 pprof.StartCPUProfile(f)
 // 忘记defer pprof.StopCPUProfile()和f.Close()
}

// ❌ 错误示例3：在剖析期间修改程序行为
func badAlterBehavior() {
 pprof.StartCPUProfile(f)
 // 程序行为与正常情况不同，剖析结果不准确
}

// ❌ 错误示例4：忽略剖析开销
func badIgnoreOverhead() {
 // 没有考虑剖析本身的开销
 for {
  pprof.StartCPUProfile(f)
  doWork()
  pprof.StopCPUProfile()
 }
}

// ❌ 错误示例5：不正确的内存剖析时机
func badWrongTiming() {
 // 在GC之前记录内存，数据不准确
 pprof.WriteHeapProfile(f)
 runtime.GC()
}
```

### 6.9 最佳实践

```go
// ✅ 正确示例1：合理的采样率
func goodSamplingRate() {
 // CPU剖析使用默认100Hz
 // Block剖析根据需要设置
 runtime.SetBlockProfileRate(100) // 每100纳秒采样一次阻塞
 // Mutex剖析使用合理的比例
 runtime.SetMutexProfileFraction(100) // 1%采样
}

// ✅ 正确示例2：资源清理
func goodResourceCleanup() {
 f, err := os.Create("cpu.prof")
 if err != nil {
  log.Fatal(err)
 }
 defer func() {
  pprof.StopCPUProfile()
  f.Close()
 }()

 pprof.StartCPUProfile(f)
 // ...
}

// ✅ 正确示例3：条件性剖析
var cpuProfileEnabled = flag.Bool("cpuprofile", false, "write cpu profile")

func goodConditionalProfiling() {
 if *cpuProfileEnabled {
  f, err := os.Create("cpu.prof")
  if err != nil {
   log.Fatal(err)
  }
  defer f.Close()
  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()
 }
}

// ✅ 正确示例4：准确的内存剖析
func goodAccurateMemoryProfile() {
 runtime.GC() // 先执行GC
 time.Sleep(100 * time.Millisecond) // 等待GC完成
 pprof.WriteHeapProfile(f)
}

// ✅ 正确示例5：剖析数据上传
func goodProfileUpload() {
 // 定期保存并上传剖析数据
 ticker := time.NewTicker(1 * time.Hour)
 for range ticker.C {
  filename := fmt.Sprintf("heap_%d.prof", time.Now().Unix())
  if err := writeHeapProfile(filename); err != nil {
   log.Printf("Failed to write profile: %v", err)
   continue
  }
  // 上传到分析服务
  uploadProfile(filename)
 }
}

// ✅ 正确示例6：HTTP pprof安全配置
func goodSecurePprof() {
 // 只在内部网络暴露pprof端点
 go func() {
  mux := http.NewServeMux()
  mux.HandleFunc("/debug/pprof/", pprof.Index)
  mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
  // ... 其他端点

  // 绑定到localhost
  log.Fatal(http.ListenAndServe("localhost:6060", mux))
 }()
}
```

---

## 7. 健康检查

### 7.1 概念定义

**健康检查（Health Check）** 是一种监控机制，用于确定应用程序或服务是否正常运行，以及是否准备好接收流量。

**Kubernetes探针类型**：

| 探针类型 | 用途 | 失败行为 |
|----------|------|----------|
| Liveness Probe | 检测应用是否存活 | 重启容器 |
| Readiness Probe | 检测应用是否就绪 | 从Service端点移除 |
| Startup Probe | 检测应用是否启动完成 | 禁用其他探针 |

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      健康检查探针对比                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Liveness Probe (存活探针)                                                   │
│  ┌─────────────┐                                                            │
│  │ Application │                                                            │
│  │   Running   │ ──失败──→ 重启容器                                          │
│  └─────────────┘                                                            │
│                                                                              │
│  用途：检测死锁、无限循环等导致应用无响应的情况                                 │
│  注意：不要依赖外部服务，避免级联重启                                          │
│                                                                              │
│  Readiness Probe (就绪探针)                                                  │
│  ┌─────────────┐     ┌─────────────┐                                        │
│  │ Application │────→│   Service   │                                        │
│  │   Ready     │     │  Endpoints  │                                        │
│  └─────────────┘     └─────────────┘                                        │
│         │                   │                                               │
│       失败                  ↓                                               │
│         │              从端点列表移除                                         │
│         └───────────────────┘                                               │
│                                                                              │
│  用途：检测应用是否准备好处理请求                                             │
│  注意：可以检查外部依赖（数据库、缓存等）                                       │
│                                                                              │
│  Startup Probe (启动探针)                                                    │
│  ┌─────────────┐                                                            │
│  │ Application │                                                            │
│  │  Starting   │ ──成功──→ 启用Liveness/Readiness探针                       │
│  └─────────────┘                                                            │
│                                                                              │
│  用途：保护启动慢的应用，避免过早判定失败                                       │
│  注意：只影响启动阶段，成功后不再执行                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 健康检查架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      健康检查架构                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Kubernetes                        Application                              │
│  ┌─────────────┐                   ┌─────────────────────────────┐          │
│  │   Kubelet   │                   │      Health Check Server    │          │
│  │             │  HTTP GET         │  ┌─────────────────────┐    │          │
│  │ ┌─────────┐ │  /healthz/live    │  │  Liveness Handler   │    │          │
│  │ │ Probe   │─┼──────────────────→│  │  - Basic check      │    │          │
│  │ │ Config  │ │                   │  │  - Goroutine check  │    │          │
│  │ └─────────┘ │  HTTP GET         │  └─────────────────────┘    │          │
│  │             │  /healthz/ready   │  ┌─────────────────────┐    │          │
│  │ ┌─────────┐ │──────────────────→│  │  Readiness Handler  │    │          │
│  │ │ Service │ │                   │  │  - DB connection    │    │          │
│  │ │ Endpoint│←┼───────────────────│  │  - Cache check      │    │          │
│  │ └─────────┘ │  200/503          │  │  - External deps    │    │          │
│  └─────────────┘                   │  └─────────────────────┘    │          │
│                                    │  ┌─────────────────────┐    │          │
│                                    │  │  Startup Handler    │    │          │
│                                    │  │  - Init complete    │    │          │
│                                    │  └─────────────────────┘    │          │
│                                    └─────────────────────────────┘          │
│                                                                              │
│  探针配置参数：                                                               │
│  • initialDelaySeconds: 首次检查前的等待时间                                  │
│  • periodSeconds: 检查间隔                                                   │
│  • timeoutSeconds: 超时时间                                                  │
│  • successThreshold: 成功阈值（连续成功次数）                                  │
│  • failureThreshold: 失败阈值（连续失败次数）                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Go实现示例

```go
package main

import (
 "context"
 "database/sql"
 "encoding/json"
 "fmt"
 "net/http"
 "runtime"
 "sync/atomic"
 "time"

 _ "github.com/lib/pq"
 "github.com/redis/go-redis/v9"
)

// HealthStatus 健康状态
type HealthStatus string

const (
 StatusUp   HealthStatus = "UP"
 StatusDown HealthStatus = "DOWN"
)

// HealthCheck 健康检查结果
type HealthCheck struct {
 Status    HealthStatus       `json:"status"`
 Timestamp time.Time          `json:"timestamp"`
 Version   string             `json:"version"`
 Checks    map[string]*Check  `json:"checks"`
}

// Check 单个检查项
type Check struct {
 Status    HealthStatus `json:"status"`
 Message   string       `json:"message,omitempty"`
 Latency   string       `json:"latency,omitempty"`
 Timestamp time.Time    `json:"timestamp"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
 db          *sql.DB
 redisClient *redis.Client
 version     string
 startTime   time.Time
 ready       atomic.Bool
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(db *sql.DB, redisClient *redis.Client, version string) *HealthChecker {
 hc := &HealthChecker{
  db:          db,
  redisClient: redisClient,
  version:     version,
  startTime:   time.Now(),
 }
 // 默认未就绪
 hc.ready.Store(false)
 return hc
}

// SetReady 设置就绪状态
func (hc *HealthChecker) SetReady(ready bool) {
 hc.ready.Store(ready)
}

// LivenessCheck 存活检查
func (hc *HealthChecker) LivenessCheck(w http.ResponseWriter, r *http.Request) {
 start := time.Now()

 // 基本存活检查
 check := &HealthCheck{
  Status:    StatusUp,
  Timestamp: time.Now(),
  Version:   hc.version,
  Checks:    make(map[string]*Check),
 }

 // 检查Goroutine数量（防止泄漏）
 goroutineCount := runtime.NumGoroutine()
 if goroutineCount > 10000 {
  check.Status = StatusDown
  check.Checks["goroutine"] = &Check{
   Status:    StatusDown,
   Message:   fmt.Sprintf("Too many goroutines: %d", goroutineCount),
   Timestamp: time.Now(),
  }
 } else {
  check.Checks["goroutine"] = &Check{
   Status:    StatusUp,
   Message:   fmt.Sprintf("Goroutines: %d", goroutineCount),
   Timestamp: time.Now(),
  }
 }

 // 检查内存使用
 var m runtime.MemStats
 runtime.ReadMemStats(&m)
 memMB := m.Alloc / 1024 / 1024
 if memMB > 1024 { // 1GB
  check.Status = StatusDown
  check.Checks["memory"] = &Check{
   Status:    StatusDown,
   Message:   fmt.Sprintf("Memory usage too high: %d MB", memMB),
   Timestamp: time.Now(),
  }
 } else {
  check.Checks["memory"] = &Check{
   Status:    StatusUp,
   Message:   fmt.Sprintf("Memory: %d MB", memMB),
   Timestamp: time.Now(),
  }
 }

 // 响应
 w.Header().Set("Content-Type", "application/json")
 if check.Status == StatusDown {
  w.WriteHeader(http.StatusServiceUnavailable)
 } else {
  w.WriteHeader(http.StatusOK)
 }

 latency := time.Since(start)
 check.Checks["goroutine"].Latency = latency.String()

 json.NewEncoder(w).Encode(check)
}

// ReadinessCheck 就绪检查
func (hc *HealthChecker) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
 start := time.Now()

 check := &HealthCheck{
  Status:    StatusUp,
  Timestamp: time.Now(),
  Version:   hc.version,
  Checks:    make(map[string]*Check),
 }

 // 检查是否已标记为就绪
 if !hc.ready.Load() {
  check.Status = StatusDown
  check.Checks["ready"] = &Check{
   Status:    StatusDown,
   Message:   "Application not ready yet",
   Timestamp: time.Now(),
  }
 } else {
  check.Checks["ready"] = &Check{
   Status:    StatusUp,
   Message:   "Application is ready",
   Timestamp: time.Now(),
  }
 }

 // 检查数据库连接
 if hc.db != nil {
  dbStart := time.Now()
  ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
  defer cancel()

  if err := hc.db.PingContext(ctx); err != nil {
   check.Status = StatusDown
   check.Checks["database"] = &Check{
    Status:    StatusDown,
    Message:   fmt.Sprintf("Database ping failed: %v", err),
    Timestamp: time.Now(),
   }
  } else {
   check.Checks["database"] = &Check{
    Status:  StatusUp,
    Latency: time.Since(dbStart).String(),
    Timestamp: time.Now(),
   }
  }
 }

 // 检查Redis连接
 if hc.redisClient != nil {
  redisStart := time.Now()
  ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
  defer cancel()

  if err := hc.redisClient.Ping(ctx).Err(); err != nil {
   check.Status = StatusDown
   check.Checks["redis"] = &Check{
    Status:    StatusDown,
    Message:   fmt.Sprintf("Redis ping failed: %v", err),
    Timestamp: time.Now(),
   }
  } else {
   check.Checks["redis"] = &Check{
    Status:  StatusUp,
    Latency: time.Since(redisStart).String(),
    Timestamp: time.Now(),
   }
  }
 }

 // 检查运行时间
 uptime := time.Since(hc.startTime)
 check.Checks["uptime"] = &Check{
  Status:    StatusUp,
  Message:   fmt.Sprintf("Uptime: %v", uptime),
  Timestamp: time.Now(),
 }

 // 响应
 w.Header().Set("Content-Type", "application/json")
 if check.Status == StatusDown {
  w.WriteHeader(http.StatusServiceUnavailable)
 } else {
  w.WriteHeader(http.StatusOK)
 }

 latency := time.Since(start)
 for _, c := range check.Checks {
  if c.Latency == "" {
   c.Latency = latency.String()
  }
 }

 json.NewEncoder(w).Encode(check)
}

// StartupCheck 启动检查
func (hc *HealthChecker) StartupCheck(w http.ResponseWriter, r *http.Request) {
 // 启动检查：应用是否已完成初始化
 if time.Since(hc.startTime) < 10*time.Second {
  w.WriteHeader(http.StatusServiceUnavailable)
  json.NewEncoder(w).Encode(map[string]interface{}{
   "status":  "DOWN",
   "message": "Application still starting",
  })
  return
 }

 w.WriteHeader(http.StatusOK)
 json.NewEncoder(w).Encode(map[string]interface{}{
  "status":  "UP",
  "message": "Application started successfully",
 })
}

// 模拟数据库初始化
func initDB() (*sql.DB, error) {
 // 实际使用时替换为真实的数据库连接
 db, err := sql.Open("postgres", "postgres://user:pass@localhost/db?sslmode=disable")
 if err != nil {
  return nil, err
 }
 return db, nil
}

// 模拟Redis初始化
func initRedis() *redis.Client {
 return redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })
}

func main() {
 // 初始化依赖
 db, err := initDB()
 if err != nil {
  fmt.Printf("Failed to init DB: %v\n", err)
  db = nil // 允许无DB运行
 }

 redisClient := initRedis()

 // 创建健康检查器
 hc := NewHealthChecker(db, redisClient, "1.0.0")

 // 设置路由
 http.HandleFunc("/healthz/live", hc.LivenessCheck)
 http.HandleFunc("/healthz/ready", hc.ReadinessCheck)
 http.HandleFunc("/healthz/startup", hc.StartupCheck)

 // 模拟异步初始化
 go func() {
  time.Sleep(5 * time.Second)
  hc.SetReady(true)
  fmt.Println("Application is now ready")
 }()

 fmt.Println("Health check server starting on :8080")
 fmt.Println("Endpoints:")
 fmt.Println("  /healthz/live    - Liveness probe")
 fmt.Println("  /healthz/ready   - Readiness probe")
 fmt.Println("  /healthz/startup - Startup probe")

 if err := http.ListenAndServe(":8080", nil); err != nil {
  fmt.Printf("Server failed: %v\n", err)
 }
}
```

### 7.4 Kubernetes配置示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: my-app:1.0.0
        ports:
        - containerPort: 8080

        # 存活探针 - 检测应用是否存活
        livenessProbe:
          httpGet:
            path: /healthz/live
            port: 8080
          initialDelaySeconds: 10    # 首次检查前等待10秒
          periodSeconds: 10          # 每10秒检查一次
          timeoutSeconds: 5          # 超时5秒
          failureThreshold: 3        # 连续3次失败才认为失败
          successThreshold: 1        # 1次成功即认为成功

        # 就绪探针 - 检测应用是否准备好接收流量
        readinessProbe:
          httpGet:
            path: /healthz/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
          successThreshold: 1

        # 启动探针 - 保护启动慢的应用
        startupProbe:
          httpGet:
            path: /healthz/startup
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 30       # 允许最多150秒启动时间

---
apiVersion: v1
kind: Service
metadata:
  name: my-app
spec:
  selector:
    app: my-app
  ports:
  - port: 80
    targetPort: 8080
```

### 7.5 高级健康检查模式

```go
package main

import (
 "context"
 "encoding/json"
 "net/http"
 "sync"
 "time"
)

// Checkable 可检查接口
type Checkable interface {
 Name() string
 Check(ctx context.Context) error
}

// HealthRegistry 健康检查注册表
type HealthRegistry struct {
 checks map[string]Checkable
 mu     sync.RWMutex
}

// NewHealthRegistry 创建注册表
func NewHealthRegistry() *HealthRegistry {
 return &HealthRegistry{
  checks: make(map[string]Checkable),
 }
}

// Register 注册检查项
func (r *HealthRegistry) Register(check Checkable) {
 r.mu.Lock()
 defer r.mu.Unlock()
 r.checks[check.Name()] = check
}

// Unregister 注销检查项
func (r *HealthRegistry) Unregister(name string) {
 r.mu.Lock()
 defer r.mu.Unlock()
 delete(r.checks, name)
}

// RunChecks 执行所有检查
func (r *HealthRegistry) RunChecks(ctx context.Context) map[string]*CheckResult {
 r.mu.RLock()
 checks := make(map[string]Checkable, len(r.checks))
 for k, v := range r.checks {
  checks[k] = v
 }
 r.mu.RUnlock()

 results := make(map[string]*CheckResult)
 var wg sync.WaitGroup
 var mu sync.Mutex

 for name, check := range checks {
  wg.Add(1)
  go func(n string, c Checkable) {
   defer wg.Done()

   start := time.Now()
   err := c.Check(ctx)
   latency := time.Since(start)

   result := &CheckResult{
    Name:      n,
    Status:    StatusUp,
    Latency:   latency,
    Timestamp: time.Now(),
   }

   if err != nil {
    result.Status = StatusDown
    result.Message = err.Error()
   }

   mu.Lock()
   results[n] = result
   mu.Unlock()
  }(name, check)
 }

 wg.Wait()
 return results
}

// CheckResult 检查结果
type CheckResult struct {
 Name      string        `json:"name"`
 Status    HealthStatus  `json:"status"`
 Message   string        `json:"message,omitempty"`
 Latency   time.Duration `json:"latency"`
 Timestamp time.Time     `json:"timestamp"`
}

// DatabaseCheck 数据库检查
type DatabaseCheck struct {
 db      *sql.DB
 timeout time.Duration
}

func (d *DatabaseCheck) Name() string {
 return "database"
}

func (d *DatabaseCheck) Check(ctx context.Context) error {
 ctx, cancel := context.WithTimeout(ctx, d.timeout)
 defer cancel()
 return d.db.PingContext(ctx)
}

// CacheCheck 缓存检查
type CacheCheck struct {
 client  *redis.Client
 timeout time.Duration
}

func (c *CacheCheck) Name() string {
 return "cache"
}

func (c *CacheCheck) Check(ctx context.Context) error {
 ctx, cancel := context.WithTimeout(ctx, c.timeout)
 defer cancel()
 return c.client.Ping(ctx).Err()
}

// ExternalAPICheck 外部API检查
type ExternalAPICheck struct {
 client  *http.Client
 url     string
 timeout time.Duration
}

func (e *ExternalAPICheck) Name() string {
 return "external_api"
}

func (e *ExternalAPICheck) Check(ctx context.Context) error {
 ctx, cancel := context.WithTimeout(ctx, e.timeout)
 defer cancel()

 req, err := http.NewRequestWithContext(ctx, "GET", e.url, nil)
 if err != nil {
  return err
 }

 resp, err := e.client.Do(req)
 if err != nil {
  return err
 }
 defer resp.Body.Close()

 if resp.StatusCode >= 500 {
  return fmt.Errorf("external API returned status %d", resp.StatusCode)
 }

 return nil
}

// AdvancedHealthHandler 高级健康检查处理器
type AdvancedHealthHandler struct {
 registry *HealthRegistry
 version  string
}

// ServeHTTP 实现http.Handler
func (h *AdvancedHealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
 defer cancel()

 results := h.registry.RunChecks(ctx)

 // 计算总体状态
 overallStatus := StatusUp
 for _, result := range results {
  if result.Status == StatusDown {
   overallStatus = StatusDown
   break
  }
 }

 response := map[string]interface{}{
  "status":    overallStatus,
  "version":   h.version,
  "timestamp": time.Now(),
  "checks":    results,
 }

 w.Header().Set("Content-Type", "application/json")
 if overallStatus == StatusDown {
  w.WriteHeader(http.StatusServiceUnavailable)
 } else {
  w.WriteHeader(http.StatusOK)
 }

 json.NewEncoder(w).Encode(response)
}
```

### 7.6 反例说明

```go
// ❌ 错误示例1：健康检查依赖外部服务导致级联失败
func badCascadingFailure() {
 // 如果数据库挂了，应用也会被重启
 if err := db.Ping(); err != nil {
  w.WriteHeader(http.StatusServiceUnavailable)
  return
 }
}

// ❌ 错误示例2：健康检查执行时间过长
func badSlowHealthCheck() {
 // 执行复杂查询作为健康检查
 rows, _ := db.Query("SELECT * FROM large_table") // 可能执行很久
 // ...
}

// ❌ 错误示例3：没有超时控制
func badNoTimeout() {
 // 可能永远阻塞
 if err := db.Ping(); err != nil {
  // 没有超时，可能hang住
 }
}

// ❌ 错误示例4：健康检查端点无认证暴露
func badUnprotectedEndpoint() {
 http.HandleFunc("/health", healthHandler) // 任何人都可以访问
}

// ❌ 错误示例5：返回非标准格式
func badNonStandardResponse() {
 w.Write([]byte("OK")) // 不是JSON格式
}
```

### 7.7 最佳实践

```go
// ✅ 正确示例1：存活检查只检查自身
func goodLivenessCheck() {
 // 只检查应用是否响应，不依赖外部服务
 w.WriteHeader(http.StatusOK)
 json.NewEncoder(w).Encode(map[string]string{
  "status": "UP",
 })
}

// ✅ 正确示例2：就绪检查可以依赖外部服务
func goodReadinessCheck() {
 // 检查所有依赖
 checks := map[string]bool{
  "database": checkDB(),
  "cache":    checkCache(),
  "queue":    checkQueue(),
 }

 allReady := true
 for _, ready := range checks {
  if !ready {
   allReady = false
   break
  }
 }

 if allReady {
  w.WriteHeader(http.StatusOK)
 } else {
  w.WriteHeader(http.StatusServiceUnavailable)
 }
}

// ✅ 正确示例3：带超时的检查
func goodTimeoutCheck() {
 ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
 defer cancel()

 if err := db.PingContext(ctx); err != nil {
  w.WriteHeader(http.StatusServiceUnavailable)
  return
 }
}

// ✅ 正确示例4：异步缓存检查结果
var healthCache struct {
 result    map[string]interface{}
 timestamp time.Time
 mu        sync.RWMutex
}

func goodCachedHealthCheck() {
 healthCache.mu.RLock()
 if time.Since(healthCache.timestamp) < 5*time.Second {
  json.NewEncoder(w).Encode(healthCache.result)
  healthCache.mu.RUnlock()
  return
 }
 healthCache.mu.RUnlock()

 // 执行检查并缓存结果
 result := performHealthChecks()
 healthCache.mu.Lock()
 healthCache.result = result
 healthCache.timestamp = time.Now()
 healthCache.mu.Unlock()

 json.NewEncoder(w).Encode(result)
}

// ✅ 正确示例5：健康检查端点保护
func goodProtectedEndpoint() {
 // 只在内部网络暴露
 mux := http.NewServeMux()
 mux.HandleFunc("/healthz/live", livenessHandler)
 mux.HandleFunc("/healthz/ready", readinessHandler)

 // 绑定到localhost
 go http.ListenAndServe("localhost:8080", mux)

 // 公共端口只暴露业务端点
 publicMux := http.NewServeMux()
 publicMux.HandleFunc("/api/", apiHandler)
 http.ListenAndServe(":80", publicMux)
}

// ✅ 正确示例6：详细的健康检查响应
func goodDetailedResponse() {
 response := map[string]interface{}{
  "status":    "UP",
  "version":   "1.0.0",
  "timestamp": time.Now().Format(time.RFC3339),
  "checks": map[string]interface{}{
   "database": map[string]interface{}{
    "status":  "UP",
    "latency": "5ms",
   },
   "cache": map[string]interface{}{
    "status":  "UP",
    "latency": "1ms",
   },
  },
 }
 json.NewEncoder(w).Encode(response)
}
```

---

## 8. 总结与最佳实践

### 8.1 可观测性三支柱整合

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    可观测性三支柱整合                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                        ┌─────────────┐                                      │
│                        │ Application │                                      │
│                        └──────┬──────┘                                      │
│                               │                                              │
│           ┌───────────────────┼───────────────────┐                         │
│           │                   │                   │                         │
│           ▼                   ▼                   ▼                         │
│    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                   │
│    │    Trace    │    │   Metrics   │    │    Logs     │                   │
│    │             │    │             │    │             │                   │
│    │ • Request   │    │ • Counter   │    │ • Events    │                   │
│    │   flow      │    │ • Gauge     │    │ • Errors    │                   │
│    │ • Latency   │    │ • Histogram │    │ • Debug     │                   │
│    │ • Context   │    │ • Summary   │    │ • Audit     │                   │
│    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                   │
│           │                   │                   │                         │
│           └───────────────────┼───────────────────┘                         │
│                               │                                              │
│                               ▼                                              │
│                    ┌─────────────────────┐                                  │
│                    │  Correlation IDs    │                                  │
│                    │  • TraceID          │                                  │
│                    │  • SpanID           │                                  │
│                    │  • Timestamp        │                                  │
│                    └─────────────────────┘                                  │
│                               │                                              │
│                               ▼                                              │
│                    ┌─────────────────────┐                                  │
│                    │  Observability      │                                  │
│                    │  Platform           │                                  │
│                    │  (Jaeger/Prometheus │                                  │
│                    │   /Grafana/Loki)    │                                  │
│                    └─────────────────────┘                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Go可观测性最佳实践清单

#### Trace最佳实践

- [ ] 始终传播上下文，避免使用`context.Background()`
- [ ] 批量操作使用单个Span，避免Span爆炸
- [ ] 正确记录错误状态和错误信息
- [ ] 使用属性记录结构化数据，而非日志
- [ ] 合理设置采样率，生产环境使用概率采样
- [ ] 使用Baggage传递跨Span的上下文信息

#### Metrics最佳实践

- [ ] 预定义标签值，避免高基数问题
- [ ] Counter只递增，使用Gauge或UpDownCounter表示可减的值
- [ ] 合理设置Histogram桶边界
- [ ] 在初始化时注册指标，不在热路径创建
- [ ] 使用中间件模式自动记录HTTP指标
- [ ] 指标命名遵循`domain_subsystem_unit`格式

#### Logging最佳实践

- [ ] 使用结构化日志（zap/logrus）
- [ ] 正确的日志级别（Debug/Info/Warn/Error）
- [ ] 敏感信息脱敏处理
- [ ] 关联Trace上下文到日志
- [ ] 使用单例Logger，避免重复创建
- [ ] 高频日志使用采样

#### eBPF最佳实践

- [ ] 解除内存限制`rlimit.RemoveMemlock()`
- [ ] 使用`defer`确保资源释放
- [ ] 使用CO-RE技术保证内核兼容性
- [ ] Ring Buffer用于高效数据传输
- [ ] 批量处理事件减少系统调用

#### Profiling最佳实践

- [ ] 条件性启用剖析，避免生产环境持续采样
- [ ] 合理设置采样率（CPU默认100Hz）
- [ ] 内存剖析前先执行GC
- [ ] 定期保存和上传剖析数据
- [ ] 使用HTTP pprof端点便于实时分析
- [ ] 剖析端点绑定到localhost保护安全

#### Health Check最佳实践

- [ ] 存活检查只检查自身，不依赖外部服务
- [ ] 就绪检查可以依赖外部依赖
- [ ] 所有检查设置超时控制
- [ ] 异步缓存检查结果
- [ ] 返回详细的JSON格式响应
- [ ] 健康检查端点限制内部访问

### 8.3 推荐工具链

| 用途 | 工具 | 说明 |
|------|------|------|
| Trace收集 | Jaeger/Tempo | 分布式追踪后端 |
| Metrics收集 | Prometheus | 时序数据库 |
| 日志收集 | Loki/ELK | 日志聚合平台 |
| 可视化 | Grafana | 统一仪表盘 |
| eBPF开发 | cilium/ebpf | Go eBPF库 |
| 日志库 | zap/logrus | 结构化日志 |
| Profiling | pprof | Go内置剖析 |

---

*文档版本: 1.0*
*最后更新: 2024年*
*作者: Go可观测性专家*
