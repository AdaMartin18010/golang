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

// TracingConfig 追踪配置
type TracingConfig struct {
	TracerName      string
	ServiceName     string
	ServiceVersion  string
	SkipPaths       []string
	AddRequestID    bool
	AddUserID       bool
}

// TracingMiddleware 请求追踪中间件
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
