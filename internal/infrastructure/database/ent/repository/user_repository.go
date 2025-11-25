package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/infrastructure/database/ent"
	entuser "github.com/yourusername/golang/internal/infrastructure/database/ent/user"
)

// UserRepository Ent 实现的用户仓库
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository 创建用户仓库
func NewUserRepository(client *ent.Client) user.Repository {
	return &UserRepository{client: client}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	id := uuid.New().String()
	entUser, err := r.client.User.
		Create().
		SetID(id).
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	u.ID = entUser.ID
	u.CreatedAt = entUser.CreatedAt
	u.UpdatedAt = entUser.UpdatedAt
	return nil
}

// FindByID 根据 ID 获取用户
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(entuser.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return toDomainUser(entUser), nil
}

// FindByEmail 根据邮箱获取用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(entuser.Email(email)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return toDomainUser(entUser), nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	update := r.client.User.
		UpdateOneID(u.ID)

	if u.Email != "" {
		update = update.SetEmail(u.Email)
	}
	if u.Name != "" {
		update = update.SetName(u.Name)
	}

	entUser, err := update.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	u.UpdatedAt = entUser.UpdatedAt
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	err := r.client.User.
		DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 列出用户
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	entUsers, err := r.client.User.
		Query().
		Order(entuser.ByCreatedAt(sql.OrderDesc())).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*user.User, len(entUsers))
	for i, entUser := range entUsers {
		users[i] = toDomainUser(entUser)
	}
	return users, nil
}

// toDomainUser 将 Ent 用户转换为领域用户
func toDomainUser(entUser *ent.User) *user.User {
	return &user.User{
		ID:        entUser.ID,
		Email:     entUser.Email,
		Name:      entUser.Name,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}
}
