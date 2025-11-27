// Package middleware provides HTTP middleware for the Chi router.
//
// HTTP 中间件包提供了各种 HTTP 请求处理中间件，包括：
// 1. 认证授权：JWT Token 认证和角色权限控制
// 2. 限流保护：请求限流（支持多种算法和分布式限流）
// 3. 熔断保护：服务熔断器（三种状态：关闭、开启、半开）
// 4. 请求追踪：基于 OpenTelemetry 的链路追踪
// 5. 性能监控：请求指标收集
// 6. 错误恢复：Panic 恢复和错误处理
// 7. CORS：跨域资源共享支持
//
// 设计原则：
// 1. 可配置：所有中间件都支持灵活的配置选项
// 2. 可组合：中间件可以组合使用
// 3. 高性能：最小化对请求处理性能的影响
// 4. 可观测：支持日志、追踪和指标收集
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourusername/golang/pkg/auth/jwt"
	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// AuthConfig 是认证中间件的配置结构。
//
// 功能说明：
// - 配置 JWT 认证和授权行为
// - 支持跳过特定路径的认证
// - 支持可选认证（不强制要求 Token）
//
// 字段说明：
// - JWT: JWT 认证器实例，用于验证 Token
// - SkipPaths: 跳过认证的路径列表（如 /health、/metrics）
// - OptionalAuth: 是否可选认证
//   - true: Token 不存在或无效时继续处理请求（不返回错误）
//   - false: Token 不存在或无效时返回 401 错误
//
// 使用示例：
//
//	config := middleware.AuthConfig{
//	    JWT:          jwtAuth,
//	    SkipPaths:    []string{"/health", "/metrics"},
//	    OptionalAuth: false,
//	}
//	router.Use(middleware.AuthMiddleware(config))
type AuthConfig struct {
	JWT          *jwt.JWT
	SkipPaths    []string // 跳过认证的路径
	OptionalAuth bool     // 是否可选认证（不强制）
}

// AuthMiddleware 创建认证中间件。
//
// 功能说明：
// - 验证 JWT Token
// - 将用户信息添加到请求上下文
// - 支持跳过特定路径
// - 支持可选认证模式
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 从请求中提取 Token（Authorization header、Cookie、查询参数）
// 3. 验证 Token（如果存在）
// 4. 将用户信息添加到上下文
// 5. 继续处理请求
//
// Token 提取顺序：
// 1. Authorization header（Bearer token）
// 2. Cookie（access_token）
// 3. 查询参数（token）
//
// 上下文键：
// - "user_id": 用户 ID
// - "username": 用户名
// - "roles": 用户角色列表
// - "email": 用户邮箱
// - "claims": JWT Claims 对象
//
// 参数：
// - config: 认证配置
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	config := middleware.AuthConfig{
//	    JWT:       jwtAuth,
//	    SkipPaths: []string{"/health", "/metrics"},
//	}
//	router.Use(middleware.AuthMiddleware(config))
//
//	// 在 Handler 中获取用户信息
//	userID := middleware.GetUserID(r.Context())
//	roles := middleware.GetRoles(r.Context())
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否跳过认证
			if shouldSkipAuth(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 提取Token
			tokenString := extractToken(r)
			if tokenString == "" {
				if config.OptionalAuth {
					next.ServeHTTP(w, r)
					return
				}
				response.Error(w, http.StatusUnauthorized,
					errors.NewUnauthorizedError("missing or invalid authorization token"))
				return
			}

			// 验证Token
			claims, err := config.JWT.ValidateToken(tokenString)
			if err != nil {
				if config.OptionalAuth {
					next.ServeHTTP(w, r)
					return
				}
				if err == jwt.ErrExpiredToken {
					response.Error(w, http.StatusUnauthorized,
						errors.NewUnauthorizedError("token expired"))
					return
				}
				response.Error(w, http.StatusUnauthorized,
					errors.NewUnauthorizedError("invalid token"))
				return
			}

			// 将用户信息添加到context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "roles", claims.Roles)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole 创建要求特定角色的中间件。
//
// 功能说明：
// - 检查用户是否具有任一指定的角色
// - 如果没有权限，返回 403 Forbidden
//
// 参数：
// - roles: 允许的角色列表（用户只需具有其中一个角色即可）
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	// 要求用户具有 admin 或 moderator 角色
//	router.With(middleware.RequireRole("admin", "moderator")).Get("/admin", adminHandler)
//
// 注意事项：
// - 必须在 AuthMiddleware 之后使用
// - 如果用户未认证，会返回 403 错误
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("roles").([]string)
			if !ok {
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("user roles not found"))
				return
			}

			// 检查用户是否有任一所需角色
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range userRoles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole 创建要求任一角色的中间件（RequireRole 的别名）。
//
// 功能说明：
// - 与 RequireRole 功能相同
// - 检查用户是否具有任一指定的角色
//
// 参数：
// - roles: 允许的角色列表
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	router.With(middleware.RequireAnyRole("admin", "moderator")).Get("/admin", adminHandler)
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return RequireRole(roles...)
}

