# OpenTelemetry W3C Trace Context 规范实现

> **分类**: 工程与云原生
> **标签**: #opentelemetry #w3c #trace-context #distributed-tracing
> **参考**: W3C Trace Context Specification, OpenTelemetry Specification

---

## W3C Trace Context 规范

### traceparent 格式

```
traceparent: version-trace_id-parent_id-flags

格式: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
       │  └────────────── trace_id ──────────────┘ │ │  │
       │                                           │ │  │
       version (2 hex) ────────────────────────────┘ │  │
                                                     │  │
       parent_id (16 hex) ───────────────────────────┘  │
                                                        │
       flags (2 hex) ───────────────────────────────────┘

字段说明:
- version (00-ff): 版本号，当前为 00
- trace_id (32 hex): 128-bit 追踪 ID
- parent_id (16 hex): 64-bit 父 Span ID
- flags (00-ff): 标志位
  - bit 0: sampled (1=采样, 0=未采样)
```

### tracestate 格式

```
tracestate: vendor1=value1,vendor2=value2,vendor3=value3

限制:
- 最大 32 个键值对
- 键名: [a-z0-9_-]{1,256}
- 值: 最大 256 字符
- 总长度: 最大 8192 字节
```

---

## 完整实现

```go
package tracecontext

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "regexp"
    "strings"
)

// TraceParent W3C traceparent 结构
type TraceParent struct {
    Version  byte
    TraceID  TraceID
    ParentID SpanID
    Flags    byte
}

// TraceID 128-bit 追踪 ID
type TraceID [16]byte

// SpanID 64-bit Span ID
type SpanID [8]byte

// TraceFlags 标志位
const (
    FlagSampled = 0x01
    FlagUnused1 = 0x02
    FlagUnused2 = 0x04
    FlagUnused3 = 0x08
    FlagUnused4 = 0x10
    FlagUnused5 = 0x20
    FlagUnused6 = 0x40
    FlagUnused7 = 0x80
)

// ParseTraceParent 解析 traceparent header
func ParseTraceParent(s string) (*TraceParent, error) {
    // 规范化：去除 whitespace
    s = strings.TrimSpace(s)

    // 格式验证正则
    const traceParentRegex = `^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$`
    matched, _ := regexp.MatchString(traceParentRegex, s)
    if !matched {
        return nil, fmt.Errorf("invalid traceparent format")
    }

    parts := strings.Split(s, "-")

    tp := &TraceParent{}

    // 解析 version
    version, err := hex.DecodeString(parts[0])
    if err != nil {
        return nil, fmt.Errorf("invalid version: %w", err)
    }
    tp.Version = version[0]

    // 版本 0x00 特殊处理：未来版本不兼容
    if tp.Version == 0xff {
        return nil, fmt.Errorf("version 0xff is invalid")
    }

    // 解析 trace_id
    traceID, err := hex.DecodeString(parts[1])
    if err != nil || len(traceID) != 16 {
        return nil, fmt.Errorf("invalid trace_id")
    }
    copy(tp.TraceID[:], traceID)

    // 验证 trace_id 非全零
    if isAllZero(tp.TraceID[:]) {
        return nil, fmt.Errorf("trace_id cannot be all zeros")
    }

    // 解析 parent_id
    parentID, err := hex.DecodeString(parts[2])
    if err != nil || len(parentID) != 8 {
        return nil, fmt.Errorf("invalid parent_id")
    }
    copy(tp.ParentID[:], parentID)

    // 验证 parent_id 非全零
    if isAllZero(tp.ParentID[:]) {
        return nil, fmt.Errorf("parent_id cannot be all zeros")
    }

    // 解析 flags
    flags, err := hex.DecodeString(parts[3])
    if err != nil || len(flags) != 1 {
        return nil, fmt.Errorf("invalid flags")
    }
    tp.Flags = flags[0]

    return tp, nil
}

// String 格式化 traceparent
func (tp *TraceParent) String() string {
    return fmt.Sprintf("%02x-%032x-%016x-%02x",
        tp.Version, tp.TraceID, tp.ParentID, tp.Flags)
}

// IsSampled 检查是否采样
func (tp *TraceParent) IsSampled() bool {
    return tp.Flags&FlagSampled == FlagSampled
}

// SetSampled 设置采样标志
func (tp *TraceParent) SetSampled(sampled bool) {
    if sampled {
        tp.Flags |= FlagSampled
    } else {
        tp.Flags &^= FlagSampled
    }
}

func isAllZero(b []byte) bool {
    for _, v := range b {
        if v != 0 {
            return false
        }
    }
    return true
}

// TraceState tracestate 结构
type TraceState struct {
    entries map[string]string
    order   []string // 保持插入顺序
}

func NewTraceState() *TraceState {
    return &TraceState{
        entries: make(map[string]string),
        order:   []string{},
    }
}

// ParseTraceState 解析 tracestate header
func ParseTraceState(s string) (*TraceState, error) {
    ts := NewTraceState()

    if strings.TrimSpace(s) == "" {
        return ts, nil
    }

    // 按逗号分割，考虑值中可能包含逗号
    entries := splitTraceState(s)

    for _, entry := range entries {
        entry = strings.TrimSpace(entry)
        if entry == "" {
            continue
        }

        parts := strings.SplitN(entry, "=", 2)
        if len(parts) != 2 {
            continue // 无效条目，跳过
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        if !isValidTraceStateKey(key) {
            continue
        }

        if !isValidTraceStateValue(value) {
            continue
        }

        ts.Set(key, value)
    }

    return ts, nil
}

func (ts *TraceState) Get(key string) (string, bool) {
    val, ok := ts.entries[key]
    return val, ok
}

func (ts *TraceState) Set(key, value string) {
    if _, exists := ts.entries[key]; !exists {
        ts.order = append(ts.order, key)
    }
    ts.entries[key] = value
}

func (ts *TraceState) Delete(key string) {
    delete(ts.entries, key)
    for i, k := range ts.order {
        if k == key {
            ts.order = append(ts.order[:i], ts.order[i+1:]...)
            break
        }
    }
}

// String 格式化 tracestate
func (ts *TraceState) String() string {
    var parts []string
    for _, key := range ts.order {
        if val, ok := ts.entries[key]; ok {
            parts = append(parts, fmt.Sprintf("%s=%s", key, val))
        }
    }
    return strings.Join(parts, ",")
}

// Len 返回条目数量
func (ts *TraceState) Len() int {
    return len(ts.entries)
}

func isValidTraceStateKey(key string) bool {
    if len(key) == 0 || len(key) > 256 {
        return false
    }
    // 键名规则: [a-z][a-z0-9_-]{0,255} 或 [a-z0-9][a-z0-9_-]{0,255}@\w+
    match, _ := regexp.MatchString(`^[a-z][a-z0-9_-]{0,255}$`, key)
    return match
}

func isValidTraceStateValue(value string) bool {
    if len(value) > 256 {
        return false
    }
    // 值规则: 可见字符，逗号需转义
    for _, r := range value {
        if r < 0x20 || r > 0x7e {
            return false
        }
    }
    return true
}

func splitTraceState(s string) []string {
    // 简单实现：按逗号分割
    return strings.Split(s, ",")
}
```

