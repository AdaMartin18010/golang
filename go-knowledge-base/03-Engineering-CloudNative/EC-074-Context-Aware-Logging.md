# 上下文感知日志系统 (Context-Aware Logging)

> **分类**: 工程与云原生
> **标签**: #logging #context #structured-logging #opentelemetry
> **参考**: Zap, Logrus, OpenTelemetry Log Bridge

---

## 架构设计

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Context-Aware Logging System                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Log Entry Structure                             │   │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐          │   │
│  │  │  Timestamp  │    Level    │   Message   │   Fields    │          │   │
│  │  │  (RFC3339)  │ (INFO/ERR)  │   (string)  │  (kv pairs) │          │   │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘          │   │
│  │                         │                                          │   │
│  │                         ▼                                          │   │
│  │  ┌─────────────────────────────────────────────────────────────┐  │   │
│  │  │                  Context-Extracted Fields                      │  │   │
│  │  │  trace_id: 550e8400-e29b-41d4-a716-446655440000               │  │   │
│  │  │  span_id:  0af7651916cd43dd                                 │  │   │
│  │  │  request_id: req-12345                                       │  │   │
│  │  │  user_id:   user-67890                                       │  │   │
│  │  │  tenant_id: tenant-abc                                       │  │   │
│  │  │  service:   order-service                                    │  │   │
│  │  └─────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       Output Adapters                                │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Stdout    │  │    File     │  │   Syslog    │  │   Remote    │ │   │
│  │  │  (JSON)     │  │  (Rotate)   │  │             │  │   (OTLP)    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package logctx

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "runtime"
    "strings"
    "sync"
    "time"
)

// Level 日志级别
type Level int

const (
    DebugLevel Level = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    FatalLevel
    PanicLevel
)

func (l Level) String() string {
    switch l {
    case DebugLevel:
        return "DEBUG"
    case InfoLevel:
        return "INFO"
    case WarnLevel:
        return "WARN"
    case ErrorLevel:
        return "ERROR"
    case FatalLevel:
        return "FATAL"
    case PanicLevel:
        return "PANIC"
    default:
        return "UNKNOWN"
    }
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// ContextHook 从上下文提取字段的钩子
type ContextHook interface {
    Extract(ctx context.Context) []Field
}

// Logger 上下文感知日志器
type Logger struct {
    level       Level
    output      io.Writer
    encoder     Encoder
    contextKeys []contextKey
    hooks       []ContextHook
    mu          sync.RWMutex

    // 默认字段
    defaultFields []Field

    // 调用者信息
    callerSkip    int
    enableCaller  bool
}

// Encoder 编码器接口
type Encoder interface {
    Encode(entry *Entry) ([]byte, error)
}

// Entry 日志条目
type Entry struct {
    Time       time.Time
    Level      Level
    Message    string
    Fields     []Field
    Caller     string
    StackTrace string
}

// JSONEncoder JSON 编码器
type JSONEncoder struct{}

func (e *JSONEncoder) Encode(entry *Entry) ([]byte, error) {
    m := make(map[string]interface{})

    m["time"] = entry.Time.Format(time.RFC3339Nano)
    m["level"] = entry.Level.String()
    m["msg"] = entry.Message

    if entry.Caller != "" {
        m["caller"] = entry.Caller
    }

    for _, f := range entry.Fields {
        m[f.Key] = f.Value
    }

    return json.Marshal(m)
}

// ConsoleEncoder 控制台编码器
type ConsoleEncoder struct{}

func (e *ConsoleEncoder) Encode(entry *Entry) ([]byte, error) {
    var sb strings.Builder

    // 时间
    sb.WriteString(entry.Time.Format("2006-01-02 15:04:05.000"))
    sb.WriteString(" ")

    // 级别
    sb.WriteString(fmt.Sprintf("%-5s", entry.Level.String()))
    sb.WriteString(" ")

    // 调用者
    if entry.Caller != "" {
        sb.WriteString(entry.Caller)
        sb.WriteString(" ")
    }

    // 消息
    sb.WriteString(entry.Message)

    // 字段
    for _, f := range entry.Fields {
        sb.WriteString(fmt.Sprintf(" %s=%v", f.Key, f.Value))
    }

    sb.WriteString("\n")

    return []byte(sb.String()), nil
}

// New 创建日志器
func New(level Level, output io.Writer) *Logger {
    return &Logger{
        level:        level,
        output:       output,
        encoder:      &JSONEncoder{},
        callerSkip:   2,
        enableCaller: true,
    }
}

// Default 默认日志器
var Default = New(InfoLevel, os.Stdout)

// WithContext 创建带上下文的日志条目
func (l *Logger) WithContext(ctx context.Context) *ContextLogger {
    fields := l.extractFieldsFromContext(ctx)
    return &ContextLogger{
        logger: l,
        ctx:    ctx,
        fields: fields,
    }
}

// AddContextKey 添加上下文键
func (l *Logger) AddContextKey(key interface{}, fieldName string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.contextKeys = append(l.contextKeys, contextKey{key: key, fieldName: fieldName})
}

// AddHook 添加上下文钩子
func (l *Logger) AddHook(hook ContextHook) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.hooks = append(l.hooks, hook)
}

