# 日志模式 (Logging Patterns)

> **分类**: 工程与云原生  
> **标签**: #logging #observability #structured

---

## 结构化日志

```go
import "go.uber.org/zap"

var logger *zap.Logger

func InitLogger() {
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{"stdout", "/var/log/app.log"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    var err error
    logger, err = config.Build()
    if err != nil {
        log.Fatal(err)
    }
}

// 使用
func ProcessRequest(ctx context.Context, req Request) {
    logger.Info("processing request",
        zap.String("request_id", GetRequestID(ctx)),
        zap.String("user_id", req.UserID),
        zap.String("method", req.Method),
        zap.Int("items_count", len(req.Items)),
        zap.Duration("latency", time.Since(start)),
    )
}
```

---

## 日志级别控制

```go
func LogWithLevel(level string, msg string, fields ...zap.Field) {
    switch level {
    case "debug":
        logger.Debug(msg, fields...)
    case "info":
        logger.Info(msg, fields...)
    case "warn":
        logger.Warn(msg, fields...)
    case "error":
        logger.Error(msg, fields...)
    }
}

// 动态调整级别
func SetLogLevel(level string) error {
    var l zap.AtomicLevel
    if err := l.UnmarshalText([]byte(level)); err != nil {
        return err
    }
    logger.Core().With([]zap.Field{}, l)
    return nil
}
```

---

## 上下文日志

```go
type contextKey string

const loggerKey contextKey = "logger"

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
    if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
        return logger
    }
    return zap.NewNop()
}

// 添加上下文字段
func WithRequestID(ctx context.Context, requestID string) context.Context {
    logger := LoggerFromContext(ctx)
    logger = logger.With(zap.String("request_id", requestID))
    return WithLogger(ctx, logger)
}
```

---

## 日志采样

```go
func SampledLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    
    // 每秒最多记录 100 条 info，保留所有 error
    config.Sampling = &zap.SamplingConfig{
        Initial:    100,
        Thereafter: 100,
    }
    
    logger, _ := config.Build()
    return logger
}
```

---

## 敏感信息过滤

```go
func SanitizeFields(data map[string]interface{}) map[string]interface{} {
    sensitive := []string{
        "password", "token", "secret", "api_key", 
        "credit_card", "ssn", "email",
    }
    
    sanitized := make(map[string]interface{})
    for k, v := range data {
        isSensitive := false
        for _, s := range sensitive {
            if strings.Contains(strings.ToLower(k), s) {
                isSensitive = true
                break
            }
        }
        
        if isSensitive {
            sanitized[k] = "[REDACTED]"
        } else {
            sanitized[k] = v
        }
    }
    
    return sanitized
}
```

---

## 日志轮转

```go
import "gopkg.in/natefinch/lumberjack.v2"

func NewRotatingLogger() *zap.Logger {
    w := zapcore.AddSync(&lumberjack.Logger{
        Filename:   "/var/log/app.log",
        MaxSize:    100,  // MB
        MaxBackups: 3,
        MaxAge:     7,    // days
        Compress:   true,
    })
    
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        w,
        zap.InfoLevel,
    )
    
    return zap.New(core)
}
```