---

## SpanContext 与传播

```go
package tracecontext

import (
    "context"
    "fmt"
)

// SpanContext Span 上下文
type SpanContext struct {
    traceID    TraceID
    spanID     SpanID
    traceFlags byte
    traceState *TraceState
    remote     bool // 是否来自远程
}

func NewSpanContext(traceID TraceID, spanID SpanID, flags byte) SpanContext {
    return SpanContext{
        traceID:    traceID,
        spanID:     spanID,
        traceFlags: flags,
        traceState: NewTraceState(),
    }
}

func (sc SpanContext) TraceID() TraceID     { return sc.traceID }
func (sc SpanContext) SpanID() SpanID       { return sc.spanID }
func (sc SpanContext) TraceFlags() byte     { return sc.traceFlags }
func (sc SpanContext) TraceState() *TraceState { return sc.traceState }
func (sc SpanContext) IsRemote() bool       { return sc.remote }
func (sc SpanContext) IsValid() bool {
    return !isAllZero(sc.traceID[:]) && !isAllZero(sc.spanID[:])
}
func (sc SpanContext) IsSampled() bool {
    return sc.traceFlags&FlagSampled == FlagSampled
}

// SpanContextKey context key
type spanContextKey struct{}

func SpanFromContext(ctx context.Context) SpanContext {
    sc, _ := ctx.Value(spanContextKey{}).(SpanContext)
    return sc
}

func ContextWithSpan(ctx context.Context, sc SpanContext) context.Context {
    return context.WithValue(ctx, spanContextKey{}, sc)
}

// TraceContextPropagator W3C Trace Context 传播器
type TraceContextPropagator struct{}

func (p *TraceContextPropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
    sc := SpanFromContext(ctx)
    if !sc.IsValid() {
        return
    }

    // 注入 traceparent
    tp := &TraceParent{
        Version:  0,
        TraceID:  sc.TraceID(),
        ParentID: sc.SpanID(),
        Flags:    sc.TraceFlags(),
    }
    carrier.Set("traceparent", tp.String())

    // 注入 tracestate
    if ts := sc.TraceState(); ts != nil && ts.Len() > 0 {
        carrier.Set("tracestate", ts.String())
    }
}

func (p *TraceContextPropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
    // 提取 traceparent
    tpStr := carrier.Get("traceparent")
    if tpStr == "" {
        return ctx
    }

    tp, err := ParseTraceParent(tpStr)
    if err != nil {
        return ctx
    }

    sc := SpanContext{
        traceID:    tp.TraceID,
        spanID:     tp.ParentID,
        traceFlags: tp.Flags,
        traceState: NewTraceState(),
        remote:     true,
    }

    // 提取 tracestate
    if tsStr := carrier.Get("tracestate"); tsStr != "" {
        if ts, err := ParseTraceState(tsStr); err == nil {
            sc.traceState = ts
        }
    }

    return ContextWithSpan(ctx, sc)
}

func (p *TraceContextPropagator) Fields() []string {
    return []string{"traceparent", "tracestate"}
}

// TextMapCarrier 载体接口
type TextMapCarrier interface {
    Get(key string) string
    Set(key, value string)
}

// MapCarrier map 载体实现
type MapCarrier map[string]string

func (m MapCarrier) Get(key string) string {
    return m[key]
}

func (m MapCarrier) Set(key, value string) {
    m[key] = value
}

// HTTPHeadersCarrier HTTP 头载体
type HTTPHeadersCarrier struct {
    http.Header
}

func (c HTTPHeadersCarrier) Get(key string) string {
    return c.Header.Get(key)
}

func (c HTTPHeadersCarrier) Set(key, value string) {
    c.Header.Set(key, value)
}
```

