// Package middleware 提供 HTTP 中间件实现
//
// 采样中间件：基于框架的采样机制，对请求进行采样
package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yourusername/golang/pkg/sampling"
	"go.opentelemetry.io/otel/trace"
)

// SamplingConfig 采样中间件配置
//
// 功能说明：
// - 配置请求采样行为
// - 支持多种采样策略
// - 支持跳过特定路径的采样
//
// 字段说明：
// - Sampler: 采样器实例（必填）
// - SkipPaths: 跳过采样的路径列表（如 /health、/metrics）
// - AddSamplingDecision: 是否在响应头中添加采样决策信息
//
// 使用示例：
//
//	sampler, _ := sampling.NewProbabilisticSampler(0.5)
//	config := middleware.SamplingConfig{
//	    Sampler:            sampler,
//	    SkipPaths:          []string{"/health", "/metrics"},
//	    AddSamplingDecision: true,
//	}
//	router.Use(middleware.SamplingMiddleware(config))
type SamplingConfig struct {
	Sampler            sampling.Sampler
	SkipPaths          []string
	AddSamplingDecision bool
}

// SamplingMiddleware 创建请求采样中间件
//
// 功能说明：
// - 基于框架的采样机制对请求进行采样
// - 如果请求被采样，则继续处理
// - 如果请求未被采样，可以选择跳过某些处理（如详细日志、追踪等）
// - 支持在响应头中添加采样决策信息
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 使用采样器决定是否采样
// 3. 将采样决策添加到上下文
// 4. 可选在响应头中添加采样决策信息
// 5. 继续处理请求
//
// 采样决策的使用：
// - 可以在后续中间件中使用采样决策
// - 例如：只有被采样的请求才记录详细日志
// - 例如：只有被采样的请求才进行详细的性能追踪
//
// 使用示例：
//
//	sampler, _ := sampling.NewProbabilisticSampler(0.5)
//	router.Use(middleware.SamplingMiddleware(middleware.SamplingConfig{
//	    Sampler:            sampler,
//	    SkipPaths:          []string{"/health"},
//	    AddSamplingDecision: true,
//	}))
func SamplingMiddleware(config SamplingConfig) func(http.Handler) http.Handler {
	if config.Sampler == nil {
		// 如果没有配置采样器，使用总是采样
		config.Sampler = sampling.NewAlwaysSampler()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否跳过采样
			if shouldSkipSampling(r.URL.Path, config.SkipPaths) {
				// 跳过采样，总是处理请求
				next.ServeHTTP(w, r)
				return
			}

			// 构建采样元数据
			metadata := map[string]interface{}{
				"http.method": r.Method,
				"http.path":   r.URL.Path,
				"http.host":   r.Host,
			}

			// 从追踪上下文中获取 TraceID（如果有）
			span := trace.SpanFromContext(r.Context())
			if span.SpanContext().IsValid() {
				metadata["trace.id"] = span.SpanContext().TraceID().String()
				metadata["span.id"] = span.SpanContext().SpanID().String()
			}

			// 决定是否采样
			// 注意：框架的采样器接口只接受 context，元数据需要通过其他方式传递
			shouldSample := config.Sampler.ShouldSample(r.Context())

			// 将采样决策添加到上下文
			ctx := context.WithValue(r.Context(), "sampling.decision", shouldSample)
			// 获取采样器类型名称
			samplerType := fmt.Sprintf("%T", config.Sampler)
			ctx = context.WithValue(ctx, "sampling.sampler", samplerType)

			// 如果配置了，在响应头中添加采样决策信息
			if config.AddSamplingDecision {
				if shouldSample {
					w.Header().Set("X-Sampling-Decision", "sampled")
				} else {
					w.Header().Set("X-Sampling-Decision", "not-sampled")
				}
				samplerType := fmt.Sprintf("%T", config.Sampler)
				w.Header().Set("X-Sampling-Sampler", samplerType)
			}

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsSampled 从上下文中获取采样决策
//
// 功能说明：
// - 检查当前请求是否被采样
// - 用于后续中间件或处理器中决定是否执行某些操作
//
// 参数：
//   - ctx: 请求上下文
//
// 返回：
//   - bool: 如果请求被采样返回 true，否则返回 false
//
// 使用示例：
//
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//	    if middleware.IsSampled(r.Context()) {
//	        // 只有被采样的请求才记录详细日志
//	        log.Debug("Detailed request info", ...)
//	    }
//	}
func IsSampled(ctx context.Context) bool {
	if decision, ok := ctx.Value("sampling.decision").(bool); ok {
		return decision
	}
	// 如果没有采样决策，默认返回 true（总是处理）
	return true
}

// GetSamplerName 从上下文中获取采样器名称
//
// 功能说明：
// - 获取当前请求使用的采样器名称
// - 用于日志记录和调试
//
// 参数：
//   - ctx: 请求上下文
//
// 返回：
//   - string: 采样器名称，如果没有则返回空字符串
func GetSamplerName(ctx context.Context) string {
	if name, ok := ctx.Value("sampling.sampler").(string); ok {
		return name
	}
	return ""
}

// shouldSkipSampling 检查路径是否应该跳过采样
func shouldSkipSampling(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
