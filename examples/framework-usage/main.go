package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
	"github.com/yourusername/golang/pkg/auth/jwt"
	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/eventbus"
	"github.com/yourusername/golang/pkg/health"
	"github.com/yourusername/golang/pkg/http/response"
	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/rbac"
	"github.com/yourusername/golang/pkg/validator"
	"log/slog"
)

func main() {
	// 1. 创建日志记录器
	log := logger.NewLogger(slog.LevelInfo)
	log.Info("Starting application")

	// 2. 创建JWT管理器
	jwtManager, err := jwt.NewJWT(jwt.Config{
		SecretKey:      "your-secret-key-change-in-production",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "example-app",
		Audience:       "example-users",
	})
	if err != nil {
		log.Error("Failed to create JWT manager", "error", err)
		return
	}

	// 3. 创建健康检查器
	healthChecker := health.NewHealthChecker()
	healthChecker.Register(health.NewSimpleCheck("app", func(ctx context.Context) error {
		// 这里可以检查数据库、缓存等
		return nil
	}))

	// 4. 创建事件总线
	eventBus := eventbus.NewEventBus(100)
	eventBus.Start()
	defer eventBus.Stop()

	// 订阅用户创建事件
	eventBus.Subscribe("user.created", func(ctx context.Context, event eventbus.Event) error {
		log.Info("User created event received", "event_type", event.Type(), "data", event.Data())
		// 这里可以发送邮件、更新缓存等
		return nil
	})

	// 5. 创建RBAC
	rbacSystem := rbac.NewRBAC()

	// 创建权限
	readUsersPerm := &rbac.Permission{
		ID:       "perm_read_users",
		Name:     "read_users",
		Resource: "users",
		Action:   "read",
	}
	writeUsersPerm := &rbac.Permission{
		ID:       "perm_write_users",
		Name:     "write_users",
		Resource: "users",
		Action:   "write",
	}
	rbacSystem.AddPermission(readUsersPerm)
	rbacSystem.AddPermission(writeUsersPerm)

	// 创建角色
	adminRole := &rbac.Role{
		ID:   "role_admin",
		Name: "admin",
	}
	userRole := &rbac.Role{
		ID:   "role_user",
		Name: "user",
	}
	rbacSystem.AddRole(adminRole)
	rbacSystem.AddRole(userRole)

	// 分配权限
	rbacSystem.AssignPermission("role_admin", "perm_read_users")
	rbacSystem.AssignPermission("role_admin", "perm_write_users")
	rbacSystem.AssignPermission("role_user", "perm_read_users")

	enforcer := rbac.NewEnforcer(rbacSystem)

	// 6. 创建路由
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.TracingMiddleware(middleware.TracingConfig{
		TracerName:     "example-app",
		ServiceName:    "example-service",
		ServiceVersion: "v1.0.0",
		SkipPaths:      []string{"/health", "/metrics"},
	}))

	metrics := middleware.NewMetrics()
	r.Use(middleware.MetricsMiddleware(metrics))

	r.Use(middleware.AuthMiddleware(middleware.AuthConfig{
		JWT:       jwtManager,
		SkipPaths: []string{"/health", "/metrics", "/api/v1/auth/login"},
	}))

	// 健康检查
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := healthChecker.OverallStatus(r.Context())
		results := healthChecker.Check(r.Context())

		result := map[string]interface{}{
			"status": status,
			"checks": results,
		}

		code := http.StatusOK
		if status == health.StatusUnhealthy {
			code = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(result)
	})

	// 指标
	r.Get("/metrics", middleware.MetricsHandler(metrics))

	// 认证路由
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", createLoginHandler(jwtManager))
	})

	// API路由
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", createGetUsersHandler(enforcer))
		r.Post("/users", createCreateUserHandler(validator.NewValidator(), eventBus, enforcer))
	})

	// 启动服务器
	log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed", "error", err)
	}
}

// createLoginHandler 创建登录处理器
func createLoginHandler(jwtManager *jwt.JWT) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest,
				errors.NewInvalidInputError("invalid request body"))
			return
		}

		// 验证用户凭据（简化示例）
		if req.Username != "admin" || req.Password != "password" {
			response.Error(w, http.StatusUnauthorized,
				errors.NewUnauthorizedError("invalid credentials"))
			return
		}

		// 生成Token
		accessToken, err := jwtManager.GenerateAccessToken(
			"user-123",
			req.Username,
			[]string{"admin"},
			"admin@example.com",
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError,
				errors.NewInternalError("failed to generate token", err))
			return
		}

		refreshToken, err := jwtManager.GenerateRefreshToken("user-123")
		if err != nil {
			response.Error(w, http.StatusInternalServerError,
				errors.NewInternalError("failed to generate refresh token", err))
			return
		}

		response.Success(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

// createGetUsersHandler 创建获取用户列表处理器
func createGetUsersHandler(enforcer *rbac.Enforcer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从context获取用户信息（从JWT中间件注入）
		userID := middleware.GetUserID(r.Context())
		roles := middleware.GetRoles(r.Context())

		if userID == "" {
			response.Error(w, http.StatusUnauthorized,
				errors.NewUnauthorizedError("user not found"))
			return
		}

		// 检查权限
		rbacUser := &rbac.User{
			ID:    userID,
			Roles: roles,
		}
		if err := enforcer.Enforce(rbacUser, "users", "read"); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		// 返回用户列表（简化示例）
		users := []map[string]interface{}{
			{"id": "1", "name": "John"},
			{"id": "2", "name": "Jane"},
		}

		response.Success(w, http.StatusOK, users)
	}
}

// createCreateUserHandler 创建用户处理器
func createCreateUserHandler(
	validator *validator.Validator,
	eventBus *eventbus.EventBus,
	enforcer *rbac.Enforcer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name  string `json:"name" validate:"required,min=2,max=50"`
			Email string `json:"email" validate:"required,email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest,
				errors.NewInvalidInputError("invalid request body"))
			return
		}

		// 验证请求
		if !validator.Validate(req) {
			validationErrors := validator.ValidateStruct(req)
			details := make(map[string]interface{})
			for _, err := range validationErrors {
				details[err.Field] = err.Message
			}
			response.Error(w, http.StatusBadRequest,
				errors.NewValidationError("validation failed", details))
			return
		}

		// 检查权限
		userID := middleware.GetUserID(r.Context())
		roles := middleware.GetRoles(r.Context())

		if userID == "" {
			response.Error(w, http.StatusUnauthorized,
				errors.NewUnauthorizedError("user not found"))
			return
		}

		rbacUser := &rbac.User{
			ID:    userID,
			Roles: roles,
		}
		if err := enforcer.Enforce(rbacUser, "users", "write"); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		// 创建用户（简化示例）
		newUser := map[string]interface{}{
			"id":    "3",
			"name":  req.Name,
			"email": req.Email,
		}

		// 发布事件
		event := eventbus.NewEvent("user.created", newUser)
		eventBus.Publish(event)

		response.Success(w, http.StatusCreated, newUser)
	}
}
