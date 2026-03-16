// Package repository provides Ent-based user repository implementation.
//
// 此文件实现了基于 Ent ORM 的用户仓储接口，将领域实体与 Ent 实体进行转换。
package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/infra/database/ent"
	"github.com/yourusername/golang/internal/infra/database/ent/repository"
	entuser "github.com/yourusername/golang/internal/infra/database/ent/user"
)

// EntUserRepository 是基于 Ent 的用户仓储实现
type EntUserRepository struct {
	*repository.BaseRepository[user.User, *ent.User]
	client *ent.Client
}

// NewEntUserRepository 创建基于 Ent 的用户仓储
func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	toDomain := func(e *ent.User) (*user.User, error) {
		return &user.User{
			ID:        e.ID,
			Email:     e.Email,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}, nil
	}

	toEnt := func(u *user.User) (*ent.User, error) {
		return &ent.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}, nil
	}

	getID := func(u *user.User) (string, error) {
		return u.ID, nil
	}

	setID := func(u *user.User, id string) error {
		u.ID = id
		return nil
	}

	baseRepo := repository.NewBaseRepository[user.User, *ent.User](client, toDomain, toEnt, getID, setID)

	return &EntUserRepository{
		BaseRepository: baseRepo,
		client:         client,
	}
}

// Save 保存用户（创建或更新）
func (r *EntUserRepository) Save(ctx context.Context, u *user.User) error {
	return r.Create(ctx, u)
}

// Create 创建用户
func (r *EntUserRepository) Create(ctx context.Context, u *user.User) error {
	entUser, err := r.client.User.Create().
		SetID(u.ID).
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 更新领域实体的创建时间和更新时间
	u.CreatedAt = entUser.CreatedAt
	u.UpdatedAt = entUser.UpdatedAt

	return nil
}

// FindByID 根据 ID 查找用户
func (r *EntUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	entUser, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return r.toDomain(entUser), nil
}

// FindByEmail 根据邮箱查找用户
func (r *EntUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	entUser, err := r.client.User.Query().
		Where(entuser.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.toDomain(entUser), nil
}

// Update 更新用户
func (r *EntUserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := r.client.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetEmail(u.Email).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *EntUserRepository) Delete(ctx context.Context, id string) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List 列出用户（支持分页）
func (r *EntUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	entUsers, err := r.client.User.Query().
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*user.User, len(entUsers))
	for i, entUser := range entUsers {
		users[i] = r.toDomain(entUser)
	}

	return users, nil
}

// toDomain 将 Ent 用户转换为领域用户
func (r *EntUserRepository) toDomain(entUser *ent.User) *user.User {
	return &user.User{
		ID:        entUser.ID,
		Email:     entUser.Email,
		Name:      entUser.Name,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}
}

