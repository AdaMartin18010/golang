package rbac

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRBAC(t *testing.T) {
	rbac := NewRBAC()
	assert.NotNil(t, rbac)
	assert.NotNil(t, rbac.roles)
	assert.NotNil(t, rbac.permissions)
}

func TestRBAC_AddRole(t *testing.T) {
	rbac := NewRBAC()

	role := &Role{
		ID:          "admin",
		Name:        "Administrator",
		Description: "Admin role",
		Permissions: []string{"all"},
	}

	err := rbac.AddRole(role)
	require.NoError(t, err)

	// 验证角色已添加
	retrieved, err := rbac.GetRole("admin")
	require.NoError(t, err)
	assert.Equal(t, role.ID, retrieved.ID)
	assert.Equal(t, role.Name, retrieved.Name)
}

func TestRBAC_AddRole_Duplicate(t *testing.T) {
	rbac := NewRBAC()

	role := &Role{
		ID:   "admin",
		Name: "Administrator",
	}

	err := rbac.AddRole(role)
	require.NoError(t, err)

	// 重复添加应该返回错误
	err = rbac.AddRole(role)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRBAC_AddPermission(t *testing.T) {
	rbac := NewRBAC()

	perm := &Permission{
		ID:          "user.create",
		Resource:    "user",
		Action:      "create",
		Description: "Create users",
	}

	err := rbac.AddPermission(perm)
	require.NoError(t, err)

	// 验证权限已添加
	retrieved, err := rbac.GetPermission("user.create")
	require.NoError(t, err)
	assert.Equal(t, perm.ID, retrieved.ID)
	assert.Equal(t, perm.Resource, retrieved.Resource)
}

func TestRBAC_AssignPermissionToRole(t *testing.T) {
	rbac := NewRBAC()

	// 添加角色
	role := &Role{ID: "user", Name: "User"}
	require.NoError(t, rbac.AddRole(role))

	// 添加权限
	perm := &Permission{ID: "user.read", Resource: "user", Action: "read"}
	require.NoError(t, rbac.AddPermission(perm))

	// 分配权限
	err := rbac.AssignPermissionToRole("user", "user.read")
	require.NoError(t, err)

	// 验证分配
	retrieved, err := rbac.GetRole("user")
	require.NoError(t, err)
	assert.Contains(t, retrieved.Permissions, "user.read")
}

func TestRBAC_CheckPermission(t *testing.T) {
	rbac := NewRBAC()
	ctx := context.Background()

	// 初始化默认角色
	require.NoError(t, rbac.InitializeDefaultRoles())

	tests := []struct {
		name       string
		userRoles  []string
		resource   string
		action     string
		expected   bool
		expectErr  bool
	}{
		{
			name:      "admin has all permissions",
			userRoles: []string{"admin"},
			resource:  "user",
			action:    "create",
			expected:  true,
		},
		{
			name:      "user can read",
			userRoles: []string{"user"},
			resource:  "user",
			action:    "read",
			expected:  true,
		},
		{
			name:      "user cannot delete",
			userRoles: []string{"user"},
			resource:  "user",
			action:    "delete",
			expected:  false,
		},
		{
			name:      "moderator can update (inherited)",
			userRoles: []string{"moderator"},
			resource:  "user",
			action:    "read",
			expected:  true,
		},
		{
			name:      "no roles",
			userRoles: []string{},
			resource:  "user",
			action:    "read",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPermission, err := rbac.CheckPermission(ctx, tt.userRoles, tt.resource, tt.action)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, hasPermission)
			}
		})
	}
}

func TestRBAC_RoleInheritance(t *testing.T) {
	rbac := NewRBAC()

	// 创建基础角色
	baseRole := &Role{
		ID:          "base",
		Name:        "Base",
		Permissions: []string{"perm1"},
	}
	require.NoError(t, rbac.AddRole(baseRole))

	// 创建继承角色
	derivedRole := &Role{
		ID:          "derived",
		Name:        "Derived",
		Permissions: []string{"perm2"},
		Inherits:    []string{"base"},
	}
	require.NoError(t, rbac.AddRole(derivedRole))

	// 添加权限
	require.NoError(t, rbac.AddPermission(&Permission{
		ID: "perm1", Resource: "res1", Action: "action1",
	}))
	require.NoError(t, rbac.AddPermission(&Permission{
		ID: "perm2", Resource: "res2", Action: "action2",
	}))

	// 验证继承的权限
	ctx := context.Background()
	hasPermission, err := rbac.CheckPermission(ctx, []string{"derived"}, "res1", "action1")
	require.NoError(t, err)
	assert.True(t, hasPermission, "Derived role should inherit permissions from base role")
}

func TestRBAC_ListRoles(t *testing.T) {
	rbac := NewRBAC()

	// 添加几个角色
	roles := []*Role{
		{ID: "role1", Name: "Role 1"},
		{ID: "role2", Name: "Role 2"},
		{ID: "role3", Name: "Role 3"},
	}

	for _, role := range roles {
		require.NoError(t, rbac.AddRole(role))
	}

	// 获取所有角色
	allRoles := rbac.ListRoles()
	assert.Len(t, allRoles, 3)
}

func TestRBAC_ListPermissions(t *testing.T) {
	rbac := NewRBAC()

	// 添加几个权限
	perms := []*Permission{
		{ID: "perm1", Resource: "res1", Action: "act1"},
		{ID: "perm2", Resource: "res2", Action: "act2"},
	}

	for _, perm := range perms {
		require.NoError(t, rbac.AddPermission(perm))
	}

	// 获取所有权限
	allPerms := rbac.ListPermissions()
	assert.Len(t, allPerms, 2)
}

func BenchmarkRBAC_CheckPermission(b *testing.B) {
	rbac := NewRBAC()
	rbac.InitializeDefaultRoles()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbac.CheckPermission(ctx, []string{"user"}, "user", "read")
	}
}