// RequireAllRoles 创建要求所有角色的中间件。
//
// 功能说明：
// - 检查用户是否具有所有指定的角色
// - 如果缺少任一角色，返回 403 Forbidden
//
// 参数：
// - roles: 必需的角色列表（用户必须具有所有角色）
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	// 要求用户同时具有 admin 和 super_admin 角色
//	router.With(middleware.RequireAllRoles("admin", "super_admin")).Get("/super", superHandler)
//
// 注意事项：
// - 必须在 AuthMiddleware 之后使用
// - 与 RequireRole 不同，这里要求用户具有所有角色
func RequireAllRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("roles").([]string)
			if !ok {
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("user roles not found"))
				return
			}

			// 检查用户是否有所有所需角色
			roleMap := make(map[string]bool)
			for _, role := range userRoles {
				roleMap[role] = true
			}

			for _, requiredRole := range roles {
				if !roleMap[requiredRole] {
					response.Error(w, http.StatusForbidden,
						errors.NewForbiddenError("insufficient permissions"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken 从 HTTP 请求中提取 JWT Token。
//
// 功能说明：
// - 支持多种 Token 提取方式
// - 按优先级顺序尝试提取
//
// 提取顺序（优先级从高到低）：
// 1. Authorization header（格式：Bearer <token>）
// 2. Cookie（名称为 access_token）
// 3. 查询参数（参数名为 token）
//
// 参数：
// - r: HTTP 请求
//
// 返回：
// - string: 提取到的 Token，如果不存在则返回空字符串
//
// 使用场景：
// - 认证中间件内部使用
// - 支持多种客户端 Token 传递方式
func extractToken(r *http.Request) string {
	// 1. 从 Authorization header 提取
	// 格式：Authorization: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 2. 从 Cookie 提取
	// Cookie 名称：access_token
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 3. 从查询参数提取
	// 参数名：token
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

// shouldSkipAuth 检查指定路径是否应该跳过认证。
//
// 功能说明：
// - 检查路径是否在跳过列表中
// - 支持精确匹配和前缀匹配
//
// 匹配规则：
// - 精确匹配：路径完全相等
// - 前缀匹配：路径以跳过路径开头
//
// 参数：
// - path: 要检查的路径
// - skipPaths: 跳过认证的路径列表
//
// 返回：
// - bool: 如果应该跳过认证返回 true，否则返回 false
//
// 使用示例：
//
//	skipPaths := []string{"/health", "/metrics", "/public"}
//	shouldSkipAuth("/health", skipPaths)        // true
//	shouldSkipAuth("/metrics", skipPaths)       // true
//	shouldSkipAuth("/public/api", skipPaths)    // true（前缀匹配）
//	shouldSkipAuth("/api/users", skipPaths)     // false
func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// GetUserID 从请求上下文中获取用户 ID。
//
// 功能说明：
// - 从上下文提取用户 ID
// - 由 AuthMiddleware 设置
//
// 参数：
// - ctx: 请求上下文
//
// 返回：
// - string: 用户 ID，如果不存在则返回空字符串
//
// 使用示例：
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    userID := middleware.GetUserID(r.Context())
//	    if userID == "" {
//	        // 用户未认证
//	    }
//	}
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUsername 从请求上下文中获取用户名。
//
// 功能说明：
// - 从上下文提取用户名
// - 由 AuthMiddleware 设置
//
// 参数：
// - ctx: 请求上下文
//
// 返回：
// - string: 用户名，如果不存在则返回空字符串
//
// 使用示例：
//
//	username := middleware.GetUsername(r.Context())
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value("username").(string); ok {
		return username
	}
	return ""
}

// GetRoles 从请求上下文中获取用户角色列表。
//
// 功能说明：
// - 从上下文提取用户角色列表
// - 由 AuthMiddleware 设置
//
// 参数：
// - ctx: 请求上下文
//
// 返回：
// - []string: 用户角色列表，如果不存在则返回 nil
//
// 使用示例：
//
//	roles := middleware.GetRoles(r.Context())
//	for _, role := range roles {
//	    // 处理角色
//	}
func GetRoles(ctx context.Context) []string {
	if roles, ok := ctx.Value("roles").([]string); ok {
		return roles
	}
	return nil
}

// GetEmail 从请求上下文中获取用户邮箱。
//
// 功能说明：
// - 从上下文提取用户邮箱
// - 由 AuthMiddleware 设置
//
// 参数：
// - ctx: 请求上下文
//
// 返回：
// - string: 用户邮箱，如果不存在则返回空字符串
//
// 使用示例：
//
//	email := middleware.GetEmail(r.Context())
func GetEmail(ctx context.Context) string {
	if email, ok := ctx.Value("email").(string); ok {
		return email
	}
	return ""
}

// GetClaims 从请求上下文中获取完整的 JWT Claims。
//
// 功能说明：
// - 从上下文提取完整的 JWT Claims 对象
// - 包含所有 Token 中的声明信息
// - 由 AuthMiddleware 设置
//
// 参数：
// - ctx: 请求上下文
//
// 返回：
// - *jwt.Claims: JWT Claims 对象，如果不存在则返回 nil
//
// 使用示例：
//
//	claims := middleware.GetClaims(r.Context())
//	if claims != nil {
//	    // 访问自定义声明
//	}
func GetClaims(ctx context.Context) *jwt.Claims {
	if claims, ok := ctx.Value("claims").(*jwt.Claims); ok {
		return claims
	}
	return nil
}
