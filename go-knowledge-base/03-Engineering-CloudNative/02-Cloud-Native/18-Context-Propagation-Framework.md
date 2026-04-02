# 上下文传播框架 (Context Propagation Framework)

> **分类**: 工程与云原生  
> **标签**: #context-propagation #distributed-tracing #observability

---

## 传播机制

```go
// Propagator 接口
type Propagator interface {
    // 注入上下文到载体
    Inject(ctx context.Context, carrier Carrier)
    // 从载体提取上下文
    Extract(ctx context.Context, carrier Carrier) context.Context
}

// Carrier 载体接口
type Carrier interface {
    Get(key string) string
    Set(key string, value string)
    Keys() []string
}

// 传播器注册表
type PropagatorRegistry struct {
    propagators map[string]Propagator
}

func (pr *PropagatorRegistry) Register(name string, p Propagator) {
    pr.propagators[name] = p
}

func (pr *PropagatorRegistry) InjectAll(ctx context.Context, carrier Carrier) {
    for _, p := range pr.propagators {
        p.Inject(ctx, carrier)
    }
}

func (pr *PropagatorRegistry) ExtractAll(ctx context.Context, carrier Carrier) context.Context {
    for _, p := range pr.propagators {
        ctx = p.Extract(ctx, carrier)
    }
    return ctx
}
```

---

## HTTP 传播

```go
type HTTPPropagator struct {
    keys []string
}

func (hp *HTTPPropagator) Inject(ctx context.Context, carrier Carrier) {
    // 注入请求ID
    if reqID := RequestIDFromContext(ctx); reqID != "" {
        carrier.Set("X-Request-ID", reqID)
    }
    
    // 注入Trace信息
    if traceID := TraceIDFromContext(ctx); traceID != "" {
        carrier.Set("X-Trace-ID", traceID)
        carrier.Set("X-Span-ID", SpanIDFromContext(ctx))
    }
    
    // 注入Baggage（业务上下文）
    if baggage := BaggageFromContext(ctx); len(baggage) > 0 {
        for k, v := range baggage {
            carrier.Set(fmt.Sprintf("X-Baggage-%s", k), v)
        }
    }
}

func (hp *HTTPPropagator) Extract(ctx context.Context, carrier Carrier) context.Context {
    // 提取请求ID
    if reqID := carrier.Get("X-Request-ID"); reqID != "" {
        ctx = WithRequestID(ctx, reqID)
    }
    
    // 提取Trace信息
    if traceID := carrier.Get("X-Trace-ID"); traceID != "" {
        ctx = WithTraceID(ctx, traceID)
        ctx = WithSpanID(ctx, carrier.Get("X-Span-ID"))
    }
    
    // 提取Baggage
    baggage := make(map[string]string)
    for _, key := range carrier.Keys() {
        if strings.HasPrefix(key, "X-Baggage-") {
            k := strings.TrimPrefix(key, "X-Baggage-")
            baggage[k] = carrier.Get(key)
        }
    }
    if len(baggage) > 0 {
        ctx = WithBaggage(ctx, baggage)
    }
    
    return ctx
}

// HTTP Carrier 实现
type HTTPHeadersCarrier http.Header

func (c HTTPHeadersCarrier) Get(key string) string {
    return http.Header(c).Get(key)
}

func (c HTTPHeadersCarrier) Set(key string, value string) {
    http.Header(c).Set(key, value)
}

func (c HTTPHeadersCarrier) Keys() []string {
    keys := make([]string, 0, len(c))
    for k := range c {
        keys = append(keys, k)
    }
    return keys
}
```

---

## gRPC 传播

```go
type GRPCPropagator struct{}

func (gp *GRPCPropagator) Inject(ctx context.Context, carrier Carrier) {
    md, ok := metadata.FromOutgoingContext(ctx)
    if !ok {
        md = metadata.New(nil)
    }
    
    // 注入上下文到 gRPC metadata
    if reqID := RequestIDFromContext(ctx); reqID != "" {
        md.Set("request-id", reqID)
    }
    
    if traceID := TraceIDFromContext(ctx); traceID != "" {
        md.Set("trace-id", traceID)
        md.Set("span-id", SpanIDFromContext(ctx))
    }
    
    // 更新 context
    ctx = metadata.NewOutgoingContext(ctx, md)
}

func (gp *GRPCPropagator) Extract(ctx context.Context, carrier Carrier) context.Context {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return ctx
    }
    
    if vals := md.Get("request-id"); len(vals) > 0 {
        ctx = WithRequestID(ctx, vals[0])
    }
    
    if vals := md.Get("trace-id"); len(vals) > 0 {
        ctx = WithTraceID(ctx, vals[0])
        if spanVals := md.Get("span-id"); len(spanVals) > 0 {
            ctx = WithSpanID(ctx, spanVals[0])
        }
    }
    
    return ctx
}

// gRPC 拦截器
func ContextPropagationInterceptor(registry *PropagatorRegistry) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        md, _ := metadata.FromOutgoingContext(ctx)
        carrier := MetadataCarrier(md)
        
        registry.InjectAll(ctx, carrier)
        
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}
```

