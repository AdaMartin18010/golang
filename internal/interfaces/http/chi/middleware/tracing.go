package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig 是请求追踪中间件的配置结构。
//
// 功能说明：
// - 配置 OpenTelemetry 追踪行为
// - 支持跳过特定路径的追踪
// - 支持添加请求 ID 和用户 ID 到追踪上下文
//
// 字段说明：
// - TracerName: Tracer 名称（默认：http-server）
// - ServiceName: 服务名称（默认：api）
// - ServiceVersion: 服务版本（可选）
// - SkipPaths: 跳过追踪的路径列表（如 /health、/metrics）
// - AddRequestID: 是否添加请求 ID 到追踪上下文
// - AddUserID: 是否添加用户 ID 到追踪上下文
//
// 使用示例：
//
//	config := middleware.TracingConfig{
//	    ServiceName:    "user-service",
//	    ServiceVersion: "1.0.0",
//	    SkipPaths:      []string{"/health", "/metrics"},
//	    AddRequestID:   true,
//	    AddUserID:      true,
//	}
//	router.Use(middleware.TracingMiddleware(config))
type TracingConfig struct {
	TracerName      string
	ServiceName     string
	ServiceVersion  string
	SkipPaths       []string
	AddRequestID    bool
	AddUserID       bool
}

// TracingMiddleware 创建请求追踪中间件。
//
// 功能说明：
// - 基于 OpenTelemetry 实现分布式追踪
// - 为每个 HTTP 请求创建 Span
// - 提取和传播追踪上下文
// - 记录请求和响应信息
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 从请求头提取追踪上下文（TraceContext、Baggage）
// 3. 创建新的 Span（Server Span）
// 4. 设置基本属性（HTTP 方法、路径、用户代理等）
// 5. 可选添加请求 ID 和用户 ID
// 6. 在响应头中添加追踪 ID
// 7. 执行下一个处理器
// 8. 记录响应信息（状态码、响应大小、耗时）
// 9. 设置 Span 状态（成功或错误）
//
// 追踪属性：
// - http.method: HTTP 方法
// - http.target: 请求路径
// - http.route: 路由路径
// - http.user_agent: 用户代理
// - http.client_ip: 客户端 IP
// - service.name: 服务名称
// - service.version: 服务版本（如果配置）
// - http.request_id: 请求 ID（如果配置）
// - user.id: 用户 ID（如果配置）
// - http.status_code: 响应状态码
// - http.response_size: 响应大小
// - http.duration_ms: 请求耗时（毫秒）
//
// 响应头：
// - X-Trace-ID: 追踪 ID
// - X-Span-ID: Span ID
//
// 参数：
// - config: 追踪配置
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	config := middleware.TracingConfig{
//	    ServiceName:  "api",
//	    SkipPaths:    []string{"/health"},
//	    AddRequestID: true,
//	}
//	router.Use(middleware.TracingMiddleware(config))
//
// 注意事项：
// - 需要先初始化 OpenTelemetry TracerProvider
// - 追踪上下文会自动传播到下游服务
// - 可以在 Handler 中使用 context 获取当前 Span
func TracingMiddleware(config TracingConfig) func(http.Handler) http.Handler {
	if config.TracerName == "" {
		config.TracerName = "http-server"
	}
	if config.ServiceName == "" {
		config.ServiceName = "api"
	}

	tracer := otel.Tracer(config.TracerName)
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否跳过追踪
			if shouldSkipTracing(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 从请求头提取追踪上下文
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// 创建span
			ctx, span := tracer.Start(ctx, r.URL.Path,
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// 设置基本属性
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.target", r.URL.Path),
				attribute.String("http.route", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.client_ip", getClientIP(r)),
				attribute.String("service.name", config.ServiceName),
			)

			if config.ServiceVersion != "" {
				span.SetAttributes(attribute.String("service.version", config.ServiceVersion))
			}

			// 添加请求ID
			if config.AddRequestID {
				if requestID := middleware.GetReqID(r.Context()); requestID != "" {
					span.SetAttributes(attribute.String("http.request_id", requestID))
					ctx = context.WithValue(ctx, "trace_id", requestID)
				}
			}

			// 添加用户ID
			if config.AddUserID {
				if userID := GetUserID(ctx); userID != "" {
					span.SetAttributes(attribute.String("user.id", userID))
				}
			}

			// 设置追踪ID到响应头
			spanContext := span.SpanContext()
			if spanContext.IsValid() {
				w.Header().Set("X-Trace-ID", spanContext.TraceID().String())
				w.Header().Set("X-Span-ID", spanContext.SpanID().String())
			}

			// 创建响应包装器
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			// 执行下一个处理器
			next.ServeHTTP(ww, r.WithContext(ctx))

			// 设置响应属性
			duration := time.Since(start)
			span.SetAttributes(
				attribute.Int("http.status_code", ww.Status()),
				attribute.Int("http.response_size", ww.BytesWritten()),
				attribute.Int64("http.duration_ms", duration.Milliseconds()),
			)

			// 设置状态
			if ww.Status() >= 400 {
				span.SetStatus(trace.StatusError, http.StatusText(ww.Status()))
			} else {
				span.SetStatus(trace.StatusOK, "")
			}
		})
	}
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	// 优先使用X-Forwarded-For
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	// 使用X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	// 使用RemoteAddr
	return r.RemoteAddr
}

// shouldSkipTracing 检查是否应该跳过追踪
func shouldSkipTracing(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || path == skipPath+"/" {
			return true
		}
	}
	return false
}
