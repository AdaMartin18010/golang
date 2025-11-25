package user

import "context"

// Repository 用户仓储接口（领域层定义）
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户
	Delete(ctx context.Context, id string) error

	// List 列出用户（支持分页）
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

