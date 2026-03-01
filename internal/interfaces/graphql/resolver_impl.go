// Package graphql provides GraphQL resolver implementations.
//
// 本文件实现了 GraphQL 查询和变更的具体逻辑。
package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/golang/internal/domain/user"
)

// User 查询单个用户。
func (r *Query) User(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// 调用应用层服务获取用户
	u, err := r.resolver.userService.GetUser(ctx, id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, nil // GraphQL 中返回 nil 表示未找到
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return domainUserToGraphQL(u), nil
}

// Users 查询用户列表。
func (r *Query) Users(ctx context.Context, limit *int, offset *int) ([]*User, error) {
	// 设置默认值
	l := 10
	if limit != nil && *limit > 0 {
		l = *limit
	}
	o := 0
	if offset != nil && *offset >= 0 {
		o = *offset
	}

	// 调用应用层服务获取用户列表
	users, err := r.resolver.userService.ListUsers(ctx, l, o)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// 转换为 GraphQL 类型
	result := make([]*User, len(users))
	for i, u := range users {
		result[i] = domainUserToGraphQL(u)
	}
	return result, nil
}

// CreateUser 创建用户。
func (r *Mutation) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	// 验证输入
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// 调用应用层服务创建用户
	u, err := r.resolver.userService.CreateUser(ctx, input.Email, input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return domainUserToGraphQL(u), nil
}

// UpdateUser 更新用户。
func (r *Mutation) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// 获取现有用户
	u, err := r.resolver.userService.GetUser(ctx, id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 更新名称
	if input.Name != nil && *input.Name != "" {
		if err := r.resolver.userService.UpdateUserName(ctx, id, *input.Name); err != nil {
			return nil, fmt.Errorf("failed to update user name: %w", err)
		}
		u.UpdateName(*input.Name)
	}

	return domainUserToGraphQL(u), nil
}

// DeleteUser 删除用户。
func (r *Mutation) DeleteUser(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("user ID is required")
	}

	if err := r.resolver.userService.DeleteUser(ctx, id); err != nil {
		if err == user.ErrUserNotFound {
			return false, nil // 用户不存在视为删除成功
		}
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	return true, nil
}

// UpdateUserInput 是更新用户的输入类型。
type UpdateUserInput struct {
	Name *string
}

// domainUserToGraphQL 将领域用户转换为 GraphQL 用户。
func domainUserToGraphQL(u *user.User) *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
