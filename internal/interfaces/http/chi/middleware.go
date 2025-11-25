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
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		tr := otel.Tracer("http")
		ctx, span := tr.Start(ctx, r.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		// 将 span 注入到 context
		r = r.WithContext(ctx)

		// 注入追踪头到响应
		propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  middleware.DefaultLogFormatter{}.Logger,
		NoColor: false,
	})(next)
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

// RecovererMiddleware 恢复中间件
func RecovererMiddleware(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
