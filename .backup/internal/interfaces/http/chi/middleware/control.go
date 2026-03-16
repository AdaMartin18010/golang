// Package middleware 提供 HTTP 中间件实现
//
// 精细控制中间件：基于框架的精细控制机制，提供功能开关、速率控制、熔断器等功能
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourusername/golang/pkg/control"
)

// ControlConfig 精细控制中间件配置
//
// 功能说明：
// - 配置功能开关、速率控制、熔断器等控制机制
// - 支持按路径配置不同的控制策略
//
// 字段说明：
// - FeatureController: 功能控制器（可选）
// - RateController: 速率控制器（可选）
// - CircuitController: 熔断器控制器（可选）
// - FeatureFlags: 路径到功能开关的映射
// - RateLimits: 路径到速率限制的映射
// - CircuitBreakers: 路径到熔断器的映射
// - SkipPaths: 跳过控制的路径列表
//
// 使用示例：
//
//	featureController := control.NewFeatureController()
//	rateController := control.NewRateController()
//	circuitController := control.NewCircuitController()
//
//	config := middleware.ControlConfig{
//	    FeatureController: featureController,
//	    RateController:    rateController,
//	    CircuitController: circuitController,
//	    FeatureFlags: map[string]string{
//	        "/api/v1/experimental": "experimental-feature",
//	    },
//	    RateLimits: map[string]string{
//	        "/api/v1/users": "user-api",
//	    },
//	    CircuitBreakers: map[string]string{
//	        "/api/v1/external": "external-api",
//	    },
//	}
//	router.Use(middleware.ControlMiddleware(config))
type ControlConfig struct {
	FeatureController *control.FeatureController
	RateController    *control.RateController
	CircuitController *control.CircuitController
	FeatureFlags      map[string]string // 路径 -> 功能开关名称
	RateLimits        map[string]string // 路径 -> 速率限制器名称
	CircuitBreakers   map[string]string // 路径 -> 熔断器名称
	SkipPaths         []string
}

// ControlMiddleware 创建精细控制中间件
//
// 功能说明：
// - 基于框架的精细控制机制对请求进行控制
// - 支持功能开关、速率控制、熔断器
// - 按路径配置不同的控制策略
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 检查功能开关（如果配置）
// 3. 检查速率限制（如果配置）
// 4. 检查熔断器（如果配置）
// 5. 如果任何检查失败，返回相应的错误响应
// 6. 继续处理请求
//
// 使用示例：
//
//	router.Use(middleware.ControlMiddleware(middleware.ControlConfig{
//	    FeatureController: featureController,
//	    RateController:    rateController,
//	    CircuitController: circuitController,
//	}))
func ControlMiddleware(config ControlConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// 检查是否跳过控制
			if shouldSkipControl(path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 将控制器添加到上下文，供后续使用
			ctx := r.Context()
			if config.FeatureController != nil {
				ctx = context.WithValue(ctx, "feature.controller", config.FeatureController)
			}
			if config.RateController != nil {
				ctx = context.WithValue(ctx, "rate.controller", config.RateController)
			}
			if config.CircuitController != nil {
				ctx = context.WithValue(ctx, "circuit.controller", config.CircuitController)
			}

			// 检查功能开关
			if config.FeatureController != nil {
				if flagName, ok := config.FeatureFlags[path]; ok {
					if !config.FeatureController.IsEnabled(flagName) {
						http.Error(w, "Feature is disabled", http.StatusServiceUnavailable)
						return
					}
				}
			}

			// 检查速率限制
			if config.RateController != nil {
				if limiterName, ok := config.RateLimits[path]; ok {
					if !config.RateController.Allow(limiterName) {
						http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
						return
					}
				}
			}

			// 检查熔断器
			if config.CircuitController != nil {
				if breakerName, ok := config.CircuitBreakers[path]; ok {
					if config.CircuitController.IsOpen(breakerName) {
						http.Error(w, "Circuit breaker is open", http.StatusServiceUnavailable)
						return
					}
				}
			}

			// 继续处理请求（使用更新后的上下文）
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// shouldSkipControl 检查路径是否应该跳过控制
func shouldSkipControl(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// GetFeatureFlag 从上下文中获取功能开关状态
//
// 功能说明：
// - 获取当前请求路径对应的功能开关状态
// - 用于在处理器中检查功能是否启用
//
// 参数：
//   - ctx: 请求上下文
//   - flagName: 功能开关名称
//
// 返回：
//   - bool: 如果功能启用返回 true，否则返回 false
//
// 注意：此函数需要从上下文中获取功能控制器
// 实际使用中应该通过依赖注入或全局实例获取控制器
func GetFeatureFlag(ctx context.Context, flagName string) bool {
	// 从上下文获取功能控制器
	if controller, ok := ctx.Value("feature.controller").(*control.FeatureController); ok {
		return controller.IsEnabled(flagName)
	}
	// 默认返回 true（功能启用）
	return true
}