func (l *Logger) extractFieldsFromContext(ctx context.Context) []Field {
    var fields []Field

    l.mu.RLock()
    keys := l.contextKeys
    hooks := l.hooks
    l.mu.RUnlock()

    // 从注册的键提取
    for _, ck := range keys {
        if val := ctx.Value(ck.key); val != nil {
            fields = append(fields, Field{Key: ck.fieldName, Value: val})
        }
    }

    // 从钩子提取
    for _, hook := range hooks {
        fields = append(fields, hook.Extract(ctx)...)
    }

    return fields
}

type contextKey struct {
    key       interface{}
    fieldName string
}

// log 内部日志方法
func (l *Logger) log(ctx context.Context, level Level, msg string, fields ...Field) {
    if level < l.level {
        return
    }

    // 合并字段
    allFields := make([]Field, 0, len(fields)+len(l.defaultFields))
    allFields = append(allFields, l.defaultFields...)

    if ctx != nil {
        ctxFields := l.extractFieldsFromContext(ctx)
        allFields = append(allFields, ctxFields...)
    }

    allFields = append(allFields, fields...)

    entry := &Entry{
        Time:    time.Now(),
        Level:   level,
        Message: msg,
        Fields:  allFields,
    }

    // 添加调用者信息
    if l.enableCaller {
        entry.Caller = getCaller(l.callerSkip)
    }

    // 编码并输出
    data, err := l.encoder.Encode(entry)
    if err != nil {
        return
    }

    l.mu.Lock()
    l.output.Write(data)
    l.output.Write([]byte("\n"))
    l.mu.Unlock()
}

func getCaller(skip int) string {
    _, file, line, ok := runtime.Caller(skip + 1)
    if !ok {
        return ""
    }

    // 简化路径
    idx := strings.LastIndex(file, "/")
    if idx != -1 {
        file = file[idx+1:]
    }

    return fmt.Sprintf("%s:%d", file, line)
}

// ContextLogger 带上下文的日志器
type ContextLogger struct {
    logger *Logger
    ctx    context.Context
    fields []Field
}

func (cl *ContextLogger) Debug(msg string, fields ...Field) {
    cl.logger.log(cl.ctx, DebugLevel, msg, append(cl.fields, fields...)...)
}

func (cl *ContextLogger) Info(msg string, fields ...Field) {
    cl.logger.log(cl.ctx, InfoLevel, msg, append(cl.fields, fields...)...)
}

func (cl *ContextLogger) Warn(msg string, fields ...Field) {
    cl.logger.log(cl.ctx, WarnLevel, msg, append(cl.fields, fields...)...)
}

func (cl *ContextLogger) Error(msg string, fields ...Field) {
    cl.logger.log(cl.ctx, ErrorLevel, msg, append(cl.fields, fields...)...)
}

func (cl *ContextLogger) Fatal(msg string, fields ...Field) {
    cl.logger.log(cl.ctx, FatalLevel, msg, append(cl.fields, fields...)...)
    os.Exit(1)
}

func (cl *ContextLogger) WithFields(fields ...Field) *ContextLogger {
    newFields := make([]Field, len(cl.fields)+len(fields))
    copy(newFields, cl.fields)
    copy(newFields[len(cl.fields):], fields)

    return &ContextLogger{
        logger: cl.logger,
        ctx:    cl.ctx,
        fields: newFields,
    }
}
```

---

## 上下文钩子实现

```go
package logctx

import (
    "context"

    "go.opentelemetry.io/otel/trace"
)

// TraceContextHook OpenTelemetry 追踪上下文钩子
type TraceContextHook struct{}

func (h *TraceContextHook) Extract(ctx context.Context) []Field {
    var fields []Field

    // 提取 Span 上下文
    span := trace.SpanFromContext(ctx)
    if span != nil {
        spanContext := span.SpanContext()
        if spanContext.IsValid() {
            fields = append(fields,
                Field{Key: "trace_id", Value: spanContext.TraceID().String()},
                Field{Key: "span_id", Value: spanContext.SpanID().String()},
                Field{Key: "trace_flags", Value: spanContext.TraceFlags().String()},
            )
        }
    }

    return fields
}

