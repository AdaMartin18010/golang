// Package chi 提供 HTTP 中间件实现
//
// 设计原理：
// 1. 中间件用于在请求处理前后执行通用逻辑
// 2. 中间件按顺序执行，形成中间件链
// 3. 每个中间件可以修改请求或响应
//
// 中间件类型：
// - 追踪中间件：OpenTelemetry 分布式追踪
// - 日志中间件：请求日志记录
// - 恢复中间件：Panic 恢复
// - 超时中间件：请求超时控制
// - CORS 中间件：跨域资源共享
package chi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware OpenTelemetry 追踪中间件
//
// 设计原理：
// 1. 为每个 HTTP 请求创建 OpenTelemetry Span
// 2. 提取和传播追踪上下文（TraceContext、Baggage）
// 3. 记录请求相关的属性（方法、路径、状态码等）
//
// 功能：
// - 创建服务器端 Span
// - 提取追踪上下文（从 HTTP Header）
// - 记录请求属性（HTTP 方法、路径、用户代理、客户端IP）
// - 记录响应状态码
//
// 追踪上下文传播：
// - TraceContext: W3C Trace Context 标准
// - Baggage: 跨服务传递的键值对
//
// 使用场景：
// - 分布式追踪
// - 性能分析
// - 请求链路追踪
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("http-server")
		propagator := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)

		// 从 HTTP Header 提取追踪上下文
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
		// 创建服务器端 Span
		ctx, span := tracer.Start(ctx, r.URL.Path, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// 记录请求属性
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.target", r.URL.Path),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.client_ip", r.RemoteAddr),
		)

		// 包装 ResponseWriter 以获取状态码
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		// 将上下文传递给下一个处理器
		next.ServeHTTP(ww, r.WithContext(ctx))

		// 记录响应状态码
		span.SetAttributes(attribute.Int("http.status_code", ww.Status()))
	})
}

// LoggingMiddleware 日志中间件
//
// 设计原理：
// 1. 记录每个 HTTP 请求的日志
// 2. 使用 Chi 的默认日志格式化器
// 3. 记录请求方法、路径、状态码、耗时等信息
//
// 功能：
// - 请求日志记录
// - 响应日志记录
// - 错误日志记录
//
// 日志格式：
// - 包含请求ID、方法、路径、状态码、耗时等
// - 支持彩色输出（开发环境）
func LoggingMiddleware(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  middleware.DefaultLogFormatter{}.Logger,
		NoColor: false,
	})(next)
}

// RecovererMiddleware 恢复中间件
//
// 设计原理：
// 1. 捕获请求处理过程中的 Panic
// 2. 记录 Panic 信息
// 3. 返回 500 错误响应，防止程序崩溃
//
// 功能：
// - Panic 恢复
// - 错误日志记录
// - 错误响应返回
//
// 使用场景：
// - 防止未处理的 Panic 导致程序崩溃
// - 记录错误信息用于调试
func RecovererMiddleware(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

// TimeoutMiddleware 超时中间件
//
// 设计原理：
// 1. 为每个请求设置超时时间
// 2. 如果请求处理超时，返回 504 Gateway Timeout
// 3. 防止长时间运行的请求占用资源
//
// 参数：
//   - timeout: 超时时间（如 60 秒）
//
// 返回：
//   - func(http.Handler) http.Handler: 中间件函数
//
// 功能：
// - 请求超时控制
// - 资源保护
// - 防止慢请求影响其他请求
//
// 使用场景：
// - 防止慢查询导致请求堆积
// - 保护服务器资源
// - 提高系统稳定性
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

// CORSMiddleware CORS 中间件
//
// 设计原理：
// 1. 处理跨域资源共享（CORS）请求
// 2. 设置 CORS 响应头
// 3. 处理 OPTIONS 预检请求
//
// 功能：
// - 允许跨域请求
// - 设置允许的 HTTP 方法
// - 设置允许的请求头
// - 处理 OPTIONS 预检请求
//
// 当前配置：
// - Access-Control-Allow-Origin: *（允许所有来源）
// - Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
// - Access-Control-Allow-Headers: Content-Type, Authorization
//
// 注意事项：
// - 生产环境应该限制允许的来源
// - 应该根据实际需求配置允许的方法和头
//
// 示例（生产环境配置）：
//   w.Header().Set("Access-Control-Allow-Origin", "https://example.com")
//   w.Header().Set("Access-Control-Allow-Credentials", "true")
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理 OPTIONS 预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
