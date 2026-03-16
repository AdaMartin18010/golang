package rbac

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrPermissionDenied 权限被拒绝
	ErrPermissionDenied = errors.New("permission denied")
	// ErrRoleNotFound 角色未找到
	ErrRoleNotFound = errors.New("role not found")
	// ErrPermissionNotFound 权限未找到
	ErrPermissionNotFound = errors.New("permission not found")
)

// Permission 权限
type Permission struct {
	ID          string
	Name        string
	Description string
	Resource    string
	Action      string
}

// Role 角色
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []*Permission
}

// User 用户
type User struct {
	ID    string
	Roles []string
}

// RBAC 基于角色的访问控制
type RBAC struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	mu          sync.RWMutex
}

// NewRBAC 创建RBAC实例
func NewRBAC() *RBAC {
	return &RBAC{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
	}
}

// AddPermission 添加权限
func (r *RBAC) AddPermission(permission *Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.permissions[permission.ID] = permission
}

// GetPermission 获取权限
func (r *RBAC) GetPermission(permissionID string) (*Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permission, exists := r.permissions[permissionID]
	if !exists {
		return nil, ErrPermissionNotFound
	}

	return permission, nil
}

// AddRole 添加角色
func (r *RBAC) AddRole(role *Role) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.roles[role.ID] = role
}

// GetRole 获取角色
func (r *RBAC) GetRole(roleID string) (*Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, exists := r.roles[roleID]
	if !exists {
		return nil, ErrRoleNotFound
	}

	return role, nil
}

// AssignPermission 为角色分配权限
func (r *RBAC) AssignPermission(roleID, permissionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role, exists := r.roles[roleID]
	if !exists {
		return ErrRoleNotFound
	}

	permission, exists := r.permissions[permissionID]
	if !exists {
		return ErrPermissionNotFound
	}

	// 检查权限是否已存在
	for _, p := range role.Permissions {
		if p.ID == permissionID {
			return nil // 已存在，直接返回
		}
	}

	role.Permissions = append(role.Permissions, permission)
	return nil
}

// RemovePermission 从角色移除权限
func (r *RBAC) RemovePermission(roleID, permissionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role, exists := r.roles[roleID]
	if !exists {
		return ErrRoleNotFound
	}

	for i, p := range role.Permissions {
		if p.ID == permissionID {
			role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
			return nil
		}
	}

	return ErrPermissionNotFound
}

// CheckPermission 检查用户是否有权限
func (r *RBAC) CheckPermission(user *User, resource, action string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 遍历用户的所有角色
	for _, roleID := range user.Roles {
		role, exists := r.roles[roleID]
		if !exists {
			continue
		}

		// 检查角色是否有权限
		for _, permission := range role.Permissions {
			if permission.Resource == resource && permission.Action == action {
				return true
			}
		}
	}

	return false
}

// CheckPermissionByID 通过权限ID检查权限
func (r *RBAC) CheckPermissionByID(user *User, permissionID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 遍历用户的所有角色
	for _, roleID := range user.Roles {
		role, exists := r.roles[roleID]
		if !exists {
			continue
		}

		// 检查角色是否有权限
		for _, permission := range role.Permissions {
			if permission.ID == permissionID {
				return true
			}
		}
	}

	return false
}

// GetUserPermissions 获取用户的所有权限
func (r *RBAC) GetUserPermissions(user *User) []*Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permissionMap := make(map[string]*Permission)

	// 遍历用户的所有角色
	for _, roleID := range user.Roles {
		role, exists := r.roles[roleID]
		if !exists {
			continue
		}

		// 收集权限
		for _, permission := range role.Permissions {
			permissionMap[permission.ID] = permission
		}
	}

	// 转换为切片
	permissions := make([]*Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions
}

// Enforcer 权限执行器
type Enforcer struct {
	rbac *RBAC
}

// NewEnforcer 创建权限执行器
func NewEnforcer(rbac *RBAC) *Enforcer {
	return &Enforcer{
		rbac: rbac,
	}
}

// Enforce 执行权限检查
func (e *Enforcer) Enforce(user *User, resource, action string) error {
	if !e.rbac.CheckPermission(user, resource, action) {
		return ErrPermissionDenied
	}
	return nil
}

// EnforceByID 通过权限ID执行权限检查
func (e *Enforcer) EnforceByID(user *User, permissionID string) error {
	if !e.rbac.CheckPermissionByID(user, permissionID) {
		return ErrPermissionDenied
	}
	return nil
}

// GetUserFromContext 从context获取用户（辅助函数）
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey{}).(*User)
	return user, ok
}

// WithUser 将用户添加到context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// userKey context key类型
type userKey struct{}
