# 上下文感知日志 (Context-Aware Logging)

> **分类**: 工程与云原生
> **标签**: #logging #context #observability

---

## 上下文注入

```go
// 从上下文提取日志字段
func ExtractLogFields(ctx context.Context) []zap.Field {
    var fields []zap.Field

    if reqID := RequestIDFromContext(ctx); reqID != "" {
        fields = append(fields, zap.String("request_id", reqID))
    }

    if traceID := TraceIDFromContext(ctx); traceID != "" {
        fields = append(fields, zap.String("trace_id", traceID))
    }

    if spanID := SpanIDFromContext(ctx); spanID != "" {
        fields = append(fields, zap.String("span_id", spanID))
    }

    if userID := UserIDFromContext(ctx); userID != "" {
        fields = append(fields, zap.String("user_id", userID))
    }

    if taskID := TaskIDFromContext(ctx); taskID != "" {
        fields = append(fields, zap.String("task_id", taskID))
    }

    return fields
}

// 创建上下文感知的 logger
func LoggerFromContext(ctx context.Context) *zap.Logger {
    baseLogger := zap.L()
    fields := ExtractLogFields(ctx)

    if len(fields) > 0 {
        return baseLogger.With(fields...)
    }

    return baseLogger
}

// HTTP 中间件注入
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 生成 request ID
        requestID := generateRequestID()

        // 注入上下文
        ctx := WithRequestID(c.Request.Context(), requestID)
        ctx = WithTraceID(ctx, c.GetHeader("X-Trace-ID"))

        // 替换请求上下文
        c.Request = c.Request.WithContext(ctx)

        // 创建专用 logger
        ctxLogger := LoggerFromContext(ctx)
        c.Set("logger", ctxLogger)

        // 记录请求
        ctxLogger.Info("request started",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.String("client_ip", c.ClientIP()),
        )

        c.Next()

        // 记录响应
        ctxLogger.Info("request completed",
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", time.Since(start)),
            zap.Int("bytes", c.Writer.Size()),
        )
    }
}
```

---

## 结构化上下文

```go
// 创建带上下文的 logger
func WithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
    // 提取所有上下文信息
    data := map[string]interface{}{
        "timestamp": time.Now().Format(time.RFC3339),
    }

    if reqID := RequestIDFromContext(ctx); reqID != "" {
        data["request_id"] = reqID
    }

    if traceID := TraceIDFromContext(ctx); traceID != "" {
        data["trace_id"] = traceID
        data["trace_url"] = fmt.Sprintf("https://tracing.example.com/trace/%s", traceID)
    }

    // 转换为 zap fields
    fields := make([]zap.Field, 0, len(data))
    for k, v := range data {
        fields = append(fields, zap.Any(k, v))
    }

    return logger.With(fields...)
}

// 使用示例
func ProcessOrder(ctx context.Context, order Order) error {
    logger := WithContext(ctx, zap.L())

    logger.Info("processing order",
        zap.String("order_id", order.ID),
        zap.Float64("amount", order.Amount),
    )

    // 所有日志都会自动包含 request_id 和 trace_id
    if err := validateOrder(order); err != nil {
        logger.Error("validation failed", zap.Error(err))
        return err
    }

    logger.Info("order processed successfully")
    return nil
}
```

---

## 异步任务日志

```go
// 异步任务上下文传递
func (w *Worker) ExecuteTask(ctx context.Context, task *Task) {
    // 序列化上下文
    ctxData := SerializeContext(ctx)

    // 存储任务
    task.ContextData = ctxData

    // 入队
    w.queue.Push(task)
}

func (w *Worker) processQueue() {
    for task := range w.queue.Pop() {
        // 反序列化上下文
        ctx, cancel := DeserializeContext(task.ContextData)

        // 创建任务 logger
        logger := WithContext(ctx, zap.L()).With(
            zap.String("task_id", task.ID),
            zap.String("task_type", task.Type),
        )

        logger.Info("task started")

        // 执行
        if err := w.execute(ctx, task); err != nil {
            logger.Error("task failed", zap.Error(err))
        } else {
            logger.Info("task completed")
        }

        cancel()
    }
}
```

---

## 日志追踪

```go
// 日志追踪器
type LogTracer struct {
    logger *zap.Logger
    spanID string
}

func NewLogTracer(ctx context.Context, operation string) *LogTracer {
    spanID := generateSpanID()

    tracer := &LogTracer{
        logger: WithContext(ctx, zap.L()).With(
            zap.String("operation", operation),
            zap.String("span_id", spanID),
        ),
        spanID: spanID,
    }

    tracer.logger.Info("operation started")
    return tracer
}

func (lt *LogTracer) Log(msg string, fields ...zap.Field) {
    lt.logger.Info(msg, fields...)
}

func (lt *LogTracer) Error(err error, fields ...zap.Field) {
    lt.logger.Error("operation error", append(fields, zap.Error(err))...)
}

func (lt *LogTracer) Finish(fields ...zap.Field) {
    lt.logger.Info("operation finished", fields...)
}

// 使用
func ComplexOperation(ctx context.Context) error {
    tracer := NewLogTracer(ctx, "complex_operation")
    defer tracer.Finish()

    tracer.Log("step 1 started")
    if err := step1(); err != nil {
        tracer.Error(err)
        return err
    }

    tracer.Log("step 2 started")
    if err := step2(); err != nil {
        tracer.Error(err)
        return err
    }

    return nil
}
```
