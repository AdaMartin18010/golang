package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"

	"example.com/golang-examples/framework-usage/middleware"
	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/eventbus"
	"github.com/yourusername/golang/pkg/health"
	"github.com/yourusername/golang/pkg/http/response"
	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/security/jwt"
	"github.com/yourusername/golang/pkg/validator"
)

func main() {
	// 1. 创建日志记录器
	log := logger.NewLogger(slog.LevelInfo)
	log.Info("Starting application")

	// 2. 创建JWT管理器
	jwtManager, err := jwt.NewTokenManager(jwt.Config{
		PrivateKeyPath:  "",
		PublicKeyPath:   "",
		Issuer:          "example-app",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		SigningMethod:   "HS256",
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

	// 5. 创建路由
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
		TokenManager: jwtManager,
		SkipPaths:    []string{"/health", "/metrics", "/api/v1/auth/login"},
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
		r.Get("/users", createGetUsersHandler())
		r.Post("/users", createCreateUserHandler(validator.NewValidator(), eventBus))
	})

	// 启动服务器
	log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed", "error", err)
	}
}

// createLoginHandler 创建登录处理器
func createLoginHandler(jwtManager *jwt.TokenManager) http.HandlerFunc {
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
			"admin@example.com",
			[]string{"admin"},
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
func createGetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从context获取用户信息（从JWT中间件注入）
		userID := middleware.GetUserID(r.Context())

		if userID == "" {
			response.Error(w, http.StatusUnauthorized,
				errors.NewUnauthorizedError("user not found"))
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
	val *validator.Validator,
	eventBus *eventbus.EventBus,
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

		// 验证请求（简化验证）
		if req.Name == "" || len(req.Name) < 2 {
			response.Error(w, http.StatusBadRequest,
				errors.NewValidationError("validation failed", map[string]interface{}{
					"name": "name is required and must be at least 2 characters",
				}))
			return
		}
		if req.Email == "" || !validator.ValidateEmail(req.Email) {
			response.Error(w, http.StatusBadRequest,
				errors.NewValidationError("validation failed", map[string]interface{}{
					"email": "valid email is required",
				}))
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
