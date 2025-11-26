# RBAC 权限控制框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [RBAC 权限控制框架](#rbac-权限控制框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心概念](#2-核心概念)
    - [2.1 Permission（权限）](#21-permission权限)
    - [2.2 Role（角色）](#22-role角色)
    - [2.3 User（用户）](#23-user用户)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 使用Enforcer](#32-使用enforcer)
    - [3.3 在Handler中使用](#33-在handler中使用)
    - [3.4 获取用户所有权限](#34-获取用户所有权限)
    - [3.5 与认证中间件集成](#35-与认证中间件集成)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

RBAC（基于角色的访问控制）框架提供了完整的权限管理功能：

- ✅ **角色管理**: 创建和管理角色
- ✅ **权限管理**: 创建和管理权限
- ✅ **权限分配**: 为角色分配权限
- ✅ **权限检查**: 检查用户权限
- ✅ **Context集成**: 与Context集成

---

## 2. 核心概念

### 2.1 Permission（权限）

权限定义了可以执行的操作：

```go
type Permission struct {
    ID          string
    Name        string
    Description string
    Resource    string  // 资源（如：users, orders）
    Action      string  // 操作（如：read, write, delete）
}
```

### 2.2 Role（角色）

角色是权限的集合：

```go
type Role struct {
    ID          string
    Name        string
    Description string
    Permissions []*Permission
}
```

### 2.3 User（用户）

用户拥有一个或多个角色：

```go
type User struct {
    ID    string
    Roles []string  // 角色ID列表
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "github.com/yourusername/golang/pkg/rbac"
)

// 创建RBAC实例
rbac := rbac.NewRBAC()

// 创建权限
readUsersPerm := &rbac.Permission{
    ID:       "perm_read_users",
    Name:     "read_users",
    Resource: "users",
    Action:   "read",
}
rbac.AddPermission(readUsersPerm)

writeUsersPerm := &rbac.Permission{
    ID:       "perm_write_users",
    Name:     "write_users",
    Resource: "users",
    Action:   "write",
}
rbac.AddPermission(writeUsersPerm)

// 创建角色
adminRole := &rbac.Role{
    ID:   "role_admin",
    Name: "admin",
}
rbac.AddRole(adminRole)

// 为角色分配权限
rbac.AssignPermission("role_admin", "perm_read_users")
rbac.AssignPermission("role_admin", "perm_write_users")

// 创建用户
user := &rbac.User{
    ID:    "user1",
    Roles: []string{"role_admin"},
}

// 检查权限
if rbac.CheckPermission(user, "users", "read") {
    // 用户有读取权限
}
```

### 3.2 使用Enforcer

```go
enforcer := rbac.NewEnforcer(rbac)

// 执行权限检查
err := enforcer.Enforce(user, "users", "write")
if err != nil {
    // 权限被拒绝
    return err
}

// 继续执行操作
```

### 3.3 在Handler中使用

```go
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    // 从context获取用户
    user, ok := rbac.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized,
            errors.NewUnauthorizedError("user not found"))
        return
    }

    // 检查权限
    err := h.enforcer.Enforce(user, "users", "delete")
    if err != nil {
        response.Error(w, http.StatusForbidden, err)
        return
    }

    // 执行删除操作
    // ...
}
```

### 3.4 获取用户所有权限

```go
permissions := rbac.GetUserPermissions(user)
for _, perm := range permissions {
    fmt.Printf("Permission: %s on %s\n", perm.Action, perm.Resource)
}
```

### 3.5 与认证中间件集成

```go
func (m *AuthMiddleware) WithRBAC(rbac *rbac.RBAC) func(http.Handler) http.Handler {
    enforcer := rbac.NewEnforcer(rbac)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 获取用户
            claims, _ := GetClaims(r.Context())
            user := &rbac.User{
                ID:    claims.UserID,
                Roles: claims.Roles,
            }

            // 将用户添加到context
            ctx := rbac.WithUser(r.Context(), user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **权限命名**: 使用清晰的权限命名（resource:action）
2. **角色层次**: 设计合理的角色层次结构
3. **最小权限**: 遵循最小权限原则
4. **权限缓存**: 对频繁检查的权限进行缓存
5. **权限审计**: 记录权限检查日志

### 4.2 DON'Ts ❌

1. **不要硬编码**: 不要硬编码权限检查
2. **不要过度分配**: 不要给角色分配过多权限
3. **不要忽略错误**: 始终检查权限检查结果
4. **不要忘记更新**: 权限变更时及时更新

---

## 5. 相关资源

- [认证授权中间件](../../internal/interfaces/http/chi/middleware/auth/README.md)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
