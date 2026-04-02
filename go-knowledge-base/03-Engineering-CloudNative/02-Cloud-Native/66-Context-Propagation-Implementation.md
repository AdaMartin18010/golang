# 上下文传播实现机制 (Context Propagation Implementation)

> **分类**: 工程与云原生
> **标签**: #context-propagation #distributed-systems #implementation #w3c
> **参考**: W3C Trace Context, OpenTelemetry Go SDK

---

## 上下文传播基础

### 传播器接口设计

```go
package propagation

import "context"

// TextMapCarrier 是用于传播的键值对载体
type TextMapCarrier interface {
    Get(key string) string
    Set(key, value string)
    Keys() []string
}

// HTTPHeadersCarrier 包装 http.Header 实现 TextMapCarrier
type HTTPHeadersCarrier struct {
    http.Header
}

func (c HTTPHeadersCarrier) Get(key string) string {
    return c.Header.Get(key)
}

func (c HTTPHeadersCarrier) Set(key, value string) {
    c.Header.Set(key, value)
}

func (c HTTPHeadersCarrier) Keys() []string {
    keys := make([]string, 0, len(c.Header))
    for k := range c.Header {
        keys = append(keys, k)
    }
    return keys
}

// Propagator 定义了上下文传播接口
type Propagator interface {
    // Inject 将上下文注入载体
    Inject(ctx context.Context, carrier TextMapCarrier)
    // Extract 从载体提取上下文
    Extract(ctx context.Context, carrier TextMapCarrier) context.Context
    // Fields 返回该传播器使用的字段名
    Fields() []string
}
```

### W3C Trace Context 实现

```go
package propagation

import (
    "context"
    "encoding/hex"
    "fmt"
    "strings"
)

const (
    TraceParentHeader = "traceparent"
    TraceStateHeader  = "tracestate"
)

// TraceParent W3C traceparent 格式
type TraceParent struct {
    Version  byte
    TraceID  [16]byte
    SpanID   [8]byte
    Flags    byte
}

func (tp *TraceParent) String() string {
    return fmt.Sprintf("%02x-%032x-%016x-%02x",
        tp.Version, tp.TraceID, tp.SpanID, tp.Flags)
}

func ParseTraceParent(s string) (*TraceParent, error) {
    parts := strings.Split(s, "-")
    if len(parts) != 4 {
        return nil, fmt.Errorf("invalid traceparent format")
    }

    tp := &TraceParent{}

    // 解析 version
    v, err := hex.DecodeString(parts[0])
    if err != nil || len(v) != 1 {
        return nil, fmt.Errorf("invalid version")
    }
    tp.Version = v[0]

    // 解析 trace_id
    traceID, err := hex.DecodeString(parts[1])
    if err != nil || len(traceID) != 16 {
        return nil, fmt.Errorf("invalid trace_id")
    }
    copy(tp.TraceID[:], traceID)

    // 解析 span_id
    spanID, err := hex.DecodeString(parts[2])
    if err != nil || len(spanID) != 8 {
        return nil, fmt.Errorf("invalid span_id")
    }
    copy(tp.SpanID[:], spanID)

    // 解析 flags
    flags, err := hex.DecodeString(parts[3])
    if err != nil || len(flags) != 1 {
        return nil, fmt.Errorf("invalid flags")
    }
    tp.Flags = flags[0]

    return tp, nil
}

// TraceContextPropagator W3C Trace Context 传播器
type TraceContextPropagator struct{}

func (p *TraceContextPropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
    span := SpanFromContext(ctx)
    if span == nil {
        return
    }

    spanContext := span.SpanContext()
    if !spanContext.IsValid() {
        return
    }

    // 构造 traceparent
    tp := &TraceParent{
        Version: 0,
        Flags:   0,
    }
    if spanContext.IsSampled() {
        tp.Flags |= 1
    }

    copy(tp.TraceID[:], spanContext.TraceID()[:])
    copy(tp.SpanID[:], spanContext.SpanID()[:])

    carrier.Set(TraceParentHeader, tp.String())

    // 注入 tracestate
    if tracestate := spanContext.TraceState(); tracestate != "" {
        carrier.Set(TraceStateHeader, tracestate)
    }
}

func (p *TraceContextPropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
    // 提取 traceparent
    tpStr := carrier.Get(TraceParentHeader)
    if tpStr == "" {
        return ctx
    }

    tp, err := ParseTraceParent(tpStr)
    if err != nil {
        return ctx
    }

    // 构造 SpanContext
    sc := SpanContext{
        traceID: tp.TraceID,
        spanID:  tp.SpanID,
        flags:   tp.Flags,
        remote:  true,
    }

    // 提取 tracestate
    if ts := carrier.Get(TraceStateHeader); ts != "" {
        sc.traceState = ts
    }

    return ContextWithSpanContext(ctx, sc)
}

func (p *TraceContextPropagator) Fields() []string {
    return []string{TraceParentHeader, TraceStateHeader}
}
```