---

## 采样策略

```go
package tracecontext

import (
    "math/rand"
)

// Sampler 采样器接口
type Sampler interface {
    ShouldSample(parameters SamplingParameters) SamplingResult
    Description() string
}

type SamplingParameters struct {
    ParentContext context.Context
    TraceID       TraceID
    Name          string
    Kind          SpanKind
}

type SamplingResult struct {
    Decision   SamplingDecision
    Attributes map[string]interface{}
    Tracestate *TraceState
}

type SamplingDecision int

const (
    Drop SamplingDecision = iota
    RecordOnly
    RecordAndSample
)

type SpanKind int

const (
    SpanKindUnspecified SpanKind = iota
    SpanKindInternal
    SpanKindServer
    SpanKindClient
    SpanKindProducer
    SpanKindConsumer
)

// AlwaysOnSampler 总是采样
var AlwaysOnSampler = &alwaysOnSampler{}

type alwaysOnSampler struct{}

func (s *alwaysOnSampler) ShouldSample(p SamplingParameters) SamplingResult {
    return SamplingResult{Decision: RecordAndSample}
}
func (s *alwaysOnSampler) Description() string { return "AlwaysOnSampler" }

// AlwaysOffSampler 从不采样
var AlwaysOffSampler = &alwaysOffSampler{}

type alwaysOffSampler struct{}

func (s *alwaysOffSampler) ShouldSample(p SamplingParameters) SamplingResult {
    return SamplingResult{Decision: Drop}
}
func (s *alwaysOffSampler) Description() string { return "AlwaysOffSampler" }

// TraceIDRatioBasedSampler 基于 TraceID 的比率采样
type TraceIDRatioBasedSampler struct {
    fraction float64
}

func NewTraceIDRatioBasedSampler(fraction float64) *TraceIDRatioBasedSampler {
    if fraction < 0 {
        fraction = 0
    }
    if fraction > 1 {
        fraction = 1
    }
    return &TraceIDRatioBasedSampler{fraction: fraction}
}

func (s *TraceIDRatioBasedSampler) ShouldSample(p SamplingParameters) SamplingResult {
    // 使用 TraceID 的最后 8 字节作为随机数源
    // 确保同一 TraceID 的采样结果一致（确定性采样）
    value := binary.BigEndian.Uint64(p.TraceID[8:])

    // 将 [0, MaxUint64) 映射到 [0, 1)
    threshold := uint64(s.fraction * float64(^uint64(0)))

    if value < threshold {
        return SamplingResult{Decision: RecordAndSample}
    }
    return SamplingResult{Decision: Drop}
}

func (s *TraceIDRatioBasedSampler) Description() string {
    return fmt.Sprintf("TraceIDRatioBased{%f}", s.fraction)
}

// ParentBasedSampler 基于父 Span 的采样决策
type ParentBasedSampler struct {
    root             Sampler
    remoteSampled    Sampler
    remoteNotSampled Sampler
    localSampled     Sampler
    localNotSampled  Sampler
}

type ParentBasedSamplerConfig struct {
    Root             Sampler
    RemoteSampled    Sampler
    RemoteNotSampled Sampler
    LocalSampled     Sampler
    LocalNotSampled  Sampler
}

func NewParentBasedSampler(cfg ParentBasedSamplerConfig) *ParentBasedSampler {
    if cfg.Root == nil {
        cfg.Root = AlwaysOnSampler
    }
    if cfg.RemoteSampled == nil {
        cfg.RemoteSampled = AlwaysOnSampler
    }
    if cfg.RemoteNotSampled == nil {
        cfg.RemoteNotSampled = AlwaysOffSampler
    }
    if cfg.LocalSampled == nil {
        cfg.LocalSampled = AlwaysOnSampler
    }
    if cfg.LocalNotSampled == nil {
        cfg.LocalNotSampled = AlwaysOffSampler
    }

    return &ParentBasedSampler{
        root:             cfg.Root,
        remoteSampled:    cfg.RemoteSampled,
        remoteNotSampled: cfg.RemoteNotSampled,
        localSampled:     cfg.LocalSampled,
        localNotSampled:  cfg.LocalNotSampled,
    }
}

func (s *ParentBasedSampler) ShouldSample(p SamplingParameters) SamplingResult {
    parent := SpanFromContext(p.ParentContext)

    // 没有父 Span，使用 root 采样器
    if !parent.IsValid() {
        return s.root.ShouldSample(p)
    }

    // 远程父 Span
    if parent.IsRemote() {
        if parent.IsSampled() {
            return s.remoteSampled.ShouldSample(p)
        }
        return s.remoteNotSampled.ShouldSample(p)
    }

    // 本地父 Span
    if parent.IsSampled() {
        return s.localSampled.ShouldSample(p)
    }
    return s.localNotSampled.ShouldSample(p)
}

func (s *ParentBasedSampler) Description() string {
    return "ParentBased"
}
```

