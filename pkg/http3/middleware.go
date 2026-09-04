package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// =============================================================================
// 中间件系统 - Middleware System
// =============================================================================

// Middleware 中间件类型
type Middleware func(http.Handler) http.Handler

// MiddlewareChain 中间件链
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain 创建中间件链
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

// Use 添加中间件
func (mc *MiddlewareChain) Use(m Middleware) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, m)
	return mc
}

// Then 应用中间件链到处理器
func (mc *MiddlewareChain) Then(h http.Handler) http.Handler {
	// 从后向前应用中间件
	for _, v := range slices.Backward(mc.middlewares) {
		h = v(h)
	}
	return h
}

// =============================================================================
// 内置中间件
// =============================================================================

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建响应记录器
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			rec.statusCode,
			duration,
		)
	})
}

// responseRecorder 响应记录器
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// RecoveryMiddleware 恢复中间件（捕获panic）
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware CORS中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			tw := &timeoutWriter{w: w, hdr: make(http.Header)}
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(tw, r)
				close(done)
			}()

			select {
			case <-done:
				// 请求完成
			case <-ctx.Done():
				// 超时，禁止处理器继续写入响应，并直接写入408
				tw.mu.Lock()
				defer tw.mu.Unlock()
				tw.timedOut = true
				tw.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				tw.w.Header().Set("X-Content-Type-Options", "nosniff")
				tw.w.WriteHeader(http.StatusRequestTimeout)
				fmt.Fprintln(tw.w, "Request Timeout")
			}
		})
	}
}

// timeoutWriter 超时保护写入器
// 超时后与处理器goroutine并发写ResponseWriter会产生DATA RACE，
// 这里通过互斥锁+超时标志保证同一时刻只有一个写入方，且超时后丢弃处理器的写入
type timeoutWriter struct {
	mu          sync.Mutex
	w           http.ResponseWriter
	hdr         http.Header
	timedOut    bool
	wroteHeader bool
}

func (tw *timeoutWriter) Header() http.Header {
	if tw.hdr == nil {
		tw.hdr = make(http.Header)
	}
	return tw.hdr
}

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.wroteHeader = true
	tw.syncHeaderLocked()
	tw.w.WriteHeader(statusCode)
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if !tw.wroteHeader {
		tw.wroteHeader = true
		tw.syncHeaderLocked()
		tw.w.WriteHeader(http.StatusOK)
	}
	return tw.w.Write(b)
}

// syncHeaderLocked 将快照头同步到底层writer（调用方必须持有锁）
func (tw *timeoutWriter) syncHeaderLocked() {
	for k, vv := range tw.hdr {
		tw.w.Header()[k] = vv
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware(next http.Handler) http.Handler {
	var requestID atomic.Uint64

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := requestID.Add(1)
		ctx := context.WithValue(r.Context(), "request_id", id)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", fmt.Sprintf("%d", id))

		next.ServeHTTP(w, r)
	})
}

// CompressionMiddleware 压缩中间件（简化版）
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查客户端是否支持gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 实际应用中应该实现gzip压缩
		// 这里仅作为示例
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(requestsPerSecond int) Middleware {
	// 使用令牌桶算法
	limiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-limiter.C:
				next.ServeHTTP(w, r)
			default:
				http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
			}
		})
	}
}

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 安全相关的HTTP头
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

// CacheMiddleware 缓存中间件
func CacheMiddleware(maxAge time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只缓存GET请求
			if r.Method == "GET" {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware 认证中间件（简化版）
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 实际应用中应该验证token
		// 这里仅作为示例
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
