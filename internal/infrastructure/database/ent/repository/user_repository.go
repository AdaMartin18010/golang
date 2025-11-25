package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/infrastructure/database/ent"
)

// UserRepository Ent 实现的用户仓库
// 注意：需要先运行 make generate-ent 生成 Ent 代码
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository 创建用户仓库
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Create 创建用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return fmt.Errorf("not implemented: run 'make generate-ent' first")
}

// FindByID 根据 ID 获取用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return nil, fmt.Errorf("not implemented: run 'make generate-ent' first")
}

// FindByEmail 根据邮箱获取用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, fmt.Errorf("not implemented: run 'make generate-ent' first")
}

// Update 更新用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return fmt.Errorf("not implemented: run 'make generate-ent' first")
}

// Delete 删除用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented: run 'make generate-ent' first")
}

// List 列出用户
// TODO: 生成 Ent 代码后实现
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	return nil, fmt.Errorf("not implemented: run 'make generate-ent' first")
}