---

## HTTP 集成中间件

```go
package tracecontext

import (
    "fmt"
    "net/http"
    "time"
)

// TracingMiddleware HTTP 追踪中间件
func TracingMiddleware(tracer Tracer, propagator Propagator, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 提取父上下文
        carrier := HTTPHeadersCarrier{r.Header}
        ctx = propagator.Extract(ctx, carrier)

        // 创建新 Span
        spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
        ctx, span := tracer.Start(ctx, spanName, WithSpanKind(SpanKindServer))
        defer span.End()

        // 设置 Span 属性
        span.SetAttributes(
            StringAttribute("http.method", r.Method),
            StringAttribute("http.url", r.URL.String()),
            StringAttribute("http.target", r.URL.Path),
            StringAttribute("http.host", r.Host),
            StringAttribute("http.scheme", r.URL.Scheme),
            StringAttribute("http.user_agent", r.UserAgent()),
        )

        // 包装 ResponseWriter 以捕获状态码
        rw := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

        // 执行处理
        start := time.Now()
        next.ServeHTTP(rw, r.WithContext(ctx))
        duration := time.Since(start)

        // 记录结果
        span.SetAttributes(
            IntAttribute("http.status_code", rw.statusCode),
            Int64Attribute("http.response_size", rw.written),
        )

        if rw.statusCode >= 400 {
            span.SetStatus(StatusError, fmt.Sprintf("HTTP %d", rw.statusCode))
        }

        span.SetAttributes(Int64Attribute("http.duration_ms", duration.Milliseconds()))
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    written    int64
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(p []byte) (n int, err error) {
    n, err = rr.ResponseWriter.Write(p)
    rr.written += int64(n)
    return
}

// HTTPClient 带追踪的 HTTP 客户端
type TracingHTTPClient struct {
    client     *http.Client
    tracer     Tracer
    propagator Propagator
}

func (c *TracingHTTPClient) Do(req *http.Request) (*http.Response, error) {
    ctx := req.Context()

    // 创建客户端 Span
    spanName := fmt.Sprintf("HTTP %s", req.Method)
    ctx, span := c.tracer.Start(ctx, spanName, WithSpanKind(SpanKindClient))
    defer span.End()

    // 设置属性
    span.SetAttributes(
        StringAttribute("http.method", req.Method),
        StringAttribute("http.url", req.URL.String()),
    )

    // 注入追踪上下文
    carrier := HTTPHeadersCarrier{req.Header}
    c.propagator.Inject(ctx, carrier)

    // 执行请求
    resp, err := c.client.Do(req.WithContext(ctx))
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    span.SetAttributes(IntAttribute("http.status_code", resp.StatusCode))

    return resp, nil
}
```