// RequestContextHook HTTP 请求上下文钩子
type RequestContextHook struct {
    RequestIDKey interface{}
    UserIDKey    interface{}
    TenantIDKey  interface{}
    SessionIDKey interface{}
}

func (h *RequestContextHook) Extract(ctx context.Context) []Field {
    var fields []Field

    if h.RequestIDKey != nil {
        if id, ok := ctx.Value(h.RequestIDKey).(string); ok && id != "" {
            fields = append(fields, Field{Key: "request_id", Value: id})
        }
    }

    if h.UserIDKey != nil {
        if id, ok := ctx.Value(h.UserIDKey).(string); ok && id != "" {
            fields = append(fields, Field{Key: "user_id", Value: id})
        }
    }

    if h.TenantIDKey != nil {
        if id, ok := ctx.Value(h.TenantIDKey).(string); ok && id != "" {
            fields = append(fields, Field{Key: "tenant_id", Value: id})
        }
    }

    if h.SessionIDKey != nil {
        if id, ok := ctx.Value(h.SessionIDKey).(string); ok && id != "" {
            fields = append(fields, Field{Key: "session_id", Value: id})
        }
    }

    return fields
}

// BaggageContextHook Baggage 传播钩子
type BaggageContextHook struct {
    Keys []string // 需要提取的 baggage key
}

func (h *BaggageContextHook) Extract(ctx context.Context) []Field {
    var fields []Field

    // 假设 baggage 存储在上下文中
    if baggage, ok := ctx.Value("baggage").(map[string]string); ok {
        for _, key := range h.Keys {
            if val, ok := baggage[key]; ok {
                fields = append(fields, Field{Key: key, Value: val})
            }
        }
    }

    return fields
}

// ServiceContextHook 服务信息钩子
type ServiceContextHook struct {
    ServiceName    string
    ServiceVersion string
    HostName       string
    Environment    string
}

func (h *ServiceContextHook) Extract(ctx context.Context) []Field {
    return []Field{
        {Key: "service_name", Value: h.ServiceName},
        {Key: "service_version", Value: h.ServiceVersion},
        {Key: "hostname", Value: h.HostName},
        {Key: "environment", Value: h.Environment},
    }
}
```

---

## 中间件集成

```go
package logctx

import (
    "net/http"
    "time"
)

// LoggingMiddleware HTTP 日志中间件
func LoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // 生成请求 ID
            requestID := generateRequestID()
            ctx := WithRequestID(r.Context(), requestID)
            r = r.WithContext(ctx)

            // 包装 ResponseWriter
            rw := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

            // 执行请求
            next.ServeHTTP(rw, r)

            // 记录访问日志
            duration := time.Since(start)

            logger.WithContext(ctx).Info("HTTP request",
                Field{Key: "method", Value: r.Method},
                Field{Key: "path", Value: r.URL.Path},
                Field{Key: "status", Value: rw.statusCode},
                Field{Key: "duration_ms", Value: duration.Milliseconds()},
                Field{Key: "bytes_written", Value: rw.written},
                Field{Key: "user_agent", Value: r.UserAgent()},
                Field{Key: "remote_addr", Value: r.RemoteAddr},
            )
        })
    }
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

func generateRequestID() string {
    // 生成唯一请求 ID
    return uuid.New().String()
}

// requestIDKey context key
type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey{}).(string); ok {
        return id
    }
    return ""
}

// gRPC 拦截器
func UnaryServerInterceptor(logger *Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        start := time.Now()

        // 从 metadata 提取追踪信息
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if vals := md.Get("x-request-id"); len(vals) > 0 {
                ctx = WithRequestID(ctx, vals[0])
            }
        }

        resp, err := handler(ctx, req)

        duration := time.Since(start)

        // 记录日志
        ctxLogger := logger.WithContext(ctx)
        if err != nil {
            ctxLogger.Error("gRPC request failed",
                Field{Key: "method", Value: info.FullMethod},
                Field{Key: "duration_ms", Value: duration.Milliseconds()},
                Field{Key: "error", Value: err.Error()},
            )
        } else {
            ctxLogger.Info("gRPC request",
                Field{Key: "method", Value: info.FullMethod},
                Field{Key: "duration_ms", Value: duration.Milliseconds()},
            )
        }

        return resp, err
    }
}
```

---

## 高级功能

```go
package logctx

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// AsyncWriter 异步写入器
type AsyncWriter struct {
    writer io.Writer
    buffer chan []byte
    wg     sync.WaitGroup
    stop   chan struct{}
}

func NewAsyncWriter(writer io.Writer, bufferSize int) *AsyncWriter {
    w := &AsyncWriter{
        writer: writer,
        buffer: make(chan []byte, bufferSize),
        stop:   make(chan struct{}),
    }

    w.wg.Add(1)
    go w.process()

    return w
}

