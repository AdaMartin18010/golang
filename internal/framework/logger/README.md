# Framework Logger - 框架统一日志系统

## 📋 概述

框架级别的统一日志系统，为整个框架提供一致的日志记录能力。

## 🎯 设计原则

1. **统一接口**：框架内所有组件使用相同的日志接口
2. **结构化日志**：基于 slog，支持 key-value 结构化日志
3. **OpenTelemetry 集成**：自动添加追踪信息（TraceID、SpanID）
4. **可配置**：支持日志级别、采样率、输出格式等配置
5. **线程安全**：支持并发写入

## 🚀 快速开始

### 基本使用

```go
import "github.com/yourusername/golang/internal/framework/logger"

// 获取默认日志实例
log := logger.GetLogger()

// 记录日志
log.Info("Application started", "port", 8080)
log.Error("Failed to connect", "error", err)
```

### 自定义配置

```go
import (
    "log/slog"
    "github.com/yourusername/golang/internal/framework/logger"
)

// 创建自定义日志实例
log := logger.NewLogger(&logger.Config{
    Level:      slog.LevelDebug,
    JSONFormat: true,
    ServiceName: "my-service",
    ServiceVersion: "1.0.0",
})

log.Info("Custom logger initialized")
```

### 带上下文的日志

```go
ctx := context.WithValue(context.Background(), "request_id", "req-123")
log.WithContext(ctx).Info("Processing request", "user_id", 123)
```

## 📚 API 参考

### Logger 接口

```go
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    WithContext(ctx context.Context) *slog.Logger
    WithFields(fields ...any) *slog.Logger
    WithError(err error) *slog.Logger
}
```

### Config 配置

```go
type Config struct {
    Level          slog.Level  // 日志级别
    Output         io.Writer   // 输出目标
    AddSource      bool        // 是否添加源代码位置
    JSONFormat     bool        // 是否使用 JSON 格式
    SampleRate     float64     // 采样率 (0.0-1.0)
    ServiceName    string      // 服务名称
    ServiceVersion string      // 服务版本
}
```

## 🔧 配置说明

### 日志级别

- `slog.LevelDebug`: 调试信息
- `slog.LevelInfo`: 一般信息（默认）
- `slog.LevelWarn`: 警告信息
- `slog.LevelError`: 错误信息

### 输出格式

- **JSON 格式**（默认）：适合生产环境，便于日志聚合和分析
- **文本格式**：适合开发环境，人类可读

### 采样率

- `1.0`: 记录所有日志（默认）
- `0.5`: 记录 50% 的日志
- `0.0`: 不记录日志（错误日志除外）

## 📝 使用示例

### 在应用中使用

```go
package main

import (
    "github.com/yourusername/golang/internal/framework/logger"
)

func main() {
    log := logger.GetLogger()
    log.Info("Application starting", "version", "1.0.0")

    // ... 应用逻辑
}
```

### 在中间件中使用

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log := logger.GetLogger()
        ctx := r.Context()

        log.WithContext(ctx).Info("Request received",
            "method", r.Method,
            "path", r.URL.Path,
        )

        next.ServeHTTP(w, r)
    })
}
```

## 🔍 最佳实践

1. **使用结构化日志**：使用 key-value 格式，而不是字符串拼接
2. **添加上下文**：使用 `WithContext` 添加追踪信息
3. **合理使用日志级别**：Debug 用于调试，Info 用于一般信息，Error 用于错误
4. **避免敏感信息**：不要在日志中记录密码、token 等敏感信息
5. **生产环境使用 JSON 格式**：便于日志系统解析和分析

## 🔗 相关文档

- pkg/logger - 底层日志实现
- OpenTelemetry 集成 - 追踪集成