---

## Baggage 传播

```go
package tracecontext

import (
    "context"
    "net/url"
    "strings"
)

// Baggage 键值对元数据传播
type Baggage struct {
    entries map[string]BaggageMember
}

type BaggageMember struct {
    Value    string
    Metadata map[string]string
}

func NewBaggage() *Baggage {
    return &Baggage{
        entries: make(map[string]BaggageMember),
    }
}

func (b *Baggage) Member(key string) (BaggageMember, bool) {
    m, ok := b.entries[key]
    return m, ok
}

func (b *Baggage) SetMember(key string, member BaggageMember) {
    b.entries[key] = member
}

func (b *Baggage) Get(key string) string {
    if m, ok := b.entries[key]; ok {
        return m.Value
    }
    return ""
}

func (b *Baggage) Set(key, value string) {
    b.entries[key] = BaggageMember{Value: value}
}

func (b *Baggage) Len() int {
    return len(b.entries)
}

// String W3C baggage 格式
func (b *Baggage) String() string {
    var parts []string
    for key, member := range b.entries {
        // URL 编码
        k := url.QueryEscape(key)
        v := url.QueryEscape(member.Value)

        part := fmt.Sprintf("%s=%s", k, v)

        // 添加元数据
        for mk, mv := range member.Metadata {
            part += fmt.Sprintf(";%s=%s", mk, url.QueryEscape(mv))
        }

        parts = append(parts, part)
    }
    return strings.Join(parts, ",")
}

// ParseBaggage 解析 baggage header
func ParseBaggage(s string) (*Baggage, error) {
    b := NewBaggage()

    if strings.TrimSpace(s) == "" {
        return b, nil
    }

    entries := strings.Split(s, ",")

    for _, entry := range entries {
        entry = strings.TrimSpace(entry)
        if entry == "" {
            continue
        }

        // 解析 key=value;metadata
        parts := strings.Split(entry, ";")

        kv := strings.SplitN(parts[0], "=", 2)
        if len(kv) != 2 {
            continue
        }

        key, _ := url.QueryUnescape(kv[0])
        value, _ := url.QueryUnescape(kv[1])

        member := BaggageMember{
            Value:    value,
            Metadata: make(map[string]string),
        }

        // 解析元数据
        for i := 1; i < len(parts); i++ {
            mkv := strings.SplitN(parts[i], "=", 2)
            if len(mkv) == 2 {
                mk, _ := url.QueryUnescape(mkv[0])
                mv, _ := url.QueryUnescape(mkv[1])
                member.Metadata[mk] = mv
            }
        }

        b.SetMember(key, member)
    }

    return b, nil
}

// baggageKey context key
type baggageKey struct{}

func BaggageFromContext(ctx context.Context) *Baggage {
    b, _ := ctx.Value(baggageKey{}).(*Baggage)
    return b
}

func ContextWithBaggage(ctx context.Context, baggage *Baggage) context.Context {
    return context.WithValue(ctx, baggageKey{}, baggage)
}

// BaggagePropagator baggage 传播器
type BaggagePropagator struct{}

func (p *BaggagePropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
    baggage := BaggageFromContext(ctx)
    if baggage == nil || baggage.Len() == 0 {
        return
    }

    carrier.Set("baggage", baggage.String())
}

func (p *BaggagePropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
    val := carrier.Get("baggage")
    if val == "" {
        return ctx
    }

    baggage, _ := ParseBaggage(val)
    return ContextWithBaggage(ctx, baggage)
}

func (p *BaggagePropagator) Fields() []string {
    return []string{"baggage"}
}
```