---

## Baggage 传播实现

```go
package propagation

import (
    "context"
    "net/url"
    "strings"
)

const BaggageHeader = "baggage"

// Baggage 键值对存储
type Baggage struct {
    values map[string]BaggageMember
}

type BaggageMember struct {
    Value    string
    Metadata map[string]string
}

func NewBaggage() *Baggage {
    return &Baggage{
        values: make(map[string]BaggageMember),
    }
}

func (b *Baggage) Set(key, value string) {
    b.values[key] = BaggageMember{Value: value}
}

func (b *Baggage) Get(key string) (string, bool) {
    m, ok := b.values[key]
    if !ok {
        return "", false
    }
    return m.Value, true
}

func (b *Baggage) String() string {
    var parts []string
    for k, v := range b.values {
        // URL 编码
        key := url.QueryEscape(k)
        val := url.QueryEscape(v.Value)
        parts = append(parts, key+"="+val)
    }
    return strings.Join(parts, ",")
}

func ParseBaggage(s string) *Baggage {
    b := NewBaggage()

    pairs := strings.Split(s, ",")
    for _, p := range pairs {
        kv := strings.SplitN(p, "=", 2)
        if len(kv) != 2 {
            continue
        }

        key, _ := url.QueryUnescape(strings.TrimSpace(kv[0]))
        val, _ := url.QueryUnescape(strings.TrimSpace(kv[1]))
        b.Set(key, val)
    }

    return b
}

// BaggagePropagator baggage 传播器
type BaggagePropagator struct{}

func (p *BaggagePropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
    baggage := BaggageFromContext(ctx)
    if baggage == nil || len(baggage.values) == 0 {
        return
    }

    carrier.Set(BaggageHeader, baggage.String())
}

func (p *BaggagePropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
    val := carrier.Get(BaggageHeader)
    if val == "" {
        return ctx
    }

    baggage := ParseBaggage(val)
    return ContextWithBaggage(ctx, baggage)
}

func (p *BaggagePropagator) Fields() []string {
    return []string{BaggageHeader}
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
```

---

## 复合传播器

```go
package propagation

import "context"

// CompositePropagator 组合多个传播器
type CompositePropagator struct {
    propagators []Propagator
}

func NewCompositePropagator(propagators ...Propagator) *CompositePropagator {
    return &CompositePropagator{
        propagators: propagators,
    }
}

func (p *CompositePropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
    for _, prop := range p.propagators {
        prop.Inject(ctx, carrier)
    }
}

func (p *CompositePropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
    for _, prop := range p.propagators {
        ctx = prop.Extract(ctx, carrier)
    }
    return ctx
}

func (p *CompositePropagator) Fields() []string {
    fields := make([]string, 0)
    seen := make(map[string]bool)

    for _, prop := range p.propagators {
        for _, f := range prop.Fields() {
            if !seen[f] {
                seen[f] = true
                fields = append(fields, f)
            }
        }
    }

    return fields
}

// 默认传播器
func DefaultPropagator() Propagator {
    return NewCompositePropagator(
        &TraceContextPropagator{},
        &BaggagePropagator{},
    )
}
```

---

## HTTP/gRPC 集成