func (w *AsyncWriter) Write(p []byte) (n int, err error) {
    select {
    case w.buffer <- p:
        return len(p), nil
    default:
        // 缓冲区满，丢弃或阻塞
        return 0, fmt.Errorf("buffer full")
    }
}

func (w *AsyncWriter) process() {
    defer w.wg.Done()

    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    var batch [][]byte

    for {
        select {
        case <-w.stop:
            // 刷新剩余数据
            w.flush(batch)
            return
        case data := <-w.buffer:
            batch = append(batch, data)
            if len(batch) >= 100 {
                w.flush(batch)
                batch = nil
            }
        case <-ticker.C:
            if len(batch) > 0 {
                w.flush(batch)
                batch = nil
            }
        }
    }
}

func (w *AsyncWriter) flush(batch [][]byte) {
    for _, data := range batch {
        w.writer.Write(data)
    }
}

func (w *AsyncWriter) Close() error {
    close(w.stop)
    w.wg.Wait()
    return nil
}

// RotatingFileWriter 轮转文件写入器
type RotatingFileWriter struct {
    filename   string
    maxSize    int64
    maxBackups int
    maxAge     int

    file       *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingFileWriter(filename string, maxSize int64, maxBackups, maxAge int) (*RotatingFileWriter, error) {
    w := &RotatingFileWriter{
        filename:   filename,
        maxSize:    maxSize,
        maxBackups: maxBackups,
        maxAge:     maxAge,
    }

    if err := w.open(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *RotatingFileWriter) Write(p []byte) (n int, err error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.size+int64(len(p)) > w.maxSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = w.file.Write(p)
    w.size += int64(n)

    return n, err
}

func (w *RotatingFileWriter) open() error {
    file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.file = file
    w.size = info.Size()

    return nil
}

func (w *RotatingFileWriter) rotate() error {
    // 关闭当前文件
    if w.file != nil {
        w.file.Close()
    }

    // 重命名文件
    backup := fmt.Sprintf("%s.%s", w.filename, time.Now().Format("20060102-150405"))
    os.Rename(w.filename, backup)

    // 清理旧备份
    w.cleanup()

    // 打开新文件
    return w.open()
}

func (w *RotatingFileWriter) cleanup() {
    // 实现备份清理逻辑
    dir := filepath.Dir(w.filename)
    pattern := filepath.Base(w.filename) + ".*"

    matches, _ := filepath.Glob(filepath.Join(dir, pattern))
    if len(matches) > w.maxBackups {
        // 删除最旧的备份
        // ...
    }
}

func (w *RotatingFileWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.file != nil {
        return w.file.Close()
    }
    return nil
}

// SamplingLogger 采样日志器
type SamplingLogger struct {
    logger   *Logger
    interval time.Duration
    lastLog  map[string]time.Time
    mu       sync.RWMutex
}

func NewSamplingLogger(logger *Logger, interval time.Duration) *SamplingLogger {
    return &SamplingLogger{
        logger:   logger,
        interval: interval,
        lastLog:  make(map[string]time.Time),
    }
}

func (s *SamplingLogger) ShouldLog(key string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    if last, ok := s.lastLog[key]; !ok || now.Sub(last) > s.interval {
        s.lastLog[key] = now
        return true
    }

    return false
}

func (s *SamplingLogger) Info(key, msg string, fields ...Field) {
    if s.ShouldLog(key) {
        s.logger.log(nil, InfoLevel, msg, fields...)
    }
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "os"

    "logctx"
)

func main() {
    // 创建日志器
    logger := logctx.New(logctx.InfoLevel, os.Stdout)

    // 添加上下文钩子
    logger.AddHook(&logctx.TraceContextHook{})
    logger.AddHook(&logctx.RequestContextHook{
        RequestIDKey: "request-id",
        UserIDKey:    "user-id",
    })

    // 注册上下文键
    logger.AddContextKey("tenant-id", "tenant_id")

    // 创建带上下文的日志器
    ctx := context.WithValue(context.Background(), "request-id", "req-123")
    ctx = context.WithValue(ctx, "user-id", "user-456")
    ctx = context.WithValue(ctx, "tenant-id", "tenant-abc")

    ctxLogger := logger.WithContext(ctx)

    // 输出日志
    ctxLogger.Info("Processing order",
        logctx.Field{Key: "order_id", Value: "ORD-789"},
        logctx.Field{Key: "amount", Value: 100.50},
    )

    // 输出：
    // {"time":"2024-01-15T10:30:00.123Z","level":"INFO","msg":"Processing order",
    //  "request_id":"req-123","user_id":"user-456","tenant_id":"tenant-abc",
    //  "order_id":"ORD-789","amount":100.5}
}
```
