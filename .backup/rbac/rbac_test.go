package rbac

import (
	"context"
	"testing"
)

func TestRBAC_AddPermission(t *testing.T) {
	rbac := NewRBAC()
	permission := &Permission{
		ID:       "perm1",
		Name:     "read_users",
		Resource: "users",
		Action:   "read",
	}

	rbac.AddPermission(permission)

	retrieved, err := rbac.GetPermission("perm1")
	if err != nil {
		t.Fatalf("Failed to get permission: %v", err)
	}

	if retrieved.Name != "read_users" {
		t.Errorf("Expected permission name 'read_users', got '%s'", retrieved.Name)
	}
}

func TestRBAC_AddRole(t *testing.T) {
	rbac := NewRBAC()
	role := &Role{
		ID:   "role1",
		Name: "admin",
	}

	rbac.AddRole(role)

	retrieved, err := rbac.GetRole("role1")
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	if retrieved.Name != "admin" {
		t.Errorf("Expected role name 'admin', got '%s'", retrieved.Name)
	}
}

func TestRBAC_AssignPermission(t *testing.T) {
	rbac := NewRBAC()

	permission := &Permission{
		ID:       "perm1",
		Name:     "read_users",
		Resource: "users",
		Action:   "read",
	}
	rbac.AddPermission(permission)

	role := &Role{
		ID:   "role1",
		Name: "admin",
	}
	rbac.AddRole(role)

	err := rbac.AssignPermission("role1", "perm1")
	if err != nil {
		t.Fatalf("Failed to assign permission: %v", err)
	}

	role, _ = rbac.GetRole("role1")
	if len(role.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(role.Permissions))
	}
}

func TestRBAC_CheckPermission(t *testing.T) {
	rbac := NewRBAC()

	permission := &Permission{
		ID:       "perm1",
		Name:     "read_users",
		Resource: "users",
		Action:   "read",
	}
	rbac.AddPermission(permission)

	role := &Role{
		ID:   "role1",
		Name: "admin",
	}
	rbac.AddRole(role)
	rbac.AssignPermission("role1", "perm1")

	user := &User{
		ID:    "user1",
		Roles: []string{"role1"},
	}

	// 检查有权限
	if !rbac.CheckPermission(user, "users", "read") {
		t.Error("Expected user to have permission")
	}

	// 检查无权限
	if rbac.CheckPermission(user, "users", "write") {
		t.Error("Expected user to not have permission")
	}
}

func TestRBAC_GetUserPermissions(t *testing.T) {
	rbac := NewRBAC()

	perm1 := &Permission{ID: "perm1", Resource: "users", Action: "read"}
	perm2 := &Permission{ID: "perm2", Resource: "users", Action: "write"}
	rbac.AddPermission(perm1)
	rbac.AddPermission(perm2)

	role1 := &Role{ID: "role1", Name: "admin"}
	role2 := &Role{ID: "role2", Name: "editor"}
	rbac.AddRole(role1)
	rbac.AddRole(role2)

	rbac.AssignPermission("role1", "perm1")
	rbac.AssignPermission("role2", "perm2")

	user := &User{
		ID:    "user1",
		Roles: []string{"role1", "role2"},
	}

	permissions := rbac.GetUserPermissions(user)
	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	}
}

func TestEnforcer_Enforce(t *testing.T) {
	rbac := NewRBAC()

	permission := &Permission{
		ID:       "perm1",
		Resource: "users",
		Action:   "read",
	}
	rbac.AddPermission(permission)

	role := &Role{ID: "role1", Name: "admin"}
	rbac.AddRole(role)
	rbac.AssignPermission("role1", "perm1")

	user := &User{
		ID:    "user1",
		Roles: []string{"role1"},
	}

	enforcer := NewEnforcer(rbac)

	// 有权限
	err := enforcer.Enforce(user, "users", "read")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 无权限
	err = enforcer.Enforce(user, "users", "write")
	if err != ErrPermissionDenied {
		t.Errorf("Expected ErrPermissionDenied, got %v", err)
	}
}

func TestGetUserFromContext(t *testing.T) {
	user := &User{
		ID:    "user1",
		Roles: []string{"role1"},
	}

	ctx := WithUser(context.Background(), user)
	retrieved, ok := GetUserFromContext(ctx)
	if !ok {
		t.Error("Expected user in context")
	}

	if retrieved.ID != "user1" {
		t.Errorf("Expected user ID 'user1', got '%s'", retrieved.ID)
	}
}
