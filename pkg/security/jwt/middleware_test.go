package jwt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMiddleware(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
		SkipPaths:    []string{"/health", "/metrics"},
	})

	assert.NotNil(t, m)
	assert.NotNil(t, m.tokenManager)
	assert.True(t, m.skipPaths["/health"])
	assert.True(t, m.skipPaths["/metrics"])
}

func TestMiddleware_Authenticate_Success(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	// 生成有效令牌
	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
	require.NoError(t, err)

	// 创建测试处理函数
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证 claims 已添加到上下文
		claims, ok := GetClaims(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "user-123", claims.UserID)
		w.WriteHeader(http.StatusOK)
	}))

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_Authenticate_MissingHeader(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMiddleware_Authenticate_InvalidFormat(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	tests := []struct {
		name  string
		token string
	}{
		{"no bearer prefix", "invalid-token"},
		{"wrong scheme", "Basic dXNlcjpwYXNz"},
		{"empty token", "Bearer "},
		{"only one part", "BearerToken"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.token)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}

func TestMiddleware_Authenticate_InvalidToken(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMiddleware_Authenticate_SkipPath(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
		SkipPaths:    []string{"/health"},
	})

	called := false
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_RequireRoles(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	tests := []struct {
		name       string
		userRoles  []string
		required   []string
		wantStatus int
	}{
		{"has required role", []string{"admin", "user"}, []string{"admin"}, http.StatusOK},
		{"has one of required", []string{"user"}, []string{"admin", "user"}, http.StatusOK},
		{"missing role", []string{"user"}, []string{"admin"}, http.StatusForbidden},
		{"empty roles", []string{}, []string{"admin"}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := m.RequireRoles(tt.required...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			ctx := WithUserRoles(req.Context(), tt.userRoles)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestMiddleware_RequireRoles_NoClaims(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	m := NewMiddleware(MiddlewareConfig{
		TokenManager: tm,
	})

	handler := m.RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetClaims(t *testing.T) {
	claims := &Claims{
		UserID:   "user-123",
		Username: "john",
		Roles:    []string{"user"},
	}

	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	retrieved, ok := GetClaims(ctx)
	assert.True(t, ok)
	assert.Equal(t, claims.UserID, retrieved.UserID)
}

func TestGetClaims_NotFound(t *testing.T) {
	_, ok := GetClaims(context.Background())
	assert.False(t, ok)
}

func TestGetUserID(t *testing.T) {
	claims := &Claims{
		UserID: "user-123",
	}

	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	userID, ok := GetUserID(ctx)
	assert.True(t, ok)
	assert.Equal(t, "user-123", userID)
}

func TestGetUserID_NotFound(t *testing.T) {
	_, ok := GetUserID(context.Background())
	assert.False(t, ok)
}

func TestGetUserRoles(t *testing.T) {
	claims := &Claims{
		Roles: []string{"admin", "user"},
	}

	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	roles, ok := GetUserRoles(ctx)
	assert.True(t, ok)
	assert.Equal(t, []string{"admin", "user"}, roles)
}

func TestGetUserRoles_NotFound(t *testing.T) {
	_, ok := GetUserRoles(context.Background())
	assert.False(t, ok)
}

// Context key 用于测试
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	claims := &Claims{
		Roles: roles,
	}
	return context.WithValue(ctx, ClaimsKey, claims)
}
