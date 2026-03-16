package rbac

import (
	"context"
	"testing"
)

func TestRBAC_Enhanced(t *testing.T) {
	rbac := NewRBAC()

	// 创建角色
	adminRole := &Role{
		ID:   "admin",
		Name: "Administrator",
		Permissions: []*Permission{
			{ID: "user:read", Name: "Read Users", Resource: "user", Action: "read"},
			{ID: "user:write", Name: "Write Users", Resource: "user", Action: "write"},
			{ID: "admin:all", Name: "Admin All", Resource: "*", Action: "*"},
		},
	}

	userRole := &Role{
		ID:   "user",
		Name: "User",
		Permissions: []*Permission{
			{ID: "user:read", Name: "Read Users", Resource: "user", Action: "read"},
		},
	}

	rbac.AddRole(adminRole)
	rbac.AddRole(userRole)

	// 创建用户
	adminUser := &User{
		ID:    "user-1",
		Roles: []string{"admin"},
	}

	regularUser := &User{
		ID:    "user-2",
		Roles: []string{"user"},
	}

	// 测试管理员权限
	if !rbac.CheckPermission(adminUser, "user", "read") {
		t.Error("Admin should have user:read permission")
	}

	if !rbac.CheckPermission(adminUser, "user", "write") {
		t.Error("Admin should have user:write permission")
	}

	if !rbac.CheckPermission(adminUser, "any", "any") {
		t.Error("Admin should have all permissions")
	}

	// 测试普通用户权限
	if !rbac.CheckPermission(regularUser, "user", "read") {
		t.Error("User should have user:read permission")
	}

	if rbac.CheckPermission(regularUser, "user", "write") {
		t.Error("User should not have user:write permission")
	}

	// 测试权限检查（通过 ID）
	if !rbac.CheckPermissionByID(adminUser, "user:read") {
		t.Error("Admin should have user:read permission by ID")
	}

	if !rbac.CheckPermissionByID(adminUser, "admin:all") {
		t.Error("Admin should have admin:all permission")
	}

	// 测试获取用户权限
	permissions := rbac.GetUserPermissions(adminUser)
	if len(permissions) != 3 {
		t.Errorf("Expected 3 permissions, got %d", len(permissions))
	}

	// 测试 Enforcer
	enforcer := NewEnforcer(rbac)

	if err := enforcer.Enforce(adminUser, "user", "read"); err != nil {
		t.Errorf("Enforce should succeed for admin, got error: %v", err)
	}

	if err := enforcer.Enforce(regularUser, "user", "write"); err == nil {
		t.Error("Enforce should fail for regular user")
	}
}

func TestRBAC_ContextIntegration(t *testing.T) {
	rbac := NewRBAC()

	role := &Role{
		ID:   "test",
		Name: "Test Role",
		Permissions: []*Permission{
			{ID: "test:read", Name: "Test Read", Resource: "test", Action: "read"},
		},
	}
	rbac.AddRole(role)

	user := &User{
		ID:    "user-1",
		Roles: []string{"test"},
	}

	// 测试 Context 集成
	ctx := context.Background()
	ctx = WithUser(ctx, user)

	retrievedUser, ok := GetUserFromContext(ctx)
	if !ok {
		t.Error("User should be in context")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	// 测试权限检查
	if !rbac.CheckPermission(retrievedUser, "test", "read") {
		t.Error("User from context should have test:read permission")
	}
}

func TestRBAC_EdgeCases(t *testing.T) {
	rbac := NewRBAC()

	// 测试空用户
	emptyUser := &User{ID: "empty", Roles: []string{}}
	if rbac.CheckPermission(emptyUser, "any", "any") {
		t.Error("Empty user should not have any permissions")
	}

	// 测试不存在的角色
	userWithInvalidRole := &User{
		ID:    "user-1",
		Roles: []string{"nonexistent"},
	}
	if rbac.CheckPermission(userWithInvalidRole, "any", "any") {
		t.Error("User with invalid role should not have permissions")
	}

	// 测试通配符权限
	wildcardRole := &Role{
		ID:   "wildcard",
		Name: "Wildcard Role",
		Permissions: []*Permission{
			{ID: "all", Name: "All", Resource: "*", Action: "*"},
		},
	}
	rbac.AddRole(wildcardRole)

	wildcardUser := &User{
		ID:    "user-2",
		Roles: []string{"wildcard"},
	}

	if !rbac.CheckPermission(wildcardUser, "any", "any") {
		t.Error("Wildcard role should grant all permissions")
	}
}

func TestRBAC_ConcurrentAccess(t *testing.T) {
	rbac := NewRBAC()

	role := &Role{
		ID:   "test",
		Name: "Test Role",
		Permissions: []*Permission{
			{ID: "test:read", Name: "Test Read", Resource: "test", Action: "read"},
		},
	}
	rbac.AddRole(role)

	user := &User{
		ID:    "user-1",
		Roles: []string{"test"},
	}

	// 并发测试
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				rbac.CheckPermission(user, "test", "read")
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}
