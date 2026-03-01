package rbac

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMiddleware(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	assert.NotNil(t, m)
	assert.Equal(t, rbac, m.rbac)
}

func TestMiddleware_RequirePermission(t *testing.T) {
	rbac := NewRBAC()
	require.NoError(t, rbac.InitializeDefaultRoles())

	m := NewMiddleware(rbac)

	tests := []struct {
		name       string
		userRoles  []string
		resource   string
		action     string
		wantStatus int
	}{
		{"admin can delete", []string{"admin"}, "user", "delete", http.StatusOK},
		{"user can read", []string{"user"}, "user", "read", http.StatusOK},
		{"user cannot delete", []string{"user"}, "user", "delete", http.StatusForbidden},
		{"moderator can update", []string{"moderator"}, "user", "update", http.StatusOK},
		{"no roles", []string{}, "user", "read", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := m.RequirePermission(tt.resource, tt.action)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestMiddleware_RequirePermission_NoRolesInContext(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	handler := m.RequirePermission("user", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMiddleware_RequirePermission_RBACError(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	// 不初始化角色，直接检查权限会导致错误
	handler := m.RequirePermission("user", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := WithUserRoles(req.Context(), []string{"nonexistent"})
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	// 由于我们修改了实现以忽略不存在的角色，这应该返回 Forbidden
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestMiddleware_RequireRole(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	tests := []struct {
		name         string
		userRoles    []string
		requiredRole string
		wantStatus   int
	}{
		{"has exact role", []string{"admin"}, "admin", http.StatusOK},
		{"has one of roles", []string{"user", "admin"}, "admin", http.StatusOK},
		{"missing role", []string{"user"}, "admin", http.StatusForbidden},
		{"empty user roles", []string{}, "admin", http.StatusUnauthorized},
		{"multiple required one match", []string{"user"}, "admin", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := m.RequireRole(tt.requiredRole)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestMiddleware_RequireMultipleRoles(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	handler := m.RequireRole("admin", "moderator")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 测试只有一个匹配的角色
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := WithUserRoles(req.Context(), []string{"user", "moderator"})
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_RequireRole_NoContext(t *testing.T) {
	rbac := NewRBAC()
	m := NewMiddleware(rbac)

	handler := m.RequireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestContextHelpers(t *testing.T) {
	t.Run("WithUserRoles and GetUserRoles", func(t *testing.T) {
		roles := []string{"admin", "user"}
		ctx := WithUserRoles(context.Background(), roles)

		retrieved, ok := GetUserRoles(ctx)
		assert.True(t, ok)
		assert.Equal(t, roles, retrieved)
	})

	t.Run("GetUserRoles not found", func(t *testing.T) {
		_, ok := GetUserRoles(context.Background())
		assert.False(t, ok)
	})

	t.Run("GetUserRoles wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserRolesKey, "not-a-slice")
		_, ok := GetUserRoles(ctx)
		assert.False(t, ok)
	})

	t.Run("WithUserID and GetUserID", func(t *testing.T) {
		userID := "user-123"
		ctx := WithUserID(context.Background(), userID)

		retrieved, ok := GetUserID(ctx)
		assert.True(t, ok)
		assert.Equal(t, userID, retrieved)
	})

	t.Run("GetUserID not found", func(t *testing.T) {
		_, ok := GetUserID(context.Background())
		assert.False(t, ok)
	})

	t.Run("GetUserID wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, 123)
		_, ok := GetUserID(ctx)
		assert.False(t, ok)
	})
}

func TestContextKeys(t *testing.T) {
	// 验证上下文键的唯一性
	assert.NotEqual(t, string(UserRolesKey), string(UserIDKey))
	assert.Equal(t, "user_roles", string(UserRolesKey))
	assert.Equal(t, "user_id", string(UserIDKey))
}

func BenchmarkMiddleware_CheckPermission(b *testing.B) {
	rbac := NewRBAC()
	rbac.InitializeDefaultRoles()
	m := NewMiddleware(rbac)

	handler := m.RequirePermission("user", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := WithUserRoles(req.Context(), []string{"user"})
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