```go
package propagation

import (
    "net/http"

    "google.golang.org/grpc/metadata"
)

// HTTPMiddleware HTTP 传播中间件
func HTTPMiddleware(propagator Propagator) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 提取上游上下文
            carrier := HTTPHeadersCarrier{r.Header}
            ctx := propagator.Extract(r.Context(), carrier)

            // 创建响应包装器以注入传播字段
            rw := &responseWriter{
                ResponseWriter: w,
                header:         w.Header(),
            }

            // 注入响应头
            defer func() {
                propagator.Inject(ctx, HTTPHeadersCarrier{rw.header})
            }()

            next.ServeHTTP(rw, r.WithContext(ctx))
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    header http.Header
}

func (rw *responseWriter) Header() http.Header {
    return rw.header
}

// GRPCUnaryInterceptor gRPC 一元拦截器
func GRPCUnaryInterceptor(propagator Propagator) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // 从 gRPC metadata 提取
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            md = metadata.New(nil)
        }

        carrier := GRPCMetadataCarrier(md)
        ctx = propagator.Extract(ctx, carrier)

        return handler(ctx, req)
    }
}

// GRPCUnaryClientInterceptor gRPC 客户端拦截器
func GRPCUnaryClientInterceptor(propagator Propagator) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        // 注入到 metadata
        md, ok := metadata.FromOutgoingContext(ctx)
        if !ok {
            md = metadata.New(nil)
        }

        md = md.Copy()
        carrier := GRPCMetadataCarrier(md)
        propagator.Inject(ctx, carrier)

        ctx = metadata.NewOutgoingContext(ctx, md)
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}

// GRPCMetadataCarrier 适配 gRPC metadata
type GRPCMetadataCarrier metadata.MD

func (m GRPCMetadataCarrier) Get(key string) string {
    vals := metadata.MD(m).Get(key)
    if len(vals) == 0 {
        return ""
    }
    return vals[0]
}

func (m GRPCMetadataCarrier) Set(key, value string) {
    metadata.MD(m).Set(key, value)
}

func (m GRPCMetadataCarrier) Keys() []string {
    keys := make([]string, 0, len(m))
    for k := range metadata.MD(m) {
        keys = append(keys, k)
    }
    return keys
}
```

---

## 消息队列上下文传播

```go
package propagation

import (
    "context"
    "encoding/base64"
)

// MessageCarrier 消息载体接口
type MessageCarrier interface {
    GetHeader(key string) string
    SetHeader(key, value string)
    GetProperties() map[string]string
}

// KafkaCarrier Kafka 消息载体
type KafkaCarrier struct {
    Headers []KafkaHeader
}

type KafkaHeader struct {
    Key   string
    Value []byte
}

func (c *KafkaCarrier) GetHeader(key string) string {
    for _, h := range c.Headers {
        if h.Key == key {
            return string(h.Value)
        }
    }
    return ""
}

func (c *KafkaCarrier) SetHeader(key, value string) {
    // 查找是否已存在
    for i, h := range c.Headers {
        if h.Key == key {
            c.Headers[i].Value = []byte(value)
            return
        }
    }
    // 添加新 header
    c.Headers = append(c.Headers, KafkaHeader{
        Key:   key,
        Value: []byte(value),
    })
}

// AMQPCarrier AMQP (RabbitMQ) 载体
type AMQPCarrier struct {
    Headers map[string]interface{}
}

func (c *AMQPCarrier) GetHeader(key string) string {
    if c.Headers == nil {
        return ""
    }
    val, ok := c.Headers[key]
    if !ok {
        return ""
    }
    switch v := val.(type) {
    case string:
        return v
    case []byte:
        return string(v)
    default:
        return ""
    }
}

func (c *AMQPCarrier) SetHeader(key, value string) {
    if c.Headers == nil {
        c.Headers = make(map[string]interface{})
    }
    c.Headers[key] = value
}

// 生产消息时注入上下文
func InjectIntoMessage(ctx context.Context, propagator Propagator, carrier MessageCarrier) {
    // 使用 TextMapCarrier 适配
    textCarrier := &messageAdapter{carrier}
    propagator.Inject(ctx, textCarrier)
}

// 消费消息时提取上下文
func ExtractFromMessage(propagator Propagator, carrier MessageCarrier) context.Context {
    ctx := context.Background()
    textCarrier := &messageAdapter{carrier}
    return propagator.Extract(ctx, textCarrier)
}

type messageAdapter struct {
    MessageCarrier
}

func (m *messageAdapter) Get(key string) string {
    return m.GetHeader(key)
}

func (m *messageAdapter) Set(key, value string) {
    m.SetHeader(key, value)
}

func (m *messageAdapter) Keys() []string {
    props := m.GetProperties()
    keys := make([]string, 0, len(props))
    for k := range props {
        keys = append(keys, k)
    }
    return keys
}
```

