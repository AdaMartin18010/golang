package rbac

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// RBAC 基于角色的访问控制
// 实现标准的 RBAC 模型
type RBAC struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	mu          sync.RWMutex
}

// Role 角色
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []string // Permission IDs
	Inherits    []string // Parent Role IDs (角色继承)
}

// Permission 权限
type Permission struct {
	ID          string
	Resource    string // 资源类型 (e.g., "user", "post")
	Action      string // 操作 (e.g., "create", "read", "update", "delete")
	Description string
}

// NewRBAC 创建 RBAC 实例
func NewRBAC() *RBAC {
	return &RBAC{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
	}
}

// AddRole 添加角色
func (r *RBAC) AddRole(role *Role) error {
	if role == nil || role.ID == "" {
		return errors.New("invalid role")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[role.ID]; exists {
		return fmt.Errorf("role %s already exists", role.ID)
	}

	r.roles[role.ID] = role
	return nil
}

// AddPermission 添加权限
func (r *RBAC) AddPermission(perm *Permission) error {
	if perm == nil || perm.ID == "" {
		return errors.New("invalid permission")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.permissions[perm.ID]; exists {
		return fmt.Errorf("permission %s already exists", perm.ID)
	}

	r.permissions[perm.ID] = perm
	return nil
}

// AssignPermissionToRole 为角色分配权限
func (r *RBAC) AssignPermissionToRole(roleID, permID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role, exists := r.roles[roleID]
	if !exists {
		return fmt.Errorf("role %s not found", roleID)
	}

	if _, exists := r.permissions[permID]; !exists {
		return fmt.Errorf("permission %s not found", permID)
	}

	// 检查是否已分配
	for _, p := range role.Permissions {
		if p == permID {
			return nil // 已存在
		}
	}

	role.Permissions = append(role.Permissions, permID)
	return nil
}

// CheckPermission 检查用户是否有权限
func (r *RBAC) CheckPermission(ctx context.Context, userRoles []string, resource, action string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 收集用户所有权限（包括继承的）
	userPermissions := make(map[string]bool)
	for _, roleID := range userRoles {
		perms, err := r.getRolePermissions(roleID, make(map[string]bool))
		if err != nil {
			return false, err
		}
		for permID := range perms {
			userPermissions[permID] = true
		}
	}

	// 检查是否有匹配的权限
	for permID := range userPermissions {
		perm, exists := r.permissions[permID]
		if !exists {
			continue
		}

		// 检查资源和操作是否匹配
		if (perm.Resource == resource || perm.Resource == "*") &&
			(perm.Action == action || perm.Action == "*") {
			return true, nil
		}
	}

	return false, nil
}

// getRolePermissions 获取角色的所有权限（递归获取继承的权限）
func (r *RBAC) getRolePermissions(roleID string, visited map[string]bool) (map[string]bool, error) {
	// 防止循环继承
	if visited[roleID] {
		return nil, fmt.Errorf("circular role inheritance detected: %s", roleID)
	}
	visited[roleID] = true

	role, exists := r.roles[roleID]
	if !exists {
		return nil, fmt.Errorf("role %s not found", roleID)
	}

	permissions := make(map[string]bool)

	// 添加当前角色的权限
	for _, permID := range role.Permissions {
		permissions[permID] = true
	}

	// 递归添加继承角色的权限
	for _, parentRoleID := range role.Inherits {
		parentPerms, err := r.getRolePermissions(parentRoleID, visited)
		if err != nil {
			return nil, err
		}
		for permID := range parentPerms {
			permissions[permID] = true
		}
	}

	return permissions, nil
}

// GetRole 获取角色
func (r *RBAC) GetRole(roleID string) (*Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, exists := r.roles[roleID]
	if !exists {
		return nil, fmt.Errorf("role %s not found", roleID)
	}

	return role, nil
}

// GetPermission 获取权限
func (r *RBAC) GetPermission(permID string) (*Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perm, exists := r.permissions[permID]
	if !exists {
		return nil, fmt.Errorf("permission %s not found", permID)
	}

	return perm, nil
}

// ListRoles 列出所有角色
func (r *RBAC) ListRoles() []*Role {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roles := make([]*Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles
}

// ListPermissions 列出所有权限
func (r *RBAC) ListPermissions() []*Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perms := make([]*Permission, 0, len(r.permissions))
	for _, perm := range r.permissions {
		perms = append(perms, perm)
	}
	return perms
}

// 预定义的常用角色和权限

// InitializeDefaultRoles 初始化默认角色
func (r *RBAC) InitializeDefaultRoles() error {
	// 定义权限
	permissions := []*Permission{
		{ID: "user.create", Resource: "user", Action: "create", Description: "Create users"},
		{ID: "user.read", Resource: "user", Action: "read", Description: "Read users"},
		{ID: "user.update", Resource: "user", Action: "update", Description: "Update users"},
		{ID: "user.delete", Resource: "user", Action: "delete", Description: "Delete users"},
		{ID: "admin.all", Resource: "*", Action: "*", Description: "All permissions"},
	}

	for _, perm := range permissions {
		if err := r.AddPermission(perm); err != nil {
			return err
		}
	}

	// 定义角色
	roles := []*Role{
		{
			ID:          "admin",
			Name:        "Administrator",
			Description: "Full system access",
			Permissions: []string{"admin.all"},
		},
		{
			ID:          "user",
			Name:        "Regular User",
			Description: "Basic user access",
			Permissions: []string{"user.read"},
		},
		{
			ID:          "moderator",
			Name:        "Moderator",
			Description: "User management access",
			Permissions: []string{"user.read", "user.update"},
			Inherits:    []string{"user"}, // 继承 user 角色
		},
	}

	for _, role := range roles {
		if err := r.AddRole(role); err != nil {
			return err
		}
	}

	return nil
}
