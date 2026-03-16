package user

import (
	"context"
	"fmt"

	"github.com/yourusername/golang/internal/domain/user"
)

// 错误类型常量，用于统一错误标识
const (
	ErrTypeUserAlreadyExists = "user_already_exists"
	ErrTypeUserNotFound      = "user_not_found"
	ErrTypeInvalidInput      = "invalid_input"
	ErrTypeInternal          = "internal_error"
)

// Service 用户应用服务
type Service struct {
	repo UserRepository
}

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Save(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*user.User, error)
}

// NewService 创建用户服务
func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetUser 获取用户
func (s *Service) GetUser(ctx context.Context, id string) (*user.User, error) {
	return s.repo.FindByID(ctx, id)
}

// CreateUser 创建用户
func (s *Service) CreateUser(ctx context.Context, email, name string) (*user.User, error) {
	// 检查邮箱是否已存在
	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: user with email %s already exists", ErrUserAlreadyExists, email)
	}

	// 创建新用户
	newUser := user.NewUser(email, name)

	// 验证
	if err := newUser.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid user: %w", err)
	}

	// 保存
	if err := s.repo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return newUser, nil
}

// UpdateUserName 更新用户名称
func (s *Service) UpdateUserName(ctx context.Context, id, name string) error {
	// 获取用户
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 更新名称
	u.UpdateName(name)

	// 保存
	if err := s.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ListUsers 列出用户
func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("%w: limit must be positive", ErrInvalidInput)
	}
	if offset < 0 {
		return nil, fmt.Errorf("%w: offset cannot be negative", ErrInvalidInput)
	}

	return s.repo.List(ctx, limit, offset)
}
