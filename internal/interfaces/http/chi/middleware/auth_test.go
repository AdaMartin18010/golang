package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/golang/pkg/security/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	config := jwt.Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := jwt.NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name           string
		authHeader     string
		skipPaths      []string
		optionalAuth   bool
		expectedStatus int
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "skip path",
			authHeader:     "",
			skipPaths:      []string{"/public"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "optional auth without token",
			authHeader:     "",
			optionalAuth:   true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			
			// 使用 jwt 包的中间件
			jwtMiddleware := jwt.NewMiddleware(jwt.MiddlewareConfig{
				JWT:          j,
				SkipPaths:    tt.skipPaths,
				OptionalAuth: tt.optionalAuth,
			})
			r.Use(jwtMiddleware.Authenticate)
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			r.Get("/public", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			if len(tt.skipPaths) > 0 {
				req = httptest.NewRequest("GET", "/public", nil)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name:           "has required role",
			userRoles:      []string{"user", "admin"},
			requiredRoles:  []string{"admin"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "does not have required role",
			userRoles:      []string{"user"},
			requiredRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "has any of required roles",
			userRoles:      []string{"user"},
			requiredRoles:  []string{"user", "admin"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := jwt.Config{
				SecretKey:      "test-secret-key",
				SigningMethod:  "HS256",
				AccessTokenTTL: 15 * time.Minute,
			}
			j, _ := jwt.NewJWT(config)
			jwtMiddleware := jwt.NewMiddleware(jwt.MiddlewareConfig{
				JWT: j,
			})

			r := chi.NewRouter()
			// 先设置角色到上下文，然后检查角色
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), "roles", tt.userRoles)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			})
			r.Use(jwtMiddleware.RequireRoles(tt.requiredRoles...))
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", "user-123")
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", userID)
	}
}

func TestGetRoles(t *testing.T) {
	roles := []string{"user", "admin"}
	ctx := context.WithValue(context.Background(), "roles", roles)
	gotRoles, ok := ctx.Value("roles").([]string)
	if !ok || len(gotRoles) != len(roles) {
		t.Errorf("Expected %d roles, got %d", len(roles), len(gotRoles))
	}
}