---

## 上下文值传递最佳实践

```go
package propagation

import (
    "context"
    "fmt"
)

// 类型安全的上下文键
type contextKey struct {
    name string
}

func (k contextKey) String() string {
    return k.name
}

// 预定义键
var (
    requestIDKey = contextKey{"request-id"}
    tenantIDKey  = contextKey{"tenant-id"}
    userIDKey    = contextKey{"user-id"}
    traceIDKey   = contextKey{"trace-id"}
)

// RequestID 操作
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}

// 安全获取（带默认值）
func RequestIDOrEmpty(ctx context.Context) string {
    if id, ok := RequestIDFromContext(ctx); ok {
        return id
    }
    return ""
}

// TenantID 操作
func WithTenantID(ctx context.Context, tenantID string) context.Context {
    return context.WithValue(ctx, tenantIDKey, tenantID)
}

func TenantIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(tenantIDKey).(string)
    return id, ok
}

// 上下文克隆（用于异步操作）
func CloneContext(ctx context.Context) context.Context {
    newCtx := context.Background()

    // 复制关键值
    if id, ok := RequestIDFromContext(ctx); ok {
        newCtx = WithRequestID(newCtx, id)
    }
    if tenant, ok := TenantIDFromContext(ctx); ok {
        newCtx = WithTenantID(newCtx, tenant)
    }
    if baggage := BaggageFromContext(ctx); baggage != nil {
        newCtx = ContextWithBaggage(newCtx, baggage)
    }

    return newCtx
}

// 上下文转 map（用于日志）
func ContextToMap(ctx context.Context) map[string]string {
    m := make(map[string]string)

    if id, ok := RequestIDFromContext(ctx); ok {
        m["request_id"] = id
    }
    if tenant, ok := TenantIDFromContext(ctx); ok {
        m["tenant_id"] = tenant
    }

    return m
}

// 验证必需字段
func ValidateContext(ctx context.Context, required ...contextKey) error {
    for _, key := range required {
        if ctx.Value(key) == nil {
            return fmt.Errorf("missing required context value: %s", key)
        }
    }
    return nil
}
```

---

## 性能优化

```go
package propagation

import (
    "sync"
)

// 缓存 TraceParent 解析结果
type traceParentCache struct {
    mu    sync.RWMutex
    cache map[string]*TraceParent
}

var tpCache = &traceParentCache{
    cache: make(map[string]*TraceParent),
}

func (c *traceParentCache) Get(s string) (*TraceParent, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    tp, ok := c.cache[s]
    return tp, ok
}

func (c *traceParentCache) Put(s string, tp *TraceParent) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[s] = tp
}

// 优化的 TraceParent 解析
func FastParseTraceParent(s string) (*TraceParent, error) {
    // 先查缓存
    if tp, ok := tpCache.Get(s); ok {
        return tp, nil
    }

    tp, err := ParseTraceParent(s)
    if err != nil {
        return nil, err
    }

    // 放入缓存
    tpCache.Put(s, tp)
    return tp, nil
}

// 对象池复用 carrier
var carrierPool = sync.Pool{
    New: func() interface{} {
        return make(map[string]string)
    },
}

func GetCarrier() map[string]string {
    return carrierPool.Get().(map[string]string)
}

func PutCarrier(m map[string]string) {
    for k := range m {
        delete(m, k)
    }
    carrierPool.Put(m)
}
```