---

## 消息队列传播

```go
type MessageCarrier map[string]string

func (mc MessageCarrier) Get(key string) string {
    return mc[key]
}

func (mc MessageCarrier) Set(key string, value string) {
    mc[key] = value
}

func (mc MessageCarrier) Keys() []string {
    keys := make([]string, 0, len(mc))
    for k := range mc {
        keys = append(keys, k)
    }
    return keys
}

// Kafka 消息上下文传播
func PublishWithContext(ctx context.Context, producer sarama.SyncProducer, topic string, msg *sarama.ProducerMessage) error {
    carrier := make(MessageCarrier)
    
    // 注入上下文
    propagator.Inject(ctx, carrier)
    
    // 设置消息 Headers
    for k, v := range carrier {
        msg.Headers = append(msg.Headers, sarama.RecordHeader{
            Key:   []byte(k),
            Value: []byte(v),
        })
    }
    
    _, _, err := producer.SendMessage(msg)
    return err
}

func ConsumeWithContext(msg *sarama.ConsumerMessage) context.Context {
    ctx := context.Background()
    carrier := make(MessageCarrier)
    
    // 从 Headers 提取
    for _, h := range msg.Headers {
        carrier[string(h.Key)] = string(h.Value)
    }
    
    // 提取上下文
    return propagator.Extract(ctx, carrier)
}
```

---

## 异步任务上下文传递

```go
// 任务上下文序列化
type SerializableContext struct {
    RequestID string
    TraceID   string
    SpanID    string
    Baggage   map[string]string
    Deadline  *time.Time
}

func SerializeContext(ctx context.Context) ([]byte, error) {
    sc := SerializableContext{
        RequestID: RequestIDFromContext(ctx),
        TraceID:   TraceIDFromContext(ctx),
        SpanID:    SpanIDFromContext(ctx),
        Baggage:   BaggageFromContext(ctx),
    }
    
    if deadline, ok := ctx.Deadline(); ok {
        sc.Deadline = &deadline
    }
    
    return json.Marshal(sc)
}

func DeserializeContext(data []byte) (context.Context, context.CancelFunc) {
    var sc SerializableContext
    json.Unmarshal(data, &sc)
    
    ctx := context.Background()
    var cancel context.CancelFunc
    
    // 恢复 deadline
    if sc.Deadline != nil && time.Now().Before(*sc.Deadline) {
        ctx, cancel = context.WithDeadline(ctx, *sc.Deadline)
    }
    
    // 恢复上下文值
    ctx = WithRequestID(ctx, sc.RequestID)
    ctx = WithTraceID(ctx, sc.TraceID)
    ctx = WithSpanID(ctx, sc.SpanID)
    ctx = WithBaggage(ctx, sc.Baggage)
    
    return ctx, cancel
}

// 在任务执行时恢复上下文
func (w *Worker) executeTask(task *Task) {
    ctx, cancel := DeserializeContext(task.ContextData)
    defer cancel()
    
    // 现在 ctx 包含了原始请求的所有上下文信息
    w.executor.Execute(ctx, task.Payload)
}
```

---

## 上下文清理与验证

```go
// 清理敏感信息
func SanitizeContext(ctx context.Context) context.Context {
    // 创建新的上下文，移除敏感信息
    newCtx := context.Background()
    
    // 保留安全的传播信息
    if reqID := RequestIDFromContext(ctx); reqID != "" {
        newCtx = WithRequestID(newCtx, reqID)
    }
    
    if traceID := TraceIDFromContext(ctx); traceID != "" {
        newCtx = WithTraceID(newCtx, traceID)
    }
    
    // 清理 baggage 中的敏感字段
    baggage := BaggageFromContext(ctx)
    sanitized := make(map[string]string)
    for k, v := range baggage {
        if !isSensitiveKey(k) {
            sanitized[k] = v
        }
    }
    newCtx = WithBaggage(newCtx, sanitized)
    
    return newCtx
}

func isSensitiveKey(key string) bool {
    sensitive := []string{"password", "token", "secret", "key"}
    lower := strings.ToLower(key)
    for _, s := range sensitive {
        if strings.Contains(lower, s) {
            return true
        }
    }
    return false
}
```
