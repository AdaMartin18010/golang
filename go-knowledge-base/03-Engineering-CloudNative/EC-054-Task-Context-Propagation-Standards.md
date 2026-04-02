# 任务上下文传播标准 (Task Context Propagation Standards)

> **分类**: 工程与云原生
> **标签**: #standards #w3c #opentelemetry #interop

---

## W3C Trace Context 实现

```go
// W3C Trace Context 标准实现
// https://www.w3.org/TR/trace-context/

package tracecontext

const (
    TraceParentHeader = "traceparent"
    TraceStateHeader  = "tracestate"
)

// TraceParent 格式: version-trace_id-parent_id-flags
type TraceParent struct {
    Version  string
    TraceID  string  // 16 bytes (32 hex chars)
    ParentID string  // 8 bytes (16 hex chars)
    Flags    byte    // 1 byte (2 hex chars)
}

func (tp *TraceParent) String() string {
    return fmt.Sprintf("%s-%s-%s-%02x",
        tp.Version,
        tp.TraceID,
        tp.ParentID,
        tp.Flags,
    )
}

func ParseTraceParent(s string) (*TraceParent, error) {
    parts := strings.Split(s, "-")
    if len(parts) != 4 {
        return nil, fmt.Errorf("invalid traceparent format")
    }

    flags, err := strconv.ParseUint(parts[3], 16, 8)
    if err != nil {
        return nil, fmt.Errorf("invalid flags: %w", err)
    }

    return &TraceParent{
        Version:  parts[0],
        TraceID:  parts[1],
        ParentID: parts[2],
        Flags:    byte(flags),
    }, nil
}

// 注入 HTTP 头
func InjectHTTP(ctx context.Context, req *http.Request) {
    if span := trace.SpanFromContext(ctx); span != nil {
        spanContext := span.SpanContext()

        if spanContext.IsValid() {
            tp := TraceParent{
                Version:  "00",
                TraceID:  spanContext.TraceID().String(),
                ParentID: spanContext.SpanID().String(),
                Flags:    byte(spanContext.TraceFlags()),
            }

            req.Header.Set(TraceParentHeader, tp.String())

            // 注入 Trace State
            if tracestate := spanContext.TraceState().String(); tracestate != "" {
                req.Header.Set(TraceStateHeader, tracestate)
            }
        }
    }
}

// 从 HTTP 头提取
func ExtractHTTP(ctx context.Context, req *http.Request) context.Context {
    traceparent := req.Header.Get(TraceParentHeader)
    if traceparent == "" {
        return ctx
    }

    tp, err := ParseTraceParent(traceparent)
    if err != nil {
        return ctx
    }

    traceID, _ := trace.TraceIDFromHex(tp.TraceID)
    spanID, _ := trace.SpanIDFromHex(tp.ParentID)

    spanContext := trace.NewSpanContext(trace.SpanContextConfig{
        TraceID:    traceID,
        SpanID:     spanID,
        TraceFlags: trace.TraceFlags(tp.Flags),
        Remote:     true,
    })

    // 提取 Trace State
    tracestate := req.Header.Get(TraceStateHeader)
    if tracestate != "" {
        ts, _ := trace.ParseTraceState(tracestate)
        spanContext = spanContext.WithTraceState(ts)
    }

    return trace.ContextWithSpanContext(ctx, spanContext)
}
```

---

## Baggage 标准

```go
// W3C Baggage 标准实现
// https://www.w3.org/TR/baggage/

const BaggageHeader = "baggage"

// Baggage 管理器
type BaggageManager struct {
    maxEntries int
    maxBytes   int
}

func (bm *BaggageManager) Inject(ctx context.Context, req *http.Request) {
    baggage := baggage.FromContext(ctx)
    if baggage.Len() == 0 {
        return
    }

    var parts []string
    baggage.Iterate(func(m baggage.Member) bool {
        part := fmt.Sprintf("%s=%s", url.QueryEscape(m.Key()), url.QueryEscape(m.Value()))
        if m.Metadata() != "" {
            part += ";" + m.Metadata()
        }
        parts = append(parts, part)
        return true
    })

    req.Header.Set(BaggageHeader, strings.Join(parts, ","))
}

func (bm *BaggageManager) Extract(ctx context.Context, req *http.Request) context.Context {
    header := req.Header.Get(BaggageHeader)
    if header == "" {
        return ctx
    }

    members := []baggage.Member{}

    for _, part := range strings.Split(header, ",") {
        part = strings.TrimSpace(part)
        if part == "" {
            continue
        }

        // 解析 key=value;metadata
        kv, _, _ := strings.Cut(part, ";")
        key, value, ok := strings.Cut(kv, "=")
        if !ok {
            continue
        }

        key, _ = url.QueryUnescape(key)
        value, _ = url.QueryUnescape(value)

        member, _ := baggage.NewMember(key, value)
        members = append(members, member)
    }

    b, _ := baggage.New(members...)
    return baggage.ContextWithBaggage(ctx, b)
}
```

