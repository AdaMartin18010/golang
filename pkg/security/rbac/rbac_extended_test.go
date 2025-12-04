package rbac

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRBAC_GetNonExistentRole 测试获取不存在的角色
func TestRBAC_GetNonExistentRole(t *testing.T) {
	rbac := NewRBAC()

	// 获取不存在的角色应该返回错误
	_, err := rbac.GetRole("non-existent")
	assert.Error(t, err)
}

// TestRBAC_GetNonExistentPermission 测试获取不存在的权限
func TestRBAC_GetNonExistentPermission(t *testing.T) {
	rbac := NewRBAC()

	// 获取不存在的权限应该返回错误
	_, err := rbac.GetPermission("non-existent")
	assert.Error(t, err)
}

// TestRBAC_DuplicatePermission 测试重复添加权限
func TestRBAC_DuplicatePermission(t *testing.T) {
	rbac := NewRBAC()

	perm := &Permission{ID: "test-perm", Resource: "test", Action: "read"}
	require.NoError(t, rbac.AddPermission(perm))

	// 重复添加应该返回错误
	err := rbac.AddPermission(perm)
	assert.Error(t, err)
}

// TestRBAC_ComplexInheritance 测试复杂继承关系
func TestRBAC_ComplexInheritance(t *testing.T) {
	rbac := NewRBAC()

	// 创建角色层次: admin -> moderator -> user
	userRole := &Role{
		ID:          "user",
		Name:        "User",
		Permissions: []string{"read"},
	}
	require.NoError(t, rbac.AddRole(userRole))

	moderatorRole := &Role{
		ID:          "moderator",
		Name:        "Moderator",
		Permissions: []string{"write"},
		Inherits:    []string{"user"},
	}
	require.NoError(t, rbac.AddRole(moderatorRole))

	adminRole := &Role{
		ID:          "admin",
		Name:        "Admin",
		Permissions: []string{"delete"},
		Inherits:    []string{"moderator"},
	}
	require.NoError(t, rbac.AddRole(adminRole))

	// 添加权限
	require.NoError(t, rbac.AddPermission(&Permission{ID: "read", Resource: "data", Action: "read"}))
	require.NoError(t, rbac.AddPermission(&Permission{ID: "write", Resource: "data", Action: "write"}))
	require.NoError(t, rbac.AddPermission(&Permission{ID: "delete", Resource: "data", Action: "delete"}))

	ctx := context.Background()

	// Admin应该有所有权限
	hasRead, err := rbac.CheckPermission(ctx, []string{"admin"}, "data", "read")
	require.NoError(t, err)
	assert.True(t, hasRead)

	hasWrite, err := rbac.CheckPermission(ctx, []string{"admin"}, "data", "write")
	require.NoError(t, err)
	assert.True(t, hasWrite)

	hasDelete, err := rbac.CheckPermission(ctx, []string{"admin"}, "data", "delete")
	require.NoError(t, err)
	assert.True(t, hasDelete)
}

// TestRBAC_MultipleRoles 测试多角色权限检查
func TestRBAC_MultipleRoles(t *testing.T) {
	rbac := NewRBAC()
	require.NoError(t, rbac.InitializeDefaultRoles())

	ctx := context.Background()

	// 用户同时拥有user和moderator角色
	hasPermission, err := rbac.CheckPermission(ctx, []string{"user", "moderator"}, "user", "update")
	require.NoError(t, err)
	assert.True(t, hasPermission, "User with multiple roles should have combined permissions")
}

// TestRBAC_NonExistentRole 测试不存在的角色
func TestRBAC_NonExistentRole(t *testing.T) {
	rbac := NewRBAC()
	ctx := context.Background()

	hasPermission, err := rbac.CheckPermission(ctx, []string{"non-existent"}, "user", "read")
	require.NoError(t, err)
	assert.False(t, hasPermission)
}

// TestRBAC_WildcardPermission 测试通配符权限
func TestRBAC_WildcardPermission(t *testing.T) {
	rbac := NewRBAC()
	require.NoError(t, rbac.InitializeDefaultRoles())

	ctx := context.Background()

	// Admin有"all"权限，应该匹配所有资源和操作
	hasPermission, err := rbac.CheckPermission(ctx, []string{"admin"}, "any-resource", "any-action")
	require.NoError(t, err)
	assert.True(t, hasPermission)
}

// TestRBAC_EmptyRolesList 测试空角色列表
func TestRBAC_EmptyRolesList(t *testing.T) {
	rbac := NewRBAC()
	ctx := context.Background()

	// 空角色列表不应该有任何权限
	hasPermission, err := rbac.CheckPermission(ctx, []string{}, "user", "read")
	require.NoError(t, err)
	assert.False(t, hasPermission)
}

// TestRBAC_RoleWithNoPermissions 测试没有权限的角色
func TestRBAC_RoleWithNoPermissions(t *testing.T) {
	rbac := NewRBAC()

	// 添加没有权限的角色
	role := &Role{
		ID:          "empty-role",
		Name:        "Empty Role",
		Permissions: []string{},
	}
	require.NoError(t, rbac.AddRole(role))

	ctx := context.Background()
	hasPermission, err := rbac.CheckPermission(ctx, []string{"empty-role"}, "user", "read")
	require.NoError(t, err)
	assert.False(t, hasPermission)
}

// TestRBAC_ConcurrentAccess 测试并发访问
func TestRBAC_ConcurrentAccess(t *testing.T) {
	rbac := NewRBAC()
	require.NoError(t, rbac.InitializeDefaultRoles())

	ctx := context.Background()
	const numGoroutines = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := rbac.CheckPermission(ctx, []string{"user"}, "user", "read")
			assert.NoError(t, err)
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

// BenchmarkRBAC_CheckPermissionWithInheritance 性能测试 - 带继承的权限检查
func BenchmarkRBAC_CheckPermissionWithInheritance(b *testing.B) {
	rbac := NewRBAC()
	rbac.InitializeDefaultRoles()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbac.CheckPermission(ctx, []string{"moderator"}, "user", "read")
	}
}

// BenchmarkRBAC_CheckPermissionMultipleRoles 性能测试 - 多角色权限检查
func BenchmarkRBAC_CheckPermissionMultipleRoles(b *testing.B) {
	rbac := NewRBAC()
	rbac.InitializeDefaultRoles()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbac.CheckPermission(ctx, []string{"user", "moderator", "admin"}, "user", "read")
	}
}