---

## 跨语言兼容性

```go
// 与不同语言的互操作

type InteropPropagator struct {
    format PropagatorFormat
}

type PropagatorFormat int

const (
    FormatW3C PropagatorFormat = iota
    FormatB3
    FormatJaeger
    FormatOTTrace
)

func (ip *InteropPropagator) Inject(ctx context.Context, req *http.Request) {
    switch ip.format {
    case FormatW3C:
        ip.injectW3C(ctx, req)
    case FormatB3:
        ip.injectB3(ctx, req)
    case FormatJaeger:
        ip.injectJaeger(ctx, req)
    }
}

// B3 格式 (Zipkin)
func (ip *InteropPropagator) injectB3(ctx context.Context, req *http.Request) {
    span := trace.SpanFromContext(ctx)
    sc := span.SpanContext()

    if sc.IsValid() {
        req.Header.Set("X-B3-TraceId", sc.TraceID().String())
        req.Header.Set("X-B3-SpanId", sc.SpanID().String())
        if sc.IsSampled() {
            req.Header.Set("X-B3-Sampled", "1")
        }
    }
}

// Jaeger 格式
func (ip *InteropPropagator) injectJaeger(ctx context.Context, req *http.Request) {
    span := trace.SpanFromContext(ctx)
    sc := span.SpanContext()

    if sc.IsValid() {
        // Jaeger 使用 uber-trace-id 头
        // 格式: {trace-id}:{span-id}:{parent-span-id}:{flags}
        uberTraceID := fmt.Sprintf("%s:%s:0:%d",
            sc.TraceID().String(),
            sc.SpanID().String(),
            sc.TraceFlags(),
        )
        req.Header.Set("uber-trace-id", uberTraceID)
    }
}
```

---

## 协议桥接

```go
// 在不同协议间传播上下文

type ProtocolBridge struct {
    propagators map[string]Propagator
}

type Propagator interface {
    Inject(ctx context.Context, carrier interface{}) error
    Extract(ctx context.Context, carrier interface{}) (context.Context, error)
}

// HTTP 到 gRPC 桥接
func (pb *ProtocolBridge) HTTPToGRPC(ctx context.Context, req *http.Request) context.Context {
    // 从 HTTP 提取
    ctx = ExtractHTTP(ctx, req)

    // 转换为 gRPC metadata
    md := metadata.MD{}

    if tp := req.Header.Get(TraceParentHeader); tp != "" {
        md.Set(TraceParentHeader, tp)
    }

    if ts := req.Header.Get(TraceStateHeader); ts != "" {
        md.Set(TraceStateHeader, ts)
    }

    return metadata.NewOutgoingContext(ctx, md)
}

// gRPC 到消息队列桥接
func (pb *ProtocolBridge) GRPCToMessage(ctx context.Context, msg *Message) context.Context {
    md, _ := metadata.FromIncomingContext(ctx)

    // 提取 trace 信息
    if tp := md.Get(TraceParentHeader); len(tp) > 0 {
        msg.Headers[TraceParentHeader] = tp[0]
    }

    // 提取 baggage
    if baggage := baggage.FromContext(ctx); baggage.Len() > 0 {
        var parts []string
        baggage.Iterate(func(m baggage.Member) bool {
            parts = append(parts, fmt.Sprintf("%s=%s", m.Key(), m.Value()))
            return true
        })
        msg.Headers[BaggageHeader] = strings.Join(parts, ",")
    }

    return ctx
}
```
